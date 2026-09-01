package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/ui"
)

var noColor bool

// NewRootCmd creates and returns a new root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "githooks",
		Short:         "githooks manages monorepo Git hook operations",
		Long:          "githooks provides tooling for Git hooks (such as Lefthook) to operate efficiently across monorepo Go modules.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color and styled terminal output")

	rootCmd.AddCommand(newLintCmd())
	rootCmd.AddCommand(newTestCmd())
	rootCmd.AddCommand(newTidyCmd())
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newEOFCmd())
	rootCmd.AddCommand(newWhitespaceCmd())
	rootCmd.AddCommand(newPRSizeCmd())
	rootCmd.AddCommand(newSampleOutputsCmd())
	return rootCmd
}

func uiForCmd(cmd *cobra.Command) *ui.UI {
	if noColor {
		return ui.NewWithMode(cmd.OutOrStdout(), ui.ColorNever)
	}
	return ui.New(cmd.OutOrStdout())
}

// SilentError wraps an error whose details have already been rendered to the user.
// When returned to main.go, it exits with code 1 without re-printing "Error: ...".
type SilentError struct {
	Err error
}

func (e SilentError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e SilentError) Unwrap() error {
	return e.Err
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
