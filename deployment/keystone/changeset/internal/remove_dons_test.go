package internal_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func Test_RemoveDONsRequest_validate(t *testing.T) {
	type fields struct {
		DONs                 []uint32
		chain                cldf_evm.Chain
		capabilitiesRegistry *kcr.CapabilitiesRegistry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "missing capabilities registry",
			fields: fields{
				DONs:                 []uint32{1},
				chain:                cldf_evm.Chain{},
				capabilitiesRegistry: nil,
			},
			wantErr: true,
		},
		{
			name: "empty DONs list",
			fields: fields{
				DONs:                 []uint32{},
				chain:                cldf_evm.Chain{},
				capabilitiesRegistry: &kcr.CapabilitiesRegistry{},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				DONs:                 []uint32{1},
				chain:                cldf_evm.Chain{},
				capabilitiesRegistry: &kcr.CapabilitiesRegistry{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &internal.RemoveDONsRequest{
				DONs:                 tt.fields.DONs,
				Chain:                tt.fields.chain,
				CapabilitiesRegistry: tt.fields.capabilitiesRegistry,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("internal.RemoveDONsRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemoveDONs tests the RemoveDONs function
func TestRemoveDONs(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	var (
		p2p1     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(100))
		pubKey1  = "11114981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e1098e7" // valid csa key
		admin1   = common.HexToAddress("0x1111567890123456789012345678901234567890")  // valid eth address
		signing1 = "11117293a4cc2621b61193135a95928735e4795f"                         // valid eth address
		node1    = newNode(t, minimalNodeCfg{
			id:            "test node 1",
			pubKey:        pubKey1,
			registryChain: registryChain,
			p2p:           p2p1,
			signingAddr:   signing1,
			admin:         admin1,
		})

		p2p2     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(200))
		pubKey2  = "22224981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109000" // valid csa key
		admin2   = common.HexToAddress("0x2222567890123456789012345678901234567891")  // valid eth address
		signing2 = "22227293a4cc2621b61193135a95928735e4ffff"                         // valid eth address
		node2    = newNode(t, minimalNodeCfg{
			id:            "test node 2",
			pubKey:        pubKey2,
			registryChain: registryChain,
			p2p:           p2p2,
			signingAddr:   signing2,
			admin:         admin2,
		})

		p2p3     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(300))
		pubKey3  = "33334981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109111" // valid csa key
		admin3   = common.HexToAddress("0x3333567890123456789012345678901234567892")  // valid eth address
		signing3 = "33337293a4cc2621b61193135a959287aaaaffff"                         // valid eth address
		node3    = newNode(t, minimalNodeCfg{
			id:            "test node 3",
			pubKey:        pubKey3,
			registryChain: registryChain,
			p2p:           p2p3,
			signingAddr:   signing3,
			admin:         admin3,
		})

		p2p4     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(400))
		pubKey4  = "44444981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109222" // valid csa key
		admin4   = common.HexToAddress("0x4444567890123456789012345678901234567893")  // valid eth address
		signing4 = "44447293a4cc2621b61193135a959287aaaaffff"                         // valid eth address
		node4    = newNode(t, minimalNodeCfg{
			id:            "test node 4",
			pubKey:        pubKey4,
			registryChain: registryChain,
			p2p:           p2p4,
			signingAddr:   signing4,
			admin:         admin4,
		})

		initialCap = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "test",
			Version:        "1.0.0",
			CapabilityType: 0,
		}

		initialCapCfg = kstest.GetDefaultCapConfig(t, initialCap)
	)

	setupEnv := func(t *testing.T) (cldf_evm.Chain, *kcr.CapabilitiesRegistry, []kcr.CapabilitiesRegistryDONInfo) {
		// 1) Set up chain + registry
		cfg := setupUpdateDonTestConfig{
			dons: []internal.DonInfo{
				{
					Name:         "don 1",
					Nodes:        []deployment.Node{node1, node2, node3, node4},
					Capabilities: []internal.DONCapabilityWithConfig{{Capability: initialCap, Config: initialCapCfg}},
				},
			},
			nops: []internal.NOP{
				{
					Name:  "nop 1",
					Nodes: []string{node1.NodeID, node2.NodeID, node3.NodeID, node4.NodeID},
				},
			},
		}

		testResp := registerTestDon(t, lggr, cfg)

		// 2) Grab the existing DONs in the test registry
		dons, err := testResp.CapabilitiesRegistry.GetDONs(&bind.CallOpts{})
		require.NoError(t, err)
		require.NotZero(t, dons, "failed to find newly created DON")

		return testResp.Chain, testResp.CapabilitiesRegistry, dons
	}

	tests := []struct {
		name            string
		useMCMS         bool
		wantErr         bool
		validateRemoval bool // if true, we will confirm on-chain that the DON is removed
		donToRemove     *uint32
	}{
		{
			name:            "remove single DON - no MCMS",
			useMCMS:         false,
			wantErr:         false,
			validateRemoval: true, // we finalize on-chain => check removal
		},
		{
			name:    "remove single DON - with MCMS",
			useMCMS: true,
			wantErr: false,
			// In MCMS mode, the transaction is batched and not finalized on-chain,
			// so skip on-chain validation that it was actually removed.
			validateRemoval: false,
		},
		{
			name:            "remove - error if DON not found in registry",
			useMCMS:         false,
			wantErr:         true,
			validateRemoval: false,
			donToRemove:     pointer.To(uint32(2839748937)),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chain, registry, dons := setupEnv(t)

			donIDs := []uint32{dons[0].Id}

			// Build the RemoveDONs request
			req := &internal.RemoveDONsRequest{
				Chain:                chain,
				CapabilitiesRegistry: registry,
				DONs:                 donIDs,
				UseMCMS:              tt.useMCMS,
			}
			if tt.useMCMS {
				req.Ops = &mcmstypes.BatchOperation{}
			}
			if tt.donToRemove != nil {
				req.DONs = []uint32{*tt.donToRemove}
			}

			// Invoke RemoveDONs
			resp, err := internal.RemoveDONs(lggr, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err, "RemoveDONs should not fail in this scenario")
			require.NotNil(t, resp)
			assert.NotEmpty(t, resp.TxHash, "Expected a valid TxHash")

			if tt.useMCMS {
				require.NotNil(t, resp.Ops, "Ops should be set in MCMS mode")
				require.NotEmpty(t, resp.Ops.Transactions, "Expected at least one transaction in MCMS batch")
			} else {
				assert.Empty(t, resp.Ops.Transactions, "No MCMS operation should be returned if not in MCMS mode")
			}

			// Optionally confirm the removal on-chain if this test scenario does a direct finalization
			if tt.validateRemoval {
				don, err := registry.GetDON(&bind.CallOpts{}, dons[0].Id)
				require.NoError(t, err)
				require.NotEqual(t, dons[0].Id, don.Id, "Expect donID to not be found")
				require.Equal(t, uint32(0), don.Id, "Expected returned donID to be 0")
			}
		})
	}
}
