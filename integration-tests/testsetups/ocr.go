// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

const (
	saveFileLocation    = "/persistence/ocr-soak-test-state.toml"
	interruptedExitCode = 3
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs                *OCRSoakTestInputs
	TestReporter          testreporters.OCRSoakTestReporter
	OperatorForwarderFlow bool

	t                *testing.T
	startTime        time.Time
	timeLeft         time.Duration
	startingBlockNum uint64
	testEnvironment  *environment.Environment
	namespace        string
	log              zerolog.Logger
	bootstrapNode    *client.ChainlinkK8sClient
	workerNodes      []*client.ChainlinkK8sClient
	chainClient      blockchain.EVMClient
	mockServer       *ctfClient.MockserverClient
	filterQuery      geth.FilterQuery

	ocrRoundStates []*testreporters.OCRRoundState
	testIssues     []*testreporters.TestIssue

	ocrInstances   []contracts.OffchainAggregator
	ocrInstanceMap map[string]contracts.OffchainAggregator // address : instance
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	TestDuration            time.Duration `envconfig:"TEST_DURATION" default:"15m"`         // How long to run the test for
	NumberOfContracts       int           `envconfig:"NUMBER_CONTRACTS" default:"2"`        // Number of OCR contracts to launch
	ChainlinkNodeFunding    float64       `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"` // Amount of native currency to fund each chainlink node with
	bigChainlinkNodeFunding *big.Float    // Convenience conversions for funding
	TimeBetweenRounds       time.Duration `envconfig:"TIME_BETWEEN_ROUNDS" default:"1m"` // How long to wait before starting a new round; controls frequency of rounds
}

func (i OCRSoakTestInputs) setForRemoteRunner() {
	os.Setenv("TEST_OCR_TEST_DURATION", i.TestDuration.String())
	os.Setenv("TEST_OCR_NUMBER_CONTRACTS", fmt.Sprint(i.NumberOfContracts))
	os.Setenv("TEST_OCR_CHAINLINK_NODE_FUNDING", strconv.FormatFloat(i.ChainlinkNodeFunding, 'f', -1, 64))
	os.Setenv("TEST_OCR_TIME_BETWEEN_ROUNDS", i.TimeBetweenRounds.String())

	selectedNetworks := strings.Split(os.Getenv("SELECTED_NETWORKS"), ",")
	for _, networkPrefix := range selectedNetworks {
		urlEnv := fmt.Sprintf("%s_URLS", networkPrefix)
		httpEnv := fmt.Sprintf("%s_HTTP_URLS", networkPrefix)
		os.Setenv(fmt.Sprintf("TEST_%s", urlEnv), os.Getenv(urlEnv))
		os.Setenv(fmt.Sprintf("TEST_%s", httpEnv), os.Getenv(httpEnv))
	}
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(t *testing.T, forwarderFlow bool) (*OCRSoakTest, error) {
	var testInputs OCRSoakTestInputs
	err := envconfig.Process("OCR", &testInputs)
	if err != nil {
		return nil, err
	}
	testInputs.setForRemoteRunner()

	test := &OCRSoakTest{
		Inputs:                &testInputs,
		OperatorForwarderFlow: forwarderFlow,
		TestReporter: testreporters.OCRSoakTestReporter{
			StartTime: time.Now(),
		},
		t:              t,
		startTime:      time.Now(),
		timeLeft:       testInputs.TestDuration,
		log:            logging.GetTestLogger(t),
		ocrRoundStates: make([]*testreporters.OCRRoundState, 0),
		ocrInstanceMap: make(map[string]contracts.OffchainAggregator),
	}
	return test, test.ensureInputValues()
}

// DeployEnvironment deploys the test environment, starting all Chainlink nodes and other components for the test
func (o *OCRSoakTest) DeployEnvironment(customChainlinkNetworkTOML string) {
	network := networks.SelectedNetwork // Environment currently being used to soak test on
	nsPre := "soak-ocr-"
	if o.OperatorForwarderFlow {
		nsPre = fmt.Sprintf("%sforwarder-", nsPre)
	}
	nsPre = fmt.Sprintf("%s%s", nsPre, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"))
	baseEnvironmentConfig := &environment.Config{
		TTL:                time.Hour * 720, // 30 days,
		NamespacePrefix:    nsPre,
		Test:               o.t,
		PreventPodEviction: true,
	}

	cd := chainlink.New(0, map[string]any{
		"replicas": 6,
		"toml":     client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customChainlinkNetworkTOML, network),
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
	})

	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})).
		AddHelm(cd)
	err := testEnvironment.Run()
	require.NoError(o.t, err, "Error launching test environment")
	o.testEnvironment = testEnvironment
	o.namespace = testEnvironment.Cfg.Namespace
}

// LoadEnvironment loads an existing test environment using the provided URLs
func (o *OCRSoakTest) LoadEnvironment(chainlinkURLs []string, chainURL, mockServerURL string) {
	var (
		network = networks.SelectedNetwork
		err     error
	)
	o.chainClient, err = blockchain.ConnectEVMClient(network, o.log)
	require.NoError(o.t, err, "Error connecting to EVM client")
	chainlinkNodes, err := client.ConnectChainlinkNodeURLs(chainlinkURLs)
	require.NoError(o.t, err, "Error connecting to chainlink nodes")
	o.bootstrapNode, o.workerNodes = chainlinkNodes[0], chainlinkNodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServerURL(mockServerURL)
	require.NoError(o.t, err, "Error connecting to mockserver")
}

// Environment returns the full K8s test environment
func (o *OCRSoakTest) Environment() *environment.Environment {
	return o.testEnvironment
}

func (o *OCRSoakTest) Setup() {
	var (
		err     error
		network = networks.SelectedNetwork
	)

	// Environment currently being used to soak test on
	// Make connections to soak test resources
	o.chainClient, err = blockchain.NewEVMClient(network, o.testEnvironment, o.log)
	require.NoError(o.t, err, "Error creating EVM client")
	contractDeployer, err := contracts.NewContractDeployer(o.chainClient, o.log)
	require.NoError(o.t, err, "Unable to create contract deployer")
	require.NotNil(o.t, contractDeployer, "Contract deployer shouldn't be nil")
	nodes, err := client.ConnectChainlinkNodes(o.testEnvironment)
	require.NoError(o.t, err, "Connecting to chainlink nodes shouldn't fail")
	o.bootstrapNode, o.workerNodes = nodes[0], nodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServer(o.testEnvironment)
	require.NoError(o.t, err, "Creating mockserver clients shouldn't fail")
	o.chainClient.ParallelTransactions(true)
	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(o.t, err, "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(o.workerNodes, o.chainClient, o.Inputs.bigChainlinkNodeFunding)
	require.NoError(o.t, err, "Error funding Chainlink nodes")

	if o.OperatorForwarderFlow {
		contractLoader, err := contracts.NewContractLoader(o.chainClient, o.log)
		require.NoError(o.t, err, "Loading contracts shouldn't fail")

		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
			o.t, contractDeployer, linkTokenContract, o.chainClient, len(o.workerNodes),
		)
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.workerNodes)
		require.NoError(o.t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range o.workerNodes {
			actions.AcceptAuthorizedReceiversOperator(
				o.t, operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, o.chainClient, contractLoader,
			)
			require.NoError(o.t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
			actions.TrackForwarder(o.t, o.chainClient, authorizedForwarders[i], o.workerNodes[i])
			err = o.chainClient.WaitForEvents()
		}

		o.ocrInstances = actions.DeployOCRContractsForwarderFlow(
			o.t,
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
		require.NoError(o.t, err)
	}

	err = o.chainClient.WaitForEvents()
	require.NoError(o.t, err, "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range o.ocrInstances {
		o.ocrInstanceMap[ocrInstance.Address()] = ocrInstance
	}
	o.log.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	latestBlockNum, err := o.chainClient.LatestBlockNumber(ctx)
	cancel()
	require.NoError(o.t, err, "Error getting current block number")
	o.startingBlockNum = latestBlockNum

	startingValue := 5
	if o.OperatorForwarderFlow {
		actions.CreateOCRJobsWithForwarder(o.t, o.ocrInstances, o.bootstrapNode, o.workerNodes, startingValue, o.mockServer, o.chainClient.GetChainID().String())
	} else {
		err := actions.CreateOCRJobs(o.ocrInstances, o.bootstrapNode, o.workerNodes, startingValue, o.mockServer, o.chainClient.GetChainID().String())
		require.NoError(o.t, err, "Error creating OCR jobs")
	}

	o.log.Info().
		Str("Test Duration", o.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Number of OCR Contracts", len(o.ocrInstances)).
		Msg("Starting OCR Soak Test")

	o.testLoop(o.Inputs.TestDuration, startingValue)
	o.complete()
}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	string,
	[]*client.ChainlinkK8sClient,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, o.namespace, append(o.workerNodes, o.bootstrapNode), &o.TestReporter, o.chainClient
}

// *********************
// Recovery if the test is shut-down/rebalanced by K8s
// *********************

// OCRSoakTestState contains all the info needed by the test to recover from a K8s rebalance, assuming the test was in a running state
type OCRSoakTestState struct {
	Namespace            string                         `toml:"namespace"`
	OCRRoundStates       []*testreporters.OCRRoundState `toml:"ocrRoundStates"`
	TestIssues           []*testreporters.TestIssue     `toml:"testIssues"`
	StartingBlockNum     uint64                         `toml:"startingBlockNum"`
	StartTime            time.Time                      `toml:"startTime"`
	TimeRunning          time.Duration                  `toml:"timeRunning"`
	TestDuration         time.Duration                  `toml:"testDuration"`
	OCRContractAddresses []string                       `toml:"ocrContractAddresses"`

	BootStrapNodeURL string   `toml:"bootstrapNodeURL"`
	WorkerNodeURLs   []string `toml:"workerNodeURLs"`
	ChainURL         string   `toml:"chainURL"`
	MockServerURL    string   `toml:"mockServerURL"`
}

// SaveState saves the current state of the test to a TOML file
func (o *OCRSoakTest) SaveState() error {
	ocrAddresses := make([]string, len(o.ocrInstances))
	for i, ocrInstance := range o.ocrInstances {
		ocrAddresses[i] = ocrInstance.Address()
	}
	workerNodeURLs := make([]string, len(o.workerNodes))
	for i, workerNode := range o.workerNodes {
		workerNodeURLs[i] = workerNode.URL()
	}

	testState := &OCRSoakTestState{
		Namespace:            o.namespace,
		OCRRoundStates:       o.ocrRoundStates,
		TestIssues:           o.testIssues,
		StartingBlockNum:     o.startingBlockNum,
		StartTime:            o.startTime,
		TimeRunning:          time.Since(o.startTime),
		TestDuration:         o.Inputs.TestDuration,
		OCRContractAddresses: ocrAddresses,

		ChainURL:         o.chainClient.GetNetworkConfig().URL,
		MockServerURL:    "http://mockserver:1080", // TODO: Make this dynamic
		BootStrapNodeURL: o.bootstrapNode.URL(),
		WorkerNodeURLs:   workerNodeURLs,
	}
	data, err := toml.Marshal(testState)
	if err != nil {
		return err
	}
	// #nosec G306 - let everyone read
	if err = os.WriteFile(saveFileLocation, data, 0644); err != nil {
		return err
	}
	fmt.Println("---Saved State---")
	fmt.Println(saveFileLocation)
	fmt.Println("-----------------")
	fmt.Println(string(data))
	fmt.Println("-----------------")
	return nil
}

// LoadState loads the test state from a TOML file
func (o *OCRSoakTest) LoadState() error {
	if !o.Interrupted() {
		return fmt.Errorf("no save file found at '%s'", saveFileLocation)
	}

	testState := &OCRSoakTestState{}
	saveData, err := os.ReadFile(saveFileLocation)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(saveData, testState)
	if err != nil {
		return err
	}
	fmt.Println("---Loaded State---")
	fmt.Println(saveFileLocation)
	fmt.Println("------------------")
	fmt.Println(string(saveData))
	fmt.Println("------------------")

	o.namespace = testState.Namespace
	o.TestReporter = testreporters.OCRSoakTestReporter{
		StartTime: testState.StartTime,
	}
	o.ocrRoundStates = testState.OCRRoundStates
	o.testIssues = testState.TestIssues
	o.Inputs.TestDuration = testState.TestDuration
	o.timeLeft = testState.TestDuration - testState.TimeRunning
	o.startTime = testState.StartTime
	o.startingBlockNum = testState.StartingBlockNum

	network := networks.SelectedNetwork
	o.chainClient, err = blockchain.ConnectEVMClient(network, o.log)
	if err != nil {
		return err
	}
	contractDeployer, err := contracts.NewContractDeployer(o.chainClient, o.log)
	if err != nil {
		return err
	}
	o.bootstrapNode, err = client.ConnectChainlinkNodeURL(testState.BootStrapNodeURL)
	if err != nil {
		return err
	}
	o.workerNodes, err = client.ConnectChainlinkNodeURLs(testState.WorkerNodeURLs)
	if err != nil {
		return err
	}

	o.ocrInstances = make([]contracts.OffchainAggregator, len(testState.OCRContractAddresses))
	for i, addr := range testState.OCRContractAddresses {
		address := common.HexToAddress(addr)
		instance, err := contractDeployer.LoadOffChainAggregator(&address)
		if err != nil {
			return err
		}
		o.ocrInstances[i] = instance
	}
	o.mockServer, err = ctfClient.ConnectMockServerURL(testState.MockServerURL)
	if err != nil {
		return err
	}

	return err
}

func (o *OCRSoakTest) Resume() {
	o.testIssues = append(o.testIssues, &testreporters.TestIssue{
		StartTime: time.Now(),
		Message:   "Test Resumed",
	})
	o.log.Info().
		Str("Total Duration", o.Inputs.TestDuration.String()).
		Str("Time Left", o.timeLeft.String()).
		Msg("Resuming OCR Soak Test")

	ocrAddresses := make([]common.Address, len(o.ocrInstances))
	for i, ocrInstance := range o.ocrInstances {
		ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
	}
	contractABI, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	require.NoError(o.t, err, "Error retrieving OCR contract ABI")
	o.filterQuery = geth.FilterQuery{
		Addresses: ocrAddresses,
		Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
		FromBlock: big.NewInt(0).SetUint64(o.startingBlockNum),
	}

	startingValue := 5
	o.testLoop(o.timeLeft, startingValue)

	o.log.Info().Msg("Test Complete, collecting on-chain events")

	err = o.collectEvents()
	o.log.Error().Err(err).Interface("Query", o.filterQuery).Msg("Error collecting on-chain events, expect malformed report")
	o.TestReporter.RecordEvents(o.ocrRoundStates, o.testIssues)
}

// Interrupted indicates whether the test was interrupted by something like a K8s rebalance or not
func (o *OCRSoakTest) Interrupted() bool {
	_, err := os.Stat(saveFileLocation)
	return err == nil
}

// *********************
// ****** Helpers ******
// *********************

// testLoop is the primary test loop that will trigger new rounds and watch events
func (o *OCRSoakTest) testLoop(testDuration time.Duration, newValue int) {
	endTest := time.After(testDuration)
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Kill, os.Interrupt, syscall.SIGTERM)
	lastValue := 0
	newRoundTrigger := time.NewTimer(0) // Want to trigger a new round ASAP
	defer newRoundTrigger.Stop()
	o.setFilterQuery()
	err := o.observeOCREvents()
	require.NoError(o.t, err, "Error subscribing to OCR events")

	for {
		select {
		case <-interruption:
			saveStart := time.Now()
			o.log.Warn().Msg("Test interrupted, saving state before shut down")
			o.testIssues = append(o.testIssues, &testreporters.TestIssue{
				StartTime: time.Now(),
				Message:   "Test Interrupted",
			})
			if err := o.SaveState(); err != nil {
				o.log.Error().Err(err).Msg("Error saving state")
			}
			o.log.Warn().Str("Time Taken", time.Since(saveStart).String()).Msg("Saved state")
			os.Exit(interruptedExitCode) // Exit with interrupted code to indicate test was interrupted, not just a normal failure
		case <-endTest:
			return
		case <-newRoundTrigger.C:
			err := o.triggerNewRound(newValue)
			timerReset := o.Inputs.TimeBetweenRounds
			if err != nil {
				timerReset = time.Second * 5
				o.log.Error().Err(err).
					Str("Waiting", timerReset.String()).
					Msg("Error triggering new round, waiting and trying again. Possible connection issues with mockserver")
			}
			newRoundTrigger.Reset(timerReset)

			// Change value for the next round
			newValue = rand.Intn(256) + 1 // #nosec G404 - not everything needs to be cryptographically secure
			for newValue == lastValue {
				newValue = rand.Intn(256) + 1 // #nosec G404 - kudos to you if you actually find a way to exploit this
			}
			lastValue = newValue
		case t := <-o.chainClient.ConnectionIssue():
			o.testIssues = append(o.testIssues, &testreporters.TestIssue{
				StartTime: t,
				Message:   "RPC Connection Lost",
			})
		case t := <-o.chainClient.ConnectionRestored():
			o.testIssues = append(o.testIssues, &testreporters.TestIssue{
				StartTime: t,
				Message:   "RPC Connection Restored",
			})
		}
	}
}

// completes the test
func (o *OCRSoakTest) complete() {
	o.log.Info().Msg("Test Complete, collecting on-chain events")

	err := o.collectEvents()
	if err != nil {
		o.log.Error().Err(err).Interface("Query", o.filterQuery).Msg("Error collecting on-chain events, expect malformed report")
	}
	o.TestReporter.RecordEvents(o.ocrRoundStates, o.testIssues)
}

// setFilterQuery to look for all events that happened
func (o *OCRSoakTest) setFilterQuery() {
	ocrAddresses := make([]common.Address, len(o.ocrInstances))
	for i, ocrInstance := range o.ocrInstances {
		ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
	}
	contractABI, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	require.NoError(o.t, err, "Error retrieving OCR contract ABI")
	o.filterQuery = geth.FilterQuery{
		Addresses: ocrAddresses,
		Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
		FromBlock: big.NewInt(0).SetUint64(o.startingBlockNum),
	}
	o.log.Debug().
		Interface("Addresses", ocrAddresses).
		Str("Topic", contractABI.Events["AnswerUpdated"].ID.Hex()).
		Uint64("Starting Block", o.startingBlockNum).
		Msg("Filter Query Set")
}

// observeOCREvents subscribes to OCR events and logs them to the test logger
// WARNING: Should only be used for observation and logging. This is not a reliable way to collect events.
func (o *OCRSoakTest) observeOCREvents() error {
	eventLogs := make(chan types.Log)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventSub, err := o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
	cancel()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-eventLogs:
				answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
				if err != nil {
					o.log.Warn().
						Err(err).
						Str("Address", event.Address.Hex()).
						Uint64("Block Number", event.BlockNumber).
						Msg("Error parsing event as AnswerUpdated")
					continue
				}
				o.log.Info().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Answer Updated Event")
			case err = <-eventSub.Err():
				backoff := time.Second
				for err != nil {
					o.log.Info().
						Err(err).
						Str("Backoff", backoff.String()).
						Interface("Query", o.filterQuery).
						Msg("Error while subscribed to OCR Logs. Resubscribing")
					ctx, cancel = context.WithTimeout(context.Background(), backoff)
					eventSub, err = o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
					cancel()
					if err != nil {
						time.Sleep(backoff)
						backoff = time.Duration(math.Min(float64(backoff)*2, float64(30*time.Second)))
					}
				}
			}
		}
	}()

	return nil
}

// triggers a new OCR round by setting a new mock adapter value
func (o *OCRSoakTest) triggerNewRound(newValue int) error {
	if len(o.ocrRoundStates) > 0 {
		o.ocrRoundStates[len(o.ocrRoundStates)-1].EndTime = time.Now()
	}

	err := actions.SetAllAdapterResponsesToTheSameValue(newValue, o.ocrInstances, o.workerNodes, o.mockServer)
	if err != nil {
		return err
	}

	expectedState := &testreporters.OCRRoundState{
		StartTime:   time.Now(),
		Answer:      int64(newValue),
		FoundEvents: make(map[string][]*testreporters.FoundEvent),
	}
	for _, ocrInstance := range o.ocrInstances {
		expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
	}
	o.ocrRoundStates = append(o.ocrRoundStates, expectedState)
	o.log.Info().
		Int("Value", newValue).
		Msg("Starting a New OCR Round")
	return nil
}

func (o *OCRSoakTest) collectEvents() error {
	start := time.Now()
	if len(o.ocrRoundStates) == 0 {
		return fmt.Errorf("error collecting on-chain events, no rounds have been started")
	}
	o.ocrRoundStates[len(o.ocrRoundStates)-1].EndTime = start // Set end time for last expected event
	o.log.Info().Msg("Collecting on-chain events")

	// We must retrieve the events, use exponential backoff for timeout to retry
	timeout := time.Second * 15
	o.log.Info().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Retrieving on-chain events")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	contractEvents, err := o.chainClient.FilterLogs(ctx, o.filterQuery)
	cancel()
	for err != nil {
		o.log.Info().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Retrieving on-chain events")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		contractEvents, err = o.chainClient.FilterLogs(ctx, o.filterQuery)
		cancel()
		if err != nil {
			o.log.Warn().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Error collecting on-chain events, trying again")
			timeout *= 2
		}
	}

	sortedFoundEvents := make([]*testreporters.FoundEvent, 0)
	for _, event := range contractEvents {
		answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
		if err != nil {
			return fmt.Errorf("error parsing EventAnswerUpdated for event: %v, %w", event, err)
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
		if !event.StartTime.Before(o.ocrRoundStates[expectedIndex].EndTime) {
			expectedIndex++
			if expectedIndex >= len(o.ocrRoundStates) {
				o.log.Warn().
					Str("Event Time", event.StartTime.String()).
					Str("Expected End Time", o.ocrRoundStates[expectedIndex].EndTime.String()).
					Msg("Found events after last expected end time, adding event to that final report, things might be weird")
			}
		}
		o.ocrRoundStates[expectedIndex].FoundEvents[event.Address] = append(o.ocrRoundStates[expectedIndex].FoundEvents[event.Address], event)
		o.ocrRoundStates[expectedIndex].TimeLineEvents = append(o.ocrRoundStates[expectedIndex].TimeLineEvents, event)
	}

	o.log.Info().
		Str("Time", time.Since(start).String()).
		Msg("Collected on-chain events")
	return nil
}

// ensureValues ensures that all values needed to run the test are present
func (o *OCRSoakTest) ensureInputValues() error {
	inputs := o.Inputs
	if inputs.NumberOfContracts <= 0 {
		return fmt.Errorf("Number of OCR contracts must be greater than 0, found %d", inputs.NumberOfContracts)
	}
	if inputs.ChainlinkNodeFunding <= 0 {
		return fmt.Errorf("Chainlink node funding must be greater than 0, found %f", inputs.ChainlinkNodeFunding)
	}
	if inputs.TestDuration <= time.Minute {
		return fmt.Errorf("Test duration must be greater than 1 minute, found %s", inputs.TestDuration.String())
	}
	if inputs.TimeBetweenRounds >= time.Hour {
		return fmt.Errorf("Time between rounds must be less than 1 hour, found %s", inputs.TimeBetweenRounds.String())
	}
	if inputs.TimeBetweenRounds < time.Second*30 {
		return fmt.Errorf("Time between rounds must be greater or equal to 30 seconds, found %s", inputs.TimeBetweenRounds.String())
	}
	o.Inputs.bigChainlinkNodeFunding = big.NewFloat(inputs.ChainlinkNodeFunding)
	return nil
}
