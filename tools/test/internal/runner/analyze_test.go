package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
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

func TestAnalyzeHandlesLongLines(t *testing.T) {
	t.Parallel()
	// Create a line that exceeds bufio.MaxScanTokenSize (64 KiB).
	over := strings.Repeat("x", bufio.MaxScanTokenSize+1) + "\n"
	iter := `{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}` + "\n" + over +
		`{"Action":"pass","Package":"p","Test":"T2","Elapsed":0.01}` + "\n"
	rep, _, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.NotNil(t, rep)
	require.Len(t, rep.IterationSummaries, 1)
	assert.Equal(t, "pass", rep.IterationSummaries[0].Result)
	assert.NotNil(t, rep.Summary)
	assert.Equal(t, 2, rep.Summary.DistinctNamedTests)
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

func TestAnalyzeTestdataFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		file         string
		wantResult   string
		wantFailPkg  int
		wantFailTest []string
		wantTimeouts []string
	}{
		{
			name:         "build failure",
			file:         "build-failure.log.jsonl",
			wantResult:   "fail",
			wantFailPkg:  1,
			wantFailTest: nil,
			wantTimeouts: nil,
		},
		{
			name:         "panic",
			file:         "panic.log.jsonl",
			wantResult:   "fail",
			wantFailPkg:  1, // package-level failure
			wantFailTest: []string{"TestPanic", "TestPanic/test5"},
			wantTimeouts: nil,
		},
		{
			name:         "timeout",
			file:         "timeout.log.jsonl",
			wantResult:   "timeout",
			wantFailPkg:  1, // package-level failure
			wantFailTest: nil,
			wantTimeouts: []string{"TestTimeout/test5"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(filepath.Join("testdata", tc.file))
			require.NoError(t, err)
			defer f.Close()

			rep, _, err := Analyze([]io.Reader{f}, 30*time.Second)
			require.NoError(t, err)

			require.Len(t, rep.IterationSummaries, 1)
			assert.Equal(t, tc.wantResult, rep.IterationSummaries[0].Result)

			if tc.wantResult == "fail" && len(tc.wantFailTest) > 0 {
				assert.ElementsMatch(t, tc.wantFailTest, rep.IterationSummaries[0].FailingTests)
			} else if tc.wantResult == "fail" && len(tc.wantFailTest) == 0 && tc.wantFailPkg > 0 {
				// build failure typically puts package name in FailingTests
				require.NotEmpty(t, rep.IterationSummaries[0].FailingTests)
				assert.Contains(t, rep.IterationSummaries[0].FailingTests[0], "github.com/smartcontractkit/chainlink")
			}

			var failTests []string
			var failPkgs int
			for _, fail := range rep.Failures {
				if fail.Test == "" {
					failPkgs++
				} else {
					failTests = append(failTests, fail.Test)
				}
			}
			assert.Equal(t, tc.wantFailPkg, failPkgs, "package level failures")
			assert.ElementsMatch(t, tc.wantFailTest, failTests, "test level failures")

			var timeoutTests []string
			for _, to := range rep.Timeouts {
				timeoutTests = append(timeoutTests, to.Test)
			}
			assert.ElementsMatch(t, tc.wantTimeouts, timeoutTests, "timeouts")
		})
	}
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

func TestDigestIterationJSONL(t *testing.T) {
	t.Parallel()

	t.Run("pass no slow", func(t *testing.T) {
		t.Parallel()
		iter := `{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}` + "\n"
		d, err := DigestIterationJSONL(strings.NewReader(iter), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, "pass", d.Result)
		assert.Equal(t, 1, d.RanTests)
		assert.Equal(t, 0, d.FailTests)
		assert.Equal(t, 0, d.SlowTests)
		assert.Equal(t, 0, d.TimeoutTests)
		assert.False(t, d.BuildFailure)
	})

	t.Run("slow test", func(t *testing.T) {
		t.Parallel()
		slowJSON := `{"Action":"run","Package":"pkg/a","Test":"TestSlow"}
{"Action":"pass","Package":"pkg/a","Test":"TestSlow","Elapsed":45.0}
`
		d, err := DigestIterationJSONL(strings.NewReader(slowJSON), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 1, d.SlowTests)
		assert.Equal(t, 1, d.RanTests)
		assert.Equal(t, "pass", d.Result)
	})

	t.Run("package fail", func(t *testing.T) {
		t.Parallel()
		failJSON := `{"Action":"fail","Package":"pkg/build","Elapsed":0.0}` + "\n"
		d, err := DigestIterationJSONL(strings.NewReader(failJSON), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, "fail", d.Result)
		assert.Equal(t, 0, d.RanTests)
		assert.Equal(t, 1, d.FailTests)
		assert.True(t, d.BuildFailure)
	})

	t.Run("failed_build_field", func(t *testing.T) {
		t.Parallel()
		jsonl := `{"Action":"fail","Package":"example.com/badpkg","Elapsed":0,"FailedBuild":"example.com/badpkg.test"}` + "\n"
		d, err := DigestIterationJSONL(strings.NewReader(jsonl), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, "fail", d.Result)
		assert.True(t, d.BuildFailure)
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()
		toJSON := `{"Action":"output","Package":"pkg/hang","Output":"panic: test timed out after 2m0s\n"}
{"Action":"fail","Package":"pkg/hang","Elapsed":120.0}
`
		d, err := DigestIterationJSONL(strings.NewReader(toJSON), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, "timeout", d.Result)
		assert.Equal(t, 0, d.RanTests)
		assert.GreaterOrEqual(t, d.TimeoutTests, 1)
	})

	t.Run("two named tests", func(t *testing.T) {
		t.Parallel()
		jsonl := `{"Action":"pass","Package":"p","Test":"TA","Elapsed":0.01}
{"Action":"pass","Package":"p","Test":"TB","Elapsed":0.02}
`
		d, err := DigestIterationJSONL(strings.NewReader(jsonl), 30*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 2, d.RanTests)
		assert.Equal(t, "pass", d.Result)
	})
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
			assert.Equal(t, tc.wantFlakes, publicTestEntries(rep.Flakes), "flakes")
			assert.Equal(t, tc.wantFailures, publicTestEntries(rep.Failures), "failures")
			assert.Equal(t, tc.wantTimeouts, publicTestEntries(rep.Timeouts), "timeouts")
			assert.Equal(t, tc.wantSlow, publicTestEntries(rep.Slow), "slow")
		})
	}
}

func TestReportSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		iterations    []string
		slowThreshold time.Duration
		check         func(t *testing.T, s *ReportSummary)
	}{
		{
			name:          "no iterations yields nil summary",
			iterations:    []string{},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				assert.Nil(t, s)
			},
		},
		{
			name: "flake prevalence and execution rate",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
				`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
			},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 1, s.DistinctNamedTests)
				assert.Equal(t, 1, s.FlakeNamedCount)
				require.NotNil(t, s.FlakePrevalence)
				assert.InDelta(t, 1.0, *s.FlakePrevalence, 1e-9)
				assert.Equal(t, 1, s.FlakeFailRuns)
				assert.Equal(t, 2, s.FlakeTotalRuns)
				require.NotNil(t, s.FlakeExecutionFailRate)
				assert.InDelta(t, 0.5, *s.FlakeExecutionFailRate, 1e-9)
				assert.Equal(t, 1, s.FlakeFailingIterations)
				assert.Equal(t, 2, s.FlakeIterationTotal)
				require.NotNil(t, s.FlakeIterationFailRate)
				assert.InDelta(t, 0.5, *s.FlakeIterationFailRate, 1e-9)
				require.NotNil(t, s.SlowPrevalence)
				assert.InDelta(t, 0.0, *s.SlowPrevalence, 1e-9)
			},
		},
		{
			name: "deterministic failure has zero flake buckets",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/bar","Test":"TestBroken","Elapsed":0.1}` + "\n",
				`{"Action":"fail","Package":"pkg/bar","Test":"TestBroken","Elapsed":0.1}` + "\n",
			},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 1, s.DistinctNamedTests)
				assert.Equal(t, 0, s.FlakeNamedCount)
				require.NotNil(t, s.FlakePrevalence)
				assert.InDelta(t, 0.0, *s.FlakePrevalence, 1e-9)
				assert.Equal(t, 0, s.FlakeFailRuns)
				assert.Equal(t, 0, s.FlakeTotalRuns)
				assert.Nil(t, s.FlakeExecutionFailRate)
				assert.Equal(t, 0, s.FlakeFailingIterations)
				assert.Equal(t, 0, s.FlakeIterationTotal)
				assert.Nil(t, s.FlakeIterationFailRate)
			},
		},
		{
			name: "slow prevalence",
			iterations: []string{
				`{"Action":"pass","Package":"pkg/a","Test":"TestSlow","Elapsed":45.0}`,
			},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 1, s.DistinctNamedTests)
				assert.Equal(t, 1, s.SlowCount)
				require.NotNil(t, s.SlowPrevalence)
				assert.InDelta(t, 1.0, *s.SlowPrevalence, 1e-9)
			},
		},
		{
			name: "clean pass zero flake and slow rates",
			iterations: []string{
				`{"Action":"pass","Package":"pkg/c","Test":"TestOK","Elapsed":0.01}`,
			},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 1, s.DistinctNamedTests)
				assert.Equal(t, 0, s.FlakeNamedCount)
				require.NotNil(t, s.FlakePrevalence)
				assert.InDelta(t, 0.0, *s.FlakePrevalence, 1e-9)
				require.NotNil(t, s.SlowPrevalence)
				assert.InDelta(t, 0.0, *s.SlowPrevalence, 1e-9)
				assert.Equal(t, 0, s.FlakeFailingIterations)
				assert.Equal(t, 0, s.FlakeIterationTotal)
				assert.Nil(t, s.FlakeIterationFailRate)
			},
		},
		{
			name: "slow prevalence omitted when threshold disabled",
			iterations: []string{
				`{"Action":"pass","Package":"pkg/a","Test":"TestSlow","Elapsed":45.0}`,
			},
			slowThreshold: 0,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 0, s.SlowCount)
				assert.Nil(t, s.SlowPrevalence)
			},
		},
		{
			name: "package-level flake execution without named tests",
			iterations: []string{
				`{"Action":"fail","Package":"pkg/build","Elapsed":0.0}` + "\n",
				`{"Action":"pass","Package":"pkg/build","Elapsed":0.0}` + "\n",
			},
			slowThreshold: 30 * time.Second,
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 0, s.DistinctNamedTests)
				assert.Equal(t, 0, s.FlakeNamedCount)
				assert.Nil(t, s.FlakePrevalence)
				assert.Equal(t, 1, s.FlakeFailRuns)
				assert.Equal(t, 2, s.FlakeTotalRuns)
				require.NotNil(t, s.FlakeExecutionFailRate)
				assert.InDelta(t, 0.5, *s.FlakeExecutionFailRate, 1e-9)
				assert.Equal(t, 1, s.FlakeFailingIterations)
				assert.Equal(t, 2, s.FlakeIterationTotal)
				require.NotNil(t, s.FlakeIterationFailRate)
				assert.InDelta(t, 0.5, *s.FlakeIterationFailRate, 1e-9)
			},
		},
		{
			name: "multiple flakes prevalence over distinct named tests",
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
			check: func(t *testing.T, s *ReportSummary) {
				require.NotNil(t, s)
				assert.Equal(t, 3, s.DistinctNamedTests)
				assert.Equal(t, 2, s.FlakeNamedCount)
				require.NotNil(t, s.FlakePrevalence)
				assert.InDelta(t, 2.0/3.0, *s.FlakePrevalence, 1e-9)
				assert.Equal(t, 2, s.FlakeFailRuns)
				assert.Equal(t, 4, s.FlakeTotalRuns)
				assert.Equal(t, 1, s.FlakeFailingIterations)
				assert.Equal(t, 2, s.FlakeIterationTotal)
				require.NotNil(t, s.FlakeIterationFailRate)
				assert.InDelta(t, 0.5, *s.FlakeIterationFailRate, 1e-9)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep, _, err := Analyze(readers(tc.iterations...), tc.slowThreshold)
			require.NoError(t, err)
			tc.check(t, rep.Summary)
		})
	}
}

func TestPrintSummaryOverallContains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		prep   func(t *testing.T) *Report
		needle []string
	}{
		{
			name: "flake_rates_and_slow_line",
			prep: func(t *testing.T) *Report {
				rep, _, err := Analyze(readers(
					`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
					`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
				), 30*time.Second)
				require.NoError(t, err)
				return rep
			},
			needle: []string{"Overall", "Broken tests:", "Flaky tests:", "Flaky Iterations: 1/2 (50.0%)", "Slow tests:"},
		},
		{
			name: "iteration_wall_clock_runtimes",
			prep: func(t *testing.T) *Report {
				rep, _, err := Analyze(readers(`{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}`), 30*time.Second)
				require.NoError(t, err)
				require.NotNil(t, rep.Summary)
				rep.IterationSummaries[0].Duration = 5 * time.Second
				fillIterationRuntimeSummary(rep)
				return rep
			},
			needle: []string{"Overall", "Iteration runtimes:", "min=5s", "Broken tests:"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf strings.Builder
			PrintSummary(&buf, tc.prep(t))
			out := buf.String()
			for _, s := range tc.needle {
				assert.Contains(t, out, s)
			}
		})
	}
}

func TestPrintSummaryOverall_usesSeverityColors(t *testing.T) {
	t.Parallel()
	rep, _, err := Analyze(readers(
		`{"Action":"fail","Package":"pkg/foo","Test":"TestX","Elapsed":0.5}`,
		`{"Action":"pass","Package":"pkg/foo","Test":"TestX","Elapsed":0.4}`,
	), 30*time.Second)
	require.NoError(t, err)
	require.NotNil(t, rep.Summary)
	s := rep.Summary

	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()

	brokenN := len(rep.Failures)
	pctBroken := float64(brokenN) / float64(s.DistinctNamedTests) * 100
	brokenLine := fmt.Sprintf("  Broken tests: %d/%d (%.1f%%)", brokenN, s.DistinctNamedTests, pctBroken)
	if brokenN > 0 {
		assert.Contains(t, out, termstyle.Bad.Render(brokenLine))
	} else {
		assert.Contains(t, out, termstyle.OK.Render(brokenLine))
	}

	pctFlake := 0.0
	if s.FlakePrevalence != nil {
		pctFlake = *s.FlakePrevalence * 100
	}
	flakyLine := fmt.Sprintf("  Flaky tests: %d/%d (%.1f%%)", s.FlakeNamedCount, s.DistinctNamedTests, pctFlake)
	if s.FlakeNamedCount > 0 {
		assert.Contains(t, out, termstyle.Bad.Render(flakyLine))
	} else {
		assert.Contains(t, out, termstyle.OK.Render(flakyLine))
	}

	if len(rep.Flakes) > 0 && s.FlakeIterationFailRate != nil {
		pctFI := *s.FlakeIterationFailRate * 100
		fiLine := fmt.Sprintf("  Flaky Iterations: %d/%d (%.1f%%)", s.FlakeFailingIterations, s.FlakeIterationTotal, pctFI)
		if s.FlakeFailingIterations > 0 {
			assert.Contains(t, out, termstyle.Bad.Render(fiLine))
		} else {
			assert.Contains(t, out, termstyle.OK.Render(fiLine))
		}
	}

	if rep.SlowThreshold > 0 && s.DistinctNamedTests > 0 && s.SlowPrevalence != nil {
		pctSlow := *s.SlowPrevalence * 100
		slowLine := fmt.Sprintf("  Slow tests: %d/%d (%.1f%%)", s.SlowCount, s.DistinctNamedTests, pctSlow)
		if s.SlowCount > 0 {
			assert.Contains(t, out, termstyle.Accent.Render(slowLine))
		} else {
			assert.Contains(t, out, termstyle.OK.Render(slowLine))
		}
	}
}

func publicTestEntries(entries []TestEntry) []TestEntry {
	out := append([]TestEntry(nil), entries...)
	for i := range out {
		out[i].Logs = nil
		out[i].FailIters = nil
		out[i].TimeoutIters = nil
		out[i].SlowIters = nil
	}
	return out
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

func TestPrintSummaryPackageLevelFlakeDoesNotPrintPackageAsTest(t *testing.T) {
	t.Parallel()
	rep := &Report{
		Iterations:    50,
		SlowThreshold: 30 * time.Second,
		Flakes: []TestEntry{
			{Package: "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2", Runs: 50, Fails: 4},
			{Package: "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2", Test: "TestVRFV2Integration_SingleConsumer_ForceFulfillment", Runs: 48, Fails: 1},
		},
	}
	var buf strings.Builder
	PrintSummary(&buf, rep)
	out := buf.String()
	assert.Contains(t, out, "Flaky (2)")
	assert.Contains(t, out, "|-- v2/ (4/50) 8.0%")
	assert.Contains(t, out, "|---- TestVRFV2Integration_SingleConsumer_ForceFulfillment (1/48) 2.1%")
	assert.NotContains(t, out, "|---- github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2")
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
		{
			name: "package fail summary suppressed when tests fail",
			iterations: []string{
				`{"Action":"fail","Package":"p","Test":"TestA","Elapsed":0.1}` + "\n" +
					`{"Action":"fail","Package":"p","Elapsed":0.1}` + "\n",
			},
			want: []IterationSummary{
				{Index: 0, Result: "fail", FailingTests: []string{"TestA"}},
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

func TestDurationSampleStats(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		samples []time.Duration
		wantMin time.Duration
		wantMax time.Duration
		wantP50 time.Duration
	}{
		{
			name:    "empty",
			samples: nil,
			wantMin: 0,
			wantMax: 0,
			wantP50: 0,
		},
		{
			name:    "single",
			samples: []time.Duration{5 * time.Second},
			wantMin: 5 * time.Second,
			wantMax: 5 * time.Second,
			wantP50: 5 * time.Second,
		},
		{
			name:    "odd_count_unsorted",
			samples: []time.Duration{3, 1, 2},
			wantMin: 1,
			wantMax: 3,
			wantP50: 2,
		},
		{
			name: "even_count_averages_middle_two",
			samples: []time.Duration{
				1 * time.Second, 3 * time.Second, 5 * time.Second, 9 * time.Second,
			},
			wantMin: 1 * time.Second,
			wantMax: 9 * time.Second,
			wantP50: 4 * time.Second,
		},
		{
			name: "three_spread_values",
			samples: []time.Duration{
				100 * time.Millisecond,
				300 * time.Millisecond,
				200 * time.Millisecond,
			},
			wantMin: 100 * time.Millisecond,
			wantMax: 300 * time.Millisecond,
			wantP50: 200 * time.Millisecond,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotMin, gotMax, gotP50 := sortedDurationStats(tc.samples)
			assert.Equal(t, tc.wantMin, gotMin, "sortedDurationStats min")
			assert.Equal(t, tc.wantMax, gotMax, "sortedDurationStats max")
			assert.Equal(t, tc.wantP50, gotP50, "sortedDurationStats p50")
			min2, p502 := stats(tc.samples)
			assert.Equal(t, tc.wantMin, min2, "stats min")
			assert.Equal(t, tc.wantP50, p502, "stats p50")
		})
	}
}

func TestFillIterationRuntimeSummaryTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		iters   []IterationSummary
		wantMin time.Duration
		wantMax time.Duration
		wantP50 time.Duration
	}{
		{
			name: "three_wall_clock_samples",
			iters: []IterationSummary{
				{Index: 0, Duration: 10 * time.Second},
				{Index: 1, Duration: 30 * time.Second},
				{Index: 2, Duration: 20 * time.Second},
			},
			wantMin: 10 * time.Second,
			wantMax: 30 * time.Second,
			wantP50: 20 * time.Second,
		},
		{
			name: "skips_zero_duration",
			iters: []IterationSummary{
				{Index: 0, Duration: 0},
				{Index: 1, Duration: 10 * time.Second},
			},
			wantMin: 10 * time.Second,
			wantMax: 10 * time.Second,
			wantP50: 10 * time.Second,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep := &Report{
				Summary:            &ReportSummary{},
				IterationSummaries: tc.iters,
			}
			fillIterationRuntimeSummary(rep)
			require.NotNil(t, rep.Summary)
			assert.Equal(t, tc.wantMin, rep.Summary.IterationDurationMin)
			assert.Equal(t, tc.wantMax, rep.Summary.IterationDurationMax)
			assert.Equal(t, tc.wantP50, rep.Summary.IterationDurationP50)
		})
	}
}

func TestMarshalAISummaryJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		build func(t *testing.T) *Report
		check func(t *testing.T, raw []byte)
	}{
		{
			name: "nil_report",
			build: func(t *testing.T) *Report {
				return nil
			},
			check: func(t *testing.T, raw []byte) {
				assert.Equal(t, "null", string(raw))
			},
		},
		{
			name: "from_analyze_flake",
			build: func(t *testing.T) *Report {
				rep, _, err := Analyze(readers(
					`{"Action":"fail","Package":"p","Test":"T","Elapsed":0.1}`,
					`{"Action":"pass","Package":"p","Test":"T","Elapsed":0.1}`,
				), 30*time.Second)
				require.NoError(t, err)
				return rep
			},
			check: func(t *testing.T, raw []byte) {
				var sum ReportSummary
				require.NoError(t, json.Unmarshal(raw, &sum))
				assert.Equal(t, 1, sum.DistinctNamedTests)
				assert.Equal(t, 1, sum.FlakeNamedCount)
				assert.Equal(t, 1, sum.FlakeFailingIterations)
				assert.Equal(t, 2, sum.FlakeIterationTotal)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rep := tc.build(t)
			b, err := marshalAISummaryJSON(rep)
			require.NoError(t, err)
			tc.check(t, b)
		})
	}
}

func TestAnalyzeSlowTestsNoDuplication(t *testing.T) {
	t.Parallel()
	iter := `{"Action":"pass","Package":"pkg/slow","Test":"TestSlow","Elapsed":10.0}
{"Action":"pass","Package":"pkg/slow","Elapsed":10.0}
`
	rep, _, err := Analyze([]io.Reader{strings.NewReader(iter)}, 1*time.Second)
	require.NoError(t, err)

	pkgSlowCount := 0
	testSlowCount := 0
	for _, s := range rep.Slow {
		if s.Package == "pkg/slow" {
			switch s.Test {
			case "":
				pkgSlowCount++
			case "TestSlow":
				testSlowCount++
			}
		}
	}

	assert.Equal(t, 1, pkgSlowCount, "Package should appear exactly once in slow reports")
	assert.Equal(t, 1, testSlowCount, "Slow test should appear exactly once")
}
