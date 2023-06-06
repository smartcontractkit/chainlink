package pipeline

import (
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/v2/core/null"
)

type BaseTask struct {
	outputs []Task
	inputs  []TaskDependency

	id        int
	dotID     string
	Index     int32          `mapstructure:"index" json:"-" `
	Timeout   *time.Duration `mapstructure:"timeout"`
	FailEarly bool           `mapstructure:"failEarly"`

	Retries    null.Uint32   `mapstructure:"retries"`
	MinBackoff time.Duration `mapstructure:"minBackoff"`
	MaxBackoff time.Duration `mapstructure:"maxBackoff"`

	uuid uuid.UUID
}

func NewBaseTask(id int, dotID string, inputs []TaskDependency, outputs []Task, index int32) BaseTask {
	return BaseTask{id: id, dotID: dotID, inputs: inputs, outputs: outputs, Index: index}
}

func (t *BaseTask) Base() *BaseTask {
	return t
}

func (t BaseTask) ID() int {
	return t.id
}

func (t BaseTask) DotID() string {
	return t.dotID
}

func (t BaseTask) OutputIndex() int32 {
	return t.Index
}

func (t BaseTask) Outputs() []Task {
	return t.outputs
}

func (t BaseTask) Inputs() []TaskDependency {
	return t.inputs
}

func (t BaseTask) TaskTimeout() (time.Duration, bool) {
	if t.Timeout == nil {
		return time.Duration(0), false
	}
	return *t.Timeout, true
}

func (t BaseTask) TaskRetries() uint32 {
	return t.Retries.Uint32
}

func (t BaseTask) TaskMinBackoff() time.Duration {
	if t.MinBackoff > 0 {
		return t.MinBackoff
	}
	return time.Second * 5
}

func (t BaseTask) TaskMaxBackoff() time.Duration {
	if t.MinBackoff > 0 {
		return t.MaxBackoff
	}
	return time.Minute
}
