package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

// Fixed clock for diagnoseResultsDirName assertions (timestamp 20240601123045).
var diagnoseResultsDirNameAt = time.Date(2024, 6, 1, 12, 30, 45, 0, time.UTC)

func testArgHash8(goTestArgs []string) string {
	h := sha256.Sum256([]byte(strings.Join(goTestArgs, "\x00")))
	return hex.EncodeToString(h[:4])
}

// When ctx is already canceled before Diagnose starts, no iterations run but
// analysis still produces a report.json — this is the path a user hits after
// Ctrl+C'ing a long-running diagnose run.
func TestDiagnoseCanceledCtxRunsNoIterationsButStillWritesReport(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	conf := &config.App{
		RepoRoot:   repoRoot,
		AIOutput:   true,
		Iterations: 3,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Diagnose(ctx, conf, []string{"./..."}, nil, nil)
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(repoRoot, diagnoseResultsNamePrefix+"*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
	resultsDir := matches[0]

	// No iteration jsonl files because the loop guard tripped on entry.
	iterFiles, err := filepath.Glob(filepath.Join(resultsDir, "iteration-*.log.jsonl"))
	require.NoError(t, err)
	assert.Empty(t, iterFiles)

	reportBytes, err := os.ReadFile(filepath.Join(resultsDir, "report.json"))
	require.NoError(t, err)
	var rep Report
	require.NoError(t, json.Unmarshal(reportBytes, &rep))
	assert.Equal(t, 0, rep.Iterations)
}

func TestParseDiagnoseGoTestCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantSet bool
		wantN   int
		wantErr bool
	}{
		{name: "no count", args: []string{"-v", "./..."}, wantSet: false, wantN: 0, wantErr: false},
		{name: "count 1", args: []string{"-count=1", "./..."}, wantSet: true, wantN: 1, wantErr: false},
		{name: "count 1 spaced", args: []string{"-count", "1", "./..."}, wantSet: true, wantN: 1, wantErr: false},
		{name: "count 2", args: []string{"-count=2", "./..."}, wantSet: true, wantN: 2, wantErr: false},
		{name: "count 99", args: []string{"-count", "99"}, wantSet: true, wantN: 99, wantErr: false},
		{name: "last count wins", args: []string{"-count=1", "-count=3"}, wantSet: true, wantN: 3, wantErr: false},
		{name: "count after -args ignored", args: []string{"-v", "-args", "-count=50"}, wantSet: false, wantN: 0, wantErr: false},
		{name: "invalid count value", args: []string{"-count=maybe"}, wantErr: true},
		{name: "-count without value", args: []string{"-count"}, wantErr: true},
		{name: "count zero", args: []string{"-count=0", "./..."}, wantErr: true},
		{name: "count negative", args: []string{"-count=-1", "./..."}, wantErr: true},
		{name: "count zero spaced", args: []string{"-count", "0"}, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			set, n, err := parseDiagnoseGoTestCount(tc.args)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantSet, set)
			assert.Equal(t, tc.wantN, n)
		})
	}
}

func TestWarnDiagnoseGoTestCount(t *testing.T) {
	t.Parallel()

	t.Run("count 1", func(t *testing.T) {
		t.Parallel()
		var buf strings.Builder
		require.NoError(t, WarnDiagnoseGoTestCount(&buf, []string{"-count=1", "./pkg"}))
		assert.Contains(t, buf.String(), "unnecessary")
	})

	t.Run("count greater than 1", func(t *testing.T) {
		t.Parallel()
		var buf strings.Builder
		require.NoError(t, WarnDiagnoseGoTestCount(&buf, []string{"-count=5"}))
		assert.Contains(t, buf.String(), "prefer")
		assert.Contains(t, buf.String(), "iterations")
	})

	t.Run("no count", func(t *testing.T) {
		t.Parallel()
		var buf strings.Builder
		require.NoError(t, WarnDiagnoseGoTestCount(&buf, []string{"./..."}))
		assert.Empty(t, strings.TrimSpace(buf.String()))
	})

	t.Run("invalid non positive count", func(t *testing.T) {
		t.Parallel()
		var buf strings.Builder
		err := WarnDiagnoseGoTestCount(&buf, []string{"-count=0", "./..."})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "positive integer")
	})
}

func TestBuildDiagnoseArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		goTestArgs  []string
		shuffleSeed int64
		want        []string
		wantErr     bool
	}{
		{
			name:       "passthrough flags and package",
			goTestArgs: []string{"-timeout=5m", "./pkg"},
			want:       []string{"test", "-json", "-timeout=5m", "./pkg", "-count=1"},
		},
		{
			name:        "shuffle seed appended",
			goTestArgs:  []string{"./pkg"},
			shuffleSeed: 12345,
			want:        []string{"test", "-json", "./pkg", "-shuffle=12345", "-count=1"},
		},
		{
			name:        "zero shuffle seed omitted",
			goTestArgs:  []string{"./pkg"},
			shuffleSeed: 0,
			want:        []string{"test", "-json", "./pkg", "-count=1"},
		},
		{
			name:       "strips duplicate -json; keeps -count greater than 1",
			goTestArgs: []string{"-json", "-count=3", "-race", "-run=^X$", "./pkg"},
			want:       []string{"test", "-json", "-count=3", "-race", "-run=^X$", "./pkg"},
		},
		{
			name:       "passes through -count with separate value when greater than 1",
			goTestArgs: []string{"-count", "99", "./a"},
			want:       []string{"test", "-json", "-count", "99", "./a"},
		},
		{
			name:       "explicit -count=1 gets default appended",
			goTestArgs: []string{"-count=1", "./pkg"},
			want:       []string{"test", "-json", "-count=1", "./pkg", "-count=1"},
		},
		{
			name:       "reject count zero",
			goTestArgs: []string{"-count=0", "./pkg"},
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := buildDiagnoseArgs(tc.goTestArgs, tc.shuffleSeed)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDiagnoseShuffleSeedsAbsentWhenNoIterationsRun(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	conf := &config.App{
		RepoRoot:   repoRoot,
		AIOutput:   true,
		Iterations: 3,
		Shuffle:    true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.NoError(t, Diagnose(ctx, conf, []string{"./..."}, nil, nil))

	matches, err := filepath.Glob(filepath.Join(repoRoot, diagnoseResultsNamePrefix+"*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	reportBytes, err := os.ReadFile(filepath.Join(matches[0], "report.json"))
	require.NoError(t, err)
	var rep Report
	require.NoError(t, json.Unmarshal(reportBytes, &rep))
	assert.Empty(t, rep.IterationSummaries)
}

func TestDiagnoseResultsDirName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		conf       *config.App
		goTestArgs []string
		want       string
	}{
		{
			name: "repo root pattern",
			conf: &config.App{
				Iterations: 1,
			},
			goTestArgs: []string{"./..."},
			want:       diagnoseResultsNamePrefix + "allpkgs-it1-h" + testArgHash8([]string{"./..."}) + "-20240601123045",
		},
		{
			name: "nested package with ellipsis",
			conf: &config.App{
				Iterations: 10,
			},
			goTestArgs: []string{"./core/..."},
			want:       diagnoseResultsNamePrefix + "core_allpkgs-it10-h" + testArgHash8([]string{"./core/..."}) + "-20240601123045",
		},
		{
			name: "fail-fast and shuffle and non-default slow",
			conf: &config.App{
				Iterations:    2,
				SlowThreshold: 45 * time.Second,
				FailFast:      true,
				Shuffle:       true,
			},
			goTestArgs: []string{"-race", "-run=^TestFoo$", "./pkg"},
			want: diagnoseResultsNamePrefix + "pkg-it2-h" + testArgHash8([]string{"-race", "-run=^TestFoo$", "./pkg"}) +
				"-ff-shuffle-slow45s-20240601123045",
		},
		{
			name: "default slow threshold omitted",
			conf: &config.App{
				Iterations:    3,
				SlowThreshold: 30 * time.Second,
			},
			goTestArgs: []string{"./a"},
			want:       diagnoseResultsNamePrefix + "a-it3-h" + testArgHash8([]string{"./a"}) + "-20240601123045",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := diagnoseResultsDirName(tc.conf, tc.goTestArgs, diagnoseResultsDirNameAt)
			assert.Equal(t, tc.want, got)
			assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
		})
	}
}

func TestDiagnoseResultsDirNameLongRunAndPath(t *testing.T) {
	t.Parallel()
	longRun := strings.Repeat("Xy", 80)
	goTestArgs := []string{"-run=" + longRun, "./p"}
	conf := &config.App{
		Iterations: 1,
	}
	got := diagnoseResultsDirName(conf, goTestArgs, diagnoseResultsDirNameAt)
	assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
	assert.Contains(t, got, "-it1-h")
	assert.Contains(t, got, testArgHash8(goTestArgs))
	assert.Regexp(t, `diagnose-p-it1-h[0-9a-f]{8}-20240601123045`, got)

	longTarget := "./" + strings.Repeat("seg/", 60) + "z"
	goTestArgs2 := []string{longTarget}
	got2 := diagnoseResultsDirName(conf, goTestArgs2, diagnoseResultsDirNameAt)
	assert.LessOrEqual(t, len(got2), maxDiagnoseResultsBasename)
	assert.True(t, strings.HasPrefix(got2, diagnoseResultsNamePrefix))
}

func TestDiagnoseDumpDBCalledWithResultsDir(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	conf := &config.App{
		RepoRoot:   repoRoot,
		AIOutput:   true,
		Iterations: 2,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type call struct {
		dir  string
		iter int
	}
	var calls []call
	dumpDB := func(_ context.Context, dir string, iteration int) error {
		calls = append(calls, call{dir, iteration})
		return nil
	}

	// pre-canceled ctx → no iterations run → dumpDB never called
	require.NoError(t, Diagnose(ctx, conf, []string{"./..."}, nil, dumpDB))
	assert.Empty(t, calls)
}

func TestTruncateUTF8MaxBytes(t *testing.T) {
	t.Parallel()
	s := "ééé" // 6 bytes, 3 runes
	assert.Empty(t, truncateUTF8MaxBytes(s, 0))
	assert.Empty(t, truncateUTF8MaxBytes(s, 1))
	assert.Equal(t, "é", truncateUTF8MaxBytes(s, 2))
	assert.Equal(t, "éé", truncateUTF8MaxBytes(s, 4))
	assert.Equal(t, "ééé", truncateUTF8MaxBytes(s, 6))
	assert.Equal(t, "ééé", truncateUTF8MaxBytes(s, 10))
	// U+FFFD is utf8.RuneError's value; truncation must not strip a valid final replacement character.
	assert.Equal(t, "abc\uFFFD", truncateUTF8MaxBytes("abc\uFFFD"+"x", 6))
}

func TestPackagePatternsFromEnd(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"./core/...", "./foo"}, packagePatternsFromEnd([]string{"-race", "-timeout=5m", "./core/...", "./foo"}))
	assert.Nil(t, packagePatternsFromEnd([]string{"-v", "-race"}))
}
