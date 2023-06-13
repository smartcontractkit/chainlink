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
			common.LeftPadBytes([]byte{0}, 16),
			conditionTrigger,
		},
		{
			"log trigger",
			common.LeftPadBytes([]byte{1}, 16),
			logTrigger,
		},
		{
			"cron trigger",
			common.LeftPadBytes([]byte{2}, 16),
			cronTrigger,
		},
		{
			"ready trigger",
			common.LeftPadBytes([]byte{3}, 16),
			readyTrigger,
		},
		{
			"log trigger id",
			func() ocr2keepers.UpkeepIdentifier {
				id, ok := big.NewInt(0).SetString("32329108151019397958065800113404894502874153543356521479058624064899121404671", 10)
				if !ok {
					panic("failed to parse id")
				}
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
