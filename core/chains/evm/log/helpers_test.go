package log_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
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

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type broadcasterHelper struct {
	t            *testing.T
	lb           log.BroadcasterInTest
	db           *sqlx.DB
	mockEth      *evmtest.MockEth
	globalConfig config.BasicConfig
	config       evmconfig.ChainScopedConfig

	// each received channel corresponds to one eth subscription
	chchRawLogs    chan evmtest.RawSub[types.Log]
	toUnsubscribe  []func()
	pipelineHelper cltest.JobPipelineV2TestHelper
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	return broadcasterHelperCfg{}.new(t, blockHeight, timesSubscribe, nil, overridesFn)
}

type broadcasterHelperCfg struct {
	highestSeenHead *evmtypes.Head
	db              *sqlx.DB
}

func (c broadcasterHelperCfg) new(t *testing.T, blockHeight int64, timesSubscribe int, filterLogsResult []types.Log, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	if c.db == nil {
		// ensure we check before registering any mock Cleanup assertions
		testutils.SkipShortDB(t)
	}
	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: timesSubscribe,
		HeaderByNumber:      1,
		FilterLogs:          1,
		FilterLogsResult:    filterLogsResult,
	}

	chchRawLogs := make(chan evmtest.RawSub[types.Log], timesSubscribe)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := c.newWithEthClient(t, mockEth.EthClient, overridesFn)
	helper.chchRawLogs = chchRawLogs
	helper.mockEth = mockEth
	return helper
}

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient evmclient.Client, highestSeenHead *evmtypes.Head, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	return broadcasterHelperCfg{highestSeenHead: highestSeenHead}.newWithEthClient(t, ethClient, overridesFn)
}

func (c broadcasterHelperCfg) newWithEthClient(t *testing.T, ethClient evmclient.Client, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	if c.db == nil {
		c.db = pgtest.NewSqlxDB(t)
	}

	globalConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.LogQueries = ptr(true)
		finality := uint32(10)
		c.EVM[0].FinalityDepth = &finality

		if overridesFn != nil {
			overridesFn(c, s)
		}
	})
	config := evmtest.NewChainScopedConfig(t, globalConfig)
	lggr := logger.TestLogger(t)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	orm := log.NewORM(c.db, lggr, config, cltest.FixtureChainID)
	lb := log.NewTestBroadcaster(orm, ethClient, config, lggr, c.highestSeenHead, mailMon)
	kst := cltest.NewKeyStore(t, c.db, globalConfig)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{
		Client:         ethClient,
		GeneralConfig:  globalConfig,
		DB:             c.db,
		KeyStore:       kst.Eth(),
		LogBroadcaster: &log.NullBroadcaster{},
		MailMon:        mailMon,
	})

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

func (helper *broadcasterHelper) register(listener log.Listener, contract log.AbigenContract, numConfirmations uint32) {
	logs := []generated.AbigenLog{
		flux_aggregator_wrapper.FluxAggregatorNewRound{},
		flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
	}
	helper.registerWithTopics(listener, contract, logs, numConfirmations)
}

func (helper *broadcasterHelper) registerWithTopics(listener log.Listener, contract log.AbigenContract, logs []generated.AbigenLog, numConfirmations uint32) {
	logsWithTopics := make(map[common.Hash][][]log.Topic)
	for _, log := range logs {
		logsWithTopics[log.Topic()] = nil
	}
	helper.registerWithTopicValues(listener, contract, numConfirmations, logsWithTopics)
}

func (helper *broadcasterHelper) registerWithTopicValues(listener log.Listener, contract log.AbigenContract, numConfirmations uint32,
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
	lggr                logger.Logger
	cfg                 pg.QConfig
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
		lggr:     logger.TestLogger(t),
		cfg:      helper.config,
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
	listener.received.Lock()
	defer listener.received.Unlock()
	listener.lggr.Tracef("Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())

	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handleLogBroadcast(lb)

	if !consumed {
		listener.received.uniqueLogs = append(listener.received.uniqueLogs, lb.RawLog())
	} else {
		listener.lggr.Warnf("Listener %v: Log was already consumed!", listener.name)
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
	}, testutils.WaitTimeout(t), time.Second, "len(received.uniqueLogs): %v is not equal len(expectedState.uniqueLogs): %v", len(received.getUniqueLogs()), len(expectedState.getUniqueLogs()))

	received.Lock()
	defer received.Unlock()
	for i, ul := range expectedState.getUniqueLogs() {
		assert.Equal(t, ul, received.uniqueLogs[i])
	}
}

func (listener *simpleLogListener) handleLogBroadcast(lb log.Broadcast) bool {
	t := listener.t
	consumed, err := listener.WasAlreadyConsumed(lb)
	if !assert.NoError(t, err) {
		return false
	}
	if !consumed && !listener.skipMarkingConsumed.Load() {

		err = listener.MarkConsumed(lb)
		if assert.NoError(t, err) {

			consumed2, err := listener.WasAlreadyConsumed(lb)
			if assert.NoError(t, err) {
				assert.True(t, consumed2)
			}
		}
	}
	return consumed
}

func (listener *simpleLogListener) WasAlreadyConsumed(broadcast log.Broadcast) (bool, error) {
	return log.NewORM(listener.db, listener.lggr, listener.cfg, cltest.FixtureChainID).WasBroadcastConsumed(broadcast.RawLog().BlockHash, broadcast.RawLog().Index, listener.jobID)
}

func (listener *simpleLogListener) MarkConsumed(broadcast log.Broadcast) error {
	return log.NewORM(listener.db, listener.lggr, listener.cfg, cltest.FixtureChainID).MarkBroadcastConsumed(broadcast.RawLog().BlockHash, broadcast.RawLog().BlockNumber, broadcast.RawLog().Index, listener.jobID)
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
	ethClient := evmclimocks.NewClient(t)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	mockEth.EthClient.On("ConfiguredChainID", mock.Anything).Return(&cltest.FixtureChainID)
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
