package memory

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	chainsel "github.com/smartcontractkit/chain-selectors"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
)

type JobClient struct {
	Nodes           map[string]Node
	RegisteredNodes map[string]Node
}

func (j JobClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	// TODO CCIP-3108  implement me
	panic("implement me")
}

func (j JobClient) UpdateJob(ctx context.Context, in *jobv1.UpdateJobRequest, opts ...grpc.CallOption) (*jobv1.UpdateJobResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DisableNode(ctx context.Context, in *nodev1.DisableNodeRequest, opts ...grpc.CallOption) (*nodev1.DisableNodeResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) EnableNode(ctx context.Context, in *nodev1.EnableNodeRequest, opts ...grpc.CallOption) (*nodev1.EnableNodeResponse, error) {
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
	for _, node := range j.Nodes {
		if node.Keys.CSA.ID() == in.GetPublicKey() {
			foundNode = &node
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

func (j JobClient) GetKeypair(ctx context.Context, in *csav1.GetKeypairRequest, opts ...grpc.CallOption) (*csav1.GetKeypairResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (j JobClient) ListKeypairs(ctx context.Context, in *csav1.ListKeypairsRequest, opts ...grpc.CallOption) (*csav1.ListKeypairsResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetNode(ctx context.Context, in *nodev1.GetNodeRequest, opts ...grpc.CallOption) (*nodev1.GetNodeResponse, error) {
	n, ok := j.Nodes[in.Id]
	if !ok {
		return nil, errors.New("node not found")
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
	for id, n := range j.Nodes {
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
	n, ok := j.Nodes[in.Filter.NodeIds[0]]
	if !ok {
		return nil, fmt.Errorf("node id not found: %s", in.Filter.NodeIds[0])
	}
	var chainConfigs []*nodev1.ChainConfig
	for _, selector := range n.Chains {
		family, err := chainsel.GetSelectorFamily(selector)
		if err != nil {
			return nil, err
		}

		// NOTE: this supports non-EVM too
		chainID, err := chainsel.GetChainIDFromSelector(selector)
		if err != nil {
			return nil, err
		}

		var ocrtype chaintype.ChainType
		switch family {
		case chainsel.FamilyEVM:
			ocrtype = chaintype.EVM
		case chainsel.FamilySolana:
			ocrtype = chaintype.Solana
		case chainsel.FamilyStarknet:
			ocrtype = chaintype.StarkNet
		case chainsel.FamilyCosmos:
			ocrtype = chaintype.Cosmos
		case chainsel.FamilyAptos:
			ocrtype = chaintype.Aptos
		default:
			return nil, fmt.Errorf("Unsupported chain family %v", family)
		}

		bundle := n.Keys.OCRKeyBundles[ocrtype]
		offpk := bundle.OffchainPublicKey()
		cpk := bundle.ConfigEncryptionPublicKey()

		keyBundle := &nodev1.OCR2Config_OCRKeyBundle{
			BundleId:              bundle.ID(),
			ConfigPublicKey:       common.Bytes2Hex(cpk[:]),
			OffchainPublicKey:     common.Bytes2Hex(offpk[:]),
			OnchainSigningAddress: bundle.OnChainPublicKey(),
		}

		var ctype nodev1.ChainType
		switch family {
		case chainsel.FamilyEVM:
			ctype = nodev1.ChainType_CHAIN_TYPE_EVM
		case chainsel.FamilySolana:
			ctype = nodev1.ChainType_CHAIN_TYPE_SOLANA
		case chainsel.FamilyStarknet:
			ctype = nodev1.ChainType_CHAIN_TYPE_STARKNET
		case chainsel.FamilyAptos:
			ctype = nodev1.ChainType_CHAIN_TYPE_APTOS
		default:
			panic(fmt.Sprintf("Unsupported chain family %v", family))
		}

		transmitter := n.Keys.Transmitters[selector]

		chainConfigs = append(chainConfigs, &nodev1.ChainConfig{
			Chain: &nodev1.Chain{
				Id:   chainID,
				Type: ctype,
			},
			AccountAddress: transmitter,
			AdminAddress:   transmitter,
			Ocr1Config:     nil,
			Ocr2Config: &nodev1.OCR2Config{
				Enabled:     true,
				IsBootstrap: n.IsBoostrap,
				P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
					PeerId: n.Keys.PeerID.String(),
				},
				OcrKeyBundle:     keyBundle,
				Multiaddr:        n.Addr.String(),
				Plugins:          nil,
				ForwarderAddress: ptr(""),
			},
		})
	}
	return &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: chainConfigs,
	}, nil
}

func (j JobClient) GetJob(ctx context.Context, in *jobv1.GetJobRequest, opts ...grpc.CallOption) (*jobv1.GetJobResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) GetProposal(ctx context.Context, in *jobv1.GetProposalRequest, opts ...grpc.CallOption) (*jobv1.GetProposalResponse, error) {
	// we are using proposal id as job id
	// refer to ListJobs and ProposeJobs for the assignment of proposal id
	for _, node := range j.Nodes {
		jobs, _, err := node.App.JobORM().FindJobs(ctx, 0, 1000)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs {
			if job.ExternalJobID.String() == in.Id {
				specBytes, err := toml.Marshal(job.CCIPSpec)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal job spec: %w", err)
				}
				return &jobv1.GetProposalResponse{
					Proposal: &jobv1.Proposal{
						Id:     job.ExternalJobID.String(),
						Status: jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED,
						Spec:   string(specBytes),
						JobId:  job.ExternalJobID.String(),
					},
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("job not found: %s", in.Id)
}

func (j JobClient) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest, opts ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	jobResponse := make([]*jobv1.Job, 0)
	for _, req := range in.Filter.NodeIds {
		if _, ok := j.Nodes[req]; !ok {
			return nil, fmt.Errorf("node not found: %s", req)
		}
		n := j.Nodes[req]
		jobs, _, err := n.App.JobORM().FindJobs(ctx, 0, 1000)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs {
			jobResponse = append(jobResponse, &jobv1.Job{
				Id:     string(job.ID),
				Uuid:   job.ExternalJobID.String(),
				NodeId: req,
				// based on the current implementation, there is only one proposal per job
				// see ProposeJobs for ProposalId assignment
				ProposalIds: []string{job.ExternalJobID.String()},
				CreatedAt:   timestamppb.New(job.CreatedAt),
				UpdatedAt:   timestamppb.New(job.CreatedAt),
			})
		}
	}
	return &jobv1.ListJobsResponse{
		Jobs: jobResponse,
	}, nil
}

func (j JobClient) ListProposals(ctx context.Context, in *jobv1.ListProposalsRequest, opts ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	n := j.Nodes[in.NodeId]
	// TODO: Use FMS
	jb, err := validate.ValidatedCCIPSpec(in.Spec)
	if err != nil {
		if !strings.Contains(err.Error(), "the only supported type is currently 'ccip'") {
			return nil, err
		}
		// check if it's offchainreporting2 job
		jb, err = ocr2validate.ValidatedOracleSpecToml(
			ctx,
			n.App.GetConfig().OCR2(),
			n.App.GetConfig().Insecure(),
			in.Spec,
			nil, // not required for validation
		)
		if err != nil {
			if !strings.Contains(err.Error(), "the only supported type is currently 'offchainreporting2'") {
				return nil, err
			}
			// check if it's bootstrap job
			jb, err = ocrbootstrap.ValidatedBootstrapSpecToml(in.Spec)
			if err != nil {
				return nil, fmt.Errorf("failed to validate job spec only ccip, bootstrap and offchainreporting2 are supported: %w", err)
			}
		}
	}
	err = n.App.AddJobV2(ctx, &jb)
	if err != nil {
		return nil, err
	}
	return &jobv1.ProposeJobResponse{Proposal: &jobv1.Proposal{
		// make the proposal id the same as the job id for further reference
		// if you are changing this make sure to change the GetProposal and ListJobs method implementation
		Id: jb.ExternalJobID.String(),
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
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	// TODO CCIP-3108 implement me
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
	return &JobClient{nodesByPeerID, make(map[string]Node)}
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
