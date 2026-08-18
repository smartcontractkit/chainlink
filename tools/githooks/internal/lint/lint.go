package lint

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	RepoRoot string
	Targets  []modules.ModulePackages
	Fix      bool
	Rev      string
	Executor Executor
	Stdout   io.Writer
	Stderr   io.Writer
}

// Run iterates over affected module targets and runs golangci-lint on the changed packages.
func Run(ctx context.Context, cfg Config) error {
	if len(cfg.Targets) == 0 {
		return nil
	}

	execRunner := cfg.Executor
	if execRunner == nil {
		execRunner = &osExecutor{stdout: cfg.Stdout, stderr: cfg.Stderr}
	}

	for _, target := range cfg.Targets {
		modDir := cfg.RepoRoot
		if target.Module != "." && target.Module != "" {
			modDir = filepath.Join(cfg.RepoRoot, target.Module)
		}

		var args []string
		args = append(args, "run")
		if cfg.Rev != "" {
			args = append(args, "--new-from-rev="+cfg.Rev)
		}
		if cfg.Fix {
			args = append(args, "--fix")
		}
		args = append(args, target.Packages...)

		fmt.Fprintf(cfg.Stdout, "==> Linting module '%s' packages: %s\n", target.Module, strings.Join(target.Packages, " "))
		if err := execRunner.Run(ctx, modDir, "golangci-lint", args...); err != nil {
			return fmt.Errorf("golangci-lint failed in %s: %w", target.Module, err)
		}
	}

	return nil
}
