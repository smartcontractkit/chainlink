package test

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

const PluginLoggerTestName = "logger-test"

const (
	PANIC = iota
	FATAL
	CRITICAL
	ERROR
	INFO
	WARN
	DEBUG
)

// NOTE: This is part of the test package because it needs to be imported by the test binary at `./internal/test/cmd`
// as well as the test at `./pkg/loop/logger_loop_test.go`
type GRPCPluginLoggerTest struct {
	plugin.NetRPCUnsupportedPlugin
	logger.SugaredLogger
	ErrorType int
}

func (g *GRPCPluginLoggerTest) GRPCServer(*plugin.GRPCBroker, *grpc.Server) (err error) {
	//Simulate panic/error/log after GRPC is started, if a panic is thrown before the GRPC server is initialized
	//it will not be caught as stderr will be closed before HashiCorp plugin will have a change to read from it
	go func() {
		time.Sleep(time.Second)
		switch g.ErrorType {
		case PANIC:
			panic("random panic")
		case FATAL:
			g.Fatalw("some panic log", "custom-name-panic", "custom-value-panic")
		case CRITICAL:
			g.Criticalw("some critical error log", "custom-name-critical", "custom-value-critical")
		case ERROR:
			g.Errorw("some error log", "custom-name-error", "custom-value-error")
		case INFO:
			g.Infow("some info log", "custom-name-info", "custom-value-info")
		case WARN:
			g.Warnw("some warn log", "custom-name-warn", "custom-value-warn")
		case DEBUG:
			g.Debugw("some debug log", "custom-name-debug", "custom-value-debug")
		}
	}()
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
		Logger:           loop.HCLogLogger(g.SugaredLogger),
	}
}

func PluginLoggerTestHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_LOGGER_TEST_MAGIC_COOKIE",
		MagicCookieValue: "272d1867cdc8042f9405d7c1da3762ec",
	}
}
