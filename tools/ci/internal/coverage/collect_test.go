package coverage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type fakeCollectClient struct {
	runBatches   [][]WorkflowRun
	runs         []WorkflowRun
	runCalls     int
	artifacts    map[int64][]Artifact
	artifactErrs map[int64]error
	downloads    []string
	download     func(context.Context, Artifact, string) error
	onRuns       func()
}

func (f *fakeCollectClient) RunsForCommit(_ context.Context, _ string) ([]WorkflowRun, error) {
	f.runCalls++
	if f.onRuns != nil {
		f.onRuns()
	}
	if len(f.runBatches) > 0 {
		index := f.runCalls - 1
		if index >= len(f.runBatches) {
			index = len(f.runBatches) - 1
		}
		return f.runBatches[index], nil
	}
	return f.runs, nil
}

func (f *fakeCollectClient) ArtifactsForRun(_ context.Context, runID int64) ([]Artifact, error) {
	if err := f.artifactErrs[runID]; err != nil {
		return nil, err
	}
	return f.artifacts[runID], nil
}

func (f *fakeCollectClient) DownloadArtifact(ctx context.Context, artifact Artifact, destination string) error {
	f.downloads = append(f.downloads, artifact.Name)
	if f.download != nil {
		return f.download(ctx, artifact, destination)
	}
	return nil
}

func TestCoverageCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		artifact string
		want     string
	}{
		{name: "core unit", artifact: "coverage-core-go_core_tests-sha-1", want: "unit"},
		{name: "scripts unit", artifact: "coverage-core-scripts-sha-1", want: "unit"},
		{name: "tools unit", artifact: "coverage-tools-sha-1", want: "unit"},
		{name: "core integration takes precedence", artifact: "coverage-core-go_core_tests_integration-sha-1", want: "integration"},
		{name: "deployment integration", artifact: "coverage-deployment-sha-1", want: "integration"},
		{name: "in memory integration", artifact: "coverage-inmemory-sha-1", want: "integration"},
		{name: "e2e", artifact: "coverage-e2e-system-1", want: "e2e"},
		{name: "race is never normal coverage", artifact: "coverage-core-go_core_race_tests-sha-1", want: ""},
		{name: "unrelated", artifact: "coverage-report-sha-1", want: ""},
		{name: "test results have no coverage category", artifact: "test-results-core-1", want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := coverageCategory(test.artifact); got != test.want {
				t.Fatalf("coverageCategory(%q) = %q, want %q", test.artifact, got, test.want)
			}
		})
	}
	if !isRaceCoverageArtifact("coverage-core-go_core_race_tests-sha-1") {
		t.Fatal("race artifact was not identified for source linking")
	}
}

func TestCollectAggregatesJUnitAndLinksAllSources(t *testing.T) {
	t.Parallel()

	covdata := makeTestCovdata(t)
	unit := sourceArtifact(1, "coverage-core-go_core_tests-sha-1")
	race := sourceArtifact(2, "coverage-core-go_core_race_tests-sha-1")
	results := sourceArtifact(3, "test-results-core-1")
	broken := sourceArtifact(4, "coverage-e2e-broken-1")
	unrelated := sourceArtifact(5, "unrelated-1")
	client := &fakeCollectClient{
		runs:      []WorkflowRun{{ID: 42, Name: "CI Core", Status: "completed", Conclusion: "success"}},
		artifacts: map[int64][]Artifact{42: {unit, race, results, broken, unrelated}},
	}
	client.download = func(_ context.Context, artifact Artifact, destination string) error {
		switch artifact.ID {
		case unit.ID:
			copyCovdata(t, covdata, destination)
		case results.ID:
			writeTestFile(t, filepath.Join(destination, "one.xml"), `<testsuites><testsuite tests="3" failures="1" errors="0" skipped="1"/><testsuite tests="2" failures="0" errors="1" skipped="0"/></testsuites>`)
			writeTestFile(t, filepath.Join(destination, "nested", "two.xml"), `<testsuite tests="4" failures="1" errors="1" skipped="2"/>`)
		case broken.ID:
			return errors.New("download failed")
		case race.ID:
			t.Error("race coverage must not be downloaded")
		}
		return nil
	}
	outDir := filepath.Join(t.TempDir(), "out")
	artifactsDir := filepath.Join(t.TempDir(), "artifacts")
	result, err := Collect(context.Background(), CollectOptions{
		RepoRoot:     t.TempDir(),
		TargetSHA:    "0123456789abcdef",
		OutDir:       outDir,
		ArtifactsDir: artifactsDir,
		Client:       client,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Full {
		t.Fatal("unit-only report was marked full")
	}
	if want := []string{"unit", "any"}; !reflect.DeepEqual(result.Categories, want) {
		t.Fatalf("Categories = %v, want %v", result.Categories, want)
	}
	if want := (TestResults{Tests: 9, Failures: 2, Errors: 2, Skipped: 3, Files: 2}); result.TestResults != want {
		t.Fatalf("TestResults = %+v, want %+v", result.TestResults, want)
	}
	if len(result.SourceArtifacts) != 4 {
		t.Fatalf("SourceArtifacts = %+v, want unit, race, test results, and failed e2e download", result.SourceArtifacts)
	}
	if !containsArtifact(result.SourceArtifacts, race.ID) {
		t.Fatal("race artifact is missing from SourceArtifacts")
	}
	if containsArtifact(result.SourceArtifacts, unrelated.ID) {
		t.Fatal("unrelated artifact was included in SourceArtifacts")
	}
	if want := []string{unit.Name, broken.Name, results.Name}; !reflect.DeepEqual(client.downloads, want) {
		t.Fatalf("downloads = %v, want %v; race must be excluded", client.downloads, want)
	}
	if len(result.Warnings) == 0 || !strings.Contains(strings.Join(result.Warnings, "\n"), broken.Name) {
		t.Fatalf("download failure did not become a warning: %v", result.Warnings)
	}

	summary := readTestFile(t, filepath.Join(outDir, "summary.md"))
	for _, want := range []string{
		"0123456789abcdef",
		"## Coverage — unit",
		"[" + race.Name + "](" + race.URL + ")",
		"[run](" + race.RunURL + ")",
		`"tests": 9`,
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary.md does not contain %q:\n%s", want, summary)
		}
	}
	var decoded CollectResult
	data := readTestFile(t, filepath.Join(outDir, "result.json"))
	if err := json.Unmarshal([]byte(data), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.TargetSHA != result.TargetSHA || decoded.TestResults != result.TestResults || decoded.Full {
		t.Fatalf("unexpected result.json: %+v", decoded)
	}
	for _, field := range []string{`"target_sha"`, `"test_results"`, `"source_artifacts"`, `"warnings"`} {
		if !strings.Contains(data, field) {
			t.Errorf("result.json does not contain %s: %s", field, data)
		}
	}
}

func TestCollectFullGuard(t *testing.T) {
	t.Parallel()

	covdata := makeTestCovdata(t)
	artifacts := []Artifact{
		sourceArtifact(1, "coverage-core-go_core_tests-sha-1"),
		sourceArtifact(2, "coverage-core-go_core_tests_integration-sha-1"),
		sourceArtifact(3, "coverage-e2e-system-1"),
	}
	client := &fakeCollectClient{
		runs:      []WorkflowRun{{ID: 42, Name: "CI Core", Status: "completed"}},
		artifacts: map[int64][]Artifact{42: artifacts},
		download: func(_ context.Context, _ Artifact, destination string) error {
			copyCovdata(t, covdata, destination)
			return nil
		},
	}
	outDir := filepath.Join(t.TempDir(), "out")
	result, err := Collect(context.Background(), CollectOptions{
		RepoRoot:     t.TempDir(),
		TargetSHA:    "0123456789abcdef",
		OutDir:       outDir,
		ArtifactsDir: filepath.Join(t.TempDir(), "artifacts"),
		Client:       client,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Full {
		t.Fatalf("complete report was not marked full; warnings: %v", result.Warnings)
	}
	if want := []string{"unit", "integration", "e2e", "any"}; !reflect.DeepEqual(result.Categories, want) {
		t.Fatalf("Categories = %v, want %v", result.Categories, want)
	}
	for _, name := range []string{"unit.txt", "integration.txt", "e2e.txt", "any.json", "any.md"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
}

func TestCollectAllowsPartialOutputAfterMergeFailure(t *testing.T) {
	t.Parallel()

	unit := sourceArtifact(1, "coverage-core-go_core_tests-sha-1")
	client := &fakeCollectClient{
		runs:      []WorkflowRun{{ID: 42, Name: "CI Core", Status: "completed"}},
		artifacts: map[int64][]Artifact{42: {unit}},
		download: func(_ context.Context, _ Artifact, destination string) error {
			writeTestFile(t, filepath.Join(destination, "covmeta.invalid"), "not binary coverage")
			return nil
		},
	}
	outDir := filepath.Join(t.TempDir(), "out")
	result, err := Collect(context.Background(), CollectOptions{
		TargetSHA:    "0123456789abcdef",
		OutDir:       outDir,
		ArtifactsDir: filepath.Join(t.TempDir(), "artifacts"),
		Client:       client,
	})
	if err != nil {
		t.Fatalf("merge failure should allow partial output: %v", err)
	}
	if result.Full || len(result.Categories) != 0 {
		t.Fatalf("failed merge produced a full/category result: %+v", result)
	}
	if len(result.Warnings) == 0 || !strings.Contains(strings.Join(result.Warnings, "\n"), "unit") {
		t.Fatalf("merge failure did not become a warning: %v", result.Warnings)
	}
	if _, err := os.Stat(filepath.Join(outDir, "result.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "summary.md")); err != nil {
		t.Fatal(err)
	}
}

func TestCollectRefusesToCleanRepositoryAncestor(t *testing.T) {
	t.Parallel()

	parent := t.TempDir()
	repoRoot := filepath.Join(parent, "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	client := &fakeCollectClient{artifacts: map[int64][]Artifact{}}
	_, err := Collect(context.Background(), CollectOptions{
		RepoRoot:     repoRoot,
		TargetSHA:    "0123456789abcdef",
		OutDir:       parent,
		ArtifactsDir: filepath.Join(t.TempDir(), "artifacts"),
		Client:       client,
	})
	if err == nil {
		t.Fatal("Collect accepted an output directory containing the repository")
	}
	if client.runCalls != 0 {
		t.Fatal("Collect enumerated runs before rejecting an unsafe output directory")
	}
}

func TestCollectWaitsOnlyForWatchedWorkflows(t *testing.T) {
	t.Parallel()

	client := &fakeCollectClient{
		runBatches: [][]WorkflowRun{
			{
				{ID: 1, Name: "CI Core", Status: "in_progress"},
				{ID: 2, Name: "Unwatched Workflow", Status: "queued"},
			},
			{
				{ID: 1, Name: "CI Core", Status: "completed"},
				{ID: 2, Name: "Unwatched Workflow", Status: "queued"},
			},
		},
		artifacts: map[int64][]Artifact{},
	}
	_, err := Collect(context.Background(), CollectOptions{
		TargetSHA:    "0123456789abcdef",
		OutDir:       filepath.Join(t.TempDir(), "out"),
		ArtifactsDir: filepath.Join(t.TempDir(), "artifacts"),
		Client:       client,
		Wait:         true,
		WaitTimeout:  time.Minute,
		PollInterval: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.runCalls != 2 {
		t.Fatalf("RunsForCommit calls = %d, want 2; unwatched queued runs must not block", client.runCalls)
	}
}

func TestCollectWaitHonorsContextWithoutSleeping(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	client := &fakeCollectClient{
		runs:      []WorkflowRun{{ID: 1, Name: "Deployment Tests", Status: "queued"}},
		artifacts: map[int64][]Artifact{},
		onRuns:    cancel,
	}
	_, err := Collect(ctx, CollectOptions{
		TargetSHA:    "0123456789abcdef",
		OutDir:       filepath.Join(t.TempDir(), "out"),
		ArtifactsDir: filepath.Join(t.TempDir(), "artifacts"),
		Client:       client,
		Wait:         true,
		WaitTimeout:  time.Hour,
		PollInterval: time.Hour,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Collect error = %v, want context.Canceled", err)
	}
	if client.runCalls != 1 {
		t.Fatalf("RunsForCommit calls = %d, want 1", client.runCalls)
	}
}

func sourceArtifact(id int64, name string) Artifact {
	return Artifact{
		ID:      id,
		Name:    name,
		RunID:   42,
		RunName: "CI Core",
		SHA:     "0123456789abcdef",
		URL:     "https://github.example/runs/42/artifacts/" + strconv.FormatInt(id, 10),
		RunURL:  "https://github.example/runs/42",
	}
}

func containsArtifact(artifacts []Artifact, id int64) bool {
	for _, artifact := range artifacts {
		if artifact.ID == id {
			return true
		}
	}
	return false
}

func makeTestCovdata(t *testing.T) string {
	t.Helper()
	moduleDir := t.TempDir()
	writeTestFile(t, filepath.Join(moduleDir, "go.mod"), "module coveragefixture\n\ngo 1.23\n")
	writeTestFile(t, filepath.Join(moduleDir, "fixture.go"), "package coveragefixture\n\nfunc Value() int { return 1 }\n")
	writeTestFile(t, filepath.Join(moduleDir, "fixture_test.go"), "package coveragefixture\n\nimport \"testing\"\n\nfunc TestValue(t *testing.T) { if Value() != 1 { t.Fatal() } }\n")
	covdata := t.TempDir()
	cmd := exec.Command("go", "test", "-cover", "-args", "-test.gocoverdir="+covdata)
	cmd.Dir = moduleDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create coverage fixture: %v\n%s", err, output)
	}
	return covdata
}

func copyCovdata(t *testing.T, source, destination string) {
	t.Helper()
	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		writeTestFile(t, filepath.Join(destination, entry.Name()), string(data))
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
