package clo

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/environment/clo/models"
)

type JobClient struct {
	NodeOperators []*models.NodeOperator `json:"nodeOperators"`
	nodesByID     map[string]*models.Node
	lggr          logger.Logger
}

func (j JobClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	//TODO implement me
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
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
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
	for _, nop := range j.NodeOperators {
		for _, n := range nop.Nodes {
			p2pId, err := NodeP2PId(n)
			if err != nil {
				return nil, fmt.Errorf("failed to get p2p id for node %s: %w", n.ID, err)
			}
			node := &nodev1.Node{
				Id:          n.ID,
				Name:        n.Name,
				PublicKey:   *n.PublicKey,
				IsEnabled:   n.Enabled,
				IsConnected: n.Connected,
				Labels: []*ptypes.Label{
					{
						Key:   "p2p_id",
						Value: &p2pId, // here n.ID is also peer ID
					},
				},
			}
			if include(node) {
				nodes = append(nodes, node)
			}
		}
	}
	return &nodev1.ListNodesResponse{
		Nodes: nodes,
	}, nil
}

func (j JobClient) ListNodeChainConfigs(ctx context.Context, in *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {

	resp := &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: make([]*nodev1.ChainConfig, 0),
	}
	// no filter, return all
	if in.Filter == nil || len(in.Filter.NodeIds) == 0 {
		for _, n := range j.nodesByID {
			ccfg := cloNodeToChainConfigs(n)
			resp.ChainConfigs = append(resp.ChainConfigs, ccfg...)
		}
	} else {
		for _, want := range in.Filter.NodeIds {
			n, ok := j.nodesByID[want]
			if !ok {
				j.lggr.Warn("node not found", zap.String("node_id", want))
				continue
			}
			ccfg := cloNodeToChainConfigs(n)
			resp.ChainConfigs = append(resp.ChainConfigs, ccfg...)
		}
	}
	return resp, nil

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
	panic("implement me")

}

func (j JobClient) RevokeJob(ctx context.Context, in *jobv1.RevokeJobRequest, opts ...grpc.CallOption) (*jobv1.RevokeJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

func (j JobClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	//TODO CCIP-3108 implement me
	panic("implement me")
}

type GetNodeOperatorsResponse struct {
	NodeOperators []*models.NodeOperator `json:"nodeOperators"`
}

type JobClientConfig struct {
	Nops []*models.NodeOperator
}

func NewJobClient(lggr logger.Logger, cfg JobClientConfig) *JobClient {

	c := &JobClient{
		NodeOperators: cfg.Nops,
		nodesByID:     make(map[string]*models.Node),
		lggr:          lggr,
	}
	for _, nop := range c.NodeOperators {
		for _, n := range nop.Nodes {
			node := n
			c.nodesByID[n.ID] = node // maybe should use the public key instead?
		}
	}
	return c
}

func cloNodeToChainConfigs(n *models.Node) []*nodev1.ChainConfig {
	out := make([]*nodev1.ChainConfig, 0)
	for _, ccfg := range n.ChainConfigs {
		out = append(out, cloChainCfgToJDChainCfg(ccfg))
	}
	return out
}

func cloChainCfgToJDChainCfg(ccfg *models.NodeChainConfig) *nodev1.ChainConfig {
	var ctype nodev1.ChainType
	switch ccfg.Network.ChainType {
	case models.ChainTypeEvm:
		ctype = nodev1.ChainType_CHAIN_TYPE_EVM
	case models.ChainTypeSolana:
		ctype = nodev1.ChainType_CHAIN_TYPE_SOLANA
	case models.ChainTypeStarknet:
		ctype = nodev1.ChainType_CHAIN_TYPE_STARKNET
	case models.ChainTypeAptos:
		ctype = nodev1.ChainType_CHAIN_TYPE_APTOS
	default:
		panic(fmt.Sprintf("Unsupported chain family %v", ccfg.Network.ChainType))
	}

	return &nodev1.ChainConfig{
		Chain: &nodev1.Chain{
			Id:   ccfg.Network.ChainID,
			Type: ctype,
		},
		AccountAddress: ccfg.AccountAddress,
		AdminAddress:   ccfg.AdminAddress,
		// only care about ocr2 for now
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     ccfg.Ocr2Config.Enabled,
			IsBootstrap: ccfg.Ocr2Config.IsBootstrap,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    ccfg.Ocr2Config.P2pKeyBundle.PeerID,
				PublicKey: ccfg.Ocr2Config.P2pKeyBundle.PublicKey,
			},
			OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
				BundleId:              ccfg.Ocr2Config.OcrKeyBundle.BundleID,
				ConfigPublicKey:       ccfg.Ocr2Config.OcrKeyBundle.ConfigPublicKey,
				OffchainPublicKey:     ccfg.Ocr2Config.OcrKeyBundle.OffchainPublicKey,
				OnchainSigningAddress: ccfg.Ocr2Config.OcrKeyBundle.OnchainSigningAddress,
			},
			// TODO: the clo cli does not serialize this field, so it will always be nil
			//Multiaddr:        *ccfg.Ocr2Config.Multiaddr,
			//ForwarderAddress: ccfg.Ocr2Config.ForwarderAddress,
		},
	}
}
