package dkg

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type DKGContainer struct {
	logger logger.Logger
}

var _ plugins.OraclePlugin = &DKGContainer{}

func NewDKG(l logger.Logger) (*DKGContainer, error) {
	return &DKGContainer{
		logger: l,
	}, nil
}

func (d *DKGContainer) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	return DKGFactory{
		logger: d.logger,
	}, nil

}

func (d *DKGContainer) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
