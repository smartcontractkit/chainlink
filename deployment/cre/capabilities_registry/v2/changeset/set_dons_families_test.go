package changeset_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v2_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestSetDONsFamilies_VerifyPreconditions(t *testing.T) {
	cs := changeset.SetDONsFamilies{}

	h := test.NewTestHarness(t)
	chainSelector := h.RegistrySelector

	t.Run("invalid registry selector", func(t *testing.T) {
		err := cs.VerifyPreconditions(h.Runtime.Environment(), changeset.SetDONsFamiliesInput{
			RegistrySelector:    0, // invalid
			RegistryQualifier:   "qual",
			DONsFamiliesChanges: []sequences.DONFamiliesChange{{DonName: "don-1", AddToFamilies: []string{"fam-1"}}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistrySelector")
	})

	t.Run("empty qualifier", func(t *testing.T) {
		err := cs.VerifyPreconditions(h.Runtime.Environment(), changeset.SetDONsFamiliesInput{
			RegistrySelector:    chainSelector,
			RegistryQualifier:   "",
			DONsFamiliesChanges: []sequences.DONFamiliesChange{{DonName: "don-1", AddToFamilies: []string{"fam-1"}}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryQualifier")
	})

	t.Run("no DON family changes", func(t *testing.T) {
		err := cs.VerifyPreconditions(h.Runtime.Environment(), changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: "test",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must specify at least one DON family change")
	})
}

func TestSetDONsFamilies_Apply(t *testing.T) {
	cs := changeset.SetDONsFamilies{}

	h := test.NewTestHarness(t)
	chainSelector := h.RegistrySelector

	chain, ok := h.Runtime.Environment().BlockChains.EVMChains()[h.RegistrySelector]
	require.True(t, ok, "chain not found for selector")

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		h.RegistryAddress, chain.Client,
	)
	require.NoError(t, err)

	originalDON, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)
	require.Len(t, originalDON.DonFamilies, 1)
	require.Contains(t, originalDON.DonFamilies, "test-family")

	t.Run("validates DONs Families Changes input", func(t *testing.T) {
		testErr := h.Runtime.Exec(
			runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
				RegistrySelector:  chainSelector,
				RegistryQualifier: test.RegistryQualifier,
				DONsFamiliesChanges: []sequences.DONFamiliesChange{
					{
						DonName:            test.DONName,
						AddToFamilies:      []string{},
						RemoveFromFamilies: []string{},
					},
				},
			}),
		)
		require.Error(t, testErr)
		require.ErrorContains(t, testErr, "must specify at least one family to add or remove")
	})

	t.Run("set families for existing DON", func(t *testing.T) {
		testErr := h.Runtime.Exec(
			runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
				RegistrySelector:  chainSelector,
				RegistryQualifier: test.RegistryQualifier,
				DONsFamiliesChanges: []sequences.DONFamiliesChange{
					{
						DonName:       test.DONName,
						AddToFamilies: []string{"family-new", "family-common"},
					},
				},
			}),
		)
		require.NoError(t, testErr)

		updatedDON, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON.DonFamilies, 3)
		assert.Contains(t, updatedDON.DonFamilies, "family-new", "family-common")
	})

	t.Run("set families for existing DON - MCMS", func(t *testing.T) {
		mcmsEnv := test.NewTestHarness(t, test.WithMCMS())

		duration := mcmstypes.NewDuration(1 * time.Second)

		task := runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
			RegistrySelector:  chainSelector,
			RegistryQualifier: test.RegistryQualifier,
			DONsFamiliesChanges: []sequences.DONFamiliesChange{
				{
					DonName:       test.DONName,
					AddToFamilies: []string{"family-new", "family-common"},
				},
			},
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					chainSelector: "",
				},
				ValidDuration: &duration,
			},
		})

		testErr := mcmsEnv.Runtime.Exec(task)
		require.NoError(t, testErr)

		csOut := mcmsEnv.Runtime.State().Outputs[task.ID()]

		// Verify the changeset output
		assert.NotNil(t, csOut.Reports, "reports should be present")
		assert.NotEmpty(t, csOut.MCMSTimelockProposals, "should have MCMS proposals when using MCMS")
	})

	t.Run("remove families for existing DON", func(t *testing.T) {
		testErr := h.Runtime.Exec(
			runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
				RegistrySelector:  chainSelector,
				RegistryQualifier: test.RegistryQualifier,
				DONsFamiliesChanges: []sequences.DONFamiliesChange{
					{
						DonName:            test.DONName,
						RemoveFromFamilies: []string{"family-common"},
					},
				},
			}),
		)
		require.NoError(t, testErr)

		updatedDON, testErr := capReg.GetDONByName(nil, test.DONName)
		require.NoError(t, testErr)
		assert.Len(t, updatedDON.DonFamilies, 2)
		assert.Contains(t, updatedDON.DonFamilies, "test-family", "family-new")
	})

	t.Run("remove ALL families for existing DON", func(t *testing.T) {
		testErr := h.Runtime.Exec(
			runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
				RegistrySelector:  chainSelector,
				RegistryQualifier: test.RegistryQualifier,
				DONsFamiliesChanges: []sequences.DONFamiliesChange{
					{
						DonName:            test.DONName,
						RemoveFromFamilies: []string{"test-family", "family-new", "family-common"},
					},
				},
			}),
		)
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

		testErr = h.Runtime.Exec(
			runtime.ChangesetTask(changeset.ConfigureCapabilitiesRegistry{}, configureInput),
		)
		require.NoError(t, testErr)

		testErr = h.Runtime.Exec(
			runtime.ChangesetTask(cs, changeset.SetDONsFamiliesInput{
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
			}),
		)
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
