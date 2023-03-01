package envvar

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"go.uber.org/zap/zapcore"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ConfigSchema records the schema of configuration at the type level
//
// # A note on Feature Flags
//
// Feature flags should be used during development of large features that might
// span more than one release cycle. Most changes that are not considered "complete"
// when a PR is merged and might affect node operation should be put behind a
// feature flag.
//
// This also allows to disable large parts of the code that may not be needed
// for all deployments that could introduce attack surface area if it is not
// needed.
//
// Good example usage is for alternative blockchain support, new services like
// Feeds Manager, external initiators and so on.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type ConfigSchema struct {
	// ESSENTIAL
	DatabaseURL string `env:"DATABASE_URL"`

	// General/misc
	ChainType                    string          `env:"CHAIN_TYPE"`
	Dev                          bool            `env:"CHAINLINK_DEV" default:"false"`
	ExplorerAccessKey            string          `env:"EXPLORER_ACCESS_KEY"`
	ExplorerSecret               string          `env:"EXPLORER_SECRET"`
	ExplorerURL                  *url.URL        `env:"EXPLORER_URL"`
	FlagsContractAddress         string          `env:"FLAGS_CONTRACT_ADDRESS"`               //nodoc
	InsecureFastScrypt           bool            `env:"INSECURE_FAST_SCRYPT" default:"false"` //nodoc
	ReaperExpiration             models.Duration `env:"REAPER_EXPIRATION" default:"240h"`     //nodoc
	RootDir                      string          `env:"ROOT" default:"~/.chainlink"`
	TelemetryIngressUniConn      bool            `env:"TELEMETRY_INGRESS_UNICONN" default:"true"`
	TelemetryIngressLogging      bool            `env:"TELEMETRY_INGRESS_LOGGING" default:"false"`
	TelemetryIngressServerPubKey string          `env:"TELEMETRY_INGRESS_SERVER_PUB_KEY"`
	TelemetryIngressURL          *url.URL        `env:"TELEMETRY_INGRESS_URL"`
	TelemetryIngressBufferSize   uint            `env:"TELEMETRY_INGRESS_BUFFER_SIZE" default:"100"`
	TelemetryIngressMaxBatchSize uint            `env:"TELEMETRY_INGRESS_MAX_BATCH_SIZE" default:"50"`
	TelemetryIngressSendInterval time.Duration   `env:"TELEMETRY_INGRESS_SEND_INTERVAL" default:"500ms"`
	TelemetryIngressSendTimeout  time.Duration   `env:"TELEMETRY_INGRESS_SEND_TIMEOUT" default:"10s"`
	TelemetryIngressUseBatchSend bool            `env:"TELEMETRY_INGRESS_USE_BATCH_SEND" default:"true"`
	ShutdownGracePeriod          time.Duration   `env:"SHUTDOWN_GRACE_PERIOD" default:"5s"`

	// Audit Logger
	AuditLoggerEnabled        bool   `env:"AUDIT_LOGGER_ENABLED" default:"false"`
	AuditLoggerForwardToUrl   string `env:"AUDIT_LOGGER_FORWARD_TO_URL" default:""`
	AuditLoggerHeaders        string `env:"AUDIT_LOGGER_HEADERS" default:""`
	AuditLoggerJsonWrapperKey string `env:"AUDIT_LOGGER_JSON_WRAPPER_KEY" default:""`

	// Database
	DatabaseListenerMaxReconnectDuration time.Duration `env:"DATABASE_LISTENER_MAX_RECONNECT_DURATION" default:"10m"` //nodoc
	DatabaseListenerMinReconnectInterval time.Duration `env:"DATABASE_LISTENER_MIN_RECONNECT_INTERVAL" default:"1m"`  //nodoc
	MigrateDatabase                      bool          `env:"MIGRATE_DATABASE" default:"true"`
	ORMMaxIdleConns                      int           `env:"ORM_MAX_IDLE_CONNS" default:"10"`
	ORMMaxOpenConns                      int           `env:"ORM_MAX_OPEN_CONNS" default:"20"`
	TriggerFallbackDBPollInterval        time.Duration `env:"TRIGGER_FALLBACK_DB_POLL_INTERVAL" default:"30s"` //nodoc
	// Database Global Lock
	AdvisoryLockCheckInterval time.Duration `env:"ADVISORY_LOCK_CHECK_INTERVAL" default:"1s"`
	AdvisoryLockID            int64         `env:"ADVISORY_LOCK_ID" default:"1027321974924625846"`
	DatabaseLockingMode       string        `env:"DATABASE_LOCKING_MODE" default:"dual"`
	LeaseLockDuration         time.Duration `env:"LEASE_LOCK_DURATION" default:"10s"`
	LeaseLockRefreshInterval  time.Duration `env:"LEASE_LOCK_REFRESH_INTERVAL" default:"1s"`
	// Database Autobackups
	DatabaseBackupDir              string        `env:"DATABASE_BACKUP_DIR"`
	DatabaseBackupFrequency        time.Duration `env:"DATABASE_BACKUP_FREQUENCY" default:"1h"`
	DatabaseBackupMode             string        `env:"DATABASE_BACKUP_MODE" default:"none"`
	DatabaseBackupOnVersionUpgrade bool          `env:"DATABASE_BACKUP_ON_VERSION_UPGRADE" default:"true"`
	DatabaseBackupURL              *url.URL      `env:"DATABASE_BACKUP_URL"`

	// Logging
	JSONConsole       bool           `env:"JSON_CONSOLE" default:"false"`
	LogFileDir        string         `env:"LOG_FILE_DIR"`
	LogLevel          zapcore.Level  `env:"LOG_LEVEL"`
	LogSQL            bool           `env:"LOG_SQL" default:"false"`
	LogFileMaxSize    utils.FileSize `env:"LOG_FILE_MAX_SIZE" default:"5120mb"` // 5120mb was determined based on previously collected logs, in which a daily log would be ~2.5GB and compressed would be ~210MB
	LogFileMaxAge     int64          `env:"LOG_FILE_MAX_AGE" default:"0"`
	LogFileMaxBackups int64          `env:"LOG_FILE_MAX_BACKUPS" default:"1"`
	LogUnixTS         bool           `env:"LOG_UNIX_TS" default:"false"`

	// Web Server
	AllowOrigins                   string          `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	AuthenticatedRateLimit         int64           `env:"AUTHENTICATED_RATE_LIMIT" default:"1000"`
	AuthenticatedRateLimitPeriod   time.Duration   `env:"AUTHENTICATED_RATE_LIMIT_PERIOD" default:"1m"`
	BridgeResponseURL              url.URL         `env:"BRIDGE_RESPONSE_URL"`
	BridgeCacheTTL                 time.Duration   `env:"BRIDGE_CACHE_TTL" default:"0s"`
	HTTPServerWriteTimeout         time.Duration   `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"10s"`
	Port                           uint16          `env:"CHAINLINK_PORT" default:"6688"`
	SecureCookies                  bool            `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout                 models.Duration `env:"SESSION_TIMEOUT" default:"15m"`
	UnAuthenticatedRateLimit       int64           `env:"UNAUTHENTICATED_RATE_LIMIT" default:"5"`
	UnAuthenticatedRateLimitPeriod time.Duration   `env:"UNAUTHENTICATED_RATE_LIMIT_PERIOD" default:"20s"`

	// Web Server MFA
	RPID     string `env:"MFA_RPID"`
	RPOrigin string `env:"MFA_RPORIGIN"`

	// Web Server TLS
	TLSCertPath string `env:"TLS_CERT_PATH"`
	TLSHost     string `env:"CHAINLINK_TLS_HOST"`
	TLSKeyPath  string `env:"TLS_KEY_PATH"`
	TLSPort     uint16 `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect bool   `env:"CHAINLINK_TLS_REDIRECT" default:"false"`

	// Feeds manager
	FeatureFeedsManager bool `env:"FEATURE_FEEDS_MANAGER" default:"true"` //nodoc
	FeatureUICSAKeys    bool `env:"FEATURE_UI_CSA_KEYS" default:"false"`  //nodoc

	// LogPoller
	FeatureLogPoller bool `env:"FEATURE_LOG_POLLER" default:"false"` //nodoc

	// General chains/RPC
	EVMEnabled      bool   `env:"EVM_ENABLED" default:"true"`
	EVMRPCEnabled   bool   `env:"EVM_RPC_ENABLED" default:"true"`
	SolanaEnabled   bool   `env:"SOLANA_ENABLED" default:"false"`
	SolanaNodes     string `env:"SOLANA_NODES"`
	StarknetEnabled bool   `env:"STARKNET_ENABLED" default:"false"`
	StarknetNodes   string `env:"STARKNET_NODES"`

	// EVM/Ethereum
	// Legacy Eth ENV vars
	EthereumHTTPURL       string `env:"ETH_HTTP_URL"`
	EthereumNodes         string `env:"EVM_NODES"`
	EthereumSecondaryURL  string `env:"ETH_SECONDARY_URL"` //nodoc
	EthereumSecondaryURLs string `env:"ETH_SECONDARY_URLS"`
	EthereumURL           string `env:"ETH_URL"`
	// Global
	DefaultChainID *big.Int `env:"ETH_CHAIN_ID"`
	// Per-chain overrides
	BalanceMonitorEnabled             bool          `env:"BALANCE_MONITOR_ENABLED"`
	BlockBackfillDepth                uint64        `env:"BLOCK_BACKFILL_DEPTH" default:"10"`
	BlockBackfillSkip                 bool          `env:"BLOCK_BACKFILL_SKIP" default:"false"`
	BlockEmissionIdleWarningThreshold time.Duration `env:"BLOCK_EMISSION_IDLE_WARNING_THRESHOLD"` //nodoc
	EthTxReaperInterval               time.Duration `env:"ETH_TX_REAPER_INTERVAL"`
	EthTxReaperThreshold              time.Duration `env:"ETH_TX_REAPER_THRESHOLD"`
	EthTxResendAfterThreshold         time.Duration `env:"ETH_TX_RESEND_AFTER_THRESHOLD"`
	EvmFinalityDepth                  uint32        `env:"ETH_FINALITY_DEPTH"`
	EvmHeadTrackerHistoryDepth        uint          `env:"ETH_HEAD_TRACKER_HISTORY_DEPTH"`
	EvmHeadTrackerMaxBufferSize       uint          `env:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE"`
	EvmHeadTrackerSamplingInterval    time.Duration `env:"ETH_HEAD_TRACKER_SAMPLING_INTERVAL"`
	EvmLogBackfillBatchSize           uint32        `env:"ETH_LOG_BACKFILL_BATCH_SIZE"`
	EvmLogPollInterval                time.Duration `env:"ETH_LOG_POLL_INTERVAL"`
	EvmLogKeepBlocksDepth             uint32        `env:"ETH_LOG_KEEP_BLOCKS_DEPTH"`
	EvmRPCDefaultBatchSize            uint32        `env:"ETH_RPC_DEFAULT_BATCH_SIZE"`
	LinkContractAddress               string        `env:"LINK_CONTRACT_ADDRESS"`
	OCR2AutomationGasLimit            uint32        `env:"OCR2_AUTOMATION_GAS_LIMIT"`
	OperatorFactoryAddress            string        `env:"OPERATOR_FACTORY_ADDRESS"`
	MinIncomingConfirmations          uint32        `env:"MIN_INCOMING_CONFIRMATIONS"`
	MinimumContractPayment            assets.Link   `env:"MINIMUM_CONTRACT_PAYMENT_LINK_JUELS"`
	// Node liveness checking
	NodeNoNewHeadsThreshold  time.Duration `env:"NODE_NO_NEW_HEADS_THRESHOLD"`
	NodePollFailureThreshold uint32        `env:"NODE_POLL_FAILURE_THRESHOLD"`
	NodePollInterval         time.Duration `env:"NODE_POLL_INTERVAL"`
	NodeSelectionMode        string        `env:"NODE_SELECTION_MODE"`
	NodeSyncThreshold        uint32        `env:"NODE_SYNC_THRESHOLD"`

	// EVM Gas Controls
	EvmEIP1559DynamicFees bool     `env:"EVM_EIP1559_DYNAMIC_FEES"`
	EvmGasBumpPercent     uint16   `env:"ETH_GAS_BUMP_PERCENT"`
	EvmGasBumpThreshold   uint64   `env:"ETH_GAS_BUMP_THRESHOLD"`
	EvmGasBumpWei         *big.Int `env:"ETH_GAS_BUMP_WEI"`
	EvmGasFeeCapDefault   *big.Int `env:"EVM_GAS_FEE_CAP_DEFAULT"`
	EvmGasLimitDefault    uint32   `env:"ETH_GAS_LIMIT_DEFAULT"`
	EvmGasLimitMax        uint32   `env:"ETH_GAS_LIMIT_MAX"`
	EvmGasLimitMultiplier float32  `env:"ETH_GAS_LIMIT_MULTIPLIER"`
	EvmGasLimitTransfer   uint32   `env:"ETH_GAS_LIMIT_TRANSFER"`
	EvmGasPriceDefault    *big.Int `env:"ETH_GAS_PRICE_DEFAULT"`
	EvmGasTipCapDefault   *big.Int `env:"EVM_GAS_TIP_CAP_DEFAULT"`
	EvmGasTipCapMinimum   *big.Int `env:"EVM_GAS_TIP_CAP_MINIMUM"`
	EvmMaxGasPriceWei     *big.Int `env:"ETH_MAX_GAS_PRICE_WEI"`
	EvmMinGasPriceWei     *big.Int `env:"ETH_MIN_GAS_PRICE_WEI"`
	// Gas limits per job type
	EvmGasLimitOCRJobType    *uint32 `env:"ETH_GAS_LIMIT_OCR_JOB_TYPE"`
	EvmGasLimitDRJobType     *uint32 `env:"ETH_GAS_LIMIT_DR_JOB_TYPE"`
	EvmGasLimitVRFJobType    *uint32 `env:"ETH_GAS_LIMIT_VRF_JOB_TYPE"`
	EvmGasLimitFMJobType     *uint32 `env:"ETH_GAS_LIMIT_FM_JOB_TYPE"`
	EvmGasLimitKeeperJobType *uint32 `env:"ETH_GAS_LIMIT_KEEPER_JOB_TYPE"`
	// Gas Estimation
	GasEstimatorMode                               string `env:"GAS_ESTIMATOR_MODE"`
	BlockHistoryEstimatorBatchSize                 uint32 `env:"BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE"`
	BlockHistoryEstimatorBlockDelay                uint16 `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY"`
	BlockHistoryEstimatorBlockHistorySize          uint16 `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE"`
	BlockHistoryEstimatorCheckInclusionBlocks      uint16 `env:"BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_BLOCKS"`
	BlockHistoryEstimatorCheckInclusionPercentile  uint16 `env:"BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_PERCENTILE"`
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks uint16 `env:"BLOCK_HISTORY_ESTIMATOR_EIP1559_FEE_CAP_BUFFER_BLOCKS"`
	BlockHistoryEstimatorTransactionPercentile     uint16 `env:"BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE"`
	// Txm
	EvmGasBumpTxDepth          uint16 `env:"ETH_GAS_BUMP_TX_DEPTH"`
	EvmMaxInFlightTransactions uint32 `env:"ETH_MAX_IN_FLIGHT_TRANSACTIONS"`
	EvmMaxQueuedTransactions   uint64 `env:"ETH_MAX_QUEUED_TRANSACTIONS"`
	EvmNonceAutoSync           bool   `env:"ETH_NONCE_AUTO_SYNC"`
	EvmUseForwarders           bool   `env:"ETH_USE_FORWARDERS"`

	// Job Pipeline and tasks
	DefaultHTTPLimit                 int64           `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout               models.Duration `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	FeatureExternalInitiators        bool            `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	JobPipelineMaxRunDuration        time.Duration   `env:"JOB_PIPELINE_MAX_RUN_DURATION" default:"10m"`
	JobPipelineMaxSuccessfulRuns     uint64          `env:"JOB_PIPELINE_MAX_SUCCESSFUL_RUNS" default:"10000"`
	JobPipelineReaperInterval        time.Duration   `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold       time.Duration   `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"24h"`
	JobPipelineResultWriteQueueDepth uint64          `env:"JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH" default:"100"`

	// Flux Monitor
	FMDefaultTransactionQueueDepth uint32 `env:"FM_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"` //nodoc
	FMSimulateTransactions         bool   `env:"FM_SIMULATE_TRANSACTIONS" default:"false"`

	// OCR V2
	FeatureOffchainReporting2 bool `env:"FEATURE_OFFCHAIN_REPORTING2" default:"false"` //nodoc
	// Global defaults
	OCR2ContractConfirmations              uint          `env:"OCR2_CONTRACT_CONFIRMATIONS" default:"3"`                  //nodoc
	OCR2BlockchainTimeout                  time.Duration `env:"OCR2_BLOCKCHAIN_TIMEOUT" default:"20s"`                    //nodoc
	OCR2ContractPollInterval               time.Duration `env:"OCR2_CONTRACT_POLL_INTERVAL" default:"1m"`                 //nodoc
	OCR2ContractSubscribeInterval          time.Duration `env:"OCR2_CONTRACT_SUBSCRIBE_INTERVAL" default:"2m"`            //nodoc
	OCR2ContractTransmitterTransmitTimeout time.Duration `env:"OCR2_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT" default:"10s"` //nodoc
	OCR2DatabaseTimeout                    time.Duration `env:"OCR2_DATABASE_TIMEOUT" default:"10s"`                      //nodoc
	OCR2KeyBundleID                        string        `env:"OCR2_KEY_BUNDLE_ID"`                                       //nodoc

	// OCR V1
	FeatureOffchainReporting bool `env:"FEATURE_OFFCHAIN_REPORTING" default:"false"`
	// Per-chain defaults
	OCRContractConfirmations              uint          `env:"OCR_CONTRACT_CONFIRMATIONS"`                //nodoc
	OCRContractTransmitterTransmitTimeout time.Duration `env:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT"` //nodoc
	OCRDatabaseTimeout                    time.Duration `env:"OCR_DATABASE_TIMEOUT"`                      //nodoc
	OCRObservationGracePeriod             time.Duration `env:"OCR_OBSERVATION_GRACE_PERIOD"`              //nodoc
	// Global defaults
	OCRObservationTimeout           time.Duration `env:"OCR_OBSERVATION_TIMEOUT" default:"5s"`            //nodoc
	OCRBlockchainTimeout            time.Duration `env:"OCR_BLOCKCHAIN_TIMEOUT" default:"20s"`            //nodoc
	OCRContractPollInterval         time.Duration `env:"OCR_CONTRACT_POLL_INTERVAL" default:"1m"`         //nodoc
	OCRContractSubscribeInterval    time.Duration `env:"OCR_CONTRACT_SUBSCRIBE_INTERVAL" default:"2m"`    //nodoc
	OCRDefaultTransactionQueueDepth uint32        `env:"OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"` //nodoc
	// Optional
	OCRKeyBundleID          string `env:"OCR_KEY_BUNDLE_ID"`
	OCRSimulateTransactions bool   `env:"OCR_SIMULATE_TRANSACTIONS" default:"false"`
	OCRTraceLogging         bool   `env:"OCR_TRACE_LOGGING" default:"false"` //nodoc
	OCRTransmitterAddress   string `env:"OCR_TRANSMITTER_ADDRESS"`

	// P2P Networking
	// V1 and V2
	P2PNetworkingStack           ocrnetworking.NetworkingStack `env:"P2P_NETWORKING_STACK" default:"V1"`
	P2PIncomingMessageBufferSize int                           `env:"P2P_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"` //nodoc
	P2POutgoingMessageBufferSize int                           `env:"P2P_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"` //nodoc
	// V1 Only
	P2PAnnounceIP                       net.IP        `env:"P2P_ANNOUNCE_IP"`
	P2PAnnouncePort                     uint16        `env:"P2P_ANNOUNCE_PORT"`
	P2PBootstrapCheckInterval           time.Duration `env:"P2P_BOOTSTRAP_CHECK_INTERVAL" default:"20s"` //nodoc
	P2PBootstrapPeers                   []string      `env:"P2P_BOOTSTRAP_PEERS"`
	P2PDHTAnnouncementCounterUserPrefix uint32        `env:"P2P_DHT_ANNOUNCEMENT_COUNTER_USER_PREFIX" default:"0"` //nodoc
	P2PDHTLookupInterval                int           `env:"P2P_DHT_LOOKUP_INTERVAL" default:"10"`                 //nodoc
	P2PListenIP                         net.IP        `env:"P2P_LISTEN_IP" default:"0.0.0.0"`
	P2PListenPort                       uint16        `env:"P2P_LISTEN_PORT"`
	P2PNewStreamTimeout                 time.Duration `env:"P2P_NEW_STREAM_TIMEOUT" default:"10s"`
	P2PPeerID                           p2pkey.PeerID `env:"P2P_PEER_ID"`
	P2PPeerstoreWriteInterval           time.Duration `env:"P2P_PEERSTORE_WRITE_INTERVAL" default:"5m"` //nodoc
	// V2 Only
	P2PV2AnnounceAddresses []string        `env:"P2PV2_ANNOUNCE_ADDRESSES"`
	P2PV2Bootstrappers     []string        `env:"P2PV2_BOOTSTRAPPERS"`
	P2PV2DeltaDial         models.Duration `env:"P2PV2_DELTA_DIAL" default:"15s"`     //nodoc
	P2PV2DeltaReconcile    models.Duration `env:"P2PV2_DELTA_RECONCILE" default:"1m"` //nodoc
	P2PV2ListenAddresses   []string        `env:"P2PV2_LISTEN_ADDRESSES"`
	// DEPRECATED
	OCROutgoingMessageBufferSize int           `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE"` //nodoc
	OCRIncomingMessageBufferSize int           `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE"` //nodoc
	OCRDHTLookupInterval         int           `env:"OCR_DHT_LOOKUP_INTERVAL"`          //nodoc
	OCRBootstrapCheckInterval    time.Duration `env:"OCR_BOOTSTRAP_CHECK_INTERVAL"`     //nodoc
	OCRNewStreamTimeout          time.Duration `env:"OCR_NEW_STREAM_TIMEOUT"`           //nodoc

	// Keeper
	KeeperDefaultTransactionQueueDepth uint32        `env:"KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"` //nodoc
	KeeperGasPriceBufferPercent        uint16        `env:"KEEPER_GAS_PRICE_BUFFER_PERCENT" default:"20"`
	KeeperGasTipCapBufferPercent       uint16        `env:"KEEPER_GAS_TIP_CAP_BUFFER_PERCENT" default:"20"`
	KeeperBaseFeeBufferPercent         uint16        `env:"KEEPER_BASE_FEE_BUFFER_PERCENT" default:"20"`
	KeeperMaximumGracePeriod           int64         `env:"KEEPER_MAXIMUM_GRACE_PERIOD" default:"100"`
	KeeperRegistryCheckGasOverhead     uint64        `env:"KEEPER_REGISTRY_CHECK_GAS_OVERHEAD" default:"200000"`
	KeeperRegistryPerformGasOverhead   uint64        `env:"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD" default:"300000"`
	KeeperRegistryMaxPerformDataSize   uint64        `env:"KEEPER_REGISTRY_MAX_PERFORM_DATA_SIZE" default:"5000"`
	KeeperRegistrySyncInterval         time.Duration `env:"KEEPER_REGISTRY_SYNC_INTERVAL" default:"30m"`
	KeeperRegistrySyncUpkeepQueueSize  uint32        `env:"KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE" default:"10"`
	KeeperTurnLookBack                 int64         `env:"KEEPER_TURN_LOOK_BACK" default:"1000"`

	// Debugging
	AutoPprofEnabled              bool            `env:"AUTO_PPROF_ENABLED" default:"false"`            //nodoc
	AutoPprofProfileRoot          string          `env:"AUTO_PPROF_PROFILE_ROOT"`                       //nodoc (defaults to $CHAINLINK_ROOT)
	AutoPprofPollInterval         models.Duration `env:"AUTO_PPROF_POLL_INTERVAL" default:"10s"`        //nodoc
	AutoPprofGatherDuration       models.Duration `env:"AUTO_PPROF_GATHER_DURATION" default:"10s"`      //nodoc
	AutoPprofGatherTraceDuration  models.Duration `env:"AUTO_PPROF_GATHER_TRACE_DURATION" default:"5s"` //nodoc
	AutoPprofMaxProfileSize       utils.FileSize  `env:"AUTO_PPROF_MAX_PROFILE_SIZE" default:"100mb"`   //nodoc
	AutoPprofCPUProfileRate       int             `env:"AUTO_PPROF_CPU_PROFILE_RATE" default:"1"`       //nodoc
	AutoPprofMemProfileRate       int             `env:"AUTO_PPROF_MEM_PROFILE_RATE" default:"1"`       //nodoc
	AutoPprofBlockProfileRate     int             `env:"AUTO_PPROF_BLOCK_PROFILE_RATE" default:"1"`     //nodoc
	AutoPprofMutexProfileFraction int             `env:"AUTO_PPROF_MUTEX_PROFILE_FRACTION" default:"1"` //nodoc
	AutoPprofMemThreshold         utils.FileSize  `env:"AUTO_PPROF_MEM_THRESHOLD" default:"4gb"`        //nodoc
	AutoPprofGoroutineThreshold   int             `env:"AUTO_PPROF_GOROUTINE_THRESHOLD" default:"5000"` //nodoc

	// Pyroscope (live profiling)
	PyroscopeAuthToken     string `env:"PYROSCOPE_AUTH_TOKEN"`                    //nodoc
	PyroscopeServerAddress string `env:"PYROSCOPE_SERVER_ADDRESS"`                //nodoc
	PyroscopeEnvironment   string `env:"PYROSCOPE_ENVIRONMENT" default:"mainnet"` //nodoc
}

// Name gets the environment variable Name for a config schema field
func Name(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		log.Panicf("Invariant violated, no field of name %s found on ConfigSchema", field)
	}
	return item.Tag.Get("env")
}

// TryName gracefully tries to get the environment variable Name for a config schema field
func TryName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		return fmt.Sprintf("<%s>", field)
	}
	return item.Tag.Get("env")
}

// DefaultValue looks up the default value
func DefaultValue(name string) (string, bool) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		return item.Tag.Lookup("default")
	}
	log.Panicf("Invariant violated, no field of name %s found for DefaultValue", name)
	return "", false
}
