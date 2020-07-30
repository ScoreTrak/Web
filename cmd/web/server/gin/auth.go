package gin

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/check"
	"github.com/L1ghtman2k/ScoreTrak/pkg/host"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/property"
	"github.com/L1ghtman2k/ScoreTrak/pkg/service"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/http/handler"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type authController struct {
	log         logger.LogInfoFormat
	userService user.Serv
	ClientStore handler.ClientStore
}

func NewAuthController(l logger.LogInfoFormat, u user.Serv, c handler.ClientStore) *authController {
	return &authController{l, u, c}
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
				ID:       uuid.FromStringOrNil(claims["id"].(string)),
				Username: claims[identityKey].(string),
				TeamID:   uuid.FromStringOrNil(claims["team_id"].(string)),
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
				ID:       usr.ID,
				Username: loginVals.Username,
				TeamID:   usr.TeamID,
				Role:     usr.Role,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*user.User); ok && v.Role == "black" {
				return true
			} else if v.Role == "blue" {
				if c.Request.Method == "GET" {
					rts := []string{"/api/property/", "/api/service/", "/api/host/", "/api/check/", "/api/last_non_elapsing/", "/api/report/"}
					pre, ok := containsPrefix(c.Request.URL.String(), rts)
					if ok {
						switch pre {
						case "/api/property/":
							id, err := handler.UuidResolver(c, "id")
							if err == nil {
								tID, prop, err := a.teamIDFromProperty(id)
								if err == nil && prop != nil && tID == v.TeamID && prop.Status != property.Hide {
									c.Set("shortcut", prop)
									return true
								}
							}
						case "/api/service/":
							id, err := handler.UuidResolver(c, "id")
							if err == nil {
								tID, serv, err := a.teamIDFromService(id)
								if err == nil && serv != nil && tID == v.TeamID {
									c.Set("shortcut", serv)
									return true
								}
							}
						case "/api/host/":
							id, err := handler.UuidResolver(c, "id")
							if err == nil {
								tID, serv, err := a.teamIDFromHost(id)
								if err == nil && serv != nil && tID == v.TeamID {
									c.Set("shortcut", serv)
									return true
								}
							}
						case "/api/check/":
							rid, err := handler.UintResolver(c, "RoundID")
							sid, err2 := handler.UuidResolver(c, "ServiceID")
							if err == nil && err2 == nil {
								tID, ck, err := a.teamIDFromCheck(rid, sid)
								if err == nil && ck != nil && tID == v.TeamID {
									c.Set("shortcut", ck)
									return true
								}
							}
						case "/api/last_non_elapsing/":
							return true
						case "/api/report/":
							c.Set("team_id", v.TeamID)
							return true
						}
					}
				} else if c.Request.Method == "PATCH" {
					if strings.HasPrefix(c.Request.URL.String(), "/api/property/") {
						id, err := handler.UuidResolver(c, "id")
						if err == nil {
							tID, p, err := a.teamIDFromProperty(id)
							if err == nil && p != nil && tID == v.TeamID && p.Status == property.Edit {
								us := &property.Property{}
								err = c.BindJSON(us)
								if err == nil {
									c.Set("filtered", &property.Property{ID: us.ID, Value: us.Value})
									return true
								}
							}
						}
					}
					if strings.HasPrefix(c.Request.URL.String(), "/api/user/") {
						id, err := handler.UuidResolver(c, "id")
						if err == nil && v.ID == id {
							us := &user.User{}
							err = c.BindJSON(us)
							if err == nil {
								c.Set("filtered", &user.User{ID: us.ID, Username: us.Username, Password: us.Password})
								return true
							}
						}
					}
					if strings.HasPrefix(c.Request.URL.String(), "/api/host/") {
						id, err := handler.UuidResolver(c, "id")
						if err == nil {
							tID, p, err := a.teamIDFromHost(id)
							if err == nil && p != nil && tID == v.TeamID && p.EditHost != nil && *p.EditHost == true {
								us := &host.Host{}
								err = c.BindJSON(us)
								if err == nil {
									c.Set("filtered", &host.Host{ID: us.ID, Address: us.Address})
									return true
								}
							}
						}
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

func (a *authController) teamIDFromProperty(propertyID uuid.UUID) (teamID uuid.UUID, property *property.Property, err error) {
	property, err = a.ClientStore.PropertyClient.GetByID(propertyID)
	if err != nil || property == nil {
		return
	}
	teamID, _, err = a.teamIDFromService(property.ServiceID)
	return
}

func (a *authController) teamIDFromCheck(roundID uint, serviceID uuid.UUID) (teamID uuid.UUID, check *check.Check, err error) {
	check, err = a.ClientStore.CheckClient.GetByRoundServiceID(roundID, serviceID)
	if err != nil || check == nil {
		return
	}
	teamID, _, err = a.teamIDFromService(check.ServiceID)
	return
}

func (a *authController) teamIDFromService(serviceID uuid.UUID) (teamID uuid.UUID, service *service.Service, err error) {
	service, err = a.ClientStore.ServiceClient.GetByID(serviceID)
	if err != nil || service == nil {
		return
	}
	teamID, _, err = a.teamIDFromHost(service.HostID)
	return
}

func (a *authController) teamIDFromHost(hostID uuid.UUID) (teamID uuid.UUID, host *host.Host, err error) {
	host, err = a.ClientStore.HostClient.GetByID(hostID)
	if err != nil || host == nil {
		return
	}
	return host.TeamID, host, err
}

//TODO: Make sure patching of properties, hosts, etc only changes certain fields instead of everything
