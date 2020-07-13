package gin

import (
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	if config.GetStaticConfig().Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	return gin.Default()
}

func (ds *dserver) MapRoutes() {

	ds.router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := ds.router.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
