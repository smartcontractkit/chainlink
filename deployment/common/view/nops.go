package view

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

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

func GenerateNopsView(nodeIDs []string, oc deployment.OffchainClient) (map[string]NopView, error) {
	nv := make(map[string]NopView)
	nodes, err := deployment.NodeInfo(nodeIDs, oc)
	if errors.Is(err, deployment.ErrMissingNodeMetadata) {
		fmt.Printf("WARNING: Missing node metadata:\n%s", err.Error())
	} else if err != nil {
		return nv, err
	}
	for _, node := range nodes {
		// get node info
		nodeDetails, err := oc.GetNode(context.Background(), &nodev1.GetNodeRequest{Id: node.NodeID})
		if err != nil {
			return nv, errors.Wrapf(err, "failed to get node details from offchain client for node %s", node.NodeID)
		}
		if nodeDetails == nil || nodeDetails.Node == nil {
			return nv, fmt.Errorf("failed to get node details from offchain client for node %s", node.NodeID)
		}
		nodeName := nodeDetails.Node.Name
		if nodeName == "" {
			nodeName = node.NodeID
		}
		labels := []LabelView{}
		for _, l := range nodeDetails.Node.Labels {
			labels = append(labels, LabelView{
				Key:   l.Key,
				Value: l.Value,
			})
		}
		jobspecs, err := approvedJobspecs(context.Background(), node.NodeID, oc)
		if err != nil {
			return nv, errors.Wrapf(err, "failed to get approved jobspecs for node %s", node.NodeID)
		}
		nop := NopView{
			NodeID:           node.NodeID,
			PeerID:           node.PeerID.String(),
			IsBootstrap:      node.IsBootstrap,
			OCRKeys:          make(map[string]OCRKeyView),
			PayeeAddress:     node.AdminAddr,
			CSAKey:           nodeDetails.Node.PublicKey,
			WorkflowKey:      nodeDetails.Node.GetWorkflowKey(),
			IsConnected:      nodeDetails.Node.IsConnected,
			IsEnabled:        nodeDetails.Node.IsEnabled,
			Labels:           labels,
			ApprovedJobspecs: jobspecs,
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

// acceptedOrPendingAcceptedJobSpecs returns a map of nodeID to job specs that are either accepted or pending review
// or proposed
func approvedJobspecs(ctx context.Context, nodeID string, oc deployment.OffchainClient) (map[string]JobView, error) {
	existingSpecs := make(map[string]JobView)
	jobs, err := oc.ListJobs(ctx, &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			NodeIds: []string{nodeID},
		},
	})
	if err != nil {
		return existingSpecs, fmt.Errorf("failed to list jobs for node %s: %w", nodeID, err)
	}
	for _, j := range jobs.Jobs {
		// skip deleted jobs
		if j.DeletedAt != nil {
			continue
		}
		for _, propID := range j.ProposalIds {
			jbProposal, err := oc.GetProposal(ctx, &jobv1.GetProposalRequest{
				Id: propID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get job proposal %s on node %s: %w", propID, nodeID, err)
			}
			if jbProposal.Proposal == nil {
				return nil, fmt.Errorf("job proposal %s on node %s is nil", propID, nodeID)
			}
			if jbProposal.Proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
				existingSpecs[jbProposal.Proposal.JobId] = JobView{
					ProposalID: jbProposal.Proposal.JobId,
					UUID:       j.Uuid,
					Spec:       jbProposal.Proposal.Spec,
				}
			}
		}
	}
	return existingSpecs, nil
}
