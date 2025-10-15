package changeset

import (
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
)

// TestMultipleMCMSDeploymentsConflict demonstrates the gap where GetMCMSContracts
// cannot distinguish between multiple MCMS deployments on the same chain
func TestMultipleMCMSDeploymentsConflict(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// Create Team A's MCMS config with qualifier
	teamAQualifier := "team-a"
	teamAConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	teamAConfig.Qualifier = &teamAQualifier

	teamATimelockCfgs := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		selector: teamAConfig,
	}

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamATimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team A's MCMS infrastructure")
	t.Log("Team A's MCMS infrastructure deployed successfully")

	// Get Team A's MCMS contracts using their qualifier
	teamAMCMSContracts, err := strategies.GetMCMSContracts(rt.Environment(), selector, teamAQualifier)
	require.NoError(t, err, "should be able to get Team A's MCMS contracts")
	require.NotNil(t, teamAMCMSContracts, "Team A's MCMS contracts should not be nil")

	teamATimelockAddr := teamAMCMSContracts.Timelock.Address()
	teamAProposerAddr := teamAMCMSContracts.ProposerMcm.Address()
	t.Logf("Team A - Timelock: %s, Proposer: %s", teamATimelockAddr.Hex(), teamAProposerAddr.Hex())

	// Create Team B's MCMS config with different qualifier
	teamBQualifier := "team-b"
	teamBConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	teamBConfig.Qualifier = &teamBQualifier

	teamBTimelockCfgs := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		selector: teamBConfig,
	}

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamBTimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team B's MCMS infrastructure")
	t.Log("Team B's MCMS infrastructure deployed successfully")

	// Get Team B's MCMS contracts using their qualifier
	teamBMCMSContracts, err := strategies.GetMCMSContracts(rt.Environment(), selector, teamBQualifier)
	require.NoError(t, err, "should be able to get Team B's MCMS contracts with their qualifier")
	require.NotNil(t, teamBMCMSContracts, "Team B's MCMS contracts should not be nil")

	teamBTimelockAddr := teamBMCMSContracts.Timelock.Address()
	teamBProposerAddr := teamBMCMSContracts.ProposerMcm.Address()
	t.Logf("Team B - Timelock: %s, Proposer: %s", teamBTimelockAddr.Hex(), teamBProposerAddr.Hex())

	// Verify that each team has different MCMS contracts (true multi-tenancy)
	require.NotEqual(t, teamATimelockAddr, teamBTimelockAddr,
		"Team A and Team B should have different timelock contracts")
	require.NotEqual(t, teamAProposerAddr, teamBProposerAddr,
		"Team A and Team B should have different proposer contracts")
}
