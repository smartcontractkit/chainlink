package logprovider

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

func TestBlockWindow(t *testing.T) {
	tests := []struct {
		name      string
		block     int64
		blockRate int
		wantStart int64
		wantEnd   int64
	}{
		{
			name:      "block 0, blockRate 1",
			block:     0,
			blockRate: 1,
			wantStart: 0,
			wantEnd:   0,
		},
		{
			name:      "block 81, blockRate 1",
			block:     81,
			blockRate: 1,
			wantStart: 81,
			wantEnd:   81,
		},
		{
			name:      "block 0, blockRate 4",
			block:     0,
			blockRate: 4,
			wantStart: 0,
			wantEnd:   3,
		},
		{
			name:      "block 81, blockRate 4",
			block:     81,
			blockRate: 4,
			wantStart: 80,
			wantEnd:   83,
		},
		{
			name:      "block 83, blockRate 4",
			block:     83,
			blockRate: 4,
			wantStart: 80,
			wantEnd:   83,
		},
		{
			name:      "block 84, blockRate 4",
			block:     84,
			blockRate: 4,
			wantStart: 84,
			wantEnd:   87,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, end := BlockWindow(tc.block, tc.blockRate)
			require.Equal(t, tc.wantStart, start)
			require.Equal(t, tc.wantEnd, end)
		})
	}
}

func TestLogComparatorSorter(t *testing.T) {
	tests := []struct {
		name     string
		a        logpoller.Log
		b        logpoller.Log
		wantCmp  int
		wantSort bool
	}{
		{
			name: "a == b",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			b: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			wantCmp:  0,
			wantSort: false,
		},
		{
			name: "a < b: block number",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			b: logpoller.Log{
				BlockNumber: 2,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			wantCmp:  -1,
			wantSort: false,
		},
		{
			name: "a < b: log index",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			b: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    2,
			},
			wantCmp:  -1,
			wantSort: false,
		},
		{
			name: "a > b: block number",
			a: logpoller.Log{
				BlockNumber: 3,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			b: logpoller.Log{
				BlockNumber: 2,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    1,
			},
			wantCmp:  1,
			wantSort: true,
		},
		{
			name: "a > b: log index",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    3,
			},
			b: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    2,
			},
			wantCmp:  1,
			wantSort: true,
		},
		{
			name: "a > b: tx hash",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x21"),
				LogIndex:    2,
			},
			b: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    2,
			},
			wantCmp:  1,
			wantSort: true,
		},
		{
			name: "a < b: tx hash",
			a: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x1"),
				LogIndex:    2,
			},
			b: logpoller.Log{
				BlockNumber: 1,
				TxHash:      common.HexToHash("0x4"),
				LogIndex:    2,
			},
			wantCmp:  -1,
			wantSort: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantCmp, LogComparator(tc.a, tc.b))
			require.Equal(t, tc.wantSort, LogSorter(tc.a, tc.b))
		})
	}
}
