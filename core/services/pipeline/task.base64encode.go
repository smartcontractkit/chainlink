package pipeline

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Return types:
//
//	string
type Base64EncodeTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*Base64EncodeTask)(nil)

func (t *Base64EncodeTask) Type() TaskType {
	return TaskTypeBase64Decode
}

func (t *Base64EncodeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var stringInput StringParam
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&stringInput, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err == nil {
		// string
		return Result{Value: base64.StdEncoding.EncodeToString([]byte(stringInput.String()))}, runInfo
	}

	var bytesInput BytesParam
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&bytesInput, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err == nil {
		// bytes
		return Result{Value: base64.StdEncoding.EncodeToString(bytesInput)}, runInfo
	}

	return Result{Error: err}, runInfo
}
