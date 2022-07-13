// package vrfv2soak

// //revive:disable:dot-imports
// import (
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"time"

// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"github.com/rs/zerolog/log"
// 	"github.com/smartcontractkit/chainlink-testing-framework/actions"
// 	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
// 	"github.com/smartcontractkit/chainlink-testing-framework/config"
// 	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
// 	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
// 	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
// 	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
// 	"github.com/smartcontractkit/chainlink-testing-framework/utils"
// 	"github.com/smartcontractkit/helmenv/environment"
// )

// var _ = Describe("Vrfv2 soak test suite @soak_vrfv2", func() {
// 	var (
// 		err           error
// 		env           *environment.Environment
// 		vrfv2SoakTest *testsetups.VRFV2SoakTest
// 		coordinator   contracts.VRFCoordinatorV2
// 		consumer      contracts.VRFConsumerV2
// 		// jobInfo       []testsetups.VRFV2SoakTestJobInfo
// 		jobInfo []actions.VRFV2JobInfo
// 	)
// 	BeforeEach(func() {
// 		By("Deploying the environment", func() {
// 			// load unique namspace if provided
// 			namespace, ok := os.LookupEnv("SOAK_NAMESPACE")
// 			if !ok {
// 				namespace = "local"
// 			}
// 			c := environment.NewChainlinkConfig(
// 				config.ChainlinkVals(),
// 				fmt.Sprintf("%s-vrfv2-soak", namespace),
// 				// works only on perf Geth
// 				environment.PerformanceGeth,
// 			)
// 			chart, err := c.Charts.Get("chainlink")
// 			log.Error().Interface("charts", c.Charts).Msg("what charts are there")
// 			Expect(err).ShouldNot(HaveOccurred())
// 			chart.Values["replicas"] = 2
// 			env, err = environment.DeployOrLoadEnvironment(
// 				c,
// 			)
// 			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
// 			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
// 		})

// 		By("Setup the Vrfv2 test", func() {
// 			vrfv2SoakTest = &testsetups.VRFV2SoakTest{
// 				Inputs: &testsetups.VRFV2SoakTestInputs{
// 					TestDuration:         time.Minute * 1,
// 					ChainlinkNodeFunding: big.NewFloat(1000),
// 					StopTestOnError:      false,

// 					RequestsPerMinute: 8,

// 					// Make the test simple and just request randomness and return any errors
// 					TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) error {
// 						words := uint32(10)
// 						err := consumer.RequestRandomness(jobInfo[0].ProvingKeyHash, 1, 1, 300000, words)
// 						return err
// 					},
// 				},
// 				TestReporter: testreporters.VRFV2SoakTestReporter{
// 					Reports: make(map[string]*testreporters.VRFV2SoakTestReport),
// 				},
// 			}

// 			vrfv2SoakTest.Setup(env, []string{"chainlink"}, true)
// 		})

// 		By("Deploy Contracts", func() {
// 			// With the environment setup now we can deploy contracts and jobs
// 			contractDeployer, err := contracts.NewContractDeployer(vrfv2SoakTest.DefaultNetwork)
// 			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

// 			// Deploy LINK
// 			linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
// 			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

// 			// Fund Chainlink nodes
// 			err = actions.FundChainlinkNodes(vrfv2SoakTest.ChainlinkNodes, vrfv2SoakTest.DefaultNetwork, vrfv2SoakTest.Inputs.ChainlinkNodeFunding)
// 			Expect(err).ShouldNot(HaveOccurred())

// 			linkEthFeedResponse := big.NewInt(1e18)
// 			mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
// 			Expect(err).ShouldNot(HaveOccurred())
// 			coordinator, consumer, _ = actions.DeployVRFV2Contracts(linkTokenContract, contractDeployer, vrfv2SoakTest.Networks, mf.Address())

// 			evmClient := vrfv2SoakTest.Networks.Default.(*blockchain.EthereumMultinodeClient).DefaultClient
// 			temp := &contracts.EthereumVRFConsumerV2{}
// 			temp.LoadFromAddress(consumer.Address(), evmClient)
// 			consumer = temp

// 			err = linkTokenContract.Transfer(consumer.Address(), big.NewInt(0).Mul(big.NewInt(1e4), big.NewInt(1e18)))
// 			Expect(err).ShouldNot(HaveOccurred())
// 			err = coordinator.SetConfig(
// 				1,
// 				2.5e6,
// 				86400,
// 				33825,
// 				linkEthFeedResponse,
// 				ethereum.VRFCoordinatorV2FeeConfig{
// 					FulfillmentFlatFeeLinkPPMTier1: 1,
// 					FulfillmentFlatFeeLinkPPMTier2: 1,
// 					FulfillmentFlatFeeLinkPPMTier3: 1,
// 					FulfillmentFlatFeeLinkPPMTier4: 1,
// 					FulfillmentFlatFeeLinkPPMTier5: 1,
// 					ReqsForTier2:                   big.NewInt(10),
// 					ReqsForTier3:                   big.NewInt(20),
// 					ReqsForTier4:                   big.NewInt(30),
// 					ReqsForTier5:                   big.NewInt(40)},
// 			)
// 			Expect(err).ShouldNot(HaveOccurred())
// 			err = vrfv2SoakTest.Networks.Default.WaitForEvents()
// 			Expect(err).ShouldNot(HaveOccurred())

// 			err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(big.NewInt(30), big.NewInt(1e18)))
// 			Expect(err).ShouldNot(HaveOccurred())
// 			err = vrfv2SoakTest.Networks.Default.WaitForEvents()
// 			Expect(err).ShouldNot(HaveOccurred())

// 			jobInfo = actions.CreateVRFV2Jobs(vrfv2SoakTest.ChainlinkNodes, coordinator, vrfv2SoakTest.Networks, 1)

// 			err = vrfv2SoakTest.DefaultNetwork.WaitForEvents()
// 			Expect(err).ShouldNot(HaveOccurred())
// 		})
// 	})
// 	Describe("Run the test", func() {
// 		It("Makes requests for randomness and verifies number of jobs have been run", func() {
// 			vrfv2SoakTest.Run()
// 		})
// 	})

// 	AfterEach(func() {
// 		By("Tearing down the environment", func() {
// 			_, nets, _, _ := vrfv2SoakTest.TearDownVals()
// 			err = actions.TeardownSuite(env, nets, utils.ProjectRoot, nil, nil)
// 			Expect(err).ShouldNot(HaveOccurred())
// 		})
// 	})
// })
