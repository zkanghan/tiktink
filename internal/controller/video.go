package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/internal/response"
	"tiktink/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PublishVideo(c *gin.Context) {
	//获取文件
	title := c.Query("title")
	videoFile, err := c.FormFile("data")
	if err != nil {
		logger.L.Error("获取文件出错：", zap.Error(err))
		response.Error(c, http.StatusBadRequest, code.InvalidFile)
		return
	}
	video := &model.PublishVideoReq{
		Title: title,
		Data:  videoFile,
	}
	//  获取用户id
	userID := c.GetInt64(middleware.CtxUserIDtxKey)
	//  交付逻辑层处理
	if err := logic.PublishVideo(c, video, userID); err != nil {
		logger.L.Error("文件上传失败：", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		return
	}
	c.JSON(http.StatusOK, response.RESP{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	})
}

func badPublishListResp(c *gin.Context, resCode code.ResCode) {
	c.JSON(http.StatusOK, model.PublishListResp{
		StatusCode: resCode,
		StatusMsg:  resCode.MSG(),
	})
}

func PublishList(c *gin.Context) {
	var req model.PublishListReq
	if err := c.ShouldBind(&req); err != nil {
		badPublishListResp(c, code.ServeBusy)
		return
	}
	userID := c.GetInt64(middleware.CtxUserIDtxKey)
	videoList, err := logic.GetVideoList(userID, req.UserID)
	if err != nil {
		badPublishListResp(c, code.ServeBusy)
		logger.L.Error("获取视频列表错误：", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, model.PublishListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	})
}
