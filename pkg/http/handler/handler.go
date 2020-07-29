package handler

import (
	"errors"
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	shandler "github.com/L1ghtman2k/ScoreTrak/pkg/api/handler"
	"github.com/L1ghtman2k/ScoreTrak/pkg/check"
	"github.com/L1ghtman2k/ScoreTrak/pkg/config"
	"github.com/L1ghtman2k/ScoreTrak/pkg/host"
	"github.com/L1ghtman2k/ScoreTrak/pkg/host_group"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/property"
	"github.com/L1ghtman2k/ScoreTrak/pkg/report"
	"github.com/L1ghtman2k/ScoreTrak/pkg/round"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service_group"
	"github.com/L1ghtman2k/ScoreTrak/pkg/team"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

func genericStore(c *gin.Context, m string, svc interface{}, g interface{}, log logger.LogInfoFormat) {
	err := c.BindJSON(g)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = shandler.InvokeNoRetMethod(svc, m, g)
	if err != nil {
		log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
}

func genericGet(c *gin.Context, m string, svc interface{}, log logger.LogInfoFormat) {
	sg, err := shandler.InvokeRetMethod(svc, m)
	if err != nil {
		log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, sg)
}

func genericGetByID(c *gin.Context, m string, svc interface{}, log logger.LogInfoFormat) {
	v, ok := c.Get("shortcut")
	if ok {
		c.JSON(200, v)
		return
	}
	id, err := IdResolver(c, "id")
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sg, err := shandler.InvokeRetMethod(svc, m, id)
	if err != nil {
		log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, sg)
}

func genericDelete(c *gin.Context, m string, svc interface{}, log logger.LogInfoFormat) {
	id, err := IdResolver(c, "id")
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = shandler.InvokeNoRetMethod(svc, m, id)
	if err != nil {
		log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
}

func genericUpdate(c *gin.Context, m string, svc interface{}, g interface{}, log logger.LogInfoFormat) {

	id, err := IdResolver(c, "id")
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, ok := c.Get("filtered")
	if ok {
		g = val
	} else {
		err = c.BindJSON(g)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	v := reflect.ValueOf(g).Elem()
	f := reflect.ValueOf(id)
	v.FieldByName("ID").Set(f)
	err = shandler.InvokeNoRetMethod(svc, m, g)
	if err != nil {
		log.Error(err.Error())
		if serr, ok := err.(*client.InvalidResponse); ok {
			c.AbortWithStatusJSON(serr.ResponseCode, gin.H{"error": serr.Error(), "details": serr.ResponseBody})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
}

func IdResolver(c *gin.Context, param string) (id uint64, err error) {
	idParam := c.Param(param)
	if idParam == "" {
		return 0, errors.New(fmt.Sprintf("%s parameter was not identified", param))
	}
	id, err = strconv.ParseUint(idParam, 10, 64)
	return
}

type ClientStore struct {
	ConfigClient       config.Serv
	TeamClient         team.Serv
	HostClient         host.Serv
	ServiceClient      service.Serv
	ServiceGroupClient service_group.Serv
	HostGroupClient    host_group.Serv
	PropertyClient     property.Serv
	RoundClient        round.Serv
	CheckClient        check.Serv
	ReportClient       report.Serv
}
