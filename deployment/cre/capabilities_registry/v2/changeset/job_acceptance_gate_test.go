package changeset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

type fakeJD struct {
	cldf_offchain.Client
	nodes     []*nodev1.Node
	jobs      []*jobv1.Job
	proposals []*jobv1.Proposal
	jobsErr   error
	propsErr  error
}

func (f *fakeJD) ListNodes(context.Context, *nodev1.ListNodesRequest, ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	return &nodev1.ListNodesResponse{Nodes: f.nodes}, nil
}

func (f *fakeJD) ListJobs(context.Context, *jobv1.ListJobsRequest, ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	if f.jobsErr != nil {
		return nil, f.jobsErr
	}

	return &jobv1.ListJobsResponse{Jobs: f.jobs}, nil
}

func (f *fakeJD) ListProposals(context.Context, *jobv1.ListProposalsRequest, ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	if f.propsErr != nil {
		return nil, f.propsErr
	}

	return &jobv1.ListProposalsResponse{Proposals: f.proposals}, nil
}

func jdNode(id, name, p2p string) *nodev1.Node {
	return &nodev1.Node{Id: id, Name: name, Labels: []*ptypes.Label{{Key: jdP2PIDLabel, Value: &p2p}}}
}

func envWith(jd cldf_offchain.Client) cldf.Environment {
	return cldf.Environment{Offchain: jd, GetContext: context.Background}
}

func TestVerifyNoPendingJobProposals(t *testing.T) {
	t.Parallel()

	base := func() *fakeJD {
		return &fakeJD{
			nodes: []*nodev1.Node{jdNode("n1", "node-1", "p2p_a"), jdNode("n2", "node-2", "p2p_b")},
			jobs:  []*jobv1.Job{{Id: "j1", NodeId: "n1"}, {Id: "j2", NodeId: "n2"}},
		}
	}
	prop := func(job string, rev int64, st jobv1.ProposalStatus) *jobv1.Proposal {
		return &jobv1.Proposal{Id: job + "-p", JobId: job, Revision: rev, Status: st}
	}

	t.Run("no p2p ids is a no-op without JD", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, verifyNoPendingJobProposals(cldf.Environment{}, nil))
	})

	t.Run("missing JD client errors", func(t *testing.T) {
		t.Parallel()
		err := verifyNoPendingJobProposals(cldf.Environment{}, []string{"p2p_a"})
		require.ErrorContains(t, err, "job distributor client")
	})

	t.Run("all approved passes", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
			prop("j2", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
		}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a", "p2p_b"}))
	})

	t.Run("proposed revision on one node blocks and names it", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
			prop("j1", 2, jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED),
			prop("j2", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
		}
		err := verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a", "p2p_b"})
		require.ErrorContains(t, err, "node-1")
		require.NotContains(t, err.Error(), "node-2")
	})

	t.Run("superseded older pending revision is ignored", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING),
			prop("j1", 2, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
		}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a"}))
	})

	t.Run("pending status blocks", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{prop("j2", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING)}
		require.ErrorContains(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_b"}), "node-2")
	})

	t.Run("deleted jobs are ignored", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.jobs[0].DeletedAt = timestamppb.Now()
		jd.proposals = []*jobv1.Proposal{prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED)}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a"}))
	})

	t.Run("nodes outside the target set are ignored", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.jobs = []*jobv1.Job{{Id: "j1", NodeId: "n1"}}
		jd.proposals = []*jobv1.Proposal{prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED)}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_unknown"}))
	})

	t.Run("rejected latest revision blocks", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED),
			prop("j1", 2, jobv1.ProposalStatus_PROPOSAL_STATUS_REJECTED),
		}
		require.ErrorContains(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a"}), "node-1")
	})

	t.Run("cancelled and revoked latest revisions pass", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_CANCELLED),
			prop("j2", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_REVOKED),
		}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a", "p2p_b"}))
	})

	t.Run("unmatched p2p id is skipped", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED)}
		require.NoError(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_unlabeled"}))
	})

	t.Run("JD list errors propagate", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.jobsErr = errors.New("boom")
		require.ErrorContains(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a"}), "boom")

		jd = base()
		jd.propsErr = errors.New("bang")
		require.ErrorContains(t, verifyNoPendingJobProposals(envWith(jd), []string{"p2p_a"}), "bang")
	})

	t.Run("bare hex p2p id still matches JD label", func(t *testing.T) {
		t.Parallel()
		var id p2pkey.PeerID
		id[0] = 1
		jd := &fakeJD{
			nodes:     []*nodev1.Node{jdNode("n1", "node-1", id.String())},
			jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
			proposals: []*jobv1.Proposal{prop("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED)},
		}
		require.ErrorContains(t, verifyNoPendingJobProposals(envWith(jd), []string{id.Raw()}), "node-1")
	})
}

func TestUpdateDONPreconditions_JobAcceptanceGate(t *testing.T) {
	t.Parallel()

	e := envWith(&fakeJD{
		nodes:     []*nodev1.Node{jdNode("n1", "node-1", "p2p_a")},
		jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
		proposals: []*jobv1.Proposal{{JobId: "j1", Revision: 1, Status: jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED}},
	})
	base := UpdateDONInput{DONName: "d", Nodes: []string{"p2p_a"}, SkipDuplicateCapabilityValidation: true}

	err := UpdateDON{}.VerifyPreconditions(e, base)
	require.ErrorContains(t, err, "node-1")

	base.SkipJobAcceptanceCheck = true
	require.NoError(t, UpdateDON{}.VerifyPreconditions(e, base))
}

func TestUpdateNodesPreconditions_JobAcceptanceGate(t *testing.T) {
	t.Parallel()

	e := envWith(&fakeJD{
		nodes:     []*nodev1.Node{jdNode("n1", "node-1", "p2p_a")},
		jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
		proposals: []*jobv1.Proposal{{JobId: "j1", Revision: 1, Status: jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED}},
	})
	in := UpdateNodesInput{Nodes: []CapabilitiesRegistryNodeParams{{
		P2pID: "p2p_a", NOP: "nop", Signer: "s", EncryptionPublicKey: "e", CsaKey: "c",
	}}}

	require.ErrorContains(t, UpdateNodes{}.VerifyPreconditions(e, in), "node-1")

	in.SkipJobAcceptanceCheck = true
	require.NoError(t, UpdateNodes{}.VerifyPreconditions(cldf.Environment{}, in))
}

func TestDonNodeP2PIDs(t *testing.T) {
	t.Parallel()

	var id p2pkey.PeerID
	id[0] = 7
	dons := []capabilities_registry_v2.CapabilitiesRegistryDONInfo{
		{Name: "a", NodeP2PIds: [][32]byte{id}},
		{Name: "b"},
	}

	ids, ok := donNodeP2PIDs(dons, "a")
	require.True(t, ok)
	require.Equal(t, []string{id.String()}, ids)

	_, ok = donNodeP2PIDs(dons, "missing")
	require.False(t, ok)
}
