package engine

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gitCmd runs git in dir with a fixed identity, failing the test on error.
func gitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	full := append([]string{"-c", "user.name=Test", "-c", "user.email=test@example.com",
		"-c", "init.defaultBranch=main", "-c", "commit.gpgsign=false"}, args...)
	cmd := exec.CommandContext(context.Background(), "git", full...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

// setupOriginOnlyBranch creates two repos: "remote" with a commit on branch
// release/9.99.0, and "local" which fetches it — so release/9.99.0 exists in
// local only as refs/remotes/origin/release/9.99.0, not as a local branch.
// Returns (localDir, remoteHEAD).
func setupOriginOnlyBranch(t *testing.T) (string, string) {
	t.Helper()
	base := t.TempDir()
	remote := filepath.Join(base, "remote")
	local := filepath.Join(base, "local")

	for _, dir := range []string{remote, local} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		gitCmd(t, dir, "init")
	}
	if err := os.WriteFile(filepath.Join(remote, "go.mod"), []byte("module example.com/x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitCmd(t, remote, "add", ".")
	gitCmd(t, remote, "commit", "-m", "init")
	gitCmd(t, remote, "branch", "release/9.99.0")

	gitCmd(t, local, "remote", "add", "origin", remote)
	gitCmd(t, local, "fetch", "origin")

	// Sanity: no local branch named release/9.99.0.
	cmd := exec.CommandContext(context.Background(), "git", "rev-parse", "--verify", "release/9.99.0^{commit}")
	cmd.Dir = local
	if err := cmd.Run(); err == nil {
		t.Fatal("test setup invalid: local branch release/9.99.0 unexpectedly exists")
	}

	out, err := exec.CommandContext(context.Background(), "git", "-C", remote, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	sha := string(out[:40])
	return local, sha
}

func TestResolveRef_FallsBackToOrigin(t *testing.T) {
	t.Parallel()

	local, wantSHA := setupOriginOnlyBranch(t)
	g := gitRunner{dir: local}

	got, err := g.ResolveRef(context.Background(), "release/9.99.0")
	if err != nil {
		t.Fatalf("ResolveRef failed: %v", err)
	}
	if got != wantSHA {
		t.Errorf("ResolveRef = %s, want %s", got, wantSHA)
	}
}

func TestResolveRef_UnknownRef(t *testing.T) {
	t.Parallel()

	local, _ := setupOriginOnlyBranch(t)
	g := gitRunner{dir: local}
	if _, err := g.ResolveRef(context.Background(), "release/does-not-exist"); err == nil {
		t.Error("expected error for unknown ref")
	}
}

// TestResolveRef_FetchesMissingTag creates a tag in the remote AFTER the
// local clone's initial fetch: ResolveRef must fetch it on demand.
func TestResolveRef_FetchesMissingTag(t *testing.T) {
	t.Parallel()

	local, remoteSHA := setupOriginOnlyBranch(t)
	remote := filepath.Join(filepath.Dir(local), "remote")
	gitCmd(t, remote, "tag", "v9.99.1")

	g := gitRunner{dir: local}
	got, err := g.ResolveRef(context.Background(), "v9.99.1")
	if err != nil {
		t.Fatalf("ResolveRef failed: %v", err)
	}
	if got != remoteSHA {
		t.Errorf("ResolveRef = %s, want %s", got, remoteSHA)
	}
}

// TestResolveRef_SuggestsNearMiss checks the error hints at similarly named
// remote refs (e.g. user passed "9.99.0" while "release/9.99.0" exists).
func TestResolveRef_SuggestsNearMiss(t *testing.T) {
	t.Parallel()

	local, _ := setupOriginOnlyBranch(t)
	g := gitRunner{dir: local}
	_, err := g.ResolveRef(context.Background(), "9.99.0")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "release/9.99.0") {
		t.Errorf("error should suggest release/9.99.0, got: %v", err)
	}
}
