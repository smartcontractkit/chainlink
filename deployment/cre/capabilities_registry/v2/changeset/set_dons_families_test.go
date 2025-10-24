package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/common/view/v2_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestSetDONsFamilies_VerifyPreconditions(t *testing.T) {
	cs := changeset.SetDONsFamilies{}

	env := test.SetupEnvV2(t, false)
	chainSelector := env.RegistrySelector

	t.Run("invalid registry selector", func(t *testing.T) {
		err := cs.VerifyPreconditions(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:    0, // invalid
			RegistryQualifier:   "qual",
			DONsFamiliesChanges: []sequences.DONFamiliesChange{{DonName: "don-1", AddToFamilies: []string{"fam-1"}}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistrySelector")
	})

	t.Run("empty qualifier", func(t *testing.T) {
		err := cs.VerifyPreconditions(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:    chainSelector,
			RegistryQualifier:   "",
			DONsFamiliesChanges: []sequences.DONFamiliesChange{{DonName: "don-1", AddToFamilies: []string{"fam-1"}}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryQualifier")
	})

	t.Run("no DON family changes", func(t *testing.T) {
		err := cs.VerifyPreconditions(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: "test",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must specify at least one DON family change")
	})
}

func TestSetDONsFamilies_Apply(t *testing.T) {
	cs := changeset.SetDONsFamilies{}

	env := test.SetupEnvV2(t, false)
	chainSelector := env.RegistrySelector

	chain, ok := env.Env.BlockChains.EVMChains()[env.RegistrySelector]
	require.True(t, ok, "chain not found for selector")

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		env.RegistryAddress, chain.Client,
	)
	require.NoError(t, err)

	originalDON, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)
	require.Len(t, originalDON.DonFamilies, 1)
	require.Contains(t, originalDON.DonFamilies, "test-family")

	t.Run("validates DONs Families Changes input", func(t *testing.T) {
		_, testErr := cs.Apply(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:            test.DONName,
					AddToFamilies:      []string{},
					RemoveFromFamilies: []string{},
				},
			},
		})
		require.Error(t, testErr)
		assert.Contains(t, testErr.Error(), "must specify at least one family to add or remove")
	})

	t.Run("set families for existing DON", func(t *testing.T) {
		_, testErr := cs.Apply(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:       test.DONName,
					AddToFamilies: []string{"family-new", "family-common"},
				},
			},
		})
		require.NoError(t, testErr)

		updatedDON, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON.DonFamilies, 3)
		assert.Contains(t, updatedDON.DonFamilies, "family-new", "family-common")
	})

	t.Run("remove families for existing DON", func(t *testing.T) {
		_, testErr := cs.Apply(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:            test.DONName,
					RemoveFromFamilies: []string{"family-common"},
				},
			},
		})
		require.NoError(t, testErr)

		updatedDON, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON.DonFamilies, 2)
		assert.Contains(t, updatedDON.DonFamilies, "test-family", "family-new")
	})

	t.Run("remove ALL families for existing DON", func(t *testing.T) {
		_, testErr := cs.Apply(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:            test.DONName,
					RemoveFromFamilies: []string{"test-family", "family-new", "family-common"},
				},
			},
		})
		require.NoError(t, testErr)

		updatedDON, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Empty(t, updatedDON.DonFamilies)
	})

	t.Run("set families for multiple DONs", func(t *testing.T) {
		// First, create another DON to test with
		extendedCapReg := v2_0.ExtendedCapabilityRegistry{CapabilitiesRegistry: capReg}

		existingCaps, testErr := extendedCapReg.GetCapabilitiesSimple(nil)
		require.NoError(t, testErr)
		existingNodes, testErr := extendedCapReg.GetNodesSimple(nil)
		require.NoError(t, testErr)

		var caps []changeset.CapabilitiesRegistryCapability
		for _, c := range existingCaps {
			caps = append(caps, changeset.CapabilitiesRegistryCapability{
				CapabilityID:          c.CapabilityId,
				ConfigurationContract: c.ConfigurationContract,
			})
		}

		var p2pIDs []string
		for _, n := range existingNodes {
			p2pIDs = append(p2pIDs, p2pkey.PeerID(n.P2pId).String())
		}

		don := changeset.CapabilitiesRegistryNewDONParams{
			Name:             "second-test-don",
			DonFamilies:      []string{"family-a"},
			Nodes:            p2pIDs,
			F:                2,
			IsPublic:         true,
			AcceptsWorkflows: false,
		}

		configureInput := changeset.ConfigureCapabilitiesRegistryInput{
			ChainSelector: chainSelector,
			Qualifier:     test.RegistryQualifier,
			Capabilities:  caps,
			DONs:          []changeset.CapabilitiesRegistryNewDONParams{don},
		}

		_, testErr = changeset.ConfigureCapabilitiesRegistry{}.Apply(*env.Env, configureInput)
		require.NoError(t, testErr)

		_, testErr = cs.Apply(*env.Env, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:       test.DONName,
					AddToFamilies: []string{"test-family", "family-new", "family-common"},
				},
				{
					DonName:            don.Name,
					AddToFamilies:      []string{"test-family"},
					RemoveFromFamilies: []string{"family-a"},
				},
			},
		})
		require.NoError(t, testErr)

		updatedDON1, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON1.DonFamilies, 3)
		assert.Contains(t, updatedDON1.DonFamilies, "test-family", "family-new", "family-common")

		updatedDON2, testErr := capReg.GetDONByName(nil, don.Name)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON2.DonFamilies, 1)
		assert.Contains(t, updatedDON2.DonFamilies, "test-family")
	})
}
