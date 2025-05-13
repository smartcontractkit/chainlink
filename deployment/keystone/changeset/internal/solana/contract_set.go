package solana

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// DeployForwarder deploys the KeystoneForwarder contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployForwarder(ctx context.Context, chain deployment.SolChain, ab cldf.AddressBook) (*DeployResponse, error) {
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
