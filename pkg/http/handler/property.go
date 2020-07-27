package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/property"
	"github.com/gin-gonic/gin"
)

type propertyController struct {
	log            logger.LogInfoFormat
	propertyClient property.Serv
}

func NewPropertyController(log logger.LogInfoFormat, tc property.Serv) *propertyController {
	return &propertyController{log, tc}
}

func (u *propertyController) Store(c *gin.Context) {
	us := &property.Property{}
	genericStore(c, "Store", u.propertyClient, us, u.log)

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
