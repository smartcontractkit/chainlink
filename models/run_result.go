package models

import (
	null "gopkg.in/guregu/null.v3"
)

type RunResult struct {
	Output Output
	Error  error
}

type Output map[string]null.String

func (self RunResult) Value() string {
	return self.value().String
}

func (self RunResult) NullValue() bool {
	return !self.value().Valid
}

func (self RunResult) value() null.String {
	return self.Output["value"]
}

func RunResultWithValue(val string) RunResult {
	return RunResult{
		Output: Output{"value": null.StringFrom(val)},
	}
}
