package core

import (
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

type UpkeepType uint8

const (
	ConditionTrigger UpkeepType = iota
	LogTrigger
	CronTrigger
	ReadyTrigger
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
//
// TODO: check endianness
func GetUpkeepType(id ocr2keepers.UpkeepIdentifier) UpkeepType {
	if len(id) < upkeepTypeByteIndex+1 {
		return ConditionTrigger
	}
	idx, ok := big.NewInt(0).SetString(string(id), 10)
	if ok {
		id = ocr2keepers.UpkeepIdentifier(idx.Bytes())
	}
	for i := upkeepTypeStartIndex; i < upkeepTypeByteIndex; i++ {
		if id[i] != 0 { // old id
			return ConditionTrigger
		}
	}
	typeByte := id[upkeepTypeByteIndex]
	return UpkeepType(typeByte)
}
