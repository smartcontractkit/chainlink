package pipeline

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	*decimal.Decimal
type LengthTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*LengthTask)(nil)

func (t *LengthTask) Type() TaskType {
	return TaskTypeLength
}

func (t *LengthTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: pkgerrors.Wrap(err, "task inputs")}, runInfo
	}

	var input BytesParam

	err = errors.Join(
		pkgerrors.Wrap(ResolveParam(&input, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	return Result{Value: decimal.NewFromInt(int64(len(input)))}, runInfo
}
