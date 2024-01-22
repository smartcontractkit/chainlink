package logprovider

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestLogEventBuffer_GetBlocksInRange(t *testing.T) {
	size := 3
	maxSeenBlock := int64(4)
	buf := newLogEventBuffer(logger.TestLogger(t), size, 10, 10)

	buf.enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	)

	buf.enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 2},
		logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 2},
		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	)

	tests := []struct {
		name string
		from int64
		to   int64
		want int
	}{
		{
			name: "all",
			from: 2,
			to:   4,
			want: 3,
		},
		{
			name: "partial",
			from: 2,
			to:   3,
			want: 2,
		},
		{
			name: "circular",
			from: 3,
			to:   4,
			want: 2,
		},
		{
			name: "zero start",
			from: 0,
			to:   2,
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
		{
			name: "outside max last seen",
			from: 5,
			to:   10,
		},
		{
			name: "limited by max last seen",
			from: 2,
			to:   5,
			want: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blocks := buf.getBlocksInRange(int(tc.from), int(tc.to))
			require.Equal(t, tc.want, len(blocks))
			if tc.want > 0 {
				from := tc.from
				require.Equal(t, from, blocks[0].blockNumber)
				to := tc.to
				if to >= maxSeenBlock {
					to = maxSeenBlock
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
		buf := newLogEventBuffer(logger.TestLogger(t), 2, 2, 2)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		))
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue logs overflow with dynamic limits", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 2, 10, 2)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		))
		buf.SetLimits(10, 3)
		require.Equal(t, 3, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 1},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 2},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 3},
		))

		buf.lock.Lock()
		defer buf.lock.Unlock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		require.Equal(t, 3, len(buf.blocks[1].logs))
	})

	t.Run("enqueue logs overflow with dynamic limits", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 2, 10, 2)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 3},
		))
		buf.SetLimits(10, 3)
		require.Equal(t, 3, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 1},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 2},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x21"), LogIndex: 3},
		))

		buf.lock.Lock()
		defer buf.lock.Unlock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
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
		removed := buf.dequeueRange(int64(1), int64(2), 2, 10)
		require.Equal(t, 2, len(removed))
		results = buf.peekRange(int64(1), int64(2))
		require.Equal(t, 0, len(results))
	})

	t.Run("enqueue peek and dequeue", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 4, 10, 10)

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
		removed := buf.dequeueRange(1, 3, 5, 5)
		require.Equal(t, 4, len(removed))
		buf.lock.Lock()
		require.Equal(t, 0, len(buf.blocks[0].logs))
		require.Equal(t, int64(2), buf.blocks[1].blockNumber)
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
		results := buf.peekRange(int64(1), int64(5))
		fmt.Println(results)
		verifyBlockNumbers(t, results, 2, 3, 4, 4)
	})

	t.Run("dequeue with limits returns latest block logs", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)
		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 0},
			logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x5"), LogIndex: 0},
		), 5)

		logs := buf.dequeueRange(1, 5, 2, 10)
		require.Equal(t, 2, len(logs))
		require.Equal(t, int64(5), logs[0].log.BlockNumber)
		require.Equal(t, int64(4), logs[1].log.BlockNumber)

		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 1},
			logpoller.Log{BlockNumber: 5, TxHash: common.HexToHash("0x5"), LogIndex: 1},
		), 2)

		logs = buf.dequeueRange(1, 5, 3, 2)
		require.Equal(t, 2, len(logs))
	})

	t.Run("dequeue doesn't return same logs again", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 5, 10)
		require.Equal(t, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
		), 3)

		logs := buf.dequeueRange(3, 3, 2, 10)
		fmt.Println(logs)
		require.Equal(t, 1, len(logs))

		logs = buf.dequeueRange(3, 3, 2, 10)
		fmt.Println(logs)
		require.Equal(t, 0, len(logs))
	})
}

func TestLogEventBuffer_FetchedBlock_Append(t *testing.T) {
	type appendArgs struct {
		fl                          fetchedLog
		maxBlockLogs, maxUpkeepLogs int
		added, dropped              bool
	}

	tests := []struct {
		name        string
		blockNumber int64
		logs        []fetchedLog
		visited     []fetchedLog
		toAdd       []appendArgs
		expected    []fetchedLog
		added       bool
	}{
		{
			name:        "empty block",
			blockNumber: 1,
			logs:        []fetchedLog{},
			visited:     []fetchedLog{},
			toAdd: []appendArgs{
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    0,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         true,
				},
			},
			expected: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
		},
		{
			name:        "existing log",
			blockNumber: 1,
			logs: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
			visited: []fetchedLog{},
			toAdd: []appendArgs{
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    0,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         false,
				},
			},
			expected: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
		},
		{
			name:        "visited log",
			blockNumber: 1,
			logs:        []fetchedLog{},
			visited: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
			toAdd: []appendArgs{
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    0,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         false,
				},
			},
			expected: []fetchedLog{},
		},
		{
			name:        "upkeep log limits",
			blockNumber: 1,
			logs:        []fetchedLog{},
			visited:     []fetchedLog{},
			toAdd: []appendArgs{
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    0,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         true,
				},
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    1,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         true,
				},
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    2,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  10,
					maxUpkeepLogs: 2,
					added:         true,
					dropped:       true,
				},
			},
			expected: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    1,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    2,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
		},
		{
			name:        "block log limits",
			blockNumber: 1,
			logs:        []fetchedLog{},
			visited:     []fetchedLog{},
			toAdd: []appendArgs{
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    0,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  2,
					maxUpkeepLogs: 4,
					added:         true,
				},
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    1,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  2,
					maxUpkeepLogs: 4,
					added:         true,
				},
				{
					fl: fetchedLog{
						log: logpoller.Log{
							BlockNumber: 1,
							TxHash:      common.HexToHash("0x1"),
							LogIndex:    2,
						},
						upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
					},
					maxBlockLogs:  2,
					maxUpkeepLogs: 4,
					added:         true,
					dropped:       true,
				},
			},
			expected: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    1,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    2,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			b := fetchedBlock{
				blockNumber: tc.blockNumber,
				logs:        make([]fetchedLog, len(tc.logs)),
				visited:     make([]fetchedLog, len(tc.visited)),
			}
			copy(b.logs, tc.logs)
			copy(b.visited, tc.visited)

			for _, args := range tc.toAdd {
				dropped, added := b.Append(lggr, args.fl, args.maxBlockLogs, args.maxUpkeepLogs)
				require.Equal(t, args.added, added)
				if args.dropped {
					require.NotNil(t, dropped.upkeepID)
				} else {
					require.Nil(t, dropped.upkeepID)
				}
			}
			// clear cached logIDs
			for i := range b.logs {
				b.logs[i].cachedLogID = ""
			}
			require.Equal(t, tc.expected, b.logs)
		})
	}
}
func TestLogEventBuffer_FetchedBlock_Sort(t *testing.T) {
	tests := []struct {
		name        string
		blockNumber int64
		logs        []fetchedLog
		beforeSort  []string
		afterSort   []string
		iterations  int
	}{
		{
			name:        "no logs",
			blockNumber: 10,
			logs:        []fetchedLog{},
			beforeSort:  []string{},
			afterSort:   []string{},
		},
		{
			name:        "single log",
			blockNumber: 1,
			logs: []fetchedLog{
				{
					log: logpoller.Log{
						BlockHash:   common.HexToHash("0x111"),
						BlockNumber: 1,
						TxHash:      common.HexToHash("0x1"),
						LogIndex:    0,
					},
				},
			},
			beforeSort: []string{
				"0000000000000000000000000000000000000000000000000000000000000111000000000000000000000000000000000000000000000000000000000000000100000000",
			},
			afterSort: []string{
				"0000000000000000000000000000000000000000000000000000000000000111000000000000000000000000000000000000000000000000000000000000000100000000",
			},
		},
		{
			name:        "multiple logs with 10 iterations",
			blockNumber: 1,
			logs: []fetchedLog{
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xb711bd1103927611ee41152aa8ae27f3330"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    0,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "222").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    4,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    3,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "222").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    2,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    5,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    3,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
				{
					log: logpoller.Log{
						BlockNumber: 1,
						BlockHash:   common.HexToHash("0xa25ebae1099f3fbae2525ebae279f3ae25e"),
						TxHash:      common.HexToHash("0xa651bd1109922111ee411525ebae27f3fb6"),
						LogIndex:    1,
					},
					upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
				},
			},
			beforeSort: []string{
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000b711bd1103927611ee41152aa8ae27f333000000000",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000000",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000004",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000003",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000002",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000005",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000003",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000001",
			},
			afterSort: []string{
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000b711bd1103927611ee41152aa8ae27f333000000000",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000000",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000001",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000002",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000003",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000003",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000004",
				"00000000000000000000000000000a25ebae1099f3fbae2525ebae279f3ae25e00000000000000000000000000000a651bd1109922111ee411525ebae27f3fb600000005",
			},
			iterations: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := fetchedBlock{
				blockNumber: tc.blockNumber,
				logs:        make([]fetchedLog, len(tc.logs)),
			}
			if tc.iterations == 0 {
				tc.iterations = 1
			}
			// performing the same multiple times should yield the same result
			// default is one iteration
			for i := 0; i < tc.iterations; i++ {
				copy(b.logs, tc.logs)
				logIDs := getLogIds(b)
				require.Equal(t, len(tc.beforeSort), len(logIDs))
				require.Equal(t, tc.beforeSort, logIDs)
				b.Sort()
				logIDsAfterSort := getLogIds(b)
				require.Equal(t, len(tc.afterSort), len(logIDsAfterSort))
				require.Equal(t, tc.afterSort, logIDsAfterSort)
			}
		})
	}
}

func TestLogEventBuffer_FetchedBlock_Clone(t *testing.T) {
	b1 := fetchedBlock{
		blockNumber: 1,
		logs: []fetchedLog{
			{
				log: logpoller.Log{
					BlockNumber: 1,
					TxHash:      common.HexToHash("0x1"),
					LogIndex:    0,
				},
				upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
			},
			{
				log: logpoller.Log{
					BlockNumber: 1,
					TxHash:      common.HexToHash("0x1"),
					LogIndex:    2,
				},
				upkeepID: core.GenUpkeepID(types.LogTrigger, "111").BigInt(),
			},
		},
	}

	b2 := b1.Clone()
	require.Equal(t, b1.blockNumber, b2.blockNumber)
	require.Equal(t, len(b1.logs), len(b2.logs))
	require.Equal(t, b1.logs[0].log.BlockNumber, b2.logs[0].log.BlockNumber)

	b1.blockNumber = 2
	b1.logs[0].log.BlockNumber = 2
	require.NotEqual(t, b1.blockNumber, b2.blockNumber)
	require.NotEqual(t, b1.logs[0].log.BlockNumber, b2.logs[0].log.BlockNumber)
}

func verifyBlockNumbers(t *testing.T, logs []fetchedLog, bns ...int64) {
	require.Equal(t, len(bns), len(logs), "expected length mismatch")
	for i, log := range logs {
		require.Equal(t, bns[i], log.log.BlockNumber, "wrong block number")
	}
}

func getLogIds(b fetchedBlock) []string {
	logIDs := make([]string, len(b.logs))
	for i, l := range b.logs {
		ext := ocr2keepers.LogTriggerExtension{
			TxHash:    l.log.TxHash,
			Index:     uint32(l.log.LogIndex),
			BlockHash: l.log.BlockHash,
		}
		logIDs[i] = hex.EncodeToString(ext.LogIdentifier())
	}
	return logIDs
}
