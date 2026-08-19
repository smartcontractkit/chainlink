package main_test

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Scenario struct {
	Name          string
	TargetSummary string
	Files         []string
	AllFiles      bool
}

type ParsedHook struct {
	StepTimes map[string]time.Duration
	TotalTime time.Duration
	RawOutput string
}

type ScenarioTiming struct {
	Name              string
	TargetSummary     string
	Generate          time.Duration
	Tidy              time.Duration
	Lint              time.Duration
	ShortTests        time.Duration
	LefthookPreCommit time.Duration
	LefthookPrePush   time.Duration
	TotalWallClock    time.Duration
}

var (
	ansiRegex    = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	stepRegex    = regexp.MustCompile(`(?m)^\s*[^\w]*([a-zA-Z0-9_-]+)(?::[^(]+)?\s*\(([0-9.]+)\s*seconds?\)`)
	summaryRegex = regexp.MustCompile(`(?i)summary:\s*\(done in\s+([0-9.]+)\s*seconds?\)`)
)

func parseLefthookOutput(output string, wallClock time.Duration) ParsedHook {
	res := ParsedHook{
		StepTimes: make(map[string]time.Duration),
		TotalTime: wallClock,
		RawOutput: output,
	}

	cleanOutput := ansiRegex.ReplaceAllString(output, "")

	for line := range strings.SplitSeq(cleanOutput, "\n") {
		line = strings.TrimSpace(line)
		if sm := summaryRegex.FindStringSubmatch(line); len(sm) == 2 {
			if sec, err := strconv.ParseFloat(sm[1], 64); err == nil {
				res.TotalTime = time.Duration(sec * float64(time.Second))
			}
			continue
		}
		if m := stepRegex.FindStringSubmatch(line); len(m) == 3 {
			stepName := strings.TrimSpace(m[1])
			if sec, err := strconv.ParseFloat(m[2], 64); err == nil {
				res.StepTimes[stepName] = time.Duration(sec * float64(time.Second))
			}
		}
	}

	return res
}

func TestParseLefthookOutput(t *testing.T) {
	t.Parallel()

	sample := `summary: (done in 3.51 seconds)
✔️ prewarm-model (0.01 seconds)
✔️ trailing-whitespace (0.02 seconds)
✔️ end-of-file-fixer (0.02 seconds)
✔️ go-mod-tidy (0.76 seconds)
✔️ generate (1.27 seconds)
✔️ golangci-lint (2.58 seconds)
🥊 betterleaks: betterleaks failed or missing. Install: go install github.com/betterleaks/betterleaks@latest (0.65 seconds)
🥊 typos: typos failed or missing. Install: brew install typos-cli (or cargo install typos-cli) (3.39 seconds)`

	parsed := parseLefthookOutput(sample, 4*time.Second)
	assert.Equal(t, 3510*time.Millisecond, parsed.TotalTime)
	assert.Equal(t, 1270*time.Millisecond, parsed.StepTimes["generate"])
	assert.Equal(t, 760*time.Millisecond, parsed.StepTimes["go-mod-tidy"])
	assert.Equal(t, 2580*time.Millisecond, parsed.StepTimes["golangci-lint"])
}

func getRepoRoot(ctx context.Context, t testing.TB) string {
	t.Helper()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	require.NoError(t, err)
	return strings.TrimSpace(string(out))
}

func runLefthook(ctx context.Context, dir string, hook string, allFiles bool, files []string) ParsedHook {
	args := []string{"run", hook, "--force"}
	if allFiles {
		args = append(args, "--all-files")
	} else {
		for _, f := range files {
			args = append(args, "--file", f)
		}
	}

	cmd := exec.CommandContext(ctx, "lefthook", args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	start := time.Now()
	_ = cmd.Run()
	elapsed := time.Since(start)

	return parseLefthookOutput(buf.String(), elapsed)
}

func runE2EBenchmark(t testing.TB) string {
	ctx := context.Background()
	repoRoot := getRepoRoot(ctx, t)

	scenarios := []Scenario{
		{
			Name:          "Small Change (1 pkg)",
			TargetSummary: "`core/logger`",
			Files:         []string{"core/logger/logger.go"},
		},
		{
			Name:          "Medium Change (3 pkgs)",
			TargetSummary: "`core/logger`, `core/services/cron`, `tools/ci-testshard`",
			Files: []string{
				"core/logger/logger.go",
				"core/services/cron/cron.go",
				"tools/ci-testshard/main.go",
			},
		},
		{
			Name:          "Large Change (6 pkgs)",
			TargetSummary: "`core/{logger,cron,telemetry,arbiter}`, `tools/{ci-testshard,githooks}`",
			Files: []string{
				"core/logger/logger.go",
				"core/services/cron/cron.go",
				"core/services/telemetry/ingress.go",
				"core/services/arbiter/arbiter.go",
				"tools/ci-testshard/main.go",
				"tools/githooks/main.go",
			},
		},
		{
			Name:          "Full Repo Baseline",
			TargetSummary: "All monorepo modules / packages",
			AllFiles:      true,
		},
	}

	results := make([]ScenarioTiming, 0, len(scenarios))

	for _, sc := range scenarios {
		t.Logf("Running Lefthook E2E for %s...", sc.Name)
		startScenario := time.Now()

		preCommit := runLefthook(ctx, repoRoot, "pre-commit", sc.AllFiles, sc.Files)
		prePush := runLefthook(ctx, repoRoot, "pre-push", sc.AllFiles, sc.Files)

		totalWallClock := time.Since(startScenario)

		results = append(results, ScenarioTiming{
			Name:              sc.Name,
			TargetSummary:     sc.TargetSummary,
			Generate:          preCommit.StepTimes["generate"],
			Tidy:              preCommit.StepTimes["go-mod-tidy"],
			Lint:              preCommit.StepTimes["golangci-lint"],
			ShortTests:        prePush.StepTimes["short-tests"],
			LefthookPreCommit: preCommit.TotalTime,
			LefthookPrePush:   prePush.TotalTime,
			TotalWallClock:    totalWallClock,
		})
	}

	return buildMarkdownTable(results)
}

func TestE2EBenchmark(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping E2E benchmark in short mode")
	}

	table := runE2EBenchmark(t)
	fmt.Println("\n" + table)
}

func BenchmarkE2E(b *testing.B) {
	table := runE2EBenchmark(b)
	b.Log("\n" + table)
}

func formatDur(d time.Duration) string {
	if d == 0 {
		return "0.00s"
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func buildMarkdownTable(results []ScenarioTiming) string {
	var sb strings.Builder
	sb.WriteString("| Change Scope | Target Packages / Files | `generate` | `go-mod-tidy` | `golangci-lint` | `short-tests` | Lefthook `pre-commit` | Lefthook `pre-push` | Total Wall Clock |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |\n")

	for _, r := range results {
		fmt.Fprintf(
			&sb,
			"| **%s** | %s | %s | %s | %s | %s | **%s** | **%s** | %s |\n",
			r.Name,
			r.TargetSummary,
			formatDur(r.Generate),
			formatDur(r.Tidy),
			formatDur(r.Lint),
			formatDur(r.ShortTests),
			formatDur(r.LefthookPreCommit),
			formatDur(r.LefthookPrePush),
			formatDur(r.TotalWallClock),
		)
	}

	return sb.String()
}
