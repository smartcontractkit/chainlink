package testsetups

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	goeath "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// KeeperBenchmarkTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBenchmarkTest struct {
	Inputs       KeeperBenchmarkTestInputs
	TestReporter testreporters.KeeperBenchmarkTestReporter

	keeperRegistries        []contracts.KeeperRegistry
	keeperRegistrars        []contracts.KeeperRegistrar
	keeperConsumerContracts [][]contracts.KeeperConsumerBenchmark
	upkeepIDs               [][]*big.Int

	env            *environment.Environment
	chainlinkNodes []*client.Chainlink
	chainClient    blockchain.EVMClient
}

// KeeperBenchmarkTestInputs are all the required inputs for a Keeper Benchmark Test
type KeeperBenchmarkTestInputs struct {
	BlockchainClient       blockchain.EVMClient              // Client for the test to connect to the blockchain with
	NumberOfContracts      int                               // Number of upkeep contracts
	KeeperRegistrySettings *contracts.KeeperRegistrySettings // Settings of each keeper contract
	Timeout                time.Duration                     // Timeout for the test
	BlockRange             int64                             // How many blocks to run the test for
	BlockInterval          int64                             // Interval of blocks that upkeeps are expected to be performed
	CheckGasToBurn         int64                             // How much gas should be burned on checkUpkeep() calls
	PerformGasToBurn       int64                             // How much gas should be burned on performUpkeep() calls
	ChainlinkNodeFunding   *big.Float                        // Amount of ETH to fund each chainlink node with
	UpkeepGasLimit         int64                             // Maximum gas that can be consumed by the upkeeps
	UpkeepSLA              int64                             // SLA in number of blocks for an upkeep to be performed once it becomes eligible
	FirstEligibleBuffer    int64                             // How many blocks to add to randomised first eligible block, set to 0 to disable randomised first eligible block
	RegistryVersions       []ethereum.KeeperRegistryVersion  // Registry version to use
	PreDeployedConsumers   []string                          // PreDeployed consumer contracts to re-use in test
	UpkeepResetterAddress  string
	InitialBatchReset      bool
	BlockTime              time.Duration
}

// NewKeeperBenchmarkTest prepares a new keeper benchmark test to be run
func NewKeeperBenchmarkTest(inputs KeeperBenchmarkTestInputs) *KeeperBenchmarkTest {
	return &KeeperBenchmarkTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBenchmarkTest) Setup(t *testing.T, env *environment.Environment) {
	startTime := time.Now()
	k.TestReporter.Summary.StartTime = startTime.UnixMilli()
	k.ensureInputValues(t)
	k.env = env
	inputs := k.Inputs

	k.keeperRegistries = make([]contracts.KeeperRegistry, len(inputs.RegistryVersions))
	k.keeperRegistrars = make([]contracts.KeeperRegistrar, len(inputs.RegistryVersions))
	k.keeperConsumerContracts = make([][]contracts.KeeperConsumerBenchmark, len(inputs.RegistryVersions))
	k.upkeepIDs = make([][]*big.Int, len(inputs.RegistryVersions))

	var err error
	// Connect to networks and prepare for contract deployment
	contractDeployer, err := contracts.NewContractDeployer(k.chainClient)
	require.NoError(t, err, "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	k.chainClient.ParallelTransactions(true)

	if len(inputs.RegistryVersions) > 1 {
		for nodeIndex, node := range k.chainlinkNodes {
			for registryIndex := 1; registryIndex < len(inputs.RegistryVersions); registryIndex++ {
				log.Debug().Str("URL", node.URL()).Int("NodeIndex", nodeIndex).Int("RegistryIndex", registryIndex).Msg("Create Tx key")
				_, _, err := node.CreateTxKey("evm", k.Inputs.BlockchainClient.GetChainID().String())
				require.NoError(t, err, "Creating transaction key shouldn't fail")
			}
		}
	}

	for index := range inputs.RegistryVersions {
		log.Info().Int("Index", index).Msg("Starting Test Setup")

		linkToken, err := contractDeployer.DeployLinkTokenContract()
		require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Failed waiting for LINK Contract deployment")

		k.keeperRegistries[index], k.keeperRegistrars[index], k.keeperConsumerContracts[index], k.upkeepIDs[index] = actions.DeployBenchmarkKeeperContracts(
			t,
			inputs.RegistryVersions[index],
			inputs.NumberOfContracts,
			uint32(inputs.UpkeepGasLimit), //upkeepGasLimit
			linkToken,
			contractDeployer,
			k.chainClient,
			inputs.KeeperRegistrySettings,
			inputs.BlockRange,
			inputs.BlockInterval,
			inputs.CheckGasToBurn,
			inputs.PerformGasToBurn,
			inputs.FirstEligibleBuffer,
			inputs.PreDeployedConsumers,
			inputs.UpkeepResetterAddress,
			k.chainlinkNodes,
			inputs.BlockTime,
		)
	}

	for index := range inputs.RegistryVersions {
		// Fund chainlink nodes
		nodesToFund := k.chainlinkNodes
		if inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_0 {
			nodesToFund = k.chainlinkNodes[1:]
		}
		err = actions.FundChainlinkNodesAddress(nodesToFund, k.chainClient, k.Inputs.ChainlinkNodeFunding, index)
		require.NoError(t, err, "Funding Chainlink nodes shouldn't fail")
	}

	log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Benchmark Test Setup")
	err = k.SendSlackNotification(nil)
	if err != nil {
		log.Warn().Msg("Sending test start slack notification failed")
	}
}

// Run runs the keeper benchmark test
func (k *KeeperBenchmarkTest) Run(t *testing.T) {
	k.TestReporter.Summary.Load.TotalCheckGasPerBlock = int64(k.Inputs.NumberOfContracts) * k.Inputs.CheckGasToBurn
	k.TestReporter.Summary.Load.TotalPerformGasPerBlock = int64((float64(k.Inputs.NumberOfContracts) / float64(k.Inputs.BlockInterval)) * float64(k.Inputs.PerformGasToBurn))
	k.TestReporter.Summary.Load.AverageExpectedPerformsPerBlock = float64(k.Inputs.NumberOfContracts) / float64(k.Inputs.BlockInterval)
	k.TestReporter.Summary.TestInputs = map[string]interface{}{
		"NumberOfContracts":   k.Inputs.NumberOfContracts,
		"BlockCountPerTurn":   k.Inputs.KeeperRegistrySettings.BlockCountPerTurn,
		"CheckGasLimit":       k.Inputs.KeeperRegistrySettings.CheckGasLimit,
		"MaxPerformGas":       k.Inputs.KeeperRegistrySettings.MaxPerformGas,
		"CheckGasToBurn":      k.Inputs.CheckGasToBurn,
		"PerformGasToBurn":    k.Inputs.PerformGasToBurn,
		"BlockRange":          k.Inputs.BlockRange,
		"BlockInterval":       k.Inputs.BlockInterval,
		"UpkeepSLA":           k.Inputs.UpkeepSLA,
		"FirstEligibleBuffer": k.Inputs.FirstEligibleBuffer,
		"NumberOfRegistries":  len(k.keeperRegistries),
	}
	contractDeployer, err := contracts.NewContractDeployer(k.chainClient)
	require.NoError(t, err, "Building a new contract deployer shouldn't fail")
	inputs := k.Inputs
	startTime := time.Now()

	nodesWithoutBootstrap := k.chainlinkNodes[1:]

	for rIndex := range k.keeperRegistries {
		ocrConfig := actions.BuildAutoOCR2ConfigVars(
			t, nodesWithoutBootstrap, *inputs.KeeperRegistrySettings, k.keeperRegistrars[rIndex].Address(), k.Inputs.BlockTime*5,
		)

		// Send keeper jobs to registry and chainlink nodes
		if inputs.RegistryVersions[rIndex] == ethereum.RegistryVersion_2_0 {
			actions.CreateOCRKeeperJobs(t, k.chainlinkNodes, k.keeperRegistries[rIndex].Address(), k.chainClient.GetChainID().Int64(), rIndex)
			err = k.keeperRegistries[rIndex].SetConfig(*inputs.KeeperRegistrySettings, ocrConfig)
			require.NoError(t, err, "Registry config should be be set successfully")
			// Give time for OCR nodes to bootstrap
			time.Sleep(2 * time.Minute)
		} else {
			actions.CreateKeeperJobsWithKeyIndex(t, k.chainlinkNodes, k.keeperRegistries[rIndex], rIndex, ocrConfig)
		}
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Error waiting for registry setConfig")
		// Reset upkeeps so that they become eligible gradually after the test starts
		actions.ResetUpkeeps(t, contractDeployer, k.chainClient, inputs.NumberOfContracts, inputs.BlockRange, inputs.BlockInterval, inputs.CheckGasToBurn,
			inputs.PerformGasToBurn, inputs.FirstEligibleBuffer, k.keeperConsumerContracts[rIndex], inputs.UpkeepResetterAddress)
		for index, keeperConsumer := range k.keeperConsumerContracts[rIndex] {
			k.chainClient.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d %d", rIndex, index),
				contracts.NewKeeperConsumerBenchmarkRoundConfirmer(
					keeperConsumer,
					k.keeperRegistries[rIndex],
					k.upkeepIDs[rIndex][index],
					inputs.BlockRange+inputs.UpkeepSLA,
					inputs.UpkeepSLA,
					&k.TestReporter,
					int64(index),
				),
			)
		}
	}
	defer func() { // Cleanup the subscriptions
		for rIndex := range k.keeperRegistries {
			for index := range k.keeperConsumerContracts[rIndex] {
				k.chainClient.DeleteHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d %d", rIndex, index))
			}
		}

	}()
	logSubscriptionStop := make(chan bool)
	for rIndex := range k.keeperRegistries {
		k.subscribeToUpkeepPerformedEvent(t, logSubscriptionStop, &k.TestReporter, rIndex)
	}
	err = k.chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for keeper subscriptions")
	close(logSubscriptionStop)

	for _, chainlinkNode := range k.chainlinkNodes {
		txData, err := chainlinkNode.MustReadTransactionAttempts()
		if err != nil {
			log.Error().Err(err).Msg("Error reading transaction attempts from Chainlink Node")
		}
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}

	k.TestReporter.Summary.Config.Chainlink, err = k.env.ResourcesSummary("app=chainlink-0")
	if err != nil {
		log.Error().Err(err).Msg("Error getting resource summary of chainlink node")
	}

	k.TestReporter.Summary.Config.Geth, err = k.env.ResourcesSummary("app=geth")
	if err != nil && k.Inputs.BlockchainClient.NetworkSimulated() {
		log.Error().Err(err).Msg("Error getting resource summary of geth node")
	}

	endTime := time.Now()
	k.TestReporter.Summary.EndTime = endTime.UnixMilli() + (30 * time.Second.Milliseconds())

	for rIndex := range k.keeperRegistries {
		// Delete keeper jobs on chainlink nodes
		actions.DeleteKeeperJobsWithId(t, k.chainlinkNodes, rIndex+1)
	}

	log.Info().Str("Run Time", endTime.Sub(startTime).String()).Msg("Finished Keeper Benchmark Test")
}

// subscribeToUpkeepPerformedEvent subscribes to the event log for UpkeepPerformed event and
// counts the number of times it was unsuccessful
func (k *KeeperBenchmarkTest) subscribeToUpkeepPerformedEvent(
	t *testing.T,
	doneChan chan bool,
	metricsReporter *testreporters.KeeperBenchmarkTestReporter,
	rIndex int,
) {
	contractABI, err := ethereum.KeeperRegistry11MetaData.GetAbi()
	require.NoError(t, err, "Error getting ABI")
	switch k.Inputs.RegistryVersions[rIndex] {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		contractABI, err = ethereum.KeeperRegistry11MetaData.GetAbi()
	case ethereum.RegistryVersion_1_2:
		contractABI, err = ethereum.KeeperRegistry12MetaData.GetAbi()
	case ethereum.RegistryVersion_1_3:
		contractABI, err = ethereum.KeeperRegistry13MetaData.GetAbi()
	case ethereum.RegistryVersion_2_0:
		contractABI, err = ethereum.KeeperRegistry20MetaData.GetAbi()
	default:
		contractABI, err = ethereum.KeeperRegistry13MetaData.GetAbi()
	}

	require.NoError(t, err, "Getting contract abi for registry shouldn't fail")
	query := goeath.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(k.keeperRegistries[rIndex].Address())},
	}
	eventLogs := make(chan types.Log)
	sub, err := k.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	require.NoError(t, err, "Subscribing to upkeep performed events log shouldn't fail")
	go func() {
		var numRevertedUpkeeps int64
		for {
			select {
			case err := <-sub.Err():
				log.Error().Err(err).Msg("Error while subscribing to Keeper Event Logs. Resubscribing...")
				sub.Unsubscribe()

				sub, err = k.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				require.NoError(t, err, "Error re-subscribing to event logs")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				require.NoError(t, err, "Getting event details for subscribed log shouldn't fail")
				if eventDetails.Name != "UpkeepPerformed" {
					// Skip non upkeepPerformed Logs
					continue
				}
				parsedLog, err := k.keeperRegistries[rIndex].ParseUpkeepPerformedLog(&vLog)
				require.NoError(t, err, "Parsing upkeep performed log shouldn't fail")

				if parsedLog.Success {
					log.Info().
						Str("Upkeep ID", parsedLog.Id.String()).
						Bool("Success", parsedLog.Success).
						Str("From", parsedLog.From.String()).
						Str("Registry", k.keeperRegistries[rIndex].Address()).
						Msg("Got successful Upkeep Performed log on Registry")

				} else {
					log.Warn().
						Str("Upkeep ID", parsedLog.Id.String()).
						Bool("Success", parsedLog.Success).
						Str("From", parsedLog.From.String()).
						Str("Registry", k.keeperRegistries[rIndex].Address()).
						Msg("Got reverted Upkeep Performed log on Registry")
					numRevertedUpkeeps++
				}
			case <-doneChan:
				metricsReporter.NumRevertedUpkeeps = numRevertedUpkeeps
				return
			}
		}
	}()
}

// TearDownVals returns the networks that the test is running on
func (k *KeeperBenchmarkTest) TearDownVals(t *testing.T) (
	*testing.T,
	*environment.Environment,
	[]*client.Chainlink,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, k.env, k.chainlinkNodes, &k.TestReporter, k.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBenchmarkTest) ensureInputValues(t *testing.T) {
	inputs := k.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Need a valid blockchain client to use for the test")
	k.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.NumberOfContracts, 1, "Expecting at least 1 keeper contracts")
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
	require.NotNil(t, inputs.UpkeepSLA, "Expected UpkeepSLA to be set")
	require.NotNil(t, inputs.FirstEligibleBuffer, "You need to set FirstEligibleBuffer")
	require.NotNil(t, inputs.RegistryVersions[0], "You need to set RegistryVersion")
	require.NotNil(t, inputs.BlockTime, "You need to set BlockTime")
}

func (k *KeeperBenchmarkTest) SendSlackNotification(slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(reportModel.SlackAPIKey)
	}

	headerText := ":white_check_mark: Automation Benchmark Test STARTED :white_check_mark:"
	formattedDashboardUrl := fmt.Sprintf("%s&from=%d&to=%s&var-namespace=%s&var-cl_node=chainlink-0-0", testreporters.DashboardUrl, k.TestReporter.Summary.StartTime, "now", k.env.Cfg.Namespace)
	log.Info().Str("Dashboard", formattedDashboardUrl).Msg("Dashboard URL")

	notificationBlocks := []slack.Block{}
	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", k.env.Cfg.Namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("<%s|Test Dashboard> \nNotifying <@%s>",
			formattedDashboardUrl, reportModel.SlackUserID), false, true), nil, nil))

	ts, err := reportModel.SendSlackMessage(slackClient, slack.MsgOptionBlocks(notificationBlocks...))
	log.Debug().Str("ts", ts).Msg("Sent Slack Message")
	return err
}
