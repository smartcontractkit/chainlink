package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

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
	t.Run("success create add nodes mcms proposal", func(t *testing.T) {
		var (
			nop2Add  = testNop(t, "newNop")
			caps2Add = map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{
				testPeerID(t, "0x2"): {
					{
						LabelledName:   "new-cap",
						Version:        "1.0.0",
						CapabilityType: 0,
					},
				},
			}

			nopToNodes = map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc{
				nop2Add: {
					{
						Signer:              [32]byte{0: 1},
						P2PKey:              testPeerID(t, "0x2"),
						EncryptionPublicKey: [32]byte{3: 16, 4: 2},
					},
				},
			}

			rc, _ = kstest.MustAddCapabilities(t, lggr, caps2Add, chain, registry)

			nps = kstest.ToNodeParams(t, nopToNodes,
				kstest.ToP2PToCapabilities(t, caps2Add, registry, rc),
			)
		)

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
				"testDON": rc,
			},
			NopToNodeIDs: map[kcr.CapabilitiesRegistryNodeOperator][]string{
				nop2Add: {"node-id"},
			},
			DonToNodes: map[string][]deployment.Node{
				"testDON": {
					{
						PeerID: nps[0].P2pId,
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
				Name:           nop2Add.Name,
				Admin:          nop2Add.Admin,
				NodeOperatorId: 2,
			}},
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
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
