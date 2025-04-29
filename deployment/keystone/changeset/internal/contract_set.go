package internal

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
)

type deployContractsRequest struct {
	chain           deployment.Chain
	isRegistryChain bool
	ad              deployment.AddressBook
}

// DeployCapabilitiesRegistry deploys the CapabilitiesRegistry contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployCapabilitiesRegistry(_ context.Context, chain deployment.Chain, ab deployment.AddressBook) (*DeployResponse, error) {
	capabilitiesRegistryDeployer, err := NewCapabilitiesRegistryDeployer()
	capabilitiesRegistryResp, err := capabilitiesRegistryDeployer.Deploy(DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	err = ab.Save(chain.Selector, capabilitiesRegistryResp.Address.String(), capabilitiesRegistryResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save CapabilitiesRegistry: %w", err)
	}
	return capabilitiesRegistryResp, nil
}

// DeployOCR3 deploys the OCR3Capability contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployOCR3(_ context.Context, chain deployment.Chain, ab deployment.AddressBook) (*DeployResponse, error) {
	ocr3Deployer, err := NewOCR3Deployer()
	if err != nil {
		return nil, fmt.Errorf("failed to create OCR3Deployer: %w", err)
	}
	ocr3Resp, err := ocr3Deployer.deploy(DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	err = ab.Save(chain.Selector, ocr3Resp.Address.String(), ocr3Resp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save OCR3Capability: %w", err)
	}

	return ocr3Resp, nil
}

// DeployForwarder deploys the KeystoneForwarder contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployForwarder(ctx context.Context, chain deployment.Chain, ab deployment.AddressBook) (*DeployResponse, error) {
	forwarderDeployer, err := NewKeystoneForwarderDeployer()
	if err != nil {
		return nil, fmt.Errorf("failed to create KeystoneForwarderDeployer: %w", err)
	}
	forwarderResp, err := forwarderDeployer.deploy(ctx, DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}
	err = ab.Save(chain.Selector, forwarderResp.Address.String(), forwarderResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save KeystoneForwarder: %w", err)
	}
	return forwarderResp, nil
}

// DeployFeedsConsumer deploys the KeystoneFeedsConsumer contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployFeedsConsumer(_ context.Context, chain deployment.Chain, ab deployment.AddressBook) (*DeployResponse, error) {
	consumerDeploy, err := NewKeystoneFeedsConsumerDeployer()
	if err != nil {
		return nil, err
	}
	consumerResp, err := consumerDeploy.deploy(DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy FeedsConsumer: %w", err)
	}
	err = ab.Save(chain.Selector, consumerResp.Address.String(), consumerResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save FeedsConsumer: %w", err)
	}
	return consumerResp, nil
}

// DeployForwarder deploys the BalanceReader contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployBalanceReader(_ context.Context, chain deployment.Chain, ab deployment.AddressBook) (*DeployResponse, error) {
	balanceReaderDeployer, err := NewBalanceReaderDeployer()
	if err != nil {
		return nil, fmt.Errorf("failed to create BalanceReaderDeployer: %w", err)
	}
	balanceReaderResp, err := balanceReaderDeployer.deploy(DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy BalanceReader: %w", err)
	}
	err = ab.Save(chain.Selector, balanceReaderResp.Address.String(), balanceReaderResp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save BalanceReader: %w", err)
	}
	return balanceReaderResp, nil
}
