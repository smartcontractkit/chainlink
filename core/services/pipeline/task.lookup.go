package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Look up a field on a map
//
// Return types:
//
// interface{}
type LookupTask struct {
	BaseTask `mapstructure:",squash"`
	Key      string `json:"key"`
}

var _ Task = (*LookupTask)(nil)

func (t *LookupTask) Type() TaskType {
	return TaskTypeLookup
}

func (t *LookupTask) Run(_ context.Context, _ logger.Logger, _ Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 1, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var key StringParam
	err = errors.Wrap(ResolveParam(&key, From(t.Key)), "key")
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var val interface{}
	switch m := inputs[0].Value.(type) {
	case map[string]interface{}:
		val = m[(string)(key)]
	default:
		return Result{Error: errors.Errorf("unexpected input type: %T", inputs[0].Value)}, runInfo
	}

	return Result{Value: val}, runInfo
}
