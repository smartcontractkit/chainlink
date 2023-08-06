package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/stretchr/testify/assert"
)

func TestGetUpkeepType(t *testing.T) {
	tests := []struct {
		name       string
		upkeepID   ocr2keepers.UpkeepIdentifier
		upkeepType ocr2keepers.UpkeepType
	}{
		{
			"zeroed id",
			big.NewInt(0).Bytes(),
			ocr2keepers.ConditionTrigger,
		},
		{
			"old id",
			[]byte("5820911532554020907796191562093071158274499580927271776163559390280294438608"),
			ocr2keepers.ConditionTrigger,
		},
		{
			"condition trigger",
			genUpkeepID(ocr2keepers.ConditionTrigger, "").Bytes(),
			ocr2keepers.ConditionTrigger,
		},
		{
			"log trigger string",
			[]byte(genUpkeepID(ocr2keepers.LogTrigger, "111").String()),
			ocr2keepers.LogTrigger,
		},
		{
			"log trigger",
			genUpkeepID(ocr2keepers.LogTrigger, "111").Bytes(),
			ocr2keepers.LogTrigger,
		},
		{
			"log trigger id",
			func() ocr2keepers.UpkeepIdentifier {
				id, _ := big.NewInt(0).SetString("32329108151019397958065800113404894502874153543356521479058624064899121404671", 10)
				return id.Bytes()
			}(),
			ocr2keepers.LogTrigger,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.upkeepType, GetUpkeepType(tc.upkeepID))
		})
	}
}

func genUpkeepID(uType ocr2keepers.UpkeepType, rand string) *big.Int {
	b := append([]byte{1}, common.LeftPadBytes([]byte{uint8(uType)}, 15)...)
	b = append(b, []byte(rand)...)
	return big.NewInt(0).SetBytes(b)
}
