package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var iterations int

var surveyCmd = &cobra.Command{
	Use:   "survey [gotestsum flags] [-- go test flags and packages]",
	Short: "Re-run tests multiple times (flake hunting, timing); uses gotestsum",
	Long: `Runs gotestsum in a loop with --jsonfile per iteration, then prints a short summary.
Pass packages and go test flags after --, same as gotestsum.`,
	Example: `  chainlink-test survey --iterations 10 -- --shuffle=1 -timeout=15m ./core/...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if iterations < 1 {
			return fmt.Errorf("--iterations must be >= 1")
		}
		return runner.Survey(cmd.Context(), iterations, args)
	},
}

func init() {
	surveyCmd.Flags().IntVar(&iterations, "iterations", 1, "number of full test runs")
}
