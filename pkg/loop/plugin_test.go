package loop_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils/tests"
)

func testPlugin[I any](t *testing.T, name string, p plugin.Plugin, testFn func(*testing.T, I)) {
	ctx, cancel := context.WithCancel(tests.Context(t))
	defer cancel()

	ch := make(chan *plugin.ReattachConfig, 1)
	closeCh := make(chan struct{})
	go plugin.Serve(&plugin.ServeConfig{
		Test: &plugin.ServeTestConfig{
			Context:          ctx,
			ReattachConfigCh: ch,
			CloseCh:          closeCh,
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins:    map[string]plugin.Plugin{name: p},
	})

	// We should get a config
	var config *plugin.ReattachConfig
	select {
	case config = <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("should've received reattach")
	}
	require.NotNil(t, config)

	c := plugin.NewClient(&plugin.ClientConfig{
		Reattach: config,
		Plugins:  map[string]plugin.Plugin{name: p},
	})
	t.Cleanup(c.Kill)
	clientProtocol, err := c.Client()
	require.NoError(t, err)
	defer clientProtocol.Close()
	i, err := clientProtocol.Dispense(name)
	require.NoError(t, err)

	testFn(t, i.(I))

	// stop plugin
	cancel()
	select {
	case <-closeCh:
	case <-time.After(5 * time.Second):
		t.Fatal("should've stopped")
	}
	require.Error(t, clientProtocol.Ping())
}

func helperProcess(s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	env := []string{
		"GO_WANT_HELPER_PROCESS=1",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}

// This is not a real test. This is just a helper process kicked off by
// tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]

	limit := -1
	if len(args) > 0 {
		var err error
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse integer limit: %s\n", err)
			os.Exit(2)
		}
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
	switch cmd {
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

	case PluginLoggerTestName:
		loggerTest := &GRPCPluginLoggerTest{Logger: logger.Named(lggr, LoggerTestName)}
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: PluginLoggerTestHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				PluginLoggerTestName: loggerTest,
			},
			GRPCServer: grpcServer,
		})

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmd)
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
