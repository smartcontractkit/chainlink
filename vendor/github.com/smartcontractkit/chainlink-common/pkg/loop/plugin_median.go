package loop

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/median"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

// PluginMedianName is the name for [types.PluginMedian]/[NewGRPCPluginMedian].
const PluginMedianName = "median"

// Deprecated
type PluginMedian = core.PluginMedian

// Deprecated
type ErrorLog = core.ErrorLog

func PluginMedianHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_MEDIAN_MAGIC_COOKIE",
		MagicCookieValue: "b12a697e19748cd695dd1690c09745ee7cc03717179958e8eadd5a7ca4646728",
	}
}

// Deprecated
type ReportingPluginFactory = types.ReportingPluginFactory

type GRPCPluginMedian struct {
	plugin.NetRPCUnsupportedPlugin

	BrokerConfig

	PluginServer core.PluginMedian

	pluginClient *median.PluginMedianClient
}

func (p *GRPCPluginMedian) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return median.RegisterPluginMedianServer(server, broker, p.BrokerConfig, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.PluginMedian], updated with the new broker and conn.
func (p *GRPCPluginMedian) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = median.NewPluginMedianClient(broker, p.BrokerConfig, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return core.PluginMedian(p.pluginClient), nil
}

func (p *GRPCPluginMedian) ClientConfig() *plugin.ClientConfig {
	c := &plugin.ClientConfig{
		HandshakeConfig: PluginMedianHandshakeConfig(),
		Plugins:         map[string]plugin.Plugin{PluginMedianName: p},
	}
	return ManagedGRPCClientConfig(c, p.BrokerConfig)
}
