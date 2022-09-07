package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tracer"
)

type favoriteDealer struct {
	Context *tracer.TraceCtx
}
type favoriteFunc interface {
	GetIsLiked(userID int64, videoID int64) (bool, error)
}

var _ favoriteFunc = &favoriteDealer{}

func NewFavoriteDealer(ctx *tracer.TraceCtx) *favoriteDealer {
	return &favoriteDealer{
		Context: ctx,
	}
}

func (f *favoriteDealer) GetIsLiked(userID int64, videoID int64) (bool, error) {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).QueryIsLiked(userID, videoID)
}

func (f *favoriteDealer) DoFavorite(userID int64, videoID int64) error {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).DoFavorite(userID, videoID)
}

func (f *favoriteDealer) CancelFavorite(userID int64, videoID int64) error {
	f.Context.TraceCaller()
	return mysql.NewFavoriteDealer(f.Context).CancelFavorite(userID, videoID)
}

func (f *favoriteDealer) GetFavoriteList(userID int64) ([]*model.VideoMSG, error) {
	f.Context.TraceCaller()
	videoMsgS, err := mysql.NewFavoriteDealer(f.Context).QueryFavoriteList(userID)
	if err != nil {
		return nil, err
	}

	//  todo: 把循环去掉改为一次查询

	for _, videoMsg := range videoMsgS {
		followed, err := mysql.NewRelationDealer(f.Context).QueryIsFollow(userID, videoMsg.UserMSG.ID)
		if err != nil {
			return nil, err
		}
		liked, err := mysql.NewFavoriteDealer(f.Context).QueryIsLiked(userID, videoMsg.ID)
		if err != nil {
			return nil, err
		}
		videoMsg.UserMSG.IsFollow = followed
		videoMsg.IsFavorite = liked
	}
	return videoMsgS, nil
}
