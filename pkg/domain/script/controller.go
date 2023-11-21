package script

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

func (c *Controller) handleAddScript(ctx *gin.Context) {
	script := new(Script)
	err := ctx.ShouldBind(script)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddScript(script)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleListScript(ctx *gin.Context) {
	scriptList, err := c.service.ListScript()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, scriptList)
	}
}

func (c *Controller) handleDeleteScript(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteScript(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleUpdateScript(ctx *gin.Context) {
	script := new(Script)
	err := ctx.ShouldBind(script)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.UpdateScript(script)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/script")
	api.POST("/add", c.handleAddScript)
	api.GET("/list", c.handleListScript)
	api.DELETE("/delete/:id", c.handleDeleteScript)
	api.PUT("/update", c.handleUpdateScript)
}
