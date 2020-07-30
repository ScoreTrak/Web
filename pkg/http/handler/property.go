package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/property"
	"github.com/gin-gonic/gin"
	"net/http"
)

type propertyController struct {
	log            logger.LogInfoFormat
	propertyClient property.Serv
}

func NewPropertyController(log logger.LogInfoFormat, tc property.Serv) *propertyController {
	return &propertyController{log, tc}
}

func (u *propertyController) Store(c *gin.Context) {
	var us []*property.Property
	err := c.BindJSON(&us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = u.propertyClient.Store(us)
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

func (u *propertyController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.propertyClient, u.log)
}

func (u *propertyController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.propertyClient, u.log)
}

func (u *propertyController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.propertyClient, u.log)
}

func (u *propertyController) Update(c *gin.Context) {
	us := &property.Property{}
	genericUpdate(c, "Update", u.propertyClient, us, u.log)
}
