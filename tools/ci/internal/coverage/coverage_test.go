package coverage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBlockAndNormalize(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	patch := computePatch(blocks, files, func(path string) bool { return !excluded(path) && !strings.HasSuffix(path, "_test.go") })
	if patch.ChangedFiles != 2 || patch.ChangedLines != 5 {
		t.Fatalf("unexpected patch counts: %+v", patch)
	}
	if patch.Statements != 5 || patch.Covered != 2 || patch.Missed != 3 {
		t.Fatalf("unexpected patch stats: %+v", patch)
	}
}

func TestLoadProfilesMergesDuplicates(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "coverage.txt")
	content := `mode: atomic
github.com/smartcontractkit/chainlink/v2/core/a.go:1.1,1.10 2 0
github.com/smartcontractkit/chainlink/v2/core/a.go:1.1,1.10 2 5
github.com/smartcontractkit/chainlink/deployment/b.go:2.1,2.10 1 1
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
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

func TestGenerateReportWritesJSONAndMarkdown(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	profile := filepath.Join(dir, "coverage.txt")
	jsonPath := filepath.Join(dir, "report.json")
	markdownPath := filepath.Join(dir, "report.md")
	content := `mode: atomic
github.com/smartcontractkit/chainlink/v2/core/a.go:1.1,1.10 2 1
github.com/smartcontractkit/chainlink/v2/core/a.go:2.1,2.10 2 0
github.com/smartcontractkit/chainlink/v2/core/mocks/m.go:1.1,1.10 4 0
`
	if err := os.WriteFile(profile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := GenerateReport(context.Background(), ReportOptions{
		Profiles:     []string{profile},
		Repo:         ".",
		Category:     "unit",
		JSONPath:     jsonPath,
		MarkdownPath: markdownPath,
		MinGlobal:    -1,
		MinPatch:     -1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Report.Global.Statements != 8 || result.Report.Global.Covered != 2 || result.Report.Global.Percent != 25 {
		t.Fatalf("unexpected global stats: %+v", result.Report.Global)
	}
	if result.Report.GlobalSourceOnly.Statements != 4 || result.Report.GlobalSourceOnly.Covered != 2 {
		t.Fatalf("unexpected source-only stats: %+v", result.Report.GlobalSourceOnly)
	}
	if result.Markdown != RenderMarkdown(result.Report) {
		t.Fatal("returned markdown differs from RenderMarkdown")
	}
	markdown, err := os.ReadFile(markdownPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(markdown) != result.Markdown {
		t.Fatalf("unexpected markdown file:\n%s", markdown)
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Report
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Category != "unit" || decoded.GlobalSourceOnly.Missed != 2 {
		t.Fatalf("unexpected JSON report: %+v", decoded)
	}
	for _, field := range []string{`"Category"`, `"Profile"`, `"GlobalSourceOnly"`, `"PatchSourceOnly"`} {
		if !strings.Contains(string(data), field) {
			t.Errorf("JSON does not contain legacy field %s: %s", field, data)
		}
	}
}

func TestGenerateReportThresholds(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	profile := filepath.Join(dir, "coverage.txt")
	if err := os.WriteFile(profile, []byte("mode: atomic\ncore/a.go:1.1,1.10 1 0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := GenerateReport(context.Background(), ReportOptions{
		Profiles:  []string{profile},
		MinGlobal: 50,
		MinPatch:  50,
	})
	if err == nil || !strings.Contains(err.Error(), "global coverage 0.00% is below 50.00%") || !strings.Contains(err.Error(), "patch coverage requested but --base/--head not provided") {
		t.Fatalf("unexpected threshold error: %v", err)
	}
}

func TestMergeProfiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	one := filepath.Join(dir, "one.txt")
	two := filepath.Join(dir, "two.txt")
	out := filepath.Join(dir, "merged.txt")
	if err := os.WriteFile(one, []byte("mode: set\ncore/b.go:2.1,2.2 1 0\ncore/a.go:1.1,1.2 2 3\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(two, []byte("mode: atomic\ncore/a.go:1.1,1.2 2 8\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MergeProfiles(context.Background(), []string{one, two}, out); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	want := "mode: atomic\ncore/a.go:1.1,1.2 2 8\ncore/b.go:2.1,2.2 1 0\n"
	if string(data) != want {
		t.Fatalf("unexpected merged profile:\n%s\nwant:\n%s", data, want)
	}
}

func TestRenderMarkdown(t *testing.T) {
	t.Parallel()
	rep := Report{
		Category:         "unit",
		Base:             "base",
		Head:             "head",
		Global:           Stats{Statements: 4, Covered: 3, Missed: 1, Percent: 75},
		GlobalSourceOnly: Stats{Statements: 2, Covered: 2, Percent: 100},
		Patch: &PatchStats{
			ChangedFiles: 1,
			ChangedLines: 2,
			Statements:   2,
			Covered:      1,
			Missed:       1,
			Percent:      50,
			Files: []FilePatch{{
				Path:        "core/a.go",
				Added:       []LineRange{{Start: 1, End: 2}},
				Statements:  2,
				Covered:     1,
				Missed:      1,
				MissedSpans: []string{"2-2"},
			}},
		},
		PatchSourceOnly: &PatchStats{ChangedFiles: 1, ChangedLines: 2, Statements: 2, Covered: 1, Missed: 1, Percent: 50},
	}
	want := `## Coverage — unit

` + "`base..head`" + `

### Global

| Scope | Statements | Covered | Missed | Coverage |
|---|---:|---:|---:|---:|
| raw | 4 | 3 | 1 | **75.00%** |
| source-only | 2 | 2 | 0 | **100.00%** |

### Patch

| Scope | Changed files | Changed lines | Statements | Covered | Missed | Coverage |
|---|---:|---:|---:|---:|---:|---:|
| raw | 1 | 2 | 2 | 1 | 1 | **50.00%** |
| source-only | 1 | 2 | 2 | 1 | 1 | **50.00%** |

### Changed files

| File | Statements | Covered | Missed | Coverage | Missed spans |
|---|---:|---:|---:|---:|---|
| ` + "`core/a.go`" + ` | 2 | 1 | 1 | 50.00% | 2-2 |
`
	if got := RenderMarkdown(rep); got != want {
		t.Fatalf("unexpected markdown:\n%s\nwant:\n%s", got, want)
	}
}
