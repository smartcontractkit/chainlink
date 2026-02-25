package environment

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/stretchr/testify/require"
)

func TestValidateRemoteBlockchainInput(t *testing.T) {
	err := validateRemoteBlockchainInput(nil)
	require.Error(t, err, "expected nil input to fail validation")

	err = validateRemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeGeth})
	require.Error(t, err, "expected non-anvil input to fail validation")

	err = validateRemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeAnvil})
	require.NoError(t, err, "expected anvil input to pass validation")
}

func TestNewRemoteComponentClientPrefersEC2(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "203.0.113.10")
	t.Setenv(envEC2AgentPort, "18080")

	runtime, err := resolveRemoteRuntime(zerolog.Nop())
	require.NoError(t, err, "expected remote runtime to resolve")
	client, err := newRemoteComponentClient(runtime)
	require.NoError(t, err, "expected ec2-first client to be created")

	httpClient, ok := client.(*httpComponentClient)
	require.True(t, ok, "expected httpComponentClient, got %T", client)
	require.True(t, httpClient.checkHealth, "expected ec2 client to enable health checks")
	require.Equal(t, 3, httpClient.maxAttempts, "expected ec2 client retries to be enabled")
	require.Equal(t, "http://203.0.113.10:18080", httpClient.baseURL, "unexpected ec2 base url")
}

func TestResolveEC2AgentBaseURLRequiresHostOrInstanceInfoWhenURLMissing(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "")
	t.Setenv(runtimecfg.EnvEC2InstanceID, "")
	t.Setenv(envEC2AgentPort, "")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop())
	require.Error(t, err, "expected missing direct host resolution inputs to fail when %s is not set", envEC2AgentURL)
}

func TestResolveEC2AgentBaseURLRejectsInvalidPort(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "203.0.113.10")
	t.Setenv(envEC2AgentPort, "not-a-port")

	_, err := resolveEC2AgentBaseURL(zerolog.Nop())
	require.Error(t, err, "expected invalid %s to fail", envEC2AgentPort)
	require.Contains(t, err.Error(), envEC2AgentPort, "expected error to mention %s", envEC2AgentPort)
}

func TestResolveEC2AgentBaseURLDirectMode(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "203.0.113.10")
	t.Setenv(envEC2AgentPort, "18080")

	baseURL, err := resolveEC2AgentBaseURL(zerolog.Nop())
	require.NoError(t, err, "expected direct mode url resolution to succeed")
	require.Equal(t, "http://203.0.113.10:18080", baseURL, "unexpected direct mode base url")
}

func TestResolveRemoteRuntimeRequiresEC2Resolution(t *testing.T) {
	t.Setenv(envEC2AgentURL, "")
	t.Setenv(runtimecfg.EnvEC2HostIP, "")
	t.Setenv(runtimecfg.EnvEC2InstanceID, "")

	_, err := resolveRemoteRuntime(zerolog.Nop())
	require.Error(t, err, "expected runtime resolution without EC2 inputs to fail")
}

func TestRewriteRemoteBlockchainOutputForDirectAccess(t *testing.T) {
	t.Setenv(runtimecfg.EnvEC2HostIP, "203.0.113.10")
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
	err := rewriteRemoteBlockchainOutputForDirectAccess(out, "203.0.113.10")
	require.NoError(t, err, "expected rewrite helper to succeed")

	require.Equal(t, "http://203.0.113.10:8545", out.Nodes[0].ExternalHTTPUrl, "unexpected rewritten http url")
	require.Equal(t, "ws://203.0.113.10:8546", out.Nodes[0].ExternalWSUrl, "unexpected rewritten ws url")
	require.Equal(t, "http://anvil-1337:8545", out.Nodes[0].InternalHTTPUrl, "internal http url should remain unchanged in direct mode")
	require.Equal(t, "ws://anvil-1337:8546", out.Nodes[0].InternalWSUrl, "internal ws url should remain unchanged in direct mode")
}

func TestRewriteRemoteBlockchainOutputForDirectAccess_NilOutputNoop(t *testing.T) {
	err := rewriteRemoteBlockchainOutputForDirectAccess(nil, "203.0.113.10")
	require.NoError(t, err, "expected nil output rewrite to be a no-op")
}

func TestRewriteRemoteBlockchainOutputForDirectAccess_InvalidExternalURL(t *testing.T) {
	out := &blockchain.Output{
		Nodes: []*blockchain.Node{
			{
				ExternalHTTPUrl: "://bad-url",
				ExternalWSUrl:   "ws://anvil-1337:8546",
			},
		},
	}

	err := rewriteRemoteBlockchainOutputForDirectAccess(out, "203.0.113.10")
	require.Error(t, err, "expected invalid external URL to fail rewrite")
	require.Contains(t, err.Error(), "failed to parse url", "expected parse failure context")
}

func TestRemoteAgentErrorFormatting(t *testing.T) {
	err := remoteAgentError("deployment_failed", "failed to deploy blockchain output")
	want := "remote agent error (deployment_failed): failed to deploy blockchain output"
	require.EqualError(t, err, want, "unexpected remote agent error formatting")
}
