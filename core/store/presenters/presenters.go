// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tidwall/gjson"
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
	BalanceMonitorEnabled                      bool            `json:"BALANCE_MONITOR_ENABLED"`
	BlockBackfillDepth                         uint64          `json:"BLOCK_BACKFILL_DEPTH"`
	BlockBackfillSkip                          bool            `json:"BLOCK_BACKFILL_SKIP"`
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
	EnableExperimentalAdapters                 bool            `json:"ENABLE_EXPERIMENTAL_ADAPTERS"`
	EnableLegacyJobPipeline                    bool            `json:"ENABLE_LEGACY_JOB_PIPELINE"`
	EthBalanceMonitorBlockDelay                uint16          `json:"ETH_BALANCE_MONITOR_BLOCK_DELAY"`
	EthFinalityDepth                           uint            `json:"ETH_FINALITY_DEPTH"`
	EthGasBumpThreshold                        uint64          `json:"ETH_GAS_BUMP_THRESHOLD"`
	EthGasBumpTxDepth                          uint16          `json:"ETH_GAS_BUMP_TX_DEPTH"`
	EthGasBumpWei                              *big.Int        `json:"ETH_GAS_BUMP_WEI"`
	EthGasLimitDefault                         uint64          `json:"ETH_GAS_LIMIT_DEFAULT"`
	EthGasLimitTransfer                        uint64          `json:"ETH_GAS_LIMIT_TRANSFER"`
	EthGasPriceDefault                         *big.Int        `json:"ETH_GAS_PRICE_DEFAULT"`
	EthHeadTrackerHistoryDepth                 uint            `json:"ETH_HEAD_TRACKER_HISTORY_DEPTH"`
	EthHeadTrackerMaxBufferSize                uint            `json:"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE"`
	EthMaxGasPriceWei                          *big.Int        `json:"ETH_MAX_GAS_PRICE_WEI"`
	EthereumDisabled                           bool            `json:"ETH_DISABLED"`
	EthereumHTTPURL                            string          `json:"ETH_HTTP_URL"`
	EthereumSecondaryURLs                      []string        `json:"ETH_SECONDARY_URLS"`
	EthereumURL                                string          `json:"ETH_URL"`
	ExplorerURL                                string          `json:"EXPLORER_URL"`
	FMDefaultTransactionQueueDepth             uint32          `json:"FM_DEFAULT_TRANSACTION_QUEUE_DEPTH"`
	FeatureExternalInitiators                  bool            `json:"FEATURE_EXTERNAL_INITIATORS"`
	FeatureFluxMonitor                         bool            `json:"FEATURE_FLUX_MONITOR"`
	FeatureOffchainReporting                   bool            `json:"FEATURE_OFFCHAIN_REPORTING"`
	FlagsContractAddress                       string          `json:"FLAGS_CONTRACT_ADDRESS"`
	GasEstimatorMode                           string          `json:"GAS_ESTIMATOR_MODE"`
	InsecureFastScrypt                         bool            `json:"INSECURE_FAST_SCRYPT"`
	JSONConsole                                bool            `json:"JSON_CONSOLE"`
	JobPipelineReaperInterval                  time.Duration   `json:"JOB_PIPELINE_REAPER_INTERVAL"`
	JobPipelineReaperThreshold                 time.Duration   `json:"JOB_PIPELINE_REAPER_THRESHOLD"`
	KeeperDefaultTransactionQueueDepth         uint32          `json:"KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH"`
	LinkContractAddress                        string          `json:"LINK_CONTRACT_ADDRESS"`
	LogLevel                                   config.LogLevel `json:"LOG_LEVEL"`
	LogSQLMigrations                           bool            `json:"LOG_SQL_MIGRATIONS"`
	LogSQLStatements                           bool            `json:"LOG_SQL"`
	LogToDisk                                  bool            `json:"LOG_TO_DISK"`
	MaximumServiceDuration                     models.Duration `json:"MAXIMUM_SERVICE_DURATION"`
	MinIncomingConfirmations                   uint32          `json:"MIN_INCOMING_CONFIRMATIONS"`
	MinRequiredOutgoingConfirmations           uint64          `json:"MIN_OUTGOING_CONFIRMATIONS"`
	MinimumContractPayment                     *assets.Link    `json:"MINIMUM_CONTRACT_PAYMENT_LINK_JUELS"`
	MinimumRequestExpiration                   uint64          `json:"MINIMUM_REQUEST_EXPIRATION"`
	MinimumServiceDuration                     models.Duration `json:"MINIMUM_SERVICE_DURATION"`
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
	OperatorContractAddress                    common.Address  `json:"OPERATOR_CONTRACT_ADDRESS"`
	Port                                       uint16          `json:"CHAINLINK_PORT"`
	ReaperExpiration                           models.Duration `json:"REAPER_EXPIRATION"`
	ReplayFromBlock                            int64           `json:"REPLAY_FROM_BLOCK"`
	RootDir                                    string          `json:"ROOT"`
	SecureCookies                              bool            `json:"SECURE_COOKIES"`
	SessionTimeout                             models.Duration `json:"SESSION_TIMEOUT"`
	TLSHost                                    string          `json:"CHAINLINK_TLS_HOST"`
	TLSPort                                    uint16          `json:"CHAINLINK_TLS_PORT"`
	TLSRedirect                                bool            `json:"CHAINLINK_TLS_REDIRECT"`
}

// NewConfigPrinter creates an instance of ConfigPrinter
func NewConfigPrinter(store *store.Store) (ConfigPrinter, error) {
	config := store.Config

	explorerURL := ""
	if config.ExplorerURL() != nil {
		explorerURL = config.ExplorerURL().String()
	}
	p2pBootstrapPeers, _ := config.P2PBootstrapPeers(nil)
	ethereumHTTPURL := ""
	if config.EthereumHTTPURL() != nil {
		ethereumHTTPURL = config.EthereumHTTPURL().String()
	}
	return ConfigPrinter{
		EnvPrinter: EnvPrinter{
			AllowOrigins:                               config.AllowOrigins(),
			BalanceMonitorEnabled:                      config.BalanceMonitorEnabled(),
			BlockBackfillDepth:                         config.BlockBackfillDepth(),
			BlockBackfillSkip:                          config.BlockBackfillSkip(),
			BlockHistoryEstimatorBlockDelay:            config.BlockHistoryEstimatorBlockDelay(),
			BlockHistoryEstimatorBlockHistorySize:      config.BlockHistoryEstimatorBlockHistorySize(),
			BlockHistoryEstimatorTransactionPercentile: config.BlockHistoryEstimatorTransactionPercentile(),
			BridgeResponseURL:                          config.BridgeResponseURL().String(),
			ChainID:                                    config.ChainID(),
			ClientNodeURL:                              config.ClientNodeURL(),
			DatabaseBackupFrequency:                    config.DatabaseBackupFrequency(),
			DatabaseBackupMode:                         string(config.DatabaseBackupMode()),
			DatabaseMaximumTxDuration:                  config.DatabaseMaximumTxDuration(),
			DatabaseTimeout:                            config.DatabaseTimeout(),
			DefaultHTTPLimit:                           config.DefaultHTTPLimit(),
			DefaultHTTPTimeout:                         config.DefaultHTTPTimeout(),
			Dev:                                        config.Dev(),
			EnableExperimentalAdapters:                 config.EnableExperimentalAdapters(),
			EnableLegacyJobPipeline:                    config.EnableLegacyJobPipeline(),
			EthBalanceMonitorBlockDelay:                config.EthBalanceMonitorBlockDelay(),
			EthFinalityDepth:                           config.EthFinalityDepth(),
			EthGasBumpThreshold:                        config.EthGasBumpThreshold(),
			EthGasBumpTxDepth:                          config.EthGasBumpTxDepth(),
			EthGasBumpWei:                              config.EthGasBumpWei(),
			EthGasLimitDefault:                         config.EthGasLimitDefault(),
			EthGasLimitTransfer:                        config.EthGasLimitTransfer(),
			EthGasPriceDefault:                         config.EthGasPriceDefault(),
			EthHeadTrackerHistoryDepth:                 config.EthHeadTrackerHistoryDepth(),
			EthHeadTrackerMaxBufferSize:                config.EthHeadTrackerMaxBufferSize(),
			EthMaxGasPriceWei:                          config.EthMaxGasPriceWei(),
			EthereumDisabled:                           config.EthereumDisabled(),
			EthereumHTTPURL:                            ethereumHTTPURL,
			EthereumSecondaryURLs:                      mapToStringA(config.EthereumSecondaryURLs()),
			EthereumURL:                                config.EthereumURL(),
			ExplorerURL:                                explorerURL,
			FMDefaultTransactionQueueDepth:             config.FMDefaultTransactionQueueDepth(),
			FeatureExternalInitiators:                  config.FeatureExternalInitiators(),
			FeatureFluxMonitor:                         config.FeatureFluxMonitor(),
			FeatureOffchainReporting:                   config.FeatureOffchainReporting(),
			FlagsContractAddress:                       config.FlagsContractAddress(),
			GasEstimatorMode:                           config.GasEstimatorMode(),
			InsecureFastScrypt:                         config.InsecureFastScrypt(),
			JSONConsole:                                config.JSONConsole(),
			JobPipelineReaperInterval:                  config.JobPipelineReaperInterval(),
			JobPipelineReaperThreshold:                 config.JobPipelineReaperThreshold(),
			KeeperDefaultTransactionQueueDepth:         config.KeeperDefaultTransactionQueueDepth(),
			LinkContractAddress:                        config.LinkContractAddress(),
			LogLevel:                                   config.LogLevel(),
			LogSQLMigrations:                           config.LogSQLMigrations(),
			LogSQLStatements:                           config.LogSQLStatements(),
			LogToDisk:                                  config.LogToDisk(),
			MaximumServiceDuration:                     config.MaximumServiceDuration(),
			MinIncomingConfirmations:                   config.MinIncomingConfirmations(),
			MinRequiredOutgoingConfirmations:           config.MinRequiredOutgoingConfirmations(),
			MinimumContractPayment:                     config.MinimumContractPayment(),
			MinimumRequestExpiration:                   config.MinimumRequestExpiration(),
			MinimumServiceDuration:                     config.MinimumServiceDuration(),
			OCRBootstrapCheckInterval:                  config.OCRBootstrapCheckInterval(),
			OCRContractTransmitterTransmitTimeout:      config.OCRContractTransmitterTransmitTimeout(),
			OCRDHTLookupInterval:                       config.OCRDHTLookupInterval(),
			OCRDatabaseTimeout:                         config.OCRDatabaseTimeout(),
			OCRDefaultTransactionQueueDepth:            config.OCRDefaultTransactionQueueDepth(),
			OCRIncomingMessageBufferSize:               config.OCRIncomingMessageBufferSize(),
			OCRNewStreamTimeout:                        config.OCRNewStreamTimeout(),
			OCROutgoingMessageBufferSize:               config.OCROutgoingMessageBufferSize(),
			OCRTraceLogging:                            config.OCRTraceLogging(),
			P2PBootstrapPeers:                          p2pBootstrapPeers,
			P2PListenIP:                                config.P2PListenIP().String(),
			P2PListenPort:                              config.P2PListenPortRaw(),
			P2PNetworkingStack:                         config.P2PNetworkingStackRaw(),
			P2PPeerID:                                  config.P2PPeerIDRaw(),
			P2PV2AnnounceAddresses:                     config.P2PV2AnnounceAddressesRaw(),
			P2PV2Bootstrappers:                         config.P2PV2BootstrappersRaw(),
			P2PV2DeltaDial:                             config.P2PV2DeltaDial(),
			P2PV2DeltaReconcile:                        config.P2PV2DeltaReconcile(),
			P2PV2ListenAddresses:                       config.P2PV2ListenAddresses(),
			TriggerFallbackDBPollInterval:              config.TriggerFallbackDBPollInterval(),
			OperatorContractAddress:                    config.OperatorContractAddress(),
			Port:                                       config.Port(),
			ReaperExpiration:                           config.ReaperExpiration(),
			ReplayFromBlock:                            config.ReplayFromBlock(),
			RootDir:                                    config.RootDir(),
			SecureCookies:                              config.SecureCookies(),
			SessionTimeout:                             config.SessionTimeout(),
			TLSHost:                                    config.TLSHost(),
			TLSPort:                                    config.TLSPort(),
			TLSRedirect:                                config.TLSRedirect(),
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

// JobSpec holds the JobSpec definition together with
// the total link earned from that job
type JobSpec struct {
	models.JobSpec
	Errors   []models.JobSpecError `json:"errors"`
	Earnings *assets.Link          `json:"earnings"`
}

// MarshalJSON returns the JSON data of the Job and its Initiators.
func (job JobSpec) MarshalJSON() ([]byte, error) {
	type Alias JobSpec
	pis := make([]Initiator, len(job.Initiators))
	for i, modelInitr := range job.Initiators {
		pis[i] = Initiator{modelInitr}
	}
	return json.Marshal(&struct {
		Initiators []Initiator `json:"initiators"`
		Alias
	}{
		pis,
		Alias(job),
	})
}

// FriendlyCreatedAt returns a human-readable string of the Job's
// CreatedAt field.
func (job JobSpec) FriendlyCreatedAt() string {
	return utils.ISO8601UTC(job.CreatedAt)
}

// FriendlyStartAt returns a human-readable string of the Job's
// StartAt field.
func (job JobSpec) FriendlyStartAt() string {
	if job.StartAt.Valid {
		return utils.ISO8601UTC(job.StartAt.Time)
	}
	return ""
}

// FriendlyEndAt returns a human-readable string of the Job's
// EndAt field.
func (job JobSpec) FriendlyEndAt() string {
	if job.EndAt.Valid {
		return utils.ISO8601UTC(job.EndAt.Time)
	}
	return ""
}

// FriendlyMinPayment returns a formatted string of the Job's
// Minimum Link Payment threshold
func (job JobSpec) FriendlyMinPayment() string {
	return job.MinPayment.Text(10)
}

// FriendlyInitiators returns the list of Initiator types as
// a comma separated string.
func (job JobSpec) FriendlyInitiators() string {
	var initrs []string
	for _, i := range job.Initiators {
		initrs = append(initrs, i.Type)
	}
	return strings.Join(initrs, "\n")
}

// FriendlyTasks returns the list of Task types as a comma
// separated string.
func (job JobSpec) FriendlyTasks() string {
	var tasks []string
	for _, t := range job.Tasks {
		tasks = append(tasks, t.Type.String())
	}

	return strings.Join(tasks, "\n")
}

// Initiator holds the Job definition's Initiator.
type Initiator struct {
	models.Initiator
}

// MarshalJSON returns the JSON data of the Initiator based
// on its Initiator Type.
func (i Initiator) MarshalJSON() ([]byte, error) {
	p, err := initiatorParams(i)
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(&struct {
		ID     int64       `json:"id"`
		JobID  uuid.UUID   `json:"jobSpecId"`
		Type   string      `json:"type"`
		Params interface{} `json:"params"`
	}{i.ID, i.JobSpecID.UUID(), i.Type, p})
}

func initiatorParams(i Initiator) (interface{}, error) {
	switch i.Type {
	case models.InitiatorWeb:
		return struct{}{}, nil
	case models.InitiatorCron:
		return struct {
			Schedule models.Cron `json:"schedule"`
		}{i.Schedule}, nil
	case models.InitiatorRunAt:
		return struct {
			Time models.AnyTime `json:"time"`
			Ran  bool           `json:"ran"`
		}{models.NewAnyTime(i.Time.Time), i.Ran}, nil
	case models.InitiatorEthLog:
		fallthrough
	case models.InitiatorRunLog:
		return struct {
			Address common.Address `json:"address"`
		}{i.Address}, nil
	case models.InitiatorExternal:
		return struct {
			Name string `json:"name"`
		}{i.Name}, nil
	case models.InitiatorFluxMonitor:
		return struct {
			Address           common.Address         `json:"address"`
			RequestData       models.JSON            `json:"requestData"`
			Feeds             models.JSON            `json:"feeds"`
			Threshold         float32                `json:"threshold"`
			AbsoluteThreshold float32                `json:"absoluteThreshold"`
			Precision         int32                  `json:"precision"`
			PollTimer         models.PollTimerConfig `json:"pollTimer,omitempty"`
			IdleTimer         models.IdleTimerConfig `json:"idleTimer,omitempty"`
		}{i.Address, i.RequestData, i.Feeds, i.Threshold, i.AbsoluteThreshold,
			i.Precision, i.PollTimer, i.IdleTimer}, nil
	case models.InitiatorRandomnessLog:
		return struct {
			Address          common.Address `json:"address"`
			JobIDTopicFilter models.JobID   `json:"jobIDTopicFilter"`
		}{
			i.Address,
			i.JobIDTopicFilter,
		}, nil
	default:
		return nil, fmt.Errorf("cannot marshal unsupported initiator type '%v'", i.Type)
	}
}

// FriendlyRunAt returns a human-readable string for Cron Initiator types.
func (i Initiator) FriendlyRunAt() string {
	if i.Type == models.InitiatorRunAt {
		return utils.ISO8601UTC(i.Time.Time)
	}
	return ""
}

// FriendlyAddress returns the Ethereum address if present, and a blank
// string if not.
func (i Initiator) FriendlyAddress() string {
	if i.IsLogInitiated() {
		return utils.LogListeningAddress(i.Address)
	}
	return ""
}

// JobRun presents an API friendly version of the data.
type JobRun struct {
	models.JobRun
}

// MarshalJSON returns the JSON data of the JobRun and its Initiator.
func (jr JobRun) MarshalJSON() ([]byte, error) {
	type Alias JobRun
	return json.Marshal(&struct {
		Alias
		Initiator Initiator `json:"initiator"`
	}{
		Alias(jr),
		Initiator{jr.Initiator},
	})
}

// TaskSpec holds a task specified in the Job definition.
type TaskSpec struct {
	models.TaskSpec
}

// FriendlyParams returns a map of the TaskSpec's parameters.
func (t TaskSpec) FriendlyParams() (string, string) {
	keys := []string{}
	values := []string{}
	t.Params.ForEach(func(key, value gjson.Result) bool {
		if key.String() != "type" {
			keys = append(keys, key.String())
			values = append(values, value.String())
		}
		return true
	})
	return strings.Join(keys, "\n"), strings.Join(values, "\n")
}

// FriendlyBigInt returns a string printing the integer in both
// decimal and hexadecimal formats.
func FriendlyBigInt(n *big.Int) string {
	return fmt.Sprintf("#%[1]v (0x%[1]x)", n)
}

// ServiceAgreement presents an API friendly version of the data.
type ServiceAgreement struct {
	models.ServiceAgreement
}

type ServiceAgreementPresentation struct {
	ID            string             `json:"id"`
	CreatedAt     string             `json:"createdAt"`
	Encumbrance   models.Encumbrance `json:"encumbrance"`
	EncumbranceID int64              `json:"encumbranceID"`
	RequestBody   string             `json:"requestBody"`
	Signature     string             `json:"signature"`
	JobSpec       models.JobSpec     `json:"jobSpec"`
	JobSpecID     string             `json:"jobSpecId"`
}

// MarshalJSON presents the ServiceAgreement as public JSON data
func (sa ServiceAgreement) MarshalJSON() ([]byte, error) {
	return json.Marshal(ServiceAgreementPresentation{
		ID:            sa.ID,
		CreatedAt:     utils.ISO8601UTC(sa.CreatedAt),
		Encumbrance:   sa.Encumbrance,
		EncumbranceID: sa.EncumbranceID,
		RequestBody:   sa.RequestBody,
		Signature:     sa.Signature.String(),
		JobSpec:       sa.JobSpec,
		JobSpecID:     sa.JobSpecID.String(),
	})
}

// FriendlyCreatedAt returns the ServiceAgreement's created at time in a human
// readable format.
func (sa ServiceAgreement) FriendlyCreatedAt() string {
	return utils.ISO8601UTC(sa.CreatedAt)
}

// FriendlyExpiration returns the ServiceAgreement's Encumbrance expiration time
// in a human readable format.
func (sa ServiceAgreement) FriendlyExpiration() string {
	return fmt.Sprintf("%v seconds", sa.Encumbrance.Expiration)
}

// FriendlyPayment returns the ServiceAgreement's Encumbrance payment amount in
// a human readable format.
func (sa ServiceAgreement) FriendlyPayment() string {
	return fmt.Sprintf("%v LINK", sa.Encumbrance.Payment.Link())
}

// FriendlyAggregator returns the ServiceAgreement's aggregator address,
// in a human readable format.
func (sa ServiceAgreement) FriendlyAggregator() string {
	return sa.Encumbrance.Aggregator.String()
}

// FriendlyAggregator returns the ServiceAgreement's aggregator initialization
// method's function selector, in a human readable format.
func (sa ServiceAgreement) FriendlyAggregatorInitMethod() string {
	return sa.Encumbrance.AggInitiateJobSelector.String()
}

// FriendlyAggregatorFulfillMethod returns the ServiceAgreement's aggregator
// fulfillment (orcale reporting) method's function selector, in a human
// readable format.
func (sa ServiceAgreement) FriendlyAggregatorFulfillMethod() string {
	return sa.Encumbrance.AggFulfillSelector.String()
}
