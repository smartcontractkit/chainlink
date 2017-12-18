package adapters

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink-go/models"
)

type EthSendTx struct {
	AdapterBase
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSendTx) Perform(input models.RunResult) models.RunResult {
	eth, err := rpc.Dial(self.Config.EthereumURL)
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
