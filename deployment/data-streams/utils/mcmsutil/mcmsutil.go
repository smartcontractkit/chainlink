package mcmsutil

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
	"github.com/smartcontractkit/mcms"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
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
		proposalutils.TimelockConfig{MinDelay: mcmsMinDelay},
	)
}

// ExecuteOrPropose executes the transactions if no MCMS is configured, otherwise creates a proposal.
func ExecuteOrPropose(
	e deployment.Environment,
	txs []*txutil.PreparedTx,
	mcmsCfg *dsTypes.MCMSConfig,
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

// TransferToMCMSWithTimelockForTypeAndVersion transfers ownership of the contracts of a specific type and version to the
// MCMS timelock on that chain. The output will contain an MCMS timelock proposal for "AcceptOwnership" of those contracts
// The address book should be recently deployed addresses that are being transferred to MCMS and should not be in e.ExistingAddresses
func TransferToMCMSWithTimelockForTypeAndVersion(e deployment.Environment,
	ab deployment.AddressBook, filter deployment.TypeAndVersion,
	MCMSConfig proposalutils.TimelockConfig) (deployment.ChangesetOutput, error) {

	contractAddresses := make(map[uint64][]common.Address)
	for _, chain := range e.Chains {
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get addresses from address book: %w", err)
		}
		for address, typeAndVersion := range chainAddresses {
			if typeAndVersion.Type == filter.Type && typeAndVersion.Version == filter.Version {
				contractAddresses[chain.Selector] = append(contractAddresses[chain.Selector], common.HexToAddress(address))
			}
		}
	}

	// create a merged addressbook with the existing + new addresses. Sub-changesets will need all addresses
	// This is required when chaining together changesets
	existingAddresses := e.ExistingAddresses
	abTemp := deployment.NewMemoryAddressBook()
	if err := abTemp.Merge(ab); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed merging new addresses into temp addresses: %w", err)
	}
	if err := abTemp.Merge(e.ExistingAddresses); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed merging existing addresses into temp addresses: %w", err)
	}

	e.ExistingAddresses = abTemp

	transferCs := deployment.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2)
	transferCfg := commonchangeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contractAddresses,
		MCMSConfig:       MCMSConfig,
	}

	transferOut, err := transferCs.Apply(e, transferCfg)
	e.ExistingAddresses = existingAddresses // reset the address book to the original state regardless of errors
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer contracts to MCMS: %w", err)
	}

	return deployment.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: transferOut.MCMSTimelockProposals,
	}, nil
}

// MergeSimilarTimelockProposals merges multiple MCMS timelock proposals into a single proposal.
// It assumes that all proposals have the same action, delay, and timelock addresses, etc... just different operations.
// This is useful for combining multiple proposals into 1 when chaining together changesets.
func MergeSimilarTimelockProposals(proposals []mcms.TimelockProposal) (mcms.TimelockProposal, error) {
	var newProposal mcms.TimelockProposal

	if len(proposals) >= 1 {
		// we make an assumption that all proposals have these common settings
		newProposal = mcms.TimelockProposal{
			BaseProposal:      proposals[0].BaseProposal,
			Action:            proposals[0].Action,
			Delay:             proposals[0].Delay,
			TimelockAddresses: proposals[0].TimelockAddresses,
			SaltOverride:      proposals[0].SaltOverride,
		}
	}

	for _, prop := range proposals {
		newProposal.Operations = append(newProposal.Operations, prop.Operations...)
	}

	return newProposal, nil
}
