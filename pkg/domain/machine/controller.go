package machine

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

func (c *Controller) handleAddMachine(ctx *gin.Context) {
	machine := new(Machine)
	err := ctx.ShouldBind(machine)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddMachine(machine)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleListMachine(ctx *gin.Context) {
	machineList, err := c.service.ListMachine()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, machineList)
	}
}

func (c *Controller) handleDeleteMachine(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteMachine(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleUpdateMachine(ctx *gin.Context) {
	machine := new(Machine)
	err := ctx.ShouldBind(machine)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.UpdateMachine(machine)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleAddGroup(ctx *gin.Context) {
	group := new(Group)
	err := ctx.ShouldBind(group)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.AddGroup(group)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) handleListGroup(ctx *gin.Context) {
	groupList, err := c.service.ListGroup()
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, groupList)
	}
}

func (c *Controller) handleDeleteGroup(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	err = c.service.DeleteGroup(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
	} else {
		response.Success(ctx, nil)
	}
}

func (c *Controller) handleUpdateGroup(ctx *gin.Context) {
	group := new(Group)
	err := ctx.ShouldBind(group)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	} else {
		err := c.service.UpdateGroup(group)
		if err != nil {
			response.Error(ctx, response.CodeParamError, err.Error())
		} else {
			response.Success(ctx, nil)
		}
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/machine")
	api.POST("/add", c.handleAddMachine)
	api.GET("/list", c.handleListMachine)
	api.DELETE("/delete/:id", c.handleDeleteMachine)
	api.PUT("/update", c.handleUpdateMachine)

	api.POST("/group/add", c.handleAddGroup)
	api.GET("/group/list", c.handleListGroup)
	api.DELETE("/group/delete/:id", c.handleDeleteGroup)
	api.PUT("/group/update", c.handleUpdateGroup)

	api.GET("/terminal/:id", c.handleTerminal)
}
