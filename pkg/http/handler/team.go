package handler

import (
	"errors"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	sTeam "github.com/L1ghtman2k/ScoreTrak/pkg/team"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/team"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type teamController struct {
	log        logger.LogInfoFormat
	serv       team.Serv
	teamClient sTeam.Serv
}

func NewTeamController(log logger.LogInfoFormat, svc team.Serv, tc sTeam.Serv) *teamController {
	return &teamController{log, svc, tc}
}

func (u *teamController) Store(c *gin.Context) {
	us := &team.Team{}
	err := c.BindJSON(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = u.teamClient.Store(&sTeam.Team{ID: us.ID, Name: us.Name, Enabled: us.Enabled})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	savedTeam, _ := u.teamClient.GetByID(us.ID)
	us.ID = savedTeam.ID
	err = u.serv.Store(us)
	if err != nil {
		_ = u.teamClient.Delete(us.ID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (u *teamController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := u.serv.GetByID(id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}
	err = u.teamClient.Delete(t.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = u.serv.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (u *teamController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := u.teamClient.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}

	c.JSON(200, t)
}

func (u *teamController) GetAll(c *gin.Context) {
	t, err := u.teamClient.GetAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}
	c.JSON(200, t)
}

func (u *teamController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := &team.Team{}
	err = c.BindJSON(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	us.ID = id

	ts, err := u.teamClient.GetByID(us.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = u.teamClient.Update(&sTeam.Team{ID: ts.ID, Name: us.Name, Enabled: us.Enabled})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = u.serv.Update(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

//ToDo: Research about alternative approaches of handling the bad requests. Aka, what happens if after deleting a team from scoretrak, we then encounter an error deleting it on web backend
