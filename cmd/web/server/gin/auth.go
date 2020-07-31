package gin

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/http/handler"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/role"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type authController struct {
	log         logger.LogInfoFormat
	userService user.Serv
	ClientStore *handler.ClientStore
	db          *gorm.DB
}

func NewAuthController(l logger.LogInfoFormat, u user.Serv, c *handler.ClientStore) *authController {
	return &authController{log: l, userService: u, ClientStore: c}
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
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"username":  v.Username,
					"team_id":   v.TeamID,
					"role":      v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &user.User{
				ID:       uuid.FromStringOrNil(claims[identityKey].(string)),
				Username: claims["username"].(string),
				TeamID:   uuid.FromStringOrNil(claims["team_id"].(string)),
				Role:     claims["role"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return &user.User{
					ID:       uuid.Nil,
					Username: role.Anonymous,
					TeamID:   uuid.Nil,
					Role:     role.Anonymous,
				}, nil
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
				ID:       usr.ID,
				Username: loginVals.Username,
				TeamID:   usr.TeamID,
				Role:     usr.Role,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			v, ok := data.(*user.User)
			if ok {
				c.Set("team_id", v.TeamID)
				c.Set("role", v.Role)
				if v.Role == role.Black {
					return true
				} else if v.Role == role.Blue {
					if c.Request.Method == "GET" {
						rts := []string{"/api/property/", "/api/service/", "/api/host/", "/api/check/", "/api/last_non_elapsing/", "/api/report/"}
						pre, ok := containsPrefix(c.Request.URL.String(), rts)
						if ok {
							switch pre {
							case "/api/property/":
								_, err := handler.UuidResolver(c, "id")
								if err == nil {
									return true
								}
							case "/api/service/":
								_, err := handler.UuidResolver(c, "id")
								if err == nil {
									return true
								}
							case "/api/host/":
								_, err := handler.UuidResolver(c, "id")
								if err == nil {
									return true
								}
							case "/api/check/":
								_, err := handler.UintResolver(c, "RoundID")
								_, err2 := handler.UuidResolver(c, "ServiceID")
								if err == nil && err2 == nil {
									return true
								}
							case "/api/last_non_elapsing/":
								return true
							case "/api/report/":
								return true
							}
						}
					} else if c.Request.Method == "PATCH" {
						if strings.HasPrefix(c.Request.URL.String(), "/api/property/") {
							_, err := handler.UuidResolver(c, "id")
							if err == nil {
								return true
							}
						}
						if strings.HasPrefix(c.Request.URL.String(), "/api/user/") {
							_, err := handler.UuidResolver(c, "id")
							if err == nil {
								p := a.ClientStore.PolicyClient.GetPolicy()
								if *p.AllowChangingUsernamesAndPasswords {
									return true
								}
							}
						}
						if strings.HasPrefix(c.Request.URL.String(), "/api/host/") {
							_, err := handler.UuidResolver(c, "id")
							if err == nil {
								return true
							}
						}
					}
				} else if v.Role == role.Anonymous && c.Request.Method == "GET" && c.Request.URL.String() == "/api/report/" {
					p := a.ClientStore.PolicyClient.GetPolicy()
					if *p.AllowUnauthenticatedUsers {
						return true
					}
				}
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

func containsPrefix(s string, prefix []string) (string, bool) {
	for _, pre := range prefix {
		if strings.HasPrefix(s, pre) {
			return pre, true
		}
	}
	return "", false
}
