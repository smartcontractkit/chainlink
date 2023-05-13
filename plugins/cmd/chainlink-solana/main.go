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
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.uber.org/multierr"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgsol "github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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
	slggr := logger.Sugared(lggr)

	promServer := plugins.NewPromServer(envCfg.PrometheusPort(), lggr)
	err = promServer.Start()
	if err != nil {
		lggr.Fatalf("Unrecoverable error starting prometheus server: %s", err)
	}
	defer slggr.ErrorIfFn(promServer.Close, "Failed to close prometheus server")

	providers, err := plugins.NewTelemetryProviders("chainlink-solana", envCfg.AppID().String(), lggr)
	if err != nil {
		lggr.Fatalw("Failed to setup telemetry", "err", err)
	}
	defer slggr.ErrorIfFn(providers.Close, "Failed to close telemetry providers")
	otel.SetTracerProvider(providers)
	global.SetMeterProvider(providers)

	cp := &pluginRelayer{lggr: lggr}
	defer slggr.ErrorIfFn(cp.Close, "chainPlugin")

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
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			return grpc.NewServer(append(opts,
				grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
				grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
			)...)
		},
	})
}

type pluginRelayer struct {
	lggr logger.Logger

	mu      sync.Mutex
	closers []io.Closer
}

func (c *pluginRelayer) NewRelayer(ctx context.Context, config string, keystore loop.Keystore) (loop.Relayer, error) {
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
	ra := relay.NewRelayerAdapter(pkgsol.NewRelayer(c.lggr, chainSet), chainSet)

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
