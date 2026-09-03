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

	// Generic samples: two sources in different units, one already in target unit.
	makeSamples := func() []pipeline.Sample {
		return []pipeline.Sample{
			{Source: "srcA", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.RequireFromString("0.5"), TsMs: 1},
			{Source: "srcB", Value: decimal.NewFromInt(200), Unit: "T", Weight: decimal.RequireFromString("0.5"), TsMs: 1},
		}
	}

	t.Run("converts non-target unit by multiplying factor", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"factor":  pipeline.Sample{Source: "U1", Value: decimal.RequireFromString("3"), Unit: "T", Weight: decimal.NewFromInt(1), TsMs: 1},
			"samples": makeSamples(),
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":$(factor)}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "T", res[0].Unit)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(300)), "100 * 3 = 300")
		assert.True(t, res[1].Value.Equal(decimal.NewFromInt(200)), "already in target unit, unchanged")
	})

	t.Run("drops unconvertible sample when onMissingRate=drop", func(t *testing.T) {
		t.Parallel()
		task := pipeline.NormalizeTask{
			BaseTask:      pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit:    "T",
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
		assert.Equal(t, "srcB", res[0].Source)
	})

	t.Run("errors on missing factor when onMissingRate=error", func(t *testing.T) {
		t.Parallel()
		task := pipeline.NormalizeTask{
			BaseTask:      pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit:    "T",
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
		t.Parallel()
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			TargetUnit: "T",
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
		assert.Equal(t, "U1", res[0].Unit)
	})

	t.Run("sample with empty unit passes through unchanged", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"factor": pipeline.Sample{Source: "U1", Value: decimal.NewFromInt(3), Weight: decimal.NewFromInt(1), TsMs: 1},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    `$(noUnit)`,
			UnitMap:    `{"U1":$(factor)}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		require.NoError(t, vars.Set("noUnit", []pipeline.Sample{
			{Source: "srcC", Value: decimal.NewFromInt(42), Unit: "", Weight: decimal.NewFromInt(1), TsMs: 1},
		}))
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(42)))
	})
}

func TestNormalizeTaskUnitMap(t *testing.T) {
	t.Parallel()

	unitMapSamples := func() []pipeline.Sample {
		return []pipeline.Sample{
			{Source: "srcA", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(3), TsMs: 1},
			{Source: "srcB", Value: decimal.NewFromInt(200), Unit: "T", Weight: decimal.NewFromInt(1), TsMs: 1},
		}
	}

	t.Run("unitMap with var-ref to Sample", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"u1Factor": pipeline.Sample{Source: "U1", Value: decimal.RequireFromString("2.5"), Weight: decimal.NewFromInt(1), TsMs: 1},
			"samples":  unitMapSamples(),
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":$(u1Factor)}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "T", res[0].Unit)
		assert.True(t, res[0].Value.Equal(decimal.RequireFromString("250")), "100 * 2.5 = 250")
		assert.True(t, res[1].Value.Equal(decimal.NewFromInt(200)))
	})

	t.Run("unitMap with literal number", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": unitMapSamples(),
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":2.5}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.True(t, res[0].Value.Equal(decimal.RequireFromString("250")))
	})

	t.Run("unitMap missing unit drops sample", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "srcA", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
				{Source: "srcB", Value: decimal.NewFromInt(200), Unit: "U2", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:      pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:       "$(samples)",
			UnitMap:       `{"U1":2.5}`,
			TargetUnit:    "T",
			Enabled:       "true",
			OnMissingRate: "drop",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.Equal(t, "srcA", res[0].Source)
	})

	t.Run("unitMap with multiple units", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(10), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
				{Source: "b", Value: decimal.NewFromInt(20), Unit: "U2", Weight: decimal.NewFromInt(1), TsMs: 1},
				{Source: "c", Value: decimal.NewFromInt(30), Unit: "T", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":2,"U2":3}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 3)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(20)), "10 * 2")
		assert.True(t, res[1].Value.Equal(decimal.NewFromInt(60)), "20 * 3")
		assert.True(t, res[2].Value.Equal(decimal.NewFromInt(30)), "already in target unit, unchanged")
	})

	t.Run("all samples already in target unit", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(1), Unit: "T", Weight: decimal.NewFromInt(1), TsMs: 1},
				{Source: "b", Value: decimal.NewFromInt(2), Unit: "T", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":2.5}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(1)))
		assert.True(t, res[1].Value.Equal(decimal.NewFromInt(2)))
	})

	t.Run("unitMap zero factor multiplies value to 0", func(t *testing.T) {
		// Documents the deliberate new behavior: a zero factor is applied, not
		// treated as missing. Old factors path skipped IsZero; unitMap does not.
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":0}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.IsZero(), "100 * 0 = 0")
	})

	t.Run("unitMap from JSON string", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
			"unitMap": `{"U1":2}`,
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    "$(unitMap)",
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(200)))
	})

	t.Run("unitMap with bad value errors", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":"not-a-number"}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.Error(t, out.Error)
	})

	t.Run("fractional factor scales down", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":0.5}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.Equal(decimal.RequireFromString("50")), "100 * 0.5 = 50")
	})

	t.Run("negative factor flips sign", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(100), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":-1}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(-100)), "100 * -1 = -100")
	})

	t.Run("negative value scales through", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{
			"samples": []pipeline.Sample{
				{Source: "a", Value: decimal.NewFromInt(-50), Unit: "U1", Weight: decimal.NewFromInt(1), TsMs: 1},
			},
		})
		task := pipeline.NormalizeTask{
			BaseTask:   pipeline.NewBaseTask(0, "norm", nil, nil, 0),
			Samples:    "$(samples)",
			UnitMap:    `{"U1":2}`,
			TargetUnit: "T",
			Enabled:    "true",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Value.Equal(decimal.NewFromInt(-100)), "-50 * 2 = -100")
	})
}
