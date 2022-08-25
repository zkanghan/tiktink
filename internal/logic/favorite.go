package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
)

func GetIsLiked(userID int64, videoID int64) (bool, error) {
	return mysql.DealFavorite().QueryIsLiked(userID, videoID)
}

func DoFavorite(userID int64, videoID int64) error {
	return mysql.DealFavorite().DoFavorite(userID, videoID)
}

func CancelFavorite(userID int64, videoID int64) error {
	return mysql.DealFavorite().CancelFavorite(userID, videoID)
}

func GetFavoriteList(userID int64) ([]*model.VideoMSG, error) {
	videoMsgS, err := mysql.DealFavorite().QueryFavoriteList(userID)
	if err != nil {
		return nil, err
	}
	for _, videoMsg := range videoMsgS {
		followed, err := mysql.DealFollow().QueryIsFollow(userID, videoMsg.UserMSG.ID)
		if err != nil {
			return nil, err
		}
		liked, err := mysql.DealFavorite().QueryIsLiked(userID, videoMsg.ID)
		if err != nil {
			return nil, err
		}
		videoMsg.UserMSG.IsFollow = followed
		videoMsg.IsFavorite = liked
	}
	return videoMsgS, nil
}
