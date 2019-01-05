package main

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	BotToken       string `env:"CATBOT_BOT_TOKEN,required"`
	CommandPrefix  string `env:"CATBOT_COMMAND_PREFIX" envDefault:"c!"`
	Debug          bool   `env:"CATBOT_DEBUG" envDefault:"false"`
	DbSrc          string `env:"CATBOT_DATABASE_SOURCE,required"`
	CatCost        int64  `env:"CATBOT_CAT_COST" envDefault:"25"`
	CatMinAwayMins int64  `env:"CATBOT_MIN_AWAY" envDefault:"10"`
	CatMaxAwayMins int64  `env:"CATBOT_MAX_AWAY" envDefault:"60"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return &cfg, err
}
