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

			files := args
			if len(files) == 0 {
				var stagedErr error
				files, stagedErr = modules.GetStagedFiles(ctx, repoRoot)
				if stagedErr != nil {
					return fmt.Errorf("failed to get staged files: %w", stagedErr)
				}
				if len(files) == 0 {
					var changedErr error
					files, changedErr = modules.GetChangedFiles(ctx, repoRoot)
					if changedErr != nil {
						return fmt.Errorf("failed to get changed files: %w", changedErr)
					}
				}
			}

			return generate.Run(ctx, repoRoot, files)
		},
	}

	return generateCmd
}
