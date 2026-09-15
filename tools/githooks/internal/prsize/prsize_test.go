package prsize_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/prsize"
)

// git runs a git command in dir and fails the test on error.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()

	base := make([]string, 0, 6+len(args))
	base = append(base, "-c", "commit.gpgsign=false", "-c", "user.email=test@example.com", "-c", "user.name=test")
	cmd := exec.CommandContext(t.Context(), "git", append(base, args...)...) //nolint:gosec // test helper
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %s: %s", strings.Join(args, " "), out)
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o700))
	require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
}

func generateLines(n int, prefix string) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(prefix)
		sb.WriteString(" line\n")
	}
	return sb.String()
}

func TestDiffCalculationStrategies(t *testing.T) {
	t.Parallel()

	files := []prsize.FileStat{
		{Path: "core/services/app.go", Additions: 100, Deletions: 80},
		{Path: "core/logger/logger.go", Additions: 50, Deletions: 10},
		{Path: "README.md", Additions: 20, Deletions: 0},
	}

	tests := []struct {
		strategy prsize.Strategy
		expected int
	}{
		{
			strategy: prsize.StrategyPerFileMax,
			// max(100, 80) + max(50, 10) + max(20, 0) = 100 + 50 + 20 = 170
			expected: 170,
		},
		{
			strategy: prsize.StrategySum,
			// (100 + 80) + (50 + 10) + (20 + 0) = 180 + 60 + 20 = 260
			expected: 260,
		},
		{
			strategy: prsize.StrategyMax,
			// total adds: 170, total dels: 90 -> max(170, 90) = 170
			expected: 170,
		},
		{
			strategy: prsize.StrategyWeighted,
			// total adds: 170 + 0.5 * 90 = 215
			expected: 215,
		},
	}

	for _, tc := range tests {
		t.Run(string(tc.strategy), func(t *testing.T) {
			t.Parallel()
			got := prsize.CalculateEffectiveLines(files, tc.strategy)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestClassify(t *testing.T) {
	t.Parallel()

	cfg := prsize.Config{
		SmallLimit:  200,
		MediumLimit: 500,
	}

	tests := []struct {
		lines    int
		expected prsize.Classification
	}{
		{lines: 0, expected: prsize.SizeSmall},
		{lines: 50, expected: prsize.SizeSmall},
		{lines: 200, expected: prsize.SizeSmall},
		{lines: 201, expected: prsize.SizeMedium},
		{lines: 450, expected: prsize.SizeMedium},
		{lines: 500, expected: prsize.SizeMedium},
		{lines: 501, expected: prsize.SizeLarge},
		{lines: 1200, expected: prsize.SizeLarge},
	}

	for _, tc := range tests {
		got := prsize.Classify(tc.lines, cfg)
		assert.Equal(t, tc.expected, got, "lines=%d", tc.lines)
	}
}

func TestAnalyze_GitRepo(t *testing.T) {
	t.Parallel()

	t.Run("classifies small diff correctly", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "base.go", generateLines(10, "base"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "initial commit")
		baseSHA := git(t, dir, "rev-parse", "HEAD")

		git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

		writeFile(t, dir, "feature.go", generateLines(50, "feature"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "small feature")

		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreLockfiles: true,
		}

		stat, class, err := prsize.Analyze(t.Context(), cfg)
		require.NoError(t, err)
		assert.Equal(t, prsize.SizeSmall, class)
		assert.Equal(t, 50, stat.EffectiveLines)
		assert.Equal(t, 50, stat.Additions)
		assert.Equal(t, 0, stat.Deletions)
		assert.Equal(t, 1, stat.FilesChanged)
	})

	t.Run("classifies medium diff correctly", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "base.go", generateLines(10, "base"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "initial commit")
		baseSHA := git(t, dir, "rev-parse", "HEAD")

		git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

		writeFile(t, dir, "feature.go", generateLines(350, "feature"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "medium feature")

		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreLockfiles: true,
		}

		stat, class, err := prsize.Analyze(t.Context(), cfg)
		require.NoError(t, err)
		assert.Equal(t, prsize.SizeMedium, class)
		assert.Equal(t, 350, stat.EffectiveLines)
		assert.Equal(t, 1, stat.FilesChanged)
	})

	t.Run("classifies large diff correctly and respects lockfile ignore", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "base.go", generateLines(10, "base"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "initial commit")
		baseSHA := git(t, dir, "rev-parse", "HEAD")

		git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

		// Add 100 lines of code and 1000 lines of go.sum lockfile
		writeFile(t, dir, "feature.go", generateLines(100, "feature"))
		writeFile(t, dir, "go.sum", generateLines(1000, "checksum"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "add feature and deps")

		cfgWithIgnore := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreLockfiles: true,
		}

		stat, class, err := prsize.Analyze(t.Context(), cfgWithIgnore)
		require.NoError(t, err)
		assert.Equal(t, prsize.SizeSmall, class)
		assert.Equal(t, 100, stat.EffectiveLines)
		assert.Len(t, stat.IgnoredFiles, 1)

		cfgWithoutIgnore := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreLockfiles: false,
		}

		stat2, class2, err2 := prsize.Analyze(t.Context(), cfgWithoutIgnore)
		require.NoError(t, err2)
		assert.Equal(t, prsize.SizeLarge, class2)
		assert.Equal(t, 1100, stat2.EffectiveLines)
		assert.Empty(t, stat2.IgnoredFiles)
	})

	t.Run("ignores generated changes defined in .gitattributes", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, ".gitattributes", "*.generated.go linguist-generated=true\ngen/** linguist-generated\n")
		writeFile(t, dir, "base.go", generateLines(10, "base"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "initial commit")
		baseSHA := git(t, dir, "rev-parse", "HEAD")

		git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

		writeFile(t, dir, "regular.go", generateLines(50, "regular"))
		writeFile(t, dir, "service.generated.go", generateLines(800, "generated"))
		writeFile(t, dir, "gen/docs.txt", generateLines(400, "docs"))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "feature and generated code")

		cfgWithGeneratedIgnore := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreGenerated: true,
		}

		stat, class, err := prsize.Analyze(t.Context(), cfgWithGeneratedIgnore)
		require.NoError(t, err)
		assert.Equal(t, prsize.SizeSmall, class)
		assert.Equal(t, 50, stat.EffectiveLines)
		assert.ElementsMatch(t, []string{"gen/docs.txt", "service.generated.go"}, stat.IgnoredFiles)

		cfgWithoutGeneratedIgnore := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			IgnoreGenerated: false,
		}

		stat2, class2, err2 := prsize.Analyze(t.Context(), cfgWithoutGeneratedIgnore)
		require.NoError(t, err2)
		assert.Equal(t, prsize.SizeLarge, class2)
		assert.Equal(t, 1250, stat2.EffectiveLines)
		assert.Empty(t, stat2.IgnoredFiles)
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	git(t, dir, "init")
	writeFile(t, dir, "base.go", generateLines(10, "base"))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "initial commit")
	baseSHA := git(t, dir, "rev-parse", "HEAD")

	git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
	git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	// Commit 600 lines -> Large
	writeFile(t, dir, "huge.go", generateLines(600, "huge"))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "large commit")

	t.Run("warns on large diff when FailOnLarge is false", func(t *testing.T) {
		t.Parallel()
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     false,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.NoError(t, err)

		out := stdout.String() + stderr.String()
		assert.Contains(t, out, "LARGE PR")
		assert.Contains(t, out, "Execute the skill @tools/githooks/skills/split-pr/SKILL.md to break")
	})

	t.Run("fails on large diff when FailOnLarge is true", func(t *testing.T) {
		t.Parallel()
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     true,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PR size check failed")

		out := stdout.String() + stderr.String()
		assert.Contains(t, out, "LARGE PR")
		assert.Contains(t, out, "Execute the skill @tools/githooks/skills/split-pr/SKILL.md to break")
		assert.Contains(t, out, "ALLOW_LARGE_PR=true")
	})
}

func TestRun_AllowLargePR(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init")
	writeFile(t, dir, "base.go", generateLines(10, "base"))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "initial commit")
	baseSHA := git(t, dir, "rev-parse", "HEAD")

	git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
	git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	writeFile(t, dir, "huge.go", generateLines(600, "huge"))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "large commit")

	t.Run("bypasses failure with ALLOW_LARGE_PR=true env var", func(t *testing.T) {
		t.Setenv("ALLOW_LARGE_PR", "true")
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     true,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.NoError(t, err)

		out := stdout.String() + stderr.String()
		assert.Contains(t, out, "LARGE PR")
	})

	t.Run("bypasses failure with ALLOW_LARGE_PR=1 env var", func(t *testing.T) {
		t.Setenv("ALLOW_LARGE_PR", "1")
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     true,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.NoError(t, err)
	})

	t.Run("bypasses failure with Config.AllowLargePR = true", func(t *testing.T) {
		t.Setenv("ALLOW_LARGE_PR", "")
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     true,
			AllowLargePR:    true,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.NoError(t, err)
	})

	t.Run("fails when ALLOW_LARGE_PR=false", func(t *testing.T) {
		t.Setenv("ALLOW_LARGE_PR", "false")
		var stdout, stderr bytes.Buffer
		cfg := prsize.Config{
			RepoRoot:        dir,
			Strategy:        prsize.StrategyPerFileMax,
			SmallLimit:      200,
			MediumLimit:     500,
			FailOnLarge:     true,
			IgnoreLockfiles: true,
			Stdout:          &stdout,
			Stderr:          &stderr,
		}

		err := prsize.Run(t.Context(), cfg)
		require.Error(t, err)
	})
}
