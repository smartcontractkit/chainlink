package envvar

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ConfigSchema records the schema of configuration at the type level
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
	TelemetryIngressLogging      bool            `env:"TELEMETRY_INGRESS_LOGGING" default:"false"`
	TelemetryIngressServerPubKey string          `env:"TELEMETRY_INGRESS_SERVER_PUB_KEY"`
	TelemetryIngressURL          *url.URL        `env:"TELEMETRY_INGRESS_URL"`
	ShutdownGracePeriod          time.Duration   `env:"SHUTDOWN_GRACE_PERIOD" default:"15s"`

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
	DatabaseBackupDir       string        `env:"DATABASE_BACKUP_DIR"`
	DatabaseBackupFrequency time.Duration `env:"DATABASE_BACKUP_FREQUENCY" default:"1h"`
	DatabaseBackupMode      string        `env:"DATABASE_BACKUP_MODE" default:"none"`
	DatabaseBackupURL       *url.URL      `env:"DATABASE_BACKUP_URL"`

	// Logging
	JSONConsole bool          `env:"JSON_CONSOLE" default:"false"`
	LogFileDir  string        `env:"LOG_FILE_DIR"`
	LogLevel    zapcore.Level `env:"LOG_LEVEL"`
	LogSQL      bool          `env:"LOG_SQL" default:"false"`
	LogToDisk   bool          `env:"LOG_TO_DISK" default:"false"`
	LogUnixTS   bool          `env:"LOG_UNIX_TS" default:"false"`

	// Web Server
	AllowOrigins                   string          `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	AuthenticatedRateLimit         int64           `env:"AUTHENTICATED_RATE_LIMIT" default:"1000"`
	AuthenticatedRateLimitPeriod   time.Duration   `env:"AUTHENTICATED_RATE_LIMIT_PERIOD" default:"1m"`
	BridgeResponseURL              url.URL         `env:"BRIDGE_RESPONSE_URL"`
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
	TLSCertPath string `env:"TLS_CERT_PATH" `
	TLSHost     string `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath  string `env:"TLS_KEY_PATH" `
	TLSPort     uint16 `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect bool   `env:"CHAINLINK_TLS_REDIRECT" default:"false"`

	// Feeds manager
	FeatureFeedsManager   bool `env:"FEATURE_FEEDS_MANAGER" default:"false"`    //nodoc
	FeatureUICSAKeys      bool `env:"FEATURE_UI_CSA_KEYS" default:"false"`      //nodoc
	FeatureUIFeedsManager bool `env:"FEATURE_UI_FEEDS_MANAGER" default:"false"` //nodoc

	// EVM/Ethereum
	// Legacy Eth ENV vars
	EthereumHTTPURL       string `env:"ETH_HTTP_URL"`
	EthereumSecondaryURL  string `env:"ETH_SECONDARY_URL"` //nodoc
	EthereumSecondaryURLs string `env:"ETH_SECONDARY_URLS"`
	EthereumURL           string `env:"ETH_URL"`
	UseLegacyEthEnvVars   bool   `env:"USE_LEGACY_ETH_ENV_VARS" default:"true"`
	// Global
	DefaultChainID   *big.Int `env:"ETH_CHAIN_ID"`
	EVMDisabled      bool     `env:"EVM_DISABLED" default:"false"`
	EthereumDisabled bool     `env:"ETH_DISABLED" default:"false"`
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
	EvmRPCDefaultBatchSize            uint32        `env:"ETH_RPC_DEFAULT_BATCH_SIZE"`
	LinkContractAddress               string        `env:"LINK_CONTRACT_ADDRESS"`
	MinIncomingConfirmations          uint32        `env:"MIN_INCOMING_CONFIRMATIONS"`
	MinRequiredOutgoingConfirmations  uint64        `env:"MIN_OUTGOING_CONFIRMATIONS"`
	MinimumContractPayment            assets.Link   `env:"MINIMUM_CONTRACT_PAYMENT_LINK_JUELS"`
	// EVM Gas Controls
	EvmEIP1559DynamicFees      bool     `env:"EVM_EIP1559_DYNAMIC_FEES"`
	EvmGasBumpPercent          uint16   `env:"ETH_GAS_BUMP_PERCENT"`
	EvmGasBumpThreshold        uint64   `env:"ETH_GAS_BUMP_THRESHOLD"`
	EvmGasBumpTxDepth          uint16   `env:"ETH_GAS_BUMP_TX_DEPTH"`
	EvmGasBumpWei              *big.Int `env:"ETH_GAS_BUMP_WEI"`
	EvmGasLimitDefault         uint64   `env:"ETH_GAS_LIMIT_DEFAULT"`
	EvmGasLimitMultiplier      float32  `env:"ETH_GAS_LIMIT_MULTIPLIER"`
	EvmGasLimitTransfer        uint64   `env:"ETH_GAS_LIMIT_TRANSFER"`
	EvmGasPriceDefault         *big.Int `env:"ETH_GAS_PRICE_DEFAULT"`
	EvmGasTipCapDefault        *big.Int `env:"EVM_GAS_TIP_CAP_DEFAULT"`
	EvmGasTipCapMinimum        *big.Int `env:"EVM_GAS_TIP_CAP_MINIMUM"`
	EvmMaxGasPriceWei          *big.Int `env:"ETH_MAX_GAS_PRICE_WEI"`
	EvmMaxInFlightTransactions uint32   `env:"ETH_MAX_IN_FLIGHT_TRANSACTIONS"`
	EvmMaxQueuedTransactions   uint64   `env:"ETH_MAX_QUEUED_TRANSACTIONS"`
	EvmMinGasPriceWei          *big.Int `env:"ETH_MIN_GAS_PRICE_WEI"`
	EvmNonceAutoSync           bool     `env:"ETH_NONCE_AUTO_SYNC"`
	// Gas Estimation
	GasEstimatorMode                           string `env:"GAS_ESTIMATOR_MODE"`
	BlockHistoryEstimatorBatchSize             uint32 `env:"BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE"`
	BlockHistoryEstimatorBlockDelay            uint16 `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY"`
	BlockHistoryEstimatorBlockHistorySize      uint16 `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE"`
	BlockHistoryEstimatorTransactionPercentile uint16 `env:"BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE"`

	// Job Pipeline and tasks
	DefaultHTTPAllowUnrestrictedNetworkAccess bool            `env:"DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS" default:"false"`
	DefaultHTTPLimit                          int64           `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout                        models.Duration `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	FeatureExternalInitiators                 bool            `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	JobPipelineMaxRunDuration                 time.Duration   `env:"JOB_PIPELINE_MAX_RUN_DURATION" default:"10m"`
	JobPipelineReaperInterval                 time.Duration   `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold                time.Duration   `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"24h"`
	JobPipelineResultWriteQueueDepth          uint64          `env:"JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH" default:"100"`

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
	OCR2MonitoringEndpoint                 string        `env:"OCR2_MONITORING_ENDPOINT"`                                 //nodoc

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
	OCRMonitoringEndpoint   string `env:"OCR_MONITORING_ENDPOINT"`
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
	OCROutgoingMessageBufferSize int           `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"` //nodoc
	OCRIncomingMessageBufferSize int           `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"` //nodoc
	OCRDHTLookupInterval         int           `env:"OCR_DHT_LOOKUP_INTERVAL" default:"10"`          //nodoc
	OCRBootstrapCheckInterval    time.Duration `env:"OCR_BOOTSTRAP_CHECK_INTERVAL" default:"20s"`    //nodoc
	OCRNewStreamTimeout          time.Duration `env:"OCR_NEW_STREAM_TIMEOUT" default:"10s"`          //nodoc

	// Keeper
	KeeperDefaultTransactionQueueDepth uint32        `env:"KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"` //nodoc
	KeeperGasPriceBufferPercent        uint32        `env:"KEEPER_GAS_PRICE_BUFFER_PERCENT" default:"20"`
	KeeperGasTipCapBufferPercent       uint32        `env:"KEEPER_GAS_TIP_CAP_BUFFER_PERCENT" default:"20"`
	KeeperMaximumGracePeriod           int64         `env:"KEEPER_MAXIMUM_GRACE_PERIOD" default:"100"`
	KeeperRegistryCheckGasOverhead     uint64        `env:"KEEPER_REGISTRY_CHECK_GAS_OVERHEAD" default:"200000"`
	KeeperRegistryPerformGasOverhead   uint64        `env:"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD" default:"150000"`
	KeeperRegistrySyncInterval         time.Duration `env:"KEEPER_REGISTRY_SYNC_INTERVAL" default:"30m"`
	KeeperRegistrySyncUpkeepQueueSize  uint32        `env:"KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE" default:"10"`

	// CLI client
	AdminCredentialsFile string `env:"ADMIN_CREDENTIALS_FILE" default:"$ROOT/apicredentials"`
	ClientNodeURL        string `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	InsecureSkipVerify   bool   `env:"INSECURE_SKIP_VERIFY" default:"false"`

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

func TryName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		return fmt.Sprintf("<%s>", field)
	}
	return item.Tag.Get("env")
}

func DefaultValue(name string) (string, bool) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		return item.Tag.Lookup("default")
	}
	log.Panicf("Invariant violated, no field of name %s found for DefaultValue", name)
	return "", false
}

func ZeroValue(name string) interface{} {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		if item.Type.Kind() == reflect.Ptr {
			return nil
		}
		return reflect.New(item.Type).Interface()
	}
	log.Panicf("Invariant violated, no field of name %s found for ZeroValue", name)
	return nil
}
