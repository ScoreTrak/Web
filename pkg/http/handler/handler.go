package handler

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	shandler "github.com/L1ghtman2k/ScoreTrak/pkg/api/handler"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
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
	id, err := idResolver(c)
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
	id, err := idResolver(c)
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
	id, err := idResolver(c)
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

func idResolver(c *gin.Context) (id uint64, err error) {
	idParam := c.Param("id")
	id, err = strconv.ParseUint(idParam, 10, 64)
	return
}
