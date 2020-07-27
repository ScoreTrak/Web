package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service"
	"github.com/gin-gonic/gin"
)

type serviceController struct {
	log           logger.LogInfoFormat
	serviceClient service.Serv
}

func NewServiceController(log logger.LogInfoFormat, tc service.Serv) *serviceController {
	return &serviceController{log, tc}
}

func (u *serviceController) Store(c *gin.Context) {
	us := &service.Service{}
	genericStore(c, "Store", u.serviceClient, us, u.log)

}

func (u *serviceController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.serviceClient, u.log)
}

func (u *serviceController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.serviceClient, u.log)
}

func (u *serviceController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.serviceClient, u.log)
}

func (u *serviceController) Update(c *gin.Context) {
	us := &service.Service{}
	genericUpdate(c, "Update", u.serviceClient, us, u.log)
}
