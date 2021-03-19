package log_test

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type broadcasterHelper struct {
	t             *testing.T
	lb            log.Broadcaster
	store         *store.Store
	mockEth       *mockEth
	chchRawLogs   chan chan<- types.Log
	toUnsubscribe []func()
	storeCleanup  func()
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int) *broadcasterHelper {
	store, cleanup := cltest.NewStore(t)

	chchRawLogs := make(chan chan<- types.Log, 1)
	mockEth := newMockEthClient(chchRawLogs, blockHeight, timesSubscribe)
	store.EthClient = mockEth.ethClient

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)

	return &broadcasterHelper{
		t:             t,
		lb:            lb,
		store:         store,
		mockEth:       mockEth,
		chchRawLogs:   chchRawLogs,
		toUnsubscribe: make([]func(), 0),
		storeCleanup:  cleanup,
	}
}

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient eth.Client) *broadcasterHelper {
	store, cleanup := cltest.NewStore(t)

	store.EthClient = ethClient

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config)

	return &broadcasterHelper{
		t:             t,
		lb:            lb,
		store:         store,
		toUnsubscribe: make([]func(), 0),
		storeCleanup:  cleanup,
	}
}

func (helper *broadcasterHelper) newLogListener(name string) *simpleLogListener {
	return newLogListener(helper.t, helper.store, name)
}
func (helper *broadcasterHelper) start() {
	err := helper.lb.Start()
	require.NoError(helper.t, err)
}

func (helper *broadcasterHelper) register(listener log.Listener, contract log.AbigenContract, numConfirmations uint64) {
	logs := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener, contract, logs, numConfirmations)
}

func (helper *broadcasterHelper) registerWithTopics(listener log.Listener, contract log.AbigenContract, topics []generated.AbigenLog, numConfirmations uint64) {
	_, unsubscribe := helper.lb.Register(listener, log.ListenerOpts{
		Contract:         contract,
		Logs:             topics,
		NumConfirmations: numConfirmations,
	})

	helper.toUnsubscribe = append(helper.toUnsubscribe, unsubscribe)
}

func (helper *broadcasterHelper) unsubscribeAll() {
	for _, unsubscribe := range helper.toUnsubscribe {
		unsubscribe()
	}
	time.Sleep(100 * time.Millisecond)
}
func (helper *broadcasterHelper) stop() {
	err := helper.lb.Stop()
	require.NoError(helper.t, err)
	helper.storeCleanup()
}

func newMockContract() *logmocks.AbigenContract {
	addr := cltest.NewAddress()
	contract := new(logmocks.AbigenContract)
	contract.On("Address").Return(addr)
	return contract
}

type logOnBlock struct {
	logBlockNumber uint64
	blockNumber    uint64
	blockHash      common.Hash
}

func (l logOnBlock) String() string {
	return fmt.Sprintf("blockInfo(log:%v received on: %v %s)", l.logBlockNumber, l.blockNumber, l.blockHash)
}

type received struct {
	uniqueLogs []types.Log
	logs       []types.Log
	broadcasts []log.Broadcast
	sync.Mutex
}

func newReceived(logs []types.Log) *received {
	var rec received
	rec.logs = logs
	rec.uniqueLogs = logs
	return &rec
}

func (rec *received) logsOnBlocks() []logOnBlock {
	blocks := make([]logOnBlock, 0)
	for _, broadcast := range rec.broadcasts {
		blocks = append(blocks, logOnBlock{
			logBlockNumber: broadcast.RawLog().BlockNumber,
			blockNumber:    broadcast.LatestBlockNumber(),
			blockHash:      broadcast.LatestBlockHash(),
		})
	}
	return blocks
}

type simpleLogListener struct {
	handler    func(lb log.Broadcast) bool
	consumerID models.JobID
	name       string
	received   *received
}

func newLogListener(t *testing.T, store *store.Store, name string) *simpleLogListener {
	var rec received
	return &simpleLogListener{
		handler: func(lb log.Broadcast) bool {
			return handleLogBroadcast(t, lb)
		},
		consumerID: createJob(t, store).ID,
		name:       name,
		received:   &rec,
	}
}

func (listener *simpleLogListener) HandleLog(lb log.Broadcast) {
	logger.Warnf(">>>>>>>>>>>>>>> Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())
	listener.received.Lock()
	defer listener.received.Unlock()
	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handler(lb)

	if !consumed {
		listener.received.uniqueLogs = append(listener.received.uniqueLogs, lb.RawLog())
	}
}

func (listener *simpleLogListener) OnConnect()    {}
func (listener *simpleLogListener) OnDisconnect() {}
func (listener *simpleLogListener) JobID() models.JobID {
	return listener.consumerID
}
func (listener *simpleLogListener) IsV2Job() bool {
	return false
}
func (listener *simpleLogListener) JobIDV2() int32 {
	return 0
}

func (listener *simpleLogListener) getUniqueLogs() []types.Log {
	return listener.received.uniqueLogs
}

func (listener *simpleLogListener) requireAllReceived(t *testing.T, expectedState *received) {
	received := listener.received
	require.Eventually(t, func() bool {
		received.Lock()
		defer received.Unlock()
		return len(received.uniqueLogs) == len(expectedState.uniqueLogs)
	}, 10*time.Second, 10*time.Millisecond, "len(received.logs): %v is not equal len(expectedState.logs): %v", len(received.logs), len(expectedState.logs))

	received.Lock()
	for i := range expectedState.uniqueLogs {
		require.Equal(t, expectedState.uniqueLogs[i], received.uniqueLogs[i])
	}
	received.Unlock()
}

func handleLogBroadcast(t *testing.T, lb log.Broadcast) bool {
	t.Helper()
	consumed, err := lb.WasAlreadyConsumed()
	require.NoError(t, err)
	if !consumed {

		err = lb.MarkConsumed()
		require.NoError(t, err)

		consumed2, err := lb.WasAlreadyConsumed()
		require.NoError(t, err)
		require.True(t, consumed2)
	}
	return consumed
}

type mockListener struct {
	jobID   models.JobID
	jobIDV2 int32
}

func (l *mockListener) JobID() models.JobID     { return l.jobID }
func (l *mockListener) JobIDV2() int32          { return l.jobIDV2 }
func (l *mockListener) IsV2Job() bool           { return l.jobID.IsZero() }
func (l *mockListener) OnConnect()              {}
func (l *mockListener) OnDisconnect()           {}
func (l *mockListener) HandleLog(log.Broadcast) {}

func createJob(t *testing.T, store *store.Store) models.JobSpec {
	t.Helper()

	job := cltest.NewJob()
	err := store.ORM.CreateJob(&job)
	require.NoError(t, err)
	return job
}

type mockEth struct {
	ethClient        *mocks.Client
	sub              *mocks.Subscription
	subscribeCalls   int32
	unsubscribeCalls int32
}

func (mock *mockEth) assertExpectations(t *testing.T) {
	mock.ethClient.AssertExpectations(t)
	mock.sub.AssertExpectations(t)
}

func (mock *mockEth) subscribeCallCount() int32 {
	return atomic.LoadInt32(&(mock.subscribeCalls))
}

func (mock *mockEth) unsubscribeCallCount() int32 {
	return atomic.LoadInt32(&(mock.unsubscribeCalls))
}

func newMockEthClient(chchRawLogs chan chan<- types.Log, blockHeight int64, timesSubscribe int) *mockEth {
	mockEth := &mockEth{
		ethClient: new(mocks.Client),
		sub:       new(mocks.Subscription),
	}
	mockEth.ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			atomic.AddInt32(&(mockEth.subscribeCalls), 1)
			chchRawLogs <- args.Get(2).(chan<- types.Log)
		}).
		Return(mockEth.sub, nil).
		Times(timesSubscribe)
	mockEth.ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&models.Head{Number: blockHeight}, nil).Times(1)
	mockEth.ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return(nil, nil).Times(1)
	mockEth.sub.On("Err").Return(nil)
	mockEth.sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { atomic.AddInt32(&(mockEth.unsubscribeCalls), 1) })
	return mockEth
}

type blocks struct {
	t      *testing.T
	hashes []common.Hash
}

func (lb *blocks) logOnBlockNum(i uint64, addr common.Address) types.Log {
	return cltest.RawNewRoundLog(lb.t, addr, lb.hashes[i], i, 0, false)
}
func (lb *blocks) hashesMap() map[int64]common.Hash {
	h := make(map[int64]common.Hash)
	for i := 0; i < len(lb.hashes); i++ {
		h[int64(i)] = lb.hashes[i]
	}
	return h
}
func newBlocks(t *testing.T, numHashes int) *blocks {
	hashes := make([]common.Hash, 0)
	for i := 0; i < numHashes; i++ {
		hashes = append(hashes, cltest.NewHash())
	}
	return &blocks{
		t:      t,
		hashes: hashes,
	}
}
