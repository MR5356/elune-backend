package site

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) setKeyController(ctx *gin.Context) {
	siteConfig := new(SiteConfig)
	err := ctx.ShouldBind(siteConfig)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	}

	if err := c.service.SetKey(siteConfig.Key, siteConfig.Value); err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) getKeyController(ctx *gin.Context) {
	key := ctx.Param("key")
	value, err := c.service.GetKey(key)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, value)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/site")
	api.PUT("/config", c.setKeyController)
	api.GET("/config/:key", c.getKeyController)
}
