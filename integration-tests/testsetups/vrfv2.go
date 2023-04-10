package testsetups

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// VRFV2SoakTest defines a typical VRFV2 soak test
type VRFV2SoakTest struct {
	Inputs *VRFV2SoakTestInputs

	TestReporter testreporters.VRFV2SoakTestReporter
	mockServer   *ctfClient.MockserverClient

	testEnvironment *environment.Environment
	ChainlinkNodes  []*client.Chainlink
	chainClient     blockchain.EVMClient
	DefaultNetwork  blockchain.EVMClient

	NumberOfRequests int

	ErrorOccurred error
	ErrorCount    int
}

// VRFV2SoakTestTestFunc function type for the request and validation you want done on each iteration
type VRFV2SoakTestTestFunc func(t *VRFV2SoakTest, requestNumber int) error

// VRFV2SoakTestInputs define required inputs to run a vrfv2 soak test
type VRFV2SoakTestInputs struct {
	BlockchainClient     blockchain.EVMClient // Client for the test to connect to the blockchain with
	TestDuration         time.Duration        // How long to run the test for (assuming things pass)
	ChainlinkNodeFunding *big.Float           // Amount of ETH to fund each chainlink node with
	StopTestOnError      bool                 // Do we want the test to stop after any error or just continue on

	RequestsPerMinute int                   // Number of requests for randomness per minute
	TestFunc          VRFV2SoakTestTestFunc // The function that makes the request and validations wanted
}

// NewVRFV2SoakTest creates a new vrfv2 soak test to setup and run
func NewVRFV2SoakTest(inputs *VRFV2SoakTestInputs) *VRFV2SoakTest {
	return &VRFV2SoakTest{
		Inputs: inputs,
		TestReporter: testreporters.VRFV2SoakTestReporter{
			Reports: make(map[string]*testreporters.VRFV2SoakTestReport),
		},
	}
}

// Setup sets up the test environment
func (v *VRFV2SoakTest) Setup(t *testing.T, env *environment.Environment, isLocal bool) {
	v.ensureInputValues(t)
	v.testEnvironment = env
	var err error

	// Make connections to soak test resources
	v.ChainlinkNodes, err = client.ConnectChainlinkNodes(env)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	v.mockServer, err = ctfClient.ConnectMockServer(env)
	require.NoError(t, err, "Error connecting to mockserver")

	v.chainClient.ParallelTransactions(true)
}

// Run starts the VRFV2 soak test
func (v *VRFV2SoakTest) Run(t *testing.T) {
	l := utils.GetTestLogger(t)
	l.Info().
		Str("Test Duration", v.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Max number of requests per minute wanted", v.Inputs.RequestsPerMinute).
		Msg("Starting VRFV2 Soak Test")

	// set the requests to only run for a certain amount of time
	testContext, testCancel := context.WithTimeout(context.Background(), v.Inputs.TestDuration)
	defer testCancel()

	v.NumberOfRequests = 0

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
			v.NumberOfRequests++
			go requestAndValidate(v, v.NumberOfRequests)
		}

		if stop {
			break // breaks the for loop and stops the test
		}
	}
	l.Info().Int("Requests", v.NumberOfRequests).Msg("Total Completed Requests")
	l.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished VRFV2 Soak Test Requests")
	require.Equal(t, 0, v.ErrorCount, "Expected 0 errors")
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
func (t *VRFV2SoakTest) TearDownVals() (*environment.Environment, []*client.Chainlink, reportModel.TestReporter, blockchain.EVMClient) {
	return t.testEnvironment, t.ChainlinkNodes, &t.TestReporter, t.chainClient
}

// ensureValues ensures that all values needed to run the test are present
func (v *VRFV2SoakTest) ensureInputValues(t *testing.T) {
	inputs := v.Inputs
	require.NotNil(t, inputs.BlockchainClient, "Need a valid blockchain client for the test")
	v.chainClient = inputs.BlockchainClient
	require.GreaterOrEqual(t, inputs.RequestsPerMinute, 1, "Expecting at least 1 request per minute")
	funding, _ := inputs.ChainlinkNodeFunding.Float64()
	require.Greater(t, funding, 0, "Need some amount of funding for Chainlink nodes")
	require.GreaterOrEqual(t, inputs.TestDuration, time.Minute, "Test duration should be longer than 1 minute")
	require.NotNil(t, inputs.TestFunc, "Expected there to be test to run")
}
