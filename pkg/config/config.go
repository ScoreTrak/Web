package config

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/storage"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/jinzhu/configor"
)

type StaticConfig struct {
	WebDB     storage.Config
	Logger    logger.Config
	Policy    policy.ClientConfig
	WebPort   string `default:"44444"`
	ScoreTrak ScoreTrakConfig
	Prod      bool   `default:"false"`
	Secret    string `default:"changeme"`
}

type ScoreTrakConfig struct {
	Token  string `default:""`
	Host   string `default:"scoretrak"`
	Scheme string `default:"http"`
	Port   string `default:"33333"`
}

var staticConfig StaticConfig

func GetLoggerConfig() logger.Config {
	return staticConfig.Logger
}

func GetDBConfig() storage.Config {
	return staticConfig.WebDB
}

func GetPolicyConfig() policy.ClientConfig {
	return staticConfig.Policy
}

func GetToken() string {
	return staticConfig.ScoreTrak.Token
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
