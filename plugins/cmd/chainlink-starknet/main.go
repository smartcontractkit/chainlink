package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/plugins"

	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func main() {
	envCfg, err := plugins.GetEnvConfig()
	if err != nil {
		fmt.Printf("Failed to get environment configuration: %s\n", err)
		os.Exit(1)
	}
	lggr, closeLggr := plugins.NewLogger(envCfg)
	defer closeLggr()

	promServer := plugins.NewPromServer(envCfg.PrometheusPort(), lggr)
	err = promServer.Start()
	if err != nil {
		lggr.Fatalf("Unrecoverable error starting prometheus server: %s", err)
	}
	defer func() {
		err := promServer.Close()
		if err != nil {
			lggr.Errorf("error closing prometheus server: %s", err)
		}
	}()

	cp := &pluginRelayer{lggr: lggr}
	defer func() {
		logger.Sugared(lggr).ErrorIfFn(cp.Close, "pluginRelayer")
	}()

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				StopCh:       stopCh,
				Logger:       lggr,
				PluginServer: cp,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type pluginRelayer struct {
	lggr logger.Logger

	mu      sync.Mutex
	closers []io.Closer
}

// NewRelayer implements the Loopp factory method used by the Loopp server to instantiate a starknet relayer
func (c *pluginRelayer) NewRelayer(ctx context.Context, config string, keystore loop.Keystore) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(config))
	d.DisallowUnknownFields()
	var cfg struct {
		Starknet starknet.StarknetConfigs
	}
	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w", err)
	}

	chainSet, err := starknet.NewChainSet(starknet.ChainSetOpts{
		Logger:   c.lggr,
		KeyStore: keystore,
		Configs:  starknet.NewConfigs(cfg.Starknet),
	}, cfg.Starknet)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}
	ra := relay.NewRelayerAdapter(pkgstarknet.NewRelayer(c.lggr, chainSet), chainSet)

	c.mu.Lock()
	c.closers = append(c.closers, ra)
	c.mu.Unlock()

	return ra, nil
}

func (c *pluginRelayer) Close() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, cl := range c.closers {
		if e := cl.Close(); e != nil {
			err = multierr.Append(err, e)
		}
	}
	return
}
