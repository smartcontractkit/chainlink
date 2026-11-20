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

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
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

// consensusSpecWith renders a consensus job spec with a key bundle per family.
func consensusSpecWith(families ...string) string {
	spec := "type = \"standardcapabilities\"\nname = \"OCR-Consensus-CRE (node)\"\ncommand = \"/usr/local/bin/consensus\"\n" +
		"[oracle_factory]\nenabled = true\n[oracle_factory.onchainSigningStrategy]\nstrategyName = \"multi-chain\"\n" +
		"[oracle_factory.onchainSigningStrategy.config]\nevm = \"bundle-evm\"\n"
	for _, f := range families {
		spec += f + " = \"bundle-" + f + "\"\n"
	}

	return spec
}

func proposal(job string, rev int64, st jobv1.ProposalStatus, spec string) *jobv1.Proposal {
	return &jobv1.Proposal{Id: job + "-p", JobId: job, Revision: rev, Status: st, Spec: spec}
}

func TestVerifySignerFamilyJobs(t *testing.T) {
	t.Parallel()

	approved := jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED
	proposed := jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED
	stellar := []string{"stellar"}

	base := func() *fakeJD {
		return &fakeJD{
			nodes: []*nodev1.Node{jdNode("n1", "node-1", "p2p_a"), jdNode("n2", "node-2", "p2p_b")},
			jobs:  []*jobv1.Job{{Id: "j1", NodeId: "n1"}, {Id: "j2", NodeId: "n2"}},
		}
	}

	t.Run("no families is a no-op without JD", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, verifySignerFamilyJobs(cldf.Environment{}, []string{"p2p_a"}, nil))
	})

	t.Run("missing JD client errors", func(t *testing.T) {
		t.Parallel()
		err := verifySignerFamilyJobs(cldf.Environment{}, []string{"p2p_a"}, stellar)
		require.ErrorContains(t, err, "job distributor client")
	})

	t.Run("approved consensus job with the family passes", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			proposal("j1", 1, approved, consensusSpecWith("stellar")),
			proposal("j2", 1, approved, consensusSpecWith("stellar")),
		}
		require.NoError(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a", "p2p_b"}, stellar))
	})

	t.Run("job not proposed at all blocks (INCIDENT-2694)", func(t *testing.T) {
		t.Parallel()
		jd := base()
		err := verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar)
		require.ErrorContains(t, err, "node-1 (stellar)")
	})

	t.Run("proposed but not accepted blocks", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{proposal("j1", 1, proposed, consensusSpecWith("stellar"))}
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "node-1")
	})

	t.Run("only the failing node is named", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			proposal("j1", 1, approved, consensusSpecWith("stellar")),
			proposal("j2", 1, proposed, consensusSpecWith("stellar")),
		}
		err := verifySignerFamilyJobs(envWith(jd), []string{"p2p_a", "p2p_b"}, stellar)
		require.ErrorContains(t, err, "node-2")
		require.NotContains(t, err.Error(), "node-1")
	})

	t.Run("approved job without the family blocks", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{proposal("j1", 1, approved, consensusSpecWith("aptos"))}
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "node-1 (stellar)")
	})

	t.Run("every requested family must be covered", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{proposal("j1", 1, approved, consensusSpecWith("stellar"))}
		err := verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, []string{"aptos", "stellar"})
		require.ErrorContains(t, err, "node-1 (aptos)")
	})

	t.Run("non-consensus job does not count", func(t *testing.T) {
		t.Parallel()
		jd := base()
		spec := "type = \"standardcapabilities\"\ncommand = \"/usr/local/bin/evm\"\n" +
			"[oracle_factory.onchainSigningStrategy.config]\nstellar = \"bundle-stellar\"\n"
		jd.proposals = []*jobv1.Proposal{proposal("j1", 1, approved, spec)}
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "node-1")
	})

	t.Run("only the latest revision counts", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.proposals = []*jobv1.Proposal{
			proposal("j1", 1, approved, consensusSpecWith("stellar")),
			proposal("j1", 2, proposed, consensusSpecWith()),
		}
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "node-1")
	})

	t.Run("deleted jobs do not count", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.jobs[0].DeletedAt = timestamppb.Now()
		jd.proposals = []*jobv1.Proposal{proposal("j1", 1, approved, consensusSpecWith("stellar"))}
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "node-1")
	})

	t.Run("node unknown to JD blocks", func(t *testing.T) {
		t.Parallel()
		jd := base()
		err := verifySignerFamilyJobs(envWith(jd), []string{"p2p_missing"}, stellar)
		require.ErrorContains(t, err, "p2p_missing (not in JD)")
	})

	t.Run("any JD record of a peer may hold the job", func(t *testing.T) {
		t.Parallel()
		jd := &fakeJD{
			nodes:     []*nodev1.Node{jdNode("n1", "old-record", "p2p_a"), jdNode("n2", "new-record", "p2p_a")},
			jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
			proposals: []*jobv1.Proposal{proposal("j1", 1, approved, consensusSpecWith("stellar"))},
		}
		require.NoError(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar))
	})

	t.Run("bare hex IDs match a prefixed JD label and vice versa", func(t *testing.T) {
		t.Parallel()
		var id p2pkey.PeerID
		id[0] = 3
		jd := &fakeJD{
			nodes:     []*nodev1.Node{jdNode("n1", "node-1", id.Raw())},
			jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
			proposals: []*jobv1.Proposal{proposal("j1", 1, approved, consensusSpecWith("stellar"))},
		}
		require.NoError(t, verifySignerFamilyJobs(envWith(jd), []string{id.String()}, stellar))
		require.NoError(t, verifySignerFamilyJobs(envWith(jd), []string{id.Raw()}, stellar))
	})

	t.Run("JD list errors propagate", func(t *testing.T) {
		t.Parallel()
		jd := base()
		jd.jobsErr = errors.New("boom")
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "boom")

		jd = base()
		jd.propsErr = errors.New("bang")
		require.ErrorContains(t, verifySignerFamilyJobs(envWith(jd), []string{"p2p_a"}, stellar), "bang")
	})
}

func TestConsensusSpecFamilies(t *testing.T) {
	t.Parallel()

	require.Equal(t, map[string]bool{"evm": true, "stellar": true}, consensusSpecFamilies(consensusSpecWith("stellar")))
	require.Empty(t, consensusSpecFamilies("not toml ["))
	require.Empty(t, consensusSpecFamilies(`type = "offchainreporting2"`+"\n"+`command = "consensus"`))
	require.Empty(t, consensusSpecFamilies(`type = "standardcapabilities"`+"\n"+`command = "/bin/consensus"`))

	empty := "type = \"standardcapabilities\"\ncommand = \"consensus\"\n[oracle_factory.onchainSigningStrategy.config]\nstellar = \"\"\n"
	require.Empty(t, consensusSpecFamilies(empty))
}

func ocr3Entry(families ...any) map[string]any {
	entry := map[string]any{"offchainConfig": map[string]any{"maxFaultyOracles": 1}}
	if families != nil {
		entry["extraSignerFamilies"] = families
	}

	return map[string]any{"ocr3Configs": map[string]any{"__default__": entry}}
}

func TestRequestedSignerFamilies(t *testing.T) {
	t.Parallel()

	got, err := requestedSignerFamilies([]ocr3CapConfig{
		{CapabilityID: "consensus@1.0.0", Config: ocr3Entry("stellar", "aptos")},
		{CapabilityID: "vault@1.0.0", Config: ocr3Entry("aptos")},
		{CapabilityID: "other@1.0.0", Config: map[string]any{"k": "v"}},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"aptos", "stellar"}, got)

	got, err = requestedSignerFamilies([]ocr3CapConfig{{CapabilityID: "c", Config: ocr3Entry()}})
	require.NoError(t, err)
	require.Empty(t, got)

	processed := map[string]any{"ocr3Configs": map[string]any{"__default__": map[string]any{
		"offchainConfig": "base64==", "extraSignerFamilies": []any{"stellar"},
	}}}
	got, err = requestedSignerFamilies([]ocr3CapConfig{{CapabilityID: "c", Config: processed}})
	require.NoError(t, err)
	require.Empty(t, got)

	_, err = requestedSignerFamilies([]ocr3CapConfig{{CapabilityID: "c", Config: ocr3Entry("evm")}})
	require.Error(t, err)
}

func TestDonNodeP2PIDs(t *testing.T) {
	t.Parallel()

	var id p2pkey.PeerID
	id[0] = 7
	don := &capabilities_registry_v2.CapabilitiesRegistryDONInfo{NodeP2PIds: [][32]byte{id}}
	require.Equal(t, []string{id.String()}, donNodeP2PIDs(don))
}

func TestUpdateDONPreconditions_SignerFamilyJobCheck(t *testing.T) {
	t.Parallel()

	e := envWith(&fakeJD{nodes: []*nodev1.Node{jdNode("n1", "node-1", "p2p_a")}, jobs: []*jobv1.Job{{Id: "j1", NodeId: "n1"}}})
	in := UpdateDONInput{
		DONName: "d", Nodes: []string{"p2p_a"}, SkipDuplicateCapabilityValidation: true,
		CapabilityConfigs: []contracts.CapabilityConfig{{
			Capability: contracts.Capability{CapabilityID: "consensus@1.0.0"},
			Config:     ocr3Entry("stellar"),
		}},
	}

	require.ErrorContains(t, UpdateDON{}.VerifyPreconditions(e, in), "node-1 (stellar)")

	withJob := e
	withJob.Offchain = &fakeJD{
		nodes:     []*nodev1.Node{jdNode("n1", "node-1", "p2p_a")},
		jobs:      []*jobv1.Job{{Id: "j1", NodeId: "n1"}},
		proposals: []*jobv1.Proposal{proposal("j1", 1, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, consensusSpecWith("stellar"))},
	}
	require.NoError(t, UpdateDON{}.VerifyPreconditions(withJob, in))

	in.SkipSignerFamilyJobCheck = true
	require.NoError(t, UpdateDON{}.VerifyPreconditions(cldf.Environment{}, in))

	in.SkipSignerFamilyJobCheck = false
	in.CapabilityConfigs = []contracts.CapabilityConfig{{
		Capability: contracts.Capability{CapabilityID: "consensus@1.0.0"},
		Config:     ocr3Entry(),
	}}
	require.NoError(t, UpdateDON{}.VerifyPreconditions(cldf.Environment{}, in), "no families requested: JD is not needed")
}
