package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/config"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type configController struct {
	log          logger.LogInfoFormat
	configClient config.Serv
}

func NewConfigController(log logger.LogInfoFormat, tc config.Serv) *configController {
	return &configController{log, tc}
}
func (u *configController) Get(c *gin.Context) {
	genericGet(c, "Get", u.configClient, u.log)
}

func (u *configController) Update(c *gin.Context) {
	us := &config.DynamicConfig{}
	err := c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.configClient.Update(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
}
