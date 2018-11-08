// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
)

// LogListeningAddress returns the LogListeningAddress
func LogListeningAddress(address common.Address) string {
	if address == utils.ZeroAddress {
		return "[all]"
	}
	return address.String()
}

// ShowEthBalance returns the current Eth Balance for current Account
func ShowEthBalance(store *store.Store) (map[string]interface{}, error) {
	keysAndValues := make(map[string]interface{})
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}
	account, err := store.KeyStore.GetAccount()
	if err != nil {
		return keysAndValues, err
	}
	address := account.Address
	balance, err := store.TxManager.GetEthBalance(address)
	if err != nil {
		return keysAndValues, err
	}
	keysAndValues["message"] = fmt.Sprintf("ETH Balance for %v: %v", address.Hex(), balance)
	keysAndValues["address"] = account.Address
	keysAndValues["balance"] = balance
	if balance.IsZero() {
		return keysAndValues, errors.New("0 Balance. Chainlink node not fully functional, please deposit ETH into your address: " + address.Hex())
	}
	return keysAndValues, nil
}

// ShowLinkBalance returns the current Link Balance for current Account
func ShowLinkBalance(store *store.Store) (map[string]interface{}, error) {
	keysAndValues := make(map[string]interface{})
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}
	account, err := store.KeyStore.GetAccount()
	if err != nil {
		return keysAndValues, err
	}

	address := account.Address
	linkBalance, err := store.TxManager.GetLinkBalance(address)
	if err != nil {
		return keysAndValues, err
	}
	keysAndValues["address"] = account.Address
	keysAndValues["balance"] = linkBalance.String()
	keysAndValues["message"] = fmt.Sprintf("Link Balance for %v: %v", address.Hex(), linkBalance.String())
	return keysAndValues, nil
}

// BridgeType holds a bridge.
type BridgeType struct {
	models.BridgeType
}

// MarshalJSON returns the JSON data of the Bridge.
func (bt BridgeType) MarshalJSON() ([]byte, error) {
	type Alias BridgeType
	return json.Marshal(&struct {
		Alias
	}{
		Alias(bt),
	})
}

// AccountBalance holds the hex representation of the address plus it's ETH & LINK balances
type AccountBalance struct {
	Address     string       `json:"address"`
	EthBalance  *assets.Eth  `json:"ethBalance"`
	LinkBalance *assets.Link `json:"linkBalance"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (a AccountBalance) GetID() string {
	return a.Address
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (a *AccountBalance) SetID(value string) error {
	a.Address = value
	return nil
}

// ConfigWhitelist are the non-secret values of the node
type ConfigWhitelist struct {
	AllowOrigins             string          `json:"allowOrigins"`
	BridgeResponseURL        string          `json:"bridgeResponseURL,omitempty"`
	ChainID                  uint64          `json:"ethChainId"`
	ChainlinkDev             bool            `json:"chainlinkDev"`
	ClientNodeURL            string          `json:"clientNodeUrl"`
	DatabaseTimeout          store.Duration  `json:"databaseTimeout"`
	EthereumURL              string          `json:"ethUrl"`
	EthGasBumpThreshold      uint64          `json:"ethGasBumpThreshold"`
	EthGasBumpWei            *big.Int        `json:"ethGasBumpWei"`
	EthGasPriceDefault       *big.Int        `json:"ethGasPriceDefault"`
	LinkContractAddress      string          `json:"linkContractAddress"`
	LogLevel                 store.LogLevel  `json:"logLevel"`
	MinimumContractPayment   *assets.Link    `json:"minimumContractPayment"`
	MinimumRequestExpiration uint64          `json:"minimumRequestExpiration"`
	MinIncomingConfirmations uint64          `json:"minIncomingConfirmations"`
	MinOutgoingConfirmations uint64          `json:"minOutgoingConfirmations"`
	OracleContractAddress    *common.Address `json:"oracleContractAddress"`
	Port                     uint16          `json:"chainlinkPort"`
	ReaperExpiration         store.Duration  `json:"reaperExpiration"`
	RootDir                  string          `json:"root"`
	SessionTimeout           store.Duration  `json:"sessionTimeout"`
	TLSHost                  string          `json:"chainlinkTLSHost"`
	TLSPort                  uint16          `json:"chainlinkTLSPort"`
}

// NewConfigWhitelist creates an instance of ConfigWhitelist
func NewConfigWhitelist(config store.Config) ConfigWhitelist {
	return ConfigWhitelist{
		AllowOrigins:             config.AllowOrigins,
		BridgeResponseURL:        config.BridgeResponseURL.String(),
		ChainID:                  config.ChainID,
		ChainlinkDev:             config.Dev,
		ClientNodeURL:            config.ClientNodeURL,
		DatabaseTimeout:          config.DatabaseTimeout,
		EthereumURL:              config.EthereumURL,
		EthGasBumpThreshold:      config.EthGasBumpThreshold,
		EthGasBumpWei:            &config.EthGasBumpWei,
		EthGasPriceDefault:       &config.EthGasPriceDefault,
		LinkContractAddress:      config.LinkContractAddress,
		LogLevel:                 config.LogLevel,
		MinimumContractPayment:   &config.MinimumContractPayment,
		MinimumRequestExpiration: config.MinimumRequestExpiration,
		MinIncomingConfirmations: config.MinIncomingConfirmations,
		MinOutgoingConfirmations: config.MinOutgoingConfirmations,
		OracleContractAddress:    config.OracleContractAddress,
		Port:             config.Port,
		ReaperExpiration: config.ReaperExpiration,
		RootDir:          config.RootDir,
		SessionTimeout:   config.SessionTimeout,
		TLSHost:          config.TLSHost,
		TLSPort:          config.TLSPort,
	}
}

// String returns the values as a newline delimited string
func (c ConfigWhitelist) String() string {
	fmtConfig := "LOG_LEVEL: %v\n" +
		"ROOT: %s\n" +
		"CHAINLINK_PORT: %d\n" +
		"CHAINLINK_TLS_PORT: %d\n" +
		"CHAINLINK_TLS_HOST: %s\n" +
		"ETH_URL: %s\n" +
		"ETH_CHAIN_ID: %d\n" +
		"CLIENT_NODE_URL: %s\n" +
		"TX_MIN_CONFIRMATIONS: %d\n" +
		"TASK_MIN_CONFIRMATIONS: %d\n" +
		"ETH_GAS_BUMP_THRESHOLD: %d\n" +
		"ETH_GAS_BUMP_WEI: %s\n" +
		"ETH_GAS_PRICE_DEFAULT: %s\n" +
		"LINK_CONTRACT_ADDRESS: %s\n" +
		"MINIMUM_CONTRACT_PAYMENT: %s\n" +
		"ORACLE_CONTRACT_ADDRESS: %s\n" +
		"DATABASE_POLL_INTERVAL: %s\n" +
		"ALLOW_ORIGINS: %s\n" +
		"CHAINLINK_DEV: %v\n" +
		"SESSION_TIMEOUT: %v\n" +
		"REAPER_EXPIRATION: %v\n" +
		"BRIDGE_RESPONSE_URL: %s\n"

	oracleContractAddress := ""
	if c.OracleContractAddress != nil {
		oracleContractAddress = c.OracleContractAddress.String()
	}

	return fmt.Sprintf(
		fmtConfig,
		c.LogLevel,
		c.RootDir,
		c.Port,
		c.TLSPort,
		c.TLSHost,
		c.EthereumURL,
		c.ChainID,
		c.ClientNodeURL,
		c.MinOutgoingConfirmations,
		c.MinIncomingConfirmations,
		c.EthGasBumpThreshold,
		c.EthGasBumpWei.String(),
		c.EthGasPriceDefault.String(),
		c.LinkContractAddress,
		c.MinimumContractPayment.String(),
		oracleContractAddress,
		c.DatabaseTimeout,
		c.AllowOrigins,
		c.ChainlinkDev,
		c.SessionTimeout,
		c.ReaperExpiration,
		c.BridgeResponseURL,
	)
}

// GetID generates a new ID for jsonapi serialization.
func (c ConfigWhitelist) GetID() string {
	return utils.NewBytes32ID()
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (c *ConfigWhitelist) SetID(value string) error {
	return nil
}

// JobSpec holds the JobSpec definition and each run associated with that Job.
type JobSpec struct {
	models.JobSpec
	Runs []JobRun `json:"runs,omitempty"`
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
	return job.CreatedAt.HumanString()
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
			Time models.Time `json:"time"`
			Ran  bool        `json:"ran"`
		}{i.Time, i.Ran}, nil
	case models.InitiatorEthLog:
		fallthrough
	case models.InitiatorRunLog:
		return struct {
			Address common.Address `json:"address"`
		}{i.Address}, nil
	default:
		return nil, fmt.Errorf("Cannot marshal unsupported initiator type %v", i.Type)
	}
}

// FriendlyRunAt returns a human-readable string for Cron Initiator types.
func (i Initiator) FriendlyRunAt() string {
	if i.Type == models.InitiatorRunAt {
		return i.Time.HumanString()
	}
	return ""
}

var emptyAddress = common.Address{}.String()

// FriendlyAddress returns the Ethereum address if present, and a blank
// string if not.
func (i Initiator) FriendlyAddress() string {
	if i.IsLogInitiated() {
		return LogListeningAddress(i.Address)
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
// decimal and hexidecimal formats.
func FriendlyBigInt(n *big.Int) string {
	return fmt.Sprintf("#%[1]v (0x%[1]x)", n)
}

// ServiceAgreement presents an API friendly version of the data.
type ServiceAgreement struct {
	models.ServiceAgreement
}

// MarshalJSON returns the JSON data of the ServiceAgreement.
func (sa ServiceAgreement) MarshalJSON() ([]byte, error) {
	return []byte(sa.ServiceAgreement.RequestBody), nil
}

// FriendlyCreatedAt returns the ServiceAgreement's created at time in a human
// readable format.
func (sa ServiceAgreement) FriendlyCreatedAt() string {
	return sa.CreatedAt.HumanString()
}

// FriendlyExpiration returns the ServiceAgreement's Encumbrance expiration time
// in a human readable format.
func (sa ServiceAgreement) FriendlyExpiration() string {
	return fmt.Sprintf("%v seconds", sa.Encumbrance.Expiration)
}

// FriendlyPayment returns the ServiceAgreement's Encumbrance payment amount in
// a human readable format.
func (sa ServiceAgreement) FriendlyPayment() string {
	return fmt.Sprintf("%v LINK", (*assets.Link)(sa.Encumbrance.Payment).String())
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
		CreatedAt: u.User.CreatedAt.ISO8601(),
	})
}
