package logprovider

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventBuffer_GetBlocksInRange(t *testing.T) {
	size := 3
	buf := newLogEventBuffer(logger.TestLogger(t), size, 10, 10)

	buf.enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	)

	buf.enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 2},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 2},
	)

	tests := []struct {
		name string
		from int64
		to   int64
		want int
	}{
		{
			name: "all",
			from: 1,
			to:   3,
			want: 3,
		},
		{
			name: "partial",
			from: 1,
			to:   2,
			want: 2,
		},
		{
			name: "circular",
			from: 2,
			to:   4,
			want: 3,
		},
		{
			name: "zero start",
			from: 0,
			to:   2,
			want: 2,
		},
		{
			name: "invalid zero end",
			from: 0,
			to:   0,
		},
		{
			name: "invalid from larger than to",
			from: 4,
			to:   2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blocks := buf.getBlocksInRange(int(tc.from), int(tc.to))
			require.Equal(t, tc.want, len(blocks))
			if tc.want > 0 {
				from := tc.from
				if from == 0 {
					from++
				}
				require.Equal(t, from, blocks[0].blockNumber)
				to := tc.to
				if to == 0 {
					to++
				} else if to > int64(size) {
					to = to % int64(size)
				}
				require.Equal(t, to, blocks[len(blocks)-1].blockNumber)
			}
		})
	}
}

func TestLogEventBuffer_GetBlocksInRange_Circular(t *testing.T) {
	size := 4
	buf := newLogEventBuffer(logger.TestLogger(t), size, 10, 10)

	require.Equal(t, buf.enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	), 3)

	require.Equal(t, buf.enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 2},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 2},
		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	), 3)

	require.Equal(t, buf.enqueue(big.NewInt(3),
		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 4},
		logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x3"), LogIndex: 2},
		logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x3"), LogIndex: 5},
	), 3)

	tests := []struct {
		name           string
		from           int64
		to             int64
		expectedBlocks []int64
	}{
		{
			name:           "happy flow",
			from:           2,
			to:             5,
			expectedBlocks: []int64{2, 3, 4, 5},
		},
		{
			name:           "range overflow circular",
			from:           1,
			to:             6,
			expectedBlocks: []int64{2, 3, 4, 5},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blocks := buf.getBlocksInRange(int(tc.from), int(tc.to))
			require.Equal(t, len(tc.expectedBlocks), len(blocks))
			expectedBlockNumbers := map[int64]bool{}
			for _, b := range tc.expectedBlocks {
				expectedBlockNumbers[b] = false
			}
			for _, b := range blocks {
				expectedBlockNumbers[b.blockNumber] = true
			}
			for k, v := range expectedBlockNumbers {
				require.True(t, v, "missing block %d", k)
			}
		})
	}
}

func TestLogEventBuffer_EnqueueDequeue(t *testing.T) {
	t.Run("dequeue empty", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

		results := buf.peekRange(int64(1), int64(2))
		require.Equal(t, 0, len(results))
		results = buf.peek(2)
		require.Equal(t, 0, len(results))
	})

	t.Run("enqueue", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

		buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue logs overflow", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 2, 10)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		))
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue block overflow", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 2, 10)

		require.Equal(t, 5, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 0},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 1},
		))
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue upkeep block overflow", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 10, 10, 2)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 3},
		))
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("peek range after dequeue", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

		require.Equal(t, buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		), 2)
		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		), 2)
		results := buf.peekRange(int64(1), int64(2))
		require.Equal(t, 2, len(results))
		verifyBlockNumbers(t, results, 1, 2)
		removed := buf.dequeueRange(int64(1), int64(2))
		require.Equal(t, 2, len(removed))
		results = buf.peekRange(int64(1), int64(2))
		require.Equal(t, 0, len(results))
	})

	t.Run("enqueue peek and dequeue", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

		require.Equal(t, buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		), 2)
		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		), 2)
		results := buf.peek(8)
		require.Equal(t, 4, len(results))
		verifyBlockNumbers(t, results, 1, 2, 3, 3)
		removed := buf.dequeue(8)
		require.Equal(t, 4, len(removed))
		buf.lock.Lock()
		require.Equal(t, 0, len(buf.blocks[0].logs))
		require.Equal(t, 1, len(buf.blocks[1].visited))
		buf.lock.Unlock()
	})

	t.Run("enqueue and peek range circular", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
		), 3)
		require.Equal(t, buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		), 2)

		results := buf.peekRange(int64(1), int64(1))
		require.Equal(t, 0, len(results))

		results = buf.peekRange(int64(3), int64(5))
		require.Equal(t, 3, len(results))
		verifyBlockNumbers(t, results, 3, 4, 4)
	})

	t.Run("doesnt enqueue old blocks", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)

		require.Equal(t, buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		), 2)
		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
		), 2)
		results := buf.peekRange(int64(0), int64(5))
		fmt.Println(results)
		verifyBlockNumbers(t, results, 2, 3, 4, 4)
	})
}

func verifyBlockNumbers(t *testing.T, logs []fetchedLog, bns ...int64) {
	require.Equal(t, len(bns), len(logs), "expected length mismatch")
	for i, log := range logs {
		require.Equal(t, bns[i], log.log.BlockNumber, "wrong block number")
	}
}
