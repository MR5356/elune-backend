package executor

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (c *Controller) handleStartNewJob(ctx *gin.Context) {
	type data struct {
		ScriptId       uint   `json:"scriptId"`
		MachineIds     []uint `json:"machineIds"`
		MachineGroupId uint   `json:"machineGroupId"`
		Params         string `json:"params"`
	}
	body := new(data)
	err := ctx.ShouldBind(&body)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		var id string
		if body.MachineGroupId != 0 {
			logrus.Debugf("start new job with machine group id: %d", body.MachineGroupId)
			id, err = c.service.StartNewJobWithMachineGroup(body.ScriptId, body.MachineGroupId, body.Params)
		} else {
			logrus.Debugf("start new job with machine ids: %v", body.MachineIds)
			id, err = c.service.StartNewJob(body.ScriptId, body.MachineIds, body.Params)
		}
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, id)
		}
	}
}

func (c *Controller) handleStopJob(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.StopJob(id)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleGetJobLog(ctx *gin.Context) {
	id := ctx.Param("id")
	log, err := c.service.GetJobLog(id)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, log)
	}
}

func (c *Controller) handleListJob(ctx *gin.Context) {
	list, err := c.service.ListJob()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, list)
	}
}

func (c *Controller) handlePageJob(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageSizeStr := ctx.Query("pageSize")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}
	res, err := c.service.PageJob(pageNum, pageSize)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		response.Success(ctx, res)
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/execute")
	api.POST("/new", c.handleStartNewJob)
	api.DELETE("/stop/:id", c.handleStopJob)
	api.GET("/log/:id", c.handleGetJobLog)
	api.GET("/list", c.handleListJob)
	api.GET("/list/page", c.handlePageJob)
}
