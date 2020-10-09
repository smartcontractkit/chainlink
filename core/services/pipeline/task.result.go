package pipeline

import (
	"go.uber.org/multierr"
)

// ResultTask exists solely as a Postgres performance optimization.  It's added
// automatically to the end of every pipeline, and it receives the outputs of all
// tasks that have no successor tasks.  This allows the pipeline runner to detect
// when it has reached the end a given pipeline simply by checking the `successor_id`
// field, rather than having to try to SELECT all of the pipeline run's task runs,
// (which must be done from inside of a transaction, and causes lock contention
// and serialization anomaly issues).
type ResultTask struct {
	BaseTask `mapstructure:",squash"`
}

var _ Task = (*ResultTask)(nil)

func (t *ResultTask) Type() TaskType {
	return TaskTypeResult
}

func (t *ResultTask) Run(taskRun TaskRun, inputs []Result) (result Result) {
	var values []interface{}
	for _, input := range inputs {
		if input.Error != nil {
			result.Error = multierr.Append(result.Error, input.Error)
		} else {
			values = append(values, input.Value)
		}
	}
	result.Value = values
	return result
}
