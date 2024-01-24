package logpoller

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

// Validate that filters stored in log_filters_table match the filters stored in memory
func validateFiltersTable(t *testing.T, lp *logPoller, orm *DbORM) {
	filters, err := orm.LoadFilters()
	require.NoError(t, err)
	require.Equal(t, len(filters), len(lp.filters))
	for name, dbFilter := range filters {
		dbFilter := dbFilter
		memFilter, ok := lp.filters[name]
		require.True(t, ok)
		assert.Truef(t, memFilter.Contains(&dbFilter),
			"in-memory Filter %s is missing some addresses or events from db Filter table", name)
		assert.Truef(t, dbFilter.Contains(&memFilter), "db Filter table %s is missing some addresses or events from in-memory Filter", name)
	}
}

func TestLogPoller_RegisterFilter(t *testing.T) {
	t.Parallel()
	a1 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	a2 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	// Set up a test chain with a log emitting contract deployed.
	lp := NewLogPoller(orm, nil, lggr, time.Hour, false, 1, 1, 2, 1000)

	// We expect a zero Filter if nothing registered yet.
	f := lp.Filter(nil, nil, nil)
	require.Equal(t, 1, len(f.Addresses))
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), f.Addresses[0])

	err := lp.RegisterFilter(Filter{"Emitter Log 1", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}, 0})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe EventSigs
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}, 0})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe Addresses
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2 dupe", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}, 0})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Address required.
	err = lp.RegisterFilter(Filter{"no address", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{}, 0})
	require.Error(t, err)
	// Event required
	err = lp.RegisterFilter(Filter{"No event", []common.Hash{}, []common.Address{a1}, 0})
	require.Error(t, err)
	validateFiltersTable(t, lp, orm)

	// Removing non-existence Filter should log error but return nil
	err = lp.UnregisterFilter("Filter doesn't exist")
	require.NoError(t, err)
	require.Equal(t, observedLogs.Len(), 1)
	require.Contains(t, observedLogs.TakeAll()[0].Entry.Message, "not found")

	// Check that all filters are still there
	_, ok := lp.filters["Emitter Log 1"]
	require.True(t, ok, "'Emitter Log 1 Filter' missing")
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.True(t, ok, "'Emitter Log 1 + 2' Filter missing")
	_, ok = lp.filters["Emitter Log 1 + 2 dupe"]
	require.True(t, ok, "'Emitter Log 1 + 2 dupe' Filter missing")

	// Removing an existing Filter should remove it from both memory and db
	err = lp.UnregisterFilter("Emitter Log 1 + 2")
	require.NoError(t, err)
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.False(t, ok, "'Emitter Log 1 Filter' should have been removed by UnregisterFilter()")
	require.Len(t, lp.filters, 2)
	validateFiltersTable(t, lp, orm)

	err = lp.UnregisterFilter("Emitter Log 1 + 2 dupe")
	require.NoError(t, err)
	err = lp.UnregisterFilter("Emitter Log 1")
	require.NoError(t, err)
	assert.Len(t, lp.filters, 0)
	filters, err := lp.orm.LoadFilters()
	require.NoError(t, err)
	assert.Len(t, filters, 0)

	// Make sure cache was invalidated
	assert.Len(t, lp.Filter(nil, nil, nil).Addresses, 1)
	assert.Equal(t, lp.Filter(nil, nil, nil).Addresses[0], common.HexToAddress("0x0000000000000000000000000000000000000000"))
	assert.Len(t, lp.Filter(nil, nil, nil).Topics, 1)
	assert.Len(t, lp.Filter(nil, nil, nil).Topics[0], 0)
}

func TestLogPoller_ConvertLogs(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	topics := []common.Hash{EmitterABI.Events["Log1"].ID}

	var cases = []struct {
		name     string
		logs     []types.Log
		blocks   []LogPollerBlock
		expected int
	}{
		{"SingleBlock",
			[]types.Log{{Topics: topics}, {Topics: topics}},
			[]LogPollerBlock{{BlockTimestamp: time.Now()}},
			2},
		{"BlockList",
			[]types.Log{{Topics: topics}, {Topics: topics}, {Topics: topics}},
			[]LogPollerBlock{{BlockTimestamp: time.Now()}},
			3},
		{"EmptyList",
			[]types.Log{},
			[]LogPollerBlock{},
			0},
		{"TooManyBlocks",
			[]types.Log{{}},
			[]LogPollerBlock{{}, {}},
			0},
		{"TooFewBlocks",
			[]types.Log{{}, {}, {}},
			[]LogPollerBlock{{}, {}},
			0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			logs := convertLogs(c.logs, c.blocks, lggr, big.NewInt(53))
			require.Len(t, logs, c.expected)
			for i := 0; i < c.expected; i++ {
				if len(c.blocks) == 1 {
					assert.Equal(t, c.blocks[0].BlockTimestamp, logs[i].BlockTimestamp)
				} else {
					assert.Equal(t, logs[i].BlockTimestamp, c.blocks[i].BlockTimestamp)
				}
			}
		})
	}
}

func TestFilterName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "a - b:c:d", FilterName("a", "b", "c", "d"))
	assert.Equal(t, "empty args test", FilterName("empty args test"))
}

func TestLogPoller_BackupPollerStartup(t *testing.T) {
	addr := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")
	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	head := evmtypes.Head{Number: 3}
	events := []common.Hash{EmitterABI.Events["Log1"].ID}
	log1 := types.Log{
		Index:       0,
		BlockHash:   common.Hash{},
		BlockNumber: uint64(3),
		Topics:      events,
		Address:     addr,
		TxHash:      common.HexToHash("0x1234"),
		Data:        EvmWord(uint64(300)).Bytes(),
	}

	ec := evmclimocks.NewClient(t)
	ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)
	ec.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{log1}, nil)
	ec.On("ConfiguredChainID").Return(chainID, nil)

	ctx := testutils.Context(t)

	lp := NewLogPoller(orm, ec, lggr, 1*time.Hour, false, 2, 3, 2, 1000)
	lp.BackupPollAndSaveLogs(ctx, 100)
	assert.Equal(t, int64(0), lp.backupPollerNextBlock)
	assert.Equal(t, 1, observedLogs.FilterMessageSnippet("ran before first successful log poller run").Len())

	lp.PollAndSaveLogs(ctx, 3)

	lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(ctx))
	require.NoError(t, err)
	require.Equal(t, int64(3), lastProcessed.BlockNumber)

	lp.BackupPollAndSaveLogs(ctx, 100)
	assert.Equal(t, int64(1), lp.backupPollerNextBlock) // Ensure non-negative!
}

func TestLogPoller_Replay(t *testing.T) {
	t.Parallel()
	addr := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestObserved(t, zapcore.ErrorLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	head := evmtypes.Head{Number: 4}
	events := []common.Hash{EmitterABI.Events["Log1"].ID}
	log1 := types.Log{
		Index:       0,
		BlockHash:   common.Hash{},
		BlockNumber: uint64(head.Number),
		Topics:      events,
		Address:     addr,
		TxHash:      common.HexToHash("0x1234"),
		Data:        EvmWord(uint64(300)).Bytes(),
	}

	ec := evmclimocks.NewClient(t)
	ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)
	ec.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{log1}, nil).Once()
	ec.On("ConfiguredChainID").Return(chainID, nil)
	lp := NewLogPoller(orm, ec, lggr, time.Hour, false, 3, 3, 3, 20)

	// process 1 log in block 3
	lp.PollAndSaveLogs(testutils.Context(t), 4)
	latest, err := lp.LatestBlock()
	require.NoError(t, err)
	require.Equal(t, int64(4), latest.BlockNumber)

	t.Run("abort before replayStart received", func(t *testing.T) {
		// Replay() should abort immediately if caller's context is cancelled before request signal is read
		ctx, cancel := context.WithCancel(testutils.Context(t))
		cancel()
		err = lp.Replay(ctx, 3)
		assert.ErrorIs(t, err, ErrReplayRequestAborted)
	})

	recvStartReplay := func(ctx context.Context, block int64) {
		select {
		case fromBlock := <-lp.replayStart:
			assert.Equal(t, block, fromBlock)
		case <-ctx.Done():
			assert.NoError(t, ctx.Err(), "Timed out waiting to receive replay request from lp.replayStart")
		}
	}

	// Replay() should return error code received from replayComplete
	t.Run("returns error code on replay complete", func(t *testing.T) {
		ctx := testutils.Context(t)
		anyErr := errors.New("any error")
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(ctx, 1)
			lp.replayComplete <- anyErr
		}()
		assert.ErrorIs(t, lp.Replay(ctx, 1), anyErr)
		<-done
	})

	// Replay() should return ErrReplayInProgress if caller's context is cancelled after replay has begun
	t.Run("late abort returns ErrReplayInProgress", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), time.Second) // Intentionally abort replay after 1s
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(ctx, 4)
			cancel()
		}()
		assert.ErrorIs(t, lp.Replay(ctx, 4), ErrReplayInProgress)
		<-done
		lp.replayComplete <- nil
		lp.wg.Wait()
	})

	// Main lp.run() loop shouldn't get stuck if client aborts
	t.Run("client abort doesnt hang run loop", func(t *testing.T) {
		lp.backupPollerNextBlock = 0

		ctx := testutils.Context(t)

		pass := make(chan struct{})
		cancelled := make(chan struct{})

		rctx, rcancel := context.WithCancel(testutils.Context(t))
		var wg sync.WaitGroup
		defer func() { wg.Wait() }()
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.ErrorIs(t, lp.Replay(rctx, 4), ErrReplayInProgress)
				close(cancelled)
			}()
		})
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			rcancel()
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case lp.replayStart <- 4:
					close(pass)
				case <-ctx.Done():
					return
				}
			}()
			// We cannot return until we're sure that Replay() received the cancellation signal,
			// otherwise replayComplete<- might be sent first
			<-cancelled
		})

		ec.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{log1}, nil).Maybe() // in case task gets delayed by >= 100ms

		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		select {
		case <-ctx.Done():
			t.Errorf("timed out waiting for lp.run() to respond to second replay event")
		case <-pass:
		}
	})

	// remove Maybe expectation from prior subtest, as it will override all expected calls in future subtests
	ec.On("FilterLogs", mock.Anything, mock.Anything).Unset()

	// run() should abort if log poller shuts down while replay is in progress
	t.Run("shutdown during replay", func(t *testing.T) {
		lp.backupPollerNextBlock = 0

		pass := make(chan struct{})
		done := make(chan struct{})
		defer func() { <-done }()

		ctx := testutils.Context(t)
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			go func() {
				defer close(done)
				select {
				case lp.replayStart <- 4:
				case <-ctx.Done():
				}
			}()
		})
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			lp.cancel()
			close(pass)
		})
		ec.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{log1}, nil).Maybe() // in case task gets delayed by >= 100ms

		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		select {
		case <-ctx.Done():
			t.Error("timed out waiting for lp.run() to respond to shutdown event during replay")
		case <-pass:
		}
	})

	// ReplayAsync should return as soon as replayStart is received
	t.Run("ReplayAsync success", func(t *testing.T) {
		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		lp.ReplayAsync(1)

		recvStartReplay(testutils.Context(t), 1)
	})

	t.Run("ReplayAsync error", func(t *testing.T) {
		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		anyErr := errors.New("async error")
		observedLogs.TakeAll()

		lp.ReplayAsync(4)
		recvStartReplay(testutils.Context(t), 4)

		select {
		case lp.replayComplete <- anyErr:
			time.Sleep(2 * time.Second)
		case <-lp.ctx.Done():
			t.Error("timed out waiting to send replaceComplete")
		}
		require.Equal(t, 1, observedLogs.Len())
		assert.Equal(t, observedLogs.All()[0].Message, anyErr.Error())
	})
}

func (lp *logPoller) reset() {
	lp.StateMachine = services.StateMachine{}
	lp.ctx, lp.cancel = context.WithCancel(context.Background())
}

func Test_latestBlockAndFinalityDepth(t *testing.T) {
	lggr := logger.Test(t)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	t.Run("pick latest block from chain and use finality from config with finality disabled", func(t *testing.T) {
		head := evmtypes.Head{Number: 4}
		finalityDepth := int64(3)
		ec := evmclimocks.NewClient(t)
		ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)

		lp := NewLogPoller(orm, ec, lggr, time.Hour, false, finalityDepth, 3, 3, 20)
		latestBlock, lastFinalizedBlockNumber, err := lp.latestBlocks(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, latestBlock.Number, head.Number)
		require.Equal(t, finalityDepth, latestBlock.Number-lastFinalizedBlockNumber)
	})

	t.Run("finality tags in use", func(t *testing.T) {
		t.Run("client returns data properly", func(t *testing.T) {
			expectedLatestBlockNumber := int64(20)
			expectedLastFinalizedBlockNumber := int64(12)
			ec := evmclimocks.NewClient(t)
			ec.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 2 &&
					reflect.DeepEqual(b[0].Args, []interface{}{"latest", false}) &&
					reflect.DeepEqual(b[1].Args, []interface{}{"finalized", false})
			})).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				// Latest block details
				*(elems[0].Result.(*evmtypes.Head)) = evmtypes.Head{Number: expectedLatestBlockNumber, Hash: utils.RandomBytes32()}
				// Finalized block details
				*(elems[1].Result.(*evmtypes.Head)) = evmtypes.Head{Number: expectedLastFinalizedBlockNumber, Hash: utils.RandomBytes32()}
			})

			lp := NewLogPoller(orm, ec, lggr, time.Hour, true, 3, 3, 3, 20)

			latestBlock, lastFinalizedBlockNumber, err := lp.latestBlocks(testutils.Context(t))
			require.NoError(t, err)
			require.Equal(t, expectedLatestBlockNumber, latestBlock.Number)
			require.Equal(t, expectedLastFinalizedBlockNumber, lastFinalizedBlockNumber)
		})

		t.Run("client returns error for at least one of the calls", func(t *testing.T) {
			ec := evmclimocks.NewClient(t)
			ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				// Latest block details
				*(elems[0].Result.(*evmtypes.Head)) = evmtypes.Head{Number: 10}
				// Finalized block details
				elems[1].Error = fmt.Errorf("some error")
			})

			lp := NewLogPoller(orm, ec, lggr, time.Hour, true, 3, 3, 3, 20)
			_, _, err := lp.latestBlocks(testutils.Context(t))
			require.Error(t, err)
		})

		t.Run("BatchCall returns an error", func(t *testing.T) {
			ec := evmclimocks.NewClient(t)
			ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(fmt.Errorf("some error"))

			lp := NewLogPoller(orm, ec, lggr, time.Hour, true, 3, 3, 3, 20)
			_, _, err := lp.latestBlocks(testutils.Context(t))
			require.Error(t, err)
		})
	})
}

func benchmarkFilter(b *testing.B, nFilters, nAddresses, nEvents int) {
	lggr := logger.Test(b)
	lp := NewLogPoller(nil, nil, lggr, 1*time.Hour, false, 2, 3, 2, 1000)
	for i := 0; i < nFilters; i++ {
		var addresses []common.Address
		var events []common.Hash
		for j := 0; j < nAddresses; j++ {
			addresses = append(addresses, common.BigToAddress(big.NewInt(int64(j+1))))
		}
		for j := 0; j < nEvents; j++ {
			events = append(events, common.BigToHash(big.NewInt(int64(j+1))))
		}
		err := lp.RegisterFilter(Filter{Name: "my Filter", EventSigs: events, Addresses: addresses})
		require.NoError(b, err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lp.Filter(nil, nil, nil)
	}
}

func BenchmarkFilter10_1(b *testing.B) {
	benchmarkFilter(b, 10, 1, 1)
}
func BenchmarkFilter100_10(b *testing.B) {
	benchmarkFilter(b, 100, 10, 10)
}
func BenchmarkFilter1000_100(b *testing.B) {
	benchmarkFilter(b, 1000, 100, 100)
}
