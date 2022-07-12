package smoke

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota
	PerformanceCounter
)

const (
	defaultUpkeepGasLimit  = uint32(2500000)
	defaultLinkFunds       = int64(9e18)
	defaultUpkeepsToDeploy = 10
)

var defaultRegistryConfig = contracts.KeeperRegistrySettings{
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
}

var highBCPTRegistryConfig = contracts.KeeperRegistrySettings{
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

var (
	err                  error
	chainClient          blockchain.EVMClient
	contractDeployer     contracts.ContractDeployer
	registry             contracts.KeeperRegistry
	registrar            contracts.KeeperRegistrar
	consumers            []contracts.KeeperConsumer
	consumersPerformance []contracts.KeeperConsumerPerformance
	upkeepIDs            []*big.Int
	linkToken            contracts.LinkToken
	chainlinkNodes       []client.Chainlink
	testEnvironment      *environment.Environment
)

func basicSmokeTestFunction() {
	By("watches all the registered upkeeps perform and then cancels them from the registry")
	Eventually(func(g Gomega) {
		// Check if the upkeeps are performing by analysing their counters and checking they are greater than 0
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index "+strconv.Itoa(i))
			g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
				"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
			log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
		}
	}, "1m", "1s").Should(Succeed())

	// Cancel all the registered upkeeps via the registry
	for i := 0; i < len(upkeepIDs); i++ {
		err := registry.CancelUpkeep(upkeepIDs[i])
		Expect(err).ShouldNot(HaveOccurred(), "Could not cancel upkeep at index "+strconv.Itoa(i))
	}

	err = chainClient.WaitForEvents()
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
			g.Expect(latestCounter.Int64()).Should(Equal(countersAfterCancellation[i].Int64()),
				"Expected consumer counter to remain constant at %d, but got %d",
				countersAfterCancellation[i].Int64(), latestCounter.Int64())
		}
	}, "1m", "1s").Should(Succeed())
}

func bcptTestFunction() {
	By("tests that keeper pairs change turn every blockCountPerTurn")
	keepersPerformed := make([]string, 0)
	upkeepID := upkeepIDs[0]

	// Wait for upkeep to be performed twice by different keepers (buddies)
	Eventually(func(g Gomega) {
		counter, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

		upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
		g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

		latestKeeper := upkeepInfo.LastKeeper
		log.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
		g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
		g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

		log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
		keepersPerformed = append(keepersPerformed, latestKeeper)
	}, "1m", "1s").Should(Succeed())

	Eventually(func(g Gomega) {
		upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
		g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

		latestKeeper := upkeepInfo.LastKeeper
		g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
		g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

		log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
		keepersPerformed = append(keepersPerformed, latestKeeper)
	}, "1m", "1s").Should(Succeed())

	// Expect no new keepers to perform for a while
	Consistently(func(g Gomega) {
		upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
		g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

		latestKeeper := upkeepInfo.LastKeeper
		g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
		g.Expect(latestKeeper).Should(BeElementOf(keepersPerformed), "Existing keepers should alternate turns within BCPT")
	}, "1m", "1s").Should(Succeed())

	// Now set BCPT to be low, so keepers change turn frequently
	lowBcpt := defaultRegistryConfig
	lowBcpt.BlockCountPerTurn = big.NewInt(5)
	err = registry.SetConfig(lowBcpt)
	Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

	// Expect a new keeper to perform
	Eventually(func(g Gomega) {
		counter, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Num upkeeps performed")

		upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
		g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

		latestKeeper := upkeepInfo.LastKeeper
		log.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
		g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
		g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

		log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
		keepersPerformed = append(keepersPerformed, latestKeeper)
	}, "1m", "1s").Should(Succeed())
}

func performSimulationTestFunction() {
	By("tests that performUpkeep simulation is run before tx is broadcast")
	consumerPerformance := consumersPerformance[0]
	upkeepID := upkeepIDs[0]

	// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
	Consistently(func(g Gomega) {
		// Consumer count should remain at 0
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(
			Equal(int64(0)),
			"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
		)

		// Not even reverted upkeeps should be performed. Last keeper for the upkeep should be 0 address
		upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
		g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")
		g.Expect(upkeepInfo.LastKeeper).Should(Equal(actions.ZeroAddress.String()), "Last keeper should be zero address")
	}, "1m", "1s").Should(Succeed())

	// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
	err = consumerPerformance.SetPerformGasToBurn(context.Background(), big.NewInt(100000))
	Expect(err).ShouldNot(HaveOccurred(), "Perform gas should be set successfully on consumer")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set perform gas tx")

	// Upkeep should now start performing
	Eventually(func(g Gomega) {
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
			"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
		)
	}, "1m", "1s").Should(Succeed())
}

func checkPerformGasLimitTestFunction() {
	By("tests that check/perform gas limits are respected for upkeeps")
	consumerPerformance := consumersPerformance[0]
	upkeepID := upkeepIDs[0]

	// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
	Consistently(func(g Gomega) {
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(
			Equal(int64(0)),
			"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
		)
	}, "1m", "1s").Should(Succeed())

	// Increase gas limit for the upkeep, higher than the performGasBurn
	err = registry.SetUpkeepGasLimit(upkeepID, uint32(4500000))
	Expect(err).ShouldNot(HaveOccurred(), "Upkeep gas limit should be set successfully")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetUpkeepGasLimit tx")

	// Upkeep should now start performing
	Eventually(func(g Gomega) {
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
			"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
		)
	}, "1m", "1s").Should(Succeed())

	// Now increase the checkGasBurn on consumer, upkeep should stop performing
	err = consumerPerformance.SetCheckGasToBurn(context.Background(), big.NewInt(3000000))
	Expect(err).ShouldNot(HaveOccurred(), "Check gas burn should be set successfully on consumer")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetCheckGasToBurn tx")

	// Get existing performed count
	existingCnt, err := consumerPerformance.GetUpkeepCount(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
	log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when check gas increased")

	// In most cases count should remain constant, but there might be a straggling perform tx which
	// gets committed later. Hence, we check that the upkeep count does not increase by more than 1
	Consistently(func(g Gomega) {
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(
			BeNumerically("<=", existingCnt.Int64()+1),
			"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
		)
	}, "1m", "1s").Should(Succeed())

	// Now increase checkGasLimit on registry
	highCheckGasLimit := defaultRegistryConfig
	highCheckGasLimit.CheckGasLimit = uint32(5000000)
	err = registry.SetConfig(highCheckGasLimit)
	Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

	// Upkeep should start performing again, and it should get regularly performed
	Eventually(func(g Gomega) {
		cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
		g.Expect(cnt.Int64()).Should(BeNumerically(">", existingCnt.Int64()+1),
			"Expected consumer counter to be greater than %d, but got %d", existingCnt.Int64(), cnt.Int64(),
		)
	}, "1m", "1s").Should(Succeed())
}

func registerUpkeepTestFunction() {
	By("registers a new upkeep after the initial one was already registered and watches all of them perform")
	var initialCounters = make([]*big.Int, len(upkeepIDs))

	// Observe that the upkeeps which are initially registered are performing and
	// store the value of their initial counters in order to compare later on that the value increased.
	Eventually(func(g Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			initialCounters[i] = counter
			g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index "+strconv.Itoa(i))
			g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
				"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
			log.Info().
				Int64("Upkeep counter", counter.Int64()).
				Int64("Upkeep ID", int64(i)).
				Msg("Number of upkeeps performed")
		}
	}, "1m", "1s").Should(Succeed())

	newConsumers, _ := actions.RegisterNewUpkeeps(contractDeployer, chainClient, linkToken,
		registry, registrar, defaultUpkeepGasLimit, 1)

	// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
	newUpkeep := newConsumers[0]

	// Test that the newly registered upkeep is also performing.
	Eventually(func(g Gomega) {
		counter, err := newUpkeep.Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
		g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
			"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
		log.Info().Msg("Newly registered upkeeps performed " + strconv.Itoa(int(counter.Int64())) + " times")
	}, "1m", "1s").Should(Succeed())

	Eventually(func(g Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			currentCounter, err := consumers[i].Counter(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")

			log.Info().
				Int64("Upkeep ID", int64(i)).
				Int64("Upkeep counter", currentCounter.Int64()).
				Int64("initial counter", initialCounters[i].Int64()).
				Msg("Number of upkeeps performed")

			g.Expect(currentCounter.Int64()).Should(BeNumerically(">", initialCounters[i].Int64()),
				"Expected counter to have increased from initial value of %s, but got %s",
				initialCounters[i], currentCounter)
		}
	}, "1m", "1s").Should(Succeed())
}

func addFundsToUpkeepTestFunction() {
	By("adds funds to a new underfunded upkeep")
	// Since the upkeep is currently underfunded, check that it doesn't get executed
	Consistently(func(g Gomega) {
		counter, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(counter.Int64()).Should(Equal(int64(0)),
			"Expected consumer counter to remain zero, but got %d", counter.Int64())
	}, "1m", "1s").Should(Succeed())

	// Grant permission to the registry to fund the upkeep
	err = linkToken.Approve(registry.Address(), big.NewInt(9e18))
	Expect(err).ShouldNot(HaveOccurred(), "Could not approve permissions for the registry "+
		"on the link token contract")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

	// Add funds to the upkeep whose ID we know from above
	err = registry.AddUpkeepFunds(upkeepIDs[0], big.NewInt(9e18))
	Expect(err).ShouldNot(HaveOccurred(), "Could not fund upkeep")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

	// Now the new upkeep should be performing because we added enough funds
	Eventually(func(g Gomega) {
		counter, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
			"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
	}, "1m", "1s").Should(Succeed())
}

func removingKeeperTestFunction() {
	By("removes one keeper and makes sure the upkeeps are not affected by this and still perform")
	var initialCounters = make([]*big.Int, len(upkeepIDs))
	// Make sure the upkeeps are running before we remove a keeper
	Eventually(func(g Gomega) {
		for upkeepID := 0; upkeepID < len(upkeepIDs); upkeepID++ {
			counter, err := consumers[upkeepID].Counter(context.Background())
			initialCounters[upkeepID] = counter
			g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep with ID "+strconv.Itoa(upkeepID))
			g.Expect(counter.Cmp(big.NewInt(0)) == 1, "Expected consumer counter to be greater than 0, but got %s", counter)
		}
	}, "1m", "1s").Should(Succeed())

	keepers, err := registry.GetKeeperList(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Encountered error when getting the list of Keepers")

	// Remove the first keeper from the list
	newKeeperList := keepers[1:]

	// Construct the addresses of the payees required by the SetKeepers function
	payees := make([]string, len(keepers)-1)
	for i := 0; i < len(payees); i++ {
		payees[i], err = chainlinkNodes[0].PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Encountered error when building the payee list")
	}

	err = registry.SetKeepers(newKeeperList, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Encountered error when setting the new list of Keepers")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")
	log.Info().Msg("Successfully removed keeper at address " + keepers[0] + " from the list of Keepers")

	// The upkeeps should still perform and their counters should have increased compared to the first check
	Eventually(func(g Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index "+strconv.Itoa(i))
			g.Expect(counter.Cmp(initialCounters[i]) == 1, "Expected consumer counter to be greater "+
				"than initial counter which was %s, but got %s", initialCounters[i], counter)
		}
	}, "1m", "1s").Should(Succeed())
}

func pauseRegistryTestFunction() {
	By("pauses the registry and makes sure that the upkeeps are no longer performed")
	// Observe that the upkeeps which are initially registered are performing
	Eventually(func(g Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index "+strconv.Itoa(i))
			g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
				"Expected consumer counter to be greater than 0, but got %d")
		}
	}, "1m", "1s").Should(Succeed())

	// Pause the registry
	err := registry.Pause()
	Expect(err).ShouldNot(HaveOccurred(), "Could not pause the registry")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

	// Store how many times each upkeep performed once the registry was successfully paused
	var countersAfterPause = make([]*big.Int, len(upkeepIDs))
	for i := 0; i < len(upkeepIDs); i++ {
		countersAfterPause[i], err = consumers[i].Counter(context.Background())
		Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
			" for upkeep at index "+strconv.Itoa(i))
	}

	// After we paused the registry, the counters of all the upkeeps should stay constant
	// because they are no longer getting serviced
	Consistently(func(g Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			latestCounter, err := consumers[i].Counter(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
				" for upkeep at index "+strconv.Itoa(i))
			g.Expect(latestCounter.Int64()).Should(Equal(countersAfterPause[i].Int64()),
				"Expected consumer counter to remain constant at %d, but got %d",
				countersAfterPause[i].Int64(), latestCounter.Int64())
		}
	}, "1m", "1s").Should(Succeed())
}

func migrateUpkeepTestFunction() {
	By("creates another registry and migrates one upkeep to the new registry")
	// Deploy the second registry, second registrar, and the same number of upkeeps as the first one
	secondRegistry, _, _, _ := actions.DeployKeeperContracts(
		ethereum.RegistryVersion_1_2,
		defaultRegistryConfig,
		defaultUpkeepsToDeploy,
		defaultUpkeepGasLimit,
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(defaultLinkFunds),
	)

	// Set the jobs for the second registry
	actions.CreateKeeperJobs(chainlinkNodes, secondRegistry)
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")

	err = registry.SetMigrationPermissions(common.HexToAddress(secondRegistry.Address()), 3)
	Expect(err).ShouldNot(HaveOccurred(), "Couldn't set bidirectional permissions for first registry")
	err := secondRegistry.SetMigrationPermissions(common.HexToAddress(registry.Address()), 3)
	Expect(err).ShouldNot(HaveOccurred(), "Couldn't set bidirectional permissions for second registry")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting the permissions")

	// Check that the first upkeep from the first registry is performing (before being migrated)
	Eventually(func(g Gomega) {
		counterBeforeMigration, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(counterBeforeMigration.Int64()).Should(BeNumerically(">", int64(0)),
			"Expected consumer counter to be greater than 0, but got %s", counterBeforeMigration)
	}, "1m", "1s").Should(Succeed())

	// Migrate the upkeep with index 0 from the first to the second registry
	err = registry.Migrate([]*big.Int{upkeepIDs[0]}, common.HexToAddress(secondRegistry.Address()))
	Expect(err).ShouldNot(HaveOccurred(), "Couldn't migrate the first upkeep")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for the migration")

	// Pause the first registry, in that way we make sure that the upkeep is being performed by the second one
	err = registry.Pause()
	Expect(err).ShouldNot(HaveOccurred(), "Could not pause the registry")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for the pausing of the first registry")

	counterAfterMigration, err := consumers[0].Counter(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")

	// Check that once we migrated the upkeep, the counter has increased
	Eventually(func(g Gomega) {
		currentCounter, err := consumers[0].Counter(context.Background())
		g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
		g.Expect(currentCounter.Int64()).Should(BeNumerically(">", counterAfterMigration.Int64()),
			"Expected counter to have increased, but stayed constant at %s", counterAfterMigration)
	}, "1m", "1s").Should(Succeed())
}

var _ = DescribeTable("Keeper Suite @keeper", func(registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings, consumerContract KeeperConsumerContracts, testToRunAsFunction func(),
	linkFundsForEachUpkeep *big.Int) {

	BeforeEach(func() {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-keeper"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(eth.New(nil)).
			AddHelm(chainlink.New(0, map[string]interface{}{
				"replicas": "5",
				"env": map[string]interface{}{
					"MIN_INCOMING_CONFIRMATIONS": "1",
					"KEEPER_TURN_FLAG_ENABLED":   "true",
					"KEEPER_TURN_LOOK_BACK":      "0",
				},
			}))
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM)(testEnvironment)
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

		By("Deploy Keeper Contracts")
		linkToken, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		switch consumerContract {
		case BasicCounter:
			registry, registrar, consumers, upkeepIDs = actions.DeployKeeperContracts(
				registryVersion,
				registryConfig,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				linkFundsForEachUpkeep,
			)
		case PerformanceCounter:
			registry, registrar, consumersPerformance, upkeepIDs = actions.DeployPerformanceKeeperContracts(
				registryVersion,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				&registryConfig,
				linkFundsForEachUpkeep,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
		}

		By("Register Keeper Jobs")
		actions.CreateKeeperJobs(chainlinkNodes, registry)
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")
	})

	Describe("running specific test function", testToRunAsFunction)

	AfterEach(func() {
		By("Printing gas stats")
		chainClient.GasStats().PrintStats()

		By("Tearing down the environment")
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

},
	Entry("v1.1 basic smoke test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, basicSmokeTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 basic smoke test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, basicSmokeTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.1 BCPT test @simulated", ethereum.RegistryVersion_1_1, highBCPTRegistryConfig, BasicCounter, bcptTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 BCPT test @simulated", ethereum.RegistryVersion_1_2, highBCPTRegistryConfig, BasicCounter, bcptTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Perform simulation test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, performSimulationTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Check/Perform Gas limit test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, checkPerformGasLimitTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Register upkeep test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, registerUpkeepTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.1 Add funds to upkeep test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, addFundsToUpkeepTestFunction, big.NewInt(1)),
	Entry("v1.2 Add funds to upkeep test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, addFundsToUpkeepTestFunction, big.NewInt(1)),
	Entry("v1.1 Removing one keeper test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, removingKeeperTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Removing one keeper test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, removingKeeperTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Pause registry test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, pauseRegistryTestFunction, big.NewInt(defaultLinkFunds)),
	Entry("v1.2 Migrate upkeep from a registry to another @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, migrateUpkeepTestFunction, big.NewInt(defaultLinkFunds)),
)
