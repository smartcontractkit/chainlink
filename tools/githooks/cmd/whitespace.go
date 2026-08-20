package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/whitespace"
)

func newWhitespaceCmd() *cobra.Command {
	var (
		check   bool
		quiet   bool
		verbose bool
	)

	cmd := &cobra.Command{
		Use:     "whitespace-fixer [files...]",
		Aliases: []string{"whitespace", "ws-fixer", "trailing-whitespace"},
		Short:   "Fix erroneous trailing whitespace in eligible code and text files",
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

			cfg := whitespace.Config{
				CheckOnly: check,
			}

			result, err := whitespace.Run(ctx, repoRoot, files, cfg)
			if err != nil {
				return err
			}

			if len(result.ModifiedFiles) > 0 {
				if !quiet {
					action := "Fixing trailing whitespace in"
					if check {
						action = "Erroneous trailing whitespace found in"
					}
					for _, f := range result.ModifiedFiles {
						if verbose || check {
							fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", action, f)
						}
					}
				}
				if check {
					return fmt.Errorf("%d file(s) require whitespace fixes", len(result.ModifiedFiles))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "Check if files have trailing whitespace without modifying them")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-error output")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")

	return cmd
}
