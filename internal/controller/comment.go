package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/pkg/logger"
	"tiktink/pkg/tracer"

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
	backgroundCTX := tracer.Background() //新建上下文并追踪本函数
	//  根据视频id判断视频是否存在

	videoExist, err := logic.NewVideoDealer(backgroundCTX).GetIsVideoExist(req.VideoID)
	if err != nil {
		logger.PrintWithStack(err)
		badCommentActionResp(c, code.ServeBusy)
		return
	}
	if !videoExist {
		badCommentActionResp(c, code.VideoNotExist)
		return
	}
	userID := c.GetString(middleware.CtxUserIDtxKey)
	switch req.ActionType {
	case releaseComment:
		if req.CommentText == "" { //empty comment
			badCommentActionResp(c, code.InvalidParam)
			return
		}

		commentMsg, err := logic.NewCommentDealer(backgroundCTX).ReleaseComment(req, userID)
		if err != nil {
			badCommentActionResp(c, code.ServeBusy)
			logger.PrintWithStack(err)
			return
		}
		c.JSON(http.StatusOK, &model.CommentActionResp{
			StatusCode: 0,
			StatusMsg:  "success",
			Comment:    commentMsg,
		})

	case deleteComment:
		deleteSuccess, err := logic.NewCommentDealer(backgroundCTX).DeleteComment(req.VideoID, req.CommentID, userID)
		if err != nil {
			badCommentActionResp(c, code.ServeBusy)
			logger.PrintWithStack(err)
			return
		}
		if !deleteSuccess {
			badCommentActionResp(c, code.InvalidParam)
			return
		}
		c.JSON(http.StatusOK, &model.CommentActionResp{
			StatusCode: code.Success,
			StatusMsg:  code.Success.MSG(),
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
	background := tracer.Background()

	userID := c.GetString(middleware.CtxUserIDtxKey)
	commentList, err := logic.NewCommentDealer(background).GetCommentList(*req, userID)
	if err != nil {
		badCommentListResp(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	c.JSON(http.StatusOK, model.CommentListResp{
		StatusCode:  code.Success,
		StatusMsg:   code.Success.MSG(),
		CommentList: commentList,
	})
}
