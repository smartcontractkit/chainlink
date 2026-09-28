package actionlint

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executor abstracts command execution and output capture for testability.
type Executor interface {
	Run(ctx context.Context, dir, name string, args ...string) error
	Output(ctx context.Context, dir, name string, args ...string) ([]byte, error)
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

func (e *osExecutor) Output(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// Config holds configuration for the actionlint runner.
type Config struct {
	RepoRoot string
	Files    []string
	Executor Executor
	Stdout   io.Writer
	Stderr   io.Writer
}

// IsGithubYAML reports whether relPath is a YAML file within the .github directory.
func IsGithubYAML(relPath string) bool {
	clean := filepath.ToSlash(filepath.Clean(relPath))
	clean = strings.TrimPrefix(clean, "./")
	clean = strings.TrimPrefix(clean, "/")

	if !strings.HasPrefix(clean, ".github/") {
		return false
	}

	ext := strings.ToLower(filepath.Ext(clean))
	return ext == ".yml" || ext == ".yaml"
}

// IsWorkflowYAML reports whether relPath is a YAML file within .github/workflows.
func IsWorkflowYAML(relPath string) bool {
	clean := filepath.ToSlash(filepath.Clean(relPath))
	clean = strings.TrimPrefix(clean, "./")
	clean = strings.TrimPrefix(clean, "/")

	if !strings.HasPrefix(clean, ".github/workflows/") {
		return false
	}

	ext := strings.ToLower(filepath.Ext(clean))
	return ext == ".yml" || ext == ".yaml"
}

// filterGithubYAML classifies input files into workflow YAMLs and non-workflow .github YAMLs.
func filterGithubYAML(files []string) (workflowFiles []string, hasNonWorkflowYAML, hasAnyGithubYAML bool) {
	for _, file := range files {
		trimmed := strings.TrimSpace(file)
		if trimmed == "" {
			continue
		}

		if !IsGithubYAML(trimmed) {
			continue
		}

		hasAnyGithubYAML = true
		if IsWorkflowYAML(trimmed) {
			workflowFiles = append(workflowFiles, trimmed)
		} else {
			hasNonWorkflowYAML = true
		}
	}
	return workflowFiles, hasNonWorkflowYAML, hasAnyGithubYAML
}

// validateFork ensures actionlint is available and belongs to the maintained kjanat/actionlint fork.
func validateFork(ctx context.Context, execRunner Executor, repoRoot string) error {
	out, err := execRunner.Output(ctx, repoRoot, "actionlint", "-help")
	if err != nil {
		return fmt.Errorf("actionlint not found or failed: %w. Install maintained fork: go install actionlint.kjanat.dev/cmd/actionlint@latest (https://github.com/kjanat/actionlint)", err)
	}

	outStr := string(out)
	if strings.Contains(outStr, "rhysd/actionlint") || !strings.Contains(outStr, "kjanat") {
		return errors.New("unmaintained actionlint detected (rhysd/actionlint). Upgrade to maintained fork: go install actionlint.kjanat.dev/cmd/actionlint@latest (https://github.com/kjanat/actionlint)")
	}

	return nil
}

// Run filters files, validates the actionlint fork, and runs actionlint on changed workflows.
func Run(ctx context.Context, cfg Config) error {
	if len(cfg.Files) == 0 {
		return nil
	}

	workflowFiles, hasNonWorkflowYAML, hasAnyGithubYAML := filterGithubYAML(cfg.Files)
	if !hasAnyGithubYAML {
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

	if err := validateFork(ctx, execRunner, cfg.RepoRoot); err != nil {
		return err
	}

	var args []string
	if !hasNonWorkflowYAML && len(workflowFiles) > 0 {
		args = workflowFiles
		fmt.Fprintf(cfg.Stdout, "==> Running actionlint on: %s\n", strings.Join(args, " "))
	} else {
		fmt.Fprintf(cfg.Stdout, "==> Running actionlint on all workflows\n")
	}

	if err := execRunner.Run(ctx, cfg.RepoRoot, "actionlint", args...); err != nil {
		return fmt.Errorf("actionlint failed: %w", err)
	}

	return nil
}
