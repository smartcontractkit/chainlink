package promwrapper

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

var _ ocr3types.ReportingPluginFactory[any] = &ReportingPluginFactory[any]{}

type ReportingPluginFactory[RI any] struct {
	origin  ocr3types.ReportingPluginFactory[RI]
	lggr    logger.Logger
	chainID string
	plugin  string
}

func NewReportingPluginFactory[RI any](
	origin ocr3types.ReportingPluginFactory[RI],
	lggr logger.Logger,
	chainID string,
	plugin string,
) *ReportingPluginFactory[RI] {
	return &ReportingPluginFactory[RI]{
		origin:  origin,
		lggr:    lggr,
		chainID: chainID,
		plugin:  plugin,
	}
}

func (r ReportingPluginFactory[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[RI], ocr3types.ReportingPluginInfo, error) {
	plugin, info, err := r.origin.NewReportingPlugin(ctx, config)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, err
	}
	r.lggr.Infow("Wrapping ReportingPlugin with prometheus metrics reporter",
		"configDigest", config.ConfigDigest,
		"oracleID", config.OracleID,
	)
	wrapped := newReportingPlugin(
		plugin,
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
