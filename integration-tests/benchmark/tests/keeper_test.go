package benchmark_test

//revive:disable:dot-imports
import (
	"math/big"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Keeper benchmark suite @benchmark-keeper", func() {
	var (
		err                 error
		testEnvironment     *environment.Environment
		keeperBenchmarkTest *testsetups.KeeperBenchmarkTest
		benchmarkNetwork    *blockchain.EVMNetwork
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			benchmarkNetwork = blockchain.LoadNetworkFromEnvironment()
			testEnvironment = environment.New(&environment.Config{InsideK8s: true})
			err = testEnvironment.
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(&ethereum.Props{
					NetworkName: benchmarkNetwork.Name,
					Simulated:   benchmarkNetwork.Simulated,
				})).
				AddHelm(chainlink.New(0, nil)).
				Run()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")
		})

		By("Setup the Keeper test", func() {
			chainClient, err := blockchain.NewEthereumMultiNodeClientSetup(blockchain.SimulatedEVMNetwork)(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			keeperBenchmarkTest = testsetups.NewKeeperBenchmarkTest(
				testsetups.KeeperBenchmarkTestInputs{
					BlockchainClient:  chainClient,
					NumberOfContracts: 500, // TODO- Update to 500
					KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
						PaymentPremiumPPB:    uint32(0),
						BlockCountPerTurn:    big.NewInt(100),
						CheckGasLimit:        uint32(2000000),
						StalenessSeconds:     big.NewInt(90000),
						GasCeilingMultiplier: uint16(2),
						MaxPerformGas:        uint32(5000000),
						MinUpkeepSpend:       big.NewInt(0),
						FallbackGasPrice:     big.NewInt(2e11),
						FallbackLinkPrice:    big.NewInt(2e18),
					},
					CheckGasToBurn:       100000,
					PerformGasToBurn:     150000,
					BlockRange:           3600, // TODO- Update to 3600
					BlockInterval:        20,
					ChainlinkNodeFunding: big.NewFloat(1000),
					UpkeepGasLimit:       500000,
				},
			)
			keeperBenchmarkTest.Setup(testEnvironment)
		})
	})

	Describe("Watching the keeper contracts to ensure they reply in time", func() {
		It("Watches for Upkeep counts", func() {
			keeperBenchmarkTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if err := actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
		})
	})
})
