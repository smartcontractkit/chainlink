// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
)

type requestType int

const (
	ethRequest requestType = iota
	linkRequest
)

// ShowEthBalance returns the current Eth Balance for current Account
func ShowEthBalance(store *store.Store) ([]map[string]interface{}, error) {
	return showBalanceFor(store, ethRequest)
}

// ShowLinkBalance returns the current Link Balance for current Account
func ShowLinkBalance(store *store.Store) ([]map[string]interface{}, error) {
	return showBalanceFor(store, linkRequest)
}

func showBalanceFor(store *store.Store, balanceType requestType) ([]map[string]interface{}, error) {
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}

	var merr error
	info := []map[string]interface{}{}
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
func showBalanceForAccount(store *store.Store, account accounts.Account, balanceType requestType) (map[string]interface{}, error) {
	balance, err := getBalance(store, account, balanceType)
	if err != nil {
		return nil, err
	}
	address := account.Address
	keysAndValues := make(map[string]interface{})
	keysAndValues["message"] = fmt.Sprintf("%v Balance for %v: %v", balance.Symbol(), address.Hex(), balance.String())
	keysAndValues["balance"] = balance.String()
	keysAndValues["address"] = address
	if balance.IsZero() && balanceType == ethRequest {
		return nil, errors.New("0 ETH Balance. Chainlink node not fully functional, please deposit ETH into your address: " + address.Hex())
	}
	return keysAndValues, nil
}

func getBalance(store *store.Store, account accounts.Account, balanceType requestType) (balanceable, error) {
	switch balanceType {
	case ethRequest:
		return store.TxManager.GetEthBalance(account.Address)
	case linkRequest:
		return store.TxManager.GetLINKBalance(account.Address)
	}
	return nil, fmt.Errorf("Impossible to get balance for %T with value %v", balanceType, balanceType)
}

type balanceable interface {
	IsZero() bool
	String() string
	Symbol() string
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
//
// If you add an entry here, you should update NewConfigWhitelist and
// ConfigWhitelist#String accordingly.
type ConfigWhitelist struct {
	AccountAddress string `json:"accountAddress"`
	whitelist
}

type whitelist struct {
	AllowOrigins             string          `json:"allowOrigins"`
	BridgeResponseURL        string          `json:"bridgeResponseURL,omitempty"`
	ChainID                  uint64          `json:"ethChainId"`
	ClientNodeURL            string          `json:"clientNodeUrl"`
	Dev                      bool            `json:"chainlinkDev"`
	DatabaseTimeout          time.Duration   `json:"databaseTimeout"`
	EthereumURL              string          `json:"ethUrl"`
	EthGasBumpThreshold      uint64          `json:"ethGasBumpThreshold"`
	EthGasBumpWei            *big.Int        `json:"ethGasBumpWei"`
	EthGasPriceDefault       *big.Int        `json:"ethGasPriceDefault"`
	JSONConsole              bool            `json:"jsonConsole"`
	LinkContractAddress      string          `json:"linkContractAddress"`
	ExplorerURL              string          `json:"explorerUrl"`
	LogLevel                 store.LogLevel  `json:"logLevel"`
	LogToDisk                bool            `json:"logToDisk"`
	MinimumContractPayment   *assets.Link    `json:"minimumContractPayment"`
	MinimumRequestExpiration uint64          `json:"minimumRequestExpiration"`
	MinIncomingConfirmations uint32          `json:"minIncomingConfirmations"`
	MinOutgoingConfirmations uint64          `json:"minOutgoingConfirmations"`
	OracleContractAddress    *common.Address `json:"oracleContractAddress"`
	Port                     uint16          `json:"chainlinkPort"`
	ReaperExpiration         time.Duration   `json:"reaperExpiration"`
	RootDir                  string          `json:"root"`
	SessionTimeout           time.Duration   `json:"sessionTimeout"`
	TLSHost                  string          `json:"chainlinkTLSHost"`
	TLSPort                  uint16          `json:"chainlinkTLSPort"`
}

// NewConfigWhitelist creates an instance of ConfigWhitelist
func NewConfigWhitelist(store *store.Store) (ConfigWhitelist, error) {
	config := store.Config
	account, err := store.KeyStore.GetFirstAccount()
	if err != nil {
		return ConfigWhitelist{}, err
	}

	explorerURL := ""
	if config.ExplorerURL() != nil {
		explorerURL = config.ExplorerURL().String()
	}
	return ConfigWhitelist{
		AccountAddress: account.Address.Hex(),
		whitelist: whitelist{
			AllowOrigins:             config.AllowOrigins(),
			BridgeResponseURL:        config.BridgeResponseURL().String(),
			ChainID:                  config.ChainID(),
			ClientNodeURL:            config.ClientNodeURL(),
			Dev:                      config.Dev(),
			DatabaseTimeout:          config.DatabaseTimeout(),
			EthereumURL:              config.EthereumURL(),
			EthGasBumpThreshold:      config.EthGasBumpThreshold(),
			EthGasBumpWei:            config.EthGasBumpWei(),
			EthGasPriceDefault:       config.EthGasPriceDefault(),
			JSONConsole:              config.JSONConsole(),
			LinkContractAddress:      config.LinkContractAddress(),
			ExplorerURL:              explorerURL,
			LogLevel:                 config.LogLevel(),
			LogToDisk:                config.LogToDisk(),
			MinimumContractPayment:   config.MinimumContractPayment(),
			MinimumRequestExpiration: config.MinimumRequestExpiration(),
			MinIncomingConfirmations: config.MinIncomingConfirmations(),
			MinOutgoingConfirmations: config.MinOutgoingConfirmations(),
			OracleContractAddress:    config.OracleContractAddress(),
			Port:                     config.Port(),
			ReaperExpiration:         config.ReaperExpiration(),
			RootDir:                  config.RootDir(),
			SessionTimeout:           config.SessionTimeout(),
			TLSHost:                  config.TLSHost(),
			TLSPort:                  config.TLSPort(),
		},
	}, nil
}

// String returns the values as a newline delimited string
func (c ConfigWhitelist) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("ACCOUNT_ADDRESS: %v\n", c.AccountAddress))

	schemaT := reflect.TypeOf(store.ConfigSchema{})
	cwlT := reflect.TypeOf(c.whitelist)
	cwlV := reflect.ValueOf(c.whitelist)
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
func (c ConfigWhitelist) GetID() string {
	return utils.NewBytes32ID()
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (c *ConfigWhitelist) SetID(value string) error {
	return nil
}

// JobSpec holds the JobSpec definition
type JobSpec struct {
	models.JobSpec
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
	default:
		return nil, fmt.Errorf("Cannot marshal unsupported initiator type %v", i.Type)
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

// Tx is a jsonapi wrapper for an Ethereum Transaction.
type Tx struct {
	Confirmed bool            `json:"confirmed,omitempty"`
	Data      hexutil.Bytes   `json:"data,omitempty"`
	From      *common.Address `json:"from,omitempty"`
	GasLimit  string          `json:"gasLimit,omitempty"`
	GasPrice  string          `json:"gasPrice,omitempty"`
	Hash      common.Hash     `json:"hash,omitempty"`
	Hex       string          `json:"rawHex,omitempty"`
	Nonce     string          `json:"nonce,omitempty"`
	SentAt    string          `json:"sentAt,omitempty"`
	To        *common.Address `json:"to,omitempty"`
	Value     string          `json:"value,omitempty"`
}

// NewTx builds a transaction presenter.
func NewTx(tx *models.Tx) Tx {
	return Tx{
		Confirmed: tx.Confirmed,
		Data:      hexutil.Bytes(tx.Data),
		From:      &tx.From,
		GasLimit:  strconv.FormatUint(tx.GasLimit, 10),
		GasPrice:  tx.GasPrice.String(),
		Hash:      tx.Hash,
		Hex:       tx.SignedRawTx,
		Nonce:     strconv.FormatUint(tx.Nonce, 10),
		SentAt:    strconv.FormatUint(tx.SentAt, 10),
		To:        &tx.To,
		Value:     tx.Value.String(),
	}
}

// NewTxFromAttempt builds a transaction presenter from a TxAttempt
//
// models.Tx represents a transaction in progress, with a series of
// models.TxAttempts, each one of these represents an ethereum transaction. A
// TxAttempt only stores the unique details of an ethereum transaction, with
// the rest of the details on its related Tx.
//
// So for presenting a TxAttempt, we take its Hash, GasPrice etc. and get the
// rest of the details from its Tx.
//
// NOTE: We take a copy here as we don't want side effects.
//
func NewTxFromAttempt(txAttempt models.TxAttempt) Tx {
	tx := txAttempt.Tx
	tx.Hash = txAttempt.Hash
	tx.GasPrice = txAttempt.GasPrice
	tx.Confirmed = txAttempt.Confirmed
	tx.SentAt = txAttempt.SentAt
	tx.SignedRawTx = txAttempt.SignedRawTx
	return NewTx(tx)
}

// GetID returns the jsonapi ID.
func (t Tx) GetID() string {
	return t.Hash.String()
}

// GetName returns the collection name for jsonapi.
func (Tx) GetName() string {
	return "transactions"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (t *Tx) SetID(hex string) error {
	t.Hash = common.HexToHash(hex)
	return nil
}
