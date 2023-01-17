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
	predeployedConsumersMumbai   = []string{"0x1920a1eeD8fbc8734Aa60b858e2Ac83F5e3ceF74,0x6741269BBA28FEF74B48Fbdc8b30F5DA3A9AD1Ff,0x7a5F47E27f7A1F7CF216De4505C357B9DED7Fb13,0x5A51ba3DE0aC3EcD9CCaa566EFAD403591520d43,0xc5f492a89A32f149E6493B72a14490f8ad0a1E4F,0x6896EeeB82A2E8BC75FE2913a5C57F337106F6eC,0xd5fA88a1C4B386A3aEBb10d379521d8E8CdA4e39,0xF4e27c7E8541673df69dA652d1b5bF9BCB7E9067,0xFdf61ebC86A99A7699b07acF19bAEAe8582823fB,0x987824ECE90aaD82091F595EFf8e8529cbEce4F8,0x18e5D37584b71f5537F0016736373FEF963991fE,0x52Ad83c7aB2ffc1590e89Ac816bdbaca4036E703,0x0edc9DCB830D9AE868d75c16E14d222fAf81b6bF,0x91477c9790739d4E55C533cB7209Ead3063bFa5D,0x7288d4aF0Ac30852f2EA03de2830eE42DB4636C3,0x21EB8BFC8915236B704a3eBe56Ad864D9265Bc71,0x1dB0aec30d69a3c740bd0F37C1EA9E1cC25ca082,0x1692c019549E2377fe3d8BDd009d6B90977706A8,0xB0a0c23f7b85c21306A7b6dc9df46565A556c048,0xcA6FB01bc5942AC16095Ed5D575763575A07c8C1,0x73aFCc1D1f44548CF1C94a4BDeBE48A343546E05,0xF324B419c9b00b8f2a2cB5Cacd0D997e30d75066,0x251E77AC8E7D619932916B783E4D6998D5bCCd9a,0xb47445cdA4EBd68846d53f2C39EDD80b21C73349,0x899CC161eF8c5D9015A204f63503f28d19A3dfd1,0x47E93cEF623bC08D0DeCDeDD97bE901d1259757D,0x86379898E217d363E68601943263ed90Ee9d22A5,0x27d34be818CAE5fDb0B6e4f58c1a24a6a076A32b,0xF3494C25F8F968b988BDA6A5754FcBcbBdb5eD30,0xBE1f2664Fd9D2dec28E670935E3D7A0bf69675F8,0x66442f0F19C8E700722F872596AeF9D1F316Da32,0x3F9C327eDdB120ff8078aA77f5eEDA591cEB5BDA,0xC4799bAff895fa4141735CC8431c4a37646226EB,0x8712Edcf082B73535F4115E71B7656439B63D18C,0xF799f46304aa0B2F10CfeA1586857c5135b8f1c5,0xf8a10F4057AAF2939A6C4d8918011eFaAF3EAeF6,0x8964B38E3352D792dCAb28C0653976Eb2a044e96,0x60F27C802eb6aA50203c5EFA439AdC4A3F8f2035,0xdEb1D86B06877Bc47cC79ED8A1Fa491eD1e225ad,0xEFb01BD55b0e2756c207FEe7D8844D5e7F9Ac8Db,0x978dd672b17E56F5e0C2391929158DAee75c2Eb2,0x597DB534629057f54080F99e3353C413722CE828,0xf080aFB6616a76B146A46C8a4DEee62cC4bd0486,0xC5b8f10EaD40Bad23121029b439C2C8b2C17203B,0xC623f61344c5c06Dce46F0f679bDAf1F30d3980e,0xB94533dfcAB111eE776c83831DeA0c5B7fe5A893,0x84CA2Ca24ec13df17AB94816C60cEA9F463bD0D2,0x2a992081EEEC6b956E0b0bc9D78B84dEBe9DE16c,0xd6AC3fd44CD2afd001273B5EAebd8b82Ec498c41,0xC8A381751939672CE8fe5503b791B52F23812851,0x5b0121f62fE120dA1D84fFAaaba5F1a36d813046,0xB78bAE0D3b2c5550639cF7CA462932FadAE708Bd,0x955B6043382129677FEeFC8590052913d02994b4,0xb17e67E25f1fF3d7607Dff5bb673Ad15E2114A3A,0x2a66F326a80839AD3efA5c216b0fD81a0D712f76,0xe3AbaA41CcD2AC49e16818a9644e41D95C37cff4,0x0B25235153236A346D8e650700599a247817b222,0xc8a6265f588b463e72957f898425c1e9e1dac1Cc,0xb8Ed343D3F22C6a21849D5c3EAAa08fAc7EB1011,0x6ffea5D478f4cF978F489d9FC23849135389201d,0xC4C3dC278497335b41422C6707969F96857d7665,0x99a51EC538d0B0CB89F16ee95E5Bc910506D3E64,0x06Eb89aD5D63DDf95eDd593a7699c4B5A2052126,0x891C123DFc5600B52B40927Cbf453F27A6e3BBf6,0x5Ed33F4a46986E42192Da4EF285B7368B320424C,0x333C9D65AAb8bC6e46CF37F1ff118132464BB482,0xF9b2FaBFaeE22B3BA7193C7f05E1AD37ea6e7858,0x20e6533c33888b06206523DF70A31296f2103B81,0xa25e96ade301ADA9860644DD7269C3cAda09944a,0xD322C4F5856C230c64A0e76126F0285Cb619b321,0x06fa6b336f8E0cEB3bae501410f5611275315772,0x21c7A8811F4A8C4c776649426A8155721f7003bF,0x8eE52d3Ac27296c7f675B992C7076e8d8fF220DA,0x8fC36b5ba58D17624C579e86001BF92AC234a3C1,0x2Ddd8cdEf7966A9940D623326b597aB8e970eF77,0xAbb52e9762A3d5c58E3394931380A7Dabcb9f0Eb,0x43ada25d32D9E959c4edEa5F126e19914B903e4B,0x35A097040C75641d4d7797C39a7Ebfc9B922d2A8,0xB658691AF154881b2d8dFd4822b085e53A6f4c8a,0x760f7a27285d8f032CA1639CE35456DbA1162326,0x77C680D247AE34C0e9B4F4d216EA6c15c30255B6,0x2b737979e92e3ab7abb5Ec35a85EF8a3F1E6cEe0,0x330b38Fd020c239cE1FFe06DB3419899Ec9FB881,0xd5d86985fE0fE80C0D3c93D025ceE83093DA583B,0xBA22a0d2A43520109D18c676245723E8F85c7079,0x83A1Ab2fC42A50499DeA10947662BfD27A80e21a,0xBb3722366a6ca6B853441A2523ED9662423286fE,0x315C5c9e99792c618D361BE660d71580F490A837,0xc809009B6a9E7e28Cbc2a6Bac538e5a116dfe793,0xbaB57548E0eECAf65d02b4b0b6c66fBB61F067E9,0xFd85F43f23F6B8020DECf231EDBCe650dC9D9C46,0x8552218190B8689F6e01e27261ff0259416112C5,0x1D303893880107a3E1aA997a6a7c8d709CEB9564,0x1307e397BC8D0283126B2a9885023658Cc6f846e,0xB03392636f36C72969962Ee3184eC6F73D69a306,0x8dc6a91095208F70f966440E23b6c06080977829,0x3B007B2CAE55Bf922f1C2A3047ba4976420061E7,0x4d140ec8827AE31979d69a9f3F4374B2f6256113,0x8F08254fCB102Baff74AE5c74abD72ce0E1a139D,0x0787bf63f440ff0d24A8128815F898FaB5AB4383"}
	upkeepResetterContractEmpty  = ""
	upkeepResetterContractGoerli = "0xaeA9bD8f60C9EB1771900B9338dE8Ab52584E80e"
	upkeepResetterContractMumbai = "0x4Ef8599a41fd7b6788527E4a243d0Cf61b84f300"
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
			UpkeepGasLimit:        PerformGasToBurn + 50000,
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

func TestKeeperBenchmarkGoerliTestnetRegistry_1_3(t *testing.T) {
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

func TestKeeperBenchmarkOptimismGoerliRegistry_1_3(t *testing.T) {
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

func TestKeeperBenchmarkGoerliTestnetMulti_Registry(t *testing.T) {
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

func TestKeeperBenchmarkGoerliTestnetRegistry_1_2(t *testing.T) {
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

func TestKeeperBenchmarkGoerliTestnetRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(ChainlinkNodeFunding),
		int64(4),
		predeployedConsumersGoerli,
		upkeepResetterContractGoerli,
		12 * time.Second,
	})
}

func TestKeeperBenchmarkOptimismGoerliRegistry_2_0(t *testing.T) {
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

func TestKeeperBenchmarkPolygonMumbaiRegistry_2_0(t *testing.T) {
	keeperBenchmark(t, &BenchmarkTestEntry{
		[]eth_contracts.KeeperRegistryVersion{eth_contracts.RegistryVersion_2_0},
		big.NewFloat(ChainlinkNodeFunding),
		int64(10),
		predeployedConsumersMumbai,
		upkeepResetterContractMumbai,
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
