package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/internal/response"
	"tiktink/pkg/logger"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	doFollow     int8 = 1
	cancelFollow int8 = 2
)

// RelationAction 关注操作接口
func RelationAction(c *gin.Context) {
	req := &model.FollowActionReq{}
	if err := c.ShouldBindQuery(req); err != nil { //绑定参数错误
		response.Error(c, http.StatusBadRequest, code.InvalidParam)
		return
	}
	userID := c.GetInt64(middleware.CtxUserIDtxKey)
	if userID == req.ToUserID { // 不允许自己关注自己
		response.Error(c, http.StatusBadRequest, code.FollowSelf)
		return
	}
	//  若关注的对象不存在
	toUserExist, err := logic.GetUserExistByID(req.ToUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.L.Error("查询用户是否存在失败：", zap.Error(err))
		return
	}
	if toUserExist == false {
		response.Error(c, http.StatusBadRequest, code.UserNotExist)
		return
	}
	//  若用户正常使用，是无法做出重复关注或重复取消关注操作的，对这类请求直接返回操作失败即可
	isFollowed, err := logic.GetIsFollowed(userID, req.ToUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		logger.L.Error("查询社交是否关注失败：", zap.Error(err))
		return
	}
	if req.ActionType == doFollow {
		//若已经关注
		if isFollowed {
			response.Error(c, http.StatusBadRequest, code.RepeatFollow)
			return
		}
		//真正进行关注操作
		if err := logic.DoFollow(userID, req.ToUserID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.L.Error("关注用户操作失败：", zap.Error(err))
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
		if err := logic.DoCancelFollow(userID, req.ToUserID); err != nil {
			response.Error(c, http.StatusInternalServerError, code.ServeBusy)
			logger.L.Error("取消关注用户操作失败：", zap.Error(err))
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

	userExist, err := logic.GetUserExistByID(req.UserID)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.L.Error("查询用户是否存在失败", zap.Error(err))
		return
	}
	if !userExist {
		badFollowResponse(c, code.UserNotExist)
		return
	}
	// 下面开始查询列表，获取账号主人的ID
	aUserID := c.GetInt64(middleware.CtxUserIDtxKey)
	userMSG, err := logic.GetFollowList(aUserID, req)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		logger.L.Error("查询关注列表出错：", zap.Error(err))
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
	userExist, err := logic.GetUserExistByID(req.UserID)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		zap.L().Error("查询用户是否存在失败", zap.Error(err))
		return
	}
	if !userExist {
		badFollowResponse(c, code.UserNotExist)
		return
	}
	// 下面开始查询粉丝列表，获取账号主人的ID
	aUserID := c.GetInt64(middleware.CtxUserIDtxKey)
	userMSG, err := logic.GetFansList(aUserID, req)
	if err != nil {
		badFollowResponse(c, code.ServeBusy)
		zap.L().Error("查询关粉丝列表出错：", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, &model.FollowListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userMSG,
	})
}
