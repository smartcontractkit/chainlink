package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootFrom_chainlinkModule(t *testing.T) {
	// This test runs from package dir; walk up to module root (tools/test), then repo root.
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// internal/repo -> internal -> tools/test
	testMod := filepath.Clean(filepath.Join(here, "..", ".."))
	root, err := RootFrom(testMod)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(root) != "chainlink" && !strings.HasSuffix(root, "chainlink") {
		t.Fatalf("unexpected root %q", root)
	}
	// go.mod at root must exist
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatal(err)
	}
}
