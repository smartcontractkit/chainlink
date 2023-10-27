package testsetups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
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

	t             *testing.T
	log           zerolog.Logger
	startingBlock *big.Int

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

// UpkeepConfig dictates details of how the test's upkeep contracts should be called and configured
type UpkeepConfig struct {
	NumberOfUpkeeps     int   // Number of upkeep contracts
	BlockRange          int64 // How many blocks to run the test for
	BlockInterval       int64 // Interval of blocks that upkeeps are expected to be performed
	CheckGasToBurn      int64 // How much gas should be burned on checkUpkeep() calls
	PerformGasToBurn    int64 // How much gas should be burned on performUpkeep() calls
	UpkeepGasLimit      int64 // Maximum gas that can be consumed by the upkeeps
	FirstEligibleBuffer int64 // How many blocks to add to randomised first eligible block, set to 0 to disable randomised first eligible block
}

// PreDeployedContracts are contracts that are already deployed on a (usually) live testnet chain, so re-deployment
// in unnecessary
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
func NewKeeperBenchmarkTest(t *testing.T, inputs KeeperBenchmarkTestInputs) *KeeperBenchmarkTest {
	return &KeeperBenchmarkTest{
		Inputs: inputs,
		t:      t,
		log:    logging.GetTestLogger(t),
	}
}

// Setup prepares contracts for the test
func (k *KeeperBenchmarkTest) Setup(env *environment.Environment) {
	startTime := time.Now()
	k.TestReporter.Summary.StartTime = startTime.UnixMilli()
	k.ensureInputValues()
	k.env = env
	k.namespace = k.env.Cfg.Namespace
	inputs := k.Inputs

	k.keeperRegistries = make([]contracts.KeeperRegistry, len(inputs.RegistryVersions))
	k.keeperRegistrars = make([]contracts.KeeperRegistrar, len(inputs.RegistryVersions))
	k.keeperConsumerContracts = make([]contracts.AutomationConsumerBenchmark, len(inputs.RegistryVersions))
	k.upkeepIDs = make([][]*big.Int, len(inputs.RegistryVersions))
	k.log.Debug().Interface("TestInputs", inputs).Msg("Setting up benchmark test")

	var err error
	// Connect to networks and prepare for contract deployment
	k.contractDeployer, err = contracts.NewContractDeployer(k.chainClient, k.log)
	require.NoError(k.t, err, "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	require.NoError(k.t, err, "Connecting to chainlink nodes shouldn't fail")
	k.chainClient.ParallelTransactions(true)

	if len(inputs.RegistryVersions) > 1 && !inputs.ForceSingleTxnKey {
		for nodeIndex, node := range k.chainlinkNodes {
			for registryIndex := 1; registryIndex < len(inputs.RegistryVersions); registryIndex++ {
				k.log.Debug().Str("URL", node.URL()).Int("NodeIndex", nodeIndex).Int("RegistryIndex", registryIndex).Msg("Create Tx key")
				_, _, err := node.CreateTxKey("evm", k.Inputs.BlockchainClient.GetChainID().String())
				require.NoError(k.t, err, "Creating transaction key shouldn't fail")
			}
		}
	}

	c := inputs.Contracts

	if common.IsHexAddress(c.LinkTokenAddress) {
		k.linkToken, err = k.contractDeployer.LoadLinkToken(common.HexToAddress(c.LinkTokenAddress))
		require.NoError(k.t, err, "Loading Link Token Contract shouldn't fail")
	} else {
		k.linkToken, err = k.contractDeployer.DeployLinkTokenContract()
		require.NoError(k.t, err, "Deploying Link Token Contract shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(k.t, err, "Failed waiting for LINK Contract deployment")
	}

	if common.IsHexAddress(c.EthFeedAddress) {
		k.ethFeed, err = k.contractDeployer.LoadETHLINKFeed(common.HexToAddress(c.EthFeedAddress))
		require.NoError(k.t, err, "Loading ETH-Link feed Contract shouldn't fail")
	} else {
		k.ethFeed, err = k.contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
		require.NoError(k.t, err, "Deploying mock ETH-Link feed shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(k.t, err, "Failed waiting for ETH-Link feed Contract deployment")
	}

	if common.IsHexAddress(c.GasFeedAddress) {
		k.gasFeed, err = k.contractDeployer.LoadGasFeed(common.HexToAddress(c.GasFeedAddress))
		require.NoError(k.t, err, "Loading Gas feed Contract shouldn't fail")
	} else {
		k.gasFeed, err = k.contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
		require.NoError(k.t, err, "Deploying mock gas feed shouldn't fail")
		err = k.chainClient.WaitForEvents()
		require.NoError(k.t, err, "Failed waiting for mock gas feed Contract deployment")
	}

	err = k.chainClient.WaitForEvents()
	require.NoError(k.t, err, "Failed waiting for mock feeds to deploy")

	for index := range inputs.RegistryVersions {
		k.log.Info().Int("Index", index).Msg("Starting Test Setup")

		k.DeployBenchmarkKeeperContracts(index)
	}

	var keysToFund = inputs.RegistryVersions
	if inputs.ForceSingleTxnKey {
		keysToFund = inputs.RegistryVersions[0:1]
	}

	for index := range keysToFund {
		// Fund chainlink nodes
		nodesToFund := k.chainlinkNodes
		if inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_0 || inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_1 {
			nodesToFund = k.chainlinkNodes[1:]
		}
		err = actions.FundChainlinkNodesAddress(nodesToFund, k.chainClient, k.Inputs.ChainlinkNodeFunding, index)
		require.NoError(k.t, err, "Funding Chainlink nodes shouldn't fail")
	}

	k.log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Benchmark Test Setup")
	err = k.SendSlackNotification(nil)
	if err != nil {
		k.log.Warn().Msg("Sending test start slack notification failed")
	}
}

// Run runs the keeper benchmark test
func (k *KeeperBenchmarkTest) Run() {
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
	startingBlock, err := k.chainClient.LatestBlockNumber(context.Background())
	require.NoError(k.t, err, "Error getting latest block number")
	k.startingBlock = big.NewInt(0).SetUint64(startingBlock)
	startTime := time.Now()

	nodesWithoutBootstrap := k.chainlinkNodes[1:]

	for rIndex := range k.keeperRegistries {

		var txKeyId = rIndex
		if inputs.ForceSingleTxnKey {
			txKeyId = 0
		}
		ocrConfig, err := actions.BuildAutoOCR2ConfigVarsWithKeyIndex(
			k.t, nodesWithoutBootstrap, *inputs.KeeperRegistrySettings, k.keeperRegistrars[rIndex].Address(), k.Inputs.DeltaStage, txKeyId, common.Address{},
		)
		require.NoError(k.t, err, "Building OCR config shouldn't fail")

		// Send keeper jobs to registry and chainlink nodes
		if inputs.RegistryVersions[rIndex] == ethereum.RegistryVersion_2_0 || inputs.RegistryVersions[rIndex] == ethereum.RegistryVersion_2_1 {
			actions.CreateOCRKeeperJobs(k.t, k.chainlinkNodes, k.keeperRegistries[rIndex].Address(), k.chainClient.GetChainID().Int64(), txKeyId, inputs.RegistryVersions[rIndex])
			err = k.keeperRegistries[rIndex].SetConfig(*inputs.KeeperRegistrySettings, ocrConfig)
			require.NoError(k.t, err, "Registry config should be be set successfully")
			// Give time for OCR nodes to bootstrap
			time.Sleep(1 * time.Minute)
		} else {
			actions.CreateKeeperJobsWithKeyIndex(k.t, k.chainlinkNodes, k.keeperRegistries[rIndex], txKeyId, ocrConfig, k.chainClient.GetChainID().String())
		}
		err = k.chainClient.WaitForEvents()
		require.NoError(k.t, err, "Error waiting for registry setConfig")
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
					k.log,
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

	// Main test loop
	k.observeUpkeepEvents()
	err = k.chainClient.WaitForEvents()
	require.NoError(k.t, err, "Error waiting for keeper subscriptions")

	// Collect logs for each registry to calculate test metrics
	registryLogs := make([][]types.Log, len(k.keeperRegistries))
	for rIndex := range k.keeperRegistries {
		var (
			logs        []types.Log
			timeout     = 5 * time.Second
			addr        = k.keeperRegistries[rIndex].Address()
			filterQuery = geth.FilterQuery{
				Addresses: []common.Address{common.HexToAddress(addr)},
				FromBlock: k.startingBlock,
			}
			err = fmt.Errorf("initial error") // to ensure our for loop runs at least once
		)
		for err != nil { // This RPC call can possibly time out or otherwise die. Failure is not an option, keep retrying to get our stats.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			logs, err = k.chainClient.FilterLogs(ctx, filterQuery)
			cancel()
			if err != nil {
				k.log.Error().Err(err).
					Interface("Filter Query", filterQuery).
					Str("Timeout", timeout.String()).
					Msg("Error getting logs from chain, trying again")
			} else {
				k.log.Info().Int("Log Count", len(logs)).Str("Registry Address", addr).Msg("Collected logs")
			}
		}
		registryLogs[rIndex] = logs
	}

	// Count reverts and stale upkeeps
	for rIndex := range k.keeperRegistries {
		contractABI := k.contractABI(rIndex)
		for _, l := range registryLogs[rIndex] {
			log := l
			eventDetails, err := contractABI.EventByID(log.Topics[0])
			if err != nil {
				k.log.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error getting event details for log, report data inaccurate")
				break
			}
			if eventDetails.Name == "UpkeepPerformed" {
				parsedLog, err := k.keeperRegistries[rIndex].ParseUpkeepPerformedLog(&log)
				if err != nil {
					k.log.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error parsing upkeep performed log, report data inaccurate")
					break
				}
				if !parsedLog.Success {
					k.TestReporter.NumRevertedUpkeeps++
				}
			} else if eventDetails.Name == "StaleUpkeepReport" {
				k.TestReporter.NumStaleUpkeepReports++
			}
		}
	}

	for _, chainlinkNode := range k.chainlinkNodes {
		txData, err := chainlinkNode.MustReadTransactionAttempts()
		if err != nil {
			k.log.Error().Err(err).Msg("Error reading transaction attempts from Chainlink Node")
		}
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}

	k.TestReporter.Summary.Config.Chainlink, err = k.env.ResourcesSummary("app=chainlink-0")
	if err != nil {
		k.log.Error().Err(err).Msg("Error getting resource summary of chainlink node")
	}

	k.TestReporter.Summary.Config.Geth, err = k.env.ResourcesSummary("app=geth")
	if err != nil && k.Inputs.BlockchainClient.NetworkSimulated() {
		k.log.Error().Err(err).Msg("Error getting resource summary of geth node")
	}

	endTime := time.Now()
	k.TestReporter.Summary.EndTime = endTime.UnixMilli() + (30 * time.Second.Milliseconds())

	for rIndex := range k.keeperRegistries {
		if inputs.DeleteJobsOnEnd {
			// Delete keeper jobs on chainlink nodes
			actions.DeleteKeeperJobsWithId(k.t, k.chainlinkNodes, rIndex+1)
		}
	}

	k.log.Info().Str("Run Time", endTime.Sub(startTime).String()).Msg("Finished Keeper Benchmark Test")
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

// *********************
// ****** Helpers ******
// *********************

// observeUpkeepEvents subscribes to Upkeep events on deployed registries and logs them
// WARNING: This should only be used for observation and logging. This isn't a reliable way to build a final report
// due to how fragile subscriptions can be
func (k *KeeperBenchmarkTest) observeUpkeepEvents() {
	eventLogs := make(chan types.Log)
	registryAddresses := make([]common.Address, len(k.keeperRegistries))
	addressIndexMap := map[common.Address]int{}
	for index, registry := range k.keeperRegistries {
		registryAddresses[index] = common.HexToAddress(registry.Address())
		addressIndexMap[registryAddresses[index]] = index
	}
	filterQuery := geth.FilterQuery{
		Addresses: registryAddresses,
		FromBlock: k.startingBlock,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	sub, err := k.chainClient.SubscribeFilterLogs(ctx, filterQuery, eventLogs)
	cancel()
	require.NoError(k.t, err, "Subscribing to upkeep performed events log shouldn't fail")

	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-interruption:
				k.log.Warn().Msg("Received interrupt signal, test container restarting. Dashboard view will be inaccurate.")
			case err := <-sub.Err():
				backoff := time.Second
				for err != nil { // Keep retrying until we get a successful subscription
					k.log.Error().
						Err(err).
						Interface("Query", filterQuery).
						Str("Backoff", backoff.String()).
						Msg("Error while subscribing to Keeper Event Logs. Resubscribing...")

					ctx, cancel := context.WithTimeout(context.Background(), backoff)
					sub, err = k.chainClient.SubscribeFilterLogs(ctx, filterQuery, eventLogs)
					cancel()
					if err != nil {
						time.Sleep(backoff)
						backoff = time.Duration(math.Min(float64(backoff)*2, float64(30*time.Second)))
					}
				}
				log.Info().Msg("Resubscribed to Keeper Event Logs")
			case vLog := <-eventLogs:
				rIndex, ok := addressIndexMap[vLog.Address]
				if !ok {
					k.log.Error().Str("Address", vLog.Address.Hex()).Msg("Received log from unknown registry")
					continue
				}
				contractABI := k.contractABI(rIndex)
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				require.NoError(k.t, err, "Getting event details for subscribed log shouldn't fail")
				if eventDetails.Name != "UpkeepPerformed" && eventDetails.Name != "StaleUpkeepReport" {
					// Skip non upkeepPerformed Logs
					continue
				}
				if vLog.Removed {
					k.log.Warn().
						Str("Name", eventDetails.Name).
						Str("Registry", k.keeperRegistries[rIndex].Address()).
						Msg("Got removed log")
				}
				if eventDetails.Name == "UpkeepPerformed" {
					parsedLog, err := k.keeperRegistries[rIndex].ParseUpkeepPerformedLog(&vLog)
					require.NoError(k.t, err, "Parsing upkeep performed log shouldn't fail")

					if parsedLog.Success {
						k.log.Info().
							Str("Upkeep ID", parsedLog.Id.String()).
							Bool("Success", parsedLog.Success).
							Str("From", parsedLog.From.String()).
							Str("Registry", k.keeperRegistries[rIndex].Address()).
							Msg("Got successful Upkeep Performed log on Registry")
					} else {
						k.log.Warn().
							Str("Upkeep ID", parsedLog.Id.String()).
							Bool("Success", parsedLog.Success).
							Str("From", parsedLog.From.String()).
							Str("Registry", k.keeperRegistries[rIndex].Address()).
							Msg("Got reverted Upkeep Performed log on Registry")
					}
				} else if eventDetails.Name == "StaleUpkeepReport" {
					parsedLog, err := k.keeperRegistries[rIndex].ParseStaleUpkeepReportLog(&vLog)
					require.NoError(k.t, err, "Parsing stale upkeep report log shouldn't fail")
					k.log.Warn().
						Str("Upkeep ID", parsedLog.Id.String()).
						Str("Registry", k.keeperRegistries[rIndex].Address()).
						Msg("Got stale Upkeep report log on Registry")
				}
			case <-k.chainClient.ConnectionIssue():
				k.log.Warn().Msg("RPC connection issue detected.")
			case <-k.chainClient.ConnectionRestored():
				k.log.Info().Msg("RPC connection restored.")
			}
		}
	}()
}

// contractABI returns the ABI of the proper keeper registry contract
func (k *KeeperBenchmarkTest) contractABI(rIndex int) *abi.ABI {
	var (
		contractABI *abi.ABI
		err         error
	)
	switch k.Inputs.RegistryVersions[rIndex] {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	case ethereum.RegistryVersion_2_1:
		contractABI, err = iregistry21.IKeeperRegistryMasterMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}
	require.NoError(k.t, err, "Getting contract ABI shouldn't fail")
	return contractABI
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBenchmarkTest) ensureInputValues() {
	inputs := k.Inputs
	require.NotNil(k.t, inputs.BlockchainClient, "Need a valid blockchain client to use for the test")
	k.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(k.t, inputs.Upkeeps.NumberOfUpkeeps, 1, "Expecting at least 1 keeper contracts")
	if inputs.Timeout == 0 {
		require.Greater(k.t, inputs.Upkeeps.BlockRange, int64(0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.Upkeeps.BlockRange <= 0 {
		require.GreaterOrEqual(k.t, inputs.Timeout, time.Second, "If no `testBlockRange` is provided a `timeout` is required")
	}
	require.NotNil(k.t, inputs.KeeperRegistrySettings, "You need to set KeeperRegistrySettings")
	require.NotNil(k.t, k.Inputs.ChainlinkNodeFunding, "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	require.GreaterOrEqual(k.t, clFunds, 0.0, "Expecting Chainlink node funding to be more than 0 ETH")
	require.Greater(k.t, inputs.Upkeeps.CheckGasToBurn, int64(0), "You need to set an expected amount of gas to burn on checkUpkeep()")
	require.GreaterOrEqual(
		k.t, int64(inputs.KeeperRegistrySettings.CheckGasLimit), inputs.Upkeeps.CheckGasToBurn, "CheckGasLimit should be >= CheckGasToBurn",
	)
	require.Greater(k.t, inputs.Upkeeps.PerformGasToBurn, int64(0), "You need to set an expected amount of gas to burn on performUpkeep()")
	require.NotNil(k.t, inputs.UpkeepSLA, "Expected UpkeepSLA to be set")
	require.NotNil(k.t, inputs.Upkeeps.FirstEligibleBuffer, "You need to set FirstEligibleBuffer")
	require.NotNil(k.t, inputs.RegistryVersions[0], "You need to set RegistryVersion")
	require.NotNil(k.t, inputs.BlockTime, "You need to set BlockTime")

	if k.Inputs.DeltaStage == 0 {
		k.Inputs.DeltaStage = k.Inputs.BlockTime * 5
	}
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
func (k *KeeperBenchmarkTest) DeployBenchmarkKeeperContracts(index int) {
	registryVersion := k.Inputs.RegistryVersions[index]
	k.Inputs.KeeperRegistrySettings.RegistryVersion = registryVersion
	upkeep := k.Inputs.Upkeeps
	var (
		registry  contracts.KeeperRegistry
		registrar contracts.KeeperRegistrar
	)

	// Contract deployment is different for legacy keepers and OCR automation
	if registryVersion <= ethereum.RegistryVersion_1_3 { // Legacy keeper - v1.X
		registry = actions.DeployKeeperRegistry(k.t, k.contractDeployer, k.chainClient,
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
		require.NoError(k.t, err, "Funding keeper registry contract shouldn't fail")

		registrarSettings := contracts.KeeperRegistrarSettings{
			AutoApproveConfigType: 2,
			AutoApproveMaxAllowed: math.MaxUint16,
			RegistryAddr:          registry.Address(),
			MinLinkJuels:          big.NewInt(0),
		}
		registrar = actions.DeployKeeperRegistrar(k.t, registryVersion, k.linkToken, registrarSettings, k.contractDeployer, k.chainClient, registry)
	} else { // OCR automation - v2.X
		registry, registrar = actions.DeployAutoOCRRegistryAndRegistrar(
			k.t, registryVersion, *k.Inputs.KeeperRegistrySettings, k.linkToken, k.contractDeployer, k.chainClient,
		)

		// Fund the registry with LINK
		err := k.linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(k.Inputs.Upkeeps.NumberOfUpkeeps))))
		require.NoError(k.t, err, "Funding keeper registry contract shouldn't fail")
		ocrConfig, err := actions.BuildAutoOCR2ConfigVars(k.t, k.chainlinkNodes[1:], *k.Inputs.KeeperRegistrySettings, registrar.Address(), k.Inputs.DeltaStage)
		k.log.Debug().Interface("KeeperRegistrySettings", *k.Inputs.KeeperRegistrySettings).Interface("OCRConfig", ocrConfig).Msg("Config")
		require.NoError(k.t, err, "Error building OCR config vars")
		err = registry.SetConfig(*k.Inputs.KeeperRegistrySettings, ocrConfig)
		require.NoError(k.t, err, "Registry config should be be set successfully")

	}

	consumer := k.DeployKeeperConsumersBenchmark()

	var upkeepAddresses []string

	checkData := make([][]byte, 0)
	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	require.NoError(k.t, err)
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
		require.NoError(k.t, err)
		k.log.Debug().Str("checkData: ", hexutil.Encode(data)).Int("id", i).Msg("checkData computed")
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

	upkeepIds := actions.RegisterUpkeepContractsWithCheckData(k.t, k.linkToken, linkFunds, k.chainClient, uint32(upkeep.UpkeepGasLimit), registry, registrar, upkeep.NumberOfUpkeeps, upkeepAddresses, checkData, false, false)

	k.keeperRegistries[index] = registry
	k.keeperRegistrars[index] = registrar
	k.upkeepIDs[index] = upkeepIds
	k.keeperConsumerContracts[index] = consumer
}

func (k *KeeperBenchmarkTest) DeployKeeperConsumersBenchmark() contracts.AutomationConsumerBenchmark {
	// Deploy consumer
	keeperConsumerInstance, err := k.contractDeployer.DeployKeeperConsumerBenchmark()
	if err != nil {
		k.log.Error().Err(err).Msg("Deploying AutomationConsumerBenchmark instance %d shouldn't fail")
		keeperConsumerInstance, err = k.contractDeployer.DeployKeeperConsumerBenchmark()
		require.NoError(k.t, err, "Error deploying AutomationConsumerBenchmark")
	}
	k.log.Debug().
		Str("Contract Address", keeperConsumerInstance.Address()).
		Msg("Deployed Keeper Benchmark Contract")

	err = k.chainClient.WaitForEvents()
	require.NoError(k.t, err, "Failed waiting for to deploy all keeper consumer contracts")
	k.log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerInstance
}
