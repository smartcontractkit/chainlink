package directrequestocr_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	log_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	drocr_service "github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	drocr_mocks "github.com/smartcontractkit/chainlink/core/services/directrequestocr/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	job_mocks "github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type DRListenerUniverse struct {
	runner         *pipeline_mocks.Runner
	service        *drocr_service.DRListener
	jobORM         *job_mocks.ORM
	pluginORM      *drocr_mocks.ORM
	logBroadcaster *log_mocks.Broadcaster
}

func ptr[T any](t T) *T { return &t }

func NewDRListenerUniverse(t *testing.T, timeoutSec int) *DRListenerUniverse {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	runner := pipeline_mocks.NewRunner(t)
	broadcaster.On("AddDependents", 1)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	db := pgtest.NewSqlxDB(t)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, LogBroadcaster: broadcaster, MailMon: mailMon})
	chain := cc.Chains()[0]
	lggr := logger.TestLogger(t)

	jobORM := job_mocks.NewORM(t)
	pluginORM := drocr_mocks.NewORM(t)
	jb := &job.Job{
		Type:          job.OffchainReporting2,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
		PipelineSpec:  &pipeline.Spec{},
		OCR2OracleSpec: &job.OCR2OracleSpec{
			PluginConfig: job.JSONConfig{
				"requestTimeoutSec":               timeoutSec,
				"requestTimeoutCheckFrequencySec": 1,
				"requestTimeoutBatchLookupSize":   1,
				"listenerEventHandlerTimeoutSec":  1,
			},
		},
	}

	oracle, err := directrequestocr.NewDROracle(*jb, runner, jobORM, pluginORM, chain, lggr, nil, mailMon)
	require.NoError(t, err)

	serviceArray, err := oracle.GetServices()
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	return &DRListenerUniverse{
		runner:         runner,
		service:        service.(*drocr_service.DRListener),
		jobORM:         jobORM,
		pluginORM:      pluginORM,
		logBroadcaster: broadcaster,
	}
}

func PrepareAndStartDRListener(t *testing.T, expectPipelineRun bool) (*DRListenerUniverse, *log_mocks.Broadcast, cltest.Awaiter) {
	uni := NewDRListenerUniverse(t, 0)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)

	log := log_mocks.NewBroadcast(t)
	log.On("ReceiptsRoot").Return(common.Hash{})
	log.On("TransactionsRoot").Return(common.Hash{})
	log.On("StateRoot").Return(common.Hash{})

	uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logOracleRequest := ocr2dr_oracle.OCR2DROracleOracleRequest{
		RequestId:          RequestID,
		RequestingContract: common.Address{},
		RequestInitiator:   common.Address{},
		SubscriptionId:     0,
		SubscriptionOwner:  common.Address{},
		Data:               []byte("data"),
	}
	log.On("DecodedLog").Return(&logOracleRequest)
	log.On("String").Return("")

	if !expectPipelineRun {
		return uni, log, nil
	}

	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	runBeganAwaiter := cltest.NewAwaiter()
	uni.runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Return(false, nil).
		Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once()

	return uni, log, runBeganAwaiter
}

var RequestID drocr_service.RequestID = newRequestID()

const (
	CorrectResultData string = "\"0x1234\""
	CorrectErrorData  string = "\"0x424144\""
	EmptyData         string = "\"\""
)

func TestDRListener_HandleOracleRequestSuccess(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, runBeganAwaiter := PrepareAndStartDRListener(t, true)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.CBORParseTaskName, mock.Anything).Return([]byte{}, nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.ParseResultTaskName, mock.Anything).Return([]byte(CorrectResultData), nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.ParseErrorTaskName, mock.Anything).Return([]byte(EmptyData), nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, []byte{0x12, 0x34}, mock.Anything, mock.Anything).Return(nil)

	uni.service.HandleLog(log)

	runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
	uni.service.Close()
}

func TestDRListener_HandleOracleRequestComputationError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, runBeganAwaiter := PrepareAndStartDRListener(t, true)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.CBORParseTaskName, mock.Anything).Return([]byte{}, nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.ParseResultTaskName, mock.Anything).Return([]byte(EmptyData), nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.ParseErrorTaskName, mock.Anything).Return([]byte(CorrectErrorData), nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, drocr_service.USER_ERROR, []byte("BAD"), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	uni.service.HandleLog(log)

	runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
	uni.service.Close()
}

func TestDRListener_HandleOracleRequestCBORParsingError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, runBeganAwaiter := PrepareAndStartDRListener(t, true)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, drocr_service.CBORParseTaskName, mock.Anything).Return(nil, errors.New("bad cbor"))
	uni.pluginORM.On("SetError", RequestID, mock.Anything, drocr_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	uni.service.HandleLog(log)

	runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
	uni.service.Close()
}

func TestDRListener_RequestTimeout(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqId := newRequestID()
	done := make(chan bool)
	uni := NewDRListenerUniverse(t, 1)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	uni.pluginORM.On("TimeoutExpiredResults", mock.Anything, uint32(1), mock.Anything).Return([]drocr_service.RequestID{reqId}, nil).Run(func(args mock.Arguments) {
		done <- true
	})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)
	<-done

	uni.service.Close()
}

func TestDRListener_ORMDoesNotFreezeHandlersForever(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	var ormCallExited sync.WaitGroup
	ormCallExited.Add(1)
	uni, log, _ := PrepareAndStartDRListener(t, false)
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		var queryerWrapper pg.Q
		args.Get(3).(pg.QOpt)(&queryerWrapper)
		<-queryerWrapper.ParentCtx.Done()
		ormCallExited.Done()
	}).Return(errors.New("timeout!"))

	uni.service.HandleLog(log)

	ormCallExited.Wait() // should not freeze
	uni.service.Close()
}

func TestDRListener_ExtractRawBytes(t *testing.T) {
	t.Parallel()

	res, err := drocr_service.ExtractRawBytes([]byte("\"\""))
	require.NoError(t, err)
	require.Equal(t, []byte{}, res)

	res, err = drocr_service.ExtractRawBytes([]byte("\"0xabcd\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0xab, 0xcd}, res)

	res, err = drocr_service.ExtractRawBytes([]byte("\"0x0\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0x0}, res)

	res, err = drocr_service.ExtractRawBytes([]byte("\"0x00\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0x0}, res)

	_, err = drocr_service.ExtractRawBytes([]byte(""))
	require.Error(t, err)

	_, err = drocr_service.ExtractRawBytes([]byte("0xab"))
	require.Error(t, err)

	_, err = drocr_service.ExtractRawBytes([]byte("\"0x\""))
	require.Error(t, err)

	_, err = drocr_service.ExtractRawBytes([]byte("\"0xabc\""))
	require.Error(t, err)

	_, err = drocr_service.ExtractRawBytes([]byte("\"0xqwer\""))
	require.Error(t, err)

	_, err = drocr_service.ExtractRawBytes([]byte("null"))
	require.ErrorContains(t, err, "null value")
}
