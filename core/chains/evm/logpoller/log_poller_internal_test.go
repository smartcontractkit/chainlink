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
	pkgerrors "github.com/pkg/errors"
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
)

var (
	EmitterABI, _ = abi.JSON(strings.NewReader(log_emitter.LogEmitterABI))
)

// Validate that filters stored in log_filters_table match the filters stored in memory
func validateFiltersTable(t *testing.T, lp *logPoller, orm ORM) {
	ctx := testutils.Context(t)
	filters, err := orm.LoadFilters(ctx)
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
	a1 := common.HexToAddress("0x2Ab9A2Dc53736b361b72D900Cdf9f78F9406FBBb")
	a2 := common.HexToAddress("0x2Ab9A2Dc53736b361b72D900Cdf9f78F9406FBBc")

	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)

	orm := NewORM(chainID, db, lggr)

	// Set up a test chain with a log emitting contract deployed.
	lpOpts := Opts{
		PollPeriod:               time.Hour,
		BackfillBatchSize:        1,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := NewLogPoller(orm, nil, lggr, lpOpts)

	// We expect empty list of reqs if nothing registered yet.
	reqs := lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 0)
	topics1 := [][]common.Hash{{EmitterABI.Events["Log1"].ID}}
	topics2 := [][]common.Hash{{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID}}
	err := lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 1, address 1", EventSigs: topics1[0], Addresses: []common.Address{a1}})
	require.NoError(t, err)
	reqs = lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 1)
	assert.Equal(t, []common.Address{a1}, reqs[0].Addresses())
	assert.Equal(t, [][]common.Hash{{EmitterABI.Events["Log1"].ID}}, reqs[0].Topics())
	validateFiltersTable(t, lp, orm)

	err = lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 1 + 2, address 2", EventSigs: topics2[0], Addresses: []common.Address{a2}})
	require.NoError(t, err)
	reqs = lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 2) // should have 2 separate reqs, since these cannot be merged
	getLogsReq1 := reqs[0]
	getLogsReq2 := reqs[1]
	addresses1 := reqs[0].Addresses()
	require.Len(t, addresses1, 1)
	addresses2 := reqs[1].Addresses()
	require.Len(t, addresses2, 1)
	assert.Equal(t, a1, addresses1[0])
	assert.Equal(t, a2, addresses2[0])
	assert.Equal(t, topics1, getLogsReq1.Topics())
	assert.Equal(t, topics2, getLogsReq2.Topics())
	validateFiltersTable(t, lp, orm)

	// Reqs should not change when a duplicate filter is added
	err = lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 1 + 2 dupe", EventSigs: topics2[0], Addresses: []common.Address{a2}})
	require.NoError(t, err)
	reqs = lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 2)
	assert.Equal(t, []common.Address{a1}, reqs[0].Addresses())
	assert.Equal(t, topics1, reqs[0].Topics())
	assert.Equal(t, []common.Address{a2}, reqs[1].Addresses())
	assert.Equal(t, topics2, reqs[1].Topics())
	validateFiltersTable(t, lp, orm)

	// Same address as "Emitter Log 1 + 2, address 2", but only looking for Log1 (subset of topics). Recs should not change
	err = lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 2, address 2", EventSigs: evmtypes.HashArray{topics2[0][1]}, Addresses: []common.Address{a2}})
	require.NoError(t, err)
	reqs = lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 2)
	assert.Equal(t, []common.Address{a1}, reqs[0].Addresses())
	assert.Equal(t, topics1, reqs[0].Topics())
	assert.Equal(t, []common.Address{a2}, reqs[1].Addresses())
	assert.Equal(t, topics2, reqs[1].Topics())
	validateFiltersTable(t, lp, orm)

	// Same address as "Emitter Log 1, address 1", but different topics lists. Should create new rec
	err = lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 2, address 1", EventSigs: evmtypes.HashArray{topics2[0][1]}, Addresses: []common.Address{a1}})
	require.NoError(t, err)
	reqs = lp.EthGetLogsReqs(nil, nil, nil)
	require.Len(t, reqs, 3)
	assert.Equal(t, []common.Address{a1}, reqs[1].Addresses())
	assert.Equal(t, [][]common.Hash{{topics2[0][1]}}, reqs[1].Topics())
	validateFiltersTable(t, lp, orm)

	// Address required.
	err = lp.RegisterFilter(ctx, Filter{Name: "no address", EventSigs: []common.Hash{EmitterABI.Events["Log1"].ID}})
	require.Error(t, err)
	// Event required
	err = lp.RegisterFilter(ctx, Filter{Name: "No event", Addresses: []common.Address{a1}})
	require.Error(t, err)
	validateFiltersTable(t, lp, orm)

	// Removing non-existence Filter should log error but return nil
	err = lp.UnregisterFilter(ctx, "Filter doesn't exist")
	require.NoError(t, err)
	require.Equal(t, observedLogs.Len(), 1)
	require.Contains(t, observedLogs.TakeAll()[0].Entry.Message, "not found")

	// Check that all filters are still there
	_, ok := lp.filters["Emitter Log 1, address 1"]
	require.True(t, ok, "'Emitter Log 1, address 1' Filter missing")
	_, ok = lp.filters["Emitter Log 1 + 2, address 2"]
	require.True(t, ok, "'Emitter Log 1 + 2, address 2' Filter missing")
	_, ok = lp.filters["Emitter Log 1 + 2 dupe"]
	require.True(t, ok, "'Emitter Log 1 + 2 dupe' Filter missing")
	_, ok = lp.filters["Emitter Log 2, address 2"]
	require.True(t, ok, "'Emitter Log 2, address 2' Filter missing")
	_, ok = lp.filters["Emitter Log 2, address 1"]
	require.True(t, ok, "'Emitter Log 2, address 1' Filter missing")

	// Removing an existing Filter should remove it from both memory and db
	err = lp.UnregisterFilter(ctx, "Emitter Log 1 + 2, address 2")
	require.NoError(t, err)
	_, ok = lp.filters["Emitter Log 1 + 2, address 2"]
	require.False(t, ok, "'Emitter Log 1 Filter' should have been removed by UnregisterFilter()")
	require.Len(t, lp.filters, 4)
	validateFiltersTable(t, lp, orm)

	err = lp.UnregisterFilter(ctx, "Emitter Log 1 + 2 dupe")
	require.NoError(t, err)
	err = lp.UnregisterFilter(ctx, "Emitter Log 1, address 1")
	require.NoError(t, err)
	err = lp.UnregisterFilter(ctx, "Emitter Log 2, address 1")
	require.NoError(t, err)
	err = lp.UnregisterFilter(ctx, "Emitter Log 2, address 2")
	require.NoError(t, err)
	assert.Len(t, lp.filters, 0)
	filters, err := lp.orm.LoadFilters(ctx)
	require.NoError(t, err)
	assert.Len(t, filters, 0)

	require.Len(t, lp.newFilters, 0)
	require.Len(t, lp.removedFilters, 5)

	lp.EthGetLogsReqs(nil, nil, nil)

	// Make sure cache was invalidated
	assert.Len(t, lp.cachedReqsByEventsTopicsKey, 0)
	assert.Len(t, lp.cachedReqsByAddress, 0)
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
	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr)

	head := evmtypes.Head{Number: 3}

	ec := evmclimocks.NewClient(t)
	ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)
	ec.On("ConfiguredChainID").Return(chainID, nil)

	ctx := testutils.Context(t)
	lpOpts := Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := NewLogPoller(orm, ec, lggr, lpOpts)
	lp.BackupPollAndSaveLogs(ctx)
	assert.Equal(t, int64(0), lp.backupPollerNextBlock)
	assert.Equal(t, 1, observedLogs.FilterMessageSnippet("ran before first successful log poller run").Len())

	lp.PollAndSaveLogs(ctx, 3)

	lastProcessed, err := lp.orm.SelectLatestBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(3), lastProcessed.BlockNumber)

	lp.BackupPollAndSaveLogs(ctx)
	assert.Equal(t, int64(1), lp.backupPollerNextBlock) // Ensure non-negative!
}

func TestLogPoller_Replay(t *testing.T) {
	t.Parallel()
	addr := common.HexToAddress("0x2ab9a2dc53736b361b72d900cdf9f78f9406fbbc")

	lggr, observedLogs := logger.TestObserved(t, zapcore.ErrorLevel)
	chainID := testutils.FixtureChainID
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(chainID, db, lggr)
	ctx := testutils.Context(t)

	head := evmtypes.Head{Number: 4}
	events := []common.Hash{EmitterABI.Events["Log1"].ID}

	ec := evmclimocks.NewClient(t)
	ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)
	ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Twice()
	ec.On("ConfiguredChainID").Return(chainID, nil)
	lpOpts := Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            3,
		BackfillBatchSize:        3,
		RpcBatchSize:             3,
		KeepFinalizedBlocksDepth: 20,
		BackupPollerBlockDelay:   100,
	}
	lp := NewLogPoller(orm, ec, lggr, lpOpts)

	err := lp.RegisterFilter(ctx, Filter{Name: "Emitter Log 1", EventSigs: events, Addresses: []common.Address{addr}})
	require.NoError(t, err)

	// process 1 log in block 3
	lp.PollAndSaveLogs(ctx, 4)
	latest, err := lp.LatestBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(4), latest.BlockNumber)
	require.Equal(t, int64(1), latest.FinalizedBlockNumber)

	t.Run("abort before replayStart received", func(t *testing.T) {
		// Replay() should abort immediately if caller's context is cancelled before request signal is read
		cancelCtx, cancel := context.WithCancel(testutils.Context(t))
		cancel()
		err = lp.Replay(cancelCtx, 3)
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
		anyErr := pkgerrors.New("any error")
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(ctx, 2)
			lp.replayComplete <- anyErr
		}()
		assert.ErrorIs(t, lp.Replay(ctx, 1), anyErr)
		<-done
	})

	// Replay() should return ErrReplayInProgress if caller's context is cancelled after replay has begun
	t.Run("late abort returns ErrReplayInProgress", func(t *testing.T) {
		cancelCtx, cancel := context.WithTimeout(testutils.Context(t), time.Second) // Intentionally abort replay after 1s
		done := make(chan struct{})
		go func() {
			defer close(done)
			recvStartReplay(cancelCtx, 4)
			cancel()
		}()
		assert.ErrorIs(t, lp.Replay(cancelCtx, 4), ErrReplayInProgress)
		<-done
		lp.replayComplete <- nil
		lp.wg.Wait()
	})

	// Main lp.run() loop shouldn't get stuck if client aborts
	t.Run("client abort doesnt hang run loop", func(t *testing.T) {
		lp.backupPollerNextBlock = 0

		pass := make(chan struct{})
		cancelled := make(chan struct{})

		rctx, rcancel := context.WithCancel(testutils.Context(t))
		var wg sync.WaitGroup
		defer func() { wg.Wait() }()
		ec.On("BatchCallContext", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.ErrorIs(t, lp.Replay(rctx, 4), ErrReplayInProgress)
				close(cancelled)
			}()
		})
		ec.On("BatchCallContext", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
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

		ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Maybe() // in case task gets delayed by >= 100ms

		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		select {
		case <-ctx.Done():
			t.Errorf("timed out waiting for lp.run() to respond to second replay event")
		case <-pass:
		}
	})

	// remove Maybe expectation from prior subtest, as it will override all expected calls in future subtests
	ec.On("BatchCallContext", mock.Anything, mock.Anything).Unset()

	// run() should abort if log poller shuts down while replay is in progress
	t.Run("shutdown during replay", func(t *testing.T) {
		lp.backupPollerNextBlock = 0

		pass := make(chan struct{})
		done := make(chan struct{})
		defer func() { <-done }()

		ec.On("BatchCallContext", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			go func() {
				defer close(done)
				select {
				case lp.replayStart <- 4:
				case <-ctx.Done():
				}
			}()
		})
		ec.On("BatchCallContext", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			lp.cancel()
			close(pass)
		})
		ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Maybe() // in case task gets delayed by >= 100ms

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

		recvStartReplay(testutils.Context(t), 2)
	})

	t.Run("ReplayAsync error", func(t *testing.T) {
		t.Cleanup(lp.reset)
		servicetest.Run(t, lp)

		anyErr := pkgerrors.New("async error")
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

	t.Run("run regular replay when there are not blocks in db", func(t *testing.T) {
		err := lp.orm.DeleteLogsAndBlocksAfter(ctx, 0)
		require.NoError(t, err)

		lp.ReplayAsync(1)
		recvStartReplay(testutils.Context(t), 1)
	})

	t.Run("run only backfill when everything is finalized", func(t *testing.T) {
		err := lp.orm.DeleteLogsAndBlocksAfter(ctx, 0)
		require.NoError(t, err)

		err = lp.orm.InsertBlock(ctx, head.Hash, head.Number, head.Timestamp, head.Number)
		require.NoError(t, err)

		err = lp.Replay(ctx, 1)
		require.NoError(t, err)
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
	orm := NewORM(chainID, db, lggr)
	ctx := testutils.Context(t)

	lpOpts := Opts{
		PollPeriod:               time.Hour,
		BackfillBatchSize:        3,
		RpcBatchSize:             3,
		KeepFinalizedBlocksDepth: 20,
	}

	t.Run("pick latest block from chain and use finality from config with finality disabled", func(t *testing.T) {
		head := evmtypes.Head{Number: 4}

		lpOpts.UseFinalityTag = false
		lpOpts.FinalityDepth = int64(3)
		ec := evmclimocks.NewClient(t)
		ec.On("HeadByNumber", mock.Anything, mock.Anything).Return(&head, nil)

		lp := NewLogPoller(orm, ec, lggr, lpOpts)
		latestBlock, lastFinalizedBlockNumber, err := lp.latestBlocks(ctx)
		require.NoError(t, err)
		require.Equal(t, latestBlock.Number, head.Number)
		require.Equal(t, lpOpts.FinalityDepth, latestBlock.Number-lastFinalizedBlockNumber)
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

			lpOpts.UseFinalityTag = true
			lp := NewLogPoller(orm, ec, lggr, lpOpts)

			latestBlock, lastFinalizedBlockNumber, err := lp.latestBlocks(ctx)
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

			lpOpts.UseFinalityTag = true
			lp := NewLogPoller(orm, ec, lggr, lpOpts)
			_, _, err := lp.latestBlocks(ctx)
			require.Error(t, err)
		})

		t.Run("BatchCall returns an error", func(t *testing.T) {
			ec := evmclimocks.NewClient(t)
			ec.On("BatchCallContext", mock.Anything, mock.Anything).Return(fmt.Errorf("some error"))
			lpOpts.UseFinalityTag = true
			lp := NewLogPoller(orm, ec, lggr, lpOpts)
			_, _, err := lp.latestBlocks(ctx)
			require.Error(t, err)
		})
	})
}

func Test_GetLogsReqHelpers(t *testing.T) {
	t.Parallel()

	filter1 := Filter{
		Name:      "filter1",
		Addresses: evmtypes.AddressArray{testutils.NewAddress()},
		EventSigs: evmtypes.HashArray{EmitterABI.Events["Log1"].ID, EmitterABI.Events["Log2"].ID},
		Topic2:    []common.Hash{EmitterABI.Events["Log3"].ID, EmitterABI.Events["Log4"].ID},
	}

	blockHash := common.HexToHash("0x1234")

	testCases := []struct {
		name      string
		filter    Filter
		fromBlock *big.Int
		toBlock   *big.Int
		blockHash *common.Hash
	}{
		{"block range", filter1, big.NewInt(5), big.NewInt(10), nil},
		{"block hash", filter1, nil, nil, &blockHash},
	}

	for _, c := range testCases {
		var req *GetLogsBatchElem
		t.Run("createGetLogsRec", func(t *testing.T) {
			req = NewGetLogsReq(filter1)
			assert.Equal(t, req.Method, "eth_getLogs")
			require.Len(t, req.Args, 1)
			assert.Equal(t, c.filter.Addresses, evmtypes.AddressArray(req.Addresses()))
			assert.Equal(t, c.filter.EventSigs, evmtypes.HashArray(req.Topics()[0]))
			assert.Equal(t, c.filter.Topic2, evmtypes.HashArray(req.Topics()[1]))
		})
		t.Run("mergeAddressesIntoGetLogsReq", func(t *testing.T) {
			newAddresses := []common.Address{testutils.NewAddress(), testutils.NewAddress()}
			mergeAddressesIntoGetLogsReq(req, newAddresses)
			assert.Len(t, req.Addresses(), len(c.filter.Addresses)+2)
			assert.Contains(t, req.Addresses(), newAddresses[0])
			assert.Contains(t, req.Addresses(), newAddresses[1])
		})
	}
}

func benchmarkEthGetLogsReqs(b *testing.B, nFilters, nAddresses, nEvents int) {
	lggr := logger.Test(b)
	chainID := testutils.NewRandomEVMChainID()
	db := pgtest.NewSqlxDB(b)

	orm := NewORM(chainID, db, lggr)
	lpOpts := Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := NewLogPoller(orm, nil, lggr, lpOpts)
	for i := 0; i < nFilters; i++ {
		var addresses []common.Address
		var events []common.Hash
		for j := 0; j < nAddresses; j++ {
			addresses = append(addresses, common.BigToAddress(big.NewInt(int64(j+1))))
		}
		for j := 0; j < nEvents; j++ {
			events = append(events, common.BigToHash(big.NewInt(int64(j+1))))
		}
		err := lp.RegisterFilter(testutils.Context(b), Filter{Name: "my Filter", EventSigs: events, Addresses: addresses})
		require.NoError(b, err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lp.EthGetLogsReqs(nil, nil, nil)
	}
}

func BenchmarkEthGetLogsReqs10_1(b *testing.B) {
	benchmarkEthGetLogsReqs(b, 10, 1, 1)
}
func BenchmarkEthGetLogsReqs100_10(b *testing.B) {
	benchmarkEthGetLogsReqs(b, 100, 10, 10)
}
func BenchmarkEthGetLogsReqs1000_100(b *testing.B) {
	benchmarkEthGetLogsReqs(b, 1000, 100, 100)
}
