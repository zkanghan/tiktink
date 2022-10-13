package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tools"
	"tiktink/pkg/tracer"
)

type favoriteDealer struct {
	Context *tracer.TraceCtx
}
type favoriteFunc interface {
	GetIsLiked(userID string, videoID string) (bool, error)
}

var _ favoriteFunc = &favoriteDealer{}

func NewFavoriteDealer(ctx *tracer.TraceCtx) *favoriteDealer {
	return &favoriteDealer{
		Context: ctx,
	}
}

func (f *favoriteDealer) GetIsLiked(userID string, videoID string) (bool, error) {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).QueryIsLiked(userID, videoID)
}

func (f *favoriteDealer) DoFavorite(userID string, videoID string) error {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).DoFavorite(userID, videoID)
}

func (f *favoriteDealer) CancelFavorite(userID string, videoID string) error {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).CancelFavorite(userID, videoID)
}

func (f *favoriteDealer) GetFavoriteList(req model.FavoriteListReq) ([]*model.VideoMSG, error) {
	f.Context.TraceCaller()
	videoMsgS, err := mysql.NewFavoriteDealer(f.Context).QueryFavoriteList(req.UserID, req.PageNumber)
	if err != nil {
		return nil, err
	}

	// 获取需要的用户id
	var toUserIDs []string
	for _, video := range videoMsgS {
		toUserIDs = append(toUserIDs, video.UserMSG.UserID)
	}
	//  获取user在toUserID中关注了哪些
	followedUsers, err := mysql.NewRelationDealer(f.Context).QueryListIsFollow(req.UserID, toUserIDs)
	if err != nil {
		return []*model.VideoMSG{}, err
	}
	followedUserMap := tools.SliceIntToSet(followedUsers)

	//  todo: 把循环去掉改为一次查询,这里dao的favorite也要重构

	for _, videoMsg := range videoMsgS {
		_, followed := followedUserMap[videoMsg.UserMSG.UserID]
		videoMsg.UserMSG.IsFollow = followed
		videoMsg.IsFavorite = true
	}
	return videoMsgS, nil
}
