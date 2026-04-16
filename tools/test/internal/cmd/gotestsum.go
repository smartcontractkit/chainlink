package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var gotestsumCmd = &cobra.Command{
	Use:                "gotestsum [gotestsum flags] [-- go test flags]",
	DisableFlagParsing: true,
	Short:              "Run gotestsum; all flags and args are passed through",
	Example:            "  test gotestsum --format=dots -- -count=1 ./core/...",
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.Gotestsum(cmd.Context(), args)
	},
}
