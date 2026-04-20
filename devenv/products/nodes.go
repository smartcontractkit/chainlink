package products

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	ctypes "github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

// NodeConfigTemplate is a template for CL node configuration.
type NodeConfigTemplate struct {
	LinkContractAddress string
	ChainID             string
	FinalityDepth       int
	WSURL               string
	HTTPURL             string
	ForwardersEnabled   bool
}

func RestartNodes(ctx context.Context, nodeSet *ns.Input, bc *blockchain.Input, forceStop bool, waitTime time.Duration) error {
	// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", err)
	}

	nerrg := errgroup.Group{}
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

				if !forceStop {
					stopErr := dockerClient.ContainerStop(ctx, id, ctypes.StopOptions{})
					if stopErr != nil {
						return errors.Wrapf(stopErr, "failed to stop Docker container %s", id)
					}
				}

				return dockerClient.ContainerRemove(ctx, id, ctypes.RemoveOptions{Force: forceStop})
			})
		}

		if err := cerrg.Wait(); err != nil {
			return errors.Wrapf(err, "failed to remove Docker containers")
		}

		framework.L.Info().Msgf("Starting new Docker containers for DON %s", nodeSet.Name)
		nodeSet.Out = nil
		var nodesetErr error
		nodeSet.Out, nodesetErr = ns.NewSharedDBNodeSet(nodeSet, bc.Out)
		if nodesetErr != nil {
			framework.L.Error().Msgf("Failed to create node set named %s: %s", nodeSet.Name, nodesetErr)
			framework.L.Info().Msgf("Waiting %s for the containers to be removed", waitTime.String())
			time.Sleep(waitTime)

			return errors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSet.Name)
		}

		return nil
	})

	return nerrg.Wait()
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
