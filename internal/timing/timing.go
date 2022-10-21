package timing

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/dao/redis"
	"tiktink/internal/model"
	"tiktink/pkg/logger"
	"tiktink/pkg/tools"
	"tiktink/pkg/tracer"
	"time"
)

// SyncFavoriteKey 定时任务同步Redis Key
func SyncFavoriteKey() error {
	// 模糊查询获取所有的键
	dealer := redis.NewFavoriteDealer()
	keys, err := dealer.GetFavoriteKeys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		go func(s string) {
			m, err := dealer.GetFavoriteVal(s)
			if err != nil {
				logger.PrintLog("", err)
			}
			// 序列化为 结构体
			var fr model.FavoriteRedis
			if err = tools.MapToStruct(m, &fr); err != nil {
				logger.PrintLog("", err)
			}
			ctx := tracer.Background()
			//  查询是否关注
			mysqlLiked, err := mysql.NewFavoriteDealer(ctx).QueryIsLiked(fr.UserID, fr.VideoID)
			if err != nil {
				logger.PrintLogWithCTX("", err, ctx)
			}
			//  业务逻辑判断
			if fr.Status == redis.Liked && !mysqlLiked { // redis中以点赞而MySQL未点赞
				if err = mysql.NewFavoriteDealer(ctx.Clear()).DoFavorite(fr.UserID, fr.VideoID); err != nil {
					logger.PrintLogWithCTX("", err, ctx)
				}
			} else if fr.Status == redis.Unliked && mysqlLiked { //redis 中取消点赞而MySQL已点赞
				if err = mysql.NewFavoriteDealer(ctx.Clear()).CancelFavorite(fr.UserID, fr.VideoID); err != nil {
					logger.PrintLogWithCTX("", err, ctx)
				}
			}
			// 判断完毕删除键
			if err = dealer.DeleteFavoriteKey(fr); err != nil {
				logger.PrintLog("", err)
			}
		}(key)
	}
	return nil
}

func RegisterTask(t time.Duration, f func() error) {
	ticker := time.NewTicker(t)
	for {
		<-ticker.C
		go func(work func() error) {
			if err := work(); err != nil {
				logger.PrintLog("同步redis与MySQL错误", err)
			}
		}(f)
	}

}
