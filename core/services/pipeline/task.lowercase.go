package pipeline

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	string
type LowercaseTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*LowercaseTask)(nil)

func (t *LowercaseTask) Type() TaskType {
	return TaskTypeLowercase
}

func (t *LowercaseTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var input StringParam

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&input, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	return Result{Value: strings.ToLower(string(input))}, runInfo
}
