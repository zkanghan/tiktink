package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tools"
	"tiktink/pkg/tracer"
)

//判断是系统错误还是用户侧错误

type relationDealer struct {
	Context *tracer.TraceCtx
}

type relationFunc interface {
	GetIsFollowed(UserID int64, ToUserID int64) (bool, error)
	DoFollow(UserID int64, ToUserID int64) error

	DoCancelFollow(UserID int64, ToUserID int64) error
	GetFollowList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error)
	GetFansList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error)
}

//var _ relationFunc = &relationDealer{}

func NewRelationDealer(ctx *tracer.TraceCtx) *relationDealer {
	return &relationDealer{
		Context: ctx,
	}
}

func (r *relationDealer) GetIsFollowed(UserID string, ToUserID string) (bool, error) {
	return mysql.NewRelationDealer().QueryIsFollow(UserID, ToUserID)
}

func (r *relationDealer) DoFollow(UserID string, ToUserID string) error {
	return mysql.NewRelationDealer().DoFollow(UserID, ToUserID)
}

func (r *relationDealer) DoCancelFollow(UserID string, ToUserID string) error {
	return mysql.NewRelationDealer().DoCancelFollow(UserID, ToUserID)
}

func (r *relationDealer) GetFollowList(aUserID string, req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs, err := mysql.NewRelationDealer().QueryFollowList(req)
	if err != nil {
		return []*model.UserMSG{}, err
	}
	// 组装用户ID切片
	var toUserIDs []string
	for _, user := range userMSGs {
		toUserIDs = append(toUserIDs, user.UserID)
	}
	//  查询user关注了切片中的哪些人
	followedUserIDs, err := mysql.NewRelationDealer().QueryListIsFollow(aUserID, toUserIDs)
	//  切片转map便于判断
	followedUserIDsMap := tools.SliceIntToSet(followedUserIDs)
	for _, userMsg := range userMSGs {
		_, followed := followedUserIDsMap[userMsg.UserID]
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}

func (r *relationDealer) GetFansList(aUserID string, req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs, err := mysql.NewRelationDealer().QueryFansList(req)
	if err != nil {
		return []*model.UserMSG{}, err
	}

	// 组装用户ID切片
	var toUserIDs []string
	for _, user := range userMSGs {
		toUserIDs = append(toUserIDs, user.UserID)
	}
	//  查询user关注了切片中的哪些人
	followedUserIDs, err := mysql.NewRelationDealer().QueryListIsFollow(aUserID, toUserIDs)
	//  切片转map便于判断
	followedUserIDsMap := tools.SliceIntToSet(followedUserIDs)
	for _, userMsg := range userMSGs {
		_, followed := followedUserIDsMap[userMsg.UserID]
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}
