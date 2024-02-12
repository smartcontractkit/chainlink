package txm

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

const (
	TxUnstarted               = txmgrtypes.TxState("unstarted")
	TxInProgress              = txmgrtypes.TxState("in_progress")
	TxFatalError              = txmgrtypes.TxState("fatal_error")
	TxUnconfirmed             = txmgrtypes.TxState("unconfirmed")
	TxConfirmed               = txmgrtypes.TxState("confirmed")
	TxConfirmedMissingReceipt = txmgrtypes.TxState("confirmed_missing_receipt")
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
