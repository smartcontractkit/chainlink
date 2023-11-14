package core

import (
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

const (
	// upkeepTypeStartIndex is the index where the upkeep type bytes start.
	// for 2.1 we use 11 zeros (reserved bytes for future use)
	// and 1 byte to represent the type, with index equal upkeepTypeByteIndex
	upkeepTypeStartIndex = 4
	// upkeepTypeByteIndex is the index of the byte that holds the upkeep type.
	upkeepTypeByteIndex = 15
)

// GetUpkeepType returns the upkeep type from the given ID.
// it follows the same logic as the contract, but performs it locally.
func GetUpkeepType(id ocr2keepers.UpkeepIdentifier) ocr2keepers.UpkeepType {
	for i := upkeepTypeStartIndex; i < upkeepTypeByteIndex; i++ {
		if id[i] != 0 { // old id
			return ocr2keepers.ConditionTrigger
		}
	}
	typeByte := id[upkeepTypeByteIndex]
	return ocr2keepers.UpkeepType(typeByte)
}

func getUpkeepTypeFromBigInt(id *big.Int) (ocr2keepers.UpkeepType, bool) {
	uid := &ocr2keepers.UpkeepIdentifier{}
	ok := uid.FromBigInt(id)
	return GetUpkeepType(*uid), ok
}
