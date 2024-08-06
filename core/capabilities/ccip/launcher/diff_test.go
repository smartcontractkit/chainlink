package launcher

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

const (
	ccipCapVersion    = "v1.0.0"
	ccipCapNewVersion = "v1.1.0"
	ccipCapName       = "ccip"
)

var defaultCapability = getCapability(ccipCapName, ccipCapVersion)
var newCapability = getCapability(ccipCapName, ccipCapNewVersion)

func getDON(id uint32, members []ragep2ptypes.PeerID, cfgVersion uint32) capabilities.DON {
	return capabilities.DON{
		ID:               id,
		ConfigVersion:    cfgVersion,
		F:                uint8(1),
		IsPublic:         true,
		AcceptsWorkflows: true,
		Members:          members,
	}
}
func getP2PID(id uint32) ragep2ptypes.PeerID {
	return ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(id))).PeerID())
}

var p2pID1 = getP2PID(1)
var p2pID2 = getP2PID(2)

var defaultCapCfgs = map[string]capabilities.CapabilityConfiguration{
	defaultCapability.ID: {},
}
var defaultRegistryDon = registrysyncer.DON{
	DON:                      getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
	CapabilityConfigurations: defaultCapCfgs,
}

func Test_diff(t *testing.T) {
	type args struct {
		capabilityID string
		oldState     registrysyncer.LocalRegistry
		newState     registrysyncer.LocalRegistry
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
				oldState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
					IDsToNodes: map[types.PeerID]kcr.CapabilitiesRegistryNodeInfo{},
				},
				newState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
					IDsToNodes: map[types.PeerID]kcr.CapabilitiesRegistryNodeInfo{},
				},
			},
			want: diffResult{
				added:   map[registrysyncer.DonID]registrysyncer.DON{},
				removed: map[registrysyncer.DonID]registrysyncer.DON{},
				updated: map[registrysyncer.DonID]registrysyncer.DON{},
			},
		},
		{
			"capability not present",
			args{
				capabilityID: defaultCapability.ID,
				oldState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						newCapability.ID: newCapability,
					},
				},
				newState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
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
				oldState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{},
				},
				newState: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						defaultCapability.ID: defaultCapability,
					},
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			diffResult{
				added: map[registrysyncer.DonID]registrysyncer.DON{
					1: defaultRegistryDon,
				},
				removed: map[registrysyncer.DonID]registrysyncer.DON{},
				updated: map[registrysyncer.DonID]registrysyncer.DON{},
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
		currCCIPDONs map[registrysyncer.DonID]registrysyncer.DON
		newCCIPDONs  map[registrysyncer.DonID]registrysyncer.DON
	}
	tests := []struct {
		name        string
		args        args
		wantAdded   map[registrysyncer.DonID]registrysyncer.DON
		wantRemoved map[registrysyncer.DonID]registrysyncer.DON
		wantUpdated map[registrysyncer.DonID]registrysyncer.DON
		wantErr     bool
	}{
		{
			"added dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{},
				newCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{
					1: defaultRegistryDon,
				},
			},
			map[registrysyncer.DonID]registrysyncer.DON{
				1: defaultRegistryDon,
			},
			map[registrysyncer.DonID]registrysyncer.DON{},
			map[registrysyncer.DonID]registrysyncer.DON{},
			false,
		},
		{
			"removed dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{
					1: defaultRegistryDon,
				},
				newCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{},
			},
			map[registrysyncer.DonID]registrysyncer.DON{},
			map[registrysyncer.DonID]registrysyncer.DON{
				1: defaultRegistryDon,
			},
			map[registrysyncer.DonID]registrysyncer.DON{},
			false,
		},
		{
			"updated dons",
			args{
				currCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{
					1: defaultRegistryDon,
				},
				newCCIPDONs: map[registrysyncer.DonID]registrysyncer.DON{
					1: {
						DON:                      getDON(defaultRegistryDon.ID, defaultRegistryDon.Members, defaultRegistryDon.ConfigVersion+1),
						CapabilityConfigurations: defaultCapCfgs,
					},
				},
			},
			map[registrysyncer.DonID]registrysyncer.DON{},
			map[registrysyncer.DonID]registrysyncer.DON{},
			map[registrysyncer.DonID]registrysyncer.DON{
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
		ccipCapability registrysyncer.Capability
		state          registrysyncer.LocalRegistry
	}
	tests := []struct {
		name    string
		args    args
		want    map[registrysyncer.DonID]registrysyncer.DON
		wantErr bool
	}{
		{
			"one ccip don",
			args{
				ccipCapability: defaultCapability,
				state: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			map[registrysyncer.DonID]registrysyncer.DON{
				1: defaultRegistryDon,
			},
			false,
		},
		{
			"no ccip dons - different capability",
			args{
				ccipCapability: newCapability,
				state: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			map[registrysyncer.DonID]registrysyncer.DON{},
			false,
		},
		{
			"don with multiple capabilities, one of them ccip",
			args{
				ccipCapability: defaultCapability,
				state: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: {
							DON: getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
							CapabilityConfigurations: map[string]capabilities.CapabilityConfiguration{
								defaultCapability.ID: {},
								newCapability.ID:     {},
							},
						},
					},
				},
			},
			map[registrysyncer.DonID]registrysyncer.DON{
				1: {
					DON: getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
					CapabilityConfigurations: map[string]capabilities.CapabilityConfiguration{
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
		state        registrysyncer.LocalRegistry
	}
	tests := []struct {
		name    string
		args    args
		want    registrysyncer.Capability
		wantErr bool
	}{
		{
			"in registry state",
			args{
				capabilityID: defaultCapability.ID,
				state: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
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
				state: registrysyncer.LocalRegistry{
					IDsToCapabilities: map[string]registrysyncer.Capability{
						newCapability.ID: newCapability,
					},
				},
			},
			registrysyncer.Capability{},
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
	var p2pIDs []ragep2ptypes.PeerID
	for i := range [4]struct{}{} {
		p2pIDs = append(p2pIDs, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID()))
	}
	don := registrysyncer.DON{
		DON: getDON(1, p2pIDs, 0),
	}
	require.True(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID())))
	require.False(t, isMemberOfDON(don, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(5)).PeerID())))
}

func Test_isMemberOfBootstrapSubcommittee(t *testing.T) {
	var bootstrapKeys [][32]byte
	for i := range [4]struct{}{} {
		bootstrapKeys = append(bootstrapKeys, p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID())
	}
	require.True(t, isMemberOfBootstrapSubcommittee(bootstrapKeys, p2pID1))
	require.False(t, isMemberOfBootstrapSubcommittee(bootstrapKeys, getP2PID(5)))
}

func getCapability(ccipCapName, ccipCapVersion string) registrysyncer.Capability {
	id := fmt.Sprintf("%s@%s", ccipCapName, ccipCapVersion)
	return registrysyncer.Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             id,
	}
}

func mustHashedCapabilityID(capabilityLabelledName, capabilityVersion string) [32]byte {
	r, err := capcommon.HashedCapabilityID(capabilityLabelledName, capabilityVersion)
	if err != nil {
		panic(fmt.Errorf("failed to hash capability id (labelled name: %s, version: %s): %w", capabilityLabelledName, capabilityVersion, err))
	}
	return r
}
