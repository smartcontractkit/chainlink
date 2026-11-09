package changeset_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

type pendingJD struct {
	cldf_offchain.Client
	p2pID string
}

func (f pendingJD) ListNodes(context.Context, *nodev1.ListNodesRequest, ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	return &nodev1.ListNodesResponse{Nodes: []*nodev1.Node{
		{Id: "n1", Name: "node-1", Labels: []*ptypes.Label{{Key: "p2p_id", Value: &f.p2pID}}},
	}}, nil
}

func (f pendingJD) ListJobs(context.Context, *jobv1.ListJobsRequest, ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	return &jobv1.ListJobsResponse{Jobs: []*jobv1.Job{{Id: "j1", NodeId: "n1"}}}, nil
}

func (f pendingJD) ListProposals(context.Context, *jobv1.ListProposalsRequest, ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	return &jobv1.ListProposalsResponse{Proposals: []*jobv1.Proposal{
		{JobId: "j1", Revision: 1, Status: jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED},
	}}, nil
}

func TestAddCapabilities_VerifyPreconditions_JobAcceptanceGate(t *testing.T) {
	h := test.NewTestHarness(t)

	e := h.Runtime.Environment()
	e.Offchain = pendingJD{p2pID: h.Don.GetP2PIDs().Strings()[0]}

	in := changeset.AddCapabilitiesInput{
		RegistryChainSel:  h.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			test.DONName: {{Capability: contracts.Capability{CapabilityID: "new-cap@1.0.0"}, Config: map[string]any{"k": "v"}}},
		},
		SkipDuplicateCapabilityValidation: true,
	}

	err := changeset.AddCapabilities{}.VerifyPreconditions(e, in)
	require.ErrorContains(t, err, "node-1")

	in.SkipJobAcceptanceCheck = true
	require.NoError(t, changeset.AddCapabilities{}.VerifyPreconditions(e, in))
}
