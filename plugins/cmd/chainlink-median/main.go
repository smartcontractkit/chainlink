package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	slggr := logger.Sugared(lggr)

	promServer := plugins.NewPromServer(envCfg.PrometheusPort(), lggr)
	err = promServer.Start()
	if err != nil {
		lggr.Fatalf("Failed to start prometheus server: %s", err)
	}
	defer slggr.ErrorIfFn(promServer.Close, "Failed to close prometheus server")

	providers, err := plugins.NewTelemetryProviders("chainlink-median", envCfg.AppID().String(), lggr)
	if err != nil {
		lggr.Fatalw("Failed to setup telemetry", "err", err)
	}
	defer slggr.ErrorIfFn(providers.Close, "Failed to close telemetry providers")
	otel.SetTracerProvider(providers)
	global.SetMeterProvider(providers)

	mp := median.NewPlugin(lggr)
	err = mp.Start(context.Background())
	if err != nil {
		lggr.Fatalf("Failed to start median plugin: %s", err)
	}

	stop := make(chan struct{})
	defer close(stop)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginMedianHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginMedianName: &loop.GRPCPluginMedian{
				StopCh:       stop,
				Logger:       lggr,
				PluginServer: mp,
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
