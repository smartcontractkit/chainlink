package smoke

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
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
	keeperDefaultUpkeepsToDeploy      = 10
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
			if err != nil {
				t.Fatal(err)
			}

			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
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

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.CancelUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not cancel upkeep at index %d", i)
			}

			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for upkeeps to be cancelled")

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
			if err != nil {
				t.Fatal(err)
			}

			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				highBCPTRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			keepersPerformed := make([]string, 0)
			upkeepID := upkeepIDs[0]

			// Wait for upkeep to be performed twice by different keepers (buddies)
			gom.Eventually(func(g gomega.Gomega) error {
				counter, err := consumers[0].Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

				upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				l.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
				g.Expect(latestKeeper).ShouldNot(gomega.Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(gomega.BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				l.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			gom.Eventually(func(g gomega.Gomega) error {
				upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				g.Expect(latestKeeper).ShouldNot(gomega.Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(gomega.BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				l.Info().Str("Keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			// Expect no new keepers to perform for a while
			gom.Consistently(func(g gomega.Gomega) {
				upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				g.Expect(latestKeeper).ShouldNot(gomega.Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).Should(gomega.BeElementOf(keepersPerformed), "Existing keepers should alternate turns within BCPT")
			}, "1m", "1s").Should(gomega.Succeed())

			// Now set BCPT to be low, so keepers change turn frequently
			err = registry.SetConfig(lowBCPTRegistryConfig, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for set config tx")

			// Expect a new keeper to perform
			gom.Eventually(func(g gomega.Gomega) error {
				counter, err := consumers[0].Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Num upkeeps performed")

				upkeepInfo, err := registry.GetUpkeepInfo(testcontext.Get(t), upkeepID)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				l.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
				g.Expect(latestKeeper).ShouldNot(gomega.Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(gomega.BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				l.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
				return nil
			}, "1m", "1s").Should(gomega.Succeed())
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
			if err != nil {
				t.Fatal(err)
			}

			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumersPerformance, upkeepIDs := actions.DeployPerformanceKeeperContracts(
				t,
				registryVersion,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				&keeperDefaultRegistryConfig,
				big.NewInt(keeperDefaultLinkFunds),
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

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
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting to set PerformGasToBurn")

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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumersPerformance, upkeepIDs := actions.DeployPerformanceKeeperContracts(
				t,
				registryVersion,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				&keeperDefaultRegistryConfig,
				big.NewInt(keeperDefaultLinkFunds),
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "1m", "1s").Should(gomega.Succeed())

			// Increase gas limit for the upkeep, higher than the performGasBurn
			err = registry.SetUpkeepGasLimit(upkeepID, uint32(4500000))
			require.NoError(t, err, "Error setting Upkeep gas limit")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for SetUpkeepGasLimit tx")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) error {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
				return nil
			}, "1m", "1s").Should(gomega.Succeed())

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			err = consumerPerformance.SetCheckGasToBurn(testcontext.Get(t), big.NewInt(3000000))
			require.NoError(t, err, "Error setting CheckGasToBurn")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for SetCheckGasToBurn tx")

			// Get existing performed count
			existingCnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
			require.NoError(t, err, "Error calling consumer's counter")
			l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Check Gas Increased")

			// In most cases count should remain constant, but there might be a straggling perform tx which
			// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
			// is sufficient to check that the upkeep count does not increase by more than 6.
			gom.Consistently(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.BeNumerically("<=", existingCnt.Int64()+numUpkeepsAllowedForStragglingTxs),
					"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
				)
			}, "3m", "1s").Should(gomega.Succeed())

			existingCnt, err = consumerPerformance.GetUpkeepCount(testcontext.Get(t))
			require.NoError(t, err, "Error calling consumer's counter")
			existingCntInt := existingCnt.Int64()
			l.Info().Int64("Upkeep counter", existingCntInt).Msg("Upkeep counter when consistently block finished")

			// Now increase checkGasLimit on registry
			highCheckGasLimit := keeperDefaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			err = registry.SetConfig(highCheckGasLimit, contracts.OCRv2Config{})
			require.NoError(t, err, "Error setting registry config")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for set config tx")

			// Upkeep should start performing again, and it should get regularly performed
			gom.Eventually(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", existingCntInt),
					"Expected consumer counter to be greater than %d, but got %d", existingCntInt, cnt.Int64(),
				)
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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, registrar, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			var initialCounters = make([]*big.Int, len(upkeepIDs))

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

			newConsumers, _ := actions.RegisterNewUpkeeps(t, contractDeployer, chainClient, linkToken,
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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(1),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			// Since the upkeep is currently underfunded, check that it doesn't get executed
			gom.Consistently(func(g gomega.Gomega) {
				counter, err := consumers[0].Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
					"Expected consumer counter to remain zero, but got %d", counter.Int64())
			}, "1m", "1s").Should(gomega.Succeed())

			// Grant permission to the registry to fund the upkeep
			err = linkToken.Approve(registry.Address(), big.NewInt(9e18))
			require.NoError(t, err, "Error approving permissions for registry")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Add funds to the upkeep whose ID we know from above
			err = registry.AddUpkeepFunds(upkeepIDs[0], big.NewInt(9e18))
			require.NoError(t, err, "Error funding upkeep")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Now the new upkeep should be performing because we added enough funds
			gom.Eventually(func(g gomega.Gomega) {
				counter, err := consumers[0].Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			var initialCounters = make([]*big.Int, len(upkeepIDs))
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
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")
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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				keeperDefaultRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

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
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

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
	if err != nil {
		t.Fatal(err)
	}
	chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
	registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_2,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)
	gom := gomega.NewGomegaWithT(t)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error creating keeper jobs")

	// Deploy the second registry, second registrar, and the same number of upkeeps as the first one
	secondRegistry, _, _, _ := actions.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_2,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)

	// Set the jobs for the second registry
	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, secondRegistry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error creating keeper jobs")

	err = registry.SetMigrationPermissions(common.HexToAddress(secondRegistry.Address()), 3)
	require.NoError(t, err, "Error setting bidirectional permissions for first registry")
	err = secondRegistry.SetMigrationPermissions(common.HexToAddress(registry.Address()), 3)
	require.NoError(t, err, "Error setting bidirectional permissions for second registry")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting to set permissions")

	// Check that the first upkeep from the first registry is performing (before being migrated)
	gom.Eventually(func(g gomega.Gomega) error {
		counterBeforeMigration, err := consumers[0].Counter(testcontext.Get(t))
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(counterBeforeMigration.Int64()).Should(gomega.BeNumerically(">", int64(0)),
			"Expected consumer counter to be greater than 0, but got %s", counterBeforeMigration)
		return nil
	}, "1m", "1s").Should(gomega.Succeed())

	// Migrate the upkeep with index 0 from the first to the second registry
	err = registry.Migrate([]*big.Int{upkeepIDs[0]}, common.HexToAddress(secondRegistry.Address()))
	require.NoError(t, err, "Error migrating first upkeep")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for migration")

	// Pause the first registry, in that way we make sure that the upkeep is being performed by the second one
	err = registry.Pause()
	require.NoError(t, err, "Error pausing registry")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting to pause first registry")

	counterAfterMigration, err := consumers[0].Counter(testcontext.Get(t))
	require.NoError(t, err, "Error calling consumer's counter")

	// Check that once we migrated the upkeep, the counter has increased
	gom.Eventually(func(g gomega.Gomega) error {
		currentCounter, err := consumers[0].Counter(testcontext.Get(t))
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", counterAfterMigration.Int64()),
			"Expected counter to have increased, but stayed constant at %s", counterAfterMigration)
		return nil
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
			if err != nil {
				t.Fatal(err)
			}
			chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
			registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
				t,
				registryVersion,
				lowBCPTRegistryConfig,
				keeperDefaultUpkeepsToDeploy,
				keeperDefaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(keeperDefaultLinkFunds),
			)
			gom := gomega.NewGomegaWithT(t)

			jobs, err := actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
			require.NoError(t, err, "Error creating keeper jobs")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error creating keeper jobs")

			var initialCounters = make([]*big.Int, len(upkeepIDs))

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
				err = chainClient.WaitForEvents()
				require.NoError(t, err, "Error waiting for events")
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
				err = chainClient.WaitForEvents()
				require.NoError(t, err, "Error waiting for events")
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
	if err != nil {
		t.Fatal(err)
	}
	chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
	registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_3,
		lowBCPTRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)
	gom := gomega.NewGomegaWithT(t)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error creating keeper jobs")

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

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting to pause upkeeps")

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

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting to un-pause upkeeps")

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
	if err != nil {
		t.Fatal(err)
	}
	chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
	registry, _, performDataChecker, upkeepIDs := actions.DeployPerformDataCheckerContracts(
		t,
		ethereum.RegistryVersion_1_3,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		&lowBCPTRegistryConfig,
		big.NewInt(keeperDefaultLinkFunds),
		[]byte(keeperExpectedData),
	)
	gom := gomega.NewGomegaWithT(t)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error creating keeper jobs")

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

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for updated check data")

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

func setupKeeperTest(t *testing.T, config *tc.TestConfig) (
	blockchain.EVMClient,
	[]*client.ChainlinkClient,
	contracts.ContractDeployer,
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

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(config).
		WithGeth().
		WithCLNodes(5).
		WithCLNodeConfig(clNodeConfig).
		WithFunding(big.NewFloat(.5)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	env.ParallelTransactions(true)

	linkTokenContract, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	return env.EVMClient, env.ClCluster.NodeAPIs(), env.ContractDeployer, linkTokenContract, env
}

func TestKeeperJobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	registryVersion := ethereum.RegistryVersion_1_3
	config, err := tc.GetConfig("Smoke", tc.Keeper)
	if err != nil {
		t.Fatal(err)
	}

	chainClient, chainlinkNodes, contractDeployer, linkToken, _ := setupKeeperTest(t, &config)
	registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
		t,
		registryVersion,
		keeperDefaultRegistryConfig,
		keeperDefaultUpkeepsToDeploy,
		keeperDefaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(keeperDefaultLinkFunds),
	)
	gom := gomega.NewGomegaWithT(t)

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
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

	_, err = actions.CreateKeeperJobsLocal(l, chainlinkNodes, registry, contracts.OCRv2Config{}, chainClient.GetChainID().String())
	require.NoError(t, err, "Error creating keeper jobs")
	err = chainClient.WaitForEvents()
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
