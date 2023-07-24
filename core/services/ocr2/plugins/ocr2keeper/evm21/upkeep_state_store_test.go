package evm

import (
	"fmt"
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/require"
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
	payload1 = ocr2keepers.NewUpkeepPayload(upkeepId1, ConditionalType, BlockKey1, trigger1, []byte{})
	payload2 = ocr2keepers.NewUpkeepPayload(upkeepId1, ConditionalType, BlockKey2, trigger1, []byte{})
	payload3 = ocr2keepers.NewUpkeepPayload(upkeepId2, LogTriggerType, BlockKey2, trigger1, []byte{})
	payload4 = ocr2keepers.NewUpkeepPayload(upkeepId1, LogTriggerType, BlockKey1, trigger2, []byte{})
	payload5 = ocr2keepers.NewUpkeepPayload(upkeepId1, LogTriggerType, BlockKey3, trigger1, []byte{})
)

const (
	ConditionalType = 0
	LogTriggerType  = 1
	Block1          = 111
	Block2          = 112
	Block3          = 113
	BlockKey1       = "111|0x123123132132"
	BlockKey2       = "112|0x565456465465"
	BlockKey3       = "113|0x111423246546"
	InvalidBlockKey = "2220x565456465465"
)

func TestUpkeepStateStore_SetSimpleUpkeepStates(t *testing.T) {
	s := Performed

	tests := []struct {
		name             string
		payloads         []ocr2keepers.UpkeepPayload
		states           []UpkeepState
		expectedError    error
		ids              []string
		idResult         []upkeepState
		upkeepIds        []*big.Int
		upkeepIdsResult  []upkeepState
		blocks           []int64
		blockRangeResult []upkeepState
	}{
		{
			name: "set a single upkeep state",
			payloads: []ocr2keepers.UpkeepPayload{
				payload1,
			},
			states: []UpkeepState{Performed},
			ids:    []string{payload1.ID},
			idResult: []upkeepState{{
				payload: &payload1,
				state:   &s,
			}},
			upkeepIds: []*big.Int{upkeepId1},
			upkeepIdsResult: []upkeepState{
				{
					payload: &payload1,
					state:   &s,
				},
			},
			blocks: []int64{Block1},
			blockRangeResult: []upkeepState{
				{
					payload: &payload1,
					state:   &s,
				},
			},
		},
		{
			name: "sets multiple upkeep states",
			payloads: []ocr2keepers.UpkeepPayload{
				payload2,
				payload3,
				payload4,
			},
			states: []UpkeepState{Performed, Performed, Performed},
			ids:    []string{payload2.ID, payload3.ID, payload4.ID},
			idResult: []upkeepState{
				{
					payload: &payload2,
					state:   &s,
				},
				{
					payload: &payload3,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
			},
			upkeepIds: []*big.Int{upkeepId1, upkeepId2},
			upkeepIdsResult: []upkeepState{
				{
					payload: &payload2,
					state:   &s,
				},
				{
					payload: &payload3,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
			},
			blocks: []int64{Block1, Block2},
			blockRangeResult: []upkeepState{
				{
					payload: &payload2,
					state:   &s,
				},
				{
					payload: &payload3,
					state:   &s,
				},
				{
					payload: &payload4,
					state:   &s,
				},
			},
		},
		{
			name: "failed to split invalid block key",
			payloads: []ocr2keepers.UpkeepPayload{
				ocr2keepers.NewUpkeepPayload(upkeepId2, LogTriggerType, InvalidBlockKey, trigger1, []byte{}),
			},
			states:        []UpkeepState{Performed, Performed, Performed},
			expectedError: fmt.Errorf("check block %s is invalid for upkeep %s", InvalidBlockKey, upkeepId2),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewUpkeepStateStore()
			for i, p := range tc.payloads {
				err := store.SetUpkeepState(p, tc.states[i])
				require.Equal(t, err, tc.expectedError)
			}

			if len(tc.ids) > 0 {
				for i, id := range tc.ids {
					pl, us, err := store.SelectByID(id)
					require.Nil(t, err)
					require.Equal(t, tc.idResult[i].payload, pl)
					require.Equal(t, tc.idResult[i].state, us)
				}
			}

			if len(tc.upkeepIds) > 0 {
				pl, us, err := store.SelectByUpkeepIDs(tc.upkeepIds)
				require.Nil(t, err)
				require.Equal(t, len(tc.upkeepIdsResult), len(pl))
				require.Equal(t, len(tc.upkeepIdsResult), len(us))
				for j, r := range tc.upkeepIdsResult {
					require.Equal(t, r.payload, pl[j])
					require.Equal(t, r.state, us[j])
				}
			}

			if len(tc.blocks) > 0 {
				pl, us, err := store.SelectByBlockRange(tc.blocks[0], tc.blocks[len(tc.blocks)-1]+1)
				require.Nil(t, err)
				require.Equal(t, len(tc.blockRangeResult), len(pl))
				require.Equal(t, len(tc.blockRangeResult), len(us))
				for j, r := range tc.blockRangeResult {
					require.Equal(t, r.payload, pl[j])
					require.Equal(t, r.state, us[j])
				}
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
			store := NewUpkeepStateStore()
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
