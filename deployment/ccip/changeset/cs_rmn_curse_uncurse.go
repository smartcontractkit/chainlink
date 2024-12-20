package changeset

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	GLOBAL_CURSE_SUBJECT = 0
)

type RMNCurseAction struct {
	ChainSelector  uint64
	SubjectToCurse uint64
}

type CurseAction func(e deployment.Environment) []RMNCurseAction

type RMNCurseConfig struct {
	HomeChainSelector uint64
	MCMS              *MCMSConfig
	CurseActions      []CurseAction
	CurseReason       string
}

func subjectToByte16(subject uint64) [16]byte {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[:8], subject)
	return b
}

func CurseLane(sourceSelector uint64, destinationSelector uint64) CurseAction {
	// Bidirectional curse between two chains
	return func(e deployment.Environment) []RMNCurseAction {
		return []RMNCurseAction{
			{
				ChainSelector:  sourceSelector,
				SubjectToCurse: destinationSelector,
			},
			{
				ChainSelector:  destinationSelector,
				SubjectToCurse: sourceSelector,
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
					SubjectToCurse: chainSelector,
				})
			}
		}

		// Curse the chain with a global curse to prevent any onramp or offramp message from send message in and out of the chain
		curseActions = append(curseActions, RMNCurseAction{
			ChainSelector:  chainSelector,
			SubjectToCurse: GLOBAL_CURSE_SUBJECT,
		})

		return curseActions
	}
}

func groupRMNSubjectBySelector(rmnSubjects []RMNCurseAction) map[uint64][]uint64 {
	grouped := make(map[uint64][]uint64)
	for _, subject := range rmnSubjects {
		grouped[subject.ChainSelector] = append(grouped[subject.ChainSelector], subject.SubjectToCurse)
	}

	// Only keep unique subjects, preserve only global curse if present and eliminate any curse where the selector is the same as the subject
	for chainSelector, subjects := range grouped {
		uniqueSubjects := make(map[uint64]struct{})
		for _, subject := range subjects {
			if subject == chainSelector {
				continue
			}
			uniqueSubjects[subject] = struct{}{}
		}

		if _, ok := uniqueSubjects[GLOBAL_CURSE_SUBJECT]; ok {
			grouped[chainSelector] = []uint64{GLOBAL_CURSE_SUBJECT}
		} else {
			var uniqueSubjectsSlice []uint64
			for subject := range uniqueSubjects {
				uniqueSubjectsSlice = append(uniqueSubjectsSlice, subject)
			}
			grouped[chainSelector] = uniqueSubjectsSlice
		}
	}

	return grouped
}

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
			subjectsByte16 := make([][16]byte, len(curseSubjects))
			for i, subject := range curseSubjects {
				subjectsByte16[i] = subjectToByte16(subject)
			}

			_, err := chain.RMNRemote.Curse0(deployer, subjectsByte16)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to curse chain %d: %w", selector, err)
			}
		}
	}

	return deployerGroup.enact(fmt.Sprintf("proposal to curse RMNs: %s", cfg.CurseReason))
}
