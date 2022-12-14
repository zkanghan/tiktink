package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/snowid"
	"tiktink/pkg/tracer"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

const (
	commentPageRows = 20
)

type commentFunc interface {
	CreateComment(videoID string, userID string, content string) (*model.Comment, error)
	DeleteComment(videoID string, CommentID string, userID string) (int64, error)
	QueryCommentList(videoID string, pn int) ([]*model.CommentMSG, error)
}

type commentDealer struct{}

func NewCommentDealer() commentFunc {
	return &commentDealer{}
}

func (c *commentDealer) QueryCommentList(videoID string, pn int) ([]*model.CommentMSG, error) {
	var commentMsg []*model.CommentMSG
	offset := (pn - 1) * commentPageRows
	count := commentPageRows
	err := db.Raw("select `user_id`,`user_name`,`follow_count`,`follower_count`,`comment_id`,"+
		"`content`,`comments`.`create_date`"+
		"from `users` inner join   `comments` "+
		"on  `comments`.author_id=`user_id` where `comments`.video_id = ? limit ?,?", videoID, offset, count).Scan(&commentMsg).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(videoID, pn))
	}
	return commentMsg, nil
}

// DeleteComment 只有评论的所有者可以删除评论，userID加了一层校验  返回受影响的行数
func (c *commentDealer) DeleteComment(videoID string, CommentID string, userID string) (affectRows int64, err error) {
	var todoDB *gorm.DB
	err = db.Transaction(func(tx *gorm.DB) error {
		todoDB = tx.Where("comment_id = ? and author_id = ? and video_id = ?", CommentID, userID, videoID).Delete(&model.Comment{})
		if todoDB.Error != nil {
			return errors.Wrap(todoDB.Error, tracer.FormatParam(videoID, CommentID, userID))
		}
		//  如果删除成功，对应视频下的评论数减1
		if todoDB.RowsAffected == 1 {
			if err := tx.Exec("update videos set comment_count = comment_count-1 where video_id = ?", videoID).Error; err != nil {
				return errors.Wrap(err, tracer.FormatParam(videoID, CommentID, userID))
			}
		}
		return nil
	})
	if err != nil {
		return 0, errors.Wrap(todoDB.Error, tracer.FormatParam(videoID, CommentID, userID))
	}
	return todoDB.RowsAffected, nil
}

func (c *commentDealer) CreateComment(videoID string, userID string, content string) (*model.Comment, error) {
	comment := &model.Comment{
		CommentID: snowid.GenID(),
		VideoID:   videoID,
		AuthorID:  userID,
		Content:   content,
	}
	//  这里对应视频的评论数也要加1，先™解决创建记录的问题吧
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建记录
		if err := tx.Select("comment_id", "video_id", "author_id", "content").Create(comment).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(videoID, userID, content))
		}
		//视频下对应的评论数加1
		if err := tx.Exec("update videos set comment_count = comment_count+1 where video_id = ?", videoID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(videoID, userID, content))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(videoID, userID, content))
	}
	return comment, nil
}
