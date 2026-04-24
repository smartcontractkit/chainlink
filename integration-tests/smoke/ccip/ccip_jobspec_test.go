package ccip

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// It always runs in docker, it's not enabled to run in-memory as we are testing the actual job distributor
func TestDeleteCCIPJobs(t *testing.T) {
	e, _, tenv := testsetups.NewIntegrationEnvironment(t, testhelpers.WithJobsOnly())
	nopsView, err := view.GenerateNopsView(e.Env.Logger, e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// gather all the jobIDs
	jobIDs := make([]string, 0)
	jobUUIDsByNode := make(map[string][]string)
	for _, nop := range nopsView {
		jobIDs = append(jobIDs, maps.Keys(nop.ApprovedJobspecs)...)
		for _, job := range nop.ApprovedJobspecs {
			jobUUIDsByNode[nop.NodeID] = append(jobUUIDsByNode[nop.NodeID], job.UUID)
		}
	}
	// run delete JobChangeset
	_, err = commonChangesets.Apply(t, e.Env,
		commonChangesets.Configure(
			commonChangesets.DeleteJobChangeset,
			jobIDs,
		),
	)
	require.NoError(t, err)
	// now cancel the jobs from the nodes
	require.NoError(t, tenv.DeleteJobs(e.Env.GetContext(), jobUUIDsByNode))

	// check if the jobs are deleted
	nopsView, err = view.GenerateNopsView(e.Env.Logger, e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for _, nop := range nopsView {
		require.Empty(t, nop.ApprovedJobspecs)
	}
}

// It always runs in docker, it's not enabled to run in-memory as we are testing the actual job distributor
func TestRevokeJobs(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-566")

	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithJobsOnly())
	nopsView, err := view.GenerateNopsView(e.Env.Logger, e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// gather all the jobIDs
	jobIDs := make([]string, 0)
	jobUUIDsByNode := make(map[string][]string)
	for _, nop := range nopsView {
		jobIDs = append(jobIDs, maps.Keys(nop.ApprovedJobspecs)...)
		for _, job := range nop.ApprovedJobspecs {
			jobUUIDsByNode[nop.NodeID] = append(jobUUIDsByNode[nop.NodeID], job.UUID)
		}
	}
	// run RevokeJobChangeset
	_, err = commonChangesets.Apply(t, e.Env,
		commonChangesets.Configure(
			commonChangesets.RevokeJobsChangeset,
			jobIDs,
		),
	)
	// currently all test nodes auto accept proposals so we can't revoke them
	// therefore testing the error message
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not in PROPOSED or CANCELLED state")
}
