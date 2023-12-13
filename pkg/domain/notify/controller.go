package notify

import (
	"fmt"
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
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

	tmpFilePath := fmt.Sprintf("/tmp/elune/%s-%s", time.Now().Format("20060102150405"), file.Filename)
	err = ctx.SaveUploadedFile(file, tmpFilePath)
	defer os.RemoveAll(tmpFilePath)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	err = c.service.UploadNotifierPlugin(tmpFilePath)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

//func (c *Controller) handleAddNotifierPlugin(ctx *gin.Context) {
//	type params struct {
//		Name string `json:"name"`
//		Desc string `json:"desc"`
//		File string `json:"file"`
//	}
//	ps := new(params)
//	err := ctx.ShouldBind(ps)
//	if err != nil {
//		response.Error(ctx, response.CodeParamError, err.Error())
//		return
//	}
//
//	notifierPlugin := &NotifierPlugin{
//		Name: ps.Name,
//	}
//
//	err = c.service.AddNotifierPlugin(notifierPlugin, ps.File)
//	if err != nil {
//		response.Error(ctx, response.CodeParamError, err.Error())
//		return
//	}
//	response.Success(ctx, nil)
//}

func (c *Controller) handleListNotifierPlugins(ctx *gin.Context) {
	plugins, err := c.service.ListNotifierPlugins()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, plugins)
}

func (c *Controller) handleUninstallNotifierPlugin(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.RemoveNotifierPlugin(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleAddNotifierChannel(ctx *gin.Context) {
	notifierChannel := new(NotifierChannel)
	err := ctx.ShouldBind(notifierChannel)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.AddNotifierChannel(notifierChannel)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) handleListNotifierChannels(ctx *gin.Context) {
	channels, err := c.service.ListNotifierChannels()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, channels)
}

func (c *Controller) handleRemoveNotifierChannel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.RemoveNotifierChannel(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) handleSendTestMessage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.SendTestMessage(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) handleListMessageTemplates(ctx *gin.Context) {
	res, err := c.service.ListMessageTemplates()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, res)
}

func (c *Controller) handleAddMessageTemplate(ctx *gin.Context) {
	msgTemplate := new(MessageTemplate)
	err := ctx.ShouldBind(msgTemplate)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.AddMessageTemplate(msgTemplate)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) handleRemoveMessageTemplate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.RemoveMessageTemplate(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) handleUpdateMessageTemplate(ctx *gin.Context) {
	msgTemplate := new(MessageTemplate)
	err := ctx.ShouldBind(msgTemplate)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.UpdateMessageTemplate(msgTemplate)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("notify")
	//api.POST("/plugin/add", c.handleAddNotifierPlugin)
	api.POST("/plugin/upload", c.handleUploadNotifierPlugin)
	api.GET("/plugin/list", c.handleListNotifierPlugins)
	api.DELETE("/plugin/uninstall/:id", c.handleUninstallNotifierPlugin)
	api.POST("/channel/add", c.handleAddNotifierChannel)
	api.GET("/channel/list", c.handleListNotifierChannels)
	api.DELETE("/channel/remove/:id", c.handleRemoveNotifierChannel)
	api.POST("/channel/send/:id", c.handleSendTestMessage)
	api.GET("/template/list", c.handleListMessageTemplates)
	api.POST("/template/add", c.handleAddMessageTemplate)
	api.DELETE("/template/remove/:id", c.handleRemoveMessageTemplate)
	api.PUT("/template/update", c.handleUpdateMessageTemplate)
}
