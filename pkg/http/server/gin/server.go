package gin

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/image"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/team"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
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
	err := ds.cont.Invoke(func(d *gorm.DB) {
		db = d
	})
	//TODO: For gorm v2 Implement in Scoretrak Table prefixes
	if err != nil {
		return err
	}
	db.AutoMigrate(&team.Team{})
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&image.Image{})
	return nil
}
