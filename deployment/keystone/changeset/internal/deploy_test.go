package internal_test

import (
	"context"
	"maps"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func Test_RegisterNOPS(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.CapabilitiesRegistry
		chain     = setupResp.Chain
		nops      = make([]kcr.CapabilitiesRegistryNodeOperator, 0)
	)
	t.Run("success create add NOPs mcms proposal", func(t *testing.T) {
		nops = append(nops, kcr.CapabilitiesRegistryNodeOperator{
			Name: "test-nop",
		})
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterNOPS(context.TODO(), lggr, internal.RegisterNOPSRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			Nops:                  nops,
			UseMCMS:               useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
	})
}

func Test_AddCapabilities(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.CapabilitiesRegistry
		chain     = setupResp.Chain
	)

	assertExists := func(t *testing.T, registry *kcr.CapabilitiesRegistry, capabilities ...kcr.CapabilitiesRegistryCapability) {
		for _, capability := range capabilities {
			wantID, err := registry.GetHashedCapabilityId(nil, capability.LabelledName, capability.Version)
			require.NoError(t, err)
			got, err := registry.GetCapability(nil, wantID)
			require.NoError(t, err)
			require.NotEmpty(t, got)
			assert.Equal(t, capability.CapabilityType, got.CapabilityType)
			assert.Equal(t, capability.LabelledName, got.LabelledName)
			assert.Equal(t, capability.Version, got.Version)
		}
	}
	t.Run("successfully create mcms proposal", func(t *testing.T) {
		useMCMS = true
		capabilities := []kcr.CapabilitiesRegistryCapability{kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cap1",
			Version:        "1.0.0",
			CapabilityType: 0,
		}}
		ops, err := internal.AddCapabilities(lggr, registry, chain, capabilities, useMCMS)
		require.NoError(t, err)
		require.NotNil(t, ops)
		require.Len(t, ops.Batch, 1)
	})

	t.Run("does nothing if no capabilities", func(t *testing.T) {
		ops, err := internal.AddCapabilities(lggr, registry, chain, nil, useMCMS)
		require.NoError(t, err)
		require.Nil(t, ops)
	})

	t.Run("idempotent", func(t *testing.T) {
		capabilities := []kcr.CapabilitiesRegistryCapability{kcr.CapabilitiesRegistryCapability{
			LabelledName:   "idempotent-cap",
			Version:        "1.0.0",
			CapabilityType: 0,
		},
		}
		_, err := internal.AddCapabilities(lggr, registry, chain, capabilities, false)
		require.NoError(t, err)
		assertExists(t, registry, capabilities...)
		_, err = internal.AddCapabilities(lggr, registry, chain, capabilities, false)
		require.NoError(t, err)
		assertExists(t, registry, capabilities...)
	})

	t.Run("dedup", func(t *testing.T) {
		capabilities := []kcr.CapabilitiesRegistryCapability{
			{
				LabelledName:   "would-be-duplicate",
				Version:        "1.0.0",
				CapabilityType: 0,
			},
			{
				LabelledName:   "would-be-duplicate",
				Version:        "1.0.0",
				CapabilityType: 0,
			},
		}
		_, err := internal.AddCapabilities(lggr, registry, chain, capabilities, false)
		require.NoError(t, err)
		assertExists(t, registry, capabilities[0])
	})
}

func Test_RegisterNodes(t *testing.T) {
	var (
		useMCMS                  bool
		lggr                     = logger.Test(t)
		existingNOP              = testNop(t, "testNop")
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
			existingNOP: {
				{
					Signer:              [32]byte{0: 1},
					P2PKey:              testPeerID(t, "0x1"),
					EncryptionPublicKey: [32]byte{3: 16, 4: 2},
				},
			},
		}

		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{
			P2pToCapabilities: initialp2pToCapabilities,
			NopToNodes:        nopToNodes,
		})
		registry = setupResp.CapabilitiesRegistry
		chain    = setupResp.Chain

		registeredCapabilities = kstest.GetRegisteredCapabilities(t, lggr, initialp2pToCapabilities, setupResp.CapabilityCache)

		registeredNodeParams = kstest.ToNodeParams(t, nopToNodes,
			kstest.ToP2PToCapabilities(t, initialp2pToCapabilities, registry, registeredCapabilities),
		)
	)

	t.Run("success create add nodes", func(t *testing.T) {
		type input struct {
			NodeOperatorID      uint32
			Signer              [32]byte
			P2PID               [32]byte
			EncryptionPublicKey [32]byte
			Capabilities        []kcr.CapabilitiesRegistryCapability
		}

		type expected struct {
			nodes []kcr.CapabilitiesRegistryNodeParams
			nOps  int
		}
		type args struct {
			name    string
			useMCMS bool
			want    expected
			input   input
		}
		var (
			newCap = kcr.CapabilitiesRegistryCapability{
				LabelledName:   "new-cap",
				Version:        "1.0.0",
				CapabilityType: 0,
			}

			testInput = input{
				NodeOperatorID:      1, // this already exists in the test registry
				Signer:              [32]byte{0: 27},
				P2PID:               testPeerID(t, "0x2"),
				EncryptionPublicKey: [32]byte{3: 16, 4: 2},
				Capabilities:        []kcr.CapabilitiesRegistryCapability{newCap},
			}
		)
		wantHash, err := registry.GetHashedCapabilityId(nil, newCap.LabelledName, newCap.Version)
		require.NoError(t, err)
		var cases = []args{
			{
				name:    "mcms",
				useMCMS: true,
				want:    expected{nOps: 1},
				input:   testInput,
			},

			{
				name:    "no mcms",
				useMCMS: false,
				input:   testInput,
				want: expected{nodes: []kcr.CapabilitiesRegistryNodeParams{
					{
						NodeOperatorId:      1,
						P2pId:               testPeerID(t, "0x2"),
						Signer:              [32]byte{0: 27},
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
						HashedCapabilityIds: [][32]byte{wantHash},
					},
				}},
			},
		}

		for _, tc := range cases {
			caps2Add := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
			caps2Add[tc.input.P2PID] = append(caps2Add[tc.input.P2PID], testInput.Capabilities...)
			rc, _ := kstest.MustAddCapabilities(t, lggr, caps2Add, chain, registry)

			t.Run(tc.name, func(t *testing.T) {
				env := &deployment.Environment{
					Logger: lggr,
					Chains: map[uint64]deployment.Chain{
						chain.Selector: chain,
					},
					ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
						chain.Selector: {
							registry.Address().String(): deployment.TypeAndVersion{
								Type:    internal.CapabilitiesRegistry,
								Version: deployment.Version1_0_0,
							},
						},
					}),
				}
				resp, err := internal.RegisterNodes(lggr, &internal.RegisterNodesRequest{
					Env:                   env,
					RegistryChainSelector: chain.Selector,
					UseMCMS:               tc.useMCMS, // useMCMS,
					DonToCapabilities: map[string][]internal.RegisteredCapability{
						"testDON": rc,
					},
					NopToNodeIDs: map[kcr.CapabilitiesRegistryNodeOperator][]string{
						existingNOP: {"node-id"},
					},
					DonToNodes: map[string][]deployment.Node{
						"testDON": {
							{
								PeerID: tc.input.P2PID,
								NodeID: "node-id",
								SelToOCRConfig: map[chain_selectors.ChainDetails]deployment.OCRConfig{
									{
										ChainSelector: chain.Selector,
									}: {
										OnchainPublicKey:          tc.input.Signer[:],
										ConfigEncryptionPublicKey: tc.input.EncryptionPublicKey,
									},
								},
							},
						},
					},
					Nops: []*kcr.CapabilitiesRegistryNodeOperatorAdded{{
						Name:           existingNOP.Name,
						Admin:          existingNOP.Admin,
						NodeOperatorId: 1,
					}},
				})
				require.NoError(t, err)
				if tc.useMCMS {
					require.NotNil(t, resp.Ops)
					require.Len(t, resp.Ops.Batch, tc.want.nOps)
				} else {
					assertNodesExist(t, registry, tc.want.nodes...)
				}
			})
		}
	})

	t.Run("no ops in proposal if node already exists", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterNodes(lggr, &internal.RegisterNodesRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			UseMCMS:               useMCMS,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"testDON": registeredCapabilities,
			},
			NopToNodeIDs: map[kcr.CapabilitiesRegistryNodeOperator][]string{
				existingNOP: {"node-id"},
			},
			DonToNodes: map[string][]deployment.Node{
				"testDON": {
					{
						PeerID: registeredNodeParams[0].P2pId,
						NodeID: "node-id",
						SelToOCRConfig: map[chain_selectors.ChainDetails]deployment.OCRConfig{
							{
								ChainSelector: chain.Selector,
							}: {},
						},
					},
				},
			},
			Nops: []*kcr.CapabilitiesRegistryNodeOperatorAdded{{
				Name:           existingNOP.Name,
				Admin:          existingNOP.Admin,
				NodeOperatorId: 1,
			}},
		})
		require.NoError(t, err)
		require.Nil(t, resp.Ops)
	})

	t.Run("no new nodes to add results in no mcms ops", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterNodes(lggr, &internal.RegisterNodesRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			UseMCMS:               useMCMS,
		})
		require.NoError(t, err)
		require.Nil(t, resp.Ops)
	})
}

func assertNodesExist(t *testing.T, registry *kcr.CapabilitiesRegistry, nodes ...kcr.CapabilitiesRegistryNodeParams) {
	for _, node := range nodes {
		got, err := registry.GetNode(nil, node.P2pId)
		require.NoError(t, err)
		assert.Equal(t, node.EncryptionPublicKey, got.EncryptionPublicKey, "mismatch node encryption public key")
		assert.Equal(t, node.Signer, got.Signer, "mismatch node signer")
		assert.Equal(t, node.NodeOperatorId, got.NodeOperatorId, "mismatch node operator id")
		assert.Equal(t, node.HashedCapabilityIds, got.HashedCapabilityIds, "mismatch node hashed capability ids")
		assert.Equal(t, node.P2pId, got.P2pId, "mismatch node p2p id")
	}
}

func Test_RegisterDons(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.CapabilitiesRegistry
		chain     = setupResp.Chain
	)
	t.Run("success create add DONs mcms proposal", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterDons(lggr, internal.RegisterDonsRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"test-don": {},
			},
			DonsToRegister: []internal.DONToRegister{
				{
					Name: "test-don",
					F:    2,
				},
			},
			UseMCMS: useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
	})

	t.Run("no new dons to add results in no mcms ops", func(t *testing.T) {
		var (
			existingNOP              = testNop(t, "testNop")
			existingP2Pkey           = testPeerID(t, "0x1")
			initialp2pToCapabilities = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
				existingP2Pkey: {
					{
						LabelledName:   "test",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
				},
				testPeerID(t, "0x2"): {
					{
						LabelledName:   "test",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
				},
				testPeerID(t, "0x3"): {
					{
						LabelledName:   "test",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
				},
			}
			nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
				existingNOP: {
					{
						Signer:              [32]byte{0: 1},
						P2PKey:              existingP2Pkey,
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					},
					{
						Signer:              [32]byte{0: 1, 1: 1},
						P2PKey:              testPeerID(t, "0x2"),
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					},
					{
						Signer:              [32]byte{0: 1, 1: 1, 2: 1},
						P2PKey:              testPeerID(t, "0x3"),
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					},
				},
			}

			setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{
				P2pToCapabilities: initialp2pToCapabilities,
				NopToNodes:        nopToNodes,
				Dons: []kstest.Don{
					{
						Name:   "test-don",
						P2PIDs: []p2pkey.PeerID{existingP2Pkey, testPeerID(t, "0x2"), testPeerID(t, "0x3")},
					},
				},
			})
			regContract = setupResp.CapabilitiesRegistry
		)

		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				setupResp.Chain.Selector: setupResp.Chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				setupResp.Chain.Selector: {
					regContract.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterDons(lggr, internal.RegisterDonsRequest{
			Env:                   env,
			RegistryChainSelector: setupResp.Chain.Selector,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"test-don": {},
			},
			DonsToRegister: []internal.DONToRegister{
				{
					Name: "test-don",
					F:    1,
				},
			},
			UseMCMS: true,
		})
		require.NoError(t, err)
		require.Nil(t, resp.Ops)
	})

	t.Run("success create add DONs mcms proposal with multiple DONs", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterDons(lggr, internal.RegisterDonsRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"test-don-1": {},
				"test-don-2": {},
			},
			DonsToRegister: []internal.DONToRegister{
				{
					Name: "test-don-1",
					F:    2,
				},
				{
					Name: "test-don-2",
					F:    2,
				},
			},
			UseMCMS: useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 2)
	})
}

func TestAddNodes(t *testing.T) {
	var (
		lggr                     = logger.Test(t)
		existingNOP              = testNop(t, "testNop")
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
			existingNOP: {
				{
					Signer:              [32]byte{0: 1},
					P2PKey:              testPeerID(t, "0x1"),
					EncryptionPublicKey: [32]byte{3: 16, 4: 2},
				},
			},
		}

		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{
			P2pToCapabilities: initialp2pToCapabilities,
			NopToNodes:        nopToNodes,
		})
		registry = setupResp.CapabilitiesRegistry
		chain    = setupResp.Chain
		nops     = setupResp.Nops

		registeredCapabilities = kstest.GetRegisteredCapabilities(t, lggr, initialp2pToCapabilities, setupResp.CapabilityCache)
	)

	type args struct {
		lggr logger.Logger
		req  *internal.AddNodesRequest
	}
	tests := []struct {
		name     string
		args     args
		want     *internal.AddNodesResponse
		checkErr func(t *testing.T, err error)
		preCheck func(t *testing.T)
	}{
		{
			name: "error: missing node operator",
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"test-node": {
							P2pId:               testPeerID(t, "0xabc01"),
							Signer:              [32]byte{16: 16},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      99, // does not exist
						},
					},
				},
			},
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "NodeOperatorDoesNotExist", err.Error())
			},
		},
		{
			name: "error: missing capability",
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"test-node": {
							P2pId:               testPeerID(t, "0xabc01"),
							Signer:              [32]byte{16: 16},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{[32]byte{1: 1, 2: 1}},
						},
					},
				},
			},
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "InvalidNodeCapabilities", "%v does not contain InvalidNodeCapabilities", err.Error())
			},
		},
		{
			name: "add one node",
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"test-node": {
							P2pId:               testPeerID(t, "0xabc01"),
							Signer:              [32]byte{16: 16},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
					},
				},
			},
			want: &internal.AddNodesResponse{
				AddedNodes: []kcr.CapabilitiesRegistryNodeParams{
					{
						P2pId:               testPeerID(t, "0xabc01"),
						Signer:              [32]byte{16: 16},
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
						NodeOperatorId:      nops[0].NodeOperatorId,
						HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
					},
				},
			},
		},
		{
			name: "idempotent", // not order matters; we assume that previous test case ran
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"test-node": {
							P2pId:               testPeerID(t, "0xabc01"),
							Signer:              [32]byte{16: 16},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
					},
				},
			},
			preCheck: func(t *testing.T) {
				t.Logf("preCheck: adding the same node again should be idempotent, ensuring that the node is present in the registry")
				assertNodesExist(t, registry, kcr.CapabilitiesRegistryNodeParams{
					P2pId:               testPeerID(t, "0xabc01"),
					Signer:              [32]byte{16: 16},
					EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					NodeOperatorId:      nops[0].NodeOperatorId,
					HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
				},
				)
			},
			want: &internal.AddNodesResponse{},
		},
		{
			name: "duplicate error", // the same node is added twice
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"dupeA": {
							P2pId:               testPeerID(t, "0xffff"),
							Signer:              [32]byte{31: 31},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
						"dupeB": {
							P2pId:               testPeerID(t, "0xffff"),
							Signer:              [32]byte{31: 31},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
					},
				},
			},
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "NodeAlreadyExists", "%v does not contain NodeAlreadyExists", err.Error())
			},
		},
		{
			name: "add two nodes",
			args: args{
				lggr: lggr,
				req: &internal.AddNodesRequest{
					RegistryChain:        chain,
					CapabilitiesRegistry: registry,
					NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
						"test-node-1": {
							P2pId:               testPeerID(t, "0x1111"),
							Signer:              [32]byte{1: 16, 2: 32},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
						"test-node-2": {
							P2pId:               testPeerID(t, "0x2222"),
							Signer:              [32]byte{7: 16, 8: 32},
							EncryptionPublicKey: [32]byte{3: 16, 4: 2},
							NodeOperatorId:      nops[0].NodeOperatorId,
							HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
						},
					},
				},
			},
			want: &internal.AddNodesResponse{
				AddedNodes: []kcr.CapabilitiesRegistryNodeParams{
					{
						P2pId:               testPeerID(t, "0x1111"),
						Signer:              [32]byte{1: 16, 2: 32},
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
						NodeOperatorId:      nops[0].NodeOperatorId,
						HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
					},
					{
						P2pId:               testPeerID(t, "0x2222"),
						Signer:              [32]byte{7: 16, 8: 32},
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
						NodeOperatorId:      nops[0].NodeOperatorId,
						HashedCapabilityIds: [][32]byte{registeredCapabilities[0].ID},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := internal.AddNodes(tt.args.lggr, tt.args.req)
			if err != nil && tt.checkErr == nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.checkErr != nil {
				require.Error(t, err, "checkErr is not nil but no error was returned")
				tt.checkErr(t, err)
			}
			// if no err expected, assert that the nodes were added
			if tt.checkErr == nil {
				var vals []kcr.CapabilitiesRegistryNodeParams
				for v := range maps.Values(tt.args.req.NodeParams) {
					vals = append(vals, v)
				}
				assertNodesExist(t, registry, vals...)
				if len(tt.want.AddedNodes) > 0 {
					assertNodesExist(t, registry, tt.want.AddedNodes...)
				}
			}
		})
	}
}
