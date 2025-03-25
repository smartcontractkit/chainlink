package mcmsutil

import (
	"fmt"
	"time"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// CreateMCMSProposal creates a new MCMS proposal with prepared (but not sent) transactions.
func CreateMCMSProposal(e deployment.Environment, preparedTxs []*txutil.PreparedTx, mcmsMinDelay time.Duration, proposalName string) (*mcmslib.TimelockProposal, error) {
	var chainSelectors []uint64
	for _, tx := range preparedTxs {
		chainSelectors = append(chainSelectors, tx.ChainSelector)
	}
	mcmsStatePerChain, err := commonchangeset.MaybeLoadMCMSWithTimelockState(e, chainSelectors)
	if err != nil {
		return nil, err
	}
	inspectors, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return nil, err
	}

	// Get MCMS state for each chain
	timelockAddressesPerChain := map[uint64]string{}
	proposerMcmsPerChain := map[uint64]string{}
	for _, chainSelector := range chainSelectors {
		state := mcmsStatePerChain[chainSelector]
		timelockAddressesPerChain[chainSelector] = state.Timelock.Address().Hex()
		proposerMcmsPerChain[chainSelector] = state.ProposerMcm.Address().Hex()
	}

	// Create batch operations from generated transactions
	var batches []mcmstypes.BatchOperation
	for _, tx := range preparedTxs {
		batchOp, err := proposalutils.BatchOperationForChain(
			tx.ChainSelector,
			tx.Tx.To().Hex(),
			tx.Tx.Data(),
			tx.Tx.Value(),
			tx.ContractType,
			tx.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create batch operation: %w", err)
		}
		batches = append(batches, batchOp)
	}

	return proposalutils.BuildProposalFromBatchesV2(
		e,
		timelockAddressesPerChain,
		proposerMcmsPerChain,
		inspectors,
		batches,
		proposalName,
		mcmsMinDelay,
	)
}

// ExecuteOrPropose executes the transactions if no MCMS is configured, otherwise creates a proposal.
func ExecuteOrPropose(
	e deployment.Environment,
	txs []*txutil.PreparedTx,
	mcmsCfg *changeset.MCMSConfig,
	proposalName string,
) (deployment.ChangesetOutput, error) {
	if len(txs) == 0 {
		return deployment.ChangesetOutput{}, nil
	}

	if mcmsCfg != nil {
		proposal, err := CreateMCMSProposal(e, txs, mcmsCfg.MinDelay, proposalName)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error creating MCMS proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	_, err := txutil.SignAndExecute(e, txs)
	return deployment.ChangesetOutput{}, err
}
