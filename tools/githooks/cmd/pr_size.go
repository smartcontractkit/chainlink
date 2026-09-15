package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/prsize"
)

func newPRSizeCmd() *cobra.Command {
	var (
		strategy           string
		smallLimit         int
		mediumLimit        int
		failOnLarge        bool
		baseRef            string
		ignoreLockfiles    bool
		ignoreGenerated    bool
		includeUncommitted bool
	)

	cmd := &cobra.Command{
		Use:     "pr-size",
		Aliases: []string{"big-pr-guard", "pr-guard", "diff-guard", "diff-size"},
		Short:   "Check PR diff size against default branch and guard against oversized PRs",
		Long: "Calculates the git diff of the current branch against the default branch (the PR diff), " +
			"classifies the diff as small, medium, or large, and prompts developers with recommendations " +
			"to split large PRs into smaller, reviewable chunks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := findRepoRoot(ctx)
			if err != nil {
				return fmt.Errorf("could not find repo root: %w", err)
			}

			cfg := prsize.Config{
				RepoRoot:           repoRoot,
				BaseRef:            baseRef,
				Strategy:           prsize.Strategy(strategy),
				SmallLimit:         smallLimit,
				MediumLimit:        mediumLimit,
				FailOnLarge:        failOnLarge,
				IgnoreLockfiles:    ignoreLockfiles,
				IgnoreGenerated:    ignoreGenerated,
				IncludeUncommitted: includeUncommitted,
				Stdout:             cmd.OutOrStdout(),
				Stderr:             cmd.ErrOrStderr(),
			}

			return prsize.Run(ctx, cfg)
		},
	}

	cmd.Flags().StringVar(&strategy, "strategy", string(prsize.StrategyPerFileMax), "Diff size strategy: per-file-max, sum, max, weighted")
	cmd.Flags().IntVar(&smallLimit, "small-limit", prsize.DefaultSmallLimit, "Maximum effective lines for small classification")
	cmd.Flags().IntVar(&smallLimit, "max-small", prsize.DefaultSmallLimit, "Alias for --small-limit")
	cmd.Flags().IntVar(&mediumLimit, "medium-limit", prsize.DefaultMediumLimit, "Maximum effective lines for medium classification")
	cmd.Flags().IntVar(&mediumLimit, "max-medium", prsize.DefaultMediumLimit, "Alias for --medium-limit")
	cmd.Flags().BoolVar(&failOnLarge, "fail-on-large", false, "Exit with error code if PR is classified as large")
	cmd.Flags().StringVar(&baseRef, "base", "", "Base branch or ref to diff against (default: origin default branch merge-base)")
	cmd.Flags().StringVar(&baseRef, "rev", "", "Alias for --base")
	cmd.Flags().BoolVar(&ignoreLockfiles, "ignore-lockfiles", true, "Ignore lockfiles (go.sum, package-lock.json, etc.) from diff line count")
	cmd.Flags().BoolVar(&ignoreGenerated, "ignore-generated", true, "Ignore generated changes marked in .gitattributes (linguist-generated)")
	cmd.Flags().BoolVar(&includeUncommitted, "include-uncommitted", false, "Include uncommitted staged and unstaged working-tree changes")

	return cmd
}
