package log_test

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type broadcasterHelper struct {
	t       *testing.T
	lb      log.Broadcaster
	store   *store.Store
	mockEth *mockEth

	// each received channel corresponds to one eth subscription
	chchRawLogs   chan chan<- types.Log
	toUnsubscribe []func()
	storeCleanup  func()
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int) *broadcasterHelper {
	store, cleanup := cltest.NewStore(t)

	chchRawLogs := make(chan chan<- types.Log, timesSubscribe)

	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: timesSubscribe,
		HeaderByNumber:      1,
		FilterLogs:          1,
	}

	mockEth := newMockEthClient(chchRawLogs, blockHeight, expectedCalls)
	store.EthClient = mockEth.ethClient

	dborm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(dborm, store.EthClient, store.Config, nil)
	store.Config.Set(orm.EnvVarName("EthFinalityDepth"), uint64(10))
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

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient eth.Client, highestSeenHead *models.Head) *broadcasterHelper {
	store, cleanup := cltest.NewStore(t)

	store.EthClient = ethClient

	orm := log.NewORM(store.DB)
	lb := log.NewBroadcaster(orm, store.EthClient, store.Config, highestSeenHead)

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
	unsubscribe := helper.lb.Register(listener, log.ListenerOpts{
		Contract:         contract,
		Logs:             topics,
		NumConfirmations: numConfirmations,
	})

	helper.toUnsubscribe = append(helper.toUnsubscribe, unsubscribe)
}

func (helper *broadcasterHelper) registerWithTopicValues(listener log.Listener, contract log.AbigenContract, numConfirmations uint64,
	topics map[common.Hash][][]log.Topic) {

	unsubscribe := helper.lb.Register(listener, log.ListenerOpts{
		Contract:         contract,
		LogsWithTopics:   topics,
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
	err := helper.lb.Close()
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
	var blocks []logOnBlock
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
	consumerID models.JobID
	name       string
	received   *received
	t          *testing.T
	db         *gorm.DB
}

func newLogListener(t *testing.T, store *store.Store, name string) *simpleLogListener {
	var rec received
	return &simpleLogListener{
		db:         store.DB,
		consumerID: createJob(t, store).ID,
		name:       name,
		received:   &rec,
		t:          t,
	}
}

func (listener simpleLogListener) HandleLog(lb log.Broadcast) {
	logger.Warnf("Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())
	listener.received.Lock()
	defer listener.received.Unlock()
	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handleLogBroadcast(listener.t, lb)

	if !consumed {
		listener.received.uniqueLogs = append(listener.received.uniqueLogs, lb.RawLog())
	} else {
		logger.Warnf("Listener %v: Log was already consumed!", listener.name)
	}
}

func (listener simpleLogListener) OnConnect()    {}
func (listener simpleLogListener) OnDisconnect() {}
func (listener simpleLogListener) JobID() models.JobID {
	return listener.consumerID
}
func (listener simpleLogListener) IsV2Job() bool {
	return false
}
func (listener simpleLogListener) JobIDV2() int32 {
	return 0
}

func (listener simpleLogListener) getUniqueLogs() []types.Log {
	return listener.received.uniqueLogs
}

func (listener simpleLogListener) requireAllReceived(t *testing.T, expectedState *received) {
	received := listener.received
	require.Eventually(t, func() bool {
		received.Lock()
		defer received.Unlock()
		return len(received.uniqueLogs) == len(expectedState.uniqueLogs)
	}, 5*time.Second, 10*time.Millisecond, "len(received.logs): %v is not equal len(expectedState.logs): %v", len(received.logs), len(expectedState.logs))

	received.Lock()
	for i := range expectedState.uniqueLogs {
		require.Equal(t, expectedState.uniqueLogs[i], received.uniqueLogs[i])
	}
	received.Unlock()
}

func (listener simpleLogListener) handleLogBroadcast(t *testing.T, lb log.Broadcast) bool {
	t.Helper()
	consumed, err := listener.WasAlreadyConsumed(listener.db, lb)
	require.NoError(t, err)
	if !consumed {

		err = listener.MarkConsumed(listener.db, lb)
		require.NoError(t, err)

		consumed2, err := listener.WasAlreadyConsumed(listener.db, lb)
		require.NoError(t, err)
		require.True(t, consumed2)
	}
	return consumed
}

func (listener simpleLogListener) WasAlreadyConsumed(db *gorm.DB, broadcast log.Broadcast) (bool, error) {
	return log.NewORM(listener.db).WasBroadcastConsumed(db, broadcast.RawLog().BlockHash, broadcast.RawLog().Index, listener.consumerID)
}
func (listener simpleLogListener) MarkConsumed(db *gorm.DB, broadcast log.Broadcast) error {
	return log.NewORM(listener.db).MarkBroadcastConsumed(db, broadcast.RawLog().BlockHash, broadcast.RawLog().BlockNumber, broadcast.RawLog().Index, listener.consumerID)
}

type mockListener struct {
	jobID   models.JobID
	jobIDV2 int32
}

func (l *mockListener) JobID() models.JobID                                     { return l.jobID }
func (l *mockListener) JobIDV2() int32                                          { return l.jobIDV2 }
func (l *mockListener) IsV2Job() bool                                           { return l.jobID.IsZero() }
func (l *mockListener) OnConnect()                                              {}
func (l *mockListener) OnDisconnect()                                           {}
func (l *mockListener) HandleLog(log.Broadcast)                                 {}
func (l *mockListener) WasConsumed(db *gorm.DB, lb log.Broadcast) (bool, error) { return false, nil }
func (l *mockListener) MarkConsumed(db *gorm.DB, lb log.Broadcast) error        { return nil }

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
	checkFilterLogs  func(int64, int64)
}

func (mock *mockEth) assertExpectations(t *testing.T) {
	mock.ethClient.AssertExpectations(t)
	mock.sub.AssertExpectations(t)
}

func (mock *mockEth) subscribeCallCount() int32 {
	return atomic.LoadInt32(&mock.subscribeCalls)
}

func (mock *mockEth) unsubscribeCallCount() int32 {
	return atomic.LoadInt32(&mock.unsubscribeCalls)
}

type mockEthClientExpectedCalls struct {
	SubscribeFilterLogs int
	HeaderByNumber      int
	FilterLogs          int

	FilterLogsResult []types.Log
}

func newMockEthClient(chchRawLogs chan chan<- types.Log, blockHeight int64, expectedCalls mockEthClientExpectedCalls) *mockEth {
	mockEth := &mockEth{
		ethClient:       new(mocks.Client),
		sub:             new(mocks.Subscription),
		checkFilterLogs: nil,
	}
	mockEth.ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			atomic.AddInt32(&(mockEth.subscribeCalls), 1)
			chchRawLogs <- args.Get(2).(chan<- types.Log)
		}).
		Return(mockEth.sub, nil).
		Times(expectedCalls.SubscribeFilterLogs)

	mockEth.ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&models.Head{Number: blockHeight}, nil).
		Times(expectedCalls.HeaderByNumber)

	mockEth.ethClient.On("FilterLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			filterQuery := args.Get(1).(ethereum.FilterQuery)
			fromBlock := filterQuery.FromBlock.Int64()
			toBlock := filterQuery.ToBlock.Int64()
			if mockEth.checkFilterLogs != nil {
				mockEth.checkFilterLogs(fromBlock, toBlock)
			}
		}).
		Return(expectedCalls.FilterLogsResult, nil).
		Times(expectedCalls.FilterLogs)

	mockEth.sub.On("Err").
		Return(nil)

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

func (lb *blocks) logOnBlockNumWithTopics(i uint64, logIndex uint, addr common.Address, topics []common.Hash) types.Log {
	return cltest.RawNewRoundLogWithTopics(lb.t, addr, lb.hashes[i], i, logIndex, false, topics)
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
