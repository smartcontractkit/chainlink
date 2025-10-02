package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
)

// TestMultipleMCMSDeploymentsConflict demonstrates the gap where GetMCMSContracts
// cannot distinguish between multiple MCMS deployments on the same chain
func TestMultipleMCMSDeploymentsConflict(t *testing.T) {
	lggr := logger.Test(t)
	env, chainSelector := cre.BuildMinimalEnvironment(t, lggr)

	// Create Team A's MCMS config with qualifier
	teamAQualifier := "team-a"
	teamAConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	teamAConfig.Qualifier = &teamAQualifier

	teamATimelockCfgs := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		chainSelector: teamAConfig,
	}

	teamAEnv, err := commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamATimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team A's MCMS infrastructure")
	t.Log("Team A's MCMS infrastructure deployed successfully")

	// Get Team A's MCMS contracts using their qualifier
	teamAMCMSContracts, err := strategies.GetMCMSContracts(teamAEnv, chainSelector, teamAQualifier)
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
		chainSelector: teamBConfig,
	}

	teamBEnv, err := commonchangeset.Apply(t, teamAEnv, // Build on top of Team A's environment
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamBTimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team B's MCMS infrastructure")
	t.Log("Team B's MCMS infrastructure deployed successfully")

	// Get Team B's MCMS contracts using their qualifier
	teamBMCMSContracts, err := strategies.GetMCMSContracts(teamBEnv, chainSelector, teamBQualifier)
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
