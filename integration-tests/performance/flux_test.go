package performance

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Flux monitor suite @flux", func() {
	var (
		err              error
		chainClient      blockchain.EVMClient
		contractDeployer contracts.ContractDeployer
		linkToken        contracts.LinkToken
		fluxInstance     contracts.FluxAggregator
		chainlinkNodes   []client.Chainlink
		mockServer       *ctfClient.MockserverClient
		nodeAddresses    []common.Address
		adapterPath      string
		adapterUUID      string
		fluxRoundTimeout = 2 * time.Minute
		testEnvironment  *environment.Environment
		profileTest      *testsetups.ChainlinkProfileTest
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			testEnvironment = environment.New(&environment.Config{NamespacePrefix: "performance-flux"}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"replicas": "3",
					"env": map[string]interface{}{
						"HTTP_SERVER_WRITE_TIMEOUT": "300s",
					},
				}))
			err = testEnvironment.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			chainClient, err = blockchain.NewEthereumMultiNodeClientSetup(blockchain.SimulatedEVMNetwork)(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")

			contractDeployer, err = contracts.NewContractDeployer(chainClient)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			nodeAddresses, err = actions.ChainlinkNodeAddresses(chainlinkNodes)
			Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
			mockServer, err = ctfClient.ConnectMockServer(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mock server client shouldn't fail")

			chainClient.ParallelTransactions(true)
		})

		By("Setting initial adapter value", func() {
			adapterUUID = uuid.NewV4().String()
			adapterPath = fmt.Sprintf("/variable-%s", adapterUUID)
			err = mockServer.SetValuePath(adapterPath, 1e5)
			Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")
		})

		By("Deploying and funding contract", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			fluxInstance, err = contractDeployer.DeployFluxAggregatorContract(linkToken.Address(), contracts.DefaultFluxAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Flux Aggregator Contract shouldn't fail")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for deployment of flux aggregator contract")

			err = linkToken.Transfer(fluxInstance.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred(), "Funding Flux Aggregator Contract shouldn't fail")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for funding of flux aggregator contract")

			err = fluxInstance.UpdateAvailableFunds()
			Expect(err).ShouldNot(HaveOccurred(), "Updating the available funds on the Flux Aggregator Contract shouldn't fail")
		})

		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred(), "Funding chainlink nodes with ETH shouldn't fail")
		})

		By("Setting oracle options", func() {
			err = fluxInstance.SetOracles(
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				})
			Expect(err).ShouldNot(HaveOccurred(), "Setting oracle options in the Flux Aggregator contract shouldn't fail")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Getting oracle details from the Flux aggregator contract shouldn't fail")
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})

		By("Creating flux jobs", func() {
			adapterFullURL := fmt.Sprintf("%s%s", mockServer.Config.ClusterURL, adapterPath)
			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("variable-%s", adapterUUID),
				URL:  adapterFullURL,
			}
			for i, n := range chainlinkNodes {
				err = n.CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred(), "Creating bridge shouldn't fail for node %d", i+1)

				fluxSpec := &client.FluxMonitorJobSpec{
					Name:              fmt.Sprintf("flux-monitor-%s", adapterUUID),
					ContractAddress:   fluxInstance.Address(),
					Threshold:         0,
					AbsoluteThreshold: 0,
					PollTimerPeriod:   15 * time.Second, // min 15s
					IdleTimerDisabled: true,
					ObservationSource: client.ObservationSourceSpecBridge(bta),
				}
				_, err = n.CreateJob(fluxSpec)
				Expect(err).ShouldNot(HaveOccurred(), "Creating flux job shouldn't fail for node %d", i+1)
			}
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
					// Not the last node, hence not all nodes started profiling yet.
					return
				}
				// initial value set is performed before jobs creation
				fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(1), fluxRoundTimeout)
				chainClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
				err = chainClient.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
				data, err := fluxInstance.GetContractData(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Getting contract data from flux aggregator contract shouldn't fail")
				log.Info().Interface("Data", data).Msg("Round data")
				Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(1e5)), "Expected latest round answer to be %d, but found %d", int64(1e5), data.LatestRoundData.Answer.Int64())
				Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(1)), "Expected latest round id to be %d, but found %d", int64(1), data.LatestRoundData.RoundId.Int64())
				Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(1)), "Expected latest round's answered in round to be %d, but found %d", int64(1), data.LatestRoundData.AnsweredInRound.Int64())
				Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999997)), "Expected available funds to be %d, but found %d", int64(999999999999999997), data.AvailableFunds.Int64())
				Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(3)), "Expected allocated funds to be %d, but found %d", int64(3), data.AllocatedFunds.Int64())

				fluxRound = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
				chainClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
				err = mockServer.SetValuePath(adapterPath, 1e10)
				Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mock server shouldn't fail")
				err = chainClient.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
				data, err = fluxInstance.GetContractData(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Getting contract data from flux aggregator contract shouldn't fail")
				Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(1e10)), "Expected latest round answer to be %d, but found %d", int64(1e10), data.LatestRoundData.Answer.Int64())
				Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(2)), "Expected latest round id to be %d, but found %d", int64(2), data.LatestRoundData.RoundId.Int64())
				Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(2)), "Expected latest round's answered in round to be %d, but found %d", int64(2), data.LatestRoundData.AnsweredInRound.Int64())
				Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999994)), "Expected available funds to be %d, but found %d", int64(999999999999999994), data.AvailableFunds.Int64())
				Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(6)), "Expected allocated funds to be %d, but found %d", int64(6), data.AllocatedFunds.Int64())
				log.Info().Interface("data", data).Msg("Round data")

				for _, oracleAddr := range nodeAddresses {
					payment, _ := fluxInstance.WithdrawablePayment(context.Background(), oracleAddr)
					Expect(payment.Int64()).Should(Equal(int64(2)), "Expected flux aggregator contract's withdrawable payment to be %d, but found %d", int64(2), payment.Int64())
				}
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: 30 * time.Second,
				ChainlinkNodes:  chainlinkNodes,
			})
			profileTest.Setup(testEnvironment)
		})
	})

	Describe("with Flux job", func() {
		It("performs two rounds and has withdrawable payments for oracles", func() {
			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			chainClient.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, &profileTest.TestReporter, chainClient)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
