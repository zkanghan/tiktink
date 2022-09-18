package controller

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/logic"
	"tiktink/internal/middleware"
	"tiktink/internal/model"
	"tiktink/pkg/jwt"
	"tiktink/pkg/logger"
	"tiktink/pkg/tracer"

	"github.com/gin-gonic/gin"
)

func badLoginResponse(c *gin.Context, code code.ResCode, Msg string) {
	c.JSON(http.StatusOK, &model.LoginResponse{
		StatusCode: code,
		StatusMsg:  Msg,
	})
}

func badRegisterResponse(c *gin.Context, code code.ResCode, Msg string) {
	c.JSON(http.StatusOK, &model.RegisterResponse{
		StatusCode: code,
		StatusMsg:  Msg,
	})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	user := &model.UserRequest{}
	//   绑定参数
	if err := c.ShouldBind(user); err != nil {
		badLoginResponse(c, code.InvalidParam, code.InvalidParam.MSG())
		return
	}
	//  查询数据库验证账号
	background := tracer.Background().TraceCaller()
	right, id, err := logic.NewUserDealer(background).CheckUser(user.UserName, user.Password)
	if err != nil { //运行异常
		badLoginResponse(c, code.ServeBusy, err.Error())
		return
	}
	if !right {
		badLoginResponse(c, code.WrongPassword, code.WrongPassword.MSG())
		return
	}

	token, err := jwt.GenToken(id, user.UserName)
	if err != nil {
		logger.PrintLog("生成token错误:", err)
		badLoginResponse(c, code.ServeBusy, code.ServeBusy.MSG())
		return
	}
	// 返回正确响应
	c.JSON(http.StatusOK, &model.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserID:     id,
		Token:      token,
	})
}

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	user := &model.UserRequest{}
	if err := c.ShouldBind(user); err != nil {
		badRegisterResponse(c, code.InvalidParam, code.InvalidParam.MSG())
		return
	}
	background := tracer.Background().TraceCaller()
	userExit, err := logic.NewUserDealer(background).GetUserExistByName(user.UserName)
	if err != nil {
		logger.PrintLogWithCTX("查询用户存在出错:", err, background)
		badRegisterResponse(c, code.ServeBusy, code.ServeBusy.MSG())
		return
	}
	if userExit { //用户已存在
		badRegisterResponse(c, code.UserExist, code.UserExist.MSG())
		return
	}
	id, err := logic.NewUserDealer(background.Clear().TraceCaller()).CreateUser(user.UserName, user.Password)
	if err != nil {
		logger.PrintLogWithCTX("创建用户出错", err, background)
		badRegisterResponse(c, code.ServeBusy, code.ServeBusy.MSG())
		return
	}
	c.JSON(http.StatusOK, &model.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserID:     id,
	})
}

func badUserInfoResp(c *gin.Context, resCode code.ResCode) {
	c.JSON(http.StatusOK, &model.UserInfoResponse{
		StatusCode: resCode,
		StatusMsg:  resCode.MSG(),
	})
}

// UserInformation 用户信息接口
func UserInformation(c *gin.Context) {
	req := &model.UserInfoRequest{}
	if err := c.ShouldBind(req); err != nil {
		badUserInfoResp(c, code.InvalidParam)
		return
	}
	background := tracer.Background().TraceCaller()
	userExist, err := logic.NewUserDealer(background).GetUserExistByID(req.UserID)
	if err != nil {
		logger.PrintLogWithCTX("查询用户是否存在失败:", err, background)
		badUserInfoResp(c, code.ServeBusy)
		return
	}
	if !userExist {
		badUserInfoResp(c, code.UserNotExist)
		return
	}
	userID := c.GetString(middleware.CtxUserIDtxKey)
	userMsg, err := logic.NewUserDealer(background.Clear().TraceCaller()).GetUserInformation(req.UserID, userID)
	if err != nil {
		badUserInfoResp(c, code.ServeBusy)
		logger.PrintLogWithCTX("查询用户信息错误：", err, background)
		return
	}
	c.JSON(http.StatusOK, model.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserMSG:    *userMsg,
	})
}
