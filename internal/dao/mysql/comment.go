package mysql

import (
	"tiktink/internal/model"

	"gorm.io/gorm"
)

type commentFunc interface {
	CreateComment(videoID int64, userID int64, content string) (*model.Comment, error)
	DeleteComment(videoID int64, CommentID int64, userID int64) (int64, error)
	QueryCommentList(videoID int64) ([]*model.CommentMSG, error)
}

type commentDealer struct {
}

func (c commentDealer) QueryCommentList(videoID int64) ([]*model.CommentMSG, error) {
	var commentMsg []*model.CommentMSG
	err := db.Raw("select `user_id`,`user_name`,`follow_count`,`follower_count`,`comment_id`,"+
		"`content`,`comments`.`create_date`"+
		"from `users` inner join   `comments` "+
		"on  `comments`.author_id=`user_id` where `comments`.video_id = ?", videoID).Scan(&commentMsg).Error
	if err != nil {
		return nil, err
	}
	return commentMsg, nil
}

// DeleteComment 只有评论的所有者可以删除评论，userID加了一层校验  返回受影响的行数
func (c commentDealer) DeleteComment(videoID int64, CommentID int64, userID int64) (int64, error) {
	var todoDB *gorm.DB
	err := db.Transaction(func(tx *gorm.DB) error {
		todoDB = tx.Where("comment_id = ? and author_id = ? and video_id = ?", CommentID, userID, videoID).Delete(&model.Comment{})
		if todoDB.Error != nil {
			return todoDB.Error
		}
		//  如果删除成功，对应视频下的评论数减1
		if todoDB.RowsAffected == 1 {
			if err := tx.Exec("update videos set comment_count = comment_count-1 where video_id = ?", videoID).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return todoDB.RowsAffected, nil
}

func (c commentDealer) CreateComment(videoID int64, userID int64, content string) (*model.Comment, error) {
	comment := &model.Comment{
		VideoID:  videoID,
		AuthorID: userID,
		Content:  content,
	}
	//  这里对应视频的评论数也要加1，先™解决创建记录的问题吧
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建记录
		if err := tx.Select("video_id", "author_id", "content").Create(comment).Error; err != nil {
			return err
		}
		//视频下对应的评论数加1
		if err := tx.Exec("update videos set comment_count = comment_count+1 "+
			"where video_id = ?", videoID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func DealComment() commentFunc {
	return &commentDealer{}
}
