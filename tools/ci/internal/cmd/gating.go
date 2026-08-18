package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/gating"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/githuboutput"
)

func newGatingCmd(stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gating",
		Short: "Compute integration-tests gating decisions and publish them to $GITHUB_OUTPUT",
		Long: "Evaluates the integration-tests gates (CRE smoke, CRE regression, CRE mixed-env, CCIP, " +
			"image builds) from GitHub Actions context env vars, appends the results to $GITHUB_OUTPUT, " +
			"and publishes a summary table to $GITHUB_STEP_SUMMARY.",
		RunE: func(cmd *cobra.Command, args []string) error {
			in := gating.Inputs{
				EventName:           os.Getenv("EVENT_NAME"),
				RefName:             os.Getenv("REF_NAME"),
				RefType:             os.Getenv("REF_TYPE"),
				CREChanges:          envBool("CRE_CHANGES"),
				CCIPChanges:         envBool("CCIP_CHANGES"),
				RunE2ELabel:         envBool("RUN_E2E_LABEL"),
				SkipRegressionLabel: envBool("SKIP_REGRESSION_LABEL"),
				SkipMixedEnvLabel:   envBool("SKIP_MIXED_ENV_LABEL"),
			}
			decisions := gating.Evaluate(in)

			if err := githuboutput.AppendVars(decisions.OutputVars()); err != nil {
				return err
			}

			if stepSummary := os.Getenv("GITHUB_STEP_SUMMARY"); stepSummary != "" {
				if err := githuboutput.AppendToFile(stepSummary, decisions.SummaryTable(in)); err != nil {
					return err
				}
			}

			for key, value := range decisions.OutputVars() {
				if _, err := fmt.Fprintf(stdout, "%s=%s\n", key, value); err != nil {
					return err
				}
			}
			return nil
		},
	}

	return cmd
}

func envBool(key string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(key)))
	return err == nil && value
}
