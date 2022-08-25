package mysql

import (
	"tiktink/internal/model"
)

type videoFunc interface {
	PublishVideo(video *model.Video) error
	QueryVideoExist(videoID int64) (bool, error)
	QueryVideoByAuthorID(authorID int64) ([]*model.VideoMSG, error)
}

type videoDealer struct{}

func (v videoDealer) QueryVideoByAuthorID(authorID int64) ([]*model.VideoMSG, error) {
	var videoMsgs []*model.VideoMSG
	err := db.Raw("select `video_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`,"+
		"`user_id`,`user_name`,`follow_count`,`follower_count`"+
		"from `users` inner join `videos` on `videos`.`author_id`= `user_id`"+
		"where `videos`.author_id = ?", authorID).Scan(&videoMsgs).Error
	if err != nil {
		return nil, err
	}
	return videoMsgs, nil
}

func (v videoDealer) QueryVideoExist(videoID int64) (bool, error) {
	cnt := new(int64)
	err := db.Raw("select 1 from videos where video_id = ?", videoID).Scan(cnt).Error
	if err != nil {
		return false, err
	}
	return *cnt == 1, nil
}

func (v videoDealer) PublishVideo(video *model.Video) error {
	err := db.Create(video).Error
	return err
}

func DealVideo() videoFunc {
	return &videoDealer{}
}
