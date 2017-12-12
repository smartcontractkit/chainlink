package adapters

type EthSendTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSendTx) Perform(input RunResult) RunResult {
	return RunResult{}
}
