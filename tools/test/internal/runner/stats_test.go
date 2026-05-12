package runner

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

func TestWilsonScoreInterval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		k, n      int
		z         float64
		wantLower float64
		wantUpper float64
	}{
		{
			name:      "zero_runs",
			k:         0,
			n:         0,
			z:         1.96,
			wantLower: 0,
			wantUpper: 0,
		},
		{
			name:      "all_fail",
			k:         10,
			n:         10,
			z:         1.96,
			wantLower: 0.722, // ~72.2%
			wantUpper: 1.0,
		},
		{
			name:      "none_fail",
			k:         0,
			n:         10,
			z:         1.96,
			wantLower: 0.0,
			wantUpper: 0.278, // ~27.8%
		},
		{
			name:      "half_fail",
			k:         5,
			n:         10,
			z:         1.96,
			wantLower: 0.237, // ~23.7%
			wantUpper: 0.763, // ~76.3%
		},
		{
			name:      "single_fail_large_n",
			k:         1,
			n:         100,
			z:         1.96,
			wantLower: 0.002,
			wantUpper: 0.054,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			lower, upper := WilsonScoreInterval(tc.k, tc.n, tc.z)
			assert.InDelta(t, tc.wantLower, lower, 0.001, "lower bound")
			assert.InDelta(t, tc.wantUpper, upper, 0.001, "upper bound")
			assert.True(t, lower <= upper, "lower must be <= upper")
			assert.True(t, lower >= 0, "lower must be >= 0")
			assert.True(t, upper <= 1, "upper must be <= 1")
			if tc.n > 0 {
				p := float64(tc.k) / float64(tc.n)
				assert.True(t, lower <= p+1e-9, "observed rate must be within CI")
				assert.True(t, upper >= p-1e-9, "observed rate must be within CI")
			}
		})
	}
}

func TestWilsonScoreInterval_NotNaN(t *testing.T) {
	t.Parallel()
	lower, upper := WilsonScoreInterval(0, 1, 1.96)
	assert.False(t, math.IsNaN(lower))
	assert.False(t, math.IsNaN(upper))
}

func TestFormatFlakyTestLine_includesCI(t *testing.T) {
	t.Parallel()
	e := TestEntry{
		Test:      "TestFoo",
		Package:   "pkg/bar",
		Runs:      10,
		Fails:     3,
		Successes: 7,
	}
	line := formatFlakyTestLine(e)
	assert.Contains(t, line, "30.0%", "observed rate")
	assert.Contains(t, line, "[Confidence Interval:", "CI annotation")
}

func TestFormatFlakyTestLine_packageLevel_includesCI(t *testing.T) {
	t.Parallel()
	e := TestEntry{
		Package:   "pkg/bar",
		Runs:      10,
		Fails:     3,
		Successes: 7,
	}
	line := formatFlakyTestLine(e)
	assert.Contains(t, line, "30.0%", "observed rate")
	assert.Contains(t, line, "[Confidence Interval:", "CI annotation")
}

func TestReportSummary_hasCI(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)
	require.NotNil(t, rep.Summary)
	s := rep.Summary

	require.NotNil(t, s.FlakeExecutionFailRateLower, "FlakeExecutionFailRateLower should be set")
	require.NotNil(t, s.FlakeExecutionFailRateUpper, "FlakeExecutionFailRateUpper should be set")
	assert.True(t, *s.FlakeExecutionFailRateLower >= 0)
	assert.True(t, *s.FlakeExecutionFailRateUpper <= 1)
	assert.True(t, *s.FlakeExecutionFailRateLower <= *s.FlakeExecutionFailRateUpper)

	require.NotNil(t, s.FlakeIterationFailRateLower, "FlakeIterationFailRateLower should be set")
	require.NotNil(t, s.FlakeIterationFailRateUpper, "FlakeIterationFailRateUpper should be set")
	assert.True(t, *s.FlakeIterationFailRateLower >= 0)
	assert.True(t, *s.FlakeIterationFailRateUpper <= 1)
	assert.True(t, *s.FlakeIterationFailRateLower <= *s.FlakeIterationFailRateUpper)
}

func TestPrintOverallStats_includesCI(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)

	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()

	assert.Contains(t, out, "[Confidence Interval:", "overall stats should show CI")
}

func TestReportSummary_hasCI_noFlakes(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)
	require.NotNil(t, rep.Summary)
	s := rep.Summary

	require.NotNil(t, s.FlakeIterationFailRateLower, "FlakeIterationFailRateLower should be set even with no flakes")
	require.NotNil(t, s.FlakeIterationFailRateUpper, "FlakeIterationFailRateUpper should be set even with no flakes")
	assert.Equal(t, 0.0, *s.FlakeIterationFailRateLower)
	assert.True(t, *s.FlakeIterationFailRateUpper > 0, "upper bound should be non-zero for finite n")

	require.NotNil(t, s.FlakePrevalenceLower, "FlakePrevalenceLower should be set even with no flakes")
	require.NotNil(t, s.FlakePrevalenceUpper, "FlakePrevalenceUpper should be set even with no flakes")
	assert.Equal(t, 0.0, *s.FlakePrevalenceLower)
	assert.True(t, *s.FlakePrevalenceUpper > 0)
}

func TestPrintOverallStats_includesCI_noFlakes(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)

	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()

	assert.Contains(t, out, "Flaky:", "should always show flaky line")
	assert.Contains(t, out, "Flaky Iterations:", "should always show flaky iterations line when iterations > 0")
	assert.Contains(t, out, "[Confidence Interval:", "CI annotation should appear even with 0 flakes")
}

func TestReportSummary_hasFlakePrevalenceCI(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)
	require.NotNil(t, rep.Summary)
	s := rep.Summary

	require.NotNil(t, s.FlakePrevalenceLower, "FlakePrevalenceLower should be set")
	require.NotNil(t, s.FlakePrevalenceUpper, "FlakePrevalenceUpper should be set")
	assert.True(t, *s.FlakePrevalenceLower >= 0)
	assert.True(t, *s.FlakePrevalenceUpper <= 1)
	assert.True(t, *s.FlakePrevalenceLower <= *s.FlakePrevalenceUpper)
}

func TestPrintOverallStats_flakyTestsLine_includesCI(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)

	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()

	// The "Flaky:" line (prevalence) should have a CI, not just the iterations line.
	lines := strings.Split(out, "\n")
	var flakyLine string
	for _, l := range lines {
		if strings.Contains(l, "Flaky:") && !strings.Contains(l, "Iterations") {
			flakyLine = l
			break
		}
	}
	require.NotEmpty(t, flakyLine, "Flaky: line must exist")
	assert.Contains(t, flakyLine, "[Confidence Interval:", "Flaky: line should have CI")
}

func TestCIStyleForGap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		gap       float64
		wantStyle string // rendered sentinel to confirm correct style applied
	}{
		// ≤0.10 → green (OK)
		{name: "zero gap", gap: 0.0, wantStyle: termstyle.OK.Render("x")},
		{name: "tight gap 5pp", gap: 0.05, wantStyle: termstyle.OK.Render("x")},
		{name: "boundary 10pp", gap: 0.10, wantStyle: termstyle.OK.Render("x")},
		// 0.10 < gap ≤ 0.30 → yellow (Accent)
		{name: "medium gap 11pp", gap: 0.11, wantStyle: termstyle.Accent.Render("x")},
		{name: "medium gap 20pp", gap: 0.20, wantStyle: termstyle.Accent.Render("x")},
		{name: "boundary 30pp", gap: 0.30, wantStyle: termstyle.Accent.Render("x")},
		// >0.30 → red (Bad)
		{name: "large gap 31pp", gap: 0.31, wantStyle: termstyle.Bad.Render("x")},
		{name: "large gap 60pp", gap: 0.60, wantStyle: termstyle.Bad.Render("x")},
		{name: "full gap 100pp", gap: 1.0, wantStyle: termstyle.Bad.Render("x")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ciStyleForGap(tc.gap).Render("x")
			assert.Equal(t, tc.wantStyle, got)
		})
	}
}

func TestFormatFlakyTestLine_CIColoring(t *testing.T) {
	t.Parallel()

	// 1/100 → tight CI (green): gap well under 10pp
	tight := TestEntry{Test: "TestTight", Package: "p", Runs: 100, Fails: 1, Successes: 99}
	tightLine := formatFlakyTestLine(tight)
	tlo, thi := WilsonScoreInterval(tight.Fails, tight.Runs, 0)
	tCIText := fmt.Sprintf(" [Confidence Interval: %.1f%%–%.1f%%]", tlo*100, thi*100)
	assert.Contains(t, tightLine, termstyle.OK.Render(tCIText), "tight CI should use green style")

	// 1/2 → wide CI (red): gap well over 30pp
	wide := TestEntry{Test: "TestWide", Package: "p", Runs: 2, Fails: 1, Successes: 1}
	wideLine := formatFlakyTestLine(wide)
	wlo, whi := WilsonScoreInterval(wide.Fails, wide.Runs, 0)
	wCIText := fmt.Sprintf(" [Confidence Interval: %.1f%%–%.1f%%]", wlo*100, whi*100)
	assert.Contains(t, wideLine, termstyle.Bad.Render(wCIText), "wide CI should use red style")
}
