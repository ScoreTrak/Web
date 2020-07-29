package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/round"
	"github.com/gin-gonic/gin"
)

type roundController struct {
	log         logger.LogInfoFormat
	roundClient round.Serv
}

func NewRoundController(log logger.LogInfoFormat, tc round.Serv) *roundController {
	return &roundController{log, tc}
}

func (u *roundController) GetLastNonElapsingRound(c *gin.Context) {
	genericGet(c, "GetLastNonElapsingRound", u.roundClient, u.log)

}

func (u *roundController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.roundClient, u.log)

}

func (u *roundController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.roundClient, u.log)

}
