package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var surveyCmd = &cobra.Command{
	Use:   "survey [gotestsum flags] [-- go test flags and packages]",
	Short: "Re-run tests multiple times (flake hunting, timing); uses gotestsum",
	Long: `Runs gotestsum in a loop with --jsonfile per iteration, then prints a short summary.
Pass packages and go test flags after --, same as gotestsum.`,
	Example: `# Run the full core test suite 10 times, with each iteration timing out after 15 minutes. Collect statistics, debug logs, and more
go -C ./tools/test run . survey --iterations 10 --timeout=15m ./core/...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.Load(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			//nolint // ... is a valid error message
			return errors.New("specify a go test target, e.g. ./core/...")
		}
		if len(args) > 1 {
			//nolint // ... is a valid error message
			return errors.New("only one go test target can be specified, e.g. ./core/...")
		}
		targetDir := args[0]
		if conf.Iterations < 1 {
			return errors.New("--iterations must be >= 1")
		}

		defer func() {
			if err := dbHandle.Cleanup(); err != nil {
				fmt.Fprintf(os.Stderr, "error tearing down postgres: %v\n", err)
			}
		}()

		return runner.Survey(cmd.Context(), conf, targetDir, dbHandle.Reset)
	},
}

func init() {
	surveyCmd.Flags().Int("iterations", 1, "number of full test runs")
	surveyCmd.Flags().Duration("slow-threshold", 30*time.Second, "tests whose max Elapsed exceeds this are flagged slow")
}
