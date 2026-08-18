package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/gating"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/githuboutput"
)

func newGatingCmd(stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gating",
		Short: "Evaluate integration-tests workflow gates and emit decisions to GITHUB_OUTPUT and GITHUB_STEP_SUMMARY",
		Long: `Reads event context from environment variables (EVENT_NAME, REF_NAME, REF_TYPE, CRE_CHANGES, CCIP_CHANGES,
RUN_E2E_LABEL, SKIP_REGRESSION_LABEL, SKIP_MIXED_ENV_LABEL), evaluates all gating rules, and exports outputs to GITHUB_OUTPUT.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			in := gating.Inputs{
				EventName:           os.Getenv("EVENT_NAME"),
				RefName:             os.Getenv("REF_NAME"),
				RefType:             os.Getenv("REF_TYPE"),
				CREChanges:          os.Getenv("CRE_CHANGES") == "true",
				CCIPChanges:         os.Getenv("CCIP_CHANGES") == "true",
				RunE2ELabel:         os.Getenv("RUN_E2E_LABEL") == "true",
				SkipRegressionLabel: os.Getenv("SKIP_REGRESSION_LABEL") == "true",
				SkipMixedEnvLabel:   os.Getenv("SKIP_MIXED_ENV_LABEL") == "true",
			}

			decisions := gating.Evaluate(in)

			if err := githuboutput.AppendVars(decisions.OutputVars()); err != nil {
				return fmt.Errorf("failed to write gating variables to GITHUB_OUTPUT: %w", err)
			}

			summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
			if summaryFile != "" {
				if err := githuboutput.AppendToFile(summaryFile, decisions.SummaryTable(in)); err != nil {
					return fmt.Errorf("failed to write gating summary to GITHUB_STEP_SUMMARY: %w", err)
				}
			}

			for k, v := range decisions.OutputVars() {
				if _, err := fmt.Fprintf(stdout, "%s=%s\n", k, v); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
