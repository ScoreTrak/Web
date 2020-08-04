package gin

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/api/client"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/http/handler"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/storage/orm/util"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
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
	c := client.NewScoretrakClient(&url.URL{Host: conf.ScoreTrakHost + ":" + conf.ScoreTrakPort, Scheme: conf.ScoreTrakScheme}, conf.Token, http.DefaultClient)

	cStore := &handler.ClientStore{
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

	authCtrl, err := ds.authBootstrap(cStore)
	if err != nil {
		ds.logger.Error(err)
		return err
	}
	authMiddleware, err := authCtrl.JWTMiddleware()
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
			tctrl := handler.NewTeamController(ds.logger, tsvc, cStore)
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
			sctrl := handler.NewServiceController(ds.logger, cStore)
			serviceRoute.GET("/", sctrl.GetAll)
			serviceRoute.POST("/", sctrl.Store)
			serviceRoute.GET("/:id", sctrl.GetByID)
			serviceRoute.PATCH("/:id", sctrl.Update)
			serviceRoute.DELETE("/:id", sctrl.Delete)
		}

		{

			roundRoute := api.Group("/round")
			rctrl := handler.NewRoundController(ds.logger, cStore)
			roundRoute.GET("/", rctrl.GetAll)
			roundRoute.GET("/:id", rctrl.GetByID)
			lastRoundRoute := api.Group("/last_non_elapsing")
			lastRoundRoute.GET("/", rctrl.GetLastNonElapsingRound)

		}

		serviceGroupRoute := api.Group("/service_group")
		{
			sctrl := handler.NewServiceGroupController(ds.logger, cStore)
			serviceGroupRoute.GET("/", sctrl.GetAll)
			serviceGroupRoute.POST("/", sctrl.Store)
			serviceGroupRoute.GET("/:id", sctrl.GetByID)
			serviceGroupRoute.PATCH("/:id", sctrl.Update)
			serviceGroupRoute.DELETE("/:id", sctrl.Delete)
		}

		hostGroupRoute := api.Group("/host_group")
		{
			hctrl := handler.NewHostGroupController(ds.logger, cStore)
			hostGroupRoute.GET("/", hctrl.GetAll)
			hostGroupRoute.POST("/", hctrl.Store)
			hostGroupRoute.GET("/:id", hctrl.GetByID)
			hostGroupRoute.PATCH("/:id", hctrl.Update)
			hostGroupRoute.DELETE("/:id", hctrl.Delete)
		}

		hostRoute := api.Group("/host")
		{
			hctrl := handler.NewHostController(ds.logger, cStore)
			hostRoute.GET("/", hctrl.GetAll)
			hostRoute.POST("/", hctrl.Store)
			hostRoute.GET("/:id", hctrl.GetByID)
			hostRoute.PATCH("/:id", hctrl.Update)
			hostRoute.DELETE("/:id", hctrl.Delete)
		}

		{
			hctrl := handler.NewCheckController(ds.logger, cStore)
			api.GET("/check/:RoundID/:ServiceID", hctrl.GetByRoundServiceID)
			api.GET("/check_all/:id", hctrl.GetAllByRoundID)
		}

		configRoute := api.Group("/config")
		{
			hctrl := handler.NewConfigController(ds.logger, cStore)
			configRoute.GET("/", hctrl.Get)
			configRoute.PATCH("/", hctrl.Update)
		}

		propertyRoute := api.Group("/property")
		{
			hctrl := handler.NewPropertyController(ds.logger, cStore)
			propertyRoute.GET("/", hctrl.GetAll)
			propertyRoute.POST("/", hctrl.Store)
			propertyRoute.GET("/:id", hctrl.GetByID)
			propertyRoute.PATCH("/:id", hctrl.Update)
			propertyRoute.DELETE("/:id", hctrl.Delete)
		}

		reportRoute := api.Group("/report")
		{
			hctrl := handler.NewReportController(ds.logger, cStore)
			reportRoute.GET("/", hctrl.Get)
			reportRoute.GET("/:id", hctrl.GetByTeamID)
			go hctrl.LazyReportLoader(c)
		}
		policyRoute := api.Group("/policy")
		{
			var psvc policy.Serv
			err := ds.cont.Invoke(func(svc policy.Serv) {
				psvc = svc
			})
			if err != nil {
				return err
			}
			pctrl := handler.NewPolicyController(ds.logger, psvc)
			policyRoute.GET("/", pctrl.GetPolicy)
			policyRoute.PATCH("/", pctrl.UpdatePolicy)
		}

	}
	return ds.router.Run(":" + conf.WebPort)
}

func (ds *dserver) authBootstrap(clientStore *handler.ClientStore) (*authController, error) {
	var db *gorm.DB
	err := ds.cont.Invoke(func(d *gorm.DB) {
		db = d
	})
	if err != nil {
		return nil, err
	}
	err = util.CreateBlackTeam(db)
	if err != nil {
		return nil, err
	}
	err = util.CreateAdminUser(db)
	if err != nil {
		return nil, err
	}
	p, err := util.CreatePolicy(db)
	if err != nil {
		return nil, err
	}
	var policyRepo policy.Repo
	err = ds.cont.Invoke(func(u policy.Repo) {
		policyRepo = u
	})
	if err != nil {
		return nil, err
	}
	var us user.Serv
	err = ds.cont.Invoke(func(u user.Serv) {
		us = u
	})
	if err != nil {
		return nil, err
	}
	clientStore.PolicyClient = policy.NewPolicyClient(p, policyRepo, config.GetPolicyConfig())
	authCtrl := NewAuthController(ds.logger, us, clientStore)
	go clientStore.PolicyClient.LazyPolicyLoader()
	return authCtrl, nil

}
