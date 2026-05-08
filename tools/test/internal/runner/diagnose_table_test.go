package runner

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
)

var ansiEscapeSeq = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiEscapeSeq.ReplaceAllString(s, "")
}

func TestDiagnoseTableHeaderPlain(t *testing.T) {
	t.Parallel()
	got := diagnoseTableHeaderPlain()
	want := fmt.Sprintf("%5s  %-8s  %8s  %8s  %8s  %8s  %10s",
		"Iter", "Result", "Tests", "Failures", "Timeouts", "Slow", "Runtime")
	assert.Equal(t, want, got)
	assert.Len(t, got, len(want))
}

func TestPrintDiagnoseIterationTableHeader(t *testing.T) {
	t.Parallel()
	var stderr strings.Builder
	p := output.New(false, &strings.Builder{}, &stderr, output.SkipFD)
	printDiagnoseIterationTableHeader(p)
	s := strings.TrimRight(stripANSI(stderr.String()), "\n")
	lines := strings.Split(s, "\n")
	require.Len(t, lines, 2)
	assert.Equal(t, diagnoseTableHeaderPlain(), lines[0])
	assert.Equal(t, strings.Repeat("─", len(diagnoseTableHeaderPlain())), lines[1])
}

func TestFormatDiagnoseIterationTableRow(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		iter     int
		d        IterationDigest
		dur      time.Duration
		wantSans string // stripped ANSI: digit columns and plain tokens only
	}{
		{
			name: "pass_clean",
			iter: 1,
			d: IterationDigest{
				Result: "pass", RanTests: 0, FailTests: 0, TimeoutTests: 0, SlowTests: 0,
			},
			dur:      3*time.Minute + 7*time.Second,
			wantSans: "    1  pass             0         0         0         0        3m7s",
		},
		{
			name: "fail_with_counts",
			iter: 12,
			d: IterationDigest{
				Result: "fail", RanTests: 0, FailTests: 1, TimeoutTests: 2, SlowTests: 5,
			},
			dur:      90 * time.Second,
			wantSans: "   12  fail             0         1         2         5       1m30s",
		},
		{
			name: "timeout_result",
			iter: 3,
			d: IterationDigest{
				Result: "timeout", RanTests: 0, FailTests: 0, TimeoutTests: 1, SlowTests: 0,
			},
			dur:      time.Hour,
			wantSans: "    3  timeout          0         0         1         0      1h0m0s",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := stripANSI(formatDiagnoseIterationTableRow(tc.iter, tc.d, tc.dur))
			assert.Equal(t, tc.wantSans, got)
		})
	}
}
