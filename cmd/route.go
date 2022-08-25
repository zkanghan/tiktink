package main

import (
	"tiktink/internal/controller"
	"tiktink/internal/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/tiktink")
	//视频流
	apiRouter.GET("/feed/", controller.Feed)
	//用户模块
	apiRouter.GET("/user/", middleware.JWTAuthMiddleware, controller.UserInformation)
	apiRouter.POST("/user/login/", controller.UserLogin)
	apiRouter.POST("/user/register/", controller.UserRegister)
	//  视频发布
	apiRouter.POST("/publish/action/", middleware.JWTAuthMiddleware, controller.PublishVideo)
	apiRouter.GET("/publish/list/", middleware.JWTAuthMiddleware, controller.PublishList)

	//  社交模块接口
	apiRouter.POST("/relation/action", middleware.JWTAuthMiddleware, controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JWTAuthMiddleware, controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.JWTAuthMiddleware, controller.FansList)

	//  点赞
	apiRouter.POST("/favorite/action/", middleware.JWTAuthMiddleware, controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JWTAuthMiddleware, controller.FavoriteList)
	//评论
	apiRouter.POST("/comment/action/", middleware.JWTAuthMiddleware, controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.JWTAuthMiddleware, controller.CommentList)
}
