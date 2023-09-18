package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgsol "github.com/smartcontractkit/chainlink-solana/pkg/solana"

	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	loggerName = "PluginSolana"
)

func main() {
	s := plugins.StartServer(loggerName)
	defer s.Stop()

	p := &pluginRelayer{Base: plugins.Base{Logger: s.Logger}}
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
	plugins.Base
}

func (c *pluginRelayer) NewRelayer(ctx context.Context, config string, keystore loop.Keystore) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(config))
	d.DisallowUnknownFields()
	var cfg struct {
		Solana solana.SolanaConfig
	}

	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w:\n\t%s", err, config)
	}

	opts := solana.ChainOpts{
		Logger:   c.Logger,
		KeyStore: keystore,
	}
	chain, err := solana.NewChain(&cfg.Solana, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}
	ra := relay.NewRelayerAdapter(pkgsol.NewRelayer(c.Logger, chain), chain)

	c.SubService(ra)

	return ra, nil
}
