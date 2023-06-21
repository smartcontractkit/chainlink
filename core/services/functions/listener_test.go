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
	threshold_mocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold/mocks"
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
	decryptor      *threshold_mocks.Decryptor
}

func ptr[T any](t T) *T { return &t }

var (
	RequestID            functions_service.RequestID = newRequestID()
	RequestIDStr         string                      = fmt.Sprintf("0x%x", [32]byte(RequestID))
	SubscriptionOwner    common.Address              = common.BigToAddress(big.NewInt(42069))
	SubscriptionID       uint64                      = 5
	ResultBytes          []byte                      = []byte{0xab, 0xcd}
	ErrorBytes           []byte                      = []byte{0xff, 0x11}
	Domains              []string                    = []string{"github.com", "google.com"}
	EncryptedSecretsUrls []byte                      = []byte{0x5f, 0x91, 0x82, 0xcb, 0xe8, 0x31, 0x34, 0xb3, 0x2f, 0xad, 0x55, 0x9d, 0xa4, 0x40, 0x8e, 0x8f, 0x02, 0xca, 0xae, 0x7c, 0xdd, 0x1f, 0x60, 0x21, 0x55, 0xab, 0x02, 0x1d, 0x97, 0x41, 0xbd, 0x3b, 0x47, 0xf9, 0xe7, 0x5b, 0x86, 0xd7, 0x08, 0x0b, 0xbe, 0xcf, 0xed, 0xbd, 0xaf, 0x25, 0x58, 0x97, 0x60, 0xfc, 0x03, 0x48, 0x62, 0xed, 0x46, 0x34, 0x4b, 0x05, 0x97, 0xd6, 0x2c, 0x10, 0xc3, 0x42, 0x0a, 0xfa, 0xb4, 0x7b, 0x1f, 0x2e, 0xd4, 0xd7, 0x11, 0x51, 0x34, 0xb1, 0xa3, 0xae, 0xfc, 0x97, 0x7c, 0x73, 0x36, 0x38, 0xef, 0xd6, 0x65, 0xb8, 0x2c, 0x3f, 0x19, 0xfb, 0xb0, 0x5e, 0x36, 0x5f, 0x25, 0x1a, 0x5b, 0x1e, 0xe1, 0x3b, 0x21, 0x5d, 0xe5, 0x6d, 0x7a, 0xd9, 0x97, 0xbe, 0xcb, 0x03, 0xef, 0x5e, 0x49, 0x00, 0x87, 0x92, 0xdd, 0xe1, 0x23, 0x3c, 0x7a, 0x3a, 0xf0, 0x98, 0xce, 0xcc, 0x10, 0x9b, 0x4b, 0x49, 0x5f, 0xf9, 0x2e, 0xd3, 0xc8, 0xca, 0x11, 0xd0, 0x13, 0x8e}
	EncryptedSecrets     []byte                      = []byte(`{"TDH2Ctxt":"eyJHcm91cCI6IlAyNTYiLCJDIjoiTXByNFZQNEY3WHdINjZaY3BBdXFvQUF0QnBLT20zOUJGL0Q2Sk4rNDBiRT0iLCJMYWJlbCI6IkFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUE9IiwiVSI6IkJCNStpRzNsNnBidVdyY0tScTZ3a3JzcitlOU96REpiaVFZNGJNbkJWL2g3Nk0vcVdIUjY0aS9jRytwUjFtNXZvYlNSNWFPUkcxQ3VTbHE0VXNKVkxORT0iLCJVX2JhciI6IkJQcUY3b0cwVVlwWEt6TUJsbCtVWnc1cjBBVDBTVFRJeHRVS0Jva3ErTDlibEJWWGpjbU9jcG5CTXV4dkd1SjZyVURnSk9jRTEzUkZHRXlWQ0VEVmJ4MD0iLCJFIjoick9jWHhLc2FtSVEzWVdRK01qMGRqU01YT1hXMW14SjZ0U0pEdTdqbDdQMD0iLCJGIjoiNkRkd2xSdTZPUzgrd0Z6dDhuUGRhVWpIYWhIcVZzcnMzeTJiUzVmWWVoMD0ifQ==","SymCtxt":"+yHRTEBA+BlJlucXN+o1OoMtQMlZNgzXz66OIcY/1/cpdp0yj9eiMxSrTv0ZhZmnAR2UuB1xemZ+LSFkFsFRIXsu3mOLosfkinmL/CT7pTxHz7DpUV2B/6tTKT7nRqSr+SBD1MUD0tYaFLSinzY36hgwGvZa7R4ikDxnnE/KDi4JsHX7tnLwvZ5kO50FuSHyB+QomlqKqzC8kM8QsLnDyxij10FYS33PITFM5UBEc8nsSnsNhIivLKGBhLI82eIN0nfSd3ChGiTwyD4v+x4/Ktj5+AI+Xdjw9dWdbJJp9xCzuVY3KzFTPvGIFdTASJdn1uXa4iOdVmIbnE1R6PevLZssVKLgm1kiQ8ZC/5VXFMWJ2YfuUXMn938fkwRI9eMTMiXumevKYxghKTChgVT3Nw3Ow6HxX56pEPazsrgYbyHR0PLPlGxiDLCQuyefgW2a9XFtvFcg/6iowdjsPvGN5kSr9X/l/Jz/4ZpvcsJIcCh48Qs7n7Dtoulf3TNPNepndzkzHoVy","Nonce":"kgjHyT3Jar0M155E"}`)
	DecryptedSecrets     []byte                      = []byte(`{"0x0":"lhcKs1pHXQVfsJy/nPGxsgJfZ175O9wxAgCUdIJZ4nPAj3IlGDWcNYJn5OgqZiq0FLmJn6da81gSiMHGGmJ+dsGSIjBAWPMbQ16tZotriXIUj5bY8uaMP+sqsJIdNjX2myMGUxDH7rL2NUaguk1QDlobh4ygxcf3KKWC8YgCCUipU8W4BGrue4JbK9KUwLC7FLyenybS8WFJYDBNK87D4KaJMRu472u4XvqyH3R7EwHr9MwUXjdrQIUOeaIgBSn78RvTnslCnyYJ9bhc8/v4OQLmUV7H7Nkq584x1Zg9uTe2o86uT+ueXpQsZnNk28pheTb0/aYRi95M0jIcyZ76FofqwwaYPHQPuMT3TRClI6DI"}`)
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
	decryptor := threshold_mocks.NewDecryptor(t)

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	require.NoError(t, err)

	oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(common.HexToAddress("0xa"), chain.Client())
	require.NoError(t, err)

	ingressClient := sync_mocks.NewTelemetryIngressClient(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.FunctionsRequests)

	functionsListener := functions_service.NewFunctionsListener(oracleContract, jb, bridgeAccessor, pluginORM, pluginConfig, broadcaster, lggr, mailMon, monEndpoint, decryptor)

	return &FunctionsListenerUniverse{
		service:        functionsListener,
		bridgeAccessor: bridgeAccessor,
		eaClient:       eaClient,
		pluginORM:      pluginORM,
		logBroadcaster: broadcaster,
		ingressClient:  ingressClient,
		decryptor:      decryptor,
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

func TestFunctionsListener_ThresholdDecryptedSecrets(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqData := &struct {
		Source          string   `cbor:"source"`
		Language        int      `cbor:"language"`
		Args            []string `cbor:"args"`
		SecretsLocation int      `cbor:"secretsLocation"`
		Secrets         []byte   `cbor:"secrets"`
	}{
		Source:          "abcd",
		Language:        3,
		Args:            []string{"a", "b"},
		SecretsLocation: 1,
		Secrets:         EncryptedSecretsUrls,
	}
	cborBytes, err := cbor.Marshal(reqData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, cborBytes)

	uni.pluginORM.On("CreateRequest", RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("FetchEncryptedSecrets", mock.Anything, mock.Anything, RequestIDStr, mock.Anything, mock.Anything).Return(EncryptedSecrets, nil, nil)
	uni.decryptor.On("Decrypt", mock.Anything, []byte(RequestIDStr), EncryptedSecrets).Return(DecryptedSecrets, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, mock.Anything, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
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
