package helpers

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	chainlinkUpgradeImageEnvVar           = "CTF_CHAINLINK_UPGRADE_IMAGE"
	upgradeCopyCapabilityBinariesEnvVar   = "CTF_UPGRADE_COPY_CAPABILITY_BINARIES"
)

type UpgradeHooks struct {
	BeforeUpgrade func()
	AfterUpgrade  func()
}

// RunWithOptionalWorkflowUpgrade wraps test assertions with an optional workflow DON upgrade phase.
// If CTF_CHAINLINK_UPGRADE_IMAGE is unset, it is a no-op.
func RunWithOptionalWorkflowUpgrade(t *testing.T, testEnv *ttypes.TestEnvironment, hooks UpgradeHooks) bool {
	if strings.TrimSpace(os.Getenv(chainlinkUpgradeImageEnvVar)) == "" {
		return false
	}
	if hooks.BeforeUpgrade != nil {
		hooks.BeforeUpgrade()
	}
	upgraded := ApplyWorkflowDONUpgradeIfConfigured(t, testEnv)
	if hooks.AfterUpgrade != nil {
		hooks.AfterUpgrade()
	}
	return upgraded
}

// ApplyWorkflowDONUpgradeIfConfigured updates a subset of workflow DON nodes
// to CTF_CHAINLINK_UPGRADE_IMAGE and restarts the workflow DON once.
func ApplyWorkflowDONUpgradeIfConfigured(t *testing.T, testEnv *ttypes.TestEnvironment) bool {
	upgradeImage := strings.TrimSpace(os.Getenv(chainlinkUpgradeImageEnvVar))
	if upgradeImage == "" {
		return false
	}

	workflowNodeSet, err := workflowNodeSetFromConfig(testEnv.Config.NodeSets)
	require.NoError(t, err, "failed to find workflow DON nodeset")

	upgradeCount := workflowUpgradeCount(workflowNodeSet.Nodes)
	require.Greater(t, upgradeCount, 0, "workflow DON upgrade count must be > 0")

	for idx := range upgradeCount {
		workflowNodeSet.NodeSpecs[idx].Node.Image = upgradeImage
	}

	require.NoError(t, restartNodeSet(t.Context(), workflowNodeSet), "failed to restart workflow DON nodeset")
	return true
}

func workflowNodeSetFromConfig(nodeSets []*cre.NodeSet) (*cre.NodeSet, error) {
	for _, nodeSet := range nodeSets {
		if slices.Contains(nodeSet.DONTypes, cre.WorkflowDON) {
			return nodeSet, nil
		}
	}
	return nil, fmt.Errorf("workflow DON nodeset not found")
}

// For current CI topology workflow DON has 4 nodes, so this upgrades 2 nodes.
// For other even topologies, it upgrades N/2 nodes.
func workflowUpgradeCount(totalNodes int) int {
	if totalNodes <= 0 {
		return 0
	}
	if totalNodes%2 == 0 {
		return totalNodes / 2
	}
	if totalNodes == 1 {
		return 1
	}
	return (totalNodes / 2) - 1
}

func restartNodeSet(ctx context.Context, nodeSet *cre.NodeSet) error {
	containerIDs, err := findAllDockerContainerIDs(ctx, nodeSet.Name+"-node")
	if err != nil {
		return err
	}

	removeGroup := errgroup.Group{}
	for _, containerID := range containerIDs {
		containerID := containerID
		removeGroup.Go(func() error {
			docker, clientErr := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
			if clientErr != nil {
				return clientErr
			}
			return docker.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		})
	}
	if err := removeGroup.Wait(); err != nil {
		return err
	}

	nodeSet.Input.NodeSpecs = nodeSet.ExtractCTFInputs()
	if !shouldCopyCapabilityBinariesForUpgrade() {
		for _, nodeSpec := range nodeSet.Input.NodeSpecs {
			if nodeSpec == nil || nodeSpec.Node == nil {
				continue
			}
			nodeSpec.Node.CapabilitiesBinaryPaths = nil
		}
	}
	nodeSet.Out = nil
	_, err = ns.NewSharedDBNodeSet(nodeSet.Input, nil)
	return err
}

func shouldCopyCapabilityBinariesForUpgrade() bool {
	raw := strings.TrimSpace(os.Getenv(upgradeCopyCapabilityBinariesEnvVar))
	if raw == "" {
		return true
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return true
	}
	return parsed
}

func findAllDockerContainerIDs(ctx context.Context, pattern string) ([]string, error) {
	docker, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	containerIDs := make([]string, 0)
	for _, ctr := range containers {
		for _, name := range ctr.Names {
			if strings.Contains(name, pattern) {
				containerIDs = append(containerIDs, ctr.ID)
			}
		}
	}

	return containerIDs, nil
}

func DrainChannel[T any](ch <-chan T) int {
	drained := 0
	for {
		select {
		case <-ch:
			drained++
		default:
			return drained
		}
	}
}
