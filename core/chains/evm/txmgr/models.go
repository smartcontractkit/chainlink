package txmgr

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Type aliases for EVM
type (
	EvmConfirmer              = EthConfirmer[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmBroadcaster            = EthBroadcaster[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmResender               = EthResender[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee, *evmtypes.Receipt, NullableEIP2930AccessList]
	EvmTxStore                = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmKeyStore               = txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	EvmTxAttemptBuilder       = txmgrtypes.TxAttemptBuilder[*evmtypes.Head, gas.EvmFee, common.Address, common.Hash, EvmTx, EvmTxAttempt, evmtypes.Nonce]
	EvmNonceSyncer            = NonceSyncer[common.Address, common.Hash, common.Hash]
	EvmTransmitCheckerFactory = TransmitCheckerFactory[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmTxm                    = Txm[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, NullableEIP2930AccessList]
	EvmTxManager              = TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	NullEvmTxManager          = NullTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmFwdMgr                 = txmgrtypes.ForwarderManager[common.Address]
	EvmNewTx                  = txmgrtypes.NewTx[common.Address, common.Hash]
	EvmTx                     = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EthTxMeta                 = txmgrtypes.TxMeta[common.Address, common.Hash] // TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EvmTxAttempt              = txmgrtypes.TxAttempt[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, gas.EvmFee, NullableEIP2930AccessList]
	EvmPriorAttempt           = txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]
	EvmReceipt                = txmgrtypes.Receipt[*evmtypes.Receipt, common.Hash, common.Hash]
	EvmReceiptPlus            = txmgrtypes.ReceiptPlus[*evmtypes.Receipt]
)

const (
	// TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EthTxUnstarted               = txmgrtypes.TxState("unstarted")
	EthTxInProgress              = txmgrtypes.TxState("in_progress")
	EthTxFatalError              = txmgrtypes.TxState("fatal_error")
	EthTxUnconfirmed             = txmgrtypes.TxState("unconfirmed")
	EthTxConfirmed               = txmgrtypes.TxState("confirmed")
	EthTxConfirmedMissingReceipt = txmgrtypes.TxState("confirmed_missing_receipt")

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

var _ txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash] = EvmTxAttempt{}

// GetGethSignedTx decodes the SignedRawTx into a types.Transaction struct
func GetGethSignedTx(signedRawTx []byte) (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(signedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		return nil, err
	}
	return signedTx, nil
}
