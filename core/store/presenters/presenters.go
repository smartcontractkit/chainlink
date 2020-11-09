// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v3"
)

type requestType int

const (
	ethRequest requestType = iota
	linkRequest
)

// ShowEthBalance returns the current Eth Balance for current Account
func ShowEthBalance(store *store.Store) ([]map[string]string, error) {
	return showBalanceFor(store, ethRequest)
}

// ShowLinkBalance returns the current Link Balance for current Account
func ShowLinkBalance(store *store.Store) ([]map[string]string, error) {
	return showBalanceFor(store, linkRequest)
}

func showBalanceFor(store *store.Store, balanceType requestType) ([]map[string]string, error) {
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}

	var merr error
	info := []map[string]string{}
	for _, account := range store.KeyStore.Accounts() {
		b, err := showBalanceForAccount(store, account, balanceType)
		merr = multierr.Append(merr, err)
		if err == nil {
			info = append(info, b)
		}
	}
	return info, merr
}

// ShowEthBalance returns the current Eth Balance for current Account
func showBalanceForAccount(store *store.Store, account accounts.Account, balanceType requestType) (map[string]string, error) {
	balance, err := getBalance(store, account, balanceType)
	if err != nil {
		return nil, err
	}
	address := account.Address
	keysAndValues := make(map[string]string)
	keysAndValues["message"] = fmt.Sprintf("%v Balance for %v: %v", balance.Symbol(), address.Hex(), balance.String())
	keysAndValues["balance"] = balance.String()
	keysAndValues["address"] = address.String()
	if balance.IsZero() && balanceType == ethRequest {
		return nil, errors.New("0 ETH Balance. Chainlink node not fully functional, please deposit ETH into your address: " + address.Hex())
	}
	return keysAndValues, nil
}

func getBalance(store *store.Store, account accounts.Account, balanceType requestType) (balanceable, error) {
	switch balanceType {
	case ethRequest:
		bal, err := store.EthClient.BalanceAt(context.TODO(), account.Address, nil)
		if err != nil {
			return nil, err
		}
		return (*assets.Eth)(bal), nil
	case linkRequest:
		return bulletprooftxmanager.GetLINKBalance(store.Config, store.EthClient, account.Address)
	}
	return nil, fmt.Errorf("impossible to get balance for %T with value %v", balanceType, balanceType)
}

type balanceable interface {
	IsZero() bool
	String() string
	Symbol() string
}

// ETHKey holds the hex representation of the address plus it's ETH & LINK balances
type ETHKey struct {
	Address     string       `json:"address"`
	EthBalance  *assets.Eth  `json:"ethBalance"`
	LinkBalance *assets.Link `json:"linkBalance"`
	NextNonce   *int64       `json:"nextNonce"`
	LastUsed    *time.Time   `json:"lastUsed"`
	IsFunding   bool         `json:"isFunding"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	DeletedAt   null.Time    `json:"deletedAt"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (k ETHKey) GetID() string {
	return k.Address
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (k *ETHKey) SetID(value string) error {
	k.Address = value
	return nil
}

// ConfigPrinter are the non-secret values of the node
//
// If you add an entry here, you should update NewConfigPrinter and
// ConfigPrinter#String accordingly.
type ConfigPrinter struct {
	AccountAddress string `json:"accountAddress"`
	EnvPrinter
}

// EnvPrinter contains the supported environment variables
type EnvPrinter struct {
	AllowOrigins                          string          `json:"allowOrigins"`
	BalanceMonitorEnabled                 bool            `json:"balanceMonitorEnabled"`
	BlockBackfillDepth                    uint64          `json:"blockBackfillDepth"`
	BridgeResponseURL                     string          `json:"bridgeResponseURL,omitempty"`
	ChainID                               *big.Int        `json:"ethChainId"`
	ClientNodeURL                         string          `json:"clientNodeUrl"`
	DatabaseTimeout                       models.Duration `json:"databaseTimeout"`
	DatabaseMaximumTxDuration             time.Duration   `json:"databaseMaximumTxDuration"`
	DefaultHTTPLimit                      int64           `json:"defaultHttpLimit"`
	DefaultHTTPTimeout                    models.Duration `json:"defaultHttpTimeout"`
	Dev                                   bool            `json:"chainlinkDev"`
	EnableExperimentalAdapters            bool            `json:"enableExperimentalAdapters"`
	EthBalanceMonitorBlockDelay           uint16          `json:"ethBalanceMonitorBlockDelay"`
	EthereumDisabled                      bool            `json:"ethereumDisabled"`
	EthFinalityDepth                      uint            `json:"ethFinalityDepth"`
	EthGasBumpThreshold                   uint64          `json:"ethGasBumpThreshold"`
	EthGasBumpTxDepth                     uint16          `json:"ethGasBumpTxDepth"`
	EthGasBumpWei                         *big.Int        `json:"ethGasBumpWei"`
	EthGasLimitDefault                    uint64          `json:"ethGasLimitDefault"`
	EthGasPriceDefault                    *big.Int        `json:"ethGasPriceDefault"`
	EthHeadTrackerHistoryDepth            uint            `json:"ethHeadTrackerHistoryDepth"`
	EthHeadTrackerMaxBufferSize           uint            `json:"ethHeadTrackerMaxBufferSize"`
	EthMaxGasPriceWei                     *big.Int        `json:"ethMaxGasPriceWei"`
	EthereumURL                           string          `json:"ethUrl"`
	EthereumSecondaryURL                  string          `json:"ethSecondaryURL"`
	ExplorerURL                           string          `json:"explorerUrl"`
	FeatureExternalInitiators             bool            `json:"featureExternalInitiators"`
	FeatureFluxMonitor                    bool            `json:"featureFluxMonitor"`
	FeatureOffchainReporting              bool            `json:"featureOffchainReporting"`
	FlagsContractAddress                  string          `json:"flagsContractAddress"`
	GasUpdaterBlockDelay                  uint16          `json:"gasUpdaterBlockDelay"`
	GasUpdaterBlockHistorySize            uint16          `json:"gasUpdaterBlockHistorySize"`
	GasUpdaterEnabled                     bool            `json:"gasUpdaterEnabled"`
	GasUpdaterTransactionPercentile       uint16          `json:"gasUpdaterTransactionPercentile"`
	InsecureFastScrypt                    bool            `json:"insecureFastScrypt"`
	JobPipelineDBPollInterval             time.Duration   `json:"jobPipelineDBPollInterval"`
	JobPipelineMaxTaskDuration            time.Duration   `json:"jobPipelineMaxTaskDuration"`
	JobPipelineParallelism                uint8           `json:"jobPipelineParallelism"`
	JobPipelineReaperInterval             time.Duration   `json:"jobPipelineReaperInterval"`
	JobPipelineReaperThreshold            time.Duration   `json:"jobPipelineReaperThreshold"`
	JSONConsole                           bool            `json:"jsonConsole"`
	LinkContractAddress                   string          `json:"linkContractAddress"`
	LogLevel                              orm.LogLevel    `json:"logLevel"`
	LogSQLMigrations                      bool            `json:"logSqlMigrations"`
	LogSQLStatements                      bool            `json:"logSqlStatements"`
	LogToDisk                             bool            `json:"logToDisk"`
	MaximumServiceDuration                models.Duration `json:"maximumServiceDuration"`
	MinIncomingConfirmations              uint32          `json:"minIncomingConfirmations"`
	MinRequiredOutgoingConfirmations      uint64          `json:"minOutgoingConfirmations"`
	MinimumServiceDuration                models.Duration `json:"minimumServiceDuration"`
	MinimumContractPayment                *assets.Link    `json:"minimumContractPayment"`
	MinimumRequestExpiration              uint64          `json:"minimumRequestExpiration"`
	OCRBootstrapCheckInterval             time.Duration   `json:"ocrBootstrapCheckInterval"`
	OCRContractTransmitterTransmitTimeout time.Duration   `json:"ocrContractTransmitterTransmitTimeout"`
	OCRDatabaseTimeout                    time.Duration   `json:"ocrDatabaseTimeout"`
	P2PListenIP                           string          `json:"ocrListenIP"`
	P2PListenPort                         uint16          `json:"ocrListenPort"`
	OCRIncomingMessageBufferSize          int             `json:"ocrIncomingMessageBufferSize"`
	OCROutgoingMessageBufferSize          int             `json:"ocrOutgoingMessageBufferSize"`
	OCRNewStreamTimeout                   time.Duration   `json:"ocrNewStreamTimeout"`
	OCRDHTLookupInterval                  int             `json:"ocrDHTLookupInterval"`
	OCRTraceLogging                       bool            `json:"ocrTraceLogging"`
	OperatorContractAddress               common.Address  `json:"oracleContractAddress"`
	Port                                  uint16          `json:"chainlinkPort"`
	ReaperExpiration                      models.Duration `json:"reaperExpiration"`
	ReplayFromBlock                       int64           `json:"replayFromBlock"`
	RootDir                               string          `json:"root"`
	SecureCookies                         bool            `json:"secureCookies"`
	SessionTimeout                        models.Duration `json:"sessionTimeout"`
	TLSHost                               string          `json:"chainlinkTLSHost"`
	TLSPort                               uint16          `json:"chainlinkTLSPort"`
	TLSRedirect                           bool            `json:"chainlinkTLSRedirect"`
	TxAttemptLimit                        uint16          `json:"txAttemptLimit"`
}

// NewConfigPrinter creates an instance of ConfigPrinter
func NewConfigPrinter(store *store.Store) (ConfigPrinter, error) {
	config := store.Config
	account, err := store.KeyStore.GetFirstAccount()
	if err != nil {
		return ConfigPrinter{}, err
	}

	explorerURL := ""
	if config.ExplorerURL() != nil {
		explorerURL = config.ExplorerURL().String()
	}
	return ConfigPrinter{
		AccountAddress: account.Address.Hex(),
		EnvPrinter: EnvPrinter{
			AllowOrigins:                          config.AllowOrigins(),
			BalanceMonitorEnabled:                 config.BalanceMonitorEnabled(),
			BlockBackfillDepth:                    config.BlockBackfillDepth(),
			BridgeResponseURL:                     config.BridgeResponseURL().String(),
			ChainID:                               config.ChainID(),
			ClientNodeURL:                         config.ClientNodeURL(),
			DatabaseTimeout:                       config.DatabaseTimeout(),
			DefaultHTTPLimit:                      config.DefaultHTTPLimit(),
			DefaultHTTPTimeout:                    config.DefaultHTTPTimeout(),
			DatabaseMaximumTxDuration:             config.DatabaseMaximumTxDuration(),
			Dev:                                   config.Dev(),
			EnableExperimentalAdapters:            config.EnableExperimentalAdapters(),
			EthBalanceMonitorBlockDelay:           config.EthBalanceMonitorBlockDelay(),
			EthereumDisabled:                      config.EthereumDisabled(),
			EthFinalityDepth:                      config.EthFinalityDepth(),
			EthGasBumpThreshold:                   config.EthGasBumpThreshold(),
			EthGasBumpTxDepth:                     config.EthGasBumpTxDepth(),
			EthGasBumpWei:                         config.EthGasBumpWei(),
			EthGasLimitDefault:                    config.EthGasLimitDefault(),
			EthGasPriceDefault:                    config.EthGasPriceDefault(),
			EthHeadTrackerHistoryDepth:            config.EthHeadTrackerHistoryDepth(),
			EthHeadTrackerMaxBufferSize:           config.EthHeadTrackerMaxBufferSize(),
			EthMaxGasPriceWei:                     config.EthMaxGasPriceWei(),
			EthereumURL:                           config.EthereumURL(),
			EthereumSecondaryURL:                  config.EthereumSecondaryURL(),
			ExplorerURL:                           explorerURL,
			FeatureExternalInitiators:             config.FeatureExternalInitiators(),
			FeatureFluxMonitor:                    config.FeatureFluxMonitor(),
			FeatureOffchainReporting:              config.FeatureOffchainReporting(),
			FlagsContractAddress:                  config.FlagsContractAddress(),
			GasUpdaterBlockDelay:                  config.GasUpdaterBlockDelay(),
			GasUpdaterBlockHistorySize:            config.GasUpdaterBlockHistorySize(),
			GasUpdaterEnabled:                     config.GasUpdaterEnabled(),
			GasUpdaterTransactionPercentile:       config.GasUpdaterTransactionPercentile(),
			InsecureFastScrypt:                    config.InsecureFastScrypt(),
			JobPipelineDBPollInterval:             config.JobPipelineDBPollInterval(),
			JobPipelineMaxTaskDuration:            config.JobPipelineMaxTaskDuration(),
			JobPipelineParallelism:                config.JobPipelineParallelism(),
			JobPipelineReaperInterval:             config.JobPipelineReaperInterval(),
			JobPipelineReaperThreshold:            config.JobPipelineReaperThreshold(),
			JSONConsole:                           config.JSONConsole(),
			LinkContractAddress:                   config.LinkContractAddress(),
			LogLevel:                              config.LogLevel(),
			LogSQLMigrations:                      config.LogSQLMigrations(),
			LogSQLStatements:                      config.LogSQLStatements(),
			LogToDisk:                             config.LogToDisk(),
			MaximumServiceDuration:                config.MaximumServiceDuration(),
			MinIncomingConfirmations:              config.MinIncomingConfirmations(),
			MinRequiredOutgoingConfirmations:      config.MinRequiredOutgoingConfirmations(),
			MinimumServiceDuration:                config.MinimumServiceDuration(),
			MinimumContractPayment:                config.MinimumContractPayment(),
			MinimumRequestExpiration:              config.MinimumRequestExpiration(),
			OCRBootstrapCheckInterval:             config.OCRBootstrapCheckInterval(),
			OCRContractTransmitterTransmitTimeout: config.OCRContractTransmitterTransmitTimeout(),
			OCRDatabaseTimeout:                    config.OCRDatabaseTimeout(),
			P2PListenIP:                           config.P2PListenIP().String(),
			P2PListenPort:                         config.P2PListenPort(),
			OCRIncomingMessageBufferSize:          config.OCRIncomingMessageBufferSize(),
			OCROutgoingMessageBufferSize:          config.OCROutgoingMessageBufferSize(),
			OCRNewStreamTimeout:                   config.OCRNewStreamTimeout(),
			OCRDHTLookupInterval:                  config.OCRDHTLookupInterval(),
			OCRTraceLogging:                       config.OCRTraceLogging(),
			OperatorContractAddress:               config.OperatorContractAddress(),
			Port:                                  config.Port(),
			ReaperExpiration:                      config.ReaperExpiration(),
			ReplayFromBlock:                       config.ReplayFromBlock(),
			RootDir:                               config.RootDir(),
			SecureCookies:                         config.SecureCookies(),
			SessionTimeout:                        config.SessionTimeout(),
			TLSHost:                               config.TLSHost(),
			TLSPort:                               config.TLSPort(),
			TLSRedirect:                           config.TLSRedirect(),
			TxAttemptLimit:                        config.TxAttemptLimit(),
		},
	}, nil
}

// String returns the values as a newline delimited string
func (c ConfigPrinter) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("ACCOUNT_ADDRESS: %v\n", c.AccountAddress))

	schemaT := reflect.TypeOf(orm.ConfigSchema{})
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
		Type   string      `json:"type"`
		Params interface{} `json:"params"`
	}{i.Type, p})
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
		return struct{ Address common.Address }{i.Address}, nil
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
	return fmt.Sprintf("%v LINK", sa.Encumbrance.Payment.String())
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

// UserPresenter wraps the user record for shipping as a jsonapi response in
// the API.
type UserPresenter struct {
	*models.User
}

// GetID returns the jsonapi ID.
func (u UserPresenter) GetID() string {
	return u.User.Email
}

// GetName returns the collection name for jsonapi.
func (u UserPresenter) GetName() string {
	return "users"
}

// MarshalJSON returns the User as json.
func (u UserPresenter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Email     string `json:"email"`
		CreatedAt string `json:"createdAt"`
	}{
		Email:     u.User.Email,
		CreatedAt: utils.ISO8601UTC(u.User.CreatedAt),
	})
}

// NewAccount is a jsonapi wrapper for an Ethereum account.
type NewAccount struct {
	*accounts.Account
}

// GetID returns the jsonapi ID.
func (a NewAccount) GetID() string {
	return a.Address.String()
}

// GetName returns the collection name for jsonapi.
func (a NewAccount) GetName() string {
	return "keys"
}

// EthTx is a jsonapi wrapper for an Ethereum Transaction.
type EthTx struct {
	ID       int64           `json:"-"`
	State    string          `json:"state,omitempty"`
	Data     hexutil.Bytes   `json:"data,omitempty"`
	From     *common.Address `json:"from,omitempty"`
	GasLimit string          `json:"gasLimit,omitempty"`
	GasPrice string          `json:"gasPrice,omitempty"`
	Hash     common.Hash     `json:"hash,omitempty"`
	Hex      string          `json:"rawHex,omitempty"`
	Nonce    string          `json:"nonce,omitempty"`
	SentAt   string          `json:"sentAt,omitempty"`
	To       *common.Address `json:"to,omitempty"`
	Value    string          `json:"value,omitempty"`
}

func NewEthTxFromAttempt(txa models.EthTxAttempt) EthTx {
	return newEthTxWithAttempt(txa.EthTx, txa)
}

func newEthTxWithAttempt(tx models.EthTx, txa models.EthTxAttempt) EthTx {
	ethTX := EthTx{
		Data:     hexutil.Bytes(tx.EncodedPayload),
		From:     &tx.FromAddress,
		GasLimit: strconv.FormatUint(tx.GasLimit, 10),
		GasPrice: txa.GasPrice.String(),
		Hash:     txa.Hash,
		Hex:      hexutil.Encode(txa.SignedRawTx),
		ID:       tx.ID,
		State:    string(tx.State),
		To:       &tx.ToAddress,
		Value:    tx.Value.String(),
	}
	if tx.Nonce != nil {
		ethTX.Nonce = strconv.FormatUint(uint64(*tx.Nonce), 10)
	}
	if txa.BroadcastBeforeBlockNum != nil {
		ethTX.SentAt = strconv.FormatUint(uint64(*txa.BroadcastBeforeBlockNum), 10)
	}
	return ethTX
}

// GetID returns the jsonapi ID.
func (t EthTx) GetID() string {
	return t.Hash.Hex()
}

// GetName returns the collection name for jsonapi.
func (EthTx) GetName() string {
	return "transactions"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (t *EthTx) SetID(hex string) error {
	t.Hash = common.HexToHash(hex)
	return nil
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

// ExplorerStatus represents the connected server and status of the connection
type ExplorerStatus struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

// NewExplorerStatus returns an initialized ExplorerStatus from the store
func NewExplorerStatus(statsPusher synchronization.StatsPusher) ExplorerStatus {
	url := statsPusher.GetURL()

	return ExplorerStatus{
		Status: string(statsPusher.GetStatus()),
		Url:    url.String(),
	}
}
