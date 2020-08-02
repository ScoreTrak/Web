package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/service_group"
	"github.com/gin-gonic/gin"
)

type serviceGroupController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewServiceGroupController(log logger.LogInfoFormat, client *ClientStore) *serviceGroupController {
	return &serviceGroupController{log, client}
}

func (u *serviceGroupController) Store(c *gin.Context) {
	us := &service_group.ServiceGroup{}
	genericStore(c, "Store", u.client.ServiceGroupClient, us, u.log)

}

func (u *serviceGroupController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.client.ServiceGroupClient, u.log)
}

func (u *serviceGroupController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.client.ServiceGroupClient, u.log)
}

func (u *serviceGroupController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.ServiceGroupClient, u.log)
}

func (u *serviceGroupController) Update(c *gin.Context) {
	us := &service_group.ServiceGroup{}
	genericUpdate(c, "Update", u.client.ServiceGroupClient, us, u.log)
}
