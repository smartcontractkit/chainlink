package evm

import (
	"math/big"
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestActiveUpkeepList(t *testing.T) {
	logIDs := []ocr2keepers.UpkeepIdentifier{
		core.GenUpkeepID(types.LogTrigger, "0"),
		core.GenUpkeepID(types.LogTrigger, "1"),
		core.GenUpkeepID(types.LogTrigger, "2"),
		core.GenUpkeepID(types.LogTrigger, "3"),
		core.GenUpkeepID(types.LogTrigger, "4"),
	}
	conditionalIDs := []ocr2keepers.UpkeepIdentifier{
		core.GenUpkeepID(types.ConditionTrigger, "0"),
		core.GenUpkeepID(types.ConditionTrigger, "1"),
		core.GenUpkeepID(types.ConditionTrigger, "2"),
		core.GenUpkeepID(types.ConditionTrigger, "3"),
		core.GenUpkeepID(types.ConditionTrigger, "4"),
	}

	tests := []struct {
		name                   string
		initial                []*big.Int
		add                    []*big.Int
		remove                 []*big.Int
		expectedLogIds         []*big.Int
		expectedConditionalIds []*big.Int
	}{
		{
			name:                   "happy flow",
			initial:                []*big.Int{logIDs[0].BigInt(), logIDs[1].BigInt(), conditionalIDs[0].BigInt(), conditionalIDs[1].BigInt()},
			add:                    []*big.Int{logIDs[2].BigInt(), logIDs[3].BigInt(), conditionalIDs[2].BigInt(), conditionalIDs[3].BigInt()},
			remove:                 []*big.Int{logIDs[3].BigInt(), conditionalIDs[3].BigInt()},
			expectedLogIds:         []*big.Int{logIDs[0].BigInt(), logIDs[1].BigInt(), logIDs[2].BigInt()},
			expectedConditionalIds: []*big.Int{conditionalIDs[0].BigInt(), conditionalIDs[1].BigInt(), conditionalIDs[2].BigInt()},
		},
		{
			name:                   "empty",
			initial:                []*big.Int{},
			add:                    []*big.Int{},
			remove:                 []*big.Int{},
			expectedLogIds:         []*big.Int{},
			expectedConditionalIds: []*big.Int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			al := NewActiveUpkeepList()
			al.Reset(tc.initial...)
			require.Equal(t, len(tc.initial), al.Size())
			for _, id := range tc.initial {
				require.True(t, al.IsActive(id))
			}
			al.Add(tc.add...)
			for _, id := range tc.add {
				require.True(t, al.IsActive(id))
			}
			al.Remove(tc.remove...)
			for _, id := range tc.remove {
				require.False(t, al.IsActive(id))
			}
			logIds := al.View(types.LogTrigger)
			require.Equal(t, len(tc.expectedLogIds), len(logIds))
			sort.Slice(logIds, func(i, j int) bool {
				return logIds[i].Cmp(logIds[j]) < 0
			})
			for i := range logIds {
				require.Equal(t, tc.expectedLogIds[i], logIds[i])
			}
			conditionalIds := al.View(types.ConditionTrigger)
			require.Equal(t, len(tc.expectedConditionalIds), len(conditionalIds))
			sort.Slice(conditionalIds, func(i, j int) bool {
				return conditionalIds[i].Cmp(conditionalIds[j]) < 0
			})
			for i := range conditionalIds {
				require.Equal(t, tc.expectedConditionalIds[i], conditionalIds[i])
			}
		})
	}
}

func TestActiveUpkeepList_error(t *testing.T) {
	t.Run("if invalid or negative numbers are in the store, they are excluded from the view operation", func(t *testing.T) {
		al := &activeList{}
		al.items = make(map[string]bool)
		al.items["not a number"] = true
		al.items["-1"] = true
		al.items["100"] = true

		keys := al.View(types.ConditionTrigger)
		require.Equal(t, []*big.Int{big.NewInt(100)}, keys)
	})
}
