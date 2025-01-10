package changeset

import (
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func GlobalCurseSubject() Subject {
	return Subject{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
}

type RMNCurseAction struct {
	ChainSelector  uint64
	SubjectToCurse Subject
}

type CurseAction func(e deployment.Environment) []RMNCurseAction

type RMNCurseConfig struct {
	MCMS         *MCMSConfig
	CurseActions []CurseAction
	Reason       string
}

func (c RMNCurseConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)

	if err != nil {
		return errors.Errorf("failed to load onchain state: %v", err)
	}

	if len(c.CurseActions) == 0 {
		return errors.Errorf("curse actions are required")
	}

	if c.Reason == "" {
		return errors.Errorf("reason is required")
	}

	validSelectors := e.AllChainSelectors()
	validSubjects := make([]Subject, 0, len(validSelectors)+1)
	for _, selector := range validSelectors {
		validSubjects = append(validSubjects, SelectorToSubject(selector))
	}
	validSubjects = append(validSubjects, GlobalCurseSubject())

	for _, curseAction := range c.CurseActions {
		result := curseAction(e)
		for _, action := range result {
			targetChain := e.Chains[action.ChainSelector]
			targetChainState := state.Chains[action.ChainSelector]

			if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, targetChain.DeployerKey.From, targetChainState.Timelock.Address(), targetChainState.RMNRemote); err != nil {
				return fmt.Errorf("chain %s: %w", targetChain.String(), err)
			}

			if !slices.Contains(validSelectors, action.ChainSelector) {
				return errors.Errorf("invalid chain selector %d for chain %s", action.ChainSelector, targetChain.String())
			}

			if !slices.Contains(validSubjects, action.SubjectToCurse) {
				return errors.Errorf("invalid subject %x for chain %s", action.SubjectToCurse, targetChain.String())
			}
		}
	}

	return nil
}

type Subject = [16]byte

func SelectorToSubject(subject uint64) Subject {
	var b Subject
	binary.BigEndian.PutUint64(b[8:], subject)
	return b
}

// CurseLaneOnlyOnSource curses a lane only on the source chain
// This will prevent message from source to destination to be initiated
// One noteworthy behaviour is that this means that message can be sent from destination to source but will not be executed on the source
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

func CurseLane(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Bidirectional curse between two chains
	return func(e deployment.Environment) []RMNCurseAction {
		return append(
			CurseLaneOnlyOnSource(sourceSelector, destinationSelector)(e),
			CurseLaneOnlyOnSource(destinationSelector, sourceSelector)(e)...,
		)
	}
}

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
		curseActions = append(curseActions, RMNCurseAction{
			ChainSelector:  chainSelector,
			SubjectToCurse: GlobalCurseSubject(),
		})

		return curseActions
	}
}

func groupRMNSubjectBySelector(rmnSubjects []RMNCurseAction, filter bool) map[uint64][]Subject {
	grouped := make(map[uint64][]Subject)
	for _, subject := range rmnSubjects {
		grouped[subject.ChainSelector] = append(grouped[subject.ChainSelector], subject.SubjectToCurse)
	}

	// Only keep unique subjects, preserve only global curse if present and eliminate any curse where the selector is the same as the subject
	// If filter is false then only make sure that there is no duplicate subject
	for chainSelector, subjects := range grouped {
		uniqueSubjects := make(map[Subject]struct{})
		for _, subject := range subjects {
			if subject == SelectorToSubject(chainSelector) && filter {
				continue
			}
			uniqueSubjects[subject] = struct{}{}
		}

		if _, ok := uniqueSubjects[GlobalCurseSubject()]; ok && filter {
			grouped[chainSelector] = []Subject{GlobalCurseSubject()}
		} else {
			var uniqueSubjectsSlice []Subject
			for subject := range uniqueSubjects {
				uniqueSubjectsSlice = append(uniqueSubjectsSlice, subject)
			}
			grouped[chainSelector] = uniqueSubjectsSlice
		}
	}

	return grouped
}

// NewRMNCurseChangeset creates a new changeset for cursing chains or lanes on RMNRemote contracts.
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
//	output, err := NewRMNCurseChangeset(env, cfg)
func NewRMNCurseChangeset(e deployment.Environment, cfg RMNCurseConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Errorf("failed to load onchain state: %v", err)
	}
	deployerGroup := NewDeployerGroup(e, state, cfg.MCMS)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		curseActions = append(curseActions, curseAction(e)...)
	}
	// Group curse actions by chain selector
	grouped := groupRMNSubjectBySelector(curseActions, true)

	// For each chain in the environement get the RMNRemote contract and call curse
	for selector, chain := range state.Chains {
		deployer := deployerGroup.getDeployer(selector)
		if curseSubjects, ok := grouped[selector]; ok {
			// Only curse the subject that are not actually cursed
			notAlreadyCursedSubjects := make([]Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.RMNRemote.IsCursed(nil, subject)
				if err != nil {
					return deployment.ChangesetOutput{}, errors.Errorf("failed to check if chain %d is cursed: %v", selector, err)
				}

				if !cursed {
					notAlreadyCursedSubjects = append(notAlreadyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is already cursed, ignoring it while cursing", e.Chains[selector].Name(), subject)
				}
			}
			_, err := chain.RMNRemote.Curse0(deployer, notAlreadyCursedSubjects)
			if err != nil {
				return deployment.ChangesetOutput{}, errors.Errorf("failed to curse chain %d: %v", selector, err)
			}
		}
	}

	return deployerGroup.enact("proposal to curse RMNs: " + cfg.Reason)
}

// NewRMNUncurseChangeset creates a new changeset for uncursing chains or lanes on RMNRemote contracts.
// Example usage:
//
//	cfg := RMNCurseConfig{
//	    CurseActions: []CurseAction{
//	        CurseChain(SEPOLIA_CHAIN_SELECTOR),
//	        CurseLane(SEPOLIA_CHAIN_SELECTOR, AVAX_FUJI_CHAIN_SELECTOR),
//	    },
//	    MCMS: &MCMSConfig{MinDelay: 0},
//	}
//	output, err := NewRMNUncurseChangeset(env, cfg)
//
// Curse actions are reused and reverted instead of applied in this changeset
func NewRMNUncurseChangeset(e deployment.Environment, cfg RMNCurseConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Errorf("failed to load onchain state: %v", err)
	}
	deployerGroup := NewDeployerGroup(e, state, cfg.MCMS)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		curseActions = append(curseActions, curseAction(e)...)
	}
	// Group curse actions by chain selector
	grouped := groupRMNSubjectBySelector(curseActions, false)

	// For each chain in the environement get the RMNRemote contract and call uncurse
	for selector, chain := range state.Chains {
		deployer := deployerGroup.getDeployer(selector)
		if curseSubjects, ok := grouped[selector]; ok {
			// Only keep the subject that are actually cursed
			actuallyCursedSubjects := make([]Subject, 0)
			for _, subject := range curseSubjects {
				cursed, err := chain.RMNRemote.IsCursed(nil, subject)
				if err != nil {
					return deployment.ChangesetOutput{}, errors.Errorf("failed to check if chain %d is cursed: %v", selector, err)
				}

				if cursed {
					actuallyCursedSubjects = append(actuallyCursedSubjects, subject)
				} else {
					e.Logger.Warnf("chain %s subject %x is not cursed, ignoring it while uncursing", e.Chains[selector].Name(), subject)
				}
			}

			_, err := chain.RMNRemote.Uncurse0(deployer, actuallyCursedSubjects)
			if err != nil {
				return deployment.ChangesetOutput{}, errors.Errorf("failed to uncurse chain %d: %v", selector, err)
			}
		}
	}

	return deployerGroup.enact("proposal to uncurse RMNs: %s" + cfg.Reason)
}
