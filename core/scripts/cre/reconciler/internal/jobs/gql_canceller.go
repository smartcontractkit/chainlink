package jobs

import (
	"context"

	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// GQLProposalCanceller implements ProposalCanceller using GraphQL calls.
type GQLProposalCanceller struct{}

func (GQLProposalCanceller) ApprovedSpecIDs(ctx context.Context, node *cre.Node) ([]string, error) {
	fm, err := node.Clients.GQLClient.GetJobDistributor(ctx, node.JobDistributorDetails.JDID)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, p := range fm.JobProposals {
		if p.LatestSpec.Status == "APPROVED" {
			ids = append(ids, p.LatestSpec.Id)
		}
	}
	return ids, nil
}

func (GQLProposalCanceller) CancelSpec(ctx context.Context, node *cre.Node, specID string) error {
	_, err := node.Clients.GQLClient.CancelJobProposalSpec(ctx, specID)
	return err
}
