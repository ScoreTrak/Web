package gin

import (
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/http/handler"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/rbac"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/appleboy/gin-jwt/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func NewRouter() *gin.Engine {
	if config.GetStaticConfig().Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	return router
}

func (ds *dserver) MapRoutes() error {
	authMiddleware, err := ds.authBootstrap()
	if err != nil {
		ds.logger.Error(err)
		return err
	}
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
	return nil
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

//Sessions
//func (ds *dserver) MapRoutes() error {
//	var db *gorm.DB
//	err := ds.cont.Invoke(func(d *gorm.DB) {
//		db = d
//	})
//	if err != nil{
//		return err
//	}
//	adapter, err := gormadapter.NewAdapterByDBUsePrefix(db, "cas_")
//	if err != nil{
//		return err
//	}
//	var us user.Serv
//	err = ds.cont.Invoke(func(u user.Serv) {
//		us = u
//	})
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
//	if err != nil{
//		return err
//	}
//	err = us.Store(&user.User{Username: "admin", PasswordHash: string(hashedPassword)})
//	if err != nil {
//		serr, ok := err.(*pq.Error)
//		if !ok || serr.Code.Name() != "unique_violation" {
//			return err
//		}
//	}
//	if err != nil{
//		return err
//	}
//
//
//
//
//	r := rbac.NewRBAC(adapter, "configs/rbac_model.conf")
//	authCtrl := handler.NewAuthController(ds.logger, us, r)
//	ds.router.Use(sessions.Sessions("scoretrak_session", sessions.NewCookieStore([]byte(config.GetStaticConfig().Secret))))
//	ds.router.Use(static.Serve("/", static.LocalFile("./views", true)))
//	auth := ds.router.Group("/auth")
//	{
//		auth.POST("/login", authCtrl.Login)
//		auth.GET("/logout", authCtrl.Logout)
//	}
//	api := ds.router.Group("/api")
//	api.Use(authCtrl.AuthRequired)
//	{
//		api.GET("/", authCtrl.Authorize("welcome", "read"), func(c *gin.Context) {
//			c.JSON(200, gin.H{
//				"message": "pong",
//			})
//		})
//	}
//	return nil
//}
