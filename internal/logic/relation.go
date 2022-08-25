package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
)

//判断是系统错误还是用户侧错误

func GetIsFollowed(UserID int64, ToUserID int64) (bool, error) {
	return mysql.DealFollow().QueryIsFollow(UserID, ToUserID)
}

func DoFollow(UserID int64, ToUserID int64) error {
	return mysql.DealFollow().DoFollow(UserID, ToUserID)
}

func DoCancelFollow(UserID int64, ToUserID int64) error {
	return mysql.DealFollow().DoCancelFollow(UserID, ToUserID)
}

func GetFollowList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs, err := mysql.DealFollow().QueryFollowList(req)
	if err != nil {
		return nil, err
	}
	for _, userMsg := range userMSGs {
		followed, err := GetIsFollowed(aUserID, userMsg.ID)
		if err != nil {
			return nil, err
		}
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}

func GetFansList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs, err := mysql.DealFollow().QueryFansList(req)
	if err != nil {
		return nil, err
	}
	for _, userMsg := range userMSGs {
		followed, err := GetIsFollowed(aUserID, userMsg.ID)
		if err != nil {
			return nil, err
		}
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}
