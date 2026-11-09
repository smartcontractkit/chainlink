package changeset

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const jdP2PIDLabel = "p2p_id"

// verifyNoPendingJobProposals fails if any node identified by p2pIDs has a job whose latest proposal
// in JD is undecided or rejected. It only sees proposals that already exist: a cap-reg change applied
// before the job is proposed at all is not caught. Preconditions also run at proposal time, not MCMS
// execution time. Applying cap-reg config before NOPs accept updated job specs (e.g. new
// OCR key bundles) breaks OCR on those nodes (INCIDENT-2694).
func verifyNoPendingJobProposals(e cldf.Environment, p2pIDs []string) error {
	if len(p2pIDs) == 0 {
		return nil
	}
	if e.Offchain == nil {
		return errors.New("job acceptance check needs a job distributor client; set skipJobAcceptanceCheck to bypass")
	}

	jdNodes, err := offchain.FetchNodesFromJD(e.GetContext(), e.Offchain, &nodev1.ListNodesRequest_Filter{
		Selectors: []*ptypes.Selector{{Key: jdP2PIDLabel, Op: ptypes.SelectorOp_EXIST}},
	})
	if err != nil {
		return fmt.Errorf("job acceptance check: %w", err)
	}

	nodeIDByP2P := make(map[string]string, len(jdNodes))
	nodeNameByID := make(map[string]string, len(jdNodes))
	for _, n := range jdNodes {
		for _, l := range n.Labels {
			if l.Key == jdP2PIDLabel && l.Value != nil {
				nodeIDByP2P[*l.Value] = n.Id
				nodeNameByID[n.Id] = n.Name
			}
		}
	}

	// Nodes JD does not know cannot have proposals there, so there is nothing to wait on.
	var unmatched []string
	nodeIDs := make([]string, 0, len(p2pIDs))
	for _, p := range p2pIDs {
		// Normalize so a bare-hex ID in the input still matches JD's "p2p_" label instead of being skipped.
		key := p
		if id, err := p2pkey.MakePeerID(p); err == nil {
			key = id.String()
		}
		if id, ok := nodeIDByP2P[key]; ok {
			nodeIDs = append(nodeIDs, id)
		} else {
			unmatched = append(unmatched, p)
		}
	}
	if len(unmatched) > 0 && e.Logger != nil {
		e.Logger.Warnf("job acceptance check: no JD node with a %s label for %s; their jobs were not checked", jdP2PIDLabel, strings.Join(unmatched, ", "))
	}
	if len(nodeIDs) == 0 {
		return nil
	}

	pending, err := pendingProposalNodes(e, nodeIDs)
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return nil
	}
	names := make([]string, len(pending))
	for i, id := range pending {
		names[i] = nodeNameByID[id]
	}

	return fmt.Errorf(
		"node(s) %s have job proposals that are pending or rejected; all NOPs must accept updated jobs before the capabilities registry change is applied (see INCIDENT-2694); set skipJobAcceptanceCheck to bypass",
		strings.Join(names, ", "),
	)
}

// donNodeP2PIDs returns the "p2p_" IDs of the on-chain nodes of the named DON, and false if the DON
// does not exist.
func donNodeP2PIDs(dons []capabilities_registry_v2.CapabilitiesRegistryDONInfo, name string) ([]string, bool) {
	for _, don := range dons {
		if don.Name != name {
			continue
		}
		ids := make([]string, len(don.NodeP2PIds))
		for i, id := range don.NodeP2PIds {
			ids[i] = p2pkey.PeerID(id).String()
		}

		return ids, true
	}

	return nil, false
}

// blocksCapRegChange reports whether a latest-revision status means the node is not running the
// intended spec: undecided, or refused by the NOP. CANCELLED and REVOKED are deliberate retractions.
func blocksCapRegChange(st jobv1.ProposalStatus) bool {
	switch st {
	case jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED,
		jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING,
		jobv1.ProposalStatus_PROPOSAL_STATUS_REJECTED:
		return true
	default:
		return false
	}
}

// pendingProposalNodes returns the sorted IDs of nodes whose latest proposal for some live job is
// undecided or rejected.
func pendingProposalNodes(e cldf.Environment, nodeIDs []string) ([]string, error) {
	const maxRecvSize = 50 * 1024 * 1024
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{NodeIds: nodeIDs},
	}, grpc.MaxCallRecvMsgSize(maxRecvSize))
	if err != nil {
		return nil, fmt.Errorf("job acceptance check: failed to list jobs: %w", err)
	}

	nodeByJob := make(map[string]string)
	jobIDs := make([]string, 0, len(jobsResp.Jobs))
	for _, j := range jobsResp.Jobs {
		if j.DeletedAt != nil {
			continue
		}
		nodeByJob[j.Id] = j.NodeId
		jobIDs = append(jobIDs, j.Id)
	}
	if len(jobIDs) == 0 {
		return nil, nil
	}

	propResp, err := e.Offchain.ListProposals(e.GetContext(), &jobv1.ListProposalsRequest{
		Filter: &jobv1.ListProposalsRequest_Filter{JobIds: jobIDs},
	}, grpc.MaxCallRecvMsgSize(maxRecvSize))
	if err != nil {
		return nil, fmt.Errorf("job acceptance check: failed to list proposals: %w", err)
	}

	latest := make(map[string]*jobv1.Proposal)
	for _, p := range propResp.Proposals {
		if _, live := nodeByJob[p.JobId]; !live {
			continue
		}
		if cur, ok := latest[p.JobId]; !ok || p.Revision > cur.Revision {
			latest[p.JobId] = p
		}
	}

	var pending []string
	for jobID, p := range latest {
		if !blocksCapRegChange(p.Status) {
			continue
		}
		if node := nodeByJob[jobID]; !slices.Contains(pending, node) {
			pending = append(pending, node)
		}
	}
	slices.Sort(pending)

	return pending, nil
}
