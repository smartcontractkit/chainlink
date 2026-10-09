package beholderwrapper

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	return wrapPlugin(r.lggr, r.plugin, config, plugin, info)
}

var _ ocr3_1types.ReportingPluginFactory2[any] = &ReportingPluginFactory2[any]{}

// ReportingPluginFactory2 is ReportingPluginFactory for ocr3_1types.ReportingPluginFactory2.
type ReportingPluginFactory2[RI any] struct {
	wrapped ocr3_1types.ReportingPluginFactory2[RI]
	lggr    logger.Logger
	plugin  string
}

func NewReportingPluginFactory2[RI any](
	wrapped ocr3_1types.ReportingPluginFactory2[RI],
	lggr logger.Logger,
	plugin string,
) *ReportingPluginFactory2[RI] {
	return &ReportingPluginFactory2[RI]{
		wrapped: wrapped,
		lggr:    lggr,
		plugin:  plugin,
	}
}

func (r ReportingPluginFactory2[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher, kv ocr3_1types.ReadOnlyKeyValueState) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	plugin, info, err := r.wrapped.NewReportingPlugin(ctx, config, fetcher, kv)
	if err != nil {
		return nil, nil, err
	}
	return wrapPlugin(r.lggr, r.plugin, config, plugin, info)
}

func wrapPlugin[RI any](lggr logger.Logger, name string, config ocr3types.ReportingPluginConfig, plugin ocr3_1types.ReportingPlugin[RI], info ocr3_1types.ReportingPluginInfo) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	metrics, err := newPluginMetrics(name, config.ConfigDigest.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create plugin metrics: %w", err)
	}

	lggr.Infow("Wrapping OCR3_1 ReportingPlugin with beholder metrics reporter",
		"configDigest", config.ConfigDigest,
	)

	return newReportingPlugin(plugin, metrics), info, nil
}
