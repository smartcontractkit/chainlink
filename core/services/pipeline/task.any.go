package pipeline

import (
	"context"
	"math/big"

	"crypto/rand"

	"github.com/pkg/errors"
)

// AnyTask picks a value at random from the set of non-errored inputs.
// If there are zero non-errored inputs then it returns an error.
type AnyTask struct {
	BaseTask `mapstructure:",squash"`
}

var _ Task = (*AnyTask)(nil)

func (t *AnyTask) Type() TaskType {
	return TaskTypeAny
}

func (t *AnyTask) SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error {
	return nil
}

func (t *AnyTask) Run(_ context.Context, taskRun TaskRun, inputs []Result) (result Result) {
	if len(inputs) == 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "AnyTask requires at least 1 input")}
	}

	var answers []interface{}

	for _, input := range inputs {
		if input.Error != nil {
			continue
		}

		answers = append(answers, input.Value)
	}

	if len(answers) == 0 {
		return Result{Error: errors.Wrapf(ErrBadInput, "There were zero non-errored inputs")}
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(answers))))
	if err != nil {
		return Result{Error: errors.Wrapf(err, "Failed to generate random number for picking input")}
	}
	i := int(nBig.Int64())
	answer := answers[i]

	return Result{Value: answer}
}
