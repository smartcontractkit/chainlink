package proposalutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsaptossdk "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	ccipTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	DefaultValidUntil = 72 * time.Hour
)

type TimelockConfig struct {
	MinDelay                  time.Duration        `json:"minDelay"` // delay for timelock worker to execute the transfers.
	MCMSAction                types.TimelockAction `json:"mcmsAction"`
	OverrideRoot              bool                 `json:"overrideRoot"`                        // if true, override the previous root with the new one.
	TimelockQualifierPerChain map[uint64]string    `json:"timelockQualifierPerChain,omitempty"` // optional qualifier to fetch timelock address from datastore
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
		timelockID, err := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, contractType)
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

func (tc *TimelockConfig) ValidateAptos(chain cldf_aptos.Chain, mcmsAddress aptos.AccountAddress) error {
	if err := tc.validateCommon(); err != nil {
		return err
	}

	if (mcmsAddress == aptos.AccountAddress{}) {
		return fmt.Errorf("aptos MCMS contract not present on chain %s", chain)
	}

	return nil
}

// BuildProposalFromBatchesV2 uses the new MCMS library which replaces the implementation in BuildProposalFromBatches.
func BuildProposalFromBatchesV2(
	e cldf.Environment,
	timelockAddressPerChain map[uint64]string,
	mcmsAddressPerChain map[uint64]string,
	inspectorPerChain map[uint64]mcmssdk.Inspector,
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
			solanaState, err := getSolanaState(env, selector)
			if err != nil {
				return nil, err
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
		case chain_selectors.FamilyAptos:
			// Get role from action
			role, err := GetAptosRoleFromAction(mcmsAction)
			if err != nil {
				return nil, fmt.Errorf("failed to get role from action: %w", err)
			}
			jsonRole, err := json.Marshal(mcmsaptossdk.AdditionalFieldsMetadata{Role: role})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal role for chain %d: %w", selector, err)
			}
			metaDataPerChain[chainID] = types.ChainMetadata{
				StartingOpCount:  opCount,
				MCMAddress:       proposerMcms,
				AdditionalFields: jsonRole,
			}
		}
	}

	return metaDataPerChain, nil
}

func getSolanaState(env cldf.Environment, selector uint64) (*state.MCMSWithTimelockStateSolana, error) {
	solanaChains := env.BlockChains.SolanaChains()
	addresses, err := env.ExistingAddresses.AddressesForChain(selector)
	solanaState, err1 := state.MaybeLoadMCMSWithTimelockChainStateSolana(solanaChains[selector], addresses)
	if err == nil {
		return solanaState, nil
	}

	env.Logger.Info("failed to load MCMSState from address book")
	solanaState, err2 := state.MaybeLoadMCMSWithTimelockChainStateSolanaV2(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(selector)))
	if err2 != nil {
		return nil, fmt.Errorf("failed to load solana state: %w", errors.Join(err1, err2))
	}

	return solanaState, nil
}

// AggregateProposals aggregates multiple MCMS proposals into a single proposal by combining their operations, and
// setting up the proposers and inspectors for each chain. It returns a single MCMS proposal that can be executed
// and signed.
//
// Deprecated: Use extensible AggregateProposalsV2 instead. Which accepts multiple chain families.
func AggregateProposals(
	env cldf.Environment,
	mcmsEVMState map[uint64]state.MCMSWithTimelockState,
	mcmsSolanaState map[uint64]state.MCMSWithTimelockStateSolana,
	proposals []mcmslib.TimelockProposal,
	description string,
	mcmsConfig *TimelockConfig,
) (*mcmslib.TimelockProposal, error) {
	return AggregateProposalsV2(
		env,
		MCMSStates{
			MCMSEVMState:    mcmsEVMState,
			MCMSSolanaState: mcmsSolanaState,
		},
		proposals,
		description,
		mcmsConfig,
	)
}

type MCMSStates struct {
	MCMSEVMState    map[uint64]state.MCMSWithTimelockState
	MCMSSolanaState map[uint64]state.MCMSWithTimelockStateSolana
	MCMSAptosState  map[uint64]aptos.AccountAddress
}

// AggregateProposalsV2 aggregates multiple MCMS proposals into a single proposal by combining their operations, and
// setting up the proposers and inspectors for each chain. It returns a single MCMS proposal that can be executed
// and signed.
// It has an extensible signature to allow for future chain families implementations
func AggregateProposalsV2(
	env cldf.Environment,
	mcmsTimelockStates MCMSStates,
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
		var inspectorOpts []MCMSInspectorOption

		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return nil, fmt.Errorf("failed to get family for chain %d: %w", chainSel, err)
		}
		switch family {
		case chain_selectors.FamilyEVM:
			mcmsEVMState, exists := mcmsTimelockStates.MCMSEVMState[chainSel]
			if !exists {
				return nil, fmt.Errorf("missing MCMS state for chain with selector %d", chainSel)
			}
			mcmsContract, err := mcmsConfig.MCMBasedOnAction(mcmsEVMState)
			if err != nil {
				return &mcmslib.TimelockProposal{}, fmt.Errorf("failed to get MCMS contract for chain with selector %d: %w", chainSel, err)
			}
			timelocks[chainSel] = mcmsEVMState.Timelock.Address().Hex()
			mcmsPerChain[chainSel] = mcmsContract.Address().Hex()
		case chain_selectors.FamilySolana:
			solanaState, existsInSolana := mcmsTimelockStates.MCMSSolanaState[chainSel]
			if !existsInSolana {
				return nil, fmt.Errorf("missing MCMS state for chain with selector %d", chainSel)
			}
			timelocks[chainSel] = mcmssolanasdk.ContractAddress(
				solanaState.TimelockProgram,
				mcmssolanasdk.PDASeed(solanaState.TimelockSeed),
			)
			mcmsAddr, err := mcmsConfig.MCMBasedOnActionSolana(solanaState)
			if err != nil {
				return nil, err
			}
			mcmsPerChain[chainSel] = mcmsAddr
		case chain_selectors.FamilyAptos:
			// Set MCMS addresses. Aptos uses the same address for MCMS and Timelock
			aptosMCMSAddress, existsInAptos := mcmsTimelockStates.MCMSAptosState[chainSel]
			if !existsInAptos {
				return nil, fmt.Errorf("missing MCMS state for chain with selector %d", chainSel)
			}
			timelocks[chainSel] = aptosMCMSAddress.StringLong()
			mcmsPerChain[chainSel] = aptosMCMSAddress.StringLong()
			// Getting inspector parameters for Aptos
			role, err := GetAptosRoleFromAction(mcmsConfig.MCMSAction)
			if err != nil {
				return nil, fmt.Errorf("failed to get role from action: %w", err)
			}
			inspectorOpts = append(inspectorOpts, WithAptosRole(role))
		}

		inspectors[chainSel], err = McmsInspectorForChain(env, chainSel, inspectorOpts...)
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
