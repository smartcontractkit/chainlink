package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	moduleplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployPlatform(chain deployment.AptosChain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformResponse, error) {
	address, pendingTX, platformModule, err := moduleplatform.DeployToObject(chain.DeployerSigner, chain.Client, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkPlatform: %w", err)
	}

	_, err = chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkPlatform: %w", err)
	}

	// ChainlinkPlatform package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkPlatform 1.0.0"
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployPlatformResponse{
		Address:  address,
		Tx:       pendingTX.Hash,
		Tv:       tv,
		Contract: &platformModule,
	}
	return resp, nil
}
