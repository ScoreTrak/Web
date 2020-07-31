package di

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/storage"
	"github.com/L1ghtman2k/ScoreTrakWeb/cmd/web/server/gin"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/config"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/policy"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/storage/orm"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/team"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"go.uber.org/dig"
)

var container = dig.New()

func BuildMasterContainer() (*dig.Container, error) {
	var ctr []interface{}

	ctr = append(ctr,
		config.GetStaticConfig, config.GetDBConfig, config.GetLoggerConfig,
		storage.LoadDB, logger.NewLogger, gin.NewAuthController,
		policy.NewPolicyServ, orm.NewPolicyRepo,
		user.NewUserServ, orm.NewUserRepo,
		team.NewTeamServ, orm.NewTeamRepo,
	)
	for _, i := range ctr {
		err := container.Provide(i)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func Invoke(i interface{}) {
	err := container.Invoke(i)
	if err != nil {
		panic(err)
	}
}
