package bulletprooftxmanager

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/datatypes"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	cnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type EthTxMeta struct {
	JobID         int32
	RequestID     common.Hash
	RequestTxHash common.Hash
}

func (EthTxMeta) GormDataType() string {
	return "json"
}

type EthTxState string
type EthTxAttemptState string

const (
	EthTxUnstarted               = EthTxState("unstarted")
	EthTxInProgress              = EthTxState("in_progress")
	EthTxFatalError              = EthTxState("fatal_error")
	EthTxUnconfirmed             = EthTxState("unconfirmed")
	EthTxConfirmed               = EthTxState("confirmed")
	EthTxConfirmedMissingReceipt = EthTxState("confirmed_missing_receipt")

	EthTxAttemptInProgress      = EthTxAttemptState("in_progress")
	EthTxAttemptInsufficientEth = EthTxAttemptState("insufficient_eth")
	EthTxAttemptBroadcast       = EthTxAttemptState("broadcast")
)

type NullableEIP2930AccessList struct {
	AccessList types.AccessList
	Valid      bool
}

func NullableEIP2930AccessListFrom(al types.AccessList) (n NullableEIP2930AccessList) {
	if al == nil {
		return
	}
	n.AccessList = al
	n.Valid = true
	return
}

func (e NullableEIP2930AccessList) MarshalJSON() ([]byte, error) {
	if !e.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(e.AccessList)
}

func (e *NullableEIP2930AccessList) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		e.Valid = false
		return nil
	}
	if err := json.Unmarshal(input, &e.AccessList); err != nil {
		return errors.Wrap(err, "NullableEIP2930AccessList: couldn't unmarshal JSON")
	}
	e.Valid = true
	return nil
}

// Value returns this instance serialized for database storage
func (e NullableEIP2930AccessList) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, nil
	}
	return json.Marshal(e)
}

// Scan returns the selector from its serialization in the database
func (e *NullableEIP2930AccessList) Scan(value interface{}) error {
	if value == nil {
		e.Valid = false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, e)
	default:
		return errors.Errorf("unable to convert %v of %T to Big", value, value)
	}
}

// GormDataType override forces gorm to scan/dump as JSON (otherwise it
// incorrectly auto-infers the postgres type from the first element in the
// struct)
func (NullableEIP2930AccessList) GormDataType() string {
	return "jsonb"
}

type EthTx struct {
	ID             int64
	Nonce          *int64
	FromAddress    common.Address
	ToAddress      common.Address
	EncodedPayload []byte
	Value          assets.Eth
	// GasLimit on the EthTx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	GasLimit uint64
	Error    null.String
	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt   *time.Time
	CreatedAt     time.Time
	State         EthTxState
	EthTxAttempts []EthTxAttempt `gorm:"->" json:"-"`
	// Marshalled EthTxMeta
	// Used for additional context around transactions which you want to log
	// at send time.
	Meta       *datatypes.JSON
	Subject    uuid.NullUUID
	EVMChainID utils.Big `gorm:"column:evm_chain_id"`

	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  cnull.Uint32

	// AccessList is optional and only has an effect on DynamicFee transactions
	// on chains that support it (e.g. Ethereum Mainnet after London hard fork)
	AccessList NullableEIP2930AccessList

	// Simulate if set to true will cause this eth_tx to be simulated before
	// initial send and aborted on revert
	Simulate bool
}

func (e EthTx) GetError() error {
	if e.Error.Valid {
		return errors.New(e.Error.String)
	}
	return nil
}

// GetID allows EthTx to be used as jsonapi.MarshalIdentifier
func (e EthTx) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}

type EthTxAttempt struct {
	ID      int64
	EthTxID int64
	EthTx   EthTx `gorm:"foreignkey:EthTxID;->"`
	// GasPrice applies to LegacyTx
	GasPrice *utils.Big
	// GasTipCap and GasFeeCap are used instead for DynamicFeeTx
	GasTipCap *utils.Big
	GasFeeCap *utils.Big
	// ChainSpecificGasLimit on the EthTxAttempt is always the same as the on-chain encoded value for gas limit
	ChainSpecificGasLimit   uint64
	SignedRawTx             []byte
	Hash                    common.Hash
	CreatedAt               time.Time
	BroadcastBeforeBlockNum *int64
	State                   EthTxAttemptState
	EthReceipts             []EthReceipt `gorm:"foreignKey:TxHash;references:Hash;association_foreignkey:Hash;->" json:"-"`
	TxType                  int
}

// GetSignedTx decodes the SignedRawTx into a types.Transaction struct
func (a EthTxAttempt) GetSignedTx() (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(a.SignedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		logger.Error("could not decode RLP")
		return nil, err
	}
	return signedTx, nil
}

func (a EthTxAttempt) DynamicFee() gas.DynamicFee {
	return gas.DynamicFee{
		FeeCap: a.GasFeeCap.ToInt(),
		TipCap: a.GasTipCap.ToInt(),
	}
}

type EthReceipt struct {
	ID               int64
	TxHash           common.Hash
	BlockHash        common.Hash
	BlockNumber      int64
	TransactionIndex uint
	Receipt          []byte
	CreatedAt        time.Time
}
