package promwrapper

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

var _ ocr3_1types.ReportingPluginFactory[any] = &ReportingPluginFactory[any]{}

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

func (r ReportingPluginFactory[RI]) NewReportingPlugin(
	ctx context.Context,
	config ocr3types.ReportingPluginConfig,
	bbf ocr3_1types.BlobBroadcastFetcher,
) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	plugin, info, err := r.origin.NewReportingPlugin(ctx, config, bbf)
	if err != nil {
		// chainlink's libocr has ReportingPluginInfo as an interface;
		// ReportingPluginInfo1 is the concrete "v1" impl (see
		// libocr/offchainreporting2plus/ocr3_1types/plugin.go:457-477).
		// A nil interface return would trip nilness checks downstream —
		// return an empty ReportingPluginInfo1 instead.
		return nil, ocr3_1types.ReportingPluginInfo1{}, err
	}
	r.lggr.Infow("Wrapping OCR3_1 ReportingPlugin with prometheus metrics reporter",
		"configDigest", config.ConfigDigest,
		"oracleID", config.OracleID,
	)
	wrapped := newReportingPlugin(
		plugin,
		r.chainFamily,
		r.chainID,
		r.plugin,
		config.ConfigDigest.String(),
		promOCR3_1ReportsGenerated,
		promOCR3_1Durations,
		promOCR3_1Sizes,
		promOCR3_1PluginStatus,
	)
	return wrapped, info, err
}
