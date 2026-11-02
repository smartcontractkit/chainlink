package helpers

import (
	"context"
	"testing"
	"time"

	dc "github.com/moby/moby/client"
	"github.com/stretchr/testify/require"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// StopNodesetContainers stops every Chainlink node container of the named
// nodeset. A stopped shard sends nothing at all - no trigger events, no
// execution status updates, no error reports - so it is how tests simulate a
// primary shard dying silently.
func StopNodesetContainers(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) {
	t.Helper()

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, name := range nodesetContainerNames(t, testEnv, nodesetName) {
		_, err := client.ContainerStop(ctx, name, dc.ContainerStopOptions{})
		require.NoErrorf(t, err, "failed to stop container %q of nodeset %q", name, nodesetName)
	}
}

// KillNodesetContainers SIGKILLs every Chainlink node container of the named
// nodeset - the instant form of StopNodesetContainers.
func KillNodesetContainers(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) {
	t.Helper()

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, name := range nodesetContainerNames(t, testEnv, nodesetName) {
		_, err := client.ContainerKill(ctx, name, dc.ContainerKillOptions{})
		require.NoErrorf(t, err, "failed to kill container %q of nodeset %q", name, nodesetName)
	}
}

// StartNodesetContainers starts every Chainlink node container of the named
// nodeset again, standing in for a recovered shard rejoining its DON family.
func StartNodesetContainers(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) {
	t.Helper()

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, name := range nodesetContainerNames(t, testEnv, nodesetName) {
		_, err := client.ContainerStart(ctx, name, dc.ContainerStartOptions{})
		require.NoErrorf(t, err, "failed to start container %q of nodeset %q", name, nodesetName)
	}
}
