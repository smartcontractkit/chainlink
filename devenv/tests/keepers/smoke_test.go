package keepers

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
	"github.com/smartcontractkit/chainlink/devenv/products/keepers"
)

const defaultUpkeepGasLimit = uint32(2500000)

type testcase struct {
	Name string `toml:"name"`

	RegistryVersion          contracts.KeeperRegistryVersion `toml:"registryVersion"`
	UpkeepCount              int                             `toml:"upkeepCount,omitempty"`              // how many upkeeps to deploy
	ExpectedUpkeepExecutions int                             `toml:"expectedUpkeepExecutions,omitempty"` // how many times each upkeep should execute
	UpkeepExecutionTimeout   string                          `toml:"upkeepExecutionTimeout,omitempty"`   // "1s", "5m", 1h20m", etc
	UpkeepFundingLink        int64                           `toml:"upkeepFundingLink,omitempty"`

	TestKeyFundingEth float64 `toml:"testKeyFundingEth,omitempty"`

	// Chainlink Docker image to which nodes should be upgraded to
	// upgradeImage string
}

// TODO: migrate other tests from integration-tests/smoke/keeper_test.go

func TestKeeperBasic(t *testing.T) {
	testcases := []testcase{
		{
			Name:                     "registry_1_1",
			RegistryVersion:          contracts.RegistryVersion_1_1,
			UpkeepCount:              2,
			UpkeepFundingLink:        1000000000000000000,
			TestKeyFundingEth:        10,
			ExpectedUpkeepExecutions: 10,
			UpkeepExecutionTimeout:   "5m",
		},
		{
			Name:                     "registry_1_2",
			RegistryVersion:          contracts.RegistryVersion_1_2,
			UpkeepCount:              2,
			UpkeepFundingLink:        1000000000000000000,
			TestKeyFundingEth:        10,
			ExpectedUpkeepExecutions: 10,
			UpkeepExecutionTimeout:   "5m",
		},
		{
			Name:                     "registry_1_3",
			RegistryVersion:          contracts.RegistryVersion_1_3,
			UpkeepCount:              2,
			UpkeepFundingLink:        1000000000000000000,
			TestKeyFundingEth:        10,
			ExpectedUpkeepExecutions: 10,
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

			upkeeps, upkeepIDs := automation.DeployLegacyConsumers(t, chainClient, test.Registry, test.Registrar, test.LinkToken, testcase.UpkeepCount, big.NewInt(testcase.UpkeepFundingLink), defaultUpkeepGasLimit, false, false, false, nil)

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
