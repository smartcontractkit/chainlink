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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
)

// Type aliases for EVM
type (
	EvmConfirmer              = EthConfirmer[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmBroadcaster            = EthBroadcaster[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmResender               = EthResender[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee, *evmtypes.Receipt, NullableEIP2930AccessList]
	EvmTxStore                = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, EvmNewTx, *evmtypes.Receipt, EvmTx, EvmTxAttempt, evmtypes.Nonce]
	EvmKeyStore               = txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	EvmTxAttemptBuilder       = txmgrtypes.TxAttemptBuilder[*evmtypes.Head, gas.EvmFee, common.Address, common.Hash, EvmTx, EvmTxAttempt, evmtypes.Nonce]
	EvmNonceSyncer            = NonceSyncer[common.Address, common.Hash, common.Hash]
	EvmTransmitCheckerFactory = TransmitCheckerFactory[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmTxm                    = Txm[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmTxManager              = TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	NullEvmTxManager          = NullTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmFwdMgr                 = txmgrtypes.ForwarderManager[common.Address]
	EvmNewTx                  = txmgrtypes.NewTx[common.Address, common.Hash]
	EvmTx                     = Tx[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EthTxMeta                 = txmgrtypes.TxMeta[common.Address, common.Hash] // TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EvmTxAttempt              = TxAttempt[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmPriorAttempt           = txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]
	EvmReceipt                = txmgrtypes.Receipt[*evmtypes.Receipt, common.Hash, common.Hash]
	EvmReceiptPlus            = txmgrtypes.ReceiptPlus[*evmtypes.Receipt]
)

type TxState string

const (
	EthTxUnstarted               = TxState("unstarted")
	EthTxInProgress              = TxState("in_progress")
	EthTxFatalError              = TxState("fatal_error")
	EthTxUnconfirmed             = TxState("unconfirmed")
	EthTxConfirmed               = TxState("confirmed")
	EthTxConfirmedMissingReceipt = TxState("confirmed_missing_receipt")

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

// NullableEIP2930AccessList is used in the AdditionalParameters field in Tx
// NullableEIP2930AccessList is optional and only has an effect on DynamicFee transactions
// on chains that support it (e.g. Ethereum Mainnet after London hard fork)
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

type Tx[
	CHAIN_ID txmgrtypes.ID,
	ADDR commontypes.Hashable,
	TX_HASH, BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH],
	FEE txmgrtypes.Fee,
	ADD any,
] struct {
	txmgrtypes.Transaction
	ID             int64
	Sequence       *int64
	FromAddress    ADDR
	ToAddress      ADDR
	EncodedPayload []byte
	// Value          assets.Eth // TODO
	Value big.Int
	// FeeLimit on the Tx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	FeeLimit uint32
	Error    null.String
	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt *time.Time
	// InitialBroadcastAt is recorded once, the first ever time this eth_tx is sent
	InitialBroadcastAt *time.Time
	CreatedAt          time.Time
	State              TxState
	TxAttempts         []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD] `json:"-"`
	// Marshalled TxMeta
	// Used for additional context around transactions which you want to log
	// at send time.
	Meta    *datatypes.JSON
	Subject uuid.NullUUID
	ChainID CHAIN_ID

	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  cnull.Uint32

	// AdditionalParameters is generic type that supports passing miscellaneous parameters
	// as a part of the TX struct that may be used inside chain-specific components
	// example: EVM + NullableEIP2930AccessList which only affects on DynamicFee transactions
	AdditionalParameters ADD

	// TransmitChecker defines the check that should be performed before a transaction is submitted on
	// chain.
	TransmitChecker *datatypes.JSON
}

func (e Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetError() error {
	if e.Error.Valid {
		return errors.New(e.Error.String)
	}
	return nil
}

// GetID allows EthTx to be used as jsonapi.MarshalIdentifier
func (e Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}

// GetMeta returns an EthTx's meta in struct form, unmarshalling it from JSON first.
func (e Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetMeta() (*EthTxMeta, error) {
	if e.Meta == nil {
		return nil, nil
	}
	var m EthTxMeta
	return &m, errors.Wrap(json.Unmarshal(*e.Meta, &m), "unmarshalling meta")
}

// GetLogger returns a new logger with metadata fields.
func (e Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetLogger(lgr logger.Logger) logger.Logger {
	lgr = lgr.With(
		"txID", e.ID,
		"sequence", e.Sequence,
		"checker", e.TransmitChecker,
		"feeLimit", e.FeeLimit,
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
func (e Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetChecker() (txmgrtypes.TransmitCheckerSpec[ADDR], error) {
	if e.TransmitChecker == nil {
		return txmgrtypes.TransmitCheckerSpec[ADDR]{}, nil
	}
	var t txmgrtypes.TransmitCheckerSpec[ADDR]
	return t, errors.Wrap(json.Unmarshal(*e.TransmitChecker, &t), "unmarshalling transmit checker")
}

var _ txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash] = EvmTxAttempt{}

type TxAttempt[
	CHAIN_ID txmgrtypes.ID,
	ADDR commontypes.Hashable,
	TX_HASH, BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH],
	FEE txmgrtypes.Fee,
	ADD any,
] struct {
	ID    int64
	TxID  int64
	Tx    Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]
	TxFee FEE
	// ChainSpecificFeeLimit on the TxAttempt is always the same as the on-chain encoded value for gas limit
	ChainSpecificFeeLimit   uint32
	SignedRawTx             []byte
	Hash                    TX_HASH
	CreatedAt               time.Time
	BroadcastBeforeBlockNum *int64
	State                   txmgrtypes.TxAttemptState
	Receipts                []txmgrtypes.Receipt[R, TX_HASH, BLOCK_HASH] `json:"-"`
	TxType                  int
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) String() string {
	return fmt.Sprintf("TxAttempt(ID:%d,TxID:%d,Fee:%s,TxType:%d", a.ID, a.TxID, a.TxFee, a.TxType)
}

// GetSignedTx decodes the SignedRawTx into a types.Transaction struct
func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetSignedTx() (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(a.SignedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) Fee() FEE {
	return a.TxFee
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetBroadcastBeforeBlockNum() *int64 {
	return a.BroadcastBeforeBlockNum
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetChainSpecificFeeLimit() uint32 {
	return a.ChainSpecificFeeLimit
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetHash() TX_HASH {
	return a.Hash
}

func (a TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) GetTxType() int {
	return a.TxType
}
