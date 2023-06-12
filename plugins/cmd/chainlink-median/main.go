package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	loggerName = "PluginMedian"
)

func main() {
	envCfg, err := plugins.GetEnvConfig()
	if err != nil {
		fmt.Printf("Failed to get environment configuration: %s\n", err)
		os.Exit(1)
	}
	lggr, closeLggr := plugins.NewLogger(loggerName, envCfg)
	defer closeLggr()
	slggr := logger.Sugared(lggr)

	telem := loop.SetupTelemetry(nil) // default prometheus.Registerer

	promServer := plugins.NewPromServer(envCfg.PrometheusPort(), lggr)
	err = promServer.Start()
	if err != nil {
		lggr.Fatalf("Failed to start prometheus server: %s", err)
	}
	defer slggr.ErrorIfFn(promServer.Close, "error closing prometheus server")

	healthChecker := services.NewChecker()
	err = healthChecker.Start()
	if err != nil {
		lggr.Fatalf("Failed to start health checker: %v", err)
	}
	defer slggr.ErrorIfFn(healthChecker.Close, "Failed to close health checker")

	mp := median.NewPlugin(lggr)
	defer slggr.ErrorIfFn(mp.Close, "error closing pluginMedian")

	err = healthChecker.Register(mp.Name(), mp)
	if err != nil {
		lggr.Fatalf("Failed to register plugin with health checker: %v", err)
	}

	stop := make(chan struct{})
	defer close(stop)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginMedianHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginMedianName: &loop.GRPCPluginMedian{
				PluginServer: mp,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stop,
					Logger:   lggr,
					GRPCOpts: telem,
				},
			},
		},
		GRPCServer: telem.NewServer,
	})
}
