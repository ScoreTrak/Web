package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service_group"
	"github.com/gin-gonic/gin"
)

type serviceGroupController struct {
	log                logger.LogInfoFormat
	serviceGroupClient service_group.Serv
}

func NewServiceGroupController(log logger.LogInfoFormat, tc service_group.Serv) *serviceGroupController {
	return &serviceGroupController{log, tc}
}

func (u *serviceGroupController) Store(c *gin.Context) {
	us := &service_group.ServiceGroup{}
	genericStore(c, "Store", u.serviceGroupClient, us, u.log)

}

func (u *serviceGroupController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.serviceGroupClient, u.log)
}

func (u *serviceGroupController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.serviceGroupClient, u.log)
}

func (u *serviceGroupController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.serviceGroupClient, u.log)
}

func (u *serviceGroupController) Update(c *gin.Context) {
	us := &service_group.ServiceGroup{}
	genericUpdate(c, "Update", u.serviceGroupClient, us, u.log)
}
