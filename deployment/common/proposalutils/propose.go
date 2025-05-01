package proposalutils

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	ccipTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	DefaultValidUntil = 72 * time.Hour
)

type TimelockConfig struct {
	MinDelay     time.Duration        `json:"minDelay"` // delay for timelock worker to execute the transfers.
	MCMSAction   types.TimelockAction `json:"mcmsAction"`
	OverrideRoot bool                 `json:"overrideRoot"` // if true, override the previous root with the new one.
}

func (tc *TimelockConfig) MCMBasedOnActionSolana(s state.MCMSWithTimelockStateSolana) (string, error) {
	// if MCMSAction is not set, default to timelock.Schedule, this is to ensure no breaking changes for existing code
	if tc.MCMSAction == "" {
		tc.MCMSAction = types.TimelockActionSchedule
	}
	switch tc.MCMSAction {
	case types.TimelockActionSchedule:
		contractID := mcmssolanasdk.ContractAddress(s.McmProgram, mcmssolanasdk.PDASeed(s.ProposerMcmSeed))
		return contractID, nil
	case types.TimelockActionCancel:
		contractID := mcmssolanasdk.ContractAddress(s.McmProgram, mcmssolanasdk.PDASeed(s.CancellerMcmSeed))
		return contractID, nil
	case types.TimelockActionBypass:
		contractID := mcmssolanasdk.ContractAddress(s.McmProgram, mcmssolanasdk.PDASeed(s.BypasserMcmSeed))
		return contractID, nil
	default:
		return "", errors.New("invalid MCMS action")
	}
}

func (tc *TimelockConfig) MCMBasedOnAction(s state.MCMSWithTimelockState) (*gethwrappers.ManyChainMultiSig, error) {
	// if MCMSAction is not set, default to timelock.Schedule, this is to ensure no breaking changes for existing code
	if tc.MCMSAction == "" {
		tc.MCMSAction = types.TimelockActionSchedule
	}
	switch tc.MCMSAction {
	case types.TimelockActionSchedule:
		if s.ProposerMcm == nil {
			return nil, errors.New("missing proposerMcm")
		}
		return s.ProposerMcm, nil
	case types.TimelockActionCancel:
		if s.CancellerMcm == nil {
			return nil, errors.New("missing cancellerMcm")
		}
		return s.CancellerMcm, nil
	case types.TimelockActionBypass:
		if s.BypasserMcm == nil {
			return nil, errors.New("missing bypasserMcm")
		}
		return s.BypasserMcm, nil
	default:
		return nil, errors.New("invalid MCMS action")
	}
}

func (tc *TimelockConfig) validateCommon() error {
	// if MCMSAction is not set, default to timelock.Schedule
	if tc.MCMSAction == "" {
		tc.MCMSAction = types.TimelockActionSchedule
	}
	if tc.MCMSAction != types.TimelockActionSchedule &&
		tc.MCMSAction != types.TimelockActionCancel &&
		tc.MCMSAction != types.TimelockActionBypass {
		return fmt.Errorf("invalid MCMS type %s", tc.MCMSAction)
	}
	return nil
}

func (tc *TimelockConfig) Validate(chain deployment.Chain, s state.MCMSWithTimelockState) error {
	err := tc.validateCommon()
	if err != nil {
		return err
	}
	if s.Timelock == nil {
		return fmt.Errorf("missing timelock on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionSchedule && s.ProposerMcm == nil {
		return fmt.Errorf("missing proposerMcm on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionCancel && s.CancellerMcm == nil {
		return fmt.Errorf("missing cancellerMcm on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionBypass && s.BypasserMcm == nil {
		return fmt.Errorf("missing bypasserMcm on %s", chain)
	}
	if s.Timelock == nil {
		return fmt.Errorf("missing timelock on %s", chain)
	}
	if s.CallProxy == nil {
		return fmt.Errorf("missing callProxy on %s", chain)
	}
	return nil
}

func (tc *TimelockConfig) ValidateSolana(e deployment.Environment, chainSelector uint64) error {
	err := tc.validateCommon()
	if err != nil {
		return err
	}

	validateContract := func(contractType deployment.ContractType) error {
		timelockID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, contractType) //nolint:staticcheck // Uncomment above once datastore is updated to contains addresses
		if err != nil {
			return fmt.Errorf("%s not present on the chain %w", contractType, err)
		}
		// Make sure addresses are correctly parsed. Format is: "programID.PDASeed"
		_, _, err = mcmsSolana.ParseContractAddress(timelockID)
		if err != nil {
			return fmt.Errorf("failed to parse timelock address: %w", err)
		}
		return nil
	}

	err = validateContract(ccipTypes.RBACTimelock)
	if err != nil {
		return err
	}

	switch tc.MCMSAction {
	case types.TimelockActionSchedule:
		err = validateContract(ccipTypes.ProposerManyChainMultisig)
		if err != nil {
			return err
		}
	case types.TimelockActionCancel:
		err = validateContract(ccipTypes.CancellerManyChainMultisig)
		if err != nil {
			return err
		}
	case types.TimelockActionBypass:
		err = validateContract(ccipTypes.BypasserManyChainMultisig)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid MCMS action %s", tc.MCMSAction)
	}

	return nil
}

func BuildProposalMetadata(
	chainSelectors []uint64,
	proposerMcmsesPerChain map[uint64]*gethwrappers.ManyChainMultiSig,
) (map[mcms.ChainIdentifier]mcms.ChainMetadata, error) {
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	for _, selector := range chainSelectors {
		proposerMcms, ok := proposerMcmsesPerChain[selector]
		if !ok {
			return nil, fmt.Errorf("missing proposer mcm for chain %d", selector)
		}
		chainId := mcms.ChainIdentifier(selector)
		opCount, err := proposerMcms.GetOpCount(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get op count for chain %d: %w", selector, err)
		}
		metaDataPerChain[chainId] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      proposerMcms.Address(),
		}
	}
	return metaDataPerChain, nil
}

// BuildProposalFromBatches Given batches of operations, we build the metadata and timelock addresses of those opartions
// We then return a proposal that can be executed and signed.
// You can specify multiple batches for the same chain, but the only
// usecase to do that would be you have a batch that can't fit in a single
// transaction due to gas or calldata constraints of the chain.
// The batches are specified separately because we eventually intend
// to support user-specified cross chain ordering of batch execution by the tooling itself.
// TODO: Can/should merge timelocks and proposers into a single map for the chain.
// Deprecated: Use BuildProposalFromBatchesV2 instead.
func BuildProposalFromBatches(
	timelocksPerChain map[uint64]common.Address,
	proposerMcmsesPerChain map[uint64]*gethwrappers.ManyChainMultiSig,
	batches []timelock.BatchChainOperation,
	description string,
	minDelay time.Duration,
) (*timelock.MCMSWithTimelockProposal, error) {
	if len(batches) == 0 {
		return nil, errors.New("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainIdentifier))
	}

	mcmsMd, err := BuildProposalMetadata(chains.ToSlice(), proposerMcmsesPerChain)
	if err != nil {
		return nil, err
	}

	tlsPerChainId := make(map[mcms.ChainIdentifier]common.Address)
	for chainId, tl := range timelocksPerChain {
		tlsPerChainId[mcms.ChainIdentifier(chainId)] = tl
	}
	validUntil := time.Now().Unix() + int64(DefaultValidUntil.Seconds())
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		uint32(validUntil),
		[]mcms.Signature{},
		false,
		mcmsMd,
		tlsPerChainId,
		description,
		batches,
		timelock.Schedule,
		minDelay.String(),
	)
}

// BuildProposalFromBatchesV2 uses the new MCMS library which replaces the implementation in BuildProposalFromBatches.
func BuildProposalFromBatchesV2(
	e deployment.Environment,
	timelockAddressPerChain map[uint64]string,
	mcmsAddressPerChain map[uint64]string, inspectorPerChain map[uint64]mcmssdk.Inspector,
	batches []types.BatchOperation,
	description string,
	mcmsCfg TimelockConfig,
) (*mcmslib.TimelockProposal, error) {
	// default to schedule if not set, this is to be consistent with the old implementation
	// and to avoid breaking changes
	if mcmsCfg.MCMSAction == "" {
		mcmsCfg.MCMSAction = types.TimelockActionSchedule
	}
	if len(batches) == 0 {
		return nil, errors.New("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainSelector))
	}
	tlsPerChainID := make(map[types.ChainSelector]string)
	for chainID, tl := range timelockAddressPerChain {
		tlsPerChainID[types.ChainSelector(chainID)] = tl
	}
	mcmsMd, err := buildProposalMetadataV2(e, chains.ToSlice(), inspectorPerChain, mcmsAddressPerChain, mcmsCfg.MCMSAction)
	if err != nil {
		return nil, err
	}

	validUntil := time.Now().Unix() + int64(DefaultValidUntil.Seconds())

	builder := mcmslib.NewTimelockProposalBuilder()
	builder.
		SetVersion("v1").
		SetAction(mcmsCfg.MCMSAction).
		//nolint:gosec // G115
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		SetDelay(types.NewDuration(mcmsCfg.MinDelay)).
		SetOverridePreviousRoot(mcmsCfg.OverrideRoot).
		SetChainMetadata(mcmsMd).
		SetTimelockAddresses(tlsPerChainID).
		SetOperations(batches)

	build, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return build, nil
}

func buildProposalMetadataV2(
	env deployment.Environment,
	chainSelectors []uint64,
	inspectorPerChain map[uint64]mcmssdk.Inspector,
	mcmsPerChain map[uint64]string, // can be proposer, canceller or bypasser
	mcmsAction types.TimelockAction,
) (map[types.ChainSelector]types.ChainMetadata, error) {
	metaDataPerChain := make(map[types.ChainSelector]types.ChainMetadata)
	for _, selector := range chainSelectors {
		proposerMcms, ok := mcmsPerChain[selector]
		if !ok {
			return nil, fmt.Errorf("missing proposer mcm for chain %d", selector)
		}
		chainID := types.ChainSelector(selector)
		opCount, err := inspectorPerChain[selector].GetOpCount(env.GetContext(), proposerMcms)
		if err != nil {
			return nil, fmt.Errorf("failed to get op count for chain %d: %w", selector, err)
		}
		family, err := chain_selectors.GetSelectorFamily(selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get family for chain %d: %w", selector, err)
		}
		switch family {
		case chain_selectors.FamilyEVM:
			metaDataPerChain[chainID] = types.ChainMetadata{
				StartingOpCount: opCount,
				MCMAddress:      proposerMcms,
			}
		case chain_selectors.FamilySolana:
			addresses, err := env.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return nil, fmt.Errorf("failed to load addresses for chain %d: %w", selector, err)
			}
			solanaState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(env.SolChains[selector], addresses)
			if err != nil {
				return nil, fmt.Errorf("failed to load solana state: %w", err)
			}

			var instanceSeed mcmssolanasdk.PDASeed
			switch mcmsAction {
			case types.TimelockActionSchedule:
				instanceSeed = mcmssolanasdk.PDASeed(solanaState.ProposerMcmSeed)
			case types.TimelockActionCancel:
				instanceSeed = mcmssolanasdk.PDASeed(solanaState.CancellerMcmSeed)
			case types.TimelockActionBypass:
				instanceSeed = mcmssolanasdk.PDASeed(solanaState.BypasserMcmSeed)
			default:
				return nil, fmt.Errorf("invalid MCMS action %s", mcmsAction)
			}

			metaDataPerChain[chainID], err = mcmssolanasdk.NewChainMetadata(
				opCount,
				solanaState.McmProgram,
				instanceSeed,
				solanaState.ProposerAccessControllerAccount,
				solanaState.CancellerAccessControllerAccount,
				solanaState.BypasserAccessControllerAccount)
			if err != nil {
				return nil, fmt.Errorf("failed to create chain metadata: %w", err)
			}
		}
	}

	return metaDataPerChain, nil
}

// AggregateProposals aggregates proposals from the legacy and new formats into a single proposal.
// Required if you are merging multiple changesets that have different proposal formats.
func AggregateProposals(
	env deployment.Environment,
	mcmsState map[uint64]state.MCMSWithTimelockState,
	proposals []mcmslib.TimelockProposal,
	legacyProposals []timelock.MCMSWithTimelockProposal,
	description string,
	mcmsConfig *TimelockConfig,
) (*mcmslib.TimelockProposal, error) {
	if mcmsConfig == nil {
		return nil, nil
	}

	var batches []types.BatchOperation
	// Add proposals that follow the legacy format to the aggregate.
	for _, proposal := range legacyProposals {
		for _, batchTransaction := range proposal.Transactions {
			for _, transaction := range batchTransaction.Batch {
				batchOperation, err := BatchOperationForChain(
					uint64(batchTransaction.ChainIdentifier),
					transaction.To.Hex(),
					transaction.Data,
					big.NewInt(0),
					transaction.ContractType,
					transaction.Tags,
				)
				if err != nil {
					return &mcmslib.TimelockProposal{}, fmt.Errorf("failed to create batch operation on chain with selector %d: %w", batchTransaction.ChainIdentifier, err)
				}
				batches = append(batches, batchOperation)
			}
		}
	}
	// Add proposals that follow the new format to the aggregate.
	for _, proposal := range proposals {
		batches = append(batches, proposal.Operations...)
	}

	// Return early if there are no operations.
	if len(batches) == 0 {
		return nil, nil
	}

	// Store the timelocks, proposers, and inspectors for each chain.
	timelocks := make(map[uint64]string)
	mcmsPerChain := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	for _, op := range batches {
		chainSel := uint64(op.ChainSelector)
		mcmsContract, err := mcmsConfig.MCMBasedOnAction(mcmsState[chainSel])
		if err != nil {
			return &mcmslib.TimelockProposal{}, fmt.Errorf("failed to get MCMS contract for chain with selector %d: %w", chainSel, err)
		}
		timelocks[chainSel] = mcmsState[chainSel].Timelock.Address().Hex()
		mcmsPerChain[chainSel] = mcmsContract.Address().Hex()
		inspectors[chainSel], err = McmsInspectorForChain(env, chainSel)
		if err != nil {
			return &mcmslib.TimelockProposal{}, fmt.Errorf("failed to get MCMS inspector for chain with selector %d: %w", chainSel, err)
		}
	}

	return BuildProposalFromBatchesV2(
		env,
		timelocks,
		mcmsPerChain,
		inspectors,
		batches,
		description,
		*mcmsConfig,
	)
}
