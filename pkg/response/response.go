package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func New(ctx *gin.Context, httpCode int, code int, msg string, data any) {
	ctx.JSON(httpCode, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Success(ctx *gin.Context, data any) {
	New(ctx, http.StatusOK, CodeSuccess, MsgSuccess, data)
}

func Error(ctx *gin.Context, code int, msg string) {
	if len(msg) == 0 {
		if msg, ok := responseMap[code]; ok {
			New(ctx, http.StatusOK, code, msg, nil)
		} else {
			New(ctx, http.StatusOK, code, MsgUnknownError, nil)
		}
	} else {
		New(ctx, http.StatusOK, code, msg, nil)
	}
}
