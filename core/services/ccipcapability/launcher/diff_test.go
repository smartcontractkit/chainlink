package launcher

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func Test_diff(t *testing.T) {
	type args struct {
		capabilityVersion      string
		capabilityLabelledName string
		oldState               registrysyncer.State
		newState               registrysyncer.State
	}
	tests := []struct {
		name    string
		args    args
		want    diffResult
		wantErr bool
	}{
		{
			"no diff",
			args{
				capabilityVersion:      "v1.0.0",
				capabilityLabelledName: "ccip",
				oldState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.0.0"): {
							LabelledName: "ccip",
							Version:      "v1.0.0",
						},
					},
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
								},
							},
						},
					},
					IDsToNodes: map[types.PeerID]kcr.CapabilitiesRegistryNodeInfo{},
				},
				newState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.0.0"): {
							LabelledName: "ccip",
							Version:      "v1.0.0",
						},
					},
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
								},
							},
						},
					},
					IDsToNodes: map[types.PeerID]kcr.CapabilitiesRegistryNodeInfo{},
				},
			},
			diffResult{
				added:   map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				removed: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				updated: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			},
			false,
		},
		{
			"capability not present",
			args{
				capabilityVersion:      "v1.0.0",
				capabilityLabelledName: "ccip",
				oldState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.1.0"): {
							LabelledName: "ccip",
							Version:      "v1.1.0",
						},
					},
				},
				newState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.1.0"): {
							LabelledName: "ccip",
							Version:      "v1.1.0",
						},
					},
				},
			},
			diffResult{},
			true,
		},
		{
			"diff present, new don",
			args{
				capabilityVersion:      "v1.0.0",
				capabilityLabelledName: "ccip",
				oldState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.0.0"): {
							LabelledName: "ccip",
							Version:      "v1.0.0",
						},
					},
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				},
				newState: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.0.0"): {
							LabelledName: "ccip",
							Version:      "v1.0.0",
						},
					},
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
								},
							},
						},
					},
				},
			},
			diffResult{
				added: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
					1: {
						Id: 1,
						CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
							},
						},
					},
				},
				removed: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				updated: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := diff(tt.args.capabilityVersion, tt.args.capabilityLabelledName, tt.args.oldState, tt.args.newState)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_compareDONs(t *testing.T) {
	type args struct {
		currCCIPDONs map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
		newCCIPDONs  map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
	}
	tests := []struct {
		name        string
		args        args
		wantAdded   map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
		wantRemoved map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
		wantUpdated map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
		wantErr     bool
	}{
		{
			"added dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				newCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
					1: {
						Id: 1,
					},
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
				1: {
					Id: 1,
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			false,
		},
		{
			"removed dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
					1: {
						Id: 1,
					},
				},
				newCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
				1: {
					Id: 1,
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			false,
		},
		{
			"updated dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
					1: {
						Id:          1,
						ConfigCount: 1,
					},
				},
				newCCIPDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
					1: {
						Id:          1,
						ConfigCount: 2,
					},
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
				1: {
					Id:          1,
					ConfigCount: 2,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr, err := compareDONs(tt.args.currCCIPDONs, tt.args.newCCIPDONs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantAdded, dr.added)
				require.Equal(t, tt.wantRemoved, dr.removed)
				require.Equal(t, tt.wantUpdated, dr.updated)
			}
		})
	}
}

func Test_filterCCIPDONs(t *testing.T) {
	type args struct {
		ccipCapability kcr.CapabilitiesRegistryCapabilityInfo
		state          registrysyncer.State
	}
	tests := []struct {
		name    string
		args    args
		want    map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
		wantErr bool
	}{
		{
			"one ccip don",
			args{
				ccipCapability: kcr.CapabilitiesRegistryCapabilityInfo{
					LabelledName: "ccip",
					Version:      "v1.0.0",
				},
				state: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
								},
							},
						},
					},
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
				1: {
					Id: 1,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
						},
					},
				},
			},
			false,
		},
		{
			"no ccip dons",
			args{
				ccipCapability: kcr.CapabilitiesRegistryCapabilityInfo{
					LabelledName: "ccip",
					Version:      "v1.0.0",
				},
				state: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.1.0"),
								},
							},
						},
					},
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
			false,
		},
		{
			"don with multiple capabilities, one of them ccip",
			args{
				ccipCapability: kcr.CapabilitiesRegistryCapabilityInfo{
					LabelledName: "ccip",
					Version:      "v1.0.0",
				},
				state: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
								},
								{
									CapabilityId: mustHashedCapabilityId("ccip", "v1.1.0"),
								},
							},
						},
					},
				},
			},
			map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
				1: {
					Id: 1,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: mustHashedCapabilityId("ccip", "v1.0.0"),
						},
						{
							CapabilityId: mustHashedCapabilityId("ccip", "v1.1.0"),
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterCCIPDONs(tt.args.ccipCapability, tt.args.state)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_checkCapabilityPresence(t *testing.T) {
	type args struct {
		capabilityVersion      string
		capabilityLabelledName string
		state                  registrysyncer.State
	}
	tests := []struct {
		name    string
		args    args
		want    kcr.CapabilitiesRegistryCapabilityInfo
		wantErr bool
	}{
		{
			"in registry state",
			args{
				capabilityVersion:      "v1.0.0",
				capabilityLabelledName: "ccip",
				state: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.0.0"): {
							LabelledName: "ccip",
							Version:      "v1.0.0",
						},
						mustHashedCapabilityId("ccip", "v1.1.0"): {
							LabelledName: "ccip",
							Version:      "v1.1.0",
						},
					},
				},
			},
			kcr.CapabilitiesRegistryCapabilityInfo{
				LabelledName: "ccip",
				Version:      "v1.0.0",
			},
			false,
		},
		{
			"not in registry state",
			args{
				capabilityVersion:      "v1.0.0",
				capabilityLabelledName: "ccip",
				state: registrysyncer.State{
					IDsToCapabilities: map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
						mustHashedCapabilityId("ccip", "v1.1.0"): {
							LabelledName: "ccip",
							Version:      "v1.1.0",
						},
					},
				},
			},
			kcr.CapabilitiesRegistryCapabilityInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkCapabilityPresence(tt.args.capabilityVersion, tt.args.capabilityLabelledName, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkCapabilityPresence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkCapabilityPresence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashedCapabilityId(t *testing.T) {
	transactor := testutils.MustNewSimTransactor(t)
	sb := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	crAddress, _, _, err := kcr.DeployCapabilitiesRegistry(transactor, sb)
	require.NoError(t, err)
	sb.Commit()

	cr, err := kcr.NewCapabilitiesRegistry(crAddress, sb)
	require.NoError(t, err)

	// add a capability, ignore cap config for simplicity.
	_, err = cr.AddCapabilities(transactor, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          "ccip",
			Version:               "v1.0.0",
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: common.Address{},
		},
	})
	require.NoError(t, err)
	sb.Commit()

	hidExpected, err := cr.GetHashedCapabilityId(nil, "ccip", "v1.0.0")
	require.NoError(t, err)

	hid, err := hashedCapabilityId("ccip", "v1.0.0")
	require.NoError(t, err)

	require.Equal(t, hidExpected, hid)
}

func Test_isMemberOfDON(t *testing.T) {
	var p2pIDs [][32]byte
	for i := range [4]struct{}{} {
		p2pIDs = append(p2pIDs, p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID())
	}
	don := kcr.CapabilitiesRegistryDONInfo{
		Id:         1,
		NodeP2PIds: p2pIDs,
	}
	require.True(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID())))
	require.False(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(5)).PeerID())))
}

func Test_isMemberOfBootstrapSubcommittee(t *testing.T) {
	var bootstrapKeys [][32]byte
	for i := range [4]struct{}{} {
		bootstrapKeys = append(bootstrapKeys, p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID())
	}
	require.True(t, isMemberOfBootstrapSubcommittee(bootstrapKeys, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID())))
	require.False(t, isMemberOfBootstrapSubcommittee(bootstrapKeys, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(5)).PeerID())))
}
