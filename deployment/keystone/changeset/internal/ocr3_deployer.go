package internal

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment"
)

type OCR3Deployer struct {
	lggr     logger.Logger
	contract *ocr3_capability.OCR3Capability
}

func NewOCR3Deployer() (*OCR3Deployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &OCR3Deployer{lggr: lggr}, nil
}

func (c *OCR3Deployer) deploy(req DeployRequest) (*DeployResponse, error) {
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
	tvStr, err := ocr3.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	resp := &DeployResponse{
		Address: ocr3Addr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = ocr3
	return resp, nil
}
