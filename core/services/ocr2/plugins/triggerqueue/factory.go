package triggerqueue

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = (*Factory)(nil)

// Factory creates ReportingPlugin instances for the trigger queue.
type Factory struct {
	lggr   logger.Logger
	buffer *v2.ObservationBuffer[v2.EnqueuedTriggerEvent]
}

// NewFactory creates a new trigger queue plugin factory.
func NewFactory(lggr logger.Logger, buffer *v2.ObservationBuffer[v2.EnqueuedTriggerEvent]) *Factory {
	return &Factory{lggr: lggr.Named("TriggerQueueFactory"), buffer: buffer}
}

func (f *Factory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := NewReportingPlugin(f.lggr, f.buffer)
	info := ocr3_1types.ReportingPluginInfo1{
		Name: "triggerqueue",
		Limits: ocr3_1types.ReportingPluginLimits{
			MaxQueryBytes:                                   100,
			MaxObservationBytes:                             500 * 1024, // 500KB per design doc
			MaxReportsPlusPrecursorBytes:                    500 * 1024,
			MaxReportBytes:                                  500 * 1024,
			MaxReportCount:                                  1,
			MaxKeyValueModifiedKeys:                         500,
			MaxKeyValueModifiedKeysPlusValuesBytes:          1024 * 1024, // 1MB
			MaxBlobPayloadBytes:                             25 * 1024,   // 25KB per design doc
			MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 30 * 1024 * 1024,
			MaxPerOracleUnexpiredBlobCount:                  1000,
		},
	}
	return plugin, info, nil
}
