package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewRouter() *gin.Engine {
	return gin.Default()
}

func (ds *dserver) MapRoutes() {
	ds.router.GET("/", ReadResource)
}

func ReadResource(c *gin.Context) {
	c.JSON(200, "Test!")
}
