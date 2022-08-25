package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/pkg/jwt"
	"tiktink/pkg/logger"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func badFeedResp(c *gin.Context, httpStatus int, code code.ResCode) {
	c.JSON(httpStatus, &model.FeedResp{
		StatusCode: code,
		StatusMsg:  code.MSG(),
	})
}

// Feed 视频流不限制用户登录状态，但只允许token为空或合法2种状态，其余状态会被当成错误返回
func Feed(c *gin.Context) {
	var req model.FeedReq
	if err := c.ShouldBind(&req); err != nil {
		badFeedResp(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
	userLogin := false
	if req.Token != "" { //若token不为空则校验token是否合法，合法则将用户id放入上下文
		mc, valid, err := jwt.ParseToken(req.Token)
		if err != nil {
			badFeedResp(c, http.StatusInternalServerError, code.ServeBusy)
			logger.L.Error("解析token错误：", zap.Error(err))
			return
		}
		if !valid { //token不合法直接返回响应
			badFeedResp(c, http.StatusBadRequest, code.NeedLogin)
			return
		}
		userLogin = true
		c.Set(middleware.CtxUserIDtxKey, mc.UserID)
	}
	//  设置时间限制
	latestTime := time.Now().Unix()
	if req.LatestTime != 0 {
		latestTime = req.LatestTime
	}

	var nextTime time.Time
	var videoList []*model.VideoMSG
	var err error
	switch userLogin {
	case true:
		userID := c.GetInt64(middleware.CtxUserIDtxKey)
		videoList, nextTime, err = logic.GetFeed(&userID, latestTime)
		if err != nil {
			logger.L.Error("获取视频流错误：", zap.Error(err))
			badFeedResp(c, http.StatusInternalServerError, code.ServeBusy)
			return
		}
	case false:
		videoList, nextTime, err = logic.GetFeed(nil, time.Now().Unix())
		if err != nil {
			logger.L.Error("获取视频流错误：", zap.Error(err))
			badFeedResp(c, http.StatusInternalServerError, code.ServeBusy)
			return
		}
	}
	c.JSON(http.StatusOK, &model.FeedResp{
		StatusCode: 0,
		StatusMsg:  "success",
		NextTime:   nextTime.Unix(),
		VideoList:  videoList,
	})
}
