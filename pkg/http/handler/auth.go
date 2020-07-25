package handler

import (
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/rbac"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type authController struct {
	log         logger.LogInfoFormat
	userService user.Serv
	rbac        rbac.RBAC
}

func NewAuthController(l logger.LogInfoFormat, u user.Serv, rbac rbac.RBAC) *authController {
	return &authController{l, u, rbac}
}

func (a *authController) enforce(sub string, obj string, act string) (bool, error) {
	enforcer, err := casbin.NewEnforcer(a.rbac.ConfigPath, a.rbac.Adapter)
	if err != nil {
		a.log.Error(err)
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		a.log.Error(err)
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}

	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		a.log.Error(err)
	}
	return ok, err
}

func (a *authController) JWTMiddleware() (*jwt.GinJWTMiddleware, error) {

	identityKey := "id"
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
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &user.User{
				Username: claims[identityKey].(string),
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
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			v, ok := data.(*user.User)
			if !ok {
				return false
			}

			var act string
			if strings.ToLower(c.Request.Method) == "get" {
				act = "read"
			}
			if strings.ToLower(c.Request.Method) == "post" || strings.ToLower(c.Request.Method) == "patch" || strings.ToLower(c.Request.Method) == "put" {
				act = "write"
			}
			ok, err := a.enforce(v.Username, c.HandlerName(), act) //#TODO: Revisit
			if err != nil || !ok {
				return false
			} //#TODO: Ship application in 2 ways (Static rules, and Dynamic Rules)
			return true
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
