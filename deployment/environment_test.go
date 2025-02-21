package deployment

import (
	"encoding/hex"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestNode_OCRConfigForChainSelector(t *testing.T) {
	var m = map[chain_selectors.ChainDetails]OCRConfig{
		{
			ChainSelector: chain_selectors.APTOS_TESTNET.Selector,
			ChainName:     chain_selectors.APTOS_TESTNET.Name,
		}: {
			KeyBundleID: "aptos bundle 1",
		},
		{
			ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
			ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
		}: {
			KeyBundleID: "arb bundle 1",
		},
	}

	type fields struct {
		SelToOCRConfig map[chain_selectors.ChainDetails]OCRConfig
	}
	type args struct {
		chainSel uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   OCRConfig
		exist  bool
	}{
		{
			name: "aptos ok",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.APTOS_TESTNET.Selector,
			},
			want: OCRConfig{
				KeyBundleID: "aptos bundle 1",
			},
			exist: true,
		},
		{
			name: "arb ok",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
			},
			want: OCRConfig{
				KeyBundleID: "arb bundle 1",
			},
			exist: true,
		},
		{
			name: "no exist",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.WEMIX_MAINNET.Selector, // not in test data
			},
			want:  OCRConfig{},
			exist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{
				SelToOCRConfig: tt.fields.SelToOCRConfig,
			}
			got, got1 := n.OCRConfigForChainSelector(tt.args.chainSel)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.OCRConfigForChainSelector() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exist {
				t.Errorf("Node.OCRConfigForChainSelector() got1 = %v, want %v", got1, tt.exist)
			}
		})
	}
}

func TestNode_ChainConfigs(t *testing.T) {
	type fields struct {
		NodeID         string
		SelToOCRConfig map[chain_selectors.ChainDetails]OCRConfig
		PeerID         p2pkey.PeerID
		IsBootstrap    bool
		MultiAddr      string
		AdminAddr      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*nodev1.ChainConfig
		wantErr bool
	}{
		{
			name: "bad chain",
			fields: fields{
				NodeID: "node-id1",
				SelToOCRConfig: map[chain_selectors.ChainDetails]OCRConfig{
					{
						ChainSelector: math.MaxUint64,
						ChainName:     "wrong",
					}: {
						KeyBundleID:               "bundle1",
						OffchainPublicKey:         test32Byte(t, "offchain-pub-key1"),
						OnchainPublicKey:          types2.OnchainPublicKey([]byte("onchain-pub-key1")),
						ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "config-encryption-pub-key1")),
					},
				},
				PeerID: testPeerID(t, "peer-id1"),
			},
			wantErr: true,
		},

		{
			name: "one evm chain",
			fields: fields{
				NodeID: "node-id1",
				SelToOCRConfig: map[chain_selectors.ChainDetails]OCRConfig{
					{
						ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
						ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
					}: {
						KeyBundleID:               "bundle1",
						OffchainPublicKey:         test32Byte(t, "offchain-pub-key1"),
						OnchainPublicKey:          types2.OnchainPublicKey([]byte("onchain-pub-key1")),
						ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "config-encryption-pub-key1")),
					},
				},
				PeerID: testPeerID(t, "peer-id1"),
			},
			want: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
					},
					NodeId: "node-id1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
						},

						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer-id1").String(),
						},
					},
				},
			},
		},
		{
			name: "one evm chain, one aptos",
			fields: fields{
				NodeID: "node-id1",
				SelToOCRConfig: map[chain_selectors.ChainDetails]OCRConfig{
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
				PeerID: testPeerID(t, "peer-id1"),
			},
			want: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
					},
					NodeId: "node-id1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
						},

						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer-id1").String(),
						},
					},
				},
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
						Id:   strconv.FormatUint(chain_selectors.APTOS_TESTNET.ChainID, 10),
					},
					NodeId: "node-id1",
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							BundleId:              "aptos bundle1",
							OffchainPublicKey:     hexFrom32Byte(t, "aptos offchain-pub-key1"),
							OnchainSigningAddress: hex.EncodeToString([]byte("aptos onchain-pub-key1")),
							ConfigPublicKey:       hexFrom32Byte(t, "aptos config-encryption-pub-key1"),
						},

						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: testPeerID(t, "peer-id1").String(),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{
				NodeID:         tt.fields.NodeID,
				SelToOCRConfig: tt.fields.SelToOCRConfig,
				PeerID:         tt.fields.PeerID,
				IsBootstrap:    tt.fields.IsBootstrap,
				MultiAddr:      tt.fields.MultiAddr,
				AdminAddr:      tt.fields.AdminAddr,
			}
			got, err := n.ChainConfigs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.ChainConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Len(t, got, len(tt.want))
			// returned order is not deterministic
			wantMap := make(map[string]*nodev1.ChainConfig)
			for i, w := range tt.want {
				wantMap[w.Chain.String()] = tt.want[i]
			}
			found := make(map[string]struct{})

			for i, g := range got {
				want, ok := wantMap[g.Chain.String()]
				require.True(t, ok, "unexpected chain config %v %v", g.Chain, wantMap)
				assert.EqualValues(t, want.Chain, g.Chain, "mismatched chain index %d", i)
				assert.Equal(t, want.NodeId, g.NodeId, "mismatched node id index %d", i)
				assert.EqualValues(t, want.Ocr2Config.OcrKeyBundle, g.Ocr2Config.OcrKeyBundle, "mismatched ocrkey bundle index %d", i)
				assert.EqualValues(t, want.Ocr2Config.P2PKeyBundle, g.Ocr2Config.P2PKeyBundle, "mismatched p2p bundle index %d", i)
				found[g.Chain.String()] = struct{}{}
			}
			// ensure we didn't miss any
			for w := range wantMap {
				_, ok := found[w]
				require.True(t, ok, "missing chain config %v", w)
			}
		})
	}
}

func testPeerID(t *testing.T, s string) p2pkey.PeerID {
	b := []byte(s)
	p := [32]byte{}
	copy(p[:], b)
	return p2pkey.PeerID(p)
}

func test32Byte(t *testing.T, s string) [32]byte {
	b := []byte(s)
	var out [32]byte
	copy(out[:], b)
	return out
}

func hexFrom32Byte(t *testing.T, s string) string {
	b := test32Byte(t, s)
	return hex.EncodeToString(b[:])
}

func Test_isValidMultiAddr(t *testing.T) {
	// Generate a p2p piece using p2pkey.MustNewV2XXXTestingOnly()
	seed := big.NewInt(123)
	p2p := p2pkey.MustNewV2XXXTestingOnly(seed).PeerID().String()

	// Create valid and invalid multi-address strings
	validMultiAddr := p2p[4:] + "@example.com:12345" // Remove "p2p_" prefix from p2p
	invalidMultiAddr1 := "invalid@address:12345"
	invalidMultiAddr2 := p2p[4:] + "@example.com:notanumber"
	invalidMultiAddr3 := "missingatsign.com:12345"
	invalidMultiAddr4 := p2p[4:] + "@example.com"
	invalidMultiAddr5 := "@missingp2p:123"

	// Test cases
	tests := []struct {
		name     string
		addr     string
		expected bool
	}{
		{"Valid MultiAddr", validMultiAddr, true},
		{"Invalid MultiAddr - Invalid Address", invalidMultiAddr1, false},
		{"Invalid MultiAddr - Non-numeric Port", invalidMultiAddr2, false},
		{"Invalid MultiAddr - Missing @", invalidMultiAddr3, false},
		{"Invalid MultiAddr - Missing Port", invalidMultiAddr4, false},
		{"Invalid MultiAddr - Missing p2p", invalidMultiAddr5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidMultiAddr(tt.addr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestNewNodeFromJD(t *testing.T) {
	wk := "node-workflow-key"
	type args struct {
		jdNode       *nodev1.Node
		chainConfigs []*nodev1.ChainConfig
	}
	tests := []struct {
		name     string
		args     args
		want     *Node
		checkErr func(t *testing.T, err error)
	}{
		{
			name: "ok",
			args: args{
				jdNode: &nodev1.Node{
					Id:          "node-id1",
					Name:        "node1",
					PublicKey:   "node-pub-key",
					WorkflowKey: &wk,
				},
				chainConfigs: []*nodev1.ChainConfig{
					{

						Chain: &nodev1.Chain{
							Type: nodev1.ChainType_CHAIN_TYPE_EVM,
							Id:   strconv.FormatUint(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID, 10),
						},
						NodeId: "node-id1",
						Ocr2Config: &nodev1.OCR2Config{
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundle1",
								OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
								OnchainSigningAddress: gethcommon.HexToAddress("0x1").Hex(),
								ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
							},
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId: testPeerID(t, "peer-id1").String(),
							},
						},
					},
				},
			},
			want: &Node{
				NodeID:      "node-id1",
				Name:        "node1",
				CSAKey:      "node-pub-key",
				WorkflowKey: "node-workflow-key",
				SelToOCRConfig: map[chain_selectors.ChainDetails]OCRConfig{
					{
						ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
						ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
					}: {
						KeyBundleID:               "bundle1",
						OffchainPublicKey:         test32Byte(t, "offchain-pub-key1"),
						OnchainPublicKey:          types2.OnchainPublicKey(gethcommon.HexToAddress("0x1").Bytes()),
						ConfigEncryptionPublicKey: types3.ConfigEncryptionPublicKey(test32Byte(t, "config-encryption-pub-key1")),

						PeerID: testPeerID(t, "peer-id1"),
					},
				},
				PeerID: testPeerID(t, "peer-id1"),
			},
		},
		{
			name: "no evm chain",
			args: args{
				jdNode: &nodev1.Node{
					Id:        "node-id1",
					Name:      "node1",
					PublicKey: "node-pub-key",
				},
				chainConfigs: []*nodev1.ChainConfig{
					{
						Chain: &nodev1.Chain{
							Type: nodev1.ChainType_CHAIN_TYPE_APTOS,
							Id:   strconv.FormatUint(chain_selectors.APTOS_TESTNET.ChainID, 10),
						},
						NodeId: "node-id1",
						Ocr2Config: &nodev1.OCR2Config{
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundle1",
								OffchainPublicKey:     hexFrom32Byte(t, "offchain-pub-key1"),
								OnchainSigningAddress: gethcommon.HexToAddress("0x1").Hex(),
								ConfigPublicKey:       hexFrom32Byte(t, "config-encryption-pub-key1"),
							},
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId: testPeerID(t, "peer-id1").String(),
							},
						},
					},
				},
			},
			checkErr: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrMissingEVMChain)
			},
			want: &Node{
				NodeID:         "node-id1",
				Name:           "node1",
				CSAKey:         "node-pub-key",
				SelToOCRConfig: map[chain_selectors.ChainDetails]OCRConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNodeFromJD(tt.args.jdNode, tt.args.chainConfigs)
			if tt.checkErr != nil {
				tt.checkErr(t, err)
				return
			}
			assert.Equal(t, tt.want.PeerID, got.PeerID)
			assert.Equal(t, tt.want.CSAKey, got.CSAKey)
			assert.Equal(t, tt.want.NodeID, got.NodeID)
			assert.Equal(t, tt.want.WorkflowKey, got.WorkflowKey)
			assert.Equal(t, tt.want.Name, got.Name)
			for k, v := range tt.want.SelToOCRConfig {
				assert.Equal(t, v, got.SelToOCRConfig[k])
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeFromJD() = %v, want %v", got, tt.want)
			}
		})
	}
}
