package pipeline

import "time"

type BaseTask struct {
	outputTask      Task
	dotID           string
	numPredecessors int
	Index           int32         `mapstructure:"index" json:"-" `
	Timeout         time.Duration `mapstructure:"timeout"`
}

func NewBaseTask(dotID string, t Task, index int32, numPredecessors int) BaseTask {
	return BaseTask{dotID: dotID, outputTask: t, Index: index, numPredecessors: numPredecessors}
}

func (t BaseTask) NumPredecessors() int {
	return t.numPredecessors
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

func (t *BaseTask) SetOutputTask(outputTask Task) {
	t.outputTask = outputTask
}

func (t BaseTask) TaskTimeout() (time.Duration, bool) {
	if t.Timeout == time.Duration(0) {
		return time.Duration(0), false
	}
	return t.Timeout, true
}
