package helpers

import (
	"context"
	"testing"
	"time"

	dc "github.com/moby/moby/client"
	"github.com/stretchr/testify/require"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// NodeContainerSnapshot records Docker RestartCount per Chainlink node container.
type NodeContainerSnapshot map[string]int

// SnapshotNodeContainerRestarts captures restart counts after deploy/setup is complete.
func SnapshotNodeContainerRestarts(t *testing.T, testEnv *ttypes.TestEnvironment) NodeContainerSnapshot {
	t.Helper()

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	names := clNodeContainerNames(t, testEnv)
	snap := make(NodeContainerSnapshot, len(names))
	for _, name := range names {
		res, err := client.ContainerInspect(ctx, name, dc.ContainerInspectOptions{})
		require.NoError(t, err, "inspect container %q", name)
		snap[name] = res.Container.RestartCount
	}
	return snap
}

// AssertNodeContainersStable fails if any node container restarted or was OOM-killed since the snapshot.
func AssertNodeContainersStable(t *testing.T, snap NodeContainerSnapshot) {
	t.Helper()

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for name, wantRestarts := range snap {
		res, err := client.ContainerInspect(ctx, name, dc.ContainerInspectOptions{})
		require.NoError(t, err, "inspect container %q", name)
		c := res.Container
		require.NotNil(t, c.State, "container %q has no state", name)
		require.Equal(t, wantRestarts, c.RestartCount, "container %q restart count changed", name)
		require.False(t, c.State.OOMKilled, "container %q was OOM killed", name)
		require.True(t, c.State.Running, "container %q is not running (status=%s)", name, c.State.Status)
	}
}

func newDockerClient(t *testing.T) *dc.Client {
	t.Helper()
	client, err := dc.New(dc.FromEnv)
	require.NoError(t, err, "create docker client")
	t.Cleanup(func() { _ = client.Close() })
	return client
}

// RestartGatewayContainers restarts the gateway node containers. Used by tests that
// need the gateway process restarted (e.g. to reset in-process caches like the vault
// public key cache): jobs persisted in the gateway node's DB — including cresettings
// jobs — re-run on boot, so runtime CRE settings overrides survive the restart.
func RestartGatewayContainers(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	gatewayNames := make(map[string]struct{})
	for _, don := range testEnv.Dons.List() {
		if node, ok := don.Gateway(); ok {
			gatewayNames[node.Name] = struct{}{}
		}
	}
	require.NotEmpty(t, gatewayNames, "no gateway nodes found in the environment")

	names := make([]string, 0, len(gatewayNames))
	for _, name := range clNodeContainerNames(t, testEnv) {
		if _, ok := gatewayNames[name]; ok {
			names = append(names, name)
		}
	}
	require.NotEmpty(t, names, "no gateway node containers found in the environment")

	client := newDockerClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, name := range names {
		_, err := client.ContainerRestart(ctx, name, dc.ContainerRestartOptions{})
		require.NoErrorf(t, err, "failed to restart gateway container %q", name)
	}
}
