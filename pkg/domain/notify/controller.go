package notify

import (
	"fmt"
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

func (c *Controller) handleUploadNotifierPlugin(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	tmpFilePath := fmt.Sprintf("/tmp/%s", file.Filename)
	err = ctx.SaveUploadedFile(file, tmpFilePath)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, tmpFilePath)
}

func (c *Controller) handleAddNotifierPlugin(ctx *gin.Context) {
	type params struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
		File string `json:"file"`
	}
	ps := new(params)
	err := ctx.ShouldBind(ps)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	notifierPlugin := &NotifierPlugin{
		Name: ps.Name,
		Desc: ps.Desc,
	}

	err = c.service.AddNotifierPlugin(notifierPlugin, ps.File)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("notify")
	api.POST("/plugin/add", c.handleAddNotifierPlugin)
	api.POST("/plugin/upload", c.handleUploadNotifierPlugin)
}
