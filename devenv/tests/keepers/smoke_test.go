package keepers

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
	"github.com/smartcontractkit/chainlink/devenv/products/keepers"

	automation_tests "github.com/smartcontractkit/chainlink/devenv/tests/automation"
)

const (
	defaultUpkeepGasLimit           = uint32(2500000)
	defaultLinkFunds                = 9
	defaultEthFunds                 = 10.0
	defaultAmountOfUpkeeps          = 2
	defaultUpkeepExecutionTimeout   = "5m" // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer
	defaultExpectedUpkeepExecutions = 10

	numUpkeepsAllowedForStragglingTxs = 6
	expectedData                      = "expected data"
)

type testcase struct {
	Name string `toml:"name"`

	RegistryVersion          contracts.KeeperRegistryVersion `toml:"registryVersion"`
	UpkeepCount              int                             `toml:"upkeepCount,omitempty"`              // how many upkeeps to deploy
	ExpectedUpkeepExecutions int                             `toml:"expectedUpkeepExecutions,omitempty"` // how many times each upkeep should execute
	UpkeepExecutionTimeout   string                          `toml:"upkeepExecutionTimeout,omitempty"`   // "1s", "5m", 1h20m", etc
	UpkeepFundingLink        int64                           `toml:"upkeepFundingLink,omitempty"`

	TestKeyFundingEth float64 `toml:"testKeyFundingEth,omitempty"`
}

var (
	keeperDefaultRegistryConfig = contracts.KeeperRegistrySettings{
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
	lowBCPTRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(4),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
	}
	zeroAddress = common.Address{}
)

func TestKeeperBasic(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_1",
			RegistryVersion:          contracts.RegistryVersion_1_1,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
		},
		{
			Name:                     "registry_1_2",
			RegistryVersion:          contracts.RegistryVersion_1_2,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
		},
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			upkeeps, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			l.Info().Msgf("Waiting %s for %d upkeeps to be performed by %d contracts", testcase.UpkeepExecutionTimeout, testcase.ExpectedUpkeepExecutions, testcase.UpkeepCount)

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := range upkeepIDs {
					counter, err := upkeeps[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)),
						"Expected consumer counter to be greater than %d, but got %d", testcase.ExpectedUpkeepExecutions, counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Cancel all the registered upkeeps via the registry
			for i := range upkeepIDs {
				err := test.Registry.CancelUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not cancel upkeep at index %d", i)
			}

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := range upkeepIDs {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = upkeeps[i].Counter(t.Context())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Index", i).Int64("Upkeeps Performed", countersAfterCancellation[i].Int64()).Msg("Cancelled Upkeep")
			}

			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
					latestCounter, err := upkeeps[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.Equal(countersAfterCancellation[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterCancellation[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperBlockCountPerTurn(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: defaultUpkeepExecutionTimeout,
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: defaultUpkeepExecutionTimeout,
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: defaultUpkeepExecutionTimeout,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			upkeeps, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			l.Info().Msg("Waiting for 2m for upkeeps to be performed by different keepers")
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			stop := time.After(2 * time.Minute)

			keepersPerformedLowFreq := map[*big.Int][]string{}

		LOW_LOOP:
			for {
				select {
				case <-ticker.C:
					for i := range upkeepIDs {
						counter, err := upkeeps[i].Counter(t.Context())
						require.NoError(t, err, "Calling consumer's counter shouldn't fail")
						l.Info().Str("UpkeepId", upkeepIDs[i].String()).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

						upkeepInfo, err := test.Registry.GetUpkeepInfo(t.Context(), upkeepIDs[i])
						require.NoError(t, err, "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
						if latestKeeper == zeroAddress.String() {
							continue
						}

						keepersPerformedLowFreq[upkeepIDs[i]] = append(keepersPerformedLowFreq[upkeepIDs[i]], latestKeeper)
					}
				case <-stop:
					ticker.Stop()
					break LOW_LOOP
				}
			}

			require.GreaterOrEqual(t, testcase.UpkeepCount, len(keepersPerformedLowFreq), "At least %d different keepers should have been performing upkeeps", testcase.UpkeepCount)

			// Now set BCPT to be low, so keepers change turn frequently
			err = test.Registry.SetConfig(lowBCPTRegistryConfig, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")

			keepersPerformedHigherFreq := map[*big.Int][]string{}

			ticker = time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			stop = time.After(2 * time.Minute)

		HIGH_LOOP:
			for {
				select {
				case <-ticker.C:
					for i := range upkeepIDs {
						counter, err := upkeeps[i].Counter(t.Context())
						require.NoError(t, err, "Calling consumer's counter shouldn't fail")
						l.Info().Str("UpkeepId", upkeepIDs[i].String()).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

						upkeepInfo, err := test.Registry.GetUpkeepInfo(t.Context(), upkeepIDs[i])
						require.NoError(t, err, "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
						if latestKeeper == zeroAddress.String() {
							continue
						}

						keepersPerformedHigherFreq[upkeepIDs[i]] = append(keepersPerformedHigherFreq[upkeepIDs[i]], latestKeeper)
					}
				case <-stop:
					ticker.Stop()
					break HIGH_LOOP
				}
			}

			require.GreaterOrEqual(t, testcase.UpkeepCount+1, len(keepersPerformedHigherFreq), "At least %d different keepers should have been performing upkeeps after BCPT change", testcase.UpkeepCount+1)

			var countFreq = func(keepers []string, freqMap map[string]int) {
				for _, keeper := range keepers {
					freqMap[keeper]++
				}
			}

			for i := range upkeepIDs {
				lowFreqMap := make(map[string]int)
				highFreqMap := make(map[string]int)

				countFreq(keepersPerformedLowFreq[upkeepIDs[i]], lowFreqMap)
				countFreq(keepersPerformedHigherFreq[upkeepIDs[i]], highFreqMap)

				require.Greater(t, len(highFreqMap), len(lowFreqMap), "High frequency map should have more keepers than low frequency map")

				l.Info().Interface("Low BCPT", lowFreqMap).Interface("High BCPT", highFreqMap).Str("UpkeepID", upkeepIDs[i].String()).Msg("Keeper frequency map")

				for lowKeeper, lowFreq := range lowFreqMap {
					highFreq, ok := highFreqMap[lowKeeper]
					// it might happen due to fluke that a keeper is not found in high frequency map
					if !ok {
						continue
					}
					// require.True(t, ok, "Keeper %s not found in high frequency map. This should not happen", lowKeeper)
					require.GreaterOrEqual(t, lowFreq, highFreq, "Keeper %s should have performed less times with high BCPT than with low BCPT", lowKeeper)
				}
			}
		})
	}
}

func TestKeeperSimulation(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_1",
			RegistryVersion:          contracts.RegistryVersion_1_1,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   "1m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			err = automation.DeployMultiCallAndFundDeploymentAddresses(chainClient, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)))
			require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

			consumersPerformance := keepers.DeployKeeperConsumersPerformance(
				t,
				chainClient,
				testcase.UpkeepCount,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)

			var upkeepsAddresses []string
			for _, upkeep := range consumersPerformance {
				upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
			}

			upkeepIDs := automation.RegisterUpkeepContracts(t, chainClient, test.LinkToken, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, test.Registry, test.Registrar, testcase.UpkeepCount, upkeepsAddresses, false, false, false, nil)

			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			gom := gomega.NewGomegaWithT(t)
			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(t.Context())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)

				// Not even reverted upkeeps should be performed. Last keeper for the upkeep should be 0 address
				upkeepInfo, err := test.Registry.GetUpkeepInfo(t.Context(), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")
				g.Expect(upkeepInfo.LastKeeper).Should(gomega.Equal(zeroAddress.String()), "Last keeper should be zero address")
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
			err = consumerPerformance.SetPerformGasToBurn(t.Context(), big.NewInt(100000))
			require.NoError(t, err, "Error setting PerformGasToBurn")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) error {
				cnt, err := consumerPerformance.GetUpkeepCount(t.Context())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperCheckPerformGasLimit(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_2",
			RegistryVersion:          contracts.RegistryVersion_1_2,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   "3m",
		},
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   "3m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			err = automation.DeployMultiCallAndFundDeploymentAddresses(chainClient, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)))
			require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

			consumersPerformance := keepers.DeployKeeperConsumersPerformance(
				t,
				chainClient,
				testcase.UpkeepCount,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)

			var upkeepsAddresses []string
			for _, upkeep := range consumersPerformance {
				upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
			}

			upkeepIDs := automation.RegisterUpkeepContracts(t, chainClient, test.LinkToken, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, test.Registry, test.Registrar, testcase.UpkeepCount, upkeepsAddresses, false, false, false, nil)

			gom := gomega.NewGomegaWithT(t)
			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					cnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(
						gomega.Equal(int64(0)),
						"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
					)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Increase gas limit for the upkeep, higher than the performGasBurn
			l.Info().Msg("Setting upkeep gas limit higher than performGasBurn")
			for i := range upkeepIDs {
				err = test.Registry.SetUpkeepGasLimit(upkeepIDs[i], uint32(4500000))
				require.NoError(t, err, "Error setting Upkeep gas limit")
			}

			// Upkeep should now start performing
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Eventually(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					cnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
					)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			l.Info().Msg("Setting checkGasBurn higher than performGasBurn")
			for i := range upkeepIDs {
				err = consumersPerformance[i].SetCheckGasToBurn(t.Context(), big.NewInt(3000000))
				require.NoError(t, err, "Error setting CheckGasToBurn")
			}

			// Get existing performed count
			existingCnts := make(map[*big.Int]*big.Int)
			for i := range upkeepIDs {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
				existingCnts[upkeepIDs[i]] = existingCnt
				require.NoError(t, err, "Error calling consumer's counter")
				l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Check Gas Increased")
			}

			// In most cases count should remain constant, but there might be a straggling perform tx which
			// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
			// is sufficient to check that the upkeep count does not increase by more than 6.
			l.Info().Msg("Waiting for 3m to make sure no more than 6 upkeeps are performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					cnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					existingCnt := existingCnts[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(
						gomega.BeNumerically("<=", existingCnt.Int64()+numUpkeepsAllowedForStragglingTxs),
						"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
					)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			for i := range upkeepIDs {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
				existingCnts[upkeepIDs[i]] = existingCnt
				require.NoError(t, err, "Error calling consumer's counter")
				l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when consistently block finished")
			}

			// Now increase checkGasLimit on registry
			highCheckGasLimit := keeperDefaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			err = test.Registry.SetConfig(highCheckGasLimit, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")

			t.Cleanup(func() {
				err = test.Registry.SetConfig(keeperDefaultRegistryConfig, contracts.OCRv2Config{})
				require.NoError(t, err, "Error setting registry config")
			})

			// Upkeep should start performing again, and it should get regularly performed
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Eventually(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					cnt, err := consumersPerformance[i].GetUpkeepCount(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					existingCnt := existingCnts[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", existingCnt.Int64()),
						"Expected consumer counter to be greater than %d, but got %d", existingCnt.Int64(), cnt.Int64(),
					)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperRegisterUpkeep(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing and
			// store the value of their initial counters in order to compare later on that the value increased.
			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
					l.Info().
						Int64("Upkeep counter", counter.Int64()).
						Int("Upkeep ID", i).
						Msg("Number of upkeeps performed")
				}
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			newConsumers, _ := automation.RegisterNewUpkeeps(t, chainClient, test.LinkToken, test.Registry, test.Registrar, defaultUpkeepGasLimit, 1)

			// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
			newUpkeep := newConsumers[0]

			// Test that the newly registered upkeep is also performing.
			gom.Eventually(func(g gomega.Gomega) error {
				counter, err := newUpkeep.Counter(t.Context())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				l.Info().Msg("Newly registered upkeeps performed " + strconv.Itoa(int(counter.Int64())) + " times")
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					currentCounter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")

					l.Info().
						Int("Upkeep ID", i).
						Int64("Upkeep counter", currentCounter.Int64()).
						Int64("initial counter", initialCounters[i].Int64()).
						Msg("Number of upkeeps performed")

					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperAddFunds(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			// don't fund the upkeeps with LINK, we will add funds later
			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(1), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})
			// Since the upkeep is currently underfunded, check that it doesn't get executed
			gom := gomega.NewGomegaWithT(t)
			l.Info().Msg("Waiting for 1m to make sure no upkeeps are performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected consumer counter to remain zero, but got %d", counter.Int64())
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Grant permission to the registry to fund the upkeep
			err = test.LinkToken.Approve(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(9e18), big.NewInt(int64(len(upkeepIDs)))))
			require.NoError(t, err, "Error approving permissions for registry")

			// Add funds to the upkeep whose ID we know from above
			l.Info().Msg("Adding funds to upkeeps")
			for i := range upkeepIDs {
				err = test.Registry.AddUpkeepFunds(upkeepIDs[i], big.NewInt(9e18))
				require.NoError(t, err, "Error funding upkeep")
			}

			// Now the new upkeep should be performing because we added enough funds
			gom.Eventually(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperRemove(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})
			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Make sure the upkeeps are running before we remove a keeper
			gom.Eventually(func(g gomega.Gomega) error {
				for upkeepID := range upkeepIDs {
					counter, err := consumers[upkeepID].Counter(t.Context())
					initialCounters[upkeepID] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep with ID "+strconv.Itoa(upkeepID))
					g.Expect(counter.Cmp(big.NewInt(0))).To(gomega.Equal(1), "Expected consumer counter to be greater than 0, but got %s", counter)
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			keepers, err := test.Registry.GetKeeperList(t.Context())
			require.NoError(t, err, "Error getting list of Keepers")

			// Remove the first keeper from the list
			require.GreaterOrEqual(t, len(keepers), 2, "Expected there to be at least 2 keepers")
			newKeeperList := keepers[1:]

			cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
			require.NoError(t, err, "Failed to create chainlink client")

			// Construct the addresses of the payees required by the SetKeepers function
			payees := make([]string, len(keepers)-1)
			for i := range payees {
				payees[i], err = cl[0].PrimaryEthAddress()
				require.NoError(t, err, "Error building payee list")
			}

			err = test.Registry.SetKeepers(newKeeperList, payees, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting new list of Keepers")
			l.Info().Msg("Successfully removed keeper at address " + keepers[0] + " from the list of Keepers")

			// The upkeeps should still perform and their counters should have increased compared to the first check
			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Cmp(initialCounters[i])).To(gomega.Equal(1), "Expected consumer counter to be greater "+
						"than initial counter which was %s, but got %s", initialCounters[i], counter)
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperPauseRegistry(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing
			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Pause the registry
			err = test.Registry.Pause()
			require.NoError(t, err, "Error pausing the registry")

			t.Cleanup(func() {
				err = test.Registry.Unpause()
				require.NoError(t, err, "Error unpausing the registry")
			})

			// Store how many times each upkeep performed once the registry was successfully paused
			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := range upkeepIDs {
				countersAfterPause[i], err = consumers[i].Counter(t.Context())
				require.NoError(t, err, "Error retrieving consumer at index %d", i)
			}

			// After we paused the registry, the counters of all the upkeeps should stay constant
			// because they are no longer getting serviced
			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					latestCounter, err := consumers[i].Counter(t.Context())
					require.NoError(t, err, "Error retrieving consumer contract at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.Equal(countersAfterPause[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterPause[i].Int64(), latestCounter.Int64())
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperMigrateRegistry(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "1m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			// Deploy the second registry, second registrar, and the same number of upkeeps as the first one
			secondEthLinkFeed, err := contracts.DeployMockLINKETHFeed(chainClient, big.NewInt(2e18))
			require.NoError(t, err, "Failed to deploy mock ETH-LINK feed")

			secondGasFeed, err := contracts.DeployMockGASFeed(chainClient, big.NewInt(2e11))
			require.NoError(t, err, "Failed to deploy mock gas feed")

			secondTranscoder, err := contracts.DeployUpkeepTranscoder(chainClient)
			require.NoError(t, err, "Failed to deploy upkeep transcoder")

			secondRegistry, err := contracts.DeployKeeperRegistry(chainClient, &contracts.KeeperRegistryOpts{
				RegistryVersion: config.MustGetRegistryVersion(),
				LinkAddr:        test.LinkToken.Address(),
				ETHFeedAddr:     secondEthLinkFeed.Address(),
				GasFeedAddr:     secondGasFeed.Address(),
				TranscoderAddr:  secondTranscoder.Address(),
				RegistrarAddr:   zeroAddress.Hex(),
				Settings:        config.GetRegistryConfig(),
			})
			require.NoError(t, err, "Failed to deploy keeper registry")

			registrarSettings := contracts.KeeperRegistrarSettings{
				AutoApproveConfigType: 2,
				AutoApproveMaxAllowed: math.MaxUint16,
				RegistryAddr:          secondRegistry.Address(),
				MinLinkJuels:          big.NewInt(0),
			}

			secondRegistrar, err := keepers.DeployKeeperRegistrar(chainClient, config.MustGetRegistryVersion(), test.LinkToken, registrarSettings, secondRegistry)
			require.NoError(t, err, "Failed to deploy keeper registrar")

			err = secondRegistry.SetRegistrar(secondRegistrar.Address())
			require.NoError(t, err, "Failed to set registrar")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(secondRegistry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			chainlinkNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
			require.NoError(t, err, "Failed to create chainlink client")

			// create jobs for the second registry
			for _, chainlinkNode := range chainlinkNodes {
				chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
				require.NoError(t, err, "Failed to retrieve chainlink node address")
				_, err = chainlinkNode.MustCreateJob(&keepers.KeeperJobSpec{
					Name:                     "keeper-test-" + secondRegistry.Address(),
					ContractAddress:          secondRegistry.Address(),
					FromAddress:              chainlinkNodeAddress,
					EVMChainID:               strconv.FormatInt(chainClient.ChainID, 10),
					MinIncomingConfirmations: 1,
				})
				require.NoError(t, err, "Failed to create keeper job")
			}

			primaryNode := chainlinkNodes[0]
			primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
			require.NoError(t, err, "Failed to retrieve chainlink node address")

			nodeAddresses := make([]string, 0)
			for _, clNode := range chainlinkNodes {
				clNodeAddress, err := clNode.PrimaryEthAddress()
				require.NoError(t, err, "Failed to retrieve chainlink node address")
				nodeAddresses = append(nodeAddresses, clNodeAddress)
			}

			nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
			for _, cla := range nodeAddresses {
				nodeAddressesStr = append(nodeAddressesStr, cla)
				payees = append(payees, primaryNodeAddress)
			}

			err = secondRegistry.SetKeepers(nodeAddressesStr, payees, contracts.OCRv2Config{})
			require.NoError(t, err, "Failed to set keepers in the registry")

			_, _ = automation.DeployLegacyConsumers(t, chainClient, secondRegistry, secondRegistrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			err = secondRegistry.SetMigrationPermissions(common.HexToAddress(test.Registry.Address()), 3)
			require.NoError(t, err, "Error setting bidirectional permissions for first registry")
			err = test.Registry.SetMigrationPermissions(common.HexToAddress(secondRegistry.Address()), 3)
			require.NoError(t, err, "Error setting bidirectional permissions for second registry")

			gom := gomega.NewGomegaWithT(t)

			// Check that the first upkeep from the first registry is performing (before being migrated)
			l.Info().Msg("Waiting for 1m for upkeeps to be performed before migration")
			gom.Eventually(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					counterBeforeMigration, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counterBeforeMigration.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %s", counterBeforeMigration)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Migrate the upkeeps from the first to the second registry
			for i := range upkeepIDs {
				err = test.Registry.Migrate([]*big.Int{upkeepIDs[i]}, common.HexToAddress(secondRegistry.Address()))
				require.NoError(t, err, "Error migrating first upkeep")
			}

			// Pause the first registry, in that way we make sure that the upkeep is being performed by the second one
			err = test.Registry.Pause()
			require.NoError(t, err, "Error pausing registry")

			t.Cleanup(func() {
				err = test.Registry.Unpause()
				require.NoError(t, err, "Error unpausing the registry")
			})

			counterAfterMigrationPerUpkeep := make(map[*big.Int]*big.Int)

			for i := range upkeepIDs {
				counterAfterMigration, err := consumers[i].Counter(t.Context())
				require.NoError(t, err, "Error calling consumer's counter")
				counterAfterMigrationPerUpkeep[upkeepIDs[i]] = counterAfterMigration
			}

			// Check that once we migrated the upkeep, the counter has increased
			l.Info().Msg("Waiting for 1m for upkeeps to be performed after migration")
			gom.Eventually(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					currentCounter, err := consumers[i].Counter(t.Context())
					counterAfterMigration := counterAfterMigrationPerUpkeep[upkeepIDs[i]]
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", counterAfterMigration.Int64()),
						"Expected counter to have increased, but stayed constant at %s", counterAfterMigration)
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperJobReplacement(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepExecutionTimeout:   "5m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)),
						"Expected consumer counter to be greater than %d, but got %d", testcase.ExpectedUpkeepExecutions, counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			chainlinkNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
			require.NoError(t, err, "Failed to create chainlink client")

			for i, node := range chainlinkNodes {
				jobs, _, err := node.ReadJobs()
				require.NoError(t, err, "Failed to read jobs")
				var jobID string
				for _, job := range jobs.Data {
					require.IsType(t, map[string]interface{}{}, job, "job data should be a map[string]interface{}")
					if attr, ok := job["attributes"].(map[string]interface{}); ok {
						if name, ok := attr["name"].(string); ok && name == "keeper-test-"+test.Registry.Address() {
							jobID = job["id"].(string)
							break
						}
					}
				}
				require.NotEmpty(t, jobID, "Failed to find keeper job")
				err = node.MustDeleteJob(jobID)
				require.NoError(t, err, "Error deleting job from node %d", i)
			}

			// create jobs for the second registry
			for _, chainlinkNode := range chainlinkNodes {
				chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
				require.NoError(t, err, "Failed to retrieve chainlink node address")
				_, err = chainlinkNode.MustCreateJob(&keepers.KeeperJobSpec{
					Name:                     "keeper-test-" + test.Registry.Address(),
					ContractAddress:          test.Registry.Address(),
					FromAddress:              chainlinkNodeAddress,
					EVMChainID:               strconv.FormatInt(chainClient.ChainID, 10),
					MinIncomingConfirmations: 1,
				})
				require.NoError(t, err, "Failed to create keeper job")
			}

			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)),
						"Expected consumer counter to be greater than %d, but got %d", testcase.ExpectedUpkeepExecutions, counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperNodeDown(t *testing.T) {
	testcases := []testcase{
		{
			Name:                   "registry_1_1",
			RegistryVersion:        contracts.RegistryVersion_1_1,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "3m",
		},
		{
			Name:                   "registry_1_2",
			RegistryVersion:        contracts.RegistryVersion_1_2,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "3m",
		},
		{
			Name:                   "registry_1_3",
			RegistryVersion:        contracts.RegistryVersion_1_3,
			UpkeepCount:            defaultAmountOfUpkeeps,
			UpkeepFundingLink:      defaultLinkFunds,
			TestKeyFundingEth:      defaultEthFunds,
			UpkeepExecutionTimeout: "3m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Watch upkeeps being performed and store their counters in order to compare them later in the test
			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			chainlinkNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
			require.NoError(t, err, "Failed to create chainlink client")

			// Take down half of the Keeper nodes by deleting the Keeper job registered above (after registry deployment)
			cutIndex := len(chainlinkNodes)/2 + 1
			firstHalfToTakeDown := chainlinkNodes[:cutIndex]
			for i, nodeToTakeDown := range firstHalfToTakeDown {
				jobs, _, err := nodeToTakeDown.ReadJobs()
				require.NoError(t, err, "Failed to read jobs")
				var jobID string
				for _, job := range jobs.Data {
					require.IsType(t, map[string]interface{}{}, job, "job data should be a map[string]interface{}")
					if attr, ok := job["attributes"].(map[string]interface{}); ok {
						if name, ok := attr["name"].(string); ok && name == "keeper-test-"+test.Registry.Address() {
							jobID = job["id"].(string)
							break
						}
					}
				}
				require.NotEmpty(t, jobID, "Failed to find keeper job")
				err = nodeToTakeDown.MustDeleteJob(jobID)
				require.NoError(t, err, "Error deleting job from node %d", i)
			}
			l.Info().Msg("Successfully managed to take down the first half of the nodes")

			// Assert that upkeeps are still performed and their counters have increased
			gom.Eventually(func(g gomega.Gomega) error {
				for i := range upkeepIDs {
					currentCounter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// Take down the other half of the Keeper nodes
			secondHalfToTakeDown := chainlinkNodes[cutIndex:]
			for i, nodeToTakeDown := range secondHalfToTakeDown {
				jobs, _, err := nodeToTakeDown.ReadJobs()
				require.NoError(t, err, "Failed to read jobs")
				var jobID string
				for _, job := range jobs.Data {
					require.IsType(t, map[string]interface{}{}, job, "job data should be a map[string]interface{}")
					if attr, ok := job["attributes"].(map[string]interface{}); ok {
						if name, ok := attr["name"].(string); ok && name == "keeper-test-"+test.Registry.Address() {
							jobID = job["id"].(string)
							break
						}
					}
				}
				require.NotEmpty(t, jobID, "Failed to find keeper job")
				err = nodeToTakeDown.MustDeleteJob(jobID)
				require.NoError(t, err, "Error deleting job from node %d", i)
			}
			l.Info().Msg("Successfully managed to take down the second half of the nodes")

			// See how many times each upkeep was executed
			var countersAfterNoMoreNodes = make([]*big.Int, len(upkeepIDs))
			for i := range upkeepIDs {
				countersAfterNoMoreNodes[i], err = consumers[i].Counter(t.Context())
				require.NoError(t, err, "Error retrieving consumer counter %d", i)
				l.Info().
					Int("Index", i).
					Int64("Upkeeps", countersAfterNoMoreNodes[i].Int64()).
					Msg("Upkeeps Performed")
			}

			// Once all the nodes are taken down, there might be some straggling transactions which went through before
			// all the nodes were taken down. Every keeper node can have at most 1 straggling transaction per upkeep,
			// so a +6 on the upper limit side should be sufficient.
			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					latestCounter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=",
						countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs,
					),
						"Expected consumer counter to not have increased more than %d, but got %d",
						countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperPauseUnPauseUpkeep(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: 5,
			UpkeepExecutionTimeout:   "3m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)),
						"Expected consumer counter to be greater than %d, but got %d", testcase.ExpectedUpkeepExecutions, counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

			// pause all the registered upkeeps via the registry
			for i := range upkeepIDs {
				err := test.Registry.PauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Error pausing upkeep at index %d", i)
			}

			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := range upkeepIDs {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterPause[i], err = consumers[i].Counter(t.Context())
				require.NoError(t, err, "Error retrieving upkeep count at index %d", i)
				l.Info().
					Int("Index", i).
					Int64("Upkeeps", countersAfterPause[i].Int64()).
					Msg("Paused Upkeep")
			}

			gom.Consistently(func(g gomega.Gomega) {
				for i := range upkeepIDs {
					// In most cases counters should remain constant, but there might be a straggling perform tx which
					// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
					// is sufficient to check that the upkeep count does not increase by more than 6.
					latestCounter, err := consumers[i].Counter(t.Context())
					require.NoError(t, err, "Error retrieving counter at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs),
						"Expected consumer counter not have increased more than %d, but got %d",
						countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// unpause all the registered upkeeps via the registry
			for i := range upkeepIDs {
				err := test.Registry.UnpauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Error un-pausing upkeep at index %d", i)
			}

			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5 + numbers of performing before pause
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)+countersAfterPause[i].Int64()),
						"Expected consumer counter to be greater than %d, but got %d", int64(testcase.ExpectedUpkeepExecutions)+countersAfterPause[i].Int64(), counter.Int64())
					l.Info().Int64("Upkeeps", counter.Int64()).Msg("Upkeeps Performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}

// use env from TestKeeperPauseUnPauseUpkeep
func TestKeeperUpdateCheckData(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			ExpectedUpkeepExecutions: 5,
			UpkeepExecutionTimeout:   "3m",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			l := framework.L
			t.Cleanup(func() {
				err := products.ScanLogs(l, products.DefaultSettings())
				require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

				_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
				require.NoError(t, cErr)
			})

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[keepers.Configurator](outputFile)
			require.NoError(t, err)

			var config *keepers.Keepers
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
					config = candidate
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v", testcase.RegistryVersion.String())

			pks := []string{products.NetworkPrivateKey()}

			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			test, err := NewTest(l, chainClient, config)
			require.NoError(t, err, "Failed to create test")

			// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
			err = test.LinkToken.Transfer(test.Registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(testcase.UpkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			performDataChecker := keepers.DeployPerformDataChecker(t, chainClient, testcase.UpkeepCount, []byte(expectedData))

			err = automation.DeployMultiCallAndFundDeploymentAddresses(chainClient, test.LinkToken, testcase.UpkeepCount, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)))
			require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

			var upkeepsAddresses []string
			for _, upkeep := range performDataChecker {
				upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
			}

			upkeepIDs := automation.RegisterUpkeepContracts(t, chainClient, test.LinkToken, big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)), defaultUpkeepGasLimit, test.Registry, test.Registrar, testcase.UpkeepCount, upkeepsAddresses, false, false, false, nil)

			t.Cleanup(func() {
				automation_tests.GetStalenessReportCleanupFn(t, l, chainClient, sb, test.Registry, testcase.RegistryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			gom.Consistently(func(g gomega.Gomega) {
				// expect the counter to remain 0 because perform data does not match
				for i := range upkeepIDs {
					counter, err := performDataChecker[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected perform data checker counter to be 0, but got %d", counter.Int64())
					l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(gomega.Succeed())

			for i := range upkeepIDs {
				err = test.Registry.UpdateCheckData(upkeepIDs[i], []byte(expectedData))
				require.NoError(t, err, "Error updating check data at index %d", i)
			}

			// retrieve new check data for all upkeeps
			for i := range upkeepIDs {
				upkeep, err := test.Registry.GetUpkeepInfo(t.Context(), upkeepIDs[i])
				require.NoError(t, err, "Error getting upkeep info from index %d", i)
				require.Equal(t, []byte(expectedData), upkeep.CheckData, "Check data not as expected")
			}

			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than testcase.ExpectedUpkeepExecutions
				for i := range upkeepIDs {
					counter, err := performDataChecker[i].Counter(t.Context())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(testcase.ExpectedUpkeepExecutions)),
						"Expected perform data checker counter to be greater than %d, but got %d", testcase.ExpectedUpkeepExecutions, counter.Int64())
					l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
		})
	}
}
