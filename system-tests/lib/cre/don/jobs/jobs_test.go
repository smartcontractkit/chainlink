package jobs

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestApproveFetchesProposalsOncePerNode(t *testing.T) {
	restoreLoad := loadNodeProposalRefs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalRefs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	nodeA := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	nodeB := &cre.Node{Name: "node-b", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-b"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{nodeA, nodeB}}}}

	var mu sync.Mutex
	fetches := map[string]int{}
	loadNodeProposalRefs = func(_ context.Context, node *cre.Node) (map[string]proposalRef, error) {
		mu.Lock()
		fetches[node.JobDistributorDetails.NodeID]++
		mu.Unlock()
		if node.JobDistributorDetails.NodeID == "node-a" {
			return map[string]proposalRef{
				"spec-a-1": {ProposalID: "proposal-a-1", SpecID: "spec-a-1-id"},
				"spec-a-2": {ProposalID: "proposal-a-2", SpecID: "spec-a-2-id"},
			}, nil
		}
		return map[string]proposalRef{
			"spec-b-1": {ProposalID: "proposal-b-1", SpecID: "spec-b-1-id"},
		}, nil
	}
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error { return nil }

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"spec-a-1", "spec-a-2"},
		"node-b": {"spec-b-1"},
	})
	require.NoError(t, err)
	require.Equal(t, map[string]int{"node-a": 1, "node-b": 1}, fetches)
}

func TestApproveDeduplicatesJobSpecsPerNode(t *testing.T) {
	restoreLoad := loadNodeProposalRefs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalRefs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	node := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{node}}}}

	loadNodeProposalRefs = func(_ context.Context, _ *cre.Node) (map[string]proposalRef, error) {
		return map[string]proposalRef{
			"spec-a": {ProposalID: "proposal-a", SpecID: "spec-a-id"},
			"spec-b": {ProposalID: "proposal-b", SpecID: "spec-b-id"},
		}, nil
	}

	var approved []string
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, specID string) error {
		approved = append(approved, specID)
		return nil
	}

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"spec-a", "spec-a", "spec-b", "spec-a", "spec-b"},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"spec-a-id", "spec-b-id"}, approved)
}

func TestApproveRunsAcrossNodesConcurrentlyAndWithinNodeSequentially(t *testing.T) {
	restoreLoad := loadNodeProposalRefs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalRefs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	nodeA := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	nodeB := &cre.Node{Name: "node-b", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-b"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{nodeA, nodeB}}}}

	loadNodeProposalRefs = func(_ context.Context, node *cre.Node) (map[string]proposalRef, error) {
		return map[string]proposalRef{
			node.JobDistributorDetails.NodeID + "-1": {ProposalID: node.JobDistributorDetails.NodeID + "-proposal-1", SpecID: node.JobDistributorDetails.NodeID + "-spec-1"},
			node.JobDistributorDetails.NodeID + "-2": {ProposalID: node.JobDistributorDetails.NodeID + "-proposal-2", SpecID: node.JobDistributorDetails.NodeID + "-spec-2"},
		}, nil
	}

	var mu sync.Mutex
	activeGlobal := 0
	maxGlobal := 0
	activePerNode := map[string]int{}
	maxPerNode := map[string]int{}
	approveJobProposalSpec = func(_ context.Context, node *cre.Node, _ string) error {
		nodeID := node.JobDistributorDetails.NodeID

		mu.Lock()
		activeGlobal++
		activePerNode[nodeID]++
		if activeGlobal > maxGlobal {
			maxGlobal = activeGlobal
		}
		if activePerNode[nodeID] > maxPerNode[nodeID] {
			maxPerNode[nodeID] = activePerNode[nodeID]
		}
		mu.Unlock()

		time.Sleep(25 * time.Millisecond)

		mu.Lock()
		activeGlobal--
		activePerNode[nodeID]--
		mu.Unlock()
		return nil
	}

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"node-a-1", "node-a-2"},
		"node-b": {"node-b-1", "node-b-2"},
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, maxGlobal, 2)
	require.Equal(t, 1, maxPerNode["node-a"])
	require.Equal(t, 1, maxPerNode["node-b"])
}

func TestApproveMissingNodeID(t *testing.T) {
	err := Approve(context.Background(), nil, &cre.Dons{}, map[string][]string{
		"missing-node": {"spec"},
	})
	require.ErrorContains(t, err, "node with id missing-node not found")
}

func TestApproveMissingProposalMatch(t *testing.T) {
	restoreLoad := loadNodeProposalRefs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalRefs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	node := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{node}}}}

	loadNodeProposalRefs = func(_ context.Context, _ *cre.Node) (map[string]proposalRef, error) {
		return map[string]proposalRef{}, nil
	}
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error { return nil }

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"missing-spec"},
	})
	require.ErrorContains(t, err, "no job proposal found for job spec missing-spec")
}

func TestAcceptTreatsApprovedWorkflowSpecAsSuccess(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	restoreApprovedCheck := jobProposalIsApproved
	restoreExistingJob := jobProposalHasExistingJob
	restoreRetryDuration := acceptJobRetryMaxDuration
	restoreRetryBackoff := acceptJobRetryBackoff
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
		jobProposalIsApproved = restoreApprovedCheck
		jobProposalHasExistingJob = restoreExistingJob
		acceptJobRetryMaxDuration = restoreRetryDuration
		acceptJobRetryBackoff = restoreRetryBackoff
	})
	acceptJobRetryMaxDuration = 10 * time.Millisecond
	acceptJobRetryBackoff = func() retry.Backoff { return retry.NewConstant(time.Millisecond) }

	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		return errors.New("cannot approve an approved spec")
	}
	jobProposalHasExistingJob = func(_ context.Context, _ *cre.Node, _ string) (bool, error) { return false, nil }

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, proposalRef{ProposalID: "proposal-id", SpecID: "spec-id"}, `type = "workflow"`)
	require.NoError(t, err)
}

func TestAcceptTreatsApprovedSpecAfterErrorAsSuccess(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	restoreApprovedCheck := jobProposalIsApproved
	restoreExistingJob := jobProposalHasExistingJob
	restoreRetryDuration := acceptJobRetryMaxDuration
	restoreRetryBackoff := acceptJobRetryBackoff
	restoreSettleDuration := acceptJobSettleDuration
	restoreSettleInterval := acceptJobSettleInterval
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
		jobProposalIsApproved = restoreApprovedCheck
		jobProposalHasExistingJob = restoreExistingJob
		acceptJobRetryMaxDuration = restoreRetryDuration
		acceptJobRetryBackoff = restoreRetryBackoff
		acceptJobSettleDuration = restoreSettleDuration
		acceptJobSettleInterval = restoreSettleInterval
	})
	acceptJobRetryMaxDuration = 10 * time.Millisecond
	acceptJobRetryBackoff = func() retry.Backoff { return retry.NewConstant(time.Millisecond) }
	acceptJobSettleDuration = 10 * time.Millisecond
	acceptJobSettleInterval = time.Millisecond

	attempts := 0
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		attempts++
		return errors.New("input: approveJobProposalSpec could not approve job proposal: error registering vault capability: capability already exists: id vault@1.0.0 found in registry: state INVALID_STATE")
	}

	statusChecks := 0
	jobProposalIsApproved = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		statusChecks++
		return true, nil
	}
	jobProposalHasExistingJob = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		return false, nil
	}

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, proposalRef{ProposalID: "proposal-id", SpecID: "spec-id"}, `type = "offchainreporting2"`)
	require.NoError(t, err)
	require.Equal(t, 1, attempts)
	require.Equal(t, 1, statusChecks)
}

func TestAcceptTreatsExistingJobAfterErrorAsSuccess(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	restoreApprovedCheck := jobProposalIsApproved
	restoreExistingJob := jobProposalHasExistingJob
	restoreRetryDuration := acceptJobRetryMaxDuration
	restoreRetryBackoff := acceptJobRetryBackoff
	restoreSettleDuration := acceptJobSettleDuration
	restoreSettleInterval := acceptJobSettleInterval
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
		jobProposalIsApproved = restoreApprovedCheck
		jobProposalHasExistingJob = restoreExistingJob
		acceptJobRetryMaxDuration = restoreRetryDuration
		acceptJobRetryBackoff = restoreRetryBackoff
		acceptJobSettleDuration = restoreSettleDuration
		acceptJobSettleInterval = restoreSettleInterval
	})
	acceptJobRetryMaxDuration = 10 * time.Millisecond
	acceptJobRetryBackoff = func() retry.Backoff { return retry.NewConstant(time.Millisecond) }
	acceptJobSettleDuration = 10 * time.Millisecond
	acceptJobSettleInterval = time.Millisecond

	attempts := 0
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		attempts++
		return errors.New("input: approveJobProposalSpec could not approve job proposal: error registering vault capability: capability already exists: id vault@1.0.0 found in registry: state INVALID_STATE")
	}

	approvalChecks := 0
	jobProposalIsApproved = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		approvalChecks++
		return false, nil
	}

	existingJobChecks := 0
	jobProposalHasExistingJob = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		existingJobChecks++
		return true, nil
	}

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, proposalRef{ProposalID: "proposal-id", SpecID: "spec-id"}, `type = "offchainreporting2"`)
	require.NoError(t, err)
	require.Equal(t, 1, attempts)
	require.Equal(t, 1, approvalChecks)
	require.Equal(t, 1, existingJobChecks)
}

func TestAcceptRetriesTransientApprovalErrors(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	restoreApprovedCheck := jobProposalIsApproved
	restoreExistingJob := jobProposalHasExistingJob
	restoreRetryDuration := acceptJobRetryMaxDuration
	restoreRetryBackoff := acceptJobRetryBackoff
	restoreSettleDuration := acceptJobSettleDuration
	restoreSettleInterval := acceptJobSettleInterval
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
		jobProposalIsApproved = restoreApprovedCheck
		jobProposalHasExistingJob = restoreExistingJob
		acceptJobRetryMaxDuration = restoreRetryDuration
		acceptJobRetryBackoff = restoreRetryBackoff
		acceptJobSettleDuration = restoreSettleDuration
		acceptJobSettleInterval = restoreSettleInterval
	})
	acceptJobRetryMaxDuration = 50 * time.Millisecond
	acceptJobRetryBackoff = func() retry.Backoff { return retry.NewConstant(time.Millisecond) }
	acceptJobSettleDuration = 5 * time.Millisecond
	acceptJobSettleInterval = time.Millisecond

	attempts := 0
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		attempts++
		if attempts < 3 {
			return errors.New("approveJobProposalSpec panic occurred: runtime error: invalid memory address or nil pointer dereference")
		}
		return nil
	}
	jobProposalIsApproved = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		return false, nil
	}
	jobProposalHasExistingJob = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		return false, nil
	}

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, proposalRef{ProposalID: "proposal-id", SpecID: "spec-id"}, `type = "offchainreporting2"`)
	require.NoError(t, err)
	require.Equal(t, 3, attempts)
}

func TestAcceptWaitsForProposalToSettleAfterError(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	restoreApprovedCheck := jobProposalIsApproved
	restoreExistingJob := jobProposalHasExistingJob
	restoreRetryDuration := acceptJobRetryMaxDuration
	restoreRetryBackoff := acceptJobRetryBackoff
	restoreSettleDuration := acceptJobSettleDuration
	restoreSettleInterval := acceptJobSettleInterval
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
		jobProposalIsApproved = restoreApprovedCheck
		jobProposalHasExistingJob = restoreExistingJob
		acceptJobRetryMaxDuration = restoreRetryDuration
		acceptJobRetryBackoff = restoreRetryBackoff
		acceptJobSettleDuration = restoreSettleDuration
		acceptJobSettleInterval = restoreSettleInterval
	})
	acceptJobRetryMaxDuration = 20 * time.Millisecond
	acceptJobRetryBackoff = func() retry.Backoff { return retry.NewConstant(time.Millisecond) }
	acceptJobSettleDuration = 20 * time.Millisecond
	acceptJobSettleInterval = time.Millisecond

	attempts := 0
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		attempts++
		return errors.New("input: approveJobProposalSpec could not approve job proposal: error registering vault capability: capability already exists: id vault@1.0.0 found in registry: state INVALID_STATE")
	}

	approvalChecks := 0
	jobProposalIsApproved = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		approvalChecks++
		return approvalChecks >= 3, nil
	}
	jobProposalHasExistingJob = func(_ context.Context, _ *cre.Node, _ string) (bool, error) {
		return false, nil
	}

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, proposalRef{ProposalID: "proposal-id", SpecID: "spec-id"}, `type = "offchainreporting2"`)
	require.NoError(t, err)
	require.Equal(t, 1, attempts)
	require.GreaterOrEqual(t, approvalChecks, 3)
}

func TestApproveUsesLatestSpecIDForApproval(t *testing.T) {
	restoreLoad := loadNodeProposalRefs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalRefs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	node := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{node}}}}

	loadNodeProposalRefs = func(_ context.Context, _ *cre.Node) (map[string]proposalRef, error) {
		return map[string]proposalRef{
			"spec-a": {ProposalID: "proposal-id", SpecID: "latest-spec-id"},
		}, nil
	}

	var approvedID string
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, proposalID string) error {
		approvedID = proposalID
		return nil
	}

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"spec-a"},
	})
	require.NoError(t, err)
	require.Equal(t, "latest-spec-id", approvedID)
}

func TestExtractExternalJobID(t *testing.T) {
	t.Run("extracts from job spec", func(t *testing.T) {
		definition := `type = "offchainreporting2"
externalJobID = "123e4567-e89b-12d3-a456-426614174000"
name = "vault-worker"`

		require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", extractExternalJobID(definition))
	})

	t.Run("returns empty string when missing", func(t *testing.T) {
		require.Empty(t, extractExternalJobID(`type = "offchainreporting2"`))
	})
}
