package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/changeset"
	kstest "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAppendNodeCapabilities(t *testing.T) {
	var (
		initialp2pToCapabilities = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
			testPeerID(t, "0x1"): []kcr.CapabilitiesRegistryCapability{
				{
					LabelledName:   "test",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
			},
		}
		nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc{
			testNop(t, "testNop"): []*kslib.P2PSignerEnc{
				&kslib.P2PSignerEnc{
					Signer:              [32]byte{0: 1},
					P2PKey:              testPeerID(t, "0x1"),
					EncryptionPublicKey: [32]byte{7: 7, 13: 13},
				},
			},
		}
	)

	lggr := logger.Test(t)

	type args struct {
		lggr         logger.Logger
		req          *changeset.AppendNodeCapabilitiesRequest
		initialState *kstest.SetupTestRegistryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    deployment.ChangesetOutput
		wantErr bool
	}{
		{
			name: "invalid request",
			args: args{
				lggr: lggr,
				req: &changeset.AppendNodeCapabilitiesRequest{
					Chain: deployment.Chain{},
				},
				initialState: &kstest.SetupTestRegistryRequest{},
			},
			wantErr: true,
		},
		{
			name: "happy path",
			args: args{
				lggr: lggr,
				initialState: &kstest.SetupTestRegistryRequest{
					P2pToCapabilities: initialp2pToCapabilities,
					NopToNodes:        nopToNodes,
				},
				req: &changeset.AppendNodeCapabilitiesRequest{
					P2pToCapabilities: map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
						testPeerID(t, "0x1"): []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "cap2",
								Version:        "1.0.0",
								CapabilityType: 0,
							},
							{
								LabelledName:   "cap3",
								Version:        "1.0.0",
								CapabilityType: 3,
							},
						},
					},
					NopToNodes: nopToNodes,
				},
			},
			want:    deployment.ChangesetOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// chagen the name and args to be mor egeneral
			setupResp := kstest.SetupTestRegistry(t, lggr, tt.args.initialState)

			tt.args.req.Registry = setupResp.Registry
			tt.args.req.Chain = setupResp.Chain

			got, err := changeset.AppendNodeCapabilitiesImpl(tt.args.lggr, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendNodeCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			require.NotNil(t, got)
			// should be one node param for each input p2p id
			assert.Len(t, got.NodeParams, len(tt.args.req.P2pToCapabilities))
			for _, nodeParam := range got.NodeParams {
				initialCapsOnNode := tt.args.initialState.P2pToCapabilities[nodeParam.P2pId]
				appendCaps := tt.args.req.P2pToCapabilities[nodeParam.P2pId]
				assert.Len(t, nodeParam.HashedCapabilityIds, len(initialCapsOnNode)+len(appendCaps))
			}
		})
	}
}

func testPeerID(t *testing.T, s string) p2pkey.PeerID {
	var out [32]byte
	b := []byte(s)
	copy(out[:], b)
	return p2pkey.PeerID(out)
}

func testNop(t *testing.T, name string) kcr.CapabilitiesRegistryNodeOperator {
	return kcr.CapabilitiesRegistryNodeOperator{
		Admin: common.HexToAddress("0xFFFFFFFF45297A703e4508186d4C1aa1BAf80000"),
		Name:  name,
	}
}
