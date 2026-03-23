package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	moduleplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	moduleplatform_secondary "github.com/smartcontractkit/chainlink-aptos/bindings/platform_secondary"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployDataFeeds(chain cldf_aptos.Chain, owner aptos.AccountAddress, platform aptos.AccountAddress, secondaryPlatform aptos.AccountAddress, labels []string) (*types.DeployDataFeedsResponse, error) {
	address, pendingTX, feedsModule, err := modulefeeds.DeployToObject(chain.DeployerSigner, chain.Client, owner, platform, owner, secondaryPlatform)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkDataFeeds: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkDataFeeds: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, fmt.Errorf("ChainlinkDataFeeds deployment transaction failed: %s", tx.VmStatus)
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
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &feedsModule,
	}
	return resp, nil
}

func DeployPlatform(chain cldf_aptos.Chain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformResponse, error) {
	if owner == (aptos.AccountAddress{}) {
		owner = chain.DeployerSigner.AccountAddress()
	}
	address, pendingTX, platformModule, err := moduleplatform.DeployToObject(chain.DeployerSigner, chain.Client, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkPlatform: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkPlatform: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, fmt.Errorf("ChainlinkPlatform deployment transaction failed: %s", tx.Hash)
	}
	// ChainlinkPlatform package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkPlatform 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployPlatformResponse{
		Address:  address,
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &platformModule,
	}
	return resp, nil
}

func DeployPlatformSecondary(chain cldf_aptos.Chain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformSecondaryResponse, error) {
	if owner == (aptos.AccountAddress{}) {
		owner = chain.DeployerSigner.AccountAddress()
	}
	address, pendingTX, platformModule, err := moduleplatform_secondary.DeployToObject(chain.DeployerSigner, chain.Client, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkPlatformSecondary: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkPlatformSecondary: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, fmt.Errorf("ChainlinkPlatformSecondary deployment transaction failed: %s", tx.Hash)
	}
	// ChainlinkPlatformSecondary package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkPlatformSecondary 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployPlatformSecondaryResponse{
		Address:  address,
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &platformModule,
	}
	return resp, nil
}
