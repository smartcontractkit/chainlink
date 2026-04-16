package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var testCmd = &cobra.Command{
	Use:                "test [flags] [packages]",
	DisableFlagParsing: true,
	Short:              "Run go test; all flags and args are passed through",
	Example:            "  go -C ./tools/test run . test -v -count=1 -p 4 ./core/...",
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cleanup, ok := cmd.Context().Value("cleanup").(func() error)
		if ok {
			defer cleanup()
		} else {
			fmt.Println("WARNING: No cleanup function found in context")
		}
		return runner.GoTest(cmd.Context(), args)
	},
}
