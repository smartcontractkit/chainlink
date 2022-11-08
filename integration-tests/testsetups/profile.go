package testsetups

//revive:disable:dot-imports
import (
	"time"

	. "github.com/onsi/gomega"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// ChainlinkProfileTest runs a piece of code on Chainlink nodes with PPROF enabled, then downloads the PPROF results
type ChainlinkProfileTest struct {
	Inputs       ChainlinkProfileTestInputs
	TestReporter testreporters.ChainlinkProfileTestReporter
	env          *environment.Environment
	c            blockchain.EVMClient
}

// ChainlinkProfileTestInputs are the inputs necessary to run a profiling tests
type ChainlinkProfileTestInputs struct {
	ProfileFunction func(*client.Chainlink)
	ProfileDuration time.Duration
	ChainlinkNodes  []*client.Chainlink
}

// NewChainlinkProfileTest prepares a new keeper Chainlink profiling test to be run
func NewChainlinkProfileTest(inputs ChainlinkProfileTestInputs) *ChainlinkProfileTest {
	return &ChainlinkProfileTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (c *ChainlinkProfileTest) Setup(env *environment.Environment) {
	c.ensureInputValues()
	c.env = env
}

// Run runs the profiling test
func (c *ChainlinkProfileTest) Run() {
	profileGroup := new(errgroup.Group)
	for ni, cl := range c.Inputs.ChainlinkNodes {
		chainlinkNode := cl
		nodeIndex := ni
		profileGroup.Go(func() error {
			profileResults, err := chainlinkNode.Profile(c.Inputs.ProfileDuration, c.Inputs.ProfileFunction)
			profileResults.NodeIndex = nodeIndex
			if err != nil {
				return err
			}
			c.TestReporter.Results = append(c.TestReporter.Results, profileResults)
			return nil
		})
	}
	Expect(profileGroup.Wait()).ShouldNot(HaveOccurred(), "Error while gathering chainlink Profile tests")
}

// Networks returns the networks that the test is running on
func (c *ChainlinkProfileTest) TearDownVals() (*environment.Environment, []*client.Chainlink, reportModel.TestReporter, blockchain.EVMClient) {
	return c.env, c.Inputs.ChainlinkNodes, &c.TestReporter, c.c
}

// ensureValues ensures that all values needed to run the test are present
func (c *ChainlinkProfileTest) ensureInputValues() {
	Expect(c.Inputs.ProfileFunction).ShouldNot(BeNil(), "Forgot to provide a function to profile")
	Expect(c.Inputs.ProfileDuration.Seconds()).Should(BeNumerically(">=", 1), "Time to profile should be at least 1 second")
	Expect(c.Inputs.ChainlinkNodes).ShouldNot(BeNil(), "Chainlink nodes you want to profile should be provided")
	Expect(len(c.Inputs.ChainlinkNodes)).Should(BeNumerically(">", 0), "No Chainlink nodes provided to profile")
}
