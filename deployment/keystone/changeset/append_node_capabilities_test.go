package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAppendNodeCapabilities(t *testing.T) {
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
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for _, id := range te.GetP2PIDs("wfDon") {
			newCapabilities[id] = caps
		}

		t.Run("succeeds if existing capabilities not explicit", func(t *testing.T) {
			cfg := changeset.AppendNodeCapabilitiesRequest{
				RegistryChainSel:  te.RegistrySelector,
				P2pToCapabilities: newCapabilities,
				RegistryRef:       te.CapabilityRegistryAddressRef(),
			}

			csOut, err := changeset.AppendNodeCapabilities(te.Env, &cfg)
			require.NoError(t, err)
			require.Empty(t, csOut.MCMSTimelockProposals)
			require.Nil(t, csOut.AddressBook)

			validateCapabilityAppends(t, te, newCapabilities)
		})
	})
	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for _, id := range te.GetP2PIDs("wfDon") {
			newCapabilities[id] = caps
		}

		cfg := changeset.AppendNodeCapabilitiesRequest{
			RegistryChainSel:  te.RegistrySelector,
			P2pToCapabilities: newCapabilities,
			MCMSConfig:        &changeset.MCMSConfig{MinDuration: 0},
			RegistryRef:       te.CapabilityRegistryAddressRef(),
		}

		csOut, err := changeset.AppendNodeCapabilities(te.Env, &cfg)
		require.NoError(t, err)
		require.Len(t, csOut.MCMSTimelockProposals, 1)
		require.Len(t, csOut.MCMSTimelockProposals[0].Operations, 1)
		require.Len(t, csOut.MCMSTimelockProposals[0].Operations[0].Transactions, 2) // add capabilities, update nodes
		require.Nil(t, csOut.AddressBook)

		err = applyProposal(t, te, commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.AppendNodeCapabilities),
			&cfg,
		))
		require.NoError(t, err)
		validateCapabilityAppends(t, te, newCapabilities)
	})
}

// validateUpdate checks reads nodes from the registry and checks they have the expected updates
func validateCapabilityAppends(t *testing.T, te test.EnvWrapper, appended map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability) {
	registry := te.CapabilitiesRegistry()
	wfP2PIDs := te.GetP2PIDs("wfDon").Bytes32()
	nodes, err := registry.GetNodesByP2PIds(nil, wfP2PIDs)
	require.NoError(t, err)
	require.Len(t, nodes, len(wfP2PIDs))
	for _, node := range nodes {
		want := appended[node.P2pId]
		require.NotNil(t, want)
		assertContainsCapabilities(t, registry, want, node)
	}
}

func assertContainsCapabilities(t *testing.T, registry *kcr.CapabilitiesRegistry, want []kcr.CapabilitiesRegistryCapability, got kcr.INodeInfoProviderNodeInfo) {
	wantHashes := make([][32]byte, len(want))
	for i, c := range want {
		h, err := registry.GetHashedCapabilityId(nil, c.LabelledName, c.Version)
		require.NoError(t, err)
		wantHashes[i] = h
		assert.Contains(t, got.HashedCapabilityIds, h, "missing capability %v", c)
	}
	assert.LessOrEqual(t, len(want), len(got.HashedCapabilityIds))
}
