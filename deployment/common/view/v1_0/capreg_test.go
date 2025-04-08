package v1_0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	cr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

func TestCapRegView_Denormalize(t *testing.T) {
	type fields struct {
		ContractMetaData types.ContractMetaData
		Capabilities     []CapabilityView
		Nodes            []NodeView
		Dons             []DonView
		Nops             []NopView
	}
	tests := []struct {
		name    string
		fields  fields
		want    []DonDenormalizedView
		wantErr bool
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   nil,
		},
		{
			name: "one don",
			fields: fields{
				Dons: []DonView{
					NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 1,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{0: 1},
							},
						},
					}),
				},
				Nodes: []NodeView{
					NewNodeView(cr.INodeInfoProviderNodeInfo{
						CapabilitiesDONIds: []*big.Int{big.NewInt(1)},
						NodeOperatorId:     1, // 1-based index
					}),
				},
				Nops: []NopView{
					{Name: "first nop"},
				},
				Capabilities: []CapabilityView{
					NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
						HashedId:     [32]byte{0: 1},
						LabelledName: "cap1",
						Version:      "1.0.0",
					}),

					NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
						HashedId:     [32]byte{0: 2},
						LabelledName: "cap2",
						Version:      "1.0.0",
					}),
				},
			},
			want: []DonDenormalizedView{
				{
					Don: NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 1,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{0: 1},
							},
						},
					}).DonUniversalMetadata,
					Nodes: []NodeDenormalizedView{
						{
							NodeUniversalMetadata: NewNodeView(cr.INodeInfoProviderNodeInfo{
								CapabilitiesDONIds: []*big.Int{big.NewInt(1)},
								NodeOperatorId:     1, // 1-based index
							}).NodeUniversalMetadata,
							Nop: NopView{Name: "first nop"},
						},
					},
					Capabilities: []CapabilityView{
						NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
							HashedId:     [32]byte{0: 1},
							LabelledName: "cap1",
							Version:      "1.0.0",
						}),
					},
				},
			},
		},

		{
			name: "two dons, multiple capabilities",
			fields: fields{
				Dons: []DonView{
					// don1
					NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 1,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{0: 1},
							},
							{
								CapabilityId: [32]byte{1: 1},
							},
						},
					}),

					// don2
					NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 2,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{2: 2},
							},
						},
					}),
				},
				Nodes: []NodeView{
					// nodes for don1
					NewNodeView(cr.INodeInfoProviderNodeInfo{
						P2pId:              [32]byte{31: 1},
						CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
						NodeOperatorId:     1,                         // 1-based index
					}),

					NewNodeView(cr.INodeInfoProviderNodeInfo{
						P2pId:              [32]byte{31: 11},
						CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
						NodeOperatorId:     3,                         // 1-based index
					}),

					// nodes for don2
					NewNodeView(cr.INodeInfoProviderNodeInfo{
						P2pId:              [32]byte{31: 22},
						CapabilitiesDONIds: []*big.Int{big.NewInt(2)}, // matches don ID 2
						NodeOperatorId:     2,                         // 1-based index
					}),
				},
				Nops: []NopView{
					{Name: "first nop"},
					{Name: "second nop"},
					{Name: "third nop"},
				},
				Capabilities: []CapabilityView{
					// capabilities for don1
					NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
						HashedId:     [32]byte{0: 1},
						LabelledName: "cap1",
						Version:      "1.0.0",
					}),
					NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
						HashedId:     [32]byte{1: 1}, // matches don ID 1, capabitility ID 2
						LabelledName: "cap2",
						Version:      "1.0.0",
					}),

					// capabilities for don2
					NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
						HashedId:     [32]byte{2: 2}, // matches don ID 2, capabitility ID 1
						LabelledName: "other cap",
						Version:      "1.0.0",
					}),
				},
			},
			want: []DonDenormalizedView{
				{
					Don: NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 1,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{0: 1},
							},
							{
								CapabilityId: [32]byte{1: 1},
							},
						},
					}).DonUniversalMetadata,
					Nodes: []NodeDenormalizedView{
						{
							NodeUniversalMetadata: NewNodeView(cr.INodeInfoProviderNodeInfo{
								P2pId:              [32]byte{31: 1},
								CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
							}).NodeUniversalMetadata,
							Nop: NopView{Name: "first nop"},
						},
						{
							NodeUniversalMetadata: NewNodeView(cr.INodeInfoProviderNodeInfo{
								P2pId:              [32]byte{31: 11},
								CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
							}).NodeUniversalMetadata,
							Nop: NopView{Name: "third nop"},
						},
					},
					Capabilities: []CapabilityView{
						NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
							HashedId:     [32]byte{0: 1},
							LabelledName: "cap1",
							Version:      "1.0.0",
						}),

						NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
							HashedId:     [32]byte{1: 1},
							LabelledName: "cap2",
							Version:      "1.0.0",
						}),
					},
				},
				{
					Don: NewDonView(cr.CapabilitiesRegistryDONInfo{
						Id: 2,
						CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: [32]byte{2: 2},
							},
						},
					}).DonUniversalMetadata,
					Nodes: []NodeDenormalizedView{
						{
							NodeUniversalMetadata: NewNodeView(cr.INodeInfoProviderNodeInfo{
								P2pId:              [32]byte{31: 22},
								CapabilitiesDONIds: []*big.Int{big.NewInt(2)}, // matches don ID 2
							}).NodeUniversalMetadata,
							Nop: NopView{Name: "second nop"},
						},
					},
					Capabilities: []CapabilityView{
						NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
							HashedId:     [32]byte{2: 2}, // matches don ID 2, capabitility ID 1
							LabelledName: "other cap",
							Version:      "1.0.0",
						}),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CapabilityRegistryView{
				ContractMetaData: tt.fields.ContractMetaData,
				Capabilities:     tt.fields.Capabilities,
				Nodes:            tt.fields.Nodes,
				Dons:             tt.fields.Dons,
				Nops:             tt.fields.Nops,
			}
			got, err := v.DonDenormalizedView()
			if (err != nil) != tt.wantErr {
				t.Errorf("CapRegView.Denormalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				assert.Equal(t, tt.want[i].Don, got[i].Don)
				for j := range got[i].Nodes {
					assert.Equal(t, tt.want[i].Nodes[j].NodeUniversalMetadata, got[i].Nodes[j].NodeUniversalMetadata, "NodeUniversalMetadata mismatch at index %d for don %d", j, i)
					assert.Equal(t, tt.want[i].Nodes[j].Nop, got[i].Nodes[j].Nop, "Nop mismatch at index %d for don %d", j, i)
				}
				assert.Equal(t, tt.want[i].Capabilities, got[i].Capabilities)
			}
		})
	}
}
