package beholderwrapper

import (
	"context"
	"fmt"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/beholderwrapper/metrics"
)

var _ ocr3types.ReportingPluginFactory[any] = &ReportingPluginFactory[any]{}

type ReportingPluginFactory[RI any] struct {
	wrapped ocr3types.ReportingPluginFactory[RI]
	lggr    logger.Logger
	plugin  string
}

func NewReportingPluginFactory[RI any](
	wrapped ocr3types.ReportingPluginFactory[RI],
	lggr logger.Logger,
	plugin string,
) *ReportingPluginFactory[RI] {
	return &ReportingPluginFactory[RI]{
		wrapped: wrapped,
		lggr:    lggr,
		plugin:  plugin,
	}
}

func (r ReportingPluginFactory[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[RI], ocr3types.ReportingPluginInfo, error) {
	plugin, info, err := r.wrapped.NewReportingPlugin(ctx, config)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, err
	}

	m, err := metrics.NewPluginMetrics(MetricPrefix, r.plugin, config.ConfigDigest.String())
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("failed to create plugin metrics: %w", err)
	}

	r.lggr.Infow("Wrapping OCR3 ReportingPlugin with beholder metrics reporter",
		"configDigest", config.ConfigDigest,
		"oracleID", config.OracleID,
	)

	wrappedPlugin := newReportingPlugin(plugin, m)
	return wrappedPlugin, info, nil
}

// MetricViews returns the histogram bucket views for registration with beholder
func MetricViews() []sdkmetric.View {
	return metrics.MetricViews(MetricPrefix)
}
