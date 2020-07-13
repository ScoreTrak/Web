package gin

import (
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"
)

type dserver struct {
	router *gin.Engine
	cont   *dig.Container
	logger logger.LogInfoFormat
}

func NewServer(e *gin.Engine, c *dig.Container, l logger.LogInfoFormat) *dserver {
	return &dserver{
		router: e,
		cont:   c,
		logger: l,
	}
}

func (ds *dserver) SetupDB() error {
	var db *gorm.DB
	ds.cont.Invoke(func(d *gorm.DB) {
		db = d
	})

	//var tm time.Time
	//
	//db.AutoMigrate(&team.Team{})
	return nil
}

// Start start serving the application
func (ds *dserver) Start() error {
	var cfg config.StaticConfig
	ds.cont.Invoke(func(c config.StaticConfig) { cfg = c })
	return ds.router.Run(fmt.Sprintf(":%s", cfg.Port))
}
