package pipeline

import (
	"context"
)

type PanicTask struct {
	BaseTask `mapstructure:",squash"`
	Msg      string
}

var _ Task = (*PanicTask)(nil)

func (t *PanicTask) Type() TaskType {
	return TaskTypePanic
}

func (t *PanicTask) Run(_ context.Context, vars Vars, _ []Result) (result Result) {
	panic(t.Msg)
}
