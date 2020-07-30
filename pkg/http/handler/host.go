package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/host"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type hostController struct {
	log        logger.LogInfoFormat
	hostClient host.Serv
}

func NewHostController(log logger.LogInfoFormat, tc host.Serv) *hostController {
	return &hostController{log, tc}
}

func (u *hostController) Store(c *gin.Context) {
	var us []*host.Host
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.hostClient.Store(us)
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

func (u *hostController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.hostClient, u.log)
}

func (u *hostController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.hostClient, u.log)
}

func (u *hostController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.hostClient, u.log)
}

func (u *hostController) Update(c *gin.Context) {
	us := &host.Host{}
	genericUpdate(c, "Update", u.hostClient, us, u.log)
}
