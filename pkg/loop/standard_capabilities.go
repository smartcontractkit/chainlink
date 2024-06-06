package loop

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/capability"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const PluginStandardCapabilitiesName = "standardcapabilities"

func StandardCapabilitiesHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_STANDARD_CAPABILITIES_MAGIC_COOKIE",
		MagicCookieValue: "f4df86d3-3552-4231-8206-be0d245b6c67",
	}
}

type StandardCapabilitiesLoop struct {
	Logger logger.Logger
	plugin.NetRPCUnsupportedPlugin
	BrokerConfig
	PluginServer StandardCapabilities
	pluginClient *capability.StandardCapabilitiesClient
}

func (p *StandardCapabilitiesLoop) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return capability.RegisterStandardCapabilitiesServer(server, broker, p.BrokerConfig, p.PluginServer)
}

func (p *StandardCapabilitiesLoop) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	bext := &net.BrokerExt{
		BrokerConfig: p.BrokerConfig,
		Broker:       broker,
	}

	if p.pluginClient == nil {
		p.pluginClient = capability.NewStandardCapabilitiesClient(bext, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return StandardCapabilities(p.pluginClient), nil
}

func (p *StandardCapabilitiesLoop) ClientConfig() *plugin.ClientConfig {
	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: StandardCapabilitiesHandshakeConfig(),
		Plugins:         map[string]plugin.Plugin{PluginStandardCapabilitiesName: p},
	}
	return ManagedGRPCClientConfig(clientConfig, p.BrokerConfig)
}

type StandardCapabilities interface {
	services.Service
	Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore,
		capabilityRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
		pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error
	Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error)
}

type StandardCapabilitiesService struct {
	goplugin.PluginService[*StandardCapabilitiesLoop, StandardCapabilities]
}

func NewStandardCapabilitiesService(lggr logger.Logger, grpcOpts GRPCOpts, cmd func() *exec.Cmd) *StandardCapabilitiesService {
	newService := func(ctx context.Context, instance any) (StandardCapabilities, error) {
		scs, ok := instance.(StandardCapabilities)
		if !ok {
			return nil, fmt.Errorf("expected StandardCapabilities but got %T", instance)
		}
		return scs, nil
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "StandardCapabilities")
	var rs StandardCapabilitiesService
	broker := BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	rs.Init(PluginStandardCapabilitiesName, &StandardCapabilitiesLoop{Logger: lggr, BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &rs
}
