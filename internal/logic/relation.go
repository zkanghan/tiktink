package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
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
		return nil, err
	}
	//  todo: 把循环去掉改为一次查询

	for _, userMsg := range userMSGs {
		followed, err := r.GetIsFollowed(aUserID, userMsg.ID)
		if err != nil {
			return nil, err
		}
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}

func (r *relationDealer) GetFansList(aUserID int64, req *model.FollowListReq) ([]*model.UserMSG, error) {
	r.Context.TraceCaller()
	userMSGs, err := mysql.NewRelationDealer(r.Context).QueryFansList(req)
	if err != nil {
		return nil, err
	}
	//  todo: 把循环去掉改为一次查询

	for _, userMsg := range userMSGs {
		followed, err := r.GetIsFollowed(aUserID, userMsg.ID)
		if err != nil {
			return nil, err
		}
		userMsg.IsFollow = followed
	}
	return userMSGs, nil
}
