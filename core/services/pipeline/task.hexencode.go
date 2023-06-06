package pipeline

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	string
type HexEncodeTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*HexEncodeTask)(nil)

func (t *HexEncodeTask) Type() TaskType {
	return TaskTypeHexEncode
}

func addHexPrefix(val string) string {
	if len(val) > 0 {
		return "0x" + val
	}
	return ""
}

func (t *HexEncodeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
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
		return Result{Value: addHexPrefix(hex.EncodeToString([]byte(stringInput.String())))}, runInfo
	}

	var bytesInput BytesParam
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&bytesInput, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err == nil {
		// bytes
		return Result{Value: addHexPrefix(hex.EncodeToString(bytesInput))}, runInfo
	}

	var decimalInput DecimalParam
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&decimalInput, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err == nil && !decimalInput.Decimal().IsInteger() {
		// decimal
		return Result{Error: errors.New("decimal input")}, runInfo
	}

	var bigIntInput MaybeBigIntParam
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&bigIntInput, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
	)
	if err == nil {
		// one of integer types
		if bigIntInput.BigInt().Sign() == -1 {
			return Result{Error: errors.New("negative integer")}, runInfo
		}
		return Result{Value: addHexPrefix(fmt.Sprintf("%x", bigIntInput.BigInt()))}, runInfo
	}

	return Result{Error: err}, runInfo
}
