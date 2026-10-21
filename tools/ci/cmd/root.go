package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns a new root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "ci",
		Short:         "ci manages CI workflow operations",
		Long:          "ci provides tooling for CI workflows to replace inline bash and provide testable operations.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newToolsCmd())
	rootCmd.AddCommand(newImageCmd())
	rootCmd.AddCommand(newChangesetCmd())
	rootCmd.AddCommand(newCCIPCmd())
	rootCmd.AddCommand(newRunnerCmd())
	rootCmd.AddCommand(newMatrixCmd())

	return rootCmd
}

// Execute runs the root Cobra command.
func Execute(ctx context.Context) error {
	return NewRootCmd().ExecuteContext(ctx)
}

// FindRepoRoot returns the absolute path to the root of the git repository.
func FindRepoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find repo root via git: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
