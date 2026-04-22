package runner

import (
	"context"
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

	err := Diagnose(ctx, conf, "./...", nil, nil)
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

func TestBuildDiagnoseArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		conf        *config.App
		shuffleSeed int64
		want        []string
	}{
		{
			name: "default: timeout + target only",
			conf: &config.App{Timeout: 10 * time.Minute},
			want: []string{"test", "-json", "-count=1", "-timeout", "10m0s", "./pkg"},
		},
		{
			name: "with -race",
			conf: &config.App{Timeout: 5 * time.Minute, Race: true},
			want: []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-race", "./pkg"},
		},
		{
			name: "with -run regexp",
			conf: &config.App{Timeout: 5 * time.Minute, Run: "^TestFoo$"},
			want: []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-run=^TestFoo$", "./pkg"},
		},
		{
			name: "with -race and -run together",
			conf: &config.App{Timeout: 5 * time.Minute, Race: true, Run: "TestFoo/bar"},
			want: []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-race", "-run=TestFoo/bar", "./pkg"},
		},
		{
			name: "empty Run is omitted (not `-run=`)",
			conf: &config.App{Timeout: 1 * time.Minute, Run: ""},
			want: []string{"test", "-json", "-count=1", "-timeout", "1m0s", "./pkg"},
		},
		{
			name: "with -cpu list",
			conf: &config.App{Timeout: 5 * time.Minute, CPU: "1,2,4"},
			want: []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-cpu=1,2,4", "./pkg"},
		},
		{
			name: "with -parallel",
			conf: &config.App{Timeout: 5 * time.Minute, Parallel: 4},
			want: []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-parallel=4", "./pkg"},
		},
		{
			name:        "with shuffle seed",
			conf:        &config.App{Timeout: 5 * time.Minute},
			shuffleSeed: 12345,
			want:        []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-shuffle=12345", "./pkg"},
		},
		{
			name:        "zero shuffle seed omitted",
			conf:        &config.App{Timeout: 5 * time.Minute},
			shuffleSeed: 0,
			want:        []string{"test", "-json", "-count=1", "-timeout", "5m0s", "./pkg"},
		},
		{
			name: "zero parallel omitted",
			conf: &config.App{Timeout: 1 * time.Minute, Parallel: 0},
			want: []string{"test", "-json", "-count=1", "-timeout", "1m0s", "./pkg"},
		},
		{
			name: "empty cpu omitted",
			conf: &config.App{Timeout: 1 * time.Minute, CPU: ""},
			want: []string{"test", "-json", "-count=1", "-timeout", "1m0s", "./pkg"},
		},
		{
			name:        "all flags together",
			conf:        &config.App{Timeout: 5 * time.Minute, Race: true, Run: "^TestFoo$", CPU: "1,2", Parallel: 8},
			shuffleSeed: 999,
			want:        []string{"test", "-json", "-count=1", "-timeout", "5m0s", "-race", "-run=^TestFoo$", "-cpu=1,2", "-parallel=8", "-shuffle=999", "./pkg"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := buildDiagnoseArgs(tc.conf, "./pkg", tc.shuffleSeed)
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

	require.NoError(t, Diagnose(ctx, conf, "./...", nil, nil))

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
		name   string
		conf   *config.App
		target string
		want   string
	}{
		{
			name: "repo root pattern",
			conf: &config.App{
				Iterations: 1,
				Timeout:    10 * time.Minute,
			},
			target: "./...",
			want:   diagnoseResultsNamePrefix + "allpkgs-it1-to10m0s-20240601123045",
		},
		{
			name: "nested package with ellipsis",
			conf: &config.App{
				Iterations: 10,
				Timeout:    10 * time.Minute,
			},
			target: "./core/...",
			want:   diagnoseResultsNamePrefix + "core_allpkgs-it10-to10m0s-20240601123045",
		},
		{
			name: "flags and run regexp",
			conf: &config.App{
				Iterations: 2,
				Timeout:    5 * time.Minute,
				Race:       true,
				Shuffle:    true,
				FailFast:   true,
				Parallel:   8,
				CPU:        "1,2",
				Run:        "^TestFoo$",
			},
			target: "./pkg",
			want:   diagnoseResultsNamePrefix + "pkg-it2-to5m0s-race-shuffle-ff-p8-cpu-1-2-r_TestFoo_-383ceba4-20240601123045",
		},
		{
			name: "non default slow threshold",
			conf: &config.App{
				Iterations:    1,
				Timeout:       10 * time.Minute,
				SlowThreshold: 45 * time.Second,
			},
			target: "./core/services/foo",
			want:   diagnoseResultsNamePrefix + "core_services_foo-it1-to10m0s-slow45s-20240601123045",
		},
		{
			name: "default slow threshold omitted",
			conf: &config.App{
				Iterations:    3,
				Timeout:       10 * time.Minute,
				SlowThreshold: 30 * time.Second,
			},
			target: "./a",
			want:   diagnoseResultsNamePrefix + "a-it3-to10m0s-20240601123045",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := diagnoseResultsDirName(tc.conf, tc.target, diagnoseResultsDirNameAt)
			assert.Equal(t, tc.want, got)
			assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
		})
	}
}

func TestDiagnoseResultsDirNameLongRunAndPath(t *testing.T) {
	t.Parallel()
	longRun := strings.Repeat("Xy", 80) // 160 chars; token truncates to 40 runes
	conf := &config.App{
		Iterations: 1,
		Timeout:    10 * time.Minute,
		Run:        longRun,
	}
	got := diagnoseResultsDirName(conf, "./p", diagnoseResultsDirNameAt)
	assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
	assert.Contains(t, got, "-it1-to10m0s-")
	assert.Regexp(t, `r(?:Xy){20}-[0-9a-f]{8}-20240601123045`, got)

	longTarget := "./" + strings.Repeat("seg/", 60) + "z"
	got2 := diagnoseResultsDirName(conf, longTarget, diagnoseResultsDirNameAt)
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
	require.NoError(t, Diagnose(ctx, conf, "./...", nil, dumpDB))
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
}
