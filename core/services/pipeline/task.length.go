package pipeline

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	input, err := resolveParamValue(From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0)))
	if err != nil {
		return Result{Error: errors.Wrap(err, "input")}, runInfo
	}

	if input != nil {
		v := reflect.ValueOf(input)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			return Result{Value: decimal.NewFromInt(int64(v.Len()))}, runInfo
		}
	}

	var sliceInput SliceParam
	if err = sliceInput.UnmarshalPipelineParam(input); err == nil {
		return Result{Value: decimal.NewFromInt(int64(len(sliceInput)))}, runInfo
	}

	var bytesInput BytesParam
	if err = bytesInput.UnmarshalPipelineParam(input); err != nil {
		return Result{Error: errors.Wrap(err, "input")}, runInfo
	}

	return Result{Value: decimal.NewFromInt(int64(len(bytesInput)))}, runInfo
}
