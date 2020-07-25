package handler

import (
	"errors"
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
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
	err = invokeNoRetMethod(svc, m, g)
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
	sg, err := invokeRetMethod(svc, m)
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
	sg, err := invokeRetMethod(svc, m, id)
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
	err = invokeNoRetMethod(svc, m, id)
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
	err = invokeNoRetMethod(svc, m, g)
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

func invokeRetMethod(i interface{}, methodName string, args ...interface{}) (interface{}, error) {
	finalMethod := preInvoke(i, methodName)
	if finalMethod.IsValid() {
		inputs := make([]reflect.Value, len(args))
		for i, _ := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		r := finalMethod.Call(inputs)

		if err, ok := r[1].Interface().(error); ok {
			return nil, err
		}
		return r[0].Interface(), nil
	}
	return nil, errors.New(fmt.Sprintf("The method name %s does not exist in %s", methodName, reflect.TypeOf(i).Name()))
}

func invokeNoRetMethod(i interface{}, methodName string, args ...interface{}) error {
	finalMethod := preInvoke(i, methodName)
	if finalMethod.IsValid() {
		inputs := make([]reflect.Value, len(args))
		for i, _ := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		r := finalMethod.Call(inputs)

		if err, ok := r[0].Interface().(error); ok {
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("The method name %s does not exist in %s", methodName, reflect.TypeOf(i).Name()))
}

func preInvoke(i interface{}, methodName string) reflect.Value {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	return finalMethod
}

func idResolver(c *gin.Context) (id uint64, err error) {
	idParam := c.Param("id")
	id, err = strconv.ParseUint(idParam, 10, 64)
	return
}
