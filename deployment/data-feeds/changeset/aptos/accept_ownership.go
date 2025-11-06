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

// AcceptOwnershipChangeset  accepts ownership of Registry/Router contract on Aptos
var AcceptOwnershipChangeset = cldf.CreateChangeSet(acceptDataFeedsOwnershipLogic, acceptDataFeedsOwnershipPrecondition)

func acceptDataFeedsOwnershipLogic(env cldf.Environment, c types.AcceptDataFeedsAptosOwnershipConfig) (cldf.ChangesetOutput, error) {
	state, _ := changeset.LoadAptosOnchainState(env)
	chain := env.BlockChains.AptosChains()[c.ChainSelector]
	chainState := state.AptosChains[c.ChainSelector]
	contractAddress := aptos.AccountAddress{}
	_ = contractAddress.ParseStringRelaxed(c.Address)
	contract := *chainState.DataFeeds[contractAddress]

	txOps := &bind.TransactOpts{
		Signer: chain.DeployerSigner,
	}

	if c.AcceptRegistry {
		submitResult, err := contract.Registry().AcceptOwnership(txOps)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tx, err := chain.Client.WaitForTransaction(submitResult.Hash)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if !tx.Success {
			return cldf.ChangesetOutput{}, fmt.Errorf("registry Accept ownership transaction failed: %s", tx.Hash)
		}
		env.Logger.Info("Registry Accept ownership transaction succeeded", tx.Hash)
	}
	if c.AcceptRouter {
		submitResult, err := contract.Router().AcceptOwnership(txOps)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tx, err := chain.Client.WaitForTransaction(submitResult.Hash)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if !tx.Success {
			return cldf.ChangesetOutput{}, fmt.Errorf("router Accept ownership transaction failed: %s", tx.Hash)
		}
		env.Logger.Info("Router Accept ownership transaction succeeded", tx.Hash)
	}

	return cldf.ChangesetOutput{}, nil
}

func acceptDataFeedsOwnershipPrecondition(env cldf.Environment, c types.AcceptDataFeedsAptosOwnershipConfig) error {
	if !c.AcceptRouter && !c.AcceptRegistry {
		return errors.New("at least one of AcceptRouter or AcceptRouter must be true")
	}

	return changeset.ValidateCacheForAptosChain(env, c.ChainSelector, c.Address)
}
