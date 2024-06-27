package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/stretchr/testify/assert"
)

func TestWorkID(t *testing.T) {
	tests := []struct {
		name     string
		upkeepID string
		trigger  ocr2keepers.Trigger
		expected string
		errored  bool
	}{
		{
			name:     "happy flow no extension",
			upkeepID: "12345",
			trigger: ocr2keepers.Trigger{
				BlockNumber: 123,
				BlockHash:   common.HexToHash("0xabcdef"),
			},
			expected: "e546b0a52c2879744f6def0fb483d581dc6d205de83af8440456804dd8b62380",
		},
		{
			name:     "empty trigger",
			upkeepID: "12345",
			trigger:  ocr2keepers.Trigger{},
			// same as with no extension
			expected: "e546b0a52c2879744f6def0fb483d581dc6d205de83af8440456804dd8b62380",
		},
		{
			name:     "happy flow with extension",
			upkeepID: genUpkeepID(ocr2keepers.LogTrigger, "12345").String(),
			trigger: ocr2keepers.Trigger{
				BlockNumber: 123,
				BlockHash:   common.HexToHash("0xabcdef"),
				LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
					Index:  1,
					TxHash: common.HexToHash("0x12345"),
				},
			},
			expected: "91ace35299de40860e17d31adbc64bee48f437362cedd3b69ccf749a2f38d8e5",
		},
		{
			name:     "empty upkeepID",
			upkeepID: "0",
			trigger: ocr2keepers.Trigger{
				BlockNumber: 123,
				BlockHash:   common.HexToHash("0xabcdef"),
			},
			expected: "290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Convert the string to a big.Int
			var upkeepID big.Int
			_, success := upkeepID.SetString(tc.upkeepID, 10)
			if !success {
				t.Fatal("Invalid big integer value")
			}

			res, err := UpkeepWorkID(&upkeepID, tc.trigger)
			if tc.errored {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, res, "UpkeepWorkID mismatch")
		})
	}
}

func TestNewUpkeepPayload(t *testing.T) {
	tests := []struct {
		name       string
		upkeepID   *big.Int
		upkeepType ocr2keepers.UpkeepType
		trigger    ocr2keepers.Trigger
		check      []byte
		errored    bool
		workID     string
	}{
		{
			name:       "happy flow no extension",
			upkeepID:   big.NewInt(111),
			upkeepType: ocr2keepers.ConditionTrigger,
			trigger: ocr2keepers.Trigger{
				BlockNumber: 11,
				BlockHash:   common.HexToHash("0x11111"),
			},
			check:  []byte("check-data-111"),
			workID: "39f2babe526038520877fc7c33d81accf578af4a06c5fa6b0d038cae36e12711",
		},
		{
			name:       "happy flow with extension",
			upkeepID:   big.NewInt(111),
			upkeepType: ocr2keepers.LogTrigger,
			trigger: ocr2keepers.Trigger{
				BlockNumber: 11,
				BlockHash:   common.HexToHash("0x11111"),
				LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
					Index:  1,
					TxHash: common.HexToHash("0x11111"),
				},
			},
			check:  []byte("check-data-111"),
			workID: "9fd4d46e09ad25e831fdee664dbaa3b68c37034303234bf70001e3577af43a4f",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := NewUpkeepPayload(
				tc.upkeepID,
				tc.trigger,
				tc.check,
			)
			if tc.errored {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.workID, payload.WorkID)
		})
	}
}
