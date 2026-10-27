package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const extensionsTestUserOverrides = `
[[Capabilities.WorkflowRegistry.AdditionalSources]]
URL = 'grpc-source:50051'
Name = 'existing'
`

func extensionsTestNodeSet(donType, userOverrides string) *cre.NodeSet {
	return &cre.NodeSet{
		Input:     &ns.Input{Name: donType},
		DONTypes:  []string{donType},
		NodeSpecs: []*cre.NodeSpecWithRole{{Input: &clnode.Input{Node: &clnode.NodeInput{UserConfigOverrides: userOverrides}}}},
	}
}

func TestAddExtensionsNodeOverrides(t *testing.T) {
	t.Parallel()

	workflow := extensionsTestNodeSet(cre.WorkflowDON, extensionsTestUserOverrides)
	gateway := extensionsTestNodeSet(cre.GatewayDON, extensionsTestUserOverrides)
	addExtensionsNodeOverrides([]*cre.NodeSet{workflow, gateway})

	require.Equal(t, extensionsTestUserOverrides, gateway.NodeSpecs[0].Node.UserConfigOverrides)

	var c corechainlink.Config
	require.NoError(t, toml.Unmarshal([]byte(workflow.NodeSpecs[0].Node.UserConfigOverrides), &c))
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

func TestValidateExtensionsNodeOverrides(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		nodeSet   *cre.NodeSet
		wantError string
	}{
		{"additional sources merge", extensionsTestNodeSet(cre.WorkflowDON, extensionsTestUserOverrides), ""},
		{"fetcher", extensionsTestNodeSet(cre.WorkflowDON, "[CRE.WorkflowFetcher]\nURL = 'file:///home/chainlink/workflows'\n"), "[CRE.WorkflowFetcher]"},
		{"storage as dotted key", extensionsTestNodeSet(cre.WorkflowDON, "[Capabilities.WorkflowRegistry]\nWorkflowStorage.URL = 'x'\n"), "[Capabilities.WorkflowRegistry.WorkflowStorage]"},
		{"other DONs", extensionsTestNodeSet(cre.GatewayDON, "[CRE.WorkflowFetcher]\nURL = 'x'\n"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateExtensionsNodeOverrides([]*cre.NodeSet{tt.nodeSet})
			if tt.wantError == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantError)
		})
	}
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
