package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func newMatrixCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "Commands for generating dynamic GitHub Actions test matrices",
	}

	cmd.AddCommand(newMatrixSystemCmd())
	cmd.AddCommand(newMatrixInMemoryCmd())
	cmd.AddCommand(newMatrixCCIPCmd())
	cmd.AddCommand(newMatrixMixedEnvCmd())

	return cmd
}

func resolveMatrixCommon(cmd *cobra.Command, act *ghaction.Action, runID, runAttempt, spotFlag *string) {
	ctx, _ := act.Context()

	if *runID == "" {
		*runID = act.GetInput("run_id")
	}
	if *runID == "" {
		*runID = act.Getenv("GITHUB_RUN_ID")
	}
	if *runID == "" && ctx != nil && ctx.RunID != 0 {
		*runID = strconv.FormatInt(ctx.RunID, 10)
	}

	if *runAttempt == "" {
		*runAttempt = act.GetInput("run_attempt")
	}
	if *runAttempt == "" {
		*runAttempt = act.Getenv("GITHUB_RUN_ATTEMPT")
	}
	if *runAttempt == "" && ctx != nil && ctx.RunAttempt != 0 {
		*runAttempt = strconv.FormatInt(ctx.RunAttempt, 10)
	}

	if *spotFlag == "" {
		*spotFlag = act.GetInput("spot_flag")
	}
	if *spotFlag == "" {
		*spotFlag = act.Getenv("RUNNER_SPOT_FLAG")
	}
}

func outputMatrix(cmd *cobra.Command, act *ghaction.Action, res any, jsonOutput bool) error {
	jsonData, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal matrix to JSON: %w", err)
	}

	if act.Getenv("GITHUB_OUTPUT") != "" {
		if err := act.SetOutput("matrix", string(jsonData)); err != nil {
			return fmt.Errorf("failed to set GitHub Action output: %w", err)
		}
	}

	if jsonOutput || act.Getenv("GITHUB_OUTPUT") == "" {
		enc := json.NewEncoder(cmd.OutOrStdout())
		if jsonOutput {
			enc.SetIndent("", "  ")
		}
		return enc.Encode(res)
	}

	return nil
}

func newMatrixSystemCmd() *cobra.Command {
	var (
		suite      string
		dir        string
		runID      string
		runAttempt string
		spotFlag   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "system",
		Short: "Generate test matrix for Go system tests (CRE smoke or regression)",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			resolveMatrixCommon(cmd, act, &runID, &runAttempt, &spotFlag)

			if suite == "" {
				suite = act.GetInput("suite")
			}
			if suite == "" {
				suite = "cre-smoke"
			}

			if dir == "" {
				dir = act.GetInput("dir")
			}

			var (
				res any
				err error
			)

			switch suite {
			case "cre-smoke":
				if dir == "" {
					dir = "system-tests/tests/smoke/cre"
				}
				res, err = matrix.BuildCRESmokeMatrix(cmd.Context(), matrix.CRESmokeOptions{
					Dir:        dir,
					RunID:      runID,
					RunAttempt: runAttempt,
					SpotFlag:   spotFlag,
				})
			case "cre-regression":
				if dir == "" {
					dir = "system-tests/tests/regression/cre"
				}
				res, err = matrix.BuildCRERegressionMatrix(cmd.Context(), matrix.CRERegressionOptions{
					Dir:        dir,
					RunID:      runID,
					RunAttempt: runAttempt,
					SpotFlag:   spotFlag,
				})
			default:
				return fmt.Errorf("unknown system test suite: %q (expected 'cre-smoke' or 'cre-regression')", suite)
			}

			if err != nil {
				return err
			}

			return outputMatrix(cmd, act, res, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&suite, "suite", "cre-smoke", "System test suite name ('cre-smoke', 'cre-regression')")
	cmd.Flags().StringVar(&dir, "dir", "", "Directory containing Go system test files")
	cmd.Flags().StringVar(&runID, "run-id", "", "GitHub run ID (env: GITHUB_RUN_ID)")
	cmd.Flags().StringVar(&runAttempt, "run-attempt", "", "GitHub run attempt (env: GITHUB_RUN_ATTEMPT)")
	cmd.Flags().StringVar(&spotFlag, "spot-flag", "", "RunsOn spot flag (e.g. 'spot=co', 'spot=false')")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output formatted JSON to stdout")

	return cmd
}

func newMatrixInMemoryCmd() *cobra.Command {
	var (
		file       string
		runID      string
		runAttempt string
		spotFlag   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "in-memory",
		Short: "Generate test matrix for in-memory integration tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			resolveMatrixCommon(cmd, act, &runID, &runAttempt, &spotFlag)

			if file == "" {
				file = act.GetInput("file")
			}
			if file == "" {
				file = ".github/in-memory-tests.json"
			}

			res, err := matrix.BuildInMemoryMatrix(cmd.Context(), matrix.InMemoryOptions{
				ConfigFile: file,
				RunID:      runID,
				RunAttempt: runAttempt,
				SpotFlag:   spotFlag,
			})
			if err != nil {
				return err
			}

			return outputMatrix(cmd, act, res, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&file, "file", ".github/in-memory-tests.json", "Path to in-memory tests configuration JSON")
	cmd.Flags().StringVar(&runID, "run-id", "", "GitHub run ID (env: GITHUB_RUN_ID)")
	cmd.Flags().StringVar(&runAttempt, "run-attempt", "", "GitHub run attempt (env: GITHUB_RUN_ATTEMPT)")
	cmd.Flags().StringVar(&spotFlag, "spot-flag", "", "RunsOn spot flag (e.g. 'spot=co', 'spot=false')")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output formatted JSON to stdout")

	return cmd
}

func newMatrixCCIPCmd() *cobra.Command {
	var (
		runID      string
		runAttempt string
		spotFlag   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "ccip",
		Short: "Generate test matrix for CCIP system tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			resolveMatrixCommon(cmd, act, &runID, &runAttempt, &spotFlag)

			res, err := matrix.BuildCCIPSystemMatrix(cmd.Context(), matrix.CCIPSystemOptions{
				RunID:      runID,
				RunAttempt: runAttempt,
				SpotFlag:   spotFlag,
			})
			if err != nil {
				return err
			}

			return outputMatrix(cmd, act, res, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&runID, "run-id", "", "GitHub run ID (env: GITHUB_RUN_ID)")
	cmd.Flags().StringVar(&runAttempt, "run-attempt", "", "GitHub run attempt (env: GITHUB_RUN_ATTEMPT)")
	cmd.Flags().StringVar(&spotFlag, "spot-flag", "", "RunsOn spot flag (e.g. 'spot=co', 'spot=false')")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output formatted JSON to stdout")

	return cmd
}

func newMatrixMixedEnvCmd() *cobra.Command {
	var (
		runID      string
		runAttempt string
		spotFlag   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "mixed-env",
		Short: "Generate test matrix for CRE mixed-environment system tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			resolveMatrixCommon(cmd, act, &runID, &runAttempt, &spotFlag)

			res, err := matrix.BuildCREMixedEnvMatrix(cmd.Context(), matrix.CREMixedEnvOptions{
				RunID:      runID,
				RunAttempt: runAttempt,
				SpotFlag:   spotFlag,
			})
			if err != nil {
				return err
			}

			return outputMatrix(cmd, act, res, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&runID, "run-id", "", "GitHub run ID (env: GITHUB_RUN_ID)")
	cmd.Flags().StringVar(&runAttempt, "run-attempt", "", "GitHub run attempt (env: GITHUB_RUN_ATTEMPT)")
	cmd.Flags().StringVar(&spotFlag, "spot-flag", "", "RunsOn spot flag (e.g. 'spot=co', 'spot=false')")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output formatted JSON to stdout")

	return cmd
}
