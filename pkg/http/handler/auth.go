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
			var username string
			if v, ok := data.(*user.User); ok {
				username = v.Username
			} else {
				username = "anonymous"
			}
			var act string
			if strings.ToLower(c.Request.Method) == "get" {
				act = "read"
			}
			if strings.ToLower(c.Request.Method) == "post" || strings.ToLower(c.Request.Method) == "patch" || strings.ToLower(c.Request.Method) == "put" {
				act = "write"
			}
			fmt.Println(username)
			ok, err := a.enforce(username, c.HandlerName(), act) //#TODO: Revisit
			if err != nil || !ok {
				return false
			}
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

//Session
//func (a *authController) Login(c *gin.Context) {
//	session := sessions.Default(c)
//	username := c.PostForm("username")
//	password := c.PostForm("password")
//
//	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
//		return
//	}
//	u, err := a.userService.GetByUsername(username)
//	if err != nil{
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist"})
//		return
//	}
//	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
//		return
//	}
//	session.Set("user", u.ID)
//	if err := session.Save(); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
//}
//
//func (a *authController) Logout(c *gin.Context) {
//	session := sessions.Default(c)
//	u := session.Get("user")
//	if u == nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
//		return
//	}
//	session.Delete("user")
//	if err := session.Save(); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
//}
//
//
//func (a *authController) AuthRequired(c *gin.Context) {
//	session := sessions.Default(c)
//	u := session.Get("user")
//	if u == nil {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//		return
//	}
//	c.Next()
//}
//
//func (a *authController)Authorize(obj string, act string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		val, existed := c.Get("user")
//		if !existed {
//			val = "anonymous"
//		}
//		ok, err := a.enforce(val.(string), obj, act)
//		if err != nil {
//			log.Println(err)
//			c.AbortWithStatusJSON(500, gin.H{
//				"message": "error occurred when authorizing user",
//			})
//			return
//		}
//		if !ok {
//			c.AbortWithStatusJSON(403, gin.H{
//				"message": "Unauthorized!",
//			})
//			return
//		}
//		c.Next()
//	}
//}

//func (a *authController)AuthMiddleware(sub string, obj string, act string) (bool, error) {
//	enforcer, err := casbin.NewEnforcer(a.rbac.ConfigPath, a.rbac.Adapter)
//	if err != nil {
//		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
//	}
//	err = enforcer.LoadPolicy()
//	if err != nil {
//		return false, fmt.Errorf("failed to load policy from DB: %w", err)
//	}
//	ok, err := enforcer.Enforce(sub, obj, act)
//	return ok, err
//}
