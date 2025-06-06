package v1_6

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ cldf.ChangeSet[RMNCurseConfig] = RMNCurseChangeset
	_ cldf.ChangeSet[RMNCurseConfig] = RMNUncurseChangeset
)

// RMNCurseAction represent a curse action to be applied on a chain (ChainSelector) with a specific subject (SubjectToCurse)
// The curse action will by applied by calling the Curse method on the RMNRemote contract on the chain (ChainSelector)
type RMNCurseAction struct {
	ChainSelector  uint64
	SubjectToCurse globals.Subject
}

// CurseAction is a function that returns a list of RMNCurseAction to be applied on a chain
// CurseChain, CurseLane, CurseGloballyOnlyOnSource are examples of function implementing CurseAction
type CurseAction func(e cldf.Environment) ([]RMNCurseAction, error)

type RMNCurseConfig struct {
	MCMS         *proposalutils.TimelockConfig
	CurseActions []CurseAction
	// Use this if you need to include lanes that are not in sourcechain in the offramp. i.e. not yet migrated lane from 1.5
	IncludeNotConnectedLanes bool
	// Use this if you want to include curse subject even when they are already cursed (CurseChangeset) or already uncursed (UncurseChangeset)
	Force  bool
	Reason string
}

func (c RMNCurseConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	err = state.EnforceMCMSUsageIfProd(e.GetContext(), c.MCMS)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	if len(c.CurseActions) == 0 {
		return errors.New("curse actions are required")
	}

	if c.Reason == "" {
		return errors.New("reason is required")
	}

	validEVMSubjects := map[globals.Subject]struct{}{
		globals.GlobalCurseSubject(): {},
	}

	validSolanaSubjects := map[globals.Subject]struct{}{
		globals.GlobalCurseSubject(): {},
	}

	for _, selector := range GetAllCursableChainsSelector(e) {
		validEVMSubjects[globals.FamilyAwareSelectorToSubject(selector, chain_selectors.FamilyEVM)] = struct{}{}
		validSolanaSubjects[globals.FamilyAwareSelectorToSubject(selector, chain_selectors.FamilySolana)] = struct{}{}
	}

	for _, curseAction := range c.CurseActions {
		result, err := curseAction(e)
		if err != nil {
			return fmt.Errorf("failed to generate curse actions: %w", err)
		}

		for _, action := range result {
			if err = cldf.IsValidChainSelector(action.ChainSelector); err != nil {
				return fmt.Errorf("invalid chain selector %d", action.ChainSelector)
			}

			family, err := chain_selectors.GetSelectorFamily(action.ChainSelector)
			if err != nil {
				return err
			}

			// TODO: Implement chain family agnostic validation
			switch family {
			case chain_selectors.FamilyEVM:
				if _, ok := validEVMSubjects[action.SubjectToCurse]; !ok {
					return fmt.Errorf("invalid subject %x", action.SubjectToCurse)
				}

				targetChain := e.BlockChains.EVMChains()[action.ChainSelector]
				targetChainState, ok := state.Chains[action.ChainSelector]
				if !ok {
					return fmt.Errorf("chain %s not found in onchain state", targetChain.String())
				}

				if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, targetChain.DeployerKey.From, targetChainState.Timelock.Address(), targetChainState.RMNRemote); err != nil {
					return fmt.Errorf("chain %s: %w", targetChain.String(), err)
				}
			case chain_selectors.FamilySolana:
				if _, ok := validSolanaSubjects[action.SubjectToCurse]; !ok {
					return fmt.Errorf("invalid subject %x", action.SubjectToCurse)
				}

				targetChain := e.BlockChains.SolanaChains()[action.ChainSelector]
				targetChainState, ok := state.SolChains[action.ChainSelector]
				if !ok {
					return fmt.Errorf("chain %s not found in onchain state", targetChain.String())
				}
				if err := solanastateview.ValidateOwnershipSolana(&e, targetChain, c.MCMS != nil, targetChainState.RMNRemote, shared.RMNRemote, solana.PublicKey{}); err != nil {
					return fmt.Errorf("chain %s: %w", targetChain.String(), err)
				}
			}
		}
	}

	return nil
}

// CurseLaneOnlyOnSource curses a lane only on the source chain
// This will prevent message from source to destination to be initiated
// One noteworthy behaviour is that this means that message can be sent from destination to source but will not be executed on the source
// Given 3 chains A, B, C
// CurseLaneOnlyOnSource(A, B) will curse A with the curse subject of B
func CurseLaneOnlyOnSource(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Curse from source to destination
	return func(e cldf.Environment) ([]RMNCurseAction, error) {
		family, err := chain_selectors.GetSelectorFamily(sourceSelector)
		if err != nil {
			return nil, err
		}

		return []RMNCurseAction{
			{
				ChainSelector:  sourceSelector,
				SubjectToCurse: globals.FamilyAwareSelectorToSubject(destinationSelector, family),
			},
		}, nil
	}
}

// CurseGloballyOnlyOnChain curses a chain globally only on the source chain
// Given 3 chains A, B, C
// CurseGloballyOnlyOnChain(A) will curse a with the global curse subject only
func CurseGloballyOnlyOnChain(selector uint64) CurseAction {
	return func(e cldf.Environment) ([]RMNCurseAction, error) {
		return []RMNCurseAction{
			{
				ChainSelector:  selector,
				SubjectToCurse: globals.GlobalCurseSubject(),
			},
		}, nil
	}
}

// Call Curse on both RMNRemote from source and destination to prevent message from source to destination and vice versa
// Given 3 chains A, B, C
// CurseLaneBidirectionally(A, B) will curse A with the curse subject of B and B with the curse subject of A
func CurseLaneBidirectionally(sourceSelector uint64, destinationSelector uint64) CurseAction {

	// Bidirectional curse between two chains
	return func(e cldf.Environment) ([]RMNCurseAction, error) {
		curseActions1, err := CurseLaneOnlyOnSource(sourceSelector, destinationSelector)(e)
		if err != nil {
			return nil, err
		}

		curseActions2, err := CurseLaneOnlyOnSource(destinationSelector, sourceSelector)(e)
		if err != nil {
			return nil, err
		}

		return append(curseActions1, curseActions2...), nil
	}
}

// CurseChain do a global curse on chainSelector and curse chainSelector on all other chains
// Given 3 chains A, B, C
// CurseChain(A) will curse A with the global curse subject and curse B and C with the curse subject of A
func CurseChain(chainSelector uint64) CurseAction {
	return func(e cldf.Environment) ([]RMNCurseAction, error) {
		chainSelectors := GetAllCursableChainsSelector(e)

		// Curse all other chains to prevent onramp from sending message to the cursed chain
		var curseActions []RMNCurseAction
		for _, otherChainSelector := range chainSelectors {
			if otherChainSelector != chainSelector {
				family, err := chain_selectors.GetSelectorFamily(otherChainSelector)
				if err != nil {
					return nil, err
				}

				curseActions = append(curseActions, RMNCurseAction{
					ChainSelector:  otherChainSelector,
					SubjectToCurse: globals.FamilyAwareSelectorToSubject(chainSelector, family),
				})
			}
		}

		// Curse the chain with a global curse to prevent any onramp or offramp message from send message in and out of the chain
		globalCurse, err := CurseGloballyOnlyOnChain(chainSelector)(e)
		if err != nil {
			return nil, err
		}
		curseActions = append(curseActions, globalCurse...)

		return curseActions, nil
	}
}

func CurseGloballyAllChains() CurseAction {
	return func(e cldf.Environment) ([]RMNCurseAction, error) {
		chainSelectors := GetAllCursableChainsSelector(e)
		var curseActions []RMNCurseAction
		for _, chainSelector := range chainSelectors {
			actions, err := CurseGloballyOnlyOnChain(chainSelector)(e)
			if err != nil {
				return nil, err
			}
			curseActions = append(curseActions, actions...)
		}
		return curseActions, nil
	}
}

func FilterOutNotConnectedLanes(e cldf.Environment, curseActions []RMNCurseAction) ([]RMNCurseAction, error) {
	cursableChains, err := GetCursableChains(e)
	if err != nil {
		e.Logger.Errorf("failed to load cursable chains: %v", err)
		return nil, err
	}
	// Filter the curse action to only apply on the connected chains
	returnActions := make([]RMNCurseAction, 0)
	for _, action := range curseActions {
		if action.SubjectToCurse == globals.GlobalCurseSubject() {
			returnActions = append(returnActions, action)
			continue
		}

		targetChainSelector := action.ChainSelector

		targetFamily, err := chain_selectors.GetSelectorFamily(targetChainSelector)
		if err != nil {
			e.Logger.Errorf("failed to get family for chain %d: %v", targetChainSelector, err)
			return nil, err
		}

		sourceChainSelector := globals.FamilyAwareSubjectToSelector(action.SubjectToCurse, targetFamily)

		targetSourceConnected, err := cursableChains[targetChainSelector].IsConnectedToSourceChain(sourceChainSelector)
		if err != nil {
			e.Logger.Errorf("failed to check if offramp on chain %d is configured for source chain %d: %v", targetChainSelector, sourceChainSelector, err)
			return nil, err
		}

		if targetSourceConnected {
			returnActions = append(returnActions, action)
			continue
		}

		sourceTargetConnected, err := cursableChains[sourceChainSelector].IsConnectedToSourceChain(targetChainSelector)
		if err != nil {
			e.Logger.Errorf("failed to check if offramp on chain %d is configured for source chain %d: %v", sourceChainSelector, targetChainSelector, err)
			return nil, err
		}

		if sourceTargetConnected {
			returnActions = append(returnActions, action)
			continue
		}

		e.Logger.Warnf("Offramp on chain %d is not configured for source chain %d, skipping curse action", targetChainSelector, sourceChainSelector)
	}
	return returnActions, nil
}

func groupRMNSubjectBySelector(rmnSubjects []RMNCurseAction, avoidCursingSelf bool, onlyKeepGlobal bool) (map[uint64][]globals.Subject, error) {
	grouped := make(map[uint64][]globals.Subject)
	for _, s := range rmnSubjects {
		family, err := chain_selectors.GetSelectorFamily(s.ChainSelector)
		if err != nil {
			return nil, err
		}

		// Skip self-curse if needed
		if s.SubjectToCurse == globals.FamilyAwareSelectorToSubject(s.ChainSelector, family) && avoidCursingSelf {
			continue
		}
		// Initialize slice for this chain if needed
		if _, ok := grouped[s.ChainSelector]; !ok {
			grouped[s.ChainSelector] = []globals.Subject{}
		}
		// If global is already set and we only keep global, skip
		if onlyKeepGlobal && len(grouped[s.ChainSelector]) == 1 && grouped[s.ChainSelector][0] == globals.GlobalCurseSubject() {
			continue
		}
		// If subject is global and we only keep global, reset immediately
		if s.SubjectToCurse == globals.GlobalCurseSubject() && onlyKeepGlobal {
			grouped[s.ChainSelector] = []globals.Subject{globals.GlobalCurseSubject()}
			continue
		}
		// Ensure uniqueness
		duplicate := false
		for _, added := range grouped[s.ChainSelector] {
			if added == s.SubjectToCurse {
				duplicate = true
				break
			}
		}
		if !duplicate {
			grouped[s.ChainSelector] = append(grouped[s.ChainSelector], s.SubjectToCurse)
		}
	}

	return grouped, nil
}

// RMNCurseChangeset creates a new changeset for cursing chains or lanes on RMNRemote contracts.
// Example usage:
//
//	cfg := RMNCurseConfig{
//	    CurseActions: []CurseAction{
//	        CurseChain(SEPOLIA_CHAIN_SELECTOR),
//	        CurseLane(SEPOLIA_CHAIN_SELECTOR, AVAX_FUJI_CHAIN_SELECTOR),
//	    },
//	    CurseReason: "test curse",
//	    MCMS: &MCMSConfig{MinDelay: 0},
//	}
//	output, err := RMNCurseChangeset(env, cfg)
//
// This changeset is following an anti-pattern of supporting multiple chain families. Most changeset should be family specific.
// The decision to support multiple chain families here is due to the fact that curse changesets are emergency actions
// we want to keep a simple unified interface for all chain families to streamline emergency procedures.
func RMNCurseChangeset(e cldf.Environment, cfg RMNCurseConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("proposal to curse RMNs: " + cfg.Reason)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		actions, err := curseAction(e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate curse actions: %w", err)
		}

		curseActions = append(curseActions, actions...)
	}

	if !cfg.IncludeNotConnectedLanes {
		curseActions, err = FilterOutNotConnectedLanes(e, curseActions)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to filter out not connected lanes: %w", err)
		}
	}

	// Group curse actions by chain selector
	grouped, err := groupRMNSubjectBySelector(curseActions, true, true)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to group curse actions: %w", err)
	}
	// For each chain in the environment get the RMNRemote contract and call curse
	cursableChains, err := GetCursableChains(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get cursable chains: %w", err)
	}
	for selector, chain := range cursableChains {
		if curseSubjects, ok := grouped[selector]; ok {
			// Only curse the subjects that are not actually cursed
			notAlreadyCursedSubjects := make([]globals.Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.IsSubjectCursed(subject)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if chain %d is cursed: %w", selector, err)
				}

				if !cursed || cfg.Force {
					notAlreadyCursedSubjects = append(notAlreadyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is already cursed, ignoring it while cursing", cursableChains[selector].Name(), subject)
				}
			}

			if len(notAlreadyCursedSubjects) == 0 {
				e.Logger.Infof("chain %s is already cursed with all the subjects, skipping", cursableChains[selector].Name())
				continue
			}

			err := chain.Curse(deployerGroup, notAlreadyCursedSubjects)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to curse chain %d: %w", selector, err)
			}
			e.Logger.Infof("Cursed chain %d with subjects %v", selector, notAlreadyCursedSubjects)
		}
	}

	return deployerGroup.Enact()
}

// RMNUncurseChangeset creates a new changeset for uncursing chains or lanes on RMNRemote contracts.
// Curse actions are reused and reverted instead of applied in this changeset
// Example usage:
//
//	cfg := RMNCurseConfig{
//	    CurseActions: []CurseAction{
//	        CurseChain(SEPOLIA_CHAIN_SELECTOR),
//	        CurseLane(SEPOLIA_CHAIN_SELECTOR, AVAX_FUJI_CHAIN_SELECTOR),
//	    },
//	    MCMS: &MCMSConfig{MinDelay: 0},
//	}
//	output, err := RMNUncurseChangeset(env, cfg)
//
// This changeset is following an anti-pattern of supporting multiple chain families. Most changeset should be family specific.
// The decision to support multiple chain families here is due to the fact that curse changesets are emergency actions
// we want to keep a simple unified interface for all chain families to streamline emergency procedures.
func RMNUncurseChangeset(e cldf.Environment, cfg RMNCurseConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("proposal to uncurse RMNs: " + cfg.Reason)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		actions, err := curseAction(e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate curse actions: %w", err)
		}

		curseActions = append(curseActions, actions...)
	}
	// Group curse actions by chain selector
	grouped, err := groupRMNSubjectBySelector(curseActions, false, false)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to group curse actions: %w", err)
	}

	// For each chain in the environement get the RMNRemote contract and call uncurse
	cursableChains, err := GetCursableChains(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get cursable chains: %w", err)
	}
	for selector, chain := range cursableChains {
		if curseSubjects, ok := grouped[selector]; ok {
			// Only keep the subject that are actually cursed
			actuallyCursedSubjects := make([]globals.Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.IsSubjectCursed(subject)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if chain %d is cursed: %w", selector, err)
				}

				if cursed || cfg.Force {
					actuallyCursedSubjects = append(actuallyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is not cursed, ignoring it while uncursing", cursableChains[selector].Name(), subject)
				}
			}

			if len(actuallyCursedSubjects) == 0 {
				e.Logger.Infof("chain %s is not cursed with any of the subjects, skipping", cursableChains[selector].Name())
				continue
			}

			err := chain.Uncurse(deployerGroup, actuallyCursedSubjects)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to uncurse chain %d: %w", selector, err)
			}
			e.Logger.Infof("Uncursed chain %d with subjects %v", selector, actuallyCursedSubjects)
		}
	}

	return deployerGroup.Enact()
}

type CursableChain interface {
	Name() string
	IsConnectedToSourceChain(selector uint64) (bool, error)
	IsCursable() (bool, error)
	IsSubjectCursed(subject globals.Subject) (bool, error)
	Curse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error
	Uncurse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error
}

type SolanaCursableChain struct {
	selector uint64
	env      cldf.Environment
	chain    solanastateview.CCIPChainState
}

func (c SolanaCursableChain) IsSubjectCursed(subject globals.Subject) (bool, error) {
	chain := c.env.BlockChains.SolanaChains()[c.selector]
	curseSubject := solRmnRemote.CurseSubject{
		Value: subject,
	}
	rmnRemoteConfigPDA := c.chain.RMNRemoteConfigPDA
	solRmnRemote.SetProgramID(c.chain.RMNRemote)
	rmnRemoteCursesPDA := c.chain.RMNRemoteCursesPDA
	ix, err := solRmnRemote.NewVerifyNotCursedInstruction(
		curseSubject,
		rmnRemoteCursesPDA,
		rmnRemoteConfigPDA,
	).ValidateAndBuild()
	if err != nil {
		return false, fmt.Errorf("failed to generate instructions: %w", err)
	}
	_, err = solCommonUtil.SendAndConfirmWithLookupTables(context.Background(), chain.Client, []solana.Instruction{ix}, *chain.DeployerKey, rpc.CommitmentConfirmed, nil)
	if err != nil {
		c.env.Logger.Infof("Curse already exists for chain %d and curse subject %v", c.selector, curseSubject)
		return true, nil
	}

	return false, nil
}

func (c SolanaCursableChain) Curse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error {
	err := assertEndianness(subjects, chain_selectors.FamilySolana)
	if err != nil {
		return fmt.Errorf("failed to assert subject endianness: %w", err)
	}

	rmnRemoteConfigPDA := c.chain.RMNRemoteConfigPDA
	solRmnRemote.SetProgramID(c.chain.RMNRemote)
	rmnRemoteCursesPDA := c.chain.RMNRemoteCursesPDA
	deployer, err := deployerGroup.GetDeployerForSVM(c.selector)
	if err != nil {
		return fmt.Errorf("failed to get deployer for chain %d: %w", c.selector, err)
	}
	for _, subject := range subjects {
		curseSubject := solRmnRemote.CurseSubject{
			Value: subject,
		}
		_, err := deployer(func(authority solana.PublicKey) (solana.Instruction, string, cldf.ContractType, error) {
			ix, err := solRmnRemote.NewCurseInstruction(
				curseSubject,
				rmnRemoteConfigPDA,
				authority,
				rmnRemoteCursesPDA,
				solana.SystemProgramID,
			).ValidateAndBuild()

			if err != nil {
				return nil, "", "", fmt.Errorf("failed to generate instructions: %w", err)
			}

			return ix, c.chain.RMNRemote.String(), shared.RMNRemote, nil
		})
		if err != nil {
			return fmt.Errorf("failed to build curse instruction for subject %x on chain %d: %w", subject, c.selector, err)
		}
	}
	return nil
}

func (c SolanaCursableChain) Uncurse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error {
	err := assertEndianness(subjects, chain_selectors.FamilySolana)
	if err != nil {
		return fmt.Errorf("failed to assert subject endianness: %w", err)
	}

	rmnRemoteConfigPDA := c.chain.RMNRemoteConfigPDA
	solRmnRemote.SetProgramID(c.chain.RMNRemote)
	rmnRemoteCursesPDA := c.chain.RMNRemoteCursesPDA
	deployer, err := deployerGroup.GetDeployerForSVM(c.selector)
	if err != nil {
		return fmt.Errorf("failed to get deployer for chain %d: %w", c.selector, err)
	}
	for _, subject := range subjects {
		curseSubject := solRmnRemote.CurseSubject{
			Value: subject,
		}
		_, err := deployer(func(authority solana.PublicKey) (solana.Instruction, string, cldf.ContractType, error) {
			ix, err := solRmnRemote.NewUncurseInstruction(
				curseSubject,
				rmnRemoteConfigPDA,
				authority,
				rmnRemoteCursesPDA,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return nil, "", "", fmt.Errorf("failed to generate instructions: %w", err)
			}
			return ix, c.chain.RMNRemote.String(), shared.RMNRemote, nil
		})
		if err != nil {
			return fmt.Errorf("failed to build uncurse instruction for subject %x on chain %d: %w", subject, c.selector, err)
		}
	}
	return nil
}

func (c SolanaCursableChain) IsCursable() (bool, error) {
	return c.chain.RMNRemote != solana.PublicKey{}, nil
}

func (c SolanaCursableChain) IsConnectedToSourceChain(selector uint64) (bool, error) {
	state, err := stateview.LoadOnchainStateSolana(c.env)
	if err != nil {
		return false, fmt.Errorf("failed to load onchain state: %w", err)
	}

	pda, _, err := solState.FindOfframpSourceChainPDA(selector, state.SolChains[c.selector].OffRamp)
	if err != nil {
		return false, fmt.Errorf("failed to find offramp source chain pda: %w", err)
	}

	var chainStateAccount solOffRamp.SourceChain
	if err = c.env.BlockChains.SolanaChains()[c.selector].GetAccountDataBorshInto(context.Background(), pda, &chainStateAccount); err != nil {
		return false, nil
	}

	return chainStateAccount.Config.IsEnabled, nil
}

func (c SolanaCursableChain) Name() string {
	return c.env.BlockChains.SolanaChains()[c.selector].Name()
}

type EvmCursableChain struct {
	selector uint64
	env      cldf.Environment
	chain    evm.CCIPChainState
}

func (c EvmCursableChain) Name() string {
	return c.env.BlockChains.EVMChains()[c.selector].Name()
}

func (c EvmCursableChain) IsConnectedToSourceChain(sourceSelector uint64) (bool, error) {
	destChain := c.chain
	config, err := destChain.OffRamp.GetSourceChainConfig(nil, sourceSelector)
	if err != nil {
		return false, fmt.Errorf("failed to check if chain %d is connected to chain %d: %w", c.selector, sourceSelector, err)
	}
	if !config.IsEnabled {
		return false, nil
	}
	return true, nil
}

func (c EvmCursableChain) IsSubjectCursed(subject globals.Subject) (bool, error) {
	cursed, err := c.chain.RMNRemote.IsCursed(nil, subject)
	if err != nil {
		return false, fmt.Errorf("failed to check if chain %d is cursed: %w", c.selector, err)
	}
	return cursed, nil
}

func (c EvmCursableChain) IsCursable() (bool, error) {
	return c.chain.RMNRemote != nil, nil
}

func (c EvmCursableChain) Curse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error {
	err := assertEndianness(subjects, chain_selectors.FamilyEVM)
	if err != nil {
		return fmt.Errorf("failed to assert subject endianness: %w", err)
	}

	deployer, err := deployerGroup.GetDeployer(c.selector)
	if err != nil {
		return fmt.Errorf("failed to get deployer for chain %d: %w", c.selector, err)
	}

	_, err = c.chain.RMNRemote.Curse0(deployer, subjects)
	if err != nil {
		return fmt.Errorf("failed to curse chain %d: %w", c.selector, err)
	}
	return nil
}

func (c EvmCursableChain) Uncurse(deployerGroup *deployergroup.DeployerGroup, subjects []globals.Subject) error {
	err := assertEndianness(subjects, chain_selectors.FamilyEVM)
	if err != nil {
		return fmt.Errorf("failed to assert subject endianness: %w", err)
	}

	deployer, err := deployerGroup.GetDeployer(c.selector)
	if err != nil {
		return fmt.Errorf("failed to get deployer for chain %d: %w", c.selector, err)
	}

	_, err = c.chain.RMNRemote.Uncurse0(deployer, subjects)
	if err != nil {
		return fmt.Errorf("failed to uncurse chain %d: %w", c.selector, err)
	}
	return nil
}

func GetCursableChains(env cldf.Environment) (map[uint64]CursableChain, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load onchain state: %w", err)
	}
	cursableChains := make(map[uint64]CursableChain)
	for selector := range state.Chains {
		cursableChains[selector] = EvmCursableChain{
			selector: selector,
			chain:    state.Chains[selector], // Access chain state directly
			env:      env,
		}
	}

	for selector, chain := range state.SolChains {
		cursableChains[selector] = SolanaCursableChain{
			selector: selector,
			chain:    chain,
			env:      env,
		}
	}

	activeCursableChains := make(map[uint64]CursableChain)
	for selector, chain := range cursableChains {
		cursable, err := chain.IsCursable()
		if err != nil {
			return nil, fmt.Errorf("failed to check if chain %d is cursable: %w", selector, err)
		}
		if cursable {
			activeCursableChains[selector] = chain
		}
	}

	return activeCursableChains, nil
}

func GetAllCursableChainsSelector(env cldf.Environment) []uint64 {
	selectors := make([]uint64, 0)
	selectors = append(selectors, env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))...)
	solSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	selectors = append(selectors, solSelectors...)
	return selectors
}

func assertEndianness(subjects []globals.Subject, family string) error {
	for _, subject := range subjects {
		if subject == globals.GlobalCurseSubject() {
			continue
		}
		switch family {
		case chain_selectors.FamilySolana:
			// Solana uses little endian to encode the subject so we expect the last 8 bytes to be 0
			if !bytes.Equal(subject[8:], []byte{0, 0, 0, 0, 0, 0, 0, 0}) {
				return fmt.Errorf("endianness incorrect for Solana curse subject: %s", subject)
			}
		default:
			// EVM uses big endian to encode the subject so we expect the first 8 bytes to be 0
			if !bytes.Equal(subject[:8], []byte{0, 0, 0, 0, 0, 0, 0, 0}) {
				return fmt.Errorf("endianness incorrect for curse subject: %s", subject)
			}
		}
	}
	return nil
}
