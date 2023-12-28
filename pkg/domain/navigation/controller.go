package navigation

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) handleListNavigation(ctx *gin.Context) {
	navigationList, err := c.service.ListNavigation()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, navigationList)
	}
}

func (c *Controller) handleAddNavigation(ctx *gin.Context) {
	navigation := new(Navigation)
	err := ctx.ShouldBind(navigation)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddNavigation(navigation)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleUpdateNavigation(ctx *gin.Context) {
	navigation := new(Navigation)
	err := ctx.ShouldBind(navigation)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.UpdateNavigation(navigation)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleDeleteNavigation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteNavigation(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/navigation")
	api.GET("/list", c.handleListNavigation)
	api.POST("/add", c.handleAddNavigation)
	api.PUT("/update", c.handleUpdateNavigation)
	api.DELETE("/delete/:id", c.handleDeleteNavigation)
}
