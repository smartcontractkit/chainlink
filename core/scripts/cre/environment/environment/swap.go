package environment

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	ctypes "github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	"github.com/pkg/errors"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

func swapCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Short: "Swaps parts of the local CRE without restarting the environment",
	}

	cmd.AddCommand(capabilitySwapCmd())
	cmd.AddCommand(nodesSwapCmd())

	return cmd
}

func capabilitySwapCmd() *cobra.Command {
	var (
		capabilityFlag string
		binaryPath     string
	)

	cmd := &cobra.Command{
		Use:     "capability",
		Short:   "Swaps the capability binary of the Chainlink nodes in the environment",
		Long:    "Swaps the capability binary of the Chainlink nodes in the environment. Capability flag is used to find jobs with names containing the capability name, which are cancelled and approved, so that capability binary is reloaded. Only DONs that have the capability are impacted.",
		Aliases: []string{"c", "cap"},
		RunE: func(cmd *cobra.Command, args []string) error {
			swappableapabilities := flags.NewSwappableCapabilityFlagsProvider()
			if !slices.Contains(swappableapabilities.SupportedCapabilityFlags(), capabilityFlag) {
				return fmt.Errorf("capability %s cannot be hot-reloaded. Supported capabilities: %s", capabilityFlag, strings.Join(swappableapabilities.SupportedCapabilityFlags(), ", "))
			}

			content, readErr := os.ReadFile(defaultArtifactsPathFile)
			if readErr != nil {
				return errors.Wrap(readErr, "failed to read artifact paths file. Make sure that local CRE environment is running")
			}

			var paths artifactPaths
			if err := json.Unmarshal(content, &paths); err != nil {
				return errors.Wrap(err, "failed to unmarshal artifact paths file")
			}

			setErr := os.Setenv("CTF_CONFIGS", addCachePrefix(paths.EnvConfig))
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
			}

			config, loadErr := framework.Load[envconfig.Config](nil)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load CTF config")
			}

			var envArtifact creenv.EnvArtifact
			artFile, artErr := os.ReadFile(paths.EnvArtifact)
			if artErr != nil {
				return errors.Wrap(artErr, "failed to read artifact file")
			}
			unmarshalErr := json.Unmarshal(artFile, &envArtifact)
			if unmarshalErr != nil {
				return errors.Wrap(unmarshalErr, "failed to unmarshal artifact file")
			}

			cldLogger := cldlogger.NewSingleFileLogger(nil)

			fullCldEnvOutput, _, loadErr := creenv.BuildFromSavedState(cmd.Context(), cldLogger, config, envArtifact)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load environment")
			}

			// matches only the jobspecs that are related to the capability
			var jobNameContainsCapability = func(jobSpec string) bool {
				r := regexp.MustCompile(`name\s+=\s+"(.*)"`)
				matches := r.FindStringSubmatch(jobSpec)
				if len(matches) < 2 {
					return false
				}

				return strings.Contains(matches[1], capabilityFlag)
			}

			// cancel jobs for nodes that have the capability
			// donId -> nodeId -> proposalIDs
			donIdxToNodeIdToProposalIDs := map[int]map[string][]string{}
			for idx, nodeSet := range fullCldEnvOutput.DonTopology.DonsWithMetadata {
				if !flags.HasFlagForAnyChain(nodeSet.Flags, capabilityFlag) {
					continue
				}
				donIdxToNodeIdToProposalIDs[idx] = map[string][]string{}
				for _, node := range nodeSet.DON.Nodes {
					framework.L.Info().Msgf("Cancelling all job proposals for node %s", node.Name)
					proposalIDs, cancelErr := node.CancelProposals(cmd.Context(), jobNameContainsCapability)
					if cancelErr != nil {
						return errors.Wrapf(cancelErr, "failed to cancel job proposals for node %s", node.Name)
					}
					framework.L.Info().Msgf("Cancelled %d job proposals for node %s", len(proposalIDs), node.Name)
					donIdxToNodeIdToProposalIDs[idx][node.NodeID] = proposalIDs
				}
			}

			if len(donIdxToNodeIdToProposalIDs) == 0 {
				return fmt.Errorf("no nodes found with capability %s in any of the DONs. Please check your topology and make sure that the capability is enabled at least for one DON", capabilityFlag)
			}

			// copy the binary to the Docker containers that have the capability
			for donIdx, nodeIdToProposalIDs := range donIdxToNodeIdToProposalIDs {
				pattern := fullCldEnvOutput.DonTopology.DonsWithMetadata[donIdx].Name + "-node"
				capDir, dirErr := crecapabilities.DefaultContainerDirectory(config.Infra.Type)
				if dirErr != nil {
					return errors.Wrapf(dirErr, "failed to get default capabilities directory for infra type %s", config.Infra.Type)
				}

				copyErr := creworkflow.CopyArtifactToDockerContainers(binaryPath, pattern, capDir)
				if copyErr != nil {
					return errors.Wrapf(copyErr, "failed to copy %s capability binary to Docker containers with pattern %s", binaryPath, pattern)
				}

				// approve the job proposal again, so that job is restarted with the new binary
				for _, node := range fullCldEnvOutput.DonTopology.DonsWithMetadata[donIdx].DON.Nodes {
					proposalIDs, ok := nodeIdToProposalIDs[node.NodeID]
					if ok {
						framework.L.Info().Msgf("Approving %d job proposals for node %s", len(proposalIDs), node.Name)
						approveErr := node.ApproveProposals(cmd.Context(), proposalIDs)
						if approveErr != nil {
							return errors.Wrapf(approveErr, "failed to approve job proposals for node %s", node.Name)
						}
						framework.L.Info().Msgf("Approved %d job proposals for node %s", len(proposalIDs), node.Name)
					}
				}
			}

			return storeArtifacts(config)
		},
	}

	cmd.Flags().StringVarP(&capabilityFlag, "name", "n", "", "Name of the capability to swap")
	cmd.Flags().StringVarP(&binaryPath, "binary", "b", "", "Location of the binary to swap on the host machine")
	cmd.MarkFlagRequired("binary")
	cmd.MarkFlagRequired("name")

	return cmd
}

func nodesSwapCmd() *cobra.Command {
	var (
		forceFlag bool
	)

	cmd := &cobra.Command{
		Use:     "nodes",
		Short:   "Swaps the Docker images of the Chainlink nodes in the environment",
		Long:    "Swaps the Docker images of the Chainlink nodes in the environment. If environment is configured to build the Docker image, it will be rebuilt if any change is detected in the source code.",
		Aliases: []string{"n", "node"},
		RunE: func(cmd *cobra.Command, args []string) error {
			content, readErr := os.ReadFile(defaultArtifactsPathFile)
			if readErr != nil {
				return errors.Wrap(readErr, "failed to read artifact paths file. Make sure that local CRE environment is running")
			}

			var paths artifactPaths
			if err := json.Unmarshal(content, &paths); err != nil {
				return errors.Wrap(err, "failed to unmarshal artifact paths file")
			}

			setErr := os.Setenv("CTF_CONFIGS", addCachePrefix(paths.EnvConfig))
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
			}

			config, loadErr := framework.Load[envconfig.Config](nil)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load CTF config")
			}

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			for _, nodeSet := range config.NodeSets {
				framework.L.Info().Msgf("Removing Docker containers for DON %s", nodeSet.Name)
				containerIDs, containerIDsErr := findAllDockerContainerIDs(nodeSet.Name + "-node")
				if containerIDsErr != nil {
					return errors.Wrapf(containerIDsErr, "failed to find Docker containers for node set %s", nodeSet.Name)
				}

				for _, id := range containerIDs {
					framework.L.Debug().Msgf("Removing Docker container %s", id)
					dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
					if dockerClientErr != nil {
						return errors.Wrap(dockerClientErr, "failed to create Docker client")
					}

					if !forceFlag {
						stopErr := dockerClient.ContainerStop(context.Background(), id, ctypes.StopOptions{})
						if stopErr != nil {
							return errors.Wrapf(stopErr, "failed to stop Docker container %s", id)
						}
					}

					removeErr := dockerClient.ContainerRemove(context.Background(), id, ctypes.RemoveOptions{Force: forceFlag})
					if removeErr != nil {
						return errors.Wrapf(removeErr, "failed to remove Docker container %s", id)
					}
				}

				framework.L.Info().Msgf("Starting new Docker containers for DON %s", nodeSet.Name)
				nodeSet.Out = nil
				var nodesetErr error
				nodeSet.Out, nodesetErr = ns.NewSharedDBNodeSet(nodeSet.Input, config.Blockchains[0].Out)
				if nodesetErr != nil {
					time.Sleep(2 * time.Minute)
					return errors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSet.Name)
				}
			}

			return storeArtifacts(config)
		},
	}

	cmd.Flags().BoolVarP(&forceFlag, "force", "f", true, "Force removal of Docker containers. Set to false to enable graceful shutdown of the containers (be mindful that it will take longer to remove the them)")

	return cmd
}

func findAllDockerContainerIDs(pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}

	containers, containersErr := dockerClient.ContainerList(context.Background(), ctypes.ListOptions{})
	if containersErr != nil {
		return nil, errors.Wrap(containersErr, "failed to list Docker containers")
	}

	containerIDs := []string{}
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, pattern) {
				containerIDs = append(containerIDs, container.ID)
			}
		}
	}

	return containerIDs, nil
}
