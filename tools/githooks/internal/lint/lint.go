package lint

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
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

// Config holds configuration for the linter runner.
type Config struct {
	RepoRoot  string
	Targets   []modules.ModulePackages
	Fix       bool
	Rev       string
	PatchFile string
	Executor  Executor
	Stdout    io.Writer
	Stderr    io.Writer
}

// Run iterates over affected module targets and runs golangci-lint on the changed packages in parallel.
func Run(ctx context.Context, cfg Config) error {
	if len(cfg.Targets) == 0 {
		return nil
	}

	if cfg.Stdout == nil {
		cfg.Stdout = io.Discard
	}
	if cfg.Stderr == nil {
		cfg.Stderr = io.Discard
	}

	execRunner := cfg.Executor
	if execRunner == nil {
		execRunner = &osExecutor{stdout: cfg.Stdout, stderr: cfg.Stderr}
	}

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, target := range cfg.Targets {
		modDir := cfg.RepoRoot
		if target.Module != "." && target.Module != "" {
			modDir = filepath.Join(cfg.RepoRoot, target.Module)
		}

		var args []string
		args = append(args, "run", "--allow-parallel-runners")
		if cfg.PatchFile != "" {
			args = append(args, "--new-from-patch="+cfg.PatchFile)
		} else if cfg.Rev != "" {
			args = append(args, "--new-from-rev="+cfg.Rev)
		}
		if cfg.Fix {
			args = append(args, "--fix")
		}
		args = append(args, target.Packages...)

		wg.Add(1)
		go func(t modules.ModulePackages, dir string, cmdArgs []string) {
			defer wg.Done()

			var outBuf bytes.Buffer
			fmt.Fprintf(&outBuf, "==> Linting module '%s' packages: %s\n", t.Module, strings.Join(t.Packages, " "))

			err := execRunner.Run(ctx, dir, "golangci-lint", cmdArgs...)

			mu.Lock()
			cfg.Stdout.Write(outBuf.Bytes())
			if err != nil {
				errs = append(errs, fmt.Errorf("golangci-lint failed in %s: %w", t.Module, err))
			}
			mu.Unlock()
		}(target, modDir, args)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
