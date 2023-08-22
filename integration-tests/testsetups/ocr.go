// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

import (
	"context"
	"math/big"
	"sort"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs                *OCRSoakTestInputs
	TestReporter          testreporters.OCRSoakTestReporter
	OperatorForwarderFlow bool

	testEnvironment *environment.Environment
	bootstrapNode   *client.Chainlink
	workerNodes     []*client.Chainlink
	chainClient     blockchain.EVMClient
	mockServer      *ctfClient.MockserverClient
	mockPath        string
	filterQuery     geth.FilterQuery

	ocrTestStates  []*testreporters.OCRTestState
	rpcIssues      []*testreporters.RPCIssue
	combinedEvents [][]string

	ocrInstances   []contracts.OffchainAggregator
	ocrInstanceMap map[string]contracts.OffchainAggregator // address : instance
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	BlockchainClient     blockchain.EVMClient // Client for the test to connect to the blockchain with
	TestDuration         time.Duration        // How long to run the test for
	NumberOfContracts    int                  // Number of OCR contracts to launch
	ChainlinkNodeFunding *big.Float           // Amount of ETH to fund each chainlink node with
	ExpectedRoundTime    time.Duration        // How long each round is expected to take
	TimeBetweenRounds    time.Duration        // How long to wait after a completed round to start a new one, set 0 for instant
	StartingAdapterValue int
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(inputs *OCRSoakTestInputs) *OCRSoakTest {
	if inputs.StartingAdapterValue == 0 {
		inputs.StartingAdapterValue = 5
	}
	return &OCRSoakTest{
		Inputs: inputs,
		TestReporter: testreporters.OCRSoakTestReporter{
			TestDuration: inputs.TestDuration,
		},
		ocrTestStates:  make([]*testreporters.OCRTestState, 0),
		mockPath:       "ocr",
		ocrInstanceMap: make(map[string]contracts.OffchainAggregator),
	}
}

// Setup sets up the test environment, deploying contracts and funding chainlink nodes
func (o *OCRSoakTest) Setup(t *testing.T, env *environment.Environment) {
	l := utils.GetTestLogger(t)
	o.ensureInputValues(t)
	o.testEnvironment = env
	var err error

	// Make connections to soak test resources
	contractDeployer, err := contracts.NewContractDeployer(o.chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	nodes, err := client.ConnectChainlinkNodes(env)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	o.bootstrapNode, o.workerNodes = nodes[0], nodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServer(env)
	require.NoError(t, err, "Creating mockserver clients shouldn't fail")
	o.chainClient.ParallelTransactions(true)
	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(o.workerNodes, o.chainClient, o.Inputs.ChainlinkNodeFunding)
	require.NoError(t, err, "Error funding Chainlink nodes")

	if o.OperatorForwarderFlow {
		contractLoader, err := contracts.NewContractLoader(o.chainClient)
		require.NoError(t, err, "Loading contracts shouldn't fail")

		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
			t, contractDeployer, linkTokenContract, o.chainClient, len(o.workerNodes),
		)
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.workerNodes)
		require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range o.workerNodes {
			actions.AcceptAuthorizedReceiversOperator(
				t, operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, o.chainClient, contractLoader,
			)
			require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
			actions.TrackForwarder(t, o.chainClient, authorizedForwarders[i], o.workerNodes[i])
			err = o.chainClient.WaitForEvents()
		}

		o.ocrInstances = actions.DeployOCRContractsForwarderFlow(
			t,
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.workerNodes,
			authorizedForwarders,
			o.chainClient,
		)
	} else {
		o.ocrInstances, err = actions.DeployOCRContracts(
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.bootstrapNode,
			o.workerNodes,
			o.chainClient,
		)
		require.NoError(t, err)
	}

	err = o.chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range o.ocrInstances {
		o.ocrInstanceMap[ocrInstance.Address()] = ocrInstance
	}
	l.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run(t *testing.T) {
	l := utils.GetTestLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	latestBlockNum, err := o.chainClient.LatestBlockNumber(ctx)
	cancel()
	require.NoError(t, err, "Error getting current block number")

	ocrAddresses := make([]common.Address, len(o.ocrInstances))
	for i, ocrInstance := range o.ocrInstances {
		ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
	}
	contractABI, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	require.NoError(t, err, "Error retrieving OCR contract ABI")
	o.filterQuery = geth.FilterQuery{
		Addresses: ocrAddresses,
		Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
	}

	if o.OperatorForwarderFlow {
		actions.CreateOCRJobsWithForwarder(t, o.ocrInstances, o.bootstrapNode, o.workerNodes, 5, o.mockServer)
	} else {
		err := actions.CreateOCRJobs(o.ocrInstances, o.bootstrapNode, o.workerNodes, 5, o.mockServer)
		require.NoError(t, err, "Error creating OCR jobs")
	}

	l.Info().
		Str("Test Duration", o.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Number of OCR Contracts", len(o.ocrInstances)).
		Msg("Starting OCR Soak Test")

	testDuration := time.After(o.Inputs.TestDuration)

	// *********************
	// ***** Test Loop *****
	// *********************
	lastAdapterValue, currentAdapterValue := o.Inputs.StartingAdapterValue, o.Inputs.StartingAdapterValue*25
	newRoundTrigger := time.NewTimer(0)
	defer newRoundTrigger.Stop()
	err = o.observeOCREvents(l)
	require.NoError(t, err, "Error subscribing to OCR events")

testLoop:
	for {
		select {
		case <-testDuration:
			break testLoop
		case <-newRoundTrigger.C:
			lastAdapterValue, currentAdapterValue = currentAdapterValue, lastAdapterValue
			err := o.triggerNewRound(t, currentAdapterValue)

			timerReset := o.Inputs.TimeBetweenRounds
			if err != nil {
				l.Error().Err(err).Int("Seconds Waiting", 5).Msg("Error triggering new round, waiting and trying again")
				timerReset = time.Second * 5
			}
			newRoundTrigger.Reset(timerReset)
		case t := <-o.chainClient.ConnectionIssue():
			o.rpcIssues = append(o.rpcIssues, &testreporters.RPCIssue{
				StartTime: t,
				Message:   "RPC Connection Lost",
			})
		case t := <-o.chainClient.ConnectionRestored():
			o.rpcIssues = append(o.rpcIssues, &testreporters.RPCIssue{
				StartTime: t,
				Message:   "RPC Connection Restored",
			})
		}
	}

	l.Info().Msg("Test Complete, collecting on-chain events to be collected")
	// Keep trying to collect events until we get them, no exceptions
	timeout := time.Second * 5
	err = o.collectEvents(l, timeout)
	for err != nil {
		timeout *= 2
		err = o.collectEvents(l, timeout)
	}
	o.TestReporter.RecordEvents(o.ocrTestStates, o.rpcIssues)
}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	*environment.Environment,
	[]*client.Chainlink,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, o.testEnvironment, append(o.workerNodes, o.bootstrapNode), &o.TestReporter, o.chainClient
}

// *********************
// ****** Helpers ******
// *********************

// observeOCREvents subscribes to OCR events and logs them to the test logger
// WARNING: Should only be used for observation and logging. This is not a reliable way to collect events.
func (o *OCRSoakTest) observeOCREvents(logger zerolog.Logger) error {
	eventLogs := make(chan types.Log)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eventSub, err := o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
	if err != nil {
		return err
	}

	go func() {
		defer cancel()
		for {
			select {
			case event := <-eventLogs:
				answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
				if err != nil {
					log.Warn().
						Err(err).
						Str("Address", event.Address.Hex()).
						Uint64("Block Number", event.BlockNumber).
						Msg("Error parsing event as AnswerUpdated")
					continue
				}
				logger.Info().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Answer Updated Event")
			case err = <-eventSub.Err():
				for err != nil {
					logger.Trace().
						Err(err).
						Interface("Query", o.filterQuery).
						Msg("Error while subscribed to OCR Logs. Resubscribing")
					ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
					eventSub, err = o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
				}
			}
		}
	}()

	return nil
}

// triggers a new OCR round by setting a new mock adapter value
func (o *OCRSoakTest) triggerNewRound(t *testing.T, currentAdapterValue int) error {
	l := utils.GetTestLogger(t)

	if len(o.ocrTestStates) > 0 {
		o.ocrTestStates[len(o.ocrTestStates)-1].EndTime = time.Now()
	}

	var (
		err          error
		attemptCount = 5
	)

	// It's possible the adapter is temporarily down, so we try a few times if we get errors
	for attemptCount > 0 {
		attemptCount--
		err = actions.SetAllAdapterResponsesToTheSameValue(currentAdapterValue, o.ocrInstances, o.workerNodes, o.mockServer)
		if err == nil {
			break
		}
		log.Warn().Err(err).
			Int("Attempts left", attemptCount).
			Msg("Error setting adapter responses, adapter possibly temporarily down, trying again")
	}
	if err != nil {
		return err
	}

	expectedState := &testreporters.OCRTestState{
		StartTime:   time.Now(),
		Answer:      int64(currentAdapterValue),
		FoundEvents: make(map[string][]*testreporters.FoundEvent),
	}
	for _, ocrInstance := range o.ocrInstances {
		expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
	}
	o.ocrTestStates = append(o.ocrTestStates, expectedState)
	l.Info().
		Int("Value", currentAdapterValue).
		Msg("Starting a New OCR Round")
	return nil
}

func (o *OCRSoakTest) collectEvents(logger zerolog.Logger, timeout time.Duration) error {
	start := time.Now()
	o.ocrTestStates[len(o.ocrTestStates)-1].EndTime = start // Set end time for last expected event
	logger.Info().Msg("Collecting on-chain events")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	contractEvents, err := o.chainClient.FilterLogs(ctx, o.filterQuery)
	if err != nil {
		log.Error().
			Err(err).
			Str("Time", time.Since(start).String()).
			Msg("Error collecting on-chain events")
		return err
	}

	sortedFoundEvents := make([]*testreporters.FoundEvent, 0)
	for _, event := range contractEvents {
		answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
		if err != nil {
			log.Error().
				Err(err).
				Str("Time", time.Since(start).String()).
				Msg("Error collecting on-chain events")
			return err
		}
		sortedFoundEvents = append(sortedFoundEvents, &testreporters.FoundEvent{
			StartTime:   time.Unix(answerUpdated.UpdatedAt.Int64(), 0),
			Address:     event.Address.Hex(),
			Answer:      answerUpdated.Current.Int64(),
			RoundID:     answerUpdated.RoundId.Uint64(),
			BlockNumber: event.BlockNumber,
		})
	}

	// Sort our events by time to make sure they are in order (don't trust RPCs)
	sort.Slice(sortedFoundEvents, func(i, j int) bool {
		return sortedFoundEvents[i].StartTime.Before(sortedFoundEvents[j].StartTime)
	})

	// Now match each found event with the expected event time frame
	expectedIndex := 0
	for _, event := range sortedFoundEvents {
		if !event.StartTime.Before(o.ocrTestStates[expectedIndex].EndTime) {
			expectedIndex++
			if expectedIndex >= len(o.ocrTestStates) {
				logger.Warn().
					Str("Event Time", event.StartTime.String()).
					Str("Expected End Time", o.ocrTestStates[expectedIndex].EndTime.String()).
					Msg("Found events after last expected end time, adding event to that final report, things might be weird")
			}
		}
		o.ocrTestStates[expectedIndex].FoundEvents[event.Address] = append(o.ocrTestStates[expectedIndex].FoundEvents[event.Address], event)
		o.ocrTestStates[expectedIndex].TimeLineEvents = append(o.ocrTestStates[expectedIndex].TimeLineEvents, event)
	}

	logger.Info().
		Str("Time", time.Since(start).String()).
		Msg("Collected on-chain events")
	return nil
}

// ensureValues ensures that all values needed to run the test are present
func (o *OCRSoakTest) ensureInputValues(t *testing.T) {
	inputs := o.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Need a valid blockchain client to use for the test")
	o.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.NumberOfContracts, 1, "Expecting at least 1 OCR contract")
	fund, _ := inputs.ChainlinkNodeFunding.Float64()
	require.Greater(t, fund, 0.0, "Expecting non-zero chainlink node funding amount")
	require.GreaterOrEqual(t, inputs.TestDuration, time.Minute*1, "Expected test duration to be more than a minute")
	require.GreaterOrEqual(t, inputs.ExpectedRoundTime, time.Second, "Expected ExpectedRoundTime to be greater than 1 second")
	require.NotNil(t, inputs.TimeBetweenRounds, "Expected TimeBetweenRounds to be set")
	require.Less(t, inputs.TimeBetweenRounds, time.Hour, "TimeBetweenRounds must be less than 1 hour")
}
