package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/runner"
)

func newRunnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Commands for managing runner configurations and spot strategies",
	}

	cmd.AddCommand(newRunnerSpotCmd())

	return cmd
}

func newRunnerSpotCmd() *cobra.Command {
	var (
		eventName       string
		ref             string
		refType         string
		refName         string
		baseRef         string
		headRef         string
		forceOnDemand   bool
		strategy        string
		defaultStrategy string
		jsonOutput      bool
	)

	cmd := &cobra.Command{
		Use:   "spot",
		Short: "Determine the spot setting for runs-on runners",
		Long:  "Evaluates GitHub event type, branch/tag refs, merge queue status, and strategies to output the optimal spot configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			ctx, _ := act.Context()

			if eventName == "" && ctx != nil {
				eventName = ctx.EventName
			}
			if ref == "" && ctx != nil {
				ref = ctx.Ref
			}
			if refType == "" && ctx != nil {
				refType = ctx.RefType
			}
			if refName == "" && ctx != nil {
				refName = ctx.RefName
			}
			if baseRef == "" && ctx != nil {
				baseRef = ctx.BaseRef
			}
			if headRef == "" && ctx != nil {
				headRef = ctx.HeadRef
			}
			if eventName == "" {
				eventName = act.Getenv("GITHUB_EVENT_NAME")
			}
			if ref == "" {
				ref = act.Getenv("GITHUB_REF")
			}
			if refType == "" {
				refType = act.Getenv("GITHUB_REF_TYPE")
			}
			if refName == "" {
				refName = act.Getenv("GITHUB_REF_NAME")
			}
			if baseRef == "" {
				baseRef = act.Getenv("GITHUB_BASE_REF")
			}
			if headRef == "" {
				headRef = act.Getenv("GITHUB_HEAD_REF")
			}
			if !forceOnDemand {
				envForce := act.GetInput("force_on_demand")
				if envForce == "" {
					envForce = act.Getenv("RUNNER_FORCE_ON_DEMAND")
				}
				if envForce == "true" || envForce == "1" {
					forceOnDemand = true
				}
			}
			if strategy == "" {
				strategy = act.GetInput("strategy")
			}
			if strategy == "" {
				strategy = act.Getenv("RUNNER_SPOT_STRATEGY")
			}
			if defaultStrategy == "" {
				defaultStrategy = act.GetInput("default_strategy")
			}
			if defaultStrategy == "" {
				defaultStrategy = act.Getenv("RUNNER_DEFAULT_SPOT_STRATEGY")
			}

			res, err := runner.ResolveSpot(runner.SpotInput{
				EventName:        eventName,
				Ref:              ref,
				RefType:          refType,
				RefName:          refName,
				BaseRef:          baseRef,
				HeadRef:          headRef,
				ForceOnDemand:    forceOnDemand,
				StrategyOverride: runner.SpotStrategy(strategy),
				DefaultStrategy:  runner.SpotStrategy(defaultStrategy),
			})
			if err != nil {
				return err
			}

			if act.Getenv("GITHUB_OUTPUT") != "" {
				outputs := map[string]string{
					"spot":           res.Spot,
					"spot_flag":      res.SpotFlag,
					"spot_param":     res.SpotFlag,
					"enabled":        strconv.FormatBool(res.Enabled),
					"strategy":       string(res.Strategy),
					"reason":         res.Reason,
					"is_release":     strconv.FormatBool(res.IsRelease),
					"is_merge_queue": strconv.FormatBool(res.IsMergeQueue),
				}
				if err := act.SetOutputs(outputs); err != nil {
					return err
				}
			}

			if jsonOutput {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			}

			if act.Getenv("GITHUB_OUTPUT") == "" {
				fmt.Fprintln(cmd.OutOrStdout(), res.SpotFlag)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&eventName, "event", "", "GitHub event name (env: GITHUB_EVENT_NAME)")
	cmd.Flags().StringVar(&ref, "ref", "", "GitHub ref (env: GITHUB_REF)")
	cmd.Flags().StringVar(&refType, "ref-type", "", "GitHub ref type (branch/tag) (env: GITHUB_REF_TYPE)")
	cmd.Flags().StringVar(&refName, "ref-name", "", "GitHub ref name (env: GITHUB_REF_NAME)")
	cmd.Flags().StringVar(&baseRef, "base-ref", "", "Target branch for PR (env: GITHUB_BASE_REF)")
	cmd.Flags().StringVar(&headRef, "head-ref", "", "Source branch for PR (env: GITHUB_HEAD_REF)")
	cmd.Flags().BoolVar(&forceOnDemand, "force-on-demand", false, "Force on-demand / disable spot (env: RUNNER_FORCE_ON_DEMAND)")
	cmd.Flags().StringVar(&strategy, "strategy", "", "Explicit spot strategy override ('co', 'pco', 'lowest-price', 'false') (env: RUNNER_SPOT_STRATEGY)")
	cmd.Flags().StringVar(&defaultStrategy, "default-strategy", "", "Default spot strategy if enabled ('pco', 'co') (env: RUNNER_DEFAULT_SPOT_STRATEGY)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output resolved spot configuration in JSON format")

	return cmd
}
