package loop_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
)

const PluginLoggerTestName = "logger-test"

const LoggerTestName = "server-side-logger-name"

func TestHCLogLogger(t *testing.T) {
	lggr, ol := logger.TestObserved(t, zapcore.ErrorLevel)
	loggerTest := &GRPCPluginLoggerTest{Logger: lggr}
	cc := loggerTest.ClientConfig()
	cc.Cmd = helperProcess(PluginLoggerTestName)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	_, err := c.Client()
	require.Error(t, err)

	// Some logs should come through with plugin-side names
	require.NotEmpty(t, ol.Filter(func(entry observer.LoggedEntry) bool {
		return entry.LoggerName == LoggerTestName
	}), ol.All())
}

type GRPCPluginLoggerTest struct {
	plugin.NetRPCUnsupportedPlugin

	logger.Logger
}

func (g *GRPCPluginLoggerTest) GRPCServer(*plugin.GRPCBroker, *grpc.Server) (err error) {
	err = errors.New("test error")
	g.Logger.Errorw("Error!", "err", err)
	g.Logger.Sync()
	time.Sleep(time.Second)
	return err
}

func (g *GRPCPluginLoggerTest) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (g *GRPCPluginLoggerTest) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginLoggerTestHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginLoggerTestName: g},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           loop.HCLogLogger(g.Logger),
	}
}

func PluginLoggerTestHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_LOGGER_TEST_MAGIC_COOKIE",
		MagicCookieValue: "272d1867cdc8042f9405d7c1da3762ec",
	}
}
