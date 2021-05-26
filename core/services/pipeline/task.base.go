package pipeline

import "time"

type BaseTask struct {
	outputTask      Task
	dotID           string        `mapstructure:"-"`
	numPredecessors int           `mapstructure:"-"`
	Index           int32         `mapstructure:"index" json:"-" `
	Timeout         time.Duration `mapstructure:"timeout"`
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
