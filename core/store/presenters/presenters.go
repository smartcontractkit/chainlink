// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"bytes"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/config"
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
	ChainID                                    *big.Int        `json:"ETH_CHAIN_ID"`
	ClientNodeURL                              string          `json:"CLIENT_NODE_URL"`
	DatabaseBackupFrequency                    time.Duration   `json:"DATABASE_BACKUP_FREQUENCY"`
	DatabaseBackupMode                         string          `json:"DATABASE_BACKUP_MODE"`
	DatabaseMaximumTxDuration                  time.Duration   `json:"DATABASE_MAXIMUM_TX_DURATION"`
	DatabaseTimeout                            models.Duration `json:"DATABASE_TIMEOUT"`
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
	KeeperMaximumGracePeriod                   int64           `json:"KEEPER_MAXIMUM_GRACE_PERIOD"`
	KeeperMinimumRequiredConfirmations         uint64          `json:"KEEPER_MINIMUM_REQUIRED_CONFIRMATIONS"`
	KeeperRegistryCheckGasOverhead             uint64          `json:"KEEPER_REGISTRY_CHECK_GAS_OVERHEAD"`
	KeeperRegistryPerformGasOverhead           uint64          `json:"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD"`
	KeeperRegistrySyncInterval                 time.Duration   `json:"KEEPER_REGISTRY_SYNC_INTERVAL"`
	LinkContractAddress                        string          `json:"LINK_CONTRACT_ADDRESS"`
	FlagsContractAddress                       string          `json:"FLAGS_CONTRACT_ADDRESS"`
	Layer2Type                                 string          `json:"LAYER_2_TYPE"`
	LogLevel                                   config.LogLevel `json:"LOG_LEVEL"`
	LogSQLMigrations                           bool            `json:"LOG_SQL_MIGRATIONS"`
	LogSQLStatements                           bool            `json:"LOG_SQL"`
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
func NewConfigPrinter(config config.GeneralConfig) (ConfigPrinter, error) {
	explorerURL := ""
	if config.ExplorerURL() != nil {
		explorerURL = config.ExplorerURL().String()
	}
	p2pBootstrapPeers, _ := config.P2PBootstrapPeers()
	ethereumHTTPURL := ""
	if config.EthereumHTTPURL() != nil {
		ethereumHTTPURL = config.EthereumHTTPURL().String()
	}
	telemetryIngressURL := ""
	if config.TelemetryIngressURL() != nil {
		telemetryIngressURL = config.TelemetryIngressURL().String()
	}
	return ConfigPrinter{
		EnvPrinter: EnvPrinter{
			AllowOrigins:                          config.AllowOrigins(),
			BlockBackfillDepth:                    config.BlockBackfillDepth(),
			BridgeResponseURL:                     config.BridgeResponseURL().String(),
			ChainID:                               config.ChainID(),
			ClientNodeURL:                         config.ClientNodeURL(),
			DatabaseBackupFrequency:               config.DatabaseBackupFrequency(),
			DatabaseBackupMode:                    string(config.DatabaseBackupMode()),
			DatabaseMaximumTxDuration:             config.DatabaseMaximumTxDuration(),
			DatabaseTimeout:                       config.DatabaseTimeout(),
			DefaultHTTPLimit:                      config.DefaultHTTPLimit(),
			DefaultHTTPTimeout:                    config.DefaultHTTPTimeout(),
			Dev:                                   config.Dev(),
			EthereumDisabled:                      config.EthereumDisabled(),
			EthereumHTTPURL:                       ethereumHTTPURL,
			EthereumSecondaryURLs:                 mapToStringA(config.EthereumSecondaryURLs()),
			EthereumURL:                           config.EthereumURL(),
			ExplorerURL:                           explorerURL,
			FMDefaultTransactionQueueDepth:        config.FMDefaultTransactionQueueDepth(),
			FeatureExternalInitiators:             config.FeatureExternalInitiators(),
			FeatureOffchainReporting:              config.FeatureOffchainReporting(),
			InsecureFastScrypt:                    config.InsecureFastScrypt(),
			JSONConsole:                           config.JSONConsole(),
			JobPipelineReaperInterval:             config.JobPipelineReaperInterval(),
			JobPipelineReaperThreshold:            config.JobPipelineReaperThreshold(),
			KeeperDefaultTransactionQueueDepth:    config.KeeperDefaultTransactionQueueDepth(),
			LogLevel:                              config.LogLevel(),
			LogSQLMigrations:                      config.LogSQLMigrations(),
			LogSQLStatements:                      config.LogSQLStatements(),
			LogToDisk:                             config.LogToDisk(),
			OCRBootstrapCheckInterval:             config.OCRBootstrapCheckInterval(),
			OCRContractTransmitterTransmitTimeout: config.OCRContractTransmitterTransmitTimeout(),
			OCRDHTLookupInterval:                  config.OCRDHTLookupInterval(),
			OCRDatabaseTimeout:                    config.OCRDatabaseTimeout(),
			OCRDefaultTransactionQueueDepth:       config.OCRDefaultTransactionQueueDepth(),
			OCRIncomingMessageBufferSize:          config.OCRIncomingMessageBufferSize(),
			OCRNewStreamTimeout:                   config.OCRNewStreamTimeout(),
			OCROutgoingMessageBufferSize:          config.OCROutgoingMessageBufferSize(),
			OCRTraceLogging:                       config.OCRTraceLogging(),
			P2PBootstrapPeers:                     p2pBootstrapPeers,
			P2PListenIP:                           config.P2PListenIP().String(),
			P2PListenPort:                         config.P2PListenPortRaw(),
			P2PNetworkingStack:                    config.P2PNetworkingStackRaw(),
			P2PPeerID:                             config.P2PPeerIDRaw(),
			P2PV2AnnounceAddresses:                config.P2PV2AnnounceAddressesRaw(),
			P2PV2Bootstrappers:                    config.P2PV2BootstrappersRaw(),
			P2PV2DeltaDial:                        config.P2PV2DeltaDial(),
			P2PV2DeltaReconcile:                   config.P2PV2DeltaReconcile(),
			P2PV2ListenAddresses:                  config.P2PV2ListenAddresses(),
			Port:                                  config.Port(),
			ReaperExpiration:                      config.ReaperExpiration(),
			ReplayFromBlock:                       config.ReplayFromBlock(),
			RootDir:                               config.RootDir(),
			SecureCookies:                         config.SecureCookies(),
			SessionTimeout:                        config.SessionTimeout(),
			TLSHost:                               config.TLSHost(),
			TLSPort:                               config.TLSPort(),
			TLSRedirect:                           config.TLSRedirect(),
			TelemetryIngressLogging:               config.TelemetryIngressLogging(),
			TelemetryIngressServerPubKey:          config.TelemetryIngressServerPubKey(),
			TelemetryIngressURL:                   telemetryIngressURL,
			TriggerFallbackDBPollInterval:         config.TriggerFallbackDBPollInterval(),
		},
	}, nil
}

// String returns the values as a newline delimited string
func (c ConfigPrinter) String() string {
	var buffer bytes.Buffer

	schemaT := reflect.TypeOf(config.ConfigSchema{})
	cwlT := reflect.TypeOf(c.EnvPrinter)
	cwlV := reflect.ValueOf(c.EnvPrinter)

	for index := 0; index < cwlT.NumField(); index++ {
		item := cwlT.FieldByIndex([]int{index})
		schemaItem, ok := schemaT.FieldByName(item.Name)
		if !ok {
			logger.Panicf("Field %s missing from store.Schema", item.Name)
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

// ExternalInitiatorAuthentication includes initiator and authentication details.
type ExternalInitiatorAuthentication struct {
	Name           string        `json:"name,omitempty"`
	URL            models.WebURL `json:"url,omitempty"`
	AccessKey      string        `json:"incomingAccessKey,omitempty"`
	Secret         string        `json:"incomingSecret,omitempty"`
	OutgoingToken  string        `json:"outgoingToken,omitempty"`
	OutgoingSecret string        `json:"outgoingSecret,omitempty"`
}

// NewExternalInitiatorAuthentication creates an instance of ExternalInitiatorAuthentication.
func NewExternalInitiatorAuthentication(
	ei models.ExternalInitiator,
	eia auth.Token,
) *ExternalInitiatorAuthentication {
	var result = &ExternalInitiatorAuthentication{
		Name:           ei.Name,
		AccessKey:      ei.AccessKey,
		Secret:         eia.Secret,
		OutgoingToken:  ei.OutgoingToken,
		OutgoingSecret: ei.OutgoingSecret,
	}
	if ei.URL != nil {
		result.URL = *ei.URL
	}
	return result
}

// GetID returns the jsonapi ID.
func (ei *ExternalInitiatorAuthentication) GetID() string {
	return ei.Name
}

// GetName returns the collection name for jsonapi.
func (*ExternalInitiatorAuthentication) GetName() string {
	return "external initiators"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (ei *ExternalInitiatorAuthentication) SetID(name string) error {
	ei.Name = name
	return nil
}
