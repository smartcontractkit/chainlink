package config

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
)

// ConfigSchema records the schema of configuration at the type level
type ConfigSchema struct {
	FeatureFeedsManager bool `env:"FEATURE_FEEDS_MANAGER" default:"false"`

	AdminCredentialsFile                       string                        `env:"ADMIN_CREDENTIALS_FILE" default:"$ROOT/apicredentials"`
	AllowOrigins                               string                        `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	AuthenticatedRateLimit                     int64                         `env:"AUTHENTICATED_RATE_LIMIT" default:"1000"`
	AuthenticatedRateLimitPeriod               time.Duration                 `env:"AUTHENTICATED_RATE_LIMIT_PERIOD" default:"1m"`
	AutoPprofEnabled                           bool                          `env:"AUTO_PPROF_ENABLED" default:"false"`
	AutoPprofProfileRoot                       string                        `env:"AUTO_PPROF_PROFILE_ROOT"` // Defaults to $CHAINLINK_ROOT
	AutoPprofPollInterval                      models.Duration               `env:"AUTO_PPROF_POLL_INTERVAL" default:"10s"`
	AutoPprofGatherDuration                    models.Duration               `env:"AUTO_PPROF_GATHER_DURATION" default:"10s"`
	AutoPprofGatherTraceDuration               models.Duration               `env:"AUTO_PPROF_GATHER_TRACE_DURATION" default:"5s"`
	AutoPprofMaxProfileSize                    utils.FileSize                `env:"AUTO_PPROF_MAX_PROFILE_SIZE" default:"100mb"`
	AutoPprofCPUProfileRate                    int                           `env:"AUTO_PPROF_CPU_PROFILE_RATE" default:"1"`
	AutoPprofMemProfileRate                    int                           `env:"AUTO_PPROF_MEM_PROFILE_RATE" default:"1"`
	AutoPprofBlockProfileRate                  int                           `env:"AUTO_PPROF_BLOCK_PROFILE_RATE" default:"1"`
	AutoPprofMutexProfileFraction              int                           `env:"AUTO_PPROF_MUTEX_PROFILE_FRACTION" default:"1"`
	AutoPprofMemThreshold                      utils.FileSize                `env:"AUTO_PPROF_MEM_THRESHOLD" default:"4gb"`
	AutoPprofGoroutineThreshold                int                           `env:"AUTO_PPROF_GOROUTINE_THRESHOLD" default:"5000"`
	BalanceMonitorEnabled                      bool                          `env:"BALANCE_MONITOR_ENABLED"`
	BlockBackfillDepth                         uint64                        `env:"BLOCK_BACKFILL_DEPTH" default:"10"`
	BlockBackfillSkip                          bool                          `env:"BLOCK_BACKFILL_SKIP" default:"false"`
	BlockEmissionIdleWarningThreshold          time.Duration                 `env:"BLOCK_EMISSION_IDLE_WARNING_THRESHOLD"`
	BlockHistoryEstimatorBatchSize             uint32                        `env:"BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE"`
	BlockHistoryEstimatorBlockDelay            uint16                        `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY"`
	BlockHistoryEstimatorBlockHistorySize      uint16                        `env:"BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE"`
	BlockHistoryEstimatorTransactionPercentile uint16                        `env:"BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE"`
	BridgeResponseURL                          url.URL                       `env:"BRIDGE_RESPONSE_URL"`
	ChainType                                  string                        `env:"CHAIN_TYPE"`
	ClientNodeURL                              string                        `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	DatabaseBackupDir                          string                        `env:"DATABASE_BACKUP_DIR" default:""`
	DatabaseBackupFrequency                    time.Duration                 `env:"DATABASE_BACKUP_FREQUENCY" default:"1h"`
	DatabaseBackupMode                         string                        `env:"DATABASE_BACKUP_MODE" default:"none"`
	DatabaseBackupURL                          *url.URL                      `env:"DATABASE_BACKUP_URL" default:""`
	DatabaseListenerMaxReconnectDuration       time.Duration                 `env:"DATABASE_LISTENER_MAX_RECONNECT_DURATION" default:"10m"`
	DatabaseListenerMinReconnectInterval       time.Duration                 `env:"DATABASE_LISTENER_MIN_RECONNECT_INTERVAL" default:"1m"`
	DatabaseLockingMode                        string                        `env:"DATABASE_LOCKING_MODE" default:"advisorylock"`
	DatabaseMaximumTxDuration                  time.Duration                 `env:"DATABASE_MAXIMUM_TX_DURATION" default:"30m"`
	DatabaseTimeout                            models.Duration               `env:"DATABASE_TIMEOUT" default:"0"`
	DatabaseURL                                string                        `env:"DATABASE_URL"`
	DefaultChainID                             *big.Int                      `env:"ETH_CHAIN_ID"`
	DefaultHTTPAllowUnrestrictedNetworkAccess  bool                          `env:"DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS" default:"false"`
	DefaultHTTPLimit                           int64                         `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout                         models.Duration               `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	DefaultMaxHTTPAttempts                     uint                          `env:"MAX_HTTP_ATTEMPTS" default:"5"`
	Dev                                        bool                          `env:"CHAINLINK_DEV" default:"false"`
	EVMDisabled                                bool                          `env:"EVM_DISABLED" default:"false"`
	EthTxReaperInterval                        time.Duration                 `env:"ETH_TX_REAPER_INTERVAL"`
	EthTxReaperThreshold                       time.Duration                 `env:"ETH_TX_REAPER_THRESHOLD"`
	EthTxResendAfterThreshold                  time.Duration                 `env:"ETH_TX_RESEND_AFTER_THRESHOLD"`
	EthereumDisabled                           bool                          `env:"ETH_DISABLED" default:"false"`
	EthereumHTTPURL                            string                        `env:"ETH_HTTP_URL"`
	EthereumSecondaryURL                       string                        `env:"ETH_SECONDARY_URL" default:""`
	EthereumSecondaryURLs                      string                        `env:"ETH_SECONDARY_URLS" default:""`
	EthereumURL                                string                        `env:"ETH_URL"`
	EvmDefaultBatchSize                        uint32                        `env:"ETH_DEFAULT_BATCH_SIZE"`
	EvmEIP1559DynamicFees                      bool                          `env:"EVM_EIP1559_DYNAMIC_FEES"`
	EvmFinalityDepth                           uint32                        `env:"ETH_FINALITY_DEPTH"`
	EvmGasBumpPercent                          uint16                        `env:"ETH_GAS_BUMP_PERCENT"`
	EvmGasBumpThreshold                        uint64                        `env:"ETH_GAS_BUMP_THRESHOLD"`
	EvmGasBumpTxDepth                          uint16                        `env:"ETH_GAS_BUMP_TX_DEPTH"`
	EvmGasBumpWei                              *big.Int                      `env:"ETH_GAS_BUMP_WEI"`
	EvmGasLimitDefault                         uint64                        `env:"ETH_GAS_LIMIT_DEFAULT"`
	EvmGasLimitMultiplier                      float32                       `env:"ETH_GAS_LIMIT_MULTIPLIER"`
	EvmGasLimitTransfer                        uint64                        `env:"ETH_GAS_LIMIT_TRANSFER"`
	EvmGasPriceDefault                         *big.Int                      `env:"ETH_GAS_PRICE_DEFAULT"`
	EvmGasTipCapDefault                        *big.Int                      `env:"EVM_GAS_TIP_CAP_DEFAULT"`
	EvmGasTipCapMinimum                        *big.Int                      `env:"EVM_GAS_TIP_CAP_MINIMUM"`
	EvmHeadTrackerHistoryDepth                 uint                          `env:"ETH_HEAD_TRACKER_HISTORY_DEPTH"`
	EvmHeadTrackerMaxBufferSize                uint                          `env:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE"`
	EvmHeadTrackerSamplingInterval             time.Duration                 `env:"ETH_HEAD_TRACKER_SAMPLING_INTERVAL"`
	EvmLogBackfillBatchSize                    uint32                        `env:"ETH_LOG_BACKFILL_BATCH_SIZE"`
	EvmMaxGasPriceWei                          *big.Int                      `env:"ETH_MAX_GAS_PRICE_WEI"`
	EvmMaxInFlightTransactions                 uint32                        `env:"ETH_MAX_IN_FLIGHT_TRANSACTIONS"`
	EvmMaxQueuedTransactions                   uint64                        `env:"ETH_MAX_QUEUED_TRANSACTIONS"`
	EvmMinGasPriceWei                          *big.Int                      `env:"ETH_MIN_GAS_PRICE_WEI"`
	EvmNonceAutoSync                           bool                          `env:"ETH_NONCE_AUTO_SYNC"`
	EvmRPCDefaultBatchSize                     uint32                        `env:"ETH_RPC_DEFAULT_BATCH_SIZE"`
	ExplorerAccessKey                          string                        `env:"EXPLORER_ACCESS_KEY"`
	ExplorerSecret                             string                        `env:"EXPLORER_SECRET"`
	ExplorerURL                                *url.URL                      `env:"EXPLORER_URL"`
	FMDefaultTransactionQueueDepth             uint32                        `env:"FM_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"`
	FMSimulateTransactions                     bool                          `env:"FM_SIMULATE_TRANSACTIONS" default:"false"`
	FeatureExternalInitiators                  bool                          `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	FeatureOffchainReporting                   bool                          `env:"FEATURE_OFFCHAIN_REPORTING" default:"false"`
	FeatureUICSAKeys                           bool                          `env:"FEATURE_UI_CSA_KEYS" default:"false"`
	FeatureUIFeedsManager                      bool                          `env:"FEATURE_UI_FEEDS_MANAGER" default:"false"`
	FlagsContractAddress                       string                        `env:"FLAGS_CONTRACT_ADDRESS"`
	GasEstimatorMode                           string                        `env:"GAS_ESTIMATOR_MODE"`
	GlobalLockRetryInterval                    models.Duration               `env:"GLOBAL_LOCK_RETRY_INTERVAL" default:"1s"`
	HTTPServerWriteTimeout                     time.Duration                 `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"10s"`
	InsecureFastScrypt                         bool                          `env:"INSECURE_FAST_SCRYPT" default:"false"`
	InsecureSkipVerify                         bool                          `env:"INSECURE_SKIP_VERIFY" default:"false"`
	JSONConsole                                bool                          `env:"JSON_CONSOLE" default:"false"`
	JobPipelineMaxRunDuration                  time.Duration                 `env:"JOB_PIPELINE_MAX_RUN_DURATION" default:"10m"`
	JobPipelineReaperInterval                  time.Duration                 `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold                 time.Duration                 `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"24h"`
	JobPipelineResultWriteQueueDepth           uint64                        `env:"JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH" default:"100"`
	KeeperDefaultTransactionQueueDepth         uint32                        `env:"KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"`
	KeeperGasPriceBufferPercent                uint32                        `env:"KEEPER_GAS_PRICE_BUFFER_PERCENT" default:"20"`
	KeeperGasTipCapBufferPercent               uint32                        `env:"KEEPER_GAS_TIP_CAP_BUFFER_PERCENT" default:"20"`
	KeeperMaximumGracePeriod                   int64                         `env:"KEEPER_MAXIMUM_GRACE_PERIOD" default:"100"`
	KeeperRegistryCheckGasOverhead             uint64                        `env:"KEEPER_REGISTRY_CHECK_GAS_OVERHEAD" default:"200000"`
	KeeperRegistryPerformGasOverhead           uint64                        `env:"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD" default:"150000"`
	KeeperRegistrySyncInterval                 time.Duration                 `env:"KEEPER_REGISTRY_SYNC_INTERVAL" default:"30m"`
	KeeperRegistrySyncUpkeepQueueSize          uint32                        `env:"KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE" default:"10"`
	LeaseLockRefreshInterval                   time.Duration                 `env:"LEASE_LOCK_REFRESH_INTERVAL" default:"1s"`
	LeaseLockDuration                          time.Duration                 `env:"LEASE_LOCK_DURATION" default:"10s"`
	LinkContractAddress                        string                        `env:"LINK_CONTRACT_ADDRESS"`
	LogLevel                                   LogLevel                      `env:"LOG_LEVEL"`
	LogSQLMigrations                           bool                          `env:"LOG_SQL_MIGRATIONS" default:"true"`
	LogSQL                                     bool                          `env:"LOG_SQL" default:"false"`
	LogToDisk                                  bool                          `env:"LOG_TO_DISK" default:"false"`
	LogUnixTS                                  bool                          `env:"LOG_UNIX_TS" default:"false"`
	MigrateDatabase                            bool                          `env:"MIGRATE_DATABASE" default:"true"`
	MinIncomingConfirmations                   uint32                        `env:"MIN_INCOMING_CONFIRMATIONS"`
	MinRequiredOutgoingConfirmations           uint64                        `env:"MIN_OUTGOING_CONFIRMATIONS"`
	MinimumContractPayment                     assets.Link                   `env:"MINIMUM_CONTRACT_PAYMENT_LINK_JUELS"`
	OCRBlockchainTimeout                       time.Duration                 `env:"OCR_BLOCKCHAIN_TIMEOUT" default:"20s"`
	OCRBootstrapCheckInterval                  time.Duration                 `env:"OCR_BOOTSTRAP_CHECK_INTERVAL" default:"20s"`
	OCRContractConfirmations                   uint                          `env:"OCR_CONTRACT_CONFIRMATIONS"`
	OCRContractPollInterval                    time.Duration                 `env:"OCR_CONTRACT_POLL_INTERVAL" default:"1m"`
	OCRContractSubscribeInterval               time.Duration                 `env:"OCR_CONTRACT_SUBSCRIBE_INTERVAL" default:"2m"`
	OCRContractTransmitterTransmitTimeout      time.Duration                 `env:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT" default:"10s"`
	OCRDHTLookupInterval                       int                           `env:"OCR_DHT_LOOKUP_INTERVAL" default:"10"`
	OCRDatabaseTimeout                         time.Duration                 `env:"OCR_DATABASE_TIMEOUT" default:"10s"`
	OCRDefaultTransactionQueueDepth            uint32                        `env:"OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH" default:"1"`
	OCRIncomingMessageBufferSize               int                           `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRKeyBundleID                             string                        `env:"OCR_KEY_BUNDLE_ID"`
	OCRMonitoringEndpoint                      string                        `env:"OCR_MONITORING_ENDPOINT"`
	OCRNewStreamTimeout                        time.Duration                 `env:"OCR_NEW_STREAM_TIMEOUT" default:"10s"`
	OCRObservationGracePeriod                  time.Duration                 `env:"OCR_OBSERVATION_GRACE_PERIOD" default:"1s"`
	OCRObservationTimeout                      time.Duration                 `env:"OCR_OBSERVATION_TIMEOUT" default:"5s"`
	OCROutgoingMessageBufferSize               int                           `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRSimulateTransactions                    bool                          `env:"OCR_SIMULATE_TRANSACTIONS" default:"false"`
	OCRTraceLogging                            bool                          `env:"OCR_TRACE_LOGGING" default:"false"`
	OCRTransmitterAddress                      string                        `env:"OCR_TRANSMITTER_ADDRESS"`
	ORMMaxIdleConns                            int                           `env:"ORM_MAX_IDLE_CONNS" default:"10"`
	ORMMaxOpenConns                            int                           `env:"ORM_MAX_OPEN_CONNS" default:"20"`
	P2PAnnounceIP                              net.IP                        `env:"P2P_ANNOUNCE_IP"`
	P2PAnnouncePort                            uint16                        `env:"P2P_ANNOUNCE_PORT"`
	P2PBootstrapPeers                          []string                      `env:"P2P_BOOTSTRAP_PEERS"`
	P2PDHTAnnouncementCounterUserPrefix        uint32                        `env:"P2P_DHT_ANNOUNCEMENT_COUNTER_USER_PREFIX" default:"0"`
	P2PListenIP                                net.IP                        `env:"P2P_LISTEN_IP" default:"0.0.0.0"`
	P2PListenPort                              uint16                        `env:"P2P_LISTEN_PORT"`
	P2PNetworkingStack                         ocrnetworking.NetworkingStack `env:"P2P_NETWORKING_STACK" default:"V1"`
	P2PPeerID                                  p2pkey.PeerID                 `env:"P2P_PEER_ID"`
	P2PPeerstoreWriteInterval                  time.Duration                 `env:"P2P_PEERSTORE_WRITE_INTERVAL" default:"5m"`
	P2PV2AnnounceAddresses                     []string                      `env:"P2PV2_ANNOUNCE_ADDRESSES"`
	P2PV2Bootstrappers                         []string                      `env:"P2PV2_BOOTSTRAPPERS"`
	P2PV2DeltaDial                             models.Duration               `env:"P2PV2_DELTA_DIAL" default:"15s"`
	P2PV2DeltaReconcile                        models.Duration               `env:"P2PV2_DELTA_RECONCILE" default:"1m"`
	P2PV2ListenAddresses                       []string                      `env:"P2PV2_LISTEN_ADDRESSES"`
	Port                                       uint16                        `env:"CHAINLINK_PORT" default:"6688"`
	RPID                                       string                        `env:"MFA_RPID"`
	RPOrigin                                   string                        `env:"MFA_RPORIGIN"`
	ReaperExpiration                           models.Duration               `env:"REAPER_EXPIRATION" default:"240h"`
	ReplayFromBlock                            int64                         `env:"REPLAY_FROM_BLOCK" default:"-1"`
	RootDir                                    string                        `env:"ROOT" default:"~/.chainlink"`
	SecureCookies                              bool                          `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout                             models.Duration               `env:"SESSION_TIMEOUT" default:"15m"`
	StatsPusherLogging                         string                        `env:"STATS_PUSHER_LOGGING" default:"false"`
	TLSCertPath                                string                        `env:"TLS_CERT_PATH" `
	TLSHost                                    string                        `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath                                 string                        `env:"TLS_KEY_PATH" `
	TLSPort                                    uint16                        `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect                                bool                          `env:"CHAINLINK_TLS_REDIRECT" default:"false"`
	TelemetryIngressLogging                    bool                          `env:"TELEMETRY_INGRESS_LOGGING" default:"false"`
	TelemetryIngressServerPubKey               string                        `env:"TELEMETRY_INGRESS_SERVER_PUB_KEY"`
	TelemetryIngressURL                        *url.URL                      `env:"TELEMETRY_INGRESS_URL"`
	TriggerFallbackDBPollInterval              time.Duration                 `env:"TRIGGER_FALLBACK_DB_POLL_INTERVAL" default:"30s"`
	UnAuthenticatedRateLimit                   int64                         `env:"UNAUTHENTICATED_RATE_LIMIT" default:"5"`
	UnAuthenticatedRateLimitPeriod             time.Duration                 `env:"UNAUTHENTICATED_RATE_LIMIT_PERIOD" default:"20s"`
}

// EnvVarName gets the environment variable name for a config schema field
func EnvVarName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		log.Panicf("Invariant violated, no field of name %s found on ConfigSchema", field)
	}
	return item.Tag.Get("env")
}

func TryEnvVarName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		return fmt.Sprintf("<%s>", field)
	}
	return item.Tag.Get("env")
}

func defaultValue(name string) (string, bool) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		return item.Tag.Lookup("default")
	}
	log.Panicf("Invariant violated, no field of name %s found for defaultValue", name)
	return "", false
}

func zeroValue(name string) interface{} {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		if item.Type.Kind() == reflect.Ptr {
			return nil
		}
		return reflect.New(item.Type).Interface()
	}
	log.Panicf("Invariant violated, no field of name %s found for zeroValue", name)
	return nil
}
