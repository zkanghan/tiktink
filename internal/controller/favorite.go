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
	background := tracer.Background().TraceCaller()
	// 使用redis不要查询直接set
	userID := c.GetString(middleware.CtxUserIDtxKey)

	//  安全验证videoID是存在的
	exsit, err := logic.NewVideoDealer(background.Clear()).GetIsVideoExist(req.VideoID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintLogWithCTX("查询视频存在失败:", err, background)
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
		logger.PrintLogWithCTX("redis点赞失败:", err, background)
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
	background := tracer.Background().TraceCaller()
	exist, err := logic.NewUserDealer(background).GetUserExistByID(req.UserID)
	if err != nil {
		badRespFavoriteList(c, code.ServeBusy)
		logger.PrintLogWithCTX("查询用户是否存在错误出错:", err, background)
		return
	}
	if !exist {
		badRespFavoriteList(c, code.UserNotExist)
		return
	}
	videoList, err := logic.NewFavoriteDealer(background.Clear().TraceCaller()).GetMySQLFavoriteList(*req)
	if err != nil {
		badRespFavoriteList(c, code.ServeBusy)
		logger.PrintLogWithCTX("获取点赞列表出错:", err, background)
		return
	}
	c.JSON(http.StatusOK, model.FavoriteListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	})
}
