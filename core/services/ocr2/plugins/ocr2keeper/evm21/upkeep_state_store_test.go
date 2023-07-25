package evm

import (
	"fmt"
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	upkeepId1 = big.NewInt(100)
	upkeepId2 = big.NewInt(200)
	trigger1  = ocr2keepers.Trigger{
		BlockNumber: 95,
		BlockHash:   "0x1231eqwe12eqwd",
	}
	trigger2 = ocr2keepers.Trigger{
		BlockNumber: 125,
		BlockHash:   "0x1231eqwe12eqwd",
	}
	payload2 = ocr2keepers.NewUpkeepPayload(upkeepId1, ConditionalType, BlockKey2, trigger1, []byte{})
	payload3 = ocr2keepers.NewUpkeepPayload(upkeepId2, LogTriggerType, BlockKey2, trigger1, []byte{})
	payload4 = ocr2keepers.NewUpkeepPayload(upkeepId1, LogTriggerType, BlockKey1, trigger2, []byte{})
	payload5 = ocr2keepers.NewUpkeepPayload(upkeepId1, LogTriggerType, BlockKey3, trigger1, []byte{})
)

const (
	ConditionalType = 0
	LogTriggerType  = 1
	Block1          = 111
	Block3          = 113
	BlockKey1       = "111|0x123123132132"
	BlockKey2       = "112|0x565456465465"
	BlockKey3       = "113|0x111423246546"
	InvalidBlockKey = "2220x565456465465"
)

func TestUpkeepStateStore_InvalidBlockKey(t *testing.T) {
	tests := []struct {
		name          string
		payloads      []ocr2keepers.UpkeepPayload
		states        []UpkeepState
		expectedError error
	}{
		{
			name: "failed to split invalid block key",
			payloads: []ocr2keepers.UpkeepPayload{
				ocr2keepers.NewUpkeepPayload(upkeepId2, LogTriggerType, InvalidBlockKey, trigger1, []byte{}),
			},
			states:        []UpkeepState{Performed},
			expectedError: fmt.Errorf("check block %s is invalid for upkeep %s", InvalidBlockKey, upkeepId2),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewUpkeepStateStore(logger.TestLogger(t))
			for i, p := range tc.payloads {
				err := store.SetUpkeepState(p, tc.states[i])
				require.Equal(t, err, tc.expectedError)
			}
		})
	}
}

func TestUpkeepStateStore_OverrideUpkeepStates(t *testing.T) {
	s := Performed

	tests := []struct {
		name          string
		payloads      []ocr2keepers.UpkeepPayload
		states        []UpkeepState
		expectedError error
		oldIds        []string
		oldIdResult   []upkeepState
		newIds        []string
		newIdResult   []upkeepState
		upkeepIds     []*big.Int
		endBlock      int64
		startBlock    int64
		result        []upkeepState
	}{
		{
			name: "overrides existing upkeep states",
			payloads: []ocr2keepers.UpkeepPayload{
				payload2,
				payload3,
				payload4,
				payload5, // this overrides payload 2 bc they have the same payload ID
			},
			states: []UpkeepState{Performed, Performed, Performed, Performed},
			oldIds: []string{payload2.ID, payload3.ID, payload4.ID},
			oldIdResult: []upkeepState{
				{
					payload: &payload3,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
			},
			newIds: []string{payload3.ID, payload4.ID, payload5.ID},
			newIdResult: []upkeepState{
				{
					payload: &payload3,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
				{
					payload: &payload5,
					state:   &s,
				},
			},

			upkeepIds:  []*big.Int{upkeepId1},
			endBlock:   Block3 + 1,
			startBlock: Block1,
			result: []upkeepState{
				{
					payload: &payload5,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewUpkeepStateStore(logger.TestLogger(t))
			for i, p := range tc.payloads {
				err := store.SetUpkeepState(p, tc.states[i])
				require.Equal(t, err, tc.expectedError)
			}

			pl, us, err := store.SelectByUpkeepIDsAndBlockRange(tc.upkeepIds, tc.startBlock, tc.endBlock)
			require.Nil(t, err)
			require.Equal(t, len(tc.result), len(pl))
			require.Equal(t, len(tc.result), len(us))
			for j, r := range tc.result {
				require.Equal(t, r.payload, pl[j])
				require.Equal(t, r.state, us[j])
			}

		})
	}
}
