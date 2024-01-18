package pipeline

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	map[string]interface{}
type MergeTask struct {
	BaseTask `mapstructure:",squash"`
	Left     string `json:"left"`
	Right    string `json:"right"`
}

var _ Task = (*MergeTask)(nil)

func (t *MergeTask) Type() TaskType {
	return TaskTypeMerge
}

func (t *MergeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: pkgerrors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		lMap MapParam
		rMap MapParam
	)
	err = errors.Join(
		pkgerrors.Wrap(ResolveParam(&lMap, From(VarExpr(t.Left, vars), JSONWithVarExprs(t.Left, vars, false), Input(inputs, 0))), "left-side"),
		pkgerrors.Wrap(ResolveParam(&rMap, From(VarExpr(t.Right, vars), JSONWithVarExprs(t.Right, vars, false))), "right-side"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	// clobber lMap with rMap values
	// "nil" values on the right will clobber
	for key, value := range rMap {
		lMap[key] = value
	}

	return Result{Value: lMap.Map()}, runInfo
}
