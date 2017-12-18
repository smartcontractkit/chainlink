package models

import (
	"fmt"

	null "gopkg.in/guregu/null.v3"
)

type RunResult struct {
	Output       Output
	ErrorMessage null.String
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

func RunResultWithError(err error) RunResult {
	return RunResult{
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

func (self RunResult) GetError() error {
	if self.HasError() {
		return fmt.Errorf("Run Result: ", self.Error())
	} else {
		return nil
	}
}

func (self RunResult) HasError() bool {
	return self.ErrorMessage.Valid
}

func (self RunResult) Error() string {
	return self.ErrorMessage.String
}

func (self RunResult) SetError(err error) {
	self.ErrorMessage = null.StringFrom(err.Error())
}
