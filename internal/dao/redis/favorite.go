package redis

import (
	"fmt"
	"tiktink/internal/model"
	"tiktink/pkg/tools"
	"time"
)

const (
	favoriteKey = "favorite"
	Liked       = "1"
	Unliked     = "2"
	expiration  = time.Hour * 2
)

type favoriteFunc interface {
	DeleteFavoriteKey(m model.FavoriteRedis) error
	SetFavoriteKey(fr model.FavoriteRedis) error
	GetFavoriteKeys() ([]string, error)
	GetFavoriteVal(key string) (map[string]string, error)
}

type favoriteDealer struct {
}

var _ favoriteFunc = &favoriteDealer{}

//  key值不变，value改为字典

func NewFavoriteDealer() favoriteDealer {
	return favoriteDealer{}
}

func GetFavoriteKey(m model.FavoriteRedis) string {
	return fmt.Sprintf("%s:userID{%s}:videoID{%s}", favoriteKey, m.UserID, m.VideoID)
}

// DeleteFavoriteKey 删除redis的key
func (f *favoriteDealer) DeleteFavoriteKey(m model.FavoriteRedis) error {
	key := GetFavoriteKey(m)
	return redisDB.Del(key).Err()
}

// SetFavoriteKey 设置对应key的value
func (f *favoriteDealer) SetFavoriteKey(fr model.FavoriteRedis) error {
	m, err := tools.StructToMap(fr)
	if err != nil {
		return err
	}
	key := GetFavoriteKey(fr)
	if err = redisDB.HMSet(key, m).Err(); err != nil {
		return err
	}
	if err = redisDB.Expire(key, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (f *favoriteDealer) GetFavoriteKeys() ([]string, error) {
	return redisDB.Keys(favoriteKey + "*").Result()
}

func (f *favoriteDealer) GetFavoriteVal(key string) (map[string]string, error) {
	return redisDB.HGetAll(key).Result()
}
