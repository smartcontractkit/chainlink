package core

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

	assert.Equal(t, "fe466794c97e8b54ca25b696ff3ee448a7d03e4a82a2e45d9d84de62ef4cc260", res, "UpkeepTriggerID mismatch")
}

func TestWorkID(t *testing.T) {
	id := big.NewInt(12345)
	trigger := ocr2keepers.Trigger{
		BlockNumber: 123,
		BlockHash:   common.HexToHash("0xabcdef"),
	}

	res, err := UpkeepWorkID(id, trigger)
	if err != nil {
		t.Fatalf("Error computing UpkeepWorkID: %s", err)
	}

	assert.Equal(t, "6fddaeddb48a3d7a21ef8c22fbd214192dba02f674dc0128ca474aeff038c7c3", res, "UpkeepWorkID mismatch")
}

func TestNewUpkeepPayload(t *testing.T) {
	payload, err := NewUpkeepPayload(
		big.NewInt(111),
		1,
		ocr2keepers.Trigger{
			BlockNumber: 11,
			BlockHash:   common.HexToHash("0x11111"),
			LogTriggerExtension: &ocr2keepers.LogTriggerExtenstion{
				Index:     1,
				LogTxHash: common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
			},
		},
		[]byte("check-data-111"),
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "42e165594d25cdc74d80b2e38cbb0e9857ee6c6a607dfbec1983348087d34593", payload.WorkID)
}
