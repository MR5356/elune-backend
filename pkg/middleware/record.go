package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Record() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()

		defer func() {
			if path == "/api/v1/health" {
				return
			}
			cost := time.Since(start)
			httpCode := c.Writer.Status()
			clientIP := c.ClientIP()
			clientUserAgent := c.Request.UserAgent()
			method := c.Request.Method

			entry := logrus.WithFields(logrus.Fields{
				"cost":   cost,
				"method": method,
				"code":   httpCode,
				"ip":     clientIP,
				"path":   path,
			})

			if len(c.Errors) > 0 {
				entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
			} else {
				msg := fmt.Sprintf("user-agent: %s", clientUserAgent)
				switch {
				case httpCode >= http.StatusInternalServerError:
					entry.Error(msg)
				case httpCode >= http.StatusBadRequest:
					entry.Warn(msg)
				default:
					entry.Info(msg)
				}
			}
		}()
		c.Next()
	}
}
