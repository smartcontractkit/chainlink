package memory

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
)

type JobClient struct {
	Nodes map[string]Node
}

func (j JobClient) UpdateJob(ctx context.Context, in *jobv1.UpdateJobRequest, opts ...grpc.CallOption) (*jobv1.UpdateJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DisableNode(ctx context.Context, in *nodev1.DisableNodeRequest, opts ...grpc.CallOption) (*nodev1.DisableNodeResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) EnableNode(ctx context.Context, in *nodev1.EnableNodeRequest, opts ...grpc.CallOption) (*nodev1.EnableNodeResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) RegisterNode(ctx context.Context, in *nodev1.RegisterNodeRequest, opts ...grpc.CallOption) (*nodev1.RegisterNodeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobClient) UpdateNode(ctx context.Context, in *nodev1.UpdateNodeRequest, opts ...grpc.CallOption) (*nodev1.UpdateNodeResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetKeypair(ctx context.Context, in *csav1.GetKeypairRequest, opts ...grpc.CallOption) (*csav1.GetKeypairResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobClient) ListKeypairs(ctx context.Context, in *csav1.ListKeypairsRequest, opts ...grpc.CallOption) (*csav1.ListKeypairsResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetNode(ctx context.Context, in *nodev1.GetNodeRequest, opts ...grpc.CallOption) (*nodev1.GetNodeResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	//TODO CCIP-3108
	var fiterIds map[string]struct{}
	include := func(id string) bool {
		if in.Filter == nil || len(in.Filter.Ids) == 0 {
			return true
		}
		// lazy init
		if len(fiterIds) == 0 {
			for _, id := range in.Filter.Ids {
				fiterIds[id] = struct{}{}
			}
		}
		_, ok := fiterIds[id]
		return ok
	}
	var nodes []*nodev1.Node
	for id, n := range j.Nodes {
		if include(id) {
			nodes = append(nodes, &nodev1.Node{
				Id:          id,
				PublicKey:   n.Keys.OCRKeyBundle.ID(), // is this the correct val?
				IsEnabled:   true,
				IsConnected: true,
			})
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
	n, ok := j.Nodes[in.Filter.NodeIds[0]]
	if !ok {
		return nil, fmt.Errorf("node id not found: %s", in.Filter.NodeIds[0])
	}
	offpk := n.Keys.OCRKeyBundle.OffchainPublicKey()
	cpk := n.Keys.OCRKeyBundle.ConfigEncryptionPublicKey()
	var chainConfigs []*nodev1.ChainConfig
	for evmChainID, transmitter := range n.Keys.TransmittersByEVMChainID {
		chainConfigs = append(chainConfigs, &nodev1.ChainConfig{
			Chain: &nodev1.Chain{
				Id:   strconv.Itoa(int(evmChainID)),
				Type: nodev1.ChainType_CHAIN_TYPE_EVM,
			},
			AccountAddress: transmitter.String(),
			AdminAddress:   "",
			Ocr1Config:     nil,
			Ocr2Config: &nodev1.OCR2Config{
				Enabled:     true,
				IsBootstrap: n.IsBoostrap,
				P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
					PeerId: n.Keys.PeerID.String(),
				},
				OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
					BundleId:              n.Keys.OCRKeyBundle.ID(),
					ConfigPublicKey:       common.Bytes2Hex(cpk[:]),
					OffchainPublicKey:     common.Bytes2Hex(offpk[:]),
					OnchainSigningAddress: n.Keys.OCRKeyBundle.OnChainPublicKey(),
				},
				Multiaddr:        n.Addr.String(),
				Plugins:          nil,
				ForwarderAddress: ptr(""),
			},
		})
	}

	// TODO: I think we can pull it from the feeds manager.
	return &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: chainConfigs,
	}, nil
}

func (j JobClient) GetJob(ctx context.Context, in *jobv1.GetJobRequest, opts ...grpc.CallOption) (*jobv1.GetJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetProposal(ctx context.Context, in *jobv1.GetProposalRequest, opts ...grpc.CallOption) (*jobv1.GetProposalResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest, opts ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ListProposals(ctx context.Context, in *jobv1.ListProposalsRequest, opts ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	n := j.Nodes[in.NodeId]
	// TODO: Use FMS
	jb, err := validate.ValidatedCCIPSpec(in.Spec)
	if err != nil {
		return nil, err
	}
	err = n.App.AddJobV2(ctx, &jb)
	if err != nil {
		return nil, err
	}
	return &jobv1.ProposeJobResponse{Proposal: &jobv1.Proposal{
		Id: "",
		// Auto approve for now
		Status:             jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED,
		DeliveryStatus:     jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED,
		Spec:               in.Spec,
		JobId:              jb.ExternalJobID.String(),
		CreatedAt:          nil,
		UpdatedAt:          nil,
		AckedAt:            nil,
		ResponseReceivedAt: nil,
	}}, nil
}

func (j JobClient) RevokeJob(ctx context.Context, in *jobv1.RevokeJobRequest, opts ...grpc.CallOption) (*jobv1.RevokeJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ReplayLogs(selectorToBlock map[uint64]uint64) error {
	for _, node := range j.Nodes {
		if err := node.ReplayLogs(selectorToBlock); err != nil {
			return err
		}
	}
	return nil
}

func NewMemoryJobClient(nodesByPeerID map[string]Node) *JobClient {
	return &JobClient{nodesByPeerID}
}
