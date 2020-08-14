package di

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/storage"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/storage/orm"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
	"go.uber.org/dig"
)

var container = dig.New()

func BuildMasterContainer() (*dig.Container, error) {
	var ctr []interface{}

	ctr = append(ctr,
		config.GetStaticConfig, config.GetDBConfig, config.GetLoggerConfig,
		storage.LoadDB, logger.NewLogger,
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
