package config

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/config"
	"github.com/jinzhu/configor"
)

type StaticConfig struct {
	DB            config.DB
	Logger        config.Logger
	Port          string `default:"44444"`
	ScoreTrakPort string `default:"33333"`
	Token         string `default:""`
	ScoreTrakURL  string `default:"http://scoretrak/"`
}

var staticConfig StaticConfig

func GetLoggerConfig() config.Logger {
	return staticConfig.Logger
}

func GetDBConfig() config.DB {
	return staticConfig.DB
}

func GetToken() string {
	return staticConfig.Token
}

func GetStaticConfig() StaticConfig {
	return staticConfig
}

func NewStaticConfig(f string) error {
	err := configor.Load(&staticConfig, f)
	if err != nil {
		return err
	}
	return nil
}
