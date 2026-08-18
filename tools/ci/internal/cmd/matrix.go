package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func newMatrixCmd(stdout io.Writer) *cobra.Command {
	var (
		suite        string
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
		Example: `  go run ./tools/ci matrix --suite=cre-smoke --run-id=123 --attempt=1 --github-output
  go run ./tools/ci matrix --dir=system-tests/tests/smoke/cre --pattern='^TestCRE_.*_E2E$' --run-id=123 --attempt=1 --github-output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if suite != "" {
				entries, err := matrix.BuildSuiteMatrix(matrix.SuiteType(suite), matrix.SuiteOptions{
					Dir:        dir,
					Pattern:    pattern,
					RunID:      runID,
					RunAttempt: attempt,
					RunnerSpec: runner,
				})
				if err != nil {
					return err
				}
				return matrix.WriteOutput(stdout, entries, githubOutput)
			}

			testNames, err := matrix.ScanDir(dir, pattern)
			if err != nil {
				return err
			}

			entries := matrix.BuildMatrix(testNames, runID, attempt, runner)
			return matrix.WriteOutput(stdout, entries, githubOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&suite, "suite", "", "Named test suite (cre-smoke, cre-regression, cre-mixed-env, ccip)")
	flags.StringVar(&dir, "dir", "system-tests/tests/smoke/cre", "Target directory to scan for test files")
	flags.StringVar(&pattern, "pattern", `^TestCRE_.*_E2E$`, "Regex pattern to match test function names")
	flags.StringVar(&runID, "run-id", "0", "GitHub Actions run ID for unique runner labels")
	flags.StringVar(&attempt, "attempt", "1", "GitHub Actions run attempt for unique runner labels")
	flags.StringVar(&runner, "runner", "cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", "Runner hardware and capability spec")
	flags.BoolVar(&githubOutput, "github-output", false, "Append matrix=<json> to file specified in $GITHUB_OUTPUT")

	cmd.AddCommand(newMatrixSetupCmd(stdout))

	return cmd
}

func newMatrixSetupCmd(stdout io.Writer) *cobra.Command {
	var (
		runID            string
		attempt          string
		creSmoke         bool
		creSmokeDir      string
		creRegression    bool
		creRegressionDir string
		creMixedEnv      bool
		ccip             bool
		githubOutput     bool
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Generate all enabled test matrices for integration-tests workflow setup",
		Example: `  go run ./tools/ci matrix setup --cre=true --ccip=true --run-id=123 --attempt=1 --github-output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			matrices, err := matrix.GenerateSetupMatrices(matrix.SetupOptions{
				RunID:            runID,
				RunAttempt:       attempt,
				CRESmoke:         creSmoke,
				CRESmokeDir:      creSmokeDir,
				CRERegression:    creRegression,
				CRERegressionDir: creRegressionDir,
				CREMixedEnv:      creMixedEnv,
				CCIP:             ccip,
			})
			if err != nil {
				return err
			}
			return matrix.WriteMultiOutput(stdout, matrices, githubOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&runID, "run-id", "0", "GitHub Actions run ID for unique runner labels")
	flags.StringVar(&attempt, "attempt", "1", "GitHub Actions run attempt for unique runner labels")
	flags.BoolVar(&creSmoke, "cre", false, "Generate cre-matrix for CRE smoke tests")
	flags.StringVar(&creSmokeDir, "cre-smoke-dir", "", "Custom dir for CRE smoke tests")
	flags.BoolVar(&creRegression, "cre-regression", false, "Generate cre-regression-matrix for CRE regression tests")
	flags.StringVar(&creRegressionDir, "cre-regression-dir", "", "Custom dir for CRE regression tests")
	flags.BoolVar(&creMixedEnv, "cre-mixed-env", false, "Generate cre-mixed-env-matrix for CRE mixed-env tests")
	flags.BoolVar(&ccip, "ccip", false, "Generate ccip-matrix for CCIP v1.6 tests")
	flags.BoolVar(&githubOutput, "github-output", false, "Append outputs to file specified in $GITHUB_OUTPUT")

	return cmd
}
