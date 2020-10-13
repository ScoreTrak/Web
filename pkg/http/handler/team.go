package handler

import (
	"errors"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	sTeam "github.com/ScoreTrak/ScoreTrak/pkg/team"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type teamController struct {
	log    logger.LogInfoFormat
	serv   team.Serv
	client *ClientStore
}

func NewTeamController(log logger.LogInfoFormat, svc team.Serv, client *ClientStore) *teamController {
	return &teamController{log, svc, client}
}

func (u *teamController) Store(c *gin.Context) {
	var us []*team.Team
	err := c.BindJSON(&us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		u.log.Error(err)
		return
	}

	err = u.serv.Store(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
		return
	}

	var usCopy []*sTeam.Team
	for i := range us {
		usCopy = append(usCopy, &sTeam.Team{ID: us[i].ID, Name: us[i].Name, Enabled: us[i].Enabled})
	}
	err = u.client.TeamClient.Store(usCopy)
	if err != nil {
		for i := range us {
			_ = u.serv.Delete(us[i].ID)
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
		return
	}

}

func (u *teamController) Delete(c *gin.Context) {
	idParam, err := UuidResolver(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := u.serv.GetByID(idParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}
	err = u.client.TeamClient.Delete(t.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = u.serv.Delete(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (u *teamController) GetByID(c *gin.Context) {
	idParam, err := UuidResolver(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tScoreTrak, err := u.client.TeamClient.GetByID(idParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}

	tWeb, err := u.serv.GetByID(idParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}
	tWeb.Enabled = tScoreTrak.Enabled

	c.JSON(200, tWeb)
}

func (u *teamController) GetAll(c *gin.Context) {
	tScoreTrak, err := u.client.TeamClient.GetAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}

	tWeb, err := u.serv.GetAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}

	var response []*team.Team

	for i := range tScoreTrak {
		for j := range tWeb {
			if tScoreTrak[i].ID == tWeb[j].ID {
				tWeb[j].Enabled = tScoreTrak[i].Enabled
				response = append(response, tWeb[j])
			}
		}
	}

	c.JSON(200, response)
}

func (u *teamController) Update(c *gin.Context) {
	idParam, err := UuidResolver(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := &team.Team{}
	err = c.BindJSON(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	us.ID = idParam

	ts, err := u.client.TeamClient.GetByID(us.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = u.client.TeamClient.Update(&sTeam.Team{ID: ts.ID, Name: us.Name, Enabled: us.Enabled})
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
