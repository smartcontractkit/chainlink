package jdtestutils

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	"github.com/smartcontractkit/chainlink/deployment/environment/test"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
	"github.com/smartcontractkit/chainlink/deployment/utils/nodetestutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

var _ cldf_offchain.Client = &JobClient{}

type JobClient struct {
	RegisteredNodes map[string]nodetestutils.Node
	nodeStore
	*test.JobServiceClient
}

func NewMemoryJobClient(nodesByPeerID map[string]nodetestutils.Node) *JobClient {
	m := make(map[string]*nodetestutils.Node)
	for id, node := range nodesByPeerID {
		m[id] = &node
	}
	ns := newMapNodeStore(m)
	jg := &jobApproverGetter{s: ns}
	return &JobClient{
		RegisteredNodes:  make(map[string]nodetestutils.Node),
		JobServiceClient: test.NewJobServiceClient(jg),
		nodeStore:        ns,
	}
}

func (j JobClient) GetKeypair(ctx context.Context, in *csav1.GetKeypairRequest, opts ...grpc.CallOption) (*csav1.GetKeypairResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (j JobClient) ListKeypairs(ctx context.Context, in *csav1.ListKeypairsRequest, opts ...grpc.CallOption) (*csav1.ListKeypairsResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ReplayLogs(ctx context.Context, selectorToBlock map[uint64]uint64) error {
	for _, node := range j.list() {
		if err := node.ReplayLogs(ctx, selectorToBlock); err != nil {
			return err
		}
	}
	return nil
}

// Checks if a filter exists in DB for event name in all nodes
func (j JobClient) IsLogFilterRegistered(ctx context.Context, chainSel uint64, eventName string, address []byte) (bool, error) {
	for _, node := range j.list() {
		if node.IsBoostrap {
			continue
		}
		registered, err := node.IsLogFilterRegistered(ctx, chainSel, eventName, address)
		if err != nil || !registered {
			return false, err
		}
	}
	return true, nil
}

func ApplyNodeFilter(filter *nodev1.ListNodesRequest_Filter, node *nodev1.Node) bool {
	if filter == nil {
		return true
	}
	if len(filter.Ids) > 0 {
		idx := slices.IndexFunc(filter.Ids, func(id string) bool {
			return node.Id == id
		})
		if idx < 0 {
			return false
		}
	}
	if len(filter.PublicKeys) > 0 {
		idx := slices.IndexFunc(filter.PublicKeys, func(pk string) bool {
			return node.PublicKey == pk
		})
		if idx < 0 {
			return false
		}
	}
	for _, selector := range filter.Selectors {
		idx := slices.IndexFunc(node.Labels, func(label *ptypes.Label) bool {
			return label.Key == selector.Key
		})
		if idx < 0 {
			return false
		}
		label := node.Labels[idx]

		switch selector.Op {
		case ptypes.SelectorOp_IN:
			values := strings.Split(*selector.Value, ",")
			found := slices.Contains(values, *label.Value)
			if !found {
				return false
			}
		case ptypes.SelectorOp_EQ:
			if *label.Value != *selector.Value {
				return false
			}
		case ptypes.SelectorOp_EXIST:
			// do nothing
		default:
			panic("unimplemented selector")
		}
	}
	return true
}

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

	var foundNode *nodetestutils.Node
	for _, node := range j.list() {
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
			Name:        foundNode.Name,
			PublicKey:   in.GetPublicKey(),
			IsEnabled:   true,
			IsConnected: true,
			Labels:      in.Labels,
		},
	}, nil
}

// UpdateNode only updates the labels of the node.
// WARNING: The provided input will *overwrite* the existing fields, it won't extend them.
// TODO: Updating the PublicKey is not supported in this implementation.
func (j JobClient) UpdateNode(ctx context.Context, in *nodev1.UpdateNodeRequest, opts ...grpc.CallOption) (*nodev1.UpdateNodeResponse, error) {
	node, err := j.get(in.Id)
	if err != nil {
		return nil, fmt.Errorf("node with ID %s not found", in.Id)
	}

	node.ID = in.Id
	node.Name = in.Name
	node.Labels = in.Labels
	err = j.put(in.Id, node)
	if err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	return &nodev1.UpdateNodeResponse{
		Node: &nodev1.Node{
			Id:          in.Id,
			Name:        in.Name,
			PublicKey:   node.Keys.CSA.ID(),
			IsEnabled:   true,
			IsConnected: true,
			Labels:      in.Labels,
		},
	}, nil
}

func (j JobClient) GetNode(ctx context.Context, in *nodev1.GetNodeRequest, opts ...grpc.CallOption) (*nodev1.GetNodeResponse, error) {
	n, err := j.get(in.Id)
	if err != nil {
		return nil, err
	}
	return &nodev1.GetNodeResponse{
		Node: &nodev1.Node{
			Id:          in.Id,
			Name:        n.Name,
			PublicKey:   n.Keys.CSA.PublicKeyString(),
			IsEnabled:   true,
			IsConnected: true,
			Labels:      n.Labels,
		},
	}, nil
}

func (j JobClient) ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	var nodes []*nodev1.Node
	for id, n := range j.asMap() {
		p2pIDLabel := &ptypes.Label{
			Key:   "p2p_id",
			Value: pointer.To(n.Keys.PeerID.String()),
		}
		node := &nodev1.Node{
			Id:          id,
			Name:        n.Name,
			PublicKey:   n.Keys.CSA.ID(),
			IsEnabled:   true,
			IsConnected: true,
			Labels:      append(n.Labels, p2pIDLabel),
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
	n, err := j.get(in.Filter.NodeIds[0]) // j.Nodes[in.Filter.NodeIds[0]]
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

type JobApprover interface {
	AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error
}

type autoApprovalNode struct {
	*nodetestutils.Node
}

var _ JobApprover = &autoApprovalNode{}

func (q *autoApprovalNode) AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error {
	appProposalID, err := q.App.GetFeedsService().ProposeJob(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to propose job: %w", err)
	}
	// auto approve
	proposedSpec, err := q.App.GetFeedsService().ListSpecsByJobProposalIDs(ctx, []int64{appProposalID})
	if err != nil {
		return fmt.Errorf("failed to list specs: %w", err)
	}
	// possible to have multiple specs for the same job proposal id; take the last one
	if len(proposedSpec) == 0 {
		return fmt.Errorf("no specs found for job proposal id: %d", appProposalID)
	}
	err = q.App.GetFeedsService().ApproveSpec(ctx, proposedSpec[len(proposedSpec)-1].ID, true)
	if err != nil {
		return fmt.Errorf("failed to approve job: %w", err)
	}
	return nil
}

type jobApproverGetter struct {
	s nodeStore
}

func (w *jobApproverGetter) Get(nodeID string) (test.JobApprover, error) {
	node, err := w.s.get(nodeID)
	if err != nil {
		return nil, err
	}
	return &autoApprovalNode{node}, nil
}
