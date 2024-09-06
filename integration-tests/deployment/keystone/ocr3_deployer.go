package keystone

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

type OCR3Deployer struct {
	lggr     logger.Logger
	contract *ocr3_capability.OCR3Capability
}

func (c *OCR3Deployer) deploy(req deployRequest) (*deployResponse, error) {
	ocr3Addr, tx, ocr3, err := ocr3_capability.DeployOCR3Capability(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}

	_, err = req.Chain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm transaction %s: %w", tx.Hash(), err)
	}
	resp := &deployResponse{
		Address: ocr3Addr,
		Tx:      tx.Hash(),
		Tv: deployment.TypeAndVersion{
			Type:    OCR3Capability,
			Version: deployment.Version1_0_0,
		},
	}
	c.contract = ocr3
	return resp, nil
}
