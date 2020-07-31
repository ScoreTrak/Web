package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/property"
	"github.com/gin-gonic/gin"
	"net/http"
)

type propertyController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewPropertyController(log logger.LogInfoFormat, client *ClientStore) *propertyController {
	return &propertyController{log, client}
}

func (u *propertyController) Store(c *gin.Context) {
	var us []*property.Property
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.client.PropertyClient.Store(us)
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

func (u *propertyController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.client.PropertyClient, u.log)
}

func (u *propertyController) GetByID(c *gin.Context) {
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, prop, err := teamIDFromProperty(u.client, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if tID != TeamID || prop.Status == property.Hide {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not access this object"})
			return
		}
		c.Set("shortcut", prop)
	}
	genericGetByID(c, "GetByID", u.client.PropertyClient, u.log)
}

func (u *propertyController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.PropertyClient, u.log)
}

func (u *propertyController) Update(c *gin.Context) {
	us := &property.Property{}
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	err := c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if role == "blue" {
		tID, prop, err := teamIDFromProperty(u.client, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tID != TeamID || prop.Status != property.Edit {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not edit this object"})
			return
		}
		us = &property.Property{Value: us.Value}
	}
	us.ID = id
	err = u.client.PropertyClient.Update(us)
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
