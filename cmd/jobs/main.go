package main

import (
	"flag"
	"fmt"
	cutil "github.com/ScoreTrak/ScoreTrak/pkg/config/util"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/cmd/web/server/gin"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/di"
	"github.com/ScoreTrak/Web/pkg/storage/orm/util"
	"gorm.io/gorm"
	"os"
)

func main() {
	flag.String("config", "configs/config.yml", "Please enter a path to config file")
	flag.String("encoded-config", "", "Please enter encoded config")
	skipBootstrap := flag.Bool("skip-bootstrap", false, "Specify this flag if you want to skip the setup of tables, users, and teams(This operation is idempotent)")
	preloadData := flag.Bool("preload-data", false, "Specify this flag if you want to preload sample data into the database. This should be used AFTER preload-data on ScoreTrak")
	flag.Parse()
	path, err := cutil.ConfigFlagParser()
	handleErr(config.NewStaticConfig(path))
	d, err := di.BuildMasterContainer()
	handleErr(err)
	var l logger.LogInfoFormat
	di.Invoke(func(log logger.LogInfoFormat) {
		l = log
	})
	svr := gin.NewServer(nil, d, l)
	var db *gorm.DB
	handleErr(d.Invoke(func(d *gorm.DB) {
		db = d
	}))
	if !*skipBootstrap {
		handleErr(svr.LoadTables(db))
		handleErr(util.CreateBlackTeam(db))
		handleErr(util.CreateAdminUser(db))
	}

	if *preloadData {
		db.Exec("INSERT INTO teams (id, name, enabled, team_index) VALUES ('11111111-1111-1111-1111-111111111111', 'TeamOne', 1)")
		db.Exec("INSERT INTO teams (id, name, enabled, team_index) VALUES ('22222222-2222-2222-2222-222222222222', 'TeamTwo', 2)")
		db.Exec("INSERT INTO teams (id, name, enabled, team_index) VALUES ('33333333-3333-3333-3333-333333333333', 'TeamThree', 3)")
		db.Exec("INSERT INTO teams (id, name, enabled, team_index) VALUES ('44444444-4444-4444-4444-444444444444', 'TeamFour', 4)")
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	} else {
		return
	}
}
