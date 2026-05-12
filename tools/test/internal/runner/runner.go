package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

// failFastReasonBuildFailure is used when go test reports a compile/build failure
// (FailedBuild or package-level fail with no named tests run); diagnose stops immediately.
const failFastReasonBuildFailure = "build-failure"

type diagnoseIterationResource = db.Resource

// diagnoseIterationParams is per-iteration state for diagnose (see diagnoseIterationRunner;
// context is passed separately to avoid storing context in a struct).
type diagnoseIterationParams struct {
	Conf             *config.App
	Out              *output.Printer
	ResultsDir       string
	GoTestArgs       []string
	Iteration        int
	ShuffleSeed      int64
	Env              []string
	LiveProgress     bool
	ParallelProgress *parallelDiagnoseProgress
	DiagnoseRunStart time.Time
	SerialProgressMu *sync.Mutex
}

type diagnoseIterationRunner func(ctx context.Context, p diagnoseIterationParams) error

type diagnoseRunHooks struct {
	runIteration diagnoseIterationRunner
	seed         func() int64
}

type diagnoseRunState struct {
	completed           int
	failedFast          bool
	failedFastReason    string
	failedFastIteration int // 0-based diagnose iteration index; -1 if unset
	iterDurations       []time.Duration
	shuffleSeeds        map[int]int64
	liveProgress        bool
}

// GoTest runs `go test` with the given args (repo root as working directory).
func GoTest(ctx context.Context, conf *config.App, args []string) error {
	//nolint:gosec // it's fine
	cmd := exec.CommandContext(ctx, "go", append([]string{"test"}, args...)...)
	cmd.Dir = conf.RepoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// Gotestsum runs `gotestsum` with the given args (repo root as working directory).
func Gotestsum(ctx context.Context, conf *config.App, args []string) error {
	if _, err := exec.LookPath("gotestsum"); err != nil {
		return fmt.Errorf("gotestsum not on PATH: install with go install gotest.tools/gotestsum@latest: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gotestsum", args...)
	cmd.Dir = conf.RepoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// Diagnose runs go test -json once per iteration, writing each stream to
// iteration-<n>.log.jsonl, then analyzes and writes report.json.
// With --ai-output, stdout is three lines: the results directory path, the
// path to report.json, and one line of JSON (the report's summary object, or
// the JSON keyword null when there is no summary).
// Test iteration failures do not stop later runs (unless --fail-fast); compile/build
// failures stop immediately. Results are reflected in report.json. Diagnose returns a non-nil error for setup failures
// (e.g. mkdir, database reset), analyze/write report failures, or ctx errors
// bubbling from dependencies — not for failing tests alone.
// resources supplies the prepared per-worker database state. Each resource runs
// one iteration at a time; when a worker is reused, its Reset hook restores the
// database to its freshly-prepared state.
func Diagnose(ctx context.Context, conf *config.App, out *output.Printer, goTestArgs []string, resources []diagnoseIterationResource) error {
	if out == nil {
		out = output.NewFromApp(conf)
	}
	start := time.Now()

	resultsDir, err := makeDiagnoseResultsDir(conf, goTestArgs, start)
	if err != nil {
		return err
	}
	printDiagnoseResultsDirHeader(out, resultsDir)
	if err := printDiagnoseRunTimeEstimate(out, conf, goTestArgs, len(resources)); err != nil {
		return err
	}

	state, runErr := runDiagnoseIterations(ctx, conf, out, resultsDir, goTestArgs, resources, diagnoseRunHooks{})
	if runErr != nil {
		if ctx.Err() == nil {
			return runErr
		}
	}

	interrupted := ctx.Err() != nil
	if interrupted && !out.AIOutput() {
		out.HumanStderr(
			termstyle.Accent.Render(fmt.Sprintf("interrupted after %d/%d iterations", state.completed, conf.Iterations)) +
				termstyle.Muted.Render(" — analyzing partial results…"))
	}

	if state.failedFast && !out.AIOutput() {
		if state.failedFastReason == failFastReasonBuildFailure && state.failedFastIteration >= 0 {
			iter := state.failedFastIteration + 1
			out.HumanStderr(termstyle.Bad.Render(
				fmt.Sprintf("Build failed — stopping diagnose run (iteration %d/%d).", iter, conf.Iterations)))
		} else {
			msg := "--fail-fast set, stopping early"
			if state.failedFastReason != "" {
				msg = fmt.Sprintf("fail-fast matched %s, stopping early", state.failedFastReason)
			}
			out.HumanStderr(termstyle.Accent.Render(msg))
		}
	}

	stopAnalyzing := startDiagnoseAnalyzingProgress(out, state.liveProgress)
	var report *Report
	var logs LogMap
	var analyzeErr error
	defer func() { stopAnalyzing(analyzeErr) }()
	report, logs, analyzeErr = AnalyzeResults(resultsDir, conf.SlowThreshold)
	if analyzeErr != nil {
		out.Stderrf("analyze results: %v\n", analyzeErr)
		return analyzeErr
	}
	if out.AIOutput() && state.failedFastReason == failFastReasonBuildFailure && state.failedFastIteration >= 0 {
		var pkgs []string
		if report != nil {
			for _, sum := range report.IterationSummaries {
				if sum.Index == state.failedFastIteration {
					pkgs = append(pkgs, sum.FailingTests...)
					break
				}
			}
		}
		out.Stderrf("bf_stop iter=%d pkgs=%s\n", state.failedFastIteration+1, strings.Join(pkgs, ","))
	}
	if report != nil {
		for i, d := range state.iterDurations {
			if i >= len(report.IterationSummaries) {
				break
			}
			report.IterationSummaries[i].Duration = d
			if state.shuffleSeeds != nil {
				report.IterationSummaries[i].ShuffleSeed = state.shuffleSeeds[i]
			}
		}
		finished := time.Now()
		report.Run = newRunMeta(conf, goTestArgs, resultsDir, start, &finished, len(resources) > 0)
		fillIterationRuntimeSummary(report)
	}
	if err := WriteLogFiles(resultsDir, report, logs); err != nil {
		out.Stderrf("write log files: %v\n", err)
		return err
	}
	if err := WriteReport(resultsDir, report); err != nil {
		out.Stderrf("write report: %v\n", err)
		return err
	}
	if err := WriteCSV(resultsDir, report); err != nil {
		out.Stderrf("write csv: %v\n", err)
		return err
	}

	reportPath := filepath.Join(resultsDir, "report.json")
	if out.AIOutput() {
		out.SparseStdoutln(reportPath)
		summaryJSON, err := marshalAISummaryJSON(report)
		if err != nil {
			out.Stderrf("marshal ai summary: %v\n", err)
			return err
		}
		out.SparseStdoutln(string(summaryJSON))
		return nil
	}

	out.HumanStderr(
		termstyle.Label.Render("diagnosis complete") + " " +
			termstyle.Muted.Render("["+formatDiagnoseWallClock(time.Since(start))+"]"))
	if report != nil {
		PrintSummary(out.HumanStderrWriter(), report)
	}
	out.HumanStderr(termstyle.Muted.Render("report.json: ") + termstyle.Label.Render(reportPath))
	return nil
}

// fillIterationRuntimeSummary sets summary iteration_duration_* from each
// iteration's wall-clock Duration (min / max / p50 across the run).
func fillIterationRuntimeSummary(rep *Report) {
	if rep == nil || rep.Summary == nil {
		return
	}
	var samples []time.Duration
	for _, s := range rep.IterationSummaries {
		if s.Duration > 0 {
			samples = append(samples, s.Duration)
		}
	}
	if len(samples) == 0 {
		return
	}
	minD, maxD, p50 := sortedDurationStats(samples)
	rep.Summary.IterationDurationMin = minD
	rep.Summary.IterationDurationMax = maxD
	rep.Summary.IterationDurationP50 = p50
}

// marshalAISummaryJSON returns one line of JSON for --ai-output: the report's
// summary block, or the JSON keyword null when absent (no aggregate stats).
func marshalAISummaryJSON(report *Report) ([]byte, error) {
	if report == nil || report.Summary == nil {
		return []byte("null"), nil
	}
	return json.Marshal(report.Summary)
}

// EffectiveParallelIterations returns the bounded diagnose worker count.
func EffectiveParallelIterations(conf *config.App) int {
	if conf == nil {
		return 1
	}
	parallel := max(conf.ParallelIterations, 1)
	if conf.Iterations > 0 && parallel > conf.Iterations {
		parallel = conf.Iterations
	}
	return parallel
}

func makeDiagnoseResultsDir(conf *config.App, goTestArgs []string, now time.Time) (string, error) {
	base := filepath.Join(conf.RepoRoot, diagnoseResultsDirName(goTestArgs, now))
	for i := 0; ; i++ {
		dir := base
		if i > 0 {
			dir = fmt.Sprintf("%s-%d", base, i)
		}
		err := os.Mkdir(dir, 0700)
		if err == nil {
			return dir, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}
}

type diagnoseIterationResult struct {
	iteration  int
	duration   time.Duration
	shuffle    int64
	iterErr    error
	fatalErr   error
	dumpErr    error
	failedFast bool
	failReason string
}

func runDiagnoseIterations(ctx context.Context, conf *config.App, out *output.Printer, resultsDir string, goTestArgs []string, resources []diagnoseIterationResource, hooks diagnoseRunHooks) (diagnoseRunState, error) {
	if hooks.runIteration == nil {
		hooks.runIteration = diagnoseIteration
	}
	if hooks.seed == nil {
		hooks.seed = func() int64 { return rand.Int64N(1<<62) + 1 }
	}
	parallel := EffectiveParallelIterations(conf)
	if len(resources) == 0 {
		resources = make([]diagnoseIterationResource, parallel)
	}
	if len(resources) < parallel {
		parallel = len(resources)
	}
	resources = resources[:parallel]
	state := diagnoseRunState{
		iterDurations:       make([]time.Duration, conf.Iterations),
		failedFastIteration: -1,
	}
	if conf.Shuffle {
		state.shuffleSeeds = make(map[int]int64)
	}

	if !out.AIOutput() {
		printDiagnoseIterationTableHeader(out)
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan int)
	results := make(chan diagnoseIterationResult)
	var parallelProgress *parallelDiagnoseProgress
	var diagnoseRunStart time.Time
	var progressTickWG sync.WaitGroup
	progressTickDone := make(chan struct{})
	if !out.AIOutput() && out.LiveInlineProgress() {
		state.liveProgress = true
		if parallel > 1 {
			parallelProgress = newParallelDiagnoseProgress(conf.Iterations)
			progressTickWG.Go(func() {
				tick := time.NewTicker(250 * time.Millisecond)
				defer tick.Stop()
				for {
					select {
					case <-progressTickDone:
						return
					case <-tick.C:
						parallelProgress.withRenderLock(func() {
							renderParallelDiagnoseProgressLine(out.HumanStderrWriter(), parallelProgress, time.Now(), true)
						})
					}
				}
			})
		} else {
			diagnoseRunStart = time.Now()
		}
	}
	defer func() {
		close(progressTickDone)
		progressTickWG.Wait()
	}()

	var serialProgressMu *sync.Mutex
	if parallel == 1 && state.liveProgress {
		serialProgressMu = new(sync.Mutex)
	}

	var wg sync.WaitGroup
	for _, resource := range resources {
		wg.Go(func() {
			executeSingleIteration(runCtx, conf, out, resultsDir, goTestArgs, resource, hooks, parallel, parallelProgress, diagnoseRunStart, serialProgressMu, jobs, results, cancel)
		})
	}

	wg.Go(func() {
		defer close(jobs)
		for i := range conf.Iterations {
			select {
			case <-runCtx.Done():
				return
			case jobs <- i:
			}
		}
	})

	go func() {
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for result := range results {
		if result.fatalErr != nil {
			if firstErr == nil {
				firstErr = result.fatalErr
				cancel()
			}
			continue
		}
		state.completed++
		state.iterDurations[result.iteration] = result.duration
		if state.shuffleSeeds != nil {
			state.shuffleSeeds[result.iteration] = result.shuffle
		}
		if parallelProgress != nil {
			parallelProgress.withRenderLock(func() {
				out.ClearInline()
				printDiagnoseIterationDigest(out, result.iteration, conf.Iterations, conf, resultsDir, result.duration)
			})
		} else {
			// Serial TTY: worker may redraw the next iteration's \r line before this
			// loop runs. Hold the same lock as redraw so ClearInline+Fprintln is not
			// interleaved with another progress draw (which would leave the cursor at EOL).
			if serialProgressMu != nil {
				serialProgressMu.Lock()
				out.ClearInline()
				printDiagnoseIterationDigest(out, result.iteration, conf.Iterations, conf, resultsDir, result.duration)
				serialProgressMu.Unlock()
			} else {
				printDiagnoseIterationDigest(out, result.iteration, conf.Iterations, conf, resultsDir, result.duration)
			}
		}
		if result.dumpErr != nil && !out.AIOutput() {
			out.Stderrf("postgres state dump iteration %d: %v\n", result.iteration, result.dumpErr)
		}
		if result.failedFast {
			state.failedFast = true
			if state.failedFastReason == "" {
				state.failedFastReason = result.failReason
			}
			if state.failedFastIteration < 0 {
				state.failedFastIteration = result.iteration
			}
		}
	}
	return state, firstErr
}

func executeSingleIteration(runCtx context.Context, conf *config.App, out *output.Printer, resultsDir string, goTestArgs []string, resource diagnoseIterationResource, hooks diagnoseRunHooks, parallel int, parallelProgress *parallelDiagnoseProgress, diagnoseRunStart time.Time, serialProgressMu *sync.Mutex, jobs <-chan int, results chan<- diagnoseIterationResult, cancel context.CancelFunc) {
	used := false
	for iteration := range jobs {
		if runCtx.Err() != nil {
			return
		}
		if used && resource.Reset != nil {
			if err := resource.Reset(runCtx); err != nil {
				if runCtx.Err() != nil {
					return
				}
				select {
				case results <- diagnoseIterationResult{iteration: iteration, fatalErr: fmt.Errorf("reset database before iteration %d: %w", iteration, err)}:
				case <-runCtx.Done():
				}
				cancel()
				return
			}
		}
		used = true
		var seed int64
		if conf.Shuffle {
			seed = hooks.seed()
		}
		iterStart := time.Now()
		iterErr := hooks.runIteration(runCtx, diagnoseIterationParams{
			Conf:             conf,
			Out:              out,
			ResultsDir:       resultsDir,
			GoTestArgs:       goTestArgs,
			Iteration:        iteration,
			ShuffleSeed:      seed,
			Env:              resource.Env,
			LiveProgress:     parallel == 1,
			ParallelProgress: parallelProgress,
			DiagnoseRunStart: diagnoseRunStart,
			SerialProgressMu: serialProgressMu,
		})
		iterDur := time.Since(iterStart)
		var dumpErr error
		if resource.DumpDiagnostics != nil {
			dumpErr = resource.DumpDiagnostics(runCtx, resultsDir, iteration)
		}
		failedFast, failReason := shouldFailFastIteration(conf, resultsDir, iteration, iterErr)
		failedFast = failedFast && runCtx.Err() == nil
		if failedFast {
			cancel()
		}
		select {
		case results <- diagnoseIterationResult{
			iteration:  iteration,
			duration:   iterDur,
			shuffle:    seed,
			iterErr:    iterErr,
			dumpErr:    dumpErr,
			failedFast: failedFast,
			failReason: failReason,
		}:
		case <-runCtx.Done():
			if failedFast {
				results <- diagnoseIterationResult{
					iteration:  iteration,
					duration:   iterDur,
					shuffle:    seed,
					iterErr:    iterErr,
					dumpErr:    dumpErr,
					failedFast: true,
					failReason: failReason,
				}
			}
			return
		}
	}
}

func shouldFailFastIteration(conf *config.App, resultsDir string, iteration int, iterErr error) (bool, string) {
	if conf == nil {
		return false, ""
	}
	if iterErr == nil && !conf.FailFast && len(conf.FailFastOn) == 0 {
		return false, ""
	}
	jsonPath := filepath.Join(resultsDir, fmt.Sprintf("iteration-%d.log.jsonl", iteration))
	f, err := os.Open(jsonPath)
	if err != nil {
		if iterErr != nil && conf.FailFast {
			return true, "failure"
		}
		return false, ""
	}
	defer f.Close()
	d, digestErr := DigestIterationJSONL(f, conf.SlowThreshold)
	if digestErr != nil {
		if iterErr != nil && conf.FailFast {
			return true, "failure"
		}
		return false, ""
	}
	if d.BuildFailure {
		return true, failFastReasonBuildFailure
	}
	if iterErr != nil && conf.FailFast {
		return true, "failure"
	}
	if len(conf.FailFastOn) == 0 {
		return false, ""
	}
	return failFastDigestMatch(d, conf.FailFastOn)
}

func failFastDigestMatch(d IterationDigest, categories []string) (bool, string) {
	for _, category := range categories {
		switch strings.ToLower(strings.TrimSpace(category)) {
		case config.FailFastOnAny:
			if d.Result == "fail" || d.Result == "timeout" || d.SlowTests > 0 {
				return true, config.FailFastOnAny
			}
		case config.FailFastOnFailure:
			if d.Result == "fail" && d.FailTests > 0 {
				return true, config.FailFastOnFailure
			}
		case config.FailFastOnTimeout:
			if d.Result == "timeout" || d.TimeoutTests > 0 {
				return true, config.FailFastOnTimeout
			}
		case config.FailFastOnSlow:
			if d.SlowTests > 0 {
				return true, config.FailFastOnSlow
			}
		}
	}
	return false, ""
}

// goTestFlagsBeforeArgs returns the portion of argv that belongs to `go test`
// itself, stopping before -args (flags after -args are passed to the test binary).
func goTestFlagsBeforeArgs(args []string) []string {
	for i, a := range args {
		if a == "-args" {
			return args[:i]
		}
	}
	return args
}

// parseDiagnoseGoTestCount returns the last -count in the portion of argv that
// belongs to `go test` itself (before -args). If no -count appears, set is false.
func parseDiagnoseGoTestCount(goTestArgs []string) (set bool, n int, err error) {
	args := goTestFlagsBeforeArgs(goTestArgs)
	for i := 0; i < len(args); i++ {
		a := args[i]
		if after, ok := strings.CutPrefix(a, "-count="); ok {
			v := after
			num, e := strconv.Atoi(strings.TrimSpace(v))
			if e != nil {
				return false, 0, fmt.Errorf("invalid -count value %q: %w", v, e)
			}
			if num < 1 {
				return false, 0, fmt.Errorf("invalid go test arguments: -count must be a positive integer, got %d", num)
			}
			set = true
			n = num
			continue
		}
		if a == "-count" {
			if i+1 >= len(args) {
				return false, 0, errors.New("invalid go test arguments: -count must be followed by a value")
			}
			i++
			num, e := strconv.Atoi(strings.TrimSpace(args[i]))
			if e != nil {
				return false, 0, fmt.Errorf("invalid -count value %q: %w", args[i], e)
			}
			if num < 1 {
				return false, 0, fmt.Errorf("invalid go test arguments: -count must be a positive integer, got %d", num)
			}
			set = true
			n = num
		}
	}
	return set, n, nil
}

// defaultGoTestTimeout matches `go help testflag` when -timeout is omitted.
const defaultGoTestTimeout = 10 * time.Minute

// parseGoTestTimeout returns the last -timeout in the go test flag section (before -args).
// If the value parses to 0, disabled is true (timeout disabled for the test binary).
// If no -timeout appears, set is false and callers should assume defaultGoTestTimeout for estimates.
func parseGoTestTimeout(goTestArgs []string) (set bool, d time.Duration, disabled bool, err error) {
	args := goTestFlagsBeforeArgs(goTestArgs)
	for i := 0; i < len(args); i++ {
		a := args[i]
		if after, ok := strings.CutPrefix(a, "-timeout="); ok {
			v := strings.TrimSpace(after)
			if v == "" {
				return false, 0, false, errors.New("invalid go test arguments: -timeout= requires a duration")
			}
			dur, e := time.ParseDuration(v)
			if e != nil {
				return false, 0, false, fmt.Errorf("invalid -timeout value %q: %w", v, e)
			}
			set = true
			if dur == 0 {
				disabled = true
			} else {
				disabled = false
				d = dur
			}
			continue
		}
		if a == "-timeout" {
			if i+1 >= len(args) {
				return false, 0, false, errors.New("invalid go test arguments: -timeout must be followed by a duration")
			}
			i++
			v := strings.TrimSpace(args[i])
			dur, e := time.ParseDuration(v)
			if e != nil {
				return false, 0, false, fmt.Errorf("invalid -timeout value %q: %w", args[i], e)
			}
			set = true
			if dur == 0 {
				disabled = true
			} else {
				disabled = false
				d = dur
			}
		}
	}
	return set, d, disabled, nil
}

// diagnoseIterationWaves returns ceil(iterations/workers) for scheduling diagnose iterations.
func diagnoseIterationWaves(iterations, workers int) int {
	w := max(workers, 1)
	if iterations < 1 {
		return 0
	}
	return (iterations + w - 1) / w
}

// diagnoseWallUpperBoundDetails holds the inputs used for the human-facing estimate line.
type diagnoseWallUpperBoundDetails struct {
	Bound       time.Duration
	Workers     int
	Waves       int
	PerInv      time.Duration
	UsedDefault bool // true when -timeout was omitted (10m assumed)
}

// diagnoseWallUpperBound returns worst-case wall clock if each go test invocation runs
// for the full per-invocation test timeout (Go default 10m when -timeout is unset).
// ok is false when -timeout=0 disables the test binary timeout (no finite bound).
func diagnoseWallUpperBound(conf *config.App, goTestArgs []string, resourceCount int) (diag diagnoseWallUpperBoundDetails, ok bool, err error) {
	if conf == nil {
		return diagnoseWallUpperBoundDetails{}, false, errors.New("config is nil")
	}
	set, d, disabled, err := parseGoTestTimeout(goTestArgs)
	if err != nil {
		return diagnoseWallUpperBoundDetails{}, false, err
	}
	if set && disabled {
		return diagnoseWallUpperBoundDetails{}, false, nil
	}
	perInv := defaultGoTestTimeout
	usedDefault := !set
	if set && !disabled {
		perInv = d
		usedDefault = false
	}
	parallel := EffectiveParallelIterations(conf)
	if resourceCount > 0 && resourceCount < parallel {
		parallel = resourceCount
	}
	if parallel < 1 {
		parallel = 1
	}
	waves := diagnoseIterationWaves(conf.Iterations, parallel)
	bound := time.Duration(waves) * perInv
	return diagnoseWallUpperBoundDetails{
		Bound:       bound,
		Workers:     parallel,
		Waves:       waves,
		PerInv:      perInv,
		UsedDefault: usedDefault,
	}, true, nil
}

func printDiagnoseRunTimeEstimate(out *output.Printer, conf *config.App, goTestArgs []string, resourceCount int) error {
	if out == nil || conf == nil {
		return nil
	}
	diag, ok, err := diagnoseWallUpperBound(conf, goTestArgs, resourceCount)
	if err != nil {
		return err
	}
	if out.AIOutput() {
		if !ok {
			out.Stderrf("lpr_s:inf\n")
			return nil
		}
		sec := max(diag.Bound.Round(time.Second)/time.Second, 0)
		out.Stderrf("lpr_s:%d\n", sec)
		return nil
	}
	est := "∞"
	if ok {
		est = formatDiagnoseWallClock(diag.Bound)
	}
	out.HumanStderr(termstyle.Muted.Render("Longest Possible Runtime: " + est))
	return nil
}

// WarnDiagnoseGoTestCount prints hints when the user sets -count on go test, and
// returns an error if -count values in the go test flag section are malformed.
func WarnDiagnoseGoTestCount(w io.Writer, goTestArgs []string) error {
	set, n, err := parseDiagnoseGoTestCount(goTestArgs)
	if err != nil {
		return err
	}
	if !set {
		return nil
	}
	if n == 1 {
		fmt.Fprintln(w, termstyle.Muted.Render(
			"note: -count=1 is unnecessary; diagnose adds -count=1 when you omit it."))
		return nil
	}
	fmt.Fprintln(w, termstyle.Muted.Render(
		"note: prefer diagnose --iterations for repetition; use -count>1 only if you want to avoid overhead between diagnose iterations (e.g. DB setup/teardown)."))
	return nil
}

// filterDiagnoseUserGoTestArgs removes -json/--json from the go test flag
// section so the harness can inject -json; arguments after -args are unchanged.
func filterDiagnoseUserGoTestArgs(args []string) []string {
	split := len(args)
	for i, a := range args {
		if a == "-args" {
			split = i
			break
		}
	}
	prefix := args[:split]
	suffix := args[split:]
	var out []string
	for _, a := range prefix {
		if a == "-json" || a == "--json" {
			continue
		}
		out = append(out, a)
	}
	return append(out, suffix...)
}

// buildDiagnoseArgs constructs the `go test` argv for a single diagnose iteration.
func buildDiagnoseArgs(goTestArgs []string, shuffleSeed int64) ([]string, error) {
	filtered := filterDiagnoseUserGoTestArgs(goTestArgs)
	set, n, err := parseDiagnoseGoTestCount(goTestArgs)
	if err != nil {
		return nil, err
	}
	args := []string{"test", "-json"}
	args = append(args, filtered...)
	if shuffleSeed != 0 {
		args = append(args, fmt.Sprintf("-shuffle=%d", shuffleSeed))
	}
	if !set || n <= 1 {
		args = append(args, "-count=1")
	}
	return args, nil
}

func printDiagnoseResultsDirHeader(out *output.Printer, resultsDir string) {
	if out.AIOutput() {
		out.Stdoutln(resultsDir)
		return
	}
	out.HumanStderr(termstyle.Muted.Render("results directory: ") + termstyle.Label.Render(resultsDir))
}

// formatDiagnoseWallClock formats total wall time for human diagnose footers (0.1s resolution, seconds show two decimals when fractional).
func formatDiagnoseWallClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	d = d.Round(100 * time.Millisecond)
	if d == 0 {
		return "0s"
	}
	h := d / time.Hour
	d -= h * time.Hour
	mins := d / time.Minute
	d -= mins * time.Minute
	s := float64(d) / float64(time.Second)
	secStr := formatDiagnoseSecondsFragment(s)
	if h > 0 {
		return fmt.Sprintf("%dh%dm%s", h, mins, secStr)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm%s", mins, secStr)
	}
	return secStr
}

func formatDiagnoseSecondsFragment(s float64) string {
	cs := max(int(math.Round(s*100)), 0)
	w := cs / 100
	f := cs % 100
	if f == 0 {
		return fmt.Sprintf("%ds", w)
	}
	return fmt.Sprintf("%d.%02ds", w, f)
}

// startDiagnoseAnalyzingProgress prints a live "analyzing [duration]" line and returns stop.
// Call stop once analysis is done (and defer stop as well — it is idempotent) so the final
// "analyzing [duration] ✅" line is written before diagnosis complete, not at process exit.
func startDiagnoseAnalyzingProgress(out *output.Printer, afterLiveProgress bool) (stop func(error)) {
	if out.AIOutput() {
		return func(error) {}
	}
	if afterLiveProgress {
		_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K\n")
	}

	analyzeStart := time.Now()
	var once sync.Once
	finalize := func(live bool, err error) {
		elapsed := max(time.Since(analyzeStart).Round(time.Second), 0)
		mark := termstyle.OK.Render("✅")
		if err != nil {
			mark = termstyle.Bad.Render("❌")
		}
		line := termstyle.Label.Render("analyzing") + " " + termstyle.Muted.Render("["+elapsed.String()+"]") + " " + mark
		if live {
			_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K"+line+"\n")
			return
		}
		out.HumanStderr(line)
	}

	if !out.LiveInlineProgress() {
		return func(err error) {
			once.Do(func() { finalize(false, err) })
		}
	}

	renderProgress := func() {
		elapsed := max(time.Since(analyzeStart).Round(time.Second), 0)
		line := termstyle.Label.Render("analyzing") + " " + termstyle.Muted.Render("["+elapsed.String()+"]")
		_, _ = fmt.Fprint(out.HumanStderrWriter(), "\r\033[K"+line)
	}
	renderProgress()

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		tick := time.NewTicker(250 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				renderProgress()
			}
		}
	})

	return func(err error) {
		once.Do(func() {
			close(done)
			wg.Wait()
			finalize(true, err)
		})
	}
}

func printDiagnoseIterationDigest(out *output.Printer, iterationIdx0, totalIters int, conf *config.App, resultsDir string, iterDur time.Duration) {
	jsonPath := filepath.Join(resultsDir, fmt.Sprintf("iteration-%d.log.jsonl", iterationIdx0))
	f, err := os.Open(jsonPath)
	if err != nil {
		out.Stderrf("diagnose iteration %d summary: %v\n", iterationIdx0+1, err)
		return
	}
	defer f.Close()
	d, err := DigestIterationJSONL(f, conf.SlowThreshold)
	if err != nil {
		out.Stderrf("diagnose iteration %d summary: %v\n", iterationIdx0+1, err)
		return
	}
	iter := iterationIdx0 + 1
	if out.AIOutput() {
		out.Stdoutln(formatIterationDigestAI(iter, totalIters, d, iterDur))
		return
	}
	printIterationDigestHuman(out, iter, totalIters, d, iterDur)
}

// formatIterationDigestAI prints one line for --ai-output diagnose progress.
// Tokens: d iter/total; p|f|t result; wall seconds; r named tests that ran;
// f failing-test entries; t timeouts; s slow tests.
func formatIterationDigestAI(iter, total int, d IterationDigest, dur time.Duration) string {
	rs := "?"
	switch d.Result {
	case "pass":
		rs = "p"
	case "fail":
		rs = "f"
	case "timeout":
		rs = "t"
	}
	sec := max(int(dur.Round(time.Second)/time.Second), 0)
	return fmt.Sprintf("d %d/%d %s %ds r%d f%d t%d s%d", iter, total, rs, sec, d.RanTests, d.FailTests, d.TimeoutTests, d.SlowTests)
}

func printIterationDigestHuman(out *output.Printer, iter, total int, d IterationDigest, dur time.Duration) {
	_ = total
	out.HumanStderr(formatDiagnoseIterationTableRow(iter, d, dur))
}

func renderIterationResultHuman(r string) string {
	switch r {
	case "pass":
		return termstyle.OK.Render("pass")
	case "fail":
		return termstyle.Bad.Render("fail")
	case "timeout":
		return termstyle.Accent.Render("timeout")
	default:
		return termstyle.Muted.Render(r)
	}
}

// syncedWriter serializes writes to w so stdout and stderr from `go test` can
// share one JSONL file without interleaved corrupt lines.
type syncedWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (sw *syncedWriter) Write(p []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.w.Write(p)
}

func diagnoseIteration(ctx context.Context, p diagnoseIterationParams) error {
	conf, out := p.Conf, p.Out
	resultsDir, goTestArgs := p.ResultsDir, p.GoTestArgs
	iteration, shuffleSeed := p.Iteration, p.ShuffleSeed
	env := p.Env
	liveProgress, parallelProgress := p.LiveProgress, p.ParallelProgress
	diagnoseRunStart, serialProgressMu := p.DiagnoseRunStart, p.SerialProgressMu

	start := time.Now()
	jsonPath := filepath.Join(resultsDir, fmt.Sprintf("iteration-%d.log.jsonl", iteration))
	resultsFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	args, err := buildDiagnoseArgs(goTestArgs, shuffleSeed)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = conf.RepoRoot
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), env...)
	// Soft-cancel on ctx cancellation so `go test -json` gets a chance to flush
	// its final events before we escalate to SIGKILL after WaitDelay.
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second

	if out.AIOutput() {
		sw := &syncedWriter{w: resultsFile}
		cmd.Stdout = sw
		cmd.Stderr = sw
		return cmd.Run()
	}

	sw := &syncedWriter{w: resultsFile}
	cmd.Stderr = sw

	totalPkgs := -1
	if n, listErr := listTestPackageCount(ctx, conf.RepoRoot, goTestArgs); listErr == nil {
		totalPkgs = n
	}
	prog := newDiagnoseProgress(totalPkgs)
	if parallelProgress != nil {
		parallelProgress.start(iteration)
		defer parallelProgress.finish(iteration)
	}

	pr, pw := io.Pipe()
	cmd.Stdout = pw

	live := liveProgress && out.LiveInlineProgress()
	iter, iters := iteration+1, conf.Iterations
	if liveProgress && !live {
		out.HumanStderr(termstyle.Muted.Render(fmt.Sprintf("iteration %d/%d started", iter, iters)))
	}

	redraw := func(liveInline bool) {
		if serialProgressMu != nil {
			serialProgressMu.Lock()
			defer serialProgressMu.Unlock()
		}
		renderDiagnoseProgressLine(out.HumanStderrWriter(), iter, iters, time.Since(start), diagnoseRunStart, time.Now(), liveInline)
	}

	var readWG sync.WaitGroup
	var scanErr error
	readWG.Go(func() {
		r := bufio.NewReaderSize(pr, 1024*1024)
		for {
			line, err := r.ReadBytes('\n')
			if len(line) > 0 {
				if _, werr := sw.Write(line); werr != nil {
					break
				}
				completedIncreased := prog.onTestJSONLine(line)
				if completedIncreased && !live {
					redraw(false)
				}
			}
			if err != nil {
				if err != io.EOF {
					scanErr = err
				}
				break
			}
		}
	})

	tickDone := make(chan struct{})
	var tickWG sync.WaitGroup
	if live {
		tickWG.Go(func() {
			tick := time.NewTicker(250 * time.Millisecond)
			defer tick.Stop()
			for {
				select {
				case <-tickDone:
					return
				case <-tick.C:
					redraw(true)
				}
			}
		})
		redraw(true)
	}

	runErr := cmd.Start()
	started := runErr == nil
	if started {
		runErr = cmd.Wait()
		_ = pw.Close()
	} else {
		_ = pw.CloseWithError(runErr)
	}
	readWG.Wait()
	close(tickDone)
	tickWG.Wait()

	if live {
		out.ClearInline()
	}
	if scanErr != nil {
		return fmt.Errorf("reading go test output: %w", scanErr)
	}
	return runErr
}

func newRunMeta(conf *config.App, goTestArgs []string, resultsDir string, started time.Time, finished *time.Time, hasDatabase bool) *RunMeta {
	if conf == nil {
		return nil
	}
	target := guessPackagePatternForSlug(goTestArgs)
	slug := diagnoseTargetSlug(target)
	args := append([]string(nil), goTestArgs...)
	slow := conf.SlowThreshold
	if slow == 0 {
		slow = 30 * time.Second
	}
	var ffo []string
	if n, err := config.NormalizeFailFastOn(conf.FailFastOn); err == nil && len(n) > 0 {
		ffo = n
	}
	par := max(conf.ParallelIterations, 1)
	var fin *time.Time
	if finished != nil {
		t := finished.UTC()
		fin = &t
	}
	return &RunMeta{
		ResultsDirBasename: filepath.Base(resultsDir),
		StartedAt:          started.UTC(),
		FinishedAt:         fin,
		GoTestArgs:         args,
		TargetSlug:         slug,
		DiagnoseIterations: conf.Iterations,
		ParallelIterations: par,
		SlowThreshold:      slow,
		FailFast:           conf.FailFast,
		FailFastOn:         ffo,
		Shuffle:            conf.Shuffle,
		PostgresVersion:    conf.PostgresVersion,
		HasDatabase:        hasDatabase,
	}
}
