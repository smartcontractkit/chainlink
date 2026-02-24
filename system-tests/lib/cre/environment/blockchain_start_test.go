package environment

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestValidatePhase2ARemoteBlockchainInput(t *testing.T) {
	if err := validatePhase2ARemoteBlockchainInput(nil); err == nil {
		t.Fatalf("expected nil input to fail validation")
	}

	if err := validatePhase2ARemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeGeth}); err == nil {
		t.Fatalf("expected non-anvil input to fail validation")
	}

	if err := validatePhase2ARemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeAnvil}); err != nil {
		t.Fatalf("expected anvil input to pass validation, got %v", err)
	}
}

func TestNewStartComponentClientPrefersEC2(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "10.193.28.183")
	t.Setenv(envEC2AgentPort, "18080")

	client, err := newStartComponentClient(zerolog.Nop(), tunnel.NewNoopManager())
	if err != nil {
		t.Fatalf("expected ec2-first client to be created, got %v", err)
	}

	httpClient, ok := client.(*httpComponentClient)
	if !ok {
		t.Fatalf("expected httpComponentClient, got %T", client)
	}
	if !httpClient.checkHealth {
		t.Fatalf("expected ec2 client to enable health checks")
	}
	if httpClient.maxAttempts != 3 {
		t.Fatalf("expected ec2 client retries to be enabled")
	}
	if httpClient.baseURL != "http://10.193.28.183:18080" {
		t.Fatalf("unexpected ec2 base url: %s", httpClient.baseURL)
	}
}

func TestResolveEC2AgentBaseURLRequiresHostOrInstanceInfoWhenURLMissing(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "")
	t.Setenv(runtimecfg.EnvEC2InstanceID, "")
	t.Setenv(envEC2AgentPort, "")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop())
	if err == nil {
		t.Fatalf("expected missing direct host resolution inputs to fail when %s is not set", envEC2AgentURL)
	}
}

func TestResolveEC2AgentBaseURLRejectsInvalidPort(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "10.193.28.183")
	t.Setenv(envEC2AgentPort, "not-a-port")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop())
	if err == nil {
		t.Fatalf("expected invalid %s to fail", envEC2AgentPort)
	}
	if !strings.Contains(err.Error(), envEC2AgentPort) {
		t.Fatalf("expected error to mention %s, got: %v", envEC2AgentPort, err)
	}
}

func TestResolveEC2AgentBaseURLDirectMode(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "10.193.28.183")
	t.Setenv(envEC2AgentPort, "18080")

	baseURL, err := resolveEC2AgentBaseURL(zerolog.Nop())
	if err != nil {
		t.Fatalf("expected direct mode url resolution to succeed, got %v", err)
	}
	if baseURL != "http://10.193.28.183:18080" {
		t.Fatalf("unexpected direct mode base url: %s", baseURL)
	}
}

func TestNewStartComponentClientRequiresEC2Resolution(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "")
	t.Setenv(runtimecfg.EnvEC2InstanceID, "")

	if _, err := newStartComponentClient(zerolog.Nop(), tunnel.NewNoopManager()); err == nil {
		t.Fatalf("expected client creation without EC2 resolution to fail")
	}
}

func TestRewriteRemoteBlockchainOutputForLocalAccess(t *testing.T) {
	t.Setenv(runtimecfg.EnvEC2HostIP, "10.193.28.183")
	out := &blockchain.Output{
		Nodes: []*blockchain.Node{
			{
				ExternalHTTPUrl: "http://anvil-1337:8545",
				ExternalWSUrl:   "ws://anvil-1337:8546",
				InternalHTTPUrl: "http://anvil-1337:8545",
				InternalWSUrl:   "ws://anvil-1337:8546",
			},
		},
	}
	if err := rewriteRemoteBlockchainOutputForLocalAccess(
		t.Context(),
		zerolog.Nop(),
		tunnel.NewNoopManager(),
		0,
		&blockchain.Input{Type: blockchain.TypeAnvil},
		out,
		true,
	); err != nil {
		t.Fatalf("expected rewrite helper to succeed: %v", err)
	}

	if out.Nodes[0].ExternalHTTPUrl != "http://10.193.28.183:8545" {
		t.Fatalf("unexpected rewritten http url: %s", out.Nodes[0].ExternalHTTPUrl)
	}
	if out.Nodes[0].ExternalWSUrl != "ws://10.193.28.183:8546" {
		t.Fatalf("unexpected rewritten ws url: %s", out.Nodes[0].ExternalWSUrl)
	}
	if out.Nodes[0].InternalHTTPUrl != "http://anvil-1337:8545" {
		t.Fatalf("expected internal http url unchanged in direct mode, got %s", out.Nodes[0].InternalHTTPUrl)
	}
	if out.Nodes[0].InternalWSUrl != "ws://anvil-1337:8546" {
		t.Fatalf("expected internal ws url unchanged in direct mode, got %s", out.Nodes[0].InternalWSUrl)
	}
}

func TestNewEC2TunnelManagerAlwaysReturnsNoop(t *testing.T) {
	manager, err := newEC2TunnelManager(zerolog.Nop())
	if err != nil {
		t.Fatalf("expected noop manager, got error: %v", err)
	}
	if manager.IsStarted() {
		t.Fatalf("expected noop manager to report not started")
	}
}

func TestRemoteAgentErrorFormatting(t *testing.T) {
	err := remoteAgentError("deployment_failed", "failed to deploy blockchain output")
	want := "remote agent error (deployment_failed): failed to deploy blockchain output"
	if err == nil || err.Error() != want {
		t.Fatalf("expected %q, got %v", want, err)
	}
}
