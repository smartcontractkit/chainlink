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
	Use:   "survey [flags] <go test package pattern>",
	Short: "Run /chainlink unit tests multiple times to hunt down flakes, races, timeouts, and more",
	Long: `Runs /chainlink unit tests multiple times to hunt down flakes, races, timeouts, and more.

Accepts exactly one positional argument: the same package pattern you would pass
to go test (e.g. ./core/...).`,
	Example: `# Run the full core test suite 10 times, with each iteration timing out after 15 minutes.
go -C ./tools/test run . survey --iterations 10 --timeout=15m ./core/...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.Load(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			//nolint:revive,staticcheck // ... is a valid error message
			return errors.New("specify a go test target, e.g. ./core/...")
		}
		if len(args) > 1 {
			//nolint:revive,staticcheck // ... is a valid error message
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

		return runner.Survey(cmd.Context(), conf, targetDir, dbHandle.Reset, dbHandle.DumpState)
	},
}

func init() {
	surveyCmd.Flags().Int("iterations", 1, "number of full test runs")
	surveyCmd.Flags().Duration("slow-threshold", 30*time.Second, "tests whose max Elapsed exceeds this are flagged slow")
	surveyCmd.Flags().Duration("timeout", 10*time.Minute, "go test -timeout for each iteration")
	surveyCmd.Flags().Bool("fail-fast", false, "fail the survey immediately if any iteration fails")
	surveyCmd.Flags().Bool("race", false, "run tests with -race")
	surveyCmd.Flags().String("run", "", "passed through as `go test -run=<regex>`; narrow the survey to specific tests (e.g. `^TestFoo$` or `TestFoo/bar`)")
	surveyCmd.Flags().String("cpu", "", "passed as `go test -cpu=<list>`; comma-separated GOMAXPROCS values (e.g. `1,2,4`)")
	surveyCmd.Flags().Int("parallel", 0, "passed as `go test -parallel=<n>`; 0 omits the flag")
	surveyCmd.Flags().Bool("shuffle-seed", false, "randomize test order each iteration; a unique seed is generated per iteration and recorded in report.json for reproduction")
}
