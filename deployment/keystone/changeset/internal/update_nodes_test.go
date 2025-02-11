package internal_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func Test_UpdateNodesRequest_validate(t *testing.T) {
	type fields struct {
		p2pToUpdates         map[p2pkey.PeerID]internal.NodeUpdate
		nopToNodes           map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc
		chain                deployment.Chain
		capabilitiesRegistry *kcr.CapabilitiesRegistry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "err",
			fields: fields{
				p2pToUpdates:         map[p2pkey.PeerID]internal.NodeUpdate{},
				nopToNodes:           nil,
				chain:                deployment.Chain{},
				capabilitiesRegistry: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid encryption key -- cannot decode",
			fields: fields{
				p2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
					p2pkey.PeerID{}: {
						EncryptionPublicKey: "jk",
					},
				},
				nopToNodes:           nil,
				chain:                deployment.Chain{},
				capabilitiesRegistry: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid encryption key -- invalid length",
			fields: fields{
				p2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
					testPeerID(t, "peerID_1"): {
						EncryptionPublicKey: "aabb",
					},
				},
				nopToNodes:           nil,
				chain:                deployment.Chain{},
				capabilitiesRegistry: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &internal.UpdateNodesRequest{
				P2pToUpdates:         tt.fields.p2pToUpdates,
				Chain:                tt.fields.chain,
				CapabilitiesRegistry: tt.fields.capabilitiesRegistry,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("internal.UpdateNodesRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newEncryptionKey(t *testing.T) [32]byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return [32]byte(key)
}

func TestUpdateNodes(t *testing.T) {
	chain := testChain(t)
	require.NotNil(t, chain)
	lggr := logger.Test(t)
	newKey := newEncryptionKey(t)
	newKeyStr := hex.EncodeToString(newKey[:])

	type args struct {
		lggr        logger.Logger
		req         *internal.UpdateNodesRequest
		nopsToNodes map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc
	}
	tests := []struct {
		name    string
		args    args
		want    *internal.UpdateNodesResponse
		wantErr bool
	}{
		{
			name: "one node, one capability",
			args: args{
				lggr: lggr,
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							Capabilities: []kcr.CapabilitiesRegistryCapability{
								{
									LabelledName:   "cap1",
									Version:        "1.0.0",
									CapabilityType: 0,
								},
							},
						},
					},
					Chain: chain,
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nop1"): []*internal.P2PSignerEnc{
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 1: 2},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
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
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): internal.NodeUpdate{
							Capabilities: []kcr.CapabilitiesRegistryCapability{
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
					},
					Chain: chain,
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nop1"): []*internal.P2PSignerEnc{
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 1: 2},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
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
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							NodeOperatorID: 1,
							Capabilities: []kcr.CapabilitiesRegistryCapability{
								{
									LabelledName:   "cap1",
									Version:        "1.0.0",
									CapabilityType: 0,
								},
							},
						},
						testPeerID(t, "peerID_2"): {
							NodeOperatorID: 2,
							Capabilities: []kcr.CapabilitiesRegistryCapability{
								{
									LabelledName:   "cap1",
									Version:        "1.0.0",
									CapabilityType: 0,
								},
							},
						},
					},
					Chain: chain,
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nopA"): {
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 31: 1},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
					testNop(t, "nopB"): {
						{
							P2PKey:              testPeerID(t, "peerID_2"),
							Signer:              [32]byte{0: 2, 31: 2},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
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
			name: "twos nodes with different capabilities",
			args: args{
				lggr: lggr,
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							NodeOperatorID: 1,
							Capabilities: []kcr.CapabilitiesRegistryCapability{
								{
									LabelledName:   "cap1",
									Version:        "1.0.0",
									CapabilityType: 0,
								},
							},
						},
						testPeerID(t, "peerID_2"): {
							NodeOperatorID: 2,
							Capabilities: []kcr.CapabilitiesRegistryCapability{
								{
									LabelledName:   "cap2",
									Version:        "1.0.1",
									CapabilityType: 0,
								},
							},
						},
					},
					Chain: chain,
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nopA"): {
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 31: 1},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
					testNop(t, "nopB"): {
						{
							P2PKey:              testPeerID(t, "peerID_2"),
							Signer:              [32]byte{0: 2, 31: 2},
							EncryptionPublicKey: [32]byte{0: 7, 1: 7},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
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
			name: "one node, updated encryption key",
			args: args{
				lggr: lggr,
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							EncryptionPublicKey: newKeyStr,
						},
					},
					Chain:                chain,
					CapabilitiesRegistry: nil, // set in test to ensure no conflicts
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nop1"): []*internal.P2PSignerEnc{
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 1: 2},
							EncryptionPublicKey: [32]byte{0: 1, 1: 2},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: newKey,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one node, updated signer",
			args: args{
				lggr: lggr,
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							Signer: [32]byte{0: 2, 1: 3},
						},
					},
					Chain:                chain,
					CapabilitiesRegistry: nil, // set in test to ensure no conflicts
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nop1"): []*internal.P2PSignerEnc{
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 1: 2},
							EncryptionPublicKey: [32]byte{0: 1, 1: 2},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "peerID_1"),
						Signer:              [32]byte{0: 2, 1: 3},
						EncryptionPublicKey: [32]byte{0: 1, 1: 2},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one node, updated nodeOperatorID",
			args: args{
				lggr: lggr,
				req: &internal.UpdateNodesRequest{
					P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
						testPeerID(t, "peerID_1"): {
							NodeOperatorID: 2,
						},
					},
					Chain: chain,
				},
				nopsToNodes: map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
					testNop(t, "nop1"): []*internal.P2PSignerEnc{
						{
							P2PKey:              testPeerID(t, "peerID_1"),
							Signer:              [32]byte{0: 1, 1: 2},
							EncryptionPublicKey: [32]byte{0: 1, 1: 2},
						},
					},
				},
			},
			want: &internal.UpdateNodesResponse{
				NodeParams: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      2,
						P2pId:               testPeerID(t, "peerID_1"),
						Signer:              [32]byte{0: 1, 1: 2},
						EncryptionPublicKey: [32]byte{0: 1, 1: 2},
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
			phonyCapCfg := &capabilitiespb.CapabilityConfig{
				DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
				RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh: durationpb.New(20 * time.Second),
						RegistrationExpiry:  durationpb.New(60 * time.Second),
						// F + 1; assuming n = 3f+1
						MinResponsesToAggregate: uint32(10),
					},
				},
			}
			initMap := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
			for p2pID := range tt.args.req.P2pToUpdates {
				initMap[p2pID] = []kcr.CapabilitiesRegistryCapability{phonyCap}
			}
			setupResp := kstest.SetupTestRegistry(t, tt.args.lggr, &kstest.SetupTestRegistryRequest{
				P2pToCapabilities: initMap,
				NopToNodes:        tt.args.nopsToNodes,
			})
			registry := setupResp.CapabilitiesRegistry
			tt.args.req.CapabilitiesRegistry = setupResp.CapabilitiesRegistry
			tt.args.req.Chain = setupResp.Chain

			id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, phonyCap.LabelledName, phonyCap.Version)
			require.NoError(t, err)

			// register the capabilities that the Update will use
			expectedUpdatedCaps := make(map[p2pkey.PeerID][]internal.RegisteredCapability)
			capCache := kstest.NewCapabiltyCache(t, registry)
			for p2p, update := range tt.args.req.P2pToUpdates {
				if len(update.Capabilities) > 0 {
					expectedCaps := capCache.AddCapabilities(tt.args.lggr, tt.args.req.Chain, registry, update.Capabilities)
					expectedUpdatedCaps[p2p] = expectedCaps
				} else {
					expectedUpdatedCaps[p2p] = []internal.RegisteredCapability{
						{CapabilitiesRegistryCapability: phonyCap, ID: id, Config: phonyCapCfg},
					}
				}
			}
			got, err := internal.UpdateNodes(tt.args.lggr, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, p := range got.NodeParams {
				// check the node params
				expected := findNodeParams(t, tt.want.NodeParams, p.P2pId)
				assertNodeParams(t, expected, p)

				// check the capabilities
				expectedCaps := expectedUpdatedCaps[p.P2pId]
				var wantHashedIDs [][32]byte
				for _, cap := range expectedCaps {
					wantHashedIDs = append(wantHashedIDs, cap.ID)
				}
				sort.Slice(wantHashedIDs, func(i, j int) bool {
					return bytes.Compare(wantHashedIDs[i][:], wantHashedIDs[j][:]) < 0
				})
				gotHashedIDs := p.HashedCapabilityIds
				sort.Slice(gotHashedIDs, func(i, j int) bool {
					return bytes.Compare(gotHashedIDs[i][:], gotHashedIDs[j][:]) < 0
				})
				require.Len(t, gotHashedIDs, len(wantHashedIDs))
				for j, gotCap := range gotHashedIDs {
					assert.Equal(t, wantHashedIDs[j], gotCap)
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
			nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
				testNop(t, "nopA"): []*internal.P2PSignerEnc{
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
		registry := setupResp.CapabilitiesRegistry
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
			err2 := deployment.DecodeErr(kcr.CapabilitiesRegistryABI, err)
			require.Fail(t, fmt.Sprintf("failed to call AddCapabilities: %s:  %s", err, err2))
		}
		_, err = chain.Confirm(tx)
		require.NoError(t, err)

		var req = &internal.UpdateNodesRequest{
			P2pToUpdates: map[p2pkey.PeerID]internal.NodeUpdate{
				testPeerID(t, "peerID_1"): internal.NodeUpdate{
					Capabilities: toRegister,
				},
			},
			Chain:                chain,
			CapabilitiesRegistry: setupResp.CapabilitiesRegistry,
		}
		_, err = internal.UpdateNodes(lggr, req)
		require.NoError(t, err)
		info, err = registry.GetNode(&bind.CallOpts{}, testPeerID(t, "peerID_1"))
		require.NoError(t, err)
		require.Len(t, info.HashedCapabilityIds, 1)
		want := info.HashedCapabilityIds[0]

		// update again and ensure the result is the same
		_, err = internal.UpdateNodes(lggr, req)
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
		nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
			testNop(t, "nop"): []*internal.P2PSignerEnc{
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
	registry := setupResp.CapabilitiesRegistry
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
	appendedResp, err := internal.AppendCapabilities(lggr, registry, chain, []p2pkey.PeerID{testPeerID(t, "peerID_1")}, newCaps)
	require.NoError(t, err)
	require.Len(t, appendedResp, 1)
	gotCaps := appendedResp[testPeerID(t, "peerID_1")]
	require.Len(t, gotCaps, 3)
	wantCaps := capMap[testPeerID(t, "peerID_1")]
	wantCaps = append(wantCaps, newCaps...)

	for i, got := range gotCaps {
		assert.Equal(t, internal.CapabilityID(wantCaps[i]), internal.CapabilityID(got))
	}

	// trying to append an existing capability should not change the result
	appendedResp2, err := internal.AppendCapabilities(lggr, registry, chain, []p2pkey.PeerID{testPeerID(t, "peerID_1")}, newCaps)
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
	chains, _ := memory.NewMemoryChains(t, 1, 5)
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

func findNodeParams(t *testing.T, nodes []kcr.CapabilitiesRegistryNodeParams, p2p p2pkey.PeerID) kcr.CapabilitiesRegistryNodeParams {
	for _, n := range nodes {
		if n.P2pId == p2p {
			return n
		}
	}
	require.Failf(t, "could not find node %s", p2p.String())
	return kcr.CapabilitiesRegistryNodeParams{}
}

func assertNodeParams(t *testing.T, expected, got kcr.CapabilitiesRegistryNodeParams) {
	t.Helper()
	assert.Equal(t, expected.P2pId, got.P2pId, "p2p ID failed : expected %v, got %v", expected.P2pId, got.P2pId)
	assert.Equal(t, expected.NodeOperatorId, got.NodeOperatorId, "nop ID failed : expected %d, got %d", expected.NodeOperatorId, got.NodeOperatorId)
	assert.Equal(t, expected.Signer, got.Signer, "signer failed : expected %v, got %v", expected.Signer, got.Signer)
	assert.Equal(t, expected.EncryptionPublicKey, got.EncryptionPublicKey, "encryption key failed : expected %v, got %v", expected.EncryptionPublicKey, got.EncryptionPublicKey)
}
