package view

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment"
)

type NopViewV2 struct {
	Nodes []NopNodeInfoV2 `json:"nodes"`
}

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
	Version          string                `json:"version"`
	Labels           []LabelView           `json:"labels,omitempty"`
	ApprovedJobspecs map[string]JobView    `json:"approvedJobspecs,omitempty"` // jobID => jobSpec
	ProposedJobspecs map[string]JobView    `json:"proposedJobspecs,omitempty"` // jobID => jobSpec
}

type NopNodeInfoV2 struct {
	NopView
	NodeName   string   `json:"nodeName"`
	Networks   []string `json:"networks"`
	Deployment string   `json:"deployment"`
}

type JobView struct {
	ProposalID string `json:"proposal_id"`
	UUID       string `json:"uuid"`
	Spec       string `json:"spec"`
	Revision   int64  `json:"revision,omitempty"`
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
func GenerateNopsView(lggr logger.Logger, nodeIDs []string, oc cldf_offchain.Client) (map[string]NopView, error) {
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
	jobspecs, proposedSpecs, err := ApprovedJobspecs(context.Background(), lggr, nodeIDs, oc)
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

		labels := make([]LabelView, 0, len(nodeDetails.Labels))
		for _, l := range nodeDetails.Labels {
			labels = append(labels, LabelView{
				Key:   l.Key,
				Value: l.Value,
			})
		}

		fullNodeInfo := NopView{
			NodeID:           node.NodeID,
			PeerID:           node.PeerID.String(),
			IsBootstrap:      node.IsBootstrap,
			OCRKeys:          make(map[string]OCRKeyView),
			PayeeAddress:     node.AdminAddr,
			CSAKey:           nodeDetails.PublicKey,
			WorkflowKey:      nodeDetails.GetWorkflowKey(),
			IsConnected:      nodeDetails.IsConnected,
			IsEnabled:        nodeDetails.IsEnabled,
			Version:          nodeDetails.Version,
			Labels:           labels,
			ApprovedJobspecs: jobspecs[node.NodeID],
			ProposedJobspecs: proposedSpecs[node.NodeID],
		}
		for details, ocrConfig := range node.SelToOCRConfig {
			fullNodeInfo.OCRKeys[details.ChainName] = OCRKeyView{
				OffchainPublicKey:         hex.EncodeToString(ocrConfig.OffchainPublicKey[:]),
				OnchainPublicKey:          fmt.Sprintf("%x", ocrConfig.OnchainPublicKey[:]),
				PeerID:                    ocrConfig.PeerID.String(),
				TransmitAccount:           string(ocrConfig.TransmitAccount),
				ConfigEncryptionPublicKey: hex.EncodeToString(ocrConfig.ConfigEncryptionPublicKey[:]),
				KeyBundleID:               ocrConfig.KeyBundleID,
			}
		}
		nv[nodeName] = fullNodeInfo
	}

	return nv, nil
}

type NopNameRemapper func(nodeName string) string

// GenerateNOPsViewV2 generates a view of nodes with their details in a new format.
// `deploymentKey` refers to the deployment identifier (e.g., "keystone", "cre", "ccip", "data-feeds"), which usually refers to the CLD domain.
func GenerateNOPsViewV2(ctx context.Context, lggr logger.Logger, nodeIDs []string, oc cldf_offchain.Client, deploymentKey string, nopNameRemapperFunc NopNameRemapper) (map[string]NopViewV2, error) {
	nodes, err := deployment.NodeInfo(nodeIDs, oc)
	if errors.Is(err, deployment.ErrMissingNodeMetadata) {
		lggr.Warnf("Missing node metadata: %s", err.Error())
	} else if err != nil {
		return nil, fmt.Errorf("failed to get node info: %w", err)
	}
	nodesResp, err := oc.ListNodes(ctx, &nodev1.ListNodesRequest{
		Filter: &nodev1.ListNodesRequest_Filter{
			Ids: nodeIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes from JD: %w", err)
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
	jobspecs, proposedSpecs, err := ApprovedJobspecs(ctx, lggr, nodeIDs, oc)
	if err != nil {
		// best effort on job specs
		lggr.Warnf("Failed to get approved jobspecs: %v", err)
	}

	groupedNops := make(map[string]NopViewV2)
	for _, node := range nodes {
		nodeName := node.Name
		if nodeName == "" {
			nodeName = node.NodeID
		}

		var nopName string
		if nopNameRemapperFunc != nil {
			nopName = nopNameRemapperFunc(nodeName)
		} else {
			nopName = defaultNopNameRemapper(nodeName, deploymentKey)
		}

		nodeDetails := details(node.NodeID)
		if nodeDetails == nil {
			return groupedNops, fmt.Errorf("failed to get node details for node %s", node.NodeID)
		}
		var labels []LabelView
		for _, l := range nodeDetails.Labels {
			labels = append(labels, LabelView{
				Key:   l.Key,
				Value: l.Value,
			})
		}

		networks, networksErr := getNodeNetworks(node)
		if networksErr != nil {
			// best effort on networks
			lggr.Warnf("Failed to get networks: %v", networksErr)
		}

		fullNodeInfo := NopNodeInfoV2{
			NopView: NopView{
				NodeID:           node.NodeID,
				PeerID:           node.PeerID.String(),
				IsBootstrap:      node.IsBootstrap,
				OCRKeys:          make(map[string]OCRKeyView),
				PayeeAddress:     node.AdminAddr,
				CSAKey:           nodeDetails.PublicKey,
				WorkflowKey:      nodeDetails.GetWorkflowKey(),
				IsConnected:      nodeDetails.IsConnected,
				IsEnabled:        nodeDetails.IsEnabled,
				Version:          nodeDetails.Version,
				Labels:           labels,
				ApprovedJobspecs: jobspecs[node.NodeID],
				ProposedJobspecs: proposedSpecs[node.NodeID],
			},
			NodeName:   nodeName,
			Deployment: deploymentKey,
			Networks:   networks,
		}
		for details, ocrConfig := range node.SelToOCRConfig {
			fullNodeInfo.OCRKeys[details.ChainName] = OCRKeyView{
				OffchainPublicKey:         hex.EncodeToString(ocrConfig.OffchainPublicKey[:]),
				OnchainPublicKey:          fmt.Sprintf("%x", ocrConfig.OnchainPublicKey[:]),
				PeerID:                    ocrConfig.PeerID.String(),
				TransmitAccount:           string(ocrConfig.TransmitAccount),
				ConfigEncryptionPublicKey: hex.EncodeToString(ocrConfig.ConfigEncryptionPublicKey[:]),
				KeyBundleID:               ocrConfig.KeyBundleID,
			}
		}

		var nop NopViewV2
		var ok bool
		if nop, ok = groupedNops[nopName]; !ok {
			nop = NopViewV2{
				Nodes: make([]NopNodeInfoV2, 0),
			}
		}

		nop.Nodes = append(nop.Nodes, fullNodeInfo)
		groupedNops[nopName] = nop
	}

	// Sort the nodes within each group by NodeID for deterministic output.
	for key := range groupedNops {
		sort.Slice(groupedNops[key].Nodes, func(i, j int) bool {
			return groupedNops[key].Nodes[i].NodeID < groupedNops[key].Nodes[j].NodeID
		})
	}

	return groupedNops, nil
}

// defaultNopNameRemapper groups by node name, by extracting the first word before any hyphen,
// unless the word is cll-<deployment>, cl-<deployment> or clp-<deployment>, if so, treat that as the first word.
// It assumes that <deployment> could be one of the following: cre, keystone, ccip, data-feeds, etc. (CLD domain)
func defaultNopNameRemapper(nodeName, deploymentKey string) string {
	var nopName string

	switch {
	case strings.HasPrefix(nodeName, "cl-"+deploymentKey):
		nopName = "cll-" + deploymentKey
	case strings.HasPrefix(nodeName, "cll-"+deploymentKey):
		nopName = "cll-" + deploymentKey
	case strings.HasPrefix(nodeName, "clp-"+deploymentKey):
		nopName = "clp-" + deploymentKey
	default:
		parts := strings.Split(nodeName, "-")
		nopName = parts[0]
	}

	return nopName
}

func ApprovedJobspecs(ctx context.Context, lggr logger.Logger, nodeIDs []string, oc cldf_offchain.Client) (nodeJobsView map[string]map[string]JobView, proposedJobsView map[string]map[string]JobView, verr error) {
	nodeJobsView = make(map[string]map[string]JobView)
	proposedJobsView = make(map[string]map[string]JobView)

	jobs, err := oc.ListJobs(ctx, &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			NodeIds: nodeIDs,
		},
	})
	if err != nil {
		return nodeJobsView, proposedJobsView, fmt.Errorf("failed to list jobs for nodes %v: %w", nodeIDs, err)
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
		proposed := make(map[string]JobView)
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
			if p.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED {
				if _, exists := jv[p.JobId]; exists && p.Revision < jv[p.JobId].Revision {
					// skip older revisions
					continue
				}
				proposed[p.JobId] = JobView{
					ProposalID: p.Id,
					UUID:       jobs[p.JobId].Uuid,
					Spec:       p.Spec,
					Revision:   p.Revision,
				}
			}
			if p.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
				if _, exists := jv[p.JobId]; exists && p.Revision < jv[p.JobId].Revision {
					// skip older revisions
					continue
				}
				jv[p.JobId] = JobView{
					ProposalID: p.Id,
					UUID:       jobs[p.JobId].Uuid,
					Spec:       p.Spec,
					Revision:   p.Revision,
				}
			}
		}
		nodeJobsView[nodeID] = jv
		proposedJobsView[nodeID] = proposed
	}
	return nodeJobsView, proposedJobsView, verr
}

// getNodeNetworks returns the list of networks a node is connected to.
// This function mimics the logic of the CLD command `jd node inspect`
func getNodeNetworks(node deployment.Node) ([]string, error) {
	nodeChainCfgs, nodeErr := node.ChainConfigs()
	if nodeErr != nil {
		return nil, fmt.Errorf("failed to get chain configs for node %s: %w", node.NodeID, nodeErr)
	}

	var networks []string
	for _, cfg := range nodeChainCfgs {
		var family string
		switch cfg.Chain.Type {
		case nodev1.ChainType_CHAIN_TYPE_EVM:
			family = chain_selectors.FamilyEVM
		case nodev1.ChainType_CHAIN_TYPE_APTOS:
			family = chain_selectors.FamilyAptos
		case nodev1.ChainType_CHAIN_TYPE_SOLANA:
			family = chain_selectors.FamilySolana
		case nodev1.ChainType_CHAIN_TYPE_STARKNET:
			family = chain_selectors.FamilyStarknet
		case nodev1.ChainType_CHAIN_TYPE_TRON:
			family = chain_selectors.FamilyTron
		case nodev1.ChainType_CHAIN_TYPE_TON:
			family = chain_selectors.FamilyTon
		case nodev1.ChainType_CHAIN_TYPE_SUI:
			family = chain_selectors.FamilySui
		case nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED:
			return nil, errors.New("chain type must be specified")
		default:
			return nil, fmt.Errorf("unsupported chain type %s", cfg.Chain.Type)
		}

		chainDetails, chainErr := chain_selectors.GetChainDetailsByChainIDAndFamily(cfg.Chain.Id, family)
		if chainErr != nil {
			return nil, fmt.Errorf("failed to get chain details for chain ID %s and family %s: %w", cfg.Chain.Id, family, chainErr)
		}

		networks = append(networks, chainDetails.ChainName)
	}

	return networks, nil
}
