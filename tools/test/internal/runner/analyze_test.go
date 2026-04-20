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
				{Package: "pkg/qux", Test: "TestHang", Fails: 1, MaxElapsed: 600 * time.Second, Iterations: []int{0}},
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
				{Package: "pkg/hang", Fails: 1, MaxElapsed: 120 * time.Second, Iterations: []int{0}},
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

func TestAnalyzeResultsRoundtrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// two iterations: same test flakes (fail, pass)
	must(t, os.WriteFile(filepath.Join(dir, "test-0.log.jsonl"),
		[]byte(`{"Action":"fail","Package":"pkg/z","Test":"TestFlaky","Elapsed":0.2}`+"\n"), 0600))
	must(t, os.WriteFile(filepath.Join(dir, "test-1.log.jsonl"),
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
