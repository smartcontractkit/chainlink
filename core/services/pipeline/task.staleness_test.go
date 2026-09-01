package pipeline_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestStalenessTask(t *testing.T) {
	t.Parallel()

	now := time.Now().UnixMilli()
	fresh := now - 1000   // 1 second old
	stale := now - 30000  // 30 seconds old

	makeSamples := func() []pipeline.Result {
		return []pipeline.Result{
			{Value: pipeline.Sample{Source: "a", Value: decimal.NewFromInt(100), Weight: decimal.RequireFromString("0.5"), TsMs: fresh}},
			{Value: pipeline.Sample{Source: "b", Value: decimal.NewFromInt(101), Weight: decimal.RequireFromString("0.3"), TsMs: stale}},
			{Value: pipeline.Sample{Source: "c", Value: decimal.NewFromInt(102), Weight: decimal.RequireFromString("0.2"), TsMs: fresh}},
		}
	}

	t.Run("cutoff drops stale samples", func(t *testing.T) {
		task := pipeline.StalenessTask{
			BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:   "cutoff",
			Threshold: "10s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		for _, s := range res {
			assert.NotEqual(t, "b", s.Source)
		}
	})

	t.Run("linear scales weight continuously", func(t *testing.T) {
		task := pipeline.StalenessTask{
			BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:   "linear",
			Threshold: "10s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		for _, s := range res {
			if s.Source == "b" {
				// 30s / 10s => multiplier = 0, so dropped
				t.Fatalf("expected b to be dropped")
			}
			if s.Source == "a" {
				// 1s old, multiplier is close to 0.9
				assert.True(t, s.Weight.GreaterThan(decimal.RequireFromString("0.44")) && s.Weight.LessThan(decimal.RequireFromString("0.46")), "got %s", s.Weight)
			}
			if s.Source == "c" {
				assert.True(t, s.Weight.GreaterThan(decimal.RequireFromString("0.17")) && s.Weight.LessThan(decimal.RequireFromString("0.19")), "got %s", s.Weight)
			}
		}
	})

	t.Run("exponential decay", func(t *testing.T) {
		task := pipeline.StalenessTask{
			BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:   "exp",
			Threshold: "60s",
			HalfLife: "5s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 3)
		for _, s := range res {
			if s.Source == "b" {
				// 30s old with half-life 5s -> 2^-6 = 0.015625
				assert.True(t, s.Weight.LessThan(decimal.RequireFromString("0.02")))
			}
		}
	})

	t.Run("cooldown holds flat then decays", func(t *testing.T) {
		task := pipeline.StalenessTask{
			BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:   "cooldown",
			Threshold: "10s",
			HalfLife: "5s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		for _, s := range res {
			if s.Source == "a" || s.Source == "c" {
				assert.True(t, s.Weight.Equal(decimal.RequireFromString("0.5")) || s.Weight.Equal(decimal.RequireFromString("0.2")))
			}
			if s.Source == "b" {
				// 30s old, threshold 10s, extra 20s at half-life 5s -> 2^-4 = 0.0625 * 0.3 = 0.01875
				assert.True(t, s.Weight.LessThan(decimal.RequireFromString("0.02")), "got %s", s.Weight)
			}
		}
	})

	t.Run("piecewise interpolation", func(t *testing.T) {
		task := pipeline.StalenessTask{
			BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:   "piecewise",
			Threshold: "60s",
			Points:   "0s:1;5s:1;10s:0.5;30s:0",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		for _, s := range res {
			if s.Source == "b" {
				assert.True(t, s.Weight.IsZero())
			}
		}
	})
}
