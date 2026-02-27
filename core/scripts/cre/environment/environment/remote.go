package environment

import (
	"context"

	"github.com/spf13/cobra"

	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
)

func remoteCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "remote",
		Short:            "Remote execution helpers",
		Long:             "Helpers for controlling and inspecting the remote execution agent.",
		PersistentPreRun: globalPreRunFunc,
	}

	cmd.AddCommand(stopRemoteCmd())
	cmd.AddCommand(remoteStatusCmd())
	cmd.AddCommand(remoteDebugCmds())
	return cmd
}

func remoteStatusCmd() *cobra.Command {
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
