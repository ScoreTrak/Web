package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/cmd/web/server/gin"
	"github.com/ScoreTrak/Web/pkg/config"
	"github.com/ScoreTrak/Web/pkg/di"
	"os"
)

func main() {
	path := flag.String("config", "configs/config.yml", "Please enter a path to config file")
	flag.Parse()
	if !configExists(*path) {
		handleErr(errors.New("you need to provide config"))
	}
	handleErr(config.NewStaticConfig(*path))
	r := gin.NewRouter()
	d, err := di.BuildMasterContainer()
	handleErr(err)
	var l logger.LogInfoFormat
	di.Invoke(func(log logger.LogInfoFormat) {
		l = log
	})
	svr := gin.NewServer(r, d, l)
	handleErr(svr.SetupDB())
	handleErr(svr.MapRoutesAndStart())
}

func configExists(f string) bool {
	file, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return !file.IsDir()
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	} else {
		return
	}
}
