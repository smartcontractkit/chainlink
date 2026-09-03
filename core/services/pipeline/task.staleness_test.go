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

	// Helper running Run with a single sample at the given ageMs. Uses approximate
	// bounds because the sample's TsMs is captured here but Run calls time.Now()
	// itself, introducing sub-ms drift between capture and use.
	runAt := func(method string, ageMs int64, extra ...func(*pipeline.StalenessTask)) (decimal.Decimal, error) {
		t.Helper()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    method,
			Threshold: "10s",
			HalfLife:  "5s",
		}
		for _, f := range extra {
			f(&task)
		}
		sampleAt := time.Now().UnixMilli() - ageMs
		sample := pipeline.Sample{Source: "x", Value: decimal.NewFromInt(1), Weight: decimal.NewFromInt(1), TsMs: sampleAt}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: sample}})
		if out.Error != nil {
			return decimal.Zero, out.Error
		}
		res := out.Value.([]pipeline.Sample)
		if len(res) == 0 {
			return decimal.Zero, nil
		}
		return res[0].Weight, nil
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
		subNow := time.Now().UnixMilli()
		// 50s old: exp_cooldown T=10s, H=5s -> decay = 2^(-(50-10)/5) = 2^-8 = 0.0039 < 0.03
		inputs := []pipeline.Result{
			{Value: pipeline.Sample{Source: "fresh", Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1), TsMs: subNow - 1000}},
			{Value: pipeline.Sample{Source: "stale", Value: decimal.NewFromInt(101), Weight: decimal.NewFromInt(1), TsMs: subNow - 50000}},
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
		subNow := time.Now().UnixMilli()
		// 15s old: exp_cooldown T=10s, H=5s -> decay = 2^-((15-10)/5) = 2^-1 = 0.5 > 0.03
		inputs := []pipeline.Result{
			{Value: pipeline.Sample{Source: "a", Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1), TsMs: subNow - 15000}},
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

	t.Run("exp fails without halfLife", func(t *testing.T) {
		t.Parallel()
		_, err := runAt("exp", 1000, func(t *pipeline.StalenessTask) { t.HalfLife = "" })
		require.Error(t, err)
	})

	t.Run("exp_cooldown fails without halfLife", func(t *testing.T) {
		t.Parallel()
		_, err := runAt("exp_cooldown", 15000, func(t *pipeline.StalenessTask) { t.HalfLife = "" })
		require.Error(t, err)
	})

	t.Run("unknown method returns error", func(t *testing.T) {
		t.Parallel()
		_, err := runAt("mystery", 1000)
		require.Error(t, err)
	})

	t.Run("piecewise invalid points returns error", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "piecewise",
			Threshold: "60s",
			Points:    "5s,1",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.Error(t, out.Error)
	})

	t.Run("piecewise empty points returns error", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "piecewise",
			Threshold: "60s",
			Points:    "",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
		require.Error(t, out.Error)
	})

	t.Run("exp at exact threshold not zero", func(t *testing.T) {
		t.Parallel()
		// age=10s, halfLife=5s -> 2^-2 = 0.25 (approximate, sub-ms drift)
		w, err := runAt("exp", 10000)
		require.NoError(t, err)
		assert.True(t, w.GreaterThan(decimal.RequireFromString("0.24")) && w.LessThan(decimal.RequireFromString("0.26")), "got %s", w)
	})

	t.Run("exp at exact halfLife is 0.5", func(t *testing.T) {
		t.Parallel()
		// age=5s, halfLife=5s -> 2^-1 = 0.5
		w, err := runAt("exp", 5000, func(t *pipeline.StalenessTask) { t.Threshold = "60s" })
		require.NoError(t, err)
		assert.True(t, w.GreaterThan(decimal.RequireFromString("0.49")) && w.LessThan(decimal.RequireFromString("0.51")), "got %s", w)
	})

	t.Run("exp_cooldown at exact threshold is near 1", func(t *testing.T) {
		t.Parallel()
		w, err := runAt("exp_cooldown", 10000)
		require.NoError(t, err)
		assert.True(t, w.GreaterThan(decimal.RequireFromString("0.99")), "got %s", w)
	})

	t.Run("linear at age 0 is 1", func(t *testing.T) {
		t.Parallel()
		w, err := runAt("linear", 0)
		require.NoError(t, err)
		assert.True(t, w.GreaterThan(decimal.RequireFromString("0.99")), "got %s", w)
	})

	t.Run("linear at exact threshold is 0", func(t *testing.T) {
		t.Parallel()
		w, err := runAt("linear", 10000)
		require.NoError(t, err)
		assert.True(t, w.IsZero(), "got %s", w)
	})

	t.Run("piecewise interpolates between points", func(t *testing.T) {
		t.Parallel()
		// points: 0s:1; 10s:0.5; 20s:0 -> 5s old should be ~0.75
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "piecewise",
			Threshold: "60s",
			Points:    "0s:1;10s:0.5;20s:0",
		}
		sampleAt := time.Now().UnixMilli() - 5000
		sample := pipeline.Sample{Source: "x", Value: decimal.NewFromInt(1), Weight: decimal.NewFromInt(1), TsMs: sampleAt}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: sample}})
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Weight.GreaterThan(decimal.RequireFromString("0.74")) && res[0].Weight.LessThan(decimal.RequireFromString("0.76")), "got %s", res[0].Weight)
	})

	t.Run("piecewise clamps past last point", func(t *testing.T) {
		t.Parallel()
		// points: 0s:1; 5s:0.5 -> after 5s, value stays 0.5
		task := pipeline.StalenessTask{
			BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:    "piecewise",
			Threshold: "60s",
			Points:    "0s:1;5s:0.5",
		}
		sample := pipeline.Sample{Source: "x", Value: decimal.NewFromInt(1), Weight: decimal.NewFromInt(1), TsMs: now - 60000}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: sample}})
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1)
		assert.True(t, res[0].Weight.Equal(decimal.RequireFromString("0.5")), "got %s", res[0].Weight)
	})

	t.Run("decayThreshold ignored for linear", func(t *testing.T) {
		t.Parallel()
		// linear at 5s with 10s threshold -> 0.5, K=10 should NOT drop it
		task := pipeline.StalenessTask{
			BaseTask:       pipeline.NewBaseTask(0, "stale", nil, nil, 0),
			Method:         "linear",
			Threshold:      "10s",
			DecayThreshold: "10",
		}
		sample := pipeline.Sample{Source: "x", Value: decimal.NewFromInt(1), Weight: decimal.NewFromInt(1), TsMs: now - 5000}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: sample}})
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 1, "linear must not honor decayThreshold")
	})

	t.Run("cutoff zero is disabled", func(t *testing.T) {
		t.Parallel()
		// Use piecewise so the method alone wouldn't drop anything. cutoff=0 must not
		// drop any sample; cutoff=20s must drop the 30s-old `b`.
		mk := func(c string) []pipeline.Sample {
			task := pipeline.StalenessTask{
				BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
				Method:    "piecewise",
				Threshold: "60s",
				Points:    "0s:1;120s:1",
				Cutoff:    c,
			}
			out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), makeSamples())
			require.NoError(t, out.Error)
			return out.Value.([]pipeline.Sample)
		}
		require.Len(t, mk("0s"), 3, "cutoff=0 disabled")
		require.Len(t, mk("20s"), 2, "cutoff=20s drops the 30s-old sample")
	})
}
