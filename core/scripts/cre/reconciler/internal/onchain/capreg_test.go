package onchain

import (
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestCapRegWorkerNodes_ExcludesBootstrap(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	workers := capRegWorkerNodes(topology, nil)
	require.Len(t, workers, 4)
	for _, w := range workers {
		require.Equal(t, "workflow", w.donName)
	}
}

func TestVerifyCapRegNodeInfo_MissingJDClient(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	err := verifyCapRegNodeInfo(nil, 1337, topology, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "JD client is required")
}

func TestVerifyCapRegNodeInfo_MissingRegistryChainConfig(t *testing.T) {
	t.Parallel()

	topology := testSingleWorkerTopology(t)
	workers := capRegWorkerNodes(topology, nil)
	require.Len(t, workers, 1)

	wfKey := "5193f72fc7b4323a86088fb0acb4e4494ae351920b3944bd726a59e8dbcdd45f"
	p2pID := workers[0].worker.Keys.PeerID()
	chainConfigsByP2P := map[string][]*node.ChainConfig{
		p2pID: {{
			Chain: &node.Chain{Type: node.ChainType_CHAIN_TYPE_EVM, Id: "1"},
			Ocr2Config: &node.OCR2Config{
				OcrKeyBundle: &node.OCR2Config_OCRKeyBundle{
					OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
					OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
					ConfigPublicKey:       wfKey,
					BundleId:              "665a101d79d310cb0a5ebf695b06e8fc8082b5cbe62d7d362d80d47447a31fea",
				},
				P2PKeyBundle: &node.OCR2Config_P2PKeyBundle{PeerId: p2pID},
			},
			AccountAddress: "0x2877F08d9c5Cc9F401F730Fa418fAE563A9a2FF3",
		}},
	}
	fake := newCapRegVerifyOffchainClient(chainConfigsByP2P, map[string]string{p2pID: wfKey})

	err := verifyCapRegNodeInfo(fake, 999001, topology, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "chain selector 999001")
	require.Contains(t, err.Error(), "missing registry-chain OCR config")
}

type capRegVerifyOffchainClient struct {
	offchain.Client
	chainConfigsByP2P map[string][]*node.ChainConfig
	workflowKeys      map[string]string
}

func newCapRegVerifyOffchainClient(chainConfigsByP2P map[string][]*node.ChainConfig, workflowKeys map[string]string) *capRegVerifyOffchainClient {
	return &capRegVerifyOffchainClient{
		chainConfigsByP2P: chainConfigsByP2P,
		workflowKeys:      workflowKeys,
	}
}

func (f *capRegVerifyOffchainClient) ListNodes(_ context.Context, in *node.ListNodesRequest, _ ...grpc.CallOption) (*node.ListNodesResponse, error) {
	var wantP2P map[string]bool
	if in.Filter != nil {
		for _, sel := range in.Filter.Selectors {
			if sel.Key == "p2p_id" && sel.Value != nil {
				wantP2P = make(map[string]bool)
				for v := range strings.SplitSeq(*sel.Value, ",") {
					wantP2P[v] = true
				}
			}
		}
	}

	out := make([]*node.Node, 0, len(f.chainConfigsByP2P))
	for p2pID := range f.chainConfigsByP2P {
		if wantP2P != nil && !wantP2P[p2pID] {
			continue
		}
		wfKey := f.workflowKeys[p2pID]
		out = append(out, &node.Node{
			Id:          "node_" + p2pID,
			Name:        p2pID,
			PublicKey:   p2pID,
			WorkflowKey: &wfKey,
			IsEnabled:   true,
		})
	}
	return &node.ListNodesResponse{Nodes: out}, nil
}

func (f *capRegVerifyOffchainClient) ListNodeChainConfigs(_ context.Context, in *node.ListNodeChainConfigsRequest, _ ...grpc.CallOption) (*node.ListNodeChainConfigsResponse, error) {
	if in.Filter == nil || len(in.Filter.NodeIds) == 0 {
		return nil, errors.New("node IDs required")
	}
	for _, configs := range f.chainConfigsByP2P {
		return &node.ListNodeChainConfigsResponse{ChainConfigs: configs}, nil
	}
	return &node.ListNodeChainConfigsResponse{}, nil
}
