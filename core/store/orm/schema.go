package orm

import (
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// ConfigSchema records the schema of configuration at the type level
type ConfigSchema struct {
	AdminCredentialsFile                      string          `env:"ADMIN_CREDENTIALS_FILE" default:"$ROOT/apicredentials"`
	AllowOrigins                              string          `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	AuthenticatedRateLimit                    int64           `env:"AUTHENTICATED_RATE_LIMIT" default:"1000"`
	AuthenticatedRateLimitPeriod              time.Duration   `env:"AUTHENTICATED_RATE_LIMIT_PERIOD" default:"1m"`
	BalanceMonitorEnabled                     bool            `env:"BALANCE_MONITOR_ENABLED" default:"true"`
	BlockBackfillDepth                        string          `env:"BLOCK_BACKFILL_DEPTH" default:"10"`
	BridgeResponseURL                         url.URL         `env:"BRIDGE_RESPONSE_URL"`
	ChainID                                   big.Int         `env:"ETH_CHAIN_ID" default:"1"`
	ClientNodeURL                             string          `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	DatabaseBackupFrequency                   time.Duration   `env:"DATABASE_BACKUP_FREQUENCY" default:"0m"`
	DatabaseBackupMode                        string          `env:"DATABASE_BACKUP_MODE" default:"none"`
	DatabaseBackupURL                         *url.URL        `env:"DATABASE_BACKUP_URL" default:""`
	DatabaseBackupDir                         string          `env:"DATABASE_BACKUP_DIR" default:""`
	DatabaseListenerMaxReconnectDuration      time.Duration   `env:"DATABASE_LISTENER_MAX_RECONNECT_DURATION" default:"10m"`
	DatabaseListenerMinReconnectInterval      time.Duration   `env:"DATABASE_LISTENER_MIN_RECONNECT_INTERVAL" default:"1m"`
	DatabaseMaximumTxDuration                 time.Duration   `env:"DATABASE_MAXIMUM_TX_DURATION" default:"30m"`
	DatabaseTimeout                           models.Duration `env:"DATABASE_TIMEOUT" default:"0"`
	DatabaseURL                               string          `env:"DATABASE_URL"`
	DefaultHTTPAllowUnrestrictedNetworkAccess bool            `env:"DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS" default:"false"`
	DefaultHTTPLimit                          int64           `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout                        models.Duration `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	DefaultMaxHTTPAttempts                    uint            `env:"MAX_HTTP_ATTEMPTS" default:"5"`
	Dev                                       bool            `env:"CHAINLINK_DEV" default:"false"`
	EnableExperimentalAdapters                bool            `env:"ENABLE_EXPERIMENTAL_ADAPTERS" default:"false"`
	EnableLegacyJobPipeline                   bool            `env:"ENABLE_LEGACY_JOB_PIPELINE" default:"true"`
	EthBalanceMonitorBlockDelay               uint16          `env:"ETH_BALANCE_MONITOR_BLOCK_DELAY"`
	EthFinalityDepth                          uint            `env:"ETH_FINALITY_DEPTH"`
	EthGasBumpPercent                         uint16          `env:"ETH_GAS_BUMP_PERCENT" default:"20"`
	EthGasBumpThreshold                       uint64          `env:"ETH_GAS_BUMP_THRESHOLD"`
	EthGasBumpTxDepth                         uint16          `env:"ETH_GAS_BUMP_TX_DEPTH" default:"10"`
	EthGasBumpWei                             big.Int         `env:"ETH_GAS_BUMP_WEI"`
	EthGasLimitDefault                        uint64          `env:"ETH_GAS_LIMIT_DEFAULT"`
	EthGasLimitMultiplier                     float32         `env:"ETH_GAS_LIMIT_MULTIPLIER" default:"1.0"`
	EthGasPriceDefault                        big.Int         `env:"ETH_GAS_PRICE_DEFAULT"`
	EthHeadTrackerHistoryDepth                uint            `env:"ETH_HEAD_TRACKER_HISTORY_DEPTH"`
	EthHeadTrackerMaxBufferSize               uint            `env:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE" default:"3"`
	EthHeadTrackerSamplingInterval            time.Duration   `env:"ETH_HEAD_TRACKER_SAMPLING_INTERVAL" default:"1s"`
	EthLogBackfillBatchSize                   uint32          `env:"ETH_LOG_BACKFILL_BATCH_SIZE" default:"100"`
	EthMaxGasPriceWei                         big.Int         `env:"ETH_MAX_GAS_PRICE_WEI"`
	EthMaxInFlightTransactions                uint64          `env:"ETH_MAX_IN_FLIGHT_TRANSACTIONS" default:"16"`
	EthMaxQueuedTransactions                  uint64          `env:"ETH_MAX_QUEUED_TRANSACTIONS" default:"250"`
	EthMinGasPriceWei                         big.Int         `env:"ETH_MIN_GAS_PRICE_WEI"`
	EthNonceAutoSync                          bool            `env:"ETH_NONCE_AUTO_SYNC" default:"true"`
	EthRPCDefaultBatchSize                    uint32          `env:"ETH_RPC_DEFAULT_BATCH_SIZE" default:"100"`
	EthTxReaperInterval                       time.Duration   `env:"ETH_TX_REAPER_INTERVAL" default:"1h"`
	EthTxReaperThreshold                      time.Duration   `env:"ETH_TX_REAPER_THRESHOLD" default:"168h"`
	EthTxResendAfterThreshold                 time.Duration   `env:"ETH_TX_RESEND_AFTER_THRESHOLD"`
	EthereumDisabled                          bool            `env:"ETH_DISABLED" default:"false"`
	EthereumHTTPURL                           string          `env:"ETH_HTTP_URL"`
	EthereumSecondaryURL                      string          `env:"ETH_SECONDARY_URL" default:""`
	EthereumSecondaryURLs                     string          `env:"ETH_SECONDARY_URLS" default:""`
	EthereumURL                               string          `env:"ETH_URL" default:"ws://localhost:8546"`
	ExplorerAccessKey                         string          `env:"EXPLORER_ACCESS_KEY"`
	ExplorerSecret                            string          `env:"EXPLORER_SECRET"`
	ExplorerURL                               *url.URL        `env:"EXPLORER_URL"`
	FeatureCronV2                             bool            `env:"FEATURE_CRON_V2" default:"true"`
	FeatureExternalInitiators                 bool            `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	FeatureFluxMonitor                        bool            `env:"FEATURE_FLUX_MONITOR" default:"true"`
	FeatureFluxMonitorV2                      bool            `env:"FEATURE_FLUX_MONITOR_V2" default:"true"`
	FeatureOffchainReporting                  bool            `env:"FEATURE_OFFCHAIN_REPORTING" default:"false"`
	FeatureWebhookV2                          bool            `env:"FEATURE_WEBHOOK_V2" default:"false"`
	FlagsContractAddress                      string          `env:"FLAGS_CONTRACT_ADDRESS"`
	GasUpdaterBatchSize                       uint32          `env:"GAS_UPDATER_BATCH_SIZE"`
	GasUpdaterBlockDelay                      uint16          `env:"GAS_UPDATER_BLOCK_DELAY"`
	GasUpdaterBlockHistorySize                uint16          `env:"GAS_UPDATER_BLOCK_HISTORY_SIZE"`
	GasUpdaterEnabled                         bool            `env:"GAS_UPDATER_ENABLED"`
	GasUpdaterTransactionPercentile           uint16          `env:"GAS_UPDATER_TRANSACTION_PERCENTILE" default:"60"`
	GlobalLockRetryInterval                   models.Duration `env:"GLOBAL_LOCK_RETRY_INTERVAL" default:"1s"`
	HTTPServerWriteTimeout                    time.Duration   `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"10s"`
	InsecureFastScrypt                        bool            `env:"INSECURE_FAST_SCRYPT" default:"false"`
	JSONConsole                               bool            `env:"JSON_CONSOLE" default:"false"`
	JobPipelineMaxRunDuration                 time.Duration   `env:"JOB_PIPELINE_MAX_RUN_DURATION" default:"10m"`
	JobPipelineReaperInterval                 time.Duration   `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold                time.Duration   `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"24h"`
	JobPipelineResultWriteQueueDepth          uint64          `env:"JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH" default:"100"`
	KeeperMaximumGracePeriod                  int64           `env:"KEEPER_MAXIMUM_GRACE_PERIOD" default:"100"`
	KeeperMinimumRequiredConfirmations        uint64          `env:"KEEPER_MINIMUM_REQUIRED_CONFIRMATIONS" default:"12"`
	KeeperRegistrySyncInterval                time.Duration   `env:"KEEPER_REGISTRY_SYNC_INTERVAL" default:"30m"`
	LinkContractAddress                       string          `env:"LINK_CONTRACT_ADDRESS" default:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	LogLevel                                  LogLevel        `env:"LOG_LEVEL" default:"info"`
	LogSQLMigrations                          bool            `env:"LOG_SQL_MIGRATIONS" default:"true"`
	LogSQLStatements                          bool            `env:"LOG_SQL" default:"false"`
	LogToDisk                                 bool            `env:"LOG_TO_DISK" default:"true"`
	MaximumServiceDuration                    models.Duration `env:"MAXIMUM_SERVICE_DURATION" default:"8760h" `
	MigrateDatabase                           bool            `env:"MIGRATE_DATABASE" default:"true"`
	MinIncomingConfirmations                  uint32          `env:"MIN_INCOMING_CONFIRMATIONS"`
	MinRequiredOutgoingConfirmations          uint64          `env:"MIN_OUTGOING_CONFIRMATIONS"`
	MinimumContractPayment                    assets.Link     `env:"MINIMUM_CONTRACT_PAYMENT" default:"1000000000000000000"`
	MinimumRequestExpiration                  uint64          `env:"MINIMUM_REQUEST_EXPIRATION" default:"300"`
	MinimumServiceDuration                    models.Duration `env:"MINIMUM_SERVICE_DURATION" default:"0s" `
	OCRBlockchainTimeout                      time.Duration   `env:"OCR_BLOCKCHAIN_TIMEOUT" default:"20s"`
	OCRBootstrapCheckInterval                 time.Duration   `env:"OCR_BOOTSTRAP_CHECK_INTERVAL" default:"20s"`
	OCRContractConfirmations                  uint            `env:"OCR_CONTRACT_CONFIRMATIONS" default:"3"`
	OCRContractPollInterval                   time.Duration   `env:"OCR_CONTRACT_POLL_INTERVAL" default:"1m"`
	OCRContractSubscribeInterval              time.Duration   `env:"OCR_CONTRACT_SUBSCRIBE_INTERVAL" default:"2m"`
	OCRContractTransmitterTransmitTimeout     time.Duration   `env:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT" default:"10s"`
	OCRDHTLookupInterval                      int             `env:"OCR_DHT_LOOKUP_INTERVAL" default:"10"`
	OCRDatabaseTimeout                        time.Duration   `env:"OCR_DATABASE_TIMEOUT" default:"10s"`
	OCRIncomingMessageBufferSize              int             `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRKeyBundleID                            string          `env:"OCR_KEY_BUNDLE_ID"`
	OCRMonitoringEndpoint                     string          `env:"OCR_MONITORING_ENDPOINT"`
	OCRNewStreamTimeout                       time.Duration   `env:"OCR_NEW_STREAM_TIMEOUT" default:"10s"`
	OCRObservationGracePeriod                 time.Duration   `env:"OCR_OBSERVATION_GRACE_PERIOD" default:"1s"`
	OCRObservationTimeout                     time.Duration   `env:"OCR_OBSERVATION_TIMEOUT" default:"12s"`
	OCROutgoingMessageBufferSize              int             `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRTraceLogging                           bool            `env:"OCR_TRACE_LOGGING" default:"false"`
	OCRTransmitterAddress                     string          `env:"OCR_TRANSMITTER_ADDRESS"`
	ORMMaxIdleConns                           int             `env:"ORM_MAX_IDLE_CONNS" default:"10"`
	ORMMaxOpenConns                           int             `env:"ORM_MAX_OPEN_CONNS" default:"20"`
	OperatorContractAddress                   common.Address  `env:"OPERATOR_CONTRACT_ADDRESS"`
	OptimismGasFees                           bool            `env:"OPTIMISM_GAS_FEES"`
	P2PAnnounceIP                             net.IP          `env:"P2P_ANNOUNCE_IP"`
	P2PAnnouncePort                           uint16          `env:"P2P_ANNOUNCE_PORT"`
	P2PBootstrapPeers                         []string        `env:"P2P_BOOTSTRAP_PEERS"`
	P2PDHTAnnouncementCounterUserPrefix       uint32          `env:"P2P_DHT_ANNOUNCEMENT_COUNTER_USER_PREFIX" default:"0"`
	P2PListenIP                               net.IP          `env:"P2P_LISTEN_IP" default:"0.0.0.0"`
	P2PListenPort                             uint16          `env:"P2P_LISTEN_PORT"`
	P2PPeerID                                 models.PeerID   `env:"P2P_PEER_ID"`
	P2PPeerstoreWriteInterval                 time.Duration   `env:"P2P_PEERSTORE_WRITE_INTERVAL" default:"5m"`
	Port                                      uint16          `env:"CHAINLINK_PORT" default:"6688"`
	ReaperExpiration                          models.Duration `env:"REAPER_EXPIRATION" default:"240h"`
	ReplayFromBlock                           int64           `env:"REPLAY_FROM_BLOCK" default:"-1"`
	RootDir                                   string          `env:"ROOT" default:"~/.chainlink"`
	SecureCookies                             bool            `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout                            models.Duration `env:"SESSION_TIMEOUT" default:"15m"`
	StatsPusherLogging                        string          `env:"STATS_PUSHER_LOGGING" default:"false"`
	TLSCertPath                               string          `env:"TLS_CERT_PATH" `
	TLSHost                                   string          `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath                                string          `env:"TLS_KEY_PATH" `
	TLSPort                                   uint16          `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect                               bool            `env:"CHAINLINK_TLS_REDIRECT" default:"false"`
	TriggerFallbackDBPollInterval             time.Duration   `env:"TRIGGER_FALLBACK_DB_POLL_INTERVAL" default:"30s"`
	UnAuthenticatedRateLimit                  int64           `env:"UNAUTHENTICATED_RATE_LIMIT" default:"5"`
	UnAuthenticatedRateLimitPeriod            time.Duration   `env:"UNAUTHENTICATED_RATE_LIMIT_PERIOD" default:"20s"`
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
