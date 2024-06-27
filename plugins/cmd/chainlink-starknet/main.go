package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	loggerName = "PluginStarknet"
)

func main() {
	s := plugins.StartServer(loggerName)
	defer s.Stop()

	p := &pluginRelayer{Base: plugins.Base{Logger: s.Logger}}
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p.Name(), p)

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stopCh,
					Logger:   s.Logger,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
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
