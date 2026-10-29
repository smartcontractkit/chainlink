package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBlockAndNormalize(t *testing.T) {
	b, err := parseBlock("github.com/smartcontractkit/chainlink/v2/core/foo.go:10.2,12.4 3 7")
	if err != nil {
		t.Fatal(err)
	}
	if b.File != "core/foo.go" || b.StartLine != 10 || b.EndLine != 12 || b.NumStmt != 3 || b.Count != 7 {
		t.Fatalf("unexpected block: %+v", b)
	}
	if got := normalizePath("github.com/smartcontractkit/chainlink/deployment/foo.go"); got != "deployment/foo.go" {
		t.Fatalf("unexpected deployment path: %q", got)
	}
}

func TestSummarizeAndExclude(t *testing.T) {
	blocks := map[string]block{
		"a": {File: "core/a.go", NumStmt: 4, Count: 1},
		"b": {File: "core/a.go", NumStmt: 2, Count: 0},
		"c": {File: "core/mocks/m.go", NumStmt: 10, Count: 0},
	}
	raw := summarize(blocks, func(string) bool { return true })
	if raw.Statements != 16 || raw.Covered != 4 || raw.Missed != 12 {
		t.Fatalf("unexpected raw stats: %+v", raw)
	}
	source := summarize(blocks, func(path string) bool { return !excluded(path) })
	if source.Statements != 6 || source.Covered != 4 || source.Missed != 2 {
		t.Fatalf("unexpected source stats: %+v", source)
	}
}

func TestParseDiffAndPatch(t *testing.T) {
	diff := `diff --git a/core/a.go b/core/a.go
--- a/core/a.go
+++ b/core/a.go
@@ -10,2 +10,3 @@
-old
+new
+newer
 context
@@ -30 +31 @@
-old
+new
diff --git a/core/a_test.go b/core/a_test.go
--- a/core/a_test.go
+++ b/core/a_test.go
@@ -1 +1 @@
-old
+new
`
	files := parseDiff(diff)
	if len(files["core/a.go"]) != 2 {
		t.Fatalf("unexpected ranges: %+v", files)
	}
	if r := files["core/a.go"][0]; r.Start != 10 || r.End != 12 {
		t.Fatalf("unexpected first range: %+v", r)
	}

	blocks := map[string]block{
		"one":     {File: "core/a.go", StartLine: 10, EndLine: 11, NumStmt: 2, Count: 3},
		"two":     {File: "core/a.go", StartLine: 12, EndLine: 12, NumStmt: 1, Count: 0},
		"three":   {File: "core/a.go", StartLine: 31, EndLine: 31, NumStmt: 2, Count: 0},
		"outside": {File: "core/a.go", StartLine: 50, EndLine: 50, NumStmt: 9, Count: 1},
	}
	patch := computePatch(blocks, files, func(path string) bool { return !excluded(path) && !hasTestSuffix(path) })
	if patch.ChangedFiles != 2 || patch.ChangedLines != 5 {
		t.Fatalf("unexpected patch counts: %+v", patch)
	}
	if patch.Statements != 5 || patch.Covered != 2 || patch.Missed != 3 {
		t.Fatalf("unexpected patch stats: %+v", patch)
	}
}

func TestLoadProfilesMergesDuplicates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "coverage.txt")
	content := `mode: atomic
github.com/smartcontractkit/chainlink/v2/core/a.go:1.1,1.10 2 0
github.com/smartcontractkit/chainlink/v2/core/a.go:1.1,1.10 2 5
github.com/smartcontractkit/chainlink/deployment/b.go:2.1,2.10 1 1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	blocks, err := loadProfiles([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 merged blocks, got %d", len(blocks))
	}
	for _, b := range blocks {
		if b.File == "core/a.go" && b.Count != 5 {
			t.Fatalf("duplicate block was not merged with max count: %+v", b)
		}
	}
}

func hasTestSuffix(path string) bool {
	return len(path) >= 8 && path[len(path)-8:] == "_test.go"
}
