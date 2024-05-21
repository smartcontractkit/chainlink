package vrfcommon

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestLogDeduper(t *testing.T) {
	tests := []struct {
		name    string
		logs    []types.Log
		results []bool
	}{
		{
			name: "dupe",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
			},
			results: []bool{true, false},
		},
		{
			name: "different block number",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
				{
					BlockNumber: 11,
					BlockHash:   common.Hash{0x2},
					Index:       3,
				},
			},
			results: []bool{true, true},
		},
		{
			name: "same block number different hash",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x2},
					Index:       3,
				},
			},
			results: []bool{true, true},
		},
		{
			name: "same block number same hash different index",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       4,
				},
			},
			results: []bool{true, true},
		},
		{
			name: "same block number same hash different index",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x1},
					Index:       4,
				},
			},
			results: []bool{true, true},
		},
		{
			name: "multiple blocks with dupes",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       4,
				},
				{
					BlockNumber: 11,
					BlockHash:   common.Hash{0x11},
					Index:       0,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       3,
				},
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       4,
				},
				{
					BlockNumber: 12,
					BlockHash:   common.Hash{0x12},
					Index:       1,
				},
			},
			results: []bool{true, true, true, false, false, true},
		},
		{
			name: "prune",
			logs: []types.Log{
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       3,
				},
				{
					BlockNumber: 11,
					BlockHash:   common.Hash{0x11},
					Index:       11,
				},
				{
					BlockNumber: 1015,
					BlockHash:   common.Hash{0x1, 0x1, 0x5},
					Index:       0,
				},
				// Now the logs at blocks 10 and 11 should be pruned, and therefore redelivered.
				// The log at block 115 should not be redelivered.
				{
					BlockNumber: 10,
					BlockHash:   common.Hash{0x10},
					Index:       3,
				},
				{
					BlockNumber: 11,
					BlockHash:   common.Hash{0x11},
					Index:       11,
				},
				{
					BlockNumber: 1015,
					BlockHash:   common.Hash{0x1, 0x1, 0x5},
					Index:       0,
				},
			},
			results: []bool{true, true, true, true, true, false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deduper := NewLogDeduper(100)

			for i := range test.logs {
				require.Equal(t, test.results[i], deduper.ShouldDeliver(test.logs[i]),
					"expected shouldDeliver for log %d to be %t", i, test.results[i])
			}
		})
	}
}
