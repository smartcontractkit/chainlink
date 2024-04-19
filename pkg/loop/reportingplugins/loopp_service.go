package reportingplugins

import (
	"context"
	"fmt"
	"os/exec"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ ocrtypes.ReportingPluginFactory = (*LOOPPService)(nil)

// LOOPPService is a [types.Service] that maintains an internal [types.PluginClient].
type LOOPPService struct {
	goplugin.PluginService[*GRPCService[types.PluginProvider], types.ReportingPluginFactory]
}

type LOOPPServiceValidation struct {
	goplugin.PluginService[*GRPCService[types.PluginProvider], core.ValidationService]
}

// NewLOOPPService returns a new [*PluginService].
// cmd must return a new exec.Cmd each time it is called.
// We use a `conn` here rather than a provider so that we can enforce proxy providers being passed in.
func NewLOOPPService(
	lggr logger.Logger,
	grpcOpts loop.GRPCOpts,
	cmd func() *exec.Cmd,
	config core.ReportingPluginServiceConfig,
	providerConn grpc.ClientConnInterface,
	pipelineRunner core.PipelineRunnerService,
	telemetryService core.TelemetryService,
	errorLog core.ErrorLog,
	keyValueStore core.KeyValueStore,
) *LOOPPService {
	newService := func(ctx context.Context, instance any) (types.ReportingPluginFactory, error) {
		plug, ok := instance.(core.ReportingPluginClient)
		if !ok {
			return nil, fmt.Errorf("expected GenericPluginClient but got %T", instance)
		}
		return plug.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetryService, errorLog, keyValueStore)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "GenericService")
	var ps LOOPPService
	broker := net.BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	ps.Init(PluginServiceName, &GRPCService[types.PluginProvider]{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &ps
}

func (g *LOOPPService) NewReportingPlugin(config ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	if err := g.Wait(); err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}
	return g.Service.NewReportingPlugin(config)
}

func NewLOOPPServiceValidation(
	lggr logger.Logger,
	grpcOpts loop.GRPCOpts,
	cmd func() *exec.Cmd,
) *LOOPPServiceValidation {
	newService := func(ctx context.Context, instance any) (core.ValidationService, error) {
		plug, ok := instance.(core.ReportingPluginClient)
		if !ok {
			return nil, fmt.Errorf("expected ValidationServiceClient but got %T", instance)
		}
		return plug.NewValidationService(ctx)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "GenericService")
	var ps LOOPPServiceValidation
	broker := net.BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	ps.Init(PluginServiceName, &GRPCService[types.PluginProvider]{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &ps
}

func (g *LOOPPServiceValidation) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
	if err := g.Wait(); err != nil {
		return err
	}

	return g.Service.ValidateConfig(ctx, config)
}
