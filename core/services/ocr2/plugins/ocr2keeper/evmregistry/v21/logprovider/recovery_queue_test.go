package logprovider

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRecoveryQueue(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		maxSize                  int
		maxSizePerUpkeep         int
		recordsToAdd             []automation.UpkeepPayload
		hasIDs                   []string
		expectedRemoved          []automation.UpkeepPayload
		expectedRemainingRecords []string
		expectedHas              []bool
	}{
		{
			name:             "First two records are removed",
			maxSize:          2,
			maxSizePerUpkeep: 1,
			recordsToAdd: []automation.UpkeepPayload{
				{
					UpkeepID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					WorkID:   "workID0",
				},
				{
					UpkeepID: [32]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
					WorkID:   "workID1",
				},
				{
					UpkeepID: [32]byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
					WorkID:   "workID2",
				},
				{
					UpkeepID: [32]byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
					WorkID:   "workID3",
				},
			},
			hasIDs: []string{
				"workID2",
				"workID3",
			},
			expectedHas: []bool{
				true,
				true,
			},
			expectedRemoved: []automation.UpkeepPayload{
				{
					UpkeepID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					WorkID:   "workID0",
				},
				{
					UpkeepID: [32]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
					WorkID:   "workID1",
				},
			},
			expectedRemainingRecords: []string{
				"workID2",
				"workID3",
			},
		},
		{
			name:             "Second record is skipped because we've reached the per upkeep limit",
			maxSize:          2,
			maxSizePerUpkeep: 1,
			recordsToAdd: []automation.UpkeepPayload{
				{
					UpkeepID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					WorkID:   "workID0",
				},
				{
					UpkeepID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					WorkID:   "workID1",
				},
				{
					UpkeepID: [32]byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
					WorkID:   "workID2",
				},
				{
					UpkeepID: [32]byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
					WorkID:   "workID3",
				},
			},
			hasIDs: []string{
				"workID20",
				"workID3",
			},
			expectedHas: []bool{
				false,
				true,
			},
			expectedRemoved: []automation.UpkeepPayload{
				{
					UpkeepID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					WorkID:   "workID0",
				},
				{
					UpkeepID: [32]byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
					WorkID:   "workID2",
				},
			},
			expectedRemainingRecords: []string{
				"workID1",
				"workID3",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recoveryQueue := NewRecoveryQueue(tc.maxSize, tc.maxSizePerUpkeep)

			recoveryQueue.add(tc.recordsToAdd...)

			hasResults := recoveryQueue.has(tc.hasIDs...)
			assert.Equal(t, hasResults, tc.expectedHas)

			res, err := recoveryQueue.getPayloads()

			assert.NoError(t, err)
			assert.Equal(t, res, tc.expectedRemoved)
			assert.Equal(t, recoveryQueue.queue, tc.expectedRemainingRecords)
		})

	}

}
