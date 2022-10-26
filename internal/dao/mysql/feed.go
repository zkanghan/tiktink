package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"
	"time"

	"github.com/pkg/errors"
)

type feedFunc interface {
	QueryFeedWithTime(latestTime string) ([]*model.VideoMSG, error)
	QueryLatestTimeByID(videoID string) (time.Time, error)
}

type feedDealer struct{}

func NewFeedDealer() feedFunc {
	return &feedDealer{}
}

func (f *feedDealer) QueryLatestTimeByID(videoID string) (time.Time, error) {
	var latestTime time.Time
	err := db.Raw("select `create_time` from `videos` where video_id = ?", videoID).Scan(&latestTime).Error
	if err != nil {
		return time.Now(), errors.Wrap(err, tracer.FormatParam(videoID))
	}
	return latestTime, nil
}

// QueryFeedWithTime 按视频发布时间从新到久排序
func (f *feedDealer) QueryFeedWithTime(latestTime string) ([]*model.VideoMSG, error) {
	var videoMsgs []*model.VideoMSG
	err := db.Raw("select `video_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`,"+
		"`user_id`,`user_name`,`follow_count`,`follower_count`"+
		"from `users` inner join `videos` on `videos`.`author_id`= `user_id`"+
		"where `videos`.`create_time` < ? order by `videos`.`create_time` desc limit 30", latestTime).Scan(&videoMsgs).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(latestTime))
	}
	return videoMsgs, nil
}
