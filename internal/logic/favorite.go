package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/dao/redis"
	"tiktink/internal/model"
	"tiktink/pkg/tools"
	"tiktink/pkg/tracer"
)

type favoriteDealer struct {
	Context *tracer.TraceCtx
}
type favoriteFunc interface {
	GetMySQLIsLiked(userID string, videoID string) (bool, error)
}

var _ favoriteFunc = &favoriteDealer{}

func NewFavoriteDealer(ctx *tracer.TraceCtx) *favoriteDealer {
	return &favoriteDealer{
		Context: ctx,
	}
}

func (f *favoriteDealer) GetMySQLIsLiked(userID string, videoID string) (bool, error) {
	return mysql.NewFavoriteDealer().QueryIsLiked(userID, videoID)
}

func (f *favoriteDealer) DoMySQLFavorite(userID string, videoID string) error {
	return mysql.NewFavoriteDealer().DoFavorite(userID, videoID)
}

func (f *favoriteDealer) CancelMySQLFavorite(userID string, videoID string) error {
	return mysql.NewFavoriteDealer().CancelFavorite(userID, videoID)
}

func (f *favoriteDealer) GetMySQLFavoriteList(req model.FavoriteListReq) ([]*model.VideoMSG, error) {
	videoMsgS, err := mysql.NewFavoriteDealer().QueryFavoriteList(req.UserID, req.PageNumber)
	if err != nil {
		return nil, err
	}

	// 获取需要的用户id
	var toUserIDs []string
	for _, video := range videoMsgS {
		toUserIDs = append(toUserIDs, video.UserMSG.UserID)
	}
	//  获取user在toUserID中关注了哪些
	followedUsers, err := mysql.NewRelationDealer().QueryListIsFollow(req.UserID, toUserIDs)
	if err != nil {
		return []*model.VideoMSG{}, err
	}
	followedUserMap := tools.SliceIntToSet(followedUsers)

	for _, videoMsg := range videoMsgS {
		_, followed := followedUserMap[videoMsg.UserMSG.UserID]
		videoMsg.UserMSG.IsFollow = followed
		videoMsg.IsFavorite = true
	}
	return videoMsgS, nil
}

func (f *favoriteDealer) GetRedisFavoriteVal(userID string, videoID string) (model.FavoriteRedis, error) {
	fr := model.FavoriteRedis{
		UserID:  userID,
		VideoID: videoID,
	}
	dealer := redis.NewFavoriteDealer()
	m, err := dealer.GetFavoriteVal(redis.GetFavoriteKey(fr))
	if err != nil {
		return model.FavoriteRedis{}, err
	}
	if err = tools.MapToStruct(m, &fr); err != nil {
		return model.FavoriteRedis{}, err
	}
	return fr, nil
}

func (f *favoriteDealer) SetRedisKey(userID string, videoID string, status string) error {
	dealer := redis.NewFavoriteDealer()
	return dealer.SetFavoriteKey(model.FavoriteRedis{
		UserID:  userID,
		VideoID: videoID,
		Status:  status,
	})
}
