package generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// CommandRunner executes a command in a directory.
type CommandRunner func(ctx context.Context, dir string, args ...string) error

// Config holds options for running generate.
type Config struct {
	Runner CommandRunner
}

func defaultRunner(ctx context.Context, dir string, args ...string) error {
	if len(args) > 0 && args[0] == "modgraph" {
		modgraphScript := filepath.Join(dir, "tools/bin/modgraph")
		cmd := exec.CommandContext(ctx, modgraphScript)
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("modgraph failed: %w", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "go.md"), out, 0o600); err != nil {
			return fmt.Errorf("failed to write go.md: %w", err)
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go %v in %s failed: %w (output: %s)", args, dir, err, string(out))
	}
	return nil
}

// Run maps changed files to specific code generator targets and runs them.
func Run(ctx context.Context, repoRoot string, files []string, cfg ...Config) error {
	if len(files) == 0 {
		return nil
	}

	runner := defaultRunner
	if len(cfg) > 0 && cfg[0].Runner != nil {
		runner = cfg[0].Runner
	}

	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get absolute repo root: %w", err)
	}

	protoPkgs := make(map[string]struct{})
	runConfigDocs := false
	runModGraph := false

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		absFile := file
		if !filepath.IsAbs(absFile) {
			absFile = filepath.Join(absRoot, file)
		}

		relFromRoot, relErr := filepath.Rel(absRoot, absFile)
		if relErr != nil {
			return fmt.Errorf("failed to get relative path from repo root: %w", relErr)
		}
		cleanRel := filepath.ToSlash(filepath.Clean(relFromRoot))
		cleanRel = strings.TrimPrefix(cleanRel, "./")

		baseName := filepath.Base(cleanRel)

		// Check if go.mod / go.sum changed -> trigger go.md modgraph generation
		if baseName == "go.mod" || baseName == "go.sum" {
			runModGraph = true
		}

		// Check if config docs generator should run
		if strings.HasPrefix(cleanRel, "core/config") {
			runConfigDocs = true
		}

		// Check proto or generate files
		if strings.HasSuffix(cleanRel, ".proto") ||
			baseName == "generate.go" ||
			baseName == "gen.go" {
			pkgDir := filepath.Dir(cleanRel)
			pkgPattern := "./" + pkgDir
			protoPkgs[pkgPattern] = struct{}{}
		}
	}

	var targets []string
	for pkg := range protoPkgs {
		targets = append(targets, pkg)
	}
	sort.Strings(targets)

	if len(targets) == 0 && !runConfigDocs && !runModGraph {
		return nil
	}

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, pkg := range targets {
		wg.Go(func() {
			if err := runner(ctx, absRoot, "generate", pkg); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("go generate %s: %w", pkg, err))
				mu.Unlock()
			}
		})
	}

	if runConfigDocs {
		wg.Go(func() {
			if err := runner(ctx, absRoot, "run", "./core/config/docs/cmd/generate", "-o", "./docs/"); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("generate config docs: %w", err))
				mu.Unlock()
			}
		})
	}

	if runModGraph {
		wg.Go(func() {
			if err := runner(ctx, absRoot, "modgraph"); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("generate go.md: %w", err))
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
