package log_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
		numContracts       = 3
		blockHeight  int64 = 123
	)

	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()
	defer helper.stop()

	for i := 0; i < numContracts; i++ {
		listener := helper.newLogListener("")
		helper.register(listener, newMockContract(), 1)
	}

	require.Eventually(t, func() bool { return helper.mockEth.subscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(1)))
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.unsubscribeCallCount() }).Should(gomega.Equal(int32(0)))

	helper.unsubscribeAll()

	require.Eventually(t, func() bool { return helper.mockEth.unsubscribeCallCount() == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.subscribeCallCount() }).Should(gomega.Equal(int32(1)))
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return helper.mockEth.unsubscribeCallCount() }).Should(gomega.Equal(int32(1)))

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

	blocks := newBlocks(t, 7)
	addr1SentLogs := []types.Log{
		blocks.logOnBlockNum(1, contract1.Address()),
		blocks.logOnBlockNum(2, contract1.Address()),
		blocks.logOnBlockNum(3, contract1.Address()),
	}
	addr2SentLogs := []types.Log{
		blocks.logOnBlockNum(4, contract2.Address()),
		blocks.logOnBlockNum(5, contract2.Address()),
		blocks.logOnBlockNum(6, contract2.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")
	listener3 := helper.newLogListener("listener 3")
	listener4 := helper.newLogListener("listener 4")

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       10,
		BackfillDepth:  10,
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
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

	blocks := newBlocks(t, 10)
	addr1SentLogs := []types.Log{
		blocks.logOnBlockNum(1, contract1.Address()),
		blocks.logOnBlockNum(2, contract1.Address()),
		blocks.logOnBlockNum(3, contract1.Address()),
	}

	listener1 := helper.newLogListener("listener 1")
	listener2 := helper.newLogListener("listener 2")

	helper.register(listener1, contract1, 1)
	helper.register(listener2, contract1, 8)

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     0,
		EndBlock:       10,
		BackfillDepth:  10,
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
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
			blockHash:      blocks.hashes[8],
		},
		{
			logBlockNumber: 2,
			blockNumber:    9,
			blockHash:      blocks.hashes[9],
		},
	}

	require.Equal(t, logsOnBlocks, expectedLogsOnBlocks)

	helper.mockEth.assertExpectations(t)
}

func TestBroadcaster_DeletesOldLogs(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := newBlocks(t, 20)
	addr1SentLogs := []types.Log{
		blocks.logOnBlockNum(1, contract1.Address()),
		blocks.logOnBlockNum(2, contract1.Address()),
		blocks.logOnBlockNum(3, contract1.Address()),
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
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
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
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// the new listener should still receive 2 of the 3 logs
	requireBroadcastCount(t, helper.store, 8)
	require.Equal(t, 2, len(listener3.received.uniqueLogs))

	helper.register(listener4, contract1, 1)
	cleanup, headsDone = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     10,
		EndBlock:       11,
		BackfillDepth:  1,
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
		Interval:       250 * time.Millisecond,
	})
	defer cleanup()

	<-headsDone

	// but this one should receive none
	require.Equal(t, 0, len(listener4.received.uniqueLogs))

	helper.stop()
}

func TestBroadcaster_BroadcastsAtCorrectHeightsWithLogsEarlierThanHeads(t *testing.T) {
	t.Parallel()

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 1)
	helper.start()

	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(cltest.NewAddress(), nil)
	require.NoError(t, err)

	blocks := newBlocks(t, 6)
	addr1SentLogs := []types.Log{
		blocks.logOnBlockNum(1, contract1.Address()),
		blocks.logOnBlockNum(2, contract1.Address()),
		blocks.logOnBlockNum(3, contract1.Address()),
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
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         blocks.hashesMap(),
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

func TestBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	t.Parallel()

	const (
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

	chchRawLogs := make(chan chan<- types.Log, 1)
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

	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, contract0.Address())
			require.Len(t, query.Addresses, 1)
		}).
		Return(nil, nil).
		Once()
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

	helper := newBroadcasterHelperWithEthClient(t, ethClient)
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

	blockHashes := make(map[int64]common.Hash)
	logsA := make(map[uint]types.Log)
	logsB := make(map[uint]types.Log)
	for n := 1; n < 18; n++ {
		blockHash := cltest.NewHash()
		blockHashes[int64(n)] = blockHash
		logsA[uint(n)] = cltest.RawNewRoundLog(t, addrA, blockHash, uint64(n), uint(n), false)
		logsB[uint(n)] = cltest.RawNewRoundLog(t, addrB, blockHash, uint64(n), uint(100+n), false)
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
			cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock:    test.blockHeight1,
				EndBlock:      test.blockHeight2 + 1,
				BackfillDepth: backfillDepth,
				Hashes:        blockHashes,
				HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
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
			expectedA := newReceived(pickLogs(t, logsA, test.batch1))
			logListenerA.requireAllReceived(t, expectedA)

			cleanup()

			helper.mockEth.ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: test.blockHeight2}, nil).Once()

			combinedLogs := append(pickLogs(t, logsA, test.backfillableLogs), pickLogs(t, logsB, test.backfillableLogs)...)
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
				Hashes:        blockHashes,
				HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable), cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
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

			expectedA = newReceived(pickLogs(t, logsA, test.expectedFilteredA))
			expectedB := newReceived(pickLogs(t, logsB, test.expectedFilteredB))
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

	lb := log.NewBroadcaster(nil, nil, nil)
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

	logListener := helper.newLogListener("logListener")

	contract := newMockContract()
	contract.On("ParseLog", mock.Anything).Return(flux_aggregator_wrapper.FluxAggregatorNewRound{}, nil).Once()
	contract.On("ParseLog", mock.Anything).Return(flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}, nil).Once()

	helper.register(logListener, contract, uint64(5))

	hash0 := cltest.NewHash()
	hash1 := cltest.NewHash()

	cleanup, _ := cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     3,
		BackfillDepth:  10,
		HeadTrackables: []strpkg.HeadTrackable{(helper.lb).(strpkg.HeadTrackable)},
		Hashes:         map[int64]common.Hash{0: hash0, 1: hash1},
	})
	defer cleanup()

	newRoundTopic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	answerUpdatedTopic := (flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}).Topic()

	chRawLogs := <-helper.chchRawLogs
	chRawLogs <- types.Log{Address: contract.Address(), BlockHash: hash0, BlockNumber: 0, Index: 0, Topics: []common.Hash{newRoundTopic, cltest.NewHash()}}
	chRawLogs <- types.Log{Address: contract.Address(), BlockHash: hash1, BlockNumber: 1, Index: 0, Topics: []common.Hash{answerUpdatedTopic, cltest.NewHash()}}

	require.Eventually(t, func() bool { return len(logListener.received.uniqueLogs) >= 2 }, 5*time.Second, 10*time.Millisecond)
	requireBroadcastCount(t, helper.store, 2)

	helper.mockEth.ethClient.AssertExpectations(t)
}

func TestBroadcaster_ProcessesLogsFromReorgsAndMissedHead(t *testing.T) {
	const startBlockHeight int64 = 0
	helper := newBroadcasterHelper(t, startBlockHeight, 1)
	helper.start()
	defer helper.stop()

	var (
		blockHash0  = cltest.NewHash()
		blockHash1  = cltest.NewHash()
		blockHash2  = cltest.NewHash()
		blockHash3  = cltest.NewHash()
		blockHash1R = cltest.NewHash()
		blockHash2R = cltest.NewHash()
		blockHash3R = cltest.NewHash()

		addr = cltest.NewAddress()

		log0        = cltest.RawNewRoundLog(t, addr, blockHash0, 0, 0, false)
		log1        = cltest.RawNewRoundLog(t, addr, blockHash1, 1, 0, false)
		log2        = cltest.RawNewRoundLog(t, addr, blockHash2, 2, 0, false)
		log1Removed = cltest.RawNewRoundLog(t, addr, blockHash1, 1, 0, true)
		log2Removed = cltest.RawNewRoundLog(t, addr, blockHash2, 2, 0, true)
		log1R       = cltest.RawNewRoundLog(t, addr, blockHash1R, 1, 0, false)
		log2R       = cltest.RawNewRoundLog(t, addr, blockHash2R, 2, 0, false)

		head0 = models.Head{Hash: blockHash0, Number: 0}
		// head1 - missing
		head2  = models.Head{Hash: blockHash2, Number: 2, Parent: &head0}
		head3  = models.Head{Hash: blockHash3, Number: 3, Parent: &head2}
		head1R = models.Head{Hash: blockHash1R, Number: 1, Parent: &head0}
		head2R = models.Head{Hash: blockHash2R, Number: 2, Parent: &head1R}
		head3R = models.Head{Hash: blockHash3R, Number: 3, Parent: &head2R}

		events = []interface{}{
			head0, log0,
			log1,
			head2, log2,
			head3,
			head1R, log1Removed, log2Removed, log1R,
			head2R, log2R,
			head3R,
		}

		expected = []types.Log{log0, log1, log2, log1R, log2R}
	)

	contract, err := flux_aggregator_wrapper.NewFluxAggregator(addr, nil)
	require.NoError(t, err)

	listener := helper.newLogListener("listener")
	helper.register(listener, contract, 2)

	chRawLogs := <-helper.chchRawLogs
	go func() {
		for _, event := range events {
			switch x := event.(type) {
			case models.Head:
				(helper.lb).(strpkg.HeadTrackable).OnNewLongestChain(context.Background(), x)
			case types.Log:
				chRawLogs <- x
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	if !assert.Eventually(t, func() bool { return len(listener.getUniqueLogs()) == 5 },
		5*time.Second, 10*time.Millisecond,
	) {
		t.Fatalf("getUniqueLogs was: %v (not equal 5)", len(listener.getUniqueLogs()))
	}

	requireBroadcastCount(t, helper.store, 5)
	helper.unsubscribeAll()

	require.Equal(t, expected, listener.getUniqueLogs())

	helper.mockEth.ethClient.AssertExpectations(t)
}

func TestBroadcaster_BackfillsForNewListeners(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	const blockHeight int64 = 0
	helper := newBroadcasterHelper(t, blockHeight, 2)
	helper.mockEth.ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil).Times(2)
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

func pickLogs(t *testing.T, allLogs map[uint]types.Log, indices []uint) []types.Log {
	var picked []types.Log
	for _, idx := range indices {
		picked = append(picked, allLogs[idx])
	}
	return picked
}

func requireBroadcastCount(t *testing.T, store *strpkg.Store, expectedCount int) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	comparisonFunc := func() bool {
		var count struct{ Count int }
		err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		return count.Count == expectedCount
	}
	require.Eventually(t, comparisonFunc, 5*time.Second, 10*time.Millisecond)
	g.Consistently(comparisonFunc).Should(gomega.Equal(true))
}

func requireEqualLogs(t *testing.T, expectedLogs, actualLogs []types.Log) {
	t.Helper()
	require.Equalf(t, len(expectedLogs), len(actualLogs), "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	for i := range expectedLogs {
		require.Equalf(t, expectedLogs[i], actualLogs[i], "log slices are not equal (len %v vs %v): expected(%v), actual(%v)", len(expectedLogs), len(actualLogs), expectedLogs, actualLogs)
	}
}
