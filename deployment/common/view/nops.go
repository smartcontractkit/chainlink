package view

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
)

type NopView struct {
	// NodeID is the unique identifier of the node
	NodeID           string                `json:"nodeID"`
	PeerID           string                `json:"peerID"`
	IsBootstrap      bool                  `json:"isBootstrap"`
	OCRKeys          map[string]OCRKeyView `json:"ocrKeys"`
	PayeeAddress     string                `json:"payeeAddress"`
	CSAKey           string                `json:"csaKey"`
	WorkflowKey      string                `json:"workflowKey,omitempty"`
	IsConnected      bool                  `json:"isConnected"`
	IsEnabled        bool                  `json:"isEnabled"`
	Labels           []LabelView           `json:"labels,omitempty"`
	ApprovedJobspecs map[string]JobView    `json:"approvedJobspecs,omitempty"` // jobID => jobSpec
}

type JobView struct {
	ProposalID string `json:"proposal_id"`
	UUID       string `json:"uuid"`
	Spec       string `json:"spec"`
}

type LabelView struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

type OCRKeyView struct {
	OffchainPublicKey         string `json:"offchainPublicKey"`
	OnchainPublicKey          string `json:"onchainPublicKey"`
	PeerID                    string `json:"peerID"`
	TransmitAccount           string `json:"transmitAccount"`
	ConfigEncryptionPublicKey string `json:"configEncryptionPublicKey"`
	KeyBundleID               string `json:"keyBundleID"`
}

// GenerateNopsView generates a view of nodes with their details
func GenerateNopsView(lggr logger.Logger, nodeIDs []string, oc cldf.OffchainClient) (map[string]NopView, error) {
	nv := make(map[string]NopView)
	nodes, err := deployment.NodeInfo(nodeIDs, oc)
	if errors.Is(err, deployment.ErrMissingNodeMetadata) {
		lggr.Warnf("Missing node metadata: %s", err.Error())
	} else if err != nil {
		return nv, fmt.Errorf("failed to get node info: %w", err)
	}
	nodesResp, err := oc.ListNodes(context.Background(), &nodev1.ListNodesRequest{
		Filter: &nodev1.ListNodesRequest_Filter{
			Ids: nodeIDs,
		},
	})
	if err != nil {
		return nv, fmt.Errorf("failed to list nodes from JD: %w", err)
	}
	details := func(nodeID string) *nodev1.Node {
		// extract from the response
		for _, node := range nodesResp.Nodes {
			if node.Id == nodeID {
				return node
			}
		}
		return nil
	}
	jobspecs, err := approvedJobspecs(context.Background(), lggr, nodeIDs, oc)
	if err != nil {
		// best effort on job specs
		lggr.Warnf("Failed to get approved jobspecs: %v", err)
	}

	for _, node := range nodes {
		nodeName := node.Name
		if nodeName == "" {
			nodeName = node.NodeID
		}
		nodeDetails := details(node.NodeID)
		if nodeDetails == nil {
			return nv, fmt.Errorf("failed to get node details for node %s", node.NodeID)
		}
		labels := []LabelView{}
		for _, l := range nodeDetails.Labels {
			labels = append(labels, LabelView{
				Key:   l.Key,
				Value: l.Value,
			})
		}
		nop := NopView{
			NodeID:           node.NodeID,
			PeerID:           node.PeerID.String(),
			IsBootstrap:      node.IsBootstrap,
			OCRKeys:          make(map[string]OCRKeyView),
			PayeeAddress:     node.AdminAddr,
			CSAKey:           nodeDetails.PublicKey,
			WorkflowKey:      nodeDetails.GetWorkflowKey(),
			IsConnected:      nodeDetails.IsConnected,
			IsEnabled:        nodeDetails.IsEnabled,
			Labels:           labels,
			ApprovedJobspecs: jobspecs[node.NodeID],
		}
		for details, ocrConfig := range node.SelToOCRConfig {
			nop.OCRKeys[details.ChainName] = OCRKeyView{
				OffchainPublicKey:         hex.EncodeToString(ocrConfig.OffchainPublicKey[:]),
				OnchainPublicKey:          fmt.Sprintf("%x", ocrConfig.OnchainPublicKey[:]),
				PeerID:                    ocrConfig.PeerID.String(),
				TransmitAccount:           string(ocrConfig.TransmitAccount),
				ConfigEncryptionPublicKey: hex.EncodeToString(ocrConfig.ConfigEncryptionPublicKey[:]),
				KeyBundleID:               ocrConfig.KeyBundleID,
			}
		}
		nv[nodeName] = nop
	}
	return nv, nil
}

func approvedJobspecs(ctx context.Context, lggr logger.Logger, nodeIDs []string, oc cldf.OffchainClient) (nodeJobsView map[string]map[string]JobView, verr error) {
	nodeJobsView = make(map[string]map[string]JobView)

	jobs, err := oc.ListJobs(ctx, &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			NodeIds: nodeIDs,
		},
	})
	if err != nil {
		return nodeJobsView, fmt.Errorf("failed to list jobs for nodes %v: %w", nodeIDs, err)
	}
	nodeJobIDs := make(map[string]map[string]*jobv1.Job) // node id -> job id -> job
	for i, j := range jobs.Jobs {
		// skip deleted jobs
		if j.DeletedAt != nil {
			continue
		}
		if _, ok := nodeJobIDs[j.NodeId]; !ok {
			nodeJobIDs[j.NodeId] = make(map[string]*jobv1.Job)
		}
		nodeJobIDs[j.NodeId][j.Id] = jobs.Jobs[i]
	}

	// list proposals for each node
	for nodeID, jobs := range nodeJobIDs {
		jv := make(map[string]JobView) // job id -> view
		lresp, err := oc.ListProposals(ctx, &jobv1.ListProposalsRequest{
			Filter: &jobv1.ListProposalsRequest_Filter{
				JobIds: slices.Collect(maps.Keys(jobs)),
			},
		})
		if err != nil {
			// don't block on single node error
			lggr.Warnf("failed to list job proposals on node %s: %v", nodeID, err)
			verr = errors.Join(verr, fmt.Errorf("failed to list job proposals on node %s: %w", nodeID, err))
			continue
		}
		for _, p := range lresp.Proposals {
			if p.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
				jv[p.JobId] = JobView{
					ProposalID: p.Id,
					UUID:       jobs[p.JobId].Uuid,
					Spec:       p.Spec,
				}
			}
		}
		nodeJobsView[nodeID] = jv
	}
	return nodeJobsView, verr
}
