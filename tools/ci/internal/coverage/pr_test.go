package coverage

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProductionGoFile(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"core/services/foo/foo.go":             true,
		"core/services/foo/foo_test.go":        false,
		"core/services/foo/testdata/input.go":  false,
		"core/services/foo/mocks/mock.go":      false,
		"core/services/foo/types.pb.go":        false,
		"core/services/foo/types_gen.go":       false,
		"core/services/foo/types_codecgen.go":  false,
		"vendor/example.com/foo/foo.go":        false,
		"README.md":                            false,
		"core/services/foo/testhelpers/foo.go": false,
	}
	for path, want := range tests {
		if got := productionGoFile(path); got != want {
			t.Errorf("productionGoFile(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestRunPRWithoutGoChanges(t *testing.T) {
	t.Parallel()
	repo := initTestRepo(t)
	writePRTestFile(t, repo, "README.md", "base\n")
	runGitTest(t, repo, "add", "README.md")
	runGitTest(t, repo, "commit", "-m", "base")
	base := strings.TrimSpace(runGitTest(t, repo, "rev-parse", "HEAD"))
	writePRTestFile(t, repo, "README.md", "head\n")
	runGitTest(t, repo, "commit", "-am", "docs")
	head := strings.TrimSpace(runGitTest(t, repo, "rev-parse", "HEAD"))

	result, err := RunPR(t.Context(), PROptions{
		RepoRoot: repo,
		BaseSHA:  base,
		HeadSHA:  head,
		OutDir:   filepath.Join(repo, "coverage-patch"),
		Short:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.NoGoChanges || result.TestsFailed || result.Report != nil {
		t.Fatalf("unexpected result: %+v", result)
	}
	if _, err := os.Stat(filepath.Join(repo, "coverage-patch", "result.json")); err != nil {
		t.Fatalf("result was not written: %v", err)
	}
}

func TestRenderPRCommentIncludesNightlyArtifacts(t *testing.T) {
	t.Parallel()
	result := PRResult{
		HeadSHA:     "0123456789abcdef",
		NoGoChanges: true,
	}
	nightly := []Artifact{{
		Name:      "coverage-report-develop-1",
		URL:       "https://github.com/smartcontractkit/chainlink/actions/runs/1/artifacts/10",
		RunURL:    "https://github.com/smartcontractkit/chainlink/actions/runs/1",
		SHA:       "abcdef0123456789",
		CreatedAt: time.Date(2026, 11, 26, 3, 0, 0, 0, time.UTC),
	}, {
		Name: "coverage-core-go_core_race_tests-2-1",
		URL:  "https://github.com/smartcontractkit/chainlink/actions/runs/2/artifacts/20",
	}}
	comment := RenderPRComment(result, nightly, "https://example.test/pr-run", "https://example.test/pr-artifact")
	for _, want := range []string{
		"<!-- coverage-patch -->",
		"No coverable Go changes",
		"baseline only; does not include this PR",
		"Full coverage report",
		"Race coverage (separate; never merged)",
		"https://example.test/pr-artifact",
		"abcdef012345",
	} {
		if !strings.Contains(comment, want) {
			t.Errorf("comment does not contain %q:\n%s", want, comment)
		}
	}
}

func TestRenderPRCommentMarksFailedTests(t *testing.T) {
	t.Parallel()
	comment := RenderPRComment(PRResult{HeadSHA: "abc", TestsFailed: true}, nil, "", "")
	if !strings.Contains(comment, "some affected-package tests failed") {
		t.Fatalf("missing test failure warning: %s", comment)
	}
	if !strings.Contains(comment, "No unexpired nightly artifacts") {
		t.Fatalf("missing empty nightly state: %s", comment)
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "coverage@example.test")
	runGitTest(t, dir, "config", "user.name", "Coverage Test")
	runGitTest(t, dir, "config", "commit.gpgsign", "false")
	return dir
}

func runGitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func writePRTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
