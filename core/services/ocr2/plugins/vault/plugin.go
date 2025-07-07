package vault

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
)

func NewReportingPluginFactory(store *requests.Store[*Request]) *ReportingPluginFactory {
	return &ReportingPluginFactory{
		store: store,
	}
}

type ReportingPluginFactory struct {
	store *requests.Store[*Request]
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	return nil, ocr3types.ReportingPluginInfo{}, nil
}
