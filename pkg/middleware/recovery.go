package middleware

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("gin panic: %+v", err)
				debug.PrintStack()
				response.New(c, http.StatusInternalServerError, response.CodeUnknownError, response.MsgUnknownError, nil)
			}
		}()
		c.Next()
	}
}
