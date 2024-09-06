package keystone

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type deployerSet struct {
	ocr3Deployer                 *OCR3Deployer
	keystoneForwarderDeployer    *KeystoneForwarderDeployer
	capabilitiesRegistryDeployer *CapabilitiesRegistryDeployer
}

type deployContractsRequest struct {
	chain           deployment.Chain
	isRegistryChain bool
	ad              deployment.AddressBook
}

type deployContractSetResponse struct {
	*deployerSet
	deployment.AddressBook
}

func deployContracts(lggr logger.Logger, req deployContractsRequest) (*deployContractSetResponse, error) {
	if req.ad == nil {
		req.ad = deployment.NewMemoryAddressBook()
	}
	resp := &deployContractSetResponse{
		AddressBook: req.ad,
		deployerSet: &deployerSet{},
	}
	singleRequest := deployRequest{Chain: req.chain}

	// cap reg and ocr3 only deployed on registry chain
	if req.isRegistryChain {
		capabilitiesRegistryDeployer := CapabilitiesRegistryDeployer{lggr: lggr}
		capabilitiesRegistryResp, err := capabilitiesRegistryDeployer.deploy(singleRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
		}
		err = resp.AddressBook.Save(req.chain.Selector, capabilitiesRegistryResp.Address.String(), capabilitiesRegistryResp.Tv)
		if err != nil {
			return nil, fmt.Errorf("failed to save CapabilitiesRegistry: %w", err)
		}
		resp.capabilitiesRegistryDeployer = &capabilitiesRegistryDeployer

		ocr3Deployer := OCR3Deployer{lggr: lggr}
		ocr3Resp, err := ocr3Deployer.deploy(singleRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
		}
		err = resp.AddressBook.Save(req.chain.Selector, ocr3Resp.Address.String(), ocr3Resp.Tv)
		if err != nil {
			return nil, fmt.Errorf("failed to save OCR3Capability: %w", err)
		}
		resp.ocr3Deployer = &ocr3Deployer
	}
	forwarderDeployer := KeystoneForwarderDeployer{lggr: lggr}
	forwarderResp, err := forwarderDeployer.deploy(singleRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	err = resp.AddressBook.Save(req.chain.Selector, forwarderResp.Address.String(), forwarderResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save KeystoneForwarder: %w", err)
	}
	resp.keystoneForwarderDeployer = &forwarderDeployer

	return resp, nil
}
