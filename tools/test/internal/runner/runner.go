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

type surveyResult struct {
	Iteration int
	Time      time.Duration
	Error     error
}

// Survey runs gotestsum repeatedly with a JSON file per iteration.
// Failures do not stop later iterations; the process exits non-zero if any iteration failed.
// resetDB (optional) runs before each iteration after the first to restore the
// database to its freshly-prepared state.
func Survey(ctx context.Context, conf *config.App, targetDir string, resetDB func(context.Context) error) error {
	start := time.Now()

	resultsDir := filepath.Join(conf.RepoRoot, "test-survey-results-"+time.Now().Format("20060102150405"))
	err := os.MkdirAll(resultsDir, 0700)
	if err != nil {
		return err
	}

	var results []surveyResult
	for i := range conf.Iterations {
		if i > 0 && resetDB != nil {
			if err := resetDB(ctx); err != nil {
				return fmt.Errorf("reset database before iteration %d: %w", i, err)
			}
		}
		if err := surveyIteration(ctx, conf, resultsDir, targetDir, i); err != nil {
			results = append(results, surveyResult{
				Iteration: i,
				Time:      time.Since(start),
				Error:     err,
			})
		}
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

	args := []string{"test", "-json", "-count=1", targetDir}
	if !conf.AIOutput {
		fmt.Fprintf(os.Stderr, "iteration %d/%d...", iteration+1, conf.Iterations)
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = conf.RepoRoot
	cmd.Stdout = resultsFile
	cmd.Stderr = resultsFile
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	err = cmd.Run()

	if !conf.AIOutput {
		status := "✅"
		if err != nil {
			status = "❌"
		}
		fmt.Fprintf(os.Stderr, " %s (%s)\n", status, time.Since(start).Round(time.Millisecond))
	}
	return err
}
