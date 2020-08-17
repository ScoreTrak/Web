package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/host"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type hostController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewHostController(log logger.LogInfoFormat, client *ClientStore) *hostController {
	return &hostController{log, client}
}

func (u *hostController) Store(c *gin.Context) {
	var us []*host.Host
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.client.HostClient.Store(us)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}

}

func (u *hostController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.client.HostClient, u.log)
}

func (u *hostController) GetByID(c *gin.Context) {
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, prop, err := teamIDFromHost(u.client, id)
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
	genericGetByID(c, "GetByID", u.client.HostClient, u.log)
}

func (u *hostController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.HostClient, u.log)
}

func (u *hostController) Update(c *gin.Context) {
	us := &host.Host{}
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	err := c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if role == "blue" {
		tID, prop, err := teamIDFromHost(u.client, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tID != TeamID || prop.EditHost != nil && *prop.EditHost == true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not edit this object"})
			return
		}
		us = &host.Host{Address: us.Address}
	}
	us.ID = id
	err = u.client.HostClient.Update(us)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
}
