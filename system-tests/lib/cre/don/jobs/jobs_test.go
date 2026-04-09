package jobs

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestApproveFetchesProposalsOncePerNode(t *testing.T) {
	restoreLoad := loadNodeProposalIDs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalIDs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	nodeA := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	nodeB := &cre.Node{Name: "node-b", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-b"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{nodeA, nodeB}}}}

	var mu sync.Mutex
	fetches := map[string]int{}
	loadNodeProposalIDs = func(_ context.Context, node *cre.Node) (map[string]string, error) {
		mu.Lock()
		fetches[node.JobDistributorDetails.NodeID]++
		mu.Unlock()
		if node.JobDistributorDetails.NodeID == "node-a" {
			return map[string]string{
				"spec-a-1": "proposal-a-1",
				"spec-a-2": "proposal-a-2",
			}, nil
		}
		return map[string]string{
			"spec-b-1": "proposal-b-1",
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

func TestApproveRunsAcrossNodesConcurrentlyAndWithinNodeSequentially(t *testing.T) {
	restoreLoad := loadNodeProposalIDs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalIDs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	nodeA := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	nodeB := &cre.Node{Name: "node-b", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-b"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{nodeA, nodeB}}}}

	loadNodeProposalIDs = func(_ context.Context, node *cre.Node) (map[string]string, error) {
		return map[string]string{
			node.JobDistributorDetails.NodeID + "-1": node.JobDistributorDetails.NodeID + "-proposal-1",
			node.JobDistributorDetails.NodeID + "-2": node.JobDistributorDetails.NodeID + "-proposal-2",
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
	restoreLoad := loadNodeProposalIDs
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		loadNodeProposalIDs = restoreLoad
		approveJobProposalSpec = restoreApprove
	})

	node := &cre.Node{Name: "node-a", JobDistributorDetails: &cre.JobDistributorDetails{NodeID: "node-a"}}
	dons := &cre.Dons{Dons: []*cre.Don{{Nodes: []*cre.Node{node}}}}

	loadNodeProposalIDs = func(_ context.Context, _ *cre.Node) (map[string]string, error) {
		return map[string]string{}, nil
	}
	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error { return nil }

	err := Approve(context.Background(), nil, dons, map[string][]string{
		"node-a": {"missing-spec"},
	})
	require.ErrorContains(t, err, "no job proposal found for job spec missing-spec")
}

func TestAcceptTreatsApprovedWorkflowSpecAsSuccess(t *testing.T) {
	restoreApprove := approveJobProposalSpec
	t.Cleanup(func() {
		approveJobProposalSpec = restoreApprove
	})

	approveJobProposalSpec = func(_ context.Context, _ *cre.Node, _ string) error {
		return errors.New("cannot approve an approved spec")
	}

	err := accept(context.Background(), &cre.Node{Name: "node-a"}, "proposal-id", `type = "workflow"`)
	require.NoError(t, err)
}
