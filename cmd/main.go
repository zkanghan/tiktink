package main

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/dao/redis"
	"tiktink/internal/timing"
	"tiktink/pkg/logger"
	"tiktink/pkg/setting"
	"tiktink/pkg/snowid"
	"time"

	"github.com/gin-gonic/gin"
)

//user_id为5的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI1IiwiVXNlcm5hbWUiOiJ6aGFuZ2toIiwiZXhwIjoxNjk2NDk5ODE2fQ.QnrurHrjAf6t0AzUDcbl4NrhwgjIQUB3es7SU7iShMY
// 来个6的：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwidXNlcm5hbWUiOiJ6a2giLCJleHAiOjE2OTAyODI4MDB9.TP5S9b5Ka0j0Pmiz1eyjkmqrFO-xgzbeEpk69nWjf2M
func main() {
	r := gin.Default()

	initRouter(r)

	err := r.Run(":8080")
	if err != nil {
		panic("程序启动失败")
	}
}

// 初始化函数
func init() {
	//配置文件初始化,失败抛出异常
	if err := setting.Init("./config/config.yaml"); err != nil {
		panic("配置文件初始化出错:  " + err.Error())
	}

	logger.InitLogger() //初始化日志

	//MySQL配置初始化
	if err := mysql.InitMysql(); err != nil {
		panic("MySQL初始化出错:  " + err.Error())
	}

	if err := redis.InitRedis(); err != nil {
		panic("Redis初始化出错：" + err.Error())
	}

	if err := snowid.Init(); err != nil {
		panic("雪花id生成器初始化出错：" + err.Error())
	}

	// 注册定时任务
	go func() {
		timing.RegisterTask(time.Minute, timing.SyncFavoriteKey)
	}()
}
