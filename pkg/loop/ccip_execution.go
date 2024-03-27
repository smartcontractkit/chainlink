package loop

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// CCIPExecutionLOOPName is the name for [types.CCIPExecutionFactoryGenerator]/[NewExecutionLOOP].
const CCIPExecutionLOOPName = "ccip_execution"

func PluginCCIPExecutionHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_CCIP_EXEC_MAGIC_COOKIE",
		MagicCookieValue: "5a2d1527-6c0f-4c7e-8c96-00aa4bececd2",
	}
}

type ExecutionLoop struct {
	plugin.NetRPCUnsupportedPlugin

	BrokerConfig

	PluginServer types.CCIPExecutionFactoryGenerator

	pluginClient *ccip.ExecutionLOOPClient
}

func (p *ExecutionLoop) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return ccip.RegisterExecutionLOOPServer(server, broker, p.BrokerConfig, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.CCIPExecutionFactoryGenerator], updated with the new broker and conn.
func (p *ExecutionLoop) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = ccip.NewExecutionLOOPClient(broker, p.BrokerConfig, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return types.CCIPExecutionFactoryGenerator(p.pluginClient), nil
}

func (p *ExecutionLoop) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginCCIPExecutionHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{CCIPExecutionLOOPName: p},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		GRPCDialOptions:  p.DialOpts,
		Logger:           HCLogLogger(p.Logger),
	}
}

var _ ocrtypes.ReportingPluginFactory = (*ExecutionFactoryService)(nil)

// ExecutionFactoryService is a [types.Service] that maintains an internal [types.CCIPExecutionFactoryGenerator].
type ExecutionFactoryService struct {
	goplugin.PluginService[*ExecutionLoop, types.ReportingPluginFactory]
}

// NewExecutionService returns a new [*ExecutionFactoryService].
// cmd must return a new exec.Cmd each time it is called.
func NewExecutionService(lggr logger.Logger, grpcOpts GRPCOpts, cmd func() *exec.Cmd, provider types.CCIPExecProvider) *ExecutionFactoryService {
	newService := func(ctx context.Context, instance any) (types.ReportingPluginFactory, error) {
		plug, ok := instance.(types.CCIPExecutionFactoryGenerator)
		if !ok {
			return nil, fmt.Errorf("expected CCIPExecutionFactoryGenerator but got %T", instance)
		}
		return plug.NewExecutionFactory(ctx, provider)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "CCIPExecutionService")
	var efs ExecutionFactoryService
	broker := BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	efs.Init(CCIPExecutionLOOPName, &ExecutionLoop{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &efs
}

func (m *ExecutionFactoryService) NewReportingPlugin(config ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	if err := m.Wait(); err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}
	return m.Service.NewReportingPlugin(config)
}
