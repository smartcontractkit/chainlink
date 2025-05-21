package logprovider

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

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
				BlockNumber: 4,
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
				LogIndex:    4,
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
