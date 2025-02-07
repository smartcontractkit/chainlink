package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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

	// run the same tests for both mcms and non-mcms
	var mcmsConfigs = []*changeset.MCMSConfig{nil, {MinDuration: 0}}
	for _, mcmsConfig := range mcmsConfigs {
		prefix := "no mcms"
		if mcmsConfig != nil {
			prefix = "with mcms"
		}
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         mcmsConfig != nil,
		})

		t.Run(prefix, func(t *testing.T) {
			type testCase struct {
				name            string
				nodeSetToUpdate map[string]memory.Node
				checkErr        func(t *testing.T, useMCMS bool, err error)
				mcmsConfig      *changeset.MCMSConfig
			}
			var cases = []testCase{
				{
					name:            "forbid wf update",
					nodeSetToUpdate: te.WFNodes,
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						if useMCMS {
							assert.ErrorContains(t, err, "sadfasdds")
						}
						assert.ErrorContains(t, err, "refusing to update workflow don")
					},
					mcmsConfig: mcmsConfig,
				},
				{
					name:            "writer don update ok",
					nodeSetToUpdate: te.CWNodes,
					mcmsConfig:      mcmsConfig,
				},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {

					// contract set is already deployed with capabilities
					// we have to keep track of the existing capabilities to add to the new ones
					var p2pIDs []p2pkey.PeerID
					newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
					for id := range tc.nodeSetToUpdate {
						k, err := p2pkey.MakePeerID(id)
						require.NoError(t, err)
						p2pIDs = append(p2pIDs, k)
						newCapabilities[k] = caps
					}

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
						MCMSConfig: tc.mcmsConfig,
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
						require.NotNil(t, csOut.Proposals) //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
						require.Len(t, csOut.Proposals, 1) //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
						applyErr := applyProposal(t, te, []commonchangeset.ChangesetApplication{
							{
								Changeset: commonchangeset.WrapChangeSet(changeset.UpdateDon),
								Config:    &cfg,
							},
						})
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, applyErr)
							return
						}
					}
					//			require.NoError(t, err)
					//			require.Empty(t, csOut.Proposals)
					//			require.Nil(t, csOut.AddressBook)

					assertDonContainsCapabilities(t, te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry, caps, p2pIDs)
				})
			}
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

func applyProposal(t *testing.T, te test.TestEnv, applicable []commonchangeset.ChangesetApplication) error {
	// now apply the changeset such that the proposal is signed and execed
	contracts := te.ContractSets()[te.RegistrySelector]
	timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
		te.RegistrySelector: {
			Timelock:  contracts.Timelock,
			CallProxy: contracts.CallProxy,
		},
	}
	_, err := commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, applicable)
	return err
}
