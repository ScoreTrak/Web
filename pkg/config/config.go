package config

import (
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/storage"
	"github.com/jinzhu/configor"
)

type StaticConfig struct {
	DB            storage.Config
	Logger        logger.Config
	Port          string `default:"44444"`
	ScoreTrakPort string `default:"33333"`
	Token         string `default:""`
	ScoreTrakURL  string `default:"http://scoretrak/"`
	Prod          bool   `default:"false"`
	Secret        string `default:"changeme"`
}

var staticConfig StaticConfig

func GetLoggerConfig() logger.Config {
	return staticConfig.Logger
}

func GetDBConfig() storage.Config {
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
