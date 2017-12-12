package adapters

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type EthSendTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSendTx) Perform(input RunResult) RunResult {
	eth, err := rpc.Dial("http://example.com/api")
	if err != nil {
		return RunResult{Error: err}
	}
	var result string
	err = eth.Call(&result, "eth_sendRawTransaction", input.Value())
	if err != nil {
		return RunResult{Error: err}
	}

	return RunResultWithValue(result)
}
