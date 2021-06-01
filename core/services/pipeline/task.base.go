package pipeline

import "time"

type BaseTask struct {
	outputs []Task
	inputs  []Task

	id         int64
	outputTask Task
	dotID      string
	Index      int32         `mapstructure:"index" json:"-" `
	Timeout    time.Duration `mapstructure:"timeout"`
}

func NewBaseTask(id int64, dotID string, t Task, index int32) BaseTask {
	return BaseTask{dotID: dotID, outputTask: t, Index: index}
}

func (t *BaseTask) Base() *BaseTask {
	return t
}

func (t BaseTask) ID() int64 {
	return t.id
}

func (t BaseTask) DotID() string {
	return t.dotID
}

func (t BaseTask) OutputIndex() int32 {
	return t.Index
}

func (t BaseTask) OutputTask() Task {
	return t.outputTask
}

func (t BaseTask) Outputs() []Task {
	return t.outputs
}

func (t BaseTask) Inputs() []Task {
	return t.inputs
}

func (t *BaseTask) SetOutputTask(outputTask Task) {
	t.outputTask = outputTask
}

func (t BaseTask) TaskTimeout() (time.Duration, bool) {
	if t.Timeout == time.Duration(0) {
		return time.Duration(0), false
	}
	return t.Timeout, true
}
