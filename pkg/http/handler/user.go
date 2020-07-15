package handler

import (
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type userController struct {
	log  log.Logger
	serv user.Serv
}

func (u *userController) Store(c *gin.Context) {
	us := &user.User{}
	err := c.BindJSON(u)
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

//TODO: ON password update invalidate JWT(HAshing via password?). https://stackoverflow.com/questions/28759590/best-practices-to-invalidate-jwt-while-changing-passwords-and-logout-in-node-js
