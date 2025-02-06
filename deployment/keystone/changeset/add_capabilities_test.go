package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

func TestAddCapabilities(t *testing.T) {
	t.Parallel()

	capabilitiesToAdd := []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:   "test-cap",
			Version:        "0.0.1",
			CapabilityType: 1,
		},
		{
			LabelledName:   "test-cap-2",
			Version:        "0.0.1",
			CapabilityType: 1,
		},
	}
	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		csOut, err := changeset.AddCapabilities(te.Env, &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			Capabilities:     capabilitiesToAdd,
		})
		require.NoError(t, err)
		require.Empty(t, csOut.Proposals)
		require.Nil(t, csOut.AddressBook)
		assertCapabilitiesExist(t, te.CapabilitiesRegistry(), capabilitiesToAdd...)
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		req := &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			Capabilities:     capabilitiesToAdd,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
		}
		csOut, err := changeset.AddCapabilities(te.Env, req)
		require.NoError(t, err)
		require.Len(t, csOut.Proposals, 1)
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
				Changeset: commonchangeset.WrapChangeSet(changeset.AddCapabilities),
				Config:    req,
			},
		})
		require.NoError(t, err)

		assertCapabilitiesExist(t, te.CapabilitiesRegistry(), capabilitiesToAdd...)
	})
}

func assertCapabilitiesExist(t *testing.T, registry *kcr.CapabilitiesRegistry, capabilities ...kcr.CapabilitiesRegistryCapability) {
	for _, capability := range capabilities {
		wantID, err := registry.GetHashedCapabilityId(nil, capability.LabelledName, capability.Version)
		require.NoError(t, err)
		got, err := registry.GetCapability(nil, wantID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, capability.CapabilityType, got.CapabilityType)
		assert.Equal(t, capability.LabelledName, got.LabelledName)
		assert.Equal(t, capability.Version, got.Version)
	}
}
