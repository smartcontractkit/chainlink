package pipeline

import (
	"context"
	"encoding/hex"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	commonhex "github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	bytes
type HexDecodeTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
}

var _ Task = (*HexDecodeTask)(nil)

func (t *HexDecodeTask) Type() TaskType {
	return TaskTypeHexDecode
}

func (t *HexDecodeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
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

	if commonhex.HasPrefix(input.String()) {
		noHexPrefix := commonhex.TrimPrefix(input.String())
		bs, err := hex.DecodeString(noHexPrefix)
		if err == nil {
			return Result{Value: bs}, runInfo
		}
		return Result{Error: errors.Wrap(err, "failed to decode hex string")}, runInfo
	}

	return Result{Error: errors.New("hex string must have prefix 0x")}, runInfo
}
