package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventBufferV1(t *testing.T) {
	t.Run("dequeue empty", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 1, 1)

		results, remaining := buf.Dequeue(int64(1), 20, 1, 10, DefaultUpkeepSelector)
		require.Equal(t, 0, len(results))
		require.Equal(t, 0, remaining)
	})

	t.Run("enqueue", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 1, 1)

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		require.Equal(t, 2, added)
		require.Equal(t, 0, dropped)
		q, ok := buf.(*logBuffer).getUpkeepQueue(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 2, q.sizeOfRange(1, 18))
	})

	t.Run("enqueue upkeeps limits", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 1, 1)
		limit := 2
		buf.(*logBuffer).opts.logLimitHigh.Store(uint32(limit))

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x19"), LogIndex: 0},
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x19"), LogIndex: 1},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 2},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 3},
		)
		totalLimit := limit * 3 // 3 block windows
		require.Equal(t, 7, added)
		require.Equal(t, 7-totalLimit, dropped)
		q, ok := buf.(*logBuffer).getUpkeepQueue(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, totalLimit,
			q.sizeOfRange(1, 18))
	})

	t.Run("enqueue out of block range", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 5, 1, 1)

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x10"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 1},
		)
		require.Equal(t, 2, added)
		require.Equal(t, 0, dropped)
		q, ok := buf.(*logBuffer).getUpkeepQueue(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 2, q.sizeOfRange(1, 12))
	})

	t.Run("happy path", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 20, 1)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(2),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		)
		results, remaining := buf.Dequeue(int64(1), 10, 1, 2, DefaultUpkeepSelector)
		require.Equal(t, 2, len(results))
		require.Equal(t, 2, remaining)
		require.True(t, results[0].ID.Cmp(results[1].ID) != 0)
		results, remaining = buf.Dequeue(int64(1), 10, 1, 2, DefaultUpkeepSelector)
		require.Equal(t, 2, len(results))
		require.Equal(t, 0, remaining)
	})
}

func TestLogEventBufferV1_UpkeepQueue_clean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := newUpkeepLogBuffer(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		q.clean(10)
	})

	t.Run("happy path", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 5, 1)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 0},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 1},
		)

		q, ok := buf.(*logBuffer).getUpkeepQueue(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 4, q.sizeOfRange(1, 11))

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 0},
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 1},
		)

		require.Equal(t, 4, q.sizeOfRange(1, 18))
		require.Equal(t, 0, q.clean(12))
		require.Equal(t, 2, q.sizeOfRange(1, 18))
	})
}
