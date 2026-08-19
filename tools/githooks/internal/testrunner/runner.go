package testrunner

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

// Config holds configuration for the test runner.
type Config struct {
	RepoRoot string
	Modules  []modules.ModulePackages
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
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s: %w", binPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", filepath.Dir(binPath), err)
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

// Run executes test harness grouped by Go module.
func Run(ctx context.Context, cfg Config) error {
	mods := cfg.Modules
	if len(mods) == 0 && len(cfg.Packages) > 0 {
		mods = []modules.ModulePackages{
			{
				Module:   ".",
				Packages: cfg.Packages,
			},
		}
	}

	if len(mods) == 0 {
		return nil
	}

	execRunner := cfg.Executor
	if execRunner == nil {
		execRunner = &osExecutor{stdout: cfg.Stdout, stderr: cfg.Stderr}
	}

	for _, mod := range mods {
		if len(mod.Packages) == 0 {
			continue
		}

		if mod.Module == "." {
			if cfg.Executor == nil {
				if err := EnsureBinary(ctx, cfg.RepoRoot, cfg.Stdout, cfg.Stderr); err != nil {
					return err
				}
			}
			binPath := filepath.Join(cfg.RepoRoot, "tools/test/.bin/test")
			var args []string
			if cfg.Short {
				args = append(args, "-short")
			}
			args = append(args, cfg.Args...)
			args = append(args, mod.Packages...)

			fmt.Fprintf(cfg.Stdout, "==> Running tests on %s: %s\n", mod.Module, strings.Join(mod.Packages, " "))
			if err := execRunner.Run(ctx, cfg.RepoRoot, binPath, args...); err != nil {
				return fmt.Errorf("tests failed on %s: %w", mod.Module, err)
			}
		} else {
			modDir := filepath.Join(cfg.RepoRoot, mod.Module)
			var args []string
			args = append(args, "test")
			if cfg.Short {
				args = append(args, "-short")
			}
			args = append(args, cfg.Args...)
			args = append(args, mod.Packages...)

			fmt.Fprintf(cfg.Stdout, "==> Running tests on %s: %s\n", mod.Module, strings.Join(mod.Packages, " "))
			if err := execRunner.Run(ctx, modDir, "go", args...); err != nil {
				return fmt.Errorf("tests failed on %s: %w", mod.Module, err)
			}
		}
	}

	return nil
}
