package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/testrunner"
)

func newTestCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:     "test [files...]",
		Aliases: []string{"only-changed", "short-test"},
		Short:   "Run unit tests only on affected packages",
		Long:    "Discovers affected Go test packages for changed/staged files and executes tools/test with -short on only those packages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := findRepoRoot(ctx)
			if err != nil {
				return fmt.Errorf("could not find repo root: %w", err)
			}

			files := args
			if len(files) == 0 {
				files, err = modules.GetChangedFiles(ctx, repoRoot)
				if err != nil {
					return err
				}
			}

			mods, err := modules.FindTestModules(repoRoot, files)
			if err != nil {
				return err
			}

			if len(mods) == 0 {
				return nil
			}

			cfg := testrunner.Config{
				RepoRoot: repoRoot,
				Modules:  mods,
				Short:    short,
				Stdout:   cmd.OutOrStdout(),
				Stderr:   cmd.ErrOrStderr(),
			}

			return testrunner.Run(ctx, cfg)
		},
	}

	cmd.Flags().BoolVar(&short, "short", true, "Run tests in -short mode")

	return cmd
}
