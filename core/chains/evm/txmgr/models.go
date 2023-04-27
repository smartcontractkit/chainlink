package txmgr

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Type aliases for EVM
type (
	EvmConfirmer              = EthConfirmer[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	EvmBroadcaster            = EthBroadcaster[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	EvmResender               = EthResender[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce]
	EvmTxStore                = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, EvmNewTx, *evmtypes.Receipt, EvmTx, EvmTxAttempt, evmtypes.Nonce]
	EvmKeyStore               = txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	EvmTxAttemptBuilder       = txmgrtypes.TxAttemptBuilder[*evmtypes.Head, gas.EvmFee, common.Address, common.Hash, EvmTx, EvmTxAttempt, evmtypes.Nonce]
	EvmNonceSyncer            = NonceSyncer[common.Address, common.Hash, common.Hash]
	EvmTransmitCheckerFactory = TransmitCheckerFactory[common.Address, common.Hash]
	EvmTxm                    = Txm[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	EvmTxManager              = TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash]
	NullEvmTxManager          = NullTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash]
	EvmFwdMgr                 = txmgrtypes.ForwarderManager[common.Address]
	EvmNewTx                  = txmgrtypes.NewTx[common.Address, common.Hash]
	EvmTx                     = EthTx[common.Address, common.Hash]
	EthTxMeta                 = txmgrtypes.TxMeta[common.Address, common.Hash] // TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EvmTxAttempt              = EthTxAttempt[common.Address, common.Hash]
	EvmPriorAttempt           = txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]
	EvmReceipt                = txmgrtypes.Receipt[*evmtypes.Receipt, common.Hash, common.Hash]
	EvmReceiptPlus            = txmgrtypes.ReceiptPlus[*evmtypes.Receipt]
)

type EthTxState string

const (
	EthTxUnstarted               = EthTxState("unstarted")
	EthTxInProgress              = EthTxState("in_progress")
	EthTxFatalError              = EthTxState("fatal_error")
	EthTxUnconfirmed             = EthTxState("unconfirmed")
	EthTxConfirmed               = EthTxState("confirmed")
	EthTxConfirmedMissingReceipt = EthTxState("confirmed_missing_receipt")

	// TransmitCheckerTypeSimulate is a checker that simulates the transaction before executing on
	// chain.
	TransmitCheckerTypeSimulate = txmgrtypes.TransmitCheckerType("simulate")

	// TransmitCheckerTypeVRFV1 is a checker that will not submit VRF V1 fulfillment requests that
	// have already been fulfilled. This could happen if the request was fulfilled by another node.
	TransmitCheckerTypeVRFV1 = txmgrtypes.TransmitCheckerType("vrf_v1")

	// TransmitCheckerTypeVRFV2 is a checker that will not submit VRF V2 fulfillment requests that
	// have already been fulfilled. This could happen if the request was fulfilled by another node.
	TransmitCheckerTypeVRFV2 = txmgrtypes.TransmitCheckerType("vrf_v2")
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

type EthTx[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable] struct {
	txmgrtypes.Transaction
	ID             int64
	Nonce          *int64
	FromAddress    ADDR
	ToAddress      ADDR
	EncodedPayload []byte
	Value          assets.Eth
	// GasLimit on the EthTx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	GasLimit uint32
	Error    null.String
	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt *time.Time
	// InitialBroadcastAt is recorded once, the first ever time this eth_tx is sent
	InitialBroadcastAt *time.Time
	CreatedAt          time.Time
	State              EthTxState
	EthTxAttempts      []EthTxAttempt[ADDR, TX_HASH] `json:"-"`
	// Marshalled EthTxMeta
	// Used for additional context around transactions which you want to log
	// at send time.
	Meta       *datatypes.JSON
	Subject    uuid.NullUUID
	EVMChainID utils.Big

	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  cnull.Uint32

	// AccessList is optional and only has an effect on DynamicFee transactions
	// on chains that support it (e.g. Ethereum Mainnet after London hard fork)
	AccessList NullableEIP2930AccessList

	// TransmitChecker defines the check that should be performed before a transaction is submitted on
	// chain.
	TransmitChecker *datatypes.JSON
}

func (e EthTx[ADDR, TX_HASH]) GetError() error {
	if e.Error.Valid {
		return errors.New(e.Error.String)
	}
	return nil
}

// GetID allows EthTx to be used as jsonapi.MarshalIdentifier
func (e EthTx[ADDR, TX_HASH]) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}

// GetMeta returns an EthTx's meta in struct form, unmarshalling it from JSON first.
func (e EthTx[ADDR, TX_HASH]) GetMeta() (*EthTxMeta, error) {
	if e.Meta == nil {
		return nil, nil
	}
	var m EthTxMeta
	return &m, errors.Wrap(json.Unmarshal(*e.Meta, &m), "unmarshalling meta")
}

// GetLogger returns a new logger with metadata fields.
func (e EthTx[ADDR, TX_HASH]) GetLogger(lgr logger.Logger) logger.Logger {
	lgr = lgr.With(
		"ethTxID", e.ID,
		"nonce", e.Nonce,
		"checker", e.TransmitChecker,
		"gasLimit", e.GasLimit,
	)

	meta, err := e.GetMeta()
	if err != nil {
		lgr.Errorw("failed to get meta of the transaction", "err", err)
		return lgr
	}

	if meta != nil {
		lgr = lgr.With("jobID", meta.JobID)

		if meta.RequestTxHash != nil {
			lgr = lgr.With("requestTxHash", *meta.RequestTxHash)
		}

		if meta.RequestID != nil {
			lgr = lgr.With("requestID", new(big.Int).SetBytes(meta.RequestID[:]).String())
		}

		if len(meta.RequestIDs) != 0 {
			var ids []string
			for _, id := range meta.RequestIDs {
				ids = append(ids, new(big.Int).SetBytes(id[:]).String())
			}
			lgr = lgr.With("requestIDs", strings.Join(ids, ","))
		}

		if meta.UpkeepID != nil {
			lgr = lgr.With("upkeepID", *meta.UpkeepID)
		}

		if meta.SubID != nil {
			lgr = lgr.With("subID", *meta.SubID)
		}

		if meta.MaxLink != nil {
			lgr = lgr.With("maxLink", *meta.MaxLink)
		}

		if meta.FwdrDestAddress != nil {
			lgr = lgr.With("FwdrDestAddress", *meta.FwdrDestAddress)
		}

		if len(meta.MessageIDs) > 0 {
			for _, mid := range meta.MessageIDs {
				lgr = lgr.With("messageID", mid)
			}
		}

		if len(meta.SeqNumbers) > 0 {
			lgr = lgr.With("SeqNumbers", meta.SeqNumbers)
		}
	}

	return lgr
}

// GetChecker returns an EthTx's transmit checker spec in struct form, unmarshalling it from JSON
// first.
func (e EthTx[ADDR, TX_HASH]) GetChecker() (txmgrtypes.TransmitCheckerSpec[ADDR], error) {
	if e.TransmitChecker == nil {
		return txmgrtypes.TransmitCheckerSpec[ADDR]{}, nil
	}
	var t txmgrtypes.TransmitCheckerSpec[ADDR]
	return t, errors.Wrap(json.Unmarshal(*e.TransmitChecker, &t), "unmarshalling transmit checker")
}

var _ txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash] = EthTxAttempt[common.Address, common.Hash]{}

type EthTxAttempt[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable] struct {
	ID      int64
	EthTxID int64
	EthTx   EthTx[ADDR, TX_HASH]
	// GasPrice applies to LegacyTx
	GasPrice *assets.Wei
	// GasTipCap and GasFeeCap are used instead for DynamicFeeTx
	GasTipCap *assets.Wei
	GasFeeCap *assets.Wei
	// ChainSpecificGasLimit on the EthTxAttempt is always the same as the on-chain encoded value for gas limit
	ChainSpecificGasLimit   uint32
	SignedRawTx             []byte
	Hash                    TX_HASH
	CreatedAt               time.Time
	BroadcastBeforeBlockNum *int64
	State                   txmgrtypes.TxAttemptState
	EthReceipts             []EvmReceipt `json:"-"`
	TxType                  int
}

func (a EthTxAttempt[ADDR, TX_HASH]) String() string {
	return fmt.Sprintf("EthTxAttempt(ID:%d,EthTxID:%d,GasPrice:%v,GasTipCap:%v,GasFeeCap:%v,TxType:%d", a.ID, a.EthTxID, a.GasPrice, a.GasTipCap, a.GasFeeCap, a.TxType)
}

// GetSignedTx decodes the SignedRawTx into a types.Transaction struct
func (a EthTxAttempt[ADDR, TX_HASH]) GetSignedTx() (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(a.SignedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (a EthTxAttempt[ADDR, TX_HASH]) Fee() (fee gas.EvmFee) {
	fee.Legacy = a.getGasPrice()

	dynamic := a.dynamicFee()
	// add dynamic struct only if values are not nil
	if dynamic.FeeCap != nil && dynamic.TipCap != nil {
		fee.Dynamic = &dynamic
	}
	return fee
}

func (a EthTxAttempt[ADDR, TX_HASH]) dynamicFee() gas.DynamicFee {
	return gas.DynamicFee{
		FeeCap: a.GasFeeCap,
		TipCap: a.GasTipCap,
	}
}

func (a EthTxAttempt[ADDR, TX_HASH]) GetBroadcastBeforeBlockNum() *int64 {
	return a.BroadcastBeforeBlockNum
}

func (a EthTxAttempt[ADDR, TX_HASH]) GetChainSpecificGasLimit() uint32 {
	return a.ChainSpecificGasLimit
}

func (a EthTxAttempt[ADDR, TX_HASH]) getGasPrice() *assets.Wei {
	return a.GasPrice
}

func (a EthTxAttempt[ADDR, TX_HASH]) GetHash() TX_HASH {
	return a.Hash
}

func (a EthTxAttempt[ADDR, TX_HASH]) GetTxType() int {
	return a.TxType
}
