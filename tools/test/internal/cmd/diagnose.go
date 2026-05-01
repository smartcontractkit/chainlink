package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose [--diagnose flags] [-- go test flags]",
	Short: "Run /chainlink unit tests multiple times to hunt down flakes, races, timeouts, and more",
	Long: `Runs /chainlink unit tests multiple times to hunt down flakes, races, timeouts, and more.

Pass every flag and package pattern you want forwarded to go test after "--". The harness
prepends "go test -json" (duplicate -json in your arguments is ignored) and adds "-count=1"
when you omit -count or use -count=1. Prefer diagnose --iterations for repetition; you may
use -count>1 to repeat inside one go test (e.g. to reduce DB setup/teardown between diagnose
iterations). With --shuffle-seed, a per-iteration -shuffle=<seed> is appended.`,
	Example: `# Run the full core test suite 10 times.
go -C tools/test run . diagnose --iterations 10 -- ./core/...`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.Load(cmd)
		if err != nil {
			return err
		}
		out := output.NewFromApp(conf)

		if err = validateDiagnoseConfig(conf); err != nil {
			return err
		}

		if err = runner.WarnDiagnoseGoTestCount(out.WarnWriter(), args); err != nil {
			return err
		}

		pool, err := db.EnsurePool(cmd.Context(), conf, out, runner.EffectiveParallelIterations(conf))
		if err != nil {
			return err
		}
		defer func() {
			if err := pool.Cleanup(); err != nil {
				out.Stderrf("error tearing down postgres: %v\n", err)
			}
		}()

		return runner.Diagnose(cmd.Context(), conf, out, args, pool.Resources())
	},
}

func init() {
	diagnoseCmd.Flags().Int("iterations", 1, "number of full test runs")
	diagnoseCmd.Flags().Int("parallel-iterations", 1, "maximum number of diagnose iterations to run concurrently; each worker uses its own ephemeral Postgres")
	diagnoseCmd.Flags().Duration("slow-threshold", 30*time.Second, "tests whose max Elapsed exceeds this are flagged slow")
	diagnoseCmd.Flags().Bool("fail-fast", false, "stop this diagnose run immediately if any iteration fails")
	diagnoseCmd.Flags().StringSlice("fail-fast-on", nil, `stop this diagnose run immediately when an iteration matches one or more categories: "failure", "timeout", "slow", or "any"`)
	diagnoseCmd.Flags().Bool("shuffle-seed", false, "randomize test order each iteration; a unique seed is generated per iteration and recorded in report.json for reproduction")
}

func validateDiagnoseConfig(conf *config.App) error {
	if conf.Iterations < 1 {
		return errors.New("--iterations must be >= 1")
	}
	if conf.ParallelIterations < 1 {
		return errors.New("--parallel-iterations must be >= 1")
	}
	if conf.ParallelIterations > conf.Iterations {
		return errors.New("--parallel-iterations must be <= --iterations")
	}
	if conf.ParallelIterations > 1 && conf.DatabaseURL != "" {
		return errors.New("--parallel-iterations > 1 cannot be used with --database-url")
	}
	failFastOn, err := config.NormalizeFailFastOn(conf.FailFastOn)
	if err != nil {
		return err
	}
	conf.FailFastOn = failFastOn
	return nil
}
