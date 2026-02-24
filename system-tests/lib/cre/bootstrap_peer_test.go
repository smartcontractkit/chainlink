package cre

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestResolveP2PAnnounceAddresses_LocalOnly_UsesInternalHost(t *testing.T) {
	addresses, err := ResolveP2PAnnounceAddresses("local", false, "workflow-node0", 5001)
	if err != nil {
		t.Fatalf("ResolveP2PAnnounceAddresses returned error: %v", err)
	}
	if len(addresses) != 1 {
		t.Fatalf("expected a single announce address, got %d (%v)", len(addresses), addresses)
	}
	if addresses[0] != "workflow-node0:5001" {
		t.Fatalf("expected workflow-node0:5001, got %s", addresses[0])
	}
}

func TestResolveP2PAnnounceAddresses_LocalMixed_AddsBridgedHost(t *testing.T) {
	addresses, err := ResolveP2PAnnounceAddresses("local", true, "bootstrap-gateway-node0", 5001)
	if err != nil {
		t.Fatalf("ResolveP2PAnnounceAddresses returned error: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected two announce addresses for mixed mode, got %d (%v)", len(addresses), addresses)
	}
	if addresses[0] != "bootstrap-gateway-node0:5001" {
		t.Fatalf("expected first address to stay internal, got %s", addresses[0])
	}
	if addresses[1] != "host.docker.internal:5001" {
		t.Fatalf("expected bridged host.docker.internal:5001, got %s", addresses[1])
	}
}

func TestResolveP2PAnnounceAddresses_Remote_AddsDirectHostIP(t *testing.T) {
	prevMode, hadMode := os.LookupEnv(runtimecfg.EnvRemoteAccessMode)
	prevIP, hadIP := os.LookupEnv(runtimecfg.EnvEC2HostIP)
	t.Cleanup(func() {
		if hadMode {
			_ = os.Setenv(runtimecfg.EnvRemoteAccessMode, prevMode)
		} else {
			_ = os.Unsetenv(runtimecfg.EnvRemoteAccessMode)
		}
		if hadIP {
			_ = os.Setenv(runtimecfg.EnvEC2HostIP, prevIP)
		} else {
			_ = os.Unsetenv(runtimecfg.EnvEC2HostIP)
		}
	})
	_ = os.Setenv(runtimecfg.EnvRemoteAccessMode, runtimecfg.RemoteAccessModeDirect)
	_ = os.Setenv(runtimecfg.EnvEC2HostIP, "10.1.2.3")

	addresses, err := ResolveP2PAnnounceAddresses("remote", true, "workflow-node0", 5001)
	if err != nil {
		t.Fatalf("ResolveP2PAnnounceAddresses returned error: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected two announce addresses for remote node, got %d (%v)", len(addresses), addresses)
	}
	if addresses[0] != "workflow-node0:5001" {
		t.Fatalf("expected first address to stay internal, got %s", addresses[0])
	}
	if addresses[1] != "10.1.2.3:5001" {
		t.Fatalf("expected external EC2 address 10.1.2.3:5001, got %s", addresses[1])
	}
}

func TestResolveBootstrapPeerURL_RemoteCallerToLocalBootstrap_UsesBridgedHost(t *testing.T) {
	peerURL, err := ResolveBootstrapPeerURL("remote", "local", "p2p_testPeer", "bootstrap-gateway-node0", 5001)
	if err != nil {
		t.Fatalf("ResolveBootstrapPeerURL returned error: %v", err)
	}
	if peerURL != "testPeer@host.docker.internal:5001" {
		t.Fatalf("expected bridged bootstrap peer URL, got %s", peerURL)
	}
}
