package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/eof"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func newEOFCmd() *cobra.Command {
	var (
		check   bool
		quiet   bool
		verbose bool
	)

	cmd := &cobra.Command{
		Use:     "end-of-file-fixer [files...]",
		Aliases: []string{"eof", "eof-fixer"},
		Short:   "Ensure files end with a single newline (empty files remain empty)",
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

			cfg := eof.Config{
				CheckOnly: check,
			}

			result, err := eof.Run(ctx, repoRoot, files, cfg)
			if err != nil {
				return err
			}

			if len(result.ModifiedFiles) > 0 {
				if !quiet {
					action := "Fixing end-of-file for"
					if check {
						action = "Missing/excess trailing newlines found in"
					}
					for _, f := range result.ModifiedFiles {
						if verbose || check {
							fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", action, f)
						}
					}
				}
				if check {
					return fmt.Errorf("%d file(s) require end-of-file newline fixes", len(result.ModifiedFiles))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "Check if files need end-of-file fixing without modifying them")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-error output")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")

	return cmd
}
