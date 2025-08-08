package tron

import (
	"context"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedAdminChangeset is a changeset that sets/removes an admin on DataFeedsCache contract.
var SetFeedAdminChangeset = cldf.CreateChangeSet(setFeedAdminLogic, setFeedAdminPrecondition)

func setFeedAdminLogic(env cldf.Environment, c types.SetFeedAdminTronConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.TronChains()[c.ChainSelector]

	txInfo, err := chain.TriggerContractAndConfirm(context.Background(), c.CacheAddress, "setFeedAdmin(address,bool)", []interface{}{"address", c.AdminAddress, "bool", c.IsAdmin}, c.TriggerOptions)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", txInfo.ID, err)
	}

	return cldf.ChangesetOutput{}, nil
}

func setFeedAdminPrecondition(env cldf.Environment, c types.SetFeedAdminTronConfig) error {
	_, ok := env.BlockChains.TronChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	return changeset.ValidateCacheForTronChain(env, c.ChainSelector, c.CacheAddress)
}
