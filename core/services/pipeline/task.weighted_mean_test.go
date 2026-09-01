package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestWeightedMeanTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputs   []pipeline.Sample
		minMass  string
		want     string
	}{
		{
			name: "equal values, equal weights",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 100, 1),
				sample("c", 100, 1),
				sample("d", 100, 1),
				sample("e", 100, 1),
			},
			want: "100",
		},
		{
			name: "weighted mean without clamping",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 110, 1),
				sample("c", 120, 3),
			},
			want: "114",
		},
		{
			name: "outlier is NOT clamped (unlike winsorizedmean)",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 100, 1),
				sample("c", 50000, 1),
			},
			want: "16733.33333333333333333",
		},
		{
			name: "collapse: two samples same source, weight counted once",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("a", 102, 1),
				sample("b", 200, 1),
			},
			want: "150.5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.WeightedMeanTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
			}
			if test.minMass != "" {
				task.MinWeightMass = test.minMass
			}
			inputs := make([]pipeline.Result, len(test.inputs))
			for i, v := range test.inputs {
				inputs[i] = pipeline.Result{Value: v}
			}
			out, runInfo := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, out.Error)
			got, ok := out.Value.(decimal.Decimal)
			require.True(t, ok)
			want, err := decimal.NewFromString(test.want)
			require.NoError(t, err)
			assert.True(t, got.GreaterThanOrEqual(want.Mul(decimal.NewFromFloat(0.9999))) &&
				got.LessThanOrEqual(want.Mul(decimal.NewFromFloat(1.0001))),
				"want %s, got %s", want.String(), got.String())
		})
	}
}

func TestWeightedMeanTaskGate(t *testing.T) {
	t.Parallel()

	inputs := []pipeline.Result{
		{Value: sample("a", 100, 0.1)},
		{Value: sample("b", 100, 0.1)},
		{Value: sample("c", 100, 0.1)},
		{Value: sample("d", 100, 0.1)},
	}

	t.Run("mass below minimum blocks", func(t *testing.T) {
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.5",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "insufficient weight mass")
	})

	t.Run("mass above minimum passes", func(t *testing.T) {
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.3",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
	})
}

func TestWeightedMeanTaskPrecision(t *testing.T) {
	t.Parallel()

	inputs := []pipeline.Result{
		{Value: sample("a", 100, 1)},
		{Value: sample("b", 101, 1)},
		{Value: sample("c", 102, 1)},
	}
	task := pipeline.WeightedMeanTask{
		BaseTask:  pipeline.NewBaseTask(0, "task", nil, nil, 0),
		Precision: "2",
	}
	out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
	require.NoError(t, out.Error)
	got := out.Value.(decimal.Decimal)
	assert.True(t, got.Equal(decimal.RequireFromString("101.00")), "got %s", got)
}

func TestWeightedMeanTaskInputErrors(t *testing.T) {
	t.Parallel()

	inputs := []pipeline.Result{
		{Value: sample("a", 100, 1)},
		{Error: errors.New("boom")},
	}
	task := pipeline.WeightedMeanTask{
		BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
	}
	out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
	require.NoError(t, out.Error, "errored inputs should be filtered, not fail the task")
	got := out.Value.(decimal.Decimal)
	assert.True(t, got.Equal(decimal.NewFromInt(100)), "got %s", got)
}

func TestWeightedMeanTaskNoSurvivingSources(t *testing.T) {
	t.Parallel()

	inputs := []pipeline.Result{
		{Value: sample("a", 100, 0)},
		{Value: sample("b", 100, 0)},
	}
	task := pipeline.WeightedMeanTask{
		BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
	}
	out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
	require.Error(t, out.Error)
	assert.Contains(t, out.Error.Error(), "no surviving sources")
}
