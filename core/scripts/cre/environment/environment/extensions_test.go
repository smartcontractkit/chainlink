package environment

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const extensionsTestNodeConfig = `
[CRE.WorkflowFetcher]
URL = 'file:///home/chainlink/workflows'

[[Capabilities.WorkflowRegistry.AdditionalSources]]
URL = 'grpc-source:50051'
Name = 'existing'
`

func TestExtensionsNodeConfig_workflowDON(t *testing.T) {
	t.Parallel()

	out, err := extensionsNodeConfig(extensionsTestInput(cre.WorkflowDON), cre.NodeIndexToConfigOverride{0: extensionsTestNodeConfig, 1: extensionsTestNodeConfig})
	require.NoError(t, err)
	require.Len(t, out, 2)

	var c corechainlink.Config
	require.NoError(t, toml.Unmarshal([]byte(out[1]), &c))
	require.Empty(t, *c.CRE.WorkflowFetcher.URL)

	storage := c.Capabilities.WorkflowRegistry.WorkflowStorage
	require.Equal(t, "cre-artifact-api:50052", *storage.URL)
	require.Equal(t, "cre-artifact-store:4566", *storage.ArtifactStorageHost)
	require.False(t, *storage.TLSEnabled)

	sources := c.Capabilities.WorkflowRegistry.AdditionalSourcesConfig
	require.Len(t, sources, 2)
	require.Equal(t, "existing", *sources[0].Name)
	require.Equal(t, "cre-workflow-source:50052", *sources[1].URL)
	require.False(t, *sources[1].TLSEnabled)
}

func TestExtensionsNodeConfig_otherDONsUnchanged(t *testing.T) {
	t.Parallel()

	configs := cre.NodeIndexToConfigOverride{0: extensionsTestNodeConfig}
	out, err := extensionsNodeConfig(extensionsTestInput(cre.GatewayDON), configs)
	require.NoError(t, err)
	require.Equal(t, configs, out)
}

func extensionsTestInput(flag cre.CapabilityFlag) cre.GenerateConfigsInput {
	return cre.GenerateConfigsInput{DonMetadata: &cre.DonMetadata{Name: "don", Flags: []string{flag}}}
}
