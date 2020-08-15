package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/config"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	webConfig "github.com/ScoreTrak/Web/pkg/config"
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

func (u *configController) ResetCompetition(c *gin.Context) {
	genericDelete(c, "ResetCompetition", u.client.ConfigClient, u.log)
}

func (u *configController) DeleteCompetition(c *gin.Context) {
	genericDelete(c, "DeleteCompetition", u.client.ConfigClient, u.log)
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

func (u *configController) GetStaticConfig(c *gin.Context) {
	genericGet(c, "Get", u.client.StaticConfigClient, u.log)
}

func (u *configController) GetStaticWebConfig(c *gin.Context) {
	c.JSON(200, webConfig.GetStaticConfig())
}
