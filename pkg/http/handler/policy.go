package handler

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/gin-gonic/gin"
	"net/http"
)

type policyController struct {
	log  logger.LogInfoFormat
	serv policy.Serv
}

func NewPolicyController(log logger.LogInfoFormat, serv policy.Serv) *policyController {
	return &policyController{log, serv}
}

func (a *policyController) GetPolicy(c *gin.Context) { //Todo: Expose policy for everyone. this can help with a better design of front end
	p, err := a.serv.Get()
	if err != nil {
		a.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(200, p)
}

func (a *policyController) UpdatePolicy(c *gin.Context) {
	p := &policy.Policy{ID: 1}
	err := c.BindJSON(p)
	if err != nil {
		a.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = a.serv.Update(p)
	if err != nil {
		a.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
}
