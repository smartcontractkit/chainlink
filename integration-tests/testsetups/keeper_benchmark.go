package testsetups

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/slack-go/slack"

	goeath "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-env/environment"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
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
}

// NewKeeperBenchmarkTest prepares a new keeper benchmark test to be run
func NewKeeperBenchmarkTest(inputs KeeperBenchmarkTestInputs) *KeeperBenchmarkTest {
	return &KeeperBenchmarkTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBenchmarkTest) Setup(env *environment.Environment) {
	startTime := time.Now()
	k.TestReporter.Summary.StartTime = startTime.UnixMilli()
	k.ensureInputValues()
	k.env = env
	inputs := k.Inputs

	k.keeperRegistries = make([]contracts.KeeperRegistry, len(inputs.RegistryVersions))
	k.keeperConsumerContracts = make([][]contracts.KeeperConsumerBenchmark, len(inputs.RegistryVersions))
	k.upkeepIDs = make([][]*big.Int, len(inputs.RegistryVersions))

	var err error
	// Connect to networks and prepare for contract deployment
	contractDeployer, err := contracts.NewContractDeployer(k.chainClient)
	Expect(err).ShouldNot(HaveOccurred(), "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	k.chainClient.ParallelTransactions(true)

	if len(inputs.RegistryVersions) > 1 {
		for nodeIndex, node := range k.chainlinkNodes {
			for registryIndex := 1; registryIndex < len(inputs.RegistryVersions); registryIndex++ {
				log.Debug().Str("URL", node.URL()).Int("NodeIndex", nodeIndex).Int("RegistryIndex", registryIndex).Msg("Create Tx key")
				_, _, err := node.CreateTxKey("evm", k.Inputs.BlockchainClient.GetChainID().String())
				Expect(err).ShouldNot(HaveOccurred(), "Creating transaction key shouldn't fail")
			}
		}
	}

	for index := range inputs.RegistryVersions {
		log.Info().Int("Index", index).Msg("Starting Test Setup")

		linkToken, err := contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		err = k.chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for LINK Contract deployment")

		k.keeperRegistries[index], k.keeperConsumerContracts[index], k.upkeepIDs[index] = actions.DeployBenchmarkKeeperContracts(
			inputs.RegistryVersions[index],
			inputs.NumberOfContracts,
			uint32(inputs.UpkeepGasLimit), //upkeepGasLimit
			linkToken,
			contractDeployer,
			k.chainClient,
			k.Inputs.KeeperRegistrySettings,
			inputs.BlockRange,
			inputs.BlockInterval,
			inputs.CheckGasToBurn,
			inputs.PerformGasToBurn,
			inputs.FirstEligibleBuffer,
			inputs.PreDeployedConsumers,
			inputs.UpkeepResetterAddress,
		)
	}

	for index := range inputs.RegistryVersions {
		// Fund chainlink nodes
		err = actions.FundChainlinkNodesAddress(k.chainlinkNodes, k.chainClient, k.Inputs.ChainlinkNodeFunding, index)
		Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
	}

	log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Benchmark Test Setup")
	err = k.SendSlackNotification(nil)
	if err != nil {
		log.Warn().Msg("Sending test start slack notification failed")
	}
}

// Run runs the keeper benchmark test
func (k *KeeperBenchmarkTest) Run() {
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
	Expect(err).ShouldNot(HaveOccurred(), "Building a new contract deployer shouldn't fail")
	inputs := k.Inputs
	startTime := time.Now()

	rampUpBlocks := int64(k.Inputs.NumberOfContracts) / int64(k.TestReporter.Summary.Load.AverageExpectedPerformsPerBlock*2)

	for rIndex := range k.keeperRegistries {
		// Send keeper jobs to registry and chainlink nodes
		actions.CreateKeeperJobsWithKeyIndex(k.chainlinkNodes, k.keeperRegistries[rIndex], rIndex)
		// Reset upkeeps so that they become eligible gradually after the test starts
		actions.ResetUpkeeps(contractDeployer, k.chainClient, inputs.NumberOfContracts, inputs.BlockRange, inputs.BlockInterval, inputs.CheckGasToBurn,
			inputs.PerformGasToBurn, inputs.FirstEligibleBuffer, inputs.PreDeployedConsumers, inputs.UpkeepResetterAddress)
		for index, keeperConsumer := range k.keeperConsumerContracts[rIndex] {
			k.chainClient.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d %d", rIndex, index),
				contracts.NewKeeperConsumerBenchmarkRoundConfirmer(
					keeperConsumer,
					k.keeperRegistries[rIndex],
					k.upkeepIDs[rIndex][index],
					inputs.BlockRange+inputs.UpkeepSLA,
					rampUpBlocks,
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
		k.subscribeToUpkeepPerformedEvent(logSubscriptionStop, &k.TestReporter, rIndex)
	}
	err = k.chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for keeper subscriptions")
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
		panic(err)
	}

	k.TestReporter.Summary.Config.Geth, err = k.env.ResourcesSummary("app=geth")
	if err != nil && k.Inputs.BlockchainClient.NetworkSimulated() {
		panic(err)
	}

	endTime := time.Now()
	k.TestReporter.Summary.EndTime = endTime.UnixMilli() + (30 * time.Second.Milliseconds())

	for rIndex := range k.keeperRegistries {
		// Delete keeper jobs on chainlink nodes
		actions.DeleteKeeperJobsWithId(k.chainlinkNodes, rIndex+1)
	}

	log.Info().Str("Run Time", endTime.Sub(startTime).String()).Msg("Finished Keeper Benchmark Test")
}

// subscribeToUpkeepPerformedEvent subscribes to the event log for UpkeepPerformed event and
// counts the number of times it was unsuccessful
func (k *KeeperBenchmarkTest) subscribeToUpkeepPerformedEvent(doneChan chan bool, metricsReporter *testreporters.KeeperBenchmarkTestReporter, rIndex int) {
	contractABI, err := ethereum.KeeperRegistry11MetaData.GetAbi()
	Expect(err).ShouldNot(HaveOccurred(), "Error getting ABI")
	switch k.Inputs.RegistryVersions[rIndex] {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		contractABI, err = ethereum.KeeperRegistry11MetaData.GetAbi()
	case ethereum.RegistryVersion_1_2:
		contractABI, err = ethereum.KeeperRegistry12MetaData.GetAbi()
	case ethereum.RegistryVersion_1_3:
		contractABI, err = ethereum.KeeperRegistry13MetaData.GetAbi()
	default:
		contractABI, err = ethereum.KeeperRegistry13MetaData.GetAbi()
	}

	Expect(err).ShouldNot(HaveOccurred(), "Getting contract abi for registry shouldn't fail")
	query := goeath.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(k.keeperRegistries[rIndex].Address())},
	}
	eventLogs := make(chan types.Log)
	sub, err := k.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	Expect(err).ShouldNot(HaveOccurred(), "Subscribing to upkeep performed events log shouldn't fail")
	go func() {
		defer ginkgo.GinkgoRecover()
		var numRevertedUpkeeps int64
		for {
			select {
			case err := <-sub.Err():
				log.Error().Err(err).Msg("Error while subscribing to Keeper Event Logs. Resubscribing...")
				sub.Unsubscribe()

				sub, err = k.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				Expect(err).ShouldNot(HaveOccurred(), "Error re-subscribing to event logs")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				Expect(err).ShouldNot(HaveOccurred(), "Getting event details for subscribed log shouldn't fail")
				if eventDetails.Name != "UpkeepPerformed" {
					// Skip non upkeepPerformed Logs
					continue
				}
				parsedLog, err := k.keeperRegistries[rIndex].ParseUpkeepPerformedLog(&vLog)
				Expect(err).ShouldNot(HaveOccurred(), "Parsing upkeep performed log shouldn't fail")

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
func (k *KeeperBenchmarkTest) TearDownVals() (*environment.Environment, []*client.Chainlink, reportModel.TestReporter, blockchain.EVMClient) {
	return k.env, k.chainlinkNodes, &k.TestReporter, k.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBenchmarkTest) ensureInputValues() {
	inputs := k.Inputs
	Expect(inputs.BlockchainClient).ShouldNot(BeNil(), "Need a valid blockchain client to use for the test")
	k.chainClient = inputs.BlockchainClient
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 keeper contracts")
	if inputs.Timeout == 0 {
		Expect(inputs.BlockRange).Should(BeNumerically(">", 0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.BlockRange <= 0 {
		Expect(inputs.Timeout).Should(BeNumerically(">=", 1), "If no `testBlockRange` is provided a `timeout` is required")
	}
	Expect(inputs.KeeperRegistrySettings).ShouldNot(BeNil(), "You need to set KeeperRegistrySettings")
	Expect(k.Inputs.ChainlinkNodeFunding).ShouldNot(BeNil(), "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	Expect(clFunds).Should(BeNumerically(">=", 0), "Expecting Chainlink node funding to be more than 0 ETH")
	Expect(inputs.CheckGasToBurn).Should(BeNumerically(">", 0), "You need to set an expected amount of gas to burn on checkUpkeep()")
	Expect(inputs.KeeperRegistrySettings.CheckGasLimit).Should(BeNumerically(">=", inputs.CheckGasToBurn),
		"CheckGasLimit should be >= CheckGasToBurn")
	Expect(inputs.PerformGasToBurn).Should(BeNumerically(">", 0), "You need to set an expected amount of gas to burn on performUpkeep()")
	Expect(inputs.UpkeepSLA).ShouldNot(BeNil(), "You need to set UpkeepSLA")
	Expect(inputs.FirstEligibleBuffer).ShouldNot(BeNil(), "You need to set FirstEligibleBuffer")
	Expect(inputs.RegistryVersions[0]).ShouldNot(BeNil(), "You need to set RegistryVersion")
	if len(inputs.UpkeepResetterAddress) == 0 {
		k.Inputs.UpkeepResetterAddress = "0x"
	}
}

func (k *KeeperBenchmarkTest) SendSlackNotification(slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(reportModel.SlackAPIKey)
	}

	headerText := ":white_check_mark: Keeper Benchmark Test Started :white_check_mark:"
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
