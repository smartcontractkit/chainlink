package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-plugin"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
)

func main() {
	logLevelStr := os.Getenv("CL_LOG_LEVEL")
	logLevel, err := zapcore.ParseLevel(logLevelStr)
	if err != nil {
		fmt.Printf("failed to parse CL_LOG_LEVEL = %q: %s\n", logLevelStr, err)
		os.Exit(1)
	}
	cfg := logger.Config{
		LogLevel:    logLevel,
		JsonConsole: strings.EqualFold("true", os.Getenv("CL_JSON_CONSOLE")),
		UnixTS:      strings.EqualFold("true", os.Getenv("CL_UNIX_TS")),
	}
	lggr, closeLggr := cfg.New()
	lggr = lggr.Named("PluginMedian")
	defer func() {
		if err := closeLggr(); err != nil {
			fmt.Println("Failed to close logger:", err)
		}
	}()
	stop := make(chan struct{})
	defer close(stop)
	mp := median.NewPlugin(lggr, stop)
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginMedianHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginMedianName: loop.NewGRPCPluginMedian(mp, lggr),
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
