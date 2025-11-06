package tron

import (
	"context"
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// RemoveFeedProxyMappingChangeset is a changeset that removes the feed-aggregator proxy mapping from DataFeedsCache contract.
var RemoveFeedProxyMappingChangeset = cldf.CreateChangeSet(removeFeedProxyMappingLogic, removeFeedProxyMappingPrecondition)

func removeFeedProxyMappingLogic(env cldf.Environment, c types.RemoveFeedProxyTronConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.TronChains()[c.ChainSelector]

	txInfo, err := chain.TriggerContractAndConfirm(context.Background(), c.CacheAddress, "removeDataIdMappingsForProxies(address[])", []any{"address[]", c.ProxyAddresses}, c.TriggerOptions)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", txInfo.ID, err)
	}

	return cldf.ChangesetOutput{}, nil
}

func removeFeedProxyMappingPrecondition(env cldf.Environment, c types.RemoveFeedProxyTronConfig) error {
	_, ok := env.BlockChains.TronChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if len(c.ProxyAddresses) == 0 {
		return errors.New("proxy addresses must not be empty")
	}

	return changeset.ValidateCacheForTronChain(env, c.ChainSelector, c.CacheAddress)
}
