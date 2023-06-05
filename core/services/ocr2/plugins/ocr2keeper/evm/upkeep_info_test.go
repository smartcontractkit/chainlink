package evm

import (
	"fmt"
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
			common.LeftPadBytes([]byte{1}, 16),
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
			cronTrigger,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(len(tc.upkeepID))
			assert.Equal(t, tc.upkeepType, getUpkeepType(tc.upkeepID))
		})
	}
}
