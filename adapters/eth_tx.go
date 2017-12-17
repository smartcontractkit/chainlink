package adapters

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink-go/config"
	"github.com/smartcontractkit/chainlink-go/models"
)

type EthSendTx struct {
	Config     config.Config
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSendTx) Perform(input models.RunResult) models.RunResult {
	eth, err := rpc.Dial("http://example.com/api")
	if err != nil {
		return models.RunResult{Error: err}
	}
	var result string
	err = eth.Call(&result, "eth_sendRawTransaction", input.Value())
	if err != nil {
		return models.RunResult{Error: err}
	}

	return models.RunResultWithValue(result)
}
