package promwrapper

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ ocr3_1types.ReportingPluginFactory[any] = &ReportingPluginFactory[any]{}

// ReportingPluginFactory wraps an ocr3_1types.ReportingPluginFactory so the
// produced plugin reports prometheus metrics. It is the OCR3.1 counterpart of
// core/services/ocr3/promwrapper.ReportingPluginFactory.
type ReportingPluginFactory[RI any] struct {
	origin      ocr3_1types.ReportingPluginFactory[RI]
	lggr        logger.Logger
	chainFamily string
	chainID     string
	plugin      string
}

func NewReportingPluginFactory[RI any](
	origin ocr3_1types.ReportingPluginFactory[RI],
	lggr logger.Logger,
	chainFamily string,
	chainID string,
	plugin string,
) *ReportingPluginFactory[RI] {
	return &ReportingPluginFactory[RI]{
		origin:      origin,
		lggr:        lggr,
		chainFamily: chainFamily,
		chainID:     chainID,
		plugin:      plugin,
	}
}

func (r ReportingPluginFactory[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, bbf ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	plugin, info, err := r.origin.NewReportingPlugin(ctx, config, bbf)
	if err != nil {
		return nil, nil, err
	}
	r.lggr.Infow("Wrapping OCR3.1 ReportingPlugin with prometheus metrics reporter",
		"configDigest", config.ConfigDigest,
		"oracleID", config.OracleID,
	)
	wrapped := newReportingPlugin(
		plugin,
		r.chainFamily,
		r.chainID,
		r.plugin,
		config.ConfigDigest.String(),
		promOCR3ReportsGenerated,
		promOCR3Durations,
		promOCR3Sizes,
		promOCR3PluginStatus,
	)
	return wrapped, info, err
}
