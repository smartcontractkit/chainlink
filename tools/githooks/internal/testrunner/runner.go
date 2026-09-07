package testrunner

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
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

// needsBuild reports whether the test harness binary must be (re)built: the
// binary is missing, or any .go file under srcDir is newer than it.
func needsBuild(binPath, srcDir string) (bool, error) {
	binInfo, statErr := os.Stat(binPath)
	switch {
	case statErr == nil:
	case os.IsNotExist(statErr):
		return true, nil
	default:
		return false, fmt.Errorf("failed to stat %s: %w", binPath, statErr)
	}

	var newest time.Time
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if d.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return infoErr
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to scan %s: %w", srcDir, err)
	}

	return newest.After(binInfo.ModTime()), nil
}

// EnsureBinary checks if tools/test/.bin/test exists and is newer than the
// tools/test sources, and builds it if missing or stale.
func EnsureBinary(ctx context.Context, repoRoot string, stdout, stderr io.Writer) error {
	binPath := filepath.Join(repoRoot, "tools/test/.bin/test")
	srcDir := filepath.Join(repoRoot, "tools/test")

	rebuild, err := needsBuild(binPath, srcDir)
	if err != nil {
		return err
	}
	if !rebuild {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", filepath.Dir(binPath), err)
	}

	cmd := exec.CommandContext(ctx, "go", "build", "-o", ".bin/test", ".")
	cmd.Dir = srcDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build tools/test harness: %w", err)
	}
	return nil
}

// Run executes test harness grouped by Go module sequentially, letting
// `go test` parallelize test execution within each module.
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

	for _, mod := range mods {
		if len(mod.Packages) == 0 {
			continue
		}

		fmt.Fprintf(cfg.Stdout, "==> Running tests on %s: %s\n", mod.Module, strings.Join(mod.Packages, " "))

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

			if err := execRunner.Run(ctx, modDir, "go", args...); err != nil {
				return fmt.Errorf("tests failed on %s: %w", mod.Module, err)
			}
		}
	}

	return nil
}
