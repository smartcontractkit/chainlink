package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/runner"
)

var gotestsumCmd = &cobra.Command{
	Use:                "gotestsum [gotestsum flags] [-- go test flags]",
	DisableFlagParsing: true,
	Short:              "Run tests with gotestsum",
	Example:            "  go -C ./tools/test run . gotestsum --format=dots -- -count=1 ./core/...",
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("gotestsum"); err != nil {
			return fmt.Errorf("gotestsum not on PATH: install with go install gotest.tools/gotestsum@latest: %w", err)
		}
		cleanup, ok := cmd.Context().Value("cleanup").(func() error)
		if ok {
			defer cleanup()
		} else {
			fmt.Println("WARNING: No cleanup function found in context")
		}
		return runner.Gotestsum(cmd.Context(), args)
	},
}
