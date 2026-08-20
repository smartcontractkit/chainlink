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

// createModulePatch creates a temporary git patch for the given module directory and packages.
func createModulePatch(ctx context.Context, dir string, rev string, pkgs []string) (string, error) {
	if rev == "" || len(pkgs) == 0 {
		return "", nil
	}

	args := make([]string, 0, 4+len(pkgs))
	args = append(args, "diff", "--relative", rev, "--")
	args = append(args, pkgs...)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil || out.Len() == 0 {
		return "", nil
	}

	tmpFile, err := os.CreateTemp("", "githooks-lint-*.patch")
	if err != nil {
		return "", err
	}

	if _, err := tmpFile.Write(out.Bytes()); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return "", err
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
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

		wg.Add(1)
		go func(t modules.ModulePackages, dir string) {
			defer wg.Done()

			var patchPath string
			if cfg.PatchFile != "" {
				patchPath = cfg.PatchFile
			} else if cfg.Rev != "" {
				patchPath, _ = createModulePatch(ctx, dir, cfg.Rev, t.Packages)
				if patchPath != "" {
					defer func(p string) {
						_ = os.Remove(p)
					}(patchPath)
				}
			}

			var args []string
			args = append(args, "run", "--allow-parallel-runners")
			if patchPath != "" {
				args = append(args, "--new-from-patch="+patchPath)
			} else if cfg.Rev != "" {
				args = append(args, "--new-from-rev="+cfg.Rev)
			}
			if cfg.Fix {
				args = append(args, "--fix")
			}
			args = append(args, t.Packages...)

			mu.Lock()
			fmt.Fprintf(cfg.Stdout, "==> Linting module '%s' packages: %s\n", t.Module, strings.Join(t.Packages, " "))
			mu.Unlock()

			err := execRunner.Run(ctx, dir, "golangci-lint", args...)

			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("golangci-lint failed in %s: %w", t.Module, err))
				mu.Unlock()
			}
		}(target, modDir)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
