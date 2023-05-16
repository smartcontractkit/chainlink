package loop

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ PluginMedian = (*MedianService)(nil)

// MedianService is a [types.Service] that maintains an internal [PluginMedian].
type MedianService struct {
	*pluginService[*GRPCPluginMedian, PluginMedian]
}

// NewMedianService returns a new [*MedianService].
// cmd must return a new exec.Cmd each time it is called.
func NewMedianService(lggr logger.Logger, cmd func() *exec.Cmd) *MedianService {
	newService := func(ctx context.Context, instance any) (PluginMedian, error) {
		plug, ok := instance.(PluginMedian)
		if !ok {
			return nil, fmt.Errorf("expected PluginMedian but got %T", instance)
		}
		return plug, nil
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "MedianService")
	return &MedianService{newPluginService(PluginMedianName, &GRPCPluginMedian{StopCh: stopCh, Logger: lggr}, newService, lggr, cmd, stopCh)}
}

func (m *MedianService) NewMedianFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog ErrorLog) (ocrtypes.ReportingPluginFactory, error) {
	if err := m.wait(ctx); err != nil {
		return nil, err
	}
	return m.service.NewMedianFactory(ctx, provider, dataSource, juelsPerFeeCoin, errorLog)
}
