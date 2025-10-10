package tron

import (
	"context"
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// UpdateDataIDProxyChangeset is a changeset that updates the proxy-dataId mapping on DataFeedsCache contract.
var UpdateDataIDProxyChangeset = cldf.CreateChangeSet(updateDataIDProxyLogic, updateDataIDProxyPrecondition)

func updateDataIDProxyLogic(env cldf.Environment, c types.UpdateDataIDProxyTronConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.TronChains()[c.ChainSelector]

	dataIDs, err := changeset.FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data ids: %s, %w", c.DataIDs, err)
	}

	txInfo, err := chain.TriggerContractAndConfirm(context.Background(), c.CacheAddress, "updateDataIdMappingsForProxies(address[],bytes16[])", []any{"address[]", c.ProxyAddresses, "bytes16[]", dataIDs}, c.TriggerOptions)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", txInfo.ID, err)
	}

	return cldf.ChangesetOutput{}, nil
}

func updateDataIDProxyPrecondition(env cldf.Environment, c types.UpdateDataIDProxyTronConfig) error {
	_, ok := env.BlockChains.TronChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if len(c.ProxyAddresses) == 0 || len(c.DataIDs) == 0 {
		return errors.New("empty proxies or dataIds")
	}
	if len(c.DataIDs) != len(c.ProxyAddresses) {
		return errors.New("dataIds and proxies length mismatch")
	}
	_, err := changeset.FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	return changeset.ValidateCacheForTronChain(env, c.ChainSelector, c.CacheAddress)
}
