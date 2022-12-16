package benchmark

import (
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	eth_contracts "github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

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

	NumberOfContracts, _    = strconv.Atoi(getEnv("NUMBEROFCONTRACTS", "500"))
	CheckGasToBurn, _       = strconv.ParseInt(getEnv("CHECKGASTOBURN", "100000"), 0, 64)
	PerformGasToBurn, _     = strconv.ParseInt(getEnv("PERFORMGASTOBURN", "50000"), 0, 64)
	BlockRange, _           = strconv.ParseInt(getEnv("BLOCKRANGE", "3600"), 0, 64)
	BlockInterval, _        = strconv.ParseInt(getEnv("BLOCKINTERVAL", "20"), 0, 64)
	NumberOfNodes, _        = strconv.Atoi(getEnv("AUTOMATION_NUMBER_OF_NODES", "6"))
	ChainlinkNodeFunding, _ = strconv.ParseFloat(getEnv("CHAINLINKNODEFUNDING", "0.5"), 64)
)

type BenchmarkTestEntry struct {
	registryVersions      []eth_contracts.KeeperRegistryVersion
	funding               *big.Float
	upkeepSLA             int64
	predeployedConsumers  []string
	upkeepResetterAddress string
	blockTime             time.Duration
}

func keeperBenchmark(t *testing.T, benchmarkTestEntry *BenchmarkTestEntry) {
	benchmarkNetwork := blockchain.LoadNetworkFromEnvironment()
	testEnvironment := environment.New(&environment.Config{InsideK8s: true})
	testEnvironment.
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: benchmarkNetwork.Name,
			Simulated:   benchmarkNetwork.Simulated,
			WsURLs:      benchmarkNetwork.URLs,
		}))
	for _, version := range benchmarkTestEntry.registryVersions {
		if version == eth_contracts.RegistryVersion_2_0 {
			NumberOfNodes++
		}
	}
	for i := 0; i < NumberOfNodes; i++ {
		testEnvironment.AddHelm(chainlink.New(i, nil))
	}
	err := testEnvironment.Run()
	require.NoError(t, err, "Error deploying test environment")
	log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Keepers Benchmark Environment")

	chainClient, err := blockchain.NewEVMClient(benchmarkNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	keeperBenchmarkTest := testsetups.NewKeeperBenchmarkTest(
		testsetups.KeeperBenchmarkTestInputs{
			BlockchainClient:  chainClient,
			NumberOfContracts: NumberOfContracts,
			RegistryVersions:  benchmarkTestEntry.registryVersions,
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
			ChainlinkNodeFunding:  benchmarkTestEntry.funding,
			UpkeepGasLimit:        5_000_000, //5M
			UpkeepSLA:             benchmarkTestEntry.upkeepSLA,
			FirstEligibleBuffer:   1,
			PreDeployedConsumers:  benchmarkTestEntry.predeployedConsumers,
			UpkeepResetterAddress: benchmarkTestEntry.upkeepResetterAddress,
			BlockTime:             benchmarkTestEntry.blockTime,
		},
	)
	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals(t)); err != nil {
			log.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})
	keeperBenchmarkTest.Setup(t, testEnvironment)
	keeperBenchmarkTest.Run(t)
}

func TestKeeperBenchmarkSimulatedGethRegistry_1_3(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
		big.NewFloat(100_000),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkGoerliRegistry_1_3(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
		big.NewFloat(ChainlinkNodeFunding),
		int64(4),
		predeployedConsumersGoerli,
		upkeepResetterContractGoerli,
		12 * time.Second,
	})
}

func TestKeeperBenchmarkArbitrumGoerliRegistry_1_3(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
		big.NewFloat(0.5),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkOptimisticGoerliRegistry_1_3(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_3},
		big.NewFloat(ChainlinkNodeFunding),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkSimulatedGethMulti_Registry(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
		big.NewFloat(100000),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkGoerliMulti_Registry(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2, eth_contracts.RegistryVersion_1_3},
		big.NewFloat(ChainlinkNodeFunding),
		int64(4),
		predeployedConsumersGoerli,
		upkeepResetterContractGoerli,
		time.Second,
	})
}

func TestKeeperBenchmarkSimulatedGethRegistry_1_2(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
		big.NewFloat(100000),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkGoerliRegistry_1_2(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_1_2},
		big.NewFloat(ChainlinkNodeFunding),
		int64(4),
		predeployedConsumersGoerli,
		upkeepResetterContractGoerli,
		12 * time.Second,
	})
}

func TestKeeperBenchmarkSimulatedGethRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(100000),
		int64(4),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		12 * time.Second,
	})
}

func TestKeeperBenchmarkGoerliRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(ChainlinkNodeFunding),
		int64(4),
		predeployedConsumersGoerli,
		upkeepResetterContractGoerli,
		12 * time.Second,
	})
}

func TestKeeperBenchmarkOptimisticGoerliRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(ChainlinkNodeFunding),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkArbitrumGoerliRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(0.5),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		time.Second,
	})
}

func TestKeeperBenchmarkMumbaiRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(ChainlinkNodeFunding),
		int64(20),
		predeployedConsumersEmpty,
		upkeepResetterContractEmpty,
		5 * time.Second,
	})
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
