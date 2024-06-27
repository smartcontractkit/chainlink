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
)

const (
	loggerName = "PluginStarknet"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := &pluginRelayer{Plugin: loop.Plugin{Logger: s.Logger}}
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

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
	loop.Plugin
}

// NewRelayer implements the Loopp factory method used by the Loopp server to instantiate a starknet relayer
// [github.com/smartcontractkit/chainlink-relay/pkg/loop.PluginRelayer]
// loopKs must be an implementation that can construct a starknet keystore adapter
// [github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm.NewKeystoreAdapter]
func (c *pluginRelayer) NewRelayer(ctx context.Context, config string, loopKs loop.Keystore) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(config))
	d.DisallowUnknownFields()
	var cfg struct {
		Starknet starknet.StarknetConfig
	}
	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w:\n\t%s", err, config)
	}

	opts := starknet.ChainOpts{
		Logger:   c.Logger,
		KeyStore: loopKs,
	}

	chain, err := starknet.NewChain(&cfg.Starknet, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}
	ra := &loop.RelayerAdapter{Relayer: pkgstarknet.NewRelayer(c.Logger, chain), RelayerExt: chain}

	c.SubService(ra)

	return ra, nil
}
