package internal_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment"
	kscs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var (
	registryChain = chainsel.TEST_90000001
)

func TestUpdateDon(t *testing.T) {
	var (
		// nodes
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
		// capabilities
		initialCap = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "test",
			Version:        "1.0.0",
			CapabilityType: 0,
		}

		capToAdd = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cap b",
			Version:        "1.0.0",
			CapabilityType: 1,
		}
	)

	initialCapCfg := kstest.GetDefaultCapConfig(t, initialCap)
	initialCapCfgB, err := proto.Marshal(initialCapCfg)
	require.NoError(t, err)
	capToAddCfg := kstest.GetDefaultCapConfig(t, capToAdd)
	capToAddCfgB, err := proto.Marshal(capToAddCfg)
	require.NoError(t, err)

	lggr := logger.Test(t)

	t.Run("empty", func(t *testing.T) {
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

		testCfg := registerTestDon(t, lggr, cfg)
		// add the new capabilities to registry
		m := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
		for _, node := range cfg.dons[0].Nodes {
			m[node.PeerID] = append(m[node.PeerID], capToAdd)
		}

		_, err := internal.AppendNodeCapabilitiesImpl(lggr, &internal.AppendNodeCapabilitiesRequest{
			Chain:                testCfg.Chain,
			CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
			P2pToCapabilities:    m,
		})
		require.NoError(t, err)

		req := &internal.UpdateDonRequest{
			CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
			Chain:                testCfg.Chain,
			P2PIDs:               []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID()},
			CapabilityConfigs: []internal.CapabilityConfig{
				{Capability: initialCap, Config: initialCapCfgB}, {Capability: capToAdd, Config: capToAddCfgB},
			},
		}
		want := &internal.UpdateDonResponse{
			DonInfo: kcr.CapabilitiesRegistryDONInfo{
				Id:          1,
				ConfigCount: 1,
				NodeP2PIds:  internal.PeerIDsToBytes([]p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID()}),
				CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
					{CapabilityId: kstest.MustCapabilityID(t, testCfg.CapabilitiesRegistry, initialCap), Config: initialCapCfgB},
					{CapabilityId: kstest.MustCapabilityID(t, testCfg.CapabilitiesRegistry, capToAdd), Config: capToAddCfgB},
				},
			},
		}

		got, err := internal.UpdateDon(lggr, req)
		require.NoError(t, err)
		assert.Equal(t, want.DonInfo.Id, got.DonInfo.Id)
		assert.Equal(t, want.DonInfo.ConfigCount, got.DonInfo.ConfigCount)
		assert.Equal(t, sortedP2Pids(want.DonInfo.NodeP2PIds), sortedP2Pids(got.DonInfo.NodeP2PIds))
		assert.Equal(t, capIDs(want.DonInfo.CapabilityConfigurations), capIDs(got.DonInfo.CapabilityConfigurations))
	})
}

func TestUpdateDon_ChangeComposition(t *testing.T) {
	var (
		// Initial nodes (n1-n4)
		p2p1     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(100))
		pubKey1  = "11114981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e1098e7"
		admin1   = common.HexToAddress("0x1111567890123456789012345678901234567890")
		signing1 = "11117293a4cc2621b61193135a95928735e4795f"
		node1    = newNode(t, minimalNodeCfg{
			id:            "test node 1",
			pubKey:        pubKey1,
			registryChain: registryChain,
			p2p:           p2p1,
			signingAddr:   signing1,
			admin:         admin1,
		})

		p2p2     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(200))
		pubKey2  = "22224981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109000"
		admin2   = common.HexToAddress("0x2222567890123456789012345678901234567891")
		signing2 = "22227293a4cc2621b61193135a95928735e4ffff"
		node2    = newNode(t, minimalNodeCfg{
			id:            "test node 2",
			pubKey:        pubKey2,
			registryChain: registryChain,
			p2p:           p2p2,
			signingAddr:   signing2,
			admin:         admin2,
		})

		p2p3     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(300))
		pubKey3  = "33334981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109111"
		admin3   = common.HexToAddress("0x3333567890123456789012345678901234567892")
		signing3 = "33337293a4cc2621b61193135a959287aaaaffff"
		node3    = newNode(t, minimalNodeCfg{
			id:            "test node 3",
			pubKey:        pubKey3,
			registryChain: registryChain,
			p2p:           p2p3,
			signingAddr:   signing3,
			admin:         admin3,
		})

		p2p4     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(400))
		pubKey4  = "44444981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109222"
		admin4   = common.HexToAddress("0x4444567890123456789012345678901234567893")
		signing4 = "44447293a4cc2621b61193135a959287aaaaffff"
		node4    = newNode(t, minimalNodeCfg{
			id:            "test node 4",
			pubKey:        pubKey4,
			registryChain: registryChain,
			p2p:           p2p4,
			signingAddr:   signing4,
			admin:         admin4,
		})

		// Additional node (n5) to be added
		p2p5     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(500))
		pubKey5  = "55554981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109333"
		admin5   = common.HexToAddress("0x5555567890123456789012345678901234567894")
		signing5 = "55557293a4cc2621b61193135a959287aaaabbbb"
		node5    = newNode(t, minimalNodeCfg{
			id:            "test node 5",
			pubKey:        pubKey5,
			registryChain: registryChain,
			p2p:           p2p5,
			signingAddr:   signing5,
			admin:         admin5,
		})

		// Test capability
		testCap = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "test",
			Version:        "1.0.0",
			CapabilityType: 0,
		}
	)

	lggr := logger.Test(t)

	// Setup initial registry with DON containing nodes 1-4
	cfg := setupUpdateDonTestConfig{
		dons: []internal.DonInfo{
			{
				Name:         "don 1",
				Nodes:        []deployment.Node{node1, node2, node3, node4},
				Capabilities: []internal.DONCapabilityWithConfig{{Capability: testCap, Config: kstest.GetDefaultCapConfig(t, testCap)}},
			},
		},
		nops: []internal.NOP{
			{
				Name:  "nop 1",
				Nodes: []string{node1.NodeID, node2.NodeID, node3.NodeID, node4.NodeID},
			},
		},
	}

	testCfg := registerTestDon(t, lggr, cfg)

	// Verify initial DON setup
	initialDon, err := testCfg.CapabilitiesRegistry.GetDON(&bind.CallOpts{}, 1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), initialDon.Id)
	require.Len(t, initialDon.NodeP2PIds, 4)

	testCapCfg := kstest.GetDefaultCapConfig(t, testCap)
	testCapCfgB, err := proto.Marshal(testCapCfg)
	require.NoError(t, err)

	t.Run("add node to DON composition", func(t *testing.T) {
		caps, err := testCfg.CapabilitiesRegistry.GetCapabilities(nil)
		require.NoError(t, err)
		capIDs := make([][32]byte, 0, len(caps))
		for _, c := range caps {
			capIDs = append(capIDs, c.HashedId)
		}

		r, err := internal.AddNodes(lggr, &internal.AddNodesRequest{
			CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
			RegistryChain:        testCfg.Chain,
			NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
				node5.NodeID: {
					NodeOperatorId:      1,
					P2pId:               node5.PeerID,
					Signer:              [32]byte{5: 5},
					EncryptionPublicKey: [32]byte{5: 5},
					HashedCapabilityIds: capIDs,
				},
			},
		})

		require.NoError(t, err)
		t.Logf("Added nodes: %v", r.AddedNodes)
		// Update DON to include all 5 nodes
		req := &internal.UpdateDonRequest{
			CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
			Chain:                testCfg.Chain,
			DonID:                1,
			P2PIDs:               []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID(), p2p5.PeerID()},
			CapabilityConfigs: []internal.CapabilityConfig{
				{Capability: testCap, Config: testCapCfgB},
			},
		}

		resp, err := internal.UpdateDon(lggr, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify the DON now has 5 nodes
		updatedDon, err := testCfg.CapabilitiesRegistry.GetDON(&bind.CallOpts{}, 1)
		require.NoError(t, err)
		require.Equal(t, uint32(1), updatedDon.Id)
		require.Len(t, updatedDon.NodeP2PIds, 5, "DON should now have 5 nodes")

		// Verify the correct P2P IDs are present
		expectedP2PIDs := []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID(), p2p5.PeerID()}
		actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)

		require.ElementsMatch(t, expectedP2PIDs, actualP2PIDs, "DON should contain all 5 expected P2P IDs")

		// nested because we need to remove the node we added in the previous test
		t.Run("remove node from DON composition", func(t *testing.T) {
			// Update DON to only include nodes 1-4 (removing node 5)
			req := &internal.UpdateDonRequest{
				CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
				Chain:                testCfg.Chain,
				DonID:                1,
				P2PIDs:               []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID()},
				CapabilityConfigs: []internal.CapabilityConfig{
					{Capability: testCap, Config: testCapCfgB},
				},
			}

			resp, err := internal.UpdateDon(lggr, req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify the DON is back to 4 nodes
			updatedDon, err := testCfg.CapabilitiesRegistry.GetDON(&bind.CallOpts{}, 1)
			require.NoError(t, err)
			require.Equal(t, uint32(1), updatedDon.Id)
			require.Len(t, updatedDon.NodeP2PIds, 4, "DON should be back to 4 nodes")

			// Verify the correct P2P IDs are present (should match original 4 nodes)
			expectedP2PIDs := []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p4.PeerID()}
			actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)

			require.ElementsMatch(t, expectedP2PIDs, actualP2PIDs, "DON should contain original 4 P2P IDs")

			// Verify node 5 is not in the DON
			for _, actualP2PID := range actualP2PIDs {
				require.NotEqual(t, p2p5.PeerID(), actualP2PID, "Node 5 should not be in the DON")
			}
		})

		// nested we swap the node we just re-added
		t.Run("replace nodes in DON composition", func(t *testing.T) {
			// Register another node (node_6) for replacement test
			p2p6 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(600))
			pubKey6 := "66664981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109444"
			admin6 := common.HexToAddress("0x6666567890123456789012345678901234567895")
			signing6 := "66667293a4cc2621b61193135a959287aaaacccc"
			node6 := newNode(t, minimalNodeCfg{
				id:            "test node 6",
				pubKey:        pubKey6,
				registryChain: registryChain,
				p2p:           p2p6,
				signingAddr:   signing6,
				admin:         admin6,
			})

			// Register node_6 with capabilities
			caps, err := testCfg.CapabilitiesRegistry.GetCapabilities(nil)
			require.NoError(t, err)
			capIDs := make([][32]byte, 0, len(caps))
			for _, c := range caps {
				capIDs = append(capIDs, c.HashedId)
			}

			r, err := internal.AddNodes(lggr, &internal.AddNodesRequest{
				CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
				RegistryChain:        testCfg.Chain,
				NodeParams: map[string]kcr.CapabilitiesRegistryNodeParams{
					node6.NodeID: {
						NodeOperatorId:      1,
						P2pId:               node6.PeerID,
						Signer:              [32]byte{6: 6},
						EncryptionPublicKey: [32]byte{6: 6},
						HashedCapabilityIds: capIDs,
					},
				},
			})
			require.NoError(t, err)
			lggr.Debugf("Added node 6: %v", r.AddedNodes)

			// Update DON to replace node4 with node_6 (keeping nodes 1, 2, 3, and adding 6)
			req := &internal.UpdateDonRequest{
				CapabilitiesRegistry: testCfg.CapabilitiesRegistry,
				Chain:                testCfg.Chain,
				DonID:                1,
				P2PIDs:               []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p6.PeerID()},
				CapabilityConfigs: []internal.CapabilityConfig{
					{Capability: testCap, Config: testCapCfgB},
				},
			}

			resp, err := internal.UpdateDon(lggr, req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify the DON still has 4 nodes but with node_6 instead of node4
			updatedDon, err := testCfg.CapabilitiesRegistry.GetDON(&bind.CallOpts{}, 1)
			require.NoError(t, err)
			require.Equal(t, uint32(1), updatedDon.Id)
			require.Len(t, updatedDon.NodeP2PIds, 4, "DON should still have 4 nodes")

			// Verify the correct P2P IDs are present (nodes 1, 2, 3, 6)
			expectedP2PIDs := []p2pkey.PeerID{p2p1.PeerID(), p2p2.PeerID(), p2p3.PeerID(), p2p6.PeerID()}
			actualP2PIDs := internal.BytesToPeerIDs(updatedDon.NodeP2PIds)

			require.ElementsMatch(t, expectedP2PIDs, actualP2PIDs, "DON should contain nodes 1, 2, 3, and 6")

			// Verify node4 is not in the DON anymore
			for _, actualP2PID := range actualP2PIDs {
				require.NotEqual(t, p2p4.PeerID(), actualP2PID, "Node 4 should not be in the DON")
			}

			// Verify node_6 is in the DON
			foundNode6 := slices.Contains(actualP2PIDs, p2p6.PeerID())
			require.True(t, foundNode6, "Node 6 should be in the DON")
		})
	})
}

func sortedP2Pids(p2pids [][32]byte) [][32]byte {
	// sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	return p2pids
}

func capIDs(ccs []kcr.CapabilitiesRegistryCapabilityConfiguration) [][32]byte {
	out := make([][32]byte, len(ccs))
	for i, cc := range ccs {
		out[i] = cc.CapabilityId
	}
	sort.Slice(out, func(i, j int) bool {
		return bytes.Compare(out[i][:], out[j][:]) < 0
	})
	return out
}

type minimalNodeCfg struct {
	id            string
	pubKey        string
	registryChain chainsel.Chain
	p2p           p2pkey.KeyV2
	signingAddr   string
	admin         common.Address
}

func newNode(t *testing.T, cfg minimalNodeCfg) deployment.Node {
	t.Helper()

	registryChainID, err := chainsel.ChainIdFromSelector(registryChain.Selector)
	if err != nil {
		panic(err)
	}
	registryChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(registryChainID)), chainsel.FamilyEVM) //nolint:gosec // G115
	if err != nil {
		panic(err)
	}

	signingAddr, err := hex.DecodeString(cfg.signingAddr)
	require.NoError(t, err)

	var pubkey [32]byte
	if _, err := hex.Decode(pubkey[:], []byte(cfg.pubKey)); err != nil {
		panic(fmt.Sprintf("failed to decode pubkey %s: %v", pubkey, err))
	}

	return deployment.Node{
		NodeID:    cfg.id,
		PeerID:    cfg.p2p.PeerID(),
		CSAKey:    cfg.pubKey,
		AdminAddr: cfg.admin.String(),
		SelToOCRConfig: map[chainsel.ChainDetails]deployment.OCRConfig{
			registryChainDetails: {
				OnchainPublicKey:          signingAddr,
				PeerID:                    cfg.p2p.PeerID(),
				ConfigEncryptionPublicKey: pubkey,
			},
		},
	}
}

type setupUpdateDonTestConfig struct {
	dons []internal.DonInfo
	nops []internal.NOP
}

func registerTestDon(t *testing.T, lggr logger.Logger, cfg setupUpdateDonTestConfig) *kstest.SetupTestRegistryResponse {
	t.Helper()
	req := newSetupTestRegistryRequest(t, cfg.dons, cfg.nops)
	return kstest.SetupTestRegistry(t, lggr, req)
}

func newSetupTestRegistryRequest(t *testing.T, dons []internal.DonInfo, nops []internal.NOP) *kstest.SetupTestRegistryRequest {
	t.Helper()
	nodes := make(map[string]deployment.Node)
	for _, don := range dons {
		for _, node := range don.Nodes {
			nodes[node.NodeID] = node
		}
	}
	nopsToNodes := makeNopToNodes(t, nops, nodes)
	testDons := makeTestDon(t, dons)
	p2pToCapabilities := makeP2PToCapabilities(t, dons)
	req := &kstest.SetupTestRegistryRequest{
		NopToNodes:        nopsToNodes,
		Dons:              testDons,
		P2pToCapabilities: p2pToCapabilities,
	}
	return req
}

func makeNopToNodes(t *testing.T, nops []internal.NOP, nodes map[string]deployment.Node) map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc {
	nopToNodes := make(map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc)

	for _, nop := range nops {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		crnop := kcr.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: common.HexToAddress(nodes[nop.Nodes[0]].AdminAddr),
		}
		var signers []*internal.P2PSignerEnc
		for _, nodeID := range nop.Nodes {
			node := nodes[nodeID]
			require.NotNil(t, node.CSAKey, "public key is nil %s", node.NodeID)
			// all chain configs are the same wrt admin address & node keys
			p, err := kscs.NewP2PSignerEnc(&node, registryChain.Selector)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.NodeID)
			signers = append(signers, p)
		}
		nopToNodes[crnop] = signers
	}
	return nopToNodes
}

func makeP2PToCapabilities(t *testing.T, dons []internal.DonInfo) map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability {
	p2pToCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		for _, node := range don.Nodes {
			for _, cap := range don.Capabilities {
				p, err := kscs.NewP2PSignerEnc(&node, registryChain.Selector)
				require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.NodeID)
				p2pToCapabilities[p.P2PKey] = append(p2pToCapabilities[p.P2PKey], cap.Capability)
			}
		}
	}
	return p2pToCapabilities
}

func makeTestDon(t *testing.T, dons []internal.DonInfo) []kstest.Don {
	out := make([]kstest.Don, len(dons))
	for i, don := range dons {
		out[i] = testDon(t, don)
	}
	return out
}

func testDon(t *testing.T, don internal.DonInfo) kstest.Don {
	var p2pids []p2pkey.PeerID
	for _, node := range don.Nodes {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		p, err := kscs.NewP2PSignerEnc(&node, registryChain.Selector)
		require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.NodeID)
		p2pids = append(p2pids, p.P2PKey)
	}

	var capabilityConfigs []internal.CapabilityConfig
	for i := range don.Capabilities {
		donCap := &don.Capabilities[i]
		cfg, err := proto.Marshal(donCap.Config)
		require.NoError(t, err)
		capabilityConfigs = append(capabilityConfigs, internal.CapabilityConfig{
			Capability: donCap.Capability, Config: cfg,
		})
	}
	return kstest.Don{
		Name:              don.Name,
		P2PIDs:            p2pids,
		CapabilityConfigs: capabilityConfigs,
	}
}
