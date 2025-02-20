package changeset

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// MCMSConfig defines timelock duration.
type MCMSConfig struct {
	MinDelay   time.Duration
	MCMSAction timelock.TimelockOperation
}

func (mcmsConfig *MCMSConfig) Validate() error {
	// to make it backwards compatible with the old MCMSConfig , if MCMSAction is not set, default to timelock.Schedule
	// TODO remove this after all the usages are updated to reflect canceller and bypasser with new mcmslib
	if mcmsConfig.MCMSAction == "" {
		mcmsConfig.MCMSAction = timelock.Schedule
	}
	if mcmsConfig.MCMSAction != timelock.Schedule &&
		mcmsConfig.MCMSAction != timelock.Cancel &&
		mcmsConfig.MCMSAction != timelock.Bypass {
		return fmt.Errorf("invalid MCMS type %s", mcmsConfig.MCMSAction)
	}
	return nil
}

type DeployerGroup struct {
	e                 deployment.Environment
	state             CCIPOnChainState
	mcmConfig         *MCMSConfig
	deploymentContext *DeploymentContext
}

type DeploymentContext struct {
	description    string
	transactions   map[uint64][]*types.Transaction
	previousConfig *DeploymentContext
}

func NewDeploymentContext(description string) *DeploymentContext {
	return &DeploymentContext{
		description:    description,
		transactions:   make(map[uint64][]*types.Transaction),
		previousConfig: nil,
	}
}

func (d *DeploymentContext) Fork(description string) *DeploymentContext {
	return &DeploymentContext{
		description:    description,
		transactions:   make(map[uint64][]*types.Transaction),
		previousConfig: d,
	}
}

type DeployerGroupWithContext interface {
	WithDeploymentContext(description string) *DeployerGroup
}

type deployerGroupBuilder struct {
	e         deployment.Environment
	state     CCIPOnChainState
	mcmConfig *MCMSConfig
}

func (d *deployerGroupBuilder) WithDeploymentContext(description string) *DeployerGroup {
	return &DeployerGroup{
		e:                 d.e,
		mcmConfig:         d.mcmConfig,
		state:             d.state,
		deploymentContext: NewDeploymentContext(description),
	}
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
func NewDeployerGroup(e deployment.Environment, state CCIPOnChainState, mcmConfig *MCMSConfig) DeployerGroupWithContext {
	return &deployerGroupBuilder{
		e:         e,
		mcmConfig: mcmConfig,
		state:     state,
	}
}

func (d *DeployerGroup) WithDeploymentContext(description string) *DeployerGroup {
	return &DeployerGroup{
		e:                 d.e,
		mcmConfig:         d.mcmConfig,
		state:             d.state,
		deploymentContext: d.deploymentContext.Fork(description),
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

	dc := d.deploymentContext
	sim.Signer = func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
		txCount, err := d.getTransactionCount(chain)
		if err != nil {
			return nil, err
		}

		currentNonce := big.NewInt(0).Add(startingNonce, txCount)

		tx, err := oldSigner(a, t)
		if err != nil {
			return nil, err
		}
		dc.transactions[chain] = append(dc.transactions[chain], tx)
		// Update the nonce to consider the transactions that have been sent
		sim.Nonce = big.NewInt(0).Add(currentNonce, big.NewInt(1))
		return tx, nil
	}
	return sim, nil
}

func (d *DeployerGroup) getContextChainInOrder() []*DeploymentContext {
	contexts := make([]*DeploymentContext, 0)
	for c := d.deploymentContext; c != nil; c = c.previousConfig {
		contexts = append(contexts, c)
	}
	slices.Reverse(contexts)
	return contexts
}

func (d *DeployerGroup) getTransactions() map[uint64][]*types.Transaction {
	transactions := make(map[uint64][]*types.Transaction)
	for _, c := range d.getContextChainInOrder() {
		for k, v := range c.transactions {
			transactions[k] = append(transactions[k], v...)
		}
	}
	return transactions
}

func (d *DeployerGroup) getTransactionCount(chain uint64) (*big.Int, error) {
	txs := d.getTransactions()
	return big.NewInt(int64(len(txs[chain]))), nil
}

func (d *DeployerGroup) Enact() (deployment.ChangesetOutput, error) {
	if d.mcmConfig != nil {
		return d.enactMcms()
	}

	return d.enactDeployer()
}

func (d *DeployerGroup) enactMcms() (deployment.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	proposals := make([]mcmslib.TimelockProposal, 0)
	for _, dc := range contexts {
		batches := make([]mcmstypes.BatchOperation, 0)
		for selector, txs := range dc.transactions {
			mcmTransactions := make([]mcmstypes.Transaction, len(txs))
			for i, tx := range txs {
				var err error
				mcmTransactions[i], err = proposalutils.TransactionForChain(selector, tx.To().Hex(), tx.Data(), tx.Value(), "", []string{})
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to build mcms transaction: %w", err)
				}
			}

			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  mcmTransactions,
			})
		}

		if len(batches) == 0 {
			d.e.Logger.Warnf("No batch was produced from deployment context skipping proposal: %s", dc.description)
			continue
		}

		timelocks := BuildTimelockAddressPerChain(d.e, d.state)
		proposerMcms := BuildProposerMcmAddressesPerChain(d.e, d.state)
		inspectors, err := proposalutils.McmsInspectors(d.e)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain: %w", err)
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			d.e.GetContext(),
			timelocks,
			proposerMcms,
			inspectors,
			batches,
			dc.description,
			d.mcmConfig.MinDelay,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal %w", err)
		}

		// Update the proposal metadata to incorporate the startingOpCount
		// from the previous proposal
		if len(proposals) > 0 {
			previousProposal := proposals[len(proposals)-1]
			for chain, metadata := range previousProposal.ChainMetadata {
				nextStartingOp := metadata.StartingOpCount + getBatchCountForChain(chain, proposal)
				proposal.ChainMetadata[chain] = mcmstypes.ChainMetadata{
					StartingOpCount: nextStartingOp,
					MCMAddress:      proposal.ChainMetadata[chain].MCMAddress,
				}
			}
		}

		proposals = append(proposals, *proposal)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: proposals}, nil
}

func getBatchCountForChain(chain mcmstypes.ChainSelector, timelockProposal *mcmslib.TimelockProposal) uint64 {
	batches := make([]mcmstypes.BatchOperation, 0)
	for _, batchOperation := range timelockProposal.Operations {
		if batchOperation.ChainSelector == chain {
			batches = append(batches, batchOperation)
		}
	}

	return uint64(len(batches))
}

func (d *DeployerGroup) enactDeployer() (deployment.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	for _, c := range contexts {
		g := errgroup.Group{}
		for selector, txs := range c.transactions {
			selector, txs := selector, txs
			g.Go(func() error {
				for _, tx := range txs {
					err := d.e.Chains[selector].Client.SendTransaction(context.Background(), tx)
					if err != nil {
						return fmt.Errorf("failed to send transaction: %w", err)
					}
					// TODO how to pass abi here to decode error reason
					_, err = deployment.ConfirmIfNoError(d.e.Chains[selector], tx, err)
					if err != nil {
						return fmt.Errorf("waiting for tx to be mined failed: %w", err)
					}
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return deployment.ChangesetOutput{}, err
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

func BuildTimelockAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]string {
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.Chains {
		addressPerChain[chain.Selector] = state.Chains[chain.Selector].Timelock.Address().Hex()
	}
	return addressPerChain
}

func BuildProposerMcmAddressesPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]string {
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.Chains {
		addressPerChain[chain.Selector] = state.Chains[chain.Selector].ProposerMcm.Address().Hex()
	}
	return addressPerChain
}
