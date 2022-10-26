package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

const (
	favoriteVideoPageRows = 20
)

type favoriteFunc interface {
	QueryIsLiked(userID string, videoId string) (bool, error)
	DoFavorite(userID string, videoID string) error
	CancelFavorite(userID string, videoID string) error
	QueryFavoriteList(userID string, pn int) ([]*model.VideoMSG, error)
	QueryListIsLiked(userID string, videoIDs []string) ([]string, error)
}

type favoriteDealer struct{}

func NewFavoriteDealer() favoriteFunc {
	return &favoriteDealer{}
}

func (f *favoriteDealer) QueryFavoriteList(userID string, pn int) ([]*model.VideoMSG, error) {
	var videoMsgs []*model.VideoMSG
	offset := (pn - 1) * favoriteVideoPageRows
	count := favoriteVideoPageRows
	err := db.Raw("SELECT `video_id`,`user_id`,`user_name`,`follow_count`,`follower_count`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title` "+
		"FROM `users` INNER JOIN `videos` "+
		"ON `videos`.`author_id` = `user_id`"+
		"where `video_id` in (select `favorites`.`video_id` from favorites where user_id = ?) limit ?,?", userID, offset, count).Scan(&videoMsgs).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(userID, pn))
	}
	return videoMsgs, nil
}

func (f *favoriteDealer) CancelFavorite(authorID string, videoID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and video_id = ?", authorID, videoID).Delete(&model.Favorite{}).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(authorID, videoID))
		}
		if err := tx.Exec("update videos set favorite_count = favorite_count-1 where video_id = ?", videoID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(authorID, videoID))
		}
		return nil
	})
}

func (f *favoriteDealer) DoFavorite(authorID string, videoID string) error {
	favoriteModel := &model.Favorite{
		VideoID: videoID,
		UserID:  authorID,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(favoriteModel).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(authorID, videoID))
		}
		if err := tx.Exec("update videos set favorite_count = favorite_count+1 where video_id = ?", videoID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(authorID, videoID))
		}
		return nil
	})
}

func (f *favoriteDealer) QueryIsLiked(authorID string, videoId string) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from favorites where user_id = ? and video_id = ? limit 1", authorID, videoId).Scan(res).Error
	return *res == 1, errors.Wrap(err, tracer.FormatParam(authorID, videoId))
}

// QueryListIsLiked 返回videoIDs中 user点赞的部分
func (f *favoriteDealer) QueryListIsLiked(userID string, videoIDs []string) ([]string, error) {
	var likedVideoID []string
	err := db.Raw("select video_id from favorites where user_id = ? and video_id in ?", userID, videoIDs).Scan(&likedVideoID).Error
	if err != nil {
		return []string{}, errors.Wrap(err, tracer.FormatParam(userID, videoIDs))
	}
	return likedVideoID, nil
}
