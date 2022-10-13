package main

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/dao/redis"
	"tiktink/pkg/logger"
	"tiktink/pkg/setting"
	"tiktink/pkg/snowid"

	"github.com/gin-gonic/gin"
)

//    时长1年的测试token,user_id为5：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NSwidXNlcm5hbWUiOiJ6aGFuZ2toIiwiZXhwIjoxNjg5NDE3OTA2fQ.fFv3Mvusn_HgkLJSoceYug_Ae9HDUbcOa_b-PiNtHUU
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
}
