package gin

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/api/client"
	"github.com/ScoreTrak/Web/pkg/competition"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/di/repo"
	"github.com/ScoreTrak/Web/pkg/http/handler"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/storage/orm/util"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
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
	return router
}

func (ds *dserver) MapRoutesAndStart() error {
	conf := config.GetStaticConfig()
	c := client.NewScoretrakClient(&url.URL{Host: conf.ScoreTrak.Host + ":" + conf.ScoreTrak.Port, Scheme: conf.ScoreTrak.Scheme}, conf.ScoreTrak.Token, http.DefaultClient)

	cStore := &handler.ClientStore{
		StaticConfigClient: client.NewStaticConfigClient(c),
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
		CompetitionClient:  client.NewCompetitionClient(c),
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

	ds.router.Use(static.Serve("/", static.LocalFile("./views/build", true)))

	ds.router.NoRoute(func(c *gin.Context) {
		c.File("./views/build/index.html")
	})

	//ds.router.NoRoute(func(c *gin.Context) {
	//	claims := jwt.ExtractClaims(c)
	//	ds.logger.Info(claims)
	//	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	//})

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
		serviceTestRoute := api.Group("/service_test")
		{
			sctrl := handler.NewServiceController(ds.logger, cStore)
			serviceRoute.GET("/", sctrl.GetAll)
			serviceRoute.POST("/", sctrl.Store)
			serviceRoute.GET("/:id", sctrl.GetByID)
			serviceTestRoute.GET("/:id", sctrl.TestService)

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
		serviceGroupRedeploy := api.Group("/service_group_redeploy")
		{
			sctrl := handler.NewServiceGroupController(ds.logger, cStore)
			serviceGroupRoute.GET("/", sctrl.GetAll)
			serviceGroupRoute.POST("/", sctrl.Store)
			serviceGroupRoute.GET("/:id", sctrl.GetByID)
			serviceGroupRoute.PATCH("/:id", sctrl.Update)
			serviceGroupRoute.DELETE("/:id", sctrl.Delete)
			serviceGroupRedeploy.GET("/:id", sctrl.Redeploy)
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
			api.GET("/check_round/:id", hctrl.GetAllByRoundID)
			api.GET("/check_service/:id", hctrl.GetAllByServiceID)
		}

		configRoute := api.Group("/config")
		{
			hctrl := handler.NewConfigController(ds.logger, cStore)
			configRoute.GET("/", hctrl.Get)
			configRoute.PATCH("/", hctrl.Update)
			configRoute.DELETE("/reset_competition", hctrl.ResetScores)
			configRoute.DELETE("/delete_competition", hctrl.DeleteCompetition)
			configRoute.GET("/static_config", hctrl.GetStaticConfig)
			configRoute.GET("/static_web_config", hctrl.GetStaticWebConfig)

		}

		{
			hctrl := handler.NewPropertyController(ds.logger, cStore)

			api.GET("/properties", hctrl.GetAll)
			api.GET("/properties/:ServiceID", hctrl.GetAllByServiceID)
			api.POST("/property", hctrl.Store)
			api.GET("/property/:ServiceID/:Key", hctrl.GetByServiceIDKey)
			api.DELETE("/property/:ServiceID/:Key", hctrl.Delete)
			api.PATCH("/property/:ServiceID/:Key", hctrl.Update)
		}

		reportRoute := api.Group("/report")
		{
			hctrl := handler.NewReportController(ds.logger, cStore)
			reportRoute.GET("/", hctrl.Get)
			reportRoute.GET("/:id", hctrl.GetByTeamID)
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

		competitionRoute := api.Group("/competition")
		{

			competitionServ := competition.NewCompetitionServ(repo.NewStore())
			hctrl := handler.NewCompetitionController(ds.logger, cStore, competitionServ)
			competitionRoute.GET("/export_core", hctrl.FetchCoreCompetition)
			competitionRoute.GET("/export_all", hctrl.FetchEntireCompetition)
			competitionRoute.POST("/upload", hctrl.LoadCompetition)
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
