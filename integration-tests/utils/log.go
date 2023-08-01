package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func SetupGlobalLogger() {
	lvlStr := os.Getenv("GLOBAL_LOG_LEVEL")
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}
