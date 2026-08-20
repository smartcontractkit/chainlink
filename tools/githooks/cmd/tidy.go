package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/tidy"
)

func newTidyCmd() *cobra.Command {
	tidyCmd := &cobra.Command{
		Use:   "tidy [files...]",
		Short: "Run go mod tidy in parallel on all changed Go modules",
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

			affected, err := modules.FindAffectedModules(repoRoot, files)
			if err != nil {
				return fmt.Errorf("failed to resolve affected modules: %w", err)
			}

			if len(affected) == 0 {
				return nil
			}

			modDirs := make([]string, 0, len(affected))
			for _, m := range affected {
				modDirs = append(modDirs, m.Module)
			}

			return tidy.Run(ctx, repoRoot, modDirs)
		},
	}

	return tidyCmd
}
