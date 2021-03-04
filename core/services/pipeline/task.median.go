package pipeline

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type MedianTask struct {
	BaseTask      `mapstructure:",squash"`
	AllowedFaults uint64 `json:"allowedFaults"`
}

var _ Task = (*MedianTask)(nil)

func (t *MedianTask) Type() TaskType {
	return TaskTypeMedian
}

func (t *MedianTask) SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error {
	if _, exists := inputValues["allowedFaults"]; !exists {
		if len(self.inputs()) == 0 {
			return errors.Wrapf(ErrWrongInputCardinality, "MedianTask requires at least 1 input")
		}
		t.AllowedFaults = uint64(len(self.inputs()) - 1)
	}
	return nil
}

func (t *MedianTask) Run(_ context.Context, taskRun TaskRun, inputs []Result) (result Result) {
	if len(inputs) == 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "MedianTask requires at least 1 input")}
	}

	answers := []decimal.Decimal{}
	fetchErrors := []error{}

	for _, input := range inputs {
		if input.Error != nil {
			fetchErrors = append(fetchErrors, input.Error)
			continue
		}

		answer, err := utils.ToDecimal(input.Value)
		if err != nil {
			fetchErrors = append(fetchErrors, err)
			continue
		}

		answers = append(answers, answer)
	}

	if uint64(len(fetchErrors)) > t.AllowedFaults {
		return Result{Error: errors.Wrapf(ErrBadInput, "Number of faulty inputs %v to median task > number allowed faults %v. Fetch errors: %v", len(fetchErrors), t.AllowedFaults, multierr.Combine(fetchErrors...).Error())}
	}

	sort.Slice(answers, func(i, j int) bool {
		return answers[i].LessThan(answers[j])
	})
	k := len(answers) / 2
	if len(answers)%2 == 1 {
		return Result{Value: answers[k]}
	}
	median := answers[k].Add(answers[k-1]).Div(decimal.NewFromInt(2))
	return Result{Value: median}
}
