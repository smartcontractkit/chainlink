package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventBufferV1_Clean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		buf := newUpkeepLogBuffer(logger.TestLogger(t), big.NewInt(1), 10)

		buf.clean(10)
	})

	t.Run("happy path", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 10)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 0},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 1},
		)

		upkeepBuf, ok := buf.(*logBuffer).getUpkeepBuffer(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 4, upkeepBuf.size())

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 0},
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 1},
		)

		require.Equal(t, 4, upkeepBuf.size())
		require.Equal(t, 0, upkeepBuf.clean(12))
		require.Equal(t, 2, upkeepBuf.size())
	})
}

func TestLogEventBufferV1_EnqueueDequeue(t *testing.T) {
	t.Run("dequeue empty", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 10)

		results, remaining := buf.Dequeue(int64(1), 20, 1, 10, DefaultUpkeepSelector)
		require.Equal(t, 0, len(results))
		require.Equal(t, 0, remaining)
	})

	t.Run("enqueue", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 10)

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		require.Equal(t, 2, added)
		require.Equal(t, 0, dropped)
		upkeepBuf, ok := buf.(*logBuffer).getUpkeepBuffer(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 2, upkeepBuf.size())
	})

	t.Run("enqueue upkeeps limits", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 3, 2)

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 9, TxHash: common.HexToHash("0x9"), LogIndex: 0},
			logpoller.Log{BlockNumber: 9, TxHash: common.HexToHash("0x9"), LogIndex: 1},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 2},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 3},
		)
		require.Equal(t, 7, added)
		require.Equal(t, 1, dropped)
		upkeepBuf, ok := buf.(*logBuffer).getUpkeepBuffer(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 6, upkeepBuf.size())
	})

	t.Run("enqueue out of block range", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 5, 4)

		added, dropped := buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x10"), LogIndex: 0},
			logpoller.Log{BlockNumber: 10, TxHash: common.HexToHash("0x10"), LogIndex: 1},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x11"), LogIndex: 1},
		)
		require.Equal(t, 2, added)
		require.Equal(t, 0, dropped)
		upkeepBuf, ok := buf.(*logBuffer).getUpkeepBuffer(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 2, upkeepBuf.size())
	})

	t.Run("enqueue dequeue", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 10)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(2),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		)
		results, remaining := buf.Dequeue(int64(1), 20, 1, 2, DefaultUpkeepSelector)
		require.Equal(t, 2, len(results))
		require.Equal(t, 2, remaining)
		require.True(t, results[0].ID.Cmp(results[1].ID) != 0)
		results, remaining = buf.Dequeue(int64(1), 20, 1, 2, DefaultUpkeepSelector)
		require.Equal(t, 2, len(results))
		require.Equal(t, 0, remaining)
	})

	// t.Run("enqueue logs overflow", func(t *testing.T) {
	// 	buf := NewLogBuffer(logger.TestLogger(t), 2)

	// 	require.Equal(t, 2, buf.Enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	// 	))
	// 	upkeepBuf, ok := buf.(*logBuffer).getUpkeepBuffer(big.NewInt(1))
	// 	require.True(t, ok)
	// 	require.Equal(t, 2, upkeepBuf.len())
	// })

	// t.Run("enqueue dequeue with dynamic limits", func(t *testing.T) {
	// 	buf := NewLogBuffer(logger.TestLogger(t), 2)

	// 	require.Equal(t, 3, buf.Enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	// 	))
	// 	results := buf.Dequeue(int64(1), int64(20), 1, 2)
	// 	require.Equal(t, 2, len(results))
	// 	buf.SetConfig(10, 3)
	// 	require.Equal(t, 4, buf.Enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 15, TxHash: common.HexToHash("0x21"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 15, TxHash: common.HexToHash("0x21"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 15, TxHash: common.HexToHash("0x21"), LogIndex: 2},
	// 		logpoller.Log{BlockNumber: 15, TxHash: common.HexToHash("0x21"), LogIndex: 3},
	// 	))

	// 	results = buf.Dequeue(int64(1), int64(20), 1, 4)
	// 	require.Equal(t, 3, len(results))

	// 	for _, r := range results {
	// 		require.Equal(t, int64(15), r.Log.BlockNumber)
	// 	}
	// })

	// t.Run("enqueue logs overflow with dynamic limits", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 2, 10, 2)

	// 	require.Equal(t, 2, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 3},
	// 	))
	// 	buf.SetLimits(10, 3)
	// 	require.Equal(t, 3, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 2},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 3},
	// 	))

	// 	buf.lock.Lock()
	// 	defer buf.lock.Unlock()
	// 	require.Equal(t, 2, len(buf.blocks[0].logs))
	// })

	// t.Run("enqueue block overflow", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 2, 10)

	// 	require.Equal(t, 5, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 1},
	// 	))
	// 	buf.lock.Lock()
	// 	require.Equal(t, 2, len(buf.blocks[0].logs))
	// 	buf.lock.Unlock()
	// })

	// t.Run("enqueue upkeep block overflow", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 10, 10, 2)

	// 	require.Equal(t, 2, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 3},
	// 	))
	// 	buf.lock.Lock()
	// 	require.Equal(t, 2, len(buf.blocks[0].logs))
	// 	buf.lock.Unlock()
	// })

	// t.Run("peek range after dequeue", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

	// 	require.Equal(t, buf.enqueue(big.NewInt(10),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
	// 	), 2)
	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 	), 2)
	// 	results := buf.peekRange(int64(1), int64(2))
	// 	require.Equal(t, 2, len(results))
	// 	verifyBlockNumbers(t, results, 1, 2)
	// 	removed := buf.dequeueRange(int64(1), int64(2), 2, 10)
	// 	require.Equal(t, 2, len(removed))
	// 	results = buf.peekRange(int64(1), int64(2))
	// 	require.Equal(t, 0, len(results))
	// })

	// t.Run("enqueue peek and dequeue", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 4, 10, 10)

	// 	require.Equal(t, buf.enqueue(big.NewInt(10),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
	// 	), 2)
	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	// 	), 2)
	// 	results := buf.peek(8)
	// 	require.Equal(t, 4, len(results))
	// 	verifyBlockNumbers(t, results, 1, 2, 3, 3)
	// 	removed := buf.dequeueRange(1, 3, 5, 5)
	// 	require.Equal(t, 4, len(removed))
	// 	buf.lock.Lock()
	// 	require.Equal(t, 0, len(buf.blocks[0].logs))
	// 	require.Equal(t, int64(2), buf.blocks[1].blockNumber)
	// 	require.Equal(t, 1, len(buf.blocks[1].visited))
	// 	buf.lock.Unlock()
	// })

	// t.Run("enqueue and peek range circular", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 10, 10)

	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	// 	), 3)
	// 	require.Equal(t, buf.enqueue(big.NewInt(10),
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 10},
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 11},
	// 	), 2)

	// 	results := buf.peekRange(int64(1), int64(1))
	// 	require.Equal(t, 0, len(results))

	// 	results = buf.peekRange(int64(3), int64(5))
	// 	require.Equal(t, 3, len(results))
	// 	verifyBlockNumbers(t, results, 3, 4, 4)
	// })

	// t.Run("doesnt enqueue old blocks", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)

	// 	require.Equal(t, buf.enqueue(big.NewInt(10),
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 10},
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 11},
	// 	), 2)
	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	// 	), 2)
	// 	results := buf.peekRange(int64(1), int64(5))
	// 	fmt.Println(results)
	// 	verifyBlockNumbers(t, results, 2, 3, 4, 4)
	// })

	// t.Run("dequeue with limits returns latest block logs", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)
	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x5"), LogIndex: 0},
	// 	), 5)

	// 	logs := buf.dequeueRange(1, 5, 2, 10)
	// 	require.Equal(t, 2, len(logs))
	// 	require.Equal(t, int64(5), logs[0].log.BlockNumber)
	// 	require.Equal(t, int64(4), logs[1].log.BlockNumber)

	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 1},
	// 		logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x5"), LogIndex: 1},
	// 	), 2)

	// 	logs = buf.dequeueRange(1, 5, 3, 2)
	// 	require.Equal(t, 2, len(logs))
	// })

	// t.Run("dequeue doesn't return same logs again", func(t *testing.T) {
	// 	buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)
	// 	require.Equal(t, buf.enqueue(big.NewInt(1),
	// 		logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
	// 		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	// 	), 3)

	// 	logs := buf.dequeueRange(3, 3, 2, 10)
	// 	fmt.Println(logs)
	// 	require.Equal(t, 1, len(logs))

	// 	logs = buf.dequeueRange(3, 3, 2, 10)
	// 	fmt.Println(logs)
	// 	require.Equal(t, 0, len(logs))
	// })
}
