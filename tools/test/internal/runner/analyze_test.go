package runner

import (
	"bufio"
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

func TestAnalyzePackageLevelTimeoutIterationSummary(t *testing.T) {
	t.Parallel()
	iterations := []string{
		`{"Action":"output","Package":"pkg/hang","Output":"panic: test timed out after 2m0s\n"}
{"Action":"fail","Package":"pkg/hang","Elapsed":120.0}
`,
	}
	rep, _, err := Analyze(readers(iterations...), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.IterationSummaries, 1)
	assert.Equal(t, "timeout", rep.IterationSummaries[0].Result)
}

func TestAnalyzeLineExceedsDefaultScannerLimit(t *testing.T) {
	t.Parallel()
	// bufio.Scanner default max token is bufio.MaxScanTokenSize (64 KiB).
	over := strings.Repeat("x", bufio.MaxScanTokenSize+1) + "\n"
	iter := `{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}` + "\n" + over +
		`{"Action":"pass","Package":"p","Test":"T2","Elapsed":0.01}` + "\n"
	_, _, err := Analyze(readers(iter), 30*time.Second)
	require.Error(t, err)
	require.ErrorContains(t, err, "reading iteration 0")
	require.ErrorIs(t, err, bufio.ErrTooLong, "want bufio.ErrTooLong wrapped in analyze error")
}

func TestAnalyzeBuildErrorsInterleavedWithJSONL(t *testing.T) {
	t.Parallel()
	// go test -json can interleave compiler lines (non-JSON) with events; package build ends with fail, Test "".
	iter := `# example.com/badpkg
badpkg.go:1:2: undefined: MissingType
` + `{"Action":"output","Package":"example.com/badpkg","Output":"# example.com/badpkg\n"}
{"Action":"fail","Package":"example.com/badpkg","Elapsed":0.0}
`
	rep, _, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Failures, 1)
	assert.Equal(t, "example.com/badpkg", rep.Failures[0].Package)
	assert.Empty(t, rep.Failures[0].Test)
	require.Len(t, rep.IterationSummaries, 1)
	assert.Equal(t, "fail", rep.IterationSummaries[0].Result)
	assert.Equal(t, []string{"example.com/badpkg"}, rep.IterationSummaries[0].FailingTests)
}

func TestAnalyzePackageLevelFailureIterationSummary(t *testing.T) {
	t.Parallel()
	// go test -json uses Test == "" for build failures, TestMain failures, etc.
	iterations := []string{
		`{"Action":"fail","Package":"pkg/build","Elapsed":0.0}` + "\n",
	}
	rep, _, err := Analyze(readers(iterations...), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.IterationSummaries, 1)
	assert.Equal(t, "fail", rep.IterationSummaries[0].Result)
	assert.Equal(t, []string{"pkg/build"}, rep.IterationSummaries[0].FailingTests)
	require.Len(t, rep.Failures, 1)
	assert.Equal(t, "pkg/build", rep.Failures[0].Package)
	assert.Empty(t, rep.Failures[0].Test)
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
				{
					Package: "pkg/foo", Test: "TestX",
					Runs: 2, Successes: 1, Fails: 1,
					MinElapsed: 400 * time.Millisecond,
					MaxElapsed: 500 * time.Millisecond,
					P50Elapsed: 450 * time.Millisecond,
				},
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
				{
					Package: "pkg/bar", Test: "TestBroken",
					Runs: 2, Fails: 2,
					MinElapsed: 100 * time.Millisecond,
					MaxElapsed: 100 * time.Millisecond,
					P50Elapsed: 100 * time.Millisecond,
				},
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
					Package: "pkg/qux", Test: "TestHang",
					Runs: 1, Fails: 1, Timeouts: 1,
					MinElapsed: 600 * time.Second,
					MaxElapsed: 600 * time.Second,
					P50Elapsed: 600 * time.Second,
				},
			},
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
					Package: "pkg/hang",
					Runs:    1, Fails: 1, Timeouts: 1,
					MinElapsed: 120 * time.Second,
					MaxElapsed: 120 * time.Second,
					P50Elapsed: 120 * time.Second,
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
				{
					Package: "pkg/a", Test: "TestSlow",
					Runs: 1, Successes: 1,
					MinElapsed: 45 * time.Second,
					MaxElapsed: 45 * time.Second,
					P50Elapsed: 45 * time.Second,
				},
			},
		},
		{
			name: "package-level failure without test name (build/TestMain)",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/build","Elapsed":0.0}` + "\n",
			},
			slowThreshold: 30 * time.Second,
			wantFailures: []TestEntry{
				{
					Package:    "pkg/build",
					Test:       "",
					Runs:       1,
					Fails:      1,
					MinElapsed: 0,
					MaxElapsed: 0,
					P50Elapsed: 0,
				},
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
			name: "skips-only test is not flagged",
			iterations: []string{
				`{"Action":"skip","Package":"pkg/s","Test":"TestSkipped","Elapsed":0.0}` + "\n",
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
				{
					Package: "pkg/d", Test: "TestParent",
					Runs: 2, Successes: 1, Fails: 1,
					MinElapsed: 200 * time.Millisecond,
					MaxElapsed: 200 * time.Millisecond,
					P50Elapsed: 200 * time.Millisecond,
				},
				{
					Package: "pkg/d", Test: "TestParent/sub1",
					Runs: 2, Successes: 1, Fails: 1,
					MinElapsed: 100 * time.Millisecond,
					MaxElapsed: 100 * time.Millisecond,
					P50Elapsed: 100 * time.Millisecond,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, _, err := Analyze(readers(tc.iterations...), tc.slowThreshold)
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
		wantKey    testKey
		wantIter   int
		wantOutput string
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
			category:   "failures",
			wantKey:    testKey{Package: "p", Test: "T"},
			wantIter:   0,
			wantOutput: "    t.go:12: boom\n--- FAIL: T (0.00s)\n",
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
			category:   "flakes",
			wantKey:    testKey{Package: "p", Test: "T"},
			wantIter:   0,
			wantOutput: "fail-log\n",
		},
		{
			name: "timeout captures the panic output",
			iterations: []string{
				`{"Action":"output","Package":"p","Test":"T","Output":"panic: test timed out after 10m0s\n"}
{"Action":"output","Package":"p","Test":"T","Output":"\tstack trace line\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":600.0}
`,
			},
			category:   "timeouts",
			wantKey:    testKey{Package: "p", Test: "T"},
			wantIter:   0,
			wantOutput: "panic: test timed out after 10m0s\n\tstack trace line\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, logs, err := Analyze(readers(tc.iterations...), 30*time.Second)
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
			require.Contains(t, logs, tc.wantKey, "log map should contain the flagged test")
			assert.Equal(t, tc.wantOutput, logs[tc.wantKey][tc.wantIter])
		})
	}
}

func TestAnalyzeReattributesTimeoutToRunningTests(t *testing.T) {
	t.Parallel()
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
	rep, logs, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)

	names := make([]string, 0, len(rep.Timeouts))
	for _, e := range rep.Timeouts {
		names = append(names, e.Test)
	}
	assert.ElementsMatch(t, []string{"TestSlow/sub_case", "TestOther"}, names)
	for _, e := range rep.Timeouts {
		assert.NotEqual(t, "TestFast", e.Test)
	}
	for _, e := range rep.Timeouts {
		k := testKey{Package: e.Package, Test: e.Test}
		require.Contains(t, logs, k)
		assert.Contains(t, logs[k][0], "panic: test timed out after 5s")
	}
}

func TestAnalyzeKeepsTimeoutOnCulpritWhenItWasTheReportedTest(t *testing.T) {
	t.Parallel()
	iter := `{"Action":"output","Package":"p","Test":"TestSlow","Output":"panic: test timed out after 5s\n"}
{"Action":"output","Package":"p","Test":"TestSlow","Output":"\trunning tests:\n"}
{"Action":"output","Package":"p","Test":"TestSlow","Output":"\t\tTestSlow (5s)\n"}
{"Action":"fail","Package":"p","Elapsed":5.01}
`
	rep, _, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Timeouts, 1)
	assert.Equal(t, "TestSlow", rep.Timeouts[0].Test)
}

func TestPrintSummaryTimeoutShowsTestNotPassCounts(t *testing.T) {
	t.Parallel()
	rep := &Report{
		Iterations:    3,
		SlowThreshold: 30 * time.Second,
		Timeouts: []TestEntry{
			{Package: "p", Test: "TestStuck", Successes: 2},
		},
	}
	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()
	assert.Contains(t, out, "Timeout (1)")
	assert.Contains(t, out, "|-- p/")
	assert.Contains(t, out, "TestStuck")
	assert.NotContains(t, out, "(2p/0f)")
}

func TestAnalyzeResultsRoundtrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "iteration-0.log.jsonl"),
		[]byte(`{"Action":"fail","Package":"pkg/z","Test":"TestFlaky","Elapsed":0.2}`+"\n"), 0600))
	must(t, os.WriteFile(filepath.Join(dir, "iteration-1.log.jsonl"),
		[]byte(`{"Action":"pass","Package":"pkg/z","Test":"TestFlaky","Elapsed":0.1}`+"\n"), 0600))

	rep, _, err := AnalyzeResults(dir, 30*time.Second)
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

func TestAnalyzeIterationSummaries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		iterations []string
		want       []IterationSummary
	}{
		{
			name: "all pass",
			iterations: []string{
				`{"Action":"pass","Package":"p","Test":"T","Elapsed":0.1}` + "\n",
				`{"Action":"pass","Package":"p","Test":"T","Elapsed":0.2}` + "\n",
			},
			want: []IterationSummary{
				{Index: 0, Result: "pass"},
				{Index: 1, Result: "pass"},
			},
		},
		{
			name: "fail then pass",
			iterations: []string{
				`{"Action":"fail","Package":"p","Test":"TestA","Elapsed":0.1}` + "\n",
				`{"Action":"pass","Package":"p","Test":"TestA","Elapsed":0.2}` + "\n",
			},
			want: []IterationSummary{
				{Index: 0, Result: "fail", FailingTests: []string{"TestA"}},
				{Index: 1, Result: "pass"},
			},
		},
		{
			name: "timeout",
			iterations: []string{
				`{"Action":"output","Package":"p","Test":"TestHang","Output":"panic: test timed out after 10m0s\n"}` + "\n" +
					`{"Action":"fail","Package":"p","Test":"TestHang","Elapsed":600.0}` + "\n",
			},
			want: []IterationSummary{
				{Index: 0, Result: "timeout"},
			},
		},
		{
			name: "multiple failures sorted",
			iterations: []string{
				`{"Action":"fail","Package":"p","Test":"TestB","Elapsed":0.1}` + "\n" +
					`{"Action":"fail","Package":"p","Test":"TestA","Elapsed":0.1}` + "\n",
			},
			want: []IterationSummary{
				{Index: 0, Result: "fail", FailingTests: []string{"TestA", "TestB"}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, _, err := Analyze(readers(tc.iterations...), 30*time.Second)
			require.NoError(t, err)
			require.Len(t, rep.IterationSummaries, len(tc.want))
			// Strip Duration/ShuffleSeed — set by runner, not Analyze.
			got := make([]IterationSummary, len(rep.IterationSummaries))
			for i, s := range rep.IterationSummaries {
				got[i] = IterationSummary{Index: s.Index, Result: s.Result, FailingTests: s.FailingTests}
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAnalyzeSkipsMalformedLines(t *testing.T) {
	t.Parallel()
	input := `not json at all
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}
`
	rep, _, err := Analyze(readers(input), 30*time.Second)
	require.NoError(t, err)
	assert.Empty(t, rep.Flakes)
	assert.Empty(t, rep.Failures)
}

func TestStatsP50(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		samples []time.Duration
		wantMin time.Duration
		wantP50 time.Duration
	}{
		{
			name:    "empty",
			samples: nil,
			wantMin: 0,
			wantP50: 0,
		},
		{
			name:    "single",
			samples: []time.Duration{5 * time.Second},
			wantMin: 5 * time.Second,
			wantP50: 5 * time.Second,
		},
		{
			name:    "odd count",
			samples: []time.Duration{3, 1, 2},
			wantMin: 1,
			wantP50: 2,
		},
		{
			name:    "even count averages middle two",
			samples: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second, 9 * time.Second},
			wantMin: 1 * time.Second,
			wantP50: 4 * time.Second,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			minDur, p50 := stats(tc.samples)
			assert.Equal(t, tc.wantMin, minDur, "min")
			assert.Equal(t, tc.wantP50, p50, "p50")
		})
	}
}
