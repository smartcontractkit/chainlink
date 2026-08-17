package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func newMatrixCmd(stdout io.Writer) *cobra.Command {
	var (
		dir          string
		pattern      string
		runID        string
		attempt      string
		runner       string
		githubOutput bool
	)

	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "Discover E2E tests and generate GitHub Actions matrix JSON",
		Example: `  go run ./tools/ci matrix --dir=system-tests/tests/smoke/cre --pattern='^TestCRE_.*_E2E$'
  go run ./tools/ci matrix --dir=system-tests/tests/smoke/cre --run-id=123 --attempt=1 --github-output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			testNames, err := matrix.ScanDir(dir, pattern)
			if err != nil {
				return err
			}

			entries := matrix.BuildMatrix(testNames, runID, attempt, runner)
			return matrix.WriteOutput(stdout, entries, githubOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&dir, "dir", "system-tests/tests/smoke/cre", "Target directory to scan for test files")
	flags.StringVar(&pattern, "pattern", `^TestCRE_.*_E2E$`, "Regex pattern to match test function names")
	flags.StringVar(&runID, "run-id", "0", "GitHub Actions run ID for unique runner labels")
	flags.StringVar(&attempt, "attempt", "1", "GitHub Actions run attempt for unique runner labels")
	flags.StringVar(&runner, "runner", "cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", "Runner hardware and capability spec")
	flags.BoolVar(&githubOutput, "github-output", false, "Append matrix=<json> to file specified in $GITHUB_OUTPUT")

	return cmd
}
