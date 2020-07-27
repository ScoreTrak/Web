package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/check"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type checkController struct {
	log         logger.LogInfoFormat
	checkClient check.Serv
}

func NewCheckController(log logger.LogInfoFormat, tc check.Serv) *checkController {
	return &checkController{log, tc}
}

func (u *checkController) GetByRoundServiceID(c *gin.Context) {
	srid := c.Param("RoundID")
	rid, err := strconv.ParseUint(srid, 10, 64)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ssid := c.Param("ServiceID")
	sid, err := strconv.ParseUint(ssid, 10, 64)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sg, err := u.checkClient.GetByRoundServiceID(rid, sid)
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
	genericGetByID(c, "GetAllByRoundID", u.checkClient, u.log)
}
