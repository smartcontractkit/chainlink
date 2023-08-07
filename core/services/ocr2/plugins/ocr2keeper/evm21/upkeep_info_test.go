package evm

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

func TestTriggerID(t *testing.T) {
	upkeepIDStr := "82566255084862886500628610724995377215109748679571001950554849251333329872882"
	// Convert the string to a big.Int
	var upkeepID big.Int
	_, success := upkeepID.SetString(upkeepIDStr, 10)
	if !success {
		t.Fatal("Invalid big integer value")
	}

	triggerStr := "deadbeef"
	triggerBytes, err := hex.DecodeString(triggerStr)
	if err != nil {
		t.Fatalf("Error decoding hex string: %s", err)
	}

	res := UpkeepTriggerID(&upkeepID, triggerBytes)

	expectedResult := "fe466794c97e8b54ca25b696ff3ee448a7d03e4a82a2e45d9d84de62ef4cc260"
	assert.Equal(t, res, expectedResult, "UpkeepTriggerID mismatch")
}

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
