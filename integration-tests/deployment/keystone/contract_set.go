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
	// this is mutated in the Deploy* functions
	resp := &deployContractSetResponse{
		AddressBook: req.ad,
	}

	// cap reg and ocr3 only deployed on registry chain
	if req.isRegistryChain {
		err := DeployCapabilitiesRegistry(lggr, req.chain, resp.AddressBook)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
		}
		err = DeployOCR3(lggr, req.chain, resp.AddressBook)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
		}
	}
	err := DeployForwarder(lggr, req.chain, resp.AddressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	return resp, nil
}

// DeployCapabilitiesRegistry deploys the CapabilitiesRegistry contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployCapabilitiesRegistry(lggr logger.Logger, chain deployment.Chain, ab deployment.AddressBook) error {
	capabilitiesRegistryDeployer := CapabilitiesRegistryDeployer{lggr: lggr}
	capabilitiesRegistryResp, err := capabilitiesRegistryDeployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	err = ab.Save(chain.Selector, capabilitiesRegistryResp.Address.String(), capabilitiesRegistryResp.Tv)
	if err != nil {
		return fmt.Errorf("failed to save CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", CapabilityRegistryTypeVersion.String(), chain.Selector, capabilitiesRegistryResp.Address.String())
	return nil
}

// DeployOCR3 deploys the OCR3Capability contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployOCR3(lggr logger.Logger, chain deployment.Chain, ab deployment.AddressBook) error {
	ocr3Deployer := OCR3Deployer{lggr: lggr}
	ocr3Resp, err := ocr3Deployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	err = ab.Save(chain.Selector, ocr3Resp.Address.String(), ocr3Resp.Tv)
	if err != nil {
		return fmt.Errorf("failed to save OCR3Capability: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", ocr3Resp.Tv.String(), chain.Selector, ocr3Resp.Address.String())
	return nil
}

// DeployForwarder deploys the KeystoneForwarder contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployForwarder(lggr logger.Logger, chain deployment.Chain, ab deployment.AddressBook) error {
	forwarderDeployer := KeystoneForwarderDeployer{lggr: lggr}
	forwarderResp, err := forwarderDeployer.deploy(deployRequest{Chain: chain})
	if err != nil {
		return fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	err = ab.Save(chain.Selector, forwarderResp.Address.String(), forwarderResp.Tv)
	if err != nil {
		return fmt.Errorf("failed to save KeystoneForwarder: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResp.Tv.String(), chain.Selector, forwarderResp.Address.String())
	return nil
}
