package environment

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
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

func TestNewStartComponentClientEC2Mode(t *testing.T) {
	t.Setenv(envAgentMode, "ec2")
	t.Setenv(envLocalAgentURL, "")
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(envEC2InstanceID, "")

	if _, err := newStartComponentClient(zerolog.Nop(), &fakeTunnelManager{}); err == nil {
		t.Fatalf("expected ec2 mode without %s or %s to fail", envEC2AgentURL, envEC2InstanceID)
	}

	t.Setenv(envEC2AgentURL, "http://127.0.0.1:18080") // manual tunnel override
	client, err := newStartComponentClient(zerolog.Nop(), &fakeTunnelManager{})
	if err != nil {
		t.Fatalf("expected ec2 mode client to be created, got %v", err)
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
}

func TestResolveEC2AgentBaseURLRequiresInstanceIDWhenURLMissing(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(envEC2InstanceID, "")
	t.Setenv(envEC2AgentPort, "")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop(), &fakeTunnelManager{})
	if err == nil {
		t.Fatalf("expected missing %s to fail when %s is not set", envEC2InstanceID, envEC2AgentURL)
	}
}

func TestResolveEC2AgentBaseURLRejectsInvalidPort(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(envEC2InstanceID, "i-123")
	t.Setenv(envEC2AgentPort, "not-a-port")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop(), &fakeTunnelManager{})
	if err == nil {
		t.Fatalf("expected invalid %s to fail", envEC2AgentPort)
	}
	if !strings.Contains(err.Error(), envEC2AgentPort) {
		t.Fatalf("expected error to mention %s, got: %v", envEC2AgentPort, err)
	}
}

func TestNewStartComponentClientLocalMode(t *testing.T) {
	t.Setenv(envAgentMode, "")
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(envLocalAgentURL, "")

	if _, err := newStartComponentClient(zerolog.Nop(), &fakeTunnelManager{}); err == nil {
		t.Fatalf("expected local mode without %s to fail", envLocalAgentURL)
	}

	t.Setenv(envLocalAgentURL, "http://127.0.0.1:8080")
	client, err := newStartComponentClient(zerolog.Nop(), &fakeTunnelManager{})
	if err != nil {
		t.Fatalf("expected local mode client to be created, got %v", err)
	}

	httpClient, ok := client.(*httpComponentClient)
	if !ok {
		t.Fatalf("expected httpComponentClient, got %T", client)
	}
	if httpClient.checkHealth {
		t.Fatalf("expected local client health checks to be disabled")
	}
	if httpClient.maxAttempts != 1 {
		t.Fatalf("expected local client retries to be disabled")
	}

	if os.Getenv(envLocalAgentURL) == "" {
		t.Fatalf("expected local agent url to remain set")
	}
}

type fakeTunnelManager struct {
	startCalls int
}

func (f *fakeTunnelManager) Start(_ context.Context, refs []tunnel.EndpointRef) ([]tunnel.TunnelBinding, error) {
	f.startCalls++
	bindings := make([]tunnel.TunnelBinding, 0, len(refs))
	for i, ref := range refs {
		bindings = append(bindings, tunnel.TunnelBinding{
			EndpointRef: ref,
			LocalPort:   19000 + i,
			LocalURL: map[string]string{
				"http": "http://127.0.0.1:19000",
				"ws":   "ws://127.0.0.1:19001",
			}[ref.Scheme],
		})
	}
	return bindings, nil
}

func (f *fakeTunnelManager) Stop(_ context.Context) error { return nil }
func (f *fakeTunnelManager) IsStarted() bool { return f.startCalls > 0 }
func (f *fakeTunnelManager) Snapshot() []tunnel.TunnelBinding { return []tunnel.TunnelBinding{} }

func TestRewriteRemoteBlockchainOutputForLocalAccess(t *testing.T) {
	out := &blockchain.Output{
		Nodes: []*blockchain.Node{
			{
				ExternalHTTPUrl: "http://10.0.0.10:8545",
				ExternalWSUrl:   "ws://10.0.0.10:8546",
			},
		},
	}
	manager := &fakeTunnelManager{}

	if err := rewriteRemoteBlockchainOutputForLocalAccess(
		context.Background(),
		zerolog.Nop(),
		manager,
		0,
		&blockchain.Input{Type: blockchain.TypeAnvil},
		out,
		true,
	); err != nil {
		t.Fatalf("expected rewrite helper to succeed: %v", err)
	}

	if manager.startCalls != 1 {
		t.Fatalf("expected tunnel manager start to be called once, got %d", manager.startCalls)
	}
	if out.Nodes[0].ExternalHTTPUrl != "http://127.0.0.1:19000" {
		t.Fatalf("unexpected rewritten http url: %s", out.Nodes[0].ExternalHTTPUrl)
	}
	if out.Nodes[0].ExternalWSUrl != "ws://127.0.0.1:19001" {
		t.Fatalf("unexpected rewritten ws url: %s", out.Nodes[0].ExternalWSUrl)
	}
	if out.Nodes[0].InternalHTTPUrl == "" || !strings.Contains(out.Nodes[0].InternalHTTPUrl, ":19000") {
		t.Fatalf("expected internal http url to be rewritten for docker host access, got %s", out.Nodes[0].InternalHTTPUrl)
	}
	if out.Nodes[0].InternalWSUrl == "" || !strings.Contains(out.Nodes[0].InternalWSUrl, ":19001") {
		t.Fatalf("expected internal ws url to be rewritten for docker host access, got %s", out.Nodes[0].InternalWSUrl)
	}
}

func TestNewEC2TunnelManagerReturnsNoopWhenNotApplicable(t *testing.T) {
	t.Setenv(envAgentMode, "")
	t.Setenv(envEC2InstanceID, "")
	manager, err := newEC2TunnelManager(zerolog.Nop())
	if err != nil {
		t.Fatalf("expected noop manager for local mode, got error: %v", err)
	}
	if manager.IsStarted() {
		t.Fatalf("expected noop manager to report not started")
	}

	t.Setenv(envAgentMode, "ec2")
	t.Setenv(envEC2InstanceID, "")
	manager, err = newEC2TunnelManager(zerolog.Nop())
	if err != nil {
		t.Fatalf("expected noop manager for ec2 mode without instance, got error: %v", err)
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
