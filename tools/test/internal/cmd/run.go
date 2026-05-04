package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var runCmd = &cobra.Command{
	Use:                "run [go test flags]",
	DisableFlagParsing: true,
	Short:              "Run go test; all flags and args are passed through",
	Long: `Runs go test from the Chainlink repo root (with optional ephemeral Postgres).

Because this subcommand does not parse flags, global options (--database-url,
--postgres-version, --ai-output) must appear on the root command before run, for example:
  go -C tools/test run . --database-url=postgres://... run -v -count=1 ./core/...`,
	Example: `  go -C tools/test run . run -v -count=1 -p 4 ./core/...
  go -C tools/test run . run --postgres-version=16 run -count=1 ./core/...`,
	Args: cobra.ArbitraryArgs,
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
