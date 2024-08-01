package reportingplugins

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/reportingplugin/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ocr3capability"
	pluginprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
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
	core.ReportingPluginServer[T]
	ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerConfig loop.BrokerConfig) T
}

// GRPCService is the loopp interface for a plugin that can
// run an arbitrary product in the core node. By specifying
// `T`, server instances can request a specific provider type.
type GRPCService[T types.PluginProvider] struct {
	plugin.NetRPCUnsupportedPlugin

	loop.BrokerConfig

	PluginServer ProviderServer[T]

	pluginClient *ocr2.ReportingPluginServiceClient
}

type serverAdapter struct {
	NewReportingPluginFactoryFn func(
		ctx context.Context,
		config core.ReportingPluginServiceConfig,
		conn grpc.ClientConnInterface,
		pr core.PipelineRunnerService,
		ts core.TelemetryService,
		errorLog core.ErrorLog,
		kv core.KeyValueStore,
		rs core.RelayerSet,
	) (types.ReportingPluginFactory, error)

	ValidateConfigService
}

type ValidateConfigService interface {
	NewValidationService(ctx context.Context) (core.ValidationService, error)
}

func (s serverAdapter) NewValidationService(ctx context.Context) (core.ValidationService, error) {
	return s.ValidateConfigService.NewValidationService(ctx)
}
func (s serverAdapter) NewReportingPluginFactory(
	ctx context.Context,
	config core.ReportingPluginServiceConfig,
	conn grpc.ClientConnInterface,
	pr core.PipelineRunnerService,
	ts core.TelemetryService,
	errorLog core.ErrorLog,
	kv core.KeyValueStore,
	rs core.RelayerSet,
) (types.ReportingPluginFactory, error) {
	return s.NewReportingPluginFactoryFn(ctx, config, conn, pr, ts, errorLog, kv, rs)
}

func (g *GRPCService[T]) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	newReportingPluginFactoryFn := func(
		ctx context.Context,
		cfg core.ReportingPluginServiceConfig,
		conn grpc.ClientConnInterface,
		pr core.PipelineRunnerService,
		ts core.TelemetryService,
		el core.ErrorLog,
		kv core.KeyValueStore,
		rs core.RelayerSet,
	) (types.ReportingPluginFactory, error) {
		provider := g.PluginServer.ConnToProvider(conn, broker, g.BrokerConfig)
		tc := telemetry.NewTelemetryClient(ts)
		return g.PluginServer.NewReportingPluginFactory(ctx, cfg, provider, pr, tc, el, kv, rs)
	}

	return ocr2.RegisterReportingPluginServiceServer(server, broker, g.BrokerConfig, serverAdapter{
		NewReportingPluginFactoryFn: newReportingPluginFactoryFn,
		ValidateConfigService:       g.PluginServer,
	})
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.PluginClient], updated with the new broker and conn.
func (g *GRPCService[T]) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if g.pluginClient == nil {
		g.pluginClient = ocr2.NewReportingPluginServiceClient(broker, g.BrokerConfig, conn)
	} else {
		g.pluginClient.Refresh(broker, conn)
	}

	return core.ReportingPluginClient(g.pluginClient), nil
}

func (g *GRPCService[T]) ClientConfig() *plugin.ClientConfig {
	c := &plugin.ClientConfig{
		HandshakeConfig: ReportingPluginHandshakeConfig(),
		Plugins:         map[string]plugin.Plugin{PluginServiceName: g},
	}
	return loop.ManagedGRPCClientConfig(c, g.BrokerConfig)
}

// These implement `ConnToProvider` and return the conn wrapped as
// the specified provider type. They can be embedded into the server struct
// for ease of use.
type PluginProviderServer = pluginprovider.PluginProviderServer
type MedianProviderServer = median.ProviderServer
type OCR3CapabilityProviderServer = ocr3capability.ProviderServer
