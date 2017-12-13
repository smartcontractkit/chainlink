package adapters

type NoOp struct {
}

func (self *NoOp) Perform(input RunResult) RunResult {
	return RunResult{}
}
