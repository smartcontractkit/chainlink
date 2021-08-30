package log_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type broadcasterHelper struct {
	t *testing.T

	lb      log.Broadcaster
	db      *gorm.DB
	mockEth *mockEth
	config  *configtest.TestEVMConfig

	// each received channel corresponds to one eth subscription
	chchRawLogs     chan chan<- types.Log
	getLogPoolCount func() int
	toUnsubscribe   []func()
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int) *broadcasterHelper {
	db := pgtest.NewGormDB(t)
	cfg := cltest.NewTestEVMConfig(t)

	chchRawLogs := make(chan chan<- types.Log, timesSubscribe)

	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: timesSubscribe,
		HeaderByNumber:      1,
		FilterLogs:          1,
	}

	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)

	dborm := log.NewORM(db)
	lb := log.NewBroadcaster(dborm, mockEth.ethClient, cfg, logger.Default, nil)
	cfg.Overrides.EvmFinalityDepth = null.IntFrom(10)
	return &broadcasterHelper{
		t:               t,
		lb:              lb,
		db:              db,
		config:          cfg,
		mockEth:         mockEth,
		chchRawLogs:     chchRawLogs,
		getLogPoolCount: lb.ExportedGetPoolCount(),
		toUnsubscribe:   make([]func(), 0),
	}
}

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient eth.Client, highestSeenHead *models.Head) *broadcasterHelper {
	db := pgtest.NewGormDB(t)
	cfg := cltest.NewTestEVMConfig(t)

	orm := log.NewORM(db)
	lb := log.NewBroadcaster(orm, ethClient, cfg, logger.Default, highestSeenHead)

	return &broadcasterHelper{
		t:               t,
		lb:              lb,
		db:              db,
		config:          cfg,
		getLogPoolCount: lb.ExportedGetPoolCount(),
		toUnsubscribe:   make([]func(), 0),
	}
}

func (helper *broadcasterHelper) newLogListenerWithJob(name string) *simpleLogListener {
	return newLogListenerWithJob(helper.t, helper.db, name)
}

func (helper *broadcasterHelper) start() {
	err := helper.lb.Start()
	require.NoError(helper.t, err)
}

type abigenContract interface {
	Address() common.Address
	ParseLog(log types.Log) (generated.AbigenLog, error)
}

func (helper *broadcasterHelper) register(listener log.Listener, contract abigenContract, numConfirmations uint64) {
	logs := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener, contract, logs, numConfirmations)
}

func (helper *broadcasterHelper) registerWithTopics(listener log.Listener, contract abigenContract, logs []generated.AbigenLog, numConfirmations uint64) {
	logsWithTopics := make(map[common.Hash][][]log.Topic)
	for _, log := range logs {
		logsWithTopics[log.Topic()] = nil
	}
	helper.registerWithTopicValues(listener, contract, numConfirmations, logsWithTopics)
}

func (helper *broadcasterHelper) registerWithTopicValues(listener log.Listener, contract abigenContract, numConfirmations uint64,
	topics map[common.Hash][][]log.Topic) {

	unsubscribe := helper.lb.Register(listener, log.ListenerOpts{
		Contract:         contract.Address(),
		ParseLog:         contract.ParseLog,
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
	name                string
	received            *received
	t                   *testing.T
	db                  *gorm.DB
	jobID               int32
	skipMarkingConsumed bool
}

func newLogListenerWithJob(t *testing.T, db *gorm.DB, name string) *simpleLogListener {
	job := &job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
		ExternalJobID: uuid.NewV4(),
	}
	keyStore := cltest.NewKeyStore(t, db)

	pipelineHelper := cltest.NewJobPipelineV2(t, cltest.NewTestEVMConfig(t), db, nil, keyStore, nil)
	_, err := pipelineHelper.Jrm.CreateJob(context.Background(), job, job.Pipeline)
	require.NoError(t, err)

	var rec received
	return &simpleLogListener{
		db:       db,
		name:     name,
		received: &rec,
		t:        t,
		jobID:    job.ID,
	}
}

func (listener *simpleLogListener) SkipMarkingConsumed(skip bool) {
	listener.skipMarkingConsumed = skip
}

func (listener *simpleLogListener) HandleLog(lb log.Broadcast) {
	listener.received.Lock()
	defer listener.received.Unlock()
	logger.Warnf("Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())

	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handleLogBroadcast(listener.t, lb)

	if !consumed {
		listener.received.uniqueLogs = append(listener.received.uniqueLogs, lb.RawLog())
	} else {
		logger.Warnf("Listener %v: Log was already consumed!", listener.name)
	}
}

func (listener simpleLogListener) JobID() int32 {
	return listener.jobID
}

func (listener *simpleLogListener) getUniqueLogs() []types.Log {
	return listener.received.uniqueLogs
}

func (listener *simpleLogListener) getUniqueLogsBlockNumbers() []uint64 {
	var blockNums []uint64
	for _, uniqueLog := range listener.received.uniqueLogs {
		blockNums = append(blockNums, uniqueLog.BlockNumber)
	}
	return blockNums
}

func (listener *simpleLogListener) requireAllReceived(t *testing.T, expectedState *received) {
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

func (listener *simpleLogListener) handleLogBroadcast(t *testing.T, lb log.Broadcast) bool {
	t.Helper()
	consumed, err := listener.WasAlreadyConsumed(listener.db, lb)
	require.NoError(t, err)
	if !consumed && !listener.skipMarkingConsumed {

		err = listener.MarkConsumed(listener.db, lb)
		require.NoError(t, err)

		consumed2, err := listener.WasAlreadyConsumed(listener.db, lb)
		require.NoError(t, err)
		require.True(t, consumed2)
	}
	return consumed
}

func (listener *simpleLogListener) WasAlreadyConsumed(db *gorm.DB, broadcast log.Broadcast) (bool, error) {
	return log.NewORM(listener.db).WasBroadcastConsumed(db, broadcast.RawLog().BlockHash, broadcast.RawLog().Index, listener.jobID)
}
func (listener *simpleLogListener) MarkConsumed(db *gorm.DB, broadcast log.Broadcast) error {
	return log.NewORM(listener.db).MarkBroadcastConsumed(db, broadcast.RawLog().BlockHash, broadcast.RawLog().BlockNumber, broadcast.RawLog().Index, listener.jobID)
}

type mockListener struct {
	jobID int32
}

func (l *mockListener) JobID() int32                                            { return l.jobID }
func (l *mockListener) HandleLog(log.Broadcast)                                 {}
func (l *mockListener) WasConsumed(db *gorm.DB, lb log.Broadcast) (bool, error) { return false, nil }
func (l *mockListener) MarkConsumed(db *gorm.DB, lb log.Broadcast) error        { return nil }

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

func newMockEthClient(t *testing.T, chchRawLogs chan chan<- types.Log, blockHeight int64, expectedCalls mockEthClientExpectedCalls) *mockEth {
	ethClient, sub := cltest.NewEthClientAndSubMock(t)
	mockEth := &mockEth{
		ethClient:       ethClient,
		sub:             sub,
		checkFilterLogs: nil,
	}
	mockEth.ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			atomic.AddInt32(&(mockEth.subscribeCalls), 1)
			chchRawLogs <- args.Get(2).(chan<- types.Log)
		}).
		Return(mockEth.sub, nil).
		Times(expectedCalls.SubscribeFilterLogs)

	mockEth.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&models.Head{Number: blockHeight}, nil).
		Times(expectedCalls.HeaderByNumber)

	if expectedCalls.FilterLogs > 0 {
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
	}

	mockEth.sub.On("Err").
		Return(nil)

	mockEth.sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { atomic.AddInt32(&(mockEth.unsubscribeCalls), 1) })
	return mockEth
}
