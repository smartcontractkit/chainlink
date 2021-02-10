package pipeline

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

// TODO: This task type is no longer necessary, we should deprecate/remove it.
// See: https://www.pivotaltracker.com/story/show/176557536

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

func (t *ResultTask) SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error {
	return nil
}

func (t *ResultTask) Run(_ context.Context, taskRun TaskRun, inputs []Result) Result {
	values := make([]interface{}, len(inputs))
	errors := make(FinalErrors, len(inputs))
	for i, input := range inputs {
		values[i] = input.Value
		if input.Error != nil {
			errors[i] = null.StringFrom(input.Error.Error())
		}
	}
	return Result{Value: values, Error: errors}
}

// FIXME: This error/FinalErrors conflation exists solely because of the __result__ task.
// It is confusing and needs to go, making this note to remove it along with the
// special __result__ task.
// https://www.pivotaltracker.com/story/show/176557536
type FinalErrors []null.String

func (fe FinalErrors) HasErrors() bool {
	for _, err := range fe {
		if !err.IsZero() {
			return true
		}
	}
	return false
}

func (fe FinalErrors) Error() string {
	bs, err := json.Marshal(fe)
	if err != nil {
		return `["could not unmarshal final pipeline errors"]`
	}
	return string(bs)
}

func (fe FinalErrors) Value() (driver.Value, error) {
	return fe.Error(), nil
}

func (fe *FinalErrors) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, fe)
	case string:
		return json.Unmarshal([]byte(v), fe)
	default:
		return errors.New(fmt.Sprintf("%s", value))
	}
}
