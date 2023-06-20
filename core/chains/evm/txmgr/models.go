package txmgr

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// Type aliases for EVM
type (
	Confirmer              = txmgr.Confirmer[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	Broadcaster            = txmgr.Broadcaster[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	Resender               = txmgr.Resender[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	Reaper                 = txmgr.Reaper[*big.Int]
	TxStore                = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	TransactionStore       = txmgrtypes.TransactionStore[common.Address, *big.Int, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	KeyStore               = txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	TxAttemptBuilder       = txmgrtypes.TxAttemptBuilder[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	NonceSyncer            = txmgr.SequenceSyncer[common.Address, common.Hash, common.Hash]
	TransmitCheckerFactory = txmgr.TransmitCheckerFactory[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	Txm                    = txmgr.Txm[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	TxManager              = txmgr.TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	NullTxManager          = txmgr.NullTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	FwdMgr                 = txmgrtypes.ForwarderManager[common.Address]
	TxRequest              = txmgrtypes.TxRequest[common.Address, common.Hash]
	Tx                     = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	TxMeta                 = txmgrtypes.TxMeta[common.Address, common.Hash]
	TxAttempt              = txmgrtypes.TxAttempt[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	Receipt                = dbReceipt // EvmReceipt is the exported DB table model for receipts
	ReceiptPlus            = txmgrtypes.ReceiptPlus[*evmtypes.Receipt]
	TxmClient              = txmgrtypes.TxmClient[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	TransactionClient      = txmgrtypes.TransactionClient[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	ChainReceipt           = txmgrtypes.ChainReceipt[common.Hash, common.Hash]
)

var _ KeyStore = (keystore.Eth)(nil) // check interface in txmgr to avoid circular import

const (
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

// GetGethSignedTx decodes the SignedRawTx into a types.Transaction struct
func GetGethSignedTx(signedRawTx []byte) (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(signedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		return nil, err
	}
	return signedTx, nil
}
