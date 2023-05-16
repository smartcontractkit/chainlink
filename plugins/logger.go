package plugins

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func NewLogger(cfg LoggingConfig) (logger.Logger, func()) {
	lcfg := logger.Config{
		LogLevel:    cfg.LogLevel(),
		JsonConsole: cfg.JSONConsole(),
		UnixTS:      cfg.LogUnixTimestamps(),
	}
	lggr, closeLggr := lcfg.New()
	lggr = lggr.Named("PluginSolana")
	return lggr, func() {
		if err := closeLggr(); err != nil {
			fmt.Println("Failed to close logger:", err)
		}
	}
}
