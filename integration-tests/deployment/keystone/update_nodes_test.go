package keystone_test

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	kstest "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/test"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func Test_UpdateNodesRequest_validate(t *testing.T) {
	type fields struct {
		p2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
		//nopToNodes        map[uint32][]*kslib.P2PSigner
		nopToNodes map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc
		chain      deployment.Chain
		registry   *kcr.CapabilitiesRegistry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "err",
			fields: fields{
				p2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{},
				nopToNodes:        nil,
				chain:             deployment.Chain{},
				registry:          nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &kslib.UpdateNodesRequest{
				P2pToCapabilities: tt.fields.p2pToCapabilities,
				NopToNodes:        tt.fields.nopToNodes,
				Chain:             tt.fields.chain,
				Registry:          tt.fields.registry,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("kslib.UpdateNodesRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateNodes(t *testing.T) {
	chain := testChain(t)
	require.NotNil(t, chain)
	lggr := logger.Test(t)

	type args struct {
		lggr logger.Logger
		req  *kslib.UpdateNodesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *kslib.UpdateNodesResponse
		wantErr bool
	}{
		{
			name: "one node, one capability",
			args: args{
				lggr: lggr,
				req: &kslib.UpdateNodesRequest{
					P2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
						testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap1",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
						},
					},
					NopToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
						testNop(t, "nop1"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_1"),
								Signer:              [32]byte{0: 1, 1: 2},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
					},
					Chain:    chain,
					Registry: nil, // set in test to ensure no conflicts
				},
			},
			want: &kslib.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one node, two capabilities",
			args: args{
				lggr: lggr,
				req: &kslib.UpdateNodesRequest{
					P2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
						testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap1",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
							{
								LabelledName:   "cap2",
								Version:        "1.0.1",
								CapabilityType: 2,
							},
						},
					},
					NopToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
						testNop(t, "nop1"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_1"),
								Signer:              [32]byte{0: 1, 1: 2},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
					},
					Chain:    chain,
					Registry: nil, // set in test to ensure no conflicts
				},
			},
			want: &kslib.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "twos node, one shared capability",
			args: args{
				lggr: lggr,
				req: &kslib.UpdateNodesRequest{
					P2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
						testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap1",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
						},
						testPeerID(t, "peerID_2"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap1",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
						},
					},
					NopToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
						testNop(t, "nopA"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_1"),
								Signer:              [32]byte{0: 1, 31: 1},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
						testNop(t, "nopB"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_2"),
								Signer:              [32]byte{0: 2, 31: 2},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
					},
					Chain:    chain,
					Registry: nil, // set in test to ensure no conflicts
				},
			},
			want: &kslib.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 1, 31: 1},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
					{
						NodeOperatorId:      2,
						P2pId:               testPeerID(t, "peerID_2"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 2, 31: 2},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "twos node, different capabilities",
			args: args{
				lggr: lggr,
				req: &kslib.UpdateNodesRequest{
					P2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
						testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap1",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
						},
						testPeerID(t, "peerID_2"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap2",
								Version:        "1.0.1",
								CapabilityType: 0,
							},
						},
					},
					NopToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
						testNop(t, "nopA"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_1"),
								Signer:              [32]byte{0: 1, 31: 1},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
						testNop(t, "nopB"): []*kslib.P2PSignerEnc{
							{
								P2PKey:              testPeerID(t, "peerID_2"),
								Signer:              [32]byte{0: 2, 31: 2},
								EncryptionPublicKey: [32]byte{0: 7, 1: 7},
							},
						},
					},
					Chain:    chain,
					Registry: nil, // set in test to ensure no conflicts
				},
			},
			want: &kslib.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 1, 31: 1},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
					{
						NodeOperatorId:      2,
						P2pId:               testPeerID(t, "peerID_2"),
						HashedCapabilityIds: nil, // checked dynamically based on the request
						Signer:              [32]byte{0: 2, 31: 2},
						EncryptionPublicKey: [32]byte{0: 7, 1: 7},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// need to setup the registry and chain with a phony capability so that there is something to update
			var phonyCap = kcr.CapabilitiesRegistryCapability{
				LabelledName:   "phony",
				Version:        "1.0.0",
				CapabilityType: 0,
			}
			initMap := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
			for p2pID := range tt.args.req.P2pToCapabilities {
				initMap[p2pID] = []kcr.CapabilitiesRegistryCapability{phonyCap}
			}
			setupResp := kstest.SetupTestRegistry(t, tt.args.lggr, &kstest.SetupTestRegistryRequest{
				P2pToCapabilities: initMap,
				NopToNodes:        tt.args.req.NopToNodes,
			})
			registry := setupResp.Registry
			tt.args.req.Registry = setupResp.Registry
			tt.args.req.Chain = setupResp.Chain

			//registry := kstest.SetupUpdateNodes(t, tt.args.lggr, tt.args.req)
			//tt.args.req.Registry = registry
			// register the capabilities that the Update will use
			expectedUpdatedCaps := make(map[p2pkey.PeerID][]kslib.RegisteredCapability)
			capCache := kstest.NewCapabiltyCache(t)
			for p2p, newCaps := range tt.args.req.P2pToCapabilities {
				expectedCaps := capCache.AddCapabilities(tt.args.lggr, tt.args.req.Chain, registry, newCaps)
				expectedUpdatedCaps[p2p] = expectedCaps
			}
			got, err := kslib.UpdateNodes(tt.args.lggr, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, p := range got.NodeParams {
				expected := tt.want.NodeParams[i]
				require.Equal(t, expected.NodeOperatorId, p.NodeOperatorId)
				require.Equal(t, expected.P2pId, p.P2pId)
				require.Equal(t, expected.Signer, p.Signer)
				// check the capabilities
				expectedCaps := expectedUpdatedCaps[p.P2pId]
				var wantHashedIds [][32]byte
				for _, cap := range expectedCaps {
					wantHashedIds = append(wantHashedIds, cap.ID)
				}
				sort.Slice(wantHashedIds, func(i, j int) bool {
					return bytes.Compare(wantHashedIds[i][:], wantHashedIds[j][:]) < 0
				})
				gotHashedIds := p.HashedCapabilityIds
				sort.Slice(gotHashedIds, func(i, j int) bool {
					return bytes.Compare(gotHashedIds[i][:], gotHashedIds[j][:]) < 0
				})
				require.Len(t, gotHashedIds, len(wantHashedIds))
				for j, gotCap := range gotHashedIds {
					assert.Equal(t, wantHashedIds[j], gotCap)
				}
			}
		})
	}

	// unique cases
	t.Run("duplicate update idempotent", func(t *testing.T) {
		var (
			p2pToCapabilitiesInitial = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
				testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
					{
						LabelledName:   "first",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
					{
						LabelledName:   "second",
						Version:        "1.0.0",
						CapabilityType: 2,
					},
				},
			}
			p2pToCapabilitiesUpdated = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
				testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
					{
						LabelledName:   "cap1",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
				},
			}
			nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
				testNop(t, "nopA"): []*kslib.P2PSignerEnc{
					{
						P2PKey:              testPeerID(t, "peerID_1"),
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					},
				},
			}
		)

		// setup registry and add one capability
		setupResp := kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{
			P2pToCapabilities: p2pToCapabilitiesInitial,
			NopToNodes:        nopToNodes,
		})
		registry := setupResp.Registry
		chain := setupResp.Chain

		// there should be two capabilities
		info, err := registry.GetNode(&bind.CallOpts{}, testPeerID(t, "peerID_1"))
		require.NoError(t, err)
		require.Len(t, info.HashedCapabilityIds, 2)

		// update the capabilities, there should be then be one capability
		// first update registers the new capability
		toRegister := p2pToCapabilitiesUpdated[testPeerID(t, "peerID_1")]
		tx, err := registry.AddCapabilities(chain.DeployerKey, toRegister)
		if err != nil {
			err2 := kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
			require.Fail(t, fmt.Sprintf("failed to call AddCapabilities: %s:  %s", err, err2))
		}
		_, err = chain.Confirm(tx)
		require.NoError(t, err)

		var req = &kslib.UpdateNodesRequest{
			P2pToCapabilities: p2pToCapabilitiesUpdated,
			NopToNodes:        nopToNodes,
			Chain:             chain,
			Registry:          registry,
		}
		_, err = kslib.UpdateNodes(lggr, req)
		require.NoError(t, err)
		info, err = registry.GetNode(&bind.CallOpts{}, testPeerID(t, "peerID_1"))
		require.NoError(t, err)
		require.Len(t, info.HashedCapabilityIds, 1)
		want := info.HashedCapabilityIds[0]

		// update again and ensure the result is the same
		_, err = kslib.UpdateNodes(lggr, req)
		require.NoError(t, err)
		info, err = registry.GetNode(&bind.CallOpts{}, testPeerID(t, "peerID_1"))
		require.NoError(t, err)
		require.Len(t, info.HashedCapabilityIds, 1)
		got := info.HashedCapabilityIds[0]
		assert.Equal(t, want, got)
	})
}

func TestAppendCapabilities(t *testing.T) {

	var (
		capMap = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
			testPeerID(t, "peerID_1"): []kcr.CapabilitiesRegistryCapability{
				{
					LabelledName:   "cap1",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
			},
		}
		nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
			testNop(t, "nop"): []*kslib.P2PSignerEnc{
				{
					P2PKey:              testPeerID(t, "peerID_1"),
					Signer:              [32]byte{0: 1, 1: 2},
					EncryptionPublicKey: [32]byte{0: 7, 1: 7},
				},
			},
		}
	)
	lggr := logger.Test(t)

	// setup registry and add one capability
	setupResp := kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{
		P2pToCapabilities: capMap,
		NopToNodes:        nopToNodes,
	})
	registry := setupResp.Registry
	chain := setupResp.Chain

	info, err := registry.GetNode(&bind.CallOpts{}, testPeerID(t, "peerID_1"))
	require.NoError(t, err)
	require.Len(t, info.HashedCapabilityIds, 1)
	// define the new capabilities that should be appended and ensure they are merged with the existing ones
	newCaps := []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:   "cap2",
			Version:        "1.0.1",
			CapabilityType: 0,
		},
		{
			LabelledName:   "cap3",
			Version:        "1.0.2",
			CapabilityType: 0,
		},
	}
	appendedResp, err := kslib.AppendCapabilities(lggr, registry, chain, []p2pkey.PeerID{testPeerID(t, "peerID_1")}, newCaps)
	require.NoError(t, err)
	require.Len(t, appendedResp, 1)
	gotCaps := appendedResp[testPeerID(t, "peerID_1")]
	require.Len(t, gotCaps, 3)
	wantCaps := capMap[testPeerID(t, "peerID_1")]
	wantCaps = append(wantCaps, newCaps...)

	for i, got := range gotCaps {
		assert.Equal(t, kslib.CapabilityID(wantCaps[i]), kslib.CapabilityID(got))
	}

	// trying to append an existing capability should not change the result
	appendedResp2, err := kslib.AppendCapabilities(lggr, registry, chain, []p2pkey.PeerID{testPeerID(t, "peerID_1")}, newCaps)
	require.NoError(t, err)
	require.Len(t, appendedResp2, 1)
	gotCaps2 := appendedResp2[testPeerID(t, "peerID_1")]
	require.Len(t, gotCaps2, 3)
	require.EqualValues(t, gotCaps, gotCaps2)

}

func testPeerID(t *testing.T, s string) p2pkey.PeerID {
	var out [32]byte
	b := []byte(s)
	copy(out[:], b)
	return p2pkey.PeerID(out)
}

func testChain(t *testing.T) deployment.Chain {
	chains := memory.NewMemoryChains(t, 1)
	var chain deployment.Chain
	for _, c := range chains {
		chain = c
		break
	}
	require.NotEmpty(t, chain)
	return chain
}

func testNop(t *testing.T, name string) kcr.CapabilitiesRegistryNodeOperator {
	return kcr.CapabilitiesRegistryNodeOperator{
		Admin: common.HexToAddress("0xFFFFFFFF45297A703e4508186d4C1aa1BAf80000"),
		Name:  name,
	}
}
