package prsize

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

// Classification represents the size category of a PR diff.
type Classification string

const (
	SizeSmall  Classification = "SMALL"
	SizeMedium Classification = "MEDIUM"
	SizeLarge  Classification = "LARGE"
)

// Strategy defines how lines changed are calculated from additions and deletions.
type Strategy string

const (
	StrategyPerFileMax Strategy = "per-file-max"
	StrategySum        Strategy = "sum"
	StrategyMax        Strategy = "max"
	StrategyWeighted   Strategy = "weighted"
)

// Default limits for diff classification.
const (
	DefaultSmallLimit  = 200
	DefaultMediumLimit = 500
)

var lockfileNames = map[string]bool{
	"go.sum":            true,
	"package-lock.json": true,
	"pnpm-lock.yaml":    true,
	"yarn.lock":         true,
	"gemfile.lock":      true,
	"cargo.lock":        true,
	"poetry.lock":       true,
	"composer.lock":     true,
	"flake.lock":        true,
}

// FileStat contains diff statistics for a single file.
type FileStat struct {
	Path      string
	Additions int
	Deletions int
	IsBinary  bool
}

// DiffStat summarizes the diff analysis results.
type DiffStat struct {
	Files          []FileStat
	IgnoredFiles   []string
	FilesChanged   int
	Additions      int
	Deletions      int
	EffectiveLines int
	MergeBase      string
	BaseRef        string
}

// Config controls PR size calculation and reporting behavior.
type Config struct {
	RepoRoot           string
	BaseRef            string
	Strategy           Strategy
	SmallLimit         int
	MediumLimit        int
	FailOnLarge        bool
	IgnoreLockfiles    bool
	IgnoreGenerated    bool
	IgnoreGlobs        []string
	IncludeUncommitted bool
	Stdout             io.Writer
	Stderr             io.Writer
}

func isLockfile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return lockfileNames[base]
}

// CalculateEffectiveLines computes the effective lines changed based on the chosen strategy.
func CalculateEffectiveLines(files []FileStat, strategy Strategy) int {
	if len(files) == 0 {
		return 0
	}

	switch strategy {
	case StrategySum:
		total := 0
		for _, f := range files {
			total += f.Additions + f.Deletions
		}
		return total

	case StrategyMax:
		totalAdd := 0
		totalDel := 0
		for _, f := range files {
			totalAdd += f.Additions
			totalDel += f.Deletions
		}
		return max(totalAdd, totalDel)

	case StrategyWeighted:
		totalAdd := 0
		totalDel := 0
		for _, f := range files {
			totalAdd += f.Additions
			totalDel += f.Deletions
		}
		return totalAdd + int(0.5*float64(totalDel))

	case StrategyPerFileMax:
		fallthrough
	default:
		total := 0
		for _, f := range files {
			total += max(f.Additions, f.Deletions)
		}
		return total
	}
}

// Classify categorizes line count into Small, Medium, or Large.
func Classify(lines int, cfg Config) Classification {
	smallLimit := cfg.SmallLimit
	if smallLimit <= 0 {
		smallLimit = DefaultSmallLimit
	}
	mediumLimit := cfg.MediumLimit
	if mediumLimit <= 0 {
		mediumLimit = DefaultMediumLimit
	}
	if smallLimit > mediumLimit {
		smallLimit = mediumLimit
	}

	if lines <= smallLimit {
		return SizeSmall
	}
	if lines <= mediumLimit {
		return SizeMedium
	}
	return SizeLarge
}

func resolveMergeBase(ctx context.Context, repoRoot, baseRef string) (string, error) {
	if baseRef != "" {
		cmd := exec.CommandContext(ctx, "git", "merge-base", baseRef, "HEAD")
		cmd.Dir = repoRoot
		out, err := cmd.Output()
		if err == nil {
			if sha := strings.TrimSpace(string(out)); sha != "" {
				return sha, nil
			}
		}
		return baseRef, nil
	}

	mb := modules.GetMergeBase(ctx, repoRoot)
	return mb, nil
}

func checkGeneratedFiles(ctx context.Context, repoRoot string, files []string) (map[string]bool, error) {
	if len(files) == 0 {
		return nil, nil
	}

	cmd := exec.CommandContext(ctx, "git", "check-attr", "--stdin", "linguist-generated")
	cmd.Dir = repoRoot
	cmd.Stdin = strings.NewReader(strings.Join(files, "\n") + "\n")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git check-attr failed: %w", err)
	}

	generatedMap := make(map[string]bool)
	for line := range strings.SplitSeq(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		const prefix = ": linguist-generated: "
		if filePath, val, found := strings.Cut(trimmed, prefix); found {
			val = strings.TrimSpace(val)
			if val == "true" || val == "set" || val == "1" {
				generatedMap[strings.TrimSpace(filePath)] = true
			}
		}
	}
	return generatedMap, nil
}

// Analyze inspects git diff against the merge-base with the default branch and calculates stats.
func Analyze(ctx context.Context, cfg Config) (*DiffStat, Classification, error) {
	repoRoot := cfg.RepoRoot
	if repoRoot == "" {
		repoRoot = "."
	}

	mergeBase, err := resolveMergeBase(ctx, repoRoot, cfg.BaseRef)
	if err != nil {
		return nil, "", fmt.Errorf("failed to resolve merge-base: %w", err)
	}

	diffArgs := []string{"diff", "--numstat", mergeBase}
	if !cfg.IncludeUncommitted {
		diffArgs = append(diffArgs, "HEAD")
	}

	cmd := exec.CommandContext(ctx, "git", diffArgs...)
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if runErr := cmd.Run(); runErr != nil {
		// If HEAD diff failed (e.g. initial commit or uncommitted only), fallback to diffing mergeBase
		if cfg.IncludeUncommitted {
			return nil, "", fmt.Errorf("git diff failed: %w (stderr: %s)", runErr, stderr.String())
		}
		cmd = exec.CommandContext(ctx, "git", "diff", "--numstat", mergeBase)
		cmd.Dir = repoRoot
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if fallbackErr := cmd.Run(); fallbackErr != nil {
			return nil, "", fmt.Errorf("git diff failed: %w (stderr: %s)", fallbackErr, stderr.String())
		}
	}

	type rawDiffEntry struct {
		filePath string
		addStr   string
		delStr   string
	}

	var rawEntries []rawDiffEntry
	var allPaths []string

	for line := range strings.SplitSeq(stdout.String(), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		parts := strings.Split(trimmed, "\t")
		if len(parts) < 3 {
			continue
		}

		addStr, delStr, filePath := parts[0], parts[1], parts[2]

		// Handle renames formatted as "old => new" or "{dir => dir2}/file"
		if idx := strings.LastIndex(filePath, " => "); idx != -1 {
			filePath = strings.TrimSpace(filePath[idx+4:])
			filePath = strings.TrimSuffix(filePath, "}")
		}

		rawEntries = append(rawEntries, rawDiffEntry{filePath: filePath, addStr: addStr, delStr: delStr})
		allPaths = append(allPaths, filePath)
	}

	var generatedMap map[string]bool
	if cfg.IgnoreGenerated && len(allPaths) > 0 {
		generatedMap, err = checkGeneratedFiles(ctx, repoRoot, allPaths)
		if err != nil {
			return nil, "", err
		}
	}

	stat := &DiffStat{
		MergeBase: mergeBase,
		BaseRef:   cfg.BaseRef,
	}

	for _, entry := range rawEntries {
		filePath := entry.filePath

		// Check gitattributes linguist-generated ignore
		if cfg.IgnoreGenerated && generatedMap[filePath] {
			stat.IgnoredFiles = append(stat.IgnoredFiles, filePath)
			continue
		}

		// Check lockfile ignore
		if cfg.IgnoreLockfiles && isLockfile(filePath) {
			stat.IgnoredFiles = append(stat.IgnoredFiles, filePath)
			continue
		}

		// Check custom ignore globs
		matchedGlob := false
		for _, glob := range cfg.IgnoreGlobs {
			if matched, _ := filepath.Match(glob, filePath); matched {
				matchedGlob = true
				break
			}
		}
		if matchedGlob {
			stat.IgnoredFiles = append(stat.IgnoredFiles, filePath)
			continue
		}

		fileStat := FileStat{Path: filePath}
		if entry.addStr == "-" && entry.delStr == "-" {
			fileStat.IsBinary = true
		} else {
			fileStat.Additions, _ = strconv.Atoi(entry.addStr)
			fileStat.Deletions, _ = strconv.Atoi(entry.delStr)
		}

		stat.Files = append(stat.Files, fileStat)
		stat.FilesChanged++
		stat.Additions += fileStat.Additions
		stat.Deletions += fileStat.Deletions
	}

	diffStrategy := cfg.Strategy
	if diffStrategy == "" {
		diffStrategy = StrategyPerFileMax
	}
	stat.EffectiveLines = CalculateEffectiveLines(stat.Files, diffStrategy)
	classification := Classify(stat.EffectiveLines, cfg)

	return stat, classification, nil
}

// FormatReport generates the human-readable diff size summary and guidance.
func FormatReport(stat *DiffStat, class Classification, cfg Config) string {
	var sb strings.Builder

	mediumLimit := cfg.MediumLimit
	if mediumLimit <= 0 {
		mediumLimit = DefaultMediumLimit
	}

	diffStrategy := cfg.Strategy
	if diffStrategy == "" {
		diffStrategy = StrategyPerFileMax
	}

	if class != SizeLarge {
		fmt.Fprintf(&sb, "PR Diff Size: %d effective lines (+%d, -%d in %d files) [strategy: %s] -> Classification: %s [OK]\n",
			stat.EffectiveLines, stat.Additions, stat.Deletions, stat.FilesChanged, diffStrategy, class)
		if len(stat.IgnoredFiles) > 0 {
			fmt.Fprintf(&sb, "  (Excluded %d lock/ignored files)\n", len(stat.IgnoredFiles))
		}
		return sb.String()
	}

	// Large diff report & prompt
	fmt.Fprintf(&sb, "⚠️  [pr-size] LARGE PR: %d effective lines (+%d/-%d across %d files) exceeds limit (%d)\n",
		stat.EffectiveLines, stat.Additions, stat.Deletions, stat.FilesChanged, mediumLimit)
	if len(stat.IgnoredFiles) > 0 {
		fmt.Fprintf(&sb, "   (Excluded %d lock/ignored files: %s)\n", len(stat.IgnoredFiles), strings.Join(stat.IgnoredFiles, ", "))
	}
	sb.WriteString("   Split into smaller, focused PRs or a GitHub PR stack.\n")
	sb.WriteString("   AI skill to split PR: tools/githooks/skills/split-pr/SKILL.md\n")

	return sb.String()
}

// Run performs PR size analysis and reports results according to config.
func Run(ctx context.Context, cfg Config) error {
	if cfg.Stdout == nil {
		cfg.Stdout = io.Discard
	}
	if cfg.Stderr == nil {
		cfg.Stderr = io.Discard
	}

	stat, class, err := Analyze(ctx, cfg)
	if err != nil {
		return err
	}

	report := FormatReport(stat, class, cfg)
	if class == SizeLarge {
		if cfg.FailOnLarge {
			fmt.Fprint(cfg.Stderr, report)
			return fmt.Errorf("PR size check failed: diff is classified as LARGE (%d effective lines > %d limit)", stat.EffectiveLines, cfg.MediumLimit)
		}
		fmt.Fprint(cfg.Stdout, report)
		return nil
	}

	fmt.Fprint(cfg.Stdout, report)
	return nil
}
