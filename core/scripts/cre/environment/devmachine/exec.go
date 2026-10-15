package devmachine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Runner shells out to external tools (aws, kubectl, docker, rsync, ssh, go).
// It centralises dry-run handling and command logging so every phase behaves
// consistently.
type Runner struct {
	Dir     string
	DryRun  bool
	Verbose bool
	Out     io.Writer
	Err     io.Writer
}

// NewRunner constructs a Runner writing to the provided streams.
func NewRunner(out, errOut io.Writer, dryRun, verbose bool) *Runner {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &Runner{DryRun: dryRun, Verbose: verbose, Out: out, Err: errOut}
}

// format renders a command for logging.
func format(name string, args []string) string {
	parts := append([]string{name}, args...)
	quoted := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" || strings.ContainsAny(p, " \t\n\"'") {
			quoted = append(quoted, fmt.Sprintf("%q", p))
			continue
		}
		quoted = append(quoted, p)
	}
	return strings.Join(quoted, " ")
}

// Run executes the command, capturing combined output. In dry-run mode it logs
// the command and returns an empty string.
func (r *Runner) Run(ctx context.Context, name string, args ...string) (string, error) {
	return r.RunEnv(ctx, nil, name, args...)
}

// RunEnv is Run with extra environment variables.
func (r *Runner) RunEnv(ctx context.Context, env []string, name string, args ...string) (string, error) {
	if r.DryRun {
		fmt.Fprintf(r.Out, "[dry-run] %s\n", format(name, args))
		return "", nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.Dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if r.Verbose {
		fmt.Fprintf(r.Out, "$ %s\n", format(name, args))
	}
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		return stdout.String(), errors.Wrapf(err, "command failed: %s\n%s", format(name, args), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Stream executes the command, streaming stdout/stderr to the runner's output.
func (r *Runner) Stream(ctx context.Context, name string, args ...string) error {
	return r.StreamEnv(ctx, nil, name, args...)
}

// StreamEnv is Stream with extra environment variables.
func (r *Runner) StreamEnv(ctx context.Context, env []string, name string, args ...string) error {
	if r.DryRun {
		fmt.Fprintf(r.Out, "[dry-run] %s\n", format(name, args))
		return nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.Dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Stdout = r.Out
	cmd.Stderr = r.Err
	if r.Verbose {
		fmt.Fprintf(r.Out, "$ %s\n", format(name, args))
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "command failed: %s", format(name, args))
	}
	return nil
}

// StreamStdin executes the command feeding stdin from the provided string,
// streaming stdout/stderr to the runner's output.
func (r *Runner) StreamStdin(ctx context.Context, stdin, name string, args ...string) error {
	if r.DryRun {
		fmt.Fprintf(r.Out, "[dry-run] %s (with stdin)\n", format(name, args))
		return nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.Dir
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = r.Out
	cmd.Stderr = r.Err
	if r.Verbose {
		fmt.Fprintf(r.Out, "$ %s\n", format(name, args))
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "command failed: %s", format(name, args))
	}
	return nil
}

// In returns a copy of the runner that executes commands in dir.
func (r *Runner) In(dir string) *Runner {
	cp := *r
	cp.Dir = dir
	return &cp
}

// LookPath reports whether a binary is available on PATH.
func (r *Runner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}
