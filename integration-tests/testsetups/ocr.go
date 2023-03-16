// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs       *OCRSoakTestInputs
	TestReporter testreporters.OCRSoakTestReporter

	testEnvironment *environment.Environment
	chainlinkNodes  []*client.Chainlink
	chainClient     blockchain.EVMClient
	mockServer      *ctfClient.MockserverClient

	ocrInstances          []contracts.OffchainAggregator
	ocrInstanceMap        map[string]contracts.OffchainAggregator // address : instance
	OperatorForwarderFlow bool
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	BlockchainClient     blockchain.EVMClient // Client for the test to connect to the blockchain with
	TestDuration         time.Duration        // How long to run the test for (assuming things pass)
	NumberOfContracts    int                  // Number of OCR contracts to launch
	ChainlinkNodeFunding *big.Float           // Amount of ETH to fund each chainlink node with
	RoundTimeout         time.Duration        // How long to wait for a round to update before failing the test
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
			ContractReports:       make(map[string]*testreporters.OCRSoakTestReport),
			ExpectedRoundDuration: inputs.ExpectedRoundTime,
		},
		ocrInstanceMap: make(map[string]contracts.OffchainAggregator),
	}
}

// Setup sets up the test environment, deploying contracts and funding chainlink nodes
func (o *OCRSoakTest) Setup(t *testing.T, env *environment.Environment) {
	o.ensureInputValues(t)
	o.testEnvironment = env
	var err error

	// Make connections to soak test resources
	contractDeployer, err := contracts.NewContractDeployer(o.chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	o.chainlinkNodes, err = client.ConnectChainlinkNodes(env)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	o.mockServer, err = ctfClient.ConnectMockServer(env)
	require.NoError(t, err, "Creating mockserver clients shouldn't fail")
	o.chainClient.ParallelTransactions(true)
	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(o.chainlinkNodes[1:], o.chainClient, o.Inputs.ChainlinkNodeFunding)
	require.NoError(t, err, "Error funding Chainlink nodes")

	if o.OperatorForwarderFlow {
		contractLoader, err := contracts.NewContractLoader(o.chainClient)
		require.NoError(t, err, "Loading contracts shouldn't fail")

		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
			t, contractDeployer, linkTokenContract, o.chainClient, len(o.chainlinkNodes[1:]),
		)
		forwarderNodes := o.chainlinkNodes[1:]
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.chainlinkNodes[1:])
		require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range forwarderNodes {
			actions.AcceptAuthorizedReceiversOperator(
				t, operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, o.chainClient, contractLoader,
			)
			require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
			actions.TrackForwarder(t, o.chainClient, authorizedForwarders[i], forwarderNodes[i])
			err = o.chainClient.WaitForEvents()
		}

		o.ocrInstances = actions.DeployOCRContractsForwarderFlow(
			t,
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.chainlinkNodes,
			authorizedForwarders,
			o.chainClient,
		)
	} else {
		o.ocrInstances, err = actions.DeployOCRContracts(
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.chainlinkNodes,
			o.chainClient,
		)
		require.NoError(t, err)
	}

	err = o.chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range o.ocrInstances {
		o.ocrInstanceMap[ocrInstance.Address()] = ocrInstance
		o.TestReporter.ContractReports[ocrInstance.Address()] = testreporters.NewOCRSoakTestReport(
			ocrInstance.Address(),
			o.Inputs.StartingAdapterValue,
			o.Inputs.ExpectedRoundTime,
		)
	}
	log.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run(t *testing.T) {
	// Set initial value and create jobs
	err := actions.SetAllAdapterResponsesToTheSameValue(o.Inputs.StartingAdapterValue, o.ocrInstances, o.chainlinkNodes, o.mockServer)
	require.NoError(t, err, "Error setting adapter responses")
	if o.OperatorForwarderFlow {
		actions.CreateOCRJobsWithForwarder(t, o.ocrInstances, o.chainlinkNodes, o.mockServer)
	} else {
		err = actions.CreateOCRJobs(o.ocrInstances, o.chainlinkNodes, o.mockServer)
		require.NoError(t, err, "Error creating OCR jobs")
	}

	log.Info().
		Str("Test Duration", o.Inputs.TestDuration.Truncate(time.Second).String()).
		Str("Round Timeout", o.Inputs.RoundTimeout.String()).
		Int("Number of OCR Contracts", len(o.ocrInstances)).
		Msg("Starting OCR Soak Test")

	testDuration := time.NewTimer(o.Inputs.TestDuration)

	// *********************
	// ***** Test Loop *****
	// *********************
	lastAdapterValue, currentAdapterValue := o.Inputs.StartingAdapterValue, o.Inputs.StartingAdapterValue*25
	newRoundTrigger, expiredRoundTrigger := time.NewTimer(0), time.NewTimer(o.Inputs.RoundTimeout)
	answerUpdated := make(chan *ethereum.OffchainAggregatorAnswerUpdated)
	o.subscribeOCREvents(t, answerUpdated)
	remainingExpectedAnswers := len(o.ocrInstances)
	testOver := false
	for {
		select {
		case <-testDuration.C:
			testOver = true
			log.Warn().Msg("Soak Test Duration Reached. Completing Final Round")
		case answer := <-answerUpdated:
			if o.processNewAnswer(t, answer) {
				remainingExpectedAnswers--
			}
			if remainingExpectedAnswers <= 0 {
				if testOver {
					log.Info().Msg("Soak Test Complete")
					return
				}
				log.Info().
					Str("Wait time", o.Inputs.TimeBetweenRounds.String()).
					Msg("All Expected Answers Reported. Waiting to Start a New Round")
				remainingExpectedAnswers = len(o.ocrInstances)
				newRoundTrigger, expiredRoundTrigger = time.NewTimer(o.Inputs.TimeBetweenRounds), time.NewTimer(o.Inputs.RoundTimeout)
			}
		case <-newRoundTrigger.C:
			lastAdapterValue, currentAdapterValue = currentAdapterValue, lastAdapterValue
			o.triggerNewRound(t, currentAdapterValue)
		case <-expiredRoundTrigger.C:
			log.Warn().Msg("OCR round timed out")
			expiredRoundTrigger = time.NewTimer(o.Inputs.RoundTimeout)
			remainingExpectedAnswers = len(o.ocrInstances)
			o.triggerNewRound(t, rand.Intn(o.Inputs.StartingAdapterValue*25-1-o.Inputs.StartingAdapterValue)+o.Inputs.StartingAdapterValue) // #nosec G404 | Just triggering a random number
		}
	}
}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	*environment.Environment,
	[]*client.Chainlink,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, o.testEnvironment, o.chainlinkNodes, &o.TestReporter, o.chainClient
}

// *********************
// ****** Helpers ******
// *********************

func (o *OCRSoakTest) processNewEvent(
	t *testing.T,
	eventSub geth.Subscription,
	answerUpdated chan *ethereum.OffchainAggregatorAnswerUpdated,
	event *types.Log,
	eventDetails *abi.Event,
	ocrInstance contracts.OffchainAggregator,
	contractABI *abi.ABI,
) {
	errorChan := make(chan error)
	eventConfirmed := make(chan bool)
	err := o.chainClient.ProcessEvent(eventDetails.Name, event, eventConfirmed, errorChan)
	if err != nil {
		log.Error().Err(err).Str("Hash", event.TxHash.Hex()).Str("Event", eventDetails.Name).Msg("Error trying to process event")
		return
	}
	log.Debug().
		Str("Event", eventDetails.Name).
		Str("Address", event.Address.Hex()).
		Str("Hash", event.TxHash.Hex()).
		Msg("Attempting to Confirm Event")
	for {
		select {
		case err := <-errorChan:
			log.Error().Err(err).Msg("Error while confirming event")
			return
		case confirmed := <-eventConfirmed:
			if confirmed {
				if eventDetails.Name == "AnswerUpdated" { // Send AnswerUpdated events to answerUpdated channel to handle in main loop
					answer, err := ocrInstance.ParseEventAnswerUpdated(*event)
					require.NoError(t, err, "Parsing AnswerUpdated event log in OCR instance shouldn't fail")
					answerUpdated <- answer
				}
				log.Info().
					Str("Contract", event.Address.Hex()).
					Str("Event Name", eventDetails.Name).
					Uint64("Header Number", event.BlockNumber).
					Msg("Contract Event Published")
			}
			return
		}
	}
}

// marshalls new answer events into manageable Go struct for further processing and reporting
func (o *OCRSoakTest) processNewAnswer(t *testing.T, newAnswer *ethereum.OffchainAggregatorAnswerUpdated) bool {
	// Updated Info
	answerAddress := newAnswer.Raw.Address.Hex()
	_, tracked := o.TestReporter.ContractReports[answerAddress]
	if !tracked {
		log.Error().Str("Untracked Address", answerAddress).Msg("Received AnswerUpdated event on an untracked OCR instance")
		return false
	}
	processedAnswer := &testreporters.OCRAnswerUpdated{}
	processedAnswer.ContractAddress = newAnswer.Raw.Address.Hex()
	processedAnswer.UpdatedTime = time.Unix(newAnswer.UpdatedAt.Int64(), 0)
	processedAnswer.UpdatedRoundId = newAnswer.RoundId.Uint64()
	processedAnswer.UpdatedBlockNum = newAnswer.Raw.BlockNumber
	processedAnswer.UpdatedAnswer = int(newAnswer.Current.Int64())
	processedAnswer.UpdatedBlockHash = newAnswer.Raw.BlockHash.Hex()
	processedAnswer.RoundTxHash = newAnswer.Raw.TxHash.Hex()

	// On-Chain Info
	updatedOCRInstance := o.ocrInstanceMap[answerAddress]
	onChainData, err := updatedOCRInstance.GetRound(context.Background(), newAnswer.RoundId)
	require.NoError(t, err, "Error retrieving on-chain data for '%s' at round '%d'", answerAddress, processedAnswer.UpdatedRoundId)
	processedAnswer.OnChainAnswer = int(onChainData.Answer.Int64())
	processedAnswer.OnChainRoundId = onChainData.RoundId.Uint64()

	return o.TestReporter.ContractReports[answerAddress].NewAnswerUpdated(processedAnswer)
}

// triggers a new OCR round by setting a new mock adapter value
func (o *OCRSoakTest) triggerNewRound(t *testing.T, currentAdapterValue int) {
	startingBlockNum, err := o.chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error retrieving latest block number")

	for _, report := range o.TestReporter.ContractReports {
		report.NewAnswerExpected(currentAdapterValue, startingBlockNum)
	}
	err = actions.SetAllAdapterResponsesToTheSameValue(currentAdapterValue, o.ocrInstances, o.chainlinkNodes, o.mockServer)
	require.NoError(t, err, "Error setting adapter responses")
	log.Info().
		Int("Value", currentAdapterValue).
		Msg("Starting a New OCR Round")
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
	require.GreaterOrEqual(t, inputs.RoundTimeout, inputs.ExpectedRoundTime, "Expected RoundTimeout to be greater than ExpectedRoundTime")
	require.NotNil(t, inputs.TimeBetweenRounds, "Expected TimeBetweenRounds to be set")
	require.Less(t, inputs.TimeBetweenRounds, time.Hour, "TimeBetweenRounds must be less than 1 hour")
}

// subscribeToAnswerUpdatedEvent subscribes to the event log for AnswerUpdated event and
// verifies if the answer is matching with the expected value
func (o *OCRSoakTest) subscribeOCREvents(
	t *testing.T,
	answerUpdated chan *ethereum.OffchainAggregatorAnswerUpdated,
) {
	contractABI, err := ethereum.OffchainAggregatorMetaData.GetAbi()
	require.NoError(t, err, "Getting contract abi for OCR shouldn't fail")
	latestBlockNum, err := o.chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Subscribing to contract event log for OCR instance shouldn't fail")
	query := geth.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
		Addresses: []common.Address{},
	}
	for i := 0; i < len(o.ocrInstances); i++ {
		query.Addresses = append(query.Addresses, common.HexToAddress(o.ocrInstances[i].Address()))
	}
	eventLogs := make(chan types.Log)
	sub, err := o.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	require.NoError(t, err, "Subscribing to contract event log for OCR instance shouldn't fail")

	go func() {
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				log.Error().Err(err).Msg("Error while watching for new contract events. Retrying Subscription")
				sub.Unsubscribe()

				sub, err = o.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				require.NoError(t, err, "Subscribing to contract event log for OCR instance shouldn't fail")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				require.NoError(t, err, "Getting event details for OCR instances shouldn't fail")

				go o.processNewEvent(t, sub, answerUpdated, &vLog, eventDetails, o.ocrInstances[0], contractABI)
			}
		}
	}()
}
