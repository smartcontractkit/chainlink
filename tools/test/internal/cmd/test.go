package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var testCmd = &cobra.Command{
	Use:                "test [flags] [packages]",
	DisableFlagParsing: true,
	Short:              "Run go test; all flags and args are passed through",
	Example:            "  chainlink-test test -v -count=1 -p 4 ./core/...",
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cleanup := cmd.Context().Value("cleanup").(func())
		defer cleanup()
		return runner.GoTest(cmd.Context(), args)
	},
}
