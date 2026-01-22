package beholderwrapper

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

var _ ocr3_1types.ReportingPluginFactory[any] = &ReportingPluginFactory[any]{}

type ReportingPluginFactory[RI any] struct {
	wrapped ocr3_1types.ReportingPluginFactory[RI]
	lggr    logger.Logger
	plugin  string
}

func NewReportingPluginFactory[RI any](
	wrapped ocr3_1types.ReportingPluginFactory[RI],
	lggr logger.Logger,
	plugin string,
) *ReportingPluginFactory[RI] {
	return &ReportingPluginFactory[RI]{
		wrapped: wrapped,
		lggr:    lggr,
		plugin:  plugin,
	}
}

func (r ReportingPluginFactory[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	plugin, info, err := r.wrapped.NewReportingPlugin(ctx, config, fetcher)
	if err != nil {
		return nil, nil, err
	}

	metrics, err := newPluginMetrics(r.plugin, config.ConfigDigest.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create plugin metrics: %w", err)
	}

	r.lggr.Infow("Wrapping OCR3_1 ReportingPlugin with beholder metrics reporter",
		"configDigest", config.ConfigDigest,
	)

	wrappedPlugin := newReportingPlugin(plugin, metrics)
	return wrappedPlugin, info, nil
}
