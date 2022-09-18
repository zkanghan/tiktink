package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"
	"time"
)

type feedFunc interface {
	QueryFeedWithTime(latestTime string) ([]*model.VideoMSG, error)
	QueryLatestTimeByID(videoID string) (time.Time, error)
}

type feedDealer struct {
	Context *tracer.TraceCtx
}

func (f *feedDealer) QueryLatestTimeByID(videoID string) (time.Time, error) {
	f.Context.TraceCaller()
	var latestTime time.Time
	err := db.Raw("select `create_time` from `videos` where video_id = ?", videoID).Scan(&latestTime).Error
	if err != nil {
		return time.Now(), err
	}
	return latestTime, nil
}

// QueryFeedWithTime 按视频发布时间从新到久排序
func (f *feedDealer) QueryFeedWithTime(latestTime string) ([]*model.VideoMSG, error) {
	f.Context.TraceCaller()
	var videoMsgs []*model.VideoMSG
	err := db.Raw("select `video_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`,"+
		"`user_id`,`user_name`,`follow_count`,`follower_count`"+
		"from `users` inner join `videos` on `videos`.`author_id`= `user_id`"+
		"where `videos`.`create_time` < ? order by `videos`.`create_time` desc limit 30", latestTime).Scan(&videoMsgs).Error
	if err != nil {
		return nil, err
	}
	return videoMsgs, nil
}

func NewFeedDealer(ctx *tracer.TraceCtx) feedFunc {
	return &feedDealer{
		Context: ctx,
	}
}
