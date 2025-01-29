package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// ConditionalTask checks if data is false
// for now this is all we need but in the future we can
// expand this to handle more general conditional statements
type ConditionalTask struct {
	BaseTask `mapstructure:",squash"`
	Data     string `json:"data"`
}

var _ Task = (*ConditionalTask)(nil)

func (t *ConditionalTask) Type() TaskType {
	return TaskTypeConditional
}

func (t *ConditionalTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}
	var (
		boolParam BoolParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&boolParam, From(VarExpr(t.Data, vars), Input(inputs, 0), nil)), "data"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	if !boolParam {
		return Result{Error: errors.New("conditional was not satisfied")}, runInfo
	}
	return Result{Value: true}, runInfo
}
