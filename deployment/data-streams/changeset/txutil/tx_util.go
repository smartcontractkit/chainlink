package txutil

import (
	"fmt"
	"math/big"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// TransactionGenerator defines how to generate a transaction for a specific chain
type TransactionGenerator interface {
	GenerateTransactions(e deployment.Environment, chainSelector uint64, mcms bool) ([]*gethtypes.Transaction, string, error)
}

// TransactionHandler handles executing transactions either directly or via MCMS proposal
type TransactionHandler struct {
	Environment    deployment.Environment
	ChainSelectors []uint64
	ContractType   string
	ProposalName   string
	MCMSConfig     *changeset.MCMSConfig
}

// HandleTransactions processes transactions either directly or through MCMS proposal
func (h *TransactionHandler) HandleTransactions(txGen TransactionGenerator) (deployment.ChangesetOutput, error) {
	if h.MCMSConfig != nil {
		return h.createMCMSProposal(txGen)
	}
	return h.executeDirectTransactions(txGen)
}

// executeDirectTransactions executes transactions directly on the chain with the given deployer address
func (h *TransactionHandler) executeDirectTransactions(txGen TransactionGenerator) (deployment.ChangesetOutput, error) {
	for _, chainSelector := range h.ChainSelectors {
		chain := h.Environment.Chains[chainSelector]

		txs, _, err := txGen.GenerateTransactions(h.Environment, chainSelector, false)

		for _, tx := range txs {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			if _, err := chain.Confirm(tx); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}

	}
	return deployment.ChangesetOutput{}, nil
}

func (h *TransactionHandler) createMCMSProposal(txGen TransactionGenerator) (deployment.ChangesetOutput, error) {
	mcmsStatePerChain, err := commonchangeset.MaybeLoadMCMSWithTimelockState(h.Environment, h.ChainSelectors)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	inspectors, err := proposalutils.McmsInspectors(h.Environment)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var batches []mcmstypes.BatchOperation
	timelockAddressesPerChain := map[uint64]string{}
	proposerMcmsPerChain := map[uint64]string{}

	for _, chainSelector := range h.ChainSelectors {
		txs, toAddress, err := txGen.GenerateTransactions(h.Environment, chainSelector, true)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		for _, tx := range txs {
			batchOp, err := proposalutils.BatchOperationForChain(
				chainSelector,
				toAddress,
				tx.Data(),
				big.NewInt(0),
				h.ContractType,
				[]string{},
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create batch operation: %w", err)
			}

			batches = append(batches, batchOp)
		}
		state := mcmsStatePerChain[chainSelector]
		timelockAddressesPerChain[chainSelector] = state.Timelock.Address().Hex()
		proposerMcmsPerChain[chainSelector] = state.ProposerMcm.Address().Hex()
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		h.Environment,
		timelockAddressesPerChain,
		proposerMcmsPerChain,
		inspectors,
		batches,
		h.ProposalName,
		h.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
	}, nil
}
