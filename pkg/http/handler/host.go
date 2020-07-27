package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/host"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/gin-gonic/gin"
)

type hostController struct {
	log        logger.LogInfoFormat
	hostClient host.Serv
}

func NewHostController(log logger.LogInfoFormat, tc host.Serv) *hostController {
	return &hostController{log, tc}
}

func (u *hostController) Store(c *gin.Context) {
	us := &host.Host{}
	genericStore(c, "Store", u.hostClient, us, u.log)

}

func (u *hostController) Delete(c *gin.Context) {
	genericDelete(c, "Delete", u.hostClient, u.log)
}

func (u *hostController) GetByID(c *gin.Context) {
	genericGetByID(c, "GetByID", u.hostClient, u.log)
}

func (u *hostController) GetAll(c *gin.Context) {
	genericGet(c, "GetAll", u.hostClient, u.log)
}

func (u *hostController) Update(c *gin.Context) {
	us := &host.Host{}
	genericUpdate(c, "Update", u.hostClient, us, u.log)
}
