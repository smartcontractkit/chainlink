package log_test

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBroadcaster_AwaitsInitialSubscribersOnStartup(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight int64 = 123
	)

	var (
		ethClient = new(mocks.Client)
		sub       = new(mocks.Subscription)
		listener  = new(logmocks.Listener)
	)
	store.EthClient = ethClient

	listener.On("JobID").Return(models.NewJobID())
	listener.On("JobIDV2").Return(int32(123))
	listener.On("IsV2Job").Return(true)
	listener.On("OnConnect").Return()
	listener.On("OnDisconnect").Return()

	sub.On("Unsubscribe").Maybe().Return()
	sub.On("Err").Return(nil)

	chSubscribe := make(chan struct{}, 10)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(mock.Arguments) { chSubscribe <- struct{}{} })
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil)

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
	lb.AddDependents(2)
	lb.Start()
	defer lb.Stop()

	contract := new(logmocks.AbigenContract)
	contract.On("Address").Return(common.Address{})

	_, unsubscribe := lb.Register(listener, log.ListenerOpts{
		Contract: contract,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Eventually(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))

	cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, sub, 5*time.Second, 10*time.Millisecond)
}

func TestBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		numContracts       = 3
		blockHeight  int64 = 123
	)

	var (
		ethClient        = new(mocks.Client)
		sub              = new(mocks.Subscription)
		subscribeCalls   int32
		unsubscribeCalls int32
	)
	store.EthClient = ethClient

	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(args mock.Arguments) {
			atomic.AddInt32(&subscribeCalls, 1)
		})
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil)
	sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { atomic.AddInt32(&unsubscribeCalls, 1) })
	sub.On("Err").Return(nil)

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
	lb.Start()
	defer lb.Stop()

	type registration struct {
		log.AbigenContract
		log.Listener
	}
	registrations := make([]registration, numContracts)
	var unsubscribes []func()
	for i := 0; i < numContracts; i++ {
		contract := new(logmocks.AbigenContract)
		contract.On("Address").Return(cltest.NewAddress())

		listener := new(logmocks.Listener)
		listener.On("OnConnect").Return()
		listener.On("OnDisconnect").Return()
		listener.On("JobID").Return(models.NewJobID())
		listener.On("JobIDV2").Return(int32(i))
		listener.On("IsV2Job").Return(i%2 == 0)
		registrations[i] = registration{contract, listener}
		_, unsubscribe := lb.Register(listener, log.ListenerOpts{
			Contract: contract,
			Logs: []generated.AbigenLog{
				flux_aggregator_wrapper.FluxAggregatorNewRound{},
				flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
			},
		})
		unsubscribes = append(unsubscribes, unsubscribe)
		defer unsubscribe()
	}

	require.Eventually(t, func() bool { return atomic.LoadInt32(&subscribeCalls) == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return atomic.LoadInt32(&subscribeCalls) }).Should(gomega.Equal(int32(1)))
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return atomic.LoadInt32(&unsubscribeCalls) }).Should(gomega.Equal(int32(0)))

	for _, unsub := range unsubscribes {
		unsub()
	}
	require.Eventually(t, func() bool { return atomic.LoadInt32(&unsubscribeCalls) == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return atomic.LoadInt32(&subscribeCalls) }).Should(gomega.Equal(int32(1)))

	lb.Stop()
	gomega.NewGomegaWithT(t).Consistently(func() int32 { return atomic.LoadInt32(&unsubscribeCalls) }).Should(gomega.Equal(int32(1)))

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight int64 = 0

	var (
		ethClient     = new(mocks.Client)
		sub           = new(mocks.Subscription)
		addr1         = cltest.NewAddress()
		addr2         = cltest.NewAddress()
		contract1     = new(logmocks.AbigenContract)
		contract2     = new(logmocks.AbigenContract)
		addr1SentLogs = []types.Log{
			{Address: addr1, BlockNumber: 1, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr1, BlockNumber: 2, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr1, BlockNumber: 3, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
		}
		addr2SentLogs = []types.Log{
			{Address: addr2, BlockNumber: 4, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr2, BlockNumber: 5, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr2, BlockNumber: 6, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}},
		}
	)
	store.EthClient = ethClient

	contract1.On("Address").Return(addr1)
	contract2.On("Address").Return(addr2)
	contract1.On("ParseLog", mock.Anything).Return(struct{}{}, nil)
	contract2.On("ParseLog", mock.Anything).Return(struct{}{}, nil)

	chchRawLogs := make(chan chan<- types.Log, 1)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(2).(chan<- types.Log)
		}).
		Return(sub, nil).
		Once()
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
	lb.Start()
	defer lb.Stop()

	var addr1Logs1, addr1Logs2, addr2Logs1, addr2Logs2 []types.Log

	listener1 := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			addr1Logs1 = append(addr1Logs1, lb.RawLog())
			handleLogBroadcast(t, lb)
		},
		consumerID: createJob(t, store).ID,
	}
	listener2 := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			addr1Logs2 = append(addr1Logs2, lb.RawLog())
			handleLogBroadcast(t, lb)
		},
		consumerID: createJob(t, store).ID,
	}
	listener3 := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			addr2Logs1 = append(addr2Logs1, lb.RawLog())
			handleLogBroadcast(t, lb)
		},
		consumerID: createJob(t, store).ID,
	}
	listener4 := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			addr2Logs2 = append(addr2Logs2, lb.RawLog())
			handleLogBroadcast(t, lb)
		},
		consumerID: createJob(t, store).ID,
	}

	cleanup = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     6,
		HeadTrackables: []strpkg.HeadTrackable{lb},
	})

	defer cleanup()

	_, unsubscribe := lb.Register(listener1, log.ListenerOpts{
		Contract: contract1,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()
	_, unsubscribe = lb.Register(listener2, log.ListenerOpts{
		Contract: contract1,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()
	_, unsubscribe = lb.Register(listener3, log.ListenerOpts{
		Contract: contract2,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()
	_, unsubscribe = lb.Register(listener4, log.ListenerOpts{
		Contract: contract2,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	chRawLogs := <-chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}
	for _, log := range addr2SentLogs {
		chRawLogs <- log
	}

	require.Eventually(t, func() bool { return len(addr1Logs1) == len(addr1SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr1Logs2) == len(addr1SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr2Logs1) == len(addr2SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr2Logs2) == len(addr2SentLogs) }, time.Second, 10*time.Millisecond)
	requireBroadcastCount(t, store, 12)

	lb.Stop()

	for i := range addr1SentLogs {
		require.Equal(t, addr1SentLogs[i], addr1Logs1[i])
		require.Equal(t, addr1SentLogs[i], addr1Logs2[i])
	}
	for i := range addr2SentLogs {
		require.Equal(t, addr2SentLogs[i], addr2Logs1[i])
		require.Equal(t, addr2SentLogs[i], addr2Logs2[i])
	}

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight   = 15
		expectedBlock = 5
	)
	var (
		ethClient = new(mocks.Client)
		sub       = new(mocks.Subscription)
		listener0 = new(logmocks.Listener)
		listener1 = new(logmocks.Listener)
		listener2 = new(logmocks.Listener)
		addr0     = cltest.NewAddress()
		addr1     = cltest.NewAddress()
		addr2     = cltest.NewAddress()
		contract0 = new(logmocks.AbigenContract)
		contract1 = new(logmocks.AbigenContract)
		contract2 = new(logmocks.AbigenContract)
	)
	store.EthClient = ethClient

	contract0.On("Address").Return(addr0)
	contract1.On("Address").Return(addr1)
	contract2.On("Address").Return(addr2)

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
			require.Contains(t, query.Addresses, addr0)
			require.Len(t, query.Addresses, 1)
		}).
		Return(nil, nil).
		Once()
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr0)
			require.Contains(t, query.Addresses, addr1)
			require.Len(t, query.Addresses, 2)
		}).
		Return(nil, nil).
		Once()
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr0)
			require.Contains(t, query.Addresses, addr1)
			require.Contains(t, query.Addresses, addr2)
			require.Len(t, query.Addresses, 3)
		}).
		Return(nil, nil).
		Once()

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	listener0.On("JobID").Return(models.NewJobID()).Maybe()
	listener0.On("JobIDV2").Return(int32(123)).Maybe()
	listener0.On("IsV2Job").Return(true).Maybe()
	listener1.On("JobID").Return(models.NewJobID()).Maybe()
	listener1.On("JobIDV2").Return(int32(456)).Maybe()
	listener1.On("IsV2Job").Return(true).Maybe()
	listener2.On("JobID").Return(models.NewJobID()).Maybe()
	listener2.On("JobIDV2").Return(int32(789)).Maybe()
	listener2.On("IsV2Job").Return(true).Maybe()
	listener0.On("OnConnect").Return().Maybe()
	listener1.On("OnConnect").Return().Maybe()
	listener2.On("OnConnect").Return().Maybe()
	listener0.On("OnDisconnect").Return().Maybe()
	listener1.On("OnDisconnect").Return().Maybe()
	listener2.On("OnDisconnect").Return().Maybe()

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, ethClient, store.Config)
	lb.AddDependents(1)
	lb.Start() // Subscribe #0
	defer lb.Stop()

	_, unsubscribe := lb.Register(listener0, log.ListenerOpts{
		Contract: contract0,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	lb.DependentReady()
	<-chStarted // Await startup
	<-chchRawLogs

	// Subscribe #1
	_, unsubscribe = lb.Register(listener1, log.ListenerOpts{
		Contract: contract1,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	<-chchRawLogs

	// Subscribe #2
	_, unsubscribe = lb.Register(listener1, log.ListenerOpts{
		Contract: contract1,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})

	<-chchRawLogs

	cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, listener0, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, listener1, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, listener2, 5*time.Second, 10*time.Millisecond)
	cltest.EventuallyExpectationsMet(t, sub, 5*time.Second, 10*time.Millisecond)
}

func TestBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	t.Parallel()

	addrA := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	addrB := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	logsA := make(map[uint]types.Log)
	logsB := make(map[uint]types.Log)
	for n := 1; n < 18; n++ {
		blockHash := cltest.NewHash()
		logsA[uint(n)] = types.Log{
			Address:     addrA,
			BlockNumber: uint64(n),
			BlockHash:   blockHash,
			Index:       uint(n),
			Topics:      []common.Hash{},
			Data:        []byte{},
		}
		logsB[uint(n)] = types.Log{
			Address:     addrB,
			BlockNumber: uint64(n),
			BlockHash:   blockHash,
			Index:       uint(100 + n),
			Topics:      []common.Hash{},
			Data:        []byte{},
		}
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

			blockHeight2:     2,
			backfillableLogs: nil,
			batch2:           []uint{3, 4},

			expectedFilteredA: []uint{1, 2, 3, 4},
			expectedFilteredB: []uint{3, 4},
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

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			var (
				ethClient           = new(mocks.Client)
				sub                 = new(mocks.Subscription)
				backfillDepth int64 = 5 // something other than default
			)

			store.Config.Set(orm.EnvVarName("BlockBackfillDepth"), uint64(backfillDepth))
			store.EthClient = ethClient

			chchRawLogs := make(chan chan<- types.Log, 1)
			ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					chRawLogs := args.Get(2).(chan<- types.Log)
					chchRawLogs <- chRawLogs
				}).
				Return(sub, nil).
				Twice()

			ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: test.blockHeight1}, nil).Once()
			ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil).Once()

			sub.On("Err").Return(nil)
			sub.On("Unsubscribe").Return()

			orm := log.NewORM(store.DB)
			lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
			lb.Start()
			defer lb.Stop()

			var recvdA received
			var recvdB received

			logListenerA := &simpleLogListener{
				handler: func(lb log.Broadcast) {
					consumed, err := lb.WasAlreadyConsumed()
					require.NoError(t, err)

					recvdA.Lock()
					defer recvdA.Unlock()
					if !consumed {
						recvdA.logs = append(recvdA.logs, lb.RawLog())
						err = lb.MarkConsumed()
						require.NoError(t, err)
					}
				},
				consumerID: createJob(t, store).ID,
			}

			logListenerB := &simpleLogListener{
				handler: func(lb log.Broadcast) {
					consumed, err := lb.WasAlreadyConsumed()
					require.NoError(t, err)

					recvdB.Lock()
					defer recvdB.Unlock()
					if !consumed {
						recvdB.logs = append(recvdB.logs, lb.RawLog())
						err = lb.MarkConsumed()
						require.NoError(t, err)
					}
				},
				consumerID: createJob(t, store).ID,
			}

			contractA := new(logmocks.AbigenContract)
			contractB := new(logmocks.AbigenContract)
			contractA.On("Address").Return(addrA)
			contractB.On("Address").Return(addrB)
			contractA.On("ParseLog", mock.Anything).Return(struct{}{}, nil)
			contractB.On("ParseLog", mock.Anything).Return(struct{}{}, nil)

			// Register listener A
			_, unsubscribe := lb.Register(logListenerA, log.ListenerOpts{
				Contract: contractA,
				Logs: []generated.AbigenLog{
					flux_aggregator_wrapper.FluxAggregatorNewRound{},
					flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
				},
			})
			defer unsubscribe()

			// Send initial logs
			chRawLogs1 := <-chchRawLogs
			cleanup = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock: test.blockHeight1,
				EndBlock:   test.blockHeight2,
				HeadTrackables: []strpkg.HeadTrackable{lb, cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
					if _, exists := logsA[uint(head.Number)]; !exists {
						return
					} else if !batchContains(test.batch1, uint(head.Number)) {
						return
					}
					chRawLogs1 <- logsA[uint(head.Number)]
				})},
			})
			defer cleanup()

			expectedA := received{
				logs: pickLogs(t, logsA, test.batch1),
			}
			requireAllReceived(t, &expectedA, &recvdA)
			requireBroadcastCount(t, store, len(test.batch1))

			cleanup()

			ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: test.blockHeight2}, nil).Once()

			combinedLogs := append(pickLogs(t, logsA, test.backfillableLogs), pickLogs(t, logsB, test.backfillableLogs)...)
			call := ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(combinedLogs, nil).Once()
			call.Run(func(args mock.Arguments) {
				// Validate that the ethereum.FilterQuery is specified correctly for the backfill that we expect
				fromBlock := args.Get(1).(ethereum.FilterQuery).FromBlock
				expected := big.NewInt(0)
				if test.blockHeight2 > backfillDepth {
					expected = big.NewInt(int64(test.blockHeight2 - backfillDepth))
				}
				require.Equal(t, expected, fromBlock)
			})

			// Register listener B (triggers resubscription)
			_, unsubscribe = lb.Register(logListenerB, log.ListenerOpts{
				Contract: contractB,
				Logs: []generated.AbigenLog{
					flux_aggregator_wrapper.FluxAggregatorNewRound{},
					flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
				},
			})
			defer unsubscribe()

			// Send second batch of new logs
			chRawLogs2 := <-chchRawLogs
			cleanup = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
				StartBlock: test.blockHeight2,
				HeadTrackables: []strpkg.HeadTrackable{lb, cltest.HeadTrackableFunc(func(_ context.Context, head models.Head) {
					if _, exists := logsA[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						chRawLogs2 <- logsA[uint(head.Number)]
					}
					if _, exists := logsB[uint(head.Number)]; exists && batchContains(test.batch2, uint(head.Number)) {
						chRawLogs2 <- logsB[uint(head.Number)]
					}
				})},
			})
			defer cleanup()

			expectedA = received{
				logs: pickLogs(t, logsA, test.expectedFilteredA),
			}
			expectedB := received{
				logs: pickLogs(t, logsB, test.expectedFilteredB),
			}
			requireAllReceived(t, &expectedA, &recvdA)
			requireAllReceived(t, &expectedB, &recvdB)
			requireBroadcastCount(t, store, len(test.expectedFilteredA)+len(test.expectedFilteredB))

			ethClient.AssertExpectations(t)
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
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight int64 = 0

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)
	store.EthClient = ethClient

	chchRawLogs := make(chan chan<- types.Log, 1)

	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chRawLogs := args.Get(2).(chan<- types.Log)
			chchRawLogs <- chRawLogs
		}).
		Return(sub, nil).
		Once()

	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil).Once()

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)

	lb.Start()
	defer lb.Stop()

	var broadcastCount int32 = 0

	job := createJob(t, store)
	logListener := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			consumed, err := lb.WasAlreadyConsumed()
			require.NoError(t, err)
			require.False(t, consumed)
			err = lb.MarkConsumed()
			require.NoError(t, err)
			consumed, err = lb.WasAlreadyConsumed()
			require.NoError(t, err)
			require.True(t, consumed)
			atomic.AddInt32(&broadcastCount, 1)
		},
		consumerID: job.ID,
	}
	addr := cltest.NewAddress()
	contract := new(logmocks.AbigenContract)
	contract.On("Address").Return(addr)
	contract.On("ParseLog", mock.Anything).Return(struct{}{}, nil)

	_, unsubscribe := lb.Register(logListener, log.ListenerOpts{
		Contract: contract,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	cleanup = cltest.SimulateIncomingHeads(t, cltest.SimulateIncomingHeadsArgs{
		StartBlock:     3,
		HeadTrackables: []strpkg.HeadTrackable{lb},
	})
	defer cleanup()

	chRawLogs := <-chchRawLogs
	chRawLogs <- types.Log{Address: addr, BlockHash: cltest.NewHash(), BlockNumber: 0, Index: 0}
	chRawLogs <- types.Log{Address: addr, BlockHash: cltest.NewHash(), BlockNumber: 1, Index: 0}

	require.Eventually(t, func() bool { return atomic.LoadInt32(&broadcastCount) == 2 }, 5*time.Second, 10*time.Millisecond)
	requireBroadcastCount(t, store, 2)
}

func TestBroadcaster_ProcessesLogsFromReorgs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)
	store.EthClient = ethClient

	const (
		startBlockHeight int64 = 0
	)

	chchRawLogs := make(chan chan<- types.Log, 1)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(2).(chan<- types.Log) }).
		Return(sub, nil).
		Once()
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: startBlockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil).Once()
	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
	lb.Start()
	defer lb.Stop()

	var (
		blockHash0  = cltest.NewHash()
		blockHash1  = cltest.NewHash()
		blockHash2  = cltest.NewHash()
		blockHash1R = cltest.NewHash()
		blockHash2R = cltest.NewHash()

		addr = cltest.NewAddress()
		logs = []types.Log{
			{Address: addr, BlockHash: blockHash0, BlockNumber: 0, Index: 0, Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr, BlockHash: blockHash1, BlockNumber: 1, Index: 0, Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr, BlockHash: blockHash2, BlockNumber: 2, Index: 0, Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr, BlockHash: blockHash1, BlockNumber: 1, Index: 0, Topics: []common.Hash{}, Data: []byte{}, Removed: true},
			{Address: addr, BlockHash: blockHash2, BlockNumber: 2, Index: 0, Topics: []common.Hash{}, Data: []byte{}, Removed: true},
			{Address: addr, BlockHash: blockHash1R, BlockNumber: 1, Index: 0, Topics: []common.Hash{}, Data: []byte{}},
			{Address: addr, BlockHash: blockHash2R, BlockNumber: 2, Index: 0, Topics: []common.Hash{}, Data: []byte{}},
		}
	)

	job := createJob(t, store)
	var recvd []log.Broadcast
	var recvdMu sync.Mutex
	listener := &simpleLogListener{
		handler: func(lb log.Broadcast) {
			recvdMu.Lock()
			defer recvdMu.Unlock()
			recvd = append(recvd, lb)
			handleLogBroadcast(t, lb)
		},
		consumerID: job.ID,
	}

	contract := new(logmocks.AbigenContract)
	contract.On("Address").Return(addr)
	contract.On("ParseLog", mock.Anything).Return(struct{}{}, nil)
	_, unsubscribe := lb.Register(listener, log.ListenerOpts{
		Contract: contract,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
	})
	defer unsubscribe()

	chRawLogs := <-chchRawLogs
	for i := 0; i < len(logs); i++ {
		lb.OnNewLongestChain(context.Background(), models.NewHead(big.NewInt(int64(logs[i].BlockNumber)), logs[i].BlockHash, common.Hash{123}, 0))
		chRawLogs <- logs[i]
		time.Sleep(500 * time.Millisecond)
	}

	require.Eventually(t, func() bool {
		return len(recvd) == 5
	}, 5*time.Second, 10*time.Millisecond)
	requireBroadcastCount(t, store, 3)

	var nonRemoved []types.Log
	for _, log := range logs {
		if !log.Removed {
			nonRemoved = append(nonRemoved, log)
		}
	}
	require.Len(t, nonRemoved, 5)
	for idx, broadcast := range recvd {
		require.Equal(t, nonRemoved[idx], broadcast.RawLog())
	}

	ethClient.AssertExpectations(t)
}

func TestBroadcaster_BackfillsForNewListeners(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)
	store.EthClient = ethClient

	const blockHeight int64 = 0

	chchRawLogs := make(chan chan<- types.Log, 1)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(2).(chan<- types.Log) }).
		Return(sub, nil).
		Once()
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil).Once()
	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)
	lb.Start()
	defer lb.Stop()
}

func pickLogs(t *testing.T, allLogs map[uint]types.Log, indices []uint) []types.Log {
	var picked []types.Log
	for _, idx := range indices {
		picked = append(picked, allLogs[idx])
	}
	return picked
}

type received struct {
	logs []types.Log
	sync.Mutex
}

func requireAllReceived(t *testing.T, expectedState, state *received) {
	require.Eventually(t, func() bool {
		state.Lock()
		defer state.Unlock()
		return len(state.logs) == len(expectedState.logs)
	}, 10*time.Second, 10*time.Millisecond)

	state.Lock()
	for i := range expectedState.logs {
		require.Equal(t, expectedState.logs[i], state.logs[i])
	}
	state.Unlock()
}

func requireBroadcastCount(t *testing.T, store *strpkg.Store, expectedCount int) {
	t.Helper()

	comparisonFunc := func() bool {
		var count struct{ Count int }
		err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		return count.Count == expectedCount
	}
	require.Eventually(t, comparisonFunc, 5*time.Second, 10*time.Millisecond)
}

func handleLogBroadcast(t *testing.T, lb log.Broadcast) {
	t.Helper()

	consumed, err := lb.WasAlreadyConsumed()
	require.NoError(t, err)
	require.False(t, consumed)
	err = lb.MarkConsumed()
	require.NoError(t, err)
}
