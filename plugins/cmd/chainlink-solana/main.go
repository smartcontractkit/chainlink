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
	pkgsol "github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/v2/plugins"

	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
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

	cp := &chainPlugin{lggr: lggr}
	defer func() {
		logger.Sugared(lggr).ErrorIfFn(cp.Close, "chainPlugin")
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: loop.NewGRPCPluginRelayer(cp, lggr),
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type chainPlugin struct {
	lggr       logger.Logger
	promServer *plugins.PromServer

	mu      sync.Mutex
	closers []io.Closer
}

func newChainPlugin(lggr logger.Logger, promPort int) *chainPlugin {
	return &chainPlugin{
		lggr:       lggr,
		promServer: plugins.NewPromServer(promPort, lggr),
		closers:    make([]io.Closer, 0),
	}
}

func (c *chainPlugin) NewRelayer(ctx context.Context, config string, keystore loop.Keystore) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(config))
	d.DisallowUnknownFields()
	var cfg struct {
		Solana solana.SolanaConfigs
	}
	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w", err)
	}

	chainSet, err := solana.NewChainSet(solana.ChainSetOpts{
		Logger:   c.lggr,
		KeyStore: keystore,
		Configs:  solana.NewConfigs(cfg.Solana),
	}, cfg.Solana)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}
	r := pkgsol.NewRelayer(c.lggr, chainSet)

	err = c.promServer.Start()
	if err != nil {
		err = fmt.Errorf("Failed to start prometheus server: %w", err)
		c.lggr.Critical(err.Error())
		return nil, err
	}

	c.mu.Lock()
	c.closers = append(c.closers, chainSet, r, c.promServer)
	c.mu.Unlock()

	return &relay.RelayerAdapter{Relayer: r, RelayerExt: chainSet}, nil
}

func (c *chainPlugin) Close() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, cl := range c.closers {
		if e := cl.Close(); e != nil {
			err = multierr.Append(err, e)
		}
	}
	return
}
