package jobs

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// JobDistributorAPI is the subset of JD RPCs needed for job cleanup.
type JobDistributorAPI interface {
	ListJobs(ctx context.Context, in *jobv1.ListJobsRequest) (*jobv1.ListJobsResponse, error)
	DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest) (*jobv1.DeleteJobResponse, error)
}

// NewJobDistributorAdapter wraps a cldfjd.JobDistributor to satisfy JobDistributorAPI.
func NewJobDistributorAdapter(jd *cldfjd.JobDistributor) JobDistributorAPI {
	return &jdAdapter{jd}
}

type jdAdapter struct {
	jd *cldfjd.JobDistributor
}

func (a *jdAdapter) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest) (*jobv1.ListJobsResponse, error) {
	return a.jd.ListJobs(ctx, in)
}

func (a *jdAdapter) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest) (*jobv1.DeleteJobResponse, error) {
	return a.jd.DeleteJob(ctx, in)
}

// ProposalCanceller cancels approved JD job-proposal specs for a node via its GraphQL client.
type ProposalCanceller interface {
	ApprovedSpecIDs(ctx context.Context, node *cre.Node) ([]string, error)
	CancelSpec(ctx context.Context, node *cre.Node, specID string) error
}

// DeleteAllForDons cancels approved proposals then deletes every JD job for every node in dons,
// concurrently across nodes, rate-limited to protect JD (same 5 ops/sec bound local CRE uses).
//
// Job creation used to try to detect/repair conflicts with whatever already existed (matching by
// name, comparing spec content, deleting only what looked stale) — but JD's own proposal/job
// bookkeeping doesn't behave predictably enough for that: a node's job identity can drift across
// unrelated fixes (contract redeploys, corrected peering hosts, discovered OCR2 bundle IDs),
// leaving multiple never-approved Job entities with the same name, and selective repair either
// missed the real conflict or left JD itself in a stuck state (pending "New"/"Updates" proposals,
// nothing approved). Per direction: always delete everything for these nodes first, then let
// PostEnvStartup propose everything fresh. No diffing, no conditionals.
func DeleteAllForDons(ctx context.Context, log zerolog.Logger, jd JobDistributorAPI,
	proposals ProposalCanceller, dons *cre.Dons) error {
	rl := ratelimit.New(5)
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(8)

	for _, don := range dons.List() {
		for _, node := range don.Nodes {
			if node.JobDistributorDetails == nil || node.JobDistributorDetails.NodeID == "" {
				continue
			}
			g.Go(func() error {
				// Cancel approved specs (deletes running jobs + clears FK) — same ordering/reason as before
				if err := cancelApproved(gctx, log, proposals, don.Name, node); err != nil {
					return err
				}

				rl.Take()
				resp, err := jd.ListJobs(gctx, &jobv1.ListJobsRequest{
					Filter: &jobv1.ListJobsRequest_Filter{NodeIds: []string{node.JobDistributorDetails.NodeID}},
				})
				if err != nil {
					return errors.Wrapf(err, "failed to list jobs for node %s", node.Name)
				}

				for _, j := range resp.GetJobs() {
					rl.Take()
					if _, err := jd.DeleteJob(gctx, &jobv1.DeleteJobRequest{
						IdOneof: &jobv1.DeleteJobRequest_Id{Id: j.Id},
					}); err != nil {
						return errors.Wrapf(err, "failed to delete job %s for node %s", j.Id, node.Name)
					}
					log.Info().Str("don", don.Name).Str("node", node.Name).Str("jobID", j.Id).Msg("Deleted existing job")
				}
				return nil
			})
		}
	}

	return g.Wait()
}

// cancelApproved cancels every currently-approved job proposal spec for node, which deletes the
// underlying running job on the node and clears its job_proposals.external_job_id reference.
// Nodes without a GraphQL client or JD ID configured are skipped (nothing to cancel).
func cancelApproved(ctx context.Context, log zerolog.Logger, proposals ProposalCanceller, donName string, node *cre.Node) error {
	if node.Clients.GQLClient == nil || node.JobDistributorDetails.JDID == "" {
		return nil
	}

	specIDs, err := proposals.ApprovedSpecIDs(ctx, node)
	if err != nil {
		return errors.Wrapf(err, "failed to get job distributor details for node %s", node.Name)
	}

	for _, specID := range specIDs {
		if err := proposals.CancelSpec(ctx, node, specID); err != nil {
			return errors.Wrapf(err, "failed to cancel job proposal spec %s for node %s", specID, node.Name)
		}
		log.Info().Str("don", donName).Str("node", node.Name).Str("specID", specID).Msg("Cancelled approved job proposal spec")
	}

	return nil
}
