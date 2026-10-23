package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
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

func TestValidateExtensions_rejectsBeforeStart(t *testing.T) {
	t.Parallel()

	compose := filepath.Join(t.TempDir(), "compose.yaml")
	require.NoError(t, os.WriteFile(compose, []byte("services: {}\n"), 0o600))
	workflowNodeSet := func(family string) *cre.NodeSet {
		return &cre.NodeSet{DONTypes: []string{cre.WorkflowDON}, DonFamilies: []string{family}}
	}

	tests := []struct {
		name string
		in   *envconfig.Config
		want string
	}{
		{"kubernetes", &envconfig.Config{
			Infra:      &infra.Provider{Type: infra.Kubernetes},
			Extensions: &envconfig.Extensions{StorageCompose: compose},
		}, "docker infra provider"},
		{"missing compose file", &envconfig.Config{
			Infra:      &infra.Provider{Type: infra.Docker},
			Extensions: &envconfig.Extensions{StorageCompose: filepath.Join(t.TempDir(), "missing.yaml")},
		}, "extensions.storage_compose"},
		{"two workflow DON families", &envconfig.Config{
			Infra:      &infra.Provider{Type: infra.Docker},
			NodeSets:   []*cre.NodeSet{workflowNodeSet("zone-a"), workflowNodeSet("zone-b")},
			Extensions: &envconfig.Extensions{StorageCompose: compose},
		}, "exactly one workflow DON family"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := validateExtensions(t.Context(), tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
