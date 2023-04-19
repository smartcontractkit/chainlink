package functions_test

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	log_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	functions_service "github.com/smartcontractkit/chainlink/v2/core/services/functions"
	functions_mocks "github.com/smartcontractkit/chainlink/v2/core/services/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	job_mocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsListenerUniverse struct {
	runner         *pipeline_mocks.Runner
	service        *functions_service.FunctionsListener
	jobORM         *job_mocks.ORM
	pluginORM      *functions_mocks.ORM
	logBroadcaster *log_mocks.Broadcaster
}

func ptr[T any](t T) *T { return &t }

func NewFunctionsListenerUniverse(t *testing.T, timeoutSec int) *FunctionsListenerUniverse {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	runner := pipeline_mocks.NewRunner(t)
	broadcaster.On("AddDependents", 1)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, KeyStore: kst.Eth(), LogBroadcaster: broadcaster, MailMon: mailMon})
	chain := cc.Chains()[0]
	lggr := logger.TestLogger(t)

	jobORM := job_mocks.NewORM(t)
	pluginORM := functions_mocks.NewORM(t)
	jsonConfig := job.JSONConfig{
		"requestTimeoutSec":               timeoutSec,
		"requestTimeoutCheckFrequencySec": 1,
		"requestTimeoutBatchLookupSize":   1,
		"listenerEventHandlerTimeoutSec":  1,
	}
	jb := job.Job{
		Type:          job.OffchainReporting2,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
		PipelineSpec:  &pipeline.Spec{},
		OCR2OracleSpec: &job.OCR2OracleSpec{
			PluginConfig: jsonConfig,
		},
	}

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	require.NoError(t, err)

	oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(common.HexToAddress("0x0"), chain.Client())
	require.NoError(t, err)
	functionsListener := functions_service.NewFunctionsListener(oracleContract, jb, runner, jobORM, pluginORM, pluginConfig, broadcaster, lggr, mailMon)

	return &FunctionsListenerUniverse{
		runner:         runner,
		service:        functionsListener,
		jobORM:         jobORM,
		pluginORM:      pluginORM,
		logBroadcaster: broadcaster,
	}
}

func PrepareAndStartFunctionsListener(t *testing.T, cbor []byte, expectPipelineRun bool) (*FunctionsListenerUniverse, *log_mocks.Broadcast, cltest.Awaiter) {
	uni := NewFunctionsListenerUniverse(t, 0)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)

	log := log_mocks.NewBroadcast(t)
	uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logOracleRequest := ocr2dr_oracle.OCR2DROracleOracleRequest{
		RequestId:          RequestID,
		RequestingContract: common.Address{},
		RequestInitiator:   common.Address{},
		SubscriptionId:     0,
		SubscriptionOwner:  common.Address{},
		Data:               cbor,
	}
	log.On("DecodedLog").Return(&logOracleRequest)
	log.On("String").Return("")

	if !expectPipelineRun {
		return uni, log, nil
	}

	runBeganAwaiter := cltest.NewAwaiter()
	uni.runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Return(false, nil).
		Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once()

	return uni, log, runBeganAwaiter
}

var RequestID functions_service.RequestID = newRequestID()

const (
	CorrectResultData string = "\"0x1234\""
	CorrectErrorData  string = "\"0x424144\""
	EmptyData         string = "\"\""
)

func TestFunctionsListener_HandleOracleRequestSuccess(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, runBeganAwaiter := PrepareAndStartFunctionsListener(t, []byte{}, true)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, functions_service.ParseResultTaskName, mock.Anything).Return([]byte(CorrectResultData), nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, functions_service.ParseErrorTaskName, mock.Anything).Return([]byte(EmptyData), nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, []byte{0x12, 0x34}, mock.Anything, mock.Anything).Return(nil)

	uni.service.HandleLog(log)

	runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestComputationError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, runBeganAwaiter := PrepareAndStartFunctionsListener(t, []byte{}, true)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, functions_service.ParseResultTaskName, mock.Anything).Return([]byte(EmptyData), nil)
	uni.jobORM.On("FindTaskResultByRunIDAndTaskName", mock.Anything, functions_service.ParseErrorTaskName, mock.Anything).Return([]byte(CorrectErrorData), nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, functions_service.USER_ERROR, []byte("BAD"), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	uni.service.HandleLog(log)

	runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestCBORParsingError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, _ := PrepareAndStartFunctionsListener(t, []byte("invalid cbor"), false)

	done := make(chan bool)
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, functions_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done <- true
	})

	uni.service.HandleLog(log)
	<-done
	uni.service.Close()
}

func TestFunctionsListener_RequestTimeout(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqId := newRequestID()
	done := make(chan bool)
	uni := NewFunctionsListenerUniverse(t, 1)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	uni.pluginORM.On("TimeoutExpiredResults", mock.Anything, uint32(1), mock.Anything).Return([]functions_service.RequestID{reqId}, nil).Run(func(args mock.Arguments) {
		done <- true
	})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)
	<-done

	uni.service.Close()
}

func TestFunctionsListener_ORMDoesNotFreezeHandlersForever(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	var ormCallExited sync.WaitGroup
	ormCallExited.Add(1)
	uni, log, _ := PrepareAndStartFunctionsListener(t, []byte{}, false)
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

func TestFunctionsListener_ExtractRawBytes(t *testing.T) {
	t.Parallel()

	res, err := functions_service.ExtractRawBytes([]byte("\"\""))
	require.NoError(t, err)
	require.Equal(t, []byte{}, res)

	res, err = functions_service.ExtractRawBytes([]byte("\"0xabcd\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0xab, 0xcd}, res)

	res, err = functions_service.ExtractRawBytes([]byte("\"0x0\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0x0}, res)

	res, err = functions_service.ExtractRawBytes([]byte("\"0x00\""))
	require.NoError(t, err)
	require.Equal(t, []byte{0x0}, res)

	_, err = functions_service.ExtractRawBytes([]byte(""))
	require.Error(t, err)

	_, err = functions_service.ExtractRawBytes([]byte("0xab"))
	require.Error(t, err)

	_, err = functions_service.ExtractRawBytes([]byte("\"0x\""))
	require.Error(t, err)

	_, err = functions_service.ExtractRawBytes([]byte("\"0xabc\""))
	require.Error(t, err)

	_, err = functions_service.ExtractRawBytes([]byte("\"0xqwer\""))
	require.Error(t, err)

	_, err = functions_service.ExtractRawBytes([]byte("null"))
	require.ErrorContains(t, err, "null value")
}
