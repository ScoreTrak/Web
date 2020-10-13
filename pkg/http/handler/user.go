package handler

import (
	"errors"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type userController struct {
	log  logger.LogInfoFormat
	serv user.Serv
}

func NewUserController(log logger.LogInfoFormat, svc user.Serv) *userController {
	return &userController{log, svc}
}

func (u *userController) Store(c *gin.Context) {
	var us []*user.User
	err := c.BindJSON(&us)
	for i := range us {
		if us[i].Password == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		password := []byte(us[i].Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		us[i].PasswordHash = string(hashedPassword)
	}

	err = u.serv.Store(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (u *userController) GetByID(c *gin.Context) {
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
	c.JSON(200, t)
}

func (u *userController) GetAll(c *gin.Context) {
	t, err := u.serv.GetAll()
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

func (u *userController) Delete(c *gin.Context) {
	idParam, err := UuidResolver(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.serv.Delete(idParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
		}
		return
	}
}

func (u *userController) Update(c *gin.Context) {
	us := &user.User{}
	id, _ := UuidResolver(c, "id")
	role := roleResolver(c)
	err := c.BindJSON(us)
	if err != nil {
		u.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if role == "blue" {
		if id != us.ID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You can not edit this object"})
			return
		}
		us = &user.User{Username: us.Username, Password: us.Password}
	}
	us.ID = id
	if us.Password != "" {
		password := []byte(us.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		us.PasswordHash = string(hashedPassword)
	}
	err = u.serv.Update(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

//TODO: ON password update invalidate JWT(HAshing via password?). https://stackoverflow.com/questions/28759590/best-practices-to-invalidate-jwt-while-changing-passwords-and-logout-in-node-js
