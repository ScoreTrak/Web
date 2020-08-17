package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type serviceController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewServiceController(log logger.LogInfoFormat, client *ClientStore) *serviceController {
	return &serviceController{log, client}
}

func (u *serviceController) Store(c *gin.Context) {
	var us []*service.Service
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.client.ServiceClient.Store(us)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}

}

func (u *serviceController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.client.ServiceClient, u.log)
}

func (u *serviceController) GetByID(c *gin.Context) {
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, prop, err := teamIDFromService(u.client, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if tID != TeamID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not access this object"})
			return
		}
		c.Set("shortcut", prop)

	}
	genericGetByID(c, "GetByID", u.client.ServiceClient, u.log)
}

func (u *serviceController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.ServiceClient, u.log)
}

func (u *serviceController) Update(c *gin.Context) {
	us := &service.Service{}
	genericUpdate(c, "Update", u.client.ServiceClient, us, u.log)
}
