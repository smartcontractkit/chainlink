package adapters

import "gopkg.in/guregu/null.v3"

type Adapter interface {
	Perform(RunResult) RunResult
}

type RunResult struct {
	Output map[string]null.String
	Error  error
}

func (self RunResult) Value() string {
	return self.value().String
}

func (self RunResult) NullValue() bool {
	return !self.value().Valid
}

func (self RunResult) value() null.String {
	return self.Output["value"]
}
