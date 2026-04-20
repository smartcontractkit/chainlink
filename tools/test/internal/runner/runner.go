package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
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
// reflected in report.json. Survey returns a non-nil error only for setup
// failures (e.g. mkdir, database reset), not for failing tests.
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
				fmt.Fprintln(os.Stderr, "--fail-fast: true, stopping")
				failedFast = true
				break
			}
		}
		completed = i + 1
	}

	interrupted := ctx.Err() != nil
	if interrupted && !conf.AIOutput {
		fmt.Fprintf(os.Stderr, "interrupted after %d/%d iterations — analyzing partial results...\n", completed, conf.Iterations)
	}

	if failedFast {
		fmt.Fprintln(os.Stderr, "--fail-fast set, stopping early")
	}

	report, analyzeErr := AnalyzeResults(resultsDir, conf.SlowThreshold)
	if analyzeErr != nil {
		fmt.Fprintf(os.Stderr, "analyze results: %v\n", analyzeErr)
	} else if err := WriteReport(resultsDir, report); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
	}

	if !conf.AIOutput {
		fmt.Fprintf(os.Stderr, "survey complete (%s)\n", time.Since(start).Round(time.Millisecond))
		if report != nil {
			PrintSummary(os.Stderr, report)
		}
		fmt.Fprintf(os.Stderr, "results in %s\n", resultsDir)
	}
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
	if !conf.AIOutput {
		fmt.Fprintf(os.Stderr, "iteration %d/%d...", iteration+1, conf.Iterations)
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = conf.RepoRoot
	cmd.Stdout = resultsFile
	cmd.Stderr = resultsFile
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	// Soft-cancel on ctx cancellation so `go test -json` gets a chance to flush
	// its final events before we escalate to SIGKILL after WaitDelay.
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second

	if conf.AIOutput {
		err = cmd.Run()
	} else {
		// Show progress in real time.
		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					elapsed := time.Since(start).Round(time.Second)
					fmt.Fprintf(os.Stderr, "\r\033[Kiteration %d/%d... (%s)",
						iteration+1, conf.Iterations, elapsed.String())
				}
			}
		}()

		if err = cmd.Start(); err != nil {
			close(done)
			return err
		}
		err = cmd.Wait()
		close(done)
	}

	if !conf.AIOutput {
		status := "✅"
		if err != nil {
			status = "❌"
		}
		fmt.Fprintf(os.Stderr, "\r\033[Kiteration %d/%d %s (%s)\n",
			iteration+1, conf.Iterations, status, time.Since(start).Round(time.Millisecond))
	}
	return err
}
