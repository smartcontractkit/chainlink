package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateBadgeUsesSourceOnlyAndWritesState(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	reportPath := filepath.Join(dir, "any.json")
	svgPath := filepath.Join(dir, "badge.svg")
	statePath := filepath.Join(dir, "report.json")
	report := `{
		"Global":{"Percent":54.27},
		"GlobalSourceOnly":{"Percent":60.92},
		"Patch":{"Percent":43.84}
	}`
	if err := os.WriteFile(reportPath, []byte(report), 0o600); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.November, 26, 12, 34, 56, 987654321, time.FixedZone("test", 2*60*60))
	state, err := GenerateBadge(BadgeOptions{
		ReportPath: reportPath,
		SVGPath:    svgPath,
		StatePath:  statePath,
		SHA:        "abc123",
		Artifact:   "coverage-report-abc123-1",
		Repo:       "acme/project",
		RunURL:     "https://example.invalid/run/1",
		Now:        now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if state.SHA != "abc123" || state.Repo != "acme/project" || state.ReportArtifact != "coverage-report-abc123-1" || state.RunURL != "https://example.invalid/run/1" {
		t.Errorf("unexpected state metadata: %+v", state)
	}
	if state.GlobalRaw != 54.27 || state.GlobalSourceOnly != 60.92 {
		t.Errorf("unexpected global percentages: %+v", state)
	}
	if state.Patch == nil || *state.Patch != 43.84 {
		t.Errorf("patch = %v, want 43.84", state.Patch)
	}
	if state.UpdatedAt != "2026-11-26T10:34:56Z" {
		t.Errorf("updated_at = %q, want UTC timestamp truncated to seconds", state.UpdatedAt)
	}

	svg, err := os.ReadFile(svgPath)
	if err != nil {
		t.Fatal(err)
	}
	wantSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="160" height="20" role="img" aria-label="coverage: 60.9%"><rect width="160" height="20" rx="3" fill="#555"/><rect x="80" width="80" height="20" rx="3" fill="#f9a825"/><rect x="80" width="4" height="20" fill="#f9a825"/><g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11"><text x="40" y="14">coverage</text><text x="120" y="14">60.9%</text></g></svg>` + "\n"
	if string(svg) != wantSVG {
		t.Errorf("unexpected SVG:\n%s\nwant:\n%s", svg, wantSVG)
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	var written BadgeState
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatal(err)
	}
	if written.GlobalSourceOnly != 60.92 || written.UpdatedAt != state.UpdatedAt || written.Patch == nil || *written.Patch != 43.84 {
		t.Errorf("unexpected written state: %+v", written)
	}
	if !strings.HasSuffix(string(data), "\n") {
		t.Error("state JSON does not end with a newline")
	}
}

func TestGenerateBadgeFallsBackToGlobalAndUsesColorThresholds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		percent float64
		color   string
	}{
		{name: "green", percent: 70, color: "#2e7d32"},
		{name: "yellow", percent: 55, color: "#f9a825"},
		{name: "orange", percent: 45, color: "#ef6c00"},
		{name: "red", percent: 44.99, color: "#c62828"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			reportPath := filepath.Join(dir, "any.json")
			svgPath := filepath.Join(dir, "badge.svg")
			statePath := filepath.Join(dir, "state.json")
			data, err := json.Marshal(map[string]any{"Global": map[string]float64{"Percent": test.percent}})
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(reportPath, data, 0o600); err != nil {
				t.Fatal(err)
			}
			state, err := GenerateBadge(BadgeOptions{ReportPath: reportPath, SVGPath: svgPath, StatePath: statePath, Now: time.Unix(0, 0)})
			if err != nil {
				t.Fatal(err)
			}
			if state.GlobalSourceOnly != test.percent || state.Patch != nil {
				t.Errorf("unexpected state: %+v", state)
			}
			svg, err := os.ReadFile(svgPath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(svg), test.color) {
				t.Errorf("SVG for %.2f%% does not contain %s: %s", test.percent, test.color, svg)
			}
			stateData, err := os.ReadFile(statePath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(stateData), `"patch": null`) {
				t.Errorf("missing null patch in state: %s", stateData)
			}
		})
	}
}

func TestGenerateBadgeValidatesRequiredPaths(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	reportPath := filepath.Join(dir, "any.json")
	if err := os.WriteFile(reportPath, []byte(`{"Global":{"Percent":50}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		opts BadgeOptions
	}{
		{name: "report", opts: BadgeOptions{SVGPath: filepath.Join(dir, "badge.svg"), StatePath: filepath.Join(dir, "state.json")}},
		{name: "SVG", opts: BadgeOptions{ReportPath: reportPath, StatePath: filepath.Join(dir, "state.json")}},
		{name: "state", opts: BadgeOptions{ReportPath: reportPath, SVGPath: filepath.Join(dir, "badge.svg")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := GenerateBadge(test.opts); err == nil {
				t.Fatal("GenerateBadge returned nil error")
			}
		})
	}
}

func TestFullReportGuard(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	full, err := FullReport(dir)
	if err != nil {
		t.Fatal(err)
	}
	if full {
		t.Fatal("empty directory reported as full")
	}
	for _, name := range []string{"unit.txt", "integration.txt", "e2e.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("coverage"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	full, err = FullReport(dir)
	if err != nil {
		t.Fatal(err)
	}
	if full {
		t.Fatal("report without any.json reported as full")
	}
	if err := os.WriteFile(filepath.Join(dir, "any.json"), []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	full, err = FullReport(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !full {
		t.Fatal("complete report not reported as full")
	}
}
