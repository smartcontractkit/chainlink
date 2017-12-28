package adapters

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type EthSignTx struct {
	AdapterBase
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignTx) Perform(input models.RunResult) models.RunResult {
	str := self.FunctionID + input.Value()
	data := common.FromHex(str)
	keyStore := self.Store.KeyStore
	nonce, err := self.Store.Eth.GetNonce(keyStore.GetAccount())
	if err != nil {
		return models.RunResultWithError(err)
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(self.Address),
		big.NewInt(0),
		big.NewInt(500000),
		big.NewInt(20000000000),
		data,
	)
	signedTx, err := keyStore.SignTx(tx, self.Store.Config.ChainID)
	if err != nil {
		return models.RunResultWithError(err)
	}

	buffer := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(buffer); err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(common.ToHex(buffer.Bytes()))
}

type EthSendRawTx struct {
	AdapterBase
}

func (self *EthSendRawTx) Perform(input models.RunResult) models.RunResult {
	result, err := self.Store.Eth.SendRawTx(input.Value())
	if err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(result)
}

type EthConfirmTx struct {
	AdapterBase
}

func (self *EthConfirmTx) Perform(input models.RunResult) models.RunResult {
	txid := input.Value()
	for {
		receipt, err := self.Store.Eth.GetTxReceipt(txid)
		if err != nil {
			return models.RunResultWithError(err)
		} else if receipt.TxHash.Hex() == "" {
			time.Sleep(1000 * time.Millisecond)
		} else {
			return models.RunResultWithValue(txid)
		}
	}
}

type EthSignAndSendTx struct {
	AdapterBase
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignAndSendTx) Perform(input models.RunResult) models.RunResult {
	signer := &EthSignTx{
		Address:     self.Address,
		FunctionID:  self.FunctionID,
		AdapterBase: AdapterBase{self.Store},
	}
	sender := &EthSendRawTx{
		AdapterBase: AdapterBase{self.Store},
	}
	confirmer := &EthConfirmTx{
		AdapterBase: AdapterBase{self.Store},
	}

	signed := signer.Perform(input)
	if signed.HasError() {
		return signed
	}
	sent := sender.Perform(signed)
	if sent.HasError() {
		return sent
	}
	return confirmer.Perform(sent)
}
