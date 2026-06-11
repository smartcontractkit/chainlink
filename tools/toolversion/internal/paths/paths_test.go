package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepoRootFromSubdir(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n\ngo 1.26.4\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".tool-versions"), []byte("golang 1.26.4\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(root, "integration-tests")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(sub)

	got, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("findRepoRoot() = %q, want %q", got, root)
	}
}
