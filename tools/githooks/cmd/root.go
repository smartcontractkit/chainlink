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
		Use:           "githooks",
		Short:         "githooks manages monorepo Git hook operations",
		Long:          "githooks provides tooling for Git hooks (such as Lefthook) to operate efficiently across monorepo Go modules.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.AddCommand(newLintCmd())
	rootCmd.AddCommand(newTestCmd())
	rootCmd.AddCommand(newTidyCmd())
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newEOFCmd())
	rootCmd.AddCommand(newWhitespaceCmd())
	return rootCmd
}

// Execute runs the root Cobra command.
func Execute(ctx context.Context) error {
	return NewRootCmd().ExecuteContext(ctx)
}

func findRepoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find repo root via git: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
