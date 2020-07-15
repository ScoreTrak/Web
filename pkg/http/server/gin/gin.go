package gin

import (
	"fmt"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/http/handler"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/rbac"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/appleboy/gin-jwt/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/nosurf"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func NewRouter() *gin.Engine {
	if config.GetStaticConfig().Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	return router
}

func (ds *dserver) MapRoutesAndStart() error {
	authMiddleware, err := ds.authBootstrap()
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	csrf := nosurf.New(ds.router)
	csrf.SetFailureHandler(http.HandlerFunc(csrfFailHandler))

	ds.router.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"csrf_token": nosurf.Token(c.Request)})
	})
	ds.router.POST("/login", authMiddleware.LoginHandler)
	ds.router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		ds.logger.Info(claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := ds.router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	api := ds.router.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	var cfg config.StaticConfig
	ds.cont.Invoke(func(c config.StaticConfig) { cfg = c })
	return http.ListenAndServe(":"+cfg.Port, csrf)
}

func (ds *dserver) authBootstrap() (*jwt.GinJWTMiddleware, error) {
	var db *gorm.DB
	err := ds.cont.Invoke(func(d *gorm.DB) {
		db = d
	})
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDBUsePrefix(db, "cas_")
	if err != nil {
		return nil, err
	}
	var us user.Serv
	err = ds.cont.Invoke(func(u user.Serv) {
		us = u
	})
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = us.Store(&user.User{ID: 1, Username: "admin", PasswordHash: string(hashedPassword)})
	if err != nil {
		serr, ok := err.(*pq.Error)
		if !ok || serr.Code.Name() != "unique_violation" {
			return nil, err
		}
	}
	r := rbac.NewRBAC(adapter, "configs/rbac_model.conf")
	authCtrl := handler.NewAuthController(ds.logger, us, r)
	ds.router.Use(static.Serve("/", static.LocalFile("./views", true)))
	authMiddleware, err := authCtrl.JWTMiddleware()
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}

func csrfFailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}
