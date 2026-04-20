package runner

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readers(iters ...string) []io.Reader {
	rs := make([]io.Reader, len(iters))
	for i, s := range iters {
		rs[i] = strings.NewReader(s)
	}
	return rs
}

func TestAnalyze(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		iterations    []string
		slowThreshold time.Duration
		wantFlakes    []TestEntry
		wantFailures  []TestEntry
		wantTimeouts  []TestEntry
		wantSlow      []TestEntry
	}{
		{
			name: "flake: failed once, passed once",
			iterations: []string{
				`{"Action":"run","Package":"pkg/foo","Test":"TestX"}
{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}
`,
				`{"Action":"run","Package":"pkg/foo","Test":"TestX"}
{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}
`,
			},
			slowThreshold: 30 * time.Second,
			wantFlakes: []TestEntry{
				{Package: "pkg/foo", Test: "TestX", Passes: 1, Fails: 1, MaxElapsed: 500 * time.Millisecond, Iterations: []int{0, 1}},
			},
		},
		{
			name: "deterministic failure: failed in every iteration",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/bar","Test":"TestBroken","Elapsed":0.1}` + "\n",
				`{"Action":"fail","Package":"pkg/bar","Test":"TestBroken","Elapsed":0.1}` + "\n",
			},
			slowThreshold: 30 * time.Second,
			wantFailures: []TestEntry{
				{Package: "pkg/bar", Test: "TestBroken", Fails: 2, MaxElapsed: 100 * time.Millisecond, Iterations: []int{0, 1}},
			},
		},
		{
			name: "timeout: panic output attached to a test",
			iterations: []string{
				`{"Action":"run","Package":"pkg/qux","Test":"TestHang"}
{"Action":"output","Package":"pkg/qux","Test":"TestHang","Output":"panic: test timed out after 10m0s\n"}
{"Action":"fail","Package":"pkg/qux","Test":"TestHang","Elapsed":600.0}
`,
			},
			slowThreshold: 30 * time.Second,
			wantTimeouts: []TestEntry{
				{
					Package: "pkg/qux", Test: "TestHang", Fails: 1,
					MaxElapsed: 600 * time.Second, Iterations: []int{0},
					Logs: []IterationLog{{Iteration: 0, Output: "panic: test timed out after 10m0s\n"}},
				},
			},
			// timeout also exceeds slow threshold; timeouts are excluded from slow list to avoid duplication
			wantFailures: nil,
			wantSlow:     nil,
		},
		{
			name: "timeout: package-level panic without test field",
			iterations: []string{
				`{"Action":"output","Package":"pkg/hang","Output":"panic: test timed out after 2m0s\n"}
{"Action":"fail","Package":"pkg/hang","Elapsed":120.0}
`,
			},
			slowThreshold: 30 * time.Second,
			wantTimeouts: []TestEntry{
				{
					Package: "pkg/hang", Fails: 1,
					MaxElapsed: 120 * time.Second, Iterations: []int{0},
					Logs: []IterationLog{{Iteration: 0, Output: "panic: test timed out after 2m0s\n"}},
				},
			},
		},
		{
			name: "slow: passing test exceeds threshold",
			iterations: []string{
				`{"Action":"run","Package":"pkg/a","Test":"TestSlow"}
{"Action":"pass","Package":"pkg/a","Test":"TestSlow","Elapsed":45.0}
`,
			},
			slowThreshold: 30 * time.Second,
			wantSlow: []TestEntry{
				{Package: "pkg/a", Test: "TestSlow", Passes: 1, MaxElapsed: 45 * time.Second, Iterations: []int{0}},
			},
		},
		{
			name: "clean pass is not reported",
			iterations: []string{
				`{"Action":"pass","Package":"pkg/c","Test":"TestOK","Elapsed":0.01}` + "\n",
			},
			slowThreshold: 30 * time.Second,
		},
		{
			name: "subtests counted independently of parent",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/d","Test":"TestParent/sub1","Elapsed":0.1}
{"Action":"pass","Package":"pkg/d","Test":"TestParent/sub2","Elapsed":0.1}
{"Action":"fail","Package":"pkg/d","Test":"TestParent","Elapsed":0.2}
`,
				`{"Action":"pass","Package":"pkg/d","Test":"TestParent/sub1","Elapsed":0.1}
{"Action":"pass","Package":"pkg/d","Test":"TestParent/sub2","Elapsed":0.1}
{"Action":"pass","Package":"pkg/d","Test":"TestParent","Elapsed":0.2}
`,
			},
			slowThreshold: 30 * time.Second,
			wantFlakes: []TestEntry{
				{Package: "pkg/d", Test: "TestParent", Passes: 1, Fails: 1, MaxElapsed: 200 * time.Millisecond, Iterations: []int{0, 1}},
				{Package: "pkg/d", Test: "TestParent/sub1", Passes: 1, Fails: 1, MaxElapsed: 100 * time.Millisecond, Iterations: []int{0, 1}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, err := Analyze(readers(tc.iterations...), tc.slowThreshold)
			require.NoError(t, err)
			assert.Equal(t, len(tc.iterations), rep.Iterations)
			assert.Equal(t, tc.wantFlakes, rep.Flakes, "flakes")
			assert.Equal(t, tc.wantFailures, rep.Failures, "failures")
			assert.Equal(t, tc.wantTimeouts, rep.Timeouts, "timeouts")
			assert.Equal(t, tc.wantSlow, rep.Slow, "slow")
		})
	}
}

func TestAnalyzeCapturesLogsForFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		iterations []string
		category   string // "flakes","failures","timeouts"
		wantLogs   []IterationLog
	}{
		{
			name: "failure captures output from failing iteration",
			iterations: []string{
				`{"Action":"run","Package":"p","Test":"T"}
{"Action":"output","Package":"p","Test":"T","Output":"    t.go:12: boom\n"}
{"Action":"output","Package":"p","Test":"T","Output":"--- FAIL: T (0.00s)\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":0.01}
`,
			},
			category: "failures",
			wantLogs: []IterationLog{
				{Iteration: 0, Output: "    t.go:12: boom\n--- FAIL: T (0.00s)\n"},
			},
		},
		{
			name: "flake captures logs only from failing iterations",
			iterations: []string{
				`{"Action":"output","Package":"p","Test":"T","Output":"fail-log\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":0.01}
`,
				`{"Action":"output","Package":"p","Test":"T","Output":"ok-log\n"}
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}
`,
			},
			category: "flakes",
			wantLogs: []IterationLog{
				{Iteration: 0, Output: "fail-log\n"},
			},
		},
		{
			name: "timeout captures the panic output",
			iterations: []string{
				`{"Action":"output","Package":"p","Test":"T","Output":"panic: test timed out after 10m0s\n"}
{"Action":"output","Package":"p","Test":"T","Output":"\tstack trace line\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":600.0}
`,
			},
			category: "timeouts",
			wantLogs: []IterationLog{
				{Iteration: 0, Output: "panic: test timed out after 10m0s\n\tstack trace line\n"},
			},
		},
		{
			name: "panic in failing test is captured",
			iterations: []string{
				`{"Action":"output","Package":"p","Test":"T","Output":"panic: runtime error: nil pointer\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":0.01}
`,
			},
			category: "failures",
			wantLogs: []IterationLog{
				{Iteration: 0, Output: "panic: runtime error: nil pointer\n"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, err := Analyze(readers(tc.iterations...), 30*time.Second)
			require.NoError(t, err)
			var entries []TestEntry
			switch tc.category {
			case "flakes":
				entries = rep.Flakes
			case "failures":
				entries = rep.Failures
			case "timeouts":
				entries = rep.Timeouts
			}
			require.Len(t, entries, 1, "expected exactly one %s entry", tc.category)
			assert.Equal(t, tc.wantLogs, entries[0].Logs)
		})
	}
}

func TestAnalyzeReattributesTimeoutToRunningTests(t *testing.T) {
	t.Parallel()
	// go test attaches the timeout panic to whatever test emitted last —
	// here TestFast has already passed but gets the panic output. The real
	// culprit is listed in the "running tests:" block.
	iter := `{"Action":"run","Package":"p","Test":"TestFast"}
{"Action":"pass","Package":"p","Test":"TestFast","Elapsed":0.01}
{"Action":"output","Package":"p","Test":"TestFast","Output":"panic: test timed out after 5s\n"}
{"Action":"output","Package":"p","Test":"TestFast","Output":"\trunning tests:\n"}
{"Action":"output","Package":"p","Test":"TestFast","Output":"\t\tTestSlow/sub_case (5s)\n"}
{"Action":"output","Package":"p","Test":"TestFast","Output":"\t\tTestOther (4s)\n"}
{"Action":"output","Package":"p","Test":"TestFast","Output":"\n"}
{"Action":"output","Package":"p","Test":"TestFast","Output":"goroutine 1 [chan receive]:\n"}
{"Action":"fail","Package":"p","Elapsed":5.01}
`
	rep, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)

	names := make([]string, 0, len(rep.Timeouts))
	for _, e := range rep.Timeouts {
		names = append(names, e.Test)
	}
	assert.ElementsMatch(t, []string{"TestSlow/sub_case", "TestOther"}, names,
		"timeout should be attributed to the tests listed in `running tests:`")

	// The misattributed test (TestFast) must not show up as a timeout.
	for _, e := range rep.Timeouts {
		assert.NotEqual(t, "TestFast", e.Test)
	}

	// The stack trace should travel with the re-attributed entries.
	for _, e := range rep.Timeouts {
		require.Len(t, e.Logs, 1)
		assert.Contains(t, e.Logs[0].Output, "panic: test timed out after 5s")
	}
}

func TestAnalyzeKeepsTimeoutOnCulpritWhenItWasTheReportedTest(t *testing.T) {
	t.Parallel()
	// If the test whose output carried the panic IS in the running-tests list,
	// the timeout stays attributed to it (no bogus re-attribution).
	iter := `{"Action":"output","Package":"p","Test":"TestSlow","Output":"panic: test timed out after 5s\n"}
{"Action":"output","Package":"p","Test":"TestSlow","Output":"\trunning tests:\n"}
{"Action":"output","Package":"p","Test":"TestSlow","Output":"\t\tTestSlow (5s)\n"}
{"Action":"fail","Package":"p","Elapsed":5.01}
`
	rep, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Timeouts, 1)
	assert.Equal(t, "TestSlow", rep.Timeouts[0].Test)
}

func TestAnalyzePassingTestsHaveNoLogs(t *testing.T) {
	t.Parallel()
	input := `{"Action":"output","Package":"p","Test":"T","Output":"hello\n"}
{"Action":"pass","Package":"p","Test":"T","Elapsed":45.0}
`
	rep, err := Analyze(readers(input), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Slow, 1)
	assert.Nil(t, rep.Slow[0].Logs)
}

func TestPrintSummaryTimeoutShowsIterationsNotPassCounts(t *testing.T) {
	t.Parallel()
	rep := &Report{
		Iterations:    3,
		SlowThreshold: 30 * time.Second,
		Timeouts: []TestEntry{
			{Package: "p", Test: "TestStuck", Passes: 2, Iterations: []int{0, 2}},
		},
	}
	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()
	assert.Contains(t, out, "p.TestStuck")
	assert.Contains(t, out, "iter 0,2")
	// The misleading "(2p/0f)" format must not appear for timeouts.
	assert.NotContains(t, out, "(2p/0f)")
}

func TestAnalyzeResultsRoundtrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// two iterations: same test flakes (fail, pass)
	must(t, os.WriteFile(filepath.Join(dir, "iteration-0.log.jsonl"),
		[]byte(`{"Action":"fail","Package":"pkg/z","Test":"TestFlaky","Elapsed":0.2}`+"\n"), 0600))
	must(t, os.WriteFile(filepath.Join(dir, "iteration-1.log.jsonl"),
		[]byte(`{"Action":"pass","Package":"pkg/z","Test":"TestFlaky","Elapsed":0.1}`+"\n"), 0600))

	rep, err := AnalyzeResults(dir, 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Flakes, 1)
	assert.Equal(t, "TestFlaky", rep.Flakes[0].Test)

	require.NoError(t, WriteReport(dir, rep))
	b, err := os.ReadFile(filepath.Join(dir, "report.json"))
	require.NoError(t, err)
	assert.Contains(t, string(b), `"flakes"`)
	assert.Contains(t, string(b), `"TestFlaky"`)
}

func must(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

func TestAnalyzeSkipsMalformedLines(t *testing.T) {
	t.Parallel()
	input := `not json at all
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}
`
	rep, err := Analyze(readers(input), 30*time.Second)
	require.NoError(t, err)
	assert.Empty(t, rep.Flakes)
	assert.Empty(t, rep.Failures)
}
