package adapters

type Adapter interface {
	Perform(RunResult) RunResult
}

type RunResult struct {
	Output map[string]string
	Error  error
}
