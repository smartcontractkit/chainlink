package changeset

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const (
	jdP2PIDLabel = "p2p_id"

	skipSignerFamilyJobCheckHint = "set skipSignerFamilyJobCheck to bypass"
)

// requestedSignerFamilies returns the sorted, de-duplicated extraSignerFamilies of every OCR3 config
// entry that this change will generate. Entries whose offchainConfig is already processed are not
// regenerated, so they request nothing.
func requestedSignerFamilies(capConfigs []ocr3CapConfig) ([]string, error) {
	seen := map[string]bool{}
	for _, capCfg := range capConfigs {
		for key, entryCfg := range extractOCR3Configs(capCfg.Config) {
			entryMap, ok := entryCfg.(map[string]any)
			if !ok {
				continue
			}
			if _, isMap := entryMap["offchainConfig"].(map[string]any); !isMap {
				continue
			}
			families, err := parseExtraSignerFamilies(entryMap)
			if err != nil {
				return nil, fmt.Errorf("capability %q, ocr3Configs[%q].extraSignerFamilies: %w", capCfg.CapabilityID, key, err)
			}
			for _, f := range families {
				seen[f] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for f := range seen {
		out = append(out, f)
	}
	slices.Sort(out)

	return out, nil
}

// donNodeP2PIDs returns the "p2p_" IDs of a DON's on-chain nodes.
func donNodeP2PIDs(don *capabilities_registry_v2.CapabilitiesRegistryDONInfo) []string {
	ids := make([]string, len(don.NodeP2PIds))
	for i, id := range don.NodeP2PIds {
		ids[i] = p2pkey.PeerID(id).String()
	}

	return ids
}

// consensusSpec is the part of a standard-capability job spec that carries the per-family key bundles.
type consensusSpec struct {
	Type          string `toml:"type"`
	Command       string `toml:"command"`
	OracleFactory struct {
		OnchainSigningStrategy struct {
			Config map[string]string `toml:"config"`
		} `toml:"onchainSigningStrategy"`
	} `toml:"oracle_factory"`
}

// consensusSpecFamilies returns the chain families a consensus job spec has a key bundle for. Other
// job specs, and specs that do not parse, yield none. The consensus job is the standard capability
// whose command is the consensus binary.
func consensusSpecFamilies(spec string) map[string]bool {
	var s consensusSpec
	if err := toml.Unmarshal([]byte(spec), &s); err != nil {
		return nil
	}
	if s.Type != "standardcapabilities" || !strings.HasSuffix(s.Command, "consensus") {
		return nil
	}
	out := map[string]bool{}
	for family, bundleID := range s.OracleFactory.OnchainSigningStrategy.Config {
		if bundleID != "" {
			out[family] = true
		}
	}

	return out
}

// verifySignerFamilyJobs fails unless every node in p2pIDs has an APPROVED consensus job signing for
// every family in families. Adding a signer family to the consensus OCR3 config before the nodes run
// a job holding that family's key bundle halts OCR (INCIDENT-2694). Nodes JD does not know cannot
// have accepted anything, so they fail the check too. It runs at proposal time, not MCMS execution.
func verifySignerFamilyJobs(e cldf.Environment, p2pIDs, families []string) error {
	if len(families) == 0 || len(p2pIDs) == 0 {
		return nil
	}
	if e.Offchain == nil {
		return fmt.Errorf("signer family job check needs a job distributor client; %s", skipSignerFamilyJobCheckHint)
	}

	jdNodes, err := offchain.FetchNodesFromJD(e.GetContext(), e.Offchain, &nodev1.ListNodesRequest_Filter{
		Selectors: []*ptypes.Selector{{Key: jdP2PIDLabel, Op: ptypes.SelectorOp_EXIST}},
	})
	if err != nil {
		return fmt.Errorf("signer family job check: %w", err)
	}

	// A peer can have several JD records (e.g. re-registered with a new CSA key); any may hold the job.
	nodeIDsByP2P := make(map[string][]string, len(jdNodes))
	nodeNameByID := make(map[string]string, len(jdNodes))
	for _, n := range jdNodes {
		for _, l := range n.Labels {
			if l.Key == jdP2PIDLabel && l.Value != nil {
				nodeIDsByP2P[normalizeP2PID(*l.Value)] = append(nodeIDsByP2P[normalizeP2PID(*l.Value)], n.Id)
				nodeNameByID[n.Id] = n.Name
			}
		}
	}

	var allNodeIDs []string
	for _, p := range p2pIDs {
		allNodeIDs = append(allNodeIDs, nodeIDsByP2P[normalizeP2PID(p)]...)
	}
	covered, err := approvedConsensusFamilies(e, allNodeIDs)
	if err != nil {
		return err
	}

	var problems []string
	for _, p := range p2pIDs {
		ids := nodeIDsByP2P[normalizeP2PID(p)]
		if len(ids) == 0 {
			problems = append(problems, p+" (not in JD)")
			continue
		}
		var missing []string
		for _, f := range families {
			if !slices.ContainsFunc(ids, func(id string) bool { return covered[id][f] }) {
				missing = append(missing, f)
			}
		}
		if len(missing) > 0 {
			problems = append(problems, fmt.Sprintf("%s (%s)", nodeNameByID[ids[0]], strings.Join(missing, ", ")))
		}
	}
	if len(problems) == 0 {
		return nil
	}
	slices.Sort(problems)

	return fmt.Errorf(
		"consensus config requests signer families %v, but these nodes have no APPROVED consensus job with a key bundle for them: %s; "+
			"propose the consensus job and have every NOP accept it before changing the capabilities registry (see INCIDENT-2694); %s",
		families, strings.Join(problems, "; "), skipSignerFamilyJobCheckHint,
	)
}

// verifyDONSignerFamilies gates one DON: the families its capability configs request must be signed for
// by an APPROVED consensus job on each node. nodes overrides the DON's current on-chain node set.
func verifyDONSignerFamilies(
	e cldf.Environment, chainSel uint64, qualifier, donName string, capConfigs []contracts.CapabilityConfig, nodes []string,
) error {
	ocr3Configs := make([]ocr3CapConfig, len(capConfigs))
	for i, cc := range capConfigs {
		ocr3Configs[i] = ocr3CapConfig{CapabilityID: cc.Capability.CapabilityID, Config: cc.Config}
	}
	families, err := requestedSignerFamilies(ocr3Configs)
	if err != nil || len(families) == 0 {
		return err
	}

	if len(nodes) == 0 {
		capReg, err := capRegForPreconditionCheck(e, chainSel, qualifier)
		if err != nil {
			return err
		}
		// By name, not the first-32 DON list, so DONs beyond that window still resolve.
		don, _, err := sequences.GetDonNodes(donName, capReg)
		if err != nil {
			return fmt.Errorf("DON %q not found on-chain: %w", donName, err)
		}
		nodes = donNodeP2PIDs(don)
	}

	return verifySignerFamilyJobs(e, nodes, families)
}

// normalizeP2PID maps a bare-hex or "p2p_"-prefixed peer ID to the "p2p_" form JD labels use.
func normalizeP2PID(p string) string {
	if id, err := p2pkey.MakePeerID(p); err == nil {
		return id.String()
	}

	return p
}

// approvedConsensusFamilies maps each JD node ID to the families its live consensus jobs sign for,
// counting only a job's latest revision and only when that revision is APPROVED.
func approvedConsensusFamilies(e cldf.Environment, nodeIDs []string) (map[string]map[string]bool, error) {
	covered := map[string]map[string]bool{}
	if len(nodeIDs) == 0 {
		return covered, nil
	}

	const maxRecvSize = 50 * 1024 * 1024
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{NodeIds: nodeIDs},
	}, grpc.MaxCallRecvMsgSize(maxRecvSize))
	if err != nil {
		return nil, fmt.Errorf("signer family job check: failed to list jobs: %w", err)
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
		return covered, nil
	}

	propResp, err := e.Offchain.ListProposals(e.GetContext(), &jobv1.ListProposalsRequest{
		Filter: &jobv1.ListProposalsRequest_Filter{JobIds: jobIDs},
	}, grpc.MaxCallRecvMsgSize(maxRecvSize))
	if err != nil {
		return nil, fmt.Errorf("signer family job check: failed to list proposals: %w", err)
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

	for jobID, p := range latest {
		if p.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
			continue
		}
		node := nodeByJob[jobID]
		for f := range consensusSpecFamilies(p.Spec) {
			if covered[node] == nil {
				covered[node] = map[string]bool{}
			}
			covered[node][f] = true
		}
	}

	return covered, nil
}
