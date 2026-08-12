package changelog

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// gitRunner executes git commands in the given repository directory.
type gitRunner struct {
	dir string
}

func (g gitRunner) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = g.dir
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errBuf.String()))
	}
	return out.String(), nil
}

// ResolveRef resolves any git ref (SHA, tag, branch) to a commit SHA.
// If the ref does not resolve locally, it falls back to the remote-tracking
// branch origin/<ref> (release branches often exist only remotely).
func (g gitRunner) ResolveRef(ctx context.Context, ref string) (string, error) {
	out, err := g.run(ctx, "rev-parse", "--verify", ref+"^{commit}")
	if err == nil {
		return strings.TrimSpace(out), nil
	}
	if !strings.HasPrefix(ref, "origin/") && !strings.HasPrefix(ref, "refs/") {
		if out, ferr := g.run(ctx, "rev-parse", "--verify", "origin/"+ref+"^{commit}"); ferr == nil {
			return strings.TrimSpace(out), nil
		}
	}
	return "", fmt.Errorf("resolving ref %q: %w", ref, err)
}

// FileAtRef returns the contents of path at the given ref.
func (g gitRunner) FileAtRef(ctx context.Context, ref, path string) ([]byte, error) {
	out, err := g.run(ctx, "show", ref+":"+path)
	if err != nil {
		return nil, fmt.Errorf("reading %s at %s: %w", path, ref, err)
	}
	return []byte(out), nil
}

// localCommit is one commit from the local git log.
type localCommit struct {
	SHA         string
	AuthorName  string
	AuthorEmail string
	Title       string
}

// LogRange lists commits in old..new (commits reachable from new but not old).
func (g gitRunner) LogRange(ctx context.Context, oldSHA, newSHA string) ([]localCommit, error) {
	out, err := g.run(ctx, "log", "--format=%H%x00%an%x00%ae%x00%s", oldSHA+".."+newSHA)
	if err != nil {
		return nil, err
	}
	var commits []localCommit
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 4)
		if len(parts) != 4 {
			continue
		}
		commits = append(commits, localCommit{
			SHA:         parts[0],
			AuthorName:  parts[1],
			AuthorEmail: parts[2],
			Title:       parts[3],
		})
	}
	return commits, nil
}

// IsAncestor reports whether oldSHA is an ancestor of newSHA.
func (g gitRunner) IsAncestor(ctx context.Context, oldSHA, newSHA string) bool {
	cmd := exec.CommandContext(ctx, "git", "merge-base", "--is-ancestor", oldSHA, newSHA)
	cmd.Dir = g.dir
	return cmd.Run() == nil
}

// CommitFiles returns the file paths touched by a single commit.
// Merge commits are diffed against their first parent.
func (g gitRunner) CommitFiles(ctx context.Context, sha string) ([]string, error) {
	out, err := g.run(ctx, "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return nil, err
	}
	paths := splitLines(out)
	if len(paths) == 0 {
		// Possibly a merge commit: diff against first parent.
		out, err = g.run(ctx, "diff-tree", "--first-parent", "-m", "--no-commit-id", "--name-only", "-r", sha)
		if err != nil {
			return nil, err
		}
		paths = splitLines(out)
	}
	return paths, nil
}

func splitLines(s string) []string {
	var out []string
	for line := range strings.SplitSeq(strings.TrimSpace(s), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
