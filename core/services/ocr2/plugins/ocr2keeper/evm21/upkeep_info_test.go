package evm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

func TestGetUpkeepType(t *testing.T) {
	tests := []struct {
		name       string
		upkeepID   ocr2keepers.UpkeepIdentifier
		upkeepType upkeepType
	}{
		{
			"zeroed id",
			big.NewInt(0).Bytes(),
			conditionTrigger,
		},
		{
			"old id",
			[]byte("5820911532554020907796191562093071158274499580927271776163559390280294438608"),
			conditionTrigger,
		},
		{
			"condition trigger",
			genUpkeepID(conditionTrigger, "").Bytes(),
			conditionTrigger,
		},
		{
			"log trigger string",
			[]byte(genUpkeepID(logTrigger, "111").String()),
			logTrigger,
		},
		{
			"log trigger",
			genUpkeepID(logTrigger, "111").Bytes(),
			logTrigger,
		},
		{
			"cron trigger",
			genUpkeepID(cronTrigger, "222").Bytes(),
			cronTrigger,
		},
		{
			"ready trigger",
			genUpkeepID(readyTrigger, "333").Bytes(),
			readyTrigger,
		},
		{
			"log trigger id",
			func() ocr2keepers.UpkeepIdentifier {
				id, _ := big.NewInt(0).SetString("32329108151019397958065800113404894502874153543356521479058624064899121404671", 10)
				return id.Bytes()
			}(),
			logTrigger,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.upkeepType, getUpkeepType(tc.upkeepID))
		})
	}
}

func genUpkeepID(uType upkeepType, rand string) *big.Int {
	b := append([]byte{1}, common.LeftPadBytes([]byte{uint8(uType)}, 15)...)
	b = append(b, []byte(rand)...)
	return big.NewInt(0).SetBytes(b)
}
