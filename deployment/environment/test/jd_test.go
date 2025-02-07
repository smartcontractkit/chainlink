package test_test

import (
	"context"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestJDNodeService_GetNode(t *testing.T) {
	nodes := []deployment.Node{
		{
			NodeID: "node1",
			Name:   "Node 1",
			CSAKey: "csa_key_1",
			PeerID: testPeerID(t, "peer_id_1"),
			Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}},
		},
		{
			NodeID: "node2",
			Name:   "Node 2",
			PeerID: testPeerID(t, "peer_id_2"),
		},
	}

	service := test.NewJDService(nodes)

	tests := []struct {
		name      string
		nodeID    string
		expectErr bool
		want      *nodev1.GetNodeResponse
	}{
		{
			name:   "existing node",
			nodeID: "node1",
			want:   &nodev1.GetNodeResponse{Node: &nodev1.Node{Id: "node1", Name: "Node 1", PublicKey: "csa_key_1", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}}}},
		},
		{
			name:      "non-existing node",
			nodeID:    "node3",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &nodev1.GetNodeRequest{Id: tt.nodeID}
			resp, err := service.GetNode(context.Background(), req)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.nodeID, resp.Node.Id)
			}
		})
	}
}

func TestJDNodeService_ListNodes(t *testing.T) {
	nodes := []deployment.Node{
		{
			NodeID: "node1",
			Name:   "Node 1",
			CSAKey: "csa_key_1",
			PeerID: testPeerID(t, "peer_id_1"),
			Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}},
		},
		{
			NodeID: "node2",
			Name:   "Node 2",
			CSAKey: "csa_key_2",
			PeerID: testPeerID(t, "peer_id_2"),
			Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}},
		},
	}

	service := test.NewJDService(nodes)

	tests := []struct {
		name   string
		filter *nodev1.ListNodesRequest_Filter
		want   []*nodev1.Node
	}{
		{
			name:   "all nodes",
			filter: nil,
			want: []*nodev1.Node{
				{Id: "node1", Name: "Node 1", PublicKey: "csa_key_1", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}}},
				{Id: "node2", Name: "Node 2", PublicKey: "csa_key_2", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}}},
			},
		},
		{
			name:   "filter by id",
			filter: &nodev1.ListNodesRequest_Filter{Ids: []string{"node1"}},
			want: []*nodev1.Node{
				{Id: "node1", Name: "Node 1", PublicKey: "csa_key_1", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}}},
			},
		},
		{
			name:   "filter EQ common label",
			filter: &nodev1.ListNodesRequest_Filter{Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EQ, Key: "foo", Value: ptr("bar")}}},
			want: []*nodev1.Node{
				{Id: "node1", Name: "Node 1", PublicKey: "csa_key_1", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}}},
				{Id: "node2", Name: "Node 2", PublicKey: "csa_key_2", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}}},
			},
		},
		{
			name:   "filter EQ single label",
			filter: &nodev1.ListNodesRequest_Filter{Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EQ, Key: "baz", Value: ptr("qux")}}},
			want: []*nodev1.Node{
				{Id: "node2", Name: "Node 2", PublicKey: "csa_key_2", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}}},
			},
		},
		{
			name:   "filter EXIST common label name",
			filter: &nodev1.ListNodesRequest_Filter{Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EXIST, Key: "foo"}}},
			want: []*nodev1.Node{
				{Id: "node1", Name: "Node 1", PublicKey: "csa_key_1", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}}},
				{Id: "node2", Name: "Node 2", PublicKey: "csa_key_2", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}}},
			},
		},
		{
			name:   "filter EXIST single label value",
			filter: &nodev1.ListNodesRequest_Filter{Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EXIST, Key: "baz"}}},
			want: []*nodev1.Node{
				{Id: "node2", Name: "Node 2", PublicKey: "csa_key_2", Labels: []*ptypes.Label{{Key: "foo", Value: ptr("bar")}, {Key: "baz", Value: ptr("qux")}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &nodev1.ListNodesRequest{Filter: tt.filter}
			resp, err := service.ListNodes(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.ElementsMatch(t, tt.want, resp.Nodes)
		})
	}
}

func TestJDNodeService_ListNodeChainConfigs(t *testing.T) {
	nodes := []deployment.Node{
		{
			NodeID: "node1",
			Name:   "Node 1",
			CSAKey: "csa_key_1",
			PeerID: testPeerID(t, "peer_id_1"),
			SelToOCRConfig: map[chain_selectors.ChainDetails]deployment.OCRConfig{
				{
					ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
					ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
				}: {
					KeyBundleID:               "bundle1",
					OffchainPublicKey:         test32Byte(t, "offchain-pub-key1"),
					OnchainPublicKey:          types2.OnchainPublicKey([]byte("onchain-pub-key1")),
					ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "config-encryption-pub-key1")),
				},
				{
					ChainSelector: chain_selectors.APTOS_TESTNET.Selector,
					ChainName:     chain_selectors.APTOS_TESTNET.Name,
				}: {
					KeyBundleID:               "aptos bundle1",
					OffchainPublicKey:         test32Byte(t, "aptos offchain-pub-key1"),
					OnchainPublicKey:          types2.OnchainPublicKey([]byte("aptos onchain-pub-key1")),
					ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "aptos config-encryption-pub-key1")),
				},
			},
		},
		{
			NodeID: "node2",
			Name:   "Node 2",
			CSAKey: "csa_key_2",
			PeerID: testPeerID(t, "peer_id_2"),
			SelToOCRConfig: map[chain_selectors.ChainDetails]deployment.OCRConfig{
				{
					ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
					ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
				}: {
					KeyBundleID:               "bundle2",
					OffchainPublicKey:         test32Byte(t, "offchain-pub-key2"),
					OnchainPublicKey:          types2.OnchainPublicKey([]byte("onchain-pub-key2")),
					ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "config-encryption-pub-key2")),
				},
			},
		},
	}

	service := test.NewJDService(nodes)

	tests := []struct {
		name    string
		req     *nodev1.ListNodeChainConfigsRequest
		want    []*nodev1.ChainConfig
		wantErr bool
	}{
		{
			name: "all chain configs",
			req:  &nodev1.ListNodeChainConfigsRequest{},
			want: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
					},
					NodeId: "node1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer_id_1").String(),
						},
					},
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
						Id:   strconv.FormatUint(chain_selectors.APTOS_TESTNET.ChainID, 10),
					},
					NodeId: "node1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "aptos bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "aptos offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("aptos onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "aptos config-encryption-pub-key1"),
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer_id_1").String(),
						},
					},
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
					},
					NodeId: "node2",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "bundle2",
							OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key2"),
							OnchainSigningAddress: hex.EncodeToString([]byte("onchain-pub-key2")),
							ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key2"),
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer_id_2").String(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "filter by node id",
			req:  &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{NodeIds: []string{"node1"}}},
			want: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
					},
					NodeId: "node1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer_id_1").String(),
						},
					},
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
						Id:   strconv.FormatUint(chain_selectors.APTOS_TESTNET.ChainID, 10),
					},
					NodeId: "node1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "aptos bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "aptos offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("aptos onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "aptos config-encryption-pub-key1"),
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer_id_1").String(),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ListNodeChainConfigs(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.ElementsMatch(t, tt.want, resp.ChainConfigs)
			}
		})
	}
}

func TestNewJDServiceFromListNodes(t *testing.T) {
	testData := &nodev1.ListNodesResponse{
		Nodes: []*nodev1.Node{
			{
				Id:        "node1",
				Name:      "Node 1",
				PublicKey: "csa_key_1",
				Labels:    []*ptypes.Label{{Key: "foo", Value: ptr("bar")}},
			},
			{
				Id:        "node2",
				Name:      "Node 2",
				PublicKey: "csa_key_2",
				Labels:    []*ptypes.Label{{Key: "baz", Value: ptr("qux")}},
			},
			{
				Id:        "node3",
				Name:      "Node 3",
				PublicKey: "csa_key_3",
				Labels:    []*ptypes.Label{{Key: "p2p", Value: ptr(testPeerID(t, "peer_id_3").String())}},
			},
		},
	}

	service, err := test.NewJDServiceFromListNodes(testData)
	require.NoError(t, err)

	req := &nodev1.ListNodesRequest{}
	got, err := service.ListNodes(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.ElementsMatch(t, testData.Nodes, got.Nodes)
}

func testPeerID(t *testing.T, s string) p2pkey.PeerID {
	t.Helper()
	var out [32]byte
	copy(out[:], s)
	return p2pkey.PeerID(out)
}

func test32Byte(t *testing.T, s string) [32]byte {
	t.Helper()
	b := []byte(s)
	var out [32]byte
	copy(out[:], b)
	return out
}

func hexFrom32Byte(t *testing.T, s string) string {
	t.Helper()
	b := test32Byte(t, s)
	return hex.EncodeToString(b[:])
}

func ptr[T any](v T) *T {
	return &v
}
