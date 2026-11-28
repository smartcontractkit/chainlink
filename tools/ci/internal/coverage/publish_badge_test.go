package coverage

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeBadgeTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestPublishBadgeCreatesAndUpdatesOrphanBranch(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	repo := filepath.Join(root, "repo")
	for _, args := range [][]string{{"init", "--bare", remote}, {"init", repo}} {
		cmd := exec.Command("git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	for _, args := range [][]string{
		{"config", "user.email", "coverage@example.test"},
		{"config", "user.name", "Coverage Test"},
		{"config", "commit.gpgsign", "false"},
	} {
		runGitTest(t, repo, args...)
	}
	writeBadgeTestFile(t, repo, "README.md", "test\n")
	runGitTest(t, repo, "add", "README.md")
	runGitTest(t, repo, "commit", "-m", "base")
	runGitTest(t, repo, "remote", "add", "origin", remote)
	runGitTest(t, repo, "push", "-u", "origin", "HEAD:develop")

	svg := filepath.Join(root, "badge.svg")
	state := filepath.Join(root, "report.json")
	if err := os.WriteFile(svg, []byte("<svg>one</svg>\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(state, []byte("{\"sha\":\"one\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts := PublishBadgeOptions{
		RepoRoot:    repo,
		WorktreeDir: filepath.Join(root, "coverage-worktree"),
		SVGPath:     svg,
		StatePath:   state,
		SHA:         "one",
	}
	changed, err := PublishBadge(t.Context(), opts)
	if err != nil || !changed {
		t.Fatalf("first publish changed=%v err=%v", changed, err)
	}
	cmd := exec.Command("git", "--git-dir", remote, "show", "coverage:latest/report.json")
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "\"one\"") {
		t.Fatalf("published state=%s err=%v", out, err)
	}

	changed, err = PublishBadge(t.Context(), opts)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("unchanged badge should not create a commit")
	}

	if err := os.WriteFile(state, []byte("{\"sha\":\"two\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts.SHA = "two"
	changed, err = PublishBadge(t.Context(), opts)
	if err != nil || !changed {
		t.Fatalf("second publish changed=%v err=%v", changed, err)
	}
}
