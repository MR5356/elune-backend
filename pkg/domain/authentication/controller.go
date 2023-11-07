package authentication

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/MR5356/elune-backend/pkg/utils/ginutil"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	rbacService *RBACService
	jwtService  *JWTService
	userService *Service
}

func NewController(rbacService *RBACService, jwtService *JWTService, userService *Service) *Controller {
	return &Controller{
		rbacService: rbacService,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (c *Controller) handleRefreshToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if len(tokenString) == 0 {
		tokenString, _ = ctx.Cookie("token")
	}

	newTokenString, err := c.jwtService.RefreshToken(tokenString)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, map[string]string{"token": newTokenString})
	}
}

func (c *Controller) handleUpdateUserPassword(ctx *gin.Context) {
	tokenString := ginutil.GetToken(ctx)
	user, err := c.jwtService.ParseToken(tokenString)
	if err != nil {
		ctx.SetCookie("token", "", 0, "", "", false, false)
		response.Error(ctx, response.CodeUnknownError, "用户信息错误")
		return
	} else {
		json := make(map[string]string)
		err := ctx.BindJSON(&json)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
			return
		}

		oldPassword := json["oldPassword"]
		password := json["password"]
		// TODO 和数据库对比，比较安全
		if oldPassword != user.Password {
			response.Error(ctx, response.CodeParamError, "原密码错误")
			return
		}

		user.Password = password
		err = c.userService.persistence.Update(&User{ID: user.ID}, structutil.Struct2Map(user))
		if err != nil {
			response.Error(ctx, response.CodeUnknownError, err.Error())
		} else {
			token, err := c.jwtService.CreateToken(user)
			if err != nil {
				response.Error(ctx, response.CodeUnknownError, err.Error())
				return
			} else {
				ctx.SetCookie("token", token, 0, "", "", false, false)
				response.Success(ctx, map[string]string{"token": token})
			}
		}
	}
}

func (c *Controller) handleLogin(ctx *gin.Context) {
	json := make(map[string]string)
	err := ctx.BindJSON(&json)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	username := json["username"]
	password := json["password"]

	if len(username) == 0 || len(password) == 0 {
		response.Error(ctx, response.CodeParamError, "用户名或密码不能为空")
		return
	}

	logrus.Infof("login username: %s, password: %s", username, password)

	user, err := c.userService.persistence.Detail(&User{Username: username})
	if err != nil {
		response.Error(ctx, response.CodeParamError, "用户不存在")
		return
	}

	if user.Password != password {
		response.Error(ctx, response.CodeParamError, "密码错误")
		return
	}

	token, err := c.jwtService.CreateToken(user)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		ctx.SetCookie("token", token, 0, "", "", false, false)
		response.Success(ctx, map[string]string{"token": token})
	}
}

func (c *Controller) handleLogout(ctx *gin.Context) {
	ctx.SetCookie("token", "", 0, "", "", false, false)
	response.Success(ctx, nil)
}

func (c *Controller) handleUserInfo(ctx *gin.Context) {
	tokenString := ginutil.GetToken(ctx)

	user, err := c.jwtService.ParseToken(tokenString)
	if err != nil {
		response.Success(ctx, nil)
	} else {
		response.Success(ctx, user.Desensitization())
	}
}

func (c *Controller) handleTokenNeedRefresh(ctx *gin.Context) {
	tokenString := ginutil.GetToken(ctx)

	need, err := c.jwtService.GetNeedRefreshToken(tokenString)
	if err != nil {
		response.Success(ctx, map[string]bool{"need": false})
	} else {
		response.Success(ctx, map[string]bool{"need": need})
	}
}

func (c *Controller) handleGetUserRole(ctx *gin.Context) {
	tokenString := ginutil.GetToken(ctx)

	user, err := c.jwtService.ParseToken(tokenString)
	if err != nil {
		response.Success(ctx, []string{})
	} else {
		role, _ := c.rbacService.GetRolesForUser(user.Username)
		response.Success(ctx, role)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/user")
	api.POST("login", c.handleLogin)
	api.DELETE("logout", c.handleLogout)
	api.GET("info", c.handleUserInfo)
	api.PUT("password", c.handleUpdateUserPassword)
	api.GET("role", c.handleGetUserRole)
	api.GET("token/refresh", c.handleTokenNeedRefresh)
	api.PUT("token/refresh", c.handleRefreshToken)
}
