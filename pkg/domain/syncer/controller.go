package syncer

import "github.com/gin-gonic/gin"

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {

}
