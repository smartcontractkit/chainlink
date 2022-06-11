package vrf

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type VRFContainer struct {
	logger logger.Logger
}

var _ plugins.OraclePlugin = &VRFContainer{}

func NewVRF(l logger.Logger) (*VRFContainer, error) {
	return &VRFContainer{
		logger: l,
	}, nil
}

func (d *VRFContainer) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	return VRFFactory{
		logger: d.logger,
	}, nil

}

func (d *VRFContainer) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
