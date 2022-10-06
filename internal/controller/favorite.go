package controller

import (
	"fmt"
	"net/http"
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

func FavoriteAction(c *gin.Context) {
	req := new(model.FavoriteActionReq)
	if err := c.ShouldBind(req); err != nil {
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
	background := tracer.Background().TraceCaller()
	// 查询是否点赞了
	userID := c.GetString(middleware.CtxUserIDtxKey)
	liked, err := logic.NewFavoriteDealer(background).GetIsLiked(userID, req.VideoID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintLogWithCTX("查询点赞错误:", err, background)
		return
	}
	if req.ActionType == doLike {
		if liked {
			response.Error(c, http.StatusBadRequest, code.RepeatLiked)
			return
		}
		if err := logic.NewFavoriteDealer(background.Clear().TraceCaller()).DoFavorite(userID, req.VideoID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.PrintLogWithCTX("用户点赞操作失败:", err, background)
			return
		}
	} else if req.ActionType == cancelLike {
		if !liked {
			response.Error(c, http.StatusBadRequest, code.RepeatUnLiked)
			return
		}
		if err := logic.NewFavoriteDealer(background.Clear().TraceCaller()).CancelFavorite(userID, req.VideoID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.PrintLogWithCTX("用户取消点赞失败:", err, background)
			return
		}
	} else {
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
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
	videoList, err := logic.NewFavoriteDealer(background.Clear().TraceCaller()).GetFavoriteList(req.UserID)
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
