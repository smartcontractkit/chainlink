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
	"strings"
	"syscall"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	tt "github.com/smartcontractkit/chainlink/integration-tests/types"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

const (
	saveFileLocation    = "/persistence/ocr-soak-test-state.toml"
	interruptedExitCode = 3
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Config                *tc.TestConfig
	TestReporter          testreporters.OCRSoakTestReporter
	OperatorForwarderFlow bool
	seth                  *seth.Client

	t                *testing.T
	startTime        time.Time
	timeLeft         time.Duration
	startingBlockNum uint64
	testEnvironment  *environment.Environment
	namespace        string
	log              zerolog.Logger
	bootstrapNode    *client.ChainlinkK8sClient
	workerNodes      []*client.ChainlinkK8sClient
	mockServer       *ctfClient.MockserverClient
	filterQuery      geth.FilterQuery

	ocrRoundStates []*testreporters.OCRRoundState
	testIssues     []*testreporters.TestIssue

	ocrV1Instances   []contracts.OffchainAggregator
	ocrV1InstanceMap map[string]contracts.OffchainAggregator // address : instance

	ocrV2Instances   []contracts.OffchainAggregatorV2
	ocrV2InstanceMap map[string]contracts.OffchainAggregatorV2 // address : instance
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(t *testing.T, config *tc.TestConfig, forwarderFlow bool) (*OCRSoakTest, error) {
	test := &OCRSoakTest{
		Config:                config,
		OperatorForwarderFlow: forwarderFlow,
		TestReporter: testreporters.OCRSoakTestReporter{
			OCRVersion: *config.OCR.Soak.OCRVersion,
			StartTime:  time.Now(),
		},
		t:                t,
		startTime:        time.Now(),
		timeLeft:         config.OCR.Common.TestDuration.Duration,
		log:              logging.GetTestLogger(t),
		ocrRoundStates:   make([]*testreporters.OCRRoundState, 0),
		ocrV1InstanceMap: make(map[string]contracts.OffchainAggregator),
		ocrV2InstanceMap: make(map[string]contracts.OffchainAggregatorV2),
	}
	return test, test.ensureInputValues()
}

// DeployEnvironment deploys the test environment, starting all Chainlink nodes and other components for the test
func (o *OCRSoakTest) DeployEnvironment(customChainlinkNetworkTOML string, ocrTestConfig tt.OcrTestConfig) {
	network := networks.MustGetSelectedNetworkConfig(ocrTestConfig.GetNetworkConfig())[0] // Environment currently being used to soak test on
	nsPre := fmt.Sprintf("soak-ocr-v%s-", *ocrTestConfig.GetOCRConfig().Soak.OCRVersion)
	if o.OperatorForwarderFlow {
		nsPre = fmt.Sprintf("%sforwarder-", nsPre)
	}
	nsPre = fmt.Sprintf("%s%s", nsPre, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"))
	nsPre = strings.ReplaceAll(nsPre, "_", "-")
	baseEnvironmentConfig := &environment.Config{
		TTL:                time.Hour * 720, // 30 days,
		NamespacePrefix:    nsPre,
		Test:               o.t,
		PreventPodEviction: true,
	}

	var conf string
	if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "1" {
		conf = config.BaseOCR1Config
	} else if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "2" {
		conf = config.BaseOCR2Config
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(ocrTestConfig.GetChainlinkImageConfig(), target)
		ctf_config.MightConfigOverridePyroscopeKey(ocrTestConfig.GetPyroscopeConfig(), target)
	}

	cd := chainlink.NewWithOverride(0, map[string]any{
		"replicas": 6,
		"toml":     networks.AddNetworkDetailedConfig(conf, ocrTestConfig.GetPyroscopeConfig(), customChainlinkNetworkTOML, network),
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
	}, ocrTestConfig.GetChainlinkImageConfig(), overrideFn)

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

// Environment returns the full K8s test environment
func (o *OCRSoakTest) Environment() *environment.Environment {
	return o.testEnvironment
}

func (o *OCRSoakTest) Setup(ocrTestConfig tt.OcrTestConfig) {
	var (
		err     error
		network = networks.MustGetSelectedNetworkConfig(ocrTestConfig.GetNetworkConfig())[0]
	)

	network = utils.MustReplaceSimulatedNetworkUrlWithK8(o.log, network, *o.testEnvironment)
	seth, err := actions_seth.GetChainClient(o.Config, network)
	require.NoError(o.t, err, "Error creating seth client")

	o.seth = seth

	nodes, err := client.ConnectChainlinkNodes(o.testEnvironment)
	require.NoError(o.t, err, "Connecting to chainlink nodes shouldn't fail")
	o.bootstrapNode, o.workerNodes = nodes[0], nodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServer(o.testEnvironment)
	require.NoError(o.t, err, "Creating mockserver clients shouldn't fail")

	linkContract, err := contracts.DeployLinkTokenContract(o.log, seth)
	require.NoError(o.t, err, "Error deploying LINK contract")

	// Fund Chainlink nodes, excluding the bootstrap node
	o.log.Info().Float64("ETH amount per node", *o.Config.Common.ChainlinkNodeFunding).Msg("Funding Chainlink nodes")
	err = actions_seth.FundChainlinkNodesFromRootAddress(o.log, seth, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes), big.NewFloat(*o.Config.Common.ChainlinkNodeFunding))
	require.NoError(o.t, err, "Error funding Chainlink nodes")

	var forwarders []common.Address

	if o.OperatorForwarderFlow {
		var operators []common.Address
		operators, forwarders, _ = actions_seth.DeployForwarderContracts(
			o.t, o.seth, common.HexToAddress(linkContract.Address()), len(o.workerNodes),
		)
		require.Equal(o.t, len(o.workerNodes), len(operators), "Number of operators should match number of nodes")
		require.Equal(o.t, len(o.workerNodes), len(forwarders), "Number of authorized forwarders should match number of nodes")
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.workerNodes)
		require.NoError(o.t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range o.workerNodes {
			actions_seth.AcceptAuthorizedReceiversOperator(
				o.t, o.log, o.seth, operators[i], forwarders[i], []common.Address{forwarderNodesAddresses[i]})
			require.NoError(o.t, err, "Accepting Authorize Receivers on Operator shouldn't fail")

			actions_seth.TrackForwarder(o.t, o.seth, forwarders[i], o.workerNodes[i])
		}
	} else if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "1" {
		if o.OperatorForwarderFlow {
			o.ocrV1Instances, err = actions_seth.DeployOCRContractsForwarderFlow(
				o.log,
				o.seth,
				*o.Config.OCR.Soak.NumberOfContracts,
				common.HexToAddress(linkContract.Address()),
				contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes),
				forwarders,
			)
			require.NoError(o.t, err, "Error deploying OCR Forwarder contracts")
		} else {
			o.ocrV1Instances, err = actions_seth.DeployOCRv1Contracts(
				o.log,
				seth,
				*o.Config.OCR.Soak.NumberOfContracts,
				common.HexToAddress(linkContract.Address()),
				contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes),
			)
			require.NoError(o.t, err)
		}
	} else if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "2" {
		var transmitters []string

		if o.OperatorForwarderFlow {
			for _, forwarder := range forwarders {
				transmitters = append(transmitters, forwarder.Hex())
			}
		} else {
			for _, node := range o.workerNodes {
				nodeAddress, err := node.PrimaryEthAddress()
				require.NoError(o.t, err, "Error getting node's primary ETH address")
				transmitters = append(transmitters, nodeAddress)
			}
		}

		ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
		o.ocrV2Instances, err = actions_seth.DeployOCRv2Contracts(
			o.log,
			o.seth,
			*ocrTestConfig.GetOCRConfig().Soak.NumberOfContracts,
			common.HexToAddress(linkContract.Address()),
			transmitters,
			ocrOffchainOptions,
		)
		require.NoError(o.t, err, "Error deploying OCRv2 contracts")
		contractConfig, err := actions.BuildMedianOCR2Config(o.workerNodes, ocrOffchainOptions)
		require.NoError(o.t, err, "Error building median config")
		err = actions_seth.ConfigureOCRv2AggregatorContracts(contractConfig, o.ocrV2Instances)
		require.NoError(o.t, err, "Error configuring OCRv2 aggregator contracts")
	}

	if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "1" {
		for _, ocrInstance := range o.ocrV1Instances {
			o.ocrV1InstanceMap[ocrInstance.Address()] = ocrInstance
		}
	} else if *ocrTestConfig.GetOCRConfig().Soak.OCRVersion == "2" {
		for _, ocrInstance := range o.ocrV2Instances {
			o.ocrV2InstanceMap[ocrInstance.Address()] = ocrInstance
		}
	}

	o.log.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run() {
	config, err := tc.GetConfig("soak", tc.OCR)
	require.NoError(o.t, err, "Error getting config")

	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), time.Second*5)
	latestBlockNum, err := o.seth.Client.BlockNumber(ctx)
	cancel()
	require.NoError(o.t, err, "Error getting current block number")
	o.startingBlockNum = latestBlockNum

	startingValue := 5
	if o.OperatorForwarderFlow {
		actions.CreateOCRJobsWithForwarder(o.t, o.ocrV1Instances, o.bootstrapNode, o.workerNodes, startingValue, o.mockServer, o.seth.ChainID)
	} else if *config.OCR.Soak.OCRVersion == "1" {
		ctx, cancel := context.WithTimeout(testcontext.Get(o.t), time.Second*5)
		chainId, err := o.seth.Client.ChainID(ctx)
		cancel()
		require.NoError(o.t, err, "Error getting chain ID")
		err = actions.CreateOCRJobs(o.ocrV1Instances, o.bootstrapNode, o.workerNodes, startingValue, o.mockServer, chainId.String())
		require.NoError(o.t, err, "Error creating OCR jobs")
	} else if *config.OCR.Soak.OCRVersion == "2" {
		err := actions.CreateOCRv2Jobs(o.ocrV2Instances, o.bootstrapNode, o.workerNodes, o.mockServer, startingValue, o.seth.ChainID, o.OperatorForwarderFlow)
		require.NoError(o.t, err, "Error creating OCR jobs")
	}

	o.log.Info().
		Str("Test Duration", o.Config.OCR.Common.TestDuration.Duration.Truncate(time.Second).String()).
		Int("Number of OCR Contracts", *config.OCR.Soak.NumberOfContracts).
		Str("OCR Version", *config.OCR.Soak.OCRVersion).
		Msg("Starting OCR Soak Test")

	o.testLoop(o.Config.OCR.Common.TestDuration.Duration, startingValue)
	o.complete()
}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	*seth.Client,
	string,
	[]*client.ChainlinkK8sClient,
	reportModel.TestReporter,
	reportModel.GrafanaURLProvider,
) {
	return t, o.seth, o.namespace, append(o.workerNodes, o.bootstrapNode), &o.TestReporter, o.Config
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
	OCRVersion           string                         `toml:"ocrVersion"`

	BootStrapNodeURL string   `toml:"bootstrapNodeURL"`
	WorkerNodeURLs   []string `toml:"workerNodeURLs"`
	ChainURL         string   `toml:"chainURL"`
	MockServerURL    string   `toml:"mockServerURL"`
}

// SaveState saves the current state of the test to a TOML file
func (o *OCRSoakTest) SaveState() error {
	ocrAddresses := o.getContractAddressesString()
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
		TestDuration:         o.Config.OCR.Common.TestDuration.Duration,
		OCRContractAddresses: ocrAddresses,
		OCRVersion:           *o.Config.OCR.Soak.OCRVersion,

		MockServerURL:    "http://mockserver:1080", // TODO: Make this dynamic
		BootStrapNodeURL: o.bootstrapNode.URL(),
		WorkerNodeURLs:   workerNodeURLs,
	}
	data, err := toml.Marshal(testState)
	if err != nil {
		return err
	}
	//nolint:gosec // G306 - let everyone read
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
		OCRVersion: testState.OCRVersion,
		StartTime:  testState.StartTime,
	}
	duration := blockchain.StrDuration{Duration: testState.TestDuration}
	o.ocrRoundStates = testState.OCRRoundStates
	o.testIssues = testState.TestIssues
	o.Config.OCR.Common.TestDuration = &duration
	o.timeLeft = testState.TestDuration - testState.TimeRunning
	o.startTime = testState.StartTime
	o.startingBlockNum = testState.StartingBlockNum
	o.Config.OCR.Soak.OCRVersion = &testState.OCRVersion

	o.bootstrapNode, err = client.ConnectChainlinkNodeURL(testState.BootStrapNodeURL)
	if err != nil {
		return err
	}
	o.workerNodes, err = client.ConnectChainlinkNodeURLs(testState.WorkerNodeURLs)
	if err != nil {
		return err
	}

	if testState.OCRVersion == "1" {
		o.ocrV1Instances = make([]contracts.OffchainAggregator, len(testState.OCRContractAddresses))
		for i, addr := range testState.OCRContractAddresses {
			instance, err := contracts.LoadOffchainAggregator(o.log, o.seth, common.HexToAddress(addr))
			if err != nil {
				return fmt.Errorf("failed to instantiate OCR instance: %w", err)
			}
			o.ocrV1Instances[i] = &instance
		}
	} else if testState.OCRVersion == "2" {
		o.ocrV2Instances = make([]contracts.OffchainAggregatorV2, len(testState.OCRContractAddresses))
		for i, addr := range testState.OCRContractAddresses {
			instance, err := contracts.LoadOffChainAggregatorV2(o.log, o.seth, common.HexToAddress(addr))
			if err != nil {
				return err
			}
			o.ocrV2Instances[i] = &instance
		}
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
		Str("Total Duration", o.Config.OCR.Common.TestDuration.String()).
		Str("Time Left", o.timeLeft.String()).
		Msg("Resuming OCR Soak Test")

	ocrAddresses := make([]common.Address, *o.Config.OCR.Soak.NumberOfContracts)

	if *o.Config.OCR.Soak.OCRVersion == "1" {
		for i, ocrInstance := range o.ocrV1Instances {
			ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
		}
		contractABI, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
		require.NoError(o.t, err, "Error retrieving OCR contract ABI")
		o.filterQuery = geth.FilterQuery{
			Addresses: ocrAddresses,
			Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
			FromBlock: big.NewInt(0).SetUint64(o.startingBlockNum),
		}
	} else if *o.Config.OCR.Soak.OCRVersion == "2" {
		for i, ocrInstance := range o.ocrV2Instances {
			ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
		}
		contractABI, err := ocr2aggregator.AggregatorInterfaceMetaData.GetAbi()
		require.NoError(o.t, err, "Error retrieving OCR contract ABI")
		o.filterQuery = geth.FilterQuery{
			Addresses: ocrAddresses,
			Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
			FromBlock: big.NewInt(0).SetUint64(o.startingBlockNum),
		}
	}

	startingValue := 5
	o.testLoop(o.timeLeft, startingValue)

	o.log.Info().Msg("Test Complete, collecting on-chain events")

	err := o.collectEvents()
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
	//nolint:staticcheck //ignore SA1016 we need to send the os.Kill signal
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
			timerReset := o.Config.OCR.Soak.TimeBetweenRounds.Duration
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
	ocrAddresses := o.getContractAddresses()
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
	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), 5*time.Second)
	eventSub, err := o.seth.Client.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
	cancel()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-eventLogs:
				if *o.Config.OCR.Soak.OCRVersion == "1" {
					answerUpdated, err := o.ocrV1Instances[0].ParseEventAnswerUpdated(event)
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
				} else if *o.Config.OCR.Soak.OCRVersion == "2" {
					answerUpdated, err := o.ocrV2Instances[0].ParseEventAnswerUpdated(event)
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
				}
			case err = <-eventSub.Err():
				backoff := time.Second
				for err != nil {
					o.log.Info().
						Err(err).
						Str("Backoff", backoff.String()).
						Interface("Query", o.filterQuery).
						Msg("Error while subscribed to OCR Logs. Resubscribing")
					ctx, cancel = context.WithTimeout(testcontext.Get(o.t), backoff)
					eventSub, err = o.seth.Client.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
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

	var err error
	if *o.Config.OCR.Soak.OCRVersion == "1" {
		err = actions.SetAllAdapterResponsesToTheSameValue(newValue, o.ocrV1Instances, o.workerNodes, o.mockServer)
	} else if *o.Config.OCR.Soak.OCRVersion == "2" {
		err = actions.SetOCR2AllAdapterResponsesToTheSameValue(newValue, o.ocrV2Instances, o.workerNodes, o.mockServer)
	}
	if err != nil {
		return err
	}

	expectedState := &testreporters.OCRRoundState{
		StartTime:   time.Now(),
		Answer:      int64(newValue),
		FoundEvents: make(map[string][]*testreporters.FoundEvent),
	}
	if *o.Config.OCR.Soak.OCRVersion == "1" {
		for _, ocrInstance := range o.ocrV1Instances {
			expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
		}
	} else if *o.Config.OCR.Soak.OCRVersion == "2" {
		for _, ocrInstance := range o.ocrV2Instances {
			expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
		}
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

	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), timeout)
	contractEvents, err := o.seth.Client.FilterLogs(ctx, o.filterQuery)
	cancel()
	for err != nil {
		o.log.Info().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Retrieving on-chain events")
		ctx, cancel := context.WithTimeout(testcontext.Get(o.t), timeout)
		contractEvents, err = o.seth.Client.FilterLogs(ctx, o.filterQuery)
		cancel()
		if err != nil {
			o.log.Warn().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Error collecting on-chain events, trying again")
			timeout *= 2
		}
	}

	sortedFoundEvents := make([]*testreporters.FoundEvent, 0)
	for _, event := range contractEvents {
		if *o.Config.OCR.Soak.OCRVersion == "1" {
			answerUpdated, err := o.ocrV1Instances[0].ParseEventAnswerUpdated(event)
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
		} else if *o.Config.OCR.Soak.OCRVersion == "2" {
			answerUpdated, err := o.ocrV2Instances[0].ParseEventAnswerUpdated(event)
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
	ocrConfig := o.Config.OCR.Soak
	if *ocrConfig.OCRVersion != "1" && *ocrConfig.OCRVersion != "2" {
		return fmt.Errorf("OCR version must be 1 or 2, found %s", *ocrConfig.OCRVersion)
	}
	if ocrConfig.NumberOfContracts != nil && *ocrConfig.NumberOfContracts <= 0 {
		return fmt.Errorf("Number of OCR contracts must be set and greater than 0, found %d", ocrConfig.NumberOfContracts)
	}
	if o.Config.Common.ChainlinkNodeFunding != nil && *o.Config.Common.ChainlinkNodeFunding <= 0 {
		return fmt.Errorf("Chainlink node funding must be greater than 0, found %f", *o.Config.Common.ChainlinkNodeFunding)
	}
	if o.Config.OCR.Common.TestDuration != nil && o.Config.OCR.Common.TestDuration.Duration <= time.Minute {
		return fmt.Errorf("Test duration must be greater than 1 minute, found %s", o.Config.OCR.Common.TestDuration)
	}
	if ocrConfig.TimeBetweenRounds != nil && ocrConfig.TimeBetweenRounds.Duration >= time.Hour {
		return fmt.Errorf("Time between rounds must be less than 1 hour, found %s", ocrConfig.TimeBetweenRounds)
	}
	if ocrConfig.TimeBetweenRounds != nil && ocrConfig.TimeBetweenRounds.Duration < time.Second*30 {
		return fmt.Errorf("Time between rounds must be greater or equal to 30 seconds, found %s", ocrConfig.TimeBetweenRounds)
	}

	return nil
}

// getContractAddressesString returns the addresses of all OCR contracts deployed as a string slice
func (o *OCRSoakTest) getContractAddressesString() []string {
	contractAddresses := []string{}
	if len(o.ocrV1Instances) != 0 {
		for _, ocrInstance := range o.ocrV1Instances {
			contractAddresses = append(contractAddresses, ocrInstance.Address())
		}
	} else if len(o.ocrV2Instances) != 0 {
		if len(o.ocrV2Instances) != 0 {
			for _, ocrInstance := range o.ocrV2Instances {
				contractAddresses = append(contractAddresses, ocrInstance.Address())
			}
		}
	}

	return contractAddresses
}

// getContractAddresses returns the addresses of all OCR contracts deployed
func (o *OCRSoakTest) getContractAddresses() []common.Address {
	contractAddresses := []common.Address{}
	if len(o.ocrV1Instances) != 0 {
		for _, ocrInstance := range o.ocrV1Instances {
			contractAddresses = append(contractAddresses, common.HexToAddress(ocrInstance.Address()))
		}
	} else if len(o.ocrV2Instances) != 0 {
		if len(o.ocrV2Instances) != 0 {
			for _, ocrInstance := range o.ocrV2Instances {
				contractAddresses = append(contractAddresses, common.HexToAddress(ocrInstance.Address()))
			}
		}
	}

	return contractAddresses
}
