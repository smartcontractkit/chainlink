package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventProvider_GetFilters(t *testing.T) {
	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	_, f := newEntry(p, 1)
	p.filterStore.AddActiveUpkeeps(f)

	t.Run("no filters", func(t *testing.T) {
		filters := p.getFilters(0, big.NewInt(0))
		require.Len(t, filters, 1)
		require.Empty(t, filters[0].addr)
	})

	t.Run("has filter with lower lastPollBlock", func(t *testing.T) {
		filters := p.getFilters(0, f.upkeepID)
		require.Len(t, filters, 1)
		require.NotEmpty(t, filters[0].addr)
		filters = p.getFilters(10, f.upkeepID)
		require.Len(t, filters, 1)
		require.NotEmpty(t, filters[0].addr)
	})

	t.Run("has filter with higher lastPollBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.lastPollBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Empty(t, filters[0].addr)
	})

	t.Run("has filter with higher configUpdateBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.configUpdateBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Empty(t, filters[0].addr)
	})
}

func TestLogEventProvider_UpdateEntriesLastPoll(t *testing.T) {
	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	n := 10

	// entries := map[string]upkeepFilter{}
	for i := 0; i < n; i++ {
		_, f := newEntry(p, i+1)
		p.filterStore.AddActiveUpkeeps(f)
	}

	t.Run("no entries", func(t *testing.T) {
		_, f := newEntry(p, n*2)
		f.lastPollBlock = 10
		p.updateFiltersLastPoll([]upkeepFilter{f})

		filters := p.filterStore.GetFilters(nil)
		for _, f := range filters {
			require.Equal(t, int64(0), f.lastPollBlock)
		}
	})

	t.Run("update entries", func(t *testing.T) {
		_, f2 := newEntry(p, n-2)
		f2.lastPollBlock = 10
		_, f1 := newEntry(p, n-1)
		f1.lastPollBlock = 10
		p.updateFiltersLastPoll([]upkeepFilter{f1, f2})

		p.filterStore.RangeFiltersByIDs(func(_ int, f upkeepFilter) {
			require.Equal(t, int64(10), f.lastPollBlock)
		}, f1.upkeepID, f2.upkeepID)

		// update with same block
		p.updateFiltersLastPoll([]upkeepFilter{f1})

		// checking other entries are not updated
		_, f := newEntry(p, 1)
		p.filterStore.RangeFiltersByIDs(func(_ int, f upkeepFilter) {
			require.Equal(t, int64(0), f.lastPollBlock)
		}, f.upkeepID)
	})
}

func TestLogEventProvider_ScheduleReadJobs(t *testing.T) {
	mp := new(mocks.LogPoller)

	tests := []struct {
		name         string
		maxBatchSize int
		ids          []int
		addrs        []string
	}{
		{
			"no entries",
			3,
			[]int{},
			[]string{},
		},
		{
			"single entry",
			3,
			[]int{1},
			[]string{"0x1111111"},
		},
		{
			"happy flow",
			3,
			[]int{1, 2, 3},
			[]string{"0x1111111", "0x2222222", "0x3333333"},
		},
		{
			"batching",
			3,
			[]int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, 15,
				16, 17, 18,
				19, 20, 21,
			},
			[]string{
				"0x11111111",
				"0x22222222",
				"0x33333333",
				"0x111111111",
				"0x122222222",
				"0x133333333",
				"0x1111111111",
				"0x1122222222",
				"0x1133333333",
				"0x11111111111",
				"0x11122222222",
				"0x11133333333",
				"0x111111111111",
				"0x111122222222",
				"0x111133333333",
				"0x1111111111111",
				"0x1111122222222",
				"0x1111133333333",
				"0x11111111111111",
				"0x11111122222222",
				"0x11111133333333",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			readInterval := 10 * time.Millisecond
			opts := NewOptions(200, big.NewInt(1))
			opts.ReadInterval = readInterval

			p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), opts)

			var ids []*big.Int
			for i, id := range tc.ids {
				_, f := newEntry(p, id, tc.addrs[i])
				p.filterStore.AddActiveUpkeeps(f)
				ids = append(ids, f.upkeepID)
			}

			reads := make(chan []*big.Int, 100)

			go func(ctx context.Context) {
				p.scheduleReadJobs(ctx, func(ids []*big.Int) {
					select {
					case reads <- ids:
					default:
						t.Log("dropped ids")
					}
				})
			}(ctx)

			batches := (len(tc.ids) / tc.maxBatchSize) + 1

			timeoutTicker := time.NewTicker(readInterval * time.Duration(batches*10))
			defer timeoutTicker.Stop()

			got := map[string]int{}

		readLoop:
			for {
				select {
				case <-timeoutTicker.C:
					break readLoop
				case batch := <-reads:
					for _, id := range batch {
						got[id.String()]++
					}
				case <-ctx.Done():
					break readLoop
				default:
					if p.CurrentPartitionIdx() > uint64(batches+1) {
						break readLoop
					}
				}
				runtime.Gosched()
			}

			require.Equal(t, len(ids), len(got))
			for _, id := range ids {
				_, ok := got[id.String()]
				require.True(t, ok, "id not found %s", id.String())
				require.GreaterOrEqual(t, got[id.String()], 1, "id don't have schdueled job %s", id.String())
			}
		})
	}
}

func TestLogEventProvider_ReadLogs(t *testing.T) {
	ctx := testutils.Context(t)

	mp := new(mocks.LogPoller)

	mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return()
	mp.On("HasFilter", mock.Anything).Return(false)
	mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("LatestBlock", mock.Anything).Return(logpoller.Block{BlockNumber: int64(1)}, nil)
	mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 1,
			TxHash:      common.HexToHash("0x1"),
		},
	}, nil)

	filterStore := NewUpkeepFilterStore()
	p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, filterStore, NewOptions(200, big.NewInt(1)))

	for i := 0; i < 10; i++ {
		cfg, f := newEntry(p, i+1)
		require.NoError(t, p.RegisterFilter(ctx, FilterOptions{
			UpkeepID:      f.upkeepID,
			TriggerConfig: cfg,
		}))
	}

	// TODO: test rate limiting
}

func newEntry(p *logEventProvider, i int, args ...string) (LogTriggerConfig, upkeepFilter) {
	idBytes := append(common.LeftPadBytes([]byte{1}, 16), []byte(strconv.Itoa(i))...)
	id := ocr2keepers.UpkeepIdentifier{}
	copy(id[:], idBytes)
	uid := id.BigInt()
	for len(args) < 2 {
		args = append(args, "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d")
	}
	addr, topic0 := args[0], args[1]
	cfg := LogTriggerConfig{
		ContractAddress: common.HexToAddress(addr),
		FilterSelector:  0,
		Topic0:          common.HexToHash(topic0),
	}
	filter := p.newLogFilter(uid, cfg)
	topics := make([]common.Hash, len(filter.EventSigs))
	copy(topics, filter.EventSigs)
	f := upkeepFilter{
		upkeepID: uid,
		addr:     filter.Addresses[0].Bytes(),
		topics:   topics,
	}
	return cfg, f
}

func countRemainingLogs(logs map[int64][]logpoller.Log) int {
	count := 0
	for _, logList := range logs {
		count += len(logList)
	}
	return count
}

func remainingBlockWindowCounts(queues map[string]*upkeepLogQueue, blockRate int) map[int64]int {
	blockWindowCounts := map[int64]int{}

	for _, q := range queues {
		for blockNumber, logs := range q.logs {
			start, _ := getBlockWindow(blockNumber, blockRate)

			blockWindowCounts[start] += len(logs)
		}
	}

	return blockWindowCounts
}

func TestLogEventProvider_GetLatestPayloads(t *testing.T) {
	t.Run("dequeuing from an empty buffer returns 0 logs", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Empty(t, payloads)
	})

	t.Run("a single log for a single upkeep gets dequeued", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer

		buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0})

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 1)
	})

	t.Run("a log per upkeep for 4 upkeeps across 4 blocks (2 separate block windows) is dequeued, for a total of 4 payloads", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer

		buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0})
		buffer.Enqueue(big.NewInt(2), logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0})
		buffer.Enqueue(big.NewInt(3), logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0})
		buffer.Enqueue(big.NewInt(4), logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x4"), LogIndex: 0})

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 4)
	})

	t.Run("100 logs are dequeued for a single upkeep, 1 log for every block window across 100 blocks followed by best effort", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 101, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer.(*logBuffer)

		for i := 0; i < 100; i++ {
			buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x%d", i+1)), LogIndex: 0})
		}

		assert.Equal(t, 100, countRemainingLogs(buffer.queues["1"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		assert.Equal(t, 0, countRemainingLogs(buffer.queues["1"].logs))
	})

	t.Run("100 logs are dequeued for two upkeeps, 25 logs each as min commitment (50 logs total best effort), followed by best effort", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 101, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer.(*logBuffer)

		for i := 0; i < 100; i++ {
			buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x1%d", i+1)), LogIndex: 0})
			buffer.Enqueue(big.NewInt(2), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x2%d", i+1)), LogIndex: 0})
		}

		assert.Equal(t, 100, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 100, countRemainingLogs(buffer.queues["2"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		assert.Equal(t, 50, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 50, countRemainingLogs(buffer.queues["2"].logs))

		windowCount := remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 2, windowCount[0])
		assert.Equal(t, 4, windowCount[48])
		assert.Equal(t, 4, windowCount[96])

		// the second dequeue call will retrieve the remaining 100 logs and exhaust the queues
		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		assert.Equal(t, 0, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 0, countRemainingLogs(buffer.queues["2"].logs))

		windowCount = remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 0, windowCount[0])
		assert.Equal(t, 0, windowCount[48])
		assert.Equal(t, 0, windowCount[96])
	})

	t.Run("minimum guaranteed for all windows including an incomplete window followed by best effort", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 102, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer.(*logBuffer)

		for i := 0; i < 102; i++ {
			buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x1%d", i+1)), LogIndex: 0})
			buffer.Enqueue(big.NewInt(2), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x2%d", i+1)), LogIndex: 0})
		}

		assert.Equal(t, 102, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 102, countRemainingLogs(buffer.queues["2"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		windowCount := remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 6, windowCount[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		// upkeep 1 has had the minimum number of logs dequeued on the latest (incomplete) window
		assert.Equal(t, 1, buffer.queues["1"].dequeued[100])
		// upkeep 2 has had the minimum number of logs dequeued on the latest (incomplete) window
		assert.Equal(t, 1, buffer.queues["2"].dequeued[100])

		// the third dequeue call will retrieve the remaining 100 logs and exhaust the queues
		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 4)

		assert.Equal(t, 0, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 0, countRemainingLogs(buffer.queues["2"].logs))

		windowCount = remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 0, windowCount[0])
		assert.Equal(t, 0, windowCount[28])
		assert.Equal(t, 0, windowCount[32])
		assert.Equal(t, 0, windowCount[36])
		assert.Equal(t, 0, windowCount[48])
		assert.Equal(t, 0, windowCount[96])
		assert.Equal(t, 0, windowCount[100])
	})

	t.Run("min dequeue followed by best effort followed by reorg followed by best effort", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 101, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer.(*logBuffer)

		for i := 0; i < 100; i++ {
			buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x1%d", i+1)), LogIndex: 0})
			buffer.Enqueue(big.NewInt(2), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x2%d", i+1)), LogIndex: 0})
		}

		assert.Equal(t, 100, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 100, countRemainingLogs(buffer.queues["2"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		windowCount := remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 4, windowCount[28])

		// reorg block 28
		buffer.Enqueue(big.NewInt(1), logpoller.Log{BlockNumber: int64(28), TxHash: common.HexToHash(fmt.Sprintf("0xreorg1%d", 28)), LogIndex: 0, BlockHash: common.BytesToHash([]byte("reorg"))})
		buffer.Enqueue(big.NewInt(2), logpoller.Log{BlockNumber: int64(28), TxHash: common.HexToHash(fmt.Sprintf("0xreorg2%d", 28)), LogIndex: 0, BlockHash: common.BytesToHash([]byte("reorg"))})

		windowCount = remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 6, windowCount[28])

		// the second dequeue call will retrieve the remaining 100 logs and exhaust the queues
		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		windowCount = remainingBlockWindowCounts(buffer.queues, 4)

		assert.Equal(t, 0, windowCount[0])
		assert.Equal(t, 0, windowCount[28])
		assert.Equal(t, 0, windowCount[32])
		assert.Equal(t, 0, windowCount[36])
		assert.Equal(t, 0, windowCount[48])
		assert.Equal(t, 2, windowCount[96]) // these 2 remaining logs are because of the 2 re orgd logs taking up dequeue space
	})

	t.Run("sparsely populated blocks", func(t *testing.T) {
		opts := NewOptions(200, big.NewInt(42161))

		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
		}

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, nil, opts)

		ctx := context.Background()

		buffer := provider.buffer.(*logBuffer)

		upkeepOmittedOnBlocks := map[int64][]int{
			1: {5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100},                                                                                                                      // upkeep 1 won't have logs on 20 blocks
			2: {2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100}, // upkeep 2 won't have logs on 50 blocks
			3: {3, 13, 23, 33, 43, 53, 63, 73, 83, 93},                                                                                                                                                               // upkeep 3 won't appear on 10 blocks
			4: {1, 25, 50, 75, 100},                                                                                                                                                                                  // upkeep 4 won't appear on 5 blocks
			5: {},                                                                                                                                                                                                    // upkeep 5 appears on all blocks
		}

		for upkeep, skipBlocks := range upkeepOmittedOnBlocks {
		blockLoop:
			for i := 0; i < 100; i++ {
				for _, block := range skipBlocks {
					if block == i+1 {
						continue blockLoop
					}
				}
				buffer.Enqueue(big.NewInt(upkeep), logpoller.Log{BlockNumber: int64(i + 1), TxHash: common.HexToHash(fmt.Sprintf("0x1%d", i+1)), LogIndex: 0})
			}
		}

		assert.Equal(t, 80, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 50, countRemainingLogs(buffer.queues["2"].logs))
		assert.Equal(t, 90, countRemainingLogs(buffer.queues["3"].logs))
		assert.Equal(t, 95, countRemainingLogs(buffer.queues["4"].logs))
		assert.Equal(t, 100, countRemainingLogs(buffer.queues["5"].logs))

		// perform two dequeues
		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)
		assert.Len(t, payloads, 100)

		assert.Equal(t, 40, countRemainingLogs(buffer.queues["1"].logs))
		assert.Equal(t, 10, countRemainingLogs(buffer.queues["2"].logs))
		assert.Equal(t, 50, countRemainingLogs(buffer.queues["3"].logs))
		assert.Equal(t, 55, countRemainingLogs(buffer.queues["4"].logs))
		assert.Equal(t, 60, countRemainingLogs(buffer.queues["5"].logs))
	})
}

type mockedPacker struct {
}

func (p *mockedPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	return log.Data, nil
}
