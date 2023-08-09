package evm

import (
	"math/big"
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

var (
	upkeepId1 = big.NewInt(100)
	upkeepId2 = big.NewInt(200)
	trigger1  = ocr2keepers.Trigger{
		BlockNumber: block1,
		BlockHash:   "0x1231eqwe12eqwd",
		Extension: core.LogTriggerExtension{
			LogIndex: 1,
			TxHash:   "0x1234567890123456789012345678901234567890123456789012345678901234",
		},
	}
	trigger2 = ocr2keepers.Trigger{
		BlockNumber: block3,
		BlockHash:   "0x1231eqwe12eqwd",
		Extension: core.LogTriggerExtension{
			LogIndex: 1,
			TxHash:   "0x1234567890123456789012345678901234567890123456789012345678901234",
		},
	}
	payload2, _ = core.NewUpkeepPayload(upkeepId1, conditionalType, trigger1, []byte{})
	payload3, _ = core.NewUpkeepPayload(upkeepId2, logTriggerType, trigger1, []byte{})
	payload4, _ = core.NewUpkeepPayload(upkeepId1, logTriggerType, trigger2, []byte{})
	payload5, _ = core.NewUpkeepPayload(upkeepId1, logTriggerType, trigger1, []byte{})
)

const (
	conditionalType = 0
	logTriggerType  = 1
	block1          = 111
	block3          = 113
)

func TestUpkeepStateStore_OverrideUpkeepStates(t *testing.T) {
	p := ocr2keepers.Performed
	e := ocr2keepers.Eligible

	tests := []struct {
		name          string
		payloads      []ocr2keepers.UpkeepPayload
		states        []ocr2keepers.UpkeepState
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
			states: []ocr2keepers.UpkeepState{ocr2keepers.Performed, ocr2keepers.Performed, ocr2keepers.Performed, ocr2keepers.Eligible},
			oldIds: []string{payload2.ID, payload3.ID, payload4.ID},
			oldIdResult: []upkeepState{
				{
					payload: &payload3,
					state:   &p,
				},
				{
					payload: &payload4,
					state:   &p,
				},
			},
			newIds: []string{payload3.ID, payload4.ID, payload5.ID},
			newIdResult: []upkeepState{
				{
					payload: &payload3,
					state:   &p,
				},
				{
					payload: &payload4,
					state:   &p,
				},
				{
					payload: &payload5,
					state:   &e,
				},
			},

			upkeepIds:  []*big.Int{upkeepId1},
			endBlock:   block3 + 1,
			startBlock: block1,
			result: []upkeepState{
				{
					payload: &payload5,
					state:   &e,
				},
				{
					payload: &payload4,
					state:   &p,
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
