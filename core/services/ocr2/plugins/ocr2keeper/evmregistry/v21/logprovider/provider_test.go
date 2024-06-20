package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
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
		require.Equal(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with lower lastPollBlock", func(t *testing.T) {
		filters := p.getFilters(0, f.upkeepID)
		require.Len(t, filters, 1)
		require.Greater(t, len(filters[0].addr), 0)
		filters = p.getFilters(10, f.upkeepID)
		require.Len(t, filters, 1)
		require.Greater(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with higher lastPollBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.lastPollBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Equal(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with higher configUpdateBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.configUpdateBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Equal(t, len(filters[0].addr), 0)
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
	mp.On("LatestBlock", mock.Anything).Return(logpoller.LogPollerBlock{BlockNumber: int64(1)}, nil)
	mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 1,
			TxHash:      common.HexToHash("0x1"),
		},
	}, nil)

	filterStore := NewUpkeepFilterStore()
	p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, filterStore, NewOptions(200, big.NewInt(1)))

	var ids []*big.Int
	for i := 0; i < 10; i++ {
		cfg, f := newEntry(p, i+1)
		ids = append(ids, f.upkeepID)
		require.NoError(t, p.RegisterFilter(ctx, FilterOptions{
			UpkeepID:      f.upkeepID,
			TriggerConfig: cfg,
		}))
	}

	t.Run("no entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, big.NewInt(999999)))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 0)
	})

	t.Run("has entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, ids[:2]...))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 2)

		var updatedFilters []upkeepFilter
		filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
			updatedFilters = append(updatedFilters, f.Clone())
		}, ids[:2]...)
		for _, f := range updatedFilters {
			// Last poll block should be updated
			require.Equal(t, int64(1), f.lastPollBlock)
		}
	})

	// TODO: test rate limiting
}

func newEntry(p *logEventProvider, i int, args ...string) (LogTriggerConfig, upkeepFilter) {
	idBytes := append(common.LeftPadBytes([]byte{1}, 16), []byte(fmt.Sprintf("%d", i))...)
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

func countLogs(logs map[int64][]logpoller.Log) int {
	count := 0
	for _, logList := range logs {
		count += len(logList)
	}
	return count
}

func TestLogEventProvider_GetLatestPayloads(t *testing.T) {
	t.Run("5 upkeeps, 100 logs per upkeep per block for 100 blocks", func(t *testing.T) {
		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9980, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["5"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9960, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["5"].logs))
	})

	t.Run("200 upkeeps", func(t *testing.T) {
		var upkeepIDs []*big.Int

		for i := int64(1); i <= 200; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 200, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["150"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 9999, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 9999, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 9998, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["150"].logs))
	})

	t.Run("200 upkeeps, increasing to 300 upkeeps midway through the test", func(t *testing.T) {
		var upkeepIDs []*big.Int

		for i := int64(1); i <= 200; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 200, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["9"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["21"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["150"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps; with 2 iterations this means even upkeep IDs are dequeued first
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["40"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["45"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps; on the second iteration, odd upkeep IDs are dequeued
		assert.Equal(t, 9999, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["99"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["100"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps; on the third iteration, even upkeep IDs are dequeued once again
		assert.Equal(t, 9999, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["150"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["160"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["170"].logs))

		for i := int64(201); i <= 300; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		for i := 200; i < len(upkeepIDs); i++ {
			upkeepID := upkeepIDs[i]
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		err = provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 300, provider.bufferV1.NumOfUpkeeps())

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps; the new iterations
		// have not yet been recalculated despite the new logs being added; new iterations
		// are only calculated when current iteration maxes out at the total number of iterations
		assert.Equal(t, 9998, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["51"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["52"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// with the newly added logs, iterations is recalculated
		assert.Equal(t, 3, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 9998, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["11"].logs))
		assert.Equal(t, 9997, countLogs(bufV1.queues["111"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9997, countLogs(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, provider.iterations)
		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across selected upkeeps
		assert.Equal(t, 9997, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 9997, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, countLogs(bufV1.queues["101"].logs))
		assert.Equal(t, 9997, countLogs(bufV1.queues["150"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["250"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["299"].logs))
		assert.Equal(t, 9999, countLogs(bufV1.queues["300"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, provider.iterations)
		assert.Equal(t, 3, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		var remainingLogs int
		// at this point, every queue should have had at least one log dequeued
		for _, queue := range bufV1.queues {
			assert.True(t, countLogs(queue.logs) < 10000)
			remainingLogs += countLogs(queue.logs)
		}

		// check that across all 300 upkeeps, we have only dequeued 700 of the 3000000 logs (7 dequeue calls of 100 logs)
		assert.Equal(t, 2999300, remainingLogs)
	})

	t.Run("minimum guaranteed for all windows followed by best effort", func(t *testing.T) {
		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 10000, countLogs(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9980, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 9980, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the first block window
		assert.Equal(t, 400, blockWindowCounts[1])
		assert.Equal(t, 500, blockWindowCounts[2])
		assert.Equal(t, 500, blockWindowCounts[3])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9960, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 9960, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 400, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 500, blockWindowCounts[3])

		for i := 0; i < 97; i++ {
			payloads, err = provider.GetLatestPayloads(ctx)
			assert.NoError(t, err)

			// we dequeue a maximum of 10 logs
			assert.Equal(t, 100, len(payloads))
		}

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 8020, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 8020, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 8020, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 8020, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 8020, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 400, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[99])
		assert.Equal(t, 500, blockWindowCounts[100])

		// at this point, all block windows except for the latest block window will have been dequeued
		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 8000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 8000, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 8000, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 8000, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 8000, countLogs(bufV1.queues["5"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 7980, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 7980, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 7980, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 7980, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 7980, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 300, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 7960, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 7960, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 7960, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 7960, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 7960, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 200, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 7940, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 7940, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 7940, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 7940, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 7940, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 100, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 7920, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 7920, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 7920, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 7920, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 7920, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 400, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 10 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 7900, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 7900, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 7900, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 7900, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 7900, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}
		// all 10 logs should have been dequeued from the second block window, since the first block window has met it's minimum commitment
		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 300, blockWindowCounts[2])
		assert.Equal(t, 400, blockWindowCounts[3])
		assert.Equal(t, 400, blockWindowCounts[100])
	})

	t.Run("a complete window with no logs present is immediately marked as having the min logs dequeued, logs are dequeued from the next window", func(t *testing.T) {
		oldMaxPayloads := MaxPayloads
		MaxPayloads = 10
		defer func() {
			MaxPayloads = oldMaxPayloads
		}()

		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start + 4; i <= end; i++ {
				logsToAdd := 10
				if i >= 100 {
					logsToAdd = 1
				}
				for j := 0; j < logsToAdd; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i)),
						BlockNumber: i,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 99, nil // make the latest window incomplete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"
		opts.BlockRate = 4 // block window will be 4 blocks big

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		// all 10 logs should have been dequeued from the first block window
		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 200, blockWindowCounts[4])
		assert.Equal(t, 200, blockWindowCounts[8])
		assert.Equal(t, 200, blockWindowCounts[96])

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		dequeueCoordinator := provider.dequeueCoordinator.(*dequeueCoordinator)

		// the first block window does not contain any logs, so it automatically gets marked as having the minimum dequeued
		assert.True(t, true, dequeueCoordinator.dequeuedMinimum[0])

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		// all 10 logs should have been dequeued from the second block window
		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 190, blockWindowCounts[4])
		assert.Equal(t, 200, blockWindowCounts[8])
		assert.Equal(t, 200, blockWindowCounts[96])
	})

	t.Run("an incomplete window with no logs present it not dequeued from", func(t *testing.T) {
		oldMaxPayloads := MaxPayloads
		MaxPayloads = 10
		defer func() {
			MaxPayloads = oldMaxPayloads
		}()

		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start + 4; i <= end; i++ { // ensure the first block window never receives any logs
				logsToAdd := 10
				if i >= 100 {
					logsToAdd = 1
				}
				for j := 0; j < logsToAdd; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i)),
						BlockNumber: i,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 2, nil // make the latest window incomplete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(42161))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 0, len(blockWindowCounts)) // we don't have any log counts at all at this point

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(payloads))

		dequeueCoordinator := provider.dequeueCoordinator.(*dequeueCoordinator)

		assert.Equal(t, false, dequeueCoordinator.dequeuedMinimum[0])

		provider.poller = &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 3, nil // make the latest window complete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(payloads))

		provider.poller = &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 7, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		err = provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 50, blockWindowCounts[4]) // we now have 5 upkeeps x 10 logs each = 50 logs in the second window

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 40, blockWindowCounts[4]) // the dequeue process takes 5 logs as min commitment and 5 as best effort

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 30, blockWindowCounts[4]) // 10 logs taken as best effort
	})

	t.Run("an incomplete window with minimum logs already present is dequeued", func(t *testing.T) {
		oldMaxPayloads := MaxPayloads
		MaxPayloads = 10
		defer func() {
			MaxPayloads = oldMaxPayloads
		}()

		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i <= end; i++ {
				logsToAdd := 10
				for j := 0; j < logsToAdd; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i)),
						BlockNumber: i,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 2, nil // make the latest window incomplete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(42161))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 50, blockWindowCounts[0]) // a max of 10 logs per upkeep per window, 5 upkeeps = 5 x 10 = 50 logs

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		dequeueCoordinator := provider.dequeueCoordinator.(*dequeueCoordinator)

		assert.Equal(t, true, dequeueCoordinator.dequeuedMinimum[0])

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		// all 10 logs should have been dequeued from the first block window
		assert.Equal(t, 40, blockWindowCounts[0])

		logGenerator = func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start + 4; i <= end; i++ {
				logsToAdd := 10
				for j := 0; j < logsToAdd; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i)),
						BlockNumber: i,
					})
				}
			}
			return res
		}

		provider.poller = &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 7, nil // make the latest window incomplete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		err = provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 35, blockWindowCounts[0]) // 5 logs dequeued from here as best effort AFTER window[4]
		assert.Equal(t, 45, blockWindowCounts[4]) // 5 logs dequeued from here as min guarantee

		assert.Equal(t, true, dequeueCoordinator.dequeuedMinimum[0])
		assert.Equal(t, true, dequeueCoordinator.dequeuedMinimum[4])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 25, blockWindowCounts[0]) // best effort dequeue dequeues 10 logs from here
		assert.Equal(t, 45, blockWindowCounts[4])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 15, blockWindowCounts[0]) // best effort dequeue dequeues 10 logs from here
		assert.Equal(t, 45, blockWindowCounts[4])
	})

	t.Run("min dequeue followed by best effort followed by reorg followed by best effort", func(t *testing.T) {
		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 10; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 50, blockWindowCounts[1])

		// each upkeep should have 10 logs * 100 blocks = 1000 logs
		assert.Equal(t, 1000, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 1000, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 1000, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 1000, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 1000, countLogs(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 980, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 980, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 980, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 980, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 980, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[2])
		assert.Equal(t, 50, blockWindowCounts[3])
		assert.Equal(t, 50, blockWindowCounts[100])

		for i := 0; i < 49; i++ {
			_, err = provider.GetLatestPayloads(ctx)
			assert.NoError(t, err)
		}

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 0, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[18])
		assert.Equal(t, 0, blockWindowCounts[19])
		assert.Equal(t, 0, blockWindowCounts[20])

		// another 49 dequeue calls later, and the first 100 block windows have been completely dequeued, in pursuit of min
		// commitment dequeue - in this case, for chain ID 1, min commitment is 20 logs per upkeep per window, but each
		// window only has 10 logs per upkeep, so a window will be fully exhausted before minimum commitment is met

		assert.Equal(t, 0, blockWindowCounts[97])
		assert.Equal(t, 0, blockWindowCounts[98])
		assert.Equal(t, 0, blockWindowCounts[99])
		assert.Equal(t, 0, blockWindowCounts[100])

		// reorg happens
		logGenerator = func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				if i == 97 {
					for j := 0; j < 10; j++ {
						res = append(res, logpoller.Log{
							LogIndex:    int64(j),
							BlockHash:   common.HexToHash(fmt.Sprintf("%de", i+1)),
							BlockNumber: i + 1,
						})
					}
				} else {
					for j := 0; j < 10; j++ {
						res = append(res, logpoller.Log{
							LogIndex:    int64(j),
							BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
							BlockNumber: i + 1,
						})
					}
				}
			}
			return res
		}
		// use a log poller that will create logs for the queried block range
		provider.poller = &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 102, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[97])
		assert.Equal(t, 0, blockWindowCounts[98])
		assert.Equal(t, 0, blockWindowCounts[99])
		assert.Equal(t, 0, blockWindowCounts[100])

		err = provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[97])
		assert.Equal(t, 50, blockWindowCounts[98]) // reorg block window has had new logs added after reorg
		assert.Equal(t, 0, blockWindowCounts[99])
		assert.Equal(t, 0, blockWindowCounts[100])
		assert.Equal(t, 50, blockWindowCounts[101])
		assert.Equal(t, 50, blockWindowCounts[102])

		assert.Equal(t, 30, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 30, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 30, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 30, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 30, countLogs(bufV1.queues["5"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 10, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 10, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 10, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 10, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 10, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[97])
		assert.Equal(t, 0, blockWindowCounts[98])
		assert.Equal(t, 0, blockWindowCounts[99])
		assert.Equal(t, 0, blockWindowCounts[100])
		assert.Equal(t, 0, blockWindowCounts[101])
		assert.Equal(t, 50, blockWindowCounts[102])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 50, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 0, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 0, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[98])
		assert.Equal(t, 0, blockWindowCounts[101])
		assert.Equal(t, 0, blockWindowCounts[102]) // all windows exhausted
	})

	t.Run("minimum guaranteed for all windows including an incomplete window followed by best effort", func(t *testing.T) {
		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i <= end; i++ {
				logsToAdd := 100
				if i >= 100 {
					logsToAdd = 1
				}
				for j := 0; j < logsToAdd; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i)),
						BlockNumber: i,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 102, nil // make the latest window incomplete
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(42161)) // arbitrum
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(42161), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		blockWindowCounts := map[int64]int{}

		bufV1 := provider.bufferV1.(*logBuffer)

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		// window limit = 10 x 1 log limit, 10 logs per window per upkeep, 5 upkeeps = 50 logs in each window
		assert.Equal(t, 50, blockWindowCounts[0]) // block 0 is outside the block threshold of 1 and is not enqueued
		assert.Equal(t, 50, blockWindowCounts[4])
		assert.Equal(t, 50, blockWindowCounts[8])
		assert.Equal(t, 15, blockWindowCounts[100]) // the block window starting at block 100 is only 3/4 complete as of block 102

		// 25 4-block windows make up the 100 block range, 10 logs per upkeep, plus 15 logs from a partial window, so each upkeep has a total of 253 logs in the buffer
		assert.Equal(t, 253, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 253, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 253, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 253, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 253, countLogs(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 233, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 233, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 233, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 233, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 233, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		// the min commitment logs (1 per upkeep per block window) is dequeued from the first 20 block windows (5 upkeeps * 1 log each * 20 windows = 100 logs)
		assert.Equal(t, 45, blockWindowCounts[0])
		assert.Equal(t, 45, blockWindowCounts[4])
		assert.Equal(t, 45, blockWindowCounts[8])
		assert.Equal(t, 45, blockWindowCounts[76])
		assert.Equal(t, 50, blockWindowCounts[80]) // the 21st block window starts at block 80, and has not yet been dequeued
		assert.Equal(t, 15, blockWindowCounts[100])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 213, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 213, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 213, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 213, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 213, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 45, blockWindowCounts[80])
		assert.Equal(t, 45, blockWindowCounts[84])
		assert.Equal(t, 45, blockWindowCounts[88])
		assert.Equal(t, 45, blockWindowCounts[92])
		assert.Equal(t, 45, blockWindowCounts[96])
		assert.Equal(t, 10, blockWindowCounts[100]) // the 6 remaining block windows have the min dequeue happen, for 30 logs

		assert.Equal(t, 0, blockWindowCounts[0])  // the dequeue process exhausted min guaranteed logs and operated as best effort back on block window 0, dequeuing 45 logs
		assert.Equal(t, 20, blockWindowCounts[4]) // 30 logs of min commitment and 45 best effort gives us 75 logs, so the dequeue process dequeued 25 logs from the second block window
		assert.Equal(t, 45, blockWindowCounts[8])

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 193, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 193, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 193, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 193, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 193, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				startWindow, _ := getBlockWindow(blockNumber, 4)

				blockWindowCounts[startWindow] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[0])
		assert.Equal(t, 0, blockWindowCounts[4])   // best effort dequeued 20 logs from this window...
		assert.Equal(t, 0, blockWindowCounts[8])   // ..all remaining 45 logs from this window...
		assert.Equal(t, 10, blockWindowCounts[12]) // ...and 35 logs from this window
		assert.Equal(t, 45, blockWindowCounts[16])

	})

	t.Run("sparsely populated blocks", func(t *testing.T) {
		oldMaxPayloads := MaxPayloads
		MaxPayloads = 10
		defer func() {
			MaxPayloads = oldMaxPayloads
		}()

		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		upkeepOmittedOnBlocks := map[int][]int64{
			1: {5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100},                                                                                                                      // upkeep 1 won't have logs on 20 blocks
			2: {2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100}, // upkeep 2 won't have logs on 50 blocks
			3: {3, 13, 23, 33, 43, 53, 63, 73, 83, 93},                                                                                                                                                               // upkeep 3 won't appear on 10 blocks
			4: {1, 25, 50, 75, 100},                                                                                                                                                                                  // upkeep 4 won't appear on 5 blocks
			5: {},                                                                                                                                                                                                    // upkeep 5 appears on all blocks
		}

		callCount := 0
		// this gets called once per upkeep ID
		logGenerator := func(start, end int64) []logpoller.Log {
			callCount++
			var res []logpoller.Log
		outer:
			for i := start; i < end; i++ {
				for _, skip := range upkeepOmittedOnBlocks[callCount] {
					if skip == i+1 {
						continue outer
					}
				}
				res = append(res, logpoller.Log{
					LogIndex:    i,
					BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
					BlockNumber: i + 1,
				})
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		blockWindowCounts := map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 4, blockWindowCounts[1])
		assert.Equal(t, 4, blockWindowCounts[2])
		assert.Equal(t, 4, blockWindowCounts[3])
		assert.Equal(t, 4, blockWindowCounts[4])
		assert.Equal(t, 4, blockWindowCounts[5])
		assert.Equal(t, 4, blockWindowCounts[6])
		assert.Equal(t, 5, blockWindowCounts[7]) // block 7 is the first block to contain 1 log for all upkeeps

		assert.Equal(t, 80, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 50, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 90, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 95, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 100, countLogs(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps based on availability
		assert.Equal(t, 77, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 48, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 88, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 94, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 98, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[2])
		assert.Equal(t, 2, blockWindowCounts[3])
		assert.Equal(t, 4, blockWindowCounts[4])
		assert.Equal(t, 4, blockWindowCounts[5])
		assert.Equal(t, 4, blockWindowCounts[6])
		assert.Equal(t, 5, blockWindowCounts[7]) // block 7 is the first block to contain 1 log for all upkeeps

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps based on availability
		assert.Equal(t, 76, countLogs(bufV1.queues["1"].logs))
		assert.Equal(t, 47, countLogs(bufV1.queues["2"].logs))
		assert.Equal(t, 86, countLogs(bufV1.queues["3"].logs))
		assert.Equal(t, 91, countLogs(bufV1.queues["4"].logs))
		assert.Equal(t, 95, countLogs(bufV1.queues["5"].logs))

		blockWindowCounts = map[int64]int{}

		for _, q := range bufV1.queues {
			for blockNumber, logs := range q.logs {
				blockWindowCounts[blockNumber] += len(logs)
			}
		}

		assert.Equal(t, 0, blockWindowCounts[1])
		assert.Equal(t, 0, blockWindowCounts[2])
		assert.Equal(t, 0, blockWindowCounts[3])
		assert.Equal(t, 0, blockWindowCounts[4])
		assert.Equal(t, 0, blockWindowCounts[5])
		assert.Equal(t, 4, blockWindowCounts[6])
		assert.Equal(t, 5, blockWindowCounts[7]) // block 7 is the first block to contain 1 log for all upkeeps
	})
}

type mockedPacker struct {
}

func (p *mockedPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	return log.Data, nil
}
