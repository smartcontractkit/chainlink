package reportingplugins

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// PluginServiceName is the name for [types.PluginClient]/[NewGRPCService].
const PluginServiceName = "plugin-service"

func ReportingPluginHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_GENERIC_MAGIC_COOKIE",
		MagicCookieValue: "2ad981747cd86c4ab3e23170970020fd",
	}
}

type ProviderServer[T types.PluginProvider] interface {
	types.ReportingPluginServer[T]
	ConnToProvider(conn grpc.ClientConnInterface, broker internal.Broker, brokerConfig loop.BrokerConfig) T
}

// GRPCService is the loopp interface for a plugin that can
// run an arbitrary product in the core node. By specifying
// `T`, server instances can request a specific provider type.
type GRPCService[T types.PluginProvider] struct {
	plugin.NetRPCUnsupportedPlugin

	loop.BrokerConfig

	PluginServer ProviderServer[T]

	pluginClient *internal.ReportingPluginServiceClient
}

type serverAdapter func(
	context.Context,
	types.ReportingPluginServiceConfig,
	grpc.ClientConnInterface,
	types.PipelineRunnerService,
	types.TelemetryClient,
	types.ErrorLog,
) (types.ReportingPluginFactory, error)

func (s serverAdapter) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	conn grpc.ClientConnInterface,
	pr types.PipelineRunnerService,
	ts types.TelemetryClient,
	errorLog types.ErrorLog,
) (types.ReportingPluginFactory, error) {
	return s(ctx, config, conn, pr, ts, errorLog)
}

func (g *GRPCService[T]) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	adapter := func(
		ctx context.Context,
		cfg types.ReportingPluginServiceConfig,
		conn grpc.ClientConnInterface,
		pr types.PipelineRunnerService,
		ts types.TelemetryClient,
		el types.ErrorLog,
	) (types.ReportingPluginFactory, error) {
		provider := g.PluginServer.ConnToProvider(conn, broker, g.BrokerConfig)
		return g.PluginServer.NewReportingPluginFactory(ctx, cfg, provider, pr, ts, el)
	}
	return internal.RegisterReportingPluginServiceServer(server, broker, g.BrokerConfig, serverAdapter(adapter))
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.PluginClient], updated with the new broker and conn.
func (g *GRPCService[T]) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if g.pluginClient == nil {
		g.pluginClient = internal.NewReportingPluginServiceClient(broker, g.BrokerConfig, conn)
	} else {
		g.pluginClient.Refresh(broker, conn)
	}

	return types.ReportingPluginClient(g.pluginClient), nil
}

func (g *GRPCService[T]) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  ReportingPluginHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginServiceName: g},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		GRPCDialOptions:  g.BrokerConfig.DialOpts,
		Logger:           loop.HCLogLogger(g.BrokerConfig.Logger),
	}
}

// These implement `ConnToProvider` and return the conn wrapped as
// the specified provider type. They can be embedded into the server struct
// for ease of use.
type PluginProviderServer = internal.PluginProviderServer
type MedianProviderServer = internal.MedianProviderServer
