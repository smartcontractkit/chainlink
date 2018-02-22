package adapters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// EthTx holds the Address to send the result to and the FunctionSelector
// to execute.
type EthTx struct {
	Address          common.Address          `json:"address"`
	FunctionSelector models.FunctionSelector `json:"functionSelector"`
	DataPrefix       hexutil.Bytes           `json:"dataPrefix"`
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (etx *EthTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	if !input.Pending {
		return createTxRunResult(etx, input, store)
	} else {
		return ensureTxRunResult(input, store)
	}
}

func createTxRunResult(
	e *EthTx,
	input models.RunResult,
	store *store.Store,
) models.RunResult {
	val, err := input.Value()
	if err != nil {
		return input.WithError(err)
	}

	data, err := utils.HexToBytes(e.FunctionSelector.String(), e.DataPrefix.String(), val)
	if err != nil {
		return input.WithError(err)
	}

	attempt, err := store.TxManager.CreateTx(e.Address, data)
	if err != nil {
		return input.WithError(err)
	}

	sendResult := input.WithValue(attempt.Hash.String())
	return ensureTxRunResult(sendResult, store)
}

func ensureTxRunResult(input models.RunResult, store *store.Store) models.RunResult {
	val, err := input.Value()
	if err != nil {
		return input.WithError(err)
	}

	hash := common.HexToHash(val)
	if err != nil {
		return input.WithError(err)
	}

	confirmed, err := store.TxManager.EnsureTxConfirmed(hash)

	if err != nil {
		return input.WithError(err)
	} else if !confirmed {
		return input.MarkPending()
	}
	return input.WithValue(hash.String())
}
