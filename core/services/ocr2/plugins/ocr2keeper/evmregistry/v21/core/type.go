package core

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
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
func GetUpkeepType(id ocr2keepers.UpkeepIdentifier) types.UpkeepType {
	for i := upkeepTypeStartIndex; i < upkeepTypeByteIndex; i++ {
		if id[i] != 0 { // old id
			return types.ConditionTrigger
		}
	}
	typeByte := id[upkeepTypeByteIndex]
	return types.UpkeepType(typeByte)
}

func getUpkeepTypeFromBigInt(id *big.Int) (types.UpkeepType, bool) {
	uid := &ocr2keepers.UpkeepIdentifier{}
	ok := uid.FromBigInt(id)
	return GetUpkeepType(*uid), ok
}
