package internal_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

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
		p2p_1     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(100))
		pubKey_1  = "11114981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e1098e7" // valid csa key
		admin_1   = common.HexToAddress("0x1111567890123456789012345678901234567890")  // valid eth address
		signing_1 = "11117293a4cc2621b61193135a95928735e4795f"                         // valid eth address
		node_1    = newNode(t, minimalNodeCfg{
			id:            "test node 1",
			pubKey:        pubKey_1,
			registryChain: registryChain,
			p2p:           p2p_1,
			signingAddr:   signing_1,
			admin:         admin_1,
		})

		p2p_2     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(200))
		pubKey_2  = "22224981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109000" // valid csa key
		admin_2   = common.HexToAddress("0x2222567890123456789012345678901234567891")  // valid eth address
		signing_2 = "22227293a4cc2621b61193135a95928735e4ffff"                         // valid eth address
		node_2    = newNode(t, minimalNodeCfg{
			id:            "test node 2",
			pubKey:        pubKey_2,
			registryChain: registryChain,
			p2p:           p2p_2,
			signingAddr:   signing_2,
			admin:         admin_2,
		})

		p2p_3     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(300))
		pubKey_3  = "33334981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109111" // valid csa key
		admin_3   = common.HexToAddress("0x3333567890123456789012345678901234567892")  // valid eth address
		signing_3 = "33337293a4cc2621b61193135a959287aaaaffff"                         // valid eth address
		node_3    = newNode(t, minimalNodeCfg{
			id:            "test node 3",
			pubKey:        pubKey_3,
			registryChain: registryChain,
			p2p:           p2p_3,
			signingAddr:   signing_3,
			admin:         admin_3,
		})

		p2p_4     = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(400))
		pubKey_4  = "44444981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e109222" // valid csa key
		admin_4   = common.HexToAddress("0x4444567890123456789012345678901234567893")  // valid eth address
		signing_4 = "44447293a4cc2621b61193135a959287aaaaffff"                         // valid eth address
		node_4    = newNode(t, minimalNodeCfg{
			id:            "test node 4",
			pubKey:        pubKey_4,
			registryChain: registryChain,
			p2p:           p2p_4,
			signingAddr:   signing_4,
			admin:         admin_4,
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
					Nodes:        []deployment.Node{node_1, node_2, node_3, node_4},
					Capabilities: []internal.DONCapabilityWithConfig{{Capability: initialCap, Config: initialCapCfg}},
				},
			},
			nops: []internal.NOP{
				{
					Name:  "nop 1",
					Nodes: []string{node_1.NodeID, node_2.NodeID, node_3.NodeID, node_4.NodeID},
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
			P2PIDs:               []p2pkey.PeerID{p2p_1.PeerID(), p2p_2.PeerID(), p2p_3.PeerID(), p2p_4.PeerID()},
			CapabilityConfigs: []internal.CapabilityConfig{
				{Capability: initialCap, Config: initialCapCfgB}, {Capability: capToAdd, Config: capToAddCfgB},
			},
		}
		want := &internal.UpdateDonResponse{
			DonInfo: kcr.CapabilitiesRegistryDONInfo{
				Id:          1,
				ConfigCount: 1,
				NodeP2PIds:  internal.PeerIDsToBytes([]p2pkey.PeerID{p2p_1.PeerID(), p2p_2.PeerID(), p2p_3.PeerID(), p2p_4.PeerID()}),
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
		assert.Equal(t, capIds(want.DonInfo.CapabilityConfigurations), capIds(got.DonInfo.CapabilityConfigurations))
	})
}

func sortedP2Pids(p2pids [][32]byte) [][32]byte {
	// sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	return p2pids
}

func capIds(ccs []kcr.CapabilitiesRegistryCapabilityConfiguration) [][32]byte {
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
	registryChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(registryChainID)), chainsel.FamilyEVM)
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

type setupUpdateDonTestResult struct {
	registry *kcr.CapabilitiesRegistry
	chain    cldf_evm.Chain
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
