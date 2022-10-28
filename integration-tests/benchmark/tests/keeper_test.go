package benchmark

//revive:disable:dot-imports
import (
	"math/big"
	"os"
	"strconv"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	eth_contracts "github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var (
	predeployedConsumersEmpty    []string
	predeployedConsumersGoerli   = []string{""} //Copy Addresses here before you run
	upkeepResetterContractEmpty  = ""
	upkeepResetterContractGoerli = "0xaeA9bD8f60C9EB1771900B9338dE8Ab52584E80e"
)

type BenchmarkTestEntry struct {
	registryVersions      []eth_contracts.KeeperRegistryVersion
	funding               *big.Float
	upkeepSLA             int64
	predeployedConsumers  []string
	upkeepResetterAddress string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var _ = Describe("Keeper benchmark suite @benchmark-keeper", func() {
	var (
		ChainlinkNodeFunding, _ = strconv.ParseFloat(getEnv("CHAINLINKNODEFUNDING", "0.5"), 64)
		testScenarios           = []TableEntry{
			Entry("Keeper benchmark suite on Simulated Network @simulated-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(100000),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty},
			),
			Entry("Keeper benchmark suite on Goerli Network @goerli-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(0.5),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli},
			),
			Entry("Keeper benchmark suite on Arbitrum Goerli Network @arbitrum-goerli-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(0.5),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty},
			),
			Entry("Keeper benchmark suite on Optimistic Goerli Network @optimistic-goerli-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(0.5),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty},
			),
			Entry("Keeper benchmark suite on Simulated Network with Multiple Registries @simulated-multiple-registries",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
					big.NewFloat(100000),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty},
			),
			Entry("Keeper benchmark suite on Goerli Network with Multiple Registries @goerli-multiple-registries",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
					big.NewFloat(0.5),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli},
			),
			Entry("Keeper benchmark suite on Simulated Network with 1.2 registry @simulated-registry1-2",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
					big.NewFloat(100000),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty},
			),
			Entry("Keeper benchmark suite on Goerli Network with 1.2 registry @goerli-registry1-2",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
					big.NewFloat(ChainlinkNodeFunding),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli},
			),
		}
		err                 error
		testEnvironment     *environment.Environment
		keeperBenchmarkTest *testsetups.KeeperBenchmarkTest
		benchmarkNetwork    *blockchain.EVMNetwork
	)

	var NumberOfContracts, _ = strconv.Atoi(getEnv("NUMBEROFCONTRACTS", "500"))
	var CheckGasToBurn, _ = strconv.ParseInt(getEnv("CHECKGASTOBURN", "100000"), 0, 64)
	var PerformGasToBurn, _ = strconv.ParseInt(getEnv("PERFORMGASTOBURN", "50000"), 0, 64)
	var BlockRange, _ = strconv.ParseInt(getEnv("BLOCKRANGE", "3600"), 0, 64)
	var BlockInterval, _ = strconv.ParseInt(getEnv("BLOCKINTERVAL", "20"), 0, 64)

	DescribeTable("Keeper benchmark suite on different EVM networks", func(
		testEntry BenchmarkTestEntry,
	) {
		By("Deploying the environment", func() {
			benchmarkNetwork = blockchain.LoadNetworkFromEnvironment()
			testEnvironment = environment.New(&environment.Config{InsideK8s: true})
			err = testEnvironment.
				AddHelm(ethereum.New(&ethereum.Props{
					NetworkName: benchmarkNetwork.Name,
					Simulated:   benchmarkNetwork.Simulated,
					WsURLs:      benchmarkNetwork.URLs,
				})).
				AddHelm(chainlink.New(0, nil)).
				AddHelm(chainlink.New(1, nil)).
				AddHelm(chainlink.New(2, nil)).
				AddHelm(chainlink.New(3, nil)).
				AddHelm(chainlink.New(4, nil)).
				AddHelm(chainlink.New(5, nil)).
				Run()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")
		})

		By("Setup the Keeper test", func() {
			chainClient, err := blockchain.NewEVMClient(benchmarkNetwork, testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			keeperBenchmarkTest = testsetups.NewKeeperBenchmarkTest(
				testsetups.KeeperBenchmarkTestInputs{
					BlockchainClient:  chainClient,
					NumberOfContracts: NumberOfContracts,
					RegistryVersions:  testEntry.registryVersions,
					KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
						PaymentPremiumPPB:    uint32(0),
						BlockCountPerTurn:    big.NewInt(100),
						CheckGasLimit:        uint32(45000000), //45M
						StalenessSeconds:     big.NewInt(90000),
						GasCeilingMultiplier: uint16(2),
						MaxPerformGas:        uint32(5000000), //5M
						MinUpkeepSpend:       big.NewInt(0),
						FallbackGasPrice:     big.NewInt(2e11),
						FallbackLinkPrice:    big.NewInt(2e18),
					},
					CheckGasToBurn:        CheckGasToBurn,
					PerformGasToBurn:      PerformGasToBurn,
					BlockRange:            BlockRange,
					BlockInterval:         BlockInterval,
					ChainlinkNodeFunding:  testEntry.funding,
					UpkeepGasLimit:        5000000, //5M
					UpkeepSLA:             testEntry.upkeepSLA,
					FirstEligibleBuffer:   1,
					PreDeployedConsumers:  testEntry.predeployedConsumers,
					UpkeepResetterAddress: testEntry.upkeepResetterAddress,
				},
			)
			keeperBenchmarkTest.Setup(testEnvironment)
		})

		By("Watches for Upkeep counts", func() {
			keeperBenchmarkTest.Run()
		})
	}, testScenarios)

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if err := actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
			log.Info().Msg("Keepers Benchmark Test Concluded")
		})
	})
})
