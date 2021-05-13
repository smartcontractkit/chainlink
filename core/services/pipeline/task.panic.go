package pipeline

import "context"

type PanicTask struct {
	BaseTask `mapstructure:",squash"`
	Msg      string
}

var _ Task = (*PanicTask)(nil)

func (t *PanicTask) Type() TaskType {
	return TaskTypePanic
}

func (t *PanicTask) SetDefaults(_ map[string]string, _ TaskDAG, _ TaskDAGNode) error {
	return nil
}

func (t *PanicTask) Run(_ context.Context, _ JSONSerializable, _ []Result) (result Result) {
	panic(t.Msg)
}
