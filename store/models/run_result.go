package models

import (
	"fmt"

	null "gopkg.in/guregu/null.v3"
)

type RunResult struct {
	Output       Output
	ErrorMessage null.String
	Pending      bool
}

type Output map[string]null.String

func (rr RunResult) Value() string {
	return rr.value().String
}

func (rr RunResult) NullValue() bool {
	return !rr.value().Valid
}

func (rr RunResult) value() null.String {
	return rr.Output["value"]
}

func RunResultWithValue(val string) RunResult {
	return RunResult{
		Output: Output{"value": null.StringFrom(val)},
	}
}

func RunResultWithError(err error) RunResult {
	return RunResult{
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

func RunResultPending(input RunResult) RunResult {
	return RunResult{
		Output:       input.Output,
		ErrorMessage: input.ErrorMessage,
		Pending:      true,
	}
}

func (rr RunResult) GetError() error {
	if rr.HasError() {
		return fmt.Errorf("Run Result: ", rr.Error())
	} else {
		return nil
	}
}

func (rr RunResult) HasError() bool {
	return rr.ErrorMessage.Valid
}

func (rr RunResult) Error() string {
	return rr.ErrorMessage.String
}

func (rr RunResult) SetError(err error) {
	rr.ErrorMessage = null.StringFrom(err.Error())
}
