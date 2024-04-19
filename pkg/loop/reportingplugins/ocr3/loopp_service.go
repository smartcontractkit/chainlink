package ocr3

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type LOOPPService struct {
	goplugin.PluginService[*GRPCService[types.PluginProvider], core.OCR3ReportingPluginFactory]
}

func NewLOOPPService(
	lggr logger.Logger,
	grpcOpts loop.GRPCOpts,
	cmd func() *exec.Cmd,
	config core.ReportingPluginServiceConfig,
	providerConn grpc.ClientConnInterface,
	pipelineRunner core.PipelineRunnerService,
	telemetryService core.TelemetryService,
	errorLog core.ErrorLog,
	capRegistry core.CapabilitiesRegistry,
	keyValueStore core.KeyValueStore,
) *LOOPPService {
	newService := func(ctx context.Context, instance any) (core.OCR3ReportingPluginFactory, error) {
		plug, ok := instance.(core.OCR3ReportingPluginClient)
		if !ok {
			return nil, fmt.Errorf("expected OCR3ReportingPluginClient but got %T", instance)
		}
		return plug.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetryService, errorLog, capRegistry, keyValueStore)
	}

	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "OCR3GenericService")
	var ps LOOPPService
	broker := net.BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	ps.Init(reportingplugins.PluginServiceName, &GRPCService[types.PluginProvider]{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &ps
}

func (g *LOOPPService) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	if err := g.Wait(); err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, err
	}
	return g.Service.NewReportingPlugin(config)
}

func NewLOOPPServiceValidation(
	lggr logger.Logger,
	grpcOpts loop.GRPCOpts,
	cmd func() *exec.Cmd,
) *reportingplugins.LOOPPServiceValidation {
	newService := func(ctx context.Context, instance any) (core.ValidationService, error) {
		plug, ok := instance.(core.OCR3ReportingPluginClient)
		if !ok {
			return nil, fmt.Errorf("expected ValidationServiceClient but got %T", instance)
		}
		return plug.NewValidationService(ctx)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "GenericService")
	var ps reportingplugins.LOOPPServiceValidation
	broker := net.BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	ps.Init(PluginServiceName, &reportingplugins.GRPCService[types.PluginProvider]{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &ps
}
