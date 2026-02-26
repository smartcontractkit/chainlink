package environment

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
)

func debugCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "debug",
		Short:            "Debug helpers for remote execution",
		Long:             "Debug helpers for querying remote agent state and logs.",
		PersistentPreRun: globalPreRunFunc,
	}

	cmd.AddCommand(debugStatusCmd())
	cmd.AddCommand(debugLocksCmd())
	cmd.AddCommand(debugLogsCmd())
	return cmd
}

func debugStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get remote agent status snapshot",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withResolvedRemoteRuntime(cmd.Context(), func(ctx context.Context, runtime *remoteclient.Runtime) error {
				status, err := remoteclient.GetAgentStatus(ctx, runtime)
				if err != nil {
					return err
				}
				return printDebugJSON(status)
			})
		},
	}
}

func debugLocksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "locks",
		Short: "Get remote agent lock/in-flight snapshot",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withResolvedRemoteRuntime(cmd.Context(), func(ctx context.Context, runtime *remoteclient.Runtime) error {
				locks, err := remoteclient.GetAgentLocks(ctx, runtime)
				if err != nil {
					return err
				}
				return printDebugJSON(locks)
			})
		},
	}
}

func debugLogsCmd() *cobra.Command {
	var (
		componentKey string
		limit        int
	)
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get bounded agent logs for one component key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			componentKey = strings.TrimSpace(componentKey)
			if componentKey == "" {
				return errors.New("component key is required")
			}
			return withResolvedRemoteRuntime(cmd.Context(), func(ctx context.Context, runtime *remoteclient.Runtime) error {
				logs, err := remoteclient.GetComponentLogs(ctx, runtime, componentKey, limit)
				if err != nil {
					return err
				}
				return printDebugJSON(logs)
			})
		},
	}
	cmd.Flags().StringVar(&componentKey, "component-key", "", "Remote component cache key (for example: nodeset:workflow)")
	cmd.Flags().IntVar(&limit, "limit", 200, "Number of log lines to return")
	_ = cmd.MarkFlagRequired("component-key")
	return cmd
}

func withResolvedRemoteRuntime(ctx context.Context, fn func(context.Context, *remoteclient.Runtime) error) error {
	if state, err := loadRemoteAgentState(relativePathToRepoRoot); err == nil && state != nil {
		applyRemoteAgentEnvFallback(framework.L, state)
	}
	runtime, err := remoteclient.ResolveRuntime(framework.L)
	if err != nil {
		return errors.Wrap(err, "failed to resolve remote runtime (set CRE_REMOTE_AGENT_URL or CRE_REMOTE_AGENT_EC2_INSTANCE_ID/AWS profile)")
	}
	return fn(ctx, runtime)
}

func printDebugJSON(value any) error {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode debug output: %w", err)
	}
	if _, err := fmt.Fprintln(os.Stdout, string(payload)); err != nil {
		return fmt.Errorf("failed to print debug output: %w", err)
	}
	return nil
}
