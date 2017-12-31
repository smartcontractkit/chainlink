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
		data := self.FunctionID + input.Value()
		tx, err := store.Eth.CreateTx(self.Address, data)

		if err != nil {
			return models.RunResultWithError(err)
		}
		txid := tx.Hash().String()
		input = models.RunResultWithValue(txid)
	}

	confirmer := &EthConfirmTx{}
	return confirmer.Perform(input, store)
}

type EthConfirmTx struct{}

func (self *EthConfirmTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	txid := input.Value()
	confirmed, err := store.Eth.TxConfirmed(txid)
	if err != nil {
		return models.RunResultWithError(err)
	} else if !confirmed {
		return models.RunResultPending(input)
	}
	return models.RunResultWithValue(txid)
}
