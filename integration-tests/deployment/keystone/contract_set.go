package keystone

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type deployContractsRequest struct {
	chain           deployment.Chain
	isRegistryChain bool
	ad              deployment.AddressBook
}

type deployContractSetResponse struct {
	deployment.AddressBook
}

func deployContractsToChain(lggr logger.Logger, req deployContractsRequest) (*deployContractSetResponse, error) {
	if req.ad == nil {
		req.ad = deployment.NewMemoryAddressBook()
	}
	resp := &deployContractSetResponse{
		AddressBook: req.ad,
	}

	// cap reg and ocr3 only deployed on registry chain
	if req.isRegistryChain {
		addrBook, err := DeployCapabilitiesRegistry(lggr, req.chain)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
		}
		resp.AddressBook = addrBook
		addrBook, err = DeployOCR3(lggr, req.chain)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
		}
		err = resp.AddressBook.Merge(addrBook)
		if err != nil {
			return nil, fmt.Errorf("failed to merge OCR3Capability: %w", err)
		}
	}
	addrBook, err := DeployForwarder(lggr, req.chain)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	err = resp.AddressBook.Merge(addrBook)
	if err != nil {
		return nil, fmt.Errorf("failed to merge KeystoneForwarder: %w", err)
	}

	return resp, nil
}

func DeployCapabilitiesRegistry(lggr logger.Logger, chain deployment.Chain) (deployment.AddressBook, error) {
	out := deployment.NewMemoryAddressBook()
	capabilitiesRegistryDeployer := CapabilitiesRegistryDeployer{lggr: lggr}
	capabilitiesRegistryResp, err := capabilitiesRegistryDeployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	err = out.Save(chain.Selector, capabilitiesRegistryResp.Address.String(), capabilitiesRegistryResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", CapabilityRegistryTypeVersion.String(), chain.Selector, capabilitiesRegistryResp.Address.String())
	return out, nil
}

func DeployOCR3(lggr logger.Logger, chain deployment.Chain) (deployment.AddressBook, error) {
	out := deployment.NewMemoryAddressBook()
	ocr3Deployer := OCR3Deployer{lggr: lggr}
	ocr3Resp, err := ocr3Deployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	err = out.Save(chain.Selector, ocr3Resp.Address.String(), ocr3Resp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save OCR3Capability: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", ocr3Resp.Tv.String(), chain.Selector, ocr3Resp.Address.String())
	return out, nil
}

func DeployForwarder(lggr logger.Logger, chain deployment.Chain) (deployment.AddressBook, error) {
	out := deployment.NewMemoryAddressBook()
	forwarderDeployer := KeystoneForwarderDeployer{lggr: lggr}
	forwarderResp, err := forwarderDeployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	err = out.Save(chain.Selector, forwarderResp.Address.String(), forwarderResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save KeystoneForwarder: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResp.Tv.String(), chain.Selector, forwarderResp.Address.String())
	return out, nil
}
