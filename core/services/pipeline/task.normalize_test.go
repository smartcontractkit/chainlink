package pipeline_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestNormalizeTask(t *testing.T) {
	t.Parallel()

	makeSamples := func() []pipeline.Sample {
		return []pipeline.Sample{
			{Source: "binance", Value: decimal.NewFromInt(399), Unit: "USDT", Weight: decimal.RequireFromString("0.333"), TsMs: 1},
			{Source: "lighter", Value: decimal.NewFromInt(412), Unit: "USD", Weight: decimal.RequireFromString("0.444"), TsMs: 1},
		}
	}

	t.Run("converts non-target unit with rate", func(t *testing.T) {
		vars := pipeline.NewVarsFrom(map[string]any{
			"rate":   pipeline.Sample{Source: "USDT", Value: decimal.RequireFromString("0.9995"), Unit: "USD", Weight: decimal.NewFromInt(1), TsMs: 1},
			"samples": makeSamples(),
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			Rates:      "$(rate)",
			TargetUnit: "USD",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "USD", res[0].Unit)
		assert.True(t, res[0].Value.Equal(decimal.RequireFromString("398.8005")))
		assert.True(t, res[1].Value.Equal(decimal.NewFromInt(412)))
	})

	t.Run("drops unconvertible sample when onMissingRate=drop", func(t *testing.T) {
		task := pipeline.NormalizeTask{
			BaseTask:      pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit:    "USD",
			Enabled:       "true",
			OnMissingRate: "drop",
		}
		inputs := make([]pipeline.Result, len(makeSamples()))
		for i, s := range makeSamples() {
			inputs[i] = pipeline.Result{Value: s}
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.Equal(t, "lighter", res[0].Source)
	})

	t.Run("errors on missing rate when onMissingRate=error", func(t *testing.T) {
		task := pipeline.NormalizeTask{
			BaseTask:      pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit:    "USD",
			Enabled:       "true",
			OnMissingRate: "error",
		}
		inputs := make([]pipeline.Result, len(makeSamples()))
		for i, s := range makeSamples() {
			inputs[i] = pipeline.Result{Value: s}
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.Error(t, out.Error)
	})

	t.Run("disabled returns samples unchanged", func(t *testing.T) {
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit: "USD",
			Enabled:    "false",
		}
		inputs := make([]pipeline.Result, len(makeSamples()))
		for i, s := range makeSamples() {
			inputs[i] = pipeline.Result{Value: s}
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "USDT", res[0].Unit)
	})
}
