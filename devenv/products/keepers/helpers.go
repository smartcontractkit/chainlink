package keepers

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
)

// DeployKeeperConsumersPerformance sequentially deploys keeper performance consumer contracts.
func DeployKeeperConsumersPerformance(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) []contracts.KeeperConsumerPerformance {
	l := framework.L
	upkeeps := make([]contracts.KeeperConsumerPerformance, 0)

	for contractCount := range numberOfContracts {
		// Deploy consumer
		keeperConsumerInstance, err := contracts.DeployKeeperConsumerPerformance(
			client,
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
		)
		require.NoError(t, err, "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, keeperConsumerInstance)
		l.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Performance Contract")
	}

	require.Len(t, upkeeps, numberOfContracts, "Incorrect number of consumers contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

// DeployPerformDataChecker sequentially deploys keeper perform data checker contracts.
func DeployPerformDataChecker(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	expectedData []byte,
) []contracts.KeeperPerformDataChecker {
	l := framework.L
	upkeeps := make([]contracts.KeeperPerformDataChecker, 0)

	for contractCount := range numberOfContracts {
		performDataCheckerInstance, err := contracts.DeployKeeperPerformDataChecker(client, expectedData)
		require.NoError(t, err, "Deploying KeeperPerformDataChecker instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, performDataCheckerInstance)
		l.Debug().
			Str("Contract Address", performDataCheckerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed PerformDataChecker Contract")
	}
	require.Len(t, upkeeps, numberOfContracts, "Incorrect number of PerformDataChecker contracts deployed")
	l.Info().Msg("Successfully deployed all PerformDataChecker Contracts")

	return upkeeps
}
