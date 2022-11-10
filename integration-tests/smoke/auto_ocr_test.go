package smoke

import (
	"context"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Automation OCR Suite @auto-ocr", func() {
	var (
		err              error
		chainClient      blockchain.EVMClient
		contractDeployer contracts.ContractDeployer
		registry         contracts.KeeperRegistry
		registrar        contracts.KeeperRegistrar
		consumers        []contracts.KeeperConsumer
		upkeepIDs        []*big.Int
		linkToken        contracts.LinkToken
		chainlinkNodes   []*client.Chainlink
		testEnvironment  *environment.Environment

		testScenarios = []TableEntry{
			Entry("v2.0 Basic smoke test @simulated", ethereum.RegistryVersion_2_0, defaultRegistryConfig, BasicCounter, BasicSmokeTest, big.NewInt(defaultLinkFunds), 2),
		}
	)

	DescribeTable("Automation OCR Suite @auto-ocr", func(
		registryVersion ethereum.KeeperRegistryVersion,
		registryConfig contracts.KeeperRegistrySettings,
		consumerContract KeeperConsumerContracts,
		testToRun KeeperTests,
		linkFundsForEachUpkeep *big.Int,
		numberOfUpkeeps int,
	) {
		By("Deploying the environment")
		network := networks.SimulatedEVM
		chainlinkTOML := client.NewDefaultTOMLBuilder().
			AddNetworks(network).
			AddOCR2Defaults().
			AddP2PNetworkingV2().
			String()
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-auto-ocr"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(eth.New(nil)).
			AddHelm(chainlink.New(0, map[string]interface{}{
				"replicas": "5",
				"env": map[string]interface{}{
					"cl_config": chainlinkTOML,
				},
			}))
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(network, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		contractDeployer, err = contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		chainClient.ParallelTransactions(true)

		By("Funding Chainlink nodes")
		txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
		Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
		Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")

		By("Deploy Link Token Contract")
		linkToken, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		By("Deploy Registry and Registrar")
		registry, registrar = actions.DeployAutoOCRRegistryAndRegistrar(
			registryVersion,
			registryConfig,
			numberOfUpkeeps,
			linkToken,
			contractDeployer,
			chainClient,
		)

		By("Create OCR Automation Jobs")
		actions.CreateOCRKeeperJobs(chainlinkNodes, registry.Address(), network.ChainID)
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig := actions.BuildAutoOCR2ConfigVars(nodesWithoutBootstrap, registryConfig, registrar.Address())
		err = registry.SetConfig(defaultRegistryConfig, ocrConfig)
		Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")

		By("Deploy Consumers")
		switch consumerContract {
		case BasicCounter:
			consumers, upkeepIDs = actions.DeployConsumers(
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				numberOfUpkeeps,
				linkFundsForEachUpkeep,
				defaultUpkeepGasLimit,
			)
		}

		if testToRun == BasicSmokeTest {
			By("watches all the registered upkeeps perform and then cancels them from the registry")
			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">=", int64(5)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(Succeed())

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.CancelUpkeep(upkeepIDs[i])
				Expect(err).ShouldNot(HaveOccurred(), "Could not cancel upkeep at index "+strconv.Itoa(i))
			}

			err := chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be cancelled")

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				log.Info().Msg("Cancelled upkeep at index " + strconv.Itoa(i) + " which performed " +
					strconv.Itoa(int(countersAfterCancellation[i].Int64())) + " times")
			}

			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterCancellation[i].Int64()+1),
						"Expected consumer counter to remain constant at %d, but got %d", countersAfterCancellation[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(Succeed())
		}

		By("Printing gas stats")
		chainClient.GasStats().PrintStats()

		By("Tearing down the environment")
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	},
		testScenarios,
	)
})
