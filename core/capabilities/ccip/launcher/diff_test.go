package launcher

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_diff(t *testing.T) {
	type args struct {
		capabilityID string
		oldState     *registry.MetadataRegistry
		newState     *registry.MetadataRegistry
	}
	tests := []struct {
		name    string
		args    args
		want    diffResult
		wantErr bool
	}{
		{
			name: "no diff",
			args: args{
				capabilityID: defaultCapability.ID,
				oldState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registry.DonID]registry.DON{
						1: defaultRegistryDon,
					},
					IDsToNodes: map[types.PeerID]registry.NodeInfo{},
				},
				newState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registry.DonID]registry.DON{
						1: defaultRegistryDon,
					},
					IDsToNodes: map[types.PeerID]registry.NodeInfo{},
				},
			},
			want: diffResult{
				added:   map[registry.DonID]registry.DON{},
				removed: map[registry.DonID]registry.DON{},
				updated: map[registry.DonID]registry.DON{},
			},
		},
		{
			"capability not present",
			args{
				capabilityID: defaultCapability.ID,
				oldState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						newCapability.ID: newCapability,
					},
				},
				newState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						newCapability.ID: newCapability,
					},
				},
			},
			diffResult{},
			true,
		},
		{
			"diff present, new don",
			args{
				capabilityID: defaultCapability.ID,
				oldState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registry.DonID]registry.DON{},
				},
				newState: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registry.DonID]registry.DON{
						1: defaultRegistryDon,
					},
				},
			},
			diffResult{
				added: map[registry.DonID]registry.DON{
					1: defaultRegistryDon,
				},
				removed: map[registry.DonID]registry.DON{},
				updated: map[registry.DonID]registry.DON{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := diff(tt.args.capabilityID, tt.args.oldState, tt.args.newState)
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
		currCCIPDONs map[registry.DonID]registry.DON
		newCCIPDONs  map[registry.DonID]registry.DON
	}
	tests := []struct {
		name        string
		args        args
		wantAdded   map[registry.DonID]registry.DON
		wantRemoved map[registry.DonID]registry.DON
		wantUpdated map[registry.DonID]registry.DON
		wantErr     bool
	}{
		{
			"added dons",
			args{
				currCCIPDONs: map[registry.DonID]registry.DON{},
				newCCIPDONs: map[registry.DonID]registry.DON{
					1: defaultRegistryDon,
				},
			},
			map[registry.DonID]registry.DON{
				1: defaultRegistryDon,
			},
			map[registry.DonID]registry.DON{},
			map[registry.DonID]registry.DON{},
			false,
		},
		{
			"removed dons",
			args{
				currCCIPDONs: map[registry.DonID]registry.DON{
					1: defaultRegistryDon,
				},
				newCCIPDONs: map[registry.DonID]registry.DON{},
			},
			map[registry.DonID]registry.DON{},
			map[registry.DonID]registry.DON{
				1: defaultRegistryDon,
			},
			map[registry.DonID]registry.DON{},
			false,
		},
		{
			"updated dons",
			args{
				currCCIPDONs: map[registry.DonID]registry.DON{
					1: defaultRegistryDon,
				},
				newCCIPDONs: map[registry.DonID]registry.DON{
					1: {
						DON:                      getDON(defaultRegistryDon.ID, defaultRegistryDon.Members, defaultRegistryDon.ConfigVersion+1),
						CapabilityConfigurations: defaultCapCfgs,
					},
				},
			},
			map[registry.DonID]registry.DON{},
			map[registry.DonID]registry.DON{},
			map[registry.DonID]registry.DON{
				1: {
					DON:                      getDON(defaultRegistryDon.ID, defaultRegistryDon.Members, defaultRegistryDon.ConfigVersion+1),
					CapabilityConfigurations: defaultCapCfgs,
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
		ccipCapability registry.Capability
		state          *registry.MetadataRegistry
	}
	tests := []struct {
		name    string
		args    args
		want    map[registry.DonID]registry.DON
		wantErr bool
	}{
		{
			"one ccip don",
			args{
				ccipCapability: defaultCapability,
				state: &registry.MetadataRegistry{
					IDsToDONs: map[registry.DonID]registry.DON{
						1: defaultRegistryDon,
					},
				},
			},
			map[registry.DonID]registry.DON{
				1: defaultRegistryDon,
			},
			false,
		},
		{
			"no ccip dons - different capability",
			args{
				ccipCapability: newCapability,
				state: &registry.MetadataRegistry{
					IDsToDONs: map[registry.DonID]registry.DON{
						1: defaultRegistryDon,
					},
				},
			},
			map[registry.DonID]registry.DON{},
			false,
		},
		{
			"don with multiple capabilities, one of them ccip",
			args{
				ccipCapability: defaultCapability,
				state: &registry.MetadataRegistry{
					IDsToDONs: map[registry.DonID]registry.DON{
						1: {
							DON: getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
							CapabilityConfigurations: map[string]registry.CapabilityConfiguration{
								defaultCapability.ID: {},
								newCapability.ID:     {},
							},
						},
					},
				},
			},
			map[registry.DonID]registry.DON{
				1: {
					DON: getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
					CapabilityConfigurations: map[string]registry.CapabilityConfiguration{
						defaultCapability.ID: {},
						newCapability.ID:     {},
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
		capabilityID string
		state        *registry.MetadataRegistry
	}
	tests := []struct {
		name    string
		args    args
		want    registry.Capability
		wantErr bool
	}{
		{
			"in registry state",
			args{
				capabilityID: defaultCapability.ID,
				state: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						defaultCapability.ID: defaultCapability,
					},
				},
			},
			defaultCapability,
			false,
		},
		{
			"not in registry state",
			args{
				capabilityID: defaultCapability.ID,
				state: &registry.MetadataRegistry{
					IDsToCapabilities: map[string]registry.Capability{
						newCapability.ID: newCapability,
					},
				},
			},
			registry.Capability{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkCapabilityPresence(tt.args.capabilityID, tt.args.state)
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

func Test_isMemberOfDON(t *testing.T) {
	p2pIDs := make([]ragep2ptypes.PeerID, 0, 4)
	for i := range [4]struct{}{} {
		p2pIDs = append(p2pIDs, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID()))
	}
	don := registry.DON{
		DON: getDON(1, p2pIDs, 0),
	}
	require.True(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID())))
	require.False(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(5)).PeerID())))
}
