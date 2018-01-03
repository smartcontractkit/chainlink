package adapters

import (
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type EthTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	if !input.Pending {
		return createTxRunResult(self, input, store)
	} else {
		return ensureTxRunResult(input, store)
	}
}

func createTxRunResult(e *EthTx, input models.RunResult, store *store.Store) models.RunResult {
	data := e.FunctionID + input.Value()
	tx, err := store.Eth.CreateTx(e.Address, data)

	if err != nil {
		return models.RunResultWithError(err)
	}
	return ensureTxRunResult(models.RunResultWithValue(tx.TxID()), store)
}

func ensureTxRunResult(input models.RunResult, store *store.Store) models.RunResult {
	txid := input.Value()
	confirmed, err := store.Eth.EnsureTxConfirmed(txid)
	if err != nil {
		return models.RunResultWithError(err)
	} else if !confirmed {
		return models.RunResultPending(input)
	}
	return models.RunResultWithValue(txid)
}
