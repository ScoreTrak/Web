package handler

import (
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authController struct {
	log         logger.LogInfoFormat
	userService user.Serv
}

func NewAuthController(l logger.LogInfoFormat, u user.Serv) *authController {
	return &authController{l, u}
}

func (a *authController) JWTMiddleware() (*jwt.GinJWTMiddleware, error) {

	identityKey := "username"
	type login struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "scoretrak",
		Key:         []byte(config.GetStaticConfig().Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				fmt.Println(v)
				return jwt.MapClaims{
					identityKey: v.Username,
					"team_id":   v.TeamID,
					"role":      v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			tID := claims["team_id"].(float64)
			return &user.User{
				Username: claims[identityKey].(string),
				TeamID:   uint64(tID),
				Role:     claims["role"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			usr, err := a.userService.GetByUsername(loginVals.Username)
			if err != nil {
				return nil, err
			}
			err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(loginVals.Password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user.User{
				Username: loginVals.Username,
				TeamID:   usr.TeamID,
				Role:     usr.Role,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*user.User); ok && v.Role == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
