package pipeline_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
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

	t.Run("exp at exact threshold not zero", func(t *testing.T) {
		t.Parallel()
		// age=10s, halfLife=5s -> 2^-2 = 0.25. Wall-clock timing in `runAt`
		// drifts the observed age past 10s and conflates this assertion, so
		// the equivalent assertion against decayMultiplier lives in
		// task.staleness_internal_test.go.
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
		// Use cutoff method with a threshold larger than all sample ages so the
		// method alone wouldn't drop anything. cutoff=0 must not drop any sample;
		// cutoff=20s must drop the 30s-old `b`.
		mk := func(c string) []pipeline.Sample {
			task := pipeline.StalenessTask{
				BaseTask:  pipeline.NewBaseTask(0, "stale", nil, nil, 0),
				Method:    "cutoff",
				Threshold: "120s",
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

func TestStalenessTask_RelativeToFreshest(t *testing.T) {
	t.Parallel()

	// Fixed timestamps far in the past: in relativeTo=freshest mode the wall
	// clock must play no part.
	const freshestMs = int64(1725000000000)
	sample := func(source string, ageMs int64) pipeline.Sample {
		return pipeline.Sample{Source: source, Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1), TsMs: freshestMs - ageMs}
	}
	unknown := func(source string) pipeline.Sample {
		return pipeline.Sample{Source: source, Value: decimal.NewFromInt(100), Weight: decimal.NewFromInt(1)}
	}
	run := func(t *testing.T, task pipeline.StalenessTask, inputs ...pipeline.Result) pipeline.Result {
		t.Helper()
		task.BaseTask = pipeline.NewBaseTask(0, "stale", nil, nil, 0)
		if task.Method == "" {
			task.Method = "cutoff"
		}
		if task.Threshold == "" {
			task.Threshold = "10s"
		}
		task.RelativeTo = "freshest"
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		return out
	}
	sources := func(t *testing.T, out pipeline.Result) (kept []string, stale []string, errs int) {
		t.Helper()
		require.NoError(t, out.Error)
		for _, x := range out.Value.([]any) {
			switch v := x.(type) {
			case pipeline.Sample:
				kept = append(kept, v.Source)
			case pipeline.StaleSampleError:
				stale = append(stale, v.Source)
			case error:
				errs++
			default:
				t.Fatalf("unexpected element %T", x)
			}
		}
		return
	}

	t.Run("measures age against the freshest sample, not the wall clock", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("b", 5000)},
			pipeline.Result{Value: sample("c", 30000)},
		)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "a", res[0].Source)
		assert.Equal(t, "b", res[1].Source)
	})

	t.Run("unknown timestamps are never excluded and never the reference", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault"},
			pipeline.Result{Value: unknown("u")},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("c", 30000)},
		)
		kept, stale, _ := sources(t, out)
		assert.Equal(t, []string{"u", "a"}, kept)
		assert.Equal(t, []string{"c"}, stale)
	})

	t.Run("all unknown timestamps pass through", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{}, pipeline.Result{Value: unknown("u")}, pipeline.Result{Value: unknown("v")})
		require.NoError(t, out.Error)
		require.Len(t, out.Value.([]pipeline.Sample), 2)
	})

	t.Run("onStale=fault emits StaleSampleError and passes upstream errors through", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Error: errors.New("bridge down")},
			pipeline.Result{Value: sample("c", 30000)},
		)
		kept, stale, errs := sources(t, out)
		assert.Equal(t, []string{"a"}, kept)
		assert.Equal(t, []string{"c"}, stale)
		assert.Equal(t, 1, errs)

		staleErr := out.Value.([]any)[2].(pipeline.StaleSampleError)
		require.ErrorIs(t, staleErr, pipeline.ErrStaleSample)
		assert.Equal(t, freshestMs-30000, staleErr.TsMs)
		assert.Equal(t, freshestMs, staleErr.ReferenceTsMs)
		b, err := json.Marshal(staleErr)
		require.NoError(t, err)
		assert.JSONEq(t, `{"stale":true,"source":"c","ts_ms":1724999970000,"reference_ts_ms":1725000000000}`, string(b))
	})

	t.Run("onStale=drop discards upstream errors", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Error: errors.New("bridge down")},
		)
		require.NoError(t, out.Error)
		require.Len(t, out.Value.([]pipeline.Sample), 1)
	})

	t.Run("allowedFaults excludes stalest first", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault", AllowedFaults: "2"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("b", 20000)},
			pipeline.Result{Value: sample("c", 40000)},
			pipeline.Result{Value: sample("d", 30000)},
		)
		kept, stale, _ := sources(t, out)
		assert.Equal(t, []string{"a", "b"}, kept)
		assert.Equal(t, []string{"c", "d"}, stale)
	})

	t.Run("upstream errors spend the budget first", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault", AllowedFaults: "2"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Error: errors.New("bridge down")},
			pipeline.Result{Value: sample("b", 20000)},
			pipeline.Result{Value: sample("c", 40000)},
		)
		kept, stale, errs := sources(t, out)
		assert.Equal(t, []string{"a", "b"}, kept)
		assert.Equal(t, []string{"c"}, stale)
		assert.Equal(t, 1, errs)
	})

	t.Run("budget exhausted by upstream errors excludes nothing", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault", AllowedFaults: "1"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Error: errors.New("bridge down")},
			pipeline.Result{Error: errors.New("bridge down")},
			pipeline.Result{Value: sample("c", 40000)},
		)
		kept, stale, errs := sources(t, out)
		assert.Equal(t, []string{"a", "c"}, kept)
		assert.Empty(t, stale)
		assert.Equal(t, 2, errs)
	})

	t.Run("equally stale samples are excluded in input order", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "fault", AllowedFaults: "1"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("b", 30000)},
			pipeline.Result{Value: sample("c", 30000)},
		)
		kept, stale, _ := sources(t, out)
		assert.Equal(t, []string{"a", "c"}, kept)
		assert.Equal(t, []string{"b"}, stale)
	})

	t.Run("allowedFaults caps drop mode too, keeping decayed weight", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{Method: "linear", AllowedFaults: "1"},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("b", 15000)},
			pipeline.Result{Value: sample("c", 20000)},
		)
		require.NoError(t, out.Error)
		res := out.Value.([]pipeline.Sample)
		require.Len(t, res, 2)
		assert.Equal(t, "a", res[0].Source)
		assert.Equal(t, "b", res[1].Source)
		assert.True(t, res[1].Weight.IsZero(), "reinstated sample keeps its decayed weight, got %s", res[1].Weight)
	})

	t.Run("future timestamp becomes the reference", func(t *testing.T) {
		t.Parallel()
		// A source misreporting time is outside the gate's remit: pin that the
		// freshest timestamp wins even when it is implausible.
		out := run(t, pipeline.StalenessTask{OnStale: "fault"},
			pipeline.Result{Value: sample("future", -3600000)},
			pipeline.Result{Value: sample("a", 0)},
			pipeline.Result{Value: sample("b", 0)},
		)
		kept, stale, _ := sources(t, out)
		assert.Equal(t, []string{"future"}, kept)
		assert.Equal(t, []string{"a", "b"}, stale)
	})

	t.Run("invalid relativeTo errors", func(t *testing.T) {
		t.Parallel()
		task := pipeline.StalenessTask{BaseTask: pipeline.NewBaseTask(0, "stale", nil, nil, 0), Method: "cutoff", Threshold: "10s", RelativeTo: "sometime"}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		require.Error(t, out.Error)
	})

	t.Run("invalid onStale errors", func(t *testing.T) {
		t.Parallel()
		out := run(t, pipeline.StalenessTask{OnStale: "ignore"})
		require.Error(t, out.Error)
	})
}

// TestStalenessGate_ComposedPipeline runs the sample → staleness → median
// composition from ADR 0002 end to end.
func TestStalenessGate_ComposedPipeline(t *testing.T) {
	t.Parallel()

	const spec = `
s1     [type=sample input="$(r1)" source="a" valuePath="result"
        tsPath="timestamps,providerIndicatedTimeUnixMs"
        tsFallbackPath="timestamps,providerDataReceivedUnixMs" onMissingTs=unknown]
s2     [type=sample input="$(r2)" source="b" valuePath="result"
        tsPath="timestamps,providerIndicatedTimeUnixMs"
        tsFallbackPath="timestamps,providerDataReceivedUnixMs" onMissingTs=unknown]
s3     [type=sample input="$(r3)" source="c" valuePath="result"
        tsPath="timestamps,providerIndicatedTimeUnixMs"
        tsFallbackPath="timestamps,providerDataReceivedUnixMs" onMissingTs=unknown]
fresh  [type=staleness samples="[$(s1),$(s2),$(s3)]" method=cutoff threshold="30s"
        relativeTo=freshest onStale=fault allowedFaults=1]
answer [type=median values="$(fresh)" allowedFaults=1]
s1 -> fresh; s2 -> fresh; s3 -> fresh; fresh -> answer;
`
	const freshestMs = int64(1725000000000)
	body := func(result int, ageMs int64) string {
		return fmt.Sprintf(`{"result":%d,"timestamps":{"providerIndicatedTimeUnixMs":%d}}`, result, freshestMs-ageMs)
	}

	cfg := configtest.NewTestGeneralConfig(t)
	r := pipeline.NewRunner(nil, nil, cfg.JobPipeline(), cfg.WebServer(), nil, nil, nil, logger.TestLogger(t), nil, nil)

	tests := []struct {
		name       string
		r1, r2, r3 any
		want       string
	}{
		{"no stale sources", body(100, 0), body(101, 5000), body(200, 10000), "101"},
		{"stale source excluded", body(100, 0), body(101, 5000), body(200, 60000), "100.5"},
		{"only the stalest excluded when the threshold is reached", body(100, 0), body(101, 60000), body(200, 120000), "100.5"},
		{"received timestamp used when indicated is absent", body(100, 0), body(101, 5000),
			fmt.Sprintf(`{"result":200,"timestamps":{"providerDataReceivedUnixMs":%d}}`, freshestMs-60000), "100.5"},
		{"unknown timestamp is never excluded", body(100, 0), body(101, 60000), `{"result":200}`, "150"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vars := pipeline.NewVarsFrom(map[string]any{"r1": test.r1, "r2": test.r2, "r3": test.r3})
			_, trrs, err := r.ExecuteRun(t.Context(), pipeline.Spec{DotDagSource: spec}, vars)
			require.NoError(t, err)
			result, err := trrs.FinalResult().SingularResult()
			require.NoError(t, err)
			require.NoError(t, result.Error)
			assert.Equal(t, test.want, result.Value.(decimal.Decimal).String())
		})
	}

	t.Run("source fault spends the budget so stale sources are kept", func(t *testing.T) {
		t.Parallel()
		failingSpec := strings.Replace(spec, `s2     [type=sample input="$(r2)"`, `down [type=fail msg="bridge down"]
s2     [type=sample input="$(down)"`, 1) + "down -> s2;\n"
		vars := pipeline.NewVarsFrom(map[string]any{"r1": body(100, 0), "r3": body(200, 60000)})
		_, trrs, err := r.ExecuteRun(t.Context(), pipeline.Spec{DotDagSource: failingSpec}, vars)
		require.NoError(t, err)
		result, err := trrs.FinalResult().SingularResult()
		require.NoError(t, err)
		require.NoError(t, result.Error)
		assert.Equal(t, "150", result.Value.(decimal.Decimal).String())
	})

	t.Run("median fails when source faults alone exceed the threshold", func(t *testing.T) {
		t.Parallel()
		vars := pipeline.NewVarsFrom(map[string]any{"r1": body(100, 0), "r2": "not json", "r3": "not json"})
		_, trrs, err := r.ExecuteRun(t.Context(), pipeline.Spec{DotDagSource: spec}, vars)
		require.NoError(t, err)
		result, err := trrs.FinalResult().SingularResult()
		require.NoError(t, err)
		require.ErrorIs(t, result.Error, pipeline.ErrTooManyErrors)
	})
}
