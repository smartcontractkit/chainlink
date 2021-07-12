package pipeline

import (
	"context"
)

type VRFTask struct {
	BaseTask `mapstructure:",squash"`
}

var _ Task = (*VRFTask)(nil)

func (t *VRFTask) Type() TaskType {
	return TaskTypeVRF
}

func (t *VRFTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	// TODO
	return Result{}
}
