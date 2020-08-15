package handler

import (
	"encoding/json"
	"github.com/ScoreTrak/ScoreTrak/pkg/api/client"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/competition"
	"github.com/gin-gonic/gin"
	"net/http"
)

type competitionController struct {
	log    logger.LogInfoFormat
	client *ClientStore
	serv   competition.Serv
}

func NewCompetitionController(log logger.LogInfoFormat, client *ClientStore, serv competition.Serv) *competitionController {
	return &competitionController{log, client, serv}
}

func (u *competitionController) LoadCompetition(c *gin.Context) {
	var us = &competition.Web{}
	file, err := c.FormFile("file")
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileContent, _ := file.Open()
	decoder := json.NewDecoder(fileContent)
	err = decoder.Decode(us)
	if err != nil {
		u.log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.serv.LoadCompetition(us)
	if err != nil {
		u.log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	err = u.client.CompetitionClient.LoadCompetition(us.Competition)
	if err != nil {
		u.log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

}

func (u *competitionController) FetchEntireCompetition(c *gin.Context) {
	web, err := u.serv.FetchCompetition()
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	web.Competition, err = u.client.CompetitionClient.FetchEntireCompetition()
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(200, web)
}

func (u *competitionController) FetchCoreCompetition(c *gin.Context) {
	web, err := u.serv.FetchCompetition()
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	web.Competition, err = u.client.CompetitionClient.FetchCoreCompetition()
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(200, web)
}
