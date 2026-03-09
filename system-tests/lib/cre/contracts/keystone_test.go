package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink/deployment"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// fakeOffchainClient implements offchain.Client; only ListNodes and
// ListNodeChainConfigs are wired, every other method panics via the embedded
// nil interface.
type fakeOffchainClient struct {
	cldf_offchain.Client
	nodesByID map[string]*fakeNodeInfo
}

type fakeNodeInfo struct {
	id           string
	name         string
	csaKey       string
	workflowKey  string
	p2pID        string
	chainConfigs []*nodev1.ChainConfig
}

func newFakeOffchainClient(nodes []*fakeNodeInfo) *fakeOffchainClient {
	f := &fakeOffchainClient{nodesByID: make(map[string]*fakeNodeInfo)}
	for _, n := range nodes {
		f.nodesByID[n.id] = n
	}
	return f
}

func (f *fakeOffchainClient) ListNodes(_ context.Context, in *nodev1.ListNodesRequest, _ ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	var wantP2P map[string]bool
	if in.Filter != nil {
		for _, sel := range in.Filter.Selectors {
			if sel.Key == "p2p_id" && sel.Op == ptypes.SelectorOp_IN && sel.Value != nil {
				wantP2P = make(map[string]bool)
				for _, v := range strings.Split(*sel.Value, ",") {
					wantP2P[v] = true
				}
			}
		}
	}

	out := make([]*nodev1.Node, 0, len(f.nodesByID))
	for _, n := range f.nodesByID {
		if wantP2P != nil && !wantP2P[n.p2pID] {
			continue
		}
		p2pVal := n.p2pID
		wfKey := n.workflowKey
		out = append(out, &nodev1.Node{
			Id:          n.id,
			Name:        n.name,
			PublicKey:   n.csaKey,
			WorkflowKey: &wfKey,
			IsEnabled:   true,
			Labels:      []*ptypes.Label{{Key: "p2p_id", Value: &p2pVal}},
		})
	}
	return &nodev1.ListNodesResponse{Nodes: out}, nil
}

func (f *fakeOffchainClient) ListNodeChainConfigs(_ context.Context, in *nodev1.ListNodeChainConfigsRequest, _ ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {
	if in.Filter == nil || len(in.Filter.NodeIds) == 0 {
		return nil, errors.New("filter with node IDs required")
	}
	n, ok := f.nodesByID[in.Filter.NodeIds[0]]
	if !ok {
		return nil, fmt.Errorf("node not found: %s", in.Filter.NodeIds[0])
	}
	return &nodev1.ListNodeChainConfigsResponse{ChainConfigs: n.chainConfigs}, nil
}

func TestDonsOrderedByID(t *testing.T) {
	// Test donsOrderedByID sorts by id ascending
	d := dons{
		c: make(map[string]donConfig),
	}

	d.c["don3"] = donConfig{id: 3}
	d.c["don1"] = donConfig{id: 1}
	d.c["don2"] = donConfig{id: 2}

	ordered := d.donsOrderedByID()
	if len(ordered) != 3 {
		t.Fatalf("expected 3 dons, got %d", len(ordered))
	}

	if ordered[0].id != 1 || ordered[1].id != 2 || ordered[2].id != 3 {
		t.Fatalf("expected dons ordered by id 1,2,3 got %d,%d,%d", ordered[0].id, ordered[1].id, ordered[2].id)
	}
}

func TestToV2ConfigureInput(t *testing.T) {
	chainSel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	chainID, err := chainselectors.GetChainIDFromSelector(chainSel)
	require.NoError(t, err)

	key1 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	key2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	peerID1 := key1.PeerID().String()
	peerID2 := key2.PeerID().String()

	fakeNodes := []*fakeNodeInfo{
		{
			id:          "node_01",
			name:        "test-node-1",
			csaKey:      "403b72f0b1b3b5f5a91bcfedb7f28599767502a04b5b7e067fcf3782e23eeb9c",
			workflowKey: "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
			p2pID:       peerID1,
			chainConfigs: []*nodev1.ChainConfig{{
				Chain: &nodev1.Chain{
					Type: nodev1.ChainType_CHAIN_TYPE_EVM,
					Id:   chainID,
				},
				Ocr2Config: &nodev1.OCR2Config{
					OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
						OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
						OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
						ConfigPublicKey:       "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
						BundleId:              "665a101d79d310cb0a5ebf695b06e8fc8082b5cbe62d7d362d80d47447a31fea",
					},
					P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
						PeerId: peerID1,
					},
					IsBootstrap: false,
				},
				AccountAddress: "0x2877F08d9c5Cc9F401F730Fa418fAE563A9a2FF3",
			}},
		},
		{
			id:          "node_02",
			name:        "test-node-2",
			csaKey:      "28b91143ec9111796a7d63e14c1cf6bb01b4ed59667ab54f5bc72ebe49c881be",
			workflowKey: "2c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
			p2pID:       peerID2,
			chainConfigs: []*nodev1.ChainConfig{{
				Chain: &nodev1.Chain{
					Type: nodev1.ChainType_CHAIN_TYPE_EVM,
					Id:   chainID,
				},
				Ocr2Config: &nodev1.OCR2Config{
					OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
						OffchainPublicKey:     "255096a3b7ade10e29c648e0b407fc486180464f713446b1da04f013df6179c8",
						OnchainSigningAddress: "8258f4c4761cc445333017608044a204fd0c006a",
						ConfigPublicKey:       "2c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
						BundleId:              "7a9b75510b8d09932b98142419bef52436ff725dd9395469473b487ef87fdfb0",
					},
					P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
						PeerId: peerID2,
					},
					IsBootstrap: false,
				},
				AccountAddress: "0x415aa1E9a1bcB3929ed92bFa1F9735Dc0D45AD31",
			}},
		},
	}

	offchainClient := newFakeOffchainClient(fakeNodes)

	d := &dons{
		c:        make(map[string]donConfig),
		offChain: offchainClient,
	}

	d.c["test-don"] = donConfig{
		id: 1,
		DonCapabilities: keystone_changeset.DonCapabilities{
			Name: "test-don",
			F:    1,
			Nops: []keystone_changeset.NOP{
				{
					Name:  "test-nop",
					Nodes: []string{peerID1, peerID2},
				},
			},
			Capabilities: []keystone_changeset.DONCapabilityWithConfig{
				{
					Capability: kcr.CapabilitiesRegistryCapability{
						LabelledName:   "test-capability",
						Version:        "1.0.0",
						CapabilityType: 1,
					},
					Config: &capabilitiespb.CapabilityConfig{},
				},
			},
		},
	}

	result := d.mustToV2ConfigureInput(chainSel, "0x1234567890abcdef", nil)

	require.Equal(t, chainSel, result.RegistryChainSel)

	require.Len(t, result.Nops, 1)
	require.Equal(t, "test-nop", result.Nops[0].Name)

	require.Len(t, result.Nodes, 2)

	require.Len(t, result.Capabilities, 1)
	require.Equal(t, "test-capability@1.0.0", result.Capabilities[0].CapabilityID)

	require.Len(t, result.DONs, 1)
	require.Equal(t, "test-don", result.DONs[0].Name)
	require.Equal(t, uint8(1), result.DONs[0].F)
	require.Len(t, result.DONs[0].Nodes, 2)
	require.Len(t, result.DONs[0].CapabilityConfigurations, 1)
}

func TestToV2ConfigureInput_EmbedsAptosTransmittersInSpecConfig(t *testing.T) {
	registryChainSel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	registryChainID, err := chainselectors.GetChainIDFromSelector(registryChainSel)
	require.NoError(t, err)

	aptosDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily("2", chainselectors.FamilyAptos)
	require.NoError(t, err)
	aptosChainID, err := chainselectors.GetChainIDFromSelector(aptosDetails.ChainSelector)
	require.NoError(t, err)

	key1 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(11))
	key2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(12))
	peerID1 := key1.PeerID().String()
	peerID2 := key2.PeerID().String()

	aptosTransmitterA := "0xbbb"
	aptosTransmitterB := "0xaaa"

	fakeNodes := []*fakeNodeInfo{
		{
			id:          "node_aptos_1",
			name:        "aptos-node-1",
			csaKey:      "403b72f0b1b3b5f5a91bcfedb7f28599767502a04b5b7e067fcf3782e23eeb9c",
			workflowKey: "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
			p2pID:       peerID1,
			chainConfigs: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   registryChainID,
					},
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
							OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
							ConfigPublicKey:       "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
							BundleId:              "665a101d79d310cb0a5ebf695b06e8fc8082b5cbe62d7d362d80d47447a31fea",
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: peerID1,
						},
						IsBootstrap: false,
					},
					AccountAddress: "0x2877F08d9c5Cc9F401F730Fa418fAE563A9a2FF3",
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
						Id:   aptosChainID,
					},
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							OffchainPublicKey:     "13dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
							OnchainSigningAddress: "1f",
							ConfigPublicKey:       "6193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
							BundleId:              "765a101d79d310cb0a5ebf695b06e8fc8082b5cbe62d7d362d80d47447a31feb",
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: peerID1,
						},
						IsBootstrap: false,
					},
					AccountAddressPublicKey: &aptosTransmitterA,
				},
			},
		},
		{
			id:          "node_aptos_2",
			name:        "aptos-node-2",
			csaKey:      "28b91143ec9111796a7d63e14c1cf6bb01b4ed59667ab54f5bc72ebe49c881be",
			workflowKey: "2c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
			p2pID:       peerID2,
			chainConfigs: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   registryChainID,
					},
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							OffchainPublicKey:     "255096a3b7ade10e29c648e0b407fc486180464f713446b1da04f013df6179c8",
							OnchainSigningAddress: "8258f4c4761cc445333017608044a204fd0c006a",
							ConfigPublicKey:       "2c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
							BundleId:              "7a9b75510b8d09932b98142419bef52436ff725dd9395469473b487ef87fdfb0",
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: peerID2,
						},
						IsBootstrap: false,
					},
					AccountAddress: "0x415aa1E9a1bcB3929ed92bFa1F9735Dc0D45AD31",
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
						Id:   aptosChainID,
					},
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							OffchainPublicKey:     "355096a3b7ade10e29c648e0b407fc486180464f713446b1da04f013df6179c8",
							OnchainSigningAddress: "2f",
							ConfigPublicKey:       "3c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
							BundleId:              "8a9b75510b8d09932b98142419bef52436ff725dd9395469473b487ef87fdfb1",
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: peerID2,
						},
						IsBootstrap: false,
					},
					AccountAddressPublicKey: &aptosTransmitterB,
				},
			},
		},
	}

	offchainClient := newFakeOffchainClient(fakeNodes)

	d := &dons{
		c:        make(map[string]donConfig),
		offChain: offchainClient,
	}
	d.c["aptos-don"] = donConfig{
		id: 7,
		DonCapabilities: keystone_changeset.DonCapabilities{
			Name: "aptos-don",
			F:    1,
			Nops: []keystone_changeset.NOP{{
				Name:  "aptos-nop",
				Nodes: []string{peerID1, peerID2},
			}},
			Capabilities: []keystone_changeset.DONCapabilityWithConfig{{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("aptos:ChainSelector:%d", aptosDetails.ChainSelector),
					Version:        "1.0.0",
					CapabilityType: 1,
				},
				Config: &capabilitiespb.CapabilityConfig{},
			}},
		},
	}

	result := d.mustToV2ConfigureInput(registryChainSel, "0x1234567890abcdef", nil)
	require.Len(t, result.DONs, 1)
	require.Len(t, result.DONs[0].CapabilityConfigurations, 1)

	var cfg capabilitiespb.CapabilityConfig
	err = proto.Unmarshal(result.DONs[0].CapabilityConfigurations[0].Config, &cfg)
	require.NoError(t, err)
	require.NotNil(t, cfg.SpecConfig)
	specCfg, err := values.FromMapValueProto(cfg.SpecConfig)
	require.NoError(t, err)
	require.NotNil(t, specCfg)

	rawTransmitters, ok := specCfg.Underlying[aptosSpecConfigTransmittersListKey]
	require.True(t, ok)

	var transmitters []string
	err = rawTransmitters.UnwrapTo(&transmitters)
	require.NoError(t, err)
	require.Equal(t, []string{aptosTransmitterA, aptosTransmitterB}, transmitters)
}

func TestAptosTransmittersForChainSelector_ErrorsOnEmptyTransmitAccount(t *testing.T) {
	chainSelector := uint64(4457093679053095497)
	nodes := []deployment.Node{
		{
			Name: "node-1",
			SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{
				{ChainSelector: chainSelector}: {TransmitAccount: types2.Account("0xbbb")},
			},
		},
		{
			Name: "node-2",
			SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{
				{ChainSelector: chainSelector}: {TransmitAccount: types2.Account("")},
			},
		},
		{
			Name: "node-3",
			SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{
				{ChainSelector: chainSelector}: {TransmitAccount: types2.Account("0xaaa")},
			},
		},
	}

	_, err := aptosTransmittersForChainSelector(nodes, chainSelector)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty Aptos transmitter for node node-2")
}

func TestValidateAptosTransmittersForCapabilities_ErrorsOnInvalidAptosTransmitterSet(t *testing.T) {
	chainSelector := uint64(4457093679053095497)
	nodes := []deployment.Node{
		{
			Name: "node-1",
			SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{
				{ChainSelector: chainSelector}: {TransmitAccount: types2.Account("0xabc")},
			},
		},
		{
			Name: "node-2",
			SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{
				{ChainSelector: chainSelector}: {TransmitAccount: types2.Account("")},
			},
		},
	}

	caps := []keystone_changeset.DONCapabilityWithConfig{
		{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName: fmt.Sprintf("aptos:ChainSelector:%d", chainSelector),
				Version:      "1.0.0",
			},
		},
	}

	err := validateAptosTransmittersForCapabilities(nodes, caps)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty Aptos transmitter for node node-2")
}

// TestGenerateAdminAddresses contains all the test cases for the function.
func TestGenerateAdminAddresses(t *testing.T) {
	// Test Case 1: Basic Functionality
	t.Run("Basic_Functionality_10_Addresses", func(t *testing.T) {
		count := 10
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, addresses, count, "Expected slice of length %d, but got %d", count, len(addresses))

		// Check for uniqueness and validity
		addressMap := make(map[common.Address]bool)
		for _, addr := range addresses {
			require.True(t, common.IsHexAddress(addr.Hex()))
			require.NotEqual(t, 0, addr.Cmp(common.HexToAddress("0x0000000000000000000000000000000000000000")), "Generated a zero address, which should be avoided")
			addressMap[addr] = true
		}
		require.Len(t, addressMap, count, "Expected unique address count of %d, but got %d", count, len(addressMap))
	})

	// Test Case 2: Smallest Valid Input
	t.Run("Smallest_Valid_Input_1_Address", func(t *testing.T) {
		count := 1
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, addresses, count, "Expected slice of length %d, but got %d", count, len(addresses))
	})

	// Test Case 3: Invalid Input (Zero and Negative Count)
	t.Run("Invalid_Input_Zero_Count", func(t *testing.T) {
		count := 0
		_, err := generateAdminAddresses(count)
		require.Error(t, err, "Expected an error for count %d, but got none", count)
	})

	t.Run("Invalid_Input_Negative_Count", func(t *testing.T) {
		count := -5
		_, err := generateAdminAddresses(count)
		require.Error(t, err, "Expected an error for count %d, but got none", count)
	})

	// Test that 5 digit padding starts at boundary
	t.Run("Boundary_Condition_65536_Addresses", func(t *testing.T) {
		count := 65536
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, addresses, count, "Expected slice of length %d, but got %d", count, len(addresses))

		for _, addr := range addresses {
			require.True(t, common.IsHexAddress(addr.String()), "invalid address: %s", addr)
		}
	})
}
