package test

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
)

const PluginLoggerTestName = "logger-test"

const LoggerTestName = "server-side-logger-name"

// NOTE: This is part of the test package because it needs to be imported by the test binary at `./internal/test/cmd`
// as well as the test at `./pkg/loop/logger_loop_test.go`
type GRPCPluginLoggerTest struct {
	plugin.NetRPCUnsupportedPlugin

	logger.Logger
}

func (g *GRPCPluginLoggerTest) GRPCServer(*plugin.GRPCBroker, *grpc.Server) (err error) {
	err = errors.New("test error")
	g.Logger.Errorw("Error!", "err", err)
	err = errors.Join(err, g.Logger.Sync())
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
