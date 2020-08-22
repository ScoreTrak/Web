package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/property"
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
		ClientErrorHandler(c, u.log, err)
		return
	}

}

func (u *propertyController) Delete(c *gin.Context) {
	sID, err := UuidResolver(c, "ServiceID")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := ParamResolver(c, "Key")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.client.PropertyClient.Delete(sID, key)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
}

func (u *propertyController) GetByServiceIDKey(c *gin.Context) {
	sID, err := UuidResolver(c, "ServiceID")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := ParamResolver(c, "Key")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, prop, err := teamIDFromProperty(u.client, sID, key)
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

func (u *propertyController) GetAllByServiceID(c *gin.Context) {
	id, _ := UuidResolver(c, "ServiceID")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, _, err := teamIDFromService(u.client, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		props, err := u.client.PropertyClient.GetAllByServiceID(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var filteredProps []*property.Property
		for i := range props {
			if tID != TeamID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not access this object"})
				return
			} else if props[i].Status == property.Hide {
				continue
			}
			filteredProps = append(filteredProps, props[i])
		}

		c.JSON(200, filteredProps)
		return
	}
	sg, err := u.client.PropertyClient.GetAllByServiceID(id)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
	c.JSON(200, sg)
}

func (u *propertyController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.PropertyClient, u.log)
}

//ToDo: Ensure that data handled by Users is properly handeled (XSS prevention on this stage)

func (u *propertyController) Update(c *gin.Context) {
	us := &property.Property{}
	sID, err := UuidResolver(c, "ServiceID")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := ParamResolver(c, "Key")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	err = c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if role == "blue" {
		tID, prop, err := teamIDFromProperty(u.client, sID, key)
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
	us.ServiceID = sID
	us.Key = key
	err = u.client.PropertyClient.Update(us)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
}
