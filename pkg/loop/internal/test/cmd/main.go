package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func main() {
	cmdS := ""
	flag.StringVar(&cmdS, "cmd", "", "the name of the service to run")

	limitI := 0
	flag.IntVar(&limitI, "limit", 0, "the number of services to run")

	flag.Parse()
	defer os.Exit(0)

	if cmdS == "" {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	limit := -1
	if limitI != 0 {
		limit = limitI
	}

	grpcServer := func(opts []grpc.ServerOption) *grpc.Server { return grpc.NewServer(opts...) }
	if limit > -1 {
		unary, stream := limitInterceptors(limit)
		grpcServer = func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.UnaryInterceptor(unary), grpc.StreamInterceptor(stream))
			return grpc.NewServer(opts...)
		}
	}

	lggr, err := loop.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %s\n", err)
		os.Exit(2)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	switch cmdS {
	case loop.PluginRelayerName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				loop.PluginRelayerName: &loop.GRPCPluginRelayer{PluginServer: test.StaticPluginRelayer{}, BrokerConfig: loop.BrokerConfig{Logger: lggr, StopCh: stopCh}},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	case loop.PluginMedianName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: loop.PluginMedianHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				loop.PluginMedianName: &loop.GRPCPluginMedian{PluginServer: test.StaticPluginMedian{}, BrokerConfig: loop.BrokerConfig{Logger: lggr, StopCh: stopCh}},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	case test.PluginLoggerTestName:
		loggerTest := &test.GRPCPluginLoggerTest{Logger: logger.Named(lggr, test.LoggerTestName)}
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: test.PluginLoggerTestHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				test.PluginLoggerTestName: loggerTest,
			},
			GRPCServer: grpcServer,
		})

	case reportingplugins.PluginServiceName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: reportingplugins.ReportingPluginHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				reportingplugins.PluginServiceName: &reportingplugins.GRPCService[types.PluginProvider]{
					PluginServer: test.StaticReportingPluginWithPluginProvider{},
					BrokerConfig: loop.BrokerConfig{
						Logger: lggr,
						StopCh: stopCh,
					},
				},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	case test.ReportingPluginWithMedianProviderName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: reportingplugins.ReportingPluginHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				reportingplugins.PluginServiceName: &reportingplugins.GRPCService[types.MedianProvider]{
					PluginServer: test.StaticReportingPluginWithMedianProvider{},
					BrokerConfig: loop.BrokerConfig{
						Logger: lggr,
						StopCh: stopCh,
					},
				},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmdS)
		os.Exit(2)
	}
}

// limitInterceptors returns a pair of interceptors which increment a shared count for each call and exit the program
// when limit is reached.
func limitInterceptors(limit int) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	count := make(chan struct{})
	go func() {
		for i := 0; i < limit; i++ {
			<-count
		}
		os.Exit(3)
	}()
	return limitUnaryInterceptor(count), limitStreamInterceptor(count)
}

func limitUnaryInterceptor(count chan<- struct{}) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		count <- struct{}{}
		return handler(ctx, req)
	}
}

func limitStreamInterceptor(count chan<- struct{}) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		count <- struct{}{}
		return handler(srv, ss)
	}
}
