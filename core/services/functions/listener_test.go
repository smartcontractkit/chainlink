package functions_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	sync_mocks "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsListenerUniverse struct {
	service        *functions_service.FunctionsListener
	bridgeAccessor *functions_mocks.BridgeAccessor
	eaClient       *functions_mocks.ExternalAdapterClient
	pluginORM      *functions_mocks.ORM
	logBroadcaster *log_mocks.Broadcaster
	ingressClient  *sync_mocks.TelemetryIngressClient
}

func ptr[T any](t T) *T { return &t }

var (
	RequestID         functions_service.RequestID = newRequestID()
	RequestIDStr      string                      = fmt.Sprintf("0x%x", [32]byte(RequestID))
	SubscriptionOwner common.Address              = common.BigToAddress(big.NewInt(42069))
	SubscriptionID    uint64                      = 5
	ResultBytes       []byte                      = []byte{0xab, 0xcd}
	ErrorBytes        []byte                      = []byte{0xff, 0x11}
	Domains           []string                    = []string{"github.com", "google.com"}
)

func NewFunctionsListenerUniverse(t *testing.T, timeoutSec int, pruneFrequencySec int) *FunctionsListenerUniverse {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	broadcaster.On("AddDependents", 1)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg.Database())
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, KeyStore: kst.Eth(), LogBroadcaster: broadcaster, MailMon: mailMon})
	chain := cc.Chains()[0]
	lggr := logger.TestLogger(t)

	pluginORM := functions_mocks.NewORM(t)
	jsonConfig := job.JSONConfig{
		"requestTimeoutSec":               timeoutSec,
		"requestTimeoutCheckFrequencySec": 1,
		"requestTimeoutBatchLookupSize":   1,
		"listenerEventHandlerTimeoutSec":  1,
		"pruneCheckFrequencySec":          pruneFrequencySec,
	}
	jb := job.Job{
		Type:          job.OffchainReporting2,
		SchemaVersion: 1,
		ExternalJobID: uuid.New(),
		PipelineSpec:  &pipeline.Spec{},
		OCR2OracleSpec: &job.OCR2OracleSpec{
			PluginConfig: jsonConfig,
		},
	}
	eaClient := functions_mocks.NewExternalAdapterClient(t)
	bridgeAccessor := functions_mocks.NewBridgeAccessor(t)

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	require.NoError(t, err)

	oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(common.HexToAddress("0xa"), chain.Client())
	require.NoError(t, err)

	ingressClient := sync_mocks.NewTelemetryIngressClient(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.FunctionsRequests)

	functionsListener := functions_service.NewFunctionsListener(oracleContract, jb, bridgeAccessor, pluginORM, pluginConfig, broadcaster, lggr, mailMon, monEndpoint)

	return &FunctionsListenerUniverse{
		service:        functionsListener,
		bridgeAccessor: bridgeAccessor,
		eaClient:       eaClient,
		pluginORM:      pluginORM,
		logBroadcaster: broadcaster,
		ingressClient:  ingressClient,
	}
}

func PrepareAndStartFunctionsListener(t *testing.T, cbor []byte) (*FunctionsListenerUniverse, *log_mocks.Broadcast, chan struct{}) {
	uni := NewFunctionsListenerUniverse(t, 0, 1_000_000)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)

	log := log_mocks.NewBroadcast(t)
	uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logOracleRequest := ocr2dr_oracle.OCR2DROracleOracleRequest{
		RequestId:          RequestID,
		RequestingContract: common.Address{},
		RequestInitiator:   common.Address{},
		SubscriptionId:     SubscriptionID,
		SubscriptionOwner:  SubscriptionOwner,
		Data:               cbor,
	}
	log.On("DecodedLog").Return(&logOracleRequest)
	log.On("String").Return("")
	return uni, log, make(chan struct{})
}

func TestFunctionsListener_HandleOracleRequestSuccess(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, []byte{})

	uni.pluginORM.On("CreateRequest", RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestDuplicateMarkLogConsumed(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, []byte{})

	uni.pluginORM.On("CreateRequest", RequestID, mock.Anything, mock.Anything, mock.Anything).Return(functions_service.ErrDuplicateRequestID)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_ReportSourceCodeDomains(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, []byte{})

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything).Return(ResultBytes, nil, Domains, nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	var sentMessage []byte
	uni.ingressClient.On("Send", mock.Anything).Return().Run(func(args mock.Arguments) {
		sentMessage = args[0].(synchronization.TelemPayload).Telemetry
	})

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()

	assert.NotEmpty(t, sentMessage)

	var req telem.FunctionsRequest
	err := proto.Unmarshal(sentMessage, &req)
	assert.NoError(t, err)
	assert.Equal(t, RequestIDStr, req.RequestId)
	assert.EqualValues(t, Domains, req.Domains)
}

func TestFunctionsListener_HandleOracleRequestComputationError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, []byte{})

	uni.pluginORM.On("CreateRequest", RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything).Return(nil, ErrorBytes, nil, nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, functions_service.USER_ERROR, ErrorBytes, mock.Anything, true, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestCBORParsingError(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, []byte("invalid cbor"))

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, functions_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		close(doneCh)
	})

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestCBORParsingErrorInvalidFieldType(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	incomingData := &struct {
		Source       string `cbor:"source"`
		CodeLocation string `cbor:"codeLocation"` // incorrect type
	}{
		Source:       "abcd",
		CodeLocation: "inline",
	}
	cborBytes, err := cbor.Marshal(incomingData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]
	uni, log, doneCh := PrepareAndStartFunctionsListener(t, cborBytes)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, mock.Anything, functions_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		close(doneCh)
	})

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestCBORParsingCorrect(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	incomingData := &struct {
		Source             string   `cbor:"source"`
		Language           int      `cbor:"language"`
		Secrets            []byte   `cbor:"secrets"`
		Args               []string `cbor:"args"`
		ExtraUnwantedParam string   `cbor:"extraUnwantedParam"`
		// missing CodeLocation and SecretsLocation
	}{
		Source:             "abcd",
		Language:           3,
		Secrets:            []byte{0xaa, 0xbb},
		Args:               []string{"a", "b"},
		ExtraUnwantedParam: "spam",
	}
	cborBytes, err := cbor.Marshal(incomingData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, cborBytes)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		reqData := args.Get(6).(*functions_service.RequestData)
		assert.Equal(t, incomingData.Source, reqData.Source)
		assert.Equal(t, incomingData.Language, reqData.Language)
		assert.Equal(t, incomingData.Secrets, reqData.Secrets)
		assert.Equal(t, incomingData.Args, reqData.Args)
	}).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_RequestTimeout(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqId := newRequestID()
	doneCh := make(chan bool)
	uni := NewFunctionsListenerUniverse(t, 1, 1_000_000)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	uni.pluginORM.On("TimeoutExpiredResults", mock.Anything, uint32(1), mock.Anything).Return([]functions_service.RequestID{reqId}, nil).Run(func(args mock.Arguments) {
		doneCh <- true
	})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_ORMDoesNotFreezeHandlersForever(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	var ormCallExited sync.WaitGroup
	ormCallExited.Add(1)
	uni, log, _ := PrepareAndStartFunctionsListener(t, []byte{})
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

func TestFunctionsListener_PruneRequests(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	doneCh := make(chan bool)
	uni := NewFunctionsListenerUniverse(t, 0, 1)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	uni.pluginORM.On("PruneOldestRequests", functions_service.DefaultPruneMaxStoredRequests, functions_service.DefaultPruneBatchSize, mock.Anything).Return(uint32(0), uint32(0), nil).Run(func(args mock.Arguments) {
		doneCh <- true
	})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)
	<-doneCh
	uni.service.Close()
}
