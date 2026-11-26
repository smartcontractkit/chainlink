package coverage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PROptions configures the lightweight affected-package coverage run.
type PROptions struct {
	RepoRoot    string
	BaseSHA     string
	HeadSHA     string
	PRNumber    int
	OutDir      string
	Short       bool
	TestTimeout time.Duration
}

// PRResult is the untrusted data contract passed to the trusted PR publisher.
// The publisher renders Markdown from this data instead of executing or
// publishing arbitrary files produced by a fork.
type PRResult struct {
	BaseSHA     string   `json:"base_sha"`
	HeadSHA     string   `json:"head_sha"`
	PRNumber    int      `json:"pr_number,omitempty"`
	Short       bool     `json:"short"`
	NoGoChanges bool     `json:"no_go_changes"`
	NoRunnable  bool     `json:"no_runnable,omitempty"`
	TestsFailed bool     `json:"tests_failed"`
	Packages    []string `json:"packages,omitempty"`
	NotMeasured []string `json:"not_measured,omitempty"`
	Report      *Report  `json:"report,omitempty"`
}

// RunPR tests affected packages and writes a machine-readable result. It does
// not use a GitHub token or publish comments, so it is safe for fork PRs.
func RunPR(ctx context.Context, opts PROptions) (result PRResult, retErr error) {
	result = PRResult{BaseSHA: opts.BaseSHA, HeadSHA: opts.HeadSHA, PRNumber: opts.PRNumber, Short: opts.Short}
	writeResult := false
	defer func() {
		if writeResult {
			if err := writePRResult(opts.OutDir, result); retErr == nil && err != nil {
				retErr = err
			}
		}
	}()
	if opts.RepoRoot == "" || opts.BaseSHA == "" || opts.HeadSHA == "" {
		return result, errors.New("repo root, base SHA, and head SHA are required")
	}
	if opts.OutDir == "" {
		opts.OutDir = filepath.Join(opts.RepoRoot, "coverage-patch")
	} else if !filepath.IsAbs(opts.OutDir) {
		opts.OutDir = filepath.Join(opts.RepoRoot, opts.OutDir)
	}
	clean, err := filepath.Abs(opts.OutDir)
	if err != nil || clean == opts.RepoRoot || clean == string(filepath.Separator) {
		return result, errors.New("unsafe coverage output directory")
	}
	opts.OutDir = clean
	if opts.TestTimeout <= 0 {
		opts.TestTimeout = 10 * time.Minute
	}
	if err := os.RemoveAll(opts.OutDir); err != nil {
		return result, err
	}
	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return result, err
	}
	writeResult = true

	changedOut, err := commandOutput(ctx, opts.RepoRoot, "git", "diff", "--name-only", "--diff-filter=ACMR", opts.BaseSHA+"..."+opts.HeadSHA, "--", "*.go")
	if err != nil {
		return result, fmt.Errorf("list changed Go files: %w", err)
	}
	var changed []string
	for _, path := range strings.Split(changedOut, "\n") {
		path = filepath.ToSlash(strings.TrimSpace(path))
		if productionGoFile(path) {
			changed = append(changed, path)
		}
	}
	sort.Strings(changed)
	if err := os.WriteFile(filepath.Join(opts.OutDir, "changed-files.txt"), []byte(strings.Join(changed, "\n")+"\n"), 0o600); err != nil {
		return result, err
	}
	if len(changed) == 0 {
		result.NoGoChanges = true
		return result, writePRResult(opts.OutDir, result)
	}

	modules := make(map[string][]string)
	for _, path := range changed {
		module := nearestModule(opts.RepoRoot, filepath.Dir(path))
		dir := filepath.Join(opts.RepoRoot, filepath.Dir(path))
		modules[module] = append(modules[module], dir)
	}
	moduleNames := make([]string, 0, len(modules))
	for module := range modules {
		moduleNames = append(moduleNames, module)
	}
	sort.Strings(moduleNames)

	logPath := filepath.Join(opts.OutDir, "test.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		return result, err
	}
	defer logFile.Close()

	var covRoots []string
	for _, module := range moduleNames {
		packages, err := affectedPackages(ctx, module, modules[module])
		if err != nil {
			fmt.Fprintf(logFile, "go list failed in %s: %v\n", module, err)
			result.TestsFailed = true
			continue
		}
		if len(packages) == 0 {
			continue
		}
		result.Packages = append(result.Packages, packages...)
		rel, err := filepath.Rel(opts.RepoRoot, module)
		if err != nil || rel == "." {
			rel = "root"
		}
		covDir := filepath.Join(opts.OutDir, "covdata", strings.NewReplacer("/", "_", "\\", "_").Replace(rel))
		if err := os.MkdirAll(covDir, 0o755); err != nil {
			return result, err
		}
		args := []string{"test", "-count=1", "-timeout", opts.TestTimeout.String()}
		if opts.Short {
			args = append(args, "-short")
		}
		args = append(args, "-cover", "-covermode=atomic", "-coverpkg="+strings.Join(packages, ","))
		args = append(args, packages...)
		args = append(args, "-args", "-test.gocoverdir="+covDir)
		cmd := exec.CommandContext(ctx, "go", args...) //nolint:gosec // arguments are passed without a shell
		cmd.Dir = module
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		if err := cmd.Run(); err != nil {
			result.TestsFailed = true
			fmt.Fprintf(logFile, "go test failed in %s: %v\n", module, err)
		}
		if hasCovdata(covDir) {
			covRoots = append(covRoots, covDir)
		}
	}
	sort.Strings(result.Packages)
	result.Packages = uniqueStrings(result.Packages)
	if len(result.Packages) == 0 {
		result.NoRunnable = true
		return result, writePRResult(opts.OutDir, result)
	}

	if len(covRoots) > 0 {
		profile := filepath.Join(opts.OutDir, "patch.txt")
		if err := Merge(ctx, covRoots, profile); err != nil {
			return result, fmt.Errorf("merge PR coverage: %w", err)
		}
		mergeBase, err := commandOutput(ctx, opts.RepoRoot, "git", "merge-base", opts.BaseSHA, opts.HeadSHA)
		if err != nil {
			return result, fmt.Errorf("find merge base: %w", err)
		}
		mergeBase = strings.TrimSpace(mergeBase)
		reported, err := GenerateReport(ctx, ReportOptions{
			Profiles:     []string{profile},
			Base:         mergeBase,
			Head:         opts.HeadSHA,
			Repo:         opts.RepoRoot,
			Category:     "patch (affected packages)",
			JSONPath:     filepath.Join(opts.OutDir, "patch.json"),
			MarkdownPath: filepath.Join(opts.OutDir, "patch.md"),
			MinGlobal:    -1,
			MinPatch:     -1,
		})
		if err != nil {
			return result, err
		}
		result.Report = &reported.Report
		if reported.Report.Patch != nil {
			for _, file := range reported.Report.Patch.Files {
				if !file.TestFile && file.Statements == 0 {
					result.NotMeasured = append(result.NotMeasured, file.Path)
				}
			}
		}
	}
	if err := writePRResult(opts.OutDir, result); err != nil {
		return result, err
	}
	if result.TestsFailed {
		return result, errors.New("some affected-package tests failed; coverage may be incomplete")
	}
	return result, nil
}

// LoadPRResult reads and validates the small structured artifact consumed by
// the trusted publisher. The size limit prevents accidental artifact abuse.
func LoadPRResult(path string) (PRResult, error) {
	var result PRResult
	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()
	decoder := json.NewDecoder(io.LimitReader(f, 32<<20))
	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("decode PR coverage result: %w", err)
	}
	if result.HeadSHA == "" {
		return result, errors.New("PR coverage result has no head SHA")
	}
	return result, nil
}

// PublishPROptions configures trusted PR comment publication.
type PublishPROptions struct {
	PRNumber      int
	ResultPath    string
	RunID         int64
	RunURL        string
	PRArtifactURL string
	NightlyBranch string
}

// PublishPR discovers the latest nightly baseline and updates the sticky PR
// comment using trusted rendering code.
func PublishPR(ctx context.Context, client *GitHubClient, opts PublishPROptions) error {
	result, err := LoadPRResult(opts.ResultPath)
	if err != nil {
		return err
	}
	if opts.PRNumber == 0 {
		opts.PRNumber = result.PRNumber
	}
	if opts.PRNumber <= 0 {
		return errors.New("a valid PR number is required")
	}
	if opts.NightlyBranch == "" {
		opts.NightlyBranch = "develop"
	}
	if opts.PRArtifactURL == "" && opts.RunID > 0 {
		artifacts, err := client.ArtifactsForRun(ctx, opts.RunID)
		if err == nil {
			prefix := "coverage-patch-" + result.HeadSHA
			for _, artifact := range artifacts {
				if artifact.Name == prefix {
					opts.PRArtifactURL = artifact.URL
					break
				}
			}
		}
	}
	nightly, nightlyErr := client.LatestNightlyArtifacts(ctx, opts.NightlyBranch)
	if nightlyErr != nil {
		nightly = nil
	}
	body := RenderPRComment(result, nightly, opts.RunURL, opts.PRArtifactURL)
	return client.UpsertComment(ctx, opts.PRNumber, "<!-- coverage-patch -->", body)
}

// RenderPRComment renders only structured data. Nightly artifacts are always
// identified as a develop baseline, never as coverage of the PR itself.
func RenderPRComment(result PRResult, nightly []Artifact, runURL, prArtifactURL string) string {
	var b strings.Builder
	b.WriteString("<!-- coverage-patch -->\n\n### Quick PR diff coverage\n\n")
	fmt.Fprintf(&b, "PR head: `%s`.\n\n", result.HeadSHA)
	switch {
	case result.NoGoChanges:
		b.WriteString("No coverable Go changes; no package tests were needed.\n\n")
	case result.NoRunnable:
		b.WriteString("No runnable affected Go packages were found.\n\n")
	default:
		b.WriteString("Affected-package coverage")
		if result.Short {
			b.WriteString(" using `-short`")
		}
		b.WriteString(". This is a fast, partial signal—not full CI or global coverage.\n\n")
	}
	if result.TestsFailed {
		b.WriteString("> **Warning:** some affected-package tests failed. Coverage may be incomplete; inspect the test log artifact.\n\n")
	}
	if prArtifactURL != "" {
		fmt.Fprintf(&b, "[Download this PR's coverage report and test log](%s).\n\n", prArtifactURL)
	} else if runURL != "" {
		fmt.Fprintf(&b, "[Open this PR's coverage run](%s).\n\n", runURL)
	}
	if result.Report != nil {
		b.WriteString(RenderMarkdown(*result.Report))
		b.WriteString("\n")
	}
	if len(result.NotMeasured) > 0 {
		b.WriteString("### Not measured by the quick run\n\n")
		b.WriteString("These files had no coverable statements in the affected-package result. Use the full CI/nightly reports before concluding that gaps are real.\n\n")
		for _, path := range result.NotMeasured {
			fmt.Fprintf(&b, "- `%s`\n", path)
		}
		b.WriteString("\n")
	}
	b.WriteString("### Latest nightly `develop` artifacts\n\n")
	b.WriteString("This is the authoritative full **baseline only; does not include this PR**. Race coverage is stored separately and never merged into normal coverage.\n\n")
	if len(nightly) == 0 {
		b.WriteString("No unexpired nightly artifacts are currently available.\n")
		return b.String()
	}
	for _, artifact := range nightly {
		label := nightlyLabel(artifact.Name)
		if label == "" {
			continue
		}
		fmt.Fprintf(&b, "- [%s](%s) — `%s`, %s ([run](%s))\n", label, artifact.URL, shortSHA(artifact.SHA), artifact.CreatedAt.UTC().Format("2006-01-02 15:04 UTC"), artifact.RunURL)
	}
	return b.String()
}

func productionGoFile(path string) bool {
	if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.HasPrefix(path, "vendor/") {
		return false
	}
	base := filepath.Base(path)
	if strings.HasSuffix(base, ".pb.go") || strings.HasSuffix(base, "_gen.go") || strings.HasSuffix(base, "_codecgen.go") {
		return false
	}
	for part := range strings.SplitSeq(path, "/") {
		switch part {
		case "mocks", "testdata", "testutils", "testhelpers", "cltest", "fixtures":
			return false
		}
	}
	return true
}

func nearestModule(root, dir string) string {
	current := filepath.Join(root, dir)
	for {
		if info, err := os.Stat(filepath.Join(current, "go.mod")); err == nil && !info.IsDir() {
			return current
		}
		if current == root {
			return root
		}
		parent := filepath.Dir(current)
		if parent == current || !strings.HasPrefix(parent, root) {
			return root
		}
		current = parent
	}
}

func affectedPackages(ctx context.Context, module string, dirs []string) ([]string, error) {
	out, err := commandOutput(ctx, module, "go", "list", "-e", "-f", "{{.Dir}}\t{{.ImportPath}}", "./...")
	if err != nil {
		return nil, err
	}
	byDir := make(map[string]string)
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			byDir[filepath.Clean(parts[0])] = parts[1]
		}
	}
	var packages []string
	for _, dir := range dirs {
		if pkg := byDir[filepath.Clean(dir)]; pkg != "" {
			packages = append(packages, pkg)
		}
	}
	return uniqueStrings(packages), nil
}

func uniqueStrings(values []string) []string {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value != "" {
			set[value] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func hasCovdata(root string) bool {
	found := false
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && strings.HasPrefix(entry.Name(), "covmeta.") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func commandOutput(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // arguments are passed without a shell
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func writePRResult(dir string, result PRResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "result.json"), append(data, '\n'), 0o600)
}

func nightlyLabel(name string) string {
	switch {
	case strings.HasPrefix(name, "coverage-report-"):
		return "Full coverage report"
	case strings.HasPrefix(name, "coverage-core-go_core_race_tests-"):
		return "Race coverage (separate; never merged)"
	case strings.HasPrefix(name, "coverage-core-go_core_tests-"):
		return "Unit coverage"
	case strings.HasPrefix(name, "coverage-deployment-"), strings.HasPrefix(name, "coverage-inmemory-"):
		return "Integration coverage"
	case strings.HasPrefix(name, "coverage-e2e-"):
		return "E2E coverage"
	default:
		return ""
	}
}

func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}
