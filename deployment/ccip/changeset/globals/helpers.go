package globals

import (
	"encoding/binary"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

// GlobalCurseSubject as defined here: https://github.com/smartcontractkit/chainlink/blob/new-rmn-curse-changeset/contracts/src/v0.8/ccip/rmn/RMNRemote.sol#L15
func GlobalCurseSubject() Subject {
	return Subject{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
}

type Subject = [16]byte

func selectorToSubject(selector uint64) Subject {
	var b Subject
	binary.BigEndian.PutUint64(b[8:], selector)
	return b
}

func subjectToSelector(subject [16]byte) uint64 {
	if subject == GlobalCurseSubject() {
		return 0
	}

	return binary.BigEndian.Uint64(subject[8:])
}

func selectorToSolanaSubject(selector uint64) Subject {
	var b Subject
	binary.LittleEndian.PutUint64(b[0:], selector)
	return b
}

func subjectToSolanaSelector(subject [16]byte) uint64 {
	if subject == GlobalCurseSubject() {
		return 0
	}

	return binary.LittleEndian.Uint64(subject[:])
}

func FamilyAwareSubjectToSelector(subject Subject, family string) uint64 {
	switch family {
	case chainsel.FamilySolana:
		return subjectToSolanaSelector(subject)
	default:
		return subjectToSelector(subject)
	}
}

func FamilyAwareSelectorToSubject(selector uint64, family string) Subject {
	switch family {
	case chainsel.FamilySolana:
		return selectorToSolanaSubject(selector)
	default:
		return selectorToSubject(selector)
	}
}
