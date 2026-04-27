package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var gotestsumCmd = &cobra.Command{
	Use:                "gotestsum [gotestsum flags] [-- go test flags]",
	DisableFlagParsing: true,
	Short:              "Run tests with gotestsum",
	Long: `Runs gotestsum from the Chainlink repo root (with optional ephemeral Postgres).

Because this subcommand does not parse flags, global options (--database-url,
--postgres-version, --ai-output) must appear on the root command before gotestsum, for example:
  go -C tools/test run . --database-url=postgres://... gotestsum --format=dots -- -count=1 ./core/...`,
	Example: `go -C tools/test run . gotestsum --format=dots -- -count=1 ./core/...
go -C tools/test run . --ai-output gotestsum --format=testname -- -count=1 ./core/...`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGotestsum(cmd, args, exec.LookPath, func() error { return dbHandle.Cleanup() })
	},
}

// runGotestsum runs the gotestsum path. lookPath and deferCleanup are injectable for tests.
// deferCleanup must run on every exit path after PersistentPreRunE may have started Postgres.
func runGotestsum(cmd *cobra.Command, args []string, lookPath func(string) (string, error), deferCleanup func() error) error {
	defer func() {
		if err := deferCleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "error tearing down postgres: %v\n", err)
		}
	}()
	if _, err := lookPath("gotestsum"); err != nil {
		return fmt.Errorf("gotestsum not on PATH: install with go install gotest.tools/gotestsum@latest: %w", err)
	}
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}
	return runner.Gotestsum(cmd.Context(), conf, args)
}
