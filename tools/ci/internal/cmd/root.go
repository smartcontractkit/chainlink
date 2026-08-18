package cmd

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root 'ci' command with all subcommands attached.
func NewRootCmd(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "ci",
		Short:         "Chainlink CI automation tool",
		Long:          "Universal CI CLI for test matrix discovery, package sharding, changelog preview generation, and automation tooling.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.AddCommand(newMatrixCmd(stdout))
	rootCmd.AddCommand(newTestshardCmd(stdin, stdout, stderr))
	rootCmd.AddCommand(newChangelogCmd(stdout))
	rootCmd.AddCommand(newGatingCmd(stdout))

	return rootCmd
}

// Execute runs the root command bound to os.Stdin, os.Stdout, and os.Stderr.
func Execute(ctx context.Context) error {
	cmd := NewRootCmd(os.Stdin, os.Stdout, os.Stderr)
	return cmd.ExecuteContext(ctx)
}
