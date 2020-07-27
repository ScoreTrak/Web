package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/host_group"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
)

type hostGroupController struct {
	log             logger.LogInfoFormat
	hostGroupClient host_group.Serv
}

func NewHostGroupController(log logger.LogInfoFormat, tc host_group.Serv) *hostGroupController {
	return &hostGroupController{log, tc}
}

func (u *hostGroupController) Store(c *gin.Context) {
	us := &host_group.HostGroup{}
	genericStore(c, "Store", u.hostGroupClient, us, u.log)

}

func (u *hostGroupController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.hostGroupClient, u.log)
}

func (u *hostGroupController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.hostGroupClient, u.log)
}

func (u *hostGroupController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.hostGroupClient, u.log)
}

func (u *hostGroupController) Update(c *gin.Context) {
	us := &host_group.HostGroup{}
	genericUpdate(c, "Update", u.hostGroupClient, us, u.log)
}
