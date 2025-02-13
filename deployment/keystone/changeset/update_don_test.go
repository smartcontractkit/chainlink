package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestUpdateDon(t *testing.T) {
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
	capACfg := test.GetDefaultCapConfig(t, capA)
	capACfgB, err := proto.Marshal(capACfg)
	require.NoError(t, err)
	capBCfg := test.GetDefaultCapConfig(t, capB)
	capBCfgB, err := proto.Marshal(capBCfg)
	require.NoError(t, err)

	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		// contract set is already deployed with capabilities
		// we have to keep track of the existing capabilities to add to the new ones
		p2pIDs := te.GetP2PIDs("wfDon")
		newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for _, id := range p2pIDs {
			newCapabilities[id] = caps
		}

		t.Run("succeeds if update sets new and existing capabilities", func(t *testing.T) {
			cfg := changeset.UpdateDonRequest{
				RegistryChainSel: te.RegistrySelector,
				P2PIDs:           p2pIDs,
				CapabilityConfigs: []changeset.CapabilityConfig{
					{
						Capability: capA, Config: capACfgB,
					},
					{
						Capability: capB, Config: capBCfgB,
					},
				},
			}

			csOut, err := changeset.UpdateDon(te.Env, &cfg)
			require.NoError(t, err)
			require.Empty(t, csOut.Proposals)
			require.Nil(t, csOut.AddressBook)

			assertDonContainsCapabilities(t, te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry, caps, p2pIDs)
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

		// contract set is already deployed with capabilities
		// we have to keep track of the existing capabilities to add to the new ones
		p2pIDs := te.GetP2PIDs("wfDon")

		cfg := changeset.UpdateDonRequest{
			RegistryChainSel: te.RegistrySelector,
			P2PIDs:           p2pIDs,
			CapabilityConfigs: []changeset.CapabilityConfig{
				{
					Capability: capA,
					Config:     capACfgB,
				},
				{
					Capability: capB,
					Config:     capBCfgB,
				},
			},
			MCMSConfig: &changeset.MCMSConfig{MinDuration: 0},
		}

		csOut, err := changeset.UpdateDon(te.Env, &cfg)
		require.NoError(t, err)

		require.Len(t, csOut.Proposals, 1)
		require.Len(t, csOut.Proposals[0].Transactions, 1)          // append node capabilties cs, update don
		require.Len(t, csOut.Proposals[0].Transactions[0].Batch, 3) // add capabilities, update nodes, update don
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		contracts := te.ContractSets()[te.RegistrySelector]
		timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}
		_, err = commonchangeset.Apply(t, te.Env, timelockContracts,
			commonchangeset.Configure(
				deployment.CreateLegacyChangeSet(changeset.UpdateDon),
				&cfg,
			),
		)
		require.NoError(t, err)
		assertDonContainsCapabilities(t, te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry, caps, p2pIDs)
	})
}

func assertDonContainsCapabilities(t *testing.T, registry *kcr.CapabilitiesRegistry, want []kcr.CapabilitiesRegistryCapability, p2pIDs []p2pkey.PeerID) {
	dons, err := registry.GetDONs(nil)
	require.NoError(t, err)
	var got *kcr.CapabilitiesRegistryDONInfo
	for i, don := range dons {
		if internal.SortedHash(internal.PeerIDsToBytes(p2pIDs)) == internal.SortedHash(don.NodeP2PIds) {
			got = &dons[i]
			break
		}
	}
	require.NotNil(t, got, "missing don with p2pIDs %v", p2pIDs)
	wantHashes := make([][32]byte, len(want))
	for i, c := range want {
		h, err := registry.GetHashedCapabilityId(nil, c.LabelledName, c.Version)
		require.NoError(t, err)
		wantHashes[i] = h
		assert.Contains(t, capIDsFromCapCfgs(got.CapabilityConfigurations), h, "missing capability %v", c)
	}
	assert.LessOrEqual(t, len(want), len(got.CapabilityConfigurations), "too many capabilities")
}

func capIDsFromCapCfgs(cfgs []kcr.CapabilitiesRegistryCapabilityConfiguration) [][32]byte {
	out := make([][32]byte, len(cfgs))
	for i, c := range cfgs {
		out[i] = c.CapabilityId
	}
	return out
}
