package changeset

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
)

var _ deployment.ChangeSet[types.MigrationConfig] = MigrateFeedsChangeset

type MigrationSchema struct {
	Address        string                    `json:"address"`
	TypeAndVersion deployment.TypeAndVersion `json:"typeAndVersion"`
	FeedId         string                    `json:"feedId"` // without 0x prefix
	Description    string                    `json:"description"`
}

// MigrateFeedsChangeset Migrates feeds to DataFeedsCache contract.
// 1. It reads the existing Aggregator Proxy contract addresses from the input file and saves them to the address book.
// 2. It reads the data ids and descriptions from the input file and sets the feed config on the DataFeedsCache contract.
func MigrateFeedsChangeset(env deployment.Environment, c types.MigrationConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load on chain state %w", err)
	}
	chain, ok := env.Chains[c.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	chainState, ok := state.Chains[c.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in on chain state")
	}
	if chainState.DataFeedsCache == nil {
		return deployment.ChangesetOutput{}, errors.New("DataFeedsCache not found in on chain state")
	}
	contract, ok := chainState.DataFeedsCache[c.CacheAddress]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("contract not found in on chain state")
	}

	proxies, err := shared.LoadJSON[[]*MigrationSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load addresses input file: %w", err)
	}

	dataIds := make([][16]byte, len(proxies))
	addresses := make([]common.Address, len(proxies))
	descriptions := make([]string, len(proxies))
	for i, proxy := range proxies {
		dataIdBytes16, err := shared.ConvertHexToBytes16(proxy.FeedId)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("cannot convert hex to bytes %s: %w", proxy.FeedId, err)
		}

		dataIds[i] = dataIdBytes16
		addresses[i] = common.HexToAddress(proxy.Address)
		descriptions[i] = proxy.Description

		proxy.TypeAndVersion.AddLabel(proxy.Description)
		err = ab.Save(
			c.ChainSelector,
			proxy.Address,
			proxy.TypeAndVersion,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address %s: %w", proxy.Address, err)
		}
	}

	// Set the feed config
	tx, err := contract.SetDecimalFeedConfigs(chain.DeployerKey, dataIds, descriptions, c.WorkflowMetadata)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed config %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	// Set the proxy to dataId mapping
	tx, err = contract.UpdateDataIdMappingsForProxies(chain.DeployerKey, addresses, dataIds)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update feed proxy mapping %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
