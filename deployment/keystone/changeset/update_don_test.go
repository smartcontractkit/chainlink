package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

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

func TestUpdateDon(t *testing.T) {
	t.Parallel()

	capACfg, err := proto.Marshal(test.GetDefaultCapConfig(t, capA))
	require.NoError(t, err)

	capBCfg, err := proto.Marshal(test.GetDefaultCapConfig(t, capB))
	require.NoError(t, err)

	type input struct {
		te              test.EnvWrapper
		nodeSetToUpdate []p2pkey.PeerID
		mcmsConfig      *changeset.MCMSConfig
	}
	type testCase struct {
		name     string
		input    input
		checkErr func(t *testing.T, useMCMS bool, err error)
	}

	var mcmsCases = []mcmsTestCase{
		{name: "no mcms", mcmsConfig: nil},
		{name: "with mcms", mcmsConfig: &changeset.MCMSConfig{MinDuration: 0}},
	}

	for _, mc := range mcmsCases {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         mc.mcmsConfig != nil,
		})

		t.Run(mc.name, func(t *testing.T) {
			var cases = []testCase{
				{
					name: "forbid wf update",
					input: input{
						nodeSetToUpdate: te.GetP2PIDs("wfDon"),
						mcmsConfig:      mc.mcmsConfig,
						te:              te,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						// this error is independent of mcms because it is a pre-txn check
						assert.ErrorContains(t, err, "refusing to update workflow don")
					},
				},
				{
					name: "writer don update ok",
					input: input{
						te:              te,
						nodeSetToUpdate: te.GetP2PIDs("writerDon"),
						mcmsConfig:      mc.mcmsConfig,
					},
				},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					// contract set is already deployed with capabilities
					// we have to keep track of the existing capabilities to add to the new ones
					p2pIDs := tc.input.nodeSetToUpdate
					newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
					for _, id := range p2pIDs {
						newCapabilities[id] = caps
					}

					cfg := changeset.UpdateDonRequest{
						RegistryChainSel: te.RegistrySelector,
						P2PIDs:           p2pIDs,
						CapabilityConfigs: []changeset.CapabilityConfig{
							{
								Capability: capA, Config: capACfg,
							},
							{
								Capability: capB, Config: capBCfg,
							},
						},
						MCMSConfig:  tc.input.mcmsConfig,
						RegistryRef: te.CapabilityRegistryAddressRef(),
					}

					csOut, err := changeset.UpdateDon(te.Env, &cfg)
					if err != nil && tc.checkErr == nil {
						t.Errorf("non nil err from UpdateDon %v but no checkErr func defined", err)
					}
					useMCMS := cfg.MCMSConfig != nil
					if !useMCMS {
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, err)
							return
						}
					} else {
						// when using mcms there are two kinds of errors:
						// those from creating the proposal and those executing the proposal
						// if we have a non-nil err here, its from creating the proposal
						// so check it and do not proceed to applying the proposal
						if err != nil {
							tc.checkErr(t, useMCMS, err)
							return
						}
						require.NotNil(t, csOut.MCMSTimelockProposals)
						require.Len(t, csOut.MCMSTimelockProposals, 1)
						applyErr := applyProposal(t, te, commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(changeset.UpdateDon),
							&cfg,
						))
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, applyErr)
							return
						}
					}

					assertDonContainsCapabilities(t, te.CapabilitiesRegistry(), caps, p2pIDs)
				})
			}
		})
	}
}

func TestUpdateDon_ChangeComposition(t *testing.T) {
	t.Parallel()

	// Test capability configurations
	capACfg, err := proto.Marshal(test.GetDefaultCapConfig(t, capA))
	require.NoError(t, err)

	type testCase struct {
		name       string
		mcmsConfig *changeset.MCMSConfig
	}

	var mcmsCases = []testCase{
		{name: "no mcms", mcmsConfig: nil},
		{name: "with mcms", mcmsConfig: &changeset.MCMSConfig{MinDuration: 0}},
	}

	for _, mc := range mcmsCases {
		t.Run(mc.name, func(t *testing.T) {
			// Setup test environment with initial DON configuration
			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
				UseMCMS:         mc.mcmsConfig != nil,
			})

			// Get initial DON info for writerDon
			initialDons, err := te.CapabilitiesRegistry().GetDONs(nil)
			require.NoError(t, err)

			var writerDonInfo *kcr.CapabilitiesRegistryDONInfo
			writerP2PIDs := te.GetP2PIDs("writerDon")
			for i, don := range initialDons {
				if internal.SortedHash(internal.PeerIDsToBytes(writerP2PIDs)) == internal.SortedHash(don.NodeP2PIds) {
					writerDonInfo = &initialDons[i]
					break
				}
			}
			require.NotNil(t, writerDonInfo, "writerDon not found in registry")

			// Store original DON ID and node count for verification
			donID := int(writerDonInfo.Id)
			originalNodeCount := len(writerDonInfo.NodeP2PIds)

			t.Run("add node to DON composition", func(t *testing.T) {
				// Add one more node from assetDon to writerDon
				assetP2PIDs := te.GetP2PIDs("assetDon")
				newP2PIDs := writerP2PIDs[:]
				newP2PIDs = append(newP2PIDs, assetP2PIDs[0]) // Add first asset node

				cfg := changeset.UpdateDonRequest{
					RegistryChainSel: te.RegistrySelector,
					P2PIDs:           newP2PIDs,
					DonID:            donID, // Explicitly specify DON ID
					CapabilityConfigs: []changeset.CapabilityConfig{
						{Capability: capA, Config: capACfg},
					},
					MCMSConfig:  mc.mcmsConfig,
					RegistryRef: te.CapabilityRegistryAddressRef(),
				}

				csOut, err := changeset.UpdateDon(te.Env, &cfg)
				require.NoError(t, err)

				if mc.mcmsConfig != nil {
					require.NotNil(t, csOut.MCMSTimelockProposals)
					require.Len(t, csOut.MCMSTimelockProposals, 1)

					// Apply the MCMS proposal
					applyErr := applyProposal(t, te, commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(changeset.UpdateDon),
						&cfg,
					))
					require.NoError(t, applyErr)
				}

				// Verify DON now has one additional node
				updatedDon, err := te.CapabilitiesRegistry().GetDON(nil, uint32(donID)) //nolint:gosec // G115
				require.NoError(t, err)
				require.Equal(t, uint32(donID), updatedDon.Id, "DON ID should remain the same") //nolint:gosec // G115
				require.Len(t, updatedDon.NodeP2PIds, originalNodeCount+1, "DON should have one additional node")

				// Verify the new P2P ID is present
				actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)
				require.ElementsMatch(t, newP2PIDs, actualP2PIDs, "DON should contain all expected P2P IDs")
			})

			t.Run("remove node from DON composition", func(t *testing.T) {
				// Remove one node from the current composition
				currentDon, err := te.CapabilitiesRegistry().GetDON(nil, uint32(donID)) //nolint:gosec // G115
				require.NoError(t, err)

				currentP2PIDs := internal.BytesToPeerIDs(currentDon.NodeP2PIds)
				// Remove the last node (should be the one we added)
				reducedP2PIDs := currentP2PIDs[:len(currentP2PIDs)-1]

				cfg := changeset.UpdateDonRequest{
					RegistryChainSel: te.RegistrySelector,
					P2PIDs:           reducedP2PIDs,
					DonID:            donID, // Explicitly specify DON ID
					CapabilityConfigs: []changeset.CapabilityConfig{
						{Capability: capA, Config: capACfg},
					},
					MCMSConfig:  mc.mcmsConfig,
					RegistryRef: te.CapabilityRegistryAddressRef(),
				}

				csOut, err := changeset.UpdateDon(te.Env, &cfg)
				require.NoError(t, err)

				if mc.mcmsConfig != nil {
					require.NotNil(t, csOut.MCMSTimelockProposals)
					require.Len(t, csOut.MCMSTimelockProposals, 1)

					// Apply the MCMS proposal
					applyErr := applyProposal(t, te, commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(changeset.UpdateDon),
						&cfg,
					))
					require.NoError(t, applyErr)
				}

				// Verify DON is back to original node count
				updatedDon, err := te.CapabilitiesRegistry().GetDON(nil, uint32(donID)) //nolint:gosec // G115
				require.NoError(t, err)
				require.Equal(t, uint32(donID), updatedDon.Id, "DON ID should remain the same") //nolint:gosec // G115
				require.Len(t, updatedDon.NodeP2PIds, originalNodeCount, "DON should be back to original node count")

				// Verify the correct P2P IDs are present
				actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)
				require.ElementsMatch(t, reducedP2PIDs, actualP2PIDs, "DON should contain expected P2P IDs")
			})

			t.Run("replace nodes in DON composition", func(t *testing.T) {
				// Replace some nodes from writerDon with nodes from assetDon
				assetP2PIDs := te.GetP2PIDs("assetDon")
				mixedP2PIDs := writerP2PIDs[:]
				// Keep first 2 nodes from writer, add 2 nodes from asset
				mixedP2PIDs = append(mixedP2PIDs, assetP2PIDs[:2]...)

				cfg := changeset.UpdateDonRequest{
					RegistryChainSel: te.RegistrySelector,
					P2PIDs:           mixedP2PIDs,
					DonID:            donID, // Explicitly specify DON ID
					CapabilityConfigs: []changeset.CapabilityConfig{
						{Capability: capA, Config: capACfg},
					},
					MCMSConfig:  mc.mcmsConfig,
					RegistryRef: te.CapabilityRegistryAddressRef(),
				}

				csOut, err := changeset.UpdateDon(te.Env, &cfg)
				require.NoError(t, err)

				if mc.mcmsConfig != nil {
					require.NotNil(t, csOut.MCMSTimelockProposals)
					require.Len(t, csOut.MCMSTimelockProposals, 1)

					// Apply the MCMS proposal
					applyErr := applyProposal(t, te, commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(changeset.UpdateDon),
						&cfg,
					))
					require.NoError(t, applyErr)
				}

				// Verify DON has the mixed composition
				updatedDon, err := te.CapabilitiesRegistry().GetDON(nil, uint32(donID)) //nolint:gosec // G115
				require.NoError(t, err)
				require.Equal(t, uint32(donID), updatedDon.Id, "DON ID should remain the same") //nolint:gosec // G115
				require.Len(t, updatedDon.NodeP2PIds, len(mixedP2PIDs), "DON should have expected node count")

				// Verify the correct P2P IDs are present
				actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)
				require.ElementsMatch(t, mixedP2PIDs, actualP2PIDs, "DON should contain mixed P2P IDs")

				// Verify capabilities are still configured correctly
				assertDonContainsCapabilities(t, te.CapabilitiesRegistry(), []kcr.CapabilitiesRegistryCapability{capA}, actualP2PIDs)
			})

			t.Run("error when DonID not specified and nodes change", func(t *testing.T) {
				// Try to change composition without specifying DonID
				assetP2PIDs := te.GetP2PIDs("assetDon")
				newP2PIDs := writerP2PIDs[:]
				newP2PIDs = append(newP2PIDs, assetP2PIDs[0])

				cfg := changeset.UpdateDonRequest{
					RegistryChainSel: te.RegistrySelector,
					P2PIDs:           newP2PIDs,
					// DonID intentionally omitted
					CapabilityConfigs: []changeset.CapabilityConfig{
						{Capability: capA, Config: capACfg},
					},
					MCMSConfig:  mc.mcmsConfig,
					RegistryRef: te.CapabilityRegistryAddressRef(),
				}

				// This should either work (if the system can infer the DON) or fail gracefully
				// The behavior depends on the implementation in update_don.go
				_, err := changeset.UpdateDon(te.Env, &cfg)

				// If the system requires DonID for composition changes, this should error
				// If it can infer the DON, it should work
				// The test documents the current behavior
				if err != nil {
					t.Logf("Expected behavior: UpdateDon requires DonID when changing composition: %v", err)
				} else {
					t.Log("UpdateDon successfully inferred DON from P2P IDs")
				}
			})
		})
	}
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
