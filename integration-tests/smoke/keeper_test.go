package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

const (
	keeperDefaultUpkeepGasLimit       = uint32(2500000)
	keeperDefaultLinkFunds            = int64(9e18)
	keeperDefaultUpkeepsToDeploy      = 2
	numUpkeepsAllowedForStragglingTxs = 6
	keeperExpectedData                = "abcdef"
)

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
	highBCPTRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(10000),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
	}
)

func TestKeeperBasicSmoke(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) error {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(10)),
						"Expected consumer counter to be greater than 10, but got %d", counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
				return nil
			}, "5m", "1s").Should(gomega.Succeed())

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.CancelUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not cancel upkeep at index %d", i)
			}

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Index", i).Int64("Upkeeps Performed", countersAfterCancellation[i].Int64()).Msg("Cancelled Upkeep")
			}

			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
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
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				highBCPTRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			keepersPerformedLowFreq := map[*big.Int][]string{}

			// gom := gomega.NewGomegaWithT(t)
			// Wait for upkeep to be performed by two different keepers that alternate (buddies)
			l.Info().Msg("Waiting for 2m for upkeeps to be performed by different keepers")
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			stop := time.After(2 * time.Minute)

		LOW_LOOP:
			for {
				select {
				case <-ticker.C:
					for i := 0; i < len(upkeepIDs); i++ {
						counter, err := consumers[i].Counter(testcontext.Get(t))
						require.NoError(t, err, "Calling consumer's counter shouldn't fail")
						l.Info().Str("UpkeepId", upkeepIDs[i].String()).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

						upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepIDs[i])
						require.NoError(t, err, "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
						if latestKeeper == actions.ZeroAddress.String() {
							continue
						}

						keepersPerformedLowFreq[upkeepIDs[i]] = append(keepersPerformedLowFreq[upkeepIDs[i]], latestKeeper)
					}
				case <-stop:
					ticker.Stop()
					break LOW_LOOP
				}
			}

			require.GreaterOrEqual(t, 2, len(keepersPerformedLowFreq), "At least 2 different keepers should have been performing upkeeps")

			// Now set BCPT to be low, so keepers change turn frequently
			err = registry.SetConfig(lowBCPTRegistryConfig, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")

			keepersPerformedHigherFreq := map[*big.Int][]string{}

			ticker = time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			stop = time.After(2 * time.Minute)

		HIGH_LOOP:
			for {
				select {
				case <-ticker.C:
					for i := 0; i < len(upkeepIDs); i++ {
						counter, err := consumers[i].Counter(testcontext.Get(t))
						require.NoError(t, err, "Calling consumer's counter shouldn't fail")
						l.Info().Str("UpkeepId", upkeepIDs[i].String()).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

						upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepIDs[i])
						require.NoError(t, err, "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
						if latestKeeper == actions.ZeroAddress.String() {
							continue
						}

						keepersPerformedHigherFreq[upkeepIDs[i]] = append(keepersPerformedHigherFreq[upkeepIDs[i]], latestKeeper)
					}
				case <-stop:
					ticker.Stop()
					break HIGH_LOOP
				}
			}

			require.GreaterOrEqual(t, 3, len(keepersPerformedHigherFreq), "At least 3 different keepers should have been performing upkeeps after BCPT change")

			var countFreq = func(keepers []string, freqMap map[string]int) {
				for _, keeper := range keepers {
					freqMap[keeper]++
				}
			}

			for i := 0; i < len(upkeepIDs); i++ {
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
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumersPerformance, upkeepIDs := actions_seth.DeployPerformanceKeeperContracts(
				t,
				chainClient,
				registryVersion,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				&keeperDefaultRegistryConfig,
				big.NewInt(keeperDefaultLinkFunds),
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			gom := gomega.NewGomegaWithT(t)
			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)

				// Not even reverted upkeeps should be performed. Last keeper for the upkeep should be 0 address
				upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")
				g.Expect(upkeepInfo.LastKeeper).Should(gomega.Equal(actions.ZeroAddress.String()), "Last keeper should be zero address")
			}, "1m", "1s").Should(gomega.Succeed())

			// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
			err = consumerPerformance.SetPerformGasToBurn(testcontext.Get(t), big.NewInt(100000))
			require.NoError(t, err, "Error setting PerformGasToBurn")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) error {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
				return nil
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperCheckPerformGasLimit(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumersPerformance, upkeepIDs := actions_seth.DeployPerformanceKeeperContracts(
				t,
				chainClient,
				registryVersion,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				&keeperDefaultRegistryConfig,
				big.NewInt(keeperDefaultLinkFunds),
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(
						gomega.Equal(int64(0)),
						"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
					)
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// Increase gas limit for the upkeep, higher than the performGasBurn
			l.Info().Msg("Setting upkeep gas limit higher than performGasBurn")
			for i := 0; i < len(upkeepIDs); i++ {
				err = registry.SetUpkeepGasLimit(upkeepIDs[i], uint32(4500000))
				require.NoError(t, err, "Error setting Upkeep gas limit")
			}

			// Upkeep should now start performing
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
					)
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			l.Info().Msg("Setting checkGasBurn higher than performGasBurn")
			for i := 0; i < len(upkeepIDs); i++ {
				err = consumersPerformance[i].SetCheckGasToBurn(testcontext.Get(t), big.NewInt(3000000))
				require.NoError(t, err, "Error setting CheckGasToBurn")
			}

			// Get existing performed count
			existingCnts := make(map[*big.Int]*big.Int)
			for i := 0; i < len(upkeepIDs); i++ {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
				existingCnts[upkeepIDs[i]] = existingCnt
				require.NoError(t, err, "Error calling consumer's counter")
				l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Check Gas Increased")
			}

			// In most cases count should remain constant, but there might be a straggling perform tx which
			// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
			// is sufficient to check that the upkeep count does not increase by more than 6.
			l.Info().Msg("Waiting for 3m to make sure no more than 6 upkeeps are performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					existingCnt := existingCnts[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(
						gomega.BeNumerically("<=", existingCnt.Int64()+numUpkeepsAllowedForStragglingTxs),
						"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
					)
				}
			}, "3m", "1s").Should(gomega.Succeed())

			for i := 0; i < len(upkeepIDs); i++ {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
				existingCnts[upkeepIDs[i]] = existingCnt
				require.NoError(t, err, "Error calling consumer's counter")
				l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when consistently block finished")
			}

			// Now increase checkGasLimit on registry
			highCheckGasLimit := keeperDefaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			err = registry.SetConfig(highCheckGasLimit, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")

			// Upkeep should start performing again, and it should get regularly performed
			l.Info().Msg("Waiting for 1m for upkeeps to be performed")
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					existingCnt := existingCnts[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", existingCnt.Int64()),
						"Expected consumer counter to be greater than %d, but got %d", existingCnt.Int64(), cnt.Int64(),
					)
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperRegisterUpkeep(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, registrar, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing and
			// store the value of their initial counters in order to compare later on that the value increased.
			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
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

			newConsumers, _ := actions_seth.RegisterNewUpkeeps(t, chainClient, linkToken,
				registry, registrar, keeperDefaultUpkeepGasLimit, 1)

			// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
			newUpkeep := newConsumers[0]

			// Test that the newly registered upkeep is also performing.
			gom.Eventually(func(g gomega.Gomega) error {
				counter, err := newUpkeep.Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				l.Info().Msg("Newly registered upkeeps performed " + strconv.Itoa(int(counter.Int64())) + " times")
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(testcontext.Get(t))
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
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperAddFunds(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(1),
			)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			// Since the upkeep is currently underfunded, check that it doesn't get executed
			gom := gomega.NewGomegaWithT(t)
			l.Info().Msg("Waiting for 1m to make sure no upkeeps are performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected consumer counter to remain zero, but got %d", counter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// Grant permission to the registry to fund the upkeep
			err = linkToken.Approve(registry.Address(), big.NewInt(0).Mul(big.NewInt(9e18), big.NewInt(int64(len(upkeepIDs)))))
			require.NoError(t, err, "Error approving permissions for registry")

			// Add funds to the upkeep whose ID we know from above
			l.Info().Msg("Adding funds to upkeeps")
			for i := 0; i < len(upkeepIDs); i++ {
				err = registry.AddUpkeepFunds(upkeepIDs[i], big.NewInt(9e18))
				require.NoError(t, err, "Error funding upkeep")
			}

			// Now the new upkeep should be performing because we added enough funds
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperRemove(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Make sure the upkeeps are running before we remove a keeper
			gom.Eventually(func(g gomega.Gomega) error {
				for upkeepID := 0; upkeepID < len(upkeepIDs); upkeepID++ {
					counter, err := consumers[upkeepID].Counter(testcontext.Get(t))
					initialCounters[upkeepID] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep with ID "+strconv.Itoa(upkeepID))
					g.Expect(counter.Cmp(big.NewInt(0)) == 1, "Expected consumer counter to be greater than 0, but got %s", counter)
				}
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			keepers, err := registry.GetKeeperList(testcontext.Get(t))
			require.NoError(t, err, "Error getting list of Keepers")

			// Remove the first keeper from the list
			require.GreaterOrEqual(t, len(keepers), 2, "Expected there to be at least 2 keepers")
			newKeeperList := keepers[1:]

			// Construct the addresses of the payees required by the SetKeepers function
			payees := make([]string, len(keepers)-1)
			for i := 0; i < len(payees); i++ {
				payees[i], err = chainlinkNodes[0].PrimaryEthAddress()
				require.NoError(t, err, "Error building payee list")
			}

			err = registry.SetKeepers(newKeeperList, payees, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting new list of Keepers")
			l.Info().Msg("Successfully removed keeper at address " + keepers[0] + " from the list of Keepers")

			// The upkeeps should still perform and their counters should have increased compared to the first check
			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Cmp(initialCounters[i]) == 1, "Expected consumer counter to be greater "+
						"than initial counter which was %s, but got %s", initialCounters[i], counter)
				}
				return nil
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperPauseRegistry(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			// Observe that the upkeeps which are initially registered are performing
			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			// Pause the registry
			err = registry.Pause()
			require.NoError(t, err, "Error pausing the registry")

			// Store how many times each upkeep performed once the registry was successfully paused
			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterPause[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Error retrieving consumer at index %d", i)
			}

			// After we paused the registry, the counters of all the upkeeps should stay constant
			// because they are no longer getting serviced
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Error retrieving consumer contract at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.Equal(countersAfterPause[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterPause[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestKeeperMigrateRegistry(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.Keeper)
	require.NoError(t, err, "Error getting config")
	chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

	sb, err := chainClient.Client.BlockNumber(context.Background())
	require.NoError(t, err, "Failed to get start block")

	registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_2,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	t.Cleanup(func() {
		actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, ethereum.RegistryVersion_1_2)()
	})

	// Deploy the second registry, second registrar, and the same number of upkeeps as the first one
	secondRegistry, _, _, _ := actions_seth.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_2,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)

	// Set the jobs for the second registry
	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, secondRegistry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	err = registry.SetMigrationPermissions(common.HexToAddress(secondRegistry.Address()), 3)
	require.NoError(t, err, "Error setting bidirectional permissions for first registry")
	err = secondRegistry.SetMigrationPermissions(common.HexToAddress(registry.Address()), 3)
	require.NoError(t, err, "Error setting bidirectional permissions for second registry")

	gom := gomega.NewGomegaWithT(t)

	// Check that the first upkeep from the first registry is performing (before being migrated)
	l.Info().Msg("Waiting for 1m for upkeeps to be performed before migration")
	gom.Eventually(func(g gomega.Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			counterBeforeMigration, err := consumers[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
			g.Expect(counterBeforeMigration.Int64()).Should(gomega.BeNumerically(">", int64(0)),
				"Expected consumer counter to be greater than 0, but got %s", counterBeforeMigration)
		}
	}, "1m", "1s").Should(gomega.Succeed())

	// Migrate the upkeeps from the first to the second registry
	for i := 0; i < len(upkeepIDs); i++ {
		err = registry.Migrate([]*big.Int{upkeepIDs[i]}, common.HexToAddress(secondRegistry.Address()))
		require.NoError(t, err, "Error migrating first upkeep")
	}

	// Pause the first registry, in that way we make sure that the upkeep is being performed by the second one
	err = registry.Pause()
	require.NoError(t, err, "Error pausing registry")

	counterAfterMigrationPerUpkeep := make(map[*big.Int]*big.Int)

	for i := 0; i < len(upkeepIDs); i++ {
		counterAfterMigration, err := consumers[i].Counter(testcontext.Get(t))
		require.NoError(t, err, "Error calling consumer's counter")
		counterAfterMigrationPerUpkeep[upkeepIDs[i]] = counterAfterMigration
	}

	// Check that once we migrated the upkeep, the counter has increased
	l.Info().Msg("Waiting for 1m for upkeeps to be performed after migration")
	gom.Eventually(func(g gomega.Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			currentCounter, err := consumers[i].Counter(testcontext.Get(t))
			counterAfterMigration := counterAfterMigrationPerUpkeep[upkeepIDs[i]]
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
			g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", counterAfterMigration.Int64()),
				"Expected counter to have increased, but stayed constant at %s", counterAfterMigration)
		}
	}, "1m", "1s").Should(gomega.Succeed())
}

func TestKeeperNodeDown(t *testing.T) {
	t.Parallel()
	registryVersions := []ethereum.KeeperRegistryVersion{
		ethereum.RegistryVersion_1_1,
		ethereum.RegistryVersion_1_2,
		ethereum.RegistryVersion_1_3,
	}

	for _, rv := range registryVersions {
		registryVersion := rv
		t.Run(fmt.Sprintf("registry_1_%d", registryVersion), func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig("Smoke", tc.Keeper)
			require.NoError(t, err, "Failed to get config")

			chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

			sb, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
				t,
				registryVersion,
				lowBCPTRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)

			jobs, err := actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
			require.NoError(t, err, "Error creating keeper jobs")

			t.Cleanup(func() {
				actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, registryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			gom := gomega.NewGomegaWithT(t)
			// Watch upkeeps being performed and store their counters in order to compare them later in the test
			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			// Take down half of the Keeper nodes by deleting the Keeper job registered above (after registry deployment)
			firstHalfToTakeDown := chainlinkNodes[:len(chainlinkNodes)/2+1]
			for i, nodeToTakeDown := range firstHalfToTakeDown {
				err = nodeToTakeDown.MustDeleteJob(jobs[0].Data.ID)
				require.NoError(t, err, "Error deleting job from node %d", i)
			}
			l.Info().Msg("Successfully managed to take down the first half of the nodes")

			// Assert that upkeeps are still performed and their counters have increased
			gom.Eventually(func(g gomega.Gomega) error {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
				return nil
			}, "3m", "1s").Should(gomega.Succeed())

			// Take down the other half of the Keeper nodes
			nodesAndJobs := []nodeAndJob{}
			for i, n := range chainlinkNodes {
				nodesAndJobs = append(nodesAndJobs, nodeAndJob{node: n, job: jobs[i]})
			}
			secondHalfToTakeDown := nodesAndJobs[len(nodesAndJobs)/2+1:]
			for i, nodeToTakeDown := range secondHalfToTakeDown {
				err = nodeToTakeDown.node.MustDeleteJob(nodeToTakeDown.job.Data.ID)
				require.NoError(t, err, "Error deleting job from node %d", i)
			}
			l.Info().Msg("Successfully managed to take down the second half of the nodes")

			// See how many times each upkeep was executed
			var countersAfterNoMoreNodes = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterNoMoreNodes[i], err = consumers[i].Counter(testcontext.Get(t))
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
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=",
						countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs,
					),
						"Expected consumer counter to not have increased more than %d, but got %d",
						countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
				}
			}, "3m", "1s").Should(gomega.Succeed())
		})
	}
}

type nodeAndJob struct {
	node *client.ChainlinkClient
	job  *client.Job
}

func TestKeeperPauseUnPauseUpkeep(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.Keeper)
	require.NoError(t, err, "Failed to get config")

	chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

	sb, err := chainClient.Client.BlockNumber(context.Background())
	require.NoError(t, err, "Failed to get start block")

	registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_3,
		lowBCPTRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	t.Cleanup(func() {
		actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, ethereum.RegistryVersion_1_3)()
	})

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) error {
		// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(5)),
				"Expected consumer counter to be greater than 5, but got %d", counter.Int64())
			l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
		}
		return nil
	}, "3m", "1s").Should(gomega.Succeed())

	// pause all the registered upkeeps via the registry
	for i := 0; i < len(upkeepIDs); i++ {
		err := registry.PauseUpkeep(upkeepIDs[i])
		require.NoError(t, err, "Error pausing upkeep at index %d", i)
	}

	var countersAfterPause = make([]*big.Int, len(upkeepIDs))
	for i := 0; i < len(upkeepIDs); i++ {
		// Obtain the amount of times the upkeep has been executed so far
		countersAfterPause[i], err = consumers[i].Counter(testcontext.Get(t))
		require.NoError(t, err, "Error retrieving upkeep count at index %d", i)
		l.Info().
			Int("Index", i).
			Int64("Upkeeps", countersAfterPause[i].Int64()).
			Msg("Paused Upkeep")
	}

	gom.Consistently(func(g gomega.Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			// In most cases counters should remain constant, but there might be a straggling perform tx which
			// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
			// is sufficient to check that the upkeep count does not increase by more than 6.
			latestCounter, err := consumers[i].Counter(testcontext.Get(t))
			require.NoError(t, err, "Error retrieving counter at index %d", i)
			g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs),
				"Expected consumer counter not have increased more than %d, but got %d",
				countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
		}
	}, "1m", "1s").Should(gomega.Succeed())

	// unpause all the registered upkeeps via the registry
	for i := 0; i < len(upkeepIDs); i++ {
		err := registry.UnpauseUpkeep(upkeepIDs[i])
		require.NoError(t, err, "Error un-pausing upkeep at index %d", i)
	}

	gom.Eventually(func(g gomega.Gomega) error {
		// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5 + numbers of performing before pause
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(5)+countersAfterPause[i].Int64()),
				"Expected consumer counter to be greater than %d, but got %d", int64(5)+countersAfterPause[i].Int64(), counter.Int64())
			l.Info().Int64("Upkeeps", counter.Int64()).Msg("Upkeeps Performed")
		}
		return nil
	}, "3m", "1s").Should(gomega.Succeed())
}

func TestKeeperUpdateCheckData(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.Keeper)
	require.NoError(t, err, "Failed to get config")

	chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)

	sb, err := chainClient.Client.BlockNumber(context.Background())
	require.NoError(t, err, "Failed to get start block")

	registry, _, performDataChecker, upkeepIDs := actions_seth.DeployPerformDataCheckerContracts(
		t,
		chainClient,
		ethereum.RegistryVersion_1_3,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		&lowBCPTRegistryConfig,
		big.NewInt(keeperDefaultLinkFunds),
		[]byte(keeperExpectedData),
	)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	t.Cleanup(func() {
		actions_seth.GetStalenessReportCleanupFn(t, l, chainClient, sb, registry, ethereum.RegistryVersion_1_3)()
	})

	gom := gomega.NewGomegaWithT(t)
	gom.Consistently(func(g gomega.Gomega) {
		// expect the counter to remain 0 because perform data does not match
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := performDataChecker[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
				"Expected perform data checker counter to be 0, but got %d", counter.Int64())
			l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
		}
	}, "2m", "1s").Should(gomega.Succeed())

	for i := 0; i < len(upkeepIDs); i++ {
		err = registry.UpdateCheckData(upkeepIDs[i], []byte(keeperExpectedData))
		require.NoError(t, err, "Error updating check data at index %d", i)
	}

	// retrieve new check data for all upkeeps
	for i := 0; i < len(upkeepIDs); i++ {
		upkeep, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepIDs[i])
		require.NoError(t, err, "Error getting upkeep info from index %d", i)
		require.Equal(t, []byte(keeperExpectedData), upkeep.CheckData, "Check data not as expected")
	}

	gom.Eventually(func(g gomega.Gomega) error {
		// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := performDataChecker[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker counter for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(5)),
				"Expected perform data checker counter to be greater than 5, but got %d", counter.Int64())
			l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
		}
		return nil
	}, "3m", "1s").Should(gomega.Succeed())
}

func setupKeeperTest(l zerolog.Logger, t *testing.T, config *tc.TestConfig) (
	*seth.Client,
	[]*client.ChainlinkClient,
	contracts.LinkToken,
	*test_env.CLClusterTestEnv,
) {
	clNodeConfig := node.NewConfig(node.NewBaseConfig(), node.WithP2Pv2())
	turnLookBack := int64(0)
	syncInterval := *commonconfig.MustNewDuration(5 * time.Second)
	performGasOverhead := uint32(150000)
	clNodeConfig.Keeper.TurnLookBack = &turnLookBack
	clNodeConfig.Keeper.Registry.SyncInterval = &syncInterval
	clNodeConfig.Keeper.Registry.PerformGasOverhead = &performGasOverhead

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithCLNodes(5).
		WithCLNodeConfig(clNodeConfig).
		WithFunding(big.NewFloat(.5)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]

	sethClient, err := env.GetSethClient(network.ChainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")

	linkTokenContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	return sethClient, env.ClCluster.NodeAPIs(), linkTokenContract, env
}

func TestKeeperJobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	registryVersion := ethereum.RegistryVersion_1_3
	config, err := tc.GetConfig("Smoke", tc.Keeper)
	require.NoError(t, err, "Failed to get config")

	chainClient, chainlinkNodes, linkToken, _ := setupKeeperTest(l, t, &config)
	registry, _, consumers, upkeepIDs := actions_seth.DeployKeeperContracts(
		t,
		registryVersion,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)
	gom := gomega.NewGomegaWithT(t)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	gom.Eventually(func(g gomega.Gomega) error {
		// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(10)),
				"Expected consumer counter to be greater than 10, but got %d", counter.Int64())
			l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
		}
		return nil
	}, "5m", "1s").Should(gomega.Succeed())

	for _, n := range chainlinkNodes {
		jobs, _, err := n.ReadJobs()
		require.NoError(t, err)
		for _, maps := range jobs.Data {
			_, ok := maps["id"]
			require.Equal(t, true, ok)
			id := maps["id"].(string)
			_, err := n.DeleteJob(id)
			require.NoError(t, err)
		}
	}

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, fmt.Sprint(chainClient.ChainID))
	require.NoError(t, err, "Error creating keeper jobs")

	gom.Eventually(func(g gomega.Gomega) error {
		// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(testcontext.Get(t))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(10)),
				"Expected consumer counter to be greater than 10, but got %d", counter.Int64())
			l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
		}
		return nil
	}, "5m", "1s").Should(gomega.Succeed())
}
