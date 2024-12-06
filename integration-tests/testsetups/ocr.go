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
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/grafana"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/foundry"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	tt "github.com/smartcontractkit/chainlink/integration-tests/types"
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
	sethClient            *seth.Client
	OCRVersion            string

	t                *testing.T
	startTime        time.Time
	timeLeft         time.Duration
	startingBlockNum uint64
	startingValue    int
	testEnvironment  *environment.Environment
	namespace        string
	log              zerolog.Logger
	bootstrapNode    *nodeclient.ChainlinkK8sClient
	workerNodes      []*nodeclient.ChainlinkK8sClient
	mockServer       *ctf_client.MockserverClient
	filterQuery      geth.FilterQuery

	ocrRoundStates []*testreporters.OCRRoundState
	testIssues     []*testreporters.TestIssue

	ocrV1Instances   []contracts.OffchainAggregator
	ocrV1InstanceMap map[string]contracts.OffchainAggregator // address : instance

	ocrV2Instances   []contracts.OffchainAggregatorV2
	ocrV2InstanceMap map[string]contracts.OffchainAggregatorV2 // address : instance

	linkContract *contracts.EthereumLinkToken

	rpcNetwork                 blockchain.EVMNetwork // network configuration for the blockchain node
	reorgHappened              bool                  // flag to indicate if a reorg happened during the test
	gasSpikeSimulationHappened bool                  // flag to indicate if a gas spike simulation happened during the test
	gasLimitSimulationHappened bool                  // flag to indicate if a gas limit simulation happened during the test
	chaosList                  []*havoc.Chaos        // list of chaos simulations to run during the test
}

type OCRSoakTestOption = func(c *OCRSoakTest)

func WithChaos(chaosList []*havoc.Chaos) OCRSoakTestOption {
	return func(c *OCRSoakTest) {
		c.chaosList = chaosList
	}
}

func WithNamespace(ns string) OCRSoakTestOption {
	return func(c *OCRSoakTest) {
		c.namespace = ns
	}
}

func WithForwarderFlow(forwarderFlow bool) OCRSoakTestOption {
	return func(c *OCRSoakTest) {
		c.OperatorForwarderFlow = forwarderFlow
	}
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(t *testing.T, config *tc.TestConfig, opts ...OCRSoakTestOption) (*OCRSoakTest, error) {
	test := &OCRSoakTest{
		Config: config,
		TestReporter: testreporters.OCRSoakTestReporter{
			StartTime: time.Now(),
		},
		t:                t,
		startTime:        time.Now(),
		timeLeft:         config.GetActiveOCRConfig().Common.TestDuration.Duration,
		log:              logging.GetTestLogger(t),
		ocrRoundStates:   make([]*testreporters.OCRRoundState, 0),
		ocrV1InstanceMap: make(map[string]contracts.OffchainAggregator),
		ocrV2InstanceMap: make(map[string]contracts.OffchainAggregatorV2),
	}

	ocrVersion := "1"
	if config.OCR2 != nil {
		ocrVersion = "2"
	}

	test.TestReporter.OCRVersion = ocrVersion
	test.OCRVersion = ocrVersion

	for _, opt := range opts {
		opt(test)
	}
	t.Cleanup(func() {
		test.deleteChaosSimulations()
	})
	return test, test.ensureInputValues()
}

// DeployEnvironment deploys the test environment, starting all Chainlink nodes and other components for the test
func (o *OCRSoakTest) DeployEnvironment(ocrTestConfig tt.OcrTestConfig) {
	nodeNetwork := networks.MustGetSelectedNetworkConfig(ocrTestConfig.GetNetworkConfig())[0] // Environment currently being used to soak test on

	nsPre := fmt.Sprintf("soak-ocr-v%s-", o.OCRVersion)
	if o.OperatorForwarderFlow {
		nsPre = fmt.Sprintf("%sforwarder-", nsPre)
	}

	nsPre = fmt.Sprintf("%s%s", nsPre, strings.ReplaceAll(strings.ToLower(nodeNetwork.Name), " ", "-"))
	nsPre = strings.ReplaceAll(nsPre, "_", "-")

	productName := fmt.Sprintf("data-feedsv%s.0", o.OCRVersion)
	nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(productName, "soak")
	require.NoError(o.t, err, "Error creating required chain.link labels for namespace")

	workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(productName, "soak")
	require.NoError(o.t, err, "Error creating required chain.link labels for workloads and pods")

	baseEnvironmentConfig := &environment.Config{
		TTL:                time.Hour * 720, // 30 days,
		NamespacePrefix:    nsPre,
		Test:               o.t,
		PreventPodEviction: true,
		Labels:             nsLabels,
		WorkloadLabels:     workloadPodLabels,
		PodLabels:          workloadPodLabels,
	}

	testEnv := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil))

	var anvilChart *foundry.Chart
	if nodeNetwork.Name == "Anvil" {
		anvilConfig := ocrTestConfig.GetNetworkConfig().AnvilConfigs["ANVIL"]
		anvilChart = foundry.New(&foundry.Props{
			Values: map[string]interface{}{
				"fullnameOverride": "anvil",
				"anvil": map[string]interface{}{
					"chainId":                   fmt.Sprintf("%d", nodeNetwork.ChainID),
					"blockTime":                 anvilConfig.BlockTime,
					"forkURL":                   anvilConfig.URL,
					"forkBlockNumber":           anvilConfig.BlockNumber,
					"forkRetries":               anvilConfig.Retries,
					"forkTimeout":               anvilConfig.Timeout,
					"forkComputeUnitsPerSecond": anvilConfig.ComputePerSecond,
					"forkNoRateLimit":           anvilConfig.RateLimitDisabled,
				},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "4",
						"memory": "6Gi",
					},
					"limits": map[string]interface{}{
						"cpu":    "4",
						"memory": "6Gi",
					},
				},
			},
		})
		testEnv.AddHelm(anvilChart)
		nodeNetwork.URLs = []string{anvilChart.ClusterWSURL}
		nodeNetwork.HTTPURLs = []string{anvilChart.ClusterHTTPURL}
	} else {
		testEnv.AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: nodeNetwork.Name,
			Simulated:   nodeNetwork.Simulated,
			WsURLs:      nodeNetwork.URLs,
		}))
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(ocrTestConfig.GetChainlinkImageConfig(), target)
		ctf_config.MightConfigOverridePyroscopeKey(ocrTestConfig.GetPyroscopeConfig(), target)
	}

	tomlConfig, err := actions.BuildTOMLNodeConfigForK8s(ocrTestConfig, nodeNetwork)
	require.NoError(o.t, err, "Error building TOML config for Chainlink nodes")

	cd := chainlink.NewWithOverride(0, map[string]any{
		"replicas": 6,
		"toml":     tomlConfig,
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
		"prometheus": true,
	}, ocrTestConfig.GetChainlinkImageConfig(), overrideFn)
	testEnv.AddHelm(cd)

	err = testEnv.Run()
	require.NoError(o.t, err, "Error launching test environment")
	o.testEnvironment = testEnv
	o.namespace = testEnv.Cfg.Namespace

	// If the test is using the remote runner, we don't need to set the network URLs
	// as the remote runner will handle that
	if o.Environment().WillUseRemoteRunner() {
		return
	}

	o.rpcNetwork = nodeNetwork
	if o.rpcNetwork.Simulated && o.rpcNetwork.Name == "Anvil" {
		if testEnv.Cfg.InsideK8s {
			// Test is running inside K8s, set the cluster URL of Anvil blockchain node
			o.rpcNetwork.URLs = []string{anvilChart.ClusterWSURL}
		} else {
			// Test is running locally, set forwarded URL of Anvil blockchain node
			o.rpcNetwork.URLs = []string{anvilChart.ForwardedWSURL}
			o.rpcNetwork.HTTPURLs = []string{anvilChart.ForwardedHTTPURL}
		}
	} else if o.rpcNetwork.Simulated && o.rpcNetwork.Name == blockchain.SimulatedEVMNetwork.Name {
		if testEnv.Cfg.InsideK8s {
			// Test is running inside K8s
			o.rpcNetwork.URLs = blockchain.SimulatedEVMNetwork.URLs
		} else {
			// Test is running locally, set forwarded URL of Geth blockchain node
			wsURLs := o.testEnvironment.URLs[blockchain.SimulatedEVMNetwork.Name]
			httpURLs := o.testEnvironment.URLs[blockchain.SimulatedEVMNetwork.Name+"_http"]
			require.NotEmpty(o.t, wsURLs, "Forwarded Geth URLs should not be empty")
			require.NotEmpty(o.t, httpURLs, "Forwarded Geth URLs should not be empty")
			o.rpcNetwork.URLs = wsURLs
			o.rpcNetwork.HTTPURLs = httpURLs
		}
	}
}

// Environment returns the full K8s test environment
func (o *OCRSoakTest) Environment() *environment.Environment {
	return o.testEnvironment
}

// Setup initializes the OCR Soak Test by setting up clients, funding nodes, and deploying OCR contracts.
func (o *OCRSoakTest) Setup(ocrTestConfig tt.OcrTestConfig) {
	o.initializeClients()
	o.deployLinkTokenContract(ocrTestConfig)
	o.fundChainlinkNodes()

	o.startingValue = 5

	var forwarders []common.Address
	if o.OperatorForwarderFlow {
		_, forwarders = o.deployForwarderContracts()
	}

	o.setupOCRContracts(ocrTestConfig, forwarders)
	o.log.Info().Msg("OCR Soak Test Setup Complete")
}

// initializeClients sets up the Seth client, Chainlink nodes, and mock server.
func (o *OCRSoakTest) initializeClients() {
	sethClient, err := seth_utils.GetChainClient(o.Config, o.rpcNetwork)
	require.NoError(o.t, err, "Error creating seth client")
	o.sethClient = sethClient

	nodes, err := nodeclient.ConnectChainlinkNodes(o.testEnvironment)
	require.NoError(o.t, err, "Connecting to chainlink nodes shouldn't fail")
	o.bootstrapNode, o.workerNodes = nodes[0], nodes[1:]

	o.mockServer = ctf_client.ConnectMockServer(o.testEnvironment)
	require.NoError(o.t, err, "Creating mockserver clients shouldn't fail")
}

func (o *OCRSoakTest) deployLinkTokenContract(ocrTestConfig tt.OcrTestConfig) {
	linkContract, err := actions.LinkTokenContract(o.log, o.sethClient, ocrTestConfig.GetActiveOCRConfig())
	require.NoError(o.t, err, "Error loading/deploying link token contract")
	o.linkContract = linkContract
}

// fundChainlinkNodes funds the Chainlink worker nodes.
func (o *OCRSoakTest) fundChainlinkNodes() {
	o.log.Info().Float64("ETH amount per node", *o.Config.Common.ChainlinkNodeFunding).Msg("Funding Chainlink nodes")
	err := actions.FundChainlinkNodesFromRootAddress(o.log, o.sethClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes), big.NewFloat(*o.Config.Common.ChainlinkNodeFunding))
	require.NoError(o.t, err, "Error funding Chainlink nodes")
}

// deployForwarderContracts deploys forwarder contracts if OperatorForwarderFlow is enabled.
func (o *OCRSoakTest) deployForwarderContracts() (operators []common.Address, forwarders []common.Address) {
	operators, forwarders, _ = actions.DeployForwarderContracts(
		o.t, o.sethClient, common.HexToAddress(o.linkContract.Address()), len(o.workerNodes),
	)
	require.Equal(o.t, len(o.workerNodes), len(operators), "Number of operators should match number of nodes")
	require.Equal(o.t, len(o.workerNodes), len(forwarders), "Number of authorized forwarders should match number of nodes")

	forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.workerNodes)
	require.NoError(o.t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	for i := range o.workerNodes {
		actions.AcceptAuthorizedReceiversOperator(o.t, o.log, o.sethClient, operators[i], forwarders[i], []common.Address{forwarderNodesAddresses[i]})
		require.NoError(o.t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions.TrackForwarder(o.t, o.sethClient, forwarders[i], o.workerNodes[i])
	}
	return operators, forwarders
}

// setupOCRContracts deploys and configures OCR contracts based on the version and forwarder flow.
func (o *OCRSoakTest) setupOCRContracts(ocrTestConfig tt.OcrTestConfig, forwarders []common.Address) {
	if o.OCRVersion == "1" {
		o.setupOCRv1Contracts(forwarders)
	} else if o.OCRVersion == "2" {
		o.setupOCRv2Contracts(ocrTestConfig, forwarders)
	}
}

// setupOCRv1Contracts deploys and configures OCRv1 contracts based on the forwarder flow.
func (o *OCRSoakTest) setupOCRv1Contracts(forwarders []common.Address) {
	var err error
	if o.OperatorForwarderFlow {
		o.ocrV1Instances, err = actions.DeployOCRContractsForwarderFlow(
			o.log,
			o.sethClient,
			o.Config.GetActiveOCRConfig(),
			common.HexToAddress(o.linkContract.Address()),
			contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes),
			forwarders,
		)
		require.NoError(o.t, err, "Error deploying OCR Forwarder contracts")
		o.createJobsWithForwarder()
	} else {
		o.ocrV1Instances, err = actions.SetupOCRv1Contracts(
			o.log,
			o.sethClient,
			o.Config.GetActiveOCRConfig(),
			common.HexToAddress(o.linkContract.Address()),
			contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(o.workerNodes),
		)
		require.NoError(o.t, err, "Error setting up OCRv1 contracts")
		err = o.createOCRv1Jobs()
		require.NoError(o.t, err, "Error creating OCR jobs")
	}

	o.storeOCRInstancesV1()
}

// setupOCRv2Contracts sets up and configures OCRv2 contracts.
func (o *OCRSoakTest) setupOCRv2Contracts(ocrTestConfig tt.OcrTestConfig, forwarders []common.Address) {
	var err error
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
	o.ocrV2Instances, err = actions.SetupOCRv2Contracts(
		o.log, o.sethClient, ocrTestConfig.GetActiveOCRConfig(), common.HexToAddress(o.linkContract.Address()), transmitters, ocrOffchainOptions,
	)
	require.NoError(o.t, err, "Error deploying OCRv2 contracts")
	err = o.createOCRv2Jobs()
	require.NoError(o.t, err, "Error creating OCR jobs")
	if !ocrTestConfig.GetActiveOCRConfig().UseExistingOffChainAggregatorsContracts() || (ocrTestConfig.GetActiveOCRConfig().UseExistingOffChainAggregatorsContracts() && ocrTestConfig.GetActiveOCRConfig().ConfigureExistingOffChainAggregatorsContracts()) {
		contractConfig, err := actions.BuildMedianOCR2Config(o.workerNodes, ocrOffchainOptions)
		require.NoError(o.t, err, "Error building median config")
		err = actions.ConfigureOCRv2AggregatorContracts(contractConfig, o.ocrV2Instances)
		require.NoError(o.t, err, "Error configuring OCRv2 aggregator contracts")
	}
	o.storeOCRInstancesV2()
}

// storeOCRInstancesV1 stores OCRv1 contract instances by their addresses.
func (o *OCRSoakTest) storeOCRInstancesV1() {
	for _, ocrInstance := range o.ocrV1Instances {
		o.ocrV1InstanceMap[ocrInstance.Address()] = ocrInstance
	}
}

// storeOCRInstancesV2 stores OCRv2 contract instances by their addresses.
func (o *OCRSoakTest) storeOCRInstancesV2() {
	for _, ocrInstance := range o.ocrV2Instances {
		o.ocrV2InstanceMap[ocrInstance.Address()] = ocrInstance
	}
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run() {
	config, err := tc.GetConfig([]string{"soak"}, tc.OCR)
	require.NoError(o.t, err, "Error getting config")

	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), time.Second*5)
	latestBlockNum, err := o.sethClient.Client.BlockNumber(ctx)
	cancel()
	require.NoError(o.t, err, "Error getting current block number")
	o.startingBlockNum = latestBlockNum

	o.log.Info().
		Str("Test Duration", o.Config.GetActiveOCRConfig().Common.TestDuration.Duration.Truncate(time.Second).String()).
		Int("Number of OCR Contracts", *config.GetActiveOCRConfig().Common.NumberOfContracts).
		Str("OCR Version", o.OCRVersion).
		Msg("Starting OCR Soak Test")

	o.testLoop(o.Config.GetActiveOCRConfig().Common.TestDuration.Duration, o.startingValue)
	o.complete()
}

// createJobsWithForwarder creates OCR jobs with the forwarder setup.
func (o *OCRSoakTest) createJobsWithForwarder() {
	actions.CreateOCRJobsWithForwarder(o.t, o.ocrV1Instances, o.bootstrapNode, o.workerNodes, o.startingValue, o.mockServer, o.sethClient.ChainID)
}

// createOCRv1Jobs creates OCRv1 jobs.
func (o *OCRSoakTest) createOCRv1Jobs() error {
	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), time.Second*5)
	defer cancel()

	chainId, err := o.sethClient.Client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("error getting chain ID: %w", err)
	}

	err = actions.CreateOCRJobs(o.ocrV1Instances, o.bootstrapNode, o.workerNodes, o.startingValue, o.mockServer, chainId.String())
	if err != nil {
		return fmt.Errorf("error creating OCRv1 jobs: %w", err)
	}
	return nil
}

// createOCRv2Jobs creates OCRv2 jobs.
func (o *OCRSoakTest) createOCRv2Jobs() error {
	err := actions.CreateOCRv2Jobs(o.ocrV2Instances, o.bootstrapNode, o.workerNodes, o.mockServer, o.startingValue, o.sethClient.ChainID, o.OperatorForwarderFlow, o.log)
	if err != nil {
		return fmt.Errorf("error creating OCRv2 jobs: %w", err)
	}
	return nil
}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	*seth.Client,
	string,
	[]*nodeclient.ChainlinkK8sClient,
	reportModel.TestReporter,
	reportModel.GrafanaURLProvider,
) {
	return t, o.sethClient, o.namespace, append(o.workerNodes, o.bootstrapNode), &o.TestReporter, o.Config
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
	ReorgHappened    bool     `toml:"reorgHappened"`
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
		TestDuration:         o.Config.GetActiveOCRConfig().Common.TestDuration.Duration,
		OCRContractAddresses: ocrAddresses,
		OCRVersion:           o.OCRVersion,

		MockServerURL:    "http://mockserver:1080", // TODO: Make this dynamic
		BootStrapNodeURL: o.bootstrapNode.URL(),
		WorkerNodeURLs:   workerNodeURLs,
		ReorgHappened:    o.reorgHappened,
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
	o.Config.GetActiveOCRConfig().Common.TestDuration = &duration
	o.timeLeft = testState.TestDuration - testState.TimeRunning
	o.startTime = testState.StartTime
	o.startingBlockNum = testState.StartingBlockNum
	o.reorgHappened = testState.ReorgHappened
	o.OCRVersion = testState.OCRVersion

	o.bootstrapNode, err = nodeclient.ConnectChainlinkNodeURL(testState.BootStrapNodeURL)
	if err != nil {
		return err
	}
	o.workerNodes, err = nodeclient.ConnectChainlinkNodeURLs(testState.WorkerNodeURLs)
	if err != nil {
		return err
	}

	if testState.OCRVersion == "1" {
		o.ocrV1Instances = make([]contracts.OffchainAggregator, len(testState.OCRContractAddresses))
		for i, addr := range testState.OCRContractAddresses {
			instance, err := contracts.LoadOffChainAggregator(o.log, o.sethClient, common.HexToAddress(addr))
			if err != nil {
				return fmt.Errorf("failed to instantiate OCR instance: %w", err)
			}
			o.ocrV1Instances[i] = &instance
		}
	} else if testState.OCRVersion == "2" {
		o.ocrV2Instances = make([]contracts.OffchainAggregatorV2, len(testState.OCRContractAddresses))
		for i, addr := range testState.OCRContractAddresses {
			instance, err := contracts.LoadOffchainAggregatorV2(o.log, o.sethClient, common.HexToAddress(addr))
			if err != nil {
				return err
			}
			o.ocrV2Instances[i] = &instance
		}
	}

	o.mockServer = ctf_client.ConnectMockServerURL(testState.MockServerURL)
	return err
}

func (o *OCRSoakTest) Resume() {
	o.testIssues = append(o.testIssues, &testreporters.TestIssue{
		StartTime: time.Now(),
		Message:   "Test Resumed",
	})
	o.log.Info().
		Str("Total Duration", o.Config.GetActiveOCRConfig().Common.TestDuration.String()).
		Str("Time Left", o.timeLeft.String()).
		Msg("Resuming OCR Soak Test")

	ocrAddresses := make([]common.Address, *o.Config.GetActiveOCRConfig().Common.NumberOfContracts)

	if o.OCRVersion == "1" {
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
	} else if o.OCRVersion == "2" {
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

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// Channel to signal polling to reset round event counter
	resetEventCounter := make(chan struct{})
	defer close(resetEventCounter)

	lastValue := 0
	newRoundTrigger := time.NewTimer(0) // Want to trigger a new round ASAP
	defer newRoundTrigger.Stop()
	o.setFilterQuery()
	wg.Add(1)
	go o.pollingOCREvents(ctx, &wg, resetEventCounter)

	n := o.Config.GetNetworkConfig()

	// Schedule blockchain re-org if needed
	// Reorg only avaible for Simulated Geth
	if n.IsSimulatedGethSelected() && n.GethReorgConfig.Enabled {
		var reorgDelay time.Duration
		if n.GethReorgConfig.DelayCreate.Duration > testDuration {
			// This may happen when test is resumed and the reorg delay is longer than the time left
			o.log.Warn().Msg("Reorg delay is longer than test duration, reorg scheduled immediately")
			reorgDelay = 0
		} else {
			reorgDelay = n.GethReorgConfig.DelayCreate.Duration
		}
		time.AfterFunc(reorgDelay, func() {
			if !o.reorgHappened {
				o.startGethBlockchainReorg(o.rpcNetwork, n.GethReorgConfig)
			}
		})
	}

	// Schedule gas simulations if needed
	// Gas simulation only available for Anvil
	if o.rpcNetwork.Name == "Anvil" {
		ac := o.Config.GetNetworkConfig().AnvilConfigs["ANVIL"]
		if ac != nil && ac.GasSpikeSimulation.Enabled {
			var delay time.Duration
			if ac.GasSpikeSimulation.DelayCreate.Duration > testDuration {
				// This may happen when test is resumed and the reorg delay is longer than the time left
				o.log.Warn().Msg("Gas spike simulation delay is longer than test duration, gas simulation scheduled immediately")
				delay = 0
			} else {
				delay = ac.GasSpikeSimulation.DelayCreate.Duration
			}
			time.AfterFunc(delay, func() {
				if !o.gasSpikeSimulationHappened {
					o.startAnvilGasSpikeSimulation(o.rpcNetwork, ac.GasSpikeSimulation)
				}
			})
		}
		if ac != nil && ac.GasLimitSimulation.Enabled {
			var delay time.Duration
			if ac.GasLimitSimulation.DelayCreate.Duration > testDuration {
				// This may happen when test is resumed and the reorg delay is longer than the time left
				o.log.Warn().Msg("Gas limit simulation delay is longer than test duration, gas simulation scheduled immediately")
				delay = 0
			} else {
				delay = ac.GasLimitSimulation.DelayCreate.Duration
			}
			time.AfterFunc(delay, func() {
				if !o.gasLimitSimulationHappened {
					o.startAnvilGasLimitSimulation(o.rpcNetwork, ac.GasLimitSimulation)
				}
			})
		}
	}

	// Schedule chaos simulations if needed
	if len(o.chaosList) > 0 {
		for _, chaos := range o.chaosList {
			chaos.Create(context.Background())
			chaos.AddListener(havoc.NewChaosLogger(o.log))
			chaos.AddListener(ocrTestChaosListener{t: o.t})
			// Add Grafana annotation if configured
			if o.Config.Logging.Grafana != nil && o.Config.Logging.Grafana.BaseUrl != nil && o.Config.Logging.Grafana.BearerToken != nil && o.Config.Logging.Grafana.DashboardUID != nil {
				chaos.AddListener(havoc.NewSingleLineGrafanaAnnotator(*o.Config.Logging.Grafana.BaseUrl, *o.Config.Logging.Grafana.BearerToken, *o.Config.Logging.Grafana.DashboardUID, o.log))
			} else {
				o.log.Warn().Msg("Skipping Grafana annotation for chaos simulation. Grafana config is missing either BearerToken, BaseUrl or DashboardUID")
			}
		}
	}

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
			o.deleteChaosSimulations()
			os.Exit(interruptedExitCode) // Exit with interrupted code to indicate test was interrupted, not just a normal failure
		case <-endTest:
			cancel()
			wg.Wait() // Wait for polling to complete
			return
		case <-newRoundTrigger.C:
			err := o.triggerNewRound(newValue)
			timerReset := o.Config.GetActiveOCRConfig().Soak.TimeBetweenRounds.Duration
			if err != nil {
				timerReset = time.Second * 5
				o.log.Error().Err(err).
					Str("Waiting", timerReset.String()).
					Msg("Error triggering new round, waiting and trying again. Possible connection issues with mockserver")
			}
			// Signal polling to reset event counter
			resetEventCounter <- struct{}{}
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

func (o *OCRSoakTest) startGethBlockchainReorg(network blockchain.EVMNetwork, conf ctf_config.ReorgConfig) {
	client := ctf_client.NewRPCClient(network.HTTPURLs[0], nil)
	o.log.Info().
		Str("URL", client.URL).
		Int("Depth", conf.Depth).
		Msg("Starting blockchain reorg on Simulated Geth chain")
	o.postGrafanaAnnotation(fmt.Sprintf("Starting blockchain reorg on Simulated Geth chain with depth %d", conf.Depth), nil)
	err := client.GethSetHead(conf.Depth)
	require.NoError(o.t, err, "Error starting blockchain reorg on Simulated Geth chain")
	o.reorgHappened = true
}

func (o *OCRSoakTest) startAnvilGasSpikeSimulation(network blockchain.EVMNetwork, conf ctf_config.GasSpikeSimulationConfig) {
	client := ctf_client.NewRPCClient(network.HTTPURLs[0], nil)
	o.log.Info().
		Str("URL", client.URL).
		Any("GasSpikeSimulationConfig", conf).
		Msg("Starting gas spike simulation on Anvil chain")
	o.postGrafanaAnnotation(fmt.Sprintf("Starting gas spike simulation on Anvil chain. Config: %+v", conf), nil)
	err := client.ModulateBaseFeeOverDuration(o.log, conf.StartGasPrice, conf.GasRisePercentage, conf.Duration.Duration, conf.GasSpike)
	o.postGrafanaAnnotation(fmt.Sprintf("Gas spike simulation ended. Config: %+v", conf), nil)
	require.NoError(o.t, err, "Error starting gas simulation on Anvil chain")
	o.gasSpikeSimulationHappened = true
}

func (o *OCRSoakTest) startAnvilGasLimitSimulation(network blockchain.EVMNetwork, conf ctf_config.GasLimitSimulationConfig) {
	client := ctf_client.NewRPCClient(network.HTTPURLs[0], nil)
	latestBlock, err := o.sethClient.Client.BlockByNumber(context.Background(), nil)
	require.NoError(o.t, err)
	newGasLimit := int64(math.Ceil(float64(latestBlock.GasUsed()) * conf.NextGasLimitPercentage))
	o.log.Info().
		Str("URL", client.URL).
		Any("GasLimitSimulationConfig", conf).
		Uint64("LatestBlock", latestBlock.Number().Uint64()).
		Uint64("LatestGasUsed", latestBlock.GasUsed()).
		Uint64("LatestGasLimit", latestBlock.GasLimit()).
		Int64("NewGasLimit", newGasLimit).
		Msg("Starting gas limit simulation on Anvil chain")
	o.postGrafanaAnnotation(fmt.Sprintf("Starting gas limit simulation on Anvil chain. Config: %+v", conf), nil)
	err = client.AnvilSetBlockGasLimit([]interface{}{newGasLimit})
	require.NoError(o.t, err, "Error starting gas simulation on Anvil chain")
	time.Sleep(conf.Duration.Duration)
	o.log.Info().
		Str("URL", client.URL).
		Any("GasLimitSimulationConfig", conf).
		Uint64("LatestGasLimit", latestBlock.GasLimit()).
		Msg("Returning to old gas limit simulation on Anvil chain")
	o.postGrafanaAnnotation(fmt.Sprintf("Returning to old gas limit simulation on Anvil chain. Config: %+v", conf), nil)
	err = client.AnvilSetBlockGasLimit([]interface{}{latestBlock.GasLimit()})
	require.NoError(o.t, err, "Error starting gas simulation on Anvil chain")
	o.gasLimitSimulationHappened = true
}

// Delete k8s chaos objects it any of them still exist
// This is needed to clean up the chaos objects if the test is interrupted or it finishes
func (o *OCRSoakTest) deleteChaosSimulations() {
	for _, chaos := range o.chaosList {
		err := chaos.Delete(context.Background())
		// Check if the error is because the chaos object is already deleted
		if err != nil && !strings.Contains(err.Error(), "not found") {
			o.log.Error().Err(err).Msg("Error deleting chaos object")
		}
	}
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

// pollingOCREvents Polls the blocks for OCR events and logs them to the test logger
func (o *OCRSoakTest) pollingOCREvents(ctx context.Context, wg *sync.WaitGroup, resetEventCounter <-chan struct{}) {
	defer wg.Done()
	// Keep track of the last processed block number
	processedBlockNum := o.startingBlockNum - 1
	// TODO: Make this configurable
	pollInterval := time.Second * 30
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Retrieve expected number of events per round from configuration
	expectedEventsPerRound := *o.Config.GetActiveOCRConfig().Common.NumberOfContracts
	eventCounter := 0
	roundTimeout := o.Config.GetActiveOCRConfig().Soak.TimeBetweenRounds.Duration
	timeoutTimer := time.NewTimer(roundTimeout)
	round := 0
	defer timeoutTimer.Stop()

	o.log.Info().Msg("Start Polling for Answer Updated Events")

	for {
		select {
		case <-resetEventCounter:
			if round != 0 {
				if eventCounter == expectedEventsPerRound {
					o.log.Info().
						Int("Events found", eventCounter).
						Int("Events Expected", expectedEventsPerRound).
						Msg("All expected events found")
				} else if eventCounter < expectedEventsPerRound {
					o.log.Warn().
						Int("Events found", eventCounter).
						Int("Events Expected", expectedEventsPerRound).
						Msg("Expected to find more events")
				}
			}
			// Reset event counter and timer for new round
			eventCounter = 0
			// Safely stop and drain the timer if a value is present
			if !timeoutTimer.Stop() {
				<-timeoutTimer.C
			}
			timeoutTimer.Reset(roundTimeout)
			o.log.Info().Msg("Polling for new round, event counter reset")
			round++
		case <-ctx.Done():
			o.log.Info().Msg("Test duration ended, finalizing event polling")
			timeoutTimer.Reset(roundTimeout)
			// Wait until expected events are fetched or until timeout
			for eventCounter < expectedEventsPerRound {
				select {
				case <-timeoutTimer.C:
					o.log.Warn().Msg("Timeout reached while waiting for final events")
					return
				case <-ticker.C:
					o.fetchAndProcessEvents(&eventCounter, expectedEventsPerRound, &processedBlockNum)
				}
			}
			o.log.Info().
				Int("Events found", eventCounter).
				Int("Events Expected", expectedEventsPerRound).
				Msg("Stop polling.")
			return
		case <-ticker.C:
			o.fetchAndProcessEvents(&eventCounter, expectedEventsPerRound, &processedBlockNum)
		}
	}
}

// Helper function to poll events and update eventCounter
func (o *OCRSoakTest) fetchAndProcessEvents(eventCounter *int, expectedEvents int, processedBlockNum *uint64) {
	latestBlock, err := o.sethClient.Client.BlockNumber(context.Background())
	if err != nil {
		o.log.Error().Err(err).Msg("Error getting latest block number")
		return
	}

	if *processedBlockNum == latestBlock {
		o.log.Debug().
			Uint64("Latest Block", latestBlock).
			Uint64("Last Processed Block Number", *processedBlockNum).
			Msg("No new blocks since last poll")
		return
	}

	// Check if the latest block is behind processedBlockNum due to possible reorgs
	if *processedBlockNum > latestBlock {
		o.log.Error().
			Uint64("From Block", *processedBlockNum).
			Uint64("To Block", latestBlock).
			Msg("The latest block is behind the processed block. This could happen due to RPC issues or possibly a reorg")
		*processedBlockNum = latestBlock
		return
	}

	fromBlock := *processedBlockNum + 1
	o.filterQuery.FromBlock = big.NewInt(0).SetUint64(fromBlock)
	o.filterQuery.ToBlock = big.NewInt(0).SetUint64(latestBlock)

	o.log.Debug().
		Uint64("From Block", fromBlock).
		Uint64("To Block", latestBlock).
		Msg("Fetching logs for the specified range")

	logs, err := o.sethClient.Client.FilterLogs(context.Background(), o.filterQuery)
	if err != nil {
		o.log.Error().Err(err).Msg("Error fetching logs")
		return
	}

	for _, event := range logs {
		*eventCounter++
		if o.OCRVersion == "1" {
			answerUpdated, err := o.ocrV1Instances[0].ParseEventAnswerUpdated(event)
			if err != nil {
				o.log.Warn().
					Err(err).
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Msg("Error parsing event as AnswerUpdated")
				continue
			}
			if *eventCounter <= expectedEvents {
				o.log.Info().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Answer Updated Event")
			} else {
				o.log.Error().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Excess event detected, beyond expected count")
			}
		} else if o.OCRVersion == "2" {
			answerUpdated, err := o.ocrV2Instances[0].ParseEventAnswerUpdated(event)
			if err != nil {
				o.log.Warn().
					Err(err).
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Msg("Error parsing event as AnswerUpdated")
				continue
			}
			if *eventCounter <= expectedEvents {
				o.log.Info().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Answer Updated Event")
			} else {
				o.log.Error().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Excess event detected, beyond expected count")
			}
		}
	}
	*processedBlockNum = latestBlock
}

// triggers a new OCR round by setting a new mock adapter value
func (o *OCRSoakTest) triggerNewRound(newValue int) error {
	if len(o.ocrRoundStates) > 0 {
		o.ocrRoundStates[len(o.ocrRoundStates)-1].EndTime = time.Now()
	}

	var err error
	if o.OCRVersion == "1" {
		err = actions.SetAllAdapterResponsesToTheSameValue(newValue, o.ocrV1Instances, o.workerNodes, o.mockServer)
	} else if o.OCRVersion == "2" {
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
	if o.OCRVersion == "1" {
		for _, ocrInstance := range o.ocrV1Instances {
			expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
		}
	} else if o.OCRVersion == "2" {
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

	// Set from block to be starting block before filtering
	o.filterQuery.FromBlock = big.NewInt(0).SetUint64(o.startingBlockNum)

	// We must retrieve the events, use exponential backoff for timeout to retry
	timeout := time.Second * 15
	o.log.Info().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Retrieving on-chain events")

	ctx, cancel := context.WithTimeout(testcontext.Get(o.t), timeout)
	contractEvents, err := o.sethClient.Client.FilterLogs(ctx, o.filterQuery)
	cancel()
	for err != nil {
		o.log.Info().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Retrieving on-chain events")
		ctx, cancel := context.WithTimeout(testcontext.Get(o.t), timeout)
		contractEvents, err = o.sethClient.Client.FilterLogs(ctx, o.filterQuery)
		cancel()
		if err != nil {
			o.log.Warn().Interface("Filter Query", o.filterQuery).Str("Timeout", timeout.String()).Msg("Error collecting on-chain events, trying again")
			timeout *= 2
		}
	}

	sortedFoundEvents := make([]*testreporters.FoundEvent, 0)
	for _, event := range contractEvents {
		if o.OCRVersion == "1" {
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
		} else if o.OCRVersion == "2" {
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
		Int("Events collected", len(contractEvents)).
		Msg("Collected on-chain events")

	if len(contractEvents) == 0 {
		return fmt.Errorf("no events were collected")
	}

	return nil
}

// ensureValues ensures that all values needed to run the test are present
func (o *OCRSoakTest) ensureInputValues() error {
	ocrConfig := o.Config.GetActiveOCRConfig()
	if o.OCRVersion != "1" && o.OCRVersion != "2" {
		return fmt.Errorf("OCR version must be 1 or 2, found %s", o.OCRVersion)
	}
	if ocrConfig.Common.NumberOfContracts != nil && *ocrConfig.Common.NumberOfContracts <= 0 {
		return fmt.Errorf("number of OCR contracts must be set and greater than 0, found %d", ocrConfig.Common.NumberOfContracts)
	}
	if o.Config.Common.ChainlinkNodeFunding != nil && *o.Config.Common.ChainlinkNodeFunding <= 0 {
		return fmt.Errorf("chainlink node funding must be greater than 0, found %f", *o.Config.Common.ChainlinkNodeFunding)
	}
	if o.Config.GetActiveOCRConfig().Common.TestDuration != nil && o.Config.GetActiveOCRConfig().Common.TestDuration.Duration <= time.Minute {
		return fmt.Errorf("test duration must be greater than 1 minute, found %s", o.Config.GetActiveOCRConfig().Common.TestDuration)
	}
	soakConfig := ocrConfig.Soak
	if soakConfig.TimeBetweenRounds != nil && soakConfig.TimeBetweenRounds.Duration >= time.Hour {
		return fmt.Errorf("time between rounds must be less than 1 hour, found %s", soakConfig.TimeBetweenRounds)
	}
	if soakConfig.TimeBetweenRounds != nil && soakConfig.TimeBetweenRounds.Duration < time.Second*30 {
		return fmt.Errorf("time between rounds must be greater or equal to 30 seconds, found %s", soakConfig.TimeBetweenRounds)
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
	var contractAddresses []common.Address
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

func (o *OCRSoakTest) postGrafanaAnnotation(text string, tags []string) {
	var grafanaClient *grafana.Client
	var dashboardUID *string
	if o.Config.Logging.Grafana != nil {
		baseURL := o.Config.Logging.Grafana.BaseUrl
		dashboardUID = o.Config.Logging.Grafana.DashboardUID
		token := o.Config.Logging.Grafana.BearerToken
		if token == nil || baseURL == nil || dashboardUID == nil {
			o.log.Warn().Msg("Skipping Grafana annotation. Grafana config is missing either BearerToken, BaseUrl or DashboardUID")
			return
		}
		grafanaClient = grafana.NewGrafanaClient(*baseURL, *token)
	}
	_, _, err := grafanaClient.PostAnnotation(grafana.PostAnnotation{
		DashboardUID: *dashboardUID,
		Tags:         tags,
		Text:         fmt.Sprintf("<b>Test Namespace: %s<pre>%s</pre></b>", o.namespace, text),
	})
	if err != nil {
		o.log.Error().Err(err).Msg("Error posting annotation to Grafana")
	} else {
		o.log.Info().Msgf("Annotated Grafana dashboard with text: %s", text)
	}
}

type ocrTestChaosListener struct {
	t *testing.T
}

func (l ocrTestChaosListener) OnChaosCreated(_ havoc.Chaos) {
}

func (l ocrTestChaosListener) OnChaosCreationFailed(chaos havoc.Chaos, reason error) {
	// Fail the test if chaos creation fails during chaos simulation
	require.FailNow(l.t, "Error creating chaos simulation", reason.Error(), chaos)
}

func (l ocrTestChaosListener) OnChaosStarted(_ havoc.Chaos) {
}

func (l ocrTestChaosListener) OnChaosPaused(_ havoc.Chaos) {
}

func (l ocrTestChaosListener) OnChaosEnded(_ havoc.Chaos) {
}

func (l ocrTestChaosListener) OnChaosStatusUnknown(_ havoc.Chaos) {
}

func (l ocrTestChaosListener) OnScheduleCreated(_ havoc.Schedule) {
}

func (l ocrTestChaosListener) OnScheduleDeleted(_ havoc.Schedule) {
}
