package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"time"
)

// GetFeed 限制返回视频的最新时间要小于latestTime，即用户观看视频的顺序是从新到久
func GetFeed(userID *int64, latestTime int64) ([]*model.VideoMSG, time.Time, error) {
	//  获取videoList
	timeStr := time.Unix(latestTime, 0).Format("2006-01-02 15:04:05")
	videoList, err := mysql.DealFeed().QueryFeedWithTime(timeStr)
	if err != nil {
		return nil, time.Now(), err
	}
	// 判断videoList是否为空
	if len(videoList) == 0 {
		return nil, time.Now(), nil
	}
	// 判断是否关注和点赞
	if userID != nil {
		for _, video := range videoList {
			followed, err := mysql.DealFollow().QueryIsFollow(*userID, video.UserMSG.ID)
			if err != nil {
				return nil, time.Now(), err
			}
			liked, err := mysql.DealFavorite().QueryIsLiked(*userID, video.ID)
			if err != nil {
				return nil, time.Now(), err
			}
			video.UserMSG.IsFollow = followed
			video.IsFavorite = liked
		}
	} else {
		for _, video := range videoList {
			video.IsFollow = false
			video.IsFavorite = false
		}
	}

	//  获得本次视频列表中发布时间最早的
	index := len(videoList)
	newLatestTime, err := mysql.DealFeed().QueryLatestTimeByID(videoList[index-1].ID)
	if err != nil {
		return nil, time.Now(), err
	}
	return videoList, newLatestTime, nil
}
