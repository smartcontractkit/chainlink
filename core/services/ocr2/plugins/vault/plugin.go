package vault

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

func NewReportingPluginFactory() *ReportingPluginFactory {
	return &ReportingPluginFactory{}
}

type ReportingPluginFactory struct{}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	return nil, ocr3types.ReportingPluginInfo{}, nil
}
