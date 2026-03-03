package environment

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestValidateRemoteBlockchainInput(t *testing.T) {
	err := validateRemoteBlockchainInput(nil)
	require.Error(t, err, "expected nil input to fail validation")

	err = validateRemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeGeth})
	require.Error(t, err, "expected non-anvil input to fail validation")

	err = validateRemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeAnvil})
	require.NoError(t, err, "expected anvil input to pass validation")
}

func TestNewRemoteComponentClientPrefersResolvedRuntime(t *testing.T) {
	t.Setenv(remoteclient.EnvRemoteAgentURL, "")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")
	t.Setenv(remoteclient.EnvRemoteAgentPort, "18080")

	runtime, err := remoteclient.ResolveRuntime(zerolog.Nop())
	require.NoError(t, err, "expected remote runtime to resolve")
	client, err := remoteclient.NewComponentClient(runtime)
	require.NoError(t, err, "expected runtime-backed client to be created")
	require.NotNil(t, client, "expected component client to be created")
	require.Equal(t, "http://203.0.113.10:18080", runtime.AgentBaseURL, "unexpected remote base url")
}

func TestResolveRemoteAgentBaseURLRequiresHostOrInstanceInfoWhenURLMissing(t *testing.T) {
	t.Setenv(remoteclient.EnvRemoteAgentURL, "")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "")
	t.Setenv(runtimecfg.EnvRemoteAgentEC2InstanceID, "")
	t.Setenv(remoteclient.EnvRemoteAgentPort, "")

	_, err := remoteclient.ResolveRuntime(zerolog.Nop())
	require.Error(t, err, "expected missing direct host resolution inputs to fail when %s is not set", remoteclient.EnvRemoteAgentURL)
}

func TestResolveRemoteAgentBaseURLRejectsInvalidPort(t *testing.T) {
	t.Setenv(remoteclient.EnvRemoteAgentURL, "")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")
	t.Setenv(remoteclient.EnvRemoteAgentPort, "not-a-port")

	_, err := remoteclient.ResolveRuntime(zerolog.Nop())
	require.Error(t, err, "expected invalid %s to fail", remoteclient.EnvRemoteAgentPort)
	require.Contains(t, err.Error(), remoteclient.EnvRemoteAgentPort, "expected error to mention %s", remoteclient.EnvRemoteAgentPort)
}

func TestResolveRemoteAgentBaseURLDirectMode(t *testing.T) {
	t.Setenv(remoteclient.EnvRemoteAgentURL, "")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")
	t.Setenv(remoteclient.EnvRemoteAgentPort, "18080")

	runtime, err := remoteclient.ResolveRuntime(zerolog.Nop())
	require.NoError(t, err, "expected direct mode url resolution to succeed")
	require.Equal(t, "http://203.0.113.10:18080", runtime.AgentBaseURL, "unexpected direct mode base url")
}

func TestResolveRemoteRuntimeRequiresEC2DiscoveryInputsWhenNoURLOrHost(t *testing.T) {
	t.Setenv(remoteclient.EnvRemoteAgentURL, "")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "")
	t.Setenv(runtimecfg.EnvRemoteAgentEC2InstanceID, "")

	_, err := remoteclient.ResolveRuntime(zerolog.Nop())
	require.Error(t, err, "expected runtime resolution without URL/host/EC2 discovery inputs to fail")
}

func TestRewriteRemoteBlockchainOutputForDirectAccess(t *testing.T) {
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")
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
	require.Contains(t, err.Error(), "failed to parse address", "expected parse failure context")
}

func TestRemoteAgentErrorFormatting(t *testing.T) {
	err := remoteclient.RemoteAgentError("deployment_failed", "failed to deploy blockchain output")
	want := "remote agent error (deployment_failed): failed to deploy blockchain output"
	require.EqualError(t, err, want, "unexpected remote agent error formatting")
}
