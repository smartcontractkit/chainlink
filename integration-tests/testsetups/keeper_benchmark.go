package testsetups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// KeeperBenchmarkTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBenchmarkTest struct {
	Inputs       KeeperBenchmarkTestInputs
	TestReporter testreporters.KeeperBenchmarkTestReporter

	keeperRegistries        []contracts.KeeperRegistry
	keeperRegistrars        []contracts.KeeperRegistrar
	keeperConsumerContracts []contracts.AutomationConsumerBenchmark
	upkeepIDs               [][]*big.Int

	env              *environment.Environment
	namespace        string
	chainlinkNodes   []*client.ChainlinkK8sClient
	chainClient      blockchain.EVMClient
	contractDeployer contracts.ContractDeployer

	linkToken contracts.LinkToken
	ethFeed   contracts.MockETHLINKFeed
	gasFeed   contracts.MockGasFeed
}

type UpkeepConfig struct {
	NumberOfUpkeeps     int   // Number of upkeep contracts
	BlockRange          int64 // How many blocks to run the test for
	BlockInterval       int64 // Interval of blocks that upkeeps are expected to be performed
	CheckGasToBurn      int64 // How much gas should be burned on checkUpkeep() calls
	PerformGasToBurn    int64 // How much gas should be burned on performUpkeep() calls
	UpkeepGasLimit      int64 // Maximum gas that can be consumed by the upkeeps
	FirstEligibleBuffer int64 // How many blocks to add to randomised first eligible block, set to 0 to disable randomised first eligible block
}

type PreDeployedContracts struct {
	RegistryAddress  string
	RegistrarAddress string
	LinkTokenAddress string
	EthFeedAddress   string
	GasFeedAddress   string
}

// KeeperBenchmarkTestInputs are all the required inputs for a Keeper Benchmark Test
type KeeperBenchmarkTestInputs struct {
	BlockchainClient       blockchain.EVMClient              // Client for the test to connect to the blockchain with
	KeeperRegistrySettings *contracts.KeeperRegistrySettings // Settings of each keeper contract
	Upkeeps                *UpkeepConfig
	Contracts              *PreDeployedContracts
	Timeout                time.Duration                    // Timeout for the test
	ChainlinkNodeFunding   *big.Float                       // Amount of ETH to fund each chainlink node with
	UpkeepSLA              int64                            // SLA in number of blocks for an upkeep to be performed once it becomes eligible
	RegistryVersions       []ethereum.KeeperRegistryVersion // Registry version to use
	ForceSingleTxnKey      bool
	BlockTime              time.Duration
	DeltaStage             time.Duration
	DeleteJobsOnEnd        bool
}

// NewKeeperBenchmarkTest prepares a new keeper benchmark test to be run
func NewKeeperBenchmarkTest(inputs KeeperBenchmarkTestInputs) *KeeperBenchmarkTest {
	return &KeeperBenchmarkTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBenchmarkTest) Setup(t *testing.T, env *environment.Environment) {
	l := utils.GetTestLogger(t)
	startTime := time.Now()
	k.TestReporter.Summary.StartTime = startTime.UnixMilli()
	k.ensureInputValues(t)
	k.env = env
	k.namespace = k.env.Cfg.Namespace
	inputs := k.Inputs

	k.keeperRegistries = make([]contracts.KeeperRegistry, len(inputs.RegistryVersions))
	k.keeperRegistrars = make([]contracts.KeeperRegistrar, len(inputs.RegistryVersions))
	k.keeperConsumerContracts = make([]contracts.AutomationConsumerBenchmark, len(inputs.RegistryVersions))
	k.upkeepIDs = make([][]*big.Int, len(inputs.RegistryVersions))

	var err error
	// Connect to networks and prepare for contract deployment
	k.contractDeployer, err = contracts.NewContractDeployer(k.chainClient)
	require.NoError(t, err, "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	k.chainClient.ParallelTransactions(true)

	if len(inputs.RegistryVersions) > 1 && !inputs.ForceSingleTxnKey {
		for nodeIndex, node := range k.chainlinkNodes {
			for registryIndex := 1; registryIndex < len(inputs.RegistryVersions); registryIndex++ {
				l.Debug().Str("URL", node.URL()).Int("NodeIndex", nodeIndex).Int("RegistryIndex", registryIndex).Msg("Create Tx key")
				_, _, err := node.CreateTxKey("evm", k.Inputs.BlockchainClient.GetChainID().String())
				require.NoError(t, err, "Creating transaction key shouldn't fail")
			}
		}
	}

	var ()

	c := inputs.Contracts

	if common.IsHexAddress(c.LinkTokenAddress) {
		k.linkToken, err = k.contractDeployer.LoadLinkToken(common.HexToAddress(c.LinkTokenAddress))
		require.NoError(t, err, "Loading Link Token Contract shouldn't fail")
	} else {
		k.linkToken, err = k.contractDeployer.DeployLinkTokenContract()
		require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Failed waiting for LINK Contract deployment")
	}

	if common.IsHexAddress(c.EthFeedAddress) {
		k.ethFeed, err = k.contractDeployer.LoadETHLINKFeed(common.HexToAddress(c.EthFeedAddress))
		require.NoError(t, err, "Loading ETH-Link feed Contract shouldn't fail")
	} else {
		k.ethFeed, err = k.contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
		require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Failed waiting for ETH-Link feed Contract deployment")
	}

	if common.IsHexAddress(c.GasFeedAddress) {
		k.gasFeed, err = k.contractDeployer.LoadGasFeed(common.HexToAddress(c.GasFeedAddress))
		require.NoError(t, err, "Loading Gas feed Contract shouldn't fail")
	} else {
		k.gasFeed, err = k.contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
		require.NoError(t, err, "Deploying mock gas feed shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Failed waiting for mock gas feed Contract deployment")
	}

	err = k.chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for mock feeds to deploy")

	for index := range inputs.RegistryVersions {
		l.Info().Int("Index", index).Msg("Starting Test Setup")

		k.DeployBenchmarkKeeperContracts(
			t,
			index,
		)
	}

	var keysToFund = inputs.RegistryVersions
	if inputs.ForceSingleTxnKey {
		keysToFund = inputs.RegistryVersions[0:1]
	}

	for index := range keysToFund {
		// Fund chainlink nodes
		nodesToFund := k.chainlinkNodes
		if inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_0 {
			nodesToFund = k.chainlinkNodes[1:]
		}
		err = actions.FundChainlinkNodesAddress(nodesToFund, k.chainClient, k.Inputs.ChainlinkNodeFunding, index)
		require.NoError(t, err, "Funding Chainlink nodes shouldn't fail")
	}

	l.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Benchmark Test Setup")
	err = k.SendSlackNotification(nil)
	if err != nil {
		l.Warn().Msg("Sending test start slack notification failed")
	}
}

// Run runs the keeper benchmark test
func (k *KeeperBenchmarkTest) Run(t *testing.T) {
	l := utils.GetTestLogger(t)
	u := k.Inputs.Upkeeps
	k.TestReporter.Summary.Load.TotalCheckGasPerBlock = int64(u.NumberOfUpkeeps) * u.CheckGasToBurn
	k.TestReporter.Summary.Load.TotalPerformGasPerBlock = int64((float64(u.NumberOfUpkeeps) /
		float64(u.BlockInterval)) * float64(u.PerformGasToBurn))
	k.TestReporter.Summary.Load.AverageExpectedPerformsPerBlock = float64(u.NumberOfUpkeeps) /
		float64(u.BlockInterval)
	k.TestReporter.Summary.TestInputs = map[string]interface{}{
		"NumberOfUpkeeps":     u.NumberOfUpkeeps,
		"BlockCountPerTurn":   k.Inputs.KeeperRegistrySettings.BlockCountPerTurn,
		"CheckGasLimit":       k.Inputs.KeeperRegistrySettings.CheckGasLimit,
		"MaxPerformGas":       k.Inputs.KeeperRegistrySettings.MaxPerformGas,
		"CheckGasToBurn":      u.CheckGasToBurn,
		"PerformGasToBurn":    u.PerformGasToBurn,
		"BlockRange":          u.BlockRange,
		"BlockInterval":       u.BlockInterval,
		"UpkeepSLA":           k.Inputs.UpkeepSLA,
		"FirstEligibleBuffer": u.FirstEligibleBuffer,
		"NumberOfRegistries":  len(k.keeperRegistries),
	}
	inputs := k.Inputs
	startTime := time.Now()

	nodesWithoutBootstrap := k.chainlinkNodes[1:]

	for rIndex := range k.keeperRegistries {
		if k.Inputs.DeltaStage == 0 {
			k.Inputs.DeltaStage = k.Inputs.BlockTime * 5
		}

		var txKeyId = rIndex
		if inputs.ForceSingleTxnKey {
			txKeyId = 0
		}
		ocrConfig, err := actions.BuildAutoOCR2ConfigVarsWithKeyIndex(
			t, nodesWithoutBootstrap, *inputs.KeeperRegistrySettings, k.keeperRegistrars[rIndex].Address(), k.Inputs.DeltaStage, txKeyId,
		)
		require.NoError(t, err, "Building OCR config shouldn't fail")

		// Send keeper jobs to registry and chainlink nodes
		if inputs.RegistryVersions[rIndex] == ethereum.RegistryVersion_2_0 {
			actions.CreateOCRKeeperJobs(t, k.chainlinkNodes, k.keeperRegistries[rIndex].Address(), k.chainClient.GetChainID().Int64(), txKeyId, ethereum.RegistryVersion_2_0)
			err = k.keeperRegistries[rIndex].SetConfig(*inputs.KeeperRegistrySettings, ocrConfig)
			require.NoError(t, err, "Registry config should be be set successfully")
			// Give time for OCR nodes to bootstrap
			time.Sleep(1 * time.Minute)
		} else {
			actions.CreateKeeperJobsWithKeyIndex(t, k.chainlinkNodes, k.keeperRegistries[rIndex], txKeyId, ocrConfig)
		}
		err = k.chainClient.WaitForEvents()
		require.NoError(t, err, "Error waiting for registry setConfig")
	}

	for rIndex := range k.keeperRegistries {
		for index, upkeepID := range k.upkeepIDs[rIndex] {
			k.chainClient.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d %d", rIndex, index),
				contracts.NewKeeperConsumerBenchmarkRoundConfirmer(
					k.keeperConsumerContracts[rIndex],
					k.keeperRegistries[rIndex],
					upkeepID,
					inputs.Upkeeps.BlockRange+inputs.UpkeepSLA,
					inputs.UpkeepSLA,
					&k.TestReporter,
					int64(index),
					inputs.Upkeeps.FirstEligibleBuffer,
				),
			)
		}
	}
	defer func() { // Cleanup the subscriptions
		for rIndex := range k.keeperRegistries {
			for index := range k.upkeepIDs[rIndex] {
				k.chainClient.DeleteHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d %d", rIndex, index))
			}
		}

	}()
	logSubscriptionStop := make(chan bool)
	for rIndex := range k.keeperRegistries {
		k.subscribeToUpkeepPerformedEvent(t, logSubscriptionStop, &k.TestReporter, rIndex)
	}
	err := k.chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for keeper subscriptions")
	close(logSubscriptionStop)

	for _, chainlinkNode := range k.chainlinkNodes {
		txData, err := chainlinkNode.MustReadTransactionAttempts()
		if err != nil {
			l.Error().Err(err).Msg("Error reading transaction attempts from Chainlink Node")
		}
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}

	k.TestReporter.Summary.Config.Chainlink, err = k.env.ResourcesSummary("app=chainlink-0")
	if err != nil {
		l.Error().Err(err).Msg("Error getting resource summary of chainlink node")
	}

	k.TestReporter.Summary.Config.Geth, err = k.env.ResourcesSummary("app=geth")
	if err != nil && k.Inputs.BlockchainClient.NetworkSimulated() {
		l.Error().Err(err).Msg("Error getting resource summary of geth node")
	}

	endTime := time.Now()
	k.TestReporter.Summary.EndTime = endTime.UnixMilli() + (30 * time.Second.Milliseconds())

	for rIndex := range k.keeperRegistries {
		if inputs.DeleteJobsOnEnd {
			// Delete keeper jobs on chainlink nodes
			actions.DeleteKeeperJobsWithId(t, k.chainlinkNodes, rIndex+1)
		}
	}

	l.Info().Str("Run Time", endTime.Sub(startTime).String()).Msg("Finished Keeper Benchmark Test")
}

// subscribeToUpkeepPerformedEvent subscribes to the event log for UpkeepPerformed event and
// counts the number of times it was unsuccessful
func (k *KeeperBenchmarkTest) subscribeToUpkeepPerformedEvent(
	t *testing.T,
	doneChan chan bool,
	metricsReporter *testreporters.KeeperBenchmarkTestReporter,
	rIndex int,
) {
	l := utils.GetTestLogger(t)
	contractABI, err := keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	require.NoError(t, err, "Error getting ABI")
	switch k.Inputs.RegistryVersions[rIndex] {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}

	require.NoError(t, err, "Getting contract abi for registry shouldn't fail")
	query := geth.FilterQuery{
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
				l.Error().Err(err).Msg("Error while subscribing to Keeper Event Logs. Resubscribing...")
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
					l.Info().
						Str("Upkeep ID", parsedLog.Id.String()).
						Bool("Success", parsedLog.Success).
						Str("From", parsedLog.From.String()).
						Str("Registry", k.keeperRegistries[rIndex].Address()).
						Msg("Got successful Upkeep Performed log on Registry")

				} else {
					l.Warn().
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
	string,
	[]*client.ChainlinkK8sClient,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, k.namespace, k.chainlinkNodes, &k.TestReporter, k.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBenchmarkTest) ensureInputValues(t *testing.T) {
	inputs := k.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Need a valid blockchain client to use for the test")
	k.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.Upkeeps.NumberOfUpkeeps, 1, "Expecting at least 1 keeper contracts")
	if inputs.Timeout == 0 {
		require.Greater(t, inputs.Upkeeps.BlockRange, int64(0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.Upkeeps.BlockRange <= 0 {
		require.GreaterOrEqual(t, inputs.Timeout, time.Second, "If no `testBlockRange` is provided a `timeout` is required")
	}
	require.NotNil(t, inputs.KeeperRegistrySettings, "You need to set KeeperRegistrySettings")
	require.NotNil(t, k.Inputs.ChainlinkNodeFunding, "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	require.GreaterOrEqual(t, clFunds, 0.0, "Expecting Chainlink node funding to be more than 0 ETH")
	require.Greater(t, inputs.Upkeeps.CheckGasToBurn, int64(0), "You need to set an expected amount of gas to burn on checkUpkeep()")
	require.GreaterOrEqual(
		t, int64(inputs.KeeperRegistrySettings.CheckGasLimit), inputs.Upkeeps.CheckGasToBurn, "CheckGasLimit should be >= CheckGasToBurn",
	)
	require.Greater(t, inputs.Upkeeps.PerformGasToBurn, int64(0), "You need to set an expected amount of gas to burn on performUpkeep()")
	require.NotNil(t, inputs.UpkeepSLA, "Expected UpkeepSLA to be set")
	require.NotNil(t, inputs.Upkeeps.FirstEligibleBuffer, "You need to set FirstEligibleBuffer")
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

// DeployBenchmarkKeeperContracts deploys a set amount of keeper Benchmark contracts registered to a single registry
func (k *KeeperBenchmarkTest) DeployBenchmarkKeeperContracts(
	t *testing.T,
	index int,
) {
	l := utils.GetTestLogger(t)
	registryVersion := k.Inputs.RegistryVersions[index]
	upkeep := k.Inputs.Upkeeps
	registry := actions.DeployKeeperRegistry(t, k.contractDeployer, k.chainClient,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        k.linkToken.Address(),
			ETHFeedAddr:     k.ethFeed.Address(),
			GasFeedAddr:     k.gasFeed.Address(),
			TranscoderAddr:  actions.ZeroAddress.Hex(),
			RegistrarAddr:   actions.ZeroAddress.Hex(),
			Settings:        *k.Inputs.KeeperRegistrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of AutomationConsumerBenchmark contracts
	err := k.linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(k.Inputs.Upkeeps.NumberOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := actions.DeployKeeperRegistrar(t, registryVersion, k.linkToken, registrarSettings, k.contractDeployer, k.chainClient, registry)
	if registryVersion == ethereum.RegistryVersion_2_0 {
		nodesWithoutBootstrap := k.chainlinkNodes[1:]
		ocrConfig, err := actions.BuildAutoOCR2ConfigVarsWithKeyIndex(
			t, nodesWithoutBootstrap, *k.Inputs.KeeperRegistrySettings, registrar.Address(), k.Inputs.DeltaStage, 0)
		require.NoError(t, err, "OCR config should be built successfully")
		err = registry.SetConfig(*k.Inputs.KeeperRegistrySettings, ocrConfig)
		require.NoError(t, err, "Registry config should be be set successfully")
	}

	consumer := DeployKeeperConsumersBenchmark(t, k.contractDeployer, k.chainClient)

	var upkeepAddresses []string

	checkData := make([][]byte, 0)
	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	require.NoError(t, err)
	var data []byte
	checkDataAbi := abi.Arguments{
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
	}
	for i := 0; i < upkeep.NumberOfUpkeeps; i++ {
		upkeepAddresses = append(upkeepAddresses, consumer.Address())
		// Compute check data
		data, err = checkDataAbi.Pack(
			big.NewInt(int64(i)), big.NewInt(upkeep.BlockInterval), big.NewInt(upkeep.BlockRange),
			big.NewInt(upkeep.CheckGasToBurn), big.NewInt(upkeep.PerformGasToBurn), big.NewInt(upkeep.FirstEligibleBuffer))
		require.NoError(t, err)
		l.Debug().Str("checkData: ", hexutil.Encode(data)).Int("id", i).Msg("checkData computed")
		checkData = append(checkData, data)
	}
	linkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(upkeep.BlockRange/upkeep.BlockInterval))
	gasPrice := big.NewInt(0).Mul(k.Inputs.KeeperRegistrySettings.FallbackGasPrice, big.NewInt(2))
	minLinkBalance := big.NewInt(0).
		Add(big.NewInt(0).
			Mul(big.NewInt(0).
				Div(big.NewInt(0).Mul(gasPrice, big.NewInt(upkeep.UpkeepGasLimit+80000)), k.Inputs.KeeperRegistrySettings.FallbackLinkPrice),
				big.NewInt(1e18+0)),
			big.NewInt(0))

	linkFunds = big.NewInt(0).Add(linkFunds, minLinkBalance)

	upkeepIds := actions.RegisterUpkeepContractsWithCheckData(
		t, k.linkToken, linkFunds, k.chainClient, uint32(upkeep.UpkeepGasLimit), registry,
		registrar, upkeep.NumberOfUpkeeps, upkeepAddresses, checkData,
	)

	k.keeperRegistries[index] = registry
	k.keeperRegistrars[index] = registrar
	k.upkeepIDs[index] = upkeepIds
	k.keeperConsumerContracts[index] = consumer
}

func DeployKeeperConsumersBenchmark(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) contracts.AutomationConsumerBenchmark {
	l := utils.GetTestLogger(t)

	// Deploy consumer
	keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerBenchmark()
	if err != nil {
		l.Error().Err(err).Msg("Deploying AutomationConsumerBenchmark instance %d shouldn't fail")
		keeperConsumerInstance, err = contractDeployer.DeployKeeperConsumerBenchmark()
		require.NoError(t, err, "Error deploying AutomationConsumerBenchmark")
	}
	l.Debug().
		Str("Contract Address", keeperConsumerInstance.Address()).
		Msg("Deployed Keeper Benchmark Contract")

	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper consumer contracts")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerInstance
}
