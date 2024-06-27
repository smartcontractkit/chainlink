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

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"

	log_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/ocr2dr_oracle"
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
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	evmrelay_mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types/mocks"
	s4_mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	sync_mocks "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsListenerUniverse struct {
	service          *functions_service.FunctionsListener
	bridgeAccessor   *functions_mocks.BridgeAccessor
	eaClient         *functions_mocks.ExternalAdapterClient
	pluginORM        *functions_mocks.ORM
	logBroadcaster   *log_mocks.Broadcaster
	ingressClient    *sync_mocks.TelemetryService
	decryptor        *threshold_mocks.Decryptor
	logPollerWrapper *evmrelay_mocks.LogPollerWrapper
	contractVersion  uint32
}

func ptr[T any](t T) *T { return &t }

var (
	RequestID            = newRequestID()
	RequestIDStr         = fmt.Sprintf("0x%x", [32]byte(RequestID))
	SubscriptionOwner    = common.BigToAddress(big.NewInt(42069))
	SubscriptionID       = uint64(5)
	ResultBytes          = []byte{0xab, 0xcd}
	ErrorBytes           = []byte{0xff, 0x11}
	Domains              = []string{"github.com", "google.com"}
	EncryptedSecretsUrls = []byte{0x11, 0x22}
	EncryptedSecrets     = []byte(`{"TDH2Ctxt":"eyJHcm","SymCtxt":"+yHR","Nonce":"kgjHyT3Jar0M155E"}`)
	DecryptedSecrets     = []byte(`{"0x0":"lhcK"}`)
	SignedCBORRequestHex = "a666736f75726365782172657475726e2046756e6374696f6e732e656e636f646555696e743235362831296773656372657473421234686c616e6775616765006c636f64654c6f636174696f6e006f736563726574734c6f636174696f6e0170726571756573745369676e617475726558416fb6d10871aa3865b6620dc5f4594d2a9ad9166ba6b1dbc3f508362fd27aa0461babada48979092a11ecadec9c663a2ea99da4e368408b36a3fb414acfefdd2a1c"
	SubOwnerAddr         = common.HexToAddress("0x2334dE553AB93c69b0ccbe278B6f5E8350Db6204")
	NonSubOwnerAddr      = common.HexToAddress("0x60C9CF55b9de9A956d921A97575108149b758131")
)

func NewFunctionsListenerUniverse(t *testing.T, timeoutSec int, pruneFrequencySec int, setTiers bool, version uint32) *FunctionsListenerUniverse {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	broadcaster.On("AddDependents", 1)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, KeyStore: kst.Eth(), LogBroadcaster: broadcaster, MailMon: mailMon})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)

	chain := legacyChains.Slice()[0]
	lggr := logger.TestLogger(t)

	pluginORM := functions_mocks.NewORM(t)
	jsonConfig := job.JSONConfig{
		"requestTimeoutSec":               timeoutSec,
		"requestTimeoutCheckFrequencySec": 1,
		"requestTimeoutBatchLookupSize":   1,
		"listenerEventHandlerTimeoutSec":  1,
		"pruneCheckFrequencySec":          pruneFrequencySec,
		"decryptionQueueConfig": map[string]interface{}{
			"decryptRequestTimeoutSec": 100,
		},
		"contractVersion":                    version,
		"listenerEventsCheckFrequencyMillis": 100, // only applicable to v1
	}
	if setTiers {
		jsonConfig["maxRequestSizesList"] = []uint32{10, 100, 1_000}
		jsonConfig["maxSecretsSizesList"] = []uint32{10, 100, 200}
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
	require.NoError(t, json.Unmarshal(jsonConfig.Bytes(), &pluginConfig))

	contractAddress := "0xa"

	ingressClient := sync_mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monEndpoint := ingressAgent.GenMonitoringEndpoint(contractAddress, synchronization.FunctionsRequests, "test-network", "test-chainID")

	s4Storage := s4_mocks.NewStorage(t)
	client := chain.Client()
	logPollerWrapper := evmrelay_mocks.NewLogPollerWrapper(t)
	functionsListener := functions_service.NewFunctionsListener(jb, client, contractAddress, bridgeAccessor, pluginORM, pluginConfig, s4Storage, broadcaster, lggr, mailMon, monEndpoint, decryptor, logPollerWrapper)

	return &FunctionsListenerUniverse{
		service:          functionsListener,
		bridgeAccessor:   bridgeAccessor,
		eaClient:         eaClient,
		pluginORM:        pluginORM,
		logBroadcaster:   broadcaster,
		ingressClient:    ingressClient,
		decryptor:        decryptor,
		logPollerWrapper: logPollerWrapper,
		contractVersion:  pluginConfig.ContractVersion,
	}
}

func PrepareAndStartFunctionsListener(t *testing.T, cbor []byte) (*FunctionsListenerUniverse, *log_mocks.Broadcast, chan struct{}) {
	uni := NewFunctionsListenerUniverse(t, 0, 1_000_000, false, 0)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)

	log := log_mocks.NewBroadcast(t)
	uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	logOracleRequest := ocr2dr_oracle.OCR2DROracleOracleRequest{
		RequestId:          RequestID,
		RequestingContract: common.Address{},
		RequestInitiator:   common.Address{},
		SubscriptionId:     uint64(SubscriptionID),
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
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
		SecretsLocation int    `cbor:"secretsLocation"`
		Secrets         []byte `cbor:"secrets"`
	}{
		SecretsLocation: 1,
		Secrets:         EncryptedSecretsUrls,
	}
	cborBytes, err := cbor.Marshal(reqData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, cborBytes)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("FetchEncryptedSecrets", mock.Anything, mock.Anything, RequestIDStr, mock.Anything, mock.Anything).Return(EncryptedSecrets, nil, nil)
	uni.decryptor.On("Decrypt", mock.Anything, decryptionPlugin.CiphertextId(RequestID[:]), EncryptedSecrets).Return(DecryptedSecrets, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, string(DecryptedSecrets), mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	uni.service.HandleLog(log)
	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_ThresholdDecryptedSecretsFailure(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqData := &struct {
		SecretsLocation int    `cbor:"secretsLocation"`
		Secrets         []byte `cbor:"secrets"`
	}{
		SecretsLocation: 1,
		Secrets:         EncryptedSecretsUrls,
	}
	cborBytes, err := cbor.Marshal(reqData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]

	uni, log, doneCh := PrepareAndStartFunctionsListener(t, cborBytes)

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("FetchEncryptedSecrets", mock.Anything, mock.Anything, RequestIDStr, mock.Anything, mock.Anything).Return(EncryptedSecrets, nil, nil)
	uni.decryptor.On("Decrypt", mock.Anything, decryptionPlugin.CiphertextId(RequestID[:]), EncryptedSecrets).Return(nil, errors.New("threshold decryption error"))
	uni.pluginORM.On("SetError", RequestID, functions_service.USER_ERROR, []byte("threshold decryption of secrets failed"), mock.Anything, true, mock.Anything).Run(func(args mock.Arguments) {
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(functions_service.ErrDuplicateRequestID)
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Return(ResultBytes, nil, Domains, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	var sentMessage []byte
	uni.ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[1].([]byte)
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrorBytes, nil, nil)
	uni.pluginORM.On("SetError", RequestID, functions_service.USER_ERROR, ErrorBytes, mock.Anything, true, mock.Anything).Run(func(args mock.Arguments) {
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, functions_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, functions_service.USER_ERROR, []byte("CBOR parsing error"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
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

	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		reqData := args.Get(7).(*functions_service.RequestData)
		assert.Equal(t, incomingData.Source, reqData.Source)
		assert.Equal(t, incomingData.Language, reqData.Language)
		assert.Equal(t, incomingData.Secrets, reqData.Secrets)
		assert.Equal(t, incomingData.Args, reqData.Args)
	}).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
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
	uni := NewFunctionsListenerUniverse(t, 1, 1_000_000, false, 0)
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
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		var queryerWrapper pg.Q
		args.Get(1).(pg.QOpt)(&queryerWrapper)
		<-queryerWrapper.ParentCtx.Done()
		ormCallExited.Done()
	}).Return(errors.New("timeout"))

	uni.service.HandleLog(log)

	ormCallExited.Wait() // should not freeze
	uni.service.Close()
}

func TestFunctionsListener_PruneRequests(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	doneCh := make(chan bool)
	uni := NewFunctionsListenerUniverse(t, 0, 1, false, 0)
	uni.logBroadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {})
	uni.pluginORM.On("PruneOldestRequests", functions_service.DefaultPruneMaxStoredRequests, functions_service.DefaultPruneBatchSize, mock.Anything).Return(uint32(0), uint32(0), nil).Run(func(args mock.Arguments) {
		doneCh <- true
	})

	err := uni.service.Start(testutils.Context(t))
	require.NoError(t, err)
	<-doneCh
	uni.service.Close()
}

func getFlags(requestSizeTier int, secretSizeTier int) [32]byte {
	var flags [32]byte
	flags[1] = byte(requestSizeTier)
	flags[2] = byte(secretSizeTier)
	return flags
}

func TestFunctionsListener_HandleOracleRequestV1_Success(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni := NewFunctionsListenerUniverse(t, 0, 1_000_000, true, 1)
	require.NoError(t, uni.service.Start(testutils.Context(t)))
	doneCh := make(chan struct{})

	request := types.OracleRequest{
		RequestId:         RequestID,
		SubscriptionId:    uint64(SubscriptionID),
		SubscriptionOwner: SubscriptionOwner,
		Flags:             getFlags(1, 0), // tier no 1 of request size, allows up to 100 bytes
		Data:              make([]byte, 12),
	}

	uni.logPollerWrapper.On("LatestEvents").Return([]types.OracleRequest{request}, nil, nil).Once()
	uni.logPollerWrapper.On("LatestEvents").Return(nil, nil, nil)
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestV1_ThresholdDecryptedSecrets(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	reqData := &struct {
		SecretsLocation int    `cbor:"secretsLocation"`
		Secrets         []byte `cbor:"secrets"`
	}{
		SecretsLocation: 1,
		Secrets:         EncryptedSecretsUrls,
	}
	cborBytes, err := cbor.Marshal(reqData)
	require.NoError(t, err)
	// Remove first byte (map header) to make it "diet" CBOR
	cborBytes = cborBytes[1:]
	request := types.OracleRequest{
		RequestId:         RequestID,
		SubscriptionId:    uint64(SubscriptionID),
		SubscriptionOwner: SubscriptionOwner,
		Flags:             getFlags(1, 1), // tiers no 1 of request size and secrets size, allow up to 100 bytes
		Data:              cborBytes,
	}

	uni := NewFunctionsListenerUniverse(t, 0, 1_000_000, true, 1)
	require.NoError(t, uni.service.Start(testutils.Context(t)))
	doneCh := make(chan struct{})

	uni.logPollerWrapper.On("LatestEvents").Return([]types.OracleRequest{request}, nil, nil).Once()
	uni.logPollerWrapper.On("LatestEvents").Return(nil, nil, nil)
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.bridgeAccessor.On("NewExternalAdapterClient").Return(uni.eaClient, nil)
	uni.eaClient.On("FetchEncryptedSecrets", mock.Anything, mock.Anything, RequestIDStr, mock.Anything, mock.Anything).Return(EncryptedSecrets, nil, nil)
	uni.decryptor.On("Decrypt", mock.Anything, decryptionPlugin.CiphertextId(RequestID[:]), EncryptedSecrets).Return(DecryptedSecrets, nil)
	uni.eaClient.On("RunComputation", mock.Anything, RequestIDStr, mock.Anything, SubscriptionOwner.Hex(), SubscriptionID, mock.Anything, mock.Anything, mock.Anything).Return(ResultBytes, nil, nil, nil)
	uni.pluginORM.On("SetResult", RequestID, ResultBytes, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	<-doneCh
	uni.service.Close()
}

func TestFunctionsListener_HandleOracleRequestV1_CBORTooBig(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	uni := NewFunctionsListenerUniverse(t, 0, 1_000_000, true, 1)
	require.NoError(t, uni.service.Start(testutils.Context(t)))
	doneCh := make(chan struct{})

	request := types.OracleRequest{
		RequestId:         RequestID,
		SubscriptionId:    uint64(SubscriptionID),
		SubscriptionOwner: SubscriptionOwner,
		Flags:             getFlags(0, 0), // tier no 0 of request size, allows only for max 10 bytes
		Data:              make([]byte, 20),
	}

	uni.logPollerWrapper.On("LatestEvents").Return([]types.OracleRequest{request}, nil, nil).Once()
	uni.logPollerWrapper.On("LatestEvents").Return(nil, nil, nil)
	uni.pluginORM.On("CreateRequest", mock.Anything, mock.Anything).Return(nil)
	uni.pluginORM.On("SetError", RequestID, functions_service.USER_ERROR, []byte("request too big (max 10 bytes)"), mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(doneCh)
	}).Return(nil)

	<-doneCh
	uni.service.Close()
}
