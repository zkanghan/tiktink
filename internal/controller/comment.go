package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/pkg/logger"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	releaseComment int8 = 1
	deleteComment  int8 = 2
)

func badCommentActionResp(c *gin.Context, code code.ResCode) {
	c.JSON(http.StatusOK, model.CommentActionResp{
		StatusCode: code,
		StatusMsg:  code.MSG(),
		Comment:    nil,
	})
}

func badCommentListResp(c *gin.Context, code code.ResCode) {
	c.JSON(http.StatusOK, model.FavoriteListResp{
		StatusCode: code,
		StatusMsg:  code.MSG(),
	})
}

func CommentAction(c *gin.Context) {
	req := new(model.CommentActionReq)
	if err := c.ShouldBind(req); err != nil {
		badCommentActionResp(c, code.InvalidParam)
		return
	}
	//  根据视频id判断视频是否存在
	videoExist, err := logic.GetIsVideoExist(req.VideoID)
	if err != nil {
		logger.L.Error("查询视频存在失败：", zap.Error(err))
		badCommentActionResp(c, code.ServeBusy)
		return
	}
	if !videoExist {
		badCommentActionResp(c, code.VideoNotExist)
		return
	}
	userID := c.GetInt64(middleware.CtxUserIDtxKey)
	switch req.ActionType {
	case releaseComment:
		if req.CommentText == "" {
			badCommentActionResp(c, code.InvalidParam)
			return
		}
		commentMsg, err := logic.ReleaseComment(req, userID)
		if err != nil {
			badCommentActionResp(c, code.ServeBusy)
			logger.L.Error("发布评论失败：", zap.Error(err))
			return
		}
		c.JSON(http.StatusOK, &model.CommentActionResp{
			StatusCode: 0,
			StatusMsg:  "success",
			Comment:    commentMsg,
		})

	case deleteComment:
		deleteSuccess, err := logic.DeleteComment(req.VideoID, req.CommentID, userID)
		if err != nil {
			badCommentActionResp(c, code.ServeBusy)
			logger.L.Error("删除评论失败：", zap.Error(err))
			return
		}
		if !deleteSuccess {
			badCommentActionResp(c, code.InvalidParam)
			return
		}
		c.JSON(http.StatusOK, &model.CommentActionResp{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	default:
		badCommentActionResp(c, code.InvalidParam)
		return
	}
}

func CommentList(c *gin.Context) {
	req := new(model.CommentListReq)
	if err := c.ShouldBind(req); err != nil {
		badCommentListResp(c, code.InvalidParam)
		return
	}
	userID := c.GetInt64(middleware.CtxUserIDtxKey)
	commentList, err := logic.GetCommentList(req.VideoID, userID)
	if err != nil {
		badCommentListResp(c, code.ServeBusy)
		logger.L.Error("获取评论列表失败：", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, model.CommentListResp{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: commentList,
	})
}
