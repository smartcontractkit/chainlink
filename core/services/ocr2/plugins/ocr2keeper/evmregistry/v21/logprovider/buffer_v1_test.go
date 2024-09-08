package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventBufferV1(t *testing.T) {
	buf := NewLogBuffer(logger.TestLogger(t), 10, 20, 1)

	buf.Enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	)
	buf.Enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	)
	results, remaining := buf.Dequeue(int64(1), 2, true)
	require.Equal(t, 2, len(results))
	require.Equal(t, 2, remaining)
	require.True(t, results[0].ID.Cmp(results[1].ID) != 0)
	results, remaining = buf.Dequeue(int64(1), 2, true)
	require.Equal(t, 0, len(results))
	require.Equal(t, 0, remaining)
}

func TestLogEventBufferV1_SyncFilters(t *testing.T) {
	buf := NewLogBuffer(logger.TestLogger(t), 10, 20, 1)

	buf.Enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	)
	buf.Enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	)
	filterStore := NewUpkeepFilterStore()
	filterStore.AddActiveUpkeeps(upkeepFilter{upkeepID: big.NewInt(1)})

	require.Equal(t, 2, buf.NumOfUpkeeps())
	require.NoError(t, buf.SyncFilters(filterStore))
	require.Equal(t, 1, buf.NumOfUpkeeps())
}

type readableLogger struct {
	logger.Logger
	DebugwFn func(msg string, keysAndValues ...interface{})
}

func (l *readableLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.DebugwFn(msg, keysAndValues...)
}

func (l *readableLogger) Named(name string) logger.Logger {
	return l
}

func (l *readableLogger) With(args ...interface{}) logger.Logger {
	return l
}

func TestLogEventBufferV1_EnqueueViolations(t *testing.T) {
	t.Run("enqueuing logs for a block older than latest seen logs a message", func(t *testing.T) {
		logReceived := false
		readableLogger := &readableLogger{
			Logger: logger.TestLogger(t),
			DebugwFn: func(msg string, keysAndValues ...interface{}) {
				if msg == "enqueuing logs from a block older than latest seen block" {
					logReceived = true
					assert.Equal(t, "logBlock", keysAndValues[0])
					assert.Equal(t, int64(1), keysAndValues[1])
					assert.Equal(t, "lastBlockSeen", keysAndValues[2])
					assert.Equal(t, int64(2), keysAndValues[3])
				}
			},
		}

		logBufferV1 := NewLogBuffer(readableLogger, 10, 20, 1)

		buf := logBufferV1.(*logBuffer)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(2),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		)

		assert.True(t, true, logReceived)
	})

	t.Run("enqueuing logs for the same block over multiple calls logs a message", func(t *testing.T) {
		logReceived := false
		readableLogger := &readableLogger{
			Logger: logger.TestLogger(t),
			DebugwFn: func(msg string, keysAndValues ...interface{}) {
				if msg == "enqueuing logs again for a previously seen block" {
					logReceived = true
					assert.Equal(t, "blockNumber", keysAndValues[0])
					assert.Equal(t, int64(3), keysAndValues[1])
					assert.Equal(t, "numberOfEnqueues", keysAndValues[2])
					assert.Equal(t, 2, keysAndValues[3])
				}
			},
		}

		logBufferV1 := NewLogBuffer(readableLogger, 10, 20, 1)

		buf := logBufferV1.(*logBuffer)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(2),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		)
		buf.Enqueue(big.NewInt(3),
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3a"), LogIndex: 0},
		)
		buf.Enqueue(big.NewInt(3),
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3b"), LogIndex: 0},
		)

		assert.True(t, true, logReceived)
	})
}

func TestLogEventBufferV1_Dequeue(t *testing.T) {
	tests := []struct {
		name         string
		logsInBuffer map[*big.Int][]logpoller.Log
		args         dequeueArgs
		lookback     int
		results      []logpoller.Log
		remaining    int
	}{
		{
			name:         "empty",
			logsInBuffer: map[*big.Int][]logpoller.Log{},
			args:         newDequeueArgs(10, 1, 1, 10),
			lookback:     20,
			results:      []logpoller.Log{},
		},
		{
			name: "happy path",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 14, TxHash: common.HexToHash("0x15"), LogIndex: 1},
				},
			},
			args:     newDequeueArgs(10, 5, 3, 10),
			lookback: 20,
			results: []logpoller.Log{
				{}, {},
			},
		},
		{
			name: "with upkeep limits",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 1},
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 0},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 1},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 1},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 2},
				},
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 10},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 10},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 11},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 11},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 12},
				},
			},
			args:     newDequeueArgs(10, 5, 2, 10),
			lookback: 20,
			results: []logpoller.Log{
				{}, {}, {}, {},
			},
			remaining: 8,
		},
		{
			name: "with max results",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): append(createDummyLogSequence(2, 0, 12, common.HexToHash("0x12")), createDummyLogSequence(2, 0, 13, common.HexToHash("0x13"))...),
				big.NewInt(2): append(createDummyLogSequence(2, 10, 12, common.HexToHash("0x12")), createDummyLogSequence(2, 10, 13, common.HexToHash("0x13"))...),
			},
			args:     newDequeueArgs(10, 5, 3, 4),
			lookback: 20,
			results: []logpoller.Log{
				{}, {}, {}, {},
			},
			remaining: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewLogBuffer(logger.TestLogger(t), uint32(tc.lookback), uint32(tc.args.blockRate), uint32(tc.args.upkeepLimit))
			for id, logs := range tc.logsInBuffer {
				added, dropped := buf.Enqueue(id, logs...)
				require.Equal(t, len(logs), added+dropped)
			}
			start, _ := getBlockWindow(tc.args.block, tc.args.blockRate)

			results, remaining := buf.Dequeue(start, tc.args.maxResults, true)
			require.Equal(t, len(tc.results), len(results))
			require.Equal(t, tc.remaining, remaining)
		})
	}
}

func TestLogEventBufferV1_Enqueue(t *testing.T) {
	tests := []struct {
		name                             string
		logsToAdd                        map[*big.Int][]logpoller.Log
		added, dropped                   map[string]int
		sizeOfRange                      map[*big.Int]int
		rangeStart, rangeEnd             int64
		lookback, blockRate, upkeepLimit uint32
	}{
		{
			name:        "empty",
			logsToAdd:   map[*big.Int][]logpoller.Log{},
			added:       map[string]int{},
			dropped:     map[string]int{},
			sizeOfRange: map[*big.Int]int{},
			rangeStart:  0,
			rangeEnd:    10,
			blockRate:   1,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "happy path",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 14, TxHash: common.HexToHash("0x15"), LogIndex: 1},
				},
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
				},
			},
			added: map[string]int{
				big.NewInt(1).String(): 2,
				big.NewInt(2).String(): 1,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 0,
				big.NewInt(2).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 2,
				big.NewInt(2): 1,
			},
			rangeStart:  10,
			rangeEnd:    20,
			blockRate:   5,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "above limits",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): createDummyLogSequence(11, 0, 12, common.HexToHash("0x12")),
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
				},
			},
			added: map[string]int{
				big.NewInt(1).String(): 11,
				big.NewInt(2).String(): 1,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 1,
				big.NewInt(2).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 10,
				big.NewInt(2): 1,
			},
			rangeStart:  10,
			rangeEnd:    20,
			blockRate:   10,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "out of block range",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): append(createDummyLogSequence(2, 0, 1, common.HexToHash("0x1")), createDummyLogSequence(2, 0, 100, common.HexToHash("0x1"))...),
			},
			added: map[string]int{
				big.NewInt(1).String(): 2,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 2,
			},
			rangeStart:  1,
			rangeEnd:    101,
			blockRate:   10,
			upkeepLimit: 10,
			lookback:    20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewLogBuffer(logger.TestLogger(t), tc.lookback, tc.blockRate, tc.upkeepLimit)
			for id, logs := range tc.logsToAdd {
				added, dropped := buf.Enqueue(id, logs...)
				sid := id.String()
				if _, ok := tc.added[sid]; !ok {
					tc.added[sid] = 0
				}
				if _, ok := tc.dropped[sid]; !ok {
					tc.dropped[sid] = 0
				}
				require.Equal(t, tc.added[sid], added)
				require.Equal(t, tc.dropped[sid], dropped)
			}
			for id, size := range tc.sizeOfRange {
				q, ok := buf.(*logBuffer).getUpkeepQueue(id)
				require.True(t, ok)
				require.Equal(t, size, q.sizeOfRange(tc.rangeStart, tc.rangeEnd))
			}
		})
	}
}

func TestLogEventBufferV1_UpkeepQueue(t *testing.T) {
	t.Run("enqueue dequeue", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		added, dropped := q.enqueue(10, logpoller.Log{BlockNumber: 20, TxHash: common.HexToHash("0x1"), LogIndex: 0})
		require.Equal(t, 0, dropped)
		require.Equal(t, 1, added)
		require.Equal(t, 1, q.sizeOfRange(1, 20))
		logs, remaining := q.dequeue(19, 21, 10)
		require.Equal(t, 1, len(logs))
		require.Equal(t, 0, remaining)
	})

	t.Run("enqueue with limits", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		added, dropped := q.enqueue(10,
			createDummyLogSequence(15, 0, 20, common.HexToHash("0x20"))...,
		)
		require.Equal(t, 5, dropped)
		require.Equal(t, 15, added)
	})

	t.Run("dequeue with limits", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 3))

		added, dropped := q.enqueue(10,
			logpoller.Log{BlockNumber: 20, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 20, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 20, TxHash: common.HexToHash("0x1"), LogIndex: 10},
		)
		require.Equal(t, 0, dropped)
		require.Equal(t, 3, added)

		logs, remaining := q.dequeue(19, 21, 2)
		require.Equal(t, 2, len(logs))
		require.Equal(t, 1, remaining)
	})
}

func TestLogEventBufferV1_UpkeepQueue_sizeOfRange(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		require.Equal(t, 0, q.sizeOfRange(1, 10))
	})

	t.Run("happy path", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		added, dropped := q.enqueue(10, logpoller.Log{BlockNumber: 20, TxHash: common.HexToHash("0x1"), LogIndex: 0})
		require.Equal(t, 0, dropped)
		require.Equal(t, 1, added)
		require.Equal(t, 0, q.sizeOfRange(1, 10))
		require.Equal(t, 1, q.sizeOfRange(1, 20))
	})
}

func TestLogEventBufferV1_UpkeepQueue_clean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := newUpkeepLogQueue(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

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
		q.lock.Lock()
		defer q.lock.Unlock()
		require.Equal(t, 2, len(q.states))
	})
}

func TestLogEventBufferV1_BlockWindow(t *testing.T) {
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
			start, end := getBlockWindow(tc.block, tc.blockRate)
			require.Equal(t, tc.wantStart, start)
			require.Equal(t, tc.wantEnd, end)
		})
	}
}

type dequeueArgs struct {
	block       int64
	blockRate   int
	upkeepLimit int
	maxResults  int
}

func newDequeueArgs(block int64, blockRate int, upkeepLimit int, maxResults int) dequeueArgs {
	args := dequeueArgs{
		block:       block,
		blockRate:   blockRate,
		upkeepLimit: upkeepLimit,
		maxResults:  maxResults,
	}

	if blockRate == 0 {
		args.blockRate = 1
	}
	if maxResults == 0 {
		args.maxResults = 10
	}
	if upkeepLimit == 0 {
		args.upkeepLimit = 1
	}

	return args
}

func createDummyLogSequence(n, startIndex int, block int64, tx common.Hash) []logpoller.Log {
	logs := make([]logpoller.Log, n)
	for i := 0; i < n; i++ {
		logs[i] = logpoller.Log{
			BlockNumber: block,
			TxHash:      tx,
			LogIndex:    int64(i + startIndex),
		}
	}
	return logs
}
