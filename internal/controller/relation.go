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

const (
	doFollow     int8 = 1
	cancelFollow int8 = 2
)

// RelationAction 关注操作接口
func RelationAction(c *gin.Context) {
	req := &model.FollowActionReq{}
	if err := c.ShouldBind(req); err != nil { //绑定参数错误
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
	userID := c.GetString(middleware.CtxUserIDtxKey)
	if userID == req.ToUserID { // 不允许自己关注自己
		response.Error(c, http.StatusBadRequest, code.FollowSelf)
		return
	}
	background := tracer.Background()

	//  若关注的对象不存在
	toUserExist, err := logic.NewUserDealer(background).GetUserExistByID(req.ToUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if toUserExist == false {
		response.Error(c, http.StatusBadRequest, code.UserNotExist)
		return
	}
	//  若用户正常使用，是无法做出重复关注或重复取消关注操作的，对这类请求直接返回操作失败即可
	isFollowed, err := logic.NewRelationDealer(background).GetIsFollowed(userID, req.ToUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if req.ActionType == doFollow {
		//若已经关注
		if isFollowed {
			response.Error(c, http.StatusBadRequest, code.RepeatFollow)
			return
		}
		//真正进行关注操作
		if err := logic.NewRelationDealer(background).DoFollow(userID, req.ToUserID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.PrintWithStack(err)
			return
		}
		c.JSON(http.StatusOK, &response.RESP{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	} else if req.ActionType == cancelFollow {
		//若已经取消关注
		if !isFollowed {
			response.Error(c, http.StatusBadRequest, code.RepeatUnFollow)
			return
		}
		//真正进行取消关注操作
		if err := logic.NewRelationDealer(background).DoCancelFollow(userID, req.ToUserID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.PrintWithStack(err)
			return
		}
		c.JSON(http.StatusOK, &response.RESP{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	} else { //action_type 不是1 也不是2
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
}

func badFollowResponse(c *gin.Context, code code.ResCode) {
	c.JSON(http.StatusOK, &model.UserInfoResponse{
		StatusCode: code,
		StatusMsg:  code.MSG(),
	})
}

func FollowList(c *gin.Context) {
	req := &model.FollowListReq{}
	if err := c.ShouldBind(req); err != nil {
		badFollowResponse(c, code.InvalidParam)
		return
	}
	background := tracer.Background()

	userExist, err := logic.NewUserDealer(background).GetUserExistByID(req.UserID)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if !userExist {
		badFollowResponse(c, code.UserNotExist)
		return
	}
	// 下面开始查询列表，获取账号主人的ID
	currentUserID := c.GetString(middleware.CtxUserIDtxKey)
	userMSG, err := logic.NewRelationDealer(background).GetFollowList(currentUserID, req)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	c.JSON(http.StatusOK, &model.FollowListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userMSG,
	})
}

func FansList(c *gin.Context) {
	req := &model.FollowListReq{}
	if err := c.ShouldBind(req); err != nil {
		badFollowResponse(c, code.InvalidParam)
		return
	}
	background := tracer.Background()

	userExist, err := logic.NewUserDealer(background).GetUserExistByID(req.UserID)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	if !userExist {
		badFollowResponse(c, code.UserNotExist)
		return
	}
	// 下面开始查询粉丝列表，获取账号主人的ID
	aUserID := c.GetString(middleware.CtxUserIDtxKey)
	userMSG, err := logic.NewRelationDealer(background).GetFansList(aUserID, req)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.PrintWithStack(err)
		return
	}
	c.JSON(http.StatusOK, &model.FollowListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userMSG,
	})
}
