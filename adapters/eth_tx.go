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
	attempt, err := store.Eth.CreateTx(e.Address, data)

	if err != nil {
		return models.RunResultWithError(err)
	}
	return ensureTxRunResult(models.RunResultWithValue(attempt.Hash), store)
}

func ensureTxRunResult(input models.RunResult, store *store.Store) models.RunResult {
	hash := input.Value()
	confirmed, err := store.Eth.EnsureTxConfirmed(hash)

	if err != nil {
		return models.RunResultWithError(err)
	} else if !confirmed {
		return models.RunResultPending(input)
	}
	return models.RunResultWithValue(hash)
}
