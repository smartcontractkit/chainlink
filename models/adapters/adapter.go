package adapters

type Adapter interface {
	Perform(RunResult) RunResult
}

type RunResult struct {
	Output map[string]string
	Error  error
}

func (self RunResult) Value() string {
	return self.Output["value"]
}
