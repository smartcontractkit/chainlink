package testsetups

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// VRFV2SoakTest defines a typical VRFV2 soak test
type VRFV2SoakTest struct {
	Inputs *VRFV2SoakTestInputs

	TestReporter testreporters.VRFV2SoakTestReporter

	testEnvironment *environment.Environment
	namespace       string
	ChainlinkNodes  []*client.ChainlinkK8sClient
	chainClient     blockchain.EVMClient
	DefaultNetwork  blockchain.EVMClient

	NumberOfRandRequests int

	ErrorOccurred error
	ErrorCount    int
}

// VRFV2SoakTestTestFunc function type for the request and validation you want done on each iteration
type VRFV2SoakTestTestFunc func(t *VRFV2SoakTest, requestNumber int) error

// VRFV2SoakTestInputs define required inputs to run a vrfv2 soak test
type VRFV2SoakTestInputs struct {
	BlockchainClient     blockchain.EVMClient // Client for the test to connect to the blockchain with
	TestDuration         time.Duration        `envconfig:"TEST_DURATION" default:"15m"`         // How long to run the test for (assuming things pass)
	ChainlinkNodeFunding *big.Float           `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"` // Amount of ETH to fund each chainlink node with
	SubscriptionFunding  *big.Int             `envconfig:"SUBSCRIPTION_FUNDING" default:"100"`  // Amount of Link to fund VRF Coordinator subscription
	StopTestOnError      bool                 // Do we want the test to stop after any error or just continue on

	RequestsPerMinute                int `envconfig:"REQUESTS_PER_MINUTE" default:"10"` // Number of requests for randomness per minute
	RandomnessRequestCountPerRequest int `envconfig:"RANDOMNESS_REQUEST_COUNT_PER_REQUEST" default:"1"`
	ConsumerContract                 contracts.VRFv2LoadTestConsumer
	TestFunc                         VRFV2SoakTestTestFunc // The function that makes the request and validations wanted
}

// NewVRFV2SoakTest creates a new vrfv2 soak test to setup and run
func NewVRFV2SoakTest(inputs *VRFV2SoakTestInputs, chainlinkNodes []*client.ChainlinkK8sClient) *VRFV2SoakTest {
	return &VRFV2SoakTest{
		Inputs: inputs,
		TestReporter: testreporters.VRFV2SoakTestReporter{
			Reports: make(map[string]*testreporters.VRFV2SoakTestReport),
		},
		ChainlinkNodes: chainlinkNodes,
	}
}

// Setup sets up the test environment
func (v *VRFV2SoakTest) Setup(t *testing.T, env *environment.Environment) {
	v.ensureInputValues(t)
	v.testEnvironment = env
	v.namespace = v.testEnvironment.Cfg.Namespace
	v.chainClient.ParallelTransactions(true)
}

// Run starts the VRFV2 soak test
func (v *VRFV2SoakTest) Run(t *testing.T) {
	l := logging.GetTestLogger(t)
	l.Info().
		Str("Test Duration", v.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Max number of requests per minute wanted", v.Inputs.RequestsPerMinute).
		Msg("Starting VRFV2 Soak Test")

	// set the requests to only run for a certain amount of time
	testContext, testCancel := context.WithTimeout(context.Background(), v.Inputs.TestDuration)
	defer testCancel()

	v.NumberOfRandRequests = 0

	// variables dealing with how often to tick and how to stop the ticker
	stop := false
	startTime := time.Now()
	ticker := time.NewTicker(time.Minute / time.Duration(v.Inputs.RequestsPerMinute))

	for {
		// start the loop by checking to see if any of the TestFunc responses have returned an error
		if v.Inputs.StopTestOnError {
			require.NoError(t, v.ErrorOccurred, "Found error")
		}
		select {
		case <-testContext.Done():
			// stop making requests
			stop = true
			ticker.Stop()
			break // breaks the select block
		case <-ticker.C:
			// make the next request
			v.NumberOfRandRequests++
			go requestAndValidate(v, v.NumberOfRandRequests)
		}
		if stop {
			break // breaks the for loop and stops the test
		}
	}

	err := v.chainClient.WaitForEvents()
	if err != nil {
		l.Error().Err(err).Msg("Error Occurred waiting for On chain events")
	}
	//wait some buffer time for requests to be fulfilled
	//todo - need to find better way for this
	time.Sleep(1 * time.Minute)

	loadTestMetrics, err := v.Inputs.ConsumerContract.GetLoadTestMetrics(nil)
	if err != nil {
		l.Error().Err(err).Msg("Error Occurred when getting Load Test Metrics from Consumer contract")
	}

	averageFulfillmentInBlockTime := new(big.Float).Quo(new(big.Float).SetInt(loadTestMetrics.AverageFulfillmentInMillions), big.NewFloat(1e6))

	l.Info().Int("Requests", v.NumberOfRandRequests).Msg("Total Completed Requests calculated from Test")
	l.Info().Uint64("Requests", loadTestMetrics.RequestCount.Uint64()).Msg("Total Completed Requests calculated from Contract")
	l.Info().Uint64("Fulfilments", loadTestMetrics.FulfilmentCount.Uint64()).Msg("Total Completed Fulfilments")
	l.Info().Uint64("Fastest Fulfilment", loadTestMetrics.FastestFulfillment.Uint64()).Msg("Fastest Fulfilment")
	l.Info().Uint64("Slowest Fulfilment", loadTestMetrics.SlowestFulfillment.Uint64()).Msg("Slowest Fulfilment")
	l.Info().Interface("Average Fulfillment", averageFulfillmentInBlockTime).Msg("Average Fulfillment In Block Time")

	//todo - need to calculate 95th percentile response time in Block time and calculate how many requests breached 256 block time requirement

	l.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished VRFV2 Soak Test Requests")
	require.Equal(t, 0, v.ErrorCount, "Expected 0 errors")
	require.Equal(t, loadTestMetrics.RequestCount.Uint64(), loadTestMetrics.FulfilmentCount.Uint64(), "Number of Rand Requests should be equal to Number of Fulfillments")
}

func requestAndValidate(t *VRFV2SoakTest, requestNumber int) {

	log.Info().Int("Request Number", requestNumber).Msg("Making a Request")
	err := t.Inputs.TestFunc(t, requestNumber)

	// only set the error to be checked if err is not nil so we avoid race conditions with passing requests
	if err != nil {
		t.ErrorOccurred = err
		log.Error().Err(err).Msg("Error Occurred during test")
		t.ErrorCount++
	}
}

// Networks returns the networks that the test is running on
func (v *VRFV2SoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	string,
	[]*client.ChainlinkK8sClient,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, v.namespace, v.ChainlinkNodes, &v.TestReporter, v.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (v *VRFV2SoakTest) ensureInputValues(t *testing.T) {
	inputs := v.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Need a valid blockchain client for the test")
	v.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.RequestsPerMinute, 1, "Expecting at least 1 request per minute")
	chainlinkNodeFunding, _ := inputs.ChainlinkNodeFunding.Float64()
	subscriptionFunding := inputs.SubscriptionFunding.Int64()
	require.Greater(t, chainlinkNodeFunding, float64(0), "Need some amount of funding for Chainlink nodes")
	require.Greater(t, subscriptionFunding, int64(0), "Need some amount of funding for VRF V2 Coordinator Subscription nodes")
	require.GreaterOrEqual(t, inputs.TestDuration, time.Minute, "Test duration should be longer than 1 minute")
	require.NotNil(t, inputs.TestFunc, "Expected there to be test to run")
}

func (i VRFV2SoakTestInputs) SetForRemoteRunner() {
	os.Setenv("TEST_VRFV2_TEST_DURATION", i.TestDuration.String())
	os.Setenv("TEST_VRFV2_CHAINLINK_NODE_FUNDING", i.ChainlinkNodeFunding.String())
	os.Setenv("TEST_VRFV2_SUBSCRIPTION_FUNDING", i.SubscriptionFunding.String())
	os.Setenv("TEST_VRFV2_REQUESTS_PER_MINUTE", strconv.Itoa(i.RequestsPerMinute))
	os.Setenv("TEST_VRFV2_RANDOMNESS_REQUEST_COUNT_PER_REQUEST", strconv.Itoa(i.RandomnessRequestCountPerRequest))

	selectedNetworks := strings.Split(os.Getenv("SELECTED_NETWORKS"), ",")
	for _, networkPrefix := range selectedNetworks {
		urlEnv := fmt.Sprintf("%s_URLS", networkPrefix)
		httpEnv := fmt.Sprintf("%s_HTTP_URLS", networkPrefix)
		os.Setenv(fmt.Sprintf("TEST_%s", urlEnv), os.Getenv(urlEnv))
		os.Setenv(fmt.Sprintf("TEST_%s", httpEnv), os.Getenv(httpEnv))
	}
}
