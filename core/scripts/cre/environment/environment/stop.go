package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
)

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "stop",
		Short:            "Stops local environment",
		Long:             `Stops local CRE resources only (containers, tracked local tunnels, and local state file).`,
		Example:          "go run . env stop",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := stopLocalResources(cmd.Context(), relativePathToRepoRoot, false, false); err != nil {
				return err
			}
			remoteConfiguredSummary, _ := loadRemoteStopTargets(relativePathToRepoRoot)
			if remoteConfiguredSummary.Total > 0 {
				framework.L.Warn().
					Int("count", remoteConfiguredSummary.Total).
					Msgf("Remote components are still running. Use `env remote stop` to stop them. Remote stop state: %s", remoteStateFileAbsPath(relativePathToRepoRoot))
			}
			fmt.Println("Local environment stopped successfully")
			return nil
		},
	}
	return cmd
}

func stopAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "stop-all",
		Short:            "Stops local and remote resources",
		Long:             `Stops remote CRE components (when configured), then stops local CRE resources and extra local services (beholder, billing, observability), and removes local state directory.`,
		Example:          "go run . env stop-all",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteConfiguredSummary, targets := loadRemoteStopTargets(relativePathToRepoRoot)
			if remoteConfiguredSummary.Total > 0 {
				if err := stopRemoteTargets(cmd.Context(), relativePathToRepoRoot, targets); err != nil {
					return err
				}
			}
			if err := stopLocalResources(cmd.Context(), relativePathToRepoRoot, true, false); err != nil {
				return err
			}
			fmt.Println("All resources stopped successfully")
			return nil
		},
	}
	return cmd
}

func stopRemoteCmd() *cobra.Command {
	var dryRunFlag bool
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops remote components only",
		Long:  `Stops remote CRE components through the agent without performing any local cleanup.`,
		Example: strings.TrimSpace(`
go run . env remote stop
go run . env remote stop --dry-run
`),
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteConfiguredSummary, targets := loadRemoteStopTargets(relativePathToRepoRoot)
			if dryRunFlag {
				framework.L.Info().
					Int("total", remoteConfiguredSummary.Total).
					Int("blockchains", remoteConfiguredSummary.Blockchains).
					Int("nodesets", remoteConfiguredSummary.NodeSets).
					Int("jd", remoteConfiguredSummary.JD).
					Msg("Dry-run: remote components that would be stopped")
				return nil
			}
			if remoteConfiguredSummary.Total == 0 {
				framework.L.Info().Msg("No remote components recorded; nothing to stop.")
				return nil
			}

			if err := stopRemoteTargets(cmd.Context(), relativePathToRepoRoot, targets); err != nil {
				return err
			}
			fmt.Println("Remote components stopped successfully")
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Preview what remote components would be stopped")
	return cmd
}

func loadRemoteStopTargets(relativePathToRepoRoot string) (remoteComponentSummary, *envconfig.Config) {
	var (
		targets *envconfig.Config
		summary remoteComponentSummary
	)
	if envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		cached := &envconfig.Config{}
		statePath := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
		if loadErr := cached.Load(statePath); loadErr != nil {
			framework.L.Warn().Err(loadErr).Msgf("failed to load local CRE state from %s", statePath)
		} else {
			targets = cached
			summary = summarizeRemoteComponents(targets)
		}
	}

	if summary.Total == 0 && remoteStateFileExists(relativePathToRepoRoot) {
		remoteCfg, loadErr := loadRemoteStopConfig(relativePathToRepoRoot)
		if loadErr != nil {
			framework.L.Warn().Err(loadErr).Msgf("failed to load remote component stop state from %s", remoteStateFileAbsPath(relativePathToRepoRoot))
		} else {
			targets = remoteCfg
			summary = summarizeRemoteComponents(targets)
		}
	}
	return summary, targets
}

func stopRemoteTargets(ctx context.Context, relativePathToRepoRoot string, targets *envconfig.Config) error {
	agentState, agentLoadErr := loadRemoteAgentState(relativePathToRepoRoot)
	if agentLoadErr != nil {
		framework.L.Warn().Err(agentLoadErr).Msgf("failed to load remote agent state from %s", remoteStateFileAbsPath(relativePathToRepoRoot))
	} else if agentState != nil {
		applyRemoteAgentEnvFallback(framework.L, agentState)
	}

	summary, stopRemoteErr := remoteclient.StopRemoteComponents(ctx, framework.L, targets)
	framework.L.Info().
		Int("requested", summary.Requested).
		Int("stopped", summary.Stopped).
		Int("missing", summary.Missing).
		Int("failed", summary.Failed).
		Msg("Remote component stop summary")
	if summary.ResidualQueryError != "" {
		framework.L.Warn().Msgf("failed to query remote residual CTF resources: %s", summary.ResidualQueryError)
	} else {
		framework.L.Info().
			Int("containers", len(summary.ResidualContainers)).
			Int("volumes", len(summary.ResidualVolumes)).
			Msg("Remote residual CTF resources after stop")
		if len(summary.ResidualContainers) > 0 {
			framework.L.Warn().Msgf("residual remote CTF containers: %s", strings.Join(summary.ResidualContainers, ", "))
		}
		if len(summary.ResidualVolumes) > 0 {
			framework.L.Warn().Msgf("residual remote CTF volumes: %s", strings.Join(summary.ResidualVolumes, ", "))
		}
	}
	if stopRemoteErr != nil {
		return errors.Wrap(stopRemoteErr, "failed to stop one or more remote components")
	}
	if err := stopRelaySupervisor(relativePathToRepoRoot); err != nil {
		framework.L.Warn().Err(err).Msg("failed to stop relay supervisor after remote stop")
	} else {
		framework.L.Info().Msg("stopped local relay supervisor after remote stop")
	}
	if err := removeRemoteStopConfig(relativePathToRepoRoot); err != nil {
		framework.L.Warn().Err(err).Msg("failed to remove remote component stop state")
	} else {
		framework.L.Info().Msgf("removed remote state directory: %s", filepath.Join(relativePathToRepoRoot, remoteStateDirname))
	}
	if !hasLocalComponents(targets) {
		statePath := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
		if err := os.Remove(statePath); err == nil {
			framework.L.Info().Msgf("removed local CRE state file after remote-only stop: %s", statePath)
		} else if !os.IsNotExist(err) {
			framework.L.Warn().Err(err).Msgf("failed to remove local CRE state file after remote-only stop: %s", statePath)
		}
	}
	return nil
}

func stopLocalResources(ctx context.Context, relativePathToRepoRoot string, removeAllState bool, stopRelay bool) error {
	if stopRelay {
		if err := stopRelaySupervisor(relativePathToRepoRoot); err != nil {
			framework.L.Warn().Err(err).Msg("failed to stop relay supervisor")
		}
	}

	removeErr := framework.RemoveTestContainers()
	if removeErr != nil {
		return errors.Wrap(removeErr, "failed to remove environment containers. Please remove them manually")
	}

	if removeAllState {
		stopBeholderErr := stopBeholder()
		if stopBeholderErr != nil {
			framework.L.Warn().Msgf("failed to stop Beholder: %s", stopBeholderErr)
		}

		stopBillingErr := stopBilling()
		if stopBillingErr != nil {
			framework.L.Warn().Msgf("failed to stop Billing: %s", stopBillingErr)
		}

		stopObsStack := framework.ObservabilityDown()
		if stopObsStack != nil {
			framework.L.Warn().Msgf("failed to stop observability stack: %s", stopObsStack)
		}

		removeCacheErr := envconfig.RemoveAllEnvironmentStateDir(relativePathToRepoRoot)
		if removeCacheErr != nil {
			framework.L.Warn().Msgf("failed to remove local CRE state files: %s", removeCacheErr)
		}
		return nil
	}

	creStateFile := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
	cErr := os.Remove(creStateFile)
	switch {
	case cErr != nil && !os.IsNotExist(cErr):
		framework.L.Warn().Msgf("failed to remove local CRE state file: %s", cErr)
	case cErr != nil && os.IsNotExist(cErr):
		framework.L.Info().Msgf("local CRE state file already absent: %s", creStateFile)
	default:
		framework.L.Info().Msgf("removed local CRE state file: %s", creStateFile)
	}

	runningExtras := runningExtraServiceStopHints(detectServiceStatus(ctx))
	if len(runningExtras) > 0 {
		fmt.Println()
		fmt.Println("The following extra services appear to still be running:")
		for _, hint := range runningExtras {
			fmt.Printf("- %s: stop with `%s`\n", hint.serviceName, hint.stopCommand)
		}
		fmt.Print("\n- All extra services: stop with `go run . env stop --all`\n")
	}

	fmt.Print("\nLocal CRE environment stopped successfully\n")

	return nil
}

type remoteComponentSummary struct {
	Total       int
	Blockchains int
	NodeSets    int
	JD          int
}

func summarizeRemoteComponents(cfg *envconfig.Config) remoteComponentSummary {
	summary := remoteComponentSummary{}
	if cfg == nil {
		return summary
	}
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Placement == envconfig.PlacementRemote {
			summary.Blockchains++
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) == string(envconfig.PlacementRemote) {
			summary.NodeSets++
		}
	}
	if cfg.JD != nil && cfg.JD.Placement == envconfig.PlacementRemote {
		summary.JD = 1
	}
	summary.Total = summary.Blockchains + summary.NodeSets + summary.JD
	return summary
}

func hasLocalComponents(cfg *envconfig.Config) bool {
	if cfg == nil {
		return false
	}
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Placement != envconfig.PlacementRemote {
			return true
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) != string(envconfig.PlacementRemote) {
			return true
		}
	}
	if cfg.JD != nil && cfg.JD.Placement != envconfig.PlacementRemote {
		return true
	}
	return false
}
