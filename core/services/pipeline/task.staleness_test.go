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
	fresh := now - 1000  // 1 second old
	stale := now - 30000 // 30 seconds old

	makeSamples := func() []pipeline.Result {
		return []pipeline.Result{
			{Value: pipeline.Sample{Source: "a", Value: decimal.NewFromInt(100), Weight: decimal.RequireFromString("0.5"), TsMs: fresh}},
			{Value: pipeline.Sample{Source: "b", Value: decimal.NewFromInt(101), Weight: decimal.RequireFromString("0.3"), TsMs: stale}},
			{Value: pipeline.Sample{Source: "c", Value: decimal.NewFromInt(102), Weight: decimal.RequireFromString("0.2"), TsMs: fresh}},
		}
	}

	t.Run("cutoff drops stale samples", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "cutoff",
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
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "linear",
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
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "exp",
			Threshold: "60s",
			HalfLife:  "5s",
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

	t.Run("exp_cooldown holds flat then decays", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "exp_cooldown",
			Threshold: "10s",
			HalfLife:  "5s",
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
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "piecewise",
			Threshold: "60s",
			Points:    "0s:1;5s:1;10s:0.5;30s:0",
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

	t.Run("decayThreshold drops samples below K", func(t *testing.T) {
		t.Parallel()
		now := time.Now().UnixMilli()
		// 50s old: exp_cooldown T=10s, H=5s -> decay = 2^(-(50-10)/5) = 2^-8 = 0.0039 < 0.03
		inputs := []pipeline.Result{
			{Value: pipeline.Sample{Source: "fresh", Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1), TsMs: now - 1000}},
			{Value: pipeline.Sample{Source: "stale", Value: decimal.NewFromInt(101), Weight: decimal.NewFromInt(1), TsMs: now - 50000}},
		}
		task := pipeline.StalenessTask{
			BaseTask:       pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:         "exp_cooldown",
			Threshold:      "10s",
			HalfLife:       "5s",
			DecayThreshold: "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1, "50s old sample with decay=0.0039 should be dropped by K=0.03")
		assert.Equal(t, "fresh", res[0].Source)
	})

	t.Run("decayThreshold keeps samples above K", func(t *testing.T) {
		t.Parallel()
		now := time.Now().UnixMilli()
		// 15s old: exp_cooldown T=10s, H=5s -> decay = 2^-((15-10)/5) = 2^-1 = 0.5 > 0.03
		inputs := []pipeline.Result{
			{Value: pipeline.Sample{Source: "a", Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1), TsMs: now - 15000}},
		}
		task := pipeline.StalenessTask{
			BaseTask:       pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:         "exp_cooldown",
			Threshold:      "10s",
			HalfLife:       "5s",
			DecayThreshold: "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1, "15s old sample with decay=0.5 should survive K=0.03")
	})

	t.Run("cutoff drops samples older than cutoff time", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "exp_cooldown",
			Threshold: "10s",
			HalfLife:  "5s",
			Cutoff:    "20s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2, "30s old sample b should be dropped by cutoff=20s")
		for _, s := range res {
			assert.NotEqual(t, "b", s.Source)
		}
	})

	t.Run("cutoff does not affect fresh samples", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "exp_cooldown",
			Threshold: "10s",
			HalfLife:  "5s",
			Cutoff:    "60s",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		// b is 30s old, cutoff 60s -> not cutoff-dropped, but decayed
		// 30s old: exp_cooldown T=10s, H=5s -> 2^-((30-10)/5) = 2^-4 = 0.0625 * 0.3 = 0.01875
		// All 3 samples survive (decay > 0)
		require.Len(t, res, 3)
	})
}
