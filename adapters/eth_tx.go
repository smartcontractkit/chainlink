package adapters

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink-go/models"
)

type EthSendRawTx struct {
	AdapterBase
}

type EthSignTx struct {
	AdapterBase
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSendRawTx) Perform(input models.RunResult) models.RunResult {
	eth, err := rpc.Dial(self.Store.Config.EthereumURL)
	if err != nil {
		return models.RunResultWithError(err)
	}
	var result string
	err = eth.Call(&result, "eth_sendRawTransaction", input.Value())
	if err != nil {
		return models.RunResultWithError(err)
	}

	return models.RunResultWithValue(result)
}

func (self *EthSignTx) Perform(input models.RunResult) models.RunResult {
	str := self.FunctionID + input.Value()
	data := common.FromHex(str)
	tx := types.NewTransaction(
		1,
		common.HexToAddress(self.Address),
		big.NewInt(0),
		big.NewInt(500000),
		big.NewInt(20000000000),
		data,
	)

	buffer := new(bytes.Buffer)
	err := tx.EncodeRLP(buffer)
	if err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(common.ToHex(buffer.Bytes()))
}
