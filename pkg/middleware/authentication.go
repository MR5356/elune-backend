package middleware

import (
	"github.com/MR5356/elune-backend/pkg/config"
	"github.com/MR5356/elune-backend/pkg/domain/authentication"
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/MR5356/elune-backend/pkg/utils/ginutil"
	"github.com/gin-gonic/gin"
	"strings"
)

func Authentication(config *config.Config, rbacService *authentication.RBACService, jwtService *authentication.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ginutil.GetToken(ctx)

		user, err := jwtService.ParseToken(tokenString)
		if err != nil || user == nil {
			user = &authentication.User{
				Username: "unknown",
			}
		}
		username := user.Username
		if ok, err := rbacService.HasRoleForUser(username, strings.ReplaceAll(ctx.Request.URL.Path, config.Server.Prefix, ""), ctx.Request.Method); ok && err == nil {
			ctx.Next()
		} else {
			response.Error(ctx, response.CodeForbidden, "权限不足")
			ctx.Abort()
			return
		}
	}
}
