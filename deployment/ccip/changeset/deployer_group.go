package changeset

import (
	"context"
	"fmt"

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

func NewDeployerGroup(e deployment.Environment, state CCIPOnChainState, mcmConfig *MCMSConfig) *DeployerGroup {
	return &DeployerGroup{
		e:            e,
		mcmConfig:    mcmConfig,
		state:        state,
		transactions: make(map[uint64][]*types.Transaction),
	}
}

func (d *DeployerGroup) getDeployer(chain uint64) *bind.TransactOpts {
	sim := deployment.SimTransactOpts()
	oldSigner := sim.Signer
	sim.Signer = func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
		tx, err := oldSigner(a, t)
		if err != nil {
			return nil, err
		}
		d.transactions[chain] = append(d.transactions[chain], tx)
		return tx, nil
	}
	return sim
}

func (d *DeployerGroup) enact(deploymentDescription string) (deployment.ChangesetOutput, error) {
	if d.mcmConfig != nil {
		return d.enactMcms(deploymentDescription)
	} else {
		return d.enactDeployer()
	}
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

	timelocksPerChain := buildTimelockAddressPerChain(d.e, d.state)

	proposerMCMSes := buildProposerPerChain(d.e, d.state)

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
		}
	}
	return deployment.ChangesetOutput{}, nil
}

func buildTimelockPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*proposalutils.TimelockExecutionContracts {
	timelocksPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain.Selector].Timelock,
			CallProxy: state.Chains[chain.Selector].CallProxy,
		}
	}
	return timelocksPerChain
}

func buildTimelockAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]common.Address {
	timelocksPerChain := buildTimelockPerChain(e, state)
	timelockAddressPerChain := make(map[uint64]common.Address)
	for chain, timelock := range timelocksPerChain {
		timelockAddressPerChain[chain] = timelock.Timelock.Address()
	}
	return timelockAddressPerChain
}

func buildProposerPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*gethwrappers.ManyChainMultiSig {
	proposerPerChain := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range e.Chains {
		proposerPerChain[chain.Selector] = state.Chains[chain.Selector].ProposerMcm
	}
	return proposerPerChain
}
