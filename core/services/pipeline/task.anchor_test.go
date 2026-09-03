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

func TestAnchorTask(t *testing.T) {
	t.Parallel()

	// ref=100, lo=99, hi=101 for all sources.
	// disp = (101-99)/(2*100) = 0.01
	// lo_adj = q*(1-0.01) = 99,  hi_adj = q*(1+0.01) = 101
	// (No clamping: median=100, band=[97,103], all refs = 100.)
	ref := []pipeline.Sample{
		sample("a", 100, 1),
		sample("b", 100, 1),
		sample("c", 100, 1),
	}
	lo := []pipeline.Sample{
		sample("a", 99, 1),
		sample("b", 99, 1),
		sample("c", 99, 1),
	}
	hi := []pipeline.Sample{
		sample("a", 101, 1),
		sample("b", 101, 1),
		sample("c", 101, 1),
	}

	vars := pipeline.NewVarsFrom(map[string]any{
		"ref": ref,
		"lo":  lo,
		"hi":  hi,
	})

	t.Run("basic: lo_adj and hi_adj", func(t *testing.T) {
		t.Parallel()

		loTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_lo", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := loTask.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromInt(99)), "lo_adj should be 99, got %s", got)

		hiTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_hi", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "high",
			Band:      "0.03",
		}
		out, _ = hiTask.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		got = out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromInt(101)), "hi_adj should be 101, got %s", got)
	})

	t.Run("exact formula with asymmetric spread", func(t *testing.T) {
		t.Parallel()
		// ref=100, lo=98, hi=103 → disp=(103-98)/200=0.025
		// lo_adj = 100*(1-0.025) = 97.5
		// hi_adj = 100*(1+0.025) = 102.5
		refA := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, 1),
			sample("c", 100, 1),
		}
		loA := []pipeline.Sample{
			sample("a", 98, 1),
			sample("b", 98, 1),
			sample("c", 98, 1),
		}
		hiA := []pipeline.Sample{
			sample("a", 103, 1),
			sample("b", 103, 1),
			sample("c", 103, 1),
		}
		varsA := pipeline.NewVarsFrom(map[string]any{
			"ref": refA,
			"lo":  loA,
			"hi":  hiA,
		})

		loTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_lo", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := loTask.Run(t.Context(), logger.TestLogger(t), varsA, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromFloat(97.5)), "lo_adj should be 97.5, got %s", got)

		hiTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_hi", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "high",
			Band:      "0.03",
		}
		out, _ = hiTask.Run(t.Context(), logger.TestLogger(t), varsA, nil)
		require.NoError(t, out.Error)
		got = out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromFloat(102.5)), "hi_adj should be 102.5, got %s", got)
	})

	t.Run("clamping shifts q_i and scales adj", func(t *testing.T) {
		t.Parallel()
		// One ref outlier at 200 → median=100, band=[97,103], 200→103
		// Outlier: q=103, disp=(201-199)/400=0.005, lo_adj=103*(1-0.005)=102.485
		// Normal:  q=100, disp=0.01, lo_adj=99
		// Weighted mean = (99 + 99 + 102.485) / 3 = 100.1616...
		refC := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, 1),
			sample("c", 200, 1),
		}
		loC := []pipeline.Sample{
			sample("a", 99, 1),
			sample("b", 99, 1),
			sample("c", 199, 1),
		}
		hiC := []pipeline.Sample{
			sample("a", 101, 1),
			sample("b", 101, 1),
			sample("c", 201, 1),
		}
		varsC := pipeline.NewVarsFrom(map[string]any{
			"ref": refC,
			"lo":  loC,
			"hi":  hiC,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_lo", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsC, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		want, _ := decimal.NewFromString("100.16166666666666666")
		assert.True(t, got.GreaterThanOrEqual(want.Mul(decimal.NewFromFloat(0.9999))) &&
			got.LessThanOrEqual(want.Mul(decimal.NewFromFloat(1.0001))),
			"want %s, got %s", want, got)
	})

	t.Run("unmatched source in low dropped", func(t *testing.T) {
		t.Parallel()
		refU := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, 1),
			sample("c", 100, 1),
		}
		loU := []pipeline.Sample{
			sample("a", 99, 1),
			sample("b", 99, 1),
		}
		hiU := []pipeline.Sample{
			sample("a", 101, 1),
			sample("b", 101, 1),
			sample("c", 101, 1),
		}
		varsU := pipeline.NewVarsFrom(map[string]any{
			"ref": refU,
			"lo":  loU,
			"hi":  hiU,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsU, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromInt(99)), "got %s", got)
	})

	t.Run("unmatched source in high dropped", func(t *testing.T) {
		t.Parallel()
		refU := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, 1),
			sample("c", 100, 1),
		}
		loU := []pipeline.Sample{
			sample("a", 99, 1),
			sample("b", 99, 1),
			sample("c", 99, 1),
		}
		hiU := []pipeline.Sample{
			sample("a", 101, 1),
			sample("b", 101, 1),
		}
		varsU := pipeline.NewVarsFrom(map[string]any{
			"ref": refU,
			"lo":  loU,
			"hi":  hiU,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsU, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		// "c" dropped from high → only "a" and "b" survive
		assert.True(t, got.Equal(decimal.NewFromInt(99)), "got %s", got)
	})

	t.Run("gate failure on insufficient mass", func(t *testing.T) {
		t.Parallel()
		task := pipeline.AnchorTask{
			BaseTask:      pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference:     "$(ref)",
			Low:           "$(lo)",
			High:          "$(hi)",
			Select:        "low",
			Band:          "0.03",
			MinWeightMass: "10",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "insufficient weight mass")
	})

	t.Run("div-by-zero guard drops source with zero reference", func(t *testing.T) {
		t.Parallel()
		refZ := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, 1),
			sample("c", 0, 1),
		}
		loZ := []pipeline.Sample{
			sample("a", 99, 1),
			sample("b", 99, 1),
			sample("c", 0, 1),
		}
		hiZ := []pipeline.Sample{
			sample("a", 101, 1),
			sample("b", 101, 1),
			sample("c", 0, 1),
		}
		varsZ := pipeline.NewVarsFrom(map[string]any{
			"ref": refZ,
			"lo":  loZ,
			"hi":  hiZ,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsZ, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.NewFromInt(99)), "got %s", got)
	})

	t.Run("ordering invariant: lo_adj <= ref_agg <= hi_adj", func(t *testing.T) {
		t.Parallel()
		refO := []pipeline.Sample{
			sample("a", 100, 3),
			sample("b", 101, 2),
			sample("c", 50000, 1),
		}
		loO := []pipeline.Sample{
			sample("a", 99, 3),
			sample("b", 99, 2),
			sample("c", 49999, 1),
		}
		hiO := []pipeline.Sample{
			sample("a", 101, 3),
			sample("b", 105, 2),
			sample("c", 50001, 1),
		}
		varsO := pipeline.NewVarsFrom(map[string]any{
			"ref": refO,
			"lo":  loO,
			"hi":  hiO,
		})

		refTask := pipeline.WeightedMeanTask{
			BaseTask: pipeline.NewBaseTask(0, "wm_ref", nil, nil, 0),
			Samples:  "$(ref)",
			Method:   "winsor",
			Band:     "0.03",
		}
		refOut, _ := refTask.Run(t.Context(), logger.TestLogger(t), varsO, nil)
		require.NoError(t, refOut.Error)
		refVal := refOut.Value.(decimal.Decimal)

		loTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_lo", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		loOut, _ := loTask.Run(t.Context(), logger.TestLogger(t), varsO, nil)
		require.NoError(t, loOut.Error)
		loVal := loOut.Value.(decimal.Decimal)

		hiTask := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor_hi", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "high",
			Band:      "0.03",
		}
		hiOut, _ := hiTask.Run(t.Context(), logger.TestLogger(t), varsO, nil)
		require.NoError(t, hiOut.Error)
		hiVal := hiOut.Value.(decimal.Decimal)

		assert.True(t, loVal.LessThanOrEqual(refVal), "lo %s must be <= ref %s", loVal, refVal)
		assert.True(t, refVal.LessThanOrEqual(hiVal), "ref %s must be <= hi %s", refVal, hiVal)
	})

	t.Run("precision rounding", func(t *testing.T) {
		t.Parallel()
		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
			Precision: "2",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.NoError(t, out.Error)
		got := out.Value.(decimal.Decimal)
		assert.True(t, got.Equal(decimal.RequireFromString("99.00")), "got %s", got)
	})

	t.Run("rejects negative weight on reference", func(t *testing.T) {
		t.Parallel()
		refN := []pipeline.Sample{
			sample("a", 100, 1),
			sample("b", 100, -0.5),
		}
		varsN := pipeline.NewVarsFrom(map[string]any{
			"ref": refN,
			"lo":  lo,
			"hi":  hi,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsN, nil)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "negative weight")
	})

	t.Run("rejects negative band", func(t *testing.T) {
		t.Parallel()
		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "-0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), vars, nil)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "band must be non-negative")
	})

	t.Run("no surviving reference samples", func(t *testing.T) {
		t.Parallel()
		refE := []pipeline.Sample{
			sample("a", 100, 0),
			sample("b", 100, 0),
		}
		varsE := pipeline.NewVarsFrom(map[string]any{
			"ref": refE,
			"lo":  lo,
			"hi":  hi,
		})

		task := pipeline.AnchorTask{
			BaseTask:  pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
			Reference: "$(ref)",
			Low:       "$(lo)",
			High:      "$(hi)",
			Select:    "low",
			Band:      "0.03",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), varsE, nil)
		require.Error(t, out.Error)
		assert.Contains(t, out.Error.Error(), "no surviving reference samples")
	})

	t.Run("errored inputs filtered", func(t *testing.T) {
		t.Parallel()
		inputs := []pipeline.Result{
			{Value: sample("a", 100, 1)},
			{Error: errors.New("boom")},
		}
		task := pipeline.AnchorTask{
			BaseTask: pipeline.NewBaseTask(0, "anchor", nil, nil, 0),
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), inputs)
		_ = out
	})
}

// One source returns a wildly broken value while others cluster tightly.
// The anchor task must keep the rebased values within the winsor band of
// the good sources, not let the outlier drag the aggregate by its raw weight.
func TestAnchorTaskOutlierProtection(t *testing.T) {
	t.Parallel()

	ref := []pipeline.Sample{
		sample("a", 29136.5, 3),
		sample("b", 709.6, 3),
		sample("c", 709.55, 2),
		sample("d", 709.65, 1),
		sample("e", 709.6, 1),
	}
	lo := []pipeline.Sample{
		sample("a", 29136.0, 3),
		sample("b", 709.5, 3),
		sample("c", 709.5, 2),
		sample("d", 709.6, 1),
		sample("e", 709.5, 1),
	}
	hi := []pipeline.Sample{
		sample("a", 29137.0, 3),
		sample("b", 709.7, 3),
		sample("c", 709.6, 2),
		sample("d", 709.7, 1),
		sample("e", 709.7, 1),
	}

	vars := pipeline.NewVarsFrom(map[string]any{
		"ref": ref,
		"lo":  lo,
		"hi":  hi,
	})

	// Reference aggregate (winsorized weighted mean) — already protected
	refTask := pipeline.WeightedMeanTask{
		BaseTask: pipeline.NewBaseTask(0, "wm_ref", nil, nil, 0),
		Samples:  "$(ref)",
		Method:   "winsor",
		Band:     "0.03",
	}
	refOut, _ := refTask.Run(t.Context(), logger.TestLogger(t), vars, nil)
	require.NoError(t, refOut.Error)
	refVal := refOut.Value.(decimal.Decimal)

	refF, _ := refVal.Float64()
	assert.InDelta(t, 716.0, refF, 20.0, "ref should be near good sources, got %s", refVal)

	// lo via anchor
	loTask := pipeline.AnchorTask{
		BaseTask:  pipeline.NewBaseTask(0, "anchor_lo", nil, nil, 0),
		Reference: "$(ref)",
		Low:       "$(lo)",
		High:      "$(hi)",
		Select:    "low",
		Band:      "0.03",
	}
	loOut, _ := loTask.Run(t.Context(), logger.TestLogger(t), vars, nil)
	require.NoError(t, loOut.Error)
	loVal := loOut.Value.(decimal.Decimal)

	// hi via anchor
	hiTask := pipeline.AnchorTask{
		BaseTask:  pipeline.NewBaseTask(0, "anchor_hi", nil, nil, 0),
		Reference: "$(ref)",
		Low:       "$(lo)",
		High:      "$(hi)",
		Select:    "high",
		Band:      "0.03",
	}
	hiOut, _ := hiTask.Run(t.Context(), logger.TestLogger(t), vars, nil)
	require.NoError(t, hiOut.Error)
	hiVal := hiOut.Value.(decimal.Decimal)

	loF, _ := loVal.Float64()
	hiF, _ := hiVal.Float64()

	assert.InDelta(t, 716.0, loF, 20.0, "lo_adj should be near good sources, got %s", loVal)
	assert.InDelta(t, 716.0, hiF, 20.0, "hi_adj should be near good sources, got %s", hiVal)

	assert.True(t, loVal.LessThanOrEqual(refVal), "lo %s must be <= ref %s", loVal, refVal)
	assert.True(t, refVal.LessThanOrEqual(hiVal), "ref %s must be <= hi %s", refVal, hiVal)

	spread := hiVal.Sub(loVal)
	spreadF, _ := spread.Float64()
	assert.Less(t, spreadF, 1.0, "spread should be tight, got %s", spread)
}
