package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

// genUpkeepID generates an ocr2keepers.UpkeepIdentifier with a specific UpkeepType and some random string
func genUpkeepID(uType ocr2keepers.UpkeepType, rand string) ocr2keepers.UpkeepIdentifier {
	b := append([]byte{1}, common.LeftPadBytes([]byte{uint8(uType)}, 15)...)
	b = append(b, []byte(rand)...)
	b = common.RightPadBytes(b, 32-len(b))
	if len(b) > 32 {
		b = b[:32]
	}
	var id [32]byte
	copy(id[:], b)
	return ocr2keepers.UpkeepIdentifier(id)
}

// upkeepIDFromInt converts an int string to ocr2keepers.UpkeepIdentifier
func upkeepIDFromInt(id string) ocr2keepers.UpkeepIdentifier {
	uid := &ocr2keepers.UpkeepIdentifier{}
	idInt, _ := big.NewInt(0).SetString(id, 10)
	uid.FromBigInt(idInt)
	return *uid
}
