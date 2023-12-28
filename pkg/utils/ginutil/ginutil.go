package ginutil

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetId(c *gin.Context) (id uint, ok bool) {
	idParam := c.Query("id")
	i, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || i == 0 {
		return 0, false
	}
	return uint(i), true
}

func GetToken(ctx *gin.Context) string {
	tokenString := ctx.GetHeader("Authorization")
	if len(tokenString) == 0 {
		tokenString, _ = ctx.Cookie("token")
	}
	return tokenString
}
