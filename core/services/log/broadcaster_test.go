package log_test

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/log"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBroadcaster_AwaitsInitialSubscribersOnStartup(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	const blockHeight int64 = 123
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.lb.AddDependents(2)
	helper.start()
	defer helper.stop()

	var listener = helper.newLogListener("A")
	helper.register(listener, newMockContract(), 1)

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 0 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(0)))

	helper.lb.DependentReady()

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 0 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(0)))

	helper.lb.DependentReady()

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(1)))

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.unsubscribeCallCount() }).Should(gomega.Equal(int32(1)))

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	t.Parallel()

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

	chchRawLogs := make(chan chan<- types.Log, backfillTimes)
	mockEth := newMockEthClient(chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.ethClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	blockBackfillDepth := helper.store.Config.BlockBackfillDepth()

	var backfillCount int64
	var backfillCountPtr = &backfillCount

	// the first backfill should use the height of last head saved to the db,
	// minus maxNumConfirmations of subscribers and minus blockBackfillDepth
	mockEth.checkFilterLogs = func(fromBlock int64, toBlock int64) {
		atomic.StoreInt64(backfillCountPtr, 1)
		require.Equal(t, lastStoredBlockHeight-numConfirmations-int64(blockBackfillDepth), fromBlock)
	}

	listener := helper.newLogListener("initial")
	helper.register(listener, newMockContract(), numConfirmations)
	helper.start()

	for i := 0; i < numContracts; i++ {
		listener := helper.newLogListener("")
		helper.register(listener, newMockContract(), 1)
	}

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(1)))
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.unsubscribeCallCount() }).Should(gomega.Equal(int32(0)))

	require.Eventually(t, func() bool { return atomic.LoadInt64(backfillCountPtr) == 1 }, 5*time.Second, 10*time.Millisecond)
	helper.unsubscribeAll()

	// now the backfill must use the blockBackfillDepth
	mockEth.checkFilterLogs = func(fromBlock int64, toBlock int64) {
		require.Equal(t, blockHeight-int64(blockBackfillDepth), fromBlock)
		atomic.StoreInt64(backfillCountPtr, 2)
	}

	listenerLast := helper.newLogListener("last")
	helper.register(listenerLast, newMockContract(), 1)

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() >= 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(2)))
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.unsubscribeCallCount() }).Should(gomega.Equal(int32(1)))

	require.Eventually(t, func() bool { return atomic.LoadInt64(backfillCountPtr) == 2 }, 5*time.Second, 10*time.Millisecond)

	helper.stop()
	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BackfillOnNodeStart(t *testing.T) {
	t.Parallel()

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

	chchRawLogs := make(chan chan<- types.Log, backfillTimes)
	mockEth := newMockEthClient(chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.ethClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	maxNumConfirmations := int64(10)

	var backfillCount int64
	var backfillCountPtr = &backfillCount

	listener := helper.newLogListener("one")
	helper.register(listener, newMockContract(), uint64(maxNumConfirmations))

	listener2 := helper.newLogListener("two")
	helper.register(listener2, newMockContract(), uint64(2))

	blockBackfillDepth := helper.store.Config.BlockBackfillDepth()

	// the first backfill should use the height of last head saved to the db,
	// minus maxNumConfirmations of subscribers and minus blockBackfillDepth
	mockEth.checkFilterLogs = func(fromBlock int64, toBlock int64) {
		atomic.StoreInt64(backfillCountPtr, 1)
		require.Equal(t, lastStoredBlockHeight-maxNumConfirmations-int64(blockBackfillDepth), fromBlock)
	}

	helper.start()

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return atomic.LoadInt64(backfillCountPtr) == 1 }, 5*time.Second, 10*time.Millisecond)

	helper.stop()

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() >= 1 }, 5*time.Second, 10*time.Millisecond)
	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BackfillInBatches(t *testing.T) {
	t.Parallel()

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

	chchRawLogs := make(chan chan<- types.Log, backfillTimes)
	mockEth := newMockEthClient(chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.ethClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	blockBackfillDepth := helper.store.Config.BlockBackfillDepth()
	helper.store.Config.Set(orm.EnvVarName("EthLogBackfillBatchSize"), batchSize)

	var backfillCount int64
	var backfillCountPtr = &backfillCount

	backfillStart := lastStoredBlockHeight - numConfirmations - int64(blockBackfillDepth)
	// the first backfill should start from before the last stored head
	mockEth.checkFilterLogs = func(fromBlock int64, toBlock int64) {
		times := atomic.LoadInt64(backfillCountPtr)
		logger.Warnf("Log Batch: --------- times %v - %v, %v", times, fromBlock, toBlock)

		if times <= 7 {
			require.Equal(t, backfillStart+batchSize*times, fromBlock)
			require.Equal(t, backfillStart+batchSize*(times+1)-1, toBlock)
		} else {
			// last batch is for a range of 1
			require.Equal(t, int64(120), fromBlock)
			require.Equal(t, int64(120), toBlock)
		}
		atomic.StoreInt64(backfillCountPtr, times+1)
	}

	listener := helper.newLogListener("initial")
	helper.register(listener, newMockContract(), numConfirmations)
	helper.start()

	defer helper.stop()

	require.Eventually(t, func() bool { return atomic.LoadInt64(backfillCountPtr) == expectedBatches }, 5*time.Second, 10*time.Millisecond)

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() >= 1 }, 5*time.Second, 10*time.Millisecond)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BackfillALargeNumberOfLogs(t *testing.T) {
	t.Parallel()

	const (
		lastStoredBlockHeight int64 = 10

		// a large number of blocks since lastStoredBlockHeight
		blockHeight int64 = 3000

		backfillTimes          = 1
		batchSize       uint32 = 50
		expectedBatches        = 61
	)

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
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

	chchRawLogs := make(chan chan<- types.Log, backfillTimes)
	mockEth := newMockEthClient(chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.ethClient, cltest.Head(lastStoredBlockHeight))
	helper.mockEth = mockEth

	helper.store.Config.Set(orm.EnvVarName("EthLogBackfillBatchSize"), batchSize)

	var backfillCount int64
	var backfillCountPtr = &backfillCount

	mockEth.checkFilterLogs = func(fromBlock int64, toBlock int64) {
		times := atomic.LoadInt64(backfillCountPtr)
		logger.Warnf("Log Batch: --------- times %v - %v, %v", times, fromBlock, toBlock)
		atomic.StoreInt64(backfillCountPtr, times+1)
	}

	listener := helper.newLogListener("initial")
	helper.register(listener, newMockContract(), 1)
	helper.start()

	defer helper.stop()

	require.Eventually(t, func() bool { return atomic.LoadInt64(backfillCountPtr) == expectedBatches }, 5*time.Second, 10*time.Millisecond)

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() >= 1 }, 5*time.Second, 10*time.Millisecond)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)
	contract2, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
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

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")
	listener3 := helper.newLogListener("listener 3")
	listener4 := helper.newLogListener("listener 4")

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       10,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})
	defer cleanup()

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 1)
	helper.register(listener3, contract2, 1)
	helper.register(listener4, contract2, 1)

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}
	for _, log := range addr2SentLogs {
		chRawLogs <- log
	}

	requireBroadcastCount(t, helper.store, 12)

	requireEqualLogs(t, addr1SentLogs, listener1.received.uniqueLogs)
	requireEqualLogs(t, addr1SentLogs, listener2.received.uniqueLogs)

	requireEqualLogs(t, addr2SentLogs, listener3.received.uniqueLogs)
	requireEqualLogs(t, addr2SentLogs, listener4.received.uniqueLogs)

	helper.unsubscribeAll()
	helper.stop()
	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeights(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 10)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 8)

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       10,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	requireBroadcastCount(t, helper.store, 5)
	helper.stop()

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.uniqueLogs,
	)
	requireEqualLogs(t,
		[]types.Log{
			addr1SentLogs[0],
			addr1SentLogs[1],
		},
		listener2.received.uniqueLogs,
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.uniqueLogs,
		listener1.received.logs,
	)
	requireEqualLogs(t,
		listener2.received.uniqueLogs,
		listener2.received.logs,
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

	require.Equal(t, logsOnBlocks, expectedLogsOnBlocks)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_DeletesOldLogsAfterNumberOfHeads(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint(1))
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")
	listener3 := helper.newLogListener("listener 3")
	listener4 := helper.newLogListener("listener 4")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 3)

	cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	requireBroadcastCount(t, helper.store, 6)
	<-headsDone

	helper.register(listener3, contract1, 1)
	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     6,
		EndBlock:       8,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// the new listener should still receive 2 of the 3 logs
	requireBroadcastCount(t, helper.store, 8)
	require.Equal(t, 2, len(listener3.received.uniqueLogs))

	helper.register(listener4, contract1, 1)
	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     9,
		EndBlock:       11,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// but this one should receive none
	require.Equal(t, 0, len(listener4.received.uniqueLogs))

	helper.stop()
}

func TestBroadcaster_DeletesOldLogsOnlyAfterFinalityDepth(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint(4))
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")
	listener3 := helper.newLogListener("listener 3")
	listener4 := helper.newLogListener("listener 4")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 3)

	cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	requireBroadcastCount(t, helper.store, 6)
	<-headsDone

	helper.register(listener3, contract1, 1)
	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     7,
		EndBlock:       8,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// the new listener should still receive 3 logs because of finality depth being higher than max NumConfirmations
	requireBroadcastCount(t, helper.store, 9)
	require.Equal(t, 3, len(listener3.received.uniqueLogs))

	helper.register(listener4, contract1, 1)
	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     10,
		EndBlock:       11,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// but this one should receive none
	require.Equal(t, 0, len(listener4.received.uniqueLogs))

	helper.stop()
}

func TestBroadcaster_FilterByTopicValues(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint(3))
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 20)

	topic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	field1Value1 := cltest.NewHash()
	field1Value2 := cltest.NewHash()
	field2Value1 := cltest.NewHash()
	field2Value2 := cltest.NewHash()
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNumWithTopics(1, 0, contract1.Address(), []common.Hash{topic, field1Value1, field2Value1}),
		blocks.LogOnBlockNumWithTopics(1, 1, contract1.Address(), []common.Hash{topic, field1Value2, field2Value2}),
		blocks.LogOnBlockNumWithTopics(2, 0, contract1.Address(), []common.Hash{topic, cltest.NewHash(), field2Value2}),
		blocks.LogOnBlockNumWithTopics(2, 1, contract1.Address(), []common.Hash{topic, field1Value2, cltest.NewHash()}),
	}

	listener0 := helper.newLogListener("listener 0")
	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")
	listener3 := helper.newLogListener("listener 3")
	listener4 := helper.newLogListener("listener 4")

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

	cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       5,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	<-headsDone

	require.Equal(t, 4, len(listener0.received.uniqueLogs))
	require.Equal(t, 4, len(listener1.received.uniqueLogs))
	require.Equal(t, 2, len(listener2.received.uniqueLogs))
	require.Equal(t, 3, len(listener3.received.uniqueLogs))
	require.Equal(t, 1, len(listener4.received.uniqueLogs))

	helper.stop()
}

func TestBroadcaster_BroadcastsWithOneDelayedLog(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint(2))
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 12)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),

		// this log will arrive after head with block number 3 and a previous log for it were already processed
		blocks.LogOnBlockNumWithIndex(3, 1, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	chRawLogs <- addr1SentLogs[0]
	chRawLogs <- addr1SentLogs[1]
	chRawLogs <- addr1SentLogs[2]

	cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       3,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	chRawLogs <- addr1SentLogs[3]

	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     4,
		EndBlock:       8,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	requireBroadcastCount(t, helper.store, 4)
	helper.stop()

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeightsWithLogsEarlierThanHeads(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 6)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       10,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	requireBroadcastCount(t, helper.store, 3)
	helper.stop()

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.uniqueLogs,
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.uniqueLogs,
		listener1.received.logs,
	)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_BroadcastsAtCorrectHeightsWithHeadsEarlierThanLogs(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint(2))
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := cltest.NewBlocks(t, 12)
	addr1SentLogs := []types.Log{
		blocks.LogOnBlockNum(1, contract1.Address()),
		blocks.LogOnBlockNum(2, contract1.Address()),
		blocks.LogOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	helper.register(listener1, contract1, 1)

	chRawLogs := <-helper.chchRawLogs

	cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       6,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}

	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     7,
		EndBlock:       8,
		BackfillDepth:  1,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	requireBroadcastCount(t, helper.store, 3)
	helper.stop()

	requireEqualLogs(t,
		addr1SentLogs,
		listener1.received.uniqueLogs,
	)

	// unique sends should be equal to sends overall
	requireEqualLogs(t,
		listener1.received.uniqueLogs,
		listener1.received.logs,
	)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	t.Parallel()

	const (
		backfillTimes = 1
		blockHeight   = 15
		expectedBlock = 5
	)
	var (
		ethClient = new(mocks.Client)
		sub       = new(mocks.Subscription)
		contract0 = newMockContract()
		contract1 = newMockContract()
		contract2 = newMockContract()
	)

	chchRawLogs := make(chan chan<- types.Log, backfillTimes)
	chStarted := make(chan struct{})
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(2).(chan<- types.Log)
			close(chStarted)
		}).
		Return(sub, nil).
		Once()
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(2).(chan<- types.Log)
		}).
		Return(sub, nil).
		Times(2)

	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&models.Head{Number: blockHeight}, nil)
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

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	helper := newBroadcasterHelperWithEthClient(t, ethClient, nil)
	helper.lb.AddDependents(1)
	helper.start()
	defer helper.stop()

	listener0 := helper.newLogListener("0")
	listener1 := helper.newLogListener("1")
	listener2 := helper.newLogListener("2")

	// Subscribe #0
	helper.register(listener0, contract0, 1)
	helper.lb.DependentReady()

	// Await startup
	select {
	case <-chStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("never started")
	}

	select {
	case <-chchRawLogs:
	case <-time.After(5 * time.Second):
		t.Fatal("did not subscribe")
	}

	// Subscribe #1
	helper.register(listener1, contract1, 1)

	select {
	case <-chchRawLogs:
	case <-time.After(5 * time.Second):
		t.Fatal("did not subscribe")
	}

	// Subscribe #2
	helper.register(listener2, contract2, 1)

	select {
	case <-chchRawLogs:
	case <-time.After(5 * time.Second):
		t.Fatal("did not subscribe")
	}

	cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, sub, 5*time.Second, 10*time.Millisecond)
	helper.unsubscribeAll()
}

func TestBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	t.Parallel()

	addrA := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	addrB := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	blocks := cltest.NewBlocks(t, 20)

	logsA := make(map[uint]types.Log)
	logsB := make(map[uint]types.Log)
	for n := 1; n < 18; n++ {
		logsA[uint(n)] = blocks.LogOnBlockNum(uint64(n), addrA)
		logsB[uint(n)] = blocks.LogOnBlockNum(uint64(n), addrB)
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
			t.Parallel()

			helper := newBroadcasterHelper(t, test.blockHeight1, 2)
			var backfillDepth int64 = 5
			helper.store.Config.Set(orm.EnvVarName("BlockBackfillDepth"), uint64(backfillDepth)) // something other than default

			helper.start()
			defer helper.stop()

			logListenerA := helper.newLogListener("logListenerA")
			logListenerB := helper.newLogListener("logListenerB")

			contractA, err := flux_aggregator_wrapper.NewFluxAggregator(addrA, nil)
			require.NoError(t, err)
			contractB, err := flux_aggregator_wrapper.NewFluxAggregator(addrB, nil)
			require.NoError(t, err)

			// Register listener A
			helper.register(logListenerA, contractA, 1)

			// Send initial logs
			chRawLogs1 := <-helper.chchRawLogs
			cleanup, headsDone := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock:    test.blockHeight1,
				EndBlock:      test.blockHeight2 + 1,
				BackfillDepth: backfillDepth,
				Blocks:        blocks,
				HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
					logger.Warnf("------------ HEAD TRACKABLE (%v) --------------", head.Number)
					if _, exists := logsA[uint(head.Number)]; !exists {
						logger.Warnf("  ** not exists")
						return
					} else if !batchContains(test.batch1, uint(head.Number)) {
						logger.Warnf("  ** not batchContains %v %v", head.Number, test.batch1)
						return
					}
					logger.Warnf("  ** yup!")
					select {
					case chRawLogs1 <- logsA[uint(head.Number)]:
					case <-time.After(5 * time.Second):
						t.Fatal("could not send")
					}
				})},
			})

			requireBroadcastCount(t, helper.store, len(test.batch1))
			expectedA := newReceived(pickLogs(logsA, test.batch1))
			logListenerA.requireAllReceived(t, expectedA)

			<-headsDone
			cleanup()

			helper.mockEth.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: test.blockHeight2}, nil).Once()

			combinedLogs := append(pickLogs(logsA, test.backfillableLogs), pickLogs(logsB, test.backfillableLogs)...)
			call := helper.mockEth.ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(combinedLogs, nil).Once()
			call.Run(func(args mock.Arguments) {
				// Validate that the ethereum.FilterQuery is specified correctly for the backfill that we expect
				fromBlock := args.Get(1).(ethereum.FilterQuery).FromBlock
				expected := big.NewInt(0)
				if helper.lb.LatestHead() != nil && helper.lb.LatestHead().Number > test.blockHeight2-backfillDepth {
					expected = big.NewInt(helper.lb.LatestHead().Number)
				} else if test.blockHeight2 > backfillDepth {
					expected = big.NewInt(test.blockHeight2 - backfillDepth)
				}
				require.Equal(t, expected, fromBlock)
			})

			// Register listener B (triggers re-subscription)
			helper.register(logListenerB, contractB, 1)

			// Send second batch of new logs
			chRawLogs2 := <-helper.chchRawLogs
			cleanup, _ = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock:    test.blockHeight2,
				BackfillDepth: backfillDepth,
				Blocks:        blocks,
				HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
					if _, exists := logsA[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						select {
						case chRawLogs2 <- logsA[uint(head.Number)]:
						case <-time.After(5 * time.Second):
							t.Fatal("could not send")
						}
					}
					if _, exists := logsB[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						select {
						case chRawLogs2 <- logsB[uint(head.Number)]:
						case <-time.After(5 * time.Second):
							t.Fatal("could not send")
						}
					}
				})},
			})
			defer cleanup()

			expectedA = newReceived(pickLogs(logsA, test.expectedFilteredA))
			expectedB := newReceived(pickLogs(logsB, test.expectedFilteredB))
			logListenerA.requireAllReceived(t, expectedA)
			logListenerB.requireAllReceived(t, expectedB)
			requireBroadcastCount(t, helper.store, len(test.expectedFilteredA)+len(test.expectedFilteredB))

			helper.mockEth.ethClient.AssertExpectations(t)
		})
	}
}

func TestBroadcaster_AppendLogChannel(t *testing.T) {
	t.Parallel()

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

	lb := log.NewBroadcaster(nil, nil, nil, nil)
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

	blocks := cltest.NewBlocks(t, 100)

	logListener := helper.newLogListener("logListener")

	contract := newMockContract()
	contract.On("ParseLog", mock.Anything).Return(flux_aggregator_wrapper.FluxAggregatorNewRound{}, nil).Once()
	contract.On("ParseLog", mock.Anything).Return(flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}, nil).Once()

	helper.register(logListener, contract, uint64(5))

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     3,
		BackfillDepth:  10,
		HeadTrackables: []httypes.HeadTrackable{(helper.lb).(httypes.HeadTrackable)},
		Blocks:         blocks,
	})
	defer cleanup()

	chRawLogs := <-helper.chchRawLogs
	chRawLogs <- blocks.LogOnBlockNum(0, contract.Address())
	chRawLogs <- blocks.LogOnBlockNum(1, contract.Address())

	require.Eventually(t, func() bool { return len(logListener.received.uniqueLogs) >= 2 }, 5*time.Second, 10*time.Millisecond)
	requireBroadcastCount(t, helper.store, 2)

	helper.mockEth.ethClient.AssertExpectations(t)
}

func TestBroadcaster_ProcessesLogsFromReorgsAndMissedHead(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	const startBlockHeight int64 = 0
	helper := newBroadcasterHelper(t, startBlockHeight, 1)
	helper.start()
	defer helper.stop()

	blocks := cltest.NewBlocks(t, 10)
	blocksForked := blocks.ForkAt(t, 1, 5)

	var (
		addr = cltest.NewAddress()

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

	listenerA := helper.newLogListener("listenerA")
	listenerB := helper.newLogListenerWithJobV2("listenerB")
	helper.register(listenerA, contract, 1)
	helper.register(listenerB, contract, 3)

	chRawLogs := <-helper.chchRawLogs
	go func() {
		for _, event := range events {
			switch x := event.(type) {
			case *models.Head:
				(helper.lb).(httypes.HeadTrackable).OnNewLongestChain(context.Background(), *x)
			case types.Log:
				chRawLogs <- x
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	g.Eventually(func() []uint64 { return listenerA.getUniqueLogsBlockNumbers() }, 8*time.Second, cltest.DBPollingInterval).
		Should(gomega.Equal([]uint64{0, 1, 2, 1, 2, 3, 3}))
	g.Eventually(func() []uint64 { return listenerB.getUniqueLogsBlockNumbers() }, 8*time.Second, cltest.DBPollingInterval).
		Should(gomega.Equal([]uint64{0, 1, 1, 2, 2}))

	helper.unsubscribeAll()

	require.Equal(t, expectedA, listenerA.getUniqueLogs())
	require.Equal(t, expectedB, listenerB.getUniqueLogs())

	helper.mockEth.ethClient.AssertExpectations(t)
}

func TestBroadcaster_BackfillsForNewListeners(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 2)
	helper.mockEth.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil).Times(2)
	helper.mockEth.ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil).Times(2)

	helper.start()
	defer helper.stop()

	addr1 := cltest.NewAddress()
	contract, err := flux_aggregator_wrapper.NewFluxAggregator(addr1, nil)
	require.NoError(t, err)

	listener1 := helper.newLogListener("1")
	listener2 := helper.newLogListener("2")

	topics1 := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener1, contract, topics1, 1)
	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(1)))

	<-helper.chchRawLogs

	topics2 := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
	}
	helper.registerWithTopics(listener2, contract, topics2, 1)
	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 2 }, 5*time.Second, 10*time.Millisecond)
	g.Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(2)))

	helper.unsubscribeAll()
}

func pickLogs(allLogs map[uint]types.Log, indices []uint) []types.Log {
	var picked []types.Log
	for _, idx := range indices {
		picked = append(picked, allLogs[idx])
	}
	return picked
}

func requireBroadcastCount(t *testing.T, store *strpkg.Store, expectedCount int) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	comparisonFunc := func() int {
		var count struct{ Count int }
		err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		return count.Count
	}

	g.Eventually(comparisonFunc, 5*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(expectedCount))
	g.Consistently(comparisonFunc).Should(gomega.Equal(expectedCount))
}

func requireEqualLogs(t *testing.T, expectedLogs, actualLogs []types.Log) {
	t.Helper()
	require.Equalf(t, len(expectedLogs), len(actualLogs), "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	for i := range expectedLogs {
		require.Equalf(t, expectedLogs[i], actualLogs[i], "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	}
}
