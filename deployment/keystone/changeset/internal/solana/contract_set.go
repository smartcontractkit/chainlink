package solana

import (
	"context"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func DeployForwarder(ctx context.Context, chain cldf.SolChain, ab cldf.AddressBook) (*DeployResponse, error) {
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
