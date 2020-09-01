package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type checkController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewCheckController(log logger.LogInfoFormat, client *ClientStore) *checkController {
	return &checkController{log, client}
}

func (u *checkController) GetByRoundServiceID(c *gin.Context) {
	rid, _ := UintResolver(c, "RoundID")
	sid, _ := UuidResolver(c, "ServiceID")

	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, prop, err := teamIDFromCheck(u.client, rid, sid)
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

	sg, err := u.client.CheckClient.GetByRoundServiceID(rid, sid)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
	c.JSON(200, sg)
}

func (u *checkController) GetAllByRoundID(c *gin.Context) {
	genericGetByID(c, "GetAllByRoundID", u.client.CheckClient, u.log)
}

func (u *checkController) GetAllByServiceID(c *gin.Context) {
	sid, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	TeamID := teamIDResolver(c)
	if role == "blue" {
		tID, _, err := teamIDFromService(u.client, sid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tID != TeamID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not access this object"})
			return
		}
	}
	sg, err := u.client.CheckClient.GetAllByServiceID(sid)
	if err != nil {
		ClientErrorHandler(c, u.log, err)
		return
	}
	c.JSON(200, sg)
}
