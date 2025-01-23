package changeset

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type DeployerGroup struct {
	e            deployment.Environment
	state        CCIPOnChainState
	mcmConfig    *MCMSConfig
	transactions map[uint64][]*types.Transaction
}

// DeployerGroup is an abstraction that lets developers write their changeset
// without needing to know if it's executed using a DeployerKey or an MCMS proposal.
//
// Example usage:
//
//	deployerGroup := NewDeployerGroup(e, state, mcmConfig)
//	selector := 0
//	# Get the right deployer key for the chain
//	deployer := deployerGroup.GetDeployer(selector)
//	state.Chains[selector].RMNRemote.Curse()
//	# Execute the transaction or create the proposal
//	deployerGroup.Enact("Curse RMNRemote")
func NewDeployerGroup(e deployment.Environment, state CCIPOnChainState, mcmConfig *MCMSConfig) *DeployerGroup {
	return &DeployerGroup{
		e:            e,
		mcmConfig:    mcmConfig,
		state:        state,
		transactions: make(map[uint64][]*types.Transaction),
	}
}

func (d *DeployerGroup) GetDeployer(chain uint64) (*bind.TransactOpts, error) {
	txOpts := d.e.Chains[chain].DeployerKey
	if d.mcmConfig != nil {
		txOpts = deployment.SimTransactOpts()
		txOpts = &bind.TransactOpts{
			From:       d.state.Chains[chain].Timelock.Address(),
			Signer:     txOpts.Signer,
			GasLimit:   txOpts.GasLimit,
			GasPrice:   txOpts.GasPrice,
			Nonce:      txOpts.Nonce,
			Value:      txOpts.Value,
			GasFeeCap:  txOpts.GasFeeCap,
			GasTipCap:  txOpts.GasTipCap,
			Context:    txOpts.Context,
			AccessList: txOpts.AccessList,
			NoSend:     txOpts.NoSend,
		}
	}
	sim := &bind.TransactOpts{
		From:       txOpts.From,
		Signer:     txOpts.Signer,
		GasLimit:   txOpts.GasLimit,
		GasPrice:   txOpts.GasPrice,
		Nonce:      txOpts.Nonce,
		Value:      txOpts.Value,
		GasFeeCap:  txOpts.GasFeeCap,
		GasTipCap:  txOpts.GasTipCap,
		Context:    txOpts.Context,
		AccessList: txOpts.AccessList,
		NoSend:     true,
	}
	oldSigner := sim.Signer

	var startingNonce *big.Int
	if txOpts.Nonce != nil {
		startingNonce = new(big.Int).Set(txOpts.Nonce)
	} else {
		nonce, err := d.e.Chains[chain].Client.PendingNonceAt(context.Background(), txOpts.From)
		if err != nil {
			return nil, fmt.Errorf("could not get nonce for deployer: %w", err)
		}
		startingNonce = new(big.Int).SetUint64(nonce)
	}

	sim.Signer = func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
		// Update the nonce to consider the transactions that have been sent
		sim.Nonce = big.NewInt(0).Add(startingNonce, big.NewInt(int64(len(d.transactions[chain]))+1))

		tx, err := oldSigner(a, t)
		if err != nil {
			return nil, err
		}
		d.transactions[chain] = append(d.transactions[chain], tx)
		return tx, nil
	}
	return sim, nil
}

func (d *DeployerGroup) Enact(deploymentDescription string) (deployment.ChangesetOutput, error) {
	if d.mcmConfig != nil {
		return d.enactMcms(deploymentDescription)
	}

	return d.enactDeployer()
}

func (d *DeployerGroup) enactMcms(deploymentDescription string) (deployment.ChangesetOutput, error) {
	batches := make([]timelock.BatchChainOperation, 0)
	for selector, txs := range d.transactions {
		mcmOps := make([]mcms.Operation, len(txs))
		for i, tx := range txs {
			mcmOps[i] = mcms.Operation{
				To:    *tx.To(),
				Data:  tx.Data(),
				Value: tx.Value(),
			}
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(selector),
			Batch:           mcmOps,
		})
	}

	timelocksPerChain := BuildTimelockAddressPerChain(d.e, d.state)

	proposerMCMSes := BuildProposerPerChain(d.e, d.state)

	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		deploymentDescription,
		d.mcmConfig.MinDelay,
	)

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal %w", err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func (d *DeployerGroup) enactDeployer() (deployment.ChangesetOutput, error) {
	for selector, txs := range d.transactions {
		for _, tx := range txs {
			err := d.e.Chains[selector].Client.SendTransaction(context.Background(), tx)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to send transaction: %w", err)
			}

			_, err = d.e.Chains[selector].Confirm(tx)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("waiting for tx to be mined failed: %w", err)
			}
		}
	}
	return deployment.ChangesetOutput{}, nil
}

func BuildTimelockPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*proposalutils.TimelockExecutionContracts {
	timelocksPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain.Selector].Timelock,
			CallProxy: state.Chains[chain.Selector].CallProxy,
		}
	}
	return timelocksPerChain
}

func BuildTimelockAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]common.Address {
	timelocksPerChain := BuildTimelockPerChain(e, state)
	timelockAddressPerChain := make(map[uint64]common.Address)
	for chain, timelock := range timelocksPerChain {
		timelockAddressPerChain[chain] = timelock.Timelock.Address()
	}
	return timelockAddressPerChain
}

func BuildProposerPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*gethwrappers.ManyChainMultiSig {
	proposerPerChain := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range e.Chains {
		proposerPerChain[chain.Selector] = state.Chains[chain.Selector].ProposerMcm
	}
	return proposerPerChain
}
