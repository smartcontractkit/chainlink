package ocr2vrf

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type OCR2VRFContainer struct {
	logger logger.Logger
}

var _ plugins.OraclePlugin = &OCR2VRFContainer{}

func NewOCR2VRF(l logger.Logger) (*OCR2VRFContainer, error) {
	return &OCR2VRFContainer{
		logger: l,
	}, nil
}

func (d *OCR2VRFContainer) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	return OCR2VRFFactory{
		logger: d.logger,
	}, nil

}

func (d *OCR2VRFContainer) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
