package orm

import (
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
)

// ConfigSchema records the schema of configuration at the type level
type ConfigSchema struct {
	AllowOrigins                              string          `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	BalanceMonitorEnabled                     bool            `env:"BALANCE_MONITOR_ENABLED" default:"true"`
	BlockBackfillDepth                        string          `env:"BLOCK_BACKFILL_DEPTH" default:"10"`
	BridgeResponseURL                         url.URL         `env:"BRIDGE_RESPONSE_URL"`
	ChainID                                   big.Int         `env:"ETH_CHAIN_ID" default:"1"`
	ClientNodeURL                             string          `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	DatabaseTimeout                           models.Duration `env:"DATABASE_TIMEOUT" default:"500ms"`
	DatabaseURL                               string          `env:"DATABASE_URL"`
	DatabaseListenerMinReconnectInterval      time.Duration   `env:"DATABASE_LISTENER_MIN_RECONNECT_INTERVAL" default:"1m"`
	DatabaseListenerMaxReconnectDuration      time.Duration   `env:"DATABASE_LISTENER_MAX_RECONNECT_DURATION" default:"10m"`
	DatabaseMaximumTxDuration                 time.Duration   `env:"DATABASE_MAXIMUM_TX_DURATION" default:"30m"`
	DefaultHTTPLimit                          int64           `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
	DefaultHTTPTimeout                        models.Duration `env:"DEFAULT_HTTP_TIMEOUT" default:"15s"`
	DefaultHTTPAllowUnrestrictedNetworkAccess bool            `env:"DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS" default:"false"`
	Dev                                       bool            `env:"CHAINLINK_DEV" default:"false"`
	EnableExperimentalAdapters                bool            `env:"ENABLE_EXPERIMENTAL_ADAPTERS" default:"false"`
	FeatureExternalInitiators                 bool            `env:"FEATURE_EXTERNAL_INITIATORS" default:"false"`
	FeatureFluxMonitor                        bool            `env:"FEATURE_FLUX_MONITOR" default:"true"`
	FeatureOffchainReporting                  bool            `env:"FEATURE_OFFCHAIN_REPORTING" default:"false"`
	MaximumServiceDuration                    models.Duration `env:"MAXIMUM_SERVICE_DURATION" default:"8760h" `
	MinimumServiceDuration                    models.Duration `env:"MINIMUM_SERVICE_DURATION" default:"0s" `
	EthGasBumpThreshold                       uint64          `env:"ETH_GAS_BUMP_THRESHOLD" default:"3" `
	EthGasBumpWei                             big.Int         `env:"ETH_GAS_BUMP_WEI" default:"5000000000"`
	EthGasBumpPercent                         uint16          `env:"ETH_GAS_BUMP_PERCENT" default:"20"`
	EthGasBumpTxDepth                         uint16          `env:"ETH_GAS_BUMP_TX_DEPTH" default:"10"`
	EthGasLimitDefault                        uint64          `env:"ETH_GAS_LIMIT_DEFAULT" default:"500000"`
	EthGasPriceDefault                        big.Int         `env:"ETH_GAS_PRICE_DEFAULT" default:"20000000000"`
	EthMaxGasPriceWei                         uint64          `env:"ETH_MAX_GAS_PRICE_WEI" default:"1500000000000"`
	EthFinalityDepth                          uint            `env:"ETH_FINALITY_DEPTH" default:"50"`
	EthHeadTrackerHistoryDepth                uint            `env:"ETH_HEAD_TRACKER_HISTORY_DEPTH" default:"100"`
	EthHeadTrackerMaxBufferSize               uint            `env:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE" default:"3"`
	EthBalanceMonitorBlockDelay               uint16          `env:"ETH_BALANCE_MONITOR_BLOCK_DELAY" default:"1"`
	EthereumURL                               string          `env:"ETH_URL" default:"ws://localhost:8546"`
	EthereumSecondaryURL                      string          `env:"ETH_SECONDARY_URL" default:""`
	EthereumDisabled                          bool            `env:"ETH_DISABLED" default:"false"`
	FlagsContractAddress                      string          `env:"FLAGS_CONTRACT_ADDRESS"`
	GasUpdaterBlockDelay                      uint16          `env:"GAS_UPDATER_BLOCK_DELAY" default:"3"`
	GasUpdaterBlockHistorySize                uint16          `env:"GAS_UPDATER_BLOCK_HISTORY_SIZE" default:"24"`
	GasUpdaterTransactionPercentile           uint16          `env:"GAS_UPDATER_TRANSACTION_PERCENTILE" default:"60"`
	GasUpdaterEnabled                         bool            `env:"GAS_UPDATER_ENABLED" default:"true"`
	InsecureFastScrypt                        bool            `env:"INSECURE_FAST_SCRYPT" default:"false"`
	JobPipelineMaxTaskDuration                time.Duration   `env:"JOB_PIPELINE_MAX_TASK_DURATION" default:"10m"`
	JobPipelineParallelism                    uint8           `env:"JOB_PIPELINE_PARALLELISM" default:"4"`
	JobPipelineReaperInterval                 time.Duration   `env:"JOB_PIPELINE_REAPER_INTERVAL" default:"1h"`
	JobPipelineReaperThreshold                time.Duration   `env:"JOB_PIPELINE_REAPER_THRESHOLD" default:"7d"`
	JSONConsole                               bool            `env:"JSON_CONSOLE" default:"false"`
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
	MinIncomingConfirmations                  uint32          `env:"MIN_INCOMING_CONFIRMATIONS" default:"3"`
	MinRequiredOutgoingConfirmations          uint64          `env:"MIN_OUTGOING_CONFIRMATIONS" default:"12"`
	MinimumContractPayment                    assets.Link     `env:"MINIMUM_CONTRACT_PAYMENT" default:"1000000000000000000"`
	MinimumRequestExpiration                  uint64          `env:"MINIMUM_REQUEST_EXPIRATION" default:"300"`
	OCRBootstrapCheckInterval                 time.Duration   `env:"OCR_BOOTSTRAP_CHECK_INTERVAL" default:"20s"`
	OCRContractTransmitterTransmitTimeout     time.Duration   `env:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT" default:"10s"`
	OCRDatabaseTimeout                        time.Duration   `env:"OCR_DATABASE_TIMEOUT" default:"10s"`
	OCRIncomingMessageBufferSize              int             `env:"OCR_INCOMING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCROutgoingMessageBufferSize              int             `env:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE" default:"10"`
	OCRNewStreamTimeout                       time.Duration   `env:"OCR_NEW_STREAM_TIMEOUT" default:"10s"`
	OCRDHTLookupInterval                      int             `env:"OCR_DHT_LOOKUP_INTERVAL" default:"10"`
	OCRTraceLogging                           bool            `env:"OCR_TRACE_LOGGING" default:"false"`
	OperatorContractAddress                   common.Address  `env:"OPERATOR_CONTRACT_ADDRESS"`
	P2PAnnounceIP                             net.IP          `env:"P2P_ANNOUNCE_IP"`
	P2PAnnouncePort                           uint16          `env:"P2P_ANNOUNCE_PORT"`
	P2PListenIP                               net.IP          `env:"P2P_LISTEN_IP" default:"0.0.0.0"`
	P2PListenPort                             uint16          `env:"P2P_LISTEN_PORT"`
	Port                                      uint16          `env:"CHAINLINK_PORT" default:"6688"`
	ReaperExpiration                          models.Duration `env:"REAPER_EXPIRATION" default:"240h"`
	ReplayFromBlock                           int64           `env:"REPLAY_FROM_BLOCK" default:"-1"`
	RootDir                                   string          `env:"ROOT" default:"~/.chainlink"`
	SecureCookies                             bool            `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout                            models.Duration `env:"SESSION_TIMEOUT" default:"15m"`
	TriggerFallbackDBPollInterval             time.Duration   `env:"TRIGGER_FALLBACK_DB_POLL_INTERVAL" default:"30s"`
	TLSCertPath                               string          `env:"TLS_CERT_PATH" `
	TLSHost                                   string          `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath                                string          `env:"TLS_KEY_PATH" `
	TLSPort                                   uint16          `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TLSRedirect                               bool            `env:"CHAINLINK_TLS_REDIRECT" default:"false"`
	TxAttemptLimit                            uint16          `env:"CHAINLINK_TX_ATTEMPT_LIMIT" default:"10"`
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
