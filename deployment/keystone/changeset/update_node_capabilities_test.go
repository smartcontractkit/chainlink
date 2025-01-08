package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestUpdateNodeCapabilities(t *testing.T) {
	t.Parallel()

	var (
		capA = kcr.CapabilitiesRegistryCapability{
			LabelledName: "capA",
			Version:      "0.4.2",
		}
		capB = kcr.CapabilitiesRegistryCapability{
			LabelledName: "capB",
			Version:      "3.16.0",
		}
		caps = []kcr.CapabilitiesRegistryCapability{capA, capB}
	)
	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		// contract set is already deployed with capabilities
		// we have to keep track of the existing capabilities to add to the new ones
		var p2pIDs []p2pkey.PeerID
		newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for id := range te.WFNodes {
			k, err := p2pkey.MakePeerID(id)
			require.NoError(t, err)
			p2pIDs = append(p2pIDs, k)
			newCapabilities[k] = caps
		}

		t.Run("fails if update drops existing capabilities", func(t *testing.T) {
			cfg := changeset.UpdateNodeCapabilitiesRequest{
				RegistryChainSel:  te.RegistrySelector,
				P2pToCapabilities: newCapabilities,
			}

			_, err := changeset.UpdateNodeCapabilities(te.Env, &cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "CapabilityRequiredByDON")
		})
		t.Run("succeeds if update sets new and existing capabilities", func(t *testing.T) {
			existing := getNodeCapabilities(te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry, p2pIDs)

			capabiltiesToSet := existing
			for k, v := range newCapabilities {
				capabiltiesToSet[k] = append(capabiltiesToSet[k], v...)
			}
			cfg := changeset.UpdateNodeCapabilitiesRequest{
				RegistryChainSel:  te.RegistrySelector,
				P2pToCapabilities: capabiltiesToSet,
			}

			csOut, err := changeset.UpdateNodeCapabilities(te.Env, &cfg)
			require.NoError(t, err)
			require.Empty(t, csOut.Proposals)
			require.Nil(t, csOut.AddressBook)

			validateCapabilityUpdates(t, te, capabiltiesToSet)
		})
	})
	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		// contract set is already deployed with capabilities
		// we have to keep track of the existing capabilities to add to the new ones
		var p2pIDs []p2pkey.PeerID
		newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for id := range te.WFNodes {
			k, err := p2pkey.MakePeerID(id)
			require.NoError(t, err)
			p2pIDs = append(p2pIDs, k)
			newCapabilities[k] = caps
		}

		existing := getNodeCapabilities(te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry, p2pIDs)

		capabiltiesToSet := existing
		for k, v := range newCapabilities {
			capabiltiesToSet[k] = append(capabiltiesToSet[k], v...)
		}
		cfg := changeset.UpdateNodeCapabilitiesRequest{
			RegistryChainSel:  te.RegistrySelector,
			P2pToCapabilities: capabiltiesToSet,
			MCMSConfig:        &changeset.MCMSConfig{MinDuration: 0},
		}

		csOut, err := changeset.UpdateNodeCapabilities(te.Env, &cfg)
		require.NoError(t, err)
		require.Len(t, csOut.Proposals, 1)
		require.Len(t, csOut.Proposals[0].Transactions, 1)
		require.Len(t, csOut.Proposals[0].Transactions[0].Batch, 2) // add capabilities, update nodes
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		contracts := te.ContractSets()[te.RegistrySelector]
		timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}

		_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(changeset.UpdateNodeCapabilities),
				Config:    &cfg,
			},
		})
		require.NoError(t, err)
		validateCapabilityUpdates(t, te, capabiltiesToSet)
	})
}

// validateUpdate checks reads nodes from the registry and checks they have the expected updates
func validateCapabilityUpdates(t *testing.T, te test.TestEnv, expected map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability) {
	registry := te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry
	wfP2PIDs := p2pIDs(t, maps.Keys(te.WFNodes))
	nodes, err := registry.GetNodesByP2PIds(nil, wfP2PIDs)
	require.NoError(t, err)
	require.Len(t, nodes, len(wfP2PIDs))
	for _, node := range nodes {
		want := expected[node.P2pId]
		require.NotNil(t, want)
		assertEqualCapabilities(t, registry, want, node)
	}
}

func assertEqualCapabilities(t *testing.T, registry *kcr.CapabilitiesRegistry, want []kcr.CapabilitiesRegistryCapability, got kcr.INodeInfoProviderNodeInfo) {
	wantHashes := make([][32]byte, len(want))
	for i, c := range want {
		h, err := registry.GetHashedCapabilityId(nil, c.LabelledName, c.Version)
		require.NoError(t, err)
		wantHashes[i] = h
	}
	assert.Equal(t, len(want), len(got.HashedCapabilityIds))
	assert.ElementsMatch(t, wantHashes, got.HashedCapabilityIds)
}

func getNodeCapabilities(registry *kcr.CapabilitiesRegistry, p2pIDs []p2pkey.PeerID) map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability {
	m := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	caps, err := registry.GetCapabilities(nil)
	if err != nil {
		panic(err)
	}
	var capMap = make(map[[32]byte]kcr.CapabilitiesRegistryCapability)
	for _, c := range caps {
		capMap[c.HashedId] = kcr.CapabilitiesRegistryCapability{
			LabelledName:          c.LabelledName,
			Version:               c.Version,
			CapabilityType:        c.CapabilityType,
			ResponseType:          c.ResponseType,
			ConfigurationContract: c.ConfigurationContract,
		}
	}
	nodes, err := registry.GetNodesByP2PIds(nil, peerIDsToBytes(p2pIDs))
	if err != nil {
		panic(err)
	}
	for _, n := range nodes {
		caps := make([]kcr.CapabilitiesRegistryCapability, len(n.HashedCapabilityIds))
		for i, h := range n.HashedCapabilityIds {
			c, ok := capMap[h]
			if !ok {
				panic("capability not found")
			}
			caps[i] = c
		}
		m[n.P2pId] = caps
	}
	return m
}

func peerIDsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	bs := make([][32]byte, len(p2pIDs))
	for i, p := range p2pIDs {
		bs[i] = p
	}
	return bs
}
