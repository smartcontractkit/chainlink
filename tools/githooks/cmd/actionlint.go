package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/actionlint"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func newActionlintCmd() *cobra.Command {
	var rev string

	cmd := &cobra.Command{
		Use:     "actionlint [files...]",
		Aliases: []string{"al"},
		Short:   "Run actionlint on GitHub Actions workflows when .github YAML files change",
		Long: "Executes actionlint (https://github.com/kjanat/actionlint) on GitHub Actions workflow files " +
			"when any YAML files in .github are changed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := findRepoRoot(ctx)
			if err != nil {
				return fmt.Errorf("could not find repo root: %w", err)
			}

			files := args
			if len(files) == 0 {
				if rev == "" {
					rev = modules.GetMergeBase(ctx, repoRoot)
				}
				files, err = modules.GetChangedFilesSince(ctx, repoRoot, rev)
				if err != nil {
					return err
				}
			}

			cfg := actionlint.Config{
				RepoRoot: repoRoot,
				Files:    files,
				Stdout:   cmd.OutOrStdout(),
				Stderr:   cmd.ErrOrStderr(),
			}

			return actionlint.Run(ctx, cfg)
		},
	}

	cmd.Flags().StringVar(&rev, "rev", "", "Show issues introduced since rev (default: merge-base with the origin default branch)")

	return cmd
}
