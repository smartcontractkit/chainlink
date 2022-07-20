package smoke

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
	it "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
)

var _ = Describe("Flux monitor suite @flux", func() {
	DescribeTable("Flux suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		networkChart environment.ConnectedChart,
		clChart environment.ConnectedChart,
	) {
		var (
			err              error
			c                blockchain.EVMClient
			cd               contracts.ContractDeployer
			lt               contracts.LinkToken
			fluxInstance     contracts.FluxAggregator
			chainlinkNodes   []client.Chainlink
			ms               *client.MockserverClient
			nodeAddresses    []common.Address
			adapterPath      string
			adapterUUID      string
			fluxRoundTimeout = 2 * time.Minute
			env              *environment.Environment
		)
		By("Deploying the environment", func() {
			env = environment.New(nil).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"replicas": "3",
				}))
			err = env.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			c, err = blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings)(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")

			cd, err = contracts.NewContractDeployer(c)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			nodeAddresses, err = actions.ChainlinkNodeAddresses(chainlinkNodes)
			Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
			ms, err = client.ConnectMockServer(env)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mock server client shouldn't fail")

			c.ParallelTransactions(true)
		})

		By("Setting initial adapter value", func() {
			adapterUUID = uuid.NewV4().String()
			adapterPath = fmt.Sprintf("/variable-%s", adapterUUID)
			err = ms.SetValuePath(adapterPath, 1e5)
			Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")
		})

		By("Deploying and funding contract", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			fluxInstance, err = cd.DeployFluxAggregatorContract(lt.Address(), contracts.DefaultFluxAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Flux Aggregator Contract shouldn't fail")
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for deployment of flux aggregator contract")

			err = lt.Transfer(fluxInstance.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred(), "Funding Flux Aggregator Contract shouldn't fail")
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for funding of flux aggregator contract")

			err = fluxInstance.UpdateAvailableFunds()
			Expect(err).ShouldNot(HaveOccurred(), "Updating the available funds on the Flux Aggregator Contract shouldn't fail")
		})

		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(chainlinkNodes, c, big.NewFloat(1))
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
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Getting oracle details from the Flux aggregator contract shouldn't fail")
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})

		By("Creating flux jobs", func() {
			adapterFullURL := fmt.Sprintf("%s%s", ms.Config.ClusterURL, adapterPath)
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
			By("Checking flux rounds", func() {
				// initial value set is performed before jobs creation
				fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(1), fluxRoundTimeout)
				c.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
				err = c.WaitForEvents()
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
				c.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
				err = ms.SetValuePath(adapterPath, 1e10)
				Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mock server shouldn't fail")
				err = c.WaitForEvents()
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
			})
			By("Printing gas stats", func() {
				c.GasStats().PrintStats()
			})
			By("Tearing down the environment", func() {
				err = actions.TeardownSuite(env, utils.ProjectRoot, chainlinkNodes, nil, c)
				Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
			})
		})
	},
		Entry("Flux on Geth @geth",
			blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings),
			ethereum.New(nil),
			chainlink.New(0, nil),
		),
	)
})
