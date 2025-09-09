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

// SetWorkflowConfigChangeset sets workflow configuration on Registry contract on Aptos.
var SetWorkflowConfigChangeset = cldf.CreateChangeSet(setWorkflowConfigLogic, setWorkflowConfigPrecondition)

func setWorkflowConfigLogic(env cldf.Environment, c types.SetRegistryWorkflowConfig) (cldf.ChangesetOutput, error) {
	state, _ := changeset.LoadAptosOnchainState(env)
	chain := env.BlockChains.AptosChains()[c.ChainSelector]
	chainState := state.AptosChains[c.ChainSelector]
	cacheAccountAddress := aptos.AccountAddress{}
	_ = cacheAccountAddress.ParseStringRelaxed(c.CacheAddress)
	contract := *chainState.DataFeeds[cacheAccountAddress]

	var allowedWorkflowOwners [][]byte
	for _, o := range c.AllowedWorkflowOwners {
		b, _ := aptos.ParseHex(o)
		allowedWorkflowOwners = append(allowedWorkflowOwners, b)
	}

	var allowedWorkflowNames [][]byte
	for _, n := range c.AllowedWorkflowNames {
		b, _ := aptos.ParseHex(n)
		allowedWorkflowNames = append(allowedWorkflowNames, b)
	}

	txOps := &bind.TransactOpts{
		Signer: chain.DeployerSigner,
	}

	submitResult, err := contract.Registry().SetWorkflowConfig(txOps, allowedWorkflowOwners, allowedWorkflowNames)

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

func setWorkflowConfigPrecondition(env cldf.Environment, c types.SetRegistryWorkflowConfig) error {
	if len(c.AllowedWorkflowOwners) == 0 {
		return errors.New("allowedWorkflowOwners cannot be empty")
	}

	return changeset.ValidateCacheForAptosChain(env, c.ChainSelector, c.CacheAddress)
}
