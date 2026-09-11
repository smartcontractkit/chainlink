package lint

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/ui"
)

// Executor abstracts command execution for testability.
type Executor interface {
	Run(ctx context.Context, dir, name string, args ...string) error
}

type osExecutor struct {
	stdout io.Writer
	stderr io.Writer
}

func (e *osExecutor) Run(ctx context.Context, dir, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Config holds configuration for the linter runner.
type Config struct {
	RepoRoot string
	Targets  []modules.ModulePackages
	Fix      bool
	Rev      string
	Executor Executor
	Stdout   io.Writer
	Stderr   io.Writer
	UI       *ui.UI
}

// Run iterates over affected module targets and runs golangci-lint on the
// changed packages sequentially, letting golangci-lint parallelize internally
// via --allow-parallel-runners.
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

	u := cfg.UI
	if u == nil {
		u = ui.New(cfg.Stdout)
	}

	var errs []error
	for _, target := range cfg.Targets {
		modDir := cfg.RepoRoot
		if target.Module != "." && target.Module != "" {
			modDir = filepath.Join(cfg.RepoRoot, target.Module)
		}

		var args []string
		args = append(args, "run", "--allow-parallel-runners")
		if cfg.Rev != "" {
			args = append(args, "--new-from-rev="+cfg.Rev)
		}
		if cfg.Fix {
			args = append(args, "--fix")
		}
		args = append(args, target.Packages...)

		pkgStr := strings.Join(target.Packages, " ")
		if len(target.Packages) > 4 {
			pkgStr = fmt.Sprintf("%d packages", len(target.Packages))
		}
		fmt.Fprintln(cfg.Stdout, u.RunnerHeader("LINT", target.Module, pkgStr))
		if err := execRunner.Run(ctx, modDir, "golangci-lint", args...); err != nil {
			// Keep linting remaining modules so their fixes are still applied.
			errs = append(errs, fmt.Errorf("golangci-lint failed in %s: %w", target.Module, err))
		}
	}

	return errors.Join(errs...)
}
