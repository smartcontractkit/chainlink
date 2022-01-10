package logger

// Based on https://stackoverflow.com/a/52737940

import (
	"log"

	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
)

// TestLogger creates a logger that directs output to PrettyConsoleSink configured
// for test output, and to the buffer testMemoryLog. t is optional.
// Log level is derived from the LOG_LEVEL env var.
func TestLogger(t T) Logger {
	cfg := newTestConfig()
	ll, invalid := envvar.LogLevel.ParseLogLevel()
	cfg.Level.SetLevel(ll)
	l, err := newZapLogger(cfg)
	if err != nil {
		if t == nil {
			log.Fatal(err)
		}
		t.Fatal(err)
	}
	if invalid != "" {
		l.Error(invalid)
	}
	if t == nil {
		return l
	}
	return l.Named(t.Name())
}

func newTestConfig() zap.Config {
	config := newBaseConfig()
	config.OutputPaths = []string{"pretty://console", "memory://"}
	return config
}

type T interface {
	Name() string
	Fatal(...interface{})
}
