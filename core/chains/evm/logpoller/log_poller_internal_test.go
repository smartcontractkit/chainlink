package logpoller

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

// Validate that filters stored in log_filters_table match the filters stored in memory
func validateFiltersTable(t *testing.T, lp *logPoller, orm *ORM) {
	filters, err := orm.LoadFilters()
	require.NoError(t, err)
	require.Equal(t, len(filters), len(lp.filters))
	for name, dbFilter := range filters {
		dbFilter := dbFilter
		memFilter, ok := lp.filters[name]
		require.True(t, ok)
		assert.True(t, memFilter.Contains(&dbFilter),
			fmt.Sprintf("in-memory Filter %s is missing some addresses or events from db Filter table", name))
		assert.True(t, dbFilter.Contains(&memFilter),
			fmt.Sprintf("db Filter table %s is missing some addresses or events from in-memory Filter", name))
	}
}

func TestLogPoller_RegisterFilter(t *testing.T) {
	t.Parallel()
	a1 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	a2 := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)

	orm := NewORM(chainID, db, lggr, pgtest.NewQConfig(true))
	lp := NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	filter := Filter{"test Filter", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}}
	err := lp.RegisterFilter(filter)
	require.Error(t, err, "RegisterFilter failed to save Filter to db")
	require.Equal(t, 1, observedLogs.Len())
	assertForeignConstraintError(t, observedLogs.TakeAll()[0], "evm_log_poller_filters", "evm_log_poller_filters_evm_chain_id_fkey")

	db.Close()
	db = pgtest.NewSqlxDB(t)
	orm = NewORM(chainID, db, lggr, pgtest.NewQConfig(true))

	// disable check that chain id exists for rest of tests
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS evm_log_poller_filters_evm_chain_id_fkey DEFERRED`)))
	// Set up a test chain with a log emitting contract deployed.

	lp = NewLogPoller(orm, nil, lggr, 15*time.Second, 1, 1, 2, 1000)

	// We expect a zero Filter if nothing registered yet.
	f := lp.Filter(nil, nil, nil)
	require.Equal(t, 1, len(f.Addresses))
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), f.Addresses[0])

	err = lp.RegisterFilter(Filter{"Emitter Log 1", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{a1}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe EventSigs
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Should de-dupe Addresses
	err = lp.RegisterFilter(Filter{"Emitter Log 1 + 2 dupe", []common.Hash{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}, []common.Address{a2}})
	require.NoError(t, err)
	assert.Equal(t, []common.Address{a1, a2}, lp.Filter(nil, nil, nil).Addresses)
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}, lp.Filter(nil, nil, nil).Topics)
	validateFiltersTable(t, lp, orm)

	// Address required.
	err = lp.RegisterFilter(Filter{"no address", []common.Hash{EmitterABI.Events["Log1"].ID}, []common.Address{}})
	require.Error(t, err)
	// Event required
	err = lp.RegisterFilter(Filter{"No event", []common.Hash{}, []common.Address{a1}})
	require.Error(t, err)
	validateFiltersTable(t, lp, orm)

	// Removing non-existence Filter should log error but return nil
	err = lp.UnregisterFilter("Filter doesn't exist", nil)
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
	err = lp.UnregisterFilter("Emitter Log 1 + 2", nil)
	require.NoError(t, err)
	_, ok = lp.filters["Emitter Log 1 + 2"]
	require.False(t, ok, "'Emitter Log 1 Filter' should have been removed by UnregisterFilter()")
	require.Len(t, lp.filters, 2)
	validateFiltersTable(t, lp, orm)

	err = lp.UnregisterFilter("Emitter Log 1 + 2 dupe", nil)
	require.NoError(t, err)
	err = lp.UnregisterFilter("Emitter Log 1", nil)
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

func assertForeignConstraintError(t *testing.T, observedLog observer.LoggedEntry,
	table string, constraint string) {

	assert.Equal(t, "SQL ERROR", observedLog.Entry.Message)

	field := observedLog.Context[0]
	require.Equal(t, zapcore.ErrorType, field.Type)
	err, ok := field.Interface.(error)
	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	require.True(t, ok)
	assert.Equal(t, "23503", pgErr.SQLState()) // foreign key constraint violation code
	assert.Equal(t, table, pgErr.TableName)
	assert.Equal(t, constraint, pgErr.ConstraintName)
}

func TestLogPoller_ConvertLogs(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

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
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
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

	lp := NewLogPoller(orm, ec, lggr, 1*time.Hour, 2, 3, 2, 1000)
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
	tctx := testutils.Context(t)

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
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
	lp := NewLogPoller(orm, ec, lggr, time.Hour, 3, 3, 3, 20)

	// process 1 log in block 3
	lp.PollAndSaveLogs(tctx, 4)
	latest, err := lp.LatestBlock()
	require.NoError(t, err)
	require.Equal(t, int64(4), latest)

	t.Run("abort before replayStart received", func(t *testing.T) {
		// Replay() should abort immediately if caller's context is cancelled before request signal is read
		ctx, cancel := context.WithCancel(tctx)
		cancel()
		err = lp.Replay(ctx, 3)
		assert.ErrorIs(t, err, ErrReplayRequestAborted)
	})

	recvStartReplay := func(parentCtx context.Context, block int64, withTimeout bool) {
		var err error
		var ctx context.Context
		var cancel context.CancelFunc
		if withTimeout {
			ctx, cancel = context.WithTimeout(parentCtx, time.Second)
		} else {
			ctx, cancel = context.WithCancel(parentCtx)
		}
		defer cancel()
		select {
		case fromBlock := <-lp.replayStart:
			assert.Equal(t, block, fromBlock)
		case <-ctx.Done():
			err = ctx.Err()
		}
		require.NoError(t, err, "Timed out waiting to receive replay request from lp.replayStart")
	}

	// Replay() should return error code received from replayComplete
	t.Run("returns error code on replay complete", func(t *testing.T) {
		anyErr := errors.New("any error")
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(tctx, 1, true)
			lp.replayComplete <- anyErr
		}()
		assert.ErrorIs(t, lp.Replay(tctx, 1), anyErr)
		<-done
	})

	// Replay() should return ErrReplayInProgress if caller's context is cancelled after replay has begun
	t.Run("late abort returns ErrReplayInProgress", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), time.Second)
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(ctx, 4, false)
			cancel()
		}()
		assert.ErrorIs(t, lp.Replay(ctx, 4), ErrReplayInProgress)
		<-done
	})

	// Main lp.run() loop shouldn't get stuck if client aborts
	t.Run("client abort doesnt hang run loop", func(t *testing.T) {
		lp.backupPollerNextBlock = 0

		timeout := time.After(1 * time.Second)
		lp.ctx, lp.cancel = context.WithCancel(tctx)
		ctx, cancel := context.WithCancel(tctx)

		var wg sync.WaitGroup
		wg.Add(2)
		ec.On("FilterLogs", lp.ctx, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			go func() {
				defer wg.Done()
				assert.ErrorIs(t, lp.Replay(ctx, 4), ErrReplayInProgress)
			}()
		})
		ec.On("FilterLogs", lp.ctx, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			go func() {
				defer wg.Done()
				cancel()
				lp.replayStart <- 4
			}()
		})
		ec.On("FilterLogs", lp.ctx, mock.Anything).Return([]types.Log{log1}, nil).Maybe()
		lp.wg.Add(1)
		done := make(chan struct{})
		go func() {
			defer close(done)
			lp.run()
		}()
		select {
		case <-timeout:
			assert.Fail(t, "lp.run() got stuck--failed to respond to second replay event within 1s")
		case <-utils.WaitGroupChan(&wg):
			lp.cancel()
		}
		<-done
	})

	lp.wg.Wait() // ensure logpoller has exited before continuing on to next test

	// run() should abort if log poller shuts down while replay is in progress
	t.Run("shutdown during replay", func(t *testing.T) {
		lp.backupPollerNextBlock = 0
		var wg sync.WaitGroup
		wg.Add(2)
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			go func() {
				defer wg.Done()
				lp.replayStart <- 4
			}()
		})
		ec.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log1}, nil).Run(func(args mock.Arguments) {
			defer wg.Done()
			lp.cancel()
		})
		ec.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{log1}, nil).Maybe() // in case task gets delayed by >= 100ms

		timeout := time.After(1 * time.Second)
		require.NoError(t, lp.Start(tctx))
		select {
		case <-timeout:
			assert.Fail(t, "lp.run() failed to respond to shutdown event during replay within 1s")
			lp.Close()
		case <-utils.WaitGroupChan(&wg):
			lp.wg.Wait()
		}
	})

	// ReplayAsync should return success as soon as replayStart is received
	t.Run("ReplayAsync success", func(t *testing.T) {
		lp.ctx, lp.cancel = context.WithTimeout(tctx, time.Second)
		go recvStartReplay(tctx, 1, true)
		lp.ReplayAsync(1)
		lp.replayComplete <- nil
	})

	t.Run("ReplayAsync error", func(t *testing.T) {
		lp.ctx, lp.cancel = context.WithTimeout(tctx, 2*time.Second)
		defer lp.cancel()
		anyErr := errors.New("async error")
		observedLogs.TakeAll()

		lp.ReplayAsync(4)
		recvStartReplay(tctx, 4, true)

		select {
		case lp.replayComplete <- anyErr:
			time.Sleep(2 * time.Second)
		case <-lp.ctx.Done():
			require.Fail(t, "failed to receive replayComplete signal")
		}
		require.Equal(t, 1, observedLogs.Len())
		assert.Equal(t, observedLogs.All()[0].Message, anyErr.Error())
	})
}

func benchmarkFilter(b *testing.B, nFilters, nAddresses, nEvents int) {
	lggr := logger.TestLogger(b)
	lp := NewLogPoller(nil, nil, lggr, 1*time.Hour, 2, 3, 2, 1000)
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
