package transmit

import (
	"testing"

	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func TestTransmitEventCache_Sanity(t *testing.T) {
	tests := []struct {
		name        string
		cap         int64
		logIDsToAdd []string
		eventsToAdd []ocr2keepers.TransmitEvent
		toGet       []string
		blocksToGet []int64
		expected    map[string]ocr2keepers.TransmitEvent
	}{
		{
			"empty cache",
			10,
			[]string{},
			[]ocr2keepers.TransmitEvent{},
			[]string{"1"},
			[]int64{1},
			map[string]ocr2keepers.TransmitEvent{},
		},
		{
			"happy path",
			10,
			[]string{"3", "2", "4", "1"},
			[]ocr2keepers.TransmitEvent{
				{WorkID: "3", TransmitBlock: 3},
				{WorkID: "2", TransmitBlock: 2},
				{WorkID: "4", TransmitBlock: 4},
				{WorkID: "1", TransmitBlock: 1},
			},
			[]string{"1", "3"},
			[]int64{1, 3},
			map[string]ocr2keepers.TransmitEvent{
				"1": {WorkID: "1", TransmitBlock: 1},
				"3": {WorkID: "3", TransmitBlock: 3},
			},
		},
		{
			"different blocks",
			10,
			[]string{"3", "1", "2", "4"},
			[]ocr2keepers.TransmitEvent{
				{WorkID: "3", TransmitBlock: 3},
				{WorkID: "1", TransmitBlock: 1},
				{WorkID: "2", TransmitBlock: 2},
				{WorkID: "4", TransmitBlock: 4},
			},
			[]string{"1", "3"},
			[]int64{9, 9},
			map[string]ocr2keepers.TransmitEvent{},
		},
		{
			"overflow",
			3,
			[]string{"4", "1", "3", "2", "5"},
			[]ocr2keepers.TransmitEvent{
				{WorkID: "4", TransmitBlock: 4},
				{WorkID: "1", TransmitBlock: 1},
				{WorkID: "3", TransmitBlock: 3},
				{WorkID: "2", TransmitBlock: 2},
				{WorkID: "5", TransmitBlock: 5},
			},
			[]string{"1", "4", "2", "3", "5"},
			[]int64{1, 4, 2, 3, 5},
			map[string]ocr2keepers.TransmitEvent{
				"3": {WorkID: "3", TransmitBlock: 3},
				"4": {WorkID: "4", TransmitBlock: 4},
				"5": {WorkID: "5", TransmitBlock: 5},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newTransmitEventCache(tc.cap)
			require.Equal(t, len(tc.eventsToAdd), len(tc.logIDsToAdd))
			for i, e := range tc.eventsToAdd {
				c.add(tc.logIDsToAdd[i], e)
			}
			require.Equal(t, len(tc.toGet), len(tc.blocksToGet))
			for i, logID := range tc.toGet {
				e, exist := c.get(ocr2keepers.BlockNumber(tc.blocksToGet[i]), logID)
				expected, ok := tc.expected[logID]
				if !ok {
					require.False(t, exist, "expected not to find logID %s", logID)
					continue
				}
				require.True(t, exist, "expected to find logID %s", logID)
				require.Equal(t, expected.WorkID, e.WorkID)
			}
		})
	}
}
