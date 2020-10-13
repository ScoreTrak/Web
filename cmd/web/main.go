package main

import (
	"flag"
	"fmt"
	cutil "github.com/ScoreTrak/ScoreTrak/pkg/config/util"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/cmd/web/server/gin"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/di"
	"gorm.io/gorm"
	"os"
)

func main() {
	flag.String("config", "configs/config.yml", "Please enter a path to config file")
	flag.String("encoded-config", "", "Please enter encoded config")
	flag.Parse()
	path, err := cutil.ConfigFlagParser()
	if err != nil {
		handleErr(err)
	}
	handleErr(config.NewStaticConfig(path))
	r := gin.NewRouter()
	d, err := di.BuildMasterContainer()
	handleErr(err)
	var l logger.LogInfoFormat
	di.Invoke(func(log logger.LogInfoFormat) {
		l = log
	})
	svr := gin.NewServer(r, d, l)
	var db *gorm.DB
	handleErr(d.Invoke(func(d *gorm.DB) {
		db = d
	}))
	handleErr(svr.LoadTables(db))
	handleErr(svr.MapRoutesAndStart())
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	} else {
		return
	}
}
