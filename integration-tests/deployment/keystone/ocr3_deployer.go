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

var OCR3CapabilityTypeVersion = deployment.TypeAndVersion{
	Type:    OCR3Capability,
	Version: deployment.Version1_0_0,
}

func (c *OCR3Deployer) deploy(req deployRequest) (*deployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, ocr3_capability.OCR3CapabilityABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Infof("ocr3 capability estimated gas: %d", est)

	ocr3Addr, tx, ocr3, err := ocr3_capability.DeployOCR3Capability(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm transaction %s: %w", tx.Hash().String(), err)
	}
	resp := &deployResponse{
		Address: ocr3Addr,
		Tx:      tx.Hash(),
		Tv:      OCR3CapabilityTypeVersion,
	}
	c.contract = ocr3
	return resp, nil
}
