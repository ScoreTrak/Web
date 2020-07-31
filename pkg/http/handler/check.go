package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
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
		u.log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, sg)
}

func (u *checkController) GetAllByRoundID(c *gin.Context) {
	genericGetByID(c, "GetAllByRoundID", u.client.CheckClient, u.log)
}
