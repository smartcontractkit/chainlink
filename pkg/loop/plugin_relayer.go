package loop

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// PluginRelayerName is the name for [types.PluginRelayer]/[NewGRPCPluginRelayer].
const PluginRelayerName = "relayer"

type PluginRelayer = internal.PluginRelayer

func PluginRelayerHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_RELAYER_MAGIC_COOKIE",
		MagicCookieValue: "dae753d4542311b33cf041b930db0150647e806175c2818a0c88a9ab745e45aa",
	}
}

// Deprecated
type Keystore = types.Keystore

type Relayer = internal.Relayer

type BrokerConfig = internal.BrokerConfig

var _ plugin.GRPCPlugin = (*GRPCPluginRelayer)(nil)

// GRPCPluginRelayer implements [plugin.GRPCPlugin] for [types.PluginRelayer].
type GRPCPluginRelayer struct {
	plugin.NetRPCUnsupportedPlugin

	BrokerConfig

	PluginServer PluginRelayer

	pluginClient *internal.PluginRelayerClient
}

func (p *GRPCPluginRelayer) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return internal.RegisterPluginRelayerServer(server, broker, p.BrokerConfig, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.PluginRelayer], updated with the new broker and conn.
func (p *GRPCPluginRelayer) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = internal.NewPluginRelayerClient(broker, p.BrokerConfig, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}
	return PluginRelayer(p.pluginClient), nil
}

func (p *GRPCPluginRelayer) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginRelayerHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginRelayerName: p},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		GRPCDialOptions:  p.DialOpts,
		Logger:           HCLogLogger(p.Logger),
	}
}
