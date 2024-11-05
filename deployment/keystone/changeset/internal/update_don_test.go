package internal_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/clo/models"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kscs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/test-go/testify/require"
	"gotest.tools/v3/assert"
)

func TestUpdateDon(t *testing.T) {
	var (
		registryChain = chainsel.TEST_90000001
		//nop_1         = "test nop 1"
		//nop_2         = "test nop 2"
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
		/*
			cap_B = kcr.CapabilitiesRegistryCapability{
				LabelledName:   "cap b",
				Version:        "1.0.0",
				CapabilityType: 1,
			}

		*/
	)
	//		p2p_3 = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(300))

	lggr := logger.Test(t)

	t.Run("empty", func(t *testing.T) {
		cfg := setupUpdateDonTestConfig{
			dons: []kslib.DonCapabilities{
				{
					Name: "don 1",
					Nops: []*models.NodeOperator{
						{
							Name:  "nop 1",
							Nodes: []*models.Node{node_1, node_2, node_3, node_4},
						},
					},
					Capabilities: []kcr.CapabilitiesRegistryCapability{cap_A},
				},
			},
		}

		testCfg := setupUpdateDonTest(t, lggr, cfg)
		req := &internal.UpdateDonRequest{
			Registry: testCfg.Registry,
			Chain:    testCfg.Chain,
			P2PIDs:   []p2pkey.PeerID{p2p_1.PeerID(), p2p_2.PeerID(), p2p_3.PeerID(), p2p_4.PeerID()},
			CapabilityConfigs: []internal.CapabilityConfig{
				{},
			},
		}
		want := &internal.UpdateDonResponse{}

		got, err := internal.UpdateDon(lggr, req)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

type minimalNodeCfg struct {
	id            string
	pubKey        string
	registryChain chainsel.Chain
	p2p           p2pkey.KeyV2
	signingAddr   string
	admin         common.Address
}

func newNode(t *testing.T, cfg minimalNodeCfg) *models.Node {
	t.Helper()

	return &models.Node{
		ID:        cfg.id,
		PublicKey: &cfg.pubKey,
		ChainConfigs: []*models.NodeChainConfig{
			{
				ID: "test chain",
				Network: &models.Network{
					ID:        "test network 1",
					ChainID:   strconv.FormatUint(cfg.registryChain.EvmChainID, 10),
					ChainType: models.ChainTypeEvm,
				},
				AdminAddress: cfg.admin.String(),
				Ocr2Config: &models.NodeOCR2Config{
					P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
						PeerID: cfg.p2p.PeerID().String(),
					},
					OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
						OnchainSigningAddress: cfg.signingAddr,
					},
				},
			},
		},
	}
}

type setupUpdateDonTestConfig struct {
	dons []kslib.DonCapabilities
}

type setupUpdateDonTestResult struct {
	registry *kcr.CapabilitiesRegistry
	chain    deployment.Chain
}

func setupUpdateDonTest(t *testing.T, lggr logger.Logger, cfg setupUpdateDonTestConfig) *kstest.SetupTestRegistryResponse {
	t.Helper()
	req := newSetupTestRegistryRequest(t, cfg.dons)
	return kstest.SetupTestRegistry(t, lggr, req)
}

func newSetupTestRegistryRequest(t *testing.T, dons []kslib.DonCapabilities) *kstest.SetupTestRegistryRequest {
	t.Helper()
	allNops := make(map[string]*models.NodeOperator)
	for _, don := range dons {
		for _, nop := range don.Nops {
			nop := nop
			n, exists := allNops[nop.ID]
			if exists {
				nop.Nodes = append(n.Nodes, nop.Nodes...)
			}
			allNops[nop.ID] = nop
		}
	}
	var nops []*models.NodeOperator
	for _, nop := range allNops {
		nops = append(nops, nop)
	}
	nopsToNodes := makeNopToNodes(t, nops)
	testDons := makeTestDon(t, dons)
	p2pToCapabilities := makeP2PToCapabilities(t, dons)
	req := &kstest.SetupTestRegistryRequest{
		NopToNodes:        nopsToNodes,
		Dons:              testDons,
		P2pToCapabilities: p2pToCapabilities,
	}
	return req
}

func makeNopToNodes(t *testing.T, cloNops []*models.NodeOperator) map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc {
	nopToNodes := make(map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc)

	for _, nop := range cloNops {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		crnop := kcr.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: common.HexToAddress(nop.Nodes[0].ChainConfigs[0].AdminAddress),
		}
		var nodes []*internal.P2PSignerEnc
		for _, node := range nop.Nodes {
			require.NotNil(t, node.PublicKey, "public key is nil %s", node.ID)
			// all chain configs are the same wrt admin address & node keys
			p, err := kscs.NewP2PSignerEncFromCLO(node.ChainConfigs[0], *node.PublicKey)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
			nodes = append(nodes, p)
		}
		nopToNodes[crnop] = nodes
	}
	return nopToNodes
}

func makeP2PToCapabilities(t *testing.T, dons []kslib.DonCapabilities) map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability {
	p2pToCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		for _, nop := range don.Nops {
			for _, node := range nop.Nodes {
				for _, cap := range don.Capabilities {
					p, err := kscs.NewP2PSignerEncFromCLO(node.ChainConfigs[0], *node.PublicKey)
					require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
					p2pToCapabilities[p.P2PKey] = append(p2pToCapabilities[p.P2PKey], cap)
				}
			}
		}
	}
	return p2pToCapabilities
}

func makeTestDon(t *testing.T, dons []kslib.DonCapabilities) []kstest.Don {
	out := make([]kstest.Don, len(dons))
	for i, don := range dons {
		out[i] = testDon(t, don)
	}
	return out
}

func testDon(t *testing.T, don kslib.DonCapabilities) kstest.Don {
	var p2pids []p2pkey.PeerID
	for _, nop := range don.Nops {
		for _, node := range nop.Nodes {
			// all chain configs are the same wrt admin address & node keys
			// so we can just use the first one
			p, err := kscs.NewP2PSignerEncFromCLO(node.ChainConfigs[0], *node.PublicKey)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
			p2pids = append(p2pids, p.P2PKey)
		}
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
