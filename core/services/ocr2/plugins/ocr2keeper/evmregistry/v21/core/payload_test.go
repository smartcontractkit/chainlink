package core

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func TestWorkID(t *testing.T) {
	tests := []struct {
		name     string
		upkeepID string
		trigger  ocr2keepers.Trigger
		expected string
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
			upkeepID: GenUpkeepID(types.LogTrigger, "12345").String(),
			trigger: ocr2keepers.Trigger{
				BlockNumber: 123,
				BlockHash:   common.HexToHash("0xabcdef"),
				LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
					Index:     1,
					TxHash:    common.HexToHash("0x12345"),
					BlockHash: common.HexToHash("0xabcdef"),
				},
			},
			expected: "aaa208331dfafff7a681e3358d082a2e78633dd05c8fb2817c391888cadb2912",
		},
		{
			name:     "happy path example from an actual tx",
			upkeepID: "57755329819103678328139927896464733492677608573736038892412245689671711489918",
			trigger: ocr2keepers.Trigger{
				BlockNumber: 39344455,
				BlockHash:   common.HexToHash("0xb41258d18cd44ebf7a0d70de011f2bc4a67c9b68e8b6dada864045d8543bb020"),
				LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
					Index:     41,
					TxHash:    common.HexToHash("0x44079b1b33aff337dbf17b9e12c5724ecab979c50c8201a9814a488ff3e22384"),
					BlockHash: common.HexToHash("0xb41258d18cd44ebf7a0d70de011f2bc4a67c9b68e8b6dada864045d8543bb020"),
				},
			},
			expected: "ef1b6acac8aa3682a8a08f666a13cfa165f7e811a16ea9fa0817f437fc4d110d",
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
			var id big.Int
			_, success := id.SetString(tc.upkeepID, 10)
			if !success {
				t.Fatal("Invalid big integer value")
			}
			uid := &ocr2keepers.UpkeepIdentifier{}
			ok := uid.FromBigInt(&id)
			if !ok {
				t.Fatal("Invalid upkeep identifier")
			}

			res := UpkeepWorkID(*uid, tc.trigger)
			assert.Equal(t, tc.expected, res, "UpkeepWorkID mismatch")
		})
	}
}

func TestNewUpkeepPayload(t *testing.T) {
	tests := []struct {
		name       string
		upkeepID   *big.Int
		upkeepType types.UpkeepType
		trigger    ocr2keepers.Trigger
		check      []byte
		errored    bool
		workID     string
	}{
		{
			name:       "happy flow no extension",
			upkeepID:   big.NewInt(111),
			upkeepType: types.ConditionTrigger,
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
			upkeepType: types.LogTrigger,
			trigger: ocr2keepers.Trigger{
				BlockNumber: 11,
				BlockHash:   common.HexToHash("0x11111"),
				LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
					Index:  1,
					TxHash: common.HexToHash("0x11111"),
				},
			},
			check:  []byte("check-data-111"),
			workID: "d8e7c8907a0b60b637ce71ff4f757edf076e270d52c51f6e4d46a3b0696e0a39",
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
