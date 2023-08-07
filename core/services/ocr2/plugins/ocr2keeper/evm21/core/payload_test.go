package core

import (
	"encoding/hex"
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/stretchr/testify/assert"
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

func TestWorkID(t *testing.T) {
	id := big.NewInt(12345)
	trigger := ocr2keepers.Trigger{
		BlockNumber: 123,
		BlockHash:   "0xabcdef",
		Extension: map[string]interface{}{
			"txHash":   "0xdeadbeef",
			"logIndex": "123",
		},
	}

	res, err := UpkeepWorkID(id, trigger)
	if err != nil {
		t.Fatalf("Error computing UpkeepWorkID: %s", err)
	}

	expectedResult := "182244f3e4fa6e9556a17b84718b7ffe19f0d67a811d771bd5479fe16e02bf82"
	assert.Equal(t, res, expectedResult, "UpkeepWorkID mismatch")
}

func TestNewUpkeepPayload(t *testing.T) {
	logExt := LogTriggerExtension{
		LogIndex: 1,
		TxHash:   "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	payload, err := NewUpkeepPayload(
		big.NewInt(111),
		1,
		ocr2keepers.Trigger{
			BlockNumber: 11,
			BlockHash:   "0x11111",
			Extension:   logExt,
		},
		[]byte("check-data-111"),
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "151b4fd0c9f8d97307938fb7d64d6efbc821891a0eb9ae2708ccfe98e1706b40", payload.ID)
}
