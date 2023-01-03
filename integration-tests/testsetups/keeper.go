package testsetups

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// KeeperBlockTimeTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBlockTimeTest struct {
	Inputs       KeeperBlockTimeTestInputs
	TestReporter testreporters.KeeperBlockTimeTestReporter

	keeperRegistry          contracts.KeeperRegistry
	keeperConsumerContracts []contracts.KeeperConsumerPerformance

	env            *environment.Environment
	chainlinkNodes []*client.Chainlink
	chainClient    blockchain.EVMClient
}

// KeeperBlockTimeTestInputs are all the required inputs for a Keeper Block Time Test
type KeeperBlockTimeTestInputs struct {
	BlockchainClient       blockchain.EVMClient              // Client for the test to connect to the blockchain with
	NumberOfContracts      int                               // Number of upkeep contracts
	KeeperRegistrySettings *contracts.KeeperRegistrySettings // Settings of each keeper contract
	Timeout                time.Duration                     // Timeout for the test
	BlockRange             int64                             // How many blocks to run the test for
	BlockInterval          int64                             // Interval of blocks that upkeeps are expected to be performed
	CheckGasToBurn         int64                             // How much gas should be burned on checkUpkeep() calls
	PerformGasToBurn       int64                             // How much gas should be burned on performUpkeep() calls
	ChainlinkNodeFunding   *big.Float                        // Amount of ETH to fund each chainlink node with
}

// NewKeeperBlockTimeTest prepares a new keeper block time test to be run
func NewKeeperBlockTimeTest(inputs KeeperBlockTimeTestInputs) *KeeperBlockTimeTest {
	return &KeeperBlockTimeTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBlockTimeTest) Setup(t *testing.T, env *environment.Environment) {
	startTime := time.Now()
	k.ensureInputValues(t)
	k.env = env
	inputs := k.Inputs
	var err error

	// Connect to networks and prepare for contract deployment
	contractDeployer, err := contracts.NewContractDeployer(k.chainClient)
	require.NoError(t, err, "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	k.chainClient.ParallelTransactions(true)

	// Fund chainlink nodes
	err = actions.FundChainlinkNodes(k.chainlinkNodes, k.chainClient, k.Inputs.ChainlinkNodeFunding)
	require.NoError(t, err, "Funding Chainlink nodes shouldn't fail")
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	err = k.chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for LINK Contract deployment")

	k.keeperRegistry, _, k.keeperConsumerContracts, _ = actions.DeployPerformanceKeeperContracts(
		t,
		ethereum.RegistryVersion_1_1,
		inputs.NumberOfContracts,
		uint32(2500000), //upkeepGasLimit
		linkToken,
		contractDeployer,
		k.chainClient,
		k.Inputs.KeeperRegistrySettings,
		big.NewInt(9e18),
		inputs.BlockRange,
		inputs.BlockInterval,
		inputs.CheckGasToBurn,
		inputs.PerformGasToBurn,
	)
	// Send keeper jobs to registry and chainlink nodes
	actions.CreateKeeperJobs(t, k.chainlinkNodes, k.keeperRegistry, contracts.OCRConfig{})

	log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Block Time Test Setup")
}

// Run runs the keeper block time test
func (k *KeeperBlockTimeTest) Run(t *testing.T) {
	startTime := time.Now()

	for index, keeperConsumer := range k.keeperConsumerContracts {
		k.chainClient.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index),
			contracts.NewKeeperConsumerPerformanceRoundConfirmer(
				keeperConsumer,
				k.Inputs.BlockInterval,
				k.Inputs.BlockRange,
				&k.TestReporter,
			),
		)
	}
	defer func() { // Cleanup the subscriptions
		for index := range k.keeperConsumerContracts {
			k.chainClient.DeleteHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index))
		}
	}()
	err := k.chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for keeper subscriptions")

	for _, chainlinkNode := range k.chainlinkNodes {
		txData, err := chainlinkNode.MustReadTransactionAttempts()
		require.NoError(t, err, "Error retrieving transaction data from chainlink node")
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}

	log.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished Keeper Block Time Test")
}

// Networks returns the networks that the test is running on
func (k *KeeperBlockTimeTest) TearDownVals(t *testing.T) (
	*testing.T,
	*environment.Environment,
	[]*client.Chainlink,
	reportModel.TestReporter,
	blockchain.EVMClient) {
	return t, k.env, k.chainlinkNodes, &k.TestReporter, k.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBlockTimeTest) ensureInputValues(t *testing.T) {
	inputs := k.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Expected a valid blockchain client")
	k.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.NumberOfContracts, 1, "Expecting at least 1 keeper contract")
	if inputs.Timeout == 0 {
		require.Greater(t, inputs.BlockRange, int64(0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.BlockRange <= 0 {
		require.GreaterOrEqual(t, inputs.Timeout, time.Second, "If no `testBlockRange` is provided a `timeout` is required")
	}
	require.NotNil(t, inputs.KeeperRegistrySettings, "You need to set KeeperRegistrySettings")
	require.NotNil(t, k.Inputs.ChainlinkNodeFunding, "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	require.GreaterOrEqual(t, clFunds, 0.0, "Expecting Chainlink node funding to be more than 0 ETH")
	require.Greater(t, inputs.CheckGasToBurn, int64(0), "You need to set an expected amount of gas to burn on checkUpkeep()")
	require.GreaterOrEqual(
		t, int64(inputs.KeeperRegistrySettings.CheckGasLimit), inputs.CheckGasToBurn, "CheckGasLimit should be >= CheckGasToBurn",
	)
	require.Greater(t, inputs.PerformGasToBurn, int64(0), "You need to set an expected amount of gas to burn on performUpkeep()")
}
