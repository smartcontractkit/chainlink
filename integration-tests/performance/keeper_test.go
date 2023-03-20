package performance

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

var keeperDefaultRegistryConfig = contracts.KeeperRegistrySettings{
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

func TestKeeperPerformance(t *testing.T) {
	l := utils.GetTestLogger(t)
	testEnvironment, chainClient, chainlinkNodes, contractDeployer, linkToken := setupKeeperTest(t, "basic-smoke")
	if testEnvironment.WillUseRemoteRunner() {
		return
	}
	registry, _, consumers, upkeepIDs := actions.DeployKeeperContracts(
		t,
		ethereum.RegistryVersion_1_1,
		keeperDefaultRegistryConfig,
		1,
		uint32(2500000),
		linkToken,
		contractDeployer,
		chainClient,
		big.NewInt(9e18),
	)
	gom := gomega.NewGomegaWithT(t)

	profileFunction := func(chainlinkNode *client.Chainlink) {
		if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
			// Not the last node, hence not all nodes started profiling yet.
			return
		}
		actions.CreateKeeperJobs(t, chainlinkNodes, registry, contracts.OCRConfig{})
		err := chainClient.WaitForEvents()
		require.NoError(t, err, "Error creating keeper jobs")

		gom.Eventually(func(g gomega.Gomega) {
			// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
			for i := 0; i < len(upkeepIDs); i++ {
				counter, err := consumers[i].Counter(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(10)),
					"Expected consumer counter to be greater than 10, but got %d", counter.Int64())
				l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
			}
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
			countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			l.Info().Int("Index", i).Int64("Upkeeps Performed", countersAfterCancellation[i].Int64()).Msg("Cancelled Upkeep")
		}

		gom.Consistently(func(g gomega.Gomega) {
			for i := 0; i < len(upkeepIDs); i++ {
				// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
				latestCounter, err := consumers[i].Counter(context.Background())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				g.Expect(latestCounter.Int64()).Should(gomega.Equal(countersAfterCancellation[i].Int64()),
					"Expected consumer counter to remain constant at %d, but got %d",
					countersAfterCancellation[i].Int64(), latestCounter.Int64())
			}
		}, "1m", "1s").Should(gomega.Succeed())
	}
	profileTest := testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
		ProfileFunction: profileFunction,
		ProfileDuration: 10 * time.Second,
		ChainlinkNodes:  chainlinkNodes,
	})
	// Register cleanup
	t.Cleanup(func() {
		CleanupPerformanceTest(t, testEnvironment, chainlinkNodes, profileTest.TestReporter, chainClient)
	})
	profileTest.Setup(testEnvironment)
	profileTest.Run()
}

func setupKeeperTest(
	t *testing.T,
	testName string,
) (
	*environment.Environment,
	blockchain.EVMClient,
	[]*client.Chainlink,
	contracts.ContractDeployer,
	contracts.LinkToken,
) {
	network := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !network.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	baseTOML := `[WebServer]
HTTPWriteTimout = '300s'

[Keeper]
TurnLookBack = 0

[Keeper.Registry]
SyncInterval = '5s'
PerformGasOverhead = 150_000`
	networkName := strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")
	testEnvironment := environment.New(
		&environment.Config{
			NamespacePrefix: fmt.Sprintf("performance-keeper-%s-%s", testName, networkName),
			Test:            t,
		},
	).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml":     client.AddNetworksConfig(baseTOML, network),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error deploying test environment")
	if testEnvironment.WillUseRemoteRunner() {
		return testEnvironment, nil, nil, nil, nil
	}

	chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	chainClient.ParallelTransactions(true)

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.5))
	require.NoError(t, err, "Funding Chainlink nodes shouldn't fail")

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	return testEnvironment, chainClient, chainlinkNodes, contractDeployer, linkToken
}
