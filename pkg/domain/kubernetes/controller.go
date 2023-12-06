package kubernetes

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

func (c *Controller) handleAddKubernetes(ctx *gin.Context) {
	k8s := new(Kubernetes)
	err := ctx.ShouldBind(k8s)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddKubernetes(k8s)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleListKubernetes(ctx *gin.Context) {
	k8sList, err := c.service.ListKubernetes()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, k8sList)
	}
}

func (c *Controller) handleUpdateKubernetes(ctx *gin.Context) {
	k8s := new(Kubernetes)
	err := ctx.ShouldBind(k8s)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.UpdateKubernetes(k8s)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleDeleteKubernetes(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteKubernetes(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/kubernetes")
	api.POST("add", c.handleAddKubernetes)
	api.GET("list", c.handleListKubernetes)
	api.PUT("update", c.handleUpdateKubernetes)
	api.DELETE("delete/:id", c.handleDeleteKubernetes)
}
