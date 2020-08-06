package eth_test

import (
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createJob(t *testing.T, store *store.Store) models.JobSpec {
	job := cltest.NewJob()
	err := store.ORM.CreateJob(&job)
	require.NoError(t, err)
	return job
}

func requireLogConsumptionCount(t *testing.T, store *store.Store, expectedCount int) {
	comparisonFunc := func() bool {
		observedCount, err := store.ORM.CountOf(&models.LogConsumption{})
		require.NoError(t, err)
		return observedCount == expectedCount
	}

	require.Eventually(t, comparisonFunc, 5*time.Second, 10*time.Millisecond)
}

func handleLogBroadcast(t *testing.T, lb eth.LogBroadcast) {
	consumed, err := lb.WasAlreadyConsumed()
	require.NoError(t, err)
	require.False(t, consumed)
	err = lb.MarkConsumed()
	require.NoError(t, err)
}

func TestLogBroadcaster_AwaitsInitialSubscribersOnStartup(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight uint64 = 123
	)

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	listener := new(mocks.LogListener)
	store.TxManager = txManager

	chOkayToAssert := make(chan struct{}) // avoid flaky tests

	listener.On("OnConnect").Return()
	listener.On("OnDisconnect").Return().Run(func(mock.Arguments) { close(chOkayToAssert) })

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	chSubscribe := make(chan struct{}, 10)
	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(mock.Arguments) { chSubscribe <- struct{}{} })
	txManager.On("GetLatestBlock").Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).Return([]models.Log{}, nil)

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	lb.AddDependents(2)
	lb.Start()

	lb.Register(common.Address{}, listener)

	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Eventually(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))

	lb.Stop()

	<-chOkayToAssert

	txManager.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestLogBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		numContracts        = 3
		blockHeight  uint64 = 123
	)

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	store.TxManager = txManager

	var subscribeCalls int32
	var unsubscribeCalls int32
	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(args mock.Arguments) {
			atomic.AddInt32(&subscribeCalls, 1)
		})
	txManager.On("GetLatestBlock").
		Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).
		Return(nil, nil)
	sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { atomic.AddInt32(&unsubscribeCalls, 1) })
	sub.On("Err").Return(nil)

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	lb.Start()

	type registration struct {
		common.Address
		eth.LogListener
	}
	registrations := make([]registration, numContracts)
	for i := 0; i < numContracts; i++ {
		listener := new(mocks.LogListener)
		listener.On("OnConnect").Return()
		listener.On("OnDisconnect").Return()
		registrations[i] = registration{cltest.NewAddress(), listener}
		lb.Register(registrations[i].Address, registrations[i].LogListener)
	}

	require.Eventually(t, func() bool { return atomic.LoadInt32(&subscribeCalls) == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(atomic.LoadInt32(&subscribeCalls)).Should(gomega.Equal(int32(1)))
	gomega.NewGomegaWithT(t).Consistently(atomic.LoadInt32(&unsubscribeCalls)).Should(gomega.Equal(int32(0)))

	for _, r := range registrations {
		lb.Unregister(r.Address, r.LogListener)
	}
	require.Eventually(t, func() bool { return atomic.LoadInt32(&unsubscribeCalls) == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(atomic.LoadInt32(&subscribeCalls)).Should(gomega.Equal(int32(1)))

	lb.Stop()
	gomega.NewGomegaWithT(t).Consistently(atomic.LoadInt32(&unsubscribeCalls)).Should(gomega.Equal(int32(1)))

	txManager.AssertExpectations(t)
	sub.AssertExpectations(t)
}

type simpleLogListener struct {
	handler    func(lb eth.LogBroadcast, err error)
	consumerID *models.ID
}

func (listener simpleLogListener) HandleLog(lb eth.LogBroadcast, err error) {
	listener.handler(lb, err)
}
func (listener simpleLogListener) OnConnect()    {}
func (listener simpleLogListener) OnDisconnect() {}
func (listener simpleLogListener) JobID() *models.ID {
	return listener.consumerID
}

func TestLogBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight uint64 = 0

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	store.TxManager = txManager

	chchRawLogs := make(chan chan<- models.Log, 1)
	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(1).(chan<- models.Log)
		}).
		Return(sub, nil).
		Once()
	txManager.On("GetLatestBlock").
		Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).
		Return(nil, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	lb.Start()

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()
	addr1SentLogs := []models.Log{
		{Address: addr1, BlockNumber: 1, BlockHash: cltest.NewHash()},
		{Address: addr1, BlockNumber: 2, BlockHash: cltest.NewHash()},
		{Address: addr1, BlockNumber: 3, BlockHash: cltest.NewHash()},
	}
	addr2SentLogs := []models.Log{
		{Address: addr2, BlockNumber: 4, BlockHash: cltest.NewHash()},
		{Address: addr2, BlockNumber: 5, BlockHash: cltest.NewHash()},
		{Address: addr2, BlockNumber: 6, BlockHash: cltest.NewHash()},
	}

	var addr1Logs1, addr1Logs2, addr2Logs1, addr2Logs2 []interface{}

	listener1 := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			addr1Logs1 = append(addr1Logs1, lb.Log())
			handleLogBroadcast(t, lb)
		},
		createJob(t, store).ID,
	}
	listener2 := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			addr1Logs2 = append(addr1Logs2, lb.Log())
			handleLogBroadcast(t, lb)
		},
		createJob(t, store).ID,
	}
	listener3 := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			addr2Logs1 = append(addr2Logs1, lb.Log())
			handleLogBroadcast(t, lb)
		},
		createJob(t, store).ID,
	}
	listener4 := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			addr2Logs2 = append(addr2Logs2, lb.Log())
			handleLogBroadcast(t, lb)
		},
		createJob(t, store).ID,
	}

	lb.Register(addr1, &listener1)
	lb.Register(addr1, &listener2)
	lb.Register(addr2, &listener3)
	lb.Register(addr2, &listener4)

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
	requireLogConsumptionCount(t, store, 12)

	lb.Stop()

	for i := range addr1SentLogs {
		require.Equal(t, &addr1SentLogs[i], addr1Logs1[i])
		require.Equal(t, &addr1SentLogs[i], addr1Logs2[i])
	}
	for i := range addr2SentLogs {
		require.Equal(t, &addr2SentLogs[i], addr2Logs1[i])
		require.Equal(t, &addr2SentLogs[i], addr2Logs2[i])
	}

	txManager.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestLogBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight   = 15
		expectedBlock = 5
	)

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	store.TxManager = txManager

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()

	chchRawLogs := make(chan chan<- models.Log, 1)
	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(1).(chan<- models.Log)
		}).
		Return(sub, nil).
		Twice()

	txManager.On("GetLatestBlock").
		Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(0).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr1)
			require.Len(t, query.Addresses, 1)
		}).
		Return(nil, nil).
		Once()
	txManager.On("GetLogs", mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(0).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr1)
			require.Contains(t, query.Addresses, addr2)
			require.Len(t, query.Addresses, 2)
		}).
		Return(nil, nil).
		Once()

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	listener1 := new(mocks.LogListener)
	listener2 := new(mocks.LogListener)
	listener1.On("OnConnect").Return()
	listener2.On("OnConnect").Return()
	listener1.On("OnDisconnect").Return()
	listener2.On("OnDisconnect").Return()

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	lb.Start()                    // Subscribe #1
	lb.Register(addr1, listener1) // Subscribe #2
	chRawLogs := <-chchRawLogs
	chRawLogs <- models.Log{BlockNumber: expectedBlock}
	lb.Register(addr2, listener2) // Subscribe #3
	<-chchRawLogs

	lb.Stop()

	txManager.AssertExpectations(t)
	listener1.AssertExpectations(t)
	listener2.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestDecodingLogListener(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	contract, err := eth.GetV6ContractCodec("FluxAggregator")
	require.NoError(t, err)

	type LogNewRound struct {
		models.Log
		RoundId   *big.Int
		StartedBy common.Address
		StartedAt *big.Int
	}

	logTypes := map[common.Hash]interface{}{
		eth.MustGetV6ContractEventID("FluxAggregator", "NewRound"): LogNewRound{},
	}

	var decodedLog interface{}

	listener := simpleLogListener{
		func(lb eth.LogBroadcast, innerErr error) {
			err = innerErr
			decodedLog = lb.Log()
		},
		createJob(t, store).ID,
	}

	decodingListener := eth.NewDecodingLogListener(contract, logTypes, &listener)
	rawLog := cltest.LogFromFixture(t, "../testdata/new_round_log.json")
	logBroadcast := new(mocks.LogBroadcast)

	logBroadcast.On("Log").Return(&rawLog).Once()
	logBroadcast.On("UpdateLog", mock.Anything).Run(func(args mock.Arguments) {
		logBroadcast.On("Log").Return(args.Get(0))
	})

	decodingListener.HandleLog(logBroadcast, nil)
	require.NoError(t, err)
	newRoundLog := decodedLog.(*LogNewRound)

	require.Equal(t, newRoundLog.Log, rawLog)
	require.True(t, newRoundLog.RoundId.Cmp(big.NewInt(1)) == 0)
	require.Equal(t, newRoundLog.StartedBy, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"))
	require.True(t, newRoundLog.StartedAt.Cmp(big.NewInt(15)) == 0)

	expectedErr := errors.New("oh no!")
	nilLb := new(mocks.LogBroadcast)

	logBroadcast.On("Log").Return(nil).Once()
	decodingListener.HandleLog(nilLb, expectedErr)
	require.Equal(t, err, expectedErr)
}

func TestLogBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	t.Parallel()

	logs := make(map[uint]models.Log)
	for n := 1; n < 18; n++ {
		logs[uint(n)] = models.Log{
			BlockNumber: uint64(n),
			BlockHash:   cltest.NewHash(),
			Index:       0,
		}
	}

	tests := []struct {
		name             string
		blockHeight1     uint64
		blockHeight2     uint64
		batch1           []uint
		backfillableLogs []uint
		batch2           []uint
		expectedFinal    []uint
	}{
		{
			name:             "no backfilled logs, no overlap",
			blockHeight1:     0,
			blockHeight2:     2,
			batch1:           []uint{1, 2},
			backfillableLogs: nil,
			batch2:           []uint{3, 4},
			expectedFinal:    []uint{1, 2, 3, 4},
		},
		{
			name:             "no backfilled logs, overlap",
			blockHeight1:     0,
			blockHeight2:     2,
			batch1:           []uint{1, 2},
			backfillableLogs: nil,
			batch2:           []uint{2, 3},
			expectedFinal:    []uint{1, 2, 3},
		},
		{
			name:             "backfilled logs, no overlap",
			blockHeight1:     0,
			blockHeight2:     15,
			batch1:           []uint{1, 2},
			backfillableLogs: []uint{11, 12, 15},
			batch2:           []uint{16, 17},
			expectedFinal:    []uint{1, 2, 11, 12, 15, 16, 17},
		},
		{
			name:             "backfilled logs, overlap",
			blockHeight1:     0,
			blockHeight2:     15,
			batch1:           []uint{1, 11},
			backfillableLogs: []uint{11, 12, 15},
			batch2:           []uint{16, 17},
			expectedFinal:    []uint{1, 11, 12, 15, 16, 17},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			var backfillDepth uint64 = 5 // something other than default
			store.Config.Set(orm.EnvVarName("BlockBackfillDepth"), backfillDepth)

			txManager := new(mocks.TxManager)
			sub := new(mocks.Subscription)
			store.TxManager = txManager

			chchRawLogs := make(chan chan<- models.Log, 1)

			// helper function to validate backfilled logs are being requested correctly
			checkFromBlock := func(args mock.Arguments) {
				fromBlock := args.Get(0).(ethereum.FilterQuery).FromBlock
				expected := big.NewInt(0)
				if test.blockHeight2 > backfillDepth {
					expected = big.NewInt(int64(test.blockHeight2 - backfillDepth))
				}
				require.Equal(t, expected, fromBlock)
			}

			txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					chRawLogs := args.Get(1).(chan<- models.Log)
					chchRawLogs <- chRawLogs
				}).
				Return(sub, nil).
				Twice()

			txManager.On("GetLatestBlock").Return(models.Block{Number: hexutil.Uint64(test.blockHeight1)}, nil).Once()
			txManager.On("GetLogs", mock.Anything).Return(nil, nil).Once()

			sub.On("Err").Return(nil)
			sub.On("Unsubscribe").Return()

			lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
			lb.Start()

			var recvd []*models.Log
			recvdMutex := new(sync.RWMutex)

			handleLog := func(lb eth.LogBroadcast, err error) {
				require.NoError(t, err)
				consumed, err := lb.WasAlreadyConsumed()
				require.NoError(t, err)
				if !consumed {
					recvdMutex.Lock()
					recvd = append(recvd, lb.Log().(*models.Log))
					recvdMutex.Unlock()
					err = lb.MarkConsumed()
					require.NoError(t, err)
				}
			}

			logListener := &simpleLogListener{
				handleLog,
				createJob(t, store).ID,
			}

			// Send initial logs
			lb.Register(common.Address{0}, logListener)
			chRawLogs1 := <-chchRawLogs
			for _, logNum := range test.batch1 {
				chRawLogs1 <- logs[logNum]
			}
			require.Eventually(t, func() bool {
				recvdMutex.Lock()
				defer recvdMutex.Unlock()
				return len(recvd) == len(test.batch1)
			}, 5*time.Second, 10*time.Millisecond)
			requireLogConsumptionCount(t, store, len(test.batch1))

			recvdMutex.Lock()
			for i, logNum := range test.batch1 {
				require.Equal(t, *recvd[i], logs[logNum])
			}
			recvdMutex.Unlock()

			var backfillableLogs []models.Log
			for _, logNum := range test.backfillableLogs {
				backfillableLogs = append(backfillableLogs, logs[logNum])
			}
			txManager.On("GetLatestBlock").Return(models.Block{Number: hexutil.Uint64(test.blockHeight2)}, nil).Once()
			txManager.On("GetLogs", mock.Anything).Run(checkFromBlock).Return(backfillableLogs, nil).Once()
			// Trigger resubscription
			lb.Register(common.Address{1}, &simpleLogListener{})
			chRawLogs2 := <-chchRawLogs
			for _, logNum := range test.batch2 {
				chRawLogs2 <- logs[logNum]
			}

			require.Eventually(t, func() bool {
				recvdMutex.Lock()
				defer recvdMutex.Unlock()
				return len(recvd) == len(test.expectedFinal)
			}, 5*time.Second, 10*time.Millisecond)
			requireLogConsumptionCount(t, store, len(test.expectedFinal))

			recvdMutex.Lock()
			for i, logNum := range test.expectedFinal {
				require.Equal(t, *recvd[i], logs[logNum])
			}
			recvdMutex.Unlock()

			lb.Stop()
			txManager.AssertExpectations(t)
		})
	}
}

func TestAppendLogChannel(t *testing.T) {
	t.Parallel()

	logs1 := []models.Log{
		{BlockNumber: 1},
		{BlockNumber: 2},
		{BlockNumber: 3},
		{BlockNumber: 4},
		{BlockNumber: 5},
	}

	logs2 := []models.Log{
		{BlockNumber: 6},
		{BlockNumber: 7},
		{BlockNumber: 8},
		{BlockNumber: 9},
		{BlockNumber: 10},
	}

	logs3 := []models.Log{
		{BlockNumber: 11},
		{BlockNumber: 12},
		{BlockNumber: 13},
		{BlockNumber: 14},
		{BlockNumber: 15},
	}

	ch1 := make(chan models.Log)
	ch2 := make(chan models.Log)
	ch3 := make(chan models.Log)

	chCombined := eth.ExposedAppendLogChannel(ch1, ch2)
	chCombined = eth.ExposedAppendLogChannel(chCombined, ch3)

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

func TestLogBroadcaster_InjectsLogConsumptionRecordFunctions(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight uint64 = 0

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	store.TxManager = txManager

	chchRawLogs := make(chan chan<- models.Log, 1)

	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chRawLogs := args.Get(1).(chan<- models.Log)
			chchRawLogs <- chRawLogs
		}).
		Return(sub, nil).
		Once()

	txManager.On("GetLatestBlock").Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).Return([]models.Log{}, nil).Once()

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())

	lb.Start()

	var listenerCount int32 = 0

	job := createJob(t, store)
	logListener := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			consumed, err := lb.WasAlreadyConsumed()
			require.NoError(t, err)
			require.False(t, consumed)
			err = lb.MarkConsumed()
			require.NoError(t, err)
			consumed, err = lb.WasAlreadyConsumed()
			require.NoError(t, err)
			require.True(t, consumed)
			atomic.AddInt32(&listenerCount, 1)
		},
		job.ID,
	}
	addr := common.Address{1}

	lb.Register(addr, &logListener)

	chRawLogs := <-chchRawLogs
	chRawLogs <- models.Log{Address: addr, BlockHash: cltest.NewHash(), BlockNumber: 0, Index: 0}
	chRawLogs <- models.Log{Address: addr, BlockHash: cltest.NewHash(), BlockNumber: 1, Index: 0}

	require.Eventually(t, func() bool { return atomic.LoadInt32(&listenerCount) == 2 }, 5*time.Second, 10*time.Millisecond)
	requireLogConsumptionCount(t, store, 2)
}

func TestLogBroadcaster_ProcessesLogsFromReorgs(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	sub := new(mocks.Subscription)
	store.TxManager = txManager

	const blockHeight uint64 = 0

	chchRawLogs := make(chan chan<- models.Log, 1)
	txManager.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(1).(chan<- models.Log) }).
		Return(sub, nil).
		Once()
	txManager.On("GetLatestBlock").
		Return(models.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	txManager.On("GetLogs", mock.Anything).Return([]models.Log{}, nil).Once()
	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	lb := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	lb.Start()

	blockHash0 := cltest.NewHash()
	blockHash1 := cltest.NewHash()
	blockHash2 := cltest.NewHash()
	blockHash1R := cltest.NewHash()
	blockHash2R := cltest.NewHash()

	addr := cltest.NewAddress()
	logs := []models.Log{
		{Address: addr, BlockHash: blockHash0, BlockNumber: 0, Index: 0},
		{Address: addr, BlockHash: blockHash1, BlockNumber: 1, Index: 0},
		{Address: addr, BlockHash: blockHash2, BlockNumber: 2, Index: 0},
		{Address: addr, BlockHash: blockHash1R, BlockNumber: 1, Index: 0},
		{Address: addr, BlockHash: blockHash2R, BlockNumber: 2, Index: 0},
	}

	var recvd []*models.Log
	recvdMutex := new(sync.RWMutex)

	job := createJob(t, store)
	listener := simpleLogListener{
		func(lb eth.LogBroadcast, err error) {
			require.NoError(t, err)
			ethLog := lb.Log().(*models.Log)
			recvdMutex.Lock()
			recvd = append(recvd, ethLog)
			recvdMutex.Unlock()
			handleLogBroadcast(t, lb)
		},
		job.ID,
	}

	lb.Register(addr, &listener)

	chRawLogs := <-chchRawLogs

	for i := 0; i < len(logs); i++ {
		chRawLogs <- logs[i]
	}

	require.Eventually(t, func() bool {
		recvdMutex.Lock()
		defer recvdMutex.Unlock()
		return len(recvd) == 5
	}, 5*time.Second, 10*time.Millisecond)
	requireLogConsumptionCount(t, store, 5)

	recvdMutex.Lock()
	defer recvdMutex.Unlock()
	for idx, receivedLog := range recvd {
		require.Equal(t, receivedLog, &logs[idx])
	}

	txManager.AssertExpectations(t)
}
