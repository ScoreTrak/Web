package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
)

type roundController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewRoundController(log logger.LogInfoFormat, client *ClientStore) *roundController {
	return &roundController{log, client}
}

func (u *roundController) GetLastNonElapsingRound(c *gin.Context) {
	genericGet(c, "GetLastNonElapsingRound", u.client.RoundClient, u.log)

}

func (u *roundController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.client.RoundClient, u.log)

}

func (u *roundController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.client.RoundClient, u.log)

}
