package runner

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteLogFiles(t *testing.T) {
	t.Parallel()

	iter := `{"Action":"output","Package":"github.com/foo/bar","Test":"TestFail","Output":"boom\n"}
{"Action":"fail","Package":"github.com/foo/bar","Test":"TestFail","Elapsed":0.1}
`
	dir := t.TempDir()
	rep, logs, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Failures, 1)

	require.NoError(t, WriteLogFiles(dir, rep, logs))

	// log_files path sanitization + existence
	require.Len(t, rep.Failures[0].LogFiles, 1)
	rel := rep.Failures[0].LogFiles[0]
	assert.Equal(t, "logs/foo_bar__TestFail__iter-0.log", rel)
	b, err := os.ReadFile(filepath.Join(dir, rel))
	require.NoError(t, err)
	assert.Equal(t, "boom\n", string(b))
}

func TestWriteLogFilesWritesEachIterationWithOutput(t *testing.T) {
	t.Parallel()

	// Iter 0 fails with output, iter 1 passes with output — both iterations
	// have captured output, so both get a log file for this flagged flake.
	iters := []string{
		`{"Action":"output","Package":"p","Test":"T","Output":"fail-log\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":0.01}
`,
		`{"Action":"output","Package":"p","Test":"T","Output":"ok-log\n"}
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}
`,
	}
	dir := t.TempDir()
	rep, logs, err := Analyze(readers(iters...), 30*time.Second)
	require.NoError(t, err)
	require.Len(t, rep.Flakes, 1)

	require.NoError(t, WriteLogFiles(dir, rep, logs))

	require.Len(t, rep.Flakes[0].LogFiles, 2)
	// Iteration order in report must be stable (sorted by iteration index).
	assert.Equal(t,
		[]string{"logs/p__T__iter-0.log", "logs/p__T__iter-1.log"},
		rep.Flakes[0].LogFiles)

	b0, err := os.ReadFile(filepath.Join(dir, "logs/p__T__iter-0.log"))
	require.NoError(t, err)
	assert.Equal(t, "fail-log\n", string(b0))

	b1, err := os.ReadFile(filepath.Join(dir, "logs/p__T__iter-1.log"))
	require.NoError(t, err)
	assert.Equal(t, "ok-log\n", string(b1))
}

func TestWriteLogFilesNoLogsForNonFlaggedTests(t *testing.T) {
	t.Parallel()

	// A clean pass is not flagged → no entry exists → no file written.
	iter := `{"Action":"output","Package":"p","Test":"T","Output":"hi\n"}
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}
`
	dir := t.TempDir()
	rep, logs, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	assert.Empty(t, rep.Flakes)
	assert.Empty(t, rep.Failures)
	assert.Empty(t, rep.Timeouts)

	require.NoError(t, WriteLogFiles(dir, rep, logs))

	entries, _ := os.ReadDir(filepath.Join(dir, "logs"))
	assert.Empty(t, entries, "no log files should be written for a clean-pass test")
}

func TestWriteCSV(t *testing.T) {
	t.Parallel()

	// Scenario: one flake, one failure, one timeout, one slow.
	iters := []string{
		// iter 0
		`{"Action":"fail","Package":"pkg/flake","Test":"TestFlake","Elapsed":0.1}
{"Action":"fail","Package":"pkg/fail","Test":"TestDead","Elapsed":0.2}
{"Action":"output","Package":"pkg/to","Test":"TestStuck","Output":"panic: test timed out after 5s\n"}
{"Action":"fail","Package":"pkg/to","Test":"TestStuck","Elapsed":5.0}
{"Action":"pass","Package":"pkg/slow","Test":"TestSlow","Elapsed":45.0}
`,
		// iter 1
		`{"Action":"pass","Package":"pkg/flake","Test":"TestFlake","Elapsed":0.08}
{"Action":"fail","Package":"pkg/fail","Test":"TestDead","Elapsed":0.25}
`,
	}
	dir := t.TempDir()
	rep, _, err := Analyze(readers(iters...), 30*time.Second)
	require.NoError(t, err)
	require.NoError(t, WriteCSV(dir, rep))

	f, err := os.Open(filepath.Join(dir, "report.csv"))
	require.NoError(t, err)
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(records), 5, "header + 4 rows")
	assert.Equal(t, []string{
		"package", "test", "category",
		"runs", "successes", "fails", "skips", "timeouts",
		"min", "max", "p50",
	}, records[0])

	// Worst-first: fails=2 (pkg/fail.TestDead) before timeouts=1 w/ fails=1 (pkg/to.TestStuck)
	// before fails=1 (pkg/flake.TestFlake) before slow (fails=0).
	rows := records[1:]
	categories := make([]string, 0, len(rows))
	for _, r := range rows {
		categories = append(categories, r[2])
	}
	// failure (fails=2) first
	assert.Equal(t, "failure", rows[0][2])
	assert.Equal(t, "pkg/fail", rows[0][0])
	// slow last
	assert.Equal(t, "slow", rows[len(rows)-1][2])
	// all four categories present
	assert.ElementsMatch(t, []string{"flake", "failure", "timeout", "slow"}, categories)
}

func TestWriteCSVRenamesSlowWhenAlsoTimeout(t *testing.T) {
	t.Parallel()
	// A test that's a timeout is also over the slow threshold. CSV must list
	// it once, as "timeout" not "slow" (primary signal wins).
	iter := `{"Action":"output","Package":"p","Test":"T","Output":"panic: test timed out after 10m0s\n"}
{"Action":"fail","Package":"p","Test":"T","Elapsed":600.0}
`
	dir := t.TempDir()
	rep, _, err := Analyze(readers(iter), 30*time.Second)
	require.NoError(t, err)
	require.NoError(t, WriteCSV(dir, rep))

	b, err := os.ReadFile(filepath.Join(dir, "report.csv"))
	require.NoError(t, err)
	content := string(b)
	assert.Contains(t, content, "timeout")
	// Only one data row beyond the header.
	assert.Equal(t, 2, strings.Count(content, "\n"), "header + one row")
}

func TestSanitize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"github.com/foo/bar", "github.com_foo_bar"},
		{"TestFoo/sub case", "TestFoo_sub_case"},
		{"TestName", "TestName"},
		{"a:b:c", "a_b_c"},
		{"", "_"},
		{"abc-123.go", "abc-123.go"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, sanitize(tc.in))
		})
	}
}
