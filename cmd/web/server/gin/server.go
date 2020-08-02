package gin

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/image"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
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
