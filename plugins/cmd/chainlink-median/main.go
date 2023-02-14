package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func main() {
	envCfg, err := plugins.GetEnvConfig()
	if err != nil {
		fmt.Printf("Failed to get environment configuration: %s\n", err)
		os.Exit(1)
	}
	lggr, closeLggr := plugins.NewLogger(envCfg)
	defer closeLggr()

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
