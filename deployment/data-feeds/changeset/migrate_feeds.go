package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// MigrateFeedsChangeset Migrates feeds to DataFeedsCache contract.
// 1. It reads the existing Aggregator Proxy contract addresses from the input file and saves them to the address book.
// 2. It reads the data ids and descriptions from the input file and sets the feed config on the DataFeedsCache contract.
// Returns a new addressbook with the deployed AggregatorProxy addresses.
var MigrateFeedsChangeset = cldf.CreateChangeSet(migrateFeedsLogic, migrateFeedsPrecondition)

type MigrationSchema struct {
	Address        string                    `json:"address"`
	TypeAndVersion deployment.TypeAndVersion `json:"typeAndVersion"`
	FeedID         string                    `json:"feedId"`
	Description    string                    `json:"description"`
}

func migrateFeedsLogic(env deployment.Environment, c types.MigrationConfig) (deployment.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]
	ab := deployment.NewMemoryAddressBook()

	proxies, _ := LoadJSON[[]*MigrationSchema](c.InputFileName, c.InputFS)

	var feedIDs []string
	addresses := make([]common.Address, len(proxies))
	descriptions := make([]string, len(proxies))
	for i, proxy := range proxies {
		feedIDs = append(feedIDs, proxy.FeedID)
		addresses[i] = common.HexToAddress(proxy.Address)
		descriptions[i] = proxy.Description

		proxy.TypeAndVersion.AddLabel(proxy.Description)
		err := ab.Save(
			c.ChainSelector,
			proxy.Address,
			proxy.TypeAndVersion,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address %s: %w", proxy.Address, err)
		}
	}

	dataIDs, _ := FeedIDsToBytes16(feedIDs)

	// Set the feed config
	tx, err := contract.SetDecimalFeedConfigs(chain.DeployerKey, dataIDs, descriptions, c.WorkflowMetadata)
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	// Set the proxy to dataId mapping
	tx, err = contract.UpdateDataIdMappingsForProxies(chain.DeployerKey, addresses, dataIDs)
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func migrateFeedsPrecondition(env deployment.Environment, c types.MigrationConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	proxies, err := LoadJSON[[]*MigrationSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load addresses input file: %w", err)
	}
	var feedIDs []string
	for _, proxy := range proxies {
		feedIDs = append(feedIDs, proxy.FeedID)
	}
	_, err = FeedIDsToBytes16(feedIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	if len(c.WorkflowMetadata) == 0 {
		return errors.New("workflow metadata is required")
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
