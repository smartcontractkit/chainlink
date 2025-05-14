package mcmsutil

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// CreateMCMSProposal creates a new MCMS proposal with prepared (but not sent) transactions.
func CreateMCMSProposal(e cldf.Environment, preparedTxs []*txutil.PreparedTx, mcmsMinDelay time.Duration, proposalName string) (*mcmslib.TimelockProposal, error) {
	var chainSelectors []uint64
	for _, tx := range preparedTxs {
		chainSelectors = append(chainSelectors, tx.ChainSelector)
	}

	mcmsStatePerChain, err := commonstate.MaybeLoadMCMSWithTimelockStateDataStore(e, chainSelectors)
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
	e cldf.Environment,
	txs []*txutil.PreparedTx,
	mcmsCfg *dsTypes.MCMSConfig,
	proposalName string,
) (cldf.ChangesetOutput, error) {
	if len(txs) == 0 {
		return cldf.ChangesetOutput{}, nil
	}

	if mcmsCfg != nil {
		proposal, err := CreateMCMSProposal(e, txs, mcmsCfg.MinDelay, proposalName)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error creating MCMS proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	_, err := txutil.SignAndExecute(e, txs)
	return cldf.ChangesetOutput{}, err
}

// TransferToMCMSWithTimelock transfers ownership of specified addresses to MCMS timelock contracts.
// The output will contain an MCMS timelock proposal for "AcceptOwnership" of those contracts
// The newAddresses should be recently deployed addresses that are being transferred to MCMS and
// should not be in `e` Environment
func TransferToMCMSWithTimelock(e cldf.Environment, newAddresses []ds.AddressRef,
	mcmsConfig proposalutils.TimelockConfig) (cldf.ChangesetOutput, error) {
	// Map: chainselector -> List[Address]
	contractAddressesToTransfer := make(map[uint64][]common.Address)

	for _, addressRef := range newAddresses {
		chainSelector := addressRef.ChainSelector
		contractAddressesToTransfer[chainSelector] = append(contractAddressesToTransfer[chainSelector], common.HexToAddress(addressRef.Address))
	}

	// Adapter: Convert from DataStore primitives -> AddressBook is needed for TransferToMCMSWithTimelockV2 changeset
	// This should be removed once TransferToMCMSWithTimelockV2 is updated to use the DataStore or there is a new changeset which does
	newAndExistingAddresses, err := utils.AddressRefsToAddressBook(newAddresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert new addresses to address book: %w", err)
	}
	// create a merged addressbook with the existing + new addresses. Sub-changesets will need all addresses
	// This is required when chaining together changesets
	// i.e. the MCMS timelock addresses
	existingAddrs, err := utils.DataStoreToAddressBook(e.DataStore)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert existing address book: %w", err)
	}
	if err := newAndExistingAddresses.Merge(existingAddrs); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed merging existing addresses into temp addresses: %w", err)
	}
	e.ExistingAddresses = newAndExistingAddresses

	transferCs := cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2)
	transferCfg := commonchangeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contractAddressesToTransfer,
		MCMSConfig:       mcmsConfig,
	}

	transferOut, err := transferCs.Apply(e, transferCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer contracts to MCMS: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: transferOut.MCMSTimelockProposals,
	}, nil
}

// MergeSimilarTimelockProposals merges multiple MCMS timelock proposals into a single proposal.
// It assumes that all proposals have the same action, delay, and timelock addresses, etc... just different operations.
// This is useful for combining multiple proposals into 1 when chaining together changesets.
func MergeSimilarTimelockProposals(proposals []mcms.TimelockProposal) (mcms.TimelockProposal, error) {
	var mergedProposal mcms.TimelockProposal
	if len(proposals) == 0 {
		// not considered an error, just no proposals to merge
		return mcms.TimelockProposal{}, nil
	}

	if len(proposals) == 1 {
		return proposals[0], nil
	}

	// Initialize merged proposal from first
	// we make an assumption that all proposals have these common settings and check for equality later
	mergedProposal = mcms.TimelockProposal{
		BaseProposal:      proposals[0].BaseProposal,
		Action:            proposals[0].Action,
		Delay:             proposals[0].Delay,
		TimelockAddresses: proposals[0].TimelockAddresses,
		SaltOverride:      proposals[0].SaltOverride,
	}

	for i, prop := range proposals {
		ok, err := proposalsEqualForMerge(proposals[0], prop)
		if !ok {
			return mcms.TimelockProposal{}, fmt.Errorf("cannot merge proposal index %d didn't match proposal at index 0: %w", i, err)
		}
		mergedProposal.Operations = append(mergedProposal.Operations, prop.Operations...)
	}

	return mergedProposal, nil
}

// proposalsEqualForMerge tests equality of two proposals to see if they can be merged.
// This assumes that the proposals are similar enough to be merged within the context of MergeSimilarTimelockProposals.
// Shouldn't be used outside the context of MergeSimilarTimelockProposals as the equality check is not exhaustive
// and only supports the merge capabilities of MergeSimilarTimelockProposals
func proposalsEqualForMerge(p1, p2 mcms.TimelockProposal) (bool, error) {
	if p1.Version != p2.Version {
		return false, fmt.Errorf("mismatched Version: %s vs %s", p1.Version, p2.Version)
	}
	if p1.Kind != p2.Kind {
		return false, fmt.Errorf("mismatched Kind: %s vs %s", p1.Kind, p2.Kind)
	}
	if p1.OverridePreviousRoot != p2.OverridePreviousRoot {
		return false, fmt.Errorf("mismatched OverridePreviousRoot: %v vs %v", p1.OverridePreviousRoot, p2.OverridePreviousRoot)
	}
	if p1.Action != p2.Action {
		return false, fmt.Errorf("mismatched Action: %s vs %s", p1.Action, p2.Action)
	}
	if p1.Delay != p2.Delay {
		return false, fmt.Errorf("mismatched Delay: %v vs %v", p1.Delay, p2.Delay)
	}

	// timestamps should be sufficiently close.
	t1 := time.Unix(int64(p1.BaseProposal.ValidUntil), 0)
	t2 := time.Unix(int64(p2.BaseProposal.ValidUntil), 0)
	timeDiff := t1.Sub(t2).Abs()
	if timeDiff > time.Minute {
		return false, fmt.Errorf(
			"timestamps differ too much between proposals (%d vs %d)",
			p1.BaseProposal.ValidUntil, p2.BaseProposal.ValidUntil,
		)
	}
	// Set to the same value to pass equality check
	p1.BaseProposal.ValidUntil = p2.BaseProposal.ValidUntil
	if !reflect.DeepEqual(p1.BaseProposal, p2.BaseProposal) {
		return false, errors.New("mismatched BaseProposal")
	}
	if !reflect.DeepEqual(p1.ChainMetadata, p2.ChainMetadata) {
		return false, errors.New("mismatched ChainMetadata")
	}
	if !reflect.DeepEqual(p1.TimelockAddresses, p2.TimelockAddresses) {
		return false, errors.New("mismatched TimelockAddresses")
	}

	return true, nil
}

type HasOwnershipConfig interface {
	GetOwnershipConfig() dsTypes.OwnershipSettings
}

func GetTransferOwnershipProposals[T HasOwnershipConfig](
	e cldf.Environment, cfg T, addresses []ds.AddressRef) ([]mcms.TimelockProposal, error) {
	var proposals []mcms.TimelockProposal
	if cfg.GetOwnershipConfig().ShouldTransfer && cfg.GetOwnershipConfig().MCMSProposalConfig != nil {
		res, err := TransferToMCMSWithTimelock(e, addresses, *cfg.GetOwnershipConfig().MCMSProposalConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
		}
		proposals = res.MCMSTimelockProposals
	}
	return proposals, nil
}
