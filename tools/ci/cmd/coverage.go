package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/coverage"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
)

func newCoverageCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "coverage", Short: "Compute, collect, and publish Go coverage"}
	cmd.AddCommand(
		newCoveragePRCmd(),
		newCoveragePublishPRCmd(),
		newCoverageCollectCmd(),
		newCoverageBadgeCmd(),
		newCoveragePublishBadgeCmd(),
	)
	return cmd
}

func newCoveragePRCmd() *cobra.Command {
	var opts coverage.PROptions
	var short bool
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Run lightweight affected-package coverage for a PR diff",
		RunE: func(cmd *cobra.Command, _ []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			if !cmd.Flags().Changed("base") {
				opts.BaseSHA = act.Getenv("BASE_SHA")
			}
			if !cmd.Flags().Changed("head") {
				opts.HeadSHA = act.Getenv("HEAD_SHA")
			}
			if !cmd.Flags().Changed("short") && act.Getenv("SHORT") != "" {
				value, err := strconv.ParseBool(act.Getenv("SHORT"))
				if err != nil {
					return fmt.Errorf("parse SHORT: %w", err)
				}
				short = value
			}
			opts.Short = short
			if opts.PRNumber == 0 && act.Getenv("PR_NUMBER") != "" {
				value, err := strconv.Atoi(act.Getenv("PR_NUMBER"))
				if err != nil {
					return fmt.Errorf("parse PR_NUMBER: %w", err)
				}
				opts.PRNumber = value
			}
			if opts.RepoRoot == "" {
				root, err := FindRepoRoot(cmd.Context())
				if err != nil {
					return err
				}
				opts.RepoRoot = root
			}
			if opts.BaseSHA == "" || opts.HeadSHA == "" {
				return errors.New("--base and --head (or BASE_SHA and HEAD_SHA) are required")
			}
			result, runErr := coverage.RunPR(cmd.Context(), opts)
			if jsonOutput {
				if err := writeCoverageJSON(cmd, result); err != nil {
					return err
				}
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "PR coverage result: %s\n", filepath.Join(opts.OutDir, "result.json"))
				fmt.Fprintf(cmd.OutOrStdout(), "packages=%d no_go_changes=%t tests_failed=%t\n", len(result.Packages), result.NoGoChanges, result.TestsFailed)
			}
			if act.Getenv("GITHUB_OUTPUT") != "" {
				if err := act.SetOutputs(map[string]string{
					"head_sha":      result.HeadSHA,
					"no_go_changes": strconv.FormatBool(result.NoGoChanges),
					"tests_failed":  strconv.FormatBool(result.TestsFailed),
				}); err != nil {
					return err
				}
			}
			return runErr
		},
	}
	cmd.Flags().StringVar(&opts.RepoRoot, "repo", "", "Repository root (default: current Git repository)")
	cmd.Flags().StringVar(&opts.BaseSHA, "base", "", "PR base SHA (env: BASE_SHA)")
	cmd.Flags().StringVar(&opts.HeadSHA, "head", "", "PR head SHA (env: HEAD_SHA)")
	cmd.Flags().IntVar(&opts.PRNumber, "pr", 0, "PR number (env: PR_NUMBER)")
	cmd.Flags().StringVar(&opts.OutDir, "out", "coverage-patch", "Output directory")
	cmd.Flags().BoolVar(&short, "short", true, "Run affected tests with -short (env: SHORT)")
	cmd.Flags().DurationVar(&opts.TestTimeout, "test-timeout", 10*time.Minute, "Per-package Go test timeout")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Write the structured result to stdout")
	return cmd
}

func newCoveragePublishPRCmd() *cobra.Command {
	var opts coverage.PublishPROptions
	var repo, token, apiURL, webURL string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "publish-pr",
		Short: "Publish a trusted sticky PR comment with nightly artifact links",
		RunE: func(cmd *cobra.Command, _ []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			if opts.PRNumber == 0 && act.Getenv("PR_NUMBER") != "" {
				value, err := strconv.Atoi(act.Getenv("PR_NUMBER"))
				if err != nil {
					return fmt.Errorf("parse PR_NUMBER: %w", err)
				}
				opts.PRNumber = value
			}
			if opts.RunID == 0 && act.Getenv("PR_RUN_ID") != "" {
				value, err := strconv.ParseInt(act.Getenv("PR_RUN_ID"), 10, 64)
				if err != nil {
					return fmt.Errorf("parse PR_RUN_ID: %w", err)
				}
				opts.RunID = value
			}
			if opts.RunURL == "" {
				opts.RunURL = act.Getenv("PR_RUN_URL")
			}
			repo, token = coverageGitHubConfig(act, repo, token)
			client := coverage.NewGitHubClient(apiURL, webURL, token, repo, http.DefaultClient)
			if err := coverage.PublishPR(cmd.Context(), client, opts); err != nil {
				return err
			}
			if jsonOutput {
				return writeCoverageJSON(cmd, map[string]any{"pr": opts.PRNumber, "published": true})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Published coverage comment on PR #%d\n", opts.PRNumber)
			return nil
		},
	}
	cmd.Flags().IntVar(&opts.PRNumber, "pr", 0, "PR number (env: PR_NUMBER)")
	cmd.Flags().StringVar(&opts.ResultPath, "result", "coverage-patch/result.json", "Downloaded PR result JSON")
	cmd.Flags().Int64Var(&opts.RunID, "run-id", 0, "Coverage PR workflow run ID (env: PR_RUN_ID)")
	cmd.Flags().StringVar(&opts.RunURL, "run-url", "", "Coverage PR workflow run URL (env: PR_RUN_URL)")
	cmd.Flags().StringVar(&opts.NightlyBranch, "nightly-branch", "develop", "Nightly baseline branch")
	cmd.Flags().StringVar(&repo, "repository", "", "GitHub owner/repository (env: GITHUB_REPOSITORY)")
	cmd.Flags().StringVar(&token, "token", "", "GitHub token (env: GITHUB_TOKEN or GH_TOKEN)")
	cmd.Flags().StringVar(&apiURL, "api-url", "https://api.github.com", "GitHub API URL")
	cmd.Flags().StringVar(&webURL, "web-url", "https://github.com", "GitHub web URL")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Write publication result as JSON")
	return cmd
}

func newCoverageCollectCmd() *cobra.Command {
	var opts coverage.CollectOptions
	var repo, token, apiURL, webURL string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect heavy scheduled coverage into a full report artifact",
		RunE: func(cmd *cobra.Command, _ []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			if opts.TargetSHA == "" {
				opts.TargetSHA = act.Getenv("TARGET_SHA")
			}
			if opts.BaseSHA == "" {
				opts.BaseSHA = act.Getenv("BASE_SHA")
			}
			if opts.RepoRoot == "" {
				root, err := FindRepoRoot(cmd.Context())
				if err != nil {
					return err
				}
				opts.RepoRoot = root
			}
			if opts.ArtifactsDir == "" {
				opts.ArtifactsDir = filepath.Join(os.TempDir(), "coverage-artifacts")
			}
			repo, token = coverageGitHubConfig(act, repo, token)
			opts.Client = coverage.NewGitHubClient(apiURL, webURL, token, repo, http.DefaultClient)
			result, err := coverage.Collect(cmd.Context(), opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				if err := writeCoverageJSON(cmd, result); err != nil {
					return err
				}
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Coverage report: %s\n", filepath.Join(opts.OutDir, "summary.md"))
				fmt.Fprintf(cmd.OutOrStdout(), "full=%t categories=%s tests=%d failures=%d\n", result.Full, strings.Join(result.Categories, ","), result.TestResults.Tests, result.TestResults.Failures)
				for _, warning := range result.Warnings {
					fmt.Fprintf(cmd.OutOrStdout(), "warning: %s\n", warning)
				}
			}
			if act.Getenv("GITHUB_OUTPUT") != "" {
				if err := act.SetOutputs(map[string]string{"full": strconv.FormatBool(result.Full), "categories": strings.Join(result.Categories, ",")}); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.RepoRoot, "repo", "", "Repository root (default: current Git repository)")
	cmd.Flags().StringVar(&opts.TargetSHA, "target-sha", "", "Commit whose artifacts are collected (env: TARGET_SHA)")
	cmd.Flags().StringVar(&opts.BaseSHA, "base", "", "Optional diff base SHA (env: BASE_SHA)")
	cmd.Flags().StringVar(&opts.OutDir, "out", "coverage-report", "Report output directory")
	cmd.Flags().StringVar(&opts.ArtifactsDir, "artifacts", "", "Raw artifact extraction directory")
	cmd.Flags().BoolVar(&opts.Wait, "wait", true, "Wait for watched workflows to finish")
	cmd.Flags().DurationVar(&opts.WaitTimeout, "wait-timeout", 45*time.Minute, "Maximum workflow wait")
	cmd.Flags().DurationVar(&opts.PollInterval, "poll-interval", 30*time.Second, "Workflow polling interval")
	cmd.Flags().StringVar(&repo, "repository", "", "GitHub owner/repository (env: GITHUB_REPOSITORY)")
	cmd.Flags().StringVar(&token, "token", "", "GitHub token (env: GITHUB_TOKEN or GH_TOKEN)")
	cmd.Flags().StringVar(&apiURL, "api-url", "https://api.github.com", "GitHub API URL")
	cmd.Flags().StringVar(&webURL, "web-url", "https://github.com", "GitHub web URL")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Write collection result as JSON")
	return cmd
}

func newCoverageBadgeCmd() *cobra.Command {
	var opts coverage.BadgeOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Generate the coverage SVG and latest-report JSON",
		RunE: func(cmd *cobra.Command, _ []string) error {
			full, err := coverage.FullReport(filepath.Dir(opts.ReportPath))
			if err != nil {
				return err
			}
			if !full {
				return errors.New("refusing to generate a badge from a partial report")
			}
			state, err := coverage.GenerateBadge(opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeCoverageJSON(cmd, state)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "coverage badge: %.1f%%\n", state.GlobalSourceOnly)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.ReportPath, "report", "coverage-report/any.json", "Merged any.json report")
	cmd.Flags().StringVar(&opts.SVGPath, "svg", "coverage-badge.svg", "SVG output path")
	cmd.Flags().StringVar(&opts.StatePath, "state", "coverage-state.json", "State JSON output path")
	cmd.Flags().StringVar(&opts.SHA, "sha", "", "Nightly develop SHA")
	cmd.Flags().StringVar(&opts.Artifact, "artifact", "", "Merged report artifact name")
	cmd.Flags().StringVar(&opts.Repo, "repository", "", "GitHub owner/repository")
	cmd.Flags().StringVar(&opts.RunURL, "run-url", "", "Collector run URL")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Write badge state to stdout")
	return cmd
}

func newCoveragePublishBadgeCmd() *cobra.Command {
	var opts coverage.PublishBadgeOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "publish-badge",
		Short: "Publish badge files to the orphan coverage branch",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.RepoRoot == "" {
				root, err := FindRepoRoot(cmd.Context())
				if err != nil {
					return err
				}
				opts.RepoRoot = root
			}
			if opts.WorktreeDir == "" {
				opts.WorktreeDir = filepath.Join(os.TempDir(), "coverage-branch")
			}
			changed, err := coverage.PublishBadge(cmd.Context(), opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeCoverageJSON(cmd, map[string]any{"changed": changed, "branch": opts.Branch})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "changed=%t\n", changed)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.RepoRoot, "repo", "", "Repository root (default: current Git repository)")
	cmd.Flags().StringVar(&opts.WorktreeDir, "worktree", "", "Temporary coverage-branch worktree")
	cmd.Flags().StringVar(&opts.SVGPath, "svg", "coverage-badge.svg", "Generated SVG")
	cmd.Flags().StringVar(&opts.StatePath, "state", "coverage-state.json", "Generated state JSON")
	cmd.Flags().StringVar(&opts.SHA, "sha", "", "Nightly develop SHA")
	cmd.Flags().StringVar(&opts.Branch, "branch", "coverage", "Coverage branch")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Write publication result as JSON")
	return cmd
}

func coverageGitHubConfig(act *ghaction.Action, repo, token string) (string, string) {
	if repo == "" {
		repo = act.Getenv("GITHUB_REPOSITORY")
	}
	if token == "" {
		token = act.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = act.Getenv("GH_TOKEN")
	}
	return repo, token
}

func writeCoverageJSON(cmd *cobra.Command, value any) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
