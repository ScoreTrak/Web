package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/config"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type configController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewConfigController(log logger.LogInfoFormat, client *ClientStore) *configController {
	return &configController{log, client}
}
func (u *configController) Get(c *gin.Context) {
	genericGet(c, "Get", u.client.ConfigClient, u.log)
}

func (u *configController) Update(c *gin.Context) {
	us := &config.DynamicConfig{}
	err := c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.client.ConfigClient.Update(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
}
