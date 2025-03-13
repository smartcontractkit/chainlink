package memory

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

func (j JobClient) EnableNode(ctx context.Context, in *nodev1.EnableNodeRequest, opts ...grpc.CallOption) (*nodev1.EnableNodeResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DisableNode(ctx context.Context, in *nodev1.DisableNodeRequest, opts ...grpc.CallOption) (*nodev1.DisableNodeResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j *JobClient) RegisterNode(ctx context.Context, in *nodev1.RegisterNodeRequest, opts ...grpc.CallOption) (*nodev1.RegisterNodeResponse, error) {
	if in == nil || in.GetPublicKey() == "" {
		return nil, errors.New("public key is required")
	}

	if _, exists := j.RegisteredNodes[in.GetPublicKey()]; exists {
		return nil, fmt.Errorf("node with Public Key %s is already registered", in.GetPublicKey())
	}

	var foundNode *Node
	for _, node := range j.nodeStore.list() {
		if node.Keys.CSA.ID() == in.GetPublicKey() {
			foundNode = node
			break
		}
	}

	if foundNode == nil {
		return nil, fmt.Errorf("node with Public Key %s is not known", in.GetPublicKey())
	}

	j.RegisteredNodes[in.GetPublicKey()] = *foundNode

	return &nodev1.RegisterNodeResponse{
		Node: &nodev1.Node{
			Id:          in.GetPublicKey(),
			PublicKey:   in.GetPublicKey(),
			IsEnabled:   true,
			IsConnected: true,
			Labels:      in.Labels,
		},
	}, nil
}

func (j JobClient) UpdateNode(ctx context.Context, in *nodev1.UpdateNodeRequest, opts ...grpc.CallOption) (*nodev1.UpdateNodeResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetNode(ctx context.Context, in *nodev1.GetNodeRequest, opts ...grpc.CallOption) (*nodev1.GetNodeResponse, error) {
	n, err := j.nodeStore.get(in.Id)
	if err != nil {
		return nil, err
	}
	return &nodev1.GetNodeResponse{
		Node: &nodev1.Node{
			Id:          in.Id,
			PublicKey:   n.Keys.CSA.PublicKeyString(),
			IsEnabled:   true,
			IsConnected: true,
		},
	}, nil
}

func (j JobClient) ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	var nodes []*nodev1.Node
	for id, n := range j.nodeStore.asMap() {
		node := &nodev1.Node{
			Id:          id,
			PublicKey:   n.Keys.CSA.ID(),
			IsEnabled:   true,
			IsConnected: true,
			Labels: []*ptypes.Label{
				{
					Key:   "p2p_id",
					Value: ptr(n.Keys.PeerID.String()),
				},
			},
		}
		if ApplyNodeFilter(in.Filter, node) {
			nodes = append(nodes, node)
		}
	}
	return &nodev1.ListNodesResponse{
		Nodes: nodes,
	}, nil
}

func (j JobClient) ListNodeChainConfigs(ctx context.Context, in *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {
	if in.Filter == nil {
		return nil, errors.New("filter is required")
	}
	if len(in.Filter.NodeIds) != 1 {
		return nil, errors.New("only one node id is supported")
	}
	n, err := j.nodeStore.get(in.Filter.NodeIds[0]) // j.Nodes[in.Filter.NodeIds[0]]
	if err != nil {
		return nil, fmt.Errorf("node id not found: %s", in.Filter.NodeIds[0])
	}
	chainConfigs, err := n.JDChainConfigs()
	if err != nil {
		return nil, err
	}
	return &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: chainConfigs,
	}, nil
}
