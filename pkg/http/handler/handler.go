package handler

import (
	"errors"
	"fmt"
	"github.com/ScoreTrak/ScoreTrak/pkg/api/client"
	shandler "github.com/ScoreTrak/ScoreTrak/pkg/api/handler"
	"github.com/ScoreTrak/ScoreTrak/pkg/check"
	"github.com/ScoreTrak/ScoreTrak/pkg/competition"
	"github.com/ScoreTrak/ScoreTrak/pkg/config"
	"github.com/ScoreTrak/ScoreTrak/pkg/host"
	"github.com/ScoreTrak/ScoreTrak/pkg/host_group"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/property"
	"github.com/ScoreTrak/ScoreTrak/pkg/report"
	"github.com/ScoreTrak/ScoreTrak/pkg/round"
	"github.com/ScoreTrak/ScoreTrak/pkg/service"
	"github.com/ScoreTrak/ScoreTrak/pkg/service_group"
	"github.com/ScoreTrak/ScoreTrak/pkg/team"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
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
	id, err := UuidResolver(c, "id")
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
	id, err := UuidResolver(c, "id")
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
	id, err := UuidResolver(c, "id")
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.BindJSON(g)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
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

func UintResolver(c *gin.Context, param string) (id uint, err error) {
	idParam := c.Param(param)
	if idParam == "" {
		return 0, errors.New(fmt.Sprintf("%s parameter was not identified", param))
	}
	id32, err := strconv.ParseUint(idParam, 10, 32)
	id = uint(id32)
	return
}

func UuidResolver(c *gin.Context, param string) (uuid.UUID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return uuid.Nil, errors.New(fmt.Sprintf("%s parameter was not identified", param))
	}
	return uuid.FromString(idParam)
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
	PolicyClient       *policy.Client
	CompetitionClient  competition.Serv
}

func teamIDFromProperty(c *ClientStore, propertyID uuid.UUID) (teamID uuid.UUID, property *property.Property, err error) {
	property, err = c.PropertyClient.GetByID(propertyID)
	if err != nil || property == nil {
		return
	}
	teamID, _, err = teamIDFromService(c, property.ServiceID)
	return
}

func teamIDFromCheck(c *ClientStore, roundID uint, serviceID uuid.UUID) (teamID uuid.UUID, check *check.Check, err error) {
	check, err = c.CheckClient.GetByRoundServiceID(roundID, serviceID)
	if err != nil || check == nil {
		return
	}
	teamID, _, err = teamIDFromService(c, check.ServiceID)
	return
}

func teamIDFromService(c *ClientStore, serviceID uuid.UUID) (teamID uuid.UUID, service *service.Service, err error) {
	service, err = c.ServiceClient.GetByID(serviceID)
	if err != nil || service == nil {
		return
	}
	teamID, _, err = teamIDFromHost(c, service.HostID)
	return
}

func teamIDFromHost(c *ClientStore, hostID uuid.UUID) (teamID uuid.UUID, host *host.Host, err error) {
	host, err = c.HostClient.GetByID(hostID)
	if err != nil || host == nil {
		return
	}
	return host.TeamID, host, err
}

func roleResolver(c *gin.Context) (role string) {
	if val, ok := c.Get("role"); ok && val != nil {
		role, _ = val.(string)
	}
	return
}

func teamIDResolver(c *gin.Context) (teamID uuid.UUID) {
	if val, ok := c.Get("team_id"); ok && val != nil {
		teamID, _ = val.(uuid.UUID)
	}
	return
}
