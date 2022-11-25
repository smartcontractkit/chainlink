package benchmark

//revive:disable:dot-imports
import (
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

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
	simulatedBLockTime           = time.Second
	goerliTag                    = strings.ReplaceAll(strings.ToLower(networks.GoerliTestnet.Name), " ", "-")
	arbitrumTag                  = strings.ReplaceAll(strings.ToLower(networks.ArbitrumGoerli.Name), " ", "-")
	optimismTag                  = strings.ReplaceAll(strings.ToLower(networks.OptimismGoerli.Name), " ", "-")
	mumbaiTag                    = strings.ReplaceAll(strings.ToLower(networks.PolygonMumbai.Name), " ", "-")
)

type BenchmarkTestEntry struct {
	registryVersions      []eth_contracts.KeeperRegistryVersion
	funding               *big.Float
	upkeepSLA             int64
	predeployedConsumers  []string
	upkeepResetterAddress string
	blockTime             time.Duration
}

func getEnv(key, fallback string) string {
	if inputs, ok := os.LookupEnv("TEST_INPUTS"); ok {
		values := strings.Split(inputs, ",")
		for _, value := range values {
			if strings.Contains(value, key) {
				return strings.Split(value, "=")[1]
			}
		}
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
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Goerli Network @"+goerliTag+"-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(ChainlinkNodeFunding),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli,
					12 * time.Second},
			),
			Entry("Keeper benchmark suite on Arbitrum Goerli Network @"+arbitrumTag+"-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(0.5),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Optimistic Goerli Network @"+optimismTag+"-registry-1-3",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
					big.NewFloat(ChainlinkNodeFunding),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Simulated Network with Multiple Registries @simulated-multiple-registries",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
					big.NewFloat(100000),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Goerli Network with Multiple Registries @"+goerliTag+"-multiple-registries",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
					big.NewFloat(ChainlinkNodeFunding),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli,
					time.Second},
			),
			Entry("Keeper benchmark suite on Simulated Network with 1.2 registry @simulated-registry-1-2",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
					big.NewFloat(100000),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Goerli Network with 1.2 registry @"+goerliTag+"-registry-1-2",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
					big.NewFloat(ChainlinkNodeFunding),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli,
					12 * time.Second},
			),
			Entry("Keeper benchmark suite on Simulated Network with 2.0 registry @simulated-registry-2-0",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
					big.NewFloat(100000),
					int64(4),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					12 * time.Second},
			),
			Entry("Keeper benchmark suite on Goerli Network with 2.0 registry @"+goerliTag+"-registry-2-0",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
					big.NewFloat(ChainlinkNodeFunding),
					int64(4),
					predeployedConsumersGoerli,
					upkeepResetterContractGoerli,
					12 * time.Second},
			),
			Entry("Keeper benchmark suite on Arbitrum Goerli Network with 2.0 registry @"+arbitrumTag+"-registry-2-0",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
					big.NewFloat(0.5),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Optimistic Goerli Network with 2.0 registry @"+optimismTag+"-registry-2-0",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
					big.NewFloat(ChainlinkNodeFunding),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					time.Second},
			),
			Entry("Keeper benchmark suite on Mumbai Testnet Network with 2.0 registry @"+mumbaiTag+"-registry-2-0",
				BenchmarkTestEntry{[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
					big.NewFloat(ChainlinkNodeFunding),
					int64(20),
					predeployedConsumersEmpty,
					upkeepResetterContractEmpty,
					5 * time.Second},
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

	var NumberOfNodes, _ = strconv.Atoi(getEnv("AUTOMATION_NUMBER_OF_NODES", "6"))

	DescribeTable("Keeper benchmark suite on different EVM networks", func(
		testEntry BenchmarkTestEntry,
	) {
		By("Deploying the environment", func() {
			benchmarkNetwork = blockchain.LoadNetworkFromEnvironment()
			testEnvironment = environment.New(&environment.Config{InsideK8s: true})
			testEnvironment.
				AddHelm(ethereum.New(&ethereum.Props{
					NetworkName: benchmarkNetwork.Name,
					Simulated:   benchmarkNetwork.Simulated,
					WsURLs:      benchmarkNetwork.URLs,
				}))
			for i := 0; i < NumberOfNodes; i++ {
				testEnvironment.AddHelm(chainlink.New(i, nil))
			}
			err = testEnvironment.Run()
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
						MaxCheckDataSize:     uint32(5000),
						MaxPerformDataSize:   uint32(5000),
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
					BlockTime:             testEntry.blockTime,
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
