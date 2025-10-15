package environment

import (
	"context"
	"fmt"
	"net/http"
	net "net/url"
	"os"
	"slices"
	"strings"
	"time"

	ctypes "github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	jdjob "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	ptypes "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

const relativePathToRepoRoot = "../../../../"

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
		forceFlag      bool
	)

	cmd := &cobra.Command{
		Use:              "capability",
		Short:            "Swaps the capability binary of the Chainlink nodes in the environment",
		Long:             "Swaps the capability binary of the Chainlink nodes in the environment. Capability flag is used to find jobs with names containing the capability name, which are cancelled and approved, so that capability binary is reloaded. Only DONs that have the capability are impacted.",
		Aliases:          []string{"c", "cap"},
		PersistentPreRun: joinPreRunFuncs(globalPreRunFunc, envIsRunningPreRunFunc),
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var swapErr error

			defer func() {
				metaData := map[string]any{}
				metaData["name"] = capabilityFlag
				if swapErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(swapErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track(MetricCapabilitySwap, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track capability swap: %s\n", trackingErr)
				}
			}()

			swapErr = swapCapability(cmd.Context(), capabilityFlag, binaryPath, forceFlag)

			return swapErr
		},
	}

	cmd.Flags().StringVarP(&capabilityFlag, "name", "n", "", "Name of the capability to swap")
	cmd.Flags().StringVarP(&binaryPath, "binary", "b", "", "Location of the binary to swap on the host machine")
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", true, "Force removal of Docker containers. Set to false to enable graceful shutdown of the containers (be mindful that it will take longer to remove the them)")
	_ = cmd.MarkFlagRequired("binary")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func swapCapability(ctx context.Context, capabilityFlag, binaryPath string, forceFlag bool) error {
	swappableapabilities := flags.NewSwappableCapabilityFlagsProvider()
	if !slices.Contains(swappableapabilities.SupportedCapabilityFlags(), capabilityFlag) {
		return fmt.Errorf("capability %s cannot be hot-reloaded. Supported capabilities: %s", capabilityFlag, strings.Join(swappableapabilities.SupportedCapabilityFlags(), ", "))
	}

	setErr := os.Setenv("CTF_CONFIGS", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
	if setErr != nil {
		return errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
	}

	config, loadErr := framework.Load[envconfig.Config](nil)
	if loadErr != nil {
		return errors.Wrap(loadErr, "failed to load CTF config")
	}

	envArtifact, artErr := creenv.ReadEnvArtifact(creenv.MustEnvArtifactAbsPath(relativePathToRepoRoot))
	if artErr != nil {
		return errors.Wrap(artErr, "failed to read environment artifact")
	}

	cldLogger := cldlogger.NewSingleFileLogger(nil)

	creEnvironment, _, loadErr := creenv.BuildFromSavedState(ctx, cldLogger, config, envArtifact)
	if loadErr != nil {
		return errors.Wrap(loadErr, "failed to load environment")
	}

	// cancel jobs for nodes that have the capability
	// donId -> nodeId -> proposalIDs
	donIdxToNodeIDToProposalIDs := map[int]map[string][]string{}
	for idx, don := range creEnvironment.DonTopology.Dons.List() {
		if !flags.HasFlagForAnyChain(don.Flags, capabilityFlag) {
			continue
		}

		donIdxToNodeIDToProposalIDs[idx] = map[string][]string{}
		for _, node := range creEnvironment.DonTopology.Dons.List()[idx].Nodes {
			// get all jobs that have a label named "capability" with value equal to capability name
			jobResp, jobErr := creEnvironment.CldfEnvironment.Offchain.ListJobs(ctx, &jdjob.ListJobsRequest{
				Filter: &jdjob.ListJobsRequest_Filter{
					Selectors: []*ptypes.Selector{{
						Key:   cre.CapabilityLabelKey,
						Op:    *ptypes.SelectorOp_EQ.Enum(),
						Value: &capabilityFlag,
					}},
					NodeIds: []string{node.JobDistributorDetails.NodeID},
				},
			})

			if jobErr != nil {
				return errors.Wrapf(jobErr, "failed to list jobs for node %s", node.Name)
			}

			externalJobIDs := []string{}
			for _, job := range jobResp.Jobs {
				// uuid is equal to external job ID
				externalJobIDs = append(externalJobIDs, job.GetUuid())
			}

			framework.L.Info().Msgf("Cancelling matching job proposals for node %s", node.Name)
			proposalIDs, cancelErr := node.CancelProposalsByExternalJobID(ctx, externalJobIDs)
			if cancelErr != nil {
				return errors.Wrapf(cancelErr, "failed to cancel job proposals for node %s", node.Name)
			}
			framework.L.Info().Msgf("Cancelled %d job proposals for node %s", len(proposalIDs), node.Name)
			donIdxToNodeIDToProposalIDs[idx][node.JobDistributorDetails.NodeID] = proposalIDs
		}
	}

	if len(donIdxToNodeIDToProposalIDs) == 0 {
		return fmt.Errorf("no nodes found with capability %s in any of the DONs. Please check your topology and make sure that the capability is enabled at least for one DON", capabilityFlag)
	}

	// copy the binary to the Docker containers that have the capability
	for donIdx := range donIdxToNodeIDToProposalIDs {
		pattern := ns.NodeNamePrefix(creEnvironment.DonTopology.Dons.List()[donIdx].Name)
		capDir, dirErr := crecapabilities.DefaultContainerDirectory(config.Infra.Type)
		if dirErr != nil {
			return errors.Wrapf(dirErr, "failed to get default capabilities directory for infra type %s", config.Infra.Type)
		}

		copyErr := creworkflow.CopyArtifactsToDockerContainers(capDir, pattern, binaryPath)
		if copyErr != nil {
			return errors.Wrapf(copyErr, "failed to copy %s capability binary to Docker containers with pattern %s", binaryPath, pattern)
		}
	}

	// TODO remove if clean up issues mentioned in https://smartcontract-it.atlassian.net/browse/PRODCRE-802 are fixed
	// and directly approve jobspecs without restarting nodes
	nerrg := errgroup.Group{}
	for _, nodeSet := range config.NodeSets {
		if !flags.HasFlagForAnyChain(nodeSet.ComputedCapabilities, capabilityFlag) {
			continue
		}
		nerrg.Go(func() error {
			framework.L.Info().Msgf("Removing Docker containers for DON %s", nodeSet.Name)
			containerIDs, containerIDsErr := findAllDockerContainerIDs(ctx, nodeSet.Name+"-node")
			if containerIDsErr != nil {
				return errors.Wrapf(containerIDsErr, "failed to find Docker containers for node set %s", nodeSet.Name)
			}

			cerrg := errgroup.Group{}
			for _, id := range containerIDs {
				cerrg.Go(func() error {
					framework.L.Debug().Msgf("Removing Docker container %s", id)
					dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
					if dockerClientErr != nil {
						return errors.Wrap(dockerClientErr, "failed to create Docker client")
					}

					signal := "SIGTERM"
					if forceFlag {
						signal = "SIGKILL"
					}
					return dockerClient.ContainerRestart(ctx, id, ctypes.StopOptions{Signal: signal})
				})
			}

			if err := cerrg.Wait(); err != nil {
				return errors.Wrapf(err, "failed to remove Docker containers")
			}

			// make sure that networking is up after restarting
			errg := errgroup.Group{}
			context, cancel := context.WithTimeout(ctx, 1*time.Minute)
			framework.L.Info().Msgf("Waiting for all nodes to be up")
			for _, node := range nodeSet.Out.CLNodes {
				errg.Go(func() error {
					return waitForURL(context, node.Node.ExternalURL+"/sessions", 100*time.Millisecond)
				})
			}
			if err := errg.Wait(); err != nil {
				cancel()
				return errors.Wrapf(err, "failed to wait for all nodes to be up")
			}
			cancel()

			return nil
		})
	}

	if err := nerrg.Wait(); err != nil {
		return errors.Wrapf(err, "failed to restart nodeSets")
	}

	// connect clients again after restarting
	creEnvironment, _, loadErr = creenv.BuildFromSavedState(ctx, cldLogger, config, envArtifact)
	if loadErr != nil {
		return errors.Wrap(loadErr, "failed to load environment")
	}

	// approve the job proposals again, so that the jobs are restarted with the new binary
	for donIdx, nodeIDToProposalIDs := range donIdxToNodeIDToProposalIDs {
		for _, node := range creEnvironment.DonTopology.Dons.List()[donIdx].Nodes {
			proposalIDs, ok := nodeIDToProposalIDs[node.JobDistributorDetails.NodeID]
			if ok {
				framework.L.Info().Msgf("Approving %d job proposals for node %s", len(proposalIDs), node.Name)
				approveErr := node.ApproveProposals(ctx, proposalIDs)
				if approveErr != nil {
					return errors.Wrapf(approveErr, "failed to approve job proposals for node %s", node.Name)
				}
				framework.L.Info().Msgf("Approved %d job proposals for node %s", len(proposalIDs), node.Name)
			}
		}
	}

	return config.Store(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
}

func waitForURL(ctx context.Context, url string, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	parsed, err := net.Parse(url)
	if err != nil {
		return errors.Wrapf(err, "failed to parse URL %s", url)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for %s failed: %w", url, ctx.Err())
		case <-ticker.C:
			resp, err := client.Do(&http.Request{
				Method: "GET",
				URL:    parsed,
			})
			if err == nil {
				defer resp.Body.Close()    //nolint: revive // we want to defer in the loop
				if resp.StatusCode < 500 { // service responds (even 404 means it's up)
					return nil
				}
			}
		}
	}
}

func nodesSwapCmd() *cobra.Command {
	var (
		forceFlag    bool
		waitTimeFlag time.Duration
	)

	cmd := &cobra.Command{
		Use:              "nodes",
		Short:            "Swaps the Docker images of the Chainlink nodes in the environment",
		Long:             "Swaps the Docker images of the Chainlink nodes in the environment. If environment is configured to build the Docker image, it will be rebuilt if any change is detected in the source code.",
		Aliases:          []string{"n", "node"},
		PersistentPreRun: joinPreRunFuncs(globalPreRunFunc, envIsRunningPreRunFunc),
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var swapErr error

			defer func() {
				metaData := map[string]any{}
				metaData["force"] = forceFlag
				if swapErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(swapErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track(MetricNodeSwap, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track nodes swap: %s\n", trackingErr)
				}
			}()

			swapErr = swapNodes(cmd.Context(), forceFlag, waitTimeFlag)

			return swapErr
		},
	}

	cmd.Flags().BoolVarP(&forceFlag, "force", "f", true, "Force removal of Docker containers. Set to false to enable graceful shutdown of the containers (be mindful that it will take longer to remove the them)")
	cmd.Flags().DurationVarP(&waitTimeFlag, "wait-time", "w", 2*time.Minute, "Time to wait for the containers to be removed")

	return cmd
}

func swapNodes(ctx context.Context, forceFlag bool, waitTime time.Duration) error {
	setErr := os.Setenv("CTF_CONFIGS", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
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

	nerrg := errgroup.Group{}
	for _, nodeSet := range config.NodeSets {
		nerrg.Go(func() error {
			framework.L.Info().Msgf("Removing Docker containers for DON %s", nodeSet.Name)
			containerIDs, containerIDsErr := findAllDockerContainerIDs(ctx, nodeSet.Name+"-node")
			if containerIDsErr != nil {
				return errors.Wrapf(containerIDsErr, "failed to find Docker containers for node set %s", nodeSet.Name)
			}

			cerrg := errgroup.Group{}
			for _, id := range containerIDs {
				cerrg.Go(func() error {
					framework.L.Debug().Msgf("Removing Docker container %s", id)
					dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
					if dockerClientErr != nil {
						return errors.Wrap(dockerClientErr, "failed to create Docker client")
					}

					if !forceFlag {
						stopErr := dockerClient.ContainerStop(ctx, id, ctypes.StopOptions{})
						if stopErr != nil {
							return errors.Wrapf(stopErr, "failed to stop Docker container %s", id)
						}
					}

					return dockerClient.ContainerRemove(ctx, id, ctypes.RemoveOptions{Force: forceFlag})
				})
			}

			if err := cerrg.Wait(); err != nil {
				return errors.Wrapf(err, "failed to remove Docker containers")
			}

			framework.L.Info().Msgf("Starting new Docker containers for DON %s", nodeSet.Name)
			nodeSet.Out = nil
			var nodesetErr error
			nodeSet.Out, nodesetErr = ns.NewSharedDBNodeSet(nodeSet.Input, config.Blockchains[0].Out)
			if nodesetErr != nil {
				framework.L.Error().Msgf("Failed to create node set named %s: %s", nodeSet.Name, nodesetErr)
				framework.L.Info().Msgf("Waiting %s for the containers to be removed", waitTime.String())
				time.Sleep(waitTime)

				return errors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSet.Name)
			}

			return nil
		})
	}

	if err := nerrg.Wait(); err != nil {
		return errors.Wrapf(err, "failed to restart nodeSets")
	}

	return config.Store(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
}

func findAllDockerContainerIDs(ctx context.Context, pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}

	containers, containersErr := dockerClient.ContainerList(ctx, ctypes.ListOptions{})
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

func joinPreRunFuncs(funcs ...func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, f := range funcs {
			f(cmd, args)
		}
	}
}

func envIsRunningPreRunFunc(cmd *cobra.Command, args []string) {
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		framework.L.Fatal().Str("Expected location", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)).Msg("Local CRE state file does not exist. Please start the environment first.")
	}

	if !creenv.EnvArtifactFileExists(relativePathToRepoRoot) {
		framework.L.Fatal().Str("Expected location", creenv.MustEnvArtifactAbsPath(relativePathToRepoRoot)).Msg("Environment artifact file does not exist. Please start the environment first.")
	}
}
