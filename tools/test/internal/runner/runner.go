package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
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

// Survey runs go test -json once per iteration, writing each stream to
// iteration-<n>.log.jsonl, then analyzes and writes report.json.
// Test iteration failures do not stop later runs (unless --fail-fast); they are
// reflected in report.json. Survey returns a non-nil error for setup failures
// (e.g. mkdir, database reset), analyze/write report failures, or ctx errors
// bubbling from dependencies — not for failing tests alone.
// resetDB (optional) runs before each iteration after the first to restore the
// database to its freshly-prepared state.
// dumpDB (optional) runs after each iteration to capture database state for
// per-iteration diagnosis; errors are logged but do not fail the survey.
func Survey(ctx context.Context, conf *config.App, targetDir string, resetDB func(context.Context) error, dumpDB func(context.Context, string, int) error) error {
	start := time.Now()

	resultsDir := filepath.Join(conf.RepoRoot, surveyResultsDirName(conf, targetDir, start))
	err := os.MkdirAll(resultsDir, 0700)
	if err != nil {
		return err
	}

	var (
		completed    int
		failedFast   bool
		shuffleSeeds map[int]int64
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
		iterErr := surveyIteration(ctx, conf, resultsDir, targetDir, i, seed)
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
	if report != nil && len(shuffleSeeds) > 0 {
		report.ShuffleSeeds = shuffleSeeds
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
		termstyle.Label.Render("survey complete")+
			termstyle.Muted.Render(fmt.Sprintf(" (%s)", time.Since(start).Round(time.Millisecond))))
	if report != nil {
		PrintSummary(os.Stderr, report)
	}
	fmt.Fprintln(os.Stderr,
		termstyle.Muted.Render("results in ")+termstyle.Label.Render(resultsDir))
	return nil
}

// buildSurveyArgs constructs the `go test` argv for a single survey iteration.
func buildSurveyArgs(conf *config.App, targetDir string, shuffleSeed int64) []string {
	args := []string{"test", "-json", "-count=1", "-timeout", conf.Timeout.String()}
	if conf.Race {
		args = append(args, "-race")
	}
	if conf.Run != "" {
		args = append(args, "-run="+conf.Run)
	}
	if conf.CPU != "" {
		args = append(args, "-cpu="+conf.CPU)
	}
	if conf.Parallel > 0 {
		args = append(args, fmt.Sprintf("-parallel=%d", conf.Parallel))
	}
	if shuffleSeed != 0 {
		args = append(args, fmt.Sprintf("-shuffle=%d", shuffleSeed))
	}
	args = append(args, targetDir)
	return args
}

func surveyIteration(ctx context.Context, conf *config.App, resultsDir string, targetDir string, iteration int, shuffleSeed int64) error {
	start := time.Now()
	jsonPath := filepath.Join(resultsDir, fmt.Sprintf("iteration-%d.log.jsonl", iteration))
	resultsFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	args := buildSurveyArgs(conf, targetDir, shuffleSeed)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = conf.RepoRoot
	cmd.Stderr = resultsFile
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	// Soft-cancel on ctx cancellation so `go test -json` gets a chance to flush
	// its final events before we escalate to SIGKILL after WaitDelay.
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second

	if conf.AIOutput {
		cmd.Stdout = resultsFile
		return cmd.Run()
	}

	totalPkgs := -1
	if n, listErr := listTestPackageCount(ctx, conf.RepoRoot, targetDir); listErr == nil {
		totalPkgs = n
	}
	prog := newSurveyProgress(totalPkgs)

	pr, pw := io.Pipe()
	cmd.Stdout = pw

	isTTY := term.IsTerminal(os.Stderr.Fd())
	iter, iters := iteration+1, conf.Iterations
	if !isTTY {
		fmt.Fprintln(os.Stderr,
			termstyle.Muted.Render(fmt.Sprintf("iteration %d/%d started", iter, iters)))
	}

	redraw := func(isTTYLine bool) {
		renderSurveyProgressLine(os.Stderr, iter, iters, time.Since(start), prog, isTTYLine)
	}

	var readWG sync.WaitGroup
	readWG.Go(func() {
		sc := bufio.NewScanner(pr)
		sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		for sc.Scan() {
			line := sc.Bytes()
			if _, werr := resultsFile.Write(line); werr != nil {
				break
			}
			if _, werr := resultsFile.WriteString("\n"); werr != nil {
				break
			}
			if prog.onTestJSONLine(line) && !isTTY {
				redraw(false)
			}
		}
		_ = sc.Err()
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
	return runErr
}
