package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateFile_StoreAndLoad(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "state.toml")

	original := &StateFile{
		Phase: PhaseOnChain,
		Addresses: []AddressRef{
			{ChainSelector: 1337, Address: "0xabc", Type: "CapabilitiesRegistry", Version: "v2"},
			{ChainSelector: 1337, Address: "0xdef", Type: "WorkflowRegistry", Version: "v2"},
		},
		DONIDs:    map[string]uint64{"workflow": 1, "capabilities": 2},
		JDNodeIDs: map[string]string{"node-0": "jd-id-123"},
		WorkflowReg: &WorkflowRegState{
			ChainSelector:  1337,
			AllowedDonIDs:  []uint32{1},
			WorkflowOwners: []string{"0x1111111111111111111111111111111111111111"},
		},
		NodeRuntime: map[string]NodeRuntimeInfo{
			"node-0": {
				PeerID:     "16Uiu2HA...",
				APIURL:     "https://node-0.ns.griddle.sh",
				CSAKey:     "0xabc",
				EVMAddress: map[string]string{"1337": "0x123"},
				NodeType:   "standard",
			},
		},
	}

	require.NoError(t, original.Store(path))

	loaded, err := LoadState(path)
	require.NoError(t, err)
	require.Equal(t, PhaseOnChain, loaded.Phase)
	require.Len(t, loaded.Addresses, 2)
	require.Equal(t, "0xabc", loaded.GetAddress("CapabilitiesRegistry"))
	require.Equal(t, uint64(1), loaded.DONIDs["workflow"])
	require.Equal(t, "jd-id-123", loaded.JDNodeIDs["node-0"])
	require.Equal(t, "16Uiu2HA...", loaded.NodeRuntime["node-0"].PeerID)
	require.Equal(t, "0x123", loaded.NodeRuntime["node-0"].EVMAddress["1337"])
}

func TestStateFile_LoadMissing(t *testing.T) {
	t.Parallel()

	loaded, err := LoadState("/nonexistent/path/state.toml")
	require.NoError(t, err)
	require.Nil(t, loaded)
}

func TestStateFile_SetAddress(t *testing.T) {
	t.Parallel()

	s := &StateFile{}
	s.SetAddress(AddressRef{ChainSelector: 1337, Address: "0xabc", Type: "CapabilitiesRegistry", Version: "v2"})
	require.True(t, s.HasAddress("CapabilitiesRegistry"))
	require.Equal(t, "0xabc", s.GetAddress("CapabilitiesRegistry"))

	// Update existing
	s.SetAddress(AddressRef{ChainSelector: 1337, Address: "0xnew", Type: "CapabilitiesRegistry", Version: "v2"})
	require.Equal(t, "0xnew", s.GetAddress("CapabilitiesRegistry"))
	require.Len(t, s.Addresses, 1)

	// Add different type
	s.SetAddress(AddressRef{ChainSelector: 1337, Address: "0xdef", Type: "WorkflowRegistry", Version: "v2"})
	require.Len(t, s.Addresses, 2)
}

func TestSetAddress_QualifierDistinguishesEntries(t *testing.T) {
	t.Parallel()

	s := &StateFile{}
	s.SetAddress(AddressRef{ChainSelector: 1337, Type: "KeystoneForwarder", Address: "0xA", Qualifier: "don-a"})
	s.SetAddress(AddressRef{ChainSelector: 1337, Type: "KeystoneForwarder", Address: "0xB", Qualifier: "don-b"})
	require.Len(t, s.Addresses, 2, "distinct qualifiers must not collapse into one entry")

	s.SetAddress(AddressRef{ChainSelector: 1337, Type: "KeystoneForwarder", Address: "0xA-updated", Qualifier: "don-a"})
	require.Len(t, s.Addresses, 2, "same type+chain+qualifier must update in place, not append")
	require.Equal(t, "0xA-updated", s.GetAddress("KeystoneForwarder"))
}

func TestStateFile_StoreCreatesDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	nestedDir := filepath.Join(dir, "cre", "state")
	path := filepath.Join(nestedDir, "state.toml")

	s := &StateFile{Phase: PhaseDone}
	require.NoError(t, s.Store(path))

	_, err := os.Stat(path)
	require.NoError(t, err)
}
