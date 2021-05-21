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

func (t *VRFTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *VRFTask) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
	// TODO
	return Result{}
}
