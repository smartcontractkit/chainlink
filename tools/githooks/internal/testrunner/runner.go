package testrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executor abstracts command execution for testability.
type Executor interface {
	Run(ctx context.Context, dir string, name string, args ...string) error
}

type osExecutor struct {
	stdout io.Writer
	stderr io.Writer
}

func (e *osExecutor) Run(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Config holds configuration for the test runner.
type Config struct {
	RepoRoot string
	Packages []string
	Short    bool
	Args     []string
	Executor Executor
	Stdout   io.Writer
	Stderr   io.Writer
}

// EnsureBinary checks if tools/test/.bin/test exists, and builds it if missing.
func EnsureBinary(ctx context.Context, repoRoot string, stdout, stderr io.Writer) error {
	binPath := filepath.Join(repoRoot, "tools/test/.bin/test")
	if _, err := os.Stat(binPath); err == nil {
		return nil
	}

	testModuleDir := filepath.Join(repoRoot, "tools/test")
	cmd := exec.CommandContext(ctx, "go", "build", "-o", ".bin/test", ".")
	cmd.Dir = testModuleDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build tools/test harness: %w", err)
	}
	return nil
}

// Run executes tools/test harness against specified packages.
func Run(ctx context.Context, cfg Config) error {
	if len(cfg.Packages) == 0 {
		return nil
	}

	execRunner := cfg.Executor
	if execRunner == nil {
		if err := EnsureBinary(ctx, cfg.RepoRoot, cfg.Stdout, cfg.Stderr); err != nil {
			return err
		}
		execRunner = &osExecutor{stdout: cfg.Stdout, stderr: cfg.Stderr}
	}

	binPath := filepath.Join(cfg.RepoRoot, "tools/test/.bin/test")

	var args []string
	if cfg.Short {
		args = append(args, "-short")
	}
	args = append(args, cfg.Args...)
	args = append(args, cfg.Packages...)

	fmt.Fprintf(cfg.Stdout, "==> Running tests on packages: %s\n", strings.Join(cfg.Packages, " "))
	if err := execRunner.Run(ctx, cfg.RepoRoot, binPath, args...); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	return nil
}
