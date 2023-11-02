package kubernetes

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

func (c *Controller) handleGetNodes(ctx *gin.Context) {
	res, err := c.service.GetNodes()
	if err != nil {
		response.Error(ctx, response.CodeUnknownError, "")
	} else {
		response.Success(ctx, res)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/kubernetes")
	api.GET("/self", c.handleGetNodes)
}
