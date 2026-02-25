package topologyviz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestClassifyTopology_UsesDONTypesAndShardIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dons []DONSummary
		want string
	}{
		{
			name: "single don topology",
			dons: []DONSummary{
				{Name: "workflow", DONTypes: []string{cre.WorkflowDON}},
				{Name: "gateway", DONTypes: []string{cre.GatewayDON}},
			},
			want: ClassSingleDON,
		},
		{
			name: "multi don due to dedicated capabilities don",
			dons: []DONSummary{
				{Name: "workflow", DONTypes: []string{cre.WorkflowDON}},
				{Name: "capabilities", DONTypes: []string{cre.CapabilitiesDON}},
			},
			want: ClassMultiDON,
		},
		{
			name: "sharded due to shard index",
			dons: []DONSummary{
				{Name: "workflow", DONTypes: []string{cre.WorkflowDON}},
				{Name: "shard-1", DONTypes: []string{cre.WorkflowDON}, ShardIndex: 1},
			},
			want: ClassSharded,
		},
		{
			name: "sharded due to shard don type",
			dons: []DONSummary{
				{Name: "workflow", DONTypes: []string{cre.WorkflowDON}},
				{Name: "shard-2", DONTypes: []string{cre.ShardDON}},
			},
			want: ClassSharded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := classifyTopology(tt.dons)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildCapabilityMatrix_RemoteAndChainRendering(t *testing.T) {
	t.Parallel()

	rows := buildCapabilityMatrix([]DONSummary{
		{
			Name: "capabilities",
			Capabilities: []CapabilityPlacement{
				{RawFlag: "read-contract-2337", BaseFlag: "read-contract", RemoteFrom: true, ChainID: ptrUint64(2337)},
				{RawFlag: "read-contract-1337", BaseFlag: "read-contract", RemoteFrom: true, ChainID: ptrUint64(1337)},
				{RawFlag: "read-contract-1337", BaseFlag: "read-contract", RemoteFrom: true, ChainID: ptrUint64(1337)},
			},
		},
		{
			Name: "workflow",
			Capabilities: []CapabilityPlacement{
				{RawFlag: "read-contract-1337", BaseFlag: "read-contract", RemoteFrom: false, ChainID: ptrUint64(1337)},
			},
		},
	})

	readContractRow, ok := rowByCapability(rows, "read-contract")
	require.True(t, ok)
	require.Equal(t, "remote-exposed (1337,2337)", readContractRow.ByDON["capabilities"])
	require.Equal(t, "local (1337)", readContractRow.ByDON["workflow"])
}

func TestRenderASCII_IncludesDONHeadersAndNoHint(t *testing.T) {
	t.Parallel()

	summary := &TopologySummary{
		ConfigRef: "configs/workflow-gateway-capabilities-don.toml",
		Topology:  ClassMultiDON,
		InfraType: "docker",
		DONs: []DONSummary{
			{
				Name:         "workflow",
				DONTypes:     []string{cre.WorkflowDON},
				NodeCount:    4,
				Capabilities: []CapabilityPlacement{{RawFlag: "ocr3", BaseFlag: "ocr3"}},
			},
			{
				Name:         "capabilities",
				DONTypes:     []string{cre.CapabilitiesDON},
				NodeCount:    4,
				Capabilities: []CapabilityPlacement{{RawFlag: "web-api-target", BaseFlag: "web-api-target"}},
			},
		},
	}

	rendered := RenderASCII(summary)
	require.Contains(t, rendered, "DON TOPOLOGY OVERVIEW")
	require.NotContains(t, rendered, "Hint:")
	require.Contains(t, rendered, "workflow DON")
	require.Contains(t, rendered, "capabilities DON")
	require.Contains(t, rendered, "Attributes")
}

func TestBuildSummary_PlacementMatrixShownOnlyForRemoteComponents(t *testing.T) {
	t.Parallel()

	localCfg := &envconfig.Config{
		Blockchains: []*envconfig.Blockchain{
			{Input: blockchain.Input{Type: blockchain.TypeAnvil, ChainID: "1337"}, Placement: envconfig.PlacementLocal},
		},
		JD:       &envconfig.JobDistributor{Input: jd.Input{}, Placement: envconfig.PlacementLocal},
		NodeSets: []*cre.NodeSet{{Input: &ns.Input{Name: "workflow"}, Placement: "local"}},
	}
	localSummary, err := BuildSummary(localCfg, "configs/local.toml")
	require.NoError(t, err)
	require.NotNil(t, localSummary.Placement)
	require.False(t, localSummary.Placement.HasRemote)
	require.NotContains(t, RenderASCIIStartSummary(localSummary), "Runtime Placement Matrix")

	mixedCfg := &envconfig.Config{
		Blockchains: []*envconfig.Blockchain{
			{Input: blockchain.Input{Type: blockchain.TypeAnvil, ChainID: "1337"}, Placement: envconfig.PlacementRemote},
			{Input: blockchain.Input{Type: blockchain.TypeAnvil, ChainID: "2337"}, Placement: envconfig.PlacementLocal},
		},
		JD:       &envconfig.JobDistributor{Input: jd.Input{}, Placement: envconfig.PlacementRemote},
		NodeSets: []*cre.NodeSet{{Input: &ns.Input{Name: "workflow"}, Placement: "local"}, {Input: &ns.Input{Name: "capabilities"}, Placement: "remote"}},
	}
	mixedSummary, err := BuildSummary(mixedCfg, "configs/mixed.toml")
	require.NoError(t, err)
	require.NotNil(t, mixedSummary.Placement)
	require.True(t, mixedSummary.Placement.HasRemote)

	rendered := RenderASCIIStartSummary(mixedSummary)
	require.Contains(t, rendered, "Runtime Placement Matrix")
	require.Contains(t, rendered, "nodeset:capabilities")
	require.Contains(t, rendered, "jd")
}

func TestRenderMarkdown_DropsInferredUsageSections(t *testing.T) {
	t.Parallel()

	summary := &TopologySummary{
		ConfigRef: "configs/workflow-gateway-capabilities-don.toml",
		Topology:  ClassMultiDON,
		InfraType: "docker",
		DONs: []DONSummary{
			{
				Name:         "workflow",
				DONTypes:     []string{cre.WorkflowDON},
				NodeCount:    4,
				Capabilities: []CapabilityPlacement{{RawFlag: "ocr3", BaseFlag: "ocr3"}},
			},
		},
	}

	rendered := RenderMarkdown(summary)
	require.Contains(t, rendered, "## Capability Matrix")
	require.Contains(t, rendered, "## DONs")
	require.NotContains(t, rendered, "Best for")
	require.NotContains(t, rendered, "recommended usage")
	require.NotContains(t, rendered, "Workflow additional sources")
}

func TestWriteArtifacts_WritesAsciiAndMarkdown_RemovesLegacyJSON(t *testing.T) {
	t.Parallel()

	summary := &TopologySummary{
		ConfigRef: "configs/workflow-gateway-don.toml",
		Topology:  ClassSingleDON,
		InfraType: "docker",
		DONs: []DONSummary{
			{
				Name:         "workflow",
				DONTypes:     []string{cre.WorkflowDON},
				NodeCount:    4,
				Capabilities: []CapabilityPlacement{{RawFlag: "ocr3", BaseFlag: "ocr3"}},
			},
		},
	}

	tmpDir := t.TempDir()
	legacyJSONPath := filepath.Join(tmpDir, "topology.json")
	require.NoError(t, os.WriteFile(legacyJSONPath, []byte(`{"legacy":true}`), 0o600))

	artifacts, err := WriteArtifacts(summary, tmpDir)
	require.NoError(t, err)
	require.FileExists(t, artifacts.ASCIIPath)
	require.FileExists(t, artifacts.MarkdownPath)
	_, statErr := os.Stat(legacyJSONPath)
	require.Error(t, statErr)
	require.True(t, os.IsNotExist(statErr))
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

func rowByCapability(rows []capabilityMatrixRow, capability string) (capabilityMatrixRow, bool) {
	for _, row := range rows {
		if row.Capability == capability {
			return row, true
		}
	}
	return capabilityMatrixRow{}, false
}
