package runner

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
)

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

	err := Diagnose(ctx, conf, output.New(conf.AIOutput, io.Discard, io.Discard, output.SkipFD), []string{"./..."}, nil)
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
	require.NotNil(t, rep.Run)
	assert.Equal(t, []string{"./..."}, rep.Run.GoTestArgs)
	assert.Equal(t, "allpkgs", rep.Run.TargetSlug)
	require.NotNil(t, rep.Run.FinishedAt)
}

func TestDiagnoseCanceledCtxAIStdoutTwoLines(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	conf := &config.App{
		RepoRoot:   repoRoot,
		AIOutput:   true,
		Iterations: 3,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stdout strings.Builder
	err := Diagnose(ctx, conf, output.New(conf.AIOutput, &stdout, io.Discard, output.SkipFD), []string{"./..."}, nil)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	require.Len(t, lines, 3, "stdout: %q", stdout.String())
	assert.Equal(t, filepath.Join(lines[0], "report.json"), lines[1])
	assert.Equal(t, "null", lines[2])
}

func TestDiagnoseHumanModeFooterShowsReportJSONPath(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	conf := &config.App{
		RepoRoot:   repoRoot,
		AIOutput:   false,
		Iterations: 2,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stderr strings.Builder
	err := Diagnose(ctx, conf, output.New(false, io.Discard, &stderr, output.SkipFD), []string{"./..."}, nil)
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(repoRoot, diagnoseResultsNamePrefix+"*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
	reportPath := filepath.Join(matches[0], "report.json")

	out := stderr.String()
	assert.Contains(t, out, "results directory")
	assert.Contains(t, out, matches[0])
	assert.Contains(t, out, reportPath)
	assert.Contains(t, out, "report.json:")
	assert.NotContains(t, out, "results in ")
}

func TestPrintDiagnoseAnalyzingStartsNewLineAfterLiveProgress(t *testing.T) {
	t.Parallel()
	var stderr strings.Builder
	out := output.New(false, io.Discard, &stderr, output.SkipFD)

	printDiagnoseAnalyzing(out, true)

	assert.Equal(t, "\r\u001b[K\nanalyzing...\n", stderr.String())
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

	require.NoError(t, Diagnose(ctx, conf, output.New(conf.AIOutput, io.Discard, io.Discard, output.SkipFD), []string{"./..."}, nil))

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
		goTestArgs []string
		want       string
	}{
		{
			name:       "repo root pattern",
			goTestArgs: []string{"./..."},
			want:       diagnoseResultsNamePrefix + "allpkgs-20240601123045",
		},
		{
			name:       "nested package with ellipsis",
			goTestArgs: []string{"./core/..."},
			want:       diagnoseResultsNamePrefix + "core_allpkgs-20240601123045",
		},
		{
			name:       "flags before package",
			goTestArgs: []string{"-race", "-run=^TestFoo$", "./pkg"},
			want:       diagnoseResultsNamePrefix + "pkg-20240601123045",
		},
		{
			name:       "single package",
			goTestArgs: []string{"./pkg"},
			want:       diagnoseResultsNamePrefix + "pkg-20240601123045",
		},
		{
			name:       "short path",
			goTestArgs: []string{"./a"},
			want:       diagnoseResultsNamePrefix + "a-20240601123045",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := diagnoseResultsDirName(tc.goTestArgs, diagnoseResultsDirNameAt)
			assert.Equal(t, tc.want, got)
			assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
		})
	}
}

func TestDiagnoseResultsDirNameLongRunAndPath(t *testing.T) {
	t.Parallel()
	longRun := strings.Repeat("Xy", 80)
	goTestArgs := []string{"-run=" + longRun, "./p"}
	got := diagnoseResultsDirName(goTestArgs, diagnoseResultsDirNameAt)
	assert.LessOrEqual(t, len(got), maxDiagnoseResultsBasename)
	assert.Regexp(t, `diagnose-p-20240601123045`, got)

	longTarget := "./" + strings.Repeat("seg/", 60) + "z"
	goTestArgs2 := []string{longTarget}
	got2 := diagnoseResultsDirName(goTestArgs2, diagnoseResultsDirNameAt)
	assert.LessOrEqual(t, len(got2), maxDiagnoseResultsBasename)
	assert.True(t, strings.HasPrefix(got2, diagnoseResultsNamePrefix))
}

func TestMakeDiagnoseResultsDirAvoidsExistingDirectory(t *testing.T) {
	t.Parallel()
	conf := &config.App{
		RepoRoot:   t.TempDir(),
		Iterations: 1,
	}
	first, err := makeDiagnoseResultsDir(conf, []string{"./pkg"}, diagnoseResultsDirNameAt)
	require.NoError(t, err)
	second, err := makeDiagnoseResultsDir(conf, []string{"./pkg"}, diagnoseResultsDirNameAt)
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
	assert.DirExists(t, first)
	assert.DirExists(t, second)
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
	require.NoError(t, Diagnose(ctx, conf, output.New(conf.AIOutput, io.Discard, io.Discard, output.SkipFD), []string{"./..."}, []diagnoseIterationResource{{DumpDiagnostics: dumpDB}}))
	assert.Empty(t, calls)
}

func TestPrintDiagnoseResultsDirHeader(t *testing.T) {
	t.Parallel()
	dir := "/tmp/diagnose-out-xyz"

	t.Run("human", func(t *testing.T) {
		t.Parallel()
		var stdout, stderr strings.Builder
		p := output.New(false, &stdout, &stderr, output.SkipFD)
		printDiagnoseResultsDirHeader(p, dir)
		assert.Empty(t, stdout.String())
		assert.Contains(t, stderr.String(), "results directory")
		assert.Contains(t, stderr.String(), dir)
	})

	t.Run("ai", func(t *testing.T) {
		t.Parallel()
		var stdout, stderr strings.Builder
		p := output.New(true, &stdout, &stderr, output.SkipFD)
		printDiagnoseResultsDirHeader(p, dir)
		assert.Contains(t, stdout.String(), dir)
		assert.Empty(t, stderr.String())
	})
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

func TestRunDiagnoseIterationsRunsInParallelWithWorkerIsolation(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	resultsDir := t.TempDir()
	conf := &config.App{
		RepoRoot:           repoRoot,
		AIOutput:           true,
		Iterations:         4,
		ParallelIterations: 2,
	}
	var stdout strings.Builder
	out := output.New(true, &stdout, io.Discard, output.SkipFD)

	var active int32
	var maxActive int32
	var mu sync.Mutex
	envByIter := make(map[int][]string)
	var resets, dumps int
	resources := []diagnoseIterationResource{
		{
			Env: []string{"CL_DATABASE_URL=postgres://worker-0/db"},
			Reset: func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				resets++
				return nil
			},
			DumpDiagnostics: func(_ context.Context, _ string, _ int) error {
				mu.Lock()
				defer mu.Unlock()
				dumps++
				return nil
			},
		},
		{
			Env: []string{"CL_DATABASE_URL=postgres://worker-1/db"},
			Reset: func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				resets++
				return nil
			},
			DumpDiagnostics: func(_ context.Context, _ string, _ int) error {
				mu.Lock()
				defer mu.Unlock()
				dumps++
				return nil
			},
		},
	}
	hooks := diagnoseRunHooks{
		runIteration: func(_ context.Context, _ *config.App, _ *output.Printer, dir string, _ []string, iteration int, _ int64, env []string, liveProgress bool, parallelProgress *parallelDiagnoseProgress, _ time.Time) error {
			require.False(t, liveProgress)
			require.Nil(t, parallelProgress)
			nowActive := atomic.AddInt32(&active, 1)
			for {
				seen := atomic.LoadInt32(&maxActive)
				if nowActive <= seen || atomic.CompareAndSwapInt32(&maxActive, seen, nowActive) {
					break
				}
			}
			mu.Lock()
			envByIter[iteration] = append([]string(nil), env...)
			mu.Unlock()
			time.Sleep(25 * time.Millisecond)
			err := os.WriteFile(filepath.Join(dir, "iteration-"+strconv.Itoa(iteration)+".log.jsonl"), []byte(`{"Action":"pass","Package":"p","Test":"T","Elapsed":0.01}`+"\n"), 0600)
			atomic.AddInt32(&active, -1)
			return err
		},
	}

	state, err := runDiagnoseIterations(context.Background(), conf, out, resultsDir, []string{"./pkg"}, resources, hooks)
	require.NoError(t, err)
	assert.Equal(t, 4, state.completed)
	assert.Equal(t, int32(2), atomic.LoadInt32(&maxActive))
	assert.Equal(t, 2, resets, "each worker should reset before being reused")
	assert.Equal(t, 4, dumps)
	assert.Len(t, envByIter, 4)
	for _, env := range envByIter {
		assert.Len(t, env, 1)
		assert.Contains(t, env[0], "CL_DATABASE_URL=postgres://worker-")
	}
	assert.Equal(t, 4, strings.Count(stdout.String(), "d "))
}

func TestRunDiagnoseIterationsFailFastCancelsNewWork(t *testing.T) {
	t.Parallel()
	resultsDir := t.TempDir()
	conf := &config.App{
		RepoRoot:           t.TempDir(),
		AIOutput:           true,
		Iterations:         5,
		ParallelIterations: 2,
		FailFast:           true,
	}
	out := output.New(true, io.Discard, io.Discard, output.SkipFD)
	var mu sync.Mutex
	started := make(map[int]struct{})
	hooks := diagnoseRunHooks{
		runIteration: func(ctx context.Context, _ *config.App, _ *output.Printer, dir string, _ []string, iteration int, _ int64, _ []string, _ bool, _ *parallelDiagnoseProgress, _ time.Time) error {
			mu.Lock()
			started[iteration] = struct{}{}
			mu.Unlock()
			if iteration == 0 {
				require.NoError(t, os.WriteFile(filepath.Join(dir, "iteration-0.log.jsonl"), []byte(`{"Action":"fail","Package":"p","Test":"T","Elapsed":0.01}`+"\n"), 0600))
				return errors.New("test failed")
			}
			<-ctx.Done()
			return ctx.Err()
		},
	}

	state, err := runDiagnoseIterations(context.Background(), conf, out, resultsDir, []string{"./pkg"}, []diagnoseIterationResource{{}, {}}, hooks)
	require.NoError(t, err)
	assert.True(t, state.failedFast)
	assert.LessOrEqual(t, len(started), conf.ParallelIterations)
}

func TestRunDiagnoseIterationsFailFastOnCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		failFastOn    []string
		iterationJSON string
		iterErr       error
		wantCompleted int
		wantFailed    bool
	}{
		{
			name:          "timeout stops on timeout",
			failFastOn:    []string{"timeout"},
			iterationJSON: `{"Action":"output","Package":"p","Test":"TestHang","Output":"panic: test timed out after 5s\n"}` + "\n" + `{"Action":"fail","Package":"p","Test":"TestHang","Elapsed":5}`,
			iterErr:       errors.New("test failed"),
			wantCompleted: 1,
			wantFailed:    true,
		},
		{
			name:          "timeout ignores plain failure",
			failFastOn:    []string{"timeout"},
			iterationJSON: `{"Action":"fail","Package":"p","Test":"TestFail","Elapsed":0.01}`,
			iterErr:       errors.New("test failed"),
			wantCompleted: 3,
		},
		{
			name:          "slow stops on passing slow test",
			failFastOn:    []string{"slow"},
			iterationJSON: `{"Action":"pass","Package":"p","Test":"TestSlow","Elapsed":45}`,
			wantCompleted: 1,
			wantFailed:    true,
		},
		{
			name:          "failure stops on plain failure",
			failFastOn:    []string{"failure"},
			iterationJSON: `{"Action":"fail","Package":"p","Test":"TestFail","Elapsed":0.01}`,
			iterErr:       errors.New("test failed"),
			wantCompleted: 1,
			wantFailed:    true,
		},
		{
			name:          "any stops on slow",
			failFastOn:    []string{"any"},
			iterationJSON: `{"Action":"pass","Package":"p","Test":"TestSlow","Elapsed":45}`,
			wantCompleted: 1,
			wantFailed:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resultsDir := t.TempDir()
			conf := &config.App{
				RepoRoot:      t.TempDir(),
				AIOutput:      true,
				Iterations:    3,
				SlowThreshold: 30 * time.Second,
				FailFastOn:    tc.failFastOn,
			}
			out := output.New(true, io.Discard, io.Discard, output.SkipFD)
			hooks := diagnoseRunHooks{
				runIteration: func(_ context.Context, _ *config.App, _ *output.Printer, dir string, _ []string, iteration int, _ int64, _ []string, _ bool, _ *parallelDiagnoseProgress, _ time.Time) error {
					return os.WriteFile(filepath.Join(dir, "iteration-"+strconv.Itoa(iteration)+".log.jsonl"), []byte(tc.iterationJSON+"\n"), 0600)
				},
			}
			if tc.iterErr != nil {
				hooks.runIteration = func(_ context.Context, _ *config.App, _ *output.Printer, dir string, _ []string, iteration int, _ int64, _ []string, _ bool, _ *parallelDiagnoseProgress, _ time.Time) error {
					require.NoError(t, os.WriteFile(filepath.Join(dir, "iteration-"+strconv.Itoa(iteration)+".log.jsonl"), []byte(tc.iterationJSON+"\n"), 0600))
					return tc.iterErr
				}
			}

			state, err := runDiagnoseIterations(context.Background(), conf, out, resultsDir, []string{"./pkg"}, []diagnoseIterationResource{{}}, hooks)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCompleted, state.completed)
			assert.Equal(t, tc.wantFailed, state.failedFast)
		})
	}
}

func TestFormatIterationDigestAI(t *testing.T) {
	t.Parallel()
	d := IterationDigest{
		Result: "pass", RanTests: 126, FailTests: 0, TimeoutTests: 0, SlowTests: 6,
	}
	assert.Equal(t, "d 7/100 p 90s r126 f0 t0 s6", formatIterationDigestAI(7, 100, d, 90*time.Second))
}
