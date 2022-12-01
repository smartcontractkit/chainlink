package reorg

//revive:disable:dot-imports
import (
	"context"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"math/big"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[Keeper]
TurnFlagEnabled = true
TurnLookBack = 0

[Keeper.Registry]
SyncInterval = '5m'
PerformGasOverhead = 150_000

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
	networkTOML = `FinalityDepth = 200

[EVM.HeadTracker]
HistoryDepth = 400`
	activeEVMNetwork          = networks.SelectedNetwork
	defaultAutomationSettings = map[string]interface{}{
		"toml":     client.AddNetworkDetailedConfig(baseTOML, networkTOML, activeEVMNetwork),
		"replicas": "6",
		"db": map[string]interface{}{
			"stateful": false,
			"capacity": "10Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
			},
		},
	}

	defaultReorgEthereumSettings = &reorg.Props{
		NetworkName: activeEVMNetwork.Name,
		NetworkType: "geth-reorg",
		Values: map[string]interface{}{
			"geth": map[string]interface{}{
				"genesis": map[string]interface{}{
					"networkId": "1337",
				},
			},
		},
	}

	defaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(10),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5000),
		MaxPerformDataSize:   uint32(5000),
	}
)

const (
	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
)

var _ = Describe("Automation reorg test @reorg-automation", func() {
	numberOfUpkeeps := 2
	reorgBlocks := 50
	var (
		testScenarios = []TableEntry{
			Entry("Must survive 50 block reorg for 1m @reorg-automation-50-block",
				reorg.New(defaultReorgEthereumSettings),
				chainlink.New(0, defaultAutomationSettings),
			),
		}

		testEnvironment *environment.Environment
		chainlinkNodes  []*client.Chainlink
		chainClient     blockchain.EVMClient
		registry        contracts.KeeperRegistry
		registrar       contracts.KeeperRegistrar
		consumers       []contracts.KeeperConsumer
		upkeepIDs       []*big.Int
	)

	AfterEach(func() {
		err := actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("Automation reorg on different EVM networks", func(
		networkChart environment.ConnectedChart,
		clChart environment.ConnectedChart,
	) {
		By("Deploying the environment")
		testEnvironment = environment.
			New(&environment.Config{
				NamespacePrefix: "reorg-automation",
				TTL:             time.Hour * 1}).
			AddHelm(networkChart).
			AddHelm(clChart).
			AddChart(blockscout.New(&blockscout.Props{
				Name:    "geth-blockscout",
				WsURL:   activeEVMNetwork.URL,
				HttpURL: activeEVMNetwork.HTTPURLs[0]}))
		err := testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(activeEVMNetwork, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		contractDeployer, err := contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")

		chainClient.ParallelTransactions(true)

		linkToken, err := contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		By("Funding Chainlink nodes")
		txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
		Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
		Expect(err).ShouldNot(HaveOccurred())

		By("Deploy Registry and Registrar")
		registry, registrar = actions.DeployAutoOCRRegistryAndRegistrar(
			ethereum.RegistryVersion_2_0,
			defaultOCRRegistryConfig,
			numberOfUpkeeps,
			linkToken,
			contractDeployer,
			chainClient,
		)

		By("Create OCR Automation Jobs")
		actions.CreateOCRKeeperJobs(chainlinkNodes, registry.Address(), activeEVMNetwork.ChainID, 0)
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig := actions.BuildAutoOCR2ConfigVars(nodesWithoutBootstrap, defaultOCRRegistryConfig, registrar.Address(), 5*time.Second)
		err = registry.SetConfig(defaultOCRRegistryConfig, ocrConfig)
		Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
		Expect(chainClient.WaitForEvents()).ShouldNot(HaveOccurred(), "Waiting for config to be set")

		By("Deploy Consumers")
		consumers, upkeepIDs = actions.DeployConsumers(
			registry,
			registrar,
			linkToken,
			contractDeployer,
			chainClient,
			numberOfUpkeeps,
			big.NewInt(defaultLinkFunds),
			defaultUpkeepGasLimit,
		)

		By("watches all the registered upkeeps perform and then cancels them from the registry")
		Eventually(func(g Gomega) {
			// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
			for i := 0; i < len(upkeepIDs); i++ {
				counter, err := consumers[i].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				expect := 5
				g.Expect(counter.Int64()).Should(BeNumerically(">=", int64(expect)),
					"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
			}
		}, "5m", "1s").Should(Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

		By("creating reorg for 50 blocks", func() {
			rc, err := NewReorgController(
				&ReorgConfig{
					FromPodLabel:            reorg.TXNodesAppLabel,
					ToPodLabel:              reorg.MinerNodesAppLabel,
					Network:                 chainClient,
					Env:                     testEnvironment,
					BlockConsensusThreshold: 3,
					Timeout:                 1800 * time.Second,
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			rc.ReOrg(reorgBlocks)
			rc.WaitReorgStarted()

			err = rc.WaitDepthReached()
			Expect(err).ShouldNot(HaveOccurred())
		})

		Eventually(func(g Gomega) {
			// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
			for i := 0; i < len(upkeepIDs); i++ {
				counter, err := consumers[i].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				expect := 10
				g.Expect(counter.Int64()).Should(BeNumerically(">=", int64(expect)),
					"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
			}
		}, "15m", "1s").Should(Succeed())

		//// Cancel all the registered upkeeps via the registry
		//for i := 0; i < len(upkeepIDs); i++ {
		//	err := registry.CancelUpkeep(upkeepIDs[i])
		//	Expect(err).ShouldNot(HaveOccurred(), "Could not cancel upkeep at index "+strconv.Itoa(i))
		//}
		//
		//err = chainClient.WaitForEvents()
		//Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be cancelled")
		//
		//var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))
		//
		//for i := 0; i < len(upkeepIDs); i++ {
		//	// Obtain the amount of times the upkeep has been executed so far
		//	countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
		//	Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
		//	log.Info().Msg("Cancelled upkeep at index " + strconv.Itoa(i) + " which performed " +
		//		strconv.Itoa(int(countersAfterCancellation[i].Int64())) + " times")
		//}
		//
		//Consistently(func(g Gomega) {
		//	for i := 0; i < len(upkeepIDs); i++ {
		//		// Expect the counter to remain constant (At most increase by 1 to account for stale performs) because the upkeep was cancelled
		//		latestCounter, err := consumers[i].Counter(context.Background())
		//		Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
		//		g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterCancellation[i].Int64()+1),
		//			"Expected consumer counter to remain less than equal %d, but got %d", countersAfterCancellation[i].Int64()+1, latestCounter.Int64())
		//	}
		//}, "1m", "1s").Should(Succeed())
	},
		testScenarios,
	)
})
