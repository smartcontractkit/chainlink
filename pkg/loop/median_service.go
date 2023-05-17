package loop

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

var _ ocrtypes.ReportingPluginFactory = (*MedianService)(nil)

// MedianService is a [types.Service] that maintains an internal [PluginMedian].
type MedianService struct {
	*pluginService[*GRPCPluginMedian, ReportingPluginFactory]
}

// NewMedianService returns a new [*MedianService].
// cmd must return a new exec.Cmd each time it is called.
func NewMedianService(lggr logger.Logger, cmd func() *exec.Cmd, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog ErrorLog) *MedianService {
	newService := func(ctx context.Context, instance any) (ReportingPluginFactory, error) {
		plug, ok := instance.(PluginMedian)
		if !ok {
			return nil, fmt.Errorf("expected PluginMedian but got %T", instance)
		}
		return plug.NewMedianFactory(ctx, provider, dataSource, juelsPerFeeCoin, errorLog)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "MedianService")
	return &MedianService{newPluginService(PluginMedianName, &GRPCPluginMedian{StopCh: stopCh, Logger: lggr}, newService, lggr, cmd, stopCh)}
}

func (m *MedianService) NewReportingPlugin(config ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	ctx, cancel := utils.ContextFromChan(m.pluginService.stopCh)
	defer cancel()
	if err := m.wait(ctx); err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}
	return m.service.NewReportingPlugin(config)
}
