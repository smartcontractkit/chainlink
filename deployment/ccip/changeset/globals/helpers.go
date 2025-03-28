package globals

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
)

// GlobalCurseSubject as defined here: https://github.com/smartcontractkit/chainlink/blob/new-rmn-curse-changeset/contracts/src/v0.8/ccip/rmn/RMNRemote.sol#L15
func GlobalCurseSubject() Subject {
	return Subject{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
}

type Subject = [16]byte

func SelectorToSubject(selector uint64) Subject {
	var b Subject
	binary.BigEndian.PutUint64(b[8:], selector)
	return b
}

func SubjectToSelector(subject [16]byte) uint64 {
	if subject == GlobalCurseSubject() {
		return 0
	}

	return binary.BigEndian.Uint64(subject[8:])
}

func MergeChangesetOutput(dest *deployment.ChangesetOutput, src deployment.ChangesetOutput) error {
	if dest == nil {
		return nil
	}

	if dest.AddressBook == nil {
		dest.AddressBook = src.AddressBook
	} else if src.AddressBook != nil {
		if err := dest.AddressBook.Merge(src.AddressBook); err != nil {
			return fmt.Errorf("failed to merge address book: %w", err)
		}
	}
	if dest.Jobs == nil {
		dest.Jobs = src.Jobs
	} else if src.Jobs != nil {
		dest.Jobs = append(dest.Jobs, src.Jobs...)
	}
	if dest.MCMSTimelockProposals == nil {
		dest.MCMSTimelockProposals = src.MCMSTimelockProposals
	} else if src.MCMSTimelockProposals != nil {
		dest.MCMSTimelockProposals = append(dest.MCMSTimelockProposals, src.MCMSTimelockProposals...)
	}
	if dest.DescribedTimelockProposals == nil {
		dest.DescribedTimelockProposals = src.DescribedTimelockProposals
	} else if src.DescribedTimelockProposals != nil {
		dest.DescribedTimelockProposals = append(dest.DescribedTimelockProposals, src.DescribedTimelockProposals...)
	}

	if dest.MCMSProposals == nil {
		dest.MCMSProposals = src.MCMSProposals
	} else if src.MCMSProposals != nil {
		dest.MCMSProposals = append(dest.MCMSProposals, src.MCMSProposals...)
	}

	return nil
}
