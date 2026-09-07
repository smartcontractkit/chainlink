package tidy

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
)

// CommandRunner executes a command in a directory.
type CommandRunner func(ctx context.Context, dir string, args ...string) error

// Config holds options for running tidy.
type Config struct {
	Runner CommandRunner
}

// defaultRunner executes "go" with given args in dir.
func defaultRunner(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go %v in %s failed: %w (output: %s)", args, dir, err, string(out))
	}
	return nil
}

// Run executes `go mod tidy` in parallel across all specified module directories.
func Run(ctx context.Context, repoRoot string, moduleDirs []string, cfg ...Config) error {
	if len(moduleDirs) == 0 {
		return nil
	}

	runner := defaultRunner
	if len(cfg) > 0 && cfg[0].Runner != nil {
		runner = cfg[0].Runner
	}

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, mod := range moduleDirs {
		modDir := filepath.Join(repoRoot, mod)
		wg.Add(1)
		go func(dir string, m string) {
			defer wg.Done()
			if err := runner(ctx, dir, "mod", "tidy"); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("module %s: %w", m, err))
				mu.Unlock()
			}
		}(modDir, mod)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
