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
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	evmtestutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type broadcasterHelper struct {
	t            *testing.T
	lb           log.BroadcasterInTest
	db           *sqlx.DB
	mockEth      *evmtestutils.MockEth
	globalConfig config.AppConfig
	config       evmconfig.ChainScopedConfig

	// each received channel corresponds to one eth subscription
	chchRawLogs    chan evmtestutils.RawSub[types.Log]
	toUnsubscribe  []func()
	pipelineHelper cltest.JobPipelineV2TestHelper
}

func newBroadcasterHelper(t *testing.T, blockHeight int64, timesSubscribe int, filterLogsResult []types.Log, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	// ensure we check before registering any mock Cleanup assertions
	testutils.SkipShortDB(t)

	expectedCalls := mockEthClientExpectedCalls{
		SubscribeFilterLogs: timesSubscribe,
		HeaderByNumber:      1,
		FilterLogs:          1,
		FilterLogsResult:    filterLogsResult,
	}

	chchRawLogs := make(chan evmtestutils.RawSub[types.Log], timesSubscribe)
	mockEth := newMockEthClient(t, chchRawLogs, blockHeight, expectedCalls)
	helper := newBroadcasterHelperWithEthClient(t, mockEth.EthClient, nil, overridesFn)
	helper.chchRawLogs = chchRawLogs
	helper.mockEth = mockEth
	return helper
}

func newBroadcasterHelperWithEthClient(t *testing.T, ethClient evmclient.Client, highestSeenHead *evmtypes.Head, overridesFn func(*chainlink.Config, *chainlink.Secrets)) *broadcasterHelper {
	globalConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.LogQueries = ptr(true)
		finality := uint32(10)
		c.EVM[0].FinalityDepth = &finality

		if overridesFn != nil {
			overridesFn(c, s)
		}
	})
	config := evmtest.NewChainScopedConfig(t, globalConfig)
	lggr := logger.Test(t)
	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	db := pgtest.NewSqlxDB(t)
	orm := log.NewORM(db, cltest.FixtureChainID)
	lb := log.NewTestBroadcaster(orm, ethClient, config.EVM(), lggr, highestSeenHead, mailMon)
	kst := cltest.NewKeyStore(t, db)

	cc := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{
		Client:         ethClient,
		GeneralConfig:  globalConfig,
		DB:             db,
		KeyStore:       kst.Eth(),
		LogBroadcaster: &log.NullBroadcaster{},
		MailMon:        mailMon,
	})

	m := make(map[string]legacyevm.Chain)
	for _, r := range cc.Slice() {
		m[r.Chain().ID().String()] = r.Chain()
	}
	legacyChains := legacyevm.NewLegacyChains(m, cc.AppConfig().EVMConfigs())
	pipelineHelper := cltest.NewJobPipelineV2(t, globalConfig.WebServer(), globalConfig.JobPipeline(), legacyChains, db, kst, nil, nil)

	return &broadcasterHelper{
		t:              t,
		lb:             lb,
		db:             db,
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

func newMockContract(t *testing.T) *logmocks.AbigenContract {
	addr := testutils.NewAddress()
	contract := logmocks.NewAbigenContract(t)
	contract.On("Address").Return(addr).Maybe()
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
	lggr                logger.SugaredLogger
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
		ExternalJobID: uuid.New(),
	}
	err := helper.pipelineHelper.Jrm.CreateJob(testutils.Context(t), jb)
	require.NoError(t, err)

	var rec received
	return &simpleLogListener{
		db:       db,
		lggr:     logger.Sugared(logger.Test(t)),
		name:     name,
		received: &rec,
		t:        t,
		jobID:    jb.ID,
	}
}

func (listener *simpleLogListener) SkipMarkingConsumed(skip bool) {
	listener.skipMarkingConsumed.Store(skip)
}

func (listener *simpleLogListener) HandleLog(ctx context.Context, lb log.Broadcast) {
	listener.received.Lock()
	defer listener.received.Unlock()
	listener.lggr.Tracef("Listener %v HandleLog for block %v %v received at %v %v", listener.name, lb.RawLog().BlockNumber, lb.RawLog().BlockHash, lb.LatestBlockNumber(), lb.LatestBlockHash())

	listener.received.logs = append(listener.received.logs, lb.RawLog())
	listener.received.broadcasts = append(listener.received.broadcasts, lb)
	consumed := listener.handleLogBroadcast(ctx, lb)

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
	defer func() { assert.EqualValues(t, expectedState.getUniqueLogs(), received.getUniqueLogs()) }()
	require.Eventually(t, func() bool {
		return len(received.getUniqueLogs()) == len(expectedState.getUniqueLogs())
	}, testutils.WaitTimeout(t), time.Second, "len(received.uniqueLogs): %v is not equal len(expectedState.uniqueLogs): %v", len(received.getUniqueLogs()), len(expectedState.getUniqueLogs()))
}

func (listener *simpleLogListener) handleLogBroadcast(ctx context.Context, lb log.Broadcast) bool {
	t := listener.t
	consumed, err := listener.WasAlreadyConsumed(ctx, lb)
	if !assert.NoError(t, err) {
		return false
	}
	if !consumed && !listener.skipMarkingConsumed.Load() {
		err = listener.MarkConsumed(ctx, lb)
		if assert.NoError(t, err) {
			consumed2, err := listener.WasAlreadyConsumed(ctx, lb)
			if assert.NoError(t, err) {
				assert.True(t, consumed2)
			}
		}
	}
	return consumed
}

func (listener *simpleLogListener) WasAlreadyConsumed(ctx context.Context, broadcast log.Broadcast) (bool, error) {
	return log.NewORM(listener.db, cltest.FixtureChainID).WasBroadcastConsumed(ctx, broadcast.RawLog().BlockHash, broadcast.RawLog().Index, listener.jobID)
}

func (listener *simpleLogListener) MarkConsumed(ctx context.Context, broadcast log.Broadcast) error {
	return log.NewORM(listener.db, cltest.FixtureChainID).MarkBroadcastConsumed(ctx, broadcast.RawLog().BlockHash, broadcast.RawLog().BlockNumber, broadcast.RawLog().Index, listener.jobID)
}

type mockListener struct {
	jobID int32
}

func (l *mockListener) JobID() int32                             { return l.jobID }
func (l *mockListener) HandleLog(context.Context, log.Broadcast) {}

type mockEthClientExpectedCalls struct {
	SubscribeFilterLogs int
	HeaderByNumber      int
	FilterLogs          int

	FilterLogsResult []types.Log
}

func newMockEthClient(t *testing.T, chchRawLogs chan<- evmtestutils.RawSub[types.Log], blockHeight int64, expectedCalls mockEthClientExpectedCalls) *evmtestutils.MockEth {
	ethClient := evmclimocks.NewClient(t)
	mockEth := &evmtestutils.MockEth{EthClient: ethClient}
	mockEth.EthClient.On("ConfiguredChainID", mock.Anything).Return(&cltest.FixtureChainID)
	mockEth.EthClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchRawLogs <- evmtestutils.NewRawSub(ch, sub.Err())
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
