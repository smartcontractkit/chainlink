package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	loggerName = "PluginStarknet"
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
		lggr.Fatalf("Unrecoverable error starting prometheus server: %s", err)
	}
	defer slggr.ErrorIfFn(promServer.Close, "error closing prometheus server")

	healthChecker := services.NewChecker()
	err = healthChecker.Start()
	if err != nil {
		lggr.Fatalf("Failed to start health checker: %v", err)
	}
	defer slggr.ErrorIfFn(healthChecker.Close, "Failed to close health checker")

	cp := &pluginRelayer{Base: plugins.Base{Logger: lggr}}
	defer slggr.ErrorIfFn(cp.Close, "Failed to close pluginRelayer")

	err = healthChecker.Register(cp.Name(), cp)
	if err != nil {
		lggr.Fatalf("Failed to register plugin with health checker: %v", err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				PluginServer: cp,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stopCh,
					Logger:   lggr,
					GRPCOpts: telem,
				},
			},
		},
		GRPCServer: telem.NewServer,
	})
}

type pluginRelayer struct {
	plugins.Base
}

// NewRelayer implements the Loopp factory method used by the Loopp server to instantiate a starknet relayer
// [github.com/smartcontractkit/chainlink-relay/pkg/loop.PluginRelayer]
// loopKs must be an implementation that can construct a starknet keystore adapter
// [github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm.NewKeystoreAdapter]
func (c *pluginRelayer) NewRelayer(ctx context.Context, config string, loopKs loop.Keystore) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(config))
	d.DisallowUnknownFields()
	var cfg struct {
		Starknet starknet.StarknetConfigs
	}
	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w", err)
	}

	chainSet, err := starknet.NewChainSet(starknet.ChainSetOpts{
		Logger:   c.Logger,
		KeyStore: loopKs,
		Configs:  starknet.NewConfigs(cfg.Starknet),
	}, cfg.Starknet)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}
	ra := relay.NewRelayerAdapter(pkgstarknet.NewRelayer(c.Logger, chainSet), chainSet)

	c.SubService(ra)

	return ra, nil
}
