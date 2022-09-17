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

var _ relationFunc = &relationDealer{}

func NewRelationDealer(ctx *tracer.TraceCtx) *relationDealer {
	return &relationDealer{
		Context: ctx,
	}
}

func (r *relationDealer) GetIsFollowed(UserID int64, ToUserID int64) (bool, error) {
	r.Context.TraceCaller()
	return mysql.NewRelationDealer(r.Context).QueryIsFollow(UserID, ToUserID)
}

func (r *relationDealer) DoFollow(UserID int64, ToUserID int64) error {
	r.Context.TraceCaller()
	return mysql.NewRelationDealer(r.Context).DoFollow(UserID, ToUserID)
}

func (r *relationDealer) DoCancelFollow(UserID int64, ToUserID int64) error {
	r.Context.TraceCaller()
	return mysql.NewRelationDealer(r.Context).DoCancelFollow(UserID, ToUserID)
}

func (r *relationDealer) GetFollowList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error) {
	r.Context.TraceCaller()
	userMSGs, err := mysql.NewRelationDealer(r.Context).QueryFollowList(req)
	if err != nil {
		return []*model.UserMSG{}, err
	}
	// 组装用户ID切片
	var toUserIDs []int64
	for _, user := range userMSGs {
		toUserIDs = append(toUserIDs, user.UserID)
	}
	//  查询user关注了切片中的哪些人
	followedUserIDs, err := mysql.NewRelationDealer(r.Context).QueryListIsFollow(aUserID, toUserIDs)
	//  切片转map便于判断
	followedUserIDsMap := tools.SliceIntToSet(followedUserIDs)
	for _, userMsg := range userMSGs {
		_, followed := followedUserIDsMap[userMsg.UserID]
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}

func (r *relationDealer) GetFansList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error) {
	r.Context.TraceCaller()
	userMSGs, err := mysql.NewRelationDealer(r.Context).QueryFansList(req)
	if err != nil {
		return []*model.UserMSG{}, err
	}

	// 组装用户ID切片
	var toUserIDs []int64
	for _, user := range userMSGs {
		toUserIDs = append(toUserIDs, user.UserID)
	}
	//  查询user关注了切片中的哪些人
	followedUserIDs, err := mysql.NewRelationDealer(r.Context).QueryListIsFollow(aUserID, toUserIDs)
	//  切片转map便于判断
	followedUserIDsMap := tools.SliceIntToSet(followedUserIDs)
	for _, userMsg := range userMSGs {
		_, followed := followedUserIDsMap[userMsg.UserID]
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}
