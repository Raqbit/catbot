package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

const (
	EnvPrefix    = "CATBOT"
	BotTokenEnv  = "BOT_TOKEN"
	CmdPrefixEnv = "COMMAND_PREFIX"
	DebugEnv     = "DEBUG"
	DbSrcEnv     = "DATABASE_SOURCE"
	CatCostEnv   = "CAT_COST"
)

type Config struct {
	BotToken      string
	CommandPrefix string
	Debug         bool
	DbDriver      string
	DbSrc         string
	CatCost       int64
}

func NewConfig() (*Config, error) {

	botToken := os.Getenv(getEnvName(BotTokenEnv))

	if botToken == "" {
		return nil, errors.New("No discord bot-token provided")
	}

	cmdPrefix := os.Getenv(getEnvName(CmdPrefixEnv))

	if cmdPrefix == "" {
		return nil, errors.New("No command prefix provided")
	}

	debug := false
	debugString := os.Getenv(getEnvName(DebugEnv))

	if debugString == "1" {
		debug = true
	}

	dbSrc := os.Getenv(getEnvName(DbSrcEnv))

	if dbSrc == "" {
		return nil, errors.New("No database source provided")
	}

	catCostStr := os.Getenv(getEnvName(CatCostEnv))

	if catCostStr == "" {
		return nil, errors.New("No cat cost provided")
	}

	catCost, err := strconv.ParseInt(catCostStr, 10, 64)

	if err != nil {
		return nil, errors.New("Could not parse cat cost")
	}

	return &Config{
		BotToken:      botToken,
		CommandPrefix: cmdPrefix,
		Debug:         debug,
		DbSrc:         dbSrc,
		CatCost:       catCost,
	}, nil
}

func getEnvName(name string) string {
	return EnvPrefix + "_" + name
}
