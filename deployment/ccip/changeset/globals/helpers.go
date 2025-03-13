package globals

import "encoding/binary"

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
