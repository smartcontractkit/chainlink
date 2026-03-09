package triggerqueue

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	defaultMaxQueryBytes                                   = 100
	defaultMaxObservationBytes                             = 500 * 1024
	defaultMaxReportsPlusPrecursorBytes                    = 500 * 1024
	defaultMaxReportBytes                                  = 500 * 1024
	defaultMaxReportCount                                  = 20
	defaultMaxKeyValueModifiedKeysPlusValuesBytes          = 1024 * 1024
	defaultMaxKeyValueModifiedKeys                         = 500
	defaultMaxBlobPayloadBytes                             = 25 * 1024
	defaultMaxPerOracleUnexpiredBlobCumulativePayloadBytes = 30 * 1024 * 1024
	defaultMaxPerOracleUnexpiredBlobCount                  = 1000
)

// Factory creates OCR 3.1 ReportingPlugins for the trigger queue. Draft: NewReportingPlugin returns a plugin that errors on all calls.
type Factory struct {
	lggr logger.Logger
	services.StateMachine
}

// NewFactory creates a new trigger queue plugin factory.
func NewFactory(lggr logger.Logger) (*Factory, error) {
	if lggr == nil {
		return nil, errors.New("logger is required")
	}
	return &Factory{
		lggr: lggr.Named("TriggerQueuePluginFactory"),
	}, nil
}

// NewReportingPlugin creates a new OCR 3.1 ReportingPlugin. Draft: returns plugin that errors on all OCR calls.
func (f *Factory) NewReportingPlugin(_ context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := NewReportingPlugin(f.lggr)
	_, _ = config, fetcher
	info := ocr3_1types.ReportingPluginInfo1{
		Name: "TriggerQueuePlugin",
		Limits: ocr3_1types.ReportingPluginLimits{
			MaxQueryBytes:                                   defaultMaxQueryBytes,
			MaxObservationBytes:                             defaultMaxObservationBytes,
			MaxReportsPlusPrecursorBytes:                    defaultMaxReportsPlusPrecursorBytes,
			MaxReportBytes:                                  defaultMaxReportBytes,
			MaxReportCount:                                  defaultMaxReportCount,
			MaxKeyValueModifiedKeysPlusValuesBytes:          defaultMaxKeyValueModifiedKeysPlusValuesBytes,
			MaxKeyValueModifiedKeys:                         defaultMaxKeyValueModifiedKeys,
			MaxBlobPayloadBytes:                             defaultMaxBlobPayloadBytes,
			MaxPerOracleUnexpiredBlobCumulativePayloadBytes: defaultMaxPerOracleUnexpiredBlobCumulativePayloadBytes,
			MaxPerOracleUnexpiredBlobCount:                  defaultMaxPerOracleUnexpiredBlobCount,
		},
	}
	return plugin, info, nil
}

func (f *Factory) Start(ctx context.Context) error {
	return f.StartOnce("TriggerQueuePluginFactory", func() error { return nil })
}

func (f *Factory) Close() error {
	return f.StopOnce("TriggerQueuePluginFactory", func() error { return nil })
}

func (f *Factory) Name() string { return f.lggr.Name() }

func (f *Factory) HealthReport() map[string]error {
	return map[string]error{f.Name(): f.Healthy()}
}
