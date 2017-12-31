package adapters

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type EthSignTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	data := self.FunctionID + input.Value()
	tx, err := store.Tx.NewSignedTx(self.Address, data)
	if err != nil {
		return models.RunResultWithError(err)
	}
	buffer := new(bytes.Buffer)
	if err := tx.EncodeRLP(buffer); err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(common.ToHex(buffer.Bytes()))
}

type EthSendRawTx struct{}

func (self *EthSendRawTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	result, err := store.Tx.SendRawTx(input.Value())
	if err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(result)
}

type EthConfirmTx struct{}

func (self *EthConfirmTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	txid := input.Value()
	confirmed, err := store.Tx.TxConfirmed(txid)
	if err != nil {
		return models.RunResultWithError(err)
	} else if !confirmed {
		return models.RunResultPending(input)
	}
	return models.RunResultWithValue(txid)
}

type EthSignAndSendTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignAndSendTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	signer := &EthSignTx{
		Address:    self.Address,
		FunctionID: self.FunctionID,
	}
	sender := &EthSendRawTx{}
	confirmer := &EthConfirmTx{}

	if !input.Pending {
		signed := signer.Perform(input, store)
		if signed.HasError() {
			return signed
		}
		sent := sender.Perform(signed, store)
		if sent.HasError() {
			return sent
		}
		return confirmer.Perform(sent, store)
	} else {
		return confirmer.Perform(input, store)
	}
}
