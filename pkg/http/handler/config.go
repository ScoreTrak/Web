package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/config"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
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
	genericUpdate(c, "Update", u.configClient, us, u.log)
}
