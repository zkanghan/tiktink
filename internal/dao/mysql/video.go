package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"
)

type videoFunc interface {
	PublishVideo(video *model.Video) error
	QueryVideoExist(videoID string) (bool, error)
	QueryVideoByAuthorID(authorID string) ([]*model.VideoMSG, error)
}

type videoDealer struct {
	Context *tracer.TraceCtx
}

func (v *videoDealer) QueryVideoByAuthorID(authorID string) ([]*model.VideoMSG, error) {
	v.Context.TraceCaller()
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

func (v *videoDealer) QueryVideoExist(videoID string) (bool, error) {
	v.Context.TraceCaller()
	cnt := new(int64)
	err := db.Raw("select 1 from videos where video_id = ?", videoID).Scan(cnt).Error
	if err != nil {
		return false, err
	}
	return *cnt == 1, nil
}

func (v *videoDealer) PublishVideo(video *model.Video) error {
	v.Context.TraceCaller()
	err := db.Create(video).Error
	return err
}

func NewVideoDealer(ctx *tracer.TraceCtx) videoFunc {
	return &videoDealer{
		Context: ctx,
	}
}
