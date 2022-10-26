package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/internal/response"
	"tiktink/pkg/logger"
	"tiktink/pkg/tracer"

	"github.com/gin-gonic/gin"
)

const (
	doLike     int8 = 1
	cancelLike int8 = 2
)

// FavoriteAction 点赞行为
func FavoriteAction(c *gin.Context) {
	req := new(model.FavoriteActionReq)
	if err := c.ShouldBind(req); err != nil {
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
	background := tracer.Background()
	// 使用redis不要查询直接set
	userID := c.GetString(middleware.CtxUserIDtxKey)

	//  安全验证videoID是存在的
	exsit, err := logic.NewVideoDealer(background).GetIsVideoExist(req.VideoID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if !exsit {
		c.JSON(http.StatusBadRequest, response.RESP{
			StatusCode: 0,
			StatusMsg:  "视频不存在",
		})
	}
	err = logic.NewFavoriteDealer(background).SetRedisKey(userID, req.VideoID, strconv.Itoa(int(req.ActionType)))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	c.JSON(http.StatusOK, response.RESP{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

func badRespFavoriteList(c *gin.Context, code code.ResCode) {
	c.JSON(http.StatusOK, &model.FavoriteListResp{
		StatusCode: code,
		StatusMsg:  code.MSG(),
	})
}

func FavoriteList(c *gin.Context) {
	req := new(model.FavoriteListReq)
	if err := c.ShouldBind(req); err != nil {
		badRespFavoriteList(c, code.InvalidParam)
		fmt.Println(err)
		return
	}
	background := tracer.Background()
	exist, err := logic.NewUserDealer(background).GetUserExistByID(req.UserID)
	if err != nil {
		badRespFavoriteList(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if !exist {
		badRespFavoriteList(c, code.UserNotExist)
		return
	}
	videoList, err := logic.NewFavoriteDealer(background).GetMySQLFavoriteList(*req)
	if err != nil {
		badRespFavoriteList(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	c.JSON(http.StatusOK, model.FavoriteListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	})
}
