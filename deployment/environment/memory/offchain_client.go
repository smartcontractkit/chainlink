package memory

import (
	"context"
	"slices"
	"strings"

	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/test"
)

var _ deployment.OffchainClient = &JobClient{}

type JobClient struct {
	RegisteredNodes map[string]Node
	nodeStore
	*test.JobServiceClient
}

func NewMemoryJobClient(nodesByPeerID map[string]Node) *JobClient {
	m := make(map[string]*Node)
	for id, node := range nodesByPeerID {
		m[id] = &node
	}
	ns := newMapNodeStore(m)
	jg := &jobApproverGetter{s: ns}
	return &JobClient{
		RegisteredNodes:  make(map[string]Node),
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
	for _, node := range j.nodeStore.list() {
		if err := node.ReplayLogs(ctx, selectorToBlock); err != nil {
			return err
		}
	}
	return nil
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
