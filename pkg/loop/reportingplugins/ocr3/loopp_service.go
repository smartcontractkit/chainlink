package ocr3

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type LOOPPService struct {
	internal.PluginService[*GRPCService[types.PluginProvider], types.OCR3ReportingPluginFactory]
}

func NewLOOPPService(
	lggr logger.Logger,
	grpcOpts loop.GRPCOpts,
	cmd func() *exec.Cmd,
	config types.ReportingPluginServiceConfig,
	providerConn grpc.ClientConnInterface,
	pipelineRunner types.PipelineRunnerService,
	telemetryService types.TelemetryService,
	errorLog types.ErrorLog,
	capRegistry types.CapabilitiesRegistry,
) *LOOPPService {
	newService := func(ctx context.Context, instance any) (types.OCR3ReportingPluginFactory, error) {
		plug, ok := instance.(types.OCR3ReportingPluginClient)
		if !ok {
			return nil, fmt.Errorf("expected OCR3ReportingPluginClient but got %T", instance)
		}
		return plug.NewReportingPluginFactory(ctx, config, providerConn, pipelineRunner, telemetryService, errorLog, capRegistry)
	}

	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "OCR3GenericService")
	var ps LOOPPService
	broker := net.BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	ps.Init(reportingplugins.PluginServiceName, &GRPCService[types.PluginProvider]{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &ps
}

func (g *LOOPPService) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	if err := g.Wait(); err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, err
	}
	return g.Service.NewReportingPlugin(ctx, config)
}
