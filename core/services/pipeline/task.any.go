package pipeline

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
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

func (t *AnyTask) Run(_ context.Context, _ logger.Logger, _ Vars, inputs []Result) (result Result, runInfo RunInfo) {
	if len(inputs) == 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "AnyTask requires at least 1 input")}, runInfo
	}

	var answers []interface{}

	for _, input := range inputs {
		if input.Error != nil {
			continue
		}

		answers = append(answers, input.Value)
	}

	if len(answers) == 0 {
		return Result{Error: errors.Wrapf(ErrBadInput, "There were zero non-errored inputs")}, runInfo
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(answers))))
	if err != nil {
		return Result{Error: errors.Wrapf(err, "Failed to generate random number for picking input")}, retryableRunInfo()
	}
	i := int(nBig.Int64())
	answer := answers[i]

	return Result{Value: answer}, runInfo
}
