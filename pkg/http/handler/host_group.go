package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/host_group"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type hostGroupController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewHostGroupController(log logger.LogInfoFormat, client *ClientStore) *hostGroupController {
	return &hostGroupController{log, client}
}

func (u *hostGroupController) Store(c *gin.Context) {
	var us []*host_group.HostGroup
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.client.HostGroupClient.Store(us)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}

}

func (u *hostGroupController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.client.HostGroupClient, u.log)
}

func (u *hostGroupController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.client.HostGroupClient, u.log)
}

func (u *hostGroupController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.HostGroupClient, u.log)
}

func (u *hostGroupController) Update(c *gin.Context) {
	us := &host_group.HostGroup{}
	genericUpdate(c, "Update", u.client.HostGroupClient, us, u.log)
}
