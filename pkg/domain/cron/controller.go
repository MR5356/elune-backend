package cron

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

func (c *Controller) handleListCron(ctx *gin.Context) {
	cs, err := c.service.ListCron()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, cs)
	}
}

func (c *Controller) handleSetEnableCron(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.SetEnableCron(uint(id), true)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleSetDisableCron(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.SetEnableCron(uint(id), false)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleAddCron(ctx *gin.Context) {
	cc := new(Cron)
	err := ctx.ShouldBind(cc)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddCron(cc)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleDeleteCron(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteCron(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/cron")
	api.GET("/list", c.handleListCron)
	api.PUT("/enable/:id", c.handleSetEnableCron)
	api.PUT("/disable/:id", c.handleSetDisableCron)
	api.POST("/add", c.handleAddCron)
	api.DELETE("/delete/:id", c.handleDeleteCron)
}
