package coverage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PublishBadgeOptions configures publication of the two files on the orphan
// coverage branch.
type PublishBadgeOptions struct {
	RepoRoot    string
	WorktreeDir string
	SVGPath     string
	StatePath   string
	SHA         string
	Branch      string
}

// PublishBadge commits and pushes only latest/badge.svg and latest/report.json.
// It returns false when the existing branch already contains those bytes.
func PublishBadge(ctx context.Context, opts PublishBadgeOptions) (bool, error) {
	if opts.RepoRoot == "" || opts.WorktreeDir == "" || opts.SVGPath == "" || opts.StatePath == "" || opts.SHA == "" {
		return false, errors.New("repo root, worktree, SVG, state, and SHA are required")
	}
	if opts.Branch == "" {
		opts.Branch = "coverage"
	}
	if err := os.RemoveAll(opts.WorktreeDir); err != nil {
		return false, err
	}
	if _, err := commandOutput(ctx, opts.RepoRoot, "git", "worktree", "prune"); err != nil {
		return false, err
	}

	remoteRef := "refs/remotes/origin/" + opts.Branch
	_, remoteErr := commandOutput(ctx, opts.RepoRoot, "git", "ls-remote", "--exit-code", "--heads", "origin", opts.Branch)
	if remoteErr == nil {
		if _, err := commandOutput(ctx, opts.RepoRoot, "git", "fetch", "--no-tags", "origin", opts.Branch+":"+remoteRef); err != nil {
			return false, err
		}
		if _, err := commandOutput(ctx, opts.RepoRoot, "git", "worktree", "add", "--detach", opts.WorktreeDir, remoteRef); err != nil {
			return false, err
		}
	} else {
		if _, err := commandOutput(ctx, opts.RepoRoot, "git", "worktree", "add", "--detach", opts.WorktreeDir); err != nil {
			return false, err
		}
		if _, err := commandOutput(ctx, opts.WorktreeDir, "git", "checkout", "--orphan", opts.Branch); err != nil {
			return false, err
		}
		// The initial orphan checkout can contain the caller's tracked tree.
		// git rm exits non-zero for an empty index, which is harmless here.
		cmd := exec.CommandContext(ctx, "git", "rm", "-rf", ".")
		cmd.Dir = opts.WorktreeDir
		_, _ = cmd.CombinedOutput()
	}
	defer func() {
		cmd := exec.CommandContext(context.Background(), "git", "worktree", "remove", "--force", opts.WorktreeDir)
		cmd.Dir = opts.RepoRoot
		_, _ = cmd.CombinedOutput()
	}()

	latest := filepath.Join(opts.WorktreeDir, "latest")
	if err := os.MkdirAll(latest, 0o755); err != nil {
		return false, err
	}
	if err := copyCoverageFile(opts.SVGPath, filepath.Join(latest, "badge.svg")); err != nil {
		return false, err
	}
	if err := copyCoverageFile(opts.StatePath, filepath.Join(latest, "report.json")); err != nil {
		return false, err
	}
	if _, err := commandOutput(ctx, opts.WorktreeDir, "git", "add", "latest"); err != nil {
		return false, err
	}
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--quiet")
	cmd.Dir = opts.WorktreeDir
	if err := cmd.Run(); err == nil {
		return false, nil
	} else if exit := new(exec.ExitError); !errors.As(err, &exit) || exit.ExitCode() != 1 {
		return false, fmt.Errorf("check badge diff: %w", err)
	}
	if _, err := commandOutput(ctx, opts.WorktreeDir, "git",
		"-c", "commit.gpgsign=false",
		"-c", "user.name=github-actions[bot]",
		"-c", "user.email=41898282+github-actions[bot]@users.noreply.github.com",
		"commit", "-m", "coverage: update badge for "+opts.SHA); err != nil {
		return false, err
	}
	if _, err := commandOutput(ctx, opts.WorktreeDir, "git", "push", "origin", "HEAD:refs/heads/"+opts.Branch); err != nil {
		return false, err
	}
	return true, nil
}

func copyCoverageFile(source, destination string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(filepath.Base(source), ".svg") && !strings.HasSuffix(filepath.Base(source), ".json") {
		return fmt.Errorf("unsupported badge file %q", source)
	}
	return os.WriteFile(destination, data, 0o600)
}
