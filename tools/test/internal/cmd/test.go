package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var testCmd = &cobra.Command{
	Use:                "test [flags] [packages]",
	DisableFlagParsing: true,
	Short:              "Run go test; all flags and args are passed through",
	Example:            "  go -C ./tools/test run . test -v -count=1 -p 4 ./core/...",
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.Load(cmd)
		if err != nil {
			return err
		}
		defer func() {
			if err := dbHandle.Cleanup(); err != nil {
				fmt.Fprintf(os.Stderr, "error tearing down postgres: %v\n", err)
			}
		}()
		return runner.GoTest(cmd.Context(), conf, args)
	},
}
