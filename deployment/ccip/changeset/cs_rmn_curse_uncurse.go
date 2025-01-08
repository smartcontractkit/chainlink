package changeset

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
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
	CurseReason  string
}

type Subject = [16]byte

func SelectorToSubject(subject uint64) Subject {
	var b Subject
	binary.LittleEndian.PutUint64(b[:8], subject)
	return b
}

func CurseLane(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Bidirectional curse between two chains
	return func(e deployment.Environment) []RMNCurseAction {
		return []RMNCurseAction{
			{
				ChainSelector:  sourceSelector,
				SubjectToCurse: SelectorToSubject(destinationSelector),
			},
			{
				ChainSelector:  destinationSelector,
				SubjectToCurse: SelectorToSubject(sourceSelector),
			},
		}
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

func groupRMNSubjectBySelector(rmnSubjects []RMNCurseAction) map[uint64][]Subject {
	grouped := make(map[uint64][]Subject)
	for _, subject := range rmnSubjects {
		grouped[subject.ChainSelector] = append(grouped[subject.ChainSelector], subject.SubjectToCurse)
	}

	// Only keep unique subjects, preserve only global curse if present and eliminate any curse where the selector is the same as the subject
	for chainSelector, subjects := range grouped {
		uniqueSubjects := make(map[Subject]struct{})
		for _, subject := range subjects {
			if subject == SelectorToSubject(chainSelector) {
				continue
			}
			uniqueSubjects[subject] = struct{}{}
		}

		if _, ok := uniqueSubjects[GlobalCurseSubject()]; ok {
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
//	    CurseActions: []func(deployment.Environment) []RMNCurseAction{
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
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := NewDeployerGroup(e, state, cfg.MCMS)

	// Generate curse actions
	var curseActions []RMNCurseAction
	for _, curseAction := range cfg.CurseActions {
		curseActions = append(curseActions, curseAction(e)...)
	}
	// Group curse actions by chain selector
	grouped := groupRMNSubjectBySelector(curseActions)

	// For each chain in the environement get the RMNRemote contract and call curse
	for selector, chain := range state.Chains {
		deployer := deployerGroup.getDeployer(selector)
		if curseSubjects, ok := grouped[selector]; ok {
			_, err := chain.RMNRemote.Curse0(deployer, curseSubjects)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to curse chain %d: %w", selector, err)
			}
		}
	}

	return deployerGroup.enact(fmt.Sprintf("proposal to curse RMNs: %s", cfg.CurseReason))
}
