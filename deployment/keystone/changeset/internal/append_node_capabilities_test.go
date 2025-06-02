package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAppendNodeCapabilities(t *testing.T) {
	var (
		initialp2pToCapabilities = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
			testPeerID(t, "0x1"): {
				{
					LabelledName:   "test",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
			},
		}
		nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
			testNop(t, "testNop"): {
				{
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
		req          *internal.AppendNodeCapabilitiesRequest
		initialState *kstest.SetupTestRegistryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    cldf.ChangesetOutput
		wantErr bool
	}{
		{
			name: "invalid request",
			args: args{
				lggr: lggr,
				req: &internal.AppendNodeCapabilitiesRequest{
					Chain: cldf_evm.Chain{},
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
				req: &internal.AppendNodeCapabilitiesRequest{
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
				},
			},
			want:    cldf.ChangesetOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupResp := kstest.SetupTestRegistry(t, lggr, tt.args.initialState)

			tt.args.req.Chain = setupResp.Chain
			tt.args.req.CapabilitiesRegistry = setupResp.CapabilitiesRegistry

			got, err := internal.AppendNodeCapabilitiesImpl(tt.args.lggr, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("internal.AppendNodeCapabilities() error = %v, wantErr %v", err, tt.wantErr)
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
