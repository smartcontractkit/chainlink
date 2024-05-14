package testsetups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	keepertestconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/keeper"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	tt "github.com/smartcontractkit/chainlink/integration-tests/types"
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

	env            *environment.Environment
	namespace      string
	chainlinkNodes []*client.ChainlinkK8sClient
	chainClient    *seth.Client
	testConfig     tt.KeeperBenchmarkTestConfig

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
	BlockchainClient       *seth.Client                      // Client for the test to connect to the blockchain with
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
func (k *KeeperBenchmarkTest) Setup(env *environment.Environment, config tt.KeeperBenchmarkTestConfig) {
	startTime := time.Now()
	k.TestReporter.Summary.StartTime = startTime.UnixMilli()
	k.ensureInputValues()
	k.env = env
	k.namespace = k.env.Cfg.Namespace
	inputs := k.Inputs
	k.testConfig = config

	k.keeperRegistries = make([]contracts.KeeperRegistry, len(inputs.RegistryVersions))
	k.keeperRegistrars = make([]contracts.KeeperRegistrar, len(inputs.RegistryVersions))
	k.keeperConsumerContracts = make([]contracts.AutomationConsumerBenchmark, len(inputs.RegistryVersions))
	k.upkeepIDs = make([][]*big.Int, len(inputs.RegistryVersions))
	k.log.Debug().Interface("TestInputs", inputs).Msg("Setting up benchmark test")

	// if not present disable it
	if k.testConfig.GetKeeperConfig().Resiliency == nil {
		k.testConfig.GetKeeperConfig().Resiliency = &keepertestconfig.ResiliencyConfig{
			ContractCallLimit:    ptr.Ptr(uint(0)),
			ContractCallInterval: ptr.Ptr(blockchain.StrDuration{Duration: 0 * time.Second}),
		}
	}

	var err error
	// Connect to networks and prepare for contract deployment
	k.chainlinkNodes, err = client.ConnectChainlinkNodes(k.env)
	require.NoError(k.t, err, "Connecting to chainlink nodes shouldn't fail")

	if len(inputs.RegistryVersions) > 1 && !inputs.ForceSingleTxnKey {
		for nodeIndex, node := range k.chainlinkNodes {
			for registryIndex := 1; registryIndex < len(inputs.RegistryVersions); registryIndex++ {
				k.log.Debug().Str("URL", node.URL()).Int("NodeIndex", nodeIndex).Int("RegistryIndex", registryIndex).Msg("Create Tx key")
				_, _, err := node.CreateTxKey("evm", fmt.Sprint(k.Inputs.BlockchainClient.ChainID))
				require.NoError(k.t, err, "Creating transaction key shouldn't fail")
			}
		}
	}

	c := inputs.Contracts

	if common.IsHexAddress(c.LinkTokenAddress) {
		_, err = contracts.LoadLinkTokenContract(k.log, k.chainClient, common.HexToAddress(c.LinkTokenAddress))
		require.NoError(k.t, err, "Loading Link Token Contract shouldn't fail")
	} else {
		k.linkToken, err = contracts.DeployLinkTokenContract(k.log, k.chainClient)
		require.NoError(k.t, err, "Deploying mock Link Token Contract feed shouldn't fail")
	}

	if common.IsHexAddress(c.EthFeedAddress) {
		_, err = contracts.LoadMockETHLINKFeed(k.chainClient, common.HexToAddress(c.EthFeedAddress))
		require.NoError(k.t, err, "Loading ETH-Link feed Contract shouldn't fail")
	} else {
		k.ethFeed, err = contracts.DeployMockETHLINKFeed(k.chainClient, big.NewInt(2e18))
		require.NoError(k.t, err, "Deploying mock ETH-Link feed shouldn't fail")
	}

	if common.IsHexAddress(c.GasFeedAddress) {
		k.gasFeed, err = contracts.LoadMockGASFeed(k.chainClient, common.HexToAddress(c.GasFeedAddress))
		require.NoError(k.t, err, "Loading Gas feed Contract shouldn't fail")
	} else {
		k.gasFeed, err = contracts.DeployMockGASFeed(k.chainClient, big.NewInt(2e11))
		require.NoError(k.t, err, "Deploying mock gas feed shouldn't fail")
	}

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
		if inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_0 || inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_1 || inputs.RegistryVersions[index] == ethereum.RegistryVersion_2_2 {
			nodesToFund = k.chainlinkNodes[1:]
		}
		err = actions_seth.FundChainlinkNodesAtKeyIndexFromRootAddress(k.log, k.chainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(nodesToFund), k.Inputs.ChainlinkNodeFunding, index)
		require.NoError(k.t, err, "Funding Chainlink nodes shouldn't fail")
	}

	k.log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Benchmark Test Setup")
	err = k.SendSlackNotification(nil, config)
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
	startingBlock, err := k.chainClient.Client.BlockNumber(testcontext.Get(k.t))
	require.NoError(k.t, err, "Error getting latest block number")
	k.startingBlock = big.NewInt(0).SetUint64(startingBlock)
	startTime := time.Now()

	nodesWithoutBootstrap := k.chainlinkNodes[1:]

	for rIndex := range k.keeperRegistries {
		var txKeyId = rIndex
		if inputs.ForceSingleTxnKey {
			txKeyId = 0
		}
		kr := k.keeperRegistries[rIndex]
		ocrConfig, err := actions.BuildAutoOCR2ConfigVarsWithKeyIndex(
			k.t, nodesWithoutBootstrap, *inputs.KeeperRegistrySettings, kr.Address(), k.Inputs.DeltaStage, txKeyId, common.Address{}, kr.ChainModuleAddress(), kr.ReorgProtectionEnabled(),
		)
		require.NoError(k.t, err, "Building OCR config shouldn't fail")

		rv := inputs.RegistryVersions[rIndex]
		// Send keeper jobs to registry and chainlink nodes
		if rv == ethereum.RegistryVersion_2_0 || rv == ethereum.RegistryVersion_2_1 || rv == ethereum.RegistryVersion_2_2 {
			actions.CreateOCRKeeperJobs(k.t, k.chainlinkNodes, kr.Address(), k.chainClient.ChainID, txKeyId, rv)
			if rv == ethereum.RegistryVersion_2_0 {
				err = kr.SetConfig(*inputs.KeeperRegistrySettings, ocrConfig)
			} else {
				err = kr.SetConfigTypeSafe(ocrConfig)
			}
			require.NoError(k.t, err, "Registry config should be be set successfully")
			// Give time for OCR nodes to bootstrap
			time.Sleep(1 * time.Minute)
		} else {
			actions.CreateKeeperJobsWithKeyIndex(k.t, k.chainlinkNodes, kr, txKeyId, ocrConfig, fmt.Sprint(k.chainClient.ChainID))
		}
	}

	k.log.Info().Msgf("Waiting for %d blocks for all upkeeps to be performed", inputs.Upkeeps.BlockRange+inputs.UpkeepSLA)

	errgroup, errCtx := errgroup.WithContext(context.Background())

	var startedObservations = atomic.Int32{}
	var finishedObservations = atomic.Int32{}

	// We create as many channels as listening goroutines (1 per upkeep). In the background we will be fanning out
	// headers that we get from a single channel connected to EVM node to all upkeep-specific channels.
	headerCh := make(chan *blockchain.SafeEVMHeader, 10)
	sub, err := k.chainClient.Client.Client().EthSubscribe(context.Background(), headerCh, "newHeads")
	require.NoError(k.t, err, "Subscribing to new headers for upkeep observation shouldn't fail")

	totalNumberOfChannels := 0
	for rIndex := range k.keeperRegistries {
		totalNumberOfChannels += len(k.upkeepIDs[rIndex])
	}

	contractChannels := make([]chan *blockchain.SafeEVMHeader, totalNumberOfChannels)
	for idx := 0; idx < totalNumberOfChannels; idx++ {
		contractChannels[idx] = make(chan *blockchain.SafeEVMHeader, 10) // Buffered just in case processing is slow
	}

	// signals all goroutines to stop when subscription error occurs
	stopAllGoroutinesCh := make(chan struct{})

	// this goroutine fans out headers to goroutines in the background
	// and exists when all goroutines are done or when an error occurs
	go func() {
		defer func() {
			// close all fanning out channels at the very end
			for _, ch := range contractChannels {
				close(ch)
			}
			k.log.Debug().Msg("Closed header distribution channels")
		}()
		for {
			select {
			case header := <-headerCh:
				k.log.Trace().Int64("Number", header.Number.Int64()).Msg("Fanning out new header")
				for _, ch := range contractChannels {
					ch <- header
				}
			// we don't really care if it was a success or an error, we just want to exit
			// if it was an error, we will have an error in the main goroutine
			case <-errCtx.Done():
				k.log.Debug().Msg("All goroutines finished.")
				sub.Unsubscribe()
				return
			case err := <-sub.Err():
				// no need to unsubscribe, subscription errored
				k.log.Error().Err(err).Msg("header subscription failed. Trying to reconnect...")
				connectionLostAt := time.Now()
				// we use infinite loop here on purposes, these nodes can be down for extended periods of time ¯\_(ツ)_/¯
			RECONNECT:
				for {
					sub, err = k.chainClient.Client.Client().EthSubscribe(context.Background(), headerCh, "newHeads")
					if err == nil {
						break RECONNECT
					}

					time.Sleep(5 * time.Second)
				}
				k.log.Info().Str("Reconnect Time", time.Since(connectionLostAt).String()).Msg("Reconnected to header subscription")
			}
		}
	}()

	currentChannelIndex := 0
	for rIndex := range k.keeperRegistries {
		for index, upkeepID := range k.upkeepIDs[rIndex] {
			chIndex := currentChannelIndex
			currentChannelIndex++
			upkeepIDCopy := upkeepID
			registryIndex := rIndex
			upkeepIndex := int64(index)
			errgroup.Go(func() error {
				startedObservations.Add(1)
				k.log.Info().Int("Channel index", chIndex).Str("UpkeepID", upkeepIDCopy.String()).Msg("Starting upkeep observation")

				confirmer := contracts.NewKeeperConsumerBenchmarkUpkeepObserver(
					k.keeperConsumerContracts[registryIndex],
					k.keeperRegistries[registryIndex],
					upkeepIDCopy,
					inputs.Upkeeps.BlockRange+inputs.UpkeepSLA,
					inputs.UpkeepSLA,
					&k.TestReporter,
					upkeepIndex,
					inputs.Upkeeps.FirstEligibleBuffer,
					k.log,
				)

				k.log.Debug().Str("UpkeepID", upkeepIDCopy.String()).Msg("Stared listening to new headers for upkeep observation")

				for {
					select {
					case <-stopAllGoroutinesCh: // header listening failed, exit
						return errors.New("header distribution channel closed")
					case <-errCtx.Done(): //one of goroutines errored, shut down gracefully, no need to return error
						k.log.Error().Err(errCtx.Err()).Str("UpkeepID", upkeepIDCopy.String()).Msg("Stopping obervations due to error in one of the goroutines")
						return nil
					case header := <-contractChannels[chIndex]: // new block, check if upkeep was performed
						k.log.Trace().Interface("Header number", header.Number).Str("UpkeepID", upkeepIDCopy.String()).Msg("Started processing new header")
						finished, headerErr := confirmer.ReceiveHeader(header)
						if headerErr != nil {
							k.log.Err(headerErr).Str("UpkeepID", upkeepIDCopy.String()).Msg("Error processing header")
							return errors.Wrapf(headerErr, "error processing header for upkeep %s", upkeepIDCopy.String())
						}

						if finished { // observations should be completed as we are beyond block range, if there are not there's a bug in test code
							finishedObservations.Add(1)
							k.log.Info().Str("Done/Total", fmt.Sprintf("%d/%d", finishedObservations.Load(), startedObservations.Load())).Str("UpkeepID", upkeepIDCopy.String()).Msg("Upkeep observation completed")

							if confirmer.Complete() {
								confirmer.LogDetails()
								return nil
							}
							return fmt.Errorf("confimer has finished, but without completing observation, this should never happen. Review your code. UpkdeepID: %s", upkeepIDCopy.String())
						}
						k.log.Trace().Interface("Header number", header.Number).Str("UpkeepID", upkeepIDCopy.String()).Msg("Finished processing new header")
					}
				}
			})
		}
	}

	if err := errgroup.Wait(); err != nil {
		k.t.Fatalf("errored when waiting for upkeeps: %v", err)
	}

	// Close header distribution channel once all observations are done
	close(stopAllGoroutinesCh)

	// Main test loop
	k.observeUpkeepEvents()

	// Collect logs for each registry to calculate test metrics
	// This test generates a LOT of logs, and we need to break up our reads, or risk getting rate-limited by the node
	var (
		endBlock             = big.NewInt(0).Add(k.startingBlock, big.NewInt(u.BlockRange))
		registryLogs         = make([][]types.Log, len(k.keeperRegistries))
		blockBatchSize int64 = 100
	)
	for rIndex := range k.keeperRegistries {
		// Variables for the full registry
		var (
			logs            []types.Log
			timeout         = 5 * time.Second
			addr            = k.keeperRegistries[rIndex].Address()
			queryStartBlock = big.NewInt(0).Set(k.startingBlock)
		)

		// Gather logs from the registry in 100 block chunks to avoid read limits
		for queryStartBlock.Cmp(endBlock) < 0 {
			filterQuery := geth.FilterQuery{
				Addresses: []common.Address{common.HexToAddress(addr)},
				FromBlock: queryStartBlock,
				ToBlock:   big.NewInt(0).Add(queryStartBlock, big.NewInt(blockBatchSize)),
			}

			// This RPC call can possibly time out or otherwise die. Failure is not an option, keep retrying to get our stats.
			err = fmt.Errorf("initial error") // to ensure our for loop runs at least once
			for err != nil {
				ctx, cancel := context.WithTimeout(testcontext.Get(k.t), timeout)
				logs, err = k.chainClient.Client.FilterLogs(ctx, filterQuery)
				cancel()
				if err != nil {
					k.log.Error().
						Err(err).
						Interface("Filter Query", filterQuery).
						Str("Timeout", timeout.String()).
						Msg("Error getting logs from chain, trying again")
					timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
					continue
				}
				k.log.Info().
					Uint64("From Block", queryStartBlock.Uint64()).
					Uint64("To Block", filterQuery.ToBlock.Uint64()).
					Int("Log Count", len(logs)).
					Str("Registry Address", addr).
					Msg("Collected logs")
				queryStartBlock.Add(queryStartBlock, big.NewInt(blockBatchSize))
				registryLogs[rIndex] = append(registryLogs[rIndex], logs...)
			}
		}
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
	if err != nil && k.Inputs.BlockchainClient.Cfg.IsSimulatedNetwork() {
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
	*seth.Client,
	string,
	[]*client.ChainlinkK8sClient,
	reportModel.TestReporter,
	reportModel.GrafanaURLProvider,
) {
	return t, k.chainClient, k.namespace, k.chainlinkNodes, &k.TestReporter, k.testConfig
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

	ctx, cancel := context.WithTimeout(testcontext.Get(k.t), 5*time.Second)
	sub, err := k.chainClient.Client.SubscribeFilterLogs(ctx, filterQuery, eventLogs)
	cancel()
	require.NoError(k.t, err, "Subscribing to upkeep performed events log shouldn't fail")

	interruption := make(chan os.Signal, 1)
	//nolint:staticcheck //ignore SA1016 we need to send the os.Kill signal
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

					ctx, cancel := context.WithTimeout(testcontext.Get(k.t), backoff)
					sub, err = k.chainClient.Client.SubscribeFilterLogs(ctx, filterQuery, eventLogs)
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
			}
		}
	}()
}

// contractABI returns the ABI of the proper keeper registry contract
func (k *KeeperBenchmarkTest) contractABI(rIndex int) *abi.ABI {
	contractABI, err := contracts.GetRegistryContractABI(k.Inputs.RegistryVersions[rIndex])
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

func (k *KeeperBenchmarkTest) SendSlackNotification(slackClient *slack.Client, config tt.KeeperBenchmarkTestConfig) error {
	if slackClient == nil {
		slackClient = slack.New(reportModel.SlackAPIKey)
	}

	grafanaUrl, err := config.GetGrafanaBaseURL()
	if err != nil {
		return err
	}

	dashboardUrl, err := config.GetGrafanaDashboardURL()
	if err != nil {
		return err
	}

	headerText := ":white_check_mark: Automation Benchmark Test STARTED :white_check_mark:"
	formattedDashboardUrl := fmt.Sprintf("%s%s?from=%d&to=%s&var-namespace=%s&var-cl_node=chainlink-0-0", grafanaUrl, dashboardUrl, k.TestReporter.Summary.StartTime, "now", k.env.Cfg.Namespace)
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
		err       error
	)

	// Contract deployment is different for legacy keepers and OCR automation
	if registryVersion <= ethereum.RegistryVersion_1_3 { // Legacy keeper - v1.X
		registry, err = contracts.DeployKeeperRegistry(k.chainClient, &contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        k.linkToken.Address(),
			ETHFeedAddr:     k.ethFeed.Address(),
			GasFeedAddr:     k.gasFeed.Address(),
			TranscoderAddr:  actions.ZeroAddress.Hex(),
			RegistrarAddr:   actions.ZeroAddress.Hex(),
			Settings:        *k.Inputs.KeeperRegistrySettings,
		})
		require.NoError(k.t, err, "Deploying registry contract shouldn't fail")

		// Fund the registry with 1 LINK * amount of AutomationConsumerBenchmark contracts
		err := k.linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(k.Inputs.Upkeeps.NumberOfUpkeeps))))
		require.NoError(k.t, err, "Funding keeper registry contract shouldn't fail")

		registrarSettings := contracts.KeeperRegistrarSettings{
			AutoApproveConfigType: 2,
			AutoApproveMaxAllowed: math.MaxUint16,
			RegistryAddr:          registry.Address(),
			MinLinkJuels:          big.NewInt(0),
		}

		registrar, err = contracts.DeployKeeperRegistrar(k.chainClient, registryVersion, k.linkToken.Address(), registrarSettings)
		require.NoError(k.t, err, "Funding keeper registrar contract shouldn't fail")
	} else { // OCR automation - v2.X
		registry, registrar = actions_seth.DeployAutoOCRRegistryAndRegistrar(
			k.t, k.chainClient, registryVersion, *k.Inputs.KeeperRegistrySettings, k.linkToken,
		)

		// Fund the registry with LINK
		err := k.linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(k.Inputs.Upkeeps.NumberOfUpkeeps))))
		require.NoError(k.t, err, "Funding keeper registry contract shouldn't fail")
		ocrConfig, err := actions.BuildAutoOCR2ConfigVars(k.t, k.chainlinkNodes[1:], *k.Inputs.KeeperRegistrySettings, registrar.Address(), k.Inputs.DeltaStage, registry.ChainModuleAddress(), registry.ReorgProtectionEnabled())
		require.NoError(k.t, err, "Building OCR config shouldn't fail")
		k.log.Debug().Interface("KeeperRegistrySettings", *k.Inputs.KeeperRegistrySettings).Interface("OCRConfig", ocrConfig).Msg("Config")
		require.NoError(k.t, err, "Error building OCR config vars")
		if registryVersion == ethereum.RegistryVersion_2_0 {
			err = registry.SetConfig(*k.Inputs.KeeperRegistrySettings, ocrConfig)
		} else {
			err = registry.SetConfigTypeSafe(ocrConfig)
		}
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

	err = actions_seth.DeployMultiCallAndFundDeploymentAddresses(k.chainClient, k.linkToken, upkeep.NumberOfUpkeeps, linkFunds)
	require.NoError(k.t, err, "Sending link funds to deployment addresses shouldn't fail")

	upkeepIds := actions_seth.RegisterUpkeepContractsWithCheckData(k.t, k.chainClient, k.linkToken, linkFunds, uint32(upkeep.UpkeepGasLimit), registry, registrar, upkeep.NumberOfUpkeeps, upkeepAddresses, checkData, false, false)

	k.keeperRegistries[index] = registry
	k.keeperRegistrars[index] = registrar
	k.upkeepIDs[index] = upkeepIds
	k.keeperConsumerContracts[index] = consumer
}

func (k *KeeperBenchmarkTest) DeployKeeperConsumersBenchmark() contracts.AutomationConsumerBenchmark {
	// Deploy consumer
	var err error
	var keeperConsumerInstance contracts.AutomationConsumerBenchmark
	if *k.testConfig.GetKeeperConfig().Resiliency.ContractCallLimit != 0 && k.testConfig.GetKeeperConfig().Resiliency.ContractCallInterval.Duration != 0 {
		maxRetryAttempts := *k.testConfig.GetKeeperConfig().Resiliency.ContractCallLimit
		callRetryDelay := k.testConfig.GetKeeperConfig().Resiliency.ContractCallInterval.Duration
		keeperConsumerInstance, err = contracts.DeployKeeperConsumerBenchmarkWithRetry(k.chainClient, k.log, maxRetryAttempts, callRetryDelay)
		if err != nil {
			k.log.Error().Err(err).Msg("Deploying AutomationConsumerBenchmark instance shouldn't fail")
			keeperConsumerInstance, err = contracts.DeployKeeperConsumerBenchmarkWithRetry(k.chainClient, k.log, maxRetryAttempts, callRetryDelay)
			require.NoError(k.t, err, "Error deploying AutomationConsumerBenchmark")
		}
	} else {
		keeperConsumerInstance, err = contracts.DeployKeeperConsumerBenchmark(k.chainClient)
		if err != nil {
			k.log.Error().Err(err).Msg("Deploying AutomationConsumerBenchmark instance %d shouldn't fail")
			keeperConsumerInstance, err = contracts.DeployKeeperConsumerBenchmark(k.chainClient)
			require.NoError(k.t, err, "Error deploying AutomationConsumerBenchmark")
		}
	}
	k.log.Debug().
		Str("Contract Address", keeperConsumerInstance.Address()).
		Msg("Deployed Keeper Benchmark Contract")

	k.log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerInstance
}
