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

const PluginStandardCapabilityName = "standardcapability"

func StandardCapabilityHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_STANDARD_CAPABILITY_MAGIC_COOKIE",
		MagicCookieValue: "f4df86d3-3552-4231-8206-be0d245b6c67",
	}
}

type StandardCapabilityLoop struct {
	Logger logger.Logger
	plugin.NetRPCUnsupportedPlugin
	BrokerConfig
	PluginServer StandardCapability
	pluginClient *capability.StandardCapabilityClient
}

func (p *StandardCapabilityLoop) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return capability.RegisterStandardCapabilityServer(server, broker, p.BrokerConfig, p.PluginServer)
}

func (p *StandardCapabilityLoop) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	bext := &net.BrokerExt{
		BrokerConfig: p.BrokerConfig,
		Broker:       broker,
	}

	if p.pluginClient == nil {
		p.pluginClient = capability.NewStandardCapabilityClient(bext, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return StandardCapability(p.pluginClient), nil
}

func (p *StandardCapabilityLoop) ClientConfig() *plugin.ClientConfig {
	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: StandardCapabilityHandshakeConfig(),
		Plugins:         map[string]plugin.Plugin{PluginStandardCapabilityName: p},
	}
	return ManagedGRPCClientConfig(clientConfig, p.BrokerConfig)
}

type StandardCapability interface {
	services.Service
	capabilities.BaseCapability
	Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore,
		capabilityRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
		pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error
}

type StandardCapabilityService struct {
	goplugin.PluginService[*StandardCapabilityLoop, StandardCapability]
}

func NewStandardCapability(lggr logger.Logger, grpcOpts GRPCOpts, cmd func() *exec.Cmd) *StandardCapabilityService {
	newService := func(ctx context.Context, instance any) (StandardCapability, error) {
		scs, ok := instance.(StandardCapability)
		if !ok {
			return nil, fmt.Errorf("expected StandardCapability but got %T", instance)
		}
		return scs, nil
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "StandardCapability")
	var rs StandardCapabilityService
	broker := BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	rs.Init(PluginStandardCapabilityName, &StandardCapabilityLoop{Logger: lggr, BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &rs
}
