package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// FailTask is like the Panic task but without all the drama and stack
// unwinding of a panic
type FailTask struct {
	BaseTask `mapstructure:",squash"`
	Msg      string
}

var _ Task = (*FailTask)(nil)

func (t *FailTask) Type() TaskType {
	return TaskTypeFail
}

func (t *FailTask) Run(_ context.Context, _ logger.Logger, vars Vars, _ []Result) (Result, RunInfo) {
	return Result{Error: errors.New(t.Msg)}, RunInfo{}
}
