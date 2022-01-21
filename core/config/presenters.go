// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package config

import (
	"bytes"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ConfigPrinter are the non-secret values of the node
//
// If you add an entry here, you should update NewConfigPrinter and
// ConfigPrinter#String accordingly.
type ConfigPrinter struct {
	EnvPrinter
}

// EnvPrinter contains the supported environment variables
type EnvPrinter struct {
	AllowOrigins                               string          `json:"ALLOW_ORIGINS"`
	BlockBackfillDepth                         uint64          `json:"BLOCK_BACKFILL_DEPTH"`
	BlockHistoryEstimatorBlockDelay            uint16          `json:"GAS_UPDATER_BLOCK_DELAY"`
	BlockHistoryEstimatorBlockHistorySize      uint16          `json:"GAS_UPDATER_BLOCK_HISTORY_SIZE"`
	BlockHistoryEstimatorTransactionPercentile uint16          `json:"GAS_UPDATER_TRANSACTION_PERCENTILE"`
	BridgeResponseURL                          string          `json:"BRIDGE_RESPONSE_URL,omitempty"`
	ChainType                                  string          `json:"CHAIN_TYPE"`
	ClientNodeURL                              string          `json:"CLIENT_NODE_URL"`
	DatabaseBackupFrequency                    time.Duration   `json:"DATABASE_BACKUP_FREQUENCY"`
	DatabaseBackupMode                         string          `json:"DATABASE_BACKUP_MODE"`
	DatabaseBackupOnVersionUpgrade             bool            `json:"DATABASE_BACKUP_ON_VERSION_UPGRADE"`
	DatabaseMaximumTxDuration                  time.Duration   `json:"DATABASE_MAXIMUM_TX_DURATION"`
	DatabaseTimeout                            models.Duration `json:"DATABASE_TIMEOUT"`
	DatabaseLockingMode                        string          `json:"DATABASE_LOCKING_MODE"`
	DefaultChainID                             string          `json:"ETH_CHAIN_ID"`
	DefaultHTTPLimit                           int64           `json:"DEFAULT_HTTP_LIMIT"`
	DefaultHTTPTimeout                         models.Duration `json:"DEFAULT_HTTP_TIMEOUT"`
	Dev                                        bool            `json:"CHAINLINK_DEV"`
	EthereumDisabled                           bool            `json:"ETH_DISABLED"`
	EthereumHTTPURL                            string          `json:"ETH_HTTP_URL"`
	EthereumSecondaryURLs                      []string        `json:"ETH_SECONDARY_URLS"`
	EthereumURL                                string          `json:"ETH_URL"`
	ExplorerURL                                string          `json:"EXPLORER_URL"`
	FMDefaultTransactionQueueDepth             uint32          `json:"FM_DEFAULT_TRANSACTION_QUEUE_DEPTH"`
	FeatureExternalInitiators                  bool            `json:"FEATURE_EXTERNAL_INITIATORS"`
	FeatureOffchainReporting                   bool            `json:"FEATURE_OFFCHAIN_REPORTING"`
	GasEstimatorMode                           string          `json:"GAS_ESTIMATOR_MODE"`
	InsecureFastScrypt                         bool            `json:"INSECURE_FAST_SCRYPT"`
	JSONConsole                                bool            `json:"JSON_CONSOLE"`
	JobPipelineReaperInterval                  time.Duration   `json:"JOB_PIPELINE_REAPER_INTERVAL"`
	JobPipelineReaperThreshold                 time.Duration   `json:"JOB_PIPELINE_REAPER_THRESHOLD"`
	KeeperDefaultTransactionQueueDepth         uint32          `json:"KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH"`
	KeeperGasPriceBufferPercent                uint32          `json:"KEEPER_GAS_PRICE_BUFFER_PERCENT"`
	KeeperGasTipCapBufferPercent               uint32          `json:"KEEPER_GAS_TIP_CAP_BUFFER_PERCENT"`
	KeeperMaximumGracePeriod                   int64           `json:"KEEPER_MAXIMUM_GRACE_PERIOD"`
	KeeperRegistryCheckGasOverhead             uint64          `json:"KEEPER_REGISTRY_CHECK_GAS_OVERHEAD"`
	KeeperRegistryPerformGasOverhead           uint64          `json:"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD"`
	KeeperRegistrySyncInterval                 time.Duration   `json:"KEEPER_REGISTRY_SYNC_INTERVAL"`
	KeeperRegistrySyncUpkeepQueueSize          uint32          `json:"KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE"`
	KeeperCheckUpkeepGasPriceFeatureEnabled    bool            `json:"KEEPER_CHECK_UPKEEP_GAS_PRICE_FEATURE_ENABLED"`
	LeaseLockDuration                          time.Duration   `json:"LEASE_LOCK_DURATION"`
	LeaseLockRefreshInterval                   time.Duration   `json:"LEASE_LOCK_REFRESH_INTERVAL"`
	LinkContractAddress                        string          `json:"LINK_CONTRACT_ADDRESS"`
	FlagsContractAddress                       string          `json:"FLAGS_CONTRACT_ADDRESS"`
	LogLevel                                   LogLevel        `json:"LOG_LEVEL"`
	LogSQLMigrations                           bool            `json:"LOG_SQL_MIGRATIONS"`
	LogSQL                                     bool            `json:"LOG_SQL"`
	LogToDisk                                  bool            `json:"LOG_TO_DISK"`
	OCRBootstrapCheckInterval                  time.Duration   `json:"OCR_BOOTSTRAP_CHECK_INTERVAL"`
	TriggerFallbackDBPollInterval              time.Duration   `json:"JOB_PIPELINE_DB_POLL_INTERVAL"`
	OCRContractTransmitterTransmitTimeout      time.Duration   `json:"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT"`
	OCRDatabaseTimeout                         time.Duration   `json:"OCR_DATABASE_TIMEOUT"`
	OCRDefaultTransactionQueueDepth            uint32          `json:"OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH"`
	OCRIncomingMessageBufferSize               int             `json:"OCR_INCOMING_MESSAGE_BUFFER_SIZE"`
	P2PBootstrapPeers                          []string        `json:"P2P_BOOTSTRAP_PEERS"`
	P2PListenIP                                string          `json:"P2P_LISTEN_IP"`
	P2PListenPort                              string          `json:"P2P_LISTEN_PORT"`
	P2PNetworkingStack                         string          `json:"P2P_NETWORKING_STACK"`
	P2PPeerID                                  string          `json:"P2P_PEER_ID"`
	P2PV2AnnounceAddresses                     []string        `json:"P2PV2_ANNOUNCE_ADDRESSES"`
	P2PV2Bootstrappers                         []string        `json:"P2PV2_BOOTSTRAPPERS"`
	P2PV2DeltaDial                             models.Duration `json:"P2PV2_DELTA_DIAL"`
	P2PV2DeltaReconcile                        models.Duration `json:"P2PV2_DELTA_RECONCILE"`
	P2PV2ListenAddresses                       []string        `json:"P2PV2_LISTEN_ADDRESSES"`
	OCROutgoingMessageBufferSize               int             `json:"OCR_OUTGOING_MESSAGE_BUFFER_SIZE"`
	OCRNewStreamTimeout                        time.Duration   `json:"OCR_NEW_STREAM_TIMEOUT"`
	OCRDHTLookupInterval                       int             `json:"OCR_DHT_LOOKUP_INTERVAL"`
	OCRTraceLogging                            bool            `json:"OCR_TRACE_LOGGING"`
	Port                                       uint16          `json:"CHAINLINK_PORT"`
	ReaperExpiration                           models.Duration `json:"REAPER_EXPIRATION"`
	ReplayFromBlock                            int64           `json:"REPLAY_FROM_BLOCK"`
	RootDir                                    string          `json:"ROOT"`
	SecureCookies                              bool            `json:"SECURE_COOKIES"`
	SessionTimeout                             models.Duration `json:"SESSION_TIMEOUT"`
	TelemetryIngressLogging                    bool            `json:"TELEMETRY_INGRESS_LOGGING"`
	TelemetryIngressServerPubKey               string          `json:"TELEMETRY_INGRESS_SERVER_PUB_KEY"`
	TelemetryIngressURL                        string          `json:"TELEMETRY_INGRESS_URL"`
	TLSHost                                    string          `json:"CHAINLINK_TLS_HOST"`
	TLSPort                                    uint16          `json:"CHAINLINK_TLS_PORT"`
	TLSRedirect                                bool            `json:"CHAINLINK_TLS_REDIRECT"`
}

// NewConfigPrinter creates an instance of ConfigPrinter
func NewConfigPrinter(cfg GeneralConfig) (ConfigPrinter, error) {
	explorerURL := ""
	if cfg.ExplorerURL() != nil {
		explorerURL = cfg.ExplorerURL().String()
	}
	p2pBootstrapPeers, _ := cfg.P2PBootstrapPeers()
	ethereumHTTPURL := ""
	if cfg.EthereumHTTPURL() != nil {
		ethereumHTTPURL = cfg.EthereumHTTPURL().String()
	}
	telemetryIngressURL := ""
	if cfg.TelemetryIngressURL() != nil {
		telemetryIngressURL = cfg.TelemetryIngressURL().String()
	}
	return ConfigPrinter{
		EnvPrinter: EnvPrinter{
			AllowOrigins:                          cfg.AllowOrigins(),
			BlockBackfillDepth:                    cfg.BlockBackfillDepth(),
			BridgeResponseURL:                     cfg.BridgeResponseURL().String(),
			ClientNodeURL:                         cfg.ClientNodeURL(),
			DatabaseBackupFrequency:               cfg.DatabaseBackupFrequency(),
			DatabaseBackupMode:                    string(cfg.DatabaseBackupMode()),
			DatabaseBackupOnVersionUpgrade:        cfg.DatabaseBackupOnVersionUpgrade(),
			DatabaseMaximumTxDuration:             cfg.DatabaseMaximumTxDuration(),
			DatabaseTimeout:                       cfg.DatabaseTimeout(),
			DatabaseLockingMode:                   cfg.DatabaseLockingMode(),
			DefaultChainID:                        cfg.DefaultChainID().String(),
			DefaultHTTPLimit:                      cfg.DefaultHTTPLimit(),
			DefaultHTTPTimeout:                    cfg.DefaultHTTPTimeout(),
			Dev:                                   cfg.Dev(),
			EthereumDisabled:                      cfg.EthereumDisabled(),
			EthereumHTTPURL:                       ethereumHTTPURL,
			EthereumSecondaryURLs:                 mapToStringA(cfg.EthereumSecondaryURLs()),
			EthereumURL:                           cfg.EthereumURL(),
			ExplorerURL:                           explorerURL,
			FMDefaultTransactionQueueDepth:        cfg.FMDefaultTransactionQueueDepth(),
			FeatureExternalInitiators:             cfg.FeatureExternalInitiators(),
			FeatureOffchainReporting:              cfg.FeatureOffchainReporting(),
			InsecureFastScrypt:                    cfg.InsecureFastScrypt(),
			JSONConsole:                           cfg.JSONConsole(),
			JobPipelineReaperInterval:             cfg.JobPipelineReaperInterval(),
			JobPipelineReaperThreshold:            cfg.JobPipelineReaperThreshold(),
			KeeperDefaultTransactionQueueDepth:    cfg.KeeperDefaultTransactionQueueDepth(),
			KeeperGasPriceBufferPercent:           cfg.KeeperGasPriceBufferPercent(),
			KeeperGasTipCapBufferPercent:          cfg.KeeperGasTipCapBufferPercent(),
			LeaseLockDuration:                     cfg.LeaseLockDuration(),
			LeaseLockRefreshInterval:              cfg.LeaseLockRefreshInterval(),
			LogLevel:                              LogLevel{Level: cfg.LogLevel()},
			LogSQLMigrations:                      cfg.LogSQLMigrations(),
			LogSQL:                                cfg.LogSQL(),
			LogToDisk:                             cfg.LogToDisk(),
			OCRBootstrapCheckInterval:             cfg.OCRBootstrapCheckInterval(),
			OCRContractTransmitterTransmitTimeout: cfg.OCRContractTransmitterTransmitTimeout(),
			OCRDHTLookupInterval:                  cfg.OCRDHTLookupInterval(),
			OCRDatabaseTimeout:                    cfg.OCRDatabaseTimeout(),
			OCRDefaultTransactionQueueDepth:       cfg.OCRDefaultTransactionQueueDepth(),
			OCRIncomingMessageBufferSize:          cfg.OCRIncomingMessageBufferSize(),
			OCRNewStreamTimeout:                   cfg.OCRNewStreamTimeout(),
			OCROutgoingMessageBufferSize:          cfg.OCROutgoingMessageBufferSize(),
			OCRTraceLogging:                       cfg.OCRTraceLogging(),
			P2PBootstrapPeers:                     p2pBootstrapPeers,
			P2PListenIP:                           cfg.P2PListenIP().String(),
			P2PListenPort:                         cfg.P2PListenPortRaw(),
			P2PNetworkingStack:                    cfg.P2PNetworkingStackRaw(),
			P2PPeerID:                             cfg.P2PPeerIDRaw(),
			P2PV2AnnounceAddresses:                cfg.P2PV2AnnounceAddressesRaw(),
			P2PV2Bootstrappers:                    cfg.P2PV2BootstrappersRaw(),
			P2PV2DeltaDial:                        cfg.P2PV2DeltaDial(),
			P2PV2DeltaReconcile:                   cfg.P2PV2DeltaReconcile(),
			P2PV2ListenAddresses:                  cfg.P2PV2ListenAddresses(),
			Port:                                  cfg.Port(),
			ReaperExpiration:                      cfg.ReaperExpiration(),
			ReplayFromBlock:                       cfg.ReplayFromBlock(),
			RootDir:                               cfg.RootDir(),
			SecureCookies:                         cfg.SecureCookies(),
			SessionTimeout:                        cfg.SessionTimeout(),
			TLSHost:                               cfg.TLSHost(),
			TLSPort:                               cfg.TLSPort(),
			TLSRedirect:                           cfg.TLSRedirect(),
			TelemetryIngressLogging:               cfg.TelemetryIngressLogging(),
			TelemetryIngressServerPubKey:          cfg.TelemetryIngressServerPubKey(),
			TelemetryIngressURL:                   telemetryIngressURL,
			TriggerFallbackDBPollInterval:         cfg.TriggerFallbackDBPollInterval(),
		},
	}, nil
}

// String returns the values as a newline delimited string
func (c ConfigPrinter) String() string {
	var buffer bytes.Buffer

	schemaT := reflect.TypeOf(ConfigSchema{})
	cwlT := reflect.TypeOf(c.EnvPrinter)
	cwlV := reflect.ValueOf(c.EnvPrinter)

	for index := 0; index < cwlT.NumField(); index++ {
		item := cwlT.FieldByIndex([]int{index})
		schemaItem, ok := schemaT.FieldByName(item.Name)
		if !ok {
			panic(fmt.Sprintf("Field %s missing from store.Schema", item.Name))
		}
		envName, ok := schemaItem.Tag.Lookup("env")
		if !ok {
			continue
		}

		field := cwlV.FieldByIndex(item.Index)

		buffer.WriteString(envName)
		buffer.WriteString(": ")
		if stringer, ok := field.Interface().(fmt.Stringer); ok {
			if stringer != reflect.Zero(reflect.TypeOf(stringer)).Interface() {
				buffer.WriteString(stringer.String())
			}
		} else {
			buffer.WriteString(fmt.Sprintf("%v", field))
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// GetID generates a new ID for jsonapi serialization.
func (c ConfigPrinter) GetID() string {
	return utils.NewBytes32ID()
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (c *ConfigPrinter) SetID(value string) error {
	return nil
}

func mapToStringA(in []url.URL) (out []string) {
	for _, url := range in {
		out = append(out, url.String())
	}
	return
}

// FriendlyBigInt returns a string printing the integer in both
// decimal and hexadecimal formats.
func FriendlyBigInt(n *big.Int) string {
	return fmt.Sprintf("#%[1]v (0x%[1]x)", n)
}
