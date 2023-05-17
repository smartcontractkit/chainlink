package loop

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
)

// PluginMedianName is the name for [PluginMedian]/[NewGRPCPluginMedian].
const PluginMedianName = "median"

type PluginMedian = internal.PluginMedian

type ErrorLog = internal.ErrorLog

func PluginMedianHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_MEDIAN_MAGIC_COOKIE",
		MagicCookieValue: "b12a697e19748cd695dd1690c09745ee7cc03717179958e8eadd5a7ca4646728",
	}
}

type ReportingPluginFactory = internal.ReportingPluginFactory

type GRPCPluginMedian struct {
	plugin.NetRPCUnsupportedPlugin

	StopCh <-chan struct{}
	Logger logger.Logger

	PluginServer PluginMedian

	pluginClient *internal.PluginMedianClient
}

func (p *GRPCPluginMedian) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return internal.RegisterPluginMedianServer(server, p.StopCh, p.Logger, broker, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [PluginMedian], updated with the new broker and conn.
func (p *GRPCPluginMedian) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = internal.NewPluginMedianClient(p.StopCh, p.Logger, broker, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return PluginMedian(p.pluginClient), nil
}

func (p *GRPCPluginMedian) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginMedianHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginMedianName: p},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}
