package plugins

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func NewLogger(name string, cfg LoggingConfig) (logger.Logger, func()) {
	lcfg := logger.Config{
		LogLevel:    cfg.Level(),
		JsonConsole: cfg.JSONConsole(),
		UnixTS:      cfg.UnixTimestamps(),
	}
	lggr, closeLggr := lcfg.New()
	lggr = lggr.Named(name)
	return lggr, func() {
		if err := closeLggr(); err != nil {
			fmt.Println("Failed to close logger:", err)
		}
	}
}
