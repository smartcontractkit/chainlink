package transmit

import (
	"testing"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/require"
)

func TestTransmitEventCache_Sanity(t *testing.T) {
	tests := []struct {
		name        string
		cap         int64
		logIDsToAdd []string
		eventsToAdd []ocr2keepers.TransmitEvent
		toGet       []string
		expected    map[string]ocr2keepers.TransmitEvent
	}{
		{
			"empty cache",
			10,
			[]string{},
			[]ocr2keepers.TransmitEvent{},
			[]string{"1"},
			map[string]ocr2keepers.TransmitEvent{},
		},
		{
			"happy path",
			10,
			[]string{"1", "2", "3", "4"},
			[]ocr2keepers.TransmitEvent{
				{WorkID: "1", TransmitBlock: 1},
				{WorkID: "2", TransmitBlock: 2},
				{WorkID: "3", TransmitBlock: 3},
				{WorkID: "4", TransmitBlock: 4},
			},
			[]string{"1", "3"},
			map[string]ocr2keepers.TransmitEvent{
				"1": {WorkID: "1", TransmitBlock: 1},
				"3": {WorkID: "3", TransmitBlock: 3},
			},
		},
		{
			"overflow",
			3,
			[]string{"1", "2", "3", "4", "5"},
			[]ocr2keepers.TransmitEvent{
				{WorkID: "1", TransmitBlock: 1},
				{WorkID: "2", TransmitBlock: 2},
				{WorkID: "3", TransmitBlock: 3},
				{WorkID: "4", TransmitBlock: 4},
				{WorkID: "5", TransmitBlock: 5},
			},
			[]string{"1", "2", "3", "4", "5"},
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
			for _, logID := range tc.toGet {
				e, exist := c.get(logID)
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
