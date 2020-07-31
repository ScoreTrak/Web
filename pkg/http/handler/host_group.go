package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/host_group"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
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
		u.log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
