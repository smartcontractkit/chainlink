package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/lint"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func newLintCmd() *cobra.Command {
	var fix bool
	var rev string

	cmd := &cobra.Command{
		Use:   "lint [files...]",
		Short: "Run golangci-lint on changed packages across modules",
		Long: "Discovers enclosing Go modules and packages for files changed since the merge-base with the " +
			"default branch and executes golangci-lint on only those packages, matching the CI diff base.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := findRepoRoot(ctx)
			if err != nil {
				return fmt.Errorf("could not find repo root: %w", err)
			}

			// An empty rev means "diff against the merge-base with the default
			// branch", the same base CI's only-new-issues uses. Using HEAD here
			// instead would hide issues introduced by earlier branch commits.
			if rev == "" {
				rev = modules.GetMergeBase(ctx, repoRoot)
			}

			files := args
			if len(files) == 0 {
				files, err = modules.GetChangedFilesSince(ctx, repoRoot, rev)
				if err != nil {
					return err
				}
			}

			targets, err := modules.FindAffectedModules(repoRoot, files)
			if err != nil {
				return err
			}

			if len(targets) == 0 {
				return nil
			}

			cfg := lint.Config{
				RepoRoot: repoRoot,
				Targets:  targets,
				Fix:      fix,
				Rev:      rev,
				Stdout:   cmd.OutOrStdout(),
				Stderr:   cmd.ErrOrStderr(),
			}

			return lint.Run(ctx, cfg)
		},
	}

	cmd.Flags().BoolVar(&fix, "fix", true, "Fix issues automatically where possible")
	cmd.Flags().StringVar(&rev, "rev", "", "Show issues introduced since rev (default: merge-base with the origin default branch)")

	return cmd
}
