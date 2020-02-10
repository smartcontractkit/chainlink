package adapters

import (
	"encoding/hex"

	strpkg "chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// EthTxRawCalldata holds the Address to send the result to.
type EthTxRawCalldata struct {
	Address  common.Address `json:"address"`
	GasPrice *utils.Big     `json:"gasPrice" gorm:"type:numeric"`
	GasLimit uint64         `json:"gasLimit"`
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (etx *EthTxRawCalldata) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	if !store.TxManager.Connected() {
		return pendingConfirmationsOrConnection(input)
	}

	if input.Status().PendingConfirmations() {
		return ensureTxRunResult(input, store)
	}

	result := input.Result().String()
	data, err := hex.DecodeString(result)
	if err != nil {
		err = errors.Wrap(err, "while decoding tx data from hex")
		return models.NewRunOutputError(err)
	}

	return createTxRunResult(etx.Address, etx.GasPrice, etx.GasLimit, data, input, store)
}
