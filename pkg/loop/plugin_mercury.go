package loop

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// PluginMercurynName is the name for [types.PluginMercury]/[NewGRPCPluginMercury].
const PluginMercuryName = "mercury"

func PluginMercuryHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_MERCURY_MAGIC_COOKIE",
		MagicCookieValue: "c3478720-a7a6-4d39-b506-0dd432c9be9a",
	}
}

type GRPCPluginMercury struct {
	plugin.NetRPCUnsupportedPlugin

	BrokerConfig

	PluginServer types.PluginMercury

	pluginClient *mercury.AdapterClient
}

func (p *GRPCPluginMercury) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return mercury.RegisterMercuryAdapterServer(server, broker, p.BrokerConfig, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.PluginMercury], updated with the new broker and conn.
func (p *GRPCPluginMercury) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = mercury.NewMercuryAdapterClient(broker, p.BrokerConfig, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return types.PluginMercury(p.pluginClient), nil
}

func (p *GRPCPluginMercury) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginMercuryHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginMercuryName: p},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		GRPCDialOptions:  p.DialOpts,
		Logger:           HCLogLogger(p.Logger),
	}
}
