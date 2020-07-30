package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type serviceController struct {
	log           logger.LogInfoFormat
	serviceClient service.Serv
}

func NewServiceController(log logger.LogInfoFormat, tc service.Serv) *serviceController {
	return &serviceController{log, tc}
}

func (u *serviceController) Store(c *gin.Context) {
	var us []*service.Service
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.serviceClient.Store(us)
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
