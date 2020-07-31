package gin

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/image"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/policy"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/team"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
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
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&policy.Policy{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&team.Team{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&image.Image{})
	if err != nil {
		return err
	}
	return nil
}
