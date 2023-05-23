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
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Type aliases for EVM
type (
	EvmConfirmer              = Confirmer[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList, *assets.Wei]
	EvmBroadcaster            = Broadcaster[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList, *assets.Wei]
	EvmResender               = Resender[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee, *evmtypes.Receipt, EvmAccessList]
	EvmReaper                 = Reaper[*big.Int]
	EvmTxStore                = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EvmKeyStore               = txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	EvmTxAttemptBuilder       = txmgrtypes.TxAttemptBuilder[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EvmNonceSyncer            = NonceSyncer[common.Address, common.Hash, common.Hash]
	EvmTransmitCheckerFactory = TransmitCheckerFactory[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EvmTxm                    = Txm[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList, *assets.Wei]
	EvmTxManager              = TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	NullEvmTxManager          = NullTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EvmFwdMgr                 = txmgrtypes.ForwarderManager[common.Address]
	EvmNewTx                  = txmgrtypes.NewTx[common.Address, common.Hash]
	EvmTx                     = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EthTxMeta                 = txmgrtypes.TxMeta[common.Address, common.Hash] // TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EvmTxAttempt              = txmgrtypes.TxAttempt[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
	EvmPriorAttempt           = txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]
	EvmReceipt                = txmgrtypes.Receipt[*evmtypes.Receipt, common.Hash, common.Hash]
	EvmReceiptPlus            = txmgrtypes.ReceiptPlus[*evmtypes.Receipt]
	EvmTxmClient              = txmgrtypes.TxmClient[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]
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

// EvmAccessList is a nullable EIP2930 access list
// Used in the AdditionalParameters field in Tx
// Is optional and only has an effect on DynamicFee transactions
// on chains that support it (e.g. Ethereum Mainnet after London hard fork)
type EvmAccessList struct {
	AccessList types.AccessList
	Valid      bool
}

func EvmAccessListFrom(al types.AccessList) (n EvmAccessList) {
	if al == nil {
		return
	}
	n.AccessList = al
	n.Valid = true
	return
}

func (e EvmAccessList) MarshalJSON() ([]byte, error) {
	if !e.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(e.AccessList)
}

func (e *EvmAccessList) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		e.Valid = false
		return nil
	}
	if err := json.Unmarshal(input, &e.AccessList); err != nil {
		return errors.Wrap(err, "EvmAccessList: couldn't unmarshal JSON")
	}
	e.Valid = true
	return nil
}

// Value returns this instance serialized for database storage
func (e EvmAccessList) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, nil
	}
	return json.Marshal(e)
}

// Scan returns the selector from its serialization in the database
func (e *EvmAccessList) Scan(value interface{}) error {
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

var _ txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash] = (*EvmTxAttempt)(nil)

// GetGethSignedTx decodes the SignedRawTx into a types.Transaction struct
func GetGethSignedTx(signedRawTx []byte) (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(signedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		return nil, err
	}
	return signedTx, nil
}
