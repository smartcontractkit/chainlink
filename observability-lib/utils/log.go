package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelEnvVar = "DASHBOARD_LOG_LEVEL"
)

var (
	Logger zerolog.Logger
)

func init() {
	lvlStr := os.Getenv(LogLevelEnvVar)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}
