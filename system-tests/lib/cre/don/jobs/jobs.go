package jobs

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

var (
	acceptJobRetryMaxDuration = 90 * time.Second
	acceptJobRetryBackoff     = func() retry.Backoff { return retry.NewFibonacci(1 * time.Second) }
	acceptJobSettleDuration   = 5 * time.Second
	acceptJobSettleInterval   = 250 * time.Millisecond
	externalJobIDPattern      = regexp.MustCompile(`(?m)^\s*externalJobID\s*=\s*"([^"]+)"`)
)

type proposalRef struct {
	ProposalID string
	SpecID     string
}

// defined as variables to allow for easy testing
var loadNodeProposalRefs = func(ctx context.Context, node *cre.Node) (map[string]proposalRef, error) {
	jd, err := node.Clients.GQLClient.GetJobDistributor(ctx, node.JobDistributorDetails.JDID)
	if err != nil {
		return nil, err
	}
	if jd.GetJobProposals() == nil {
		return nil, fmt.Errorf("no job proposals found for node %s", node.Name)
	}

	proposalRefsBySpec := make(map[string]proposalRef, len(jd.JobProposals))
	for _, proposal := range jd.JobProposals {
		proposalRefsBySpec[proposal.LatestSpec.Definition] = proposalRef{
			ProposalID: proposal.Id,
			SpecID:     proposal.LatestSpec.Id,
		}
	}

	return proposalRefsBySpec, nil
}

// defined as variables to allow for easy testing
var approveJobProposalSpec = func(ctx context.Context, node *cre.Node, proposalID string) error {
	approvedSpec, err := node.Clients.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", proposalID)
	}

	return nil
}

var jobProposalIsApproved = func(ctx context.Context, node *cre.Node, proposalID string) (bool, error) {
	jd, err := node.Clients.GQLClient.GetJobDistributor(ctx, node.JobDistributorDetails.JDID)
	if err != nil {
		return false, err
	}
	if jd.GetJobProposals() == nil {
		return false, fmt.Errorf("no job proposals found for node %s", node.Name)
	}

	for _, proposal := range jd.JobProposals {
		if proposal.Id != proposalID {
			continue
		}

		return string(proposal.LatestSpec.Status) == "APPROVED", nil
	}

	return false, fmt.Errorf("no job proposal found for id %s", proposalID)
}

var jobProposalHasExistingJob = func(ctx context.Context, node *cre.Node, proposalID string) (bool, error) {
	if node == nil {
		return false, errors.New("node is nil")
	}
	if node.Clients.GQLClient == nil {
		return false, fmt.Errorf("node %s does not have a GraphQL client", node.Name)
	}

	proposal, err := node.Clients.GQLClient.GetJobProposal(ctx, proposalID)
	if err != nil {
		return false, err
	}
	if proposal == nil {
		return false, fmt.Errorf("no job proposal found for id %s", proposalID)
	}

	externalJobID := proposal.ExternalJobID
	if externalJobID == "" {
		externalJobID = extractExternalJobID(proposal.LatestSpec.Definition)
	}
	if externalJobID == "" {
		return false, nil
	}

	const pageSize = 1000
	offset := 0
	for {
		jobs, err := node.Clients.GQLClient.ListJobs(ctx, offset, pageSize)
		if err != nil {
			return false, err
		}
		if jobs == nil {
			return false, nil
		}

		results := jobs.Jobs.Results
		for _, job := range results {
			if job.ExternalJobID == externalJobID {
				return true, nil
			}
		}

		if len(results) < pageSize {
			return false, nil
		}
		offset += len(results)
	}
}

func extractExternalJobID(definition string) string {
	matches := externalJobIDPattern.FindStringSubmatch(definition)
	if len(matches) != 2 {
		return ""
	}

	return matches[1]
}

func Approve(ctx context.Context, _ cldf_offchain.Client, dons *cre.Dons, nodeToSpecs map[string][]string) error {
	nodeByID := make(map[string]*cre.Node)
	for _, don := range dons.List() {
		for _, node := range don.Nodes {
			nodeByID[node.JobDistributorDetails.NodeID] = node
		}
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(4)

	for nodeID, jobSpecs := range nodeToSpecs {
		node, ok := nodeByID[nodeID]
		if !ok {
			return fmt.Errorf("node with id %s not found", nodeID)
		}

		eg.Go(func() error {
			proposalRefsBySpec, err := loadNodeProposalRefs(egCtx, node)
			if err != nil {
				return err
			}

			for _, jobSpec := range uniqueSpecs(jobSpecs) {
				ref, ok := proposalRefsBySpec[jobSpec]
				if !ok {
					return fmt.Errorf("no job proposal found for job spec %s", jobSpec)
				}
				if err := accept(egCtx, node, ref, jobSpec); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return eg.Wait()
}

func uniqueSpecs(jobSpecs []string) []string {
	if len(jobSpecs) < 2 {
		return jobSpecs
	}

	seen := make(map[string]struct{}, len(jobSpecs))
	deduped := make([]string, 0, len(jobSpecs))
	for _, jobSpec := range jobSpecs {
		if _, ok := seen[jobSpec]; ok {
			continue
		}
		seen[jobSpec] = struct{}{}
		deduped = append(deduped, jobSpec)
	}

	return deduped
}

func Create(ctx context.Context, offChainClient cldf_offchain.Client, dons *cre.Dons, jobSpecs cre.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	jobRateLimit := ratelimit.New(5)

	for _, jobReq := range jobSpecs {
		eg.Go(func() error {
			jobRateLimit.Take()
			timeout := time.Second * 60
			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			_, pErr := offChainClient.ProposeJob(ctxWithTimeout, jobReq)
			if pErr != nil {
				fmt.Println("Failed jobspec proposal:")
				fmt.Println(jobReq)
				return fmt.Errorf("failed to propose job for node %s: %w", jobReq.NodeId, pErr)
			}

			for _, don := range dons.List() {
				for _, node := range don.Nodes {
					if node.JobDistributorDetails.NodeID != jobReq.NodeId {
						continue
					}

					if err := accept(ctx, node, proposalRef{}, jobReq.Spec); err != nil {
						return err
					}
				}
			}

			if ctx.Err() != nil {
				return errors.Wrapf(pErr, "timed out after %s proposing job for node %s", timeout.String(), jobReq.NodeId)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to create at least one job for DON")
	}

	return nil
}

func accept(ctx context.Context, node *cre.Node, proposal proposalRef, jobSpec string) error {
	retryErr := retry.Do(ctx, retry.WithMaxDuration(acceptJobRetryMaxDuration, acceptJobRetryBackoff()), func(ctx context.Context) error {
		var err error
		if proposal.SpecID == "" {
			err = node.AcceptJob(ctx, jobSpec)
		} else {
			err = approveJobProposalSpec(ctx, node, proposal.SpecID)
		}
		if err == nil || strings.Contains(err.Error(), "cannot approve an approved spec") {
			return nil
		}
		if proposal.ProposalID != "" {
			approved, approvedErr := jobProposalIsApproved(ctx, node, proposal.ProposalID)
			if approvedErr == nil && approved {
				return nil
			}
			existingJob, existingJobErr := jobProposalHasExistingJob(ctx, node, proposal.ProposalID)
			if existingJobErr == nil && existingJob {
				return nil
			}
			settled, settleErr := waitForProposalToSettle(ctx, node, proposal.ProposalID)
			if settleErr == nil && settled {
				return nil
			}
		}

		return retry.RetryableError(err)
	})
	if retryErr != nil {
		fmt.Println("Failed jobspec proposal for node ", node.Name)
		fmt.Println(jobSpec)

		return fmt.Errorf("failed to accept job for node %s. err: %w", node.Name, retryErr)
	}

	return nil
}

func waitForProposalToSettle(ctx context.Context, node *cre.Node, proposalID string) (bool, error) {
	if proposalID == "" || acceptJobSettleDuration <= 0 {
		return false, nil
	}

	deadline := time.Now().Add(acceptJobSettleDuration)
	for {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		approved, approvedErr := jobProposalIsApproved(ctx, node, proposalID)
		if approvedErr == nil && approved {
			return true, nil
		}

		existingJob, existingJobErr := jobProposalHasExistingJob(ctx, node, proposalID)
		if existingJobErr == nil && existingJob {
			return true, nil
		}

		if time.Now().After(deadline) {
			return false, nil
		}

		timer := time.NewTimer(acceptJobSettleInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false, ctx.Err()
		case <-timer.C:
		}
	}
}
