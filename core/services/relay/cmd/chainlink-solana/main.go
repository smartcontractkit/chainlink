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
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgsol "github.com/smartcontractkit/chainlink-solana/pkg/solana"

	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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
	lggr = lggr.Named("PluginSolana")
	defer func() {
		if err := closeLggr(); err != nil {
			fmt.Println("Failed to close logger:", err)
		}
	}()
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
	lggr logger.Logger

	mu      sync.Mutex
	closers []io.Closer
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
	if err := chainSet.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start chain: %w", err)
	}
	r := pkgsol.NewRelayer(c.lggr, chainSet)
	c.mu.Lock()
	c.closers = append(c.closers, chainSet, r)
	c.mu.Unlock()
	return relay.NewLOOPRelayer(r, chainSet, c.lggr), nil
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
