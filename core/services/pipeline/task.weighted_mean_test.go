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

func sample(s string, value float64, weight float64) pipeline.Sample {
	return pipeline.Sample{
		Source: s,
		Value:  decimal.NewFromFloat(value),
		Weight: decimal.NewFromFloat(weight),
		TsMs:   1,
		Unit:   "USD",
	}
}

func TestWeightedMeanTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		inputs  []pipeline.Sample
		band    string
		minMass string
		want    string
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
			name: "weights shift clamped mean",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 110, 1),
				sample("c", 120, 3),
			},
			// median=110, band=[106.7, 113.3] (default 0.03)
			// 100 -> 106.7, 110 -> 110, 120 -> 113.3
			// weighted mean = (1*106.7 + 1*110 + 3*113.3) / 5 = 111.32
			want: "111.32",
		},
		{
			name: "outlier clamped to band, no weights",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 100, 1),
				sample("c", 50000, 1),
			},
			band: "0.03",
			// median = 100, band = [97, 103]; 50000 -> 103
			// weighted mean = (100 + 100 + 103) / 3 = 101
			want: "101",
		},
		{
			name: "outlier clamped but weighted low",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 101, 1),
				sample("c", 50000, 1),
			},
			band: "0.03",
			// median = 101, band = [97.97, 104.03]; 50000 -> 104.03
			// mean = (100 + 101 + 104.03) / 3 = 101.6766...
			want: "101.67666666666666666",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			task := pipeline.WeightedMeanTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Method:   "winsor",
				Band:     test.band,
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

func TestWeightedMeanTaskBoundedOutput(t *testing.T) {
	t.Parallel()
	// 100, 100, 100, 100, 10000 (corrupt source); median=100, band=[97,103]
	// Output must be in [97, 103] regardless of weights.
	tests := []struct {
		outlier float64
		weights []float64
	}{
		{10000, nil},
		{-9999, nil},
		{10000, []float64{1, 1, 1, 1, 1}},
		{10000, []float64{1, 1, 1, 1, 100}},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			inputs := []pipeline.Result{
				{Value: sample("a", 100, 1)},
				{Value: sample("b", 100, 1)},
				{Value: sample("c", 100, 1)},
				{Value: sample("d", 100, 1)},
				{Value: sample("e", tc.outlier, 1)},
			}
			if tc.weights != nil {
				for i, w := range tc.weights {
					s := inputs[i].Value.(pipeline.Sample)
					s.Weight = decimal.NewFromFloat(w)
					inputs[i].Value = s
				}
			}

			task := pipeline.WeightedMeanTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Method:   "winsor",
			}
			out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
			require.NoError(t, out.Error)
			got := out.Value.(decimal.Decimal)

			m := decimal.NewFromInt(100)
			lo := m.Mul(decimal.NewFromInt(1).Sub(decimal.NewFromFloat(0.03)))
			hi := m.Mul(decimal.NewFromInt(1).Add(decimal.NewFromFloat(0.03)))
			assert.True(t, got.GreaterThanOrEqual(lo), "got %s < lo %s", got, lo)
			assert.True(t, got.LessThanOrEqual(hi), "got %s > hi %s", got, hi)
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
		t.Parallel()
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.5",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "insufficient weight mass")
	})

	t.Run("mass above minimum passes", func(t *testing.T) {
		t.Parallel()
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.3",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
	})
}

func TestWeightedMeanTaskReference(t *testing.T) {
	t.Parallel()

	mid := []pipeline.Sample{
		sample("a", 100, 1),
		sample("b", 100, 1),
		sample("c", 101, 1),
	}
	bid := []pipeline.Sample{
		sample("a", 99, 1),
		sample("b", 99, 1),
		sample("c", 99, 1),
	}
	ask := []pipeline.Sample{
		sample("a", 101, 1),
		sample("b", 101, 1),
		sample("c", 105, 1),
	}

	vars := pipeline.NewVarsFrom(map[string]any{
		"mid": mid,
		"bid": bid,
		"ask": ask,
	})

	runOne := func(name, valExpr, refExpr string) decimal.Decimal {
		task := pipeline.WeightedMeanTask{
			BaseTask: pipeline.NewBaseTask(0, name, nil, nil, 0),
			Samples:  valExpr,
			Method:   "winsor",
		}
		if refExpr != "" {
			task.Reference = refExpr
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error, "%s failed: %v", name, out.Error)
		return out.Value.(decimal.Decimal)
	}

	midVal := runOne("mid", "$(mid)", "")
	bidVal := runOne("bid", "$(bid)", "$(mid)")
	askVal := runOne("ask", "$(ask)", "$(mid)")

	assert.True(t, bidVal.LessThanOrEqual(midVal), "bid %s must be <= mid %s", bidVal, midVal)
	assert.True(t, midVal.LessThanOrEqual(askVal), "mid %s must be <= ask %s", midVal, askVal)
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

func TestWeightedMeanTaskDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		inputs  []pipeline.Sample
		minMass string
		want    string
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
			name: "outlier is NOT clamped (default, no method)",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("b", 100, 1),
				sample("c", 50000, 1),
			},
			want: "16733.33333333333333333",
		},
		{
			name: "two samples same source, both counted independently",
			inputs: []pipeline.Sample{
				sample("a", 100, 1),
				sample("a", 102, 1),
				sample("b", 200, 1),
			},
			// No collapse: (1*100 + 1*102 + 1*200) / 3 = 134
			want: "134",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
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

func TestWeightedMeanTaskDefaultGate(t *testing.T) {
	t.Parallel()

	inputs := []pipeline.Result{
		{Value: sample("a", 100, 0.1)},
		{Value: sample("b", 100, 0.1)},
		{Value: sample("c", 100, 0.1)},
		{Value: sample("d", 100, 0.1)},
	}

	t.Run("mass below minimum blocks", func(t *testing.T) {
		t.Parallel()
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.5",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "insufficient weight mass")
	})

	t.Run("mass above minimum passes", func(t *testing.T) {
		t.Parallel()
		task := pipeline.WeightedMeanTask{
			BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
			MinWeightMass: "0.3",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
	})
}

func TestWeightedMeanTaskDefaultPrecision(t *testing.T) {
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

func TestWeightedMeanTaskDefaultInputErrors(t *testing.T) {
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

func TestWeightedMeanTaskDefaultNoSurvivingSamples(t *testing.T) {
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
	assert.Contains(t, out.Error.Error(), "no surviving samples")
}
