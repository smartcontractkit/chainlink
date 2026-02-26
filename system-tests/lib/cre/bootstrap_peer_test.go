package cre

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestResolveP2PAnnounceAddresses_LocalOnly_UsesInternalHost(t *testing.T) {
	addresses, err := ResolveP2PAnnounceAddresses("local", false, 15001)
	require.NoError(t, err, "ResolveP2PAnnounceAddresses should not fail")
	require.Len(t, addresses, 0, "expected local-only setup to leave announce addresses unset")
}

func TestResolveP2PAnnounceAddresses_LocalMixed_AddsBridgedHost(t *testing.T) {
	prevIP, hadIP := os.LookupEnv(runtimecfg.EnvRemoteHostIP)
	prevLocalIP, hadLocalIP := os.LookupEnv(runtimecfg.EnvLocalHostIP)
	t.Cleanup(func() {
		if hadIP {
			_ = os.Setenv(runtimecfg.EnvRemoteHostIP, prevIP)
		} else {
			_ = os.Unsetenv(runtimecfg.EnvRemoteHostIP)
		}
		if hadLocalIP {
			_ = os.Setenv(runtimecfg.EnvLocalHostIP, prevLocalIP)
		} else {
			_ = os.Unsetenv(runtimecfg.EnvLocalHostIP)
		}
	})
	_ = os.Setenv(runtimecfg.EnvRemoteHostIP, "10.1.2.3")
	_ = os.Setenv(runtimecfg.EnvLocalHostIP, "192.168.1.10")

	addresses, err := ResolveP2PAnnounceAddresses("local", true, 15002)
	require.NoError(t, err, "ResolveP2PAnnounceAddresses should not fail")
	require.Len(t, addresses, 2, "expected two announce addresses for mixed mode")
	require.Equal(t, "192.168.1.10:15002", addresses[0], "unexpected local host announce address")
	require.Equal(t, "10.1.2.3:15002", addresses[1], "unexpected external EC2 announce address")
}

func TestResolveP2PAnnounceAddresses_Remote_AddsDirectHostIP(t *testing.T) {
	prevIP, hadIP := os.LookupEnv(runtimecfg.EnvRemoteHostIP)
	t.Cleanup(func() {
		if hadIP {
			_ = os.Setenv(runtimecfg.EnvRemoteHostIP, prevIP)
		} else {
			_ = os.Unsetenv(runtimecfg.EnvRemoteHostIP)
		}
	})
	_ = os.Setenv(runtimecfg.EnvRemoteHostIP, "10.1.2.3")

	addresses, err := ResolveP2PAnnounceAddresses("remote", true, 16001)
	require.NoError(t, err, "ResolveP2PAnnounceAddresses should not fail")
	require.Len(t, addresses, 1, "expected one announce address for remote node")
	require.Equal(t, "10.1.2.3:16001", addresses[0], "unexpected external EC2 announce address")
}

func TestResolveBootstrapPeerURL_RemoteCallerToLocalBootstrap_UsesBridgedHost(t *testing.T) {
	peerURL, err := ResolveBootstrapPeerURL("remote", "local", "p2p_testPeer", "bootstrap-gateway-node0", 5001)
	require.NoError(t, err, "ResolveBootstrapPeerURL should not fail")
	require.Equal(t, "testPeer@host.docker.internal:5001", peerURL, "unexpected bridged bootstrap peer URL")
}

func TestResolveBootstrapAddress_Matrix(t *testing.T) {
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	tests := []struct {
		name            string
		callerTarget    string
		bootstrapTarget string
		wantAddress     string
	}{
		{
			name:            "local to local uses internal",
			callerTarget:    "local",
			bootstrapTarget: "local",
			wantAddress:     "bootstrap-gateway-node0:5001",
		},
		{
			name:            "local to remote uses external ec2",
			callerTarget:    "local",
			bootstrapTarget: "remote",
			wantAddress:     "203.0.113.10:5001",
		},
		{
			name:            "remote to local uses bridged host",
			callerTarget:    "remote",
			bootstrapTarget: "local",
			wantAddress:     "host.docker.internal:5001",
		},
		{
			name:            "remote to remote uses internal",
			callerTarget:    "remote",
			bootstrapTarget: "remote",
			wantAddress:     "bootstrap-gateway-node0:5001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := ResolveBootstrapAddress(tt.callerTarget, tt.bootstrapTarget, "bootstrap-gateway-node0", 5001)
			require.NoError(t, err, "ResolveBootstrapAddress should not fail")
			require.Equalf(t, tt.wantAddress, address, "expected ResolveBootstrapAddress() for %s", tt.name)
		})
	}
}

func TestResolveBootstrapPeerURL_RejectsEmptyPeerID(t *testing.T) {
	_, err := ResolveBootstrapPeerURL("local", "local", "", "bootstrap-gateway-node0", 5001)
	require.Error(t, err, "expected empty peer id to fail")
}
