package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedConfigChangeset  sets feed configuration on Registry contract on Aptos
var SetFeedConfigChangeset = cldf.CreateChangeSet(setFeedConfigLogic, setFeedConfigPrecondition)

func setFeedConfigLogic(env cldf.Environment, c types.SetRegistryFeedConfig) (cldf.ChangesetOutput, error) {
	state, _ := changeset.LoadAptosOnchainState(env)
	chain := env.BlockChains.AptosChains()[c.ChainSelector]
	chainState := state.AptosChains[c.ChainSelector]
	cacheAccountAddress := aptos.AccountAddress{}
	_ = cacheAccountAddress.ParseStringRelaxed(c.CacheAddress)
	contract := *chainState.DataFeeds[cacheAccountAddress]

	dataIDs, _ := changeset.FeedIDsToBytes(c.DataIDs)

	var configID []byte

	txOps := &bind.TransactOpts{
		Signer: chain.DeployerSigner,
	}

	submitResult, err := contract.Registry().SetFeeds(txOps, dataIDs, c.Descriptions, configID)

	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	tx, err := chain.Client.WaitForTransaction(submitResult.Hash)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if !tx.Success {
		fmt.Println("Transaction failed", tx.Hash)
	} else {
		fmt.Println("Transaction succeeded", tx.Hash)
	}

	return cldf.ChangesetOutput{}, nil
}

func setFeedConfigPrecondition(env cldf.Environment, c types.SetRegistryFeedConfig) error {
	if (len(c.DataIDs) == 0) || (len(c.Descriptions) == 0) {
		return errors.New("dataIDs and descriptions must not be empty")
	}
	if len(c.DataIDs) != len(c.Descriptions) {
		return errors.New("dataIDs and descriptions must have the same length")
	}

	_, err := changeset.FeedIDsToBytes(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes: %w", err)
	}

	return changeset.ValidateCacheForAptosChain(env, c.ChainSelector, c.CacheAddress)
}
