// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"math/rand"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"

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
func (t *OCRSoakTest) Setup(env *environment.Environment) {
	t.ensureInputValues()
	t.testEnvironment = env
	var err error

	// Make connections to soak test resources
	contractDeployer, err := contracts.NewContractDeployer(t.chainClient)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
	t.chainlinkNodes, err = client.ConnectChainlinkNodes(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	t.mockServer, err = ctfClient.ConnectMockServer(env)
	Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	t.chainClient.ParallelTransactions(true)
	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(t.chainlinkNodes[1:], t.chainClient, t.Inputs.ChainlinkNodeFunding)
	Expect(err).ShouldNot(HaveOccurred(), "Error funding Chainlink nodes")

	if t.OperatorForwarderFlow {
		contractLoader, err := contracts.NewContractLoader(t.chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Loading contracts shouldn't fail")

		By("Prepare forwarder contracts onchain")
		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(contractDeployer, linkTokenContract, t.chainClient, len(t.chainlinkNodes[1:]))
		forwarderNodes := t.chainlinkNodes[1:]
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(t.chainlinkNodes[1:])
		Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range forwarderNodes {
			actions.AcceptAuthorizedReceiversOperator(operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, t.chainClient, contractLoader)
			Expect(err).ShouldNot(HaveOccurred(), "Accepting Authorize Receivers on Operator shouldn't fail")
			By("Add forwarder track into DB")
			actions.TrackForwarder(t.chainClient, authorizedForwarders[i], forwarderNodes[i])
			err = t.chainClient.WaitForEvents()
		}

		t.ocrInstances = actions.DeployOCRContractsForwarderFlow(
			t.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			t.chainlinkNodes,
			authorizedForwarders,
			t.chainClient,
		)
	} else {
		t.ocrInstances = actions.DeployOCRContracts(
			t.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			t.chainlinkNodes,
			t.chainClient,
		)
	}

	err = t.chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range t.ocrInstances {
		t.ocrInstanceMap[ocrInstance.Address()] = ocrInstance
		t.TestReporter.ContractReports[ocrInstance.Address()] = testreporters.NewOCRSoakTestReport(
			ocrInstance.Address(),
			t.Inputs.StartingAdapterValue,
			t.Inputs.ExpectedRoundTime,
		)
	}
	log.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (t *OCRSoakTest) Run() {
	// Set initial value and create jobs
	By("Setting adapter responses",
		actions.SetAllAdapterResponsesToTheSameValue(t.Inputs.StartingAdapterValue, t.ocrInstances, t.chainlinkNodes, t.mockServer))
	if t.OperatorForwarderFlow {
		By("Creating OCR jobs operator forwarder flow", actions.CreateOCRJobsWithForwarder(t.ocrInstances, t.chainlinkNodes, t.mockServer))
	} else {
		By("Creating OCR jobs", actions.CreateOCRJobs(t.ocrInstances, t.chainlinkNodes, t.mockServer))
	}

	log.Info().
		Str("Test Duration", t.Inputs.TestDuration.Truncate(time.Second).String()).
		Str("Round Timeout", t.Inputs.RoundTimeout.String()).
		Int("Number of OCR Contracts", len(t.ocrInstances)).
		Msg("Starting OCR Soak Test")

	testDuration := time.NewTimer(t.Inputs.TestDuration)

	stopTestChannel := make(chan struct{}, 1)
	testsetups.StartRemoteControlServer("OCR Soak Test", stopTestChannel)

	// *********************
	// ***** Test Loop *****
	// *********************
	lastAdapterValue, currentAdapterValue := t.Inputs.StartingAdapterValue, t.Inputs.StartingAdapterValue*25
	newRoundTrigger, expiredRoundTrigger := time.NewTimer(0), time.NewTimer(t.Inputs.RoundTimeout)
	answerUpdated := make(chan *ethereum.OffchainAggregatorAnswerUpdated)
	t.subscribeOCREvents(answerUpdated)
	remainingExpectedAnswers := len(t.ocrInstances)
	testOver := false
	for {
		select {
		case <-stopTestChannel:
			t.TestReporter.UnexpectedShutdown = true
			log.Warn().Msg("Received shut down signal. Soak test stopping early")
			return
		case <-testDuration.C:
			testOver = true
			log.Warn().Msg("Soak Test Duration Reached. Completing Final Round")
		case answer := <-answerUpdated:
			if t.processNewAnswer(answer) {
				remainingExpectedAnswers--
			}
			if remainingExpectedAnswers <= 0 {
				if testOver {
					log.Info().Msg("Soak Test Complete")
					return
				}
				log.Info().
					Str("Wait time", t.Inputs.TimeBetweenRounds.String()).
					Msg("All Expected Answers Reported. Waiting to Start a New Round")
				remainingExpectedAnswers = len(t.ocrInstances)
				newRoundTrigger, expiredRoundTrigger = time.NewTimer(t.Inputs.TimeBetweenRounds), time.NewTimer(t.Inputs.RoundTimeout)
			}
		case <-newRoundTrigger.C:
			lastAdapterValue, currentAdapterValue = currentAdapterValue, lastAdapterValue
			t.triggerNewRound(currentAdapterValue)
		case <-expiredRoundTrigger.C:
			log.Warn().Msg("OCR round timed out")
			expiredRoundTrigger = time.NewTimer(t.Inputs.RoundTimeout)
			remainingExpectedAnswers = len(t.ocrInstances)
			t.triggerNewRound(rand.Int()) // #nosec G404 | Just triggering a random number
		}
	}
}

// Networks returns the networks that the test is running on
func (t *OCRSoakTest) TearDownVals() (*environment.Environment, []*client.Chainlink, reportModel.TestReporter, blockchain.EVMClient) {
	return t.testEnvironment, t.chainlinkNodes, &t.TestReporter, t.chainClient
}

// *********************
// ****** Helpers ******
// *********************

func (t *OCRSoakTest) processNewEvent(
	eventSub geth.Subscription,
	answerUpdated chan *ethereum.OffchainAggregatorAnswerUpdated,
	event *types.Log,
	eventDetails *abi.Event,
	ocrInstance contracts.OffchainAggregator,
	contractABI *abi.ABI,
) {
	defer GinkgoRecover()

	errorChan := make(chan error)
	eventConfirmed := make(chan bool)
	err := t.chainClient.ProcessEvent(eventDetails.Name, event, eventConfirmed, errorChan)
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
					Expect(err).ShouldNot(HaveOccurred(), "Parsing AnswerUpdated event log in OCR instance shouldn't fail")
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
func (t *OCRSoakTest) processNewAnswer(newAnswer *ethereum.OffchainAggregatorAnswerUpdated) bool {
	// Updated Info
	answerAddress := newAnswer.Raw.Address.Hex()
	_, tracked := t.TestReporter.ContractReports[answerAddress]
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
	updatedOCRInstance := t.ocrInstanceMap[answerAddress]
	onChainData, err := updatedOCRInstance.GetRound(context.Background(), newAnswer.RoundId)
	Expect(err).ShouldNot(HaveOccurred(), "Error retrieving on-chain data for '%s' at round '%d'", answerAddress, processedAnswer.UpdatedRoundId)
	processedAnswer.OnChainAnswer = int(onChainData.Answer.Int64())
	processedAnswer.OnChainRoundId = onChainData.RoundId.Uint64()

	return t.TestReporter.ContractReports[answerAddress].NewAnswerUpdated(processedAnswer)
}

// triggers a new OCR round by setting a new mock adapter value
func (t *OCRSoakTest) triggerNewRound(currentAdapterValue int) {
	startingBlockNum, err := t.chainClient.LatestBlockNumber(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Error retrieving latest block number")

	for _, report := range t.TestReporter.ContractReports {
		report.NewAnswerExpected(currentAdapterValue, startingBlockNum)
	}
	actions.SetAllAdapterResponsesToTheSameValue(currentAdapterValue, t.ocrInstances, t.chainlinkNodes, t.mockServer)()
	log.Info().
		Int("Value", currentAdapterValue).
		Msg("Starting a New OCR Round")
}

// ensureValues ensures that all values needed to run the test are present
func (t *OCRSoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.BlockchainClient).ShouldNot(BeNil(), "Need a valid blockchain client to use for the test")
	t.chainClient = inputs.BlockchainClient
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 OCR contract")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.ExpectedRoundTime).Should(BeNumerically(">=", time.Second*1), "Expected ExpectedRoundTime to be greater than 1 second")
	Expect(inputs.RoundTimeout).Should(BeNumerically(">=", inputs.ExpectedRoundTime), "Expected RoundTimeout to be greater than ExpectedRoundTime")
	Expect(inputs.TimeBetweenRounds).ShouldNot(BeNil(), "You forgot to set TimeBetweenRounds")
	Expect(inputs.TimeBetweenRounds).Should(BeNumerically("<", time.Hour), "TimeBetweenRounds must be less than 1 hour")
}

// subscribeToAnswerUpdatedEvent subscribes to the event log for AnswerUpdated event and
// verifies if the answer is matching with the expected value
func (t *OCRSoakTest) subscribeOCREvents(
	answerUpdated chan *ethereum.OffchainAggregatorAnswerUpdated,
) {
	contractABI, err := ethereum.OffchainAggregatorMetaData.GetAbi()
	Expect(err).ShouldNot(HaveOccurred(), "Getting contract abi for OCR shouldn't fail")
	latestBlockNum, err := t.chainClient.LatestBlockNumber(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OCR instance shouldn't fail")
	query := geth.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
		Addresses: []common.Address{},
	}
	for i := 0; i < len(t.ocrInstances); i++ {
		query.Addresses = append(query.Addresses, common.HexToAddress(t.ocrInstances[i].Address()))
	}
	eventLogs := make(chan types.Log)
	sub, err := t.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OCR instance shouldn't fail")

	go func() {
		defer GinkgoRecover()
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				log.Error().Err(err).Msg("Error while watching for new contract events. Retrying Subscription")
				sub.Unsubscribe()

				sub, err = t.chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OCR instance shouldn't fail")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				Expect(err).ShouldNot(HaveOccurred(), "Getting event details for OCR instances shouldn't fail")

				go t.processNewEvent(sub, answerUpdated, &vLog, eventDetails, t.ocrInstances[0], contractABI)
			}
		}
	}()
}
