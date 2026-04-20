package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
func Survey(ctx context.Context, conf *config.App, targetDir string, resetDB func(context.Context) error) error {
	start := time.Now()

	resultsDir := filepath.Join(conf.RepoRoot, "test-survey-results-"+time.Now().Format("20060102150405"))
	err := os.MkdirAll(resultsDir, 0700)
	if err != nil {
		return err
	}

	var (
		completed  int
		failedFast bool
	)
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
		if err := surveyIteration(ctx, conf, resultsDir, targetDir, i); err != nil {
			if conf.FailFast {
				if !conf.AIOutput {
					fmt.Fprintln(os.Stderr, termstyle.Accent.Render("--fail-fast: true, stopping"))
				}
				failedFast = true
				break
			}
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

	report, analyzeErr := AnalyzeResults(resultsDir, conf.SlowThreshold)
	if analyzeErr != nil {
		fmt.Fprintf(os.Stderr, "analyze results: %v\n", analyzeErr)
		return analyzeErr
	}
	if err := WriteReport(resultsDir, report); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
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

func surveyIteration(ctx context.Context, conf *config.App, resultsDir string, targetDir string, iteration int) error {
	start := time.Now()
	jsonPath := filepath.Join(resultsDir, fmt.Sprintf("iteration-%d.log.jsonl", iteration))
	resultsFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	args := []string{"test", "-json", "-count=1", "-timeout", conf.Timeout.String(), targetDir}
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
	if !isTTY {
		fmt.Fprintln(os.Stderr,
			termstyle.Muted.Render(fmt.Sprintf("iteration %d/%d started (stderr is not a TTY; sparse package progress)",
				iteration+1, conf.Iterations)))
	}

	var readWG sync.WaitGroup
	readWG.Add(1)
	go func() {
		defer readWG.Done()
		scanner := bufio.NewScanner(pr)
		scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			if _, werr := resultsFile.Write(line); werr != nil {
				break
			}
			if _, werr := resultsFile.WriteString("\n"); werr != nil {
				break
			}
			if prog.onTestJSONLine(line) && !isTTY {
				renderSurveyProgressLine(os.Stderr, iteration+1, conf.Iterations, time.Since(start), prog, false)
			}
		}
		_ = scanner.Err()
	}()

	tickDone := make(chan struct{})
	var tickWG sync.WaitGroup
	if isTTY {
		tickWG.Add(1)
		go func() {
			defer tickWG.Done()
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-tickDone:
					return
				case <-ticker.C:
					renderSurveyProgressLine(os.Stderr, iteration+1, conf.Iterations, time.Since(start), prog, true)
				}
			}
		}()
		renderSurveyProgressLine(os.Stderr, iteration+1, conf.Iterations, time.Since(start), prog, true)
	}

	if err = cmd.Start(); err != nil {
		_ = pw.CloseWithError(err)
		readWG.Wait()
		close(tickDone)
		tickWG.Wait()
		if isTTY {
			fmt.Fprint(os.Stderr, "\r\033[K")
		}
		return err
	}

	err = cmd.Wait()
	_ = pw.Close()
	readWG.Wait()
	close(tickDone)
	tickWG.Wait()

	status := termstyle.OK.Render("✅")
	if err != nil {
		status = termstyle.Bad.Render("❌")
	}
	if isTTY {
		fmt.Fprintf(os.Stderr, "\r\033[K")
	}
	fmt.Fprintln(os.Stderr,
		termstyle.Label.Render(fmt.Sprintf("iteration %d/%d ", iteration+1, conf.Iterations))+
			status+" "+
			termstyle.Muted.Render(fmt.Sprintf("(%s)", time.Since(start).Round(time.Millisecond))))
	return err
}
