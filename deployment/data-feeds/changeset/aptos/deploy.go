package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	moduleplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployPlatform(chain cldf.AptosChain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformResponse, error) {
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
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s, %w", pendingTX.Hash, err)
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

func DeployDataFeeds(chain cldf.AptosChain, owner aptos.AccountAddress, platform aptos.AccountAddress, labels []string) (*types.DeployDataFeedsResponse, error) {
	address, pendingTX, feedsModule, err := modulefeeds.DeployToObject(chain.DeployerSigner, chain.Client, owner, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkDataFeeds: %w", err)
	}

	_, err = chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkDataFeeds: %s, %w", pendingTX.Hash, err)
	}

	// ChainlinkDataFeeds package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkDataFeeds 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployDataFeedsResponse{
		Address:  address,
		Tx:       pendingTX.Hash,
		Tv:       tv,
		Contract: &feedsModule,
	}
	return resp, nil
}
