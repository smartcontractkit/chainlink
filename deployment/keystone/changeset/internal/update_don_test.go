package internal_test

import (
	"bytes"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kscs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateDon(t *testing.T) {
	var (
		registryChain = chainsel.TEST_90000001
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
		cap_A = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "test",
			Version:        "1.0.0",
			CapabilityType: 0,
		}

		cap_B = kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cap b",
			Version:        "1.0.0",
			CapabilityType: 1,
		}
	)

	lggr := logger.Test(t)

	t.Run("empty", func(t *testing.T) {
		cfg := setupUpdateDonTestConfig{
			dons: []kslib.DonInfo{
				{
					Name:         "don 1",
					Nodes:        []keystone.Node{node_1, node_2, node_3, node_4},
					Capabilities: []kcr.CapabilitiesRegistryCapability{cap_A},
				},
			},
			nops: []keystone.NOP{
				{
					Name:  "nop 1",
					Nodes: []string{node_1.ID, node_2.ID, node_3.ID, node_4.ID},
				},
			},
		}

		testCfg := setupUpdateDonTest(t, lggr, cfg)

		req := &internal.UpdateDonRequest{
			Registry: testCfg.Registry,
			Chain:    testCfg.Chain,
			P2PIDs:   []p2pkey.PeerID{p2p_1.PeerID(), p2p_2.PeerID(), p2p_3.PeerID(), p2p_4.PeerID()},
			CapabilityConfigs: []internal.CapabilityConfig{
				{Capability: cap_A}, {Capability: cap_B},
			},
		}
		want := &internal.UpdateDonResponse{
			DonInfo: kcr.CapabilitiesRegistryDONInfo{
				Id:          1,
				ConfigCount: 1,
				NodeP2PIds:  internal.PeerIDsToBytes([]p2pkey.PeerID{p2p_1.PeerID(), p2p_2.PeerID(), p2p_3.PeerID(), p2p_4.PeerID()}),
				CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
					{CapabilityId: kstest.MustCapabilityId(t, testCfg.Registry, cap_A)},
					{CapabilityId: kstest.MustCapabilityId(t, testCfg.Registry, cap_B)},
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

func newNode(t *testing.T, cfg minimalNodeCfg) keystone.Node {
	t.Helper()

	return keystone.Node{
		ID:        cfg.id,
		PublicKey: &cfg.pubKey,
		ChainConfigs: []*nodev1.ChainConfig{
			{
				Chain: &nodev1.Chain{
					Id:   "test chain",
					Type: nodev1.ChainType_CHAIN_TYPE_EVM,
				},
				AdminAddress: cfg.admin.String(),
				Ocr2Config: &nodev1.OCR2Config{
					P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
						PeerId: cfg.p2p.PeerID().String(),
					},
					OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
						OnchainSigningAddress: cfg.signingAddr,
					},
				},
			},
		},
	}
}

type setupUpdateDonTestConfig struct {
	dons []kslib.DonInfo
	nops []keystone.NOP
}

type setupUpdateDonTestResult struct {
	registry *kcr.CapabilitiesRegistry
	chain    deployment.Chain
}

func setupUpdateDonTest(t *testing.T, lggr logger.Logger, cfg setupUpdateDonTestConfig) *kstest.SetupTestRegistryResponse {
	t.Helper()
	req := newSetupTestRegistryRequest(t, cfg.dons, cfg.nops)
	return kstest.SetupTestRegistry(t, lggr, req)
}

func newSetupTestRegistryRequest(t *testing.T, dons []kslib.DonInfo, nops []keystone.NOP) *kstest.SetupTestRegistryRequest {
	t.Helper()
	nodes := make(map[string]keystone.Node)
	for _, don := range dons {
		for _, node := range don.Nodes {
			nodes[node.ID] = node
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

func makeNopToNodes(t *testing.T, nops []keystone.NOP, nodes map[string]keystone.Node) map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc {
	nopToNodes := make(map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc)

	for _, nop := range nops {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		crnop := kcr.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: common.HexToAddress(nodes[nop.Nodes[0]].ChainConfigs[0].AdminAddress),
		}
		var signers []*internal.P2PSignerEnc
		for _, nodeID := range nop.Nodes {
			node := nodes[nodeID]
			require.NotNil(t, node.PublicKey, "public key is nil %s", node.ID)
			// all chain configs are the same wrt admin address & node keys
			p, err := kscs.NewP2PSignerEncFromJD(node.ChainConfigs[0], *node.PublicKey)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
			signers = append(signers, p)
		}
		nopToNodes[crnop] = signers
	}
	return nopToNodes
}

func makeP2PToCapabilities(t *testing.T, dons []kslib.DonInfo) map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability {
	p2pToCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		for _, node := range don.Nodes {
			for _, cap := range don.Capabilities {
				p, err := kscs.NewP2PSignerEncFromJD(node.ChainConfigs[0], *node.PublicKey)
				require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
				p2pToCapabilities[p.P2PKey] = append(p2pToCapabilities[p.P2PKey], cap)
			}
		}
	}
	return p2pToCapabilities
}

func makeTestDon(t *testing.T, dons []kslib.DonInfo) []kstest.Don {
	out := make([]kstest.Don, len(dons))
	for i, don := range dons {
		out[i] = testDon(t, don)
	}
	return out
}

func testDon(t *testing.T, don kslib.DonInfo) kstest.Don {
	var p2pids []p2pkey.PeerID
	for _, node := range don.Nodes {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		p, err := kscs.NewP2PSignerEncFromJD(node.ChainConfigs[0], *node.PublicKey)
		require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
		p2pids = append(p2pids, p.P2PKey)
	}

	var capabilityConfigs []internal.CapabilityConfig
	for _, cap := range don.Capabilities {
		capabilityConfigs = append(capabilityConfigs, internal.CapabilityConfig{
			Capability: cap,
		})
	}
	return kstest.Don{
		Name:              don.Name,
		P2PIDs:            p2pids,
		CapabilityConfigs: capabilityConfigs,
	}
}

func newP2PSignerEnc(signer [32]byte, p2pkey p2pkey.PeerID, encryptionPublicKey [32]byte) *internal.P2PSignerEnc {
	return &internal.P2PSignerEnc{
		Signer:              signer,
		P2PKey:              p2pkey,
		EncryptionPublicKey: encryptionPublicKey,
	}
}
