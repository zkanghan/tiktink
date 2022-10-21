package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var redisDB *redis.Client

func InitRedis() (err error) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")), // redis地址
		Password: viper.GetString("redis.password"),                                               // redis密码，没有则留空
		DB:       0,                                                                               // 默认数据库，默认是0
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err = redisDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
