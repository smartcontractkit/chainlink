package pipeline

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// CoalesceTask returns the first non-nil, non-error input, or nil if there are none.
type CoalesceTask struct {
	BaseTask `mapstructure:",squash"`
}

var _ Task = (*CoalesceTask)(nil)

func (t *CoalesceTask) Type() TaskType {
	return TaskTypeCoalesce
}

func (t *CoalesceTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	for _, input := range inputs {
		if input.Error == nil && input.Value != nil {
			return input, runInfo
		}
	}
	return Result{Value: nil}, runInfo
}
