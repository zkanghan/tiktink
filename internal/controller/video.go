package controller

import (
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

func PublishVideo(c *gin.Context) {
	//获取文件
	title := c.Query("title")
	videoFile, err := c.FormFile("data")
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.InvalidFile)
		return
	}
	video := &model.PublishVideoReq{
		Title: title,
		Data:  videoFile,
	}
	//  获取用户id
	userID := c.GetString(middleware.CtxUserIDtxKey)
	//  交付逻辑层处理
	background := tracer.Background().TraceCaller()

	if err := logic.NewVideoDealer(background).PublishVideo(video, userID); err != nil {
		logger.PrintLogWithCTX("文件上传失败:", err, background)
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
	userID := c.GetString(middleware.CtxUserIDtxKey)
	background := tracer.Background().TraceCaller()
	videoList, err := logic.NewVideoDealer(background).GetVideoList(userID, req)
	if err != nil {
		badPublishListResp(c, code.ServeBusy)
		logger.PrintLogWithCTX("获取视频列表错误:", err, background)
		return
	}
	c.JSON(http.StatusOK, model.PublishListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	})
}
