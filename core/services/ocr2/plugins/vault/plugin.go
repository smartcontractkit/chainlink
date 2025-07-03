package vault

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

func NewReportingPluginFactory(store *requests.Store[*Request]) *ReportingPluginFactory {
	return &ReportingPluginFactory{
		store: store,
	}
}

type ReportingPluginFactory struct {
	store *requests.Store[*Request]
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	return nil, ocr3_1types.ReportingPluginInfo{}, nil
}
