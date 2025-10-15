package v2_0_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cr "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v2_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

type fields struct {
	ContractMetaData types.ContractMetaData
	Capabilities     []v2_0.CapabilityView
	Nodes            []v2_0.NodeView
	Dons             []v2_0.DonView
	Nops             []v2_0.NopView
}

func TestCapRegView_Denormalize(t *testing.T) {
	donConfig := map[string]any{
		"consensus": "basic",
		"timeout":   "30s",
	}
	donCfgProto := pkg.CapabilityConfig(donConfig)
	donConfigBytes, err := donCfgProto.MarshalProto()
	require.NoError(t, err)

	capMetadata := map[string]any{
		"capabilityType": 2,
		"responseType":   0,
	}
	capMetadataBytes, err := json.Marshal(capMetadata)
	require.NoError(t, err)

	capMetadataProto := pkg.CapabilityConfig(capMetadata)
	capMetadataProtoBytes, err := capMetadataProto.MarshalProto()
	require.NoError(t, err)

	t.Run("empty", func(t *testing.T) {
		assertTest(t, fields{}, nil, false)
	})

	t.Run("one don", func(t *testing.T) {
		donView, testErr := v2_0.NewDonView(cr.CapabilitiesRegistryDONInfo{
			Id:               1,
			Name:             "first don",
			IsPublic:         true,
			AcceptsWorkflows: true,
			DonFamilies:      []string{"family1", "family2"},
			Config:           donConfigBytes,
			CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: "test-cap-id",
					Config:       capMetadataProtoBytes,
				},
			},
		})
		require.NoError(t, testErr)
		cap1, testErr := v2_0.NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
			CapabilityId: "test-cap-id",
			Metadata:     capMetadataBytes,
		})
		require.NoError(t, testErr)
		cap2, testErr := v2_0.NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
			CapabilityId: "test-cap-id-2",
			Metadata:     capMetadataBytes,
		})
		require.NoError(t, testErr)

		f := fields{
			Dons: []v2_0.DonView{
				donView,
			},
			Nodes: []v2_0.NodeView{
				v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
					CapabilitiesDONIds: []*big.Int{big.NewInt(1)},
					NodeOperatorId:     1, // 1-based index
				}),
			},
			Nops: []v2_0.NopView{
				{Name: "first nop"},
			},
			Capabilities: []v2_0.CapabilityView{
				cap1, cap2,
			},
		}
		w := []v2_0.DonDenormalizedView{
			{
				Don: donView.DonUniversalMetadata,
				Nodes: []v2_0.NodeDenormalizedView{
					{
						NodeUniversalMetadata: v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
							CapabilitiesDONIds: []*big.Int{big.NewInt(1)},
							NodeOperatorId:     1, // 1-based index
						}).NodeUniversalMetadata,
						Nop: v2_0.NopView{Name: "first nop"},
					},
				},
				Capabilities: []v2_0.CapabilityView{
					cap1,
				},
			},
		}

		assertTest(t, f, w, false)
	})

	t.Run("two dons, multiple capabilities", func(t *testing.T) {
		donView, testErr := v2_0.NewDonView(cr.CapabilitiesRegistryDONInfo{
			Id:               1,
			Name:             "first don",
			IsPublic:         true,
			AcceptsWorkflows: true,
			DonFamilies:      []string{"family1", "family2"},
			Config:           donConfigBytes,
			CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: "test-cap-id",
					Config:       capMetadataProtoBytes,
				},
				{
					CapabilityId: "test-cap-id-2",
					Config:       capMetadataProtoBytes,
				},
			},
		})
		require.NoError(t, testErr)
		donView2, testErr := v2_0.NewDonView(cr.CapabilitiesRegistryDONInfo{
			Id:               2,
			Name:             "second don",
			IsPublic:         true,
			AcceptsWorkflows: true,
			DonFamilies:      []string{"family2"},
			Config:           donConfigBytes,
			CapabilityConfigurations: []cr.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: "other-cap-id",
					Config:       capMetadataProtoBytes,
				},
			},
		})
		require.NoError(t, testErr)
		cap1, testErr := v2_0.NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
			CapabilityId: "test-cap-id",
			Metadata:     capMetadataBytes,
		})
		require.NoError(t, testErr)
		cap2, testErr := v2_0.NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
			CapabilityId: "test-cap-id-2",
			Metadata:     capMetadataBytes,
		})
		require.NoError(t, testErr)
		cap3, testErr := v2_0.NewCapabilityView(cr.CapabilitiesRegistryCapabilityInfo{
			CapabilityId: "other-cap-id",
			Metadata:     capMetadataBytes,
		})
		require.NoError(t, testErr)

		f := fields{
			Dons: []v2_0.DonView{
				donView, donView2,
			},
			Nodes: []v2_0.NodeView{
				v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
					P2pId:              [32]byte{31: 1},
					CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
					NodeOperatorId:     1,                         // 1-based index
				}),

				v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
					P2pId:              [32]byte{31: 11},
					CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
					NodeOperatorId:     3,                         // 1-based index
				}),

				// nodes for don2
				v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
					P2pId:              [32]byte{31: 22},
					CapabilitiesDONIds: []*big.Int{big.NewInt(2)}, // matches don ID 2
					NodeOperatorId:     2,                         // 1-based index
				}),
			},
			Nops: []v2_0.NopView{
				{Name: "first nop"},
				{Name: "second nop"},
				{Name: "third nop"},
			},
			Capabilities: []v2_0.CapabilityView{
				cap1, cap2, cap3,
			},
		}
		w := []v2_0.DonDenormalizedView{
			{
				Don: donView.DonUniversalMetadata,
				Nodes: []v2_0.NodeDenormalizedView{
					{
						NodeUniversalMetadata: v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
							P2pId:              [32]byte{31: 1},
							CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
							NodeOperatorId:     1,                         // 1-based index
						}).NodeUniversalMetadata,
						Nop: v2_0.NopView{Name: "first nop"},
					},
					{
						NodeUniversalMetadata: v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
							P2pId:              [32]byte{31: 11},
							CapabilitiesDONIds: []*big.Int{big.NewInt(1)}, // matches don ID 1
							NodeOperatorId:     3,                         // 1-based index
						}).NodeUniversalMetadata,
						Nop: v2_0.NopView{Name: "third nop"},
					},
				},
				Capabilities: []v2_0.CapabilityView{
					cap1, cap2,
				},
			},
			{
				Don: donView2.DonUniversalMetadata,
				Nodes: []v2_0.NodeDenormalizedView{
					{
						NodeUniversalMetadata: v2_0.NewNodeView(cr.INodeInfoProviderNodeInfo{
							P2pId:              [32]byte{31: 22},
							CapabilitiesDONIds: []*big.Int{big.NewInt(2)}, // matches don ID 2
							NodeOperatorId:     2,                         // 1-based index
						}).NodeUniversalMetadata,
						Nop: v2_0.NopView{Name: "second nop"},
					},
				},
				Capabilities: []v2_0.CapabilityView{
					cap3,
				},
			},
		}

		assertTest(t, f, w, false)
	})
}

func assertTest(t *testing.T, f fields, want []v2_0.DonDenormalizedView, wantErr bool) {
	v := v2_0.CapabilityRegistryView{
		ContractMetaData: f.ContractMetaData,
		Capabilities:     f.Capabilities,
		Nodes:            f.Nodes,
		Dons:             f.Dons,
		Nops:             f.Nops,
	}
	got, err := v.DonDenormalizedView()
	if (err != nil) != wantErr {
		t.Errorf("CapRegView.Denormalize() error = %v, wantErr %v", err, wantErr)
		return
	}
	for i := range got {
		assert.Equal(t, want[i].Don, got[i].Don)
		for j := range got[i].Nodes {
			assert.Equal(t, want[i].Nodes[j].NodeUniversalMetadata, got[i].Nodes[j].NodeUniversalMetadata, "NodeUniversalMetadata mismatch at index %d for don %d", j, i)
			assert.Equal(t, want[i].Nodes[j].Nop, got[i].Nodes[j].Nop, "Nop mismatch at index %d for don %d", j, i)
		}
		assert.Equal(t, want[i].Capabilities, got[i].Capabilities)
	}
}
