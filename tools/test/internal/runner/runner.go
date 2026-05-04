package runner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/term"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

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
// Test iteration failures do not stop later runs (unless --fail-fast); they are
// reflected in report.json. Diagnose returns a non-nil error for setup failures
// (e.g. mkdir, database reset), analyze/write report failures, or ctx errors
// bubbling from dependencies — not for failing tests alone.
// resetDB (optional) runs before each iteration after the first to restore the
// database to its freshly-prepared state.
// dumpDB (optional) runs after each iteration to capture database state for
// per-iteration diagnosis; errors are logged but do not fail the diagnose run.
func Diagnose(ctx context.Context, conf *config.App, goTestArgs []string, resetDB func(context.Context) error, dumpDB func(context.Context, string, int) error) error {
	start := time.Now()

	resultsDir := filepath.Join(conf.RepoRoot, diagnoseResultsDirName(conf, goTestArgs, start))
	err := os.MkdirAll(resultsDir, 0700)
	if err != nil {
		return err
	}

	var (
		completed     int
		failedFast    bool
		iterDurations = make([]time.Duration, 0, conf.Iterations)
		shuffleSeeds  map[int]int64
	)
	if conf.Shuffle {
		shuffleSeeds = make(map[int]int64)
	}
	for i := range conf.Iterations {
		if ctx.Err() != nil {
			break
		}
		if i > 0 && resetDB != nil {
			if err := resetDB(ctx); err != nil {
				if ctx.Err() != nil {
					break
				}
				return fmt.Errorf("reset database before iteration %d: %w", i, err)
			}
		}
		var seed int64
		if conf.Shuffle {
			seed = rand.Int64N(1<<62) + 1 // always nonzero
			shuffleSeeds[i] = seed
		}
		iterStart := time.Now()
		iterErr := diagnoseIteration(ctx, conf, resultsDir, goTestArgs, i, seed)
		iterDurations = append(iterDurations, time.Since(iterStart))
		if dumpDB != nil {
			if dumpErr := dumpDB(ctx, resultsDir, i); dumpErr != nil && !conf.AIOutput {
				fmt.Fprintf(os.Stderr, "postgres state dump iteration %d: %v\n", i, dumpErr)
			}
		}
		if iterErr != nil && conf.FailFast {
			failedFast = true
			break
		}
		completed = i + 1
	}

	interrupted := ctx.Err() != nil
	if interrupted && !conf.AIOutput {
		fmt.Fprintln(os.Stderr,
			termstyle.Accent.Render(fmt.Sprintf("interrupted after %d/%d iterations", completed, conf.Iterations))+
				termstyle.Muted.Render(" — analyzing partial results…"))
	}

	if failedFast && !conf.AIOutput {
		fmt.Fprintln(os.Stderr, termstyle.Accent.Render("--fail-fast set, stopping early"))
	}

	report, logs, analyzeErr := AnalyzeResults(resultsDir, conf.SlowThreshold)
	if analyzeErr != nil {
		fmt.Fprintf(os.Stderr, "analyze results: %v\n", analyzeErr)
		return analyzeErr
	}
	if report != nil {
		for i, d := range iterDurations {
			if i >= len(report.IterationSummaries) {
				break
			}
			report.IterationSummaries[i].Duration = d
			if shuffleSeeds != nil {
				report.IterationSummaries[i].ShuffleSeed = shuffleSeeds[i]
			}
		}
	}
	if err := WriteLogFiles(resultsDir, report, logs); err != nil {
		fmt.Fprintf(os.Stderr, "write log files: %v\n", err)
		return err
	}
	if err := WriteReport(resultsDir, report); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		return err
	}
	if err := WriteCSV(resultsDir, report); err != nil {
		fmt.Fprintf(os.Stderr, "write csv: %v\n", err)
		return err
	}

	reportPath := filepath.Join(resultsDir, "report.json")
	if conf.AIOutput {
		fmt.Fprintln(os.Stdout, reportPath)
		return nil
	}

	fmt.Fprintln(os.Stderr,
		termstyle.Label.Render("diagnose complete")+
			termstyle.Muted.Render(fmt.Sprintf(" (%s)", time.Since(start).Round(time.Millisecond))))
	if report != nil {
		PrintSummary(os.Stderr, report)
	}
	fmt.Fprintln(os.Stderr,
		termstyle.Muted.Render("results in ")+termstyle.Label.Render(resultsDir))
	return nil
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

func diagnoseIteration(ctx context.Context, conf *config.App, resultsDir string, goTestArgs []string, iteration int, shuffleSeed int64) error {
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
	cmd.Env = os.Environ()
	// Soft-cancel on ctx cancellation so `go test -json` gets a chance to flush
	// its final events before we escalate to SIGKILL after WaitDelay.
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second

	if conf.AIOutput {
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

	pr, pw := io.Pipe()
	cmd.Stdout = pw

	isTTY := term.IsTerminal(os.Stderr.Fd())
	iter, iters := iteration+1, conf.Iterations
	if !isTTY {
		fmt.Fprintln(os.Stderr,
			termstyle.Muted.Render(fmt.Sprintf("iteration %d/%d started", iter, iters)))
	}

	redraw := func(isTTYLine bool) {
		renderDiagnoseProgressLine(os.Stderr, iter, iters, time.Since(start), prog, isTTYLine)
	}

	var readWG sync.WaitGroup
	var scanErr error
	readWG.Go(func() {
		sc := bufio.NewScanner(pr)
		sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		for sc.Scan() {
			line := sc.Bytes()
			out := make([]byte, len(line)+1)
			copy(out, line)
			out[len(line)] = '\n'
			if _, werr := sw.Write(out); werr != nil {
				break
			}
			if prog.onTestJSONLine(line) && !isTTY {
				redraw(false)
			}
		}
		scanErr = sc.Err()
	})

	tickDone := make(chan struct{})
	var tickWG sync.WaitGroup
	if isTTY {
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

	if isTTY {
		fmt.Fprint(os.Stderr, "\r\033[K")
	}
	if started {
		status := termstyle.OK.Render("✅")
		if runErr != nil {
			status = termstyle.Bad.Render("❌")
		}
		fmt.Fprintln(os.Stderr,
			termstyle.Label.Render(fmt.Sprintf("iteration %d/%d ", iter, iters))+
				status+" "+
				termstyle.Muted.Render(fmt.Sprintf("(%s)", time.Since(start).Round(time.Millisecond))))
	}
	if scanErr != nil {
		return fmt.Errorf("reading go test output: %w", scanErr)
	}
	return runErr
}
