package jobs

import (
	"context"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	webclient "github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// fakeAnyGQLClient is a non-nil stand-in for node.Clients.GQLClient.
type fakeAnyGQLClient struct {
	webclient.Client
}

// fakeJobDistributor fakes the JobDistributorAPI interface. DeleteAllForDons
// calls it concurrently (one goroutine per node), so deletedJobIDs is guarded
// by a mutex.
type fakeJobDistributor struct {
	jobsByNode map[string][]*jobv1.Job

	mu            sync.Mutex
	deletedJobIDs []string
}

func (f *fakeJobDistributor) ListJobs(_ context.Context, in *jobv1.ListJobsRequest) (*jobv1.ListJobsResponse, error) {
	var jobs []*jobv1.Job
	for _, nodeID := range in.GetFilter().GetNodeIds() {
		jobs = append(jobs, f.jobsByNode[nodeID]...)
	}
	return &jobv1.ListJobsResponse{Jobs: jobs}, nil
}

func (f *fakeJobDistributor) DeleteJob(_ context.Context, in *jobv1.DeleteJobRequest) (*jobv1.DeleteJobResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deletedJobIDs = append(f.deletedJobIDs, in.GetId())
	return &jobv1.DeleteJobResponse{}, nil
}

// fakeProposals implements ProposalCanceller for tests. Also called
// concurrently by DeleteAllForDons, so cancelledSpecIDs is guarded by a mutex.
type fakeProposals struct {
	approvedSpecIDsByNode map[string][]string

	mu               sync.Mutex
	cancelledSpecIDs []string
}

func (f *fakeProposals) ApprovedSpecIDs(_ context.Context, node *cre.Node) ([]string, error) {
	return f.approvedSpecIDsByNode[node.Name], nil
}

func (f *fakeProposals) CancelSpec(_ context.Context, _ *cre.Node, specID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cancelledSpecIDs = append(f.cancelledSpecIDs, specID)
	return nil
}

func TestDeleteAllForDons_CancelsApprovedSpecsBeforeDeleting(t *testing.T) {
	t.Parallel()

	fakeProposals := &fakeProposals{
		approvedSpecIDsByNode: map[string][]string{
			"node-1": {"spec-approved-1"},
		},
	}

	fake := &fakeJobDistributor{
		jobsByNode: map[string][]*jobv1.Job{"node-1": {{Id: "job-a"}}},
	}

	dons := &cre.Dons{Dons: []*cre.Don{
		{
			Name: "workflow",
			Nodes: []*cre.Node{
				{
					Name:                  "node-1",
					Clients:               cre.NodeClients{GQLClient: fakeAnyGQLClient{}},
					JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-1", JDID: "jd-1"},
				},
			},
		},
	}}

	err := DeleteAllForDons(context.Background(), zerolog.Nop(), fake, fakeProposals, dons)
	require.NoError(t, err)
	require.Equal(t, []string{"spec-approved-1"}, fakeProposals.cancelledSpecIDs)
	require.Equal(t, []string{"job-a"}, fake.deletedJobIDs)
}

func TestDeleteAllForDons_DeletesEveryJobPerNode(t *testing.T) {
	t.Parallel()

	fakeProposals := &fakeProposals{
		approvedSpecIDsByNode: map[string][]string{},
	}

	fake := &fakeJobDistributor{
		jobsByNode: map[string][]*jobv1.Job{
			"node-1": {{Id: "job-a"}, {Id: "job-b"}},
			"node-2": {{Id: "job-c"}},
		},
	}

	dons := &cre.Dons{Dons: []*cre.Don{
		{
			Name: "workflow",
			Nodes: []*cre.Node{
				{Name: "node-1", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-1"}},
				{Name: "node-2", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-2"}},
			},
		},
	}}

	err := DeleteAllForDons(context.Background(), zerolog.Nop(), fake, fakeProposals, dons)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"job-a", "job-b", "job-c"}, fake.deletedJobIDs)
}

func TestDeleteAllForDons_SkipsNodesWithoutJDDetails(t *testing.T) {
	t.Parallel()

	fakeProposals := &fakeProposals{
		approvedSpecIDsByNode: map[string][]string{},
	}

	fake := &fakeJobDistributor{
		jobsByNode: map[string][]*jobv1.Job{"node-1": {{Id: "job-a"}}},
	}

	dons := &cre.Dons{Dons: []*cre.Don{
		{
			Name: "workflow",
			Nodes: []*cre.Node{
				{Name: "node-1", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-1"}},
				{Name: "unregistered-node"},
			},
		},
	}}

	err := DeleteAllForDons(context.Background(), zerolog.Nop(), fake, fakeProposals, dons)
	require.NoError(t, err)
	require.Equal(t, []string{"job-a"}, fake.deletedJobIDs)
}
