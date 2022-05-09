package log_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestBroadcaster_AwaitsInitialSubscribersOnStartup(t *testing.T) {
	g := gomega.NewWithT(t)

	const blockHeight int64 = 123
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.lb.AddDependents(2)

	var listener = helper.newLogListenerWithJob("A")
	helper.register(listener, newMockContract(), 1)

	helper.start()
	defer helper.stop()

	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 0 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(0)))

	helper.lb.DependentReady()

	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 0 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(0)))

	helper.lb.DependentReady()

	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 1 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(1)))

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() == 1 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.UnsubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(1)))

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	const (
		numConfirmations            = 1
		numContracts                = 3
		blockHeight           int64 = 123
		lastStoredBlockHeight       = blockHeight - 25
	)

	backfillTimes := 2
	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: backfillTimes,
		HeaderByNumber:      backfillTimes,
		FilterLogs:          backfillTimes,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	blockBackfillDepth := helper.config.BlockBackfillDepth()

	var backfillCount atomic.Int64

	// the first backfill should use the height of last head saved to the db,
	// minus maxNumConfirmations of subscribers and minus blockBackfillDepth
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		backfillCount.Store(1)
		require.Equal(t, lastStoredBlockHeight-numConfirmations-int64(blockBackfillDepth), fromBlock)
	}

	listener := helper.newLogListenerWithJob("initial")

	helper.register(listener, newMockContract(), numConfirmations)

	for i := 0; i < numContracts; i++ {
		listener := helper.newLogListenerWithJob("")
		helper.register(listener, newMockContract(), 1)
	}

	helper.start()
	defer helper.stop()

	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 1 }, cltest.WaitTimeout(t), time.Second)
	gomega.NewWithT(t).Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(1)))
	gomega.NewWithT(t).Consistently(func() int32 { return helper.mockEth.UnsubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(0)))

	require.Eventually(t, func() bool { return backfillCount.Load() == 1 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	helper.unsubscribeAll()

	// now the backfill must use the blockBackfillDepth
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		require.Equal(t, blockHeight-int64(blockBackfillDepth), fromBlock)
		backfillCount.Store(2)
	}

	listenerLast := helper.newLogListenerWithJob("last")
	helper.register(listenerLast, newMockContract(), 1)

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() >= 1 }, cltest.WaitTimeout(t), time.Second)
	gomega.NewWithT(t).Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(2)))
	gomega.NewWithT(t).Consistently(func() int32 { return helper.mockEth.UnsubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(1)))

	require.Eventually(t, func() bool { return backfillCount.Load() == 2 }, cltest.WaitTimeout(t), time.Second)

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BackfillOnNodeStartAndOnReplay(t *testing.T) {
	const (
		lastStoredBlockHeight       = 100
		blockHeight           int64 = 125
		replayFrom            int64 = 40
	)

	backfillTimes := 2
	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: backfillTimes,
		HeaderByNumber:      backfillTimes,
		FilterLogs:          2,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	maxNumConfirmations := int64(10)

	var backfillCount atomic.Int64

	listener := helper.newLogListenerWithJob("one")
	helper.register(listener, newMockContract(), uint32(maxNumConfirmations))

	listener2 := helper.newLogListenerWithJob("two")
	helper.register(listener2, newMockContract(), uint32(2))

	blockBackfillDepth := helper.config.BlockBackfillDepth()

	// the first backfill should use the height of last head saved to the db,
	// minus maxNumConfirmations of subscribers and minus blockBackfillDepth
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		times := backfillCount.Inc() - 1
		if times == 0 {
			require.Equal(t, lastStoredBlockHeight-maxNumConfirmations-int64(blockBackfillDepth), fromBlock)
		} else if times == 1 {
			require.Equal(t, replayFrom, fromBlock)
		}
	}

	func() {
		helper.start()
		defer helper.stop()

		require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 1 }, cltest.WaitTimeout(t), time.Second)
		require.Eventually(t, func() bool { return backfillCount.Load() == 1 }, cltest.WaitTimeout(t), time.Second)

		helper.lb.ReplayFromBlock(replayFrom, false)

		require.Eventually(t, func() bool { return backfillCount.Load() >= 2 }, cltest.WaitTimeout(t), time.Second)
	}()

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() >= 1 }, cltest.WaitTimeout(t), time.Second)
	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_ReplaysLogs(t *testing.T) {
	const (
		blockHeight = 10
	)

	blocks := cltest.NewBlocks(t, blockHeight+3)
	contract, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)
	sentLogs := []types.Log{
		blocks.LogOnBlockNum(3, contract.Address()),
		blocks.LogOnBlockNum(7, contract.Address()),
	}

	mockEth := newMockEthClient(t, make(chan evmtest.RawSub[types.Log], 4), blockHeight, mockEthClientExpectedCalls{
		FilterLogs:       4,
		FilterLogsResult: sentLogs,
	})
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(blockHeight))
	helper.mockEth = mockEth

	listener := helper.newLogListenerWithJob("listener")
	helper.register(listener, contract, 2)

	func() {
		helper.start()
		defer helper.stop()

		// To start, no logs are sent
		require.Eventually(t, func() bool { return len(listener.getUniqueLogs()) == 0 }, cltest.WaitTimeout(t), time.Second,
			"expected unique logs to be 0 but was %d", len(listener.getUniqueLogs()))

		// Replay from block 2, the logs should be delivered. An incoming head must be simulated to
		// trigger log delivery.
		helper.lb.ReplayFromBlock(2, false)
		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     10,
			EndBlock:       10,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
			Blocks:         blocks,
		})
		require.Eventually(t, func() bool { return len(listener.getUniqueLogs()) == 2 }, cltest.WaitTimeout(t), time.Second,
			"expected unique logs to be 2 but was %d", len(listener.getUniqueLogs()))

		// Replay again, the logs are already marked consumed so they should not be included in
		// getUniqueLogs.
		helper.lb.ReplayFromBlock(2, false)
		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     11,
			EndBlock:       11,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
			Blocks:         blocks,
		})
		require.Eventually(t, func() bool { return len(listener.getUniqueLogs()) == 2 }, cltest.WaitTimeout(t), time.Second,
			"expected unique logs to be 2 but was %d", len(listener.getUniqueLogs()))

		// Replay again with forceBroadcast. The logs are consumed again.
		helper.lb.ReplayFromBlock(2, true)
		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     12,
			EndBlock:       12,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
			Blocks:         blocks,
		})
		require.Eventually(t, func() bool { return len(listener.getUniqueLogs()) == 4 }, cltest.WaitTimeout(t), time.Second,
			"expected unique logs to be 4 but was %d", len(listener.getUniqueLogs()))

	}()

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() >= 1 }, cltest.WaitTimeout(t), time.Second)
}

func TestBroadcaster_BackfillUnconsumedAfterCrash(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)

	orm := log.NewORM(db, lggr, cfg, cltest.FixtureChainID)

	helperCfg := broadcasterHelperCfg{db: db}
	contract1 := newMockContract()
	contract2 := newMockContract()

	blocks := cltest.NewBlocks(t, 10)
	const (
		log1Block = 1
		log2Block = 4

		confs = 2
	)
	log1 := blocks.LogOnBlockNum(log1Block, contract1.Address())
	log2 := blocks.LogOnBlockNum(log2Block, contract2.Address())
	logs := []types.Log{log1, log2}

	contract1.On("ParseLog", log1).Return(flux_aggregator_wrapper.FluxAggregatorNewRound{}, nil)
	contract2.On("ParseLog", log2).Return(flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}, nil)

	// Pool two logs from subscription, then shut down
	helper := helperCfg.new(t, 0, 2, logs)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(confs)
	listener := helper.newLogListenerWithJob("one")
	listener.SkipMarkingConsumed(true)
	listener2 := helper.newLogListenerWithJob("two")
	listener2.SkipMarkingConsumed(true)
	func() {
		helper.lb.AddDependents(2)
		helper.start()
		defer helper.stop()
		helper.register(listener, contract1, confs)
		helper.register(listener2, contract2, confs)
		helper.lb.DependentReady()
		helper.lb.DependentReady()

		headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     0,
			EndBlock:       1,
			Blocks:         blocks,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		})

		chRawLogs := <-helper.chchRawLogs
		chRawLogs.TrySend(log1)
		chRawLogs.TrySend(log2)

		<-headsDone

		require.Eventually(t, func() bool {
			blockNum, err := orm.GetPendingMinBlock()
			return assert.NoError(t, err) && blockNum != nil && *blockNum == int64(log1.BlockNumber)
		}, cltest.WaitTimeout(t), time.Second)
	}()

	// Pool min block in DB and neither listener received a broadcast
	blockNum, err := orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.NotNil(t, blockNum)
	require.Equal(t, int64(log1.BlockNumber), *blockNum)
	require.Empty(t, listener.getUniqueLogs())
	require.Empty(t, listener2.getUniqueLogs())
	helper.requireBroadcastCount(0)

	// Backfill pool with both, then broadcast one, but don't consume
	helper = helperCfg.new(t, 2, 2, logs)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(confs)
	listener = helper.newLogListenerWithJob("one")
	listener.SkipMarkingConsumed(true)
	listener2 = helper.newLogListenerWithJob("two")
	listener2.SkipMarkingConsumed(true)
	func() {
		helper.lb.AddDependents(2)
		helper.start()
		defer helper.stop()
		helper.register(listener, contract1, confs)
		helper.register(listener2, contract2, confs)
		helper.lb.DependentReady()
		helper.lb.DependentReady()

		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     2,
			EndBlock:       4,
			Blocks:         blocks,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		})

		require.Eventually(t, func() bool {
			blockNum, err := orm.GetPendingMinBlock()
			return assert.NoError(t, err) && blockNum != nil && *blockNum == int64(log2.BlockNumber)
		}, cltest.WaitTimeout(t), time.Second)
	}()

	// Pool min block in DB and one listener received but didn't consume
	blockNum, err = orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.NotNil(t, blockNum)
	require.Equal(t, int64(log2.BlockNumber), *blockNum)
	require.NotEmpty(t, listener.getUniqueLogs())
	require.Empty(t, listener2.getUniqueLogs())
	c, err := orm.WasBroadcastConsumed(log1.BlockHash, log1.Index, listener.JobID())
	require.NoError(t, err)
	require.False(t, c)

	// Backfill pool and broadcast two, but only consume one
	helper = helperCfg.new(t, 4, 2, logs)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(confs)
	listener = helper.newLogListenerWithJob("one")
	listener2 = helper.newLogListenerWithJob("two")
	listener2.SkipMarkingConsumed(true)
	func() {
		helper.lb.AddDependents(2)
		helper.start()
		defer helper.stop()
		helper.register(listener, contract1, confs)
		helper.register(listener2, contract2, confs)
		helper.lb.DependentReady()
		helper.lb.DependentReady()

		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     5,
			EndBlock:       7,
			Blocks:         blocks,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		})

		require.Eventually(t, func() bool {
			blockNum, err := orm.GetPendingMinBlock()
			return assert.NoError(t, err) && blockNum == nil
		}, cltest.WaitTimeout(t), time.Second)
	}()

	// Pool empty and one consumed but other didn't
	blockNum, err = orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.Nil(t, blockNum)
	require.NotEmpty(t, listener.getUniqueLogs())
	require.NotEmpty(t, listener2.getUniqueLogs())
	c, err = orm.WasBroadcastConsumed(log1.BlockHash, log1.Index, listener.JobID())
	require.NoError(t, err)
	require.True(t, c)
	c, err = orm.WasBroadcastConsumed(log2.BlockHash, log2.Index, listener2.JobID())
	require.NoError(t, err)
	require.False(t, c)

	// Backfill pool, broadcast and consume one
	helper = helperCfg.new(t, 7, 2, logs[1:])
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(confs)
	listener = helper.newLogListenerWithJob("one")
	listener2 = helper.newLogListenerWithJob("two")
	func() {
		helper.lb.AddDependents(2)
		helper.start()
		defer helper.stop()
		helper.register(listener, contract1, confs)
		helper.register(listener2, contract2, confs)
		helper.lb.DependentReady()
		helper.lb.DependentReady()

		<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     8,
			EndBlock:       8,
			Blocks:         blocks,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		})

		require.Eventually(t, func() bool {
			blockNum, err := orm.GetPendingMinBlock()
			return assert.NoError(t, err) && blockNum == nil
		}, cltest.WaitTimeout(t), time.Second)
	}()
	// Pool empty, one broadcasted and consumed
	blockNum, err = orm.GetPendingMinBlock()
	require.NoError(t, err)
	require.Nil(t, blockNum)
	require.Empty(t, listener.getUniqueLogs())
	require.NotEmpty(t, listener2.getUniqueLogs())
	c, err = orm.WasBroadcastConsumed(log2.BlockHash, log2.Index, listener2.JobID())
	require.NoError(t, err)
	require.True(t, c)
}

func TestBroadcaster_ShallowBackfillOnNodeStart(t *testing.T) {
	const (
		lastStoredBlockHeight       = 100
		blockHeight           int64 = 125
	)

	backfillTimes := 1
	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: backfillTimes,
		HeaderByNumber:      backfillTimes,
		FilterLogs:          backfillTimes,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	backfillDepth := 15

	helper.globalConfig.Overrides.BlockBackfillSkip = null.BoolFrom(true)
	helper.globalConfig.Overrides.BlockBackfillDepth = null.IntFrom(int64(backfillDepth))

	var backfillCount atomic.Int64

	listener := helper.newLogListenerWithJob("one")
	helper.register(listener, newMockContract(), uint32(10))

	listener2 := helper.newLogListenerWithJob("two")
	helper.register(listener2, newMockContract(), uint32(2))

	// the backfill does not use the height from DB because BlockBackfillSkip is true
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		backfillCount.Store(1)
		require.Equal(t, blockHeight-int64(backfillDepth), fromBlock)
	}

	func() {
		helper.start()
		defer helper.stop()

		require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 1 }, cltest.WaitTimeout(t), time.Second)
		require.Eventually(t, func() bool { return backfillCount.Load() == 1 }, cltest.WaitTimeout(t), time.Second)
	}()

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() >= 1 }, cltest.WaitTimeout(t), time.Second)
	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BackfillInBatches(t *testing.T) {
	const (
		numConfirmations            = 1
		blockHeight           int64 = 120
		lastStoredBlockHeight       = blockHeight - 29
		backfillTimes               = 1
		batchSize             int64 = 5
		expectedBatches             = 9
	)

	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: backfillTimes,
		HeaderByNumber:      backfillTimes,
		FilterLogs:          expectedBatches,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	blockBackfillDepth := helper.config.BlockBackfillDepth()
	helper.globalConfig.Overrides.GlobalEvmLogBackfillBatchSize = null.IntFrom(batchSize)

	var backfillCount atomic.Int64

	lggr := logger.TestLogger(t)
	backfillStart := lastStoredBlockHeight - numConfirmations - int64(blockBackfillDepth)
	// the first backfill should start from before the last stored head
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		times := backfillCount.Inc() - 1
		lggr.Infof("Log Batch: --------- times %v - %v, %v", times, fromBlock, toBlock)

		if times <= 7 {
			require.Equal(t, backfillStart+batchSize*times, fromBlock)
			require.Equal(t, backfillStart+batchSize*(times+1)-1, toBlock)
		} else {
			// last batch is for a range of 1
			require.Equal(t, int64(120), fromBlock)
			require.Equal(t, int64(120), toBlock)
		}
	}

	listener := helper.newLogListenerWithJob("initial")
	helper.register(listener, newMockContract(), numConfirmations)
	helper.start()

	defer helper.stop()

	require.Eventually(t, func() bool { return backfillCount.Load() == expectedBatches }, cltest.WaitTimeout(t), time.Second)

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.UnsubscribeCallCount() >= 1 }, cltest.WaitTimeout(t), time.Second)

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BackfillALargeNumberOfLogs(t *testing.T) {
	g := gomega.NewWithT(t)
	const (
		lastStoredBlockHeight int64 = 10

		// a large number of blocks since lastStoredBlockHeight
		blockHeight int64 = 3000

		backfillTimes          = 1
		batchSize       uint32 = 50
		expectedBatches        = 61
	)

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 7)
	backfilledLogs := make([]types.Log, 0)
	for i := 0; i < 50; i++ {
		aLog := blocks.LogOnBlockNum(0, contract1.Address())
		backfilledLogs = append(backfilledLogs, aLog)
	}

	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: backfillTimes,
		HeaderByNumber:      backfillTimes,
		FilterLogs:          expectedBatches,

		FilterLogsResult: backfilledLogs,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	helper.globalConfig.Overrides.GlobalEvmLogBackfillBatchSize = null.IntFrom(int64(batchSize))

	var backfillCount atomic.Int64

	lggr := logger.TestLogger(t)
	mockEth.CheckFilterLogs = func(fromBlock int64, toBlock int64) {
		times := backfillCount.Inc() - 1
		lggr.Warnf("Log Batch: --------- times %v - %v, %v", times, fromBlock, toBlock)
	}

	listener := helper.newLogListenerWithJob("initial")
	helper.register(listener, newMockContract(), 1)
	helper.start()
	defer helper.stop()
	g.Eventually(func() int64 { return backfillCount.Load() }, cltest.WaitTimeout(t), time.Second).Should(gomega.Equal(int64(expectedBatches)))

	helper.unsubscribeAll()
	g.Eventually(func() int32 { return helper.mockEth.UnsubscribeCallCount() }, cltest.WaitTimeout(t), time.Second).Should(gomega.BeNumerically(">=", int32(1)))

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)
	contract2, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 10)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}
	addr2SentLogs := []types.Log{
		blocks.LogOnBlockNum(4, contract2.Address()),
		blocks.LogOnBlockNum(5, contract2.Address()),
		blocks.LogOnBlockNum(6, contract2.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	listener2 := helper.newLogListenerWithJob("listener 2")
	listener3 := helper.newLogListenerWithJob("listener 3")
	listener4 := helper.newLogListenerWithJob("listener 4")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 1)
	helper.register(listener3, contract2, 1)
	helper.register(listener4, contract2, 1)

	func() {
		helper.start()
		defer helper.stop()

		headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
			StartBlock:     0,
			EndBlock:       9,
			HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
			Blocks:         blocks,
		})

		defer helper.unsubscribeAll()

		chRawLogs := <-helper.chchRawLogs

		for _, log := range addr1SentLogs {
			chRawLogs.TrySend(log)
		}
		for _, log := range addr2SentLogs {
			chRawLogs.TrySend(log)
		}

		<-headsDone
		helper.requireBroadcastCount(12)

		requireEqualLogs(t, addr1SentLogs, listener1.received.getUniqueLogs())
		requireEqualLogs(t, addr1SentLogs, listener2.received.getUniqueLogs())

		requireEqualLogs(t, addr2SentLogs, listener3.received.getUniqueLogs())
		requireEqualLogs(t, addr2SentLogs, listener4.received.getUniqueLogs())
	}()
	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeights(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 10)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	listener2 := helper.newLogListenerWithJob("listener 2")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 8)

	_ = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       9,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	helper.requireBroadcastCount(5)
	helper.stop()

	require.Equal(t, []uint64{1, 2, 3}, listener1.getUniqueLogsBlockNumbers())
	require.Equal(t, []uint64{1, 2}, listener2.getUniqueLogsBlockNumbers())

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.getUniqueLogs(),
	)
	requireEqualLogs(t,
		[]types.Log{
			addr1SentLogs[0],
			addr1SentLogs[1],
		},
		listener2.received.getUniqueLogs(),
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.getUniqueLogs(),
		listener1.received.getLogs(),
	)
	requireEqualLogs(t,
		listener2.received.getUniqueLogs(),
		listener2.received.getLogs(),
	)

	// the logs should have been received at much later heights
	logsOnBlocks := listener2.received.logsOnBlocks()
	expectedLogsOnBlocks := []logOnBlock{
		{
			logBlockNumber: 1,
			blockNumber:    8,
			blockHash:      blocks.Hashes[8],
		},
		{
			logBlockNumber: 2,
			blockNumber:    9,
			blockHash:      blocks.Hashes[9],
		},
	}

	assert.Equal(t, len(logsOnBlocks), len(expectedLogsOnBlocks))
	require.Equal(t, logsOnBlocks, expectedLogsOnBlocks)

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_DeletesOldLogsAfterNumberOfHeads(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(1)
	helper.start()
	defer helper.stop()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	listener2 := helper.newLogListenerWithJob("listener 2")
	listener3 := helper.newLogListenerWithJob("listener 3")
	listener4 := helper.newLogListenerWithJob("listener 4")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 3)

	headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	helper.requireBroadcastCount(6)
	<-headsDone

	helper.register(listener3, contract1, 1)
	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     6,
		EndBlock:       8,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	// the new listener should still receive 2 of the 3 logs
	helper.requireBroadcastCount(8)
	require.Equal(t, 2, len(listener3.received.getUniqueLogs()))

	helper.register(listener4, contract1, 1)
	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     9,
		EndBlock:       11,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	// but this one should receive none
	require.Equal(t, 0, len(listener4.received.getUniqueLogs()))
}

func TestBroadcaster_DeletesOldLogsOnlyAfterFinalityDepth(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(4)
	helper.start()
	defer helper.stop()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	listener2 := helper.newLogListenerWithJob("listener 2")
	listener3 := helper.newLogListenerWithJob("listener 3")
	listener4 := helper.newLogListenerWithJob("listener 4")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 3)

	headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	<-headsDone
	helper.requireBroadcastCount(6)

	helper.register(listener3, contract1, 1)
	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     7,
		EndBlock:       8,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	// the new listener should still receive 3 logs because of finality depth being higher than max NumConfirmations
	helper.requireBroadcastCount(9)
	require.Equal(t, 3, len(listener3.received.getUniqueLogs()))

	helper.register(listener4, contract1, 1)
	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     10,
		EndBlock:       11,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	// but this one should receive none
	require.Equal(t, 0, len(listener4.received.getUniqueLogs()))
}

func TestBroadcaster_FilterByTopicValues(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(3)
	helper.start()
	defer helper.stop()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)

	topic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	field1Value1 := utils.NewHash()
	field1Value2 := utils.NewHash()
	field2Value1 := utils.NewHash()
	field2Value2 := utils.NewHash()
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNumWithTopics(1, 0, contract1.Address(), []common.Hash{topic, field1Value1, field2Value1}),
		blocks.LogOnBlockNumWithTopics(1, 1, contract1.Address(), []common.Hash{topic, field1Value2, field2Value2}),
		blocks.LogOnBlockNumWithTopics(2, 0, contract1.Address(), []common.Hash{topic, utils.NewHash(), field2Value2}),
		blocks.LogOnBlockNumWithTopics(2, 1, contract1.Address(), []common.Hash{topic, field1Value2, utils.NewHash()}),
	}

	listener0 := helper.newLogListenerWithJob("listener 0")
	listener1 := helper.newLogListenerWithJob("listener 1")
	listener2 := helper.newLogListenerWithJob("listener 2")
	listener3 := helper.newLogListenerWithJob("listener 3")
	listener4 := helper.newLogListenerWithJob("listener 4")

	helper.registerWithTopicValues(listener0, contract1, 1,
		map[common.Hash][][]log.Topic{
			topic: {}, // no filters, so all values allowed
		},
	)
	helper.registerWithTopicValues(listener1, contract1, 1,
		map[common.Hash][][]log.Topic{
			topic: {{} /**/, {}}, // two empty filters, so all values allowed
		},
	)
	helper.registerWithTopicValues(listener2, contract1, 1,
		map[common.Hash][][]log.Topic{
			topic: {
				{log.Topic(field1Value1), log.Topic(field1Value2)} /**/, {log.Topic(field2Value1), log.Topic(field2Value2)}, // two values for each field allowed
			},
		},
	)
	helper.registerWithTopicValues(listener3, contract1, 1,
		map[common.Hash][][]log.Topic{
			topic: {
				{log.Topic(field1Value1), log.Topic(field1Value2)} /**/, {}, // two values allowed for field 1, and any values for field 2
			},
		},
	)
	helper.registerWithTopicValues(listener4, contract1, 1,
		map[common.Hash][][]log.Topic{
			topic: {
				{log.Topic(field1Value1)} /**/, {log.Topic(field2Value1)}, // some values allowed
			},
		},
	)

	headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	<-headsDone

	require.Eventually(t, func() bool { return len(listener0.received.getUniqueLogs()) == 4 }, cltest.WaitTimeout(t), 500*time.Millisecond)
	require.Eventually(t, func() bool { return len(listener1.received.getUniqueLogs()) == 4 }, cltest.WaitTimeout(t), 500*time.Millisecond)
	require.Eventually(t, func() bool { return len(listener2.received.getUniqueLogs()) == 2 }, cltest.WaitTimeout(t), 500*time.Millisecond)
	require.Eventually(t, func() bool { return len(listener3.received.getUniqueLogs()) == 3 }, cltest.WaitTimeout(t), 500*time.Millisecond)
	require.Eventually(t, func() bool { return len(listener4.received.getUniqueLogs()) == 1 }, cltest.WaitTimeout(t), 500*time.Millisecond)
}

func TestBroadcaster_BroadcastsWithOneDelayedLog(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(2)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 12)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),

		// this log will arrive after head with block number 3 and a previous log for it were already processed
		blocks.LogOnBlockNumWithIndex(3, 1, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	chRawLogs.TrySend(addr1SentLogs[0])
	chRawLogs.TrySend(addr1SentLogs[1])
	chRawLogs.TrySend(addr1SentLogs[2])

	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       3,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs.TrySend(addr1SentLogs[3])

	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     4,
		EndBlock:       8,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	helper.requireBroadcastCount(4)
	helper.stop()

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeightsWithLogsEarlierThanHeads(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 10)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       9,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	helper.requireBroadcastCount(3)
	helper.stop()

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.getUniqueLogs(),
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.getUniqueLogs(),
		listener1.received.getLogs(),
	)

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeightsWithHeadsEarlierThanLogs(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(2)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(testutils.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 12)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListenerWithJob("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       6,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	for _, log := range addr1SentLogs {
		chRawLogs.TrySend(log)
	}

	<-cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     7,
		EndBlock:       8,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	helper.requireBroadcastCount(3)
	helper.stop()

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.getUniqueLogs(),
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.getUniqueLogs(),
		listener1.received.getLogs(),
	)

	helper.mockEth.AssertExpectations(t)
}

func TestBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	const (
		backfillTimes = 1
		blockHeight   = 15
		expectedBlock = 5
	)
	var (
		ethClient = new(evmmocks.Client)
		contract0 = newMockContract()
		contract1 = newMockContract()
		contract2 = newMockContract()
	)
	ethClient.Test(t)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	chchRawLogs := make(chan evmtest.RawSub[types.Log], backfillTimes)
	chStarted := make(chan struct{})
	ethClient.On("ChainID", mock.Anything).Return(&cltest.FixtureChainID)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) ethereum.Subscription {
				defer close(chStarted)
				sub := mockEth.NewSub(t)
				chchRawLogs <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) error {
				return nil
			},
		).
		Once()

	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchRawLogs <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) error {
				return nil
			},
		).
		Times(3)

	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&evmtypes.Head{Number: blockHeight}, nil)

	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, contract0.Address())
			require.Len(t, query.Addresses, 1)
		}).
		Return(nil, nil).
		Times(backfillTimes)

	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, contract0.Address())
			require.Contains(t, query.Addresses, contract1.Address())
			require.Len(t, query.Addresses, 2)
		}).
		Return(nil, nil).
		Once()

	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, contract0.Address())
			require.Contains(t, query.Addresses, contract1.Address())
			require.Contains(t, query.Addresses, contract2.Address())
			require.Len(t, query.Addresses, 3)
		}).
		Return(nil, nil).
		Once()

	helper := newBroadcasterHelperWithEthClient(t, ethClient, nil)
	helper.lb.AddDependents(1)
	helper.start()
	defer helper.stop()

	listener0 := helper.newLogListenerWithJob("0")
	listener1 := helper.newLogListenerWithJob("1")
	listener2 := helper.newLogListenerWithJob("2")

	// Subscribe #0
	helper.register(listener0, contract0, 1)
	defer helper.unsubscribeAll()
	helper.lb.DependentReady()

	// Await startup
	select {
	case <-chStarted:
	case <-time.After(cltest.WaitTimeout(t)):
		t.Fatal("never started")
	}

	select {
	case <-chchRawLogs:
	case <-time.After(cltest.WaitTimeout(t)):
		t.Fatal("did not subscribe")
	}

	// Subscribe #1
	helper.register(listener1, contract1, 1)

	select {
	case <-chchRawLogs:
	case <-time.After(cltest.WaitTimeout(t)):
		t.Fatal("did not subscribe")
	}

	// Subscribe #2
	helper.register(listener2, contract2, 1)

	select {
	case <-chchRawLogs:
	case <-time.After(cltest.WaitTimeout(t)):
		t.Fatal("did not subscribe")
	}

	// ReplayFrom will not lead to backfill because the number is above current height
	helper.lb.ReplayFromBlock(125, false)

	select {
	case <-chchRawLogs:
	case <-time.After(cltest.WaitTimeout(t)):
		t.Fatal("did not subscribe")
	}

	cltest.EventuallyExpectationsMet(t, ethClient, cltest.WaitTimeout(t), time.Second)
}

func TestBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	addrA := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	addrB := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	blocks := cltest.NewBlocks(t, 20)

	logsA := make(map[uint]types.Log)
	logsB := make(map[uint]types.Log)
	for n := 1; n < 18; n++ {
		logsA[uint(n)] = blocks.LogOnBlockNumWithIndex(uint64(n), 0, addrA)
		logsB[uint(n)] = blocks.LogOnBlockNumWithIndex(uint64(n), 1, addrB)
	}

	tests := []struct {
		name              string
		blockHeight1      int64
		blockHeight2      int64
		batch1            []uint
		backfillableLogs  []uint
		batch2            []uint
		expectedFilteredA []uint
		expectedFilteredB []uint
	}{
		{
			name: "no backfilled logs, no overlap",

			blockHeight1: 0,
			batch1:       []uint{1, 2},

			blockHeight2:     3,
			backfillableLogs: nil,
			batch2:           []uint{7, 8},

			expectedFilteredA: []uint{1, 2, 7, 8},
			expectedFilteredB: []uint{7, 8},
		},
		{
			name: "no backfilled logs, overlap",

			blockHeight1: 0,
			batch1:       []uint{1, 2},

			blockHeight2:     2,
			backfillableLogs: nil,
			batch2:           []uint{2, 3},

			expectedFilteredA: []uint{1, 2, 3},
			expectedFilteredB: []uint{2, 3},
		},
		{
			name: "backfilled logs, no overlap",

			blockHeight1: 0,
			batch1:       []uint{1, 2},

			blockHeight2:     15,
			backfillableLogs: []uint{11, 12, 15},
			batch2:           []uint{16, 17},

			expectedFilteredA: []uint{1, 2, 11, 12, 15, 16, 17},
			expectedFilteredB: []uint{11, 12, 15, 16, 17},
		},
		{
			name: "backfilled logs, overlap",

			blockHeight1: 0,
			batch1:       []uint{1, 11},

			blockHeight2:     15,
			backfillableLogs: []uint{11, 12, 15},
			batch2:           []uint{16, 17},

			expectedFilteredA: []uint{1, 11, 12, 15, 16, 17},
			expectedFilteredB: []uint{11, 12, 15, 16, 17},
		},
	}

	batchContains := func(batch []uint, n uint) bool {
		for _, x := range batch {
			if x == n {
				return true
			}
		}
		return false
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			helper := newBroadcasterHelper(t, test.blockHeight1, 2)
			var backfillDepth int64 = 5
			// something other than default
			helper.globalConfig.Overrides.BlockBackfillDepth = null.IntFrom(int64(backfillDepth))

			helper.start()
			defer helper.stop()

			logListenerA := helper.newLogListenerWithJob("logListenerA")
			logListenerB := helper.newLogListenerWithJob("logListenerB")

			contractA, err := flux_aggregator_wrapper.NewFluxAggregator(addrA, nil)
			require.NoError(t, err)
			contractB, err := flux_aggregator_wrapper.NewFluxAggregator(addrB, nil)
			require.NoError(t, err)

			// Register listener A
			helper.register(logListenerA, contractA, 1)

			lggr := logger.TestLogger(t)
			// Send initial logs
			chRawLogs1 := <-helper.chchRawLogs
			headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock: test.blockHeight1,
				EndBlock:   test.blockHeight2 + 1,
				Blocks:     blocks,
				HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head *evmtypes.Head) {
					lggr.Infof("------------ HEAD TRACKABLE (%v) --------------", head.Number)
					if _, exists := logsA[uint(head.Number)]; !exists {
						lggr.Warnf("  ** not exists")
						return
					} else if !batchContains(test.batch1, uint(head.Number)) {
						lggr.Warnf("  ** not batchContains %v %v", head.Number, test.batch1)
						return
					}
					lggr.Warnf("  ** yup!")
					chRawLogs1.TrySend(logsA[uint(head.Number)])
				})},
			})

			helper.requireBroadcastCount(len(test.batch1))
			expectedA := newReceived(pickLogs(logsA, test.batch1))
			logListenerA.requireAllReceived(t, expectedA)

			<-headsDone
			helper.mockEth.EthClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&evmtypes.Head{Number: test.blockHeight2}, nil).Once()

			combinedLogs := append(pickLogs(logsA, test.backfillableLogs), pickLogs(logsB, test.backfillableLogs)...)
			call := helper.mockEth.EthClient.On("FilterLogs", mock.Anything, mock.Anything).Return(combinedLogs, nil).Once()
			call.Run(func(args mock.Arguments) {
				// Validate that the ethereum.FilterQuery is specified correctly for the backfill that we expect
				fromBlock := args.Get(1).(ethereum.FilterQuery).FromBlock
				expected := big.NewInt(0)

				blockNumber := helper.lb.BackfillBlockNumber()
				if blockNumber.Valid && blockNumber.Int64 > test.blockHeight2-backfillDepth {
					expected = big.NewInt(blockNumber.Int64)
				} else if test.blockHeight2 > backfillDepth {
					expected = big.NewInt(test.blockHeight2 - backfillDepth)
				}
				require.Equal(t, expected, fromBlock)
			})

			// Register listener B (triggers re-subscription)
			helper.register(logListenerB, contractB, 1)

			// Send second batch of new logs
			chRawLogs2 := <-helper.chchRawLogs
			headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock: test.blockHeight2,
				Blocks:     blocks,
				HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head *evmtypes.Head) {
					if _, exists := logsA[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						chRawLogs2.TrySend(logsA[uint(head.Number)])
					}
					if _, exists := logsB[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						chRawLogs2.TrySend(logsB[uint(head.Number)])
					}
				})},
			})
			defer func() { <-headsDone }()

			expectedA = newReceived(pickLogs(logsA, test.expectedFilteredA))
			expectedB := newReceived(pickLogs(logsB, test.expectedFilteredB))
			logListenerA.requireAllReceived(t, expectedA)
			logListenerB.requireAllReceived(t, expectedB)
			helper.requireBroadcastCount(len(test.expectedFilteredA) + len(test.expectedFilteredB))

			helper.mockEth.EthClient.AssertExpectations(t)
		})
	}
}

func TestBroadcaster_AppendLogChannel(t *testing.T) {
	logs1 := []types.Log{
		{BlockNumber: 1},
		{BlockNumber: 2},
		{BlockNumber: 3},
		{BlockNumber: 4},
		{BlockNumber: 5},
	}

	logs2 := []types.Log{
		{BlockNumber: 6},
		{BlockNumber: 7},
		{BlockNumber: 8},
		{BlockNumber: 9},
		{BlockNumber: 10},
	}

	logs3 := []types.Log{
		{BlockNumber: 11},
		{BlockNumber: 12},
		{BlockNumber: 13},
		{BlockNumber: 14},
		{BlockNumber: 15},
	}

	ch1 := make(chan types.Log)
	ch2 := make(chan types.Log)
	ch3 := make(chan types.Log)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	lb := log.NewBroadcaster(nil, ethClient, nil, logger.TestLogger(t), nil)
	chCombined := lb.ExportedAppendLogChannel(ch1, ch2)
	chCombined = lb.ExportedAppendLogChannel(chCombined, ch3)

	go func() {
		defer close(ch1)
		for _, log := range logs1 {
			ch1 <- log
		}
	}()
	go func() {
		defer close(ch2)
		for _, log := range logs2 {
			ch2 <- log
		}
	}()
	go func() {
		defer close(ch3)
		for _, log := range logs3 {
			ch3 <- log
		}
	}()

	expected := append(logs1, logs2...)
	expected = append(expected, logs3...)

	var i int
	for log := range chCombined {
		require.Equal(t, expected[i], log)
		i++
	}
}

func TestBroadcaster_InjectsBroadcastRecordFunctions(t *testing.T) {
	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()
	defer helper.stop()

	blocks := cltest.NewBlocks(t, 20)

	logListener := helper.newLogListenerWithJob("logListener")

	contract := newMockContract()
	log1, log2 := blocks.LogOnBlockNum(0, contract.Address()), blocks.LogOnBlockNum(1, contract.Address())
	contract.On("ParseLog", log1).Return(flux_aggregator_wrapper.FluxAggregatorNewRound{}, nil)
	contract.On("ParseLog", log2).Return(flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}, nil)

	helper.register(logListener, contract, uint32(5))

	headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     3,
		EndBlock:       20,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})

	chRawLogs := <-helper.chchRawLogs

	chRawLogs.TrySend(log1)
	chRawLogs.TrySend(log2)

	<-headsDone
	require.Eventually(t, func() bool { return len(logListener.received.getUniqueLogs()) >= 2 }, cltest.WaitTimeout(t), time.Second)
	helper.requireBroadcastCount(2)

	helper.mockEth.EthClient.AssertExpectations(t)
}

func TestBroadcaster_ProcessesLogsFromReorgsAndMissedHead(t *testing.T) {
	g := gomega.NewWithT(t)

	const startBlockHeight int64 = 0
	helper := newBroadcasterHelper(t, startBlockHeight, 1)
	helper.start()
	defer helper.stop()

	blocks := cltest.NewBlocks(t, 10)
	blocksForked := blocks.ForkAt(t, 1, 5)

	var (
		addr = testutils.NewAddress()

		log0        = blocks.LogOnBlockNum(0, addr)
		log1        = blocks.LogOnBlockNum(1, addr)
		log2        = blocks.LogOnBlockNum(2, addr)
		log1Removed = blocks.LogOnBlockNumRemoved(1, addr)
		log2Removed = blocks.LogOnBlockNumRemoved(2, addr)
		log1R       = blocksForked.LogOnBlockNum(1, addr)
		log2R       = blocksForked.LogOnBlockNum(2, addr)
		log3R1      = blocksForked.LogOnBlockNumWithIndex(3, 0, addr)
		log3R2      = blocksForked.LogOnBlockNumWithIndex(3, 1, addr) // second log on the same block

		log1RRemoved  = blocksForked.LogOnBlockNumRemoved(1, addr)
		log2RRemoved  = blocksForked.LogOnBlockNumRemoved(2, addr)
		log3R1Removed = blocksForked.LogOnBlockNumWithIndexRemoved(3, 0, addr)
		log3R2Removed = blocksForked.LogOnBlockNumWithIndexRemoved(3, 1, addr)

		events = []interface{}{
			blocks.Head(0), log0,
			log1, // head1 missing
			blocks.Head(2), log2,
			blocks.Head(3),
			blocksForked.Head(1), log1Removed, log2Removed, log1R,
			blocksForked.Head(2), log2R,
			log3R1, blocksForked.Head(3), log3R2,
			blocksForked.Head(4),
			log1RRemoved, log0, log1, blocks.Head(4), log2, log2RRemoved, log3R1Removed, log3R2Removed, // a reorg back to the previous chain
			blocks.Head(5),
			blocks.Head(6),
			blocks.Head(7),
		}

		expectedA = []types.Log{log0, log1, log2, log1R, log2R, log3R1, log3R2}

		// listenerB needs 3 confirmations, so log2 is not sent to after the first reorg,
		// but is later - after the second reorg (back to the previous chain)
		expectedB = []types.Log{log0, log1, log1R, log2R, log2}
	)

	contract, err := flux_aggregator_wrapper.NewFluxAggregator(addr, nil)
	require.NoError(t, err)

	listenerA := helper.newLogListenerWithJob("listenerA")
	listenerB := helper.newLogListenerWithJob("listenerB")
	helper.register(listenerA, contract, 1)
	helper.register(listenerB, contract, 3)

	chRawLogs := <-helper.chchRawLogs

	for _, event := range events {
		switch x := event.(type) {
		case *evmtypes.Head:
			ctx, _ := context.WithTimeout(context.Background(), cltest.WaitTimeout(t))
			(helper.lb).(httypes.HeadTrackable).OnNewLongestChain(ctx, x)
		case types.Log:
			chRawLogs.TrySend(x)
		}
		time.Sleep(250 * time.Millisecond)
	}

	g.Eventually(func() []uint64 { return listenerA.getUniqueLogsBlockNumbers() }, cltest.WaitTimeout(t), time.Second).
		Should(gomega.Equal([]uint64{0, 1, 2, 1, 2, 3, 3}))
	g.Eventually(func() []uint64 { return listenerB.getUniqueLogsBlockNumbers() }, cltest.WaitTimeout(t), time.Second).
		Should(gomega.Equal([]uint64{0, 1, 1, 2, 2}))

	helper.unsubscribeAll()

	require.Equal(t, expectedA, listenerA.getUniqueLogs())
	require.Equal(t, expectedB, listenerB.getUniqueLogs())

	helper.mockEth.EthClient.AssertExpectations(t)
}

func TestBroadcaster_BackfillsForNewListeners(t *testing.T) {
	g := gomega.NewWithT(t)

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 2)
	helper.mockEth.EthClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&evmtypes.Head{Number: blockHeight}, nil).Times(2)
	helper.mockEth.EthClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil).Times(2)

	helper.start()
	defer helper.stop()

	addr1 := testutils.NewAddress()
	contract, err := flux_aggregator_wrapper.NewFluxAggregator(addr1, nil)
	require.NoError(t, err)

	listener1 := helper.newLogListenerWithJob("1")
	listener2 := helper.newLogListenerWithJob("2")

	topics1 := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener1, contract, topics1, 1)
	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 1 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(1)))

	<-helper.chchRawLogs

	topics2 := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
	}
	helper.registerWithTopics(listener2, contract, topics2, 1)
	require.Eventually(t, func() bool { return helper.mockEth.SubscribeCallCount() == 2 }, cltest.WaitTimeout(t), 100*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.SubscribeCallCount() }, 1*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(int32(2)))

	helper.unsubscribeAll()
}

func pickLogs(allLogs map[uint]types.Log, indices []uint) []types.Log {
	var picked []types.Log
	for _, idx := range indices {
		picked = append(picked, allLogs[idx])
	}
	return picked
}

func requireEqualLogs(t *testing.T, expectedLogs, actualLogs []types.Log) {
	t.Helper()
	require.Equalf(t, len(expectedLogs), len(actualLogs), "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	for i := range expectedLogs {
		require.Equalf(t, expectedLogs[i], actualLogs[i], "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	}
}

func TestBroadcaster_BroadcastsWithZeroConfirmations(t *testing.T) {
	gm := gomega.NewWithT(t)

	ethClient := new(evmmocks.Client)
	ethClient.Test(t)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("ChainID").Return(big.NewInt(0)).Maybe()
	logsChCh := make(chan evmtest.RawSub[types.Log])
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				logsChCh <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) error {
				return nil
			},
		).
		Once()
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&evmtypes.Head{Number: 1}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Return(nil, nil)

	helper := newBroadcasterHelperWithEthClient(t, ethClient, nil)
	helper.start()
	defer helper.stop()

	addr := common.HexToAddress("0xf0d54349aDdcf704F77AE15b96510dEA15cb7952")
	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(addr, nil)
	require.NoError(t, err)

	// 3 logs all in the same block
	bh := utils.NewHash()
	addr1SentLogs := []types.Log{
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       0,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       1,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       2,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
	}

	listener1 := helper.newLogListenerWithJob("1")
	helper.register(listener1, contract1, 0)
	listener2 := helper.newLogListenerWithJob("2")
	helper.register(listener2, contract1, 0)

	logs := <-logsChCh

	for _, log := range addr1SentLogs {
		logs.TrySend(log)
	}
	// Wait until the logpool has the 3 logs
	gm.Eventually(func() bool {
		helper.lb.Pause()
		defer helper.lb.Resume()
		return helper.lb.LogsFromBlock(bh) == len(addr1SentLogs)
	}, 2*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())

	// Send a block to trigger sending the logs from the pool
	// to the subscribers
	helper.lb.OnNewLongestChain(context.Background(), &evmtypes.Head{Number: 2})

	// The subs should each get exactly 3 broadcasts each
	// If we do not receive a broadcast for 1 second
	// we assume the log broadcaster is done sending.
	gm.Eventually(func() bool {
		return len(listener1.getUniqueLogs()) == len(addr1SentLogs) && len(listener2.getUniqueLogs()) == len(addr1SentLogs)
	}, 2*time.Second, cltest.DBPollingInterval).Should(gomega.BeTrue())
	gm.Consistently(func() bool {
		return len(listener1.getUniqueLogs()) == len(addr1SentLogs) && len(listener2.getUniqueLogs()) == len(addr1SentLogs)
	}, 1*time.Second, cltest.DBPollingInterval).Should(gomega.BeTrue())
}
