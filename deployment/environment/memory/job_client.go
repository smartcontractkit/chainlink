package memory

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"

	chainsel "github.com/smartcontractkit/chain-selectors"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
)

type JobClient struct {
	Nodes map[string]Node
}

func (j JobClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	//TODO CCIP-3108  implement me
	panic("implement me")
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
	//TODO CCIP-3108
	include := func(node *nodev1.Node) bool {
		if in.Filter == nil {
			return true
		}
		if len(in.Filter.Ids) > 0 {
			idx := slices.IndexFunc(in.Filter.Ids, func(id string) bool {
				return node.Id == id
			})
			if idx < 0 {
				return false
			}
		}
		for _, selector := range in.Filter.Selectors {
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
			default:
				panic("unimplemented selector")
			}
		}
		return true
	}
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
		if include(node) {
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
	evmBundle := n.Keys.OCRKeyBundles[chaintype.EVM]
	offpk := evmBundle.OffchainPublicKey()
	cpk := evmBundle.ConfigEncryptionPublicKey()

	evmKeyBundle := &nodev1.OCR2Config_OCRKeyBundle{
		BundleId:              evmBundle.ID(),
		ConfigPublicKey:       common.Bytes2Hex(cpk[:]),
		OffchainPublicKey:     common.Bytes2Hex(offpk[:]),
		OnchainSigningAddress: evmBundle.OnChainPublicKey(),
	}

	var chainConfigs []*nodev1.ChainConfig
	for evmChainID, transmitter := range n.Keys.TransmittersByEVMChainID {
		chainConfigs = append(chainConfigs, &nodev1.ChainConfig{
			Chain: &nodev1.Chain{
				Id:   strconv.Itoa(int(evmChainID)),
				Type: nodev1.ChainType_CHAIN_TYPE_EVM,
			},
			AccountAddress: transmitter.String(),
			AdminAddress:   transmitter.String(), // TODO: custom address
			Ocr1Config:     nil,
			Ocr2Config: &nodev1.OCR2Config{
				Enabled:     true,
				IsBootstrap: n.IsBoostrap,
				P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
					PeerId: n.Keys.PeerID.String(),
				},
				OcrKeyBundle:     evmKeyBundle,
				Multiaddr:        n.Addr.String(),
				Plugins:          nil,
				ForwarderAddress: ptr(""),
			},
		})
	}
	for _, selector := range n.Chains {
		family, err := chainsel.GetSelectorFamily(selector)
		if err != nil {
			return nil, err
		}
		chainID, err := chainsel.ChainIdFromSelector(selector)
		if err != nil {
			return nil, err
		}

		if family == chainsel.FamilyEVM {
			// already handled above
			continue
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
			panic(fmt.Sprintf("Unsupported chain family %v", family))
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

		chainConfigs = append(chainConfigs, &nodev1.ChainConfig{
			Chain: &nodev1.Chain{
				Id:   strconv.Itoa(int(chainID)),
				Type: ctype,
			},
			AccountAddress: "", // TODO: support AccountAddress
			AdminAddress:   "",
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
