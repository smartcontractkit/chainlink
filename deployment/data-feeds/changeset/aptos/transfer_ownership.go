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

// TransferOwnershipChangeset  transfers ownership of Registry/Router contract on Aptos
var TransferOwnershipChangeset = cldf.CreateChangeSet(transferDataFeedsOwnershipLogic, transferDataFeedsOwnershipPrecondition)

func transferDataFeedsOwnershipLogic(env cldf.Environment, c types.TransferDataFeedsAptosOwnershipConfig) (cldf.ChangesetOutput, error) {
	state, _ := changeset.LoadAptosOnchainState(env)
	chain := env.BlockChains.AptosChains()[c.ChainSelector]
	chainState := state.AptosChains[c.ChainSelector]
	contractAddress := aptos.AccountAddress{}
	_ = contractAddress.ParseStringRelaxed(c.Address)
	contract := *chainState.DataFeeds[contractAddress]

	newOwner := aptos.AccountAddress{}
	_ = newOwner.ParseStringRelaxed(c.NewOwner)

	txOps := &bind.TransactOpts{
		Signer: chain.DeployerSigner,
	}

	if c.TransferRegistry {
		submitResult, err := contract.Registry().TransferOwnership(txOps, newOwner)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tx, err := chain.Client.WaitForTransaction(submitResult.Hash)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if !tx.Success {
			return cldf.ChangesetOutput{}, fmt.Errorf("registry Transfer ownership transaction failed: %s", tx.Hash)
		}
		env.Logger.Info("Registry Transfer ownership transaction succeeded", tx.Hash)
	}
	if c.TransferRouter {
		submitResult, err := contract.Router().TransferOwnership(txOps, newOwner)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tx, err := chain.Client.WaitForTransaction(submitResult.Hash)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if !tx.Success {
			return cldf.ChangesetOutput{}, fmt.Errorf("router Transfer ownership transaction failed: %s", tx.Hash)
		}
		env.Logger.Info("Router Transfer ownership transaction succeeded", tx.Hash)
	}

	return cldf.ChangesetOutput{}, nil
}

func transferDataFeedsOwnershipPrecondition(env cldf.Environment, c types.TransferDataFeedsAptosOwnershipConfig) error {
	if !c.TransferRegistry && !c.TransferRouter {
		return errors.New("at least one of TransferRegistry or TransferRouter must be true")
	}
	newOwner := aptos.AccountAddress{}
	err := newOwner.ParseStringRelaxed(c.NewOwner)
	if err != nil {
		return fmt.Errorf("failed to parse new owner address %w", err)
	}

	return changeset.ValidateCacheForAptosChain(env, c.ChainSelector, c.Address)
}
