package gin

import (
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/http/handler"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/role"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/team"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"net/url"
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
	conf := config.GetStaticConfig()
	c := client.NewScoretrakClient(&url.URL{Host: fmt.Sprintf("localhost:%s", conf.ScoreTrakPort), Scheme: conf.ScoreTrakScheme}, conf.Token, http.DefaultClient)

	cStore := handler.ClientStore{
		ConfigClient:       client.NewConfigClient(c),
		TeamClient:         client.NewTeamClient(c),
		HostClient:         client.NewHostClient(c),
		ServiceClient:      client.NewServiceClient(c),
		ServiceGroupClient: client.NewServiceGroupClient(c),
		HostGroupClient:    client.NewHostGroupClient(c),
		PropertyClient:     client.NewPropertyClient(c),
		RoundClient:        client.NewRoundClient(c),
		CheckClient:        client.NewCheckClient(c),
		ReportClient:       client.NewReportClient(c),
	}

	authMiddleware, err := ds.authBootstrap(cStore)
	if err != nil {
		ds.logger.Error(err)
		return err
	}

	ds.router.Use(static.Serve("/", static.LocalFile("./views", true)))
	ds.router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		ds.logger.Info(claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := ds.router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.POST("/login", authMiddleware.LoginHandler)
	api := ds.router.Group("/api")

	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		teamRoute := api.Group("/team")
		{
			var tsvc team.Serv
			err := ds.cont.Invoke(func(svc team.Serv) {
				tsvc = svc
			})
			if err != nil {
				return err
			}
			tctrl := handler.NewTeamController(ds.logger, tsvc, cStore.TeamClient)
			teamRoute.GET("/", tctrl.GetAll)
			teamRoute.POST("/", tctrl.Store)
			teamRoute.GET("/:id", tctrl.GetByID)
			teamRoute.PATCH("/:id", tctrl.Update)
			teamRoute.DELETE("/:id", tctrl.Delete)
		}

		userRoute := api.Group("/user")
		{
			var usvc user.Serv
			err := ds.cont.Invoke(func(svc user.Serv) {
				usvc = svc
			})
			if err != nil {
				return err
			}
			uctrl := handler.NewUserController(ds.logger, usvc)
			userRoute.GET("/", uctrl.GetAll)
			userRoute.POST("/", uctrl.Store)
			userRoute.GET("/:id", uctrl.GetByID)
			userRoute.PATCH("/:id", uctrl.Update)
			userRoute.DELETE("/:id", uctrl.Delete)
		}

		serviceRoute := api.Group("/service")
		{
			sctrl := handler.NewServiceController(ds.logger, cStore.ServiceClient)
			serviceRoute.GET("/", sctrl.GetAll)
			serviceRoute.POST("/", sctrl.Store)
			serviceRoute.GET("/:id", sctrl.GetByID)
			serviceRoute.PATCH("/:id", sctrl.Update)
			serviceRoute.DELETE("/:id", sctrl.Delete)
		}

		{

			roundRoute := api.Group("/round")
			rctrl := handler.NewRoundController(ds.logger, cStore.RoundClient)
			roundRoute.GET("/", rctrl.GetAll)
			roundRoute.GET("/:id", rctrl.GetByID)
			lastRoundRoute := api.Group("/last_non_elapsing")
			lastRoundRoute.GET("/", rctrl.GetLastNonElapsingRound)

		}

		serviceGroupRoute := api.Group("/service_group")
		{
			sctrl := handler.NewServiceGroupController(ds.logger, cStore.ServiceGroupClient)
			serviceGroupRoute.GET("/", sctrl.GetAll)
			serviceGroupRoute.POST("/", sctrl.Store)
			serviceGroupRoute.GET("/:id", sctrl.GetByID)
			serviceGroupRoute.PATCH("/:id", sctrl.Update)
			serviceGroupRoute.DELETE("/:id", sctrl.Delete)
		}

		hostGroupRoute := api.Group("/host_group")
		{
			hctrl := handler.NewHostGroupController(ds.logger, cStore.HostGroupClient)
			hostGroupRoute.GET("/", hctrl.GetAll)
			hostGroupRoute.POST("/", hctrl.Store)
			hostGroupRoute.GET("/:id", hctrl.GetByID)
			hostGroupRoute.PATCH("/:id", hctrl.Update)
			hostGroupRoute.DELETE("/:id", hctrl.Delete)
		}

		hostRoute := api.Group("/host")
		{
			hctrl := handler.NewHostController(ds.logger, cStore.HostClient)
			hostRoute.GET("/", hctrl.GetAll)
			hostRoute.POST("/", hctrl.Store)
			hostRoute.GET("/:id", hctrl.GetByID)
			hostRoute.PATCH("/:id", hctrl.Update)
			hostRoute.DELETE("/:id", hctrl.Delete)
		}

		{
			hctrl := handler.NewCheckController(ds.logger, cStore.CheckClient)
			api.GET("/check/:RoundID/:ServiceID", hctrl.GetByRoundServiceID)
			api.GET("/check_all/:id", hctrl.GetAllByRoundID)
		}

		configRoute := api.Group("/config")
		{
			hctrl := handler.NewConfigController(ds.logger, cStore.ConfigClient)
			configRoute.GET("/", hctrl.Get)
			configRoute.PATCH("/", hctrl.Update)
		}

		propertyRoute := api.Group("/property")
		{
			hctrl := handler.NewPropertyController(ds.logger, cStore.PropertyClient)
			propertyRoute.GET("/", hctrl.GetAll)
			propertyRoute.POST("/", hctrl.Store)
			propertyRoute.GET("/:id", hctrl.GetByID)
			propertyRoute.PATCH("/:id", hctrl.Update)
			propertyRoute.DELETE("/:id", hctrl.Delete)
		}

		reportRoute := api.Group("/report")
		{
			hctrl := NewReportController(ds.logger, cStore.ReportClient, cStore.ConfigClient, cStore.RoundClient)
			reportRoute.GET("/", hctrl.Get)
			go hctrl.LazyUpdate(c)
		}

	}
	return ds.router.Run(fmt.Sprintf(":%s", conf.WebPort))
}

func (ds *dserver) authBootstrap(clientStore handler.ClientStore) (*jwt.GinJWTMiddleware, error) {
	var db *gorm.DB
	err := ds.cont.Invoke(func(d *gorm.DB) {
		db = d
	})

	if err != nil {
		return nil, err
	}

	var ts team.Serv
	err = ds.cont.Invoke(func(u team.Serv) {
		ts = u
	})

	if err != nil {
		return nil, err
	}
	uuid1 := uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")
	err = ts.Store([]*team.Team{{ID: uuid1, Name: "Black Team"}})
	if err != nil {
		serr, ok := err.(*pgconn.PgError)
		if !ok || serr.Code != "23505" {
			return nil, err
		}
	}

	var us user.Serv
	err = ds.cont.Invoke(func(u user.Serv) {
		us = u
	})

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = us.Store([]*user.User{{ID: uuid1, TeamID: uuid1, Username: "admin", Role: role.Black, PasswordHash: string(hashedPassword)}})
	if err != nil {
		serr, ok := err.(*pgconn.PgError)
		if !ok || serr.Code != "23505" {
			return nil, err
		}
	}

	authCtrl := NewAuthController(ds.logger, us, clientStore)
	authMiddleware, err := authCtrl.JWTMiddleware()
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}
