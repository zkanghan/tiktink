package response

import (
	"net/http"
	"tiktink/internal/code"

	"github.com/gin-gonic/gin"
)

type RESP struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	Data       interface{}  `json:"data,omitempty"`
}

func Error(c *gin.Context, httpStatus int, code code.ResCode) {
	c.JSON(http.StatusOK, &RESP{
		StatusCode: code,
		StatusMsg:  code.MSG(),
		Data:       nil,
	})
}
