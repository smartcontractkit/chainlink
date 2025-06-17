package proposalutils

import (
	"errors"
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

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

func (tc *TimelockConfig) Validate(chain cldf_evm.Chain, s state.MCMSWithTimelockState) error {
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

func (tc *TimelockConfig) ValidateSolana(e cldf.Environment, chainSelector uint64) error {
	err := tc.validateCommon()
	if err != nil {
		return err
	}

	validateContract := func(contractType cldf.ContractType) error {
		timelockID, err := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, contractType) //nolint:staticcheck // Uncomment above once datastore is updated to contains addresses
		if err != nil {
			return fmt.Errorf("%s not present on the chain %w", contractType, err)
		}
		// Make sure addresses are correctly parsed. Format is: "programID.PDASeed"
		_, _, err = mcmssolanasdk.ParseContractAddress(timelockID)
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

// BuildProposalFromBatchesV2 uses the new MCMS library which replaces the implementation in BuildProposalFromBatches.
func BuildProposalFromBatchesV2(
	e cldf.Environment,
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
	env cldf.Environment,
	chainSelectors []uint64,
	inspectorPerChain map[uint64]mcmssdk.Inspector,
	mcmsPerChain map[uint64]string, // can be proposer, canceller or bypasser
	mcmsAction types.TimelockAction,
) (map[types.ChainSelector]types.ChainMetadata, error) {
	solanaChains := env.BlockChains.SolanaChains()
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
			solanaState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solanaChains[selector], addresses)
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

// AggregateProposals aggregates multiple MCMS proposals into a single proposal by combining their operations, and
// setting up the proposers and inspectors for each chain. It returns a single MCMS proposal that can be executed
// and signed.
func AggregateProposals(
	env cldf.Environment,
	mcmsEVMState map[uint64]state.MCMSWithTimelockState,
	mcmsSolanaState map[uint64]state.MCMSWithTimelockStateSolana,
	proposals []mcmslib.TimelockProposal,
	description string,
	mcmsConfig *TimelockConfig,
) (*mcmslib.TimelockProposal, error) {
	if mcmsConfig == nil {
		return nil, nil
	}

	var batches []types.BatchOperation

	// Add proposals to the aggregate.
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
		var err error
		if _, exists := mcmsEVMState[chainSel]; exists {
			mcmsContract, err := mcmsConfig.MCMBasedOnAction(mcmsEVMState[chainSel])
			if err != nil {
				return &mcmslib.TimelockProposal{}, fmt.Errorf("failed to get MCMS contract for chain with selector %d: %w", chainSel, err)
			}
			timelocks[chainSel] = mcmsEVMState[chainSel].Timelock.Address().Hex()
			mcmsPerChain[chainSel] = mcmsContract.Address().Hex()
		} else if mcmsSolanaState == nil {
			return nil, fmt.Errorf("missing MCMS state for chain with selector %d", chainSel)
		} else if solanaState, existsInSolana := mcmsSolanaState[chainSel]; existsInSolana {
			timelocks[chainSel] = mcmssolanasdk.ContractAddress(
				solanaState.TimelockProgram,
				mcmssolanasdk.PDASeed(solanaState.TimelockSeed),
			)
			mcmsAddr, err := mcmsConfig.MCMBasedOnActionSolana(solanaState)
			if err != nil {
				return nil, err
			}
			mcmsPerChain[chainSel] = mcmsAddr
		} else {
			return nil, fmt.Errorf("missing MCMS state for chain with selector %d", chainSel)
		}

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
