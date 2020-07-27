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
	us := &round.Round{}
	genericStore(c, "GetLastNonElapsingRound", u.roundClient, us, u.log)

}

func (u *roundController) GetAll(c *gin.Context) {
	us := &round.Round{}
	genericStore(c, "GetAll", u.roundClient, us, u.log)

}

func (u *roundController) GetByID(c *gin.Context) {
	us := &round.Round{}
	genericStore(c, "GetByID", u.roundClient, us, u.log)

}
