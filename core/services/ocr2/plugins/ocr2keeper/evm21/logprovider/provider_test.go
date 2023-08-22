package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"golang.org/x/time/rate"
)

func TestLogEventProvider_GetEntries(t *testing.T) {
	p := New(logger.TestLogger(t), nil, &mockedPacker{}, nil)

	_, f := newEntry(p, 1)
	p.lock.Lock()
	p.active[f.id.String()] = f
	p.lock.Unlock()

	t.Run("no entries", func(t *testing.T) {
		entries := p.getEntries(0, false, big.NewInt(0))
		require.Len(t, entries, 1)
		require.Equal(t, len(entries[0].filter.Addresses), 0)
	})

	t.Run("has entry with lower lastPollBlock", func(t *testing.T) {
		entries := p.getEntries(0, false, f.id)
		require.Len(t, entries, 1)
		require.Greater(t, len(entries[0].filter.Addresses), 0)
		entries = p.getEntries(10, false, f.id)
		require.Len(t, entries, 1)
		require.Greater(t, len(entries[0].filter.Addresses), 0)
	})

	t.Run("has entry with higher lastPollBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.lastPollBlock = 3
		p.lock.Lock()
		p.active[f.id.String()] = f
		p.lock.Unlock()

		entries := p.getEntries(1, false, f.id)
		require.Len(t, entries, 1)
		require.Equal(t, len(entries[0].filter.Addresses), 0)

		entries = p.getEntries(1, true, f.id)
		require.Len(t, entries, 1)
		require.Greater(t, len(entries[0].filter.Addresses), 0)
	})
}

func TestLogEventProvider_UpdateEntriesLastPoll(t *testing.T) {
	p := New(logger.TestLogger(t), nil, &mockedPacker{}, nil)

	n := 10

	entries := map[string]upkeepFilterEntry{}
	for i := 0; i < n; i++ {
		_, f := newEntry(p, i+1)
		entries[f.id.String()] = f
	}
	p.lock.Lock()
	p.active = entries
	p.lock.Unlock()

	t.Run("no entries", func(t *testing.T) {
		_, f := newEntry(p, n*2)
		f.lastPollBlock = 10
		p.updateEntriesLastPoll([]*upkeepFilterEntry{&f})

		p.lock.RLock()
		defer p.lock.RUnlock()
		for _, f := range p.active {
			require.Equal(t, int64(0), f.lastPollBlock)
		}
	})

	t.Run("update entries", func(t *testing.T) {
		_, f2 := newEntry(p, n-2)
		f2.lastPollBlock = 10
		_, f1 := newEntry(p, n-1)
		f1.lastPollBlock = 10
		p.updateEntriesLastPoll([]*upkeepFilterEntry{&f1, &f2})

		p.lock.RLock()
		e := p.active[f1.id.String()]
		require.Equal(t, int64(10), e.lastPollBlock)
		e = p.active[f2.id.String()]
		require.Equal(t, int64(10), e.lastPollBlock)
		p.lock.RUnlock()
		// update with same block
		p.updateEntriesLastPoll([]*upkeepFilterEntry{&f1})

		_, f := newEntry(p, 1)
		p.lock.RLock()
		defer p.lock.RUnlock()
		e = p.active[f.id.String()]
		require.Equal(t, int64(0), e.lastPollBlock)
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tick := 10 * time.Millisecond
			p := New(logger.TestLogger(t), mp, &mockedPacker{}, &LogEventProviderOptions{
				ReadMaxBatchSize: tc.maxBatchSize,
				ReadInterval:     tick,
			})

			var ids []*big.Int
			p.lock.Lock()
			for i, id := range tc.ids {
				_, f := newEntry(p, id, tc.addrs[i])
				p.active[f.id.String()] = f
				ids = append(ids, f.id)
			}
			p.lock.Unlock()

			reads := make(chan []*big.Int, 100)

			go func(ctx context.Context) {
				_ = p.scheduleReadJobs(ctx, func(ids []*big.Int) {
					select {
					case reads <- ids:
					default:
						t.Log("dropped ids")
					}
				})
			}(ctx)

			timeout := tick*time.Duration((1+len(tc.ids)/tc.maxBatchSize))*4 + 1
			<-time.After(timeout)
			timeoutTicker := time.NewTicker(timeout)
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
				}
			}

			require.Len(t, got, len(ids))
			for _, id := range ids {
				require.GreaterOrEqual(t, got[id.String()], 1, "id %s", id.String())
			}
		})
	}
}

func TestLogEventProvider_ReadLogs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mp := new(mocks.LogPoller)

	mp.On("RegisterFilter", mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("LatestBlock", mock.Anything).Return(int64(1), nil)
	mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 1,
			TxHash:      common.HexToHash("0x1"),
		},
	}, nil)

	p := New(logger.TestLogger(t), mp, &mockedPacker{}, nil)

	var ids []*big.Int
	for i := 0; i < 10; i++ {
		cfg, f := newEntry(p, i+1)
		ids = append(ids, f.id)
		require.NoError(t, p.RegisterFilter(f.id, cfg))
	}

	t.Run("no entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, false, big.NewInt(999999)))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 0)
	})

	t.Run("has entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, true, ids[:2]...))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 2)
	})

	// TODO: test rate limiting

}

func newEntry(p *logEventProvider, i int, args ...string) (LogTriggerConfig, upkeepFilterEntry) {
	id := ocr2keepers.UpkeepIdentifier(append(common.LeftPadBytes([]byte{1}, 16), []byte(fmt.Sprintf("%d", i))...))
	uid := big.NewInt(0).SetBytes(id)
	for len(args) < 2 {
		args = append(args, "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d")
	}
	addr, topic0 := args[0], args[1]
	cfg := LogTriggerConfig{
		ContractAddress: common.HexToAddress(addr),
		FilterSelector:  0,
		Topic0:          common.HexToHash(topic0),
	}
	f := upkeepFilterEntry{
		id:            uid,
		filter:        p.newLogFilter(uid, cfg),
		cfg:           cfg,
		blockLimiter:  rate.NewLimiter(p.opts.BlockRateLimit, p.opts.BlockLimitBurst),
		lastPollBlock: 0,
	}
	return cfg, f
}

type mockedPacker struct {
}

func (p *mockedPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	return log.Data, nil
}
