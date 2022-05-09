package log_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type broadcasterHelper struct {
	t            *testing.T
	lb           log.BroadcasterInTest
	db           *sqlx.DB
	mockEth      *evmtest.MockEth
	globalConfig *configtest.TestGeneralConfig
	config       evmconfig.ChainScopedConfig

	// each received channel corresponds to one eth subscription
	chchRawLogs    chan evmtest.RawSub[types.Log]
	toUnsubscribe  []func()
	pipelineHelper cltest.JobPipelineV2TestHelper
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int) *broadcasterHelper {
	return broadcasterHelperCfg{}.new(t, blockHeight, timesSubscribe, nil)
}

type broadcasterHelperCfg struct {
	highestSeenHead *evmtypes.Head
	db              *sqlx.DB
}

func (c broadcasterHelperCfg) new(t *testing.T, blockHeight int64, timesSubscribe int, filterLogsResult []types.Log) *broadcasterHelper {
	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: timesSubscribe,
		HeaderByNumber:      1,
		FilterLogs:          1,
		FilterLogsResult:    filterLogsResult,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], timesSubscribe)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := c.newWithEthClient(t, mockEth.EthClient)
	helper.chchRawLogs = chchRawLogs
	helper.mockEth = mockEth
	helper.globalConfig.Overrides.GlobalEvmFinalityDepth = null.IntFrom(10)
	return helper
}

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient evmclient.Client, highestSeenHead *evmtypes.Head) *broadcasterHelper {
	return broadcasterHelperCfg{highestSeenHead: highestSeenHead}.newWithEthClient(t, ethClient)
}

func (c broadcasterHelperCfg) newWithEthClient(t *testing.T, ethClient evmclient.Client) *broadcasterHelper {
	if c.db == nil {
		c.db = pgtest.NewSqlxDB(t)
	}

	globalConfig := cltest.NewTestGeneralConfig(t)
	globalConfig.Overrides.LogSQL = null.BoolFrom(true)
	config := evmtest.NewChainScopedConfig(t, globalConfig)
	lggr := logger.TestLogger(t)

	orm := log.NewORM(c.db, lggr, config, cltest.FixtureChainID)
	lb := log.NewTestBroadcaster(orm, ethClient, config, lggr, c.highestSeenHead)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{
		Client:         ethClient,
		GeneralConfig:  config,
		DB:             c.db,
		LogBroadcaster: &log.NullBroadcaster{},
	})
	kst := cltest.NewKeyStore(t, c.db, globalConfig)
	pipelineHelper := cltest.NewJobPipelineV2(t, config, cc, c.db, kst, nil, nil)

	return &broadcasterHelper{
		t:              t,
		lb:             lb,
		db:             c.db,
		globalConfig:   globalConfig,
		config:         config,
		pipelineHelper: pipelineHelper,
		toUnsubscribe:  make([]func(), 0),
	}
}

func (helper *broadcasterHelper) start() {
	err := helper.lb.Start(testutils.Context(helper.t))
	require.NoError(helper.t, err)
}

type abigenContract interface {
	Address() common.Address
	ParseLog(log types.Log) (generated.AbigenLog, error)
}

func (helper *broadcasterHelper) register(listener log.Listener, contract abigenContract, numConfirmations uint32) {
	logs := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener, contract, logs, numConfirmations)
}

func (helper *broadcasterHelper) registerWithTopics(listener log.Listener, contract abigenContract, logs []generated.AbigenLog, numConfirmations uint32) {
	logsWithTopics := make(map[common.Hash][][]log.Topic)
	for _, log := range logs {
		logsWithTopics[log.Topic()] = nil
	}
	helper.registerWithTopicValues(listener, contract, numConfirmations, logsWithTopics)
}

func (helper *broadcasterHelper) registerWithTopicValues(listener log.Listener, contract abigenContract, numConfirmations uint32,
	topics map[common.Hash][][]log.Topic) {

	unsubscribe := helper.lb.Register(listener, log.ListenerOpts{
		Contract:                 contract.Address(),
		ParseLog:                 contract.ParseLog,
		LogsWithTopics:           topics,
		MinIncomingConfirmations: numConfirmations,
	})

	helper.toUnsubscribe = append(helper.toUnsubscribe, unsubscribe)
}

func (helper *broadcasterHelper) requireBroadcastCount(expectedCount int) {
	helper.t.Helper()
	g := gomega.NewGomegaWithT(helper.t)

	comparisonFunc := func() (int, error) {
		var count struct{ Count int }
		err := helper.db.Get(&count, `SELECT count(*) FROM log_broadcasts`)
		return count.Count, err
	}

	g.Eventually(comparisonFunc, testutils.WaitTimeout(helper.t), time.Second).Should(gomega.Equal(expectedCount))
	g.Consistently(comparisonFunc, 1*time.Second, 200*time.Millisecond).Should(gomega.Equal(expectedCount))
}

func (helper *broadcasterHelper) unsubscribeAll() {
	for _, unsubscribe := range helper.toUnsubscribe {
		unsubscribe()
	}
	time.Sleep(100 * time.Millisecond)
}
func (helper *broadcasterHelper) stop() {
	err := helper.lb.Close()
	assert.NoError(helper.t, err)
}

func newMockContract() *logmocks.AbigenContract {
	addr := testutils.NewAddress()
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

func (rec *received) getLogs() []types.Log {
	rec.Lock()
	defer rec.Unlock()
	r := make([]types.Log, len(rec.logs))
	copy(r, rec.logs)
	return r
}

func (rec *received) getUniqueLogs() []types.Log {
	rec.Lock()
	defer rec.Unlock()
	r := make([]types.Log, len(rec.uniqueLogs))
	copy(r, rec.uniqueLogs)
	return r
}

func (rec *received) logsOnBlocks() []logOnBlock {
	rec.Lock()
	defer rec.Unlock()
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
	db                  *sqlx.DB
	jobID               int32
	skipMarkingConsumed atomic.Bool
}

func (helper *broadcasterHelper) newLogListenerWithJob(name string) *simpleLogListener {
	t := helper.t
	db := helper.db
	jb := &job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
		ExternalJobID: uuid.NewV4(),
	}
	err := helper.pipelineHelper.Jrm.CreateJob(jb)
	require.NoError(t, err)

	var rec received
	return &simpleLogListener{
		db:       db,
		name:     name,
		received: &rec,
		t:        t,
		jobID:    jb.ID,
	}
}

func (listener *simpleLogListener) SkipMarkingConsumed(skip bool) {
	listener.skipMarkingConsumed.Store(skip)
}

func (listener *simpleLogListener) HandleLog(lb log.Broadcast) {
	lggr := logger.TestLogger(listener.t)
	cfg := cltest.NewTestGeneralConfig(listener.t)
	listener.received.Lock()
	defer listener.received.Unlock()
	lggr.Tracef("Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())

	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handleLogBroadcast(listener.t, lggr, cfg, lb)

	if !consumed {
		listener.received.uniqueLogs = append(listener.received.uniqueLogs, lb.RawLog())
	} else {
		lggr.Warnf("Listener %v: Log was already consumed!", listener.name)
	}
}

func (listener *simpleLogListener) JobID() int32 {
	return listener.jobID
}

func (listener *simpleLogListener) getUniqueLogs() []types.Log {
	return listener.received.getUniqueLogs()
}

func (listener *simpleLogListener) getUniqueLogsBlockNumbers() []uint64 {
	var blockNums []uint64
	for _, uniqueLog := range listener.received.getUniqueLogs() {
		blockNums = append(blockNums, uniqueLog.BlockNumber)
	}
	return blockNums
}

func (listener *simpleLogListener) requireAllReceived(t *testing.T, expectedState *received) {
	received := listener.received
	require.Eventually(t, func() bool {
		return len(received.getUniqueLogs()) == len(expectedState.getUniqueLogs())
	}, cltest.WaitTimeout(t), time.Second, "len(received.uniqueLogs): %v is not equal len(expectedState.uniqueLogs): %v", len(received.getUniqueLogs()), len(expectedState.getUniqueLogs()))

	received.Lock()
	defer received.Unlock()
	for i, ul := range expectedState.getUniqueLogs() {
		assert.Equal(t, ul, received.uniqueLogs[i])
	}
}

func (listener *simpleLogListener) handleLogBroadcast(t *testing.T, lggr logger.Logger, cfg pg.LogConfig, lb log.Broadcast) bool {
	consumed, err := listener.WasAlreadyConsumed(listener.db, lggr, cfg, lb)
	if !assert.NoError(t, err) {
		return false
	}
	if !consumed && !listener.skipMarkingConsumed.Load() {

		err = listener.MarkConsumed(listener.db, lggr, cfg, lb)
		if assert.NoError(t, err) {

			consumed2, err := listener.WasAlreadyConsumed(listener.db, lggr, cfg, lb)
			if assert.NoError(t, err) {
				assert.True(t, consumed2)
			}
		}
	}
	return consumed
}

func (listener *simpleLogListener) WasAlreadyConsumed(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig, broadcast log.Broadcast) (bool, error) {
	return log.NewORM(listener.db, lggr, cfg, cltest.FixtureChainID).WasBroadcastConsumed(broadcast.RawLog().BlockHash, broadcast.RawLog().Index, listener.jobID)
}

func (listener *simpleLogListener) MarkConsumed(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig, broadcast log.Broadcast) error {
	return log.NewORM(listener.db, lggr, cfg, cltest.FixtureChainID).MarkBroadcastConsumed(broadcast.RawLog().BlockHash, broadcast.RawLog().BlockNumber, broadcast.RawLog().Index, listener.jobID)
}

type mockListener struct {
	jobID int32
}

func (l *mockListener) JobID() int32            { return l.jobID }
func (l *mockListener) HandleLog(log.Broadcast) {}

type mockEthClientExpectedCalls struct {
	SubscribeFilterLogs int
	HeaderByNumber      int
	FilterLogs          int

	FilterLogsResult []types.Log
}

func newMockEthClient(t *testing.T, chchRawLogs chan<- evmtest.RawSub[types.Log], blockHeight int64, expectedCalls mockEthClientExpectedCalls) *evmtest.MockEth {
	ethClient := new(evmmocks.Client)
	ethClient.Test(t)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	mockEth.EthClient.On("ChainID", mock.Anything).Return(&cltest.FixtureChainID)
	mockEth.EthClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
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
		Times(expectedCalls.SubscribeFilterLogs)

	mockEth.EthClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&evmtypes.Head{Number: blockHeight}, nil).
		Times(expectedCalls.HeaderByNumber)

	if expectedCalls.FilterLogs > 0 {
		mockEth.EthClient.On("FilterLogs", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				filterQuery := args.Get(1).(ethereum.FilterQuery)
				fromBlock := filterQuery.FromBlock.Int64()
				toBlock := filterQuery.ToBlock.Int64()
				if mockEth.CheckFilterLogs != nil {
					mockEth.CheckFilterLogs(fromBlock, toBlock)
				}
			}).
			Return(expectedCalls.FilterLogsResult, nil).
			Times(expectedCalls.FilterLogs)
	}

	return mockEth
}
