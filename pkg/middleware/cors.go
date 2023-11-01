package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWebSockets:  true,
	})
}
