package pipeline

import (
	"context"
	"encoding/base64"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	bytes
type Base64DecodeTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*Base64DecodeTask)(nil)

func (t *Base64DecodeTask) Type() TaskType {
	return TaskTypeBase64Decode
}

func (t *Base64DecodeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: pkgerrors.Wrap(err, "task inputs")}, runInfo
	}

	var input StringParam

	err = errors.Join(
		pkgerrors.Wrap(ResolveParam(&input, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	bs, err := base64.StdEncoding.DecodeString(input.String())
	if err != nil {
		return Result{Error: pkgerrors.Wrap(err, "failed to decode base64 string")}, runInfo
	}

	return Result{Value: bs}, runInfo
}
