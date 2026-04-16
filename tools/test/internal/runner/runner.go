package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/repo"
)

// GoTest runs `go test` with the given args (repo root as working directory).
func GoTest(ctx context.Context, args []string) error {
	root, err := repo.RootFromWd()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "go", append([]string{"test"}, args...)...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// Gotestsum runs `gotestsum` with the given args (repo root as working directory).
func Gotestsum(ctx context.Context, args []string) error {
	if _, err := exec.LookPath("gotestsum"); err != nil {
		return fmt.Errorf("gotestsum not on PATH: install with go install gotest.tools/gotestsum@latest: %w", err)
	}
	root, err := repo.RootFromWd()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "gotestsum", args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// Survey runs gotestsum repeatedly with a JSON file per iteration.
// Failures do not stop later iterations; the process exits non-zero if any iteration failed.
func Survey(ctx context.Context, iterations int, gotestsumArgs []string) error {
	if _, err := exec.LookPath("gotestsum"); err != nil {
		return fmt.Errorf("gotestsum not on PATH: install with go install gotest.tools/gotestsum@latest: %w", err)
	}
	root, err := repo.RootFromWd()
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "chainlink-test-survey-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	var anyErr bool
	for i := range iterations {
		jsonPath := filepath.Join(tmpDir, fmt.Sprintf("events-%d.jsonl", i))
		args := append([]string{"--jsonfile", jsonPath}, gotestsumArgs...)
		fmt.Fprintf(os.Stderr, "survey: iteration %d/%d (json -> %s)\n", i+1, iterations, jsonPath)
		cmd := exec.CommandContext(ctx, "gotestsum", args...)
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			anyErr = true
			fmt.Fprintf(os.Stderr, "survey: iteration %d: %v\n", i+1, err)
		}
	}
	if anyErr {
		return fmt.Errorf("one or more survey iterations failed")
	}
	return nil
}
