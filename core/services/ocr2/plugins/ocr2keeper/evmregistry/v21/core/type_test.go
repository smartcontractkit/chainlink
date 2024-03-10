package core

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/stretchr/testify/assert"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func TestGetUpkeepType(t *testing.T) {
	tests := []struct {
		name       string
		upkeepID   []byte
		upkeepType types.UpkeepType
	}{
		{
			"zeroed id",
			big.NewInt(0).Bytes(),
			types.ConditionTrigger,
		},
		{
			"old id",
			[]byte("5820911532554020907796191562093071158274499580927271776163559390280294438608"),
			types.ConditionTrigger,
		},
		{
			"condition trigger",
			GenUpkeepID(types.ConditionTrigger, "").BigInt().Bytes(),
			types.ConditionTrigger,
		},
		{
			"log trigger",
			GenUpkeepID(types.LogTrigger, "111").BigInt().Bytes(),
			types.LogTrigger,
		},
		{
			"log trigger id",
			func() []byte {
				id, _ := big.NewInt(0).SetString("32329108151019397958065800113404894502874153543356521479058624064899121404671", 10)
				return id.Bytes()
			}(),
			types.LogTrigger,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uid := ocr2keepers.UpkeepIdentifier{}
			copy(uid[:], tc.upkeepID)
			assert.Equal(t, tc.upkeepType, GetUpkeepType(uid))
		})
	}
}
