package mysql

import (
	"tiktink/internal/model"

	"gorm.io/gorm"
)

type favoriteFunc interface {
	QueryIsLiked(userID int64, videoId int64) (bool, error)
	DoFavorite(userID int64, videoID int64) error
	CancelFavorite(userID int64, videoID int64) error
	QueryFavoriteList(userID int64) ([]*model.VideoMSG, error)
}

type favoriteDealer struct {
}

func (f favoriteDealer) QueryFavoriteList(userID int64) ([]*model.VideoMSG, error) {
	var videoMsgs []*model.VideoMSG
	err := db.Raw("SELECT `video_id`,`user_id`,`user_name`,`follow_count`,`follower_count`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title` "+
		"FROM `users` INNER JOIN `videos` "+
		"ON `videos`.`author_id` = `user_id`"+
		"where `video_id` in (select `favorites`.`video_id` from favorites where user_id = ?)", userID).Scan(&videoMsgs).Error
	if err != nil {
		return nil, err
	}
	return videoMsgs, nil
}

func (f favoriteDealer) CancelFavorite(authorID int64, videoID int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and video_id = ?", authorID, videoID).Delete(&model.Favorite{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("update videos set favorite_count = favorite_count-1 where video_id = ?", videoID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f favoriteDealer) DoFavorite(authorID int64, videoID int64) error {
	favoriteModel := &model.Favorite{
		VideoID: videoID,
		UserID:  authorID,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(favoriteModel).Error; err != nil {
			return err
		}
		if err := tx.Exec("update videos set favorite_count = favorite_count+1 where video_id = ?", videoID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f favoriteDealer) QueryIsLiked(authorID int64, videoId int64) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from favorites where user_id = ? and video_id = ? limit 1", authorID, videoId).Scan(res).Error
	return *res == 1, err
}

func DealFavorite() favoriteFunc {
	return &favoriteDealer{}
}
