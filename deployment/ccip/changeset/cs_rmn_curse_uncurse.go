package changeset

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

var (
	_ deployment.ChangeSet[RMNCurseConfig] = RMNCurseChangeset
	_ deployment.ChangeSet[RMNCurseConfig] = RMNUncurseChangeset
)

// GlobalCurseSubject as defined here: https://github.com/smartcontractkit/chainlink/blob/new-rmn-curse-changeset/contracts/src/v0.8/ccip/rmn/RMNRemote.sol#L15
func GlobalCurseSubject() Subject {
	return Subject{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
}

// RMNCurseAction represent a curse action to be applied on a chain (ChainSelector) with a specific subject (SubjectToCurse)
// The curse action will by applied by calling the Curse method on the RMNRemote contract on the chain (ChainSelector)
type RMNCurseAction struct {
	ChainSelector  uint64
	SubjectToCurse Subject
}

// CurseAction is a function that returns a list of RMNCurseAction to be applied on a chain
// CurseChain, CurseLane, CurseGloballyOnlyOnSource are examples of function implementing CurseAction
type CurseAction func(e deployment.Environment) []RMNCurseAction

type RMNCurseConfig struct {
	MCMS         *MCMSConfig
	CurseActions []CurseAction
	Reason       string
}

func (c RMNCurseConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)

	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	if len(c.CurseActions) == 0 {
		return errors.New("curse actions are required")
	}

	if c.Reason == "" {
		return errors.New("reason is required")
	}

	validSubjects := map[Subject]struct{}{
		GlobalCurseSubject(): {},
	}
	for _, selector := range e.AllChainSelectors() {
		validSubjects[SelectorToSubject(selector)] = struct{}{}
	}

	for _, curseAction := range c.CurseActions {
		result := curseAction(e)
		for _, action := range result {
			targetChain := e.Chains[action.ChainSelector]
			targetChainState, ok := state.Chains[action.ChainSelector]
			if !ok {
				return fmt.Errorf("chain %s not found in onchain state", targetChain.String())
			}

			if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, targetChain.DeployerKey.From, targetChainState.Timelock.Address(), targetChainState.RMNRemote); err != nil {
				return fmt.Errorf("chain %s: %w", targetChain.String(), err)
			}

			if err = deployment.IsValidChainSelector(action.ChainSelector); err != nil {
				return fmt.Errorf("invalid chain selector %d for chain %s", action.ChainSelector, targetChain.String())
			}

			if _, ok := validSubjects[action.SubjectToCurse]; !ok {
				return fmt.Errorf("invalid subject %x for chain %s", action.SubjectToCurse, targetChain.String())
			}
		}
	}

	return nil
}

type Subject = [16]byte

func SelectorToSubject(selector uint64) Subject {
	var b Subject
	binary.BigEndian.PutUint64(b[8:], selector)
	return b
}

// CurseLaneOnlyOnSource curses a lane only on the source chain
// This will prevent message from source to destination to be initiated
// One noteworthy behaviour is that this means that message can be sent from destination to source but will not be executed on the source
// Given 3 chains A, B, C
// CurseLaneOnlyOnSource(A, B) will curse A with the curse subject of B
func CurseLaneOnlyOnSource(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Curse from source to destination
	return func(e deployment.Environment) []RMNCurseAction {
		return []RMNCurseAction{
			{
				ChainSelector:  sourceSelector,
				SubjectToCurse: SelectorToSubject(destinationSelector),
			},
		}
	}
}

// CurseGloballyOnlyOnChain curses a chain globally only on the source chain
// Given 3 chains A, B, C
// CurseGloballyOnlyOnChain(A) will curse a with the global curse subject only
func CurseGloballyOnlyOnChain(selector uint64) CurseAction {
	return func(e deployment.Environment) []RMNCurseAction {
		return []RMNCurseAction{
			{
				ChainSelector:  selector,
				SubjectToCurse: GlobalCurseSubject(),
			},
		}
	}
}

// Call Curse on both RMNRemote from source and destination to prevent message from source to destination and vice versa
// Given 3 chains A, B, C
// CurseLaneBidirectionally(A, B) will curse A with the curse subject of B and B with the curse subject of A
func CurseLaneBidirectionally(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Bidirectional curse between two chains
	return func(e deployment.Environment) []RMNCurseAction {
		return append(
			CurseLaneOnlyOnSource(sourceSelector, destinationSelector)(e),
			CurseLaneOnlyOnSource(destinationSelector, sourceSelector)(e)...,
		)
	}
}

// CurseChain do a global curse on chainSelector and curse chainSelector on all other chains
// Given 3 chains A, B, C
// CurseChain(A) will curse A with the global curse subject and curse B and C with the curse subject of A
func CurseChain(chainSelector uint64) CurseAction {
	return func(e deployment.Environment) []RMNCurseAction {
		chainSelectors := e.AllChainSelectors()

		// Curse all other chains to prevent onramp from sending message to the cursed chain
		var curseActions []RMNCurseAction
		for _, otherChainSelector := range chainSelectors {
			if otherChainSelector != chainSelector {
				curseActions = append(curseActions, RMNCurseAction{
					ChainSelector:  otherChainSelector,
					SubjectToCurse: SelectorToSubject(chainSelector),
				})
			}
		}

		// Curse the chain with a global curse to prevent any onramp or offramp message from send message in and out of the chain
		curseActions = append(curseActions, CurseGloballyOnlyOnChain(chainSelector)(e)...)

		return curseActions
	}
}

func groupRMNSubjectBySelector(rmnSubjects []RMNCurseAction, avoidCursingSelf bool, onlyKeepGlobal bool) map[uint64][]Subject {
	grouped := make(map[uint64][]Subject)
	for _, s := range rmnSubjects {
		// Skip self-curse if needed
		if s.SubjectToCurse == SelectorToSubject(s.ChainSelector) && avoidCursingSelf {
			continue
		}
		// Initialize slice for this chain if needed
		if _, ok := grouped[s.ChainSelector]; !ok {
			grouped[s.ChainSelector] = []Subject{}
		}
		// If global is already set and we only keep global, skip
		if onlyKeepGlobal && len(grouped[s.ChainSelector]) == 1 && grouped[s.ChainSelector][0] == GlobalCurseSubject() {
			continue
		}
		// If subject is global and we only keep global, reset immediately
		if s.SubjectToCurse == GlobalCurseSubject() && onlyKeepGlobal {
			grouped[s.ChainSelector] = []Subject{GlobalCurseSubject()}
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

	return grouped
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
func RMNCurseChangeset(e deployment.Environment, cfg RMNCurseConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("proposal to curse RMNs: " + cfg.Reason)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		curseActions = append(curseActions, curseAction(e)...)
	}
	// Group curse actions by chain selector
	grouped := groupRMNSubjectBySelector(curseActions, true, true)
	// For each chain in the environment get the RMNRemote contract and call curse
	for selector, chain := range state.Chains {
		deployer, err := deployerGroup.GetDeployer(selector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", selector, err)
		}
		if curseSubjects, ok := grouped[selector]; ok {
			// Only curse the subjects that are not actually cursed
			notAlreadyCursedSubjects := make([]Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.RMNRemote.IsCursed(nil, subject)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if chain %d is cursed: %w", selector, err)
				}

				if !cursed {
					notAlreadyCursedSubjects = append(notAlreadyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is already cursed, ignoring it while cursing", e.Chains[selector].Name(), subject)
				}
			}

			if len(notAlreadyCursedSubjects) == 0 {
				e.Logger.Infof("chain %s is already cursed with all the subjects, skipping", e.Chains[selector].Name())
				continue
			}

			_, err := chain.RMNRemote.Curse0(deployer, notAlreadyCursedSubjects)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to curse chain %d: %w", selector, err)
			}
			e.Logger.Infof("Cursed chain %d with subjects %v", selector, notAlreadyCursedSubjects)
		}
	}

	return deployerGroup.Enact()
}

// RMNUncurseChangeset creates a new changeset for uncursing chains or lanes on RMNRemote contracts.
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
// Curse actions are reused and reverted instead of applied in this changeset
func RMNUncurseChangeset(e deployment.Environment, cfg RMNCurseConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("proposal to uncurse RMNs: " + cfg.Reason)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		curseActions = append(curseActions, curseAction(e)...)
	}
	// Group curse actions by chain selector
	grouped := groupRMNSubjectBySelector(curseActions, false, false)

	// For each chain in the environement get the RMNRemote contract and call uncurse
	for selector, chain := range state.Chains {
		deployer, err := deployerGroup.GetDeployer(selector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", selector, err)
		}

		if curseSubjects, ok := grouped[selector]; ok {
			// Only keep the subject that are actually cursed
			actuallyCursedSubjects := make([]Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.RMNRemote.IsCursed(nil, subject)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if chain %d is cursed: %w", selector, err)
				}

				if cursed {
					actuallyCursedSubjects = append(actuallyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is not cursed, ignoring it while uncursing", e.Chains[selector].Name(), subject)
				}
			}

			if len(actuallyCursedSubjects) == 0 {
				e.Logger.Infof("chain %s is not cursed with any of the subjects, skipping", e.Chains[selector].Name())
				continue
			}

			_, err := chain.RMNRemote.Uncurse0(deployer, actuallyCursedSubjects)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to uncurse chain %d: %w", selector, err)
			}
			e.Logger.Infof("Uncursed chain %d with subjects %v", selector, actuallyCursedSubjects)
		}
	}

	return deployerGroup.Enact()
}
