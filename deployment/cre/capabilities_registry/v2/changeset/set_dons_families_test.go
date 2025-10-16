package changeset_test

import (
	"testing"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
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
		require.Len(t, updatedDON.DonFamilies, 3)
		require.Contains(t, updatedDON.DonFamilies, "family-new", "family-common")
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
		require.Len(t, updatedDON.DonFamilies, 2)
		require.Contains(t, updatedDON.DonFamilies, "test-family", "family-new")
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
		require.Len(t, updatedDON.DonFamilies, 0)
	})
}
