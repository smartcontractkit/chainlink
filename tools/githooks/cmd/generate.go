package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/generate"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func newGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate [files...]",
		Short: "Run targeted code generators (proto, config docs) on changed files",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			repoRoot, repoErr := findRepoRoot(ctx)
			if repoErr != nil {
				return repoErr
			}

			// Diff against the merge-base with the default branch, not staged
			// files or HEAD~1, so changes from earlier branch commits still
			// trigger their generators. Matches the lint command's diff base.
			files := args
			if len(files) == 0 {
				rev := modules.GetMergeBase(ctx, repoRoot)
				var changedErr error
				files, changedErr = modules.GetChangedFilesSince(ctx, repoRoot, rev)
				if changedErr != nil {
					return fmt.Errorf("failed to get changed files: %w", changedErr)
				}
			}

			return generate.Run(ctx, repoRoot, files)
		},
	}

	return generateCmd
}
