package middleware

import (
	"net/http"
	"tiktink/internal/code"
	"tiktink/internal/response"
	"tiktink/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const CtxUserIDtxKey = "userID"

// JWTAuthMiddleware 基于JWT的用户状态认证中间件
func JWTAuthMiddleware(c *gin.Context) {

	token := c.Query("token")

	if token == "" {
		response.Error(c, http.StatusBadRequest, code.NeedLogin)
		c.Abort()
		return
	}
	mc, tokenValid, err := jwt.ParseToken(token)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.ServeBusy)
		c.Abort()
		return
	}
	if !tokenValid {
		response.Error(c, http.StatusBadRequest, code.NeedLogin)
		c.Abort()
		return
	}
	// 将当前请求的userID信息保存到请求的上下文c上
	c.Set(CtxUserIDtxKey, mc.UserID)
	c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息

}
