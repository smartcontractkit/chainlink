package changeset_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAddDONs(t *testing.T) {
	t.Parallel()
	var testCap = kcr.CapabilitiesRegistryCapability{
		LabelledName:   "cap1",
		Version:        "1.0",
		CapabilityType: 0,
	}
	type input struct {
		te         test.EnvWrapper
		dons       []*changeset.RegisterableDon
		mcmsConfig *changeset.MCMSConfig
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
		// we need to compose a new don to register
		virtualDon := setupVirtualDONs(t, te, mc.mcmsConfig, testCap)

		t.Run(mc.name, func(t *testing.T) {
			var cases = []testCase{
				{
					name: "ok",
					input: input{
						te: te,
						dons: []*changeset.RegisterableDon{
							{
								Name:   "a_new_don",
								F:      1,
								P2PIDs: virtualDon,
								CapabilityConfigs: []changeset.CapabilityConfig{
									{
										Capability: testCap,
									},
								},
							},
						},
						mcmsConfig: mc.mcmsConfig,
					},
				},
				{
					name: "node not exists",
					input: input{
						te: te,
						dons: []*changeset.RegisterableDon{
							{
								Name:   "a_new_don",
								F:      1,
								P2PIDs: []p2pkey.PeerID{virtualDon[0], virtualDon[1], virtualDon[2], p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()},
								CapabilityConfigs: []changeset.CapabilityConfig{
									{
										Capability: testCap,
									},
								},
							},
						},
						mcmsConfig: mc.mcmsConfig,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
					},
				},
				{
					name: "cap no exists",
					input: input{
						te: te,
						dons: []*changeset.RegisterableDon{
							{
								Name:   "a_new_don",
								F:      1,
								P2PIDs: virtualDon[1:6],
								CapabilityConfigs: []changeset.CapabilityConfig{
									{
										Capability: kcr.CapabilitiesRegistryCapability{
											LabelledName:   "cap_no_exists",
											Version:        "1.0",
											CapabilityType: 0,
										},
									},
								},
							},
						},
						mcmsConfig: mc.mcmsConfig,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
						if !useMCMS { // mcms swallows error messages
							assert.ErrorContains(t, err, "CapabilityDoesNotExist")
						}
					},
				},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					// contract set is already deployed with capabilities
					// we have to keep track of the existing capabilities to add to the new ones

					cfg := changeset.AddDonsRequest{
						RegistryChainSel: te.RegistrySelector,
						DONs:             tc.input.dons,
						MCMSConfig:       tc.input.mcmsConfig,
						RegistryRef:      te.CapabilityRegistryAddressRef(),
					}

					csOut, err := changeset.AddDons(te.Env, &cfg)
					if err != nil && tc.checkErr == nil {
						t.Errorf("non nil err from AddDons %v but no checkErr func defined", err)
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
							cldf.CreateLegacyChangeSet(changeset.AddDons),
							&cfg,
						))
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, applyErr)
							return
						}
					}

					assertDonExists(t, te.CapabilitiesRegistry(), tc.input.dons[0].P2PIDs)
				})
			}
		})
	}
}

func assertDonExists(t *testing.T, registry *kcr.CapabilitiesRegistry, p2pIDs []p2pkey.PeerID) {
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
}

func setupVirtualDONs(t *testing.T, te test.EnvWrapper, mcmsConfig *changeset.MCMSConfig, capability kcr.CapabilitiesRegistryCapability) []p2pkey.PeerID {
	t.Helper()
	virtualDon := make([]p2pkey.PeerID, 0, 8)
	virtualDon = append(virtualDon, te.GetP2PIDs("wfDon")...)
	virtualDon = append(virtualDon, te.GetP2PIDs("assetDon")...)
	x := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, p2pID := range virtualDon {
		x[p2pID] = []kcr.CapabilitiesRegistryCapability{
			capability,
		}
	}
	appendReq := &changeset.AppendNodeCapabilitiesRequest{
		RegistryChainSel:  te.RegistrySelector,
		P2pToCapabilities: x,
		MCMSConfig:        mcmsConfig,
		RegistryRef:       te.CapabilityRegistryAddressRef(),
	}
	csOut, err := changeset.AppendNodeCapabilities(te.Env, appendReq)
	require.NoError(t, err)
	if mcmsConfig != nil {
		require.NotNil(t, csOut.MCMSTimelockProposals)
		require.Len(t, csOut.MCMSTimelockProposals, 1)
		applyErr := applyProposal(t, te, commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.AppendNodeCapabilities),
			appendReq,
		))
		require.NoError(t, applyErr)
	}
	return virtualDon
}
