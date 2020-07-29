package handler

import (
	"errors"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type userController struct {
	log  logger.LogInfoFormat
	serv user.Serv
}

func NewUserController(log logger.LogInfoFormat, svc user.Serv) *userController {
	return &userController{log, svc}
}

func (u *userController) Store(c *gin.Context) {
	us := &user.User{}
	err := c.BindJSON(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if us.Password != us.PasswordConfirmation {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password and password confirmation did not match"})
		return
	}
	password := []byte(us.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	us.PasswordHash = string(hashedPassword)
	err = u.serv.Store(us)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (u *userController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := u.serv.GetByID(id)
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
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.serv.Delete(id)
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
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	us := &user.User{}
	val, ok := c.Get("filtered")
	if ok {
		us = val.(*user.User)
	} else {
		err = c.BindJSON(us)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	us.ID = id

	if us.Password != "" {
		if us.Password != us.PasswordConfirmation {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password and password confirmation did not match"})
			return
		}
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
