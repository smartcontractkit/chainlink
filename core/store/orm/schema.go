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
	DatabaseTimeout                           models.Duration `env:"DATABASE_TIMEOUT" default:"0"`
	DatabaseURL                               string          `env:"DATABASE_URL"`
	DatabaseListenerMinReconnectInterval      time.Duration   `env:"DATABASE_LISTENER_MIN_RECONNECT_INTERVAL" default:"1m"`
	DatabaseListenerMaxReconnectDuration      time.Duration   `env:"DATABASE_LISTENER_MAX_RECONNECT_DURATION" default:"10m"`
	DatabaseMaximumTxDuration                 time.Duration   `env:"DATABASE_MAXIMUM_TX_DURATION" default:"30m"`
	DatabaseBackupMode                        string          `env:"DATABASE_BACKUP_MODE" default:"none"`
	DatabaseBackupFrequency                   time.Duration   `env:"DATABASE_BACKUP_FREQUENCY" default:"0m"`
	DatabaseBackupURL                         *url.URL        `env:"DATABASE_BACKUP_URL" default:""`
	DefaultHTTPLimit                          int64           `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout                        models.Duration `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	DefaultHTTPAllowUnrestrictedNetworkAccess bool            `env:"DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS" default:"false"`
	Dev                                       bool            `env:"CHAINLINK_DEV" default:"false"`
	EnableExperimentalAdapters                bool            `env:"ENABLE_EXPERIMENTAL_ADAPTERS" default:"false"`
	FeatureExternalInitiators                 bool            `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	FeatureFluxMonitor                        bool            `env:"FEATURE_FLUX_MONITOR" default:"true"`
	FeatureFluxMonitorV2                      bool            `env:"FEATURE_FLUX_MONITOR_V2" default:"false"`
	FeatureOffchainReporting                  bool            `env:"FEATURE_OFFCHAIN_REPORTING" default:"false"`
	GlobalLockRetryInterval                   models.Duration `env:"GLOBAL_LOCK_RETRY_INTERVAL" default:"1s"`
	MaximumServiceDuration                    models.Duration `env:"MAXIMUM_SERVICE_DURATION" default:"8760h" `
	MinimumServiceDuration                    models.Duration `env:"MINIMUM_SERVICE_DURATION" default:"0s" `
	EthGasBumpThreshold                       uint64          `env:"ETH_GAS_BUMP_THRESHOLD"`
	EthGasBumpWei                             big.Int         `env:"ETH_GAS_BUMP_WEI"`
	EthGasBumpPercent                         uint16          `env:"ETH_GAS_BUMP_PERCENT" default:"20"`
	EthGasBumpTxDepth                         uint16          `env:"ETH_GAS_BUMP_TX_DEPTH" default:"10"`
	EthGasLimitDefault                        uint64          `env:"ETH_GAS_LIMIT_DEFAULT" default:"500000"`
	EthGasPriceDefault                        big.Int         `env:"ETH_GAS_PRICE_DEFAULT"`
	EthMaxGasPriceWei                         big.Int         `env:"ETH_MAX_GAS_PRICE_WEI"`
	EthMaxUnconfirmedTransactions             uint64          `env:"ETH_MAX_UNCONFIRMED_TRANSACTIONS" default:"500"`
	EthFinalityDepth                          uint            `env:"ETH_FINALITY_DEPTH"`
	EthHeadTrackerHistoryDepth                uint            `env:"ETH_HEAD_TRACKER_HISTORY_DEPTH"`
	EthHeadTrackerMaxBufferSize               uint            `env:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE" default:"3"`
	EthBalanceMonitorBlockDelay               uint16          `env:"ETH_BALANCE_MONITOR_BLOCK_DELAY"`
	EthReceiptFetchBatchSize                  uint32          `env:"ETH_RECEIPT_FETCH_BATCH_SIZE" default:"100"`
	EthTxResendAfterThreshold                 time.Duration   `env:"ETH_TX_RESEND_AFTER_THRESHOLD"`
	EthLogBackfillBatchSize                   uint32          `env:"ETH_LOG_BACKFILL_BATCH_SIZE" default:"100"`
	EthereumURL                               string          `env:"ETH_URL" default:"ws://localhost:8546"`
	EthereumSecondaryURL                      string          `env:"ETH_SECONDARY_URL" default:""`
	EthereumSecondaryURLs                     string          `env:"ETH_SECONDARY_URLS" default:""`
	EthereumDisabled                          bool            `env:"ETH_DISABLED" default:"false"`
	FlagsContractAddress                      string          `env:"FLAGS_CONTRACT_ADDRESS"`
	GasUpdaterBlockDelay                      uint16          `env:"GAS_UPDATER_BLOCK_DELAY"`
	GasUpdaterBlockHistorySize                uint16          `env:"GAS_UPDATER_BLOCK_HISTORY_SIZE"`
	GasUpdaterTransactionPercentile           uint16          `env:"GAS_UPDATER_TRANSACTION_PERCENTILE" default:"60"`
	GasUpdaterEnabled                         bool            `env:"GAS_UPDATER_ENABLED" default:"true"`
	HeadTimeBudget                            time.Duration   `env:"HEAD_TIME_BUDGET"`
	InsecureFastScrypt                        bool            `env:"INSECURE_FAST_SCRYPT" default:"false"`
	JobPipelineMaxRunDuration                 time.Duration   `env:"JOB_PIPELINE_MAX_RUN_DURATION" default:"10m"`
	JobPipelineResultWriteQueueDepth          uint64          `env:"JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH" default:"100"`
	JobPipelineParallelism                    uint8           `env:"JOB_PIPELINE_PARALLELISM" default:"4"`
	JobPipelineReaperInterval                 time.Duration   `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold                time.Duration   `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"24h"`
	JSONConsole                               bool            `env:"JSON_CONSOLE" default:"false"`
	KeeperRegistrySyncInterval                time.Duration   `env:"KEEPER_REGISTRY_SYNC_INTERVAL" default:"30m"`
	KeeperMinimumRequiredConfirmations        uint64          `env:"KEEPER_MINIMUM_REQUIRED_CONFIRMATIONS" default:"12"`
	KeeperMaximumGracePeriod                  int64           `env:"KEEPER_MAXIMUM_GRACE_PERIOD" default:"100"`
	LinkContractAddress                       string          `env:"LINK_CONTRACT_ADDRESS" default:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	ExplorerURL                               *url.URL        `env:"EXPLORER_URL"`
	ExplorerAccessKey                         string          `env:"EXPLORER_ACCESS_KEY"`
	ExplorerSecret                            string          `env:"EXPLORER_SECRET"`
	LogLevel                                  LogLevel        `env:"LOG_LEVEL" default:"info"`
	LogToDisk                                 bool            `env:"LOG_TO_DISK" default:"true"`
	LogSQLStatements                          bool            `env:"LOG_SQL" default:"false"`
	LogSQLMigrations                          bool            `env:"LOG_SQL_MIGRATIONS" default:"true"`
	DefaultMaxHTTPAttempts                    uint            `env:"MAX_HTTP_ATTEMPTS" default:"5"`
	MigrateDatabase                           bool            `env:"MIGRATE_DATABASE" default:"true"`
	MinIncomingConfirmations                  uint32          `env:"MIN_INCOMING_CONFIRMATIONS"`
	MinRequiredOutgoingConfirmations          uint64          `env:"MIN_OUTGOING_CONFIRMATIONS"`
	MinimumContractPayment                    assets.Link     `env:"MINIMUM_CONTRACT_PAYMENT" default:"1000000000000000000"`
	MinimumRequestExpiration                  uint64          `env:"MINIMUM_REQUEST_EXPIRATION" default:"300"`
	OCRObservationTimeout                     time.Duration   `env:"OCR_OBSERVATION_TIMEOUT" default:"12s"`
	OCRObservationGracePeriod                 time.Duration   `env:"OCR_OBSERVATION_GRACE_PERIOD" default:"1s"`
	OCRBlockchainTimeout                      time.Duration   `env:"OCR_BLOCKCHAIN_TIMEOUT" default:"20s"`
	OCRContractSubscribeInterval              time.Duration   `env:"OCR_CONTRACT_SUBSCRIBE_INTERVAL" default:"2m"`
	OCRContractPollInterval                   time.Duration   `env:"OCR_CONTRACT_POLL_INTERVAL" default:"1m"`
	OCRContractConfirmations                  uint            `env:"OCR_CONTRACT_CONFIRMATIONS" default:"3"`
	OCRBootstrapCheckInterval                 time.Duration   `env:"OCR_BOOTSTRAP_CHECK_INTERVAL" default:"20s"`
	OCRContractTransmitterTransmitTimeout     time.Duration   `env:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT" default:"10s"`
	OCRTransmitterAddress                     string          `env:"OCR_TRANSMITTER_ADDRESS"`
	OCRKeyBundleID                            string          `env:"OCR_KEY_BUNDLE_ID"`
	OCRDatabaseTimeout                        time.Duration   `env:"OCR_DATABASE_TIMEOUT" default:"10s"`
	OCRIncomingMessageBufferSize              int             `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCROutgoingMessageBufferSize              int             `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRNewStreamTimeout                       time.Duration   `env:"OCR_NEW_STREAM_TIMEOUT" default:"10s"`
	OCRDHTLookupInterval                      int             `env:"OCR_DHT_LOOKUP_INTERVAL" default:"10"`
	OCRTraceLogging                           bool            `env:"OCR_TRACE_LOGGING" default:"false"`
	OCRMonitoringEndpoint                     string          `env:"OCR_MONITORING_ENDPOINT"`
	OperatorContractAddress                   common.Address  `env:"OPERATOR_CONTRACT_ADDRESS"`
	ORMMaxOpenConns                           int             `env:"ORM_MAX_OPEN_CONNS" default:"20"`
	ORMMaxIdleConns                           int             `env:"ORM_MAX_IDLE_CONNS" default:"10"`
	P2PAnnounceIP                             net.IP          `env:"P2P_ANNOUNCE_IP"`
	P2PAnnouncePort                           uint16          `env:"P2P_ANNOUNCE_PORT"`
	P2PDHTAnnouncementCounterUserPrefix       uint32          `env:"P2P_DHT_ANNOUNCEMENT_COUNTER_USER_PREFIX" default:"0"`
	P2PListenIP                               net.IP          `env:"P2P_LISTEN_IP" default:"0.0.0.0"`
	P2PListenPort                             uint16          `env:"P2P_LISTEN_PORT"`
	P2PPeerstoreWriteInterval                 time.Duration   `env:"P2P_PEERSTORE_WRITE_INTERVAL" default:"5m"`
	P2PPeerID                                 models.PeerID   `env:"P2P_PEER_ID"`
	P2PBootstrapPeers                         []string        `env:"P2P_BOOTSTRAP_PEERS"`
	Port                                      uint16          `env:"CHAINLINK_PORT" default:"6688"`
	ReaperExpiration                          models.Duration `env:"REAPER_EXPIRATION" default:"240h"`
	ReplayFromBlock                           int64           `env:"REPLAY_FROM_BLOCK" default:"-1"`
	RootDir                                   string          `env:"ROOT" default:"~/.chainlink"`
	SecureCookies                             bool            `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout                            models.Duration `env:"SESSION_TIMEOUT" default:"15m"`
	StatsPusherLogging                        string          `env:"STATS_PUSHER_LOGGING" default:"false"`
	TriggerFallbackDBPollInterval             time.Duration   `env:"TRIGGER_FALLBACK_DB_POLL_INTERVAL" default:"30s"`
	TLSCertPath                               string          `env:"TLS_CERT_PATH" `
	TLSHost                                   string          `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath                                string          `env:"TLS_KEY_PATH" `
	TLSPort                                   uint16          `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect                               bool            `env:"CHAINLINK_TLS_REDIRECT" default:"false"`
	HTTPServerWriteTimeout                    time.Duration   `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"10s"`
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
