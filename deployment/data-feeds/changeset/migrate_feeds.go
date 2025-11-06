package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// MigrateFeedsChangeset Migrates feeds to DataFeedsCache contract.
// 1. It reads the existing Aggregator Proxy contract addresses from the input file and saves them to the address book.
// 2. It reads the data ids and descriptions from the input file and sets the feed config on the DataFeedsCache contract.
// Returns a new datastore with the deployed AggregatorProxy addresses.
var MigrateFeedsChangeset = cldf.CreateChangeSet(migrateFeedsLogic, migrateFeedsPrecondition)

func migrateFeedsLogic(env cldf.Environment, c types.MigrationConfig) (cldf.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.BlockChains.EVMChains()[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]
	ds := datastore.NewMemoryDataStore()

	proxies := c.Proxies

	var feedIDs []string
	addresses := make([]common.Address, len(proxies))
	descriptions := make([]string, len(proxies))
	for i, proxy := range proxies {
		feedIDs = append(feedIDs, proxy.FeedID)
		addresses[i] = common.HexToAddress(proxy.Address)
		descriptions[i] = proxy.Description

		proxy.TypeAndVersion.AddLabel(proxy.Description)
		if err := ds.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: c.ChainSelector,
				Address:       proxy.Address,
				Type:          datastore.ContractType(proxy.TypeAndVersion.Type),
				Version:       &proxy.TypeAndVersion.Version,
				Qualifier:     proxy.Description,
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	dataIDs, _ := FeedIDsToBytes16(feedIDs)

	// Set the feed config
	tx, err := contract.SetDecimalFeedConfigs(chain.DeployerKey, dataIDs, descriptions, c.WorkflowMetadata)
	if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	// Set the proxy to dataId mapping
	tx, err = contract.UpdateDataIdMappingsForProxies(chain.DeployerKey, addresses, dataIDs)
	if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return cldf.ChangesetOutput{DataStore: ds}, nil
}

func migrateFeedsPrecondition(env cldf.Environment, c types.MigrationConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	proxies := c.Proxies
	var feedIDs []string
	for _, proxy := range proxies {
		feedIDs = append(feedIDs, proxy.FeedID)
	}
	_, err := FeedIDsToBytes16(feedIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	if len(c.WorkflowMetadata) == 0 {
		return errors.New("workflow metadata is required")
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
