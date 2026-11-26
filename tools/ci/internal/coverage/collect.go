package coverage

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// WorkflowRun is the workflow metadata needed while collecting coverage.
type WorkflowRun struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
}

// CollectClient enumerates and downloads workflow artifacts.
type CollectClient interface {
	RunsForCommit(context.Context, string) ([]WorkflowRun, error)
	ArtifactsForRun(context.Context, int64) ([]Artifact, error)
	DownloadArtifact(context.Context, Artifact, string) error
}

// CollectOptions configures collection of all coverage for one commit.
type CollectOptions struct {
	RepoRoot     string
	TargetSHA    string
	BaseSHA      string
	OutDir       string
	ArtifactsDir string
	Client       CollectClient
	Wait         bool
	WaitTimeout  time.Duration
	PollInterval time.Duration
}

// TestResults contains totals parsed from JUnit test-result artifacts.
type TestResults struct {
	Tests    int `json:"tests"`
	Failures int `json:"failures"`
	Errors   int `json:"errors"`
	Skipped  int `json:"skipped"`
	Files    int `json:"files"`
}

// CollectResult describes the reports and their original workflow artifacts.
type CollectResult struct {
	TargetSHA       string      `json:"target_sha"`
	Full            bool        `json:"full"`
	Categories      []string    `json:"categories"`
	TestResults     TestResults `json:"test_results"`
	SourceArtifacts []Artifact  `json:"source_artifacts"`
	Warnings        []string    `json:"warnings"`
}

type collectedArtifact struct {
	artifact    Artifact
	dir         string
	category    string
	testResults bool
}

type junitSuite struct {
	XMLName  xml.Name
	Tests    int          `xml:"tests,attr"`
	Failures int          `xml:"failures,attr"`
	Errors   int          `xml:"errors,attr"`
	Skipped  int          `xml:"skipped,attr"`
	Suites   []junitSuite `xml:"testsuite"`
}

var watchedCoverageWorkflows = map[string]struct{}{
	"CI Core":                     {},
	"Deployment Tests":            {},
	"Integration In-Memory Tests": {},
	"Integration Tests":           {},
	"CI Tools":                    {},
	"CRE System Tests":            {},
	"CRE Regression System Tests": {},
	"CRE Mixed-Env Tests":         {},
	"CRE Local Environment Tests": {},
	"CCIP System Tests":           {},
	"Legacy System Tests":         {},
	"Legacy Non Functional Tests": {},
}

// Collect downloads and merges coverage without publishing PR comments. Errors
// for individual artifacts and reports are warnings so partial output remains
// available; setup, enumeration, and final output errors are returned.
func Collect(ctx context.Context, opts CollectOptions) (CollectResult, error) {
	result := CollectResult{
		TargetSHA:       opts.TargetSHA,
		Categories:      []string{},
		SourceArtifacts: []Artifact{},
		Warnings:        []string{},
	}
	if strings.TrimSpace(opts.TargetSHA) == "" {
		return result, errors.New("target SHA is required")
	}
	if opts.Client == nil {
		return result, errors.New("collect client is required")
	}
	outDir, artifactsDir, err := prepareCollectionDirs(opts)
	if err != nil {
		return result, err
	}
	opts.OutDir = outDir
	opts.ArtifactsDir = artifactsDir

	var runs []WorkflowRun
	if opts.Wait {
		runs, err = waitForCoverageWorkflows(ctx, opts.Client, opts.TargetSHA, opts.WaitTimeout, opts.PollInterval)
	} else {
		runs, err = opts.Client.RunsForCommit(ctx, opts.TargetSHA)
	}
	if err != nil {
		return result, fmt.Errorf("list workflow runs for %s: %w", opts.TargetSHA, err)
	}
	sort.SliceStable(runs, func(i, j int) bool { return runs[i].ID < runs[j].ID })

	var downloaded []collectedArtifact
	for _, run := range runs {
		artifacts, err := opts.Client.ArtifactsForRun(ctx, run.ID)
		if err != nil {
			if ctx.Err() != nil {
				return result, ctx.Err()
			}
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to list artifacts for run %d: %v", run.ID, err))
			continue
		}
		sort.SliceStable(artifacts, func(i, j int) bool {
			if artifacts[i].Name != artifacts[j].Name {
				return artifacts[i].Name < artifacts[j].Name
			}
			return artifacts[i].ID < artifacts[j].ID
		})
		for _, artifact := range artifacts {
			category := coverageCategory(artifact.Name)
			race := isRaceCoverageArtifact(artifact.Name)
			testResults := strings.HasPrefix(artifact.Name, "test-results-")
			if category == "" && !race && !testResults {
				continue
			}
			result.SourceArtifacts = append(result.SourceArtifacts, artifact)
			if race {
				// Race-instrumented data is linked for separate inspection, but is
				// deliberately never downloaded or merged into normal coverage.
				continue
			}
			destination, err := collectionArtifactDir(artifactsDir, run.ID, artifact.Name)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("skipping artifact %s from run %d: %v", artifact.Name, run.ID, err))
				continue
			}
			if err := os.MkdirAll(destination, 0o755); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("failed to create directory for %s: %v", artifact.Name, err))
				continue
			}
			if err := opts.Client.DownloadArtifact(ctx, artifact, destination); err != nil {
				if ctx.Err() != nil {
					return result, ctx.Err()
				}
				result.Warnings = append(result.Warnings, fmt.Sprintf("failed to download %s from run %d: %v", artifact.Name, run.ID, err))
				continue
			}
			downloaded = append(downloaded, collectedArtifact{
				artifact:    artifact,
				dir:         destination,
				category:    category,
				testResults: testResults,
			})
		}
	}
	sort.SliceStable(result.SourceArtifacts, func(i, j int) bool {
		if result.SourceArtifacts[i].RunID != result.SourceArtifacts[j].RunID {
			return result.SourceArtifacts[i].RunID < result.SourceArtifacts[j].RunID
		}
		if result.SourceArtifacts[i].Name != result.SourceArtifacts[j].Name {
			return result.SourceArtifacts[i].Name < result.SourceArtifacts[j].Name
		}
		return result.SourceArtifacts[i].ID < result.SourceArtifacts[j].ID
	})

	covdata := map[string][]string{"unit": {}, "integration": {}, "e2e": {}}
	var testResultDirs []string
	for _, artifact := range downloaded {
		if artifact.testResults {
			testResultDirs = append(testResultDirs, artifact.dir)
		}
		if artifact.category == "" {
			continue
		}
		dirs, err := collectionCovdataDirs(artifact.dir)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to inspect %s: %v", artifact.artifact.Name, err))
			continue
		}
		covdata[artifact.category] = append(covdata[artifact.category], dirs...)
	}

	var profiles []string
	for _, category := range []string{"unit", "integration", "e2e"} {
		dirs := sortedUniqueCollectionPaths(covdata[category])
		if len(dirs) == 0 {
			continue
		}
		profile := filepath.Join(outDir, category+".txt")
		if err := Merge(ctx, dirs, profile); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to merge %s coverage: %v", category, err))
			continue
		}
		result.Categories = append(result.Categories, category)
		profiles = append(profiles, profile)
		result.Warnings = append(result.Warnings, generateCollectionReport(ctx, opts, category, profile)...)
	}
	if len(profiles) > 0 {
		profile := filepath.Join(outDir, "any.txt")
		if err := MergeProfiles(ctx, profiles, profile); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to merge any coverage: %v", err))
		} else {
			result.Categories = append(result.Categories, "any")
			result.Warnings = append(result.Warnings, generateCollectionReport(ctx, opts, "any", profile)...)
		}
	}

	result.TestResults = collectJUnitResults(testResultDirs, &result.Warnings)
	if err := writeCollectionJSON(filepath.Join(outDir, "test-results.json"), result.TestResults); err != nil {
		return result, err
	}
	result.Full, err = FullReport(outDir)
	if err != nil {
		return result, fmt.Errorf("check full coverage report: %w", err)
	}
	if err := writeCollectionSummary(outDir, &result); err != nil {
		return result, err
	}
	if err := writeCollectionJSON(filepath.Join(outDir, "result.json"), result); err != nil {
		return result, err
	}
	return result, nil
}

func coverageCategory(name string) string {
	if isRaceCoverageArtifact(name) {
		return ""
	}
	switch {
	case strings.HasPrefix(name, "coverage-core-go_core_tests_integration-"),
		strings.HasPrefix(name, "coverage-deployment-"),
		strings.HasPrefix(name, "coverage-inmemory-"):
		return "integration"
	case strings.HasPrefix(name, "coverage-core-go_core_tests-"),
		strings.HasPrefix(name, "coverage-core-scripts-"),
		strings.HasPrefix(name, "coverage-tools-"):
		return "unit"
	case strings.HasPrefix(name, "coverage-e2e-"):
		return "e2e"
	default:
		return ""
	}
}

func isRaceCoverageArtifact(name string) bool {
	return strings.HasPrefix(name, "coverage-core-go_core_race_tests-")
}

func waitForCoverageWorkflows(ctx context.Context, client CollectClient, sha string, timeout, pollInterval time.Duration) ([]WorkflowRun, error) {
	if timeout <= 0 {
		timeout = 45 * time.Minute
	}
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		if err := waitCtx.Err(); err != nil {
			return nil, err
		}
		runs, err := client.RunsForCommit(waitCtx, sha)
		if err != nil {
			return nil, err
		}
		pending := false
		for _, run := range runs {
			if _, watched := watchedCoverageWorkflows[run.Name]; watched && (run.Status == "queued" || run.Status == "in_progress") {
				pending = true
				break
			}
		}
		if !pending {
			return runs, nil
		}
		if pollInterval <= 0 {
			continue
		}
		timer := time.NewTimer(pollInterval)
		select {
		case <-waitCtx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return nil, waitCtx.Err()
		case <-timer.C:
		}
	}
}

func prepareCollectionDirs(opts CollectOptions) (string, string, error) {
	if opts.OutDir == "" || opts.ArtifactsDir == "" {
		return "", "", errors.New("output and artifacts directories are required")
	}
	outDir, err := filepath.Abs(opts.OutDir)
	if err != nil {
		return "", "", err
	}
	artifactsDir, err := filepath.Abs(opts.ArtifactsDir)
	if err != nil {
		return "", "", err
	}
	outDir = filepath.Clean(outDir)
	artifactsDir = filepath.Clean(artifactsDir)
	if outDir == artifactsDir {
		return "", "", errors.New("output and artifacts directories must be different")
	}
	var protectedDirs []string
	if opts.RepoRoot != "" {
		repoRoot, err := filepath.Abs(opts.RepoRoot)
		if err != nil {
			return "", "", err
		}
		protectedDirs = append(protectedDirs, filepath.Clean(repoRoot))
	}
	if workingDir, err := os.Getwd(); err == nil {
		protectedDirs = append(protectedDirs, filepath.Clean(workingDir))
	}
	for _, dir := range []string{outDir, artifactsDir} {
		if dir == string(filepath.Separator) {
			return "", "", errors.New("unsafe collection directory")
		}
		for _, protected := range protectedDirs {
			if collectionDirContains(dir, protected) {
				return "", "", fmt.Errorf("refusing to clean directory %q containing %q", dir, protected)
			}
		}
	}
	if err := os.RemoveAll(artifactsDir); err != nil {
		return "", "", err
	}
	if err := os.RemoveAll(outDir); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(artifactsDir, 0o755); err != nil {
		return "", "", err
	}
	return outDir, artifactsDir, nil
}

func collectionDirContains(dir, protected string) bool {
	relative, err := filepath.Rel(dir, protected)
	return err == nil && !filepath.IsAbs(relative) && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func collectionArtifactDir(root string, runID int64, name string) (string, error) {
	if runID <= 0 {
		return "", errors.New("workflow run ID must be positive")
	}
	if name == "" || name == "." || filepath.Base(name) != name || strings.ContainsAny(name, `/\\`) {
		return "", errors.New("unsafe artifact name")
	}
	destination := filepath.Join(root, strconv.FormatInt(runID, 10), name)
	relative, err := filepath.Rel(root, destination)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("artifact destination escapes artifacts directory")
	}
	return destination, nil
}

func collectionCovdataDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "covmeta.") {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return sortedUniqueCollectionPaths(dirs), err
}

func generateCollectionReport(ctx context.Context, opts CollectOptions, category, profile string) []string {
	reportOpts := ReportOptions{
		Profiles:     []string{profile},
		Repo:         opts.RepoRoot,
		Category:     category,
		JSONPath:     filepath.Join(opts.OutDir, category+".json"),
		MarkdownPath: filepath.Join(opts.OutDir, category+".md"),
		MinGlobal:    -1,
		MinPatch:     -1,
	}
	if opts.BaseSHA != "" {
		reportOpts.Base = opts.BaseSHA
		reportOpts.Head = opts.TargetSHA
	}
	if _, err := GenerateReport(ctx, reportOpts); err != nil {
		if reportOpts.Base == "" {
			return []string{fmt.Sprintf("failed to report %s coverage: %v", category, err)}
		}
		warnings := []string{fmt.Sprintf("failed to compute %s patch coverage; retrying without patch: %v", category, err)}
		reportOpts.Base = ""
		reportOpts.Head = ""
		if _, fallbackErr := GenerateReport(ctx, reportOpts); fallbackErr != nil {
			warnings = append(warnings, fmt.Sprintf("failed to report %s coverage: %v", category, fallbackErr))
		}
		return warnings
	}
	return nil
}

func collectJUnitResults(dirs []string, warnings *[]string) TestResults {
	var paths []string
	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".xml") {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			*warnings = append(*warnings, fmt.Sprintf("failed to inspect test results in %s: %v", dir, err))
		}
	}
	paths = sortedUniqueCollectionPaths(paths)
	totals := TestResults{Files: len(paths)}
	for _, path := range paths {
		root, err := parseCollectionJUnit(path)
		if err != nil {
			*warnings = append(*warnings, fmt.Sprintf("failed to parse test results %s: %v", path, err))
			continue
		}
		suites := []junitSuite{root}
		if root.XMLName.Local == "testsuites" {
			suites = root.Suites
		}
		for _, suite := range suites {
			totals.Tests += suite.Tests
			totals.Failures += suite.Failures
			totals.Errors += suite.Errors
			totals.Skipped += suite.Skipped
		}
	}
	return totals
}

func parseCollectionJUnit(path string) (junitSuite, error) {
	var root junitSuite
	file, err := os.Open(path)
	if err != nil {
		return root, err
	}
	defer file.Close()
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return root, err
	}
	if root.XMLName.Local != "testsuite" && root.XMLName.Local != "testsuites" {
		return root, fmt.Errorf("unexpected JUnit root element %q", root.XMLName.Local)
	}
	return root, nil
}

func writeCollectionSummary(outDir string, result *CollectResult) error {
	var summary strings.Builder
	summary.WriteString("<!-- coverage-report -->\n")
	fmt.Fprintf(&summary, "# Coverage report for `%s`\n\n", shortCollectionSHA(result.TargetSHA))
	fmt.Fprintf(&summary, "Commit: `%s`\n\n", result.TargetSHA)
	fmt.Fprintf(&summary, "Full report: `%t`\n\n", result.Full)
	summary.WriteString("### Test results\n\n```json\n")
	data, err := json.MarshalIndent(result.TestResults, "", "  ")
	if err != nil {
		return err
	}
	summary.Write(data)
	summary.WriteString("\n```\n\n")
	for _, category := range []string{"any", "unit", "integration", "e2e"} {
		data, err := os.ReadFile(filepath.Join(outDir, category+".md"))
		if err != nil {
			if !os.IsNotExist(err) {
				result.Warnings = append(result.Warnings, fmt.Sprintf("failed to read %s report: %v", category, err))
			}
			continue
		}
		summary.Write(data)
		if !strings.HasSuffix(string(data), "\n") {
			summary.WriteByte('\n')
		}
		summary.WriteByte('\n')
	}
	summary.WriteString("### Source artifacts\n\n")
	if len(result.SourceArtifacts) == 0 {
		summary.WriteString("No coverage or test-result artifacts were found.\n")
	}
	for _, artifact := range result.SourceArtifacts {
		if artifact.URL != "" {
			fmt.Fprintf(&summary, "- [%s](%s)", artifact.Name, artifact.URL)
		} else {
			fmt.Fprintf(&summary, "- `%s`", artifact.Name)
		}
		if isRaceCoverageArtifact(artifact.Name) {
			summary.WriteString(" — race coverage (separate; never merged)")
		}
		if artifact.RunURL != "" {
			fmt.Fprintf(&summary, " — %s ([run](%s))", artifact.RunName, artifact.RunURL)
		}
		summary.WriteByte('\n')
	}
	if len(result.Warnings) > 0 {
		summary.WriteString("\n### Warnings\n\n")
		for _, warning := range result.Warnings {
			fmt.Fprintf(&summary, "- %s\n", warning)
		}
	}
	return os.WriteFile(filepath.Join(outDir, "summary.md"), []byte(summary.String()), 0o600)
}

func writeCollectionJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func sortedUniqueCollectionPaths(paths []string) []string {
	set := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		set[path] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for path := range set {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func shortCollectionSHA(sha string) string {
	if len(sha) > 9 {
		return sha[:9]
	}
	return sha
}
