package onchain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestValidateDiscoveredEVMAddresses_Pass(t *testing.T) {
	t.Parallel()

	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {
			EVMAddress: map[string]string{"1337": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		},
	}
	err := validateDiscoveredEVMAddresses("workflow", []string{"node-0"}, []uint64{1337}, runtime)
	require.NoError(t, err)
}

func TestValidateDiscoveredEVMAddresses_MissingChain(t *testing.T) {
	t.Parallel()

	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {EVMAddress: map[string]string{}},
	}
	err := validateDiscoveredEVMAddresses("workflow", []string{"node-0"}, []uint64{1337}, runtime)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing EVM address for chain 1337")
}

func TestValidateDiscoveredEVMAddresses_EmptyAddress(t *testing.T) {
	t.Parallel()

	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {EVMAddress: map[string]string{"1337": "  "}},
	}
	err := validateDiscoveredEVMAddresses("workflow", []string{"node-0"}, []uint64{1337}, runtime)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing EVM address for chain 1337")
}

func TestHydrateDiscoveredEVMAddresses_SetsPublicAddress(t *testing.T) {
	t.Parallel()

	addr := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	donMeta := &cre.DonMetadata{
		Name: "workflow",
		NodesMetadata: []*cre.NodeMetadata{{
			Keys: &secrets.NodeKeys{EVM: make(map[uint64]*crypto.EVMKey)},
		}},
	}
	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {EVMAddress: map[string]string{"1337": addr}},
	}

	err := hydrateDiscoveredEVMAddresses(donMeta, []string{"node-0"}, []uint64{1337}, runtime)
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress(addr), donMeta.NodesMetadata[0].Keys.EVM[1337].PublicAddress)
}

func TestHydrateDiscoveredEVMAddresses_MissingRuntimeChain(t *testing.T) {
	t.Parallel()

	donMeta := &cre.DonMetadata{
		Name: "workflow",
		NodesMetadata: []*cre.NodeMetadata{{
			Keys: &secrets.NodeKeys{EVM: make(map[uint64]*crypto.EVMKey)},
		}},
	}
	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {EVMAddress: map[string]string{}},
	}

	err := hydrateDiscoveredEVMAddresses(donMeta, []string{"node-0"}, []uint64{1337}, runtime)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing discovered EVM address for chain 1337")
}

func TestHydrateDiscoveredOCR2BundleIDs_SetsBundleIDs(t *testing.T) {
	t.Parallel()

	donMeta := &cre.DonMetadata{
		Name: "workflow",
		NodesMetadata: []*cre.NodeMetadata{{
			Keys: &secrets.NodeKeys{},
		}},
	}
	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {OCR2BundleIDs: map[string]string{"evm": "bundle-id-1"}},
	}

	err := hydrateDiscoveredOCR2BundleIDs(donMeta, []string{"node-0"}, runtime)
	require.NoError(t, err)
	require.Equal(t, "bundle-id-1", donMeta.NodesMetadata[0].Keys.OCR2BundleIDs["evm"])
}

func TestHydrateDiscoveredOCR2BundleIDs_MissingDiscoveredBundle(t *testing.T) {
	t.Parallel()

	donMeta := &cre.DonMetadata{
		Name: "workflow",
		NodesMetadata: []*cre.NodeMetadata{{
			Keys: &secrets.NodeKeys{},
		}},
	}
	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {},
	}

	err := hydrateDiscoveredOCR2BundleIDs(donMeta, []string{"node-0"}, runtime)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no discovered OCR2 key bundles")
}

func TestHydrateDiscoveredHosts_UsesChartNodeNamespace(t *testing.T) {
	t.Parallel()

	donMeta := &cre.DonMetadata{
		Name: "zone-c-bootstrap",
		NodesMetadata: []*cre.NodeMetadata{
			{Host: "zone-c-bootstrap-bt-0"}, // synthetic host from system-tests/lib's NewDonMetadata
		},
	}
	cv := &domain.ChartValues{
		Namespace: "zone-c",
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-bt-1", NodeType: domain.RoleBootstrap, Namespace: "zone-c-bt"},
		},
	}

	err := hydrateDiscoveredHosts(donMeta, []string{"node-bt-1"}, cv)
	require.NoError(t, err)
	require.Equal(t, "node-bt-1.zone-c-bt.svc.cluster.local", donMeta.NodesMetadata[0].Host,
		"must overwrite the synthetic host with the real chart node's own namespace")
}

func TestHydrateDiscoveredHosts_MissingNodeName(t *testing.T) {
	t.Parallel()

	donMeta := &cre.DonMetadata{
		Name:          "workflow",
		NodesMetadata: []*cre.NodeMetadata{{Host: "workflow-0"}},
	}
	cv := &domain.ChartValues{Namespace: "zone-c"}

	err := hydrateDiscoveredHosts(donMeta, nil, cv)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing node name for metadata index 0")
}
