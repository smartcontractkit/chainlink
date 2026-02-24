package changeset

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
)

// ---------------------------------------------------------------------------
// Fake offchain.Client — only ListNodes and ListNodeChainConfigs are wired;
// every other method panics via the embedded nil interface.
// ---------------------------------------------------------------------------

type fakeOffchainClient struct {
	cldf_offchain.Client // satisfies the full interface; nil panics on unimplemented methods
	nodesByID            map[string]*fakeNodeInfo
}

type fakeNodeInfo struct {
	id           string
	name         string
	csaKey       string
	p2pID        string // full string, e.g. "p2p_12D3KooW…"
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
	// Build a set of wanted p2p IDs from the filter (if any).
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

	var out []*nodev1.Node
	for _, n := range f.nodesByID {
		if wantP2P != nil && !wantP2P[n.p2pID] {
			continue
		}
		p2pVal := n.p2pID
		out = append(out, &nodev1.Node{
			Id:        n.id,
			Name:      n.name,
			PublicKey: n.csaKey,
			IsEnabled: true,
			Labels:    []*ptypes.Label{{Key: "p2p_id", Value: &p2pVal}},
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

// ---------------------------------------------------------------------------
// Helper: build fake JD nodes + INodeInfoProviderNodeInfo from real testdata
// keys (first 4 Sepolia nodes from deployment/cre/ocr3/testdata).
// ---------------------------------------------------------------------------

type testNodeKeys struct {
	p2pID, offchainPubKey, onchainPubKey, transmitAccount string
	configEncryptionPubKey, keyBundleID, csaKey           string
}

func makeTestNodes(t *testing.T, chainSel uint64) ([]*fakeNodeInfo, []capabilities_registry_v2.INodeInfoProviderNodeInfo) {
	t.Helper()

	chainID, err := chainselectors.GetChainIDFromSelector(chainSel)
	require.NoError(t, err)

	// Keys sourced from deployment/cre/ocr3/testdata/testnet_wf_view.json (Sepolia).
	keys := []testNodeKeys{
		{
			p2pID:                  "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZBwkqyYv",
			offchainPubKey:         "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
			onchainPubKey:          "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
			transmitAccount:        "0x2877F08d9c5Cc9F401F730Fa418fAE563A9a2FF3",
			configEncryptionPubKey: "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f",
			keyBundleID:            "665a101d79d310cb0a5ebf695b06e8fc8082b5cbe62d7d362d80d47447a31fea",
			csaKey:                 "403b72f0b1b3b5f5a91bcfedb7f28599767502a04b5b7e067fcf3782e23eeb9c",
		},
		{
			p2pID:                  "p2p_12D3KooWCbDiL7sP9BVby5KaZqPpaVP1RBokoa9ShzH5WhkYX46v",
			offchainPubKey:         "255096a3b7ade10e29c648e0b407fc486180464f713446b1da04f013df6179c8",
			onchainPubKey:          "8258f4c4761cc445333017608044a204fd0c006a",
			transmitAccount:        "0x415aa1E9a1bcB3929ed92bFa1F9735Dc0D45AD31",
			configEncryptionPubKey: "2c45fec2320f6bcd36444529a86d9f8b4439499a5d8272dec9bcbbebb5e1bf01",
			keyBundleID:            "7a9b75510b8d09932b98142419bef52436ff725dd9395469473b487ef87fdfb0",
			csaKey:                 "28b91143ec9111796a7d63e14c1cf6bb01b4ed59667ab54f5bc72ebe49c881be",
		},
		{
			p2pID:                  "p2p_12D3KooWGDmBKZ7B3PynGrvfHTJMEecpjfHts9YK5NWk8oJuxcAo",
			offchainPubKey:         "dba3c61e5f8bec594be481bcaf67ecea0d1c2950edb15b158ce3dbc77877def3",
			onchainPubKey:          "d4dcc573e9d24a8b27a07bba670ba3a2ab36e5bb",
			transmitAccount:        "0xCea84bC1881F3cE14BA13Dc3a00DC1Ff3D553fF0",
			configEncryptionPubKey: "ee466234b3b2f65b13c848b17aa6a8d4e0aa0311d3bf8e77a64f20b04ed48d39",
			keyBundleID:            "1d20490fe469dd6af3d418cc310a6e835181fa13e8dc80156bcbe302b7afcd34",
			csaKey:                 "7a166fbc816eb4a4dcb620d11c3ccac5c085d56b1972374100116f87619debb8",
		},
		{
			p2pID:                  "p2p_12D3KooWCcVLytqinD8xMn27NvomcQhj2mqMVzyGemz6oPwv1SMT",
			offchainPubKey:         "b4c4993d6c15fee63800db901a8b35fa419057610962caab1c1d7bed55709127",
			onchainPubKey:          "6607c140e558631407f33bafbabd103863cee876",
			transmitAccount:        "0xA9eFB53c513E413762b2Be5299D161d8E6e7278e",
			configEncryptionPubKey: "63375a3d175364bd299e7cecf352cb3e469dd30116cf1418f2b7571fb46c4a4b",
			keyBundleID:            "8843b5db0608f92dac38ca56775766a08db9ee82224a19595d04bd6c58b38fbd",
			csaKey:                 "487901e0c0a9d3c66e7cfc50f3a9e3cdbfdf1b0107273d73d94a91d278545516",
		},
	}

	fakeNodes := make([]*fakeNodeInfo, 0, len(keys))
	registryNodes := make([]capabilities_registry_v2.INodeInfoProviderNodeInfo, 0, len(keys))

	for i, k := range keys {
		peerID, pErr := p2pkey.MakePeerID(k.p2pID)
		require.NoError(t, pErr)

		fakeNodes = append(fakeNodes, &fakeNodeInfo{
			id:     fmt.Sprintf("node_%02d", i),
			name:   fmt.Sprintf("test-node-%d", i),
			csaKey: k.csaKey,
			p2pID:  k.p2pID,
			chainConfigs: []*nodev1.ChainConfig{
				{
					Chain: &nodev1.Chain{
						Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						Id:   chainID,
					},
					Ocr2Config: &nodev1.OCR2Config{
						OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
							OffchainPublicKey:     k.offchainPubKey,
							OnchainSigningAddress: k.onchainPubKey,
							ConfigPublicKey:       k.configEncryptionPubKey,
							BundleId:              k.keyBundleID,
						},
						P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
							PeerId: k.p2pID,
						},
						IsBootstrap: false,
					},
					AccountAddress: k.transmitAccount,
				},
			},
		})

		registryNodes = append(registryNodes, capabilities_registry_v2.INodeInfoProviderNodeInfo{
			P2pId: peerID,
		})
	}

	return fakeNodes, registryNodes
}

// oracleConfigMap returns a minimal OracleConfig as map[string]any suitable for
// JSON roundtrip into ocr3.OracleConfig. numNodes controls TransmissionSchedule.
func oracleConfigMap(numNodes int) map[string]any {
	return map[string]any{
		"uniqueReports":                     true,
		"deltaProgressMillis":               5000,
		"deltaResendMillis":                 5000,
		"deltaInitialMillis":                5000,
		"deltaRoundMillis":                  2000,
		"deltaGraceMillis":                  500,
		"deltaCertifiedCommitRequestMillis": 1000,
		"deltaStageMillis":                  30000,
		"maxRoundsPerEpoch":                 10,
		"transmissionSchedule":              []any{numNodes},
		"maxDurationQueryMillis":            1000,
		"maxDurationObservationMillis":      1000,
		"maxDurationShouldAcceptMillis":     1000,
		"maxDurationShouldTransmitMillis":   1000,
		"maxFaultyOracles":                  1,
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestExpandOCR3Configs(t *testing.T) {
	t.Parallel()

	t.Run("no ocr3 configs to expand", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{CapabilityID: "write-chain@1.0.0", Config: map[string]any{"key": "value"}},
		}
		err := expandOCR3Configs(env, 0, nil, capConfigs, nil)
		require.NoError(t, err)
	})

	t.Run("already-processed entries are skipped", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"signers":               []any{"AQIDBA=="},
							"offchainConfig":        "base64encodedstring",
							"offchainConfigVersion": 30,
						},
					},
				},
			},
		}
		err := expandOCR3Configs(env, 0, nil, capConfigs, nil)
		require.NoError(t, err)
		entry := capConfigs[0].Config["ocr3Configs"].(map[string]any)["__default__"].(map[string]any)
		assert.Equal(t, "base64encodedstring", entry["offchainConfig"])
	})

	t.Run("entries without offchainConfig map are skipped", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"f": 1,
						},
					},
				},
			},
		}
		err := expandOCR3Configs(env, 0, nil, capConfigs, nil)
		require.NoError(t, err)
	})

	t.Run("non-map entries in ocr3Configs are skipped", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": "not-a-map",
					},
				},
			},
		}
		err := expandOCR3Configs(env, 0, nil, capConfigs, nil)
		require.NoError(t, err)
	})

	t.Run("parse error in oracle config", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"offchainConfig": map[string]any{
								"deltaProgressMillis": "not-a-number",
							},
						},
					},
				},
			},
		}
		err := expandOCR3Configs(env, 0, nil, capConfigs, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse oracle config")
	})

	t.Run("configCountFn error is propagated", func(t *testing.T) {
		t.Parallel()
		env := cldf.Environment{Logger: logger.Test(t)}
		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"offchainConfig": map[string]any{},
						},
					},
				},
			},
		}
		countErr := errors.New("config count lookup failed")
		err := expandOCR3Configs(env, 0, nil, capConfigs, func(_, _ string) (uint64, error) {
			return 0, countErr
		})
		require.ErrorIs(t, err, countErr)
	})

	t.Run("successful expansion", func(t *testing.T) {
		t.Parallel()

		chainSel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		fakeNodes, registryNodes := makeTestNodes(t, chainSel)
		client := newFakeOffchainClient(fakeNodes)

		env := cldf.Environment{
			Logger:     logger.Test(t),
			Offchain:   client,
			OCRSecrets: ocr.XXXGenerateTestOCRSecrets(),
		}

		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"offchainConfig": oracleConfigMap(len(registryNodes)),
						},
					},
				},
			},
		}

		wantConfigCount := uint64(42)
		err := expandOCR3Configs(env, chainSel, registryNodes, capConfigs, func(capID, key string) (uint64, error) {
			assert.Equal(t, "consensus@1.0.0", capID)
			assert.Equal(t, "__default__", key)
			return wantConfigCount, nil
		})
		require.NoError(t, err)

		expanded := capConfigs[0].Config["ocr3Configs"].(map[string]any)["__default__"].(map[string]any)

		signers, ok := expanded["signers"].([]string)
		require.True(t, ok, "signers should be []string")
		assert.Len(t, signers, len(registryNodes))

		transmitters, ok := expanded["transmitters"].([]string)
		require.True(t, ok, "transmitters should be []string")
		assert.Len(t, transmitters, len(registryNodes))

		assert.Equal(t, uint32(1), expanded["f"])
		assert.Equal(t, uint64(30), expanded["offchainConfigVersion"])
		assert.Equal(t, wantConfigCount, expanded["configCount"])

		offchainCfg, ok := expanded["offchainConfig"].(string)
		require.True(t, ok, "offchainConfig should be a base64 string after expansion")
		assert.NotEmpty(t, offchainCfg)
	})

	t.Run("mixed capabilities - only expandable ones are processed", func(t *testing.T) {
		t.Parallel()

		chainSel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		fakeNodes, registryNodes := makeTestNodes(t, chainSel)
		client := newFakeOffchainClient(fakeNodes)

		env := cldf.Environment{
			Logger:     logger.Test(t),
			Offchain:   client,
			OCRSecrets: ocr.XXXGenerateTestOCRSecrets(),
		}

		capConfigs := []ocr3CapConfig{
			{
				CapabilityID: "write-chain@1.0.0",
				Config:       map[string]any{"key": "value"},
			},
			{
				CapabilityID: "consensus@1.0.0",
				Config: map[string]any{
					"ocr3Configs": map[string]any{
						"__default__": map[string]any{
							"offchainConfig": oracleConfigMap(len(registryNodes)),
						},
					},
				},
			},
		}

		err := expandOCR3Configs(env, chainSel, registryNodes, capConfigs, func(_, _ string) (uint64, error) {
			return 1, nil
		})
		require.NoError(t, err)

		// write-chain config untouched
		assert.Equal(t, "value", capConfigs[0].Config["key"])

		// consensus config expanded
		expanded := capConfigs[1].Config["ocr3Configs"].(map[string]any)["__default__"].(map[string]any)
		_, hasSigners := expanded["signers"]
		assert.True(t, hasSigners, "consensus entry should have been expanded")
	})
}
